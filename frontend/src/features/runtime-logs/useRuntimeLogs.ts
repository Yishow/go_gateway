import { useEffect, useMemo, useSyncExternalStore } from 'react';
import type { RuntimeLogFilters } from '../../types/runtimeLogs';
import { RuntimeLogSession } from './runtimeLogSession';
export function useRuntimeLogs({ level, source, q }: RuntimeLogFilters) {
  const session = useMemo(() => new RuntimeLogSession({ level, source, q }), [level, source, q]);
  const view = useSyncExternalStore(session.subscribe, session.getSnapshot);
  useEffect(() => { session.start(); return session.stop; }, [session]);
  return { ...view, pause: session.pause, resume: session.resume, clear: session.clear, retry: session.retry };
}
