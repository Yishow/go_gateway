import { compareIdentity, groupAction, homeAttention, homeOwners, homeUnconfirmed } from './recordingHomeGuidanceEvidence';
import { recordingHomeProjection } from './recordingHomeProjection';
import type { RecordingNextIntent } from '../state/recordingWorkflowView';
import type { RecordingOwner, RecordingOwnerObservations } from '../state/recordingOwnerObservation';
import type { RecordingHomeGroupEvidence, RecordingHomeGroupProjection, RecordingHomeInput, RecordingHomeScope } from './recordingHomeProjection';
import type { RecordingUnavailable } from '../state/recordingWorkflowView';

/** Navigation only. Editor/operation owners must revalidate every command separately. */
export type RecordingHomeGuidanceAction = { readonly kind: 'navigate'; readonly scope: RecordingHomeScope } & (
  | { readonly target: 'connect_device' | 'evidence' }
  | { readonly target: 'configure_points' | 'choose_sources'; readonly deviceId: string }
  | { readonly target: 'recording'; readonly groupId: string; readonly groupRevision: string;
      readonly section: 'overview' | 'sources' | 'destination' | 'operation'; readonly intent: RecordingNextIntent | null }
);
export type RecordingHomeGuidanceReason = 'operation_recovery' | 'delivery_attention' | 'source_attention'
  | 'evidence_unavailable' | 'draft' | 'missing_devices' | 'missing_points' | 'missing_recordings' | 'recordings';
export interface RecordingHomeGuidanceItem {
  readonly id: string;
  readonly name: string;
  readonly evidence: RecordingHomeGroupProjection;
  readonly action: RecordingHomeGuidanceAction;
}
/** Which reads did not confirm this item. Presentation only: never health, absence or permission. */
export interface RecordingHomeUnconfirmed {
  readonly lifecycle: boolean;
  readonly delivery: boolean;
  /** Owners without a current observed result; unknown is never inferred from another owner's empty list. */
  readonly owners: readonly RecordingOwner[];
}
export interface RecordingHomeAttentionItem extends RecordingHomeGuidanceItem {
  readonly reason: 'operation_recovery' | 'delivery_attention' | 'source_attention' | 'evidence_unavailable';
  readonly unconfirmed?: RecordingHomeUnconfirmed;
}
const absent = (read: { readonly state: string; readonly reason?: string }) => read.state === 'unavailable' && read.reason === 'missing';
/** Neither detail read was ever provided for this item (e.g. off the requested page): not an observed gap. */
export const recordingHomeNotRequested = (item: RecordingHomeGuidanceItem): boolean =>
  absent(item.evidence.lifecycle) && absent(item.evidence.delivery);
/** Unknown that must not lead the page: a removed item's history, or an item whose details were never
 * requested. Both stay non-healthy; known unresolved work and actual read failures still lead. */
export const recordingHomeBackgroundAttention = (item: RecordingHomeAttentionItem): boolean =>
  item.reason === 'evidence_unavailable' && (item.evidence.group.status === 'deleted' || recordingHomeNotRequested(item));
export interface RecordingHomeGuidanceInput {
  readonly evidence: RecordingHomeInput;
  /** Retained owner observations, not absence inferred from running or successful delivery. */
  readonly owners?: Readonly<Record<string, RecordingHomeGroupEvidence<RecordingOwnerObservations>>>;
}
export interface RecordingHomeGuidance {
  readonly availability: RecordingUnavailable | { readonly state: 'available' };
  readonly reason: RecordingHomeGuidanceReason;
  readonly primaryAction: RecordingHomeGuidanceAction | null;
  readonly attentionItems: readonly RecordingHomeAttentionItem[];
  readonly draftItems: readonly RecordingHomeGuidanceItem[];
  readonly recordingItems: readonly RecordingHomeGuidanceItem[];
}
function resolveGuidance(input: RecordingHomeGuidanceInput): RecordingHomeGuidance {
  const evidence = recordingHomeProjection(input.evidence);
  const { scope } = evidence;
  const recordingItems: RecordingHomeGuidanceItem[] = evidence.groups.state === 'available' ? evidence.groups.value.map(item => ({
    id: item.group.id, name: item.group.name, evidence: item, action: groupAction(scope, item),
  })).sort((a, b) => compareIdentity(a.id, b.id)) : [];
  const draftItems = recordingItems.filter(item => item.evidence.group.status === 'draft').map(item => ({ ...item,
    action: groupAction(scope, item.evidence, !item.evidence.group.members.length ? 'sources' :
      !item.evidence.group.destination.connector_id || !item.evidence.group.destination.table_name ? 'destination' : 'overview'),
  }));
  const attentionItems = recordingItems.flatMap(item => {
    const attention = homeAttention(scope, item.evidence, homeOwners(scope, item.evidence, input.owners?.[item.id]));
    return attention ? [attention.reason === 'evidence_unavailable' ? { ...attention, unconfirmed: homeUnconfirmed(scope, item.evidence, input.owners?.[item.id]) } : attention] : [];
  });
  const priority = { operation_recovery: 0, delivery_attention: 1, source_attention: 2, evidence_unavailable: 3 };
  attentionItems.sort((a, b) => priority[a.reason] - priority[b.reason] || compareIdentity(a.id, b.id));
  const base = { availability: { state: 'available' as const }, attentionItems, draftItems, recordingItems };
  const unavailable = [evidence.workspace, evidence.sources, evidence.groups].find(item => item.state === 'unavailable');
  if (unavailable?.state === 'unavailable') return { ...base, availability: unavailable, reason: 'evidence_unavailable',
    primaryAction: scope.workspaceId && typeof scope.workspaceRevision === 'string' && Number.isSafeInteger(scope.epoch) && scope.epoch >= 0
      ? { kind: 'navigate', scope, target: 'evidence' } : null };
  const workspaceDevices = evidence.workspace.state === 'available' ? evidence.workspace.value.ordered_device_ids : [];
  const sourceCandidates = evidence.sources.state === 'available' ? evidence.sources.value.candidates : [];
  if (workspaceDevices.some(id => !id) || sourceCandidates.some(candidate => !workspaceDevices.includes(candidate.device_id))) {
    return { ...base, availability: { state: 'unavailable', reason: 'scope_mismatch' }, reason: 'evidence_unavailable',
      primaryAction: { kind: 'navigate', scope, target: 'evidence' } };
  }
  const lead = attentionItems.find(item => !recordingHomeBackgroundAttention(item));
  if (lead) return { ...base, reason: lead.reason, primaryAction: lead.action };
  if (draftItems.length) return { ...base, reason: 'draft', primaryAction: draftItems[0].action };
  const devices = [...workspaceDevices].sort(compareIdentity);
  if (!devices.length) return { ...base, reason: 'missing_devices', primaryAction: { kind: 'navigate', scope, target: 'connect_device' } };
  const candidates = sourceCandidates;
  if (!candidates.length) return { ...base, reason: 'missing_points', primaryAction: { kind: 'navigate', scope, target: 'configure_points', deviceId: devices[0] } };
  if (!recordingItems.length) return { ...base, reason: 'missing_recordings', primaryAction: { kind: 'navigate', scope, target: 'choose_sources',
    deviceId: [...candidates].sort((a, b) => compareIdentity(a.device_id, b.device_id) || compareIdentity(a.key, b.key))[0].device_id } };
  return { ...base, reason: 'recordings', primaryAction: recordingItems[0].action };
}

/** Freeze only newly owned projection output, never the input query-cache objects. */
function freezeResult<T>(value: T): T {
  if (value !== null && typeof value === 'object') {
    Object.values(value).forEach(freezeResult);
    Object.freeze(value);
  }
  return value;
}

export function recordingHomeGuidance(input: RecordingHomeGuidanceInput): RecordingHomeGuidance {
  return freezeResult(resolveGuidance(input));
}
