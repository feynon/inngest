// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: queries.sql

package sqlc

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	ulid "github.com/oklog/ulid/v2"
)

const deleteApp = `-- name: DeleteApp :exec
UPDATE apps SET deleted_at = NOW() WHERE id = ?
`

func (q *Queries) DeleteApp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteApp, id)
	return err
}

const deleteFunctionsByAppID = `-- name: DeleteFunctionsByAppID :exec
DELETE FROM functions WHERE app_id = ?
`

func (q *Queries) DeleteFunctionsByAppID(ctx context.Context, appID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFunctionsByAppID, appID)
	return err
}

const deleteFunctionsByIDs = `-- name: DeleteFunctionsByIDs :exec
DELETE FROM functions WHERE id IN (/*SLICE:ids*/?)
`

func (q *Queries) DeleteFunctionsByIDs(ctx context.Context, ids []uuid.UUID) error {
	query := deleteFunctionsByIDs
	var queryParams []interface{}
	if len(ids) > 0 {
		for _, v := range ids {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:ids*/?", "NULL", 1)
	}
	_, err := q.db.ExecContext(ctx, query, queryParams...)
	return err
}

const getAllApps = `-- name: GetAllApps :many
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps
`

func (q *Queries) GetAllApps(ctx context.Context) ([]*App, error) {
	rows, err := q.db.QueryContext(ctx, getAllApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*App
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SdkLanguage,
			&i.SdkVersion,
			&i.Framework,
			&i.Metadata,
			&i.Status,
			&i.Error,
			&i.Checksum,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApp = `-- name: GetApp :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE id = ?
`

func (q *Queries) GetApp(ctx context.Context, id uuid.UUID) (*App, error) {
	row := q.db.QueryRowContext(ctx, getApp, id)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppByChecksum = `-- name: GetAppByChecksum :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE checksum = ? LIMIT 1
`

func (q *Queries) GetAppByChecksum(ctx context.Context, checksum string) (*App, error) {
	row := q.db.QueryRowContext(ctx, getAppByChecksum, checksum)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppByURL = `-- name: GetAppByURL :one
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE url = ? LIMIT 1
`

func (q *Queries) GetAppByURL(ctx context.Context, url string) (*App, error) {
	row := q.db.QueryRowContext(ctx, getAppByURL, url)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const getAppFunctions = `-- name: GetAppFunctions :many
SELECT id, app_id, name, slug, config, created_at FROM functions WHERE app_id = ?
`

func (q *Queries) GetAppFunctions(ctx context.Context, appID uuid.UUID) ([]*Function, error) {
	rows, err := q.db.QueryContext(ctx, getAppFunctions, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Function
	for rows.Next() {
		var i Function
		if err := rows.Scan(
			&i.ID,
			&i.AppID,
			&i.Name,
			&i.Slug,
			&i.Config,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApps = `-- name: GetApps :many
SELECT id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url FROM apps WHERE deleted_at IS NULL
`

func (q *Queries) GetApps(ctx context.Context) ([]*App, error) {
	rows, err := q.db.QueryContext(ctx, getApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*App
	for rows.Next() {
		var i App
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SdkLanguage,
			&i.SdkVersion,
			&i.Framework,
			&i.Metadata,
			&i.Status,
			&i.Error,
			&i.Checksum,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventByInternalID = `-- name: GetEventByInternalID :one
SELECT internal_id, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE internal_id = ?
`

func (q *Queries) GetEventByInternalID(ctx context.Context, internalID ulid.ULID) (*Event, error) {
	row := q.db.QueryRowContext(ctx, getEventByInternalID, internalID)
	var i Event
	err := row.Scan(
		&i.InternalID,
		&i.EventID,
		&i.EventName,
		&i.EventData,
		&i.EventUser,
		&i.EventV,
		&i.EventTs,
	)
	return &i, err
}

const getEventsTimebound = `-- name: GetEventsTimebound :many
SELECT internal_id, event_id, event_name, event_data, event_user, event_v, event_ts FROM events WHERE event_ts > ? AND event_ts <= ? ORDER BY event_ts DESC LIMIT ?
`

type GetEventsTimeboundParams struct {
	After  time.Time
	Before time.Time
	Limit  int64
}

func (q *Queries) GetEventsTimebound(ctx context.Context, arg GetEventsTimeboundParams) ([]*Event, error) {
	rows, err := q.db.QueryContext(ctx, getEventsTimebound, arg.After, arg.Before, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.InternalID,
			&i.EventID,
			&i.EventName,
			&i.EventData,
			&i.EventUser,
			&i.EventV,
			&i.EventTs,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionByID = `-- name: GetFunctionByID :one
SELECT id, app_id, name, slug, config, created_at FROM functions WHERE id = ?
`

func (q *Queries) GetFunctionByID(ctx context.Context, id uuid.UUID) (*Function, error) {
	row := q.db.QueryRowContext(ctx, getFunctionByID, id)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const getFunctionRunFinishesByRunIDs = `-- name: GetFunctionRunFinishesByRunIDs :many
SELECT run_id, status, output, completed_step_count, created_at FROM function_finishes WHERE run_id IN (/*SLICE:run_ids*/?)
`

func (q *Queries) GetFunctionRunFinishesByRunIDs(ctx context.Context, runIds []ulid.ULID) ([]*FunctionFinish, error) {
	query := getFunctionRunFinishesByRunIDs
	var queryParams []interface{}
	if len(runIds) > 0 {
		for _, v := range runIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:run_ids*/?", strings.Repeat(",?", len(runIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:run_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*FunctionFinish
	for rows.Next() {
		var i FunctionFinish
		if err := rows.Scan(
			&i.RunID,
			&i.Status,
			&i.Output,
			&i.CompletedStepCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunHistory = `-- name: GetFunctionRunHistory :many
SELECT id, created_at, run_started_at, function_id, function_version, run_id, event_id, batch_id, group_id, idempotency_key, type, attempt, latency_ms, step_name, step_id, url, cancel_request, sleep, wait_for_event, wait_result, result FROM history WHERE run_id = ? ORDER BY created_at ASC
`

func (q *Queries) GetFunctionRunHistory(ctx context.Context, runID ulid.ULID) ([]*History, error) {
	rows, err := q.db.QueryContext(ctx, getFunctionRunHistory, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.RunStartedAt,
			&i.FunctionID,
			&i.FunctionVersion,
			&i.RunID,
			&i.EventID,
			&i.BatchID,
			&i.GroupID,
			&i.IdempotencyKey,
			&i.Type,
			&i.Attempt,
			&i.LatencyMs,
			&i.StepName,
			&i.StepID,
			&i.Url,
			&i.CancelRequest,
			&i.Sleep,
			&i.WaitForEvent,
			&i.WaitResult,
			&i.Result,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunsFromEvents = `-- name: GetFunctionRunsFromEvents :many
SELECT run_id, run_started_at, function_id, function_version, trigger_type, event_id, batch_id, original_run_id FROM function_runs WHERE event_id IN (/*SLICE:event_ids*/?)
`

func (q *Queries) GetFunctionRunsFromEvents(ctx context.Context, eventIds []ulid.ULID) ([]*FunctionRun, error) {
	query := getFunctionRunsFromEvents
	var queryParams []interface{}
	if len(eventIds) > 0 {
		for _, v := range eventIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:event_ids*/?", strings.Repeat(",?", len(eventIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:event_ids*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*FunctionRun
	for rows.Next() {
		var i FunctionRun
		if err := rows.Scan(
			&i.RunID,
			&i.RunStartedAt,
			&i.FunctionID,
			&i.FunctionVersion,
			&i.TriggerType,
			&i.EventID,
			&i.BatchID,
			&i.OriginalRunID,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctionRunsTimebound = `-- name: GetFunctionRunsTimebound :many
SELECT run_id, run_started_at, function_id, function_version, trigger_type, event_id, batch_id, original_run_id FROM function_runs WHERE run_started_at > ? AND run_started_at <= ? LIMIT ?
`

type GetFunctionRunsTimeboundParams struct {
	After  time.Time
	Before time.Time
	Limit  int64
}

func (q *Queries) GetFunctionRunsTimebound(ctx context.Context, arg GetFunctionRunsTimeboundParams) ([]*FunctionRun, error) {
	rows, err := q.db.QueryContext(ctx, getFunctionRunsTimebound, arg.After, arg.Before, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*FunctionRun
	for rows.Next() {
		var i FunctionRun
		if err := rows.Scan(
			&i.RunID,
			&i.RunStartedAt,
			&i.FunctionID,
			&i.FunctionVersion,
			&i.TriggerType,
			&i.EventID,
			&i.BatchID,
			&i.OriginalRunID,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFunctions = `-- name: GetFunctions :many
SELECT id, app_id, name, slug, config, created_at FROM functions
`

func (q *Queries) GetFunctions(ctx context.Context) ([]*Function, error) {
	rows, err := q.db.QueryContext(ctx, getFunctions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Function
	for rows.Next() {
		var i Function
		if err := rows.Scan(
			&i.ID,
			&i.AppID,
			&i.Name,
			&i.Slug,
			&i.Config,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const hardDeleteApp = `-- name: HardDeleteApp :exec
DELETE FROM apps WHERE id = ?
`

func (q *Queries) HardDeleteApp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, hardDeleteApp, id)
	return err
}

const insertApp = `-- name: InsertApp :one
INSERT INTO apps
	(id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, url) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type InsertAppParams struct {
	ID          uuid.UUID
	Name        string
	SdkLanguage string
	SdkVersion  string
	Framework   sql.NullString
	Metadata    string
	Status      string
	Error       sql.NullString
	Checksum    string
	Url         string
}

func (q *Queries) InsertApp(ctx context.Context, arg InsertAppParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, insertApp,
		arg.ID,
		arg.Name,
		arg.SdkLanguage,
		arg.SdkVersion,
		arg.Framework,
		arg.Metadata,
		arg.Status,
		arg.Error,
		arg.Checksum,
		arg.Url,
	)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const insertEvent = `-- name: InsertEvent :exec

INSERT INTO events
	(internal_id, event_id, event_name, event_data, event_user, event_v, event_ts) VALUES
	(?, ?, ?, ?, ?, ?, ?)
`

type InsertEventParams struct {
	InternalID ulid.ULID
	EventID    string
	EventName  string
	EventData  string
	EventUser  string
	EventV     sql.NullString
	EventTs    time.Time
}

// Events
func (q *Queries) InsertEvent(ctx context.Context, arg InsertEventParams) error {
	_, err := q.db.ExecContext(ctx, insertEvent,
		arg.InternalID,
		arg.EventID,
		arg.EventName,
		arg.EventData,
		arg.EventUser,
		arg.EventV,
		arg.EventTs,
	)
	return err
}

const insertFunction = `-- name: InsertFunction :one


INSERT INTO functions
	(id, app_id, name, slug, config, created_at) VALUES
	(?, ?, ?, ?, ?, ?) RETURNING id, app_id, name, slug, config, created_at
`

type InsertFunctionParams struct {
	ID        uuid.UUID
	AppID     uuid.UUID
	Name      string
	Slug      string
	Config    string
	CreatedAt time.Time
}

// functions
//
// note - this is very basic right now.
func (q *Queries) InsertFunction(ctx context.Context, arg InsertFunctionParams) (*Function, error) {
	row := q.db.QueryRowContext(ctx, insertFunction,
		arg.ID,
		arg.AppID,
		arg.Name,
		arg.Slug,
		arg.Config,
		arg.CreatedAt,
	)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}

const insertFunctionFinish = `-- name: InsertFunctionFinish :exec
INSERT INTO function_finishes
	(run_id, status, output, completed_step_count, created_at) VALUES 
	(?, ?, ?, ?, ?)
`

type InsertFunctionFinishParams struct {
	RunID              ulid.ULID
	Status             string
	Output             string
	CompletedStepCount int64
	CreatedAt          time.Time
}

func (q *Queries) InsertFunctionFinish(ctx context.Context, arg InsertFunctionFinishParams) error {
	_, err := q.db.ExecContext(ctx, insertFunctionFinish,
		arg.RunID,
		arg.Status,
		arg.Output,
		arg.CompletedStepCount,
		arg.CreatedAt,
	)
	return err
}

const insertFunctionRun = `-- name: InsertFunctionRun :exec

INSERT INTO function_runs
	(run_id, run_started_at, function_id, function_version, trigger_type, event_id, batch_id, original_run_id) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertFunctionRunParams struct {
	RunID           ulid.ULID
	RunStartedAt    time.Time
	FunctionID      uuid.UUID
	FunctionVersion int64
	TriggerType     string
	EventID         ulid.ULID
	BatchID         ulid.ULID
	OriginalRunID   ulid.ULID
}

// function runs
func (q *Queries) InsertFunctionRun(ctx context.Context, arg InsertFunctionRunParams) error {
	_, err := q.db.ExecContext(ctx, insertFunctionRun,
		arg.RunID,
		arg.RunStartedAt,
		arg.FunctionID,
		arg.FunctionVersion,
		arg.TriggerType,
		arg.EventID,
		arg.BatchID,
		arg.OriginalRunID,
	)
	return err
}

const insertHistory = `-- name: InsertHistory :exec

INSERT INTO history
	(id, created_at, run_started_at, function_id, function_version, run_id, event_id, batch_id, group_id, idempotency_key, type, attempt, latency_ms, step_name, step_id, url, cancel_request, sleep, wait_for_event, wait_result, result) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertHistoryParams struct {
	ID              ulid.ULID
	CreatedAt       time.Time
	RunStartedAt    time.Time
	FunctionID      uuid.UUID
	FunctionVersion int64
	RunID           ulid.ULID
	EventID         ulid.ULID
	BatchID         ulid.ULID
	GroupID         sql.NullString
	IdempotencyKey  string
	Type            string
	Attempt         int64
	LatencyMs       sql.NullInt64
	StepName        sql.NullString
	StepID          sql.NullString
	Url             sql.NullString
	CancelRequest   sql.NullString
	Sleep           sql.NullString
	WaitForEvent    sql.NullString
	WaitResult      sql.NullString
	Result          sql.NullString
}

// History
func (q *Queries) InsertHistory(ctx context.Context, arg InsertHistoryParams) error {
	_, err := q.db.ExecContext(ctx, insertHistory,
		arg.ID,
		arg.CreatedAt,
		arg.RunStartedAt,
		arg.FunctionID,
		arg.FunctionVersion,
		arg.RunID,
		arg.EventID,
		arg.BatchID,
		arg.GroupID,
		arg.IdempotencyKey,
		arg.Type,
		arg.Attempt,
		arg.LatencyMs,
		arg.StepName,
		arg.StepID,
		arg.Url,
		arg.CancelRequest,
		arg.Sleep,
		arg.WaitForEvent,
		arg.WaitResult,
		arg.Result,
	)
	return err
}

const updateAppError = `-- name: UpdateAppError :one
UPDATE apps SET error = ? WHERE id = ? RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type UpdateAppErrorParams struct {
	Error sql.NullString
	ID    uuid.UUID
}

func (q *Queries) UpdateAppError(ctx context.Context, arg UpdateAppErrorParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, updateAppError, arg.Error, arg.ID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const updateAppURL = `-- name: UpdateAppURL :one
UPDATE apps SET url = ? WHERE id = ? RETURNING id, name, sdk_language, sdk_version, framework, metadata, status, error, checksum, created_at, deleted_at, url
`

type UpdateAppURLParams struct {
	Url string
	ID  uuid.UUID
}

func (q *Queries) UpdateAppURL(ctx context.Context, arg UpdateAppURLParams) (*App, error) {
	row := q.db.QueryRowContext(ctx, updateAppURL, arg.Url, arg.ID)
	var i App
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SdkLanguage,
		&i.SdkVersion,
		&i.Framework,
		&i.Metadata,
		&i.Status,
		&i.Error,
		&i.Checksum,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.Url,
	)
	return &i, err
}

const updateFunctionConfig = `-- name: UpdateFunctionConfig :one
UPDATE functions SET config = ? WHERE id = ? RETURNING id, app_id, name, slug, config, created_at
`

type UpdateFunctionConfigParams struct {
	Config string
	ID     uuid.UUID
}

func (q *Queries) UpdateFunctionConfig(ctx context.Context, arg UpdateFunctionConfigParams) (*Function, error) {
	row := q.db.QueryRowContext(ctx, updateFunctionConfig, arg.Config, arg.ID)
	var i Function
	err := row.Scan(
		&i.ID,
		&i.AppID,
		&i.Name,
		&i.Slug,
		&i.Config,
		&i.CreatedAt,
	)
	return &i, err
}
