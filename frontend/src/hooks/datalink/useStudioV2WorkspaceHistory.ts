import { useMutation, useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceHistoryAPI } from '../../services/studioV2WorkspaceHistory';
import type { HistoryQuery } from '../../types/historyReport';

export function useStudioV2WorkspaceHistoryQuery(query: HistoryQuery, enabled: boolean) {
  const queryKeyStr = `${query.plan_id}:${(query.measurement_ids ?? []).join(',')}:${query.start_time ?? ''}:${query.end_time ?? ''}:${query.resolution ?? ''}`;

  return useQuery({
    queryKey: studioV2WorkspaceKeys.history(queryKeyStr),
    queryFn: () => studioV2WorkspaceHistoryAPI.queryHistory(query),
    enabled: enabled && !!query.plan_id,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useExportHistoryCSVMutation() {
  return useMutation({
    mutationFn: (query: HistoryQuery) => studioV2WorkspaceHistoryAPI.exportCSV(query),
  });
}
