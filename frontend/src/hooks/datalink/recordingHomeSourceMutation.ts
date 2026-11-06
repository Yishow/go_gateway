import type { QueryClient } from '@tanstack/react-query';
import { invalidateRecordingHomeSources, recordingHomeWorkspaceOwner, recordingHomeSourceContext, recordingHomeSourcePrefix } from './recordingHomeSourceCache';
/** Capture before any mutation-owner await, including optimistic preflight. */
export function captureRecordingHomeSourceMutation(client: QueryClient, workspaceId?: string) {
  const owner = workspaceId ?? recordingHomeWorkspaceOwner(client);
  // Never resolve a different workspace after awaiting acknowledgement.
  return owner ? { ...recordingHomeSourceContext(), workspaceId: owner } : undefined;
}
/** Successful configuration ACK: stale immediately, active revalidation independent
 * of all existing save/readback barriers. An inactive home reads only on return. */
export function completeRecordingHomeSourceMutation(client: QueryClient, captured: ReturnType<typeof captureRecordingHomeSourceMutation>): void {
  if (!captured) return;
  invalidateRecordingHomeSources(client, captured);
  void client.refetchQueries({ queryKey: recordingHomeSourcePrefix(captured, captured.workspaceId), type: 'active' });
}
/** For owners without asynchronous optimistic preflight. */
export async function observeRecordingHomeSourceMutation<T>(client: QueryClient, operation: () => Promise<T>, workspaceId?: string): Promise<T> {
  const captured = captureRecordingHomeSourceMutation(client, workspaceId);
  const result = await operation();
  completeRecordingHomeSourceMutation(client, captured);
  return result;
}
