import type { WriteGroupTestWriteOperation } from '../../../../types/studioV2WriteGroupTestWrite';

/** Locator only. Tokens, row values, target names and operation truth stay on the server. */
export interface RecordingTestWriteReference {
  workspace_id: string;
  group_id: string;
  group_revision: string;
  operation_id: string;
  operation_kind: 'test_write';
  created_at: string;
}
export const recordingTestWriteRecoveryKey = 'studio.recording-test-write-recovery.v1';
const fields = ['workspace_id', 'group_id', 'group_revision', 'operation_id', 'operation_kind', 'created_at'] as const;
function parse(value: unknown): RecordingTestWriteReference {
  if (!value || typeof value !== 'object') throw new Error('invalid reference');
  const entry = value as Record<string, unknown>;
  if (fields.some((field) => typeof entry[field] !== 'string' || !entry[field]) || entry.operation_kind !== 'test_write') {
    throw new Error('invalid reference');
  }
  return Object.fromEntries(fields.map((field) => [field, entry[field]])) as unknown as RecordingTestWriteReference;
}

/** No write is considered resolved solely by absence, a transport error, or a terminal status. */
export function testWriteResolved(operation: WriteGroupTestWriteOperation): boolean {
  if (operation.status === 'pending' || operation.status === 'running') return false;
  return (operation.write_outcome === 'failed' && operation.cleanup_status === 'not_attempted') ||
    (operation.cleanup_status === 'cleaned' && operation.write_outcome !== undefined && operation.write_outcome !== 'unknown');
}

export function createRecordingTestWriteRecovery(storage: () => Storage) {
  let lastReadable: RecordingTestWriteReference[] = [];
  const fallback = new Map<string, RecordingTestWriteReference>();
  const listeners = new Set<() => void>();
  const notify = () => listeners.forEach((listener) => listener());
  function read() {
    try {
      const raw = storage().getItem(recordingTestWriteRecoveryKey);
      const parsed: unknown = raw ? JSON.parse(raw) : [];
      if (!Array.isArray(parsed)) throw new Error('invalid references');
      const entries = new Map(parsed.map((value) => { const ref = parse(value); return [ref.operation_id, ref]; }));
      if (entries.size !== parsed.length) throw new Error('duplicate reference');
      // A successful read can still have lost entries (clear/remove/another tab).
      // Only an explicit exact-operation resolution is allowed to forget them.
      for (const known of lastReadable) {
        const incoming = entries.get(known.operation_id);
        if (incoming && (incoming.workspace_id !== known.workspace_id || incoming.group_id !== known.group_id ||
          incoming.group_revision !== known.group_revision)) throw new Error('conflicting reference');
        if (!incoming) entries.set(known.operation_id, known);
      }
      lastReadable = [...entries.values()];
      fallback.forEach((ref, id) => entries.set(id, ref));
      return { references: [...entries.values()], unavailable: fallback.size > 0, unreadable: false };
    } catch {
      const known = new Map(lastReadable.map((ref) => [ref.operation_id, ref]));
      fallback.forEach((ref, id) => known.set(id, ref));
      return { references: [...known.values()], unavailable: true, unreadable: true };
    }
  }
  function retain(reference: RecordingTestWriteReference): boolean {
    const current = read();
    const safe = parse(reference);
    fallback.set(safe.operation_id, safe);
    let persisted = false;
    try {
      if (current.unavailable) throw new Error('unavailable references');
      const entries = current.references.filter((ref) => ref.operation_id !== safe.operation_id);
      storage().setItem(recordingTestWriteRecoveryKey, JSON.stringify([...entries, safe]));
      lastReadable = [...entries, safe];
      fallback.delete(safe.operation_id);
      persisted = true;
    } catch { /* Keep the reference in memory and prevent dispatch. */ }
    notify();
    return persisted;
  }
  function resolve(operation: WriteGroupTestWriteOperation, workspaceId?: string) {
    if (!testWriteResolved(operation)) return;
    const current = read();
    if (current.unreadable || (workspaceId !== undefined && current.references.some((ref) =>
      ref.operation_id === operation.operation_id && ref.workspace_id !== workspaceId))) return;
    try {
      storage().setItem(recordingTestWriteRecoveryKey, JSON.stringify(current.references.filter((ref) => ref.operation_id !== operation.operation_id)));
      lastReadable = current.references.filter((ref) => ref.operation_id !== operation.operation_id);
      fallback.delete(operation.operation_id);
    } catch { /* Failure to persist resolution keeps the unresolved locator. */ }
    notify();
  }
  return { read, retain, resolve, subscribe: (listener: () => void) => {
    listeners.add(listener);
    return () => { listeners.delete(listener); };
  } };
}
export const recordingTestWriteRecovery = createRecordingTestWriteRecovery(() => localStorage);
