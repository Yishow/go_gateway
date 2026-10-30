import type { RecordingGroupView } from '../../../../types/runtimeRuleEvidence';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { RecordingHomeGroupProjection } from './recordingHomeProjection';
import type { RecordingOwnerInspection } from '../state/recordingOwnerObservation';
export interface RecordingHomeCommitHistory {
  readonly workspaceId: string; readonly groupId: string; readonly committedAt: string;
  /** Captured only when a matching receipt's destination was independently known. */
  readonly destination?: string;
}
export interface HomeResultContext {
  readonly runtime?: RecordingReadonly<RecordingGroupView>;
  readonly owner?: RecordingOwnerInspection;
  /** A previously observed SQL fact, never reused as current evidence. */
  readonly history?: RecordingHomeCommitHistory;
}
const instant = (value: string | null | undefined): value is string => typeof value === 'string' && Number.isFinite(Date.parse(value));
/** Describes independent facts, not an all-green health verdict or mutation permission.
 * Applied revision is requested configuration, not proof of currently effective runtime. */
export function recordingHomeResultSummary(item: RecordingHomeGroupProjection, context: HomeResultContext = {}) {
  const group = item.group, delivery = item.delivery.state === 'available' ? item.delivery.value : undefined;
  const history = context.history?.workspaceId === group.workspace_id && context.history.groupId === group.id && instant(context.history.committedAt) ? context.history : undefined;
  const effect = delivery?.last_sql_committed_effect;
  const matched = Boolean(effect && group.applied_revision && effect.group_revision === group.applied_revision &&
    effect.connector_revision === group.destination.connector_revision && effect.record_id && effect.effect_key &&
    /^[a-f0-9]{64}$/.test(effect.payload_digest) && instant(effect.committed_at));
  const currentRevision = group.revision === group.applied_revision;
  const runtime = context.runtime;
  const effective = runtime?.effective;
  const runtimeReceipt = runtime?.delivery?.current?.last_commit;
  const effectiveReceipt = Boolean(runtimeReceipt && effective && runtimeReceipt.group_revision === runtime.effective_revision &&
    runtimeReceipt.connector_revision === effective.destination.connector_revision && instant(runtimeReceipt.committed_at) && instant(runtime?.as_of) && Date.parse(runtimeReceipt.committed_at) <= Date.parse(runtime.as_of));
  const runtimeCurrent = runtime?.delivery?.current;
  const accepting = runtimeCurrent?.intake?.state === 'active' && runtimeCurrent.intake.group_revision === runtime?.effective_revision &&
    runtimeCurrent.binding.entity_revision === runtime.effective_revision;
  const lastCommittedAt = effectiveReceipt && runtimeReceipt ? runtimeReceipt.committed_at : matched && effect ? effect.committed_at : instant(delivery?.last_sql_committed_at) ? delivery.last_sql_committed_at : history?.committedAt ?? null;
  const sql = effectiveReceipt ? 'effective_revision_receipt' : matched ? currentRevision ? 'revision_receipt' : 'previous_revision_receipt' : lastCommittedAt ? 'historical_commit' :
    accepting && runtimeCurrent?.stages.sql_committed === 0 ? 'awaiting_first_write' : 'unknown';
  const currentBound = runtimeCurrent && effective && runtimeCurrent.binding.entity_revision === runtime?.effective_revision && instant(runtime.as_of);
  const observedWindow = (end: string | null | undefined, previous?: string | null) => instant(end) && instant(runtime?.as_of) &&
    Date.parse(end) <= Date.parse(runtime.as_of) && (!instant(previous) || Date.parse(end) >= Date.parse(previous));
  const sourceFault = currentBound && (
    runtimeCurrent.sources.some(source => source.latest && ['partial', 'missing'].includes(source.latest.state) && observedWindow(source.latest.window_end, source.last_committed?.window_end)) ||
    delivery?.source_evidence?.state === 'available' && delivery.source_evidence.group_revision === runtime.effective_revision &&
    delivery.source_evidence.sources.some(source => source.latest?.available && source.latest.group_revision === runtime.effective_revision &&
      source.latest.connector_revision === effective.destination.connector_revision && observedWindow(source.latest.window_end, source.last_sql_committed?.available ? source.last_sql_committed.window_end : undefined) &&
      ['bad', 'stale', 'invalid', 'missing', 'reset_or_uncertain', 'gap', 'partial'].includes(source.latest.status) &&
      effective.members.some(member => member.device_id === source.device_id && member.point_id === source.point_id && member.tag_id === source.tag_id)));
  const stages = currentBound ? runtimeCurrent.stages : undefined;
  const blocked = currentBound && (runtimeCurrent.intake?.state === 'blocked' || stages!.blocked > 0 || stages!.quarantined > 0 ||
    stages!.unknown > 0 || stages!.retrying > 0 || runtimeCurrent.quota?.intake_refused || runtimeCurrent.storage?.repair_state === 'blocked');
  return { saved: true as const, configuration: currentRevision && !!group.applied_revision ? 'requested' as const : 'not_requested' as const,
    effective: runtime?.effective_revision ? runtime.effective_revision === group.applied_revision ? 'requested_effective' as const : 'requested_not_effective' as const : 'unknown' as const,
    effectiveRevision: runtime?.effective_revision ?? null, observedAt: runtime?.as_of ?? null,
    requestedBoundary: runtime?.applied_effective_at ?? null,
    sql, lastCommittedAt,
    // These are historical SQL coordinates, never raw row values or readback.
    historicalDestination: effectiveReceipt && effective ? [effective.destination.schema, effective.destination.table].filter(Boolean).join('.') : matched && currentRevision ? [group.destination.table_schema, group.destination.table_name].filter(Boolean).join('.') : history?.destination ?? null,
    collection: accepting ? 'accepting' as const : 'unknown' as const,
    current: blocked ? 'delivery_attention' as const : sourceFault ? 'source_attention' as const : 'unknown' as const,
    owners: context.owner?.state === 'unresolved' ? 'unresolved' as const : context.owner?.state === 'healthy' ? 'inspected' as const : 'unknown' as const,
    readback: instant(runtime?.delivery?.test_write?.readback_at) ? 'test_write_verified' as const : 'unknown' as const,
    cleanup: instant(runtime?.delivery?.test_write?.cleanup_at) ? 'test_write_completed' as const : 'unknown' as const,
  };
}
