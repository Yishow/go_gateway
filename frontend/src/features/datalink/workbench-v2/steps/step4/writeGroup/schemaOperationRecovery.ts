const pendingOperations = new Map<string, string>();
/** Persist only an opaque ledger identity, never a DDL token or credentials. */
const key = (workspaceId: string, groupId: string) => `wbv2.schema-operation.${JSON.stringify([workspaceId, groupId])}`;
export function readSchemaOperation(workspaceId: string, groupId?: string): string | undefined {
  if (!groupId) return undefined;
  try { return sessionStorage.getItem(key(workspaceId, groupId)) || pendingOperations.get(key(workspaceId, groupId)); } catch { return pendingOperations.get(key(workspaceId, groupId)); }
}
export function retainSchemaOperation(workspaceId: string, groupId: string, operationId?: string) {
  if (operationId) pendingOperations.set(key(workspaceId, groupId), operationId);
  else pendingOperations.delete(key(workspaceId, groupId));
  try {
    if (operationId) sessionStorage.setItem(key(workspaceId, groupId), operationId);
    else sessionStorage.removeItem(key(workspaceId, groupId));
  } catch { /* The in-memory ledger fence remains active when storage is unavailable. */ }
}
