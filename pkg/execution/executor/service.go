package executor

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inngest/inngest/pkg/config"
	"github.com/inngest/inngest/pkg/cqrs"
	"github.com/inngest/inngest/pkg/event"
	"github.com/inngest/inngest/pkg/execution"
	"github.com/inngest/inngest/pkg/execution/queue"
	"github.com/inngest/inngest/pkg/execution/state"
	"github.com/inngest/inngest/pkg/inngest"
	"github.com/inngest/inngest/pkg/logger"
	"github.com/inngest/inngest/pkg/pubsub"
	"github.com/inngest/inngest/pkg/service"
	"github.com/oklog/ulid/v2"
)

type Opt func(s *svc)

func WithExecutionLoader(l cqrs.ExecutionLoader) func(s *svc) {
	return func(s *svc) {
		s.data = l
	}
}

func WithState(sm state.Manager) func(s *svc) {
	return func(s *svc) {
		s.state = sm
	}
}

func WithServiceExecutor(exec execution.Executor) func(s *svc) {
	return func(s *svc) {
		s.exec = exec
	}
}

func WithExecutorOpts(opts ...ExecutorOpt) func(s *svc) {
	return func(s *svc) {
		s.opts = opts
	}
}

func WithServiceQueue(q queue.Queue) func(s *svc) {
	return func(s *svc) {
		s.queue = q
	}
}

func NewService(c config.Config, opts ...Opt) service.Service {
	svc := &svc{config: c}
	for _, o := range opts {
		o(svc)
	}
	return svc
}

type svc struct {
	config config.Config
	// data provides the ability to load action versions when running steps.
	data cqrs.ExecutionLoader
	// state allows us to record step results
	state state.Manager
	// queue allows us to enqueue next steps.
	queue queue.Queue
	// exec runs the specific actions.
	exec execution.Executor

	wg sync.WaitGroup

	opts []ExecutorOpt
}

func (s *svc) Name() string {
	return "executor"
}

func (s *svc) Pre(ctx context.Context) error {
	var err error

	if s.state == nil {
		return fmt.Errorf("no state provided")
	}

	if s.queue == nil {
		return fmt.Errorf("no queue provided")
	}

	failureHandler, err := s.getFailureHandler(ctx)
	if err != nil {
		return fmt.Errorf("failed to create failure handler: %w", err)
	}
	s.exec.SetFailureHandler(failureHandler)

	return nil
}

func (s *svc) Executor() execution.Executor {
	return s.exec
}

func (s *svc) getFailureHandler(ctx context.Context) (func(context.Context, state.Identifier, state.State, state.DriverResponse) error, error) {
	pb, err := pubsub.NewPublisher(ctx, s.config.EventStream.Service)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	topicName := s.config.EventStream.Service.Concrete.TopicName()

	return func(ctx context.Context, id state.Identifier, s state.State, r state.DriverResponse) error {
		now := time.Now()
		evt := event.Event{
			ID:        ulid.MustNew(uint64(now.UnixMilli()), rand.Reader).String(),
			Name:      event.FnFailedName,
			Timestamp: now.UnixMilli(),
			Data: map[string]interface{}{
				"function_id": s.Function().Slug,
				"run_id":      id.RunID.String(),
				"error":       r.UserError(),
				"event":       s.Event(),
			},
		}

		byt, err := json.Marshal(evt)
		if err != nil {
			return fmt.Errorf("error marshalling failure event: %w", err)
		}

		err = pb.Publish(
			ctx,
			topicName,
			pubsub.Message{
				Name:      event.EventReceivedName,
				Data:      string(byt),
				Timestamp: now,
			},
		)
		if err != nil {
			return fmt.Errorf("error publishing failure event: %w", err)
		}

		return nil
	}, nil
}

func (s *svc) Run(ctx context.Context) error {
	logger.From(ctx).Info().Msg("subscribing to function queue")
	return s.queue.Run(ctx, func(ctx context.Context, item queue.Item) error {
		// Don't stop the service on errors.
		s.wg.Add(1)
		defer s.wg.Done()

		var err error
		switch item.Kind {
		case queue.KindEdge, queue.KindSleep:
			err = s.handleQueueItem(ctx, item)
		case queue.KindPause:
			err = s.handlePauseTimeout(ctx, item)
		default:
			err = fmt.Errorf("unknown payload type: %T", item.Payload)
		}
		return err
	})
}

func (s *svc) Stop(ctx context.Context) error {
	// Wait for all in-flight queue runs to finish
	s.wg.Wait()
	return nil
}

func (s *svc) handleQueueItem(ctx context.Context, item queue.Item) error {
	payload, err := queue.GetEdge(item)
	if err != nil {
		return fmt.Errorf("unable to get edge from queue item: %w", err)
	}
	edge := payload.Edge

	// If this is of type sleep, ensure that we save "nil" within the state store
	// for the outgoing edge ID.  This ensures that we properly increase the stack
	// for `tools.sleep` within generator functions.
	var stackIdx int
	if item.Kind == queue.KindSleep && item.Attempt == 0 {
		stackIdx, err = s.state.SaveResponse(ctx, item.Identifier, state.DriverResponse{
			Step: inngest.Step{ID: edge.Outgoing}, // XXX: Save edge name here.
		}, 0)
		if err != nil {
			return err
		}
		// After the sleep, we start a new step.  THis means we also want to start a new
		// group ID, ensuring that we correlate the next step _after_ this sleep (to be
		// scheduled in this executor run)
		ctx = state.WithGroupID(ctx, uuid.New().String())
	} else if edge.Outgoing != inngest.TriggerName {
		// Load the position within the stack for standard edges.
		stackIdx, err = s.state.StackIndex(ctx, item.Identifier.RunID, edge.Outgoing)
		if err != nil {
			return fmt.Errorf("unable to find stack index: %w", err)
		}
	}

	resp, err := s.exec.Execute(ctx, item.Identifier, item, edge, stackIdx)

	// Check if the execution is cancelled, and if so finalize and terminate early.
	// This prevents steps from scheduling children.
	if err == state.ErrFunctionCancelled {
		return nil
	}
	if err != nil || resp.Err != nil {
		// Accordingly, we check if the driver's response is retryable here;
		// this will let us know whether we can re-enqueue.
		if resp != nil && !resp.Retryable() {
			return nil
		}

		// If the error is not of type response error, we assume the step is
		// always retryable.
		if resp == nil || err != nil {
			return err
		}

		// Always retry; non-retryable is covered above.
		return fmt.Errorf("%s", resp.Error())
	}

	return nil
}

func (s *svc) handlePauseTimeout(ctx context.Context, item queue.Item) error {
	l := logger.From(ctx).With().Str("run_id", item.Identifier.RunID.String()).Logger()

	pauseTimeout, ok := item.Payload.(queue.PayloadPauseTimeout)
	if !ok {
		return fmt.Errorf("unable to get pause timeout form queue item: %T", item.Payload)
	}

	pause, err := s.state.PauseByID(ctx, pauseTimeout.PauseID)
	if err == state.ErrPauseNotFound {
		// This pause has been consumed.
		l.Debug().Interface("pause", pauseTimeout).Msg("consumed pause timeout ignored")
		return nil
	}
	if err != nil {
		return err
	}
	if pause == nil {
		return nil
	}

	return s.exec.Resume(ctx, *pause, execution.ResumeRequest{})
}
