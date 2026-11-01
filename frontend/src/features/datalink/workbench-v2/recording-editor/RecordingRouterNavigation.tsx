import { useCallback, useEffect, type ReactNode } from 'react';
import { useBlocker } from 'react-router-dom';
import { RecordingDraftNavigation, useRecordingNavigation } from './RecordingDraftNavigation';

function RecordingRouteBlocker() {
  const { blocked, request } = useRecordingNavigation();
  const blocker = useBlocker(useCallback(() => blocked(), [blocked]));
  useEffect(() => {
    if (blocker.state === 'blocked') request(blocker.proceed, blocker.reset);
  }, [blocker, request]);
  return null;
}
/** Installed only under a data router; the ordinary form boundary needs no router. */
export function RecordingRouterNavigation({ children }: { children: ReactNode }) {
  return <RecordingDraftNavigation><RecordingRouteBlocker />{children}</RecordingDraftNavigation>;
}
