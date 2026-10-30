import { inspectRecordingOwners, type RecordingOwnerInspection, type RecordingOwnerObservations } from '../state/recordingOwnerObservation';
import type { RecordingHomeGroupEvidence, RecordingHomeGroupProjection, RecordingHomeScope } from './recordingHomeProjection';
import type { RecordingHomeAttentionItem, RecordingHomeGuidanceAction } from './recordingHomeGuidance';

export const compareIdentity = (a: string, b: string): number => a < b ? -1 : a > b ? 1 : 0;

/** Inspect each actual member context through the existing owner authority. */
export function homeOwners(scope: RecordingHomeScope, item: RecordingHomeGroupProjection,
  evidence: RecordingHomeGroupEvidence<RecordingOwnerObservations> | undefined,
): { inspection: RecordingOwnerInspection; deviceId?: string; circuitId?: string } {
  const unavailable = { inspection: { state: 'unavailable', reason: 'unreadable' } } as const;
  const group = item.group;
  if (!evidence || evidence.state !== 'available') return unavailable;
  const bound = evidence.scope;
  if (bound.workspaceId !== scope.workspaceId || bound.workspaceRevision !== scope.workspaceRevision || bound.epoch !== scope.epoch ||
    bound.groupId !== group.id || bound.groupRevision !== group.revision || bound.appliedRevision !== group.applied_revision) return unavailable;
  const contexts = [...new Map(group.members.map(member => [JSON.stringify([member.device_id, member.circuit_id ?? null]),
    { deviceId: member.device_id, circuitId: member.circuit_id }])).values()];
  if (!contexts.length && group.basic_managed_device_id) contexts.push({ deviceId: group.basic_managed_device_id, circuitId: undefined });
  contexts.sort((a, b) => compareIdentity(a.deviceId, b.deviceId) || compareIdentity(a.circuitId ?? '', b.circuitId ?? ''));
  let result: ReturnType<typeof homeOwners> = contexts.length ? { inspection: { state: 'healthy' } } : unavailable;
  for (const context of contexts) {
    const observations = evidence.value.filter(observation => observation.scope.deviceId === context.deviceId && observation.scope.circuitId === context.circuitId)
      .sort((a, b) => compareIdentity(a.owner, b.owner));
    const inspection = inspectRecordingOwners(observations, { ...scope, ...context, groupId: group.id,
      groupRevision: group.revision, appliedRevision: group.applied_revision });
    if (inspection.state === 'unresolved') return { inspection, ...context };
    if (inspection.state === 'unavailable') result = { inspection, ...context };
  }
  return result;
}

export function groupAction(scope: RecordingHomeScope, item: RecordingHomeGroupProjection,
  section: Extract<RecordingHomeGuidanceAction, { target: 'recording' }>['section'] = 'overview',
): RecordingHomeGuidanceAction {
  return { kind: 'navigate', scope, target: 'recording', groupId: item.group.id,
    groupRevision: item.group.revision, section, intent: null };
}

export function homeAttention(scope: RecordingHomeScope, item: RecordingHomeGroupProjection,
  owner: ReturnType<typeof homeOwners>,
): RecordingHomeAttentionItem | undefined {
  const base = { id: item.group.id, name: item.group.name, evidence: item };
  const { inspection } = owner;
  if (inspection.state === 'unresolved' && owner.deviceId) return { ...base, reason: 'operation_recovery',
    action: { kind: 'navigate', scope, target: 'recording', groupId: item.group.id, groupRevision: item.group.revision, section: 'operation',
      intent: { kind: inspection.operationId ? 'query_operation' : 'recover_operation',
        scope: { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision, deviceId: owner.deviceId,
          ...(owner.circuitId ? { circuitId: owner.circuitId } : {}), groupId: item.group.id,
          groupRevision: item.group.revision, appliedRevision: item.group.applied_revision },
        operationOwner: inspection.owner, ...(inspection.operationId ? { operationId: inspection.operationId } : {}) } } };
  const delivery = item.delivery.state === 'available' ? item.delivery.value : undefined;
  if (delivery && (delivery.intake.state === 'blocked' || delivery.stages.blocked > 0 || delivery.stages.quarantined > 0 ||
    delivery.stages.unknown > 0 || delivery.stages.retrying > 0 || delivery.storage?.repair_state === 'blocked' || delivery.quota.intake_refused)) {
    return { ...base, reason: 'delivery_attention', action: groupAction(scope, item) };
  }
  if (delivery?.source_evidence?.state === 'available' && delivery.source_evidence.sources.some(source =>
    source.latest !== null && (!source.latest.available || ['bad', 'stale', 'invalid', 'missing', 'reset_or_uncertain', 'gap', 'partial'].includes(source.latest.status)))) {
    return { ...base, reason: 'source_attention', action: groupAction(scope, item, 'sources') };
  }
  // Missing operation/delivery evidence does not turn a saved running label into health.
  // A draft can still be continued; this is not permission to execute its commands.
  if (item.group.status !== 'draft' && (inspection.state === 'unavailable' || !delivery || item.lifecycle.state !== 'available')) {
    return { ...base, reason: 'evidence_unavailable', action: groupAction(scope, item) };
  }
}
