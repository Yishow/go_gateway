export { observeRecordingHomeSourceMutation, captureRecordingHomeSourceMutation, completeRecordingHomeSourceMutation } from './recordingHomeSourceMutation';
import { useCallback, useSyncExternalStore } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { recordingHomeSourcesAPI } from '../../services/recordingHomeSources';
import type { RecordingHomeScope, RecordingHomeEvidence } from '../../features/datalink/workbench-v2/recording-home/recordingHomeProjection';
import type { RecordingHomeSourceSummary } from '../../types/recordingHomeSources';
import { notifySources, recordingHomeSourceContext, recordingHomeSourcePrefix, sourceEntry, sourcePartition, type HomeSourceContext } from './recordingHomeSourceCache';
export type { HomeSourceContext } from './recordingHomeSourceCache';
export { invalidateRecordingHomeSources, recordingHomeSourceContext, recordingHomeSourcePrefix } from './recordingHomeSourceCache';
export function useRecordingHomeSources(scope: RecordingHomeScope, options: HomeSourceContext = {}) {
  const client = useQueryClient(), context = recordingHomeSourceContext(options), partition = sourcePartition(client, context);
  const subscribe = useCallback((listener: () => void) => { partition.listeners.add(listener); return () => { partition.listeners.delete(listener); }; }, [partition]);
  useSyncExternalStore(subscribe, () => partition.version, () => partition.version);
  // Rendering may be speculative. Only a committed query request may register
  // a visit or create/evict a cached workspace. Render reads are side-effect free.
  const observedEntry = partition.entries.get(scope.workspaceId);
  const query = useQuery({
    queryKey: [...recordingHomeSourcePrefix(context, scope.workspaceId), scope.workspaceRevision, scope.epoch],
    queryFn: async ({ signal }) => {
      const entry = sourceEntry(partition, scope.workspaceId);
      entry.epoch = Math.max(entry.epoch, scope.epoch);
      const generation = entry.generation;
      const result = await recordingHomeSourcesAPI.read(scope.workspaceId, signal, context.backend, context.apiOrigin);
      if (signal.aborted || entry.epoch !== scope.epoch || generation !== entry.generation || partition.entries.get(scope.workspaceId) !== entry) {
        throw new Error('Source observation superseded');
      }
      if (result.state === 'available') {
        entry.observation = result;
        partition.entries.delete(scope.workspaceId); partition.entries.set(scope.workspaceId, entry);
        notifySources(partition);
      }
      return { result, generation };
    }, retry: false, refetchOnWindowFocus: false, gcTime: 0,
  });
  const fresh = !query.isFetching && query.isSuccess && query.data.generation === observedEntry?.generation && observedEntry?.epoch === scope.epoch;
  const evidence: RecordingHomeEvidence<RecordingHomeSourceSummary> = query.isFetching || query.isPending
    ? { state: 'loading', scope } : query.isError ? { state: 'error', scope }
      : fresh && query.data?.result.state === 'available' ? { state: 'available', scope, value: query.data.result }
        : { state: 'unknown', scope };
  const observed = observedEntry?.observation;
  const status = query.isError ? 'error' : query.isFetching || query.isPending ? (observed ? 'refreshing' : 'loading')
    : fresh && query.data?.result.state === 'available' ? 'fresh' : observed ? 'stale' : 'unknown';
  return { ...query, observed, evidence, fresh: fresh && query.data?.result.state === 'available', status, context };
}
