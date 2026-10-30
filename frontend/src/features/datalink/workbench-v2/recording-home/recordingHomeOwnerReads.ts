import { inspectBasicRecordingDevice, type BasicRecordingDeviceInspection } from '../state/basicRecording';
import { inspectSchemaOperation } from '../steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
import type { RecordingOwnerObservations, RecordingOwnerObservationState } from '../state/recordingOwnerObservation';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { RecordingTestWriteRecoveryList } from '../../../../types/recordingTestWriteRecovery';
import type { RecordingHomeEvidence, RecordingHomeGroupEvidence, RecordingHomeGroupScope } from './recordingHomeProjection';
/** Existing owners are inspected, never mounted, mutated, cleared or replaced.
 * Empty local storage is not global operation absence. Start roles are discovered only in their original store, never derived
 * from a group/member DTO; empty inspection remains unknown. */
export function readRecordingHomeOwners(scope: RecordingHomeGroupScope, group: RecordingReadonly<WriteGroup>,
  server?: RecordingHomeEvidence<RecordingTestWriteRecoveryList>,
  starts?: ReadonlyMap<string, BasicRecordingDeviceInspection>,
): RecordingHomeGroupEvidence<RecordingOwnerObservations> {
  if (group.workspace_id !== scope.workspaceId || group.id !== scope.groupId || group.revision !== scope.groupRevision || group.applied_revision !== scope.appliedRevision) return { state: 'unknown', scope };
  const schema = inspectSchemaOperation(scope.workspaceId, group.id), local = recordingTestWriteRecovery.read();
  const serverKnown = server?.state === 'available' && server.scope.workspaceId === scope.workspaceId &&
    server.scope.workspaceRevision === scope.workspaceRevision && server.scope.epoch === scope.epoch &&
    server.value.workspace_id === scope.workspaceId && server.value.complete;
  // Retained work from an old configuration still needs its native recovery owner.
  const reference = local.references.find(item => item.workspace_id === scope.workspaceId && item.group_id === group.id);
  const remote = serverKnown ? server.value.operations.find(item => item.group_id === group.id) : undefined;
  const schemaObservation: RecordingOwnerObservationState = schema.state === 'retained'
    ? { state: 'unresolved', reason: 'unknown', operationId: schema.operationId }
    : { state: 'unavailable', reason: 'unreadable' };
  const testWrite: RecordingOwnerObservationState = reference || remote
    ? { state: 'unresolved', reason: 'cleanup_pending', operationId: (reference ?? remote)!.operation_id }
    : local.unavailable || !serverKnown ? { state: 'unavailable', reason: 'unreadable' } : { state: 'observed', result: 'none' };
  const contexts = [...new Map(group.members.map(member => [JSON.stringify([member.device_id, member.circuit_id]),
    { deviceId: member.device_id, circuitId: member.circuit_id }])).values()];
  if (!contexts.length && group.basic_managed_device_id) contexts.push({ deviceId: group.basic_managed_device_id, circuitId: undefined });
  return { state: 'available', scope, value: contexts.flatMap(context => {
    const bound = { ...scope, ...context };
    const start = (starts?.get(context.deviceId) ?? inspectBasicRecordingDevice(scope.workspaceId, context.deviceId)).retained.find(item => item.groupIds.includes(group.id));
    const startObservation: RecordingOwnerObservationState = start ? { state: 'unresolved', reason: start.operationId ? 'unknown' : 'lost_ack',
      ...(start.operationId ? { operationId: start.operationId } : {}) } : { state: 'unavailable', reason: 'unreadable' };
    return [{ ...schemaObservation, owner: 'schema' as const, scope: bound },
      { ...startObservation, owner: 'start' as const, scope: bound },
      { ...testWrite, owner: 'test_write' as const, scope: bound }];
  }) };
}
