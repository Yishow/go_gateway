/** Ephemeral inspection of an existing owner, never server-wide absence or mutation permission. */
export type RecordingOwner = 'schema' | 'start' | 'test_write';
export interface RecordingOwnerScope {
  readonly workspaceId: string;
  readonly workspaceRevision: string;
  readonly deviceId: string;
  readonly circuitId?: string;
  readonly groupId?: string;
  readonly groupRevision?: string;
  readonly appliedRevision?: string;
}
export type RecordingOwnerUnavailableReason = 'loading' | 'unreadable' | 'invalid_record' | 'scope_mismatch' | 'revision_mismatch';
export type RecordingOwnerUnresolvedReason = 'pending' | 'running' | 'unknown' | 'lost_ack' | 'query_failed' | 'cleanup_pending';
export type RecordingOwnerObservationState = (
  | { readonly state: 'unresolved'; readonly reason: RecordingOwnerUnresolvedReason; readonly operationId?: string }
  | { readonly state: 'observed'; readonly result: 'none' }
  | { readonly state: 'observed'; readonly result: 'resolved'; readonly status: 'succeeded' | 'failed'; readonly operationId?: string; readonly reason: 'terminal' }
  | { readonly state: 'unavailable'; readonly reason: RecordingOwnerUnavailableReason }
);
export type RecordingOwnerObservation = { readonly owner: RecordingOwner; readonly scope: RecordingOwnerScope } & RecordingOwnerObservationState;
export type RecordingOwnerObservations = readonly RecordingOwnerObservation[];
/** A preview already rendered by its real owner; tokens, statements and consent remain there. */
export interface RecordingSchemaReview {
  readonly scope: RecordingOwnerScope;
  readonly stage: 'idle' | 'previewed' | 'reviewed';
  readonly operationId?: string;
}
export type RecordingOwnerInspection =
  | { readonly state: 'healthy' }
  | { readonly state: 'unresolved'; readonly owner: RecordingOwner; readonly operationId?: string; readonly reason: RecordingOwnerUnresolvedReason }
  | { readonly state: 'unavailable'; readonly owner?: RecordingOwner; readonly reason: RecordingOwnerUnavailableReason };

const owners: readonly RecordingOwner[] = ['schema', 'start', 'test_write'];
const unresolvedReasons: readonly RecordingOwnerUnresolvedReason[] = ['pending', 'running', 'unknown', 'lost_ack', 'query_failed', 'cleanup_pending'];

/** Group-scoped schema evidence is exact; device owners may omit the group fields. */
export function recordingOwnerScopeReason(actual: RecordingOwnerScope, expected: Partial<RecordingOwnerScope>, groupScoped: boolean): RecordingOwnerUnavailableReason | undefined {
  if (!expected.workspaceId || !expected.deviceId || actual.workspaceId !== expected.workspaceId ||
    actual.deviceId !== expected.deviceId || actual.circuitId !== expected.circuitId) return 'scope_mismatch';
  if (!expected.workspaceRevision || actual.workspaceRevision !== expected.workspaceRevision) return 'revision_mismatch';
  if (groupScoped || actual.groupId !== undefined) {
    if (actual.groupId !== expected.groupId) return 'scope_mismatch';
    if (actual.groupRevision !== expected.groupRevision || actual.appliedRevision !== expected.appliedRevision) return 'revision_mismatch';
  }
}

/** No absence is inferred: every required owner must have inspected its current retained state. */
export function inspectRecordingOwners(observations: RecordingOwnerObservations | undefined, scope: Partial<RecordingOwnerScope>): RecordingOwnerInspection {
  const seen = new Map<RecordingOwner, RecordingOwnerObservation>();
  let unavailable: Extract<RecordingOwnerInspection, { state: 'unavailable' }> | undefined;
  let unresolved: Extract<RecordingOwnerInspection, { state: 'unresolved' }> | undefined;
  for (const observation of observations ?? []) {
    const { owner } = observation;
    if (!owners.includes(owner)) { unavailable ??= { state: 'unavailable', reason: 'invalid_record' }; continue; }
    // A complete parent report may include other actual schema groups; they do not satisfy this group's coverage.
    if (owner === 'schema' && observation.scope.groupId !== scope.groupId && observation.scope.groupId !== undefined &&
      observation.scope.workspaceId === scope.workspaceId && observation.scope.deviceId === scope.deviceId &&
      observation.scope.circuitId === scope.circuitId) continue;
    const mismatch = recordingOwnerScopeReason(observation.scope, scope, owner === 'schema');
    if (mismatch) { unavailable ??= { state: 'unavailable', owner, reason: mismatch }; continue; }
    if (seen.has(owner)) unavailable ??= { state: 'unavailable', owner, reason: 'invalid_record' };
    seen.set(owner, observation);
    if (observation.state === 'unresolved' && unresolvedReasons.includes(observation.reason) &&
      (observation.operationId === undefined || typeof observation.operationId === 'string')) {
      unresolved ??= { state: 'unresolved', owner, reason: observation.reason,
        ...(observation.operationId ? { operationId: observation.operationId } : {}) };
    } else if (observation.state === 'unavailable') {
      const reason = ['loading', 'unreadable', 'invalid_record', 'scope_mismatch', 'revision_mismatch'].includes(observation.reason)
        ? observation.reason : 'invalid_record';
      unavailable ??= { state: 'unavailable', owner, reason };
    } else if (observation.state !== 'observed' || (observation.result !== 'none' &&
      !(observation.result === 'resolved' && ['succeeded', 'failed'].includes(observation.status) && observation.reason === 'terminal'))) {
      unavailable ??= { state: 'unavailable', owner, reason: 'invalid_record' };
    }
  }
  // Known retained work must not disappear behind an unreadable/empty storage observation.
  if (unresolved) return unresolved;
  if (unavailable) return unavailable;
  const missing = owners.find(owner => !seen.has(owner));
  return missing ? { state: 'unavailable', owner: missing, reason: 'unreadable' } : { state: 'healthy' };
}
