import { EventSection } from '@/components/Event/Section';
import { FunctionRunSection } from '@/components/Function/RunSection';
import { useAppSelector } from '@/store/hooks';

export default function StreamDetails() {
  const selectedEvent = useAppSelector((state) => state.global.selectedEvent);
  const selectedRun = useAppSelector((state) => state.global.selectedRun);

  return (
    <>
      {selectedEvent && (
        <div className="grid grid-cols-2 h-full text-white">
          <EventSection eventId={selectedEvent} />
          <FunctionRunSection runId={selectedRun} />
        </div>
      )}
    </>
  );
}
