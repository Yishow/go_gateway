const pendingOperations = new Map<string, string>();
/** Persist only an opaque ledger identity, never a DDL token or credentials. */
const key = (workspaceId: string, groupId: string) => `wbv2.schema-operation.${JSON.stringify([workspaceId, groupId])}`;
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
