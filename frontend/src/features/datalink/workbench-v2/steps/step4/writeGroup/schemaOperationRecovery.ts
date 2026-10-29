const pendingOperations = new Map<string, string>();
/** Persist only an opaque ledger identity, never a DDL token or credentials. */
const key = (workspaceId: string, groupId: string) => `wbv2.schema-operation.${JSON.stringify([workspaceId, groupId])}`;
export type SchemaOperationInspection =
  | { readonly state: 'retained'; readonly operationId: string }
  | { readonly state: 'empty' }
  | { readonly state: 'unavailable'; readonly reason: 'missing_scope' | 'unreadable' };
/** Additive read-only inspection. Opaque nonempty values retain their original lookup identity. */
export function inspectSchemaOperation(workspaceId: string, groupId?: string): SchemaOperationInspection {
  if (!workspaceId || !groupId) return { state: 'unavailable', reason: 'missing_scope' };
  const pending = pendingOperations.get(key(workspaceId, groupId));
  if (pending) return { state: 'retained', operationId: pending };
  try {
    const stored = sessionStorage.getItem(key(workspaceId, groupId));
    return stored ? { state: 'retained', operationId: stored } : { state: 'empty' };
  } catch { return { state: 'unavailable', reason: 'unreadable' }; }
}
export function readSchemaOperation(workspaceId: string, groupId?: string): string | undefined {
  if (!groupId) return undefined;
  const pending = pendingOperations.get(key(workspaceId, groupId));
  if (pending) return pending;
  try { return sessionStorage.getItem(key(workspaceId, groupId)) || undefined; } catch { return undefined; }
}
export function retainSchemaOperation(workspaceId: string, groupId: string, operationId: string | undefined, expectedOperationId?: string) {
  // A delayed terminal result may release only its own operation's fence.
  if (!operationId && (!expectedOperationId || readSchemaOperation(workspaceId, groupId) !== expectedOperationId)) return;
  if (operationId) pendingOperations.set(key(workspaceId, groupId), operationId);
  else pendingOperations.delete(key(workspaceId, groupId));
  try {
    if (operationId) sessionStorage.setItem(key(workspaceId, groupId), operationId);
    else sessionStorage.removeItem(key(workspaceId, groupId));
  } catch { /* The in-memory ledger fence remains active when storage is unavailable. */ }
}
