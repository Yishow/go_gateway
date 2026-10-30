import type { StudioV2Workspace } from '../../../../types/studioV2Workspace';
import type { WriteGroup, WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '../../../../types/studioV2WriteGroupLifecycle';
import type { WriteGroupDelivery } from '../../../../types/studioV2WriteGroupDelivery';
import type { GroupCandidates } from '../state/writeGroup/candidates';
import type { RecordingEvidenceResult, RecordingNextIntent, RecordingReadonly, RecordingUnavailableReason } from '../state/recordingWorkflowView';

/** Query-owner identity, not a saved configuration or a mutation token. */
export interface RecordingHomeScope {
  readonly workspaceId: string;
  readonly workspaceRevision: string;
  /** Incremented by the query owner on scope change, including leave-and-return. */
  readonly epoch: number;
}
/** Every state belongs to its captured request scope, including failure and loading. */
export type RecordingHomeEvidence<T> =
  | { readonly state: 'available'; readonly scope: RecordingHomeScope; readonly value: RecordingReadonly<T> }
  | { readonly state: 'loading' | 'error' | 'unknown' | 'unavailable'; readonly scope: RecordingHomeScope };
export interface RecordingHomeGroupScope extends RecordingHomeScope {
  readonly groupId: string;
  readonly groupRevision: string;
  readonly appliedRevision: string;
}
export type RecordingHomeGroupEvidence<T> =
  | { readonly state: 'available'; readonly scope: RecordingHomeGroupScope; readonly value: RecordingReadonly<T> }
  | { readonly state: 'loading' | 'error' | 'unknown' | 'unavailable'; readonly scope: RecordingHomeGroupScope };
export interface RecordingHomeInput {
  readonly scope: RecordingHomeScope;
  readonly workspace?: RecordingHomeEvidence<StudioV2Workspace>;
  /** Existing saved-source candidate resolver output; never derive readiness here. */
  readonly sources?: RecordingHomeEvidence<GroupCandidates>;
  readonly groups?: RecordingHomeEvidence<WriteGroupListResponse>;
  readonly lifecycle?: Readonly<Record<string, RecordingHomeGroupEvidence<WriteGroupLifecycle>>>;
  readonly delivery?: Readonly<Record<string, RecordingHomeGroupEvidence<WriteGroupDelivery>>>;
}
/** A destination for a later owner; never a command, callback or execution permission. */
export interface RecordingHomeNavigationIntent {
  readonly kind: 'navigate';
  readonly intent: RecordingNextIntent;
}
export interface RecordingHomeGroupProjection {
  readonly group: RecordingReadonly<WriteGroup>;
  readonly lifecycle: RecordingEvidenceResult<WriteGroupLifecycle>;
  readonly delivery: RecordingEvidenceResult<WriteGroupDelivery>;
}
export interface RecordingHomeProjection {
  readonly scope: RecordingHomeScope;
  readonly workspace: RecordingEvidenceResult<StudioV2Workspace>;
  readonly sources: RecordingEvidenceResult<GroupCandidates>;
  readonly groups: RecordingEvidenceResult<readonly RecordingHomeGroupProjection[]>;
  /** Priority classification is a separate consumer of this evidence boundary. */
  readonly primaryAction: RecordingHomeNavigationIntent | null;
}

/** Copy before freezing: query-cache objects remain owned by their existing hooks. */
function snapshot<T>(value: T): RecordingReadonly<T> {
  if (Array.isArray(value)) return Object.freeze(value.map(item => snapshot(item))) as RecordingReadonly<T>;
  if (value !== null && typeof value === 'object') return Object.freeze(Object.fromEntries(
    Object.entries(value).map(([key, item]) => [key, snapshot(item)]))) as RecordingReadonly<T>;
  return value as RecordingReadonly<T>;
}
function invalidScope(scope: RecordingHomeScope): RecordingUnavailableReason | undefined {
  if (!scope.workspaceId) return 'scope_mismatch';
  if (!scope.workspaceRevision) return 'revision_mismatch';
  if (!Number.isSafeInteger(scope.epoch) || scope.epoch < 0) return 'invalid_evidence';
}
function mismatch(actual: RecordingHomeScope, expected: RecordingHomeScope): RecordingUnavailableReason | undefined {
  if (actual.workspaceId !== expected.workspaceId) return 'scope_mismatch';
  if (actual.workspaceRevision !== expected.workspaceRevision) return 'revision_mismatch';
  if (actual.epoch !== expected.epoch) return 'stale';
}
function read<T>(evidence: RecordingHomeEvidence<T> | undefined, scope: RecordingHomeScope,
  validate?: (value: RecordingReadonly<T>) => RecordingUnavailableReason | undefined,
  validateScope?: () => RecordingUnavailableReason | undefined,
): RecordingEvidenceResult<T> {
  const invalid = invalidScope(scope);
  if (invalid) return { state: 'unavailable', reason: invalid };
  if (!evidence) return { state: 'unavailable', reason: 'missing' };
  const scopeReason = mismatch(evidence.scope, scope) ?? validateScope?.();
  if (scopeReason) return { state: 'unavailable', reason: scopeReason };
  if (evidence.state === 'unavailable') return { state: 'unavailable', reason: 'missing' };
  if (evidence.state !== 'available') return { state: 'unavailable', reason: evidence.state };
  const reason = validate?.(evidence.value);
  return reason ? { state: 'unavailable', reason } : { state: 'available', value: evidence.value };
}
function readGroup<T>(evidence: RecordingHomeGroupEvidence<T> | undefined, scope: RecordingHomeScope,
  group: RecordingReadonly<WriteGroup>, validate: (value: RecordingReadonly<T>) => RecordingUnavailableReason | undefined,
): RecordingEvidenceResult<T> {
  return read(evidence, scope, validate, () => {
    if (evidence?.scope.groupId !== group.id) return 'scope_mismatch';
    if (evidence.scope.groupRevision !== group.revision || evidence.scope.appliedRevision !== group.applied_revision) return 'revision_mismatch';
  });
}
function deliveryMismatch(value: RecordingReadonly<WriteGroupDelivery>, group: RecordingReadonly<WriteGroup>): RecordingUnavailableReason | undefined {
  if (value.group_id !== group.id) return 'scope_mismatch';
  // These fields explicitly describe the applied revision; backlog is historical.
  if ((value.revision_stages && value.revision_stages.group_revision !== group.applied_revision) ||
    (value.last_sql_committed_effect && value.last_sql_committed_effect.group_revision !== group.applied_revision) ||
    (value.source_evidence?.state === 'available' && value.source_evidence.group_revision !== group.applied_revision)) return 'revision_mismatch';
}
/** Pure read boundary. Opening cannot save, start, create a table or restart anything. */
export function recordingHomeProjection(input: RecordingHomeInput): RecordingHomeProjection {
  const { scope } = input;
  const workspace = read(input.workspace, scope, value => value.id !== scope.workspaceId ? 'scope_mismatch' : undefined);
  const groups = read(input.groups, scope, value => {
    if (value.workspace_id !== scope.workspaceId) return 'scope_mismatch';
    if (value.workspace_revision !== scope.workspaceRevision) return 'revision_mismatch';
    const ids = new Set<string>();
    for (const group of value.groups) {
      if (!group.id || group.workspace_id !== scope.workspaceId) return 'scope_mismatch';
      if (!group.revision) return 'revision_mismatch';
      if (ids.has(group.id)) return 'invalid_evidence';
      ids.add(group.id);
    }
  });
  return snapshot({ scope, workspace, sources: read(input.sources, scope),
    groups: groups.state === 'available' ? { state: 'available', value: groups.value.groups.map(group => ({
      group, lifecycle: readGroup(input.lifecycle?.[group.id], scope, group, value => {
        if (value.workspace_id !== scope.workspaceId || value.group_id !== group.id) return 'scope_mismatch';
        if (value.workspace_revision !== scope.workspaceRevision || value.group_revision !== group.revision ||
          value.applied_revision !== group.applied_revision) return 'revision_mismatch';
      }),
      delivery: readGroup(input.delivery?.[group.id], scope, group, value => deliveryMismatch(value, group)),
    })) } : groups, primaryAction: null });
}

/**
 * Query adapters call this for every completion, including errors. Capture scope
 * when requesting, never relabel a response with the selection at completion.
 * No state is stored here; matching degraded evidence is as authoritative as data.
 */
export function acceptRecordingHomeResponse<T>(
  scope: RecordingHomeScope, responseScope: RecordingHomeScope, response: RecordingHomeEvidence<T>,
  current?: RecordingHomeEvidence<T>,
): RecordingHomeEvidence<T> | undefined {
  if (invalidScope(scope)) return undefined;
  if (mismatch(responseScope, scope) || mismatch(response.scope, scope)) {
    return current && !mismatch(current.scope, scope) ? current : undefined;
  }
  return response;
}
