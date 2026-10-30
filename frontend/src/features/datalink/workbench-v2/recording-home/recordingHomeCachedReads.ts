import type { QueryClient } from '@tanstack/react-query';
import { writeGroupLifecycleQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroupLifecycle';
import { writeGroupDeliveryQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import type { RecordingHomeGroupScope } from './recordingHomeProjection';
/** No observers or requests. Exact keys retain previously observed evidence only;
 * callers must label it observed, never a complete/current global health index. */
export function cachedHomeReads(queryClient: QueryClient, scopes: readonly RecordingHomeGroupScope[]) {
  const lifecycle = Object.fromEntries(scopes.map(scope => {
    const options = writeGroupLifecycleQueryOptions(scope.workspaceId, scope.groupId, false, scope);
    const cached = queryClient.getQueryState<Awaited<ReturnType<typeof options.queryFn>>>(options.queryKey);
    return [scope.groupId, cached?.status === 'success' && cached.fetchStatus === 'idle' && !cached.isInvalidated && cached.data !== undefined
      ? { state: 'available' as const, scope, value: cached.data } : { state: 'unknown' as const, scope }];
  }));
  const delivery = Object.fromEntries(scopes.map(scope => {
    const options = writeGroupDeliveryQueryOptions(scope.groupId, false, undefined, scope);
    const cached = queryClient.getQueryState<Awaited<ReturnType<typeof options.queryFn>>>(options.queryKey);
    return [scope.groupId, cached?.status === 'success' && cached.fetchStatus === 'idle' && !cached.isInvalidated && cached.data !== undefined
      ? { state: 'available' as const, scope, value: cached.data } : { state: 'unknown' as const, scope }];
  }));
  return { lifecycle, delivery };
}
