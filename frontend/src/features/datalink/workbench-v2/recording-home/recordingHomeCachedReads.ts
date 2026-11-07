import type { QueryClient } from '@tanstack/react-query';
import { writeGroupLifecycleQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroupLifecycle';
import { writeGroupDeliveryQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import type { RecordingHomeGroupScope } from './recordingHomeProjection';
/** No observers or requests. Exact keys retain previously observed evidence only;
 * callers must label it observed, never a complete/current global health index. */
/** A read that actually failed under this exact scope key (workspace, revisions, epoch). It is
 * never success; it only keeps "read and failed" distinct from "never requested". */
const failed = (state: { status: string; fetchStatus: string } | undefined) => state?.status === 'error' && state.fetchStatus === 'idle';
export function cachedHomeReads(queryClient: QueryClient, scopes: readonly RecordingHomeGroupScope[]) {
  const lifecycle = Object.fromEntries(scopes.map(scope => {
    const options = writeGroupLifecycleQueryOptions(scope.workspaceId, scope.groupId, false, scope);
    const cached = queryClient.getQueryState<Awaited<ReturnType<typeof options.queryFn>>>(options.queryKey);
    return [scope.groupId, cached?.status === 'success' && cached.fetchStatus === 'idle' && !cached.isInvalidated && cached.data !== undefined
      ? { state: 'available' as const, scope, value: cached.data } : failed(cached) ? { state: 'error' as const, scope } : { state: 'unknown' as const, scope }];
  }));
  const delivery = Object.fromEntries(scopes.map(scope => {
    const options = writeGroupDeliveryQueryOptions(scope.groupId, false, undefined, scope);
    const cached = queryClient.getQueryState<Awaited<ReturnType<typeof options.queryFn>>>(options.queryKey);
    return [scope.groupId, cached?.status === 'success' && cached.fetchStatus === 'idle' && !cached.isInvalidated && cached.data !== undefined
      ? { state: 'available' as const, scope, value: cached.data } : failed(cached) ? { state: 'error' as const, scope } : { state: 'unknown' as const, scope }];
  }));
  return { lifecycle, delivery };
}
