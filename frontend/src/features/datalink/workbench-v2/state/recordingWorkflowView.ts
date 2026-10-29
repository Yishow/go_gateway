import { inspectRecordingOwners, recordingOwnerScopeReason, type RecordingOwner, type RecordingOwnerObservations, type RecordingSchemaReview } from './recordingOwnerObservation';
export type { RecordingOwner, RecordingOwnerObservation, RecordingOwnerObservations, RecordingOwnerScope, RecordingSchemaReview } from './recordingOwnerObservation';
import type { SchemaOperation, SchemaPreviewToken } from '../../../../types/recordingPlan';
import type { RecordingStartOperation } from '../../../../types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupDraft, WriteGroupReadiness } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupTestWriteOperation } from '../../../../types/studioV2WriteGroupTestWrite';
import { parseRuntimeRuleEvidence } from '../../../../services/runtimeRuleEvidenceParser';
import type { RecordingGroupFacts, RecordingGroupView, RuleEvidenceView } from '../../../../types/runtimeRuleEvidence';
import type { CandidateCounterStatus } from './circuitProfileBasic';
import { isConfirmedCounterDefinition } from './circuitCounterConfirmation';

export type RecordingReadonly<T> = T extends object
  ? { readonly [Key in keyof T]: RecordingReadonly<T[Key]> } : T;

/**
 * Caller selection, not a locator or mutation request. IDs/revisions come from
 * the canonical workspace/group queries; names never establish ownership.
 * workspaceRevision is the WriteGroupListResponse revision, not a local counter.
 * Device/circuit identifies the viewing context, not exclusive group ownership.
 * A custom group's full membership remains visible, including other devices.
 */
export interface RecordingWorkflowScope {
  readonly workspaceId: string;
  readonly deviceId: string;
  readonly circuitId?: string;
  readonly groupId: string;
  readonly workspaceRevision: string;
  readonly groupRevision: string;
  /** Latest requested applied revision; it may be scheduled and is NOT the as-of effective revision. */
  readonly appliedRevision: string;
}

export type RecordingUnavailableReason = 'missing' | 'loading' | 'error' | 'unknown'
  | 'scope_mismatch' | 'revision_mismatch' | 'effective_snapshot_unavailable' | 'operation_unresolved' | 'stale' | 'invalid_evidence';
export type RecordingUnavailable = { readonly state: 'unavailable'; readonly reason: RecordingUnavailableReason };
export type RecordingEvidenceResult<T> = RecordingUnavailable
  | { readonly state: 'available'; readonly value: RecordingReadonly<T> };

/** Read-only query evidence, bound by its adapter to the exact selected scope. */
export type RecordingWorkflowEvidence<T> =
  | { readonly state: 'available'; readonly scope: RecordingWorkflowScope; readonly value: RecordingReadonly<T> }
  | { readonly state: 'loading' | 'error' | 'unknown' | 'unavailable' };

export type RecordingWorkflowOperation =
  | { readonly kind: 'none' }
  | { readonly kind: 'schema'; readonly operation: RecordingReadonly<SchemaOperation> }
  | { readonly kind: 'start'; readonly operation: RecordingReadonly<RecordingStartOperation> }
  | { readonly kind: 'test_write'; readonly operation: RecordingReadonly<WriteGroupTestWriteOperation> };

export type RecordingLifecycleAction = 'edit' | 'stop' | 'remove' | 'transition' | 'recreate';
/**
 * Caller projection of server lifecycle availability, NOT a backend DTO or a
 * replacement policy. Future adapters must carry server blockers unchanged;
 * missing actions are unavailable, never inferred from status/readiness booleans.
 */
export type RecordingLifecycleAvailability = Readonly<Partial<Record<RecordingLifecycleAction,
  { readonly state: 'available' | 'blocked'; readonly reasonCodes: readonly string[] }
>>>;

/** Only the persisted basic role proves exclusivity; all other configuration stays group-scoped. */
export type RecordingConfigurationScope =
  | { readonly kind: 'managed_device'; readonly deviceId: string }
  | { readonly kind: 'group'; readonly groupId: string };

export interface SavedRecordingFacts {
  readonly kind: 'saved';
  readonly configurationScope: RecordingConfigurationScope;
  readonly group: RecordingReadonly<WriteGroup>;
}
export interface EffectiveRecordingFacts {
  readonly kind: 'effective';
  /** Safe facts do not prove the historical basic role; never borrow it from a newer saved draft. */
  readonly configurationScope: { readonly kind: 'group'; readonly groupId: string };
  readonly revision: string;
  readonly asOf: string;
  readonly facts: RecordingReadonly<RecordingGroupFacts>;
}

export interface RecordingWorkflowInput {
  readonly scope: RecordingWorkflowScope;
  /** Persisted draft only; unsaved editor state must remain with its existing owner. */
  readonly savedGroup: RecordingReadonly<WriteGroup>;
  /** Runtime's safe, revision/as-of-bound DTO, never the editable GET as an effective reader. */
  readonly runtimeEvidence?: RecordingWorkflowEvidence<RuleEvidenceView>;
  /** WriteGroupReadiness is the sole readiness authority; no local ready computation. */
  readonly readiness?: RecordingWorkflowEvidence<WriteGroupReadiness>;
  readonly schema?: RecordingWorkflowEvidence<SchemaPreviewToken>;
  /** Legacy isolated operation evidence; kind:none is not global absence. Production uses ownerObservations. */
  readonly operation?: RecordingWorkflowEvidence<RecordingWorkflowOperation>;
  readonly lifecycle?: RecordingWorkflowEvidence<RecordingLifecycleAvailability>;
  readonly ownerObservations?: RecordingOwnerObservations;
}

export type RecordingWorkflowView = RecordingUnavailable | {
  readonly state: 'available';
  readonly scope: RecordingWorkflowScope;
  readonly saved: SavedRecordingFacts;
  readonly runtime: RecordingEvidenceResult<RecordingGroupView>;
  readonly effective: RecordingEvidenceResult<EffectiveRecordingFacts>;
  readonly readiness: RecordingEvidenceResult<WriteGroupReadiness>;
  readonly schema: RecordingEvidenceResult<SchemaPreviewToken>;
  readonly operation: RecordingEvidenceResult<RecordingWorkflowOperation>;
  /** Usable server evidence only, not permission to execute a mutation or a next-action decision. */
  readonly lifecycle: RecordingEvidenceResult<RecordingLifecycleAvailability>;
};

function unavailable(reason: RecordingUnavailableReason): RecordingUnavailable {
  return Object.freeze({ state: 'unavailable', reason });
}

/** Canonical inputs are JSON data; copy nested arrays/objects, never freeze query cache objects. */
function copyFrozen<T>(value: T): RecordingReadonly<T> {
  if (Array.isArray(value)) return Object.freeze(value.map((item: unknown) => copyFrozen(item))) as RecordingReadonly<T>;
  if (value !== null && typeof value === 'object') {
    return Object.freeze(Object.fromEntries(Object.entries(value).map(([key, item]: [string, unknown]) =>
      [key, copyFrozen(item)]))) as RecordingReadonly<T>;
  }
  return value as RecordingReadonly<T>;
}

function scopeMismatch(actual: RecordingWorkflowScope, expected: RecordingWorkflowScope): RecordingUnavailableReason | undefined {
  if (actual.workspaceId !== expected.workspaceId || actual.deviceId !== expected.deviceId ||
    actual.circuitId !== expected.circuitId || actual.groupId !== expected.groupId) return 'scope_mismatch';
  if (actual.workspaceRevision !== expected.workspaceRevision || actual.groupRevision !== expected.groupRevision ||
    actual.appliedRevision !== expected.appliedRevision) return 'revision_mismatch';
}

function hasGroupIdentity(group: RecordingReadonly<WriteGroup>, scope: RecordingWorkflowScope): boolean {
  if (group.workspace_id !== scope.workspaceId || group.id !== scope.groupId ||
    (group.basic_managed_device_id !== undefined && group.basic_managed_device_id !== scope.deviceId)) return false;
  return group.basic_managed_device_id === undefined ||
    group.members.every((member) => member.device_id === group.basic_managed_device_id);
}

function selectsGroup(group: RecordingReadonly<WriteGroup>, scope: RecordingWorkflowScope): boolean {
  if (group.members.length === 0) return group.basic_managed_device_id === scope.deviceId && scope.circuitId === undefined;
  return group.members.some((member) => member.device_id === scope.deviceId &&
    (scope.circuitId === undefined || member.circuit_id === scope.circuitId));
}

function configurationScope(group: RecordingReadonly<WriteGroup>): RecordingConfigurationScope {
  return group.basic_managed_device_id === undefined ? { kind: 'group', groupId: group.id }
    : { kind: 'managed_device', deviceId: group.basic_managed_device_id };
}

function readEvidence<T>(
  evidence: RecordingWorkflowEvidence<T> | undefined,
  scope: RecordingWorkflowScope,
  validate?: (value: RecordingReadonly<T>) => RecordingUnavailableReason | undefined,
): RecordingEvidenceResult<T> {
  if (!evidence || evidence.state === 'unavailable') return unavailable('missing');
  if (evidence.state !== 'available') return unavailable(evidence.state);
  const reason = scopeMismatch(evidence.scope, scope) ?? validate?.(evidence.value);
  return reason ? unavailable(reason) : { state: 'available', value: evidence.value };
}

function runtimeFacts(input: RecordingWorkflowInput): RecordingEvidenceResult<RecordingGroupView> {
  const evidence = readEvidence(input.runtimeEvidence, input.scope);
  if (evidence.state !== 'available') return evidence;
  const parsed = parseRuntimeRuleEvidence(evidence.value, { workspace_id: input.scope.workspaceId,
    workspace_revision: input.scope.workspaceRevision, device_id: input.scope.deviceId });
  if (!parsed.ok) return unavailable(parsed.reason === 'foreign-scope' ? 'scope_mismatch'
    : parsed.reason === 'stale-reference' ? 'revision_mismatch' : 'invalid_evidence');
  if (parsed.value.freshness !== 'fresh') return unavailable(parsed.value.freshness === 'stale' ? 'stale' : 'missing');
  const group = parsed.value.groups.find(item => item.group_id === input.scope.groupId);
  if (!group) return unavailable('scope_mismatch');
  if (group.saved_revision !== input.scope.groupRevision || group.applied_revision !== (input.scope.appliedRevision || null)) {
    return unavailable('revision_mismatch');
  }
  return { state: 'available', value: group };
}

function effectiveFacts(runtime: RecordingEvidenceResult<RecordingGroupView>, scope: RecordingWorkflowScope): RecordingEvidenceResult<EffectiveRecordingFacts> {
  if (runtime.state !== 'available') return runtime.reason === 'missing' ? unavailable('effective_snapshot_unavailable') : runtime;
  const value = runtime.value;
  // Runtime's unique parser validates immutable layout, revision and nanosecond as-of proof.
  // Never compare effective_revision to latest requested applied_revision or reconstruct a WriteGroup.
  if (value.state !== 'effective' || value.proof !== 'applied-snapshot-and-layout' || !value.effective || !value.effective_revision || !value.as_of) {
    return unavailable('effective_snapshot_unavailable');
  }
  return { state: 'available', value: { kind: 'effective', configurationScope: { kind: 'group', groupId: scope.groupId },
    revision: value.effective_revision, asOf: value.as_of, facts: value.effective } };
}

function selectsEffective(facts: RecordingReadonly<RecordingGroupFacts>, scope: RecordingWorkflowScope): boolean {
  return facts.members.some(member => member.device_id === scope.deviceId &&
    (scope.circuitId === undefined || member.circuit_id === scope.circuitId));
}

function readinessFacts(input: RecordingWorkflowInput): RecordingEvidenceResult<WriteGroupReadiness> {
  return readEvidence(input.readiness, input.scope, (value) => {
    if (value.workspace_id !== input.scope.workspaceId || value.group_id !== input.scope.groupId) return 'scope_mismatch';
    if (value.workspace_revision !== input.scope.workspaceRevision || value.group_revision !== input.scope.groupRevision ||
      value.applied_revision !== input.scope.appliedRevision) return 'revision_mismatch';
  });
}

function schemaFacts(input: RecordingWorkflowInput): RecordingEvidenceResult<SchemaPreviewToken> {
  return readEvidence(input.schema, input.scope, (value) => {
    if (value.workspace_id !== input.scope.workspaceId || value.plan_id !== input.scope.groupId ||
      value.connector_id !== input.savedGroup.destination.connector_id) return 'scope_mismatch';
    if (value.workspace_revision !== input.scope.workspaceRevision || value.plan_revision !== input.scope.groupRevision ||
      !value.connector_revision || value.connector_revision !== input.savedGroup.destination.connector_revision) return 'revision_mismatch';
  });
}

function operationFacts(input: RecordingWorkflowInput): RecordingEvidenceResult<RecordingWorkflowOperation> {
  return readEvidence(input.operation, input.scope, (value) => {
    if (value.kind === 'none') return;
    if (!value.operation.operation_id) return 'missing';
    // Schema/test-write DTOs lack identity fields; their query adapter owns the envelope binding.
    if (value.kind !== 'start') return;
    const { operation } = value;
    const group = operation.groups.find((item) => item.group_id === input.scope.groupId);
    if (operation.workspace_id !== input.scope.workspaceId || !operation.device_ids.includes(input.scope.deviceId) || !group) {
      return 'scope_mismatch';
    }
    if (operation.setup_revision !== input.scope.workspaceRevision || group.group_revision !== input.scope.groupRevision ||
      (group.applied_revision !== undefined && group.applied_revision !== input.scope.appliedRevision)) return 'revision_mismatch';
  });
}

function operationUnresolved(evidence: RecordingEvidenceResult<RecordingWorkflowOperation>): boolean {
  if (evidence.state !== 'available') return true;
  const value = evidence.value;
  if (value.kind === 'none') return false;
  if (value.operation.status !== 'succeeded' && value.operation.status !== 'failed') return true;
  if (value.kind !== 'test_write') return false;
  const { write_outcome, cleanup_status } = value.operation;
  // A terminal envelope never resolves an unknown write or outstanding cleanup.
  return write_outcome === undefined || write_outcome === 'unknown' ||
    !(cleanup_status === 'cleaned' || (cleanup_status === 'not_attempted' && write_outcome === 'failed'));
}

/**
 * Pure, unmounted presentation boundary: no callbacks, I/O, state, or activation.
 * It validates caller evidence, not readiness or the next primary action. An
 * unresolved operation withholds lifecycle eligibility without rewriting its
 * status. Known terminal results only permit forwarding server availability; they
 * never prove SQL commit, readback, cleanup, or that the saved draft is running.
 */
export function recordingWorkflowView(input: RecordingWorkflowInput): RecordingWorkflowView {
  const { scope, savedGroup } = input;
  if (!scope.workspaceId || !scope.deviceId || !scope.groupId || scope.circuitId === '' || !hasGroupIdentity(savedGroup, scope)) {
    return unavailable('scope_mismatch');
  }
  if (!scope.workspaceRevision || !scope.groupRevision || savedGroup.revision !== scope.groupRevision ||
    savedGroup.applied_revision !== scope.appliedRevision) return unavailable('revision_mismatch');
  const runtime = runtimeFacts(input);
  const effective = effectiveFacts(runtime, scope);
  // A draft may remove an active source. Match either proven version, retaining both complete member sets.
  if (!selectsGroup(savedGroup, scope) && !(effective.state === 'available' && selectsEffective(effective.value.facts, scope))) {
    return unavailable('scope_mismatch');
  }
  const operation = operationFacts(input);
  let lifecycle = readEvidence(input.lifecycle, scope);
  const owners = input.ownerObservations === undefined ? undefined : inspectRecordingOwners(input.ownerObservations, scope);
  if (lifecycle.state === 'available' && ((owners !== undefined && owners.state !== 'healthy') || (operationUnresolved(operation) &&
    !(operation.state === 'unavailable' && operation.reason === 'missing' && owners?.state === 'healthy')))) {
    lifecycle = unavailable('operation_unresolved');
  }
  return copyFrozen({
    state: 'available', scope, saved: { kind: 'saved', configurationScope: configurationScope(savedGroup), group: savedGroup }, runtime, effective,
    readiness: readinessFacts(input), schema: schemaFacts(input), operation, lifecycle,
  } as const);
}

export type RecordingDraftScope = Pick<RecordingWorkflowScope, 'workspaceId' | 'deviceId' | 'circuitId'> & { readonly workspaceRevision?: string };
/** save_prepare opens existing save/preparation review; it never combines saving with DDL or starts. */
export type RecordingNextIntentKind = 'choose_sources' | 'resolve_semantics' | 'save_prepare' | 'preview_schema' | 'review_schema'
  | 'confirm_schema' | 'review_scoped_start' | 'query_operation' | 'recover_operation' | 'refresh_evidence' | 'monitor' | 'inspect';
export type RecordingNextReason = 'sources_empty' | 'counter_confirmation' | 'save_required' | 'offline'
  | 'operation_recovery' | 'evidence_loading' | 'evidence_error' | 'evidence_unavailable' | 'schema_review'
  | 'schema_confirmation' | 'schema_conflict' | 'configuration_blocked' | 'readiness_blocked'
  | 'lifecycle_unavailable' | 'lifecycle_blocked' | 'start_review' | 'monitor_available' | 'scope_unavailable';

interface RecordingNextCommon {
  readonly connection: 'online' | 'offline';
  /** Legacy isolated callers may retain a clock; production preview validity stays with its owner. */
  readonly nowMs?: number;
  readonly ownerObservations?: RecordingOwnerObservations;
  readonly schemaReview?: RecordingSchemaReview;
  /** Existing authoring resolver results, not a replacement server-readiness model. */
  readonly counterStatuses?: RecordingReadonly<CandidateCounterStatus[]>;
  /** Retained/unconfirmed request from the existing recovery owner; never a new operation. */
  readonly unconfirmedOperation?: {
    readonly scope: RecordingDraftScope & { readonly groupId?: string };
    readonly operationId?: string;
  };
}

export type RecordingNextActionInput = RecordingNextCommon & (
  | { readonly kind: 'draft'; readonly scope: RecordingDraftScope; readonly draft: RecordingReadonly<WriteGroupDraft> }
  | { readonly kind: 'saved'; readonly scope: RecordingWorkflowScope; readonly workflow: RecordingWorkflowView;
    /** Unsaved editor-owned changes; presence means save/review again, never reuse saved readiness. */
    readonly draft?: RecordingReadonly<WriteGroupDraft> }
);

/** Explicit navigation/review intent only. Integration must revalidate CAS, consent and every server guard. */
export interface RecordingNextIntent {
  readonly kind: RecordingNextIntentKind;
  readonly scope: RecordingDraftScope | RecordingWorkflowScope;
  readonly operationId?: string;
  readonly operationOwner?: RecordingOwner;
  readonly activatedDeviceIds?: readonly string[];
}
export interface RecordingNextActionModel {
  readonly intent: RecordingNextIntent;
  readonly enabled: boolean;
  readonly reasons: readonly RecordingNextReason[];
  readonly configurationScope: RecordingConfigurationScope | { readonly kind: 'draft' };
}

function needsCounterConfirmation(
  draft: RecordingReadonly<WriteGroupDraft>, statuses: RecordingNextActionInput['counterStatuses'],
): boolean {
  return draft.members.some((member) => {
    if (draft.recording_profile?.kind !== 'cumulative' && member.value_role !== 'counter_raw') return false;
    const status = statuses?.find(({ candidate }) => candidate.device_id === member.device_id &&
      candidate.point_id === member.point_id && candidate.tag_id === member.tag_id);
    const proof = status?.measurement;
    return !status?.confirmed || !proof || proof.workspace_id !== draft.workspace_id ||
      proof.device_id !== member.device_id || proof.point_id !== member.point_id || proof.tag_id !== member.tag_id ||
      (!!member.measurement_id && member.measurement_id !== proof.id) || !isConfirmedCounterDefinition(proof);
  });
}

function hasMonitorHandoff(view: Extract<RecordingWorkflowView, { state: 'available' }>): boolean {
  if (view.effective.state !== 'available' || view.runtime.state !== 'available' ||
    view.saved.group.revision !== view.effective.value.revision || !selectsEffective(view.effective.value.facts, view.scope) ||
    !['ready', 'running'].includes(view.saved.group.status) || view.operation.state !== 'available' || view.operation.value.kind !== 'start') return false;
  const revision = view.effective.value.revision;
  const intake = view.runtime.value.delivery?.current?.intake;
  // A succeeded operation may only have scheduled a revision. Pipeline intake is separate authority;
  // active intake grants no SQL receipt/readback/cleanup claim and never authorizes a new start.
  const operation = view.operation.value.operation;
  return intake?.state === 'active' && intake.group_revision === revision && operation.status === 'succeeded' &&
    operation.devices.some(device => device.device_id === view.scope.deviceId && device.activated) &&
    operation.groups.some(group => group.group_id === view.scope.groupId && group.applied && group.applied_revision === revision);
}

/** One ephemeral next intent. No stored phase, default readiness, API call, or implicit consent. */
export function recordingNextAction(input: RecordingNextActionInput): RecordingNextActionModel {
  const view = input.kind === 'saved' && input.workflow.state === 'available' ? input.workflow : undefined;
  const displayScope: RecordingNextActionModel['configurationScope'] = input.kind === 'draft' ? { kind: 'draft' }
    : view?.saved.configurationScope ?? { kind: 'group', groupId: input.scope.groupId };
  const offer = (kind: RecordingNextIntentKind, reason: RecordingNextReason, enabled = true, network = false,
    extra: Pick<RecordingNextIntent, 'operationId' | 'operationOwner' | 'activatedDeviceIds'> = {}): RecordingNextActionModel => {
    const offline = network && input.connection === 'offline';
    return copyFrozen({ intent: { kind, scope: input.scope, ...extra }, enabled: enabled && !offline,
      reasons: offline ? [reason, 'offline'] : [reason], configurationScope: displayScope });
  };
  const prepare = (reason: RecordingNextReason): RecordingNextActionModel => {
    if (input.kind === 'draft') return offer('save_prepare', reason, true, true);
    const edit = view?.lifecycle.state === 'available' ? view.lifecycle.value.edit : undefined;
    return offer('save_prepare', !edit ? 'lifecycle_unavailable' : edit.state === 'blocked' ? 'lifecycle_blocked' : reason,
      edit?.state === 'available', true);
  };
  if (!input.scope.workspaceId || !input.scope.deviceId || input.scope.circuitId === '' ||
    (input.kind === 'saved' && (!input.scope.groupId || (view && scopeMismatch(view.scope, input.scope)))) ||
    (input.draft && input.draft.workspace_id !== input.scope.workspaceId)) return offer('inspect', 'scope_unavailable');
  const retained = input.unconfirmedOperation;
  if (retained) {
    if (retained.scope.workspaceId !== input.scope.workspaceId || retained.scope.deviceId !== input.scope.deviceId ||
      retained.scope.circuitId !== input.scope.circuitId || retained.scope.groupId !== (input.kind === 'saved' ? input.scope.groupId : undefined)) {
      return offer('inspect', 'scope_unavailable');
    }
    return offer(retained.operationId ? 'query_operation' : 'recover_operation', 'operation_recovery', true, Boolean(retained.operationId),
      retained.operationId ? { operationId: retained.operationId } : {});
  }
  const owners = inspectRecordingOwners(input.ownerObservations, input.scope);
  const legacyOperation = view?.operation.state === 'available' && view.operation.value.kind !== 'none' && operationUnresolved(view.operation)
    ? view.operation.value : undefined;
  if (legacyOperation) return offer(legacyOperation.operation.operation_id ? 'query_operation' : 'recover_operation', 'operation_recovery', true,
    Boolean(legacyOperation.operation.operation_id), { operationOwner: legacyOperation.kind,
      ...(legacyOperation.operation.operation_id ? { operationId: legacyOperation.operation.operation_id } : {}) });
  if (owners.state === 'unresolved') return offer(owners.operationId ? 'query_operation' : 'recover_operation', 'operation_recovery', true,
    Boolean(owners.operationId), { operationOwner: owners.owner, ...(owners.operationId ? { operationId: owners.operationId } : {}) });
  if (owners.state === 'unavailable') return offer('recover_operation', owners.reason === 'loading' ? 'evidence_loading' :
    owners.reason === 'scope_mismatch' || owners.reason === 'revision_mismatch' ? 'scope_unavailable' : 'evidence_unavailable', true, false,
    owners.owner ? { operationOwner: owners.owner } : {});
  if (view?.operation.state === 'unavailable' && view.operation.reason !== 'missing') return offer('recover_operation',
    view.operation.reason === 'loading' ? 'evidence_loading' : 'evidence_unavailable');
  if (input.kind === 'saved' && !view) return offer('inspect', 'scope_unavailable');
  const draft = input.draft ?? view?.saved.group;
  if (!draft) return offer('inspect', 'scope_unavailable');
  if (draft.members.length === 0) return offer('choose_sources', 'sources_empty');
  if ((input.kind === 'draft' && !draft.members.some((member) => member.device_id === input.scope.deviceId &&
    (input.scope.circuitId === undefined || member.circuit_id === input.scope.circuitId))) ||
    (view?.saved.configurationScope.kind === 'managed_device' && draft.members.some((member) => member.device_id !== input.scope.deviceId))) {
    return offer('inspect', 'scope_unavailable');
  }
  const unsaved = input.kind === 'draft' || input.draft !== undefined;
  const serverReadiness = !unsaved && view?.readiness.state === 'available' ? view.readiness.value : undefined;
  // finalizeWriteGroupReadiness appends counter refusals AFTER config/schema flags; neither flag proves counter semantics.
  const counterBlocked = serverReadiness?.issues.some((issue) => issue.severity === 'blocking' && [
    'counter-semantics-required', 'counter-measurement-not-counter', 'counter-measurement-foreign',
    'counter-measurement-unavailable', 'counter-measurement-unconfirmed', 'counter-type-unsupported',
  ].includes(issue.code));
  if (counterBlocked || (!serverReadiness?.ready && needsCounterConfirmation(draft, input.counterStatuses))) {
    return offer('resolve_semantics', 'counter_confirmation');
  }
  if (unsaved) return prepare('save_required');
  if (!view) return offer('inspect', 'scope_unavailable');
  if (hasMonitorHandoff(view)) return offer('monitor', 'monitor_available', true, false, { activatedDeviceIds: [view.scope.deviceId] });
  if (view.readiness.state !== 'available') {
    const reason = view.readiness.reason;
    return offer('refresh_evidence', reason === 'loading' ? 'evidence_loading' : reason === 'error' ? 'evidence_error' : 'evidence_unavailable',
      reason !== 'loading', true);
  }
  const readiness = view.readiness.value;
  if (!readiness.config_ready) return prepare('configuration_blocked');
  if (!readiness.schema_ready) {
    const conflict = readiness.issues.some((issue) => issue.severity === 'blocking' &&
      !['schema-confirmation-required', 'schema-unverified'].includes(issue.code));
    if (conflict) return offer('review_schema', 'schema_conflict');
    const review = input.schemaReview;
    if (review && (recordingOwnerScopeReason(review.scope, input.scope, true) ||
      !['idle', 'previewed', 'reviewed'].includes(review.stage) || (review.stage !== 'idle' && !review.operationId))) {
      return offer('recover_operation', 'scope_unavailable', true, false, { operationOwner: 'schema' });
    }
    if (review?.stage === 'reviewed') return offer('confirm_schema', 'schema_confirmation', true, true, { operationId: review.operationId, operationOwner: 'schema' });
    if (review?.stage === 'previewed') return offer('review_schema', 'schema_review');
    return offer('preview_schema', 'schema_review', true, true);
  }
  if (!readiness.ready || readiness.issues.some((issue) => issue.severity === 'blocking')) return offer('refresh_evidence', 'readiness_blocked', true, true);
  return offer('review_scoped_start', 'start_review');
}
