import { useQuery } from '@tanstack/react-query';
import { runtimeAPI } from '../../../services/datalink';

interface UseRuntimeStatusOptions {
  deviceId: string | null;
  pollingIntervalMs?: number | false;
}

export function useRuntimeStatus({
  deviceId,
  pollingIntervalMs = false,
}: UseRuntimeStatusOptions) {
  return useQuery({
    queryKey: ['runtime-dashboard-status', deviceId],
    queryFn: () => runtimeAPI.getStatus(deviceId ?? undefined),
    enabled: Boolean(deviceId),
    refetchInterval: pollingIntervalMs,
    refetchOnWindowFocus: false,
    retry: false,
  });
}
