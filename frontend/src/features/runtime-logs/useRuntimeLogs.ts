import { useEffect, useMemo, useSyncExternalStore } from 'react';
import type { RuntimeLogFilters } from '../../types/runtimeLogs';
import { RuntimeLogSession } from './runtimeLogSession';
export function useRuntimeLogs(filters: RuntimeLogFilters) {
  const { level, source, q, stage } = filters;
  // Only equivalent unrestricted forms normalize; invalid nonempty values reach validation.
  const workspaceId = filters.workspaceId === '' ? undefined : filters.workspaceId;
  const deviceId = filters.deviceId === '' ? undefined : filters.deviceId;
  const ruleId = filters.ruleId === '' ? undefined : filters.ruleId;
  const groupId = filters.groupId === '' ? undefined : filters.groupId;
  const actionable = filters.actionable === false ? undefined : filters.actionable;
  const session = useMemo(() => new RuntimeLogSession({ level, source, q,
    workspaceId, deviceId, ruleId, groupId, stage, actionable }),
  [level, source, q, workspaceId, deviceId, ruleId, groupId, stage, actionable]);
  const view = useSyncExternalStore(session.subscribe, session.getSnapshot);
  useEffect(() => { session.start(); return session.stop; }, [session]);
  return { ...view, pause: session.pause, resume: session.resume, clear: session.clear, retry: session.retry };
}
