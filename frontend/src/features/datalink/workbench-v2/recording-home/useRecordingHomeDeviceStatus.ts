import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { runtimeAPI } from '../../../../services/datalink';
import type { RuntimeStatusWithDiagnostics } from '../../../../types/runtimeDiagnostics';

/** One workspace-wide read, stopped while the page is hidden or unmounted. */
export function useRecordingHomeDeviceStatus(workspaceId: string) {
  const [visible, setVisible] = useState(() => document.visibilityState !== 'hidden');
  const [clock, setClock] = useState(() => Date.now());
  useEffect(() => {
    const change = () => { setClock(Date.now()); setVisible(document.visibilityState !== 'hidden'); };
    document.addEventListener('visibilitychange', change);
    return () => document.removeEventListener('visibilitychange', change);
  }, []);
  const query = useQuery({
    queryKey: ['recording-home-runtime-status', workspaceId],
    queryFn: () => runtimeAPI.getStatus(), enabled: visible, retry: false,
    refetchInterval: 5000, refetchIntervalInBackground: false, refetchOnWindowFocus: false, gcTime: 15000,
  });
  useEffect(() => {
    if (!visible || !query.dataUpdatedAt) return;
    const remaining = query.dataUpdatedAt + 15000 - Date.now();
    if (remaining <= 0) return;
    const timer = window.setTimeout(() => setClock(Date.now()), remaining);
    return () => window.clearTimeout(timer);
  }, [query.dataUpdatedAt, visible]);
  const now = Math.max(clock, Date.now());
  const data = query.data as RuntimeStatusWithDiagnostics | undefined;
  const available = Boolean(data && data.snapshot_state && !data.snapshot_state.unavailable && data.snapshot_state.state !== 'unavailable' && now - query.dataUpdatedAt < 15000);
  return { snapshot: available ? data : undefined, outdated: available && (query.isError || data?.snapshot_state?.stale === true),
    updatedAt: query.dataUpdatedAt || undefined, now: query.isError && available ? query.dataUpdatedAt : now, displayNow: now, refetch: query.refetch };
}
