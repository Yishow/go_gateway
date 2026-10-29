import { useEffect, useRef } from 'react';
import { loadBasicRecordingIntent, sameBasicRecordingIntent, subscribeBasicRecordingIntent, type BasicRecordingIntentRecord } from '../../state/basicRecording';

/** Re-read the existing record after a late owner reply; this never executes or rebuilds a request. */
export function useBasicRecordingIntentObservation(scopeKey: string, record: BasicRecordingIntentRecord | undefined,
  onObserved: (record: BasicRecordingIntentRecord) => void): void {
  const current = useRef({ scopeKey, record, onObserved }); current.current = { scopeKey, record, onObserved };
  useEffect(() => {
    let mounted = true;
    const refresh = () => {
      const owner = current.current;
      if (!mounted || owner.scopeKey !== scopeKey || !owner.record || owner.record.operation_id) return;
      const next = loadBasicRecordingIntent(scopeKey);
      if (!next?.operation_id || !sameBasicRecordingIntent(owner.record, next)) return;
      // Repeated synchronous publications cannot enqueue duplicate adoptions before React renders.
      current.current = { ...owner, record: next };
      owner.onObserved(next);
    };
    const unsubscribe = subscribeBasicRecordingIntent(scopeKey, refresh);
    refresh(); // Close the render-to-subscription gap using the same retained record.
    return () => { mounted = false; unsubscribe(); };
  }, [scopeKey, record]);
}
