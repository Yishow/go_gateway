import { useIsMutating, useMutationState, useQueryClient, type Mutation } from '@tanstack/react-query';
import { normalizeTypedEnvelope } from '../../utils/safeJson';
export const writeGroupSaveMutationKey = ['studio-v2', 'write-group-save'] as const;
/** A transport failure is not proof that the synchronous save was rejected. */
export function writeGroupSaveOutcomeUnknown(error: unknown): boolean {
  if (normalizeTypedEnvelope(error).code === 'revision_mismatch') return false;
  const response = error && typeof error === 'object' && 'response' in error ? error.response : undefined;
  const status = response && typeof response === 'object' && 'status' in response ? response.status : undefined;
  return !(typeof status === 'number' && status >= 400 && status < 500 && status !== 408);
}
function submittedScope(mutation: Mutation) {
  const variables = mutation.state.variables;
  if (!variables || typeof variables !== 'object') return undefined;
  const request = 'request' in variables ? variables.request : variables;
  if (!request || typeof request !== 'object' || !('workspace_id' in request)) return undefined;
  return { workspaceId: request.workspace_id, groupId: 'id' in variables ? variables.id : undefined };
}
/** Read the existing mutation cache, never mirror it into another submission owner. */
export function useWriteGroupSaveActivity(workspaceId: string, groupId?: string) {
  const client = useQueryClient();
  const sameScope = (mutation: Mutation) => {
    const scope = submittedScope(mutation);
    return scope?.workspaceId === workspaceId && scope.groupId === groupId;
  };
  const filters = { mutationKey: writeGroupSaveMutationKey, predicate: sameScope };
  const pending = useIsMutating(filters) > 0;
  const uncertain = useMutationState({ filters: { mutationKey: writeGroupSaveMutationKey, status: 'error',
    predicate: mutation => submittedScope(mutation)?.workspaceId === workspaceId && writeGroupSaveOutcomeUnknown(mutation.state.error) }, select: mutation => mutation.mutationId });
  return { pending, priorUnknown: uncertain.length > 0, isPending: () => client.isMutating(filters) > 0 };
}
