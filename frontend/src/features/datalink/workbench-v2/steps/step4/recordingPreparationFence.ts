import { inspectBasicRecordingIntent } from '../../state/basicRecording';
import type { RecordingOwnerObservationState } from '../../state/recordingOwnerObservation';
import { recordingTestWriteRecovery } from '../../state/recordingTestWriteRecovery';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { inspectSchemaOperation } from './writeGroup/schemaOperationRecovery';
import type { SchemaPreparationPresentation } from './writeGroup/schemaPreparationPresentation';

export interface RecordingPreparationFenceInput {
  ownerKey: string; workspaceId: string; workspaceRevision: string;
  groups: readonly WriteGroup[]; start: RecordingOwnerObservationState;
  schema: Readonly<Record<string, SchemaPreparationPresentation>>;
}
/** Fresh restrictive inspection of the existing owners. It grants no mutation permission. */
export function recordingPreparationBlocked(input: RecordingPreparationFenceInput, groupId?: string, includeTestWrite = true): boolean {
  const selected = groupId ? input.groups.filter(group => group.id === groupId) : input.groups;
  if (groupId && selected.length !== 1) return true;
  return ownersBlocked(input, selected, includeTestWrite);
}
function ownersBlocked(input: RecordingPreparationFenceInput, selected: readonly WriteGroup[], includeTestWrite: boolean, previewRetryGroupId?: string): boolean {
  if (!input.ownerKey || !input.workspaceId || !input.workspaceRevision || input.start.state !== 'observed') return true;
  const start = inspectBasicRecordingIntent(input.ownerKey);
  const resolved = input.start.result === 'resolved' ? input.start.operationId : undefined;
  if (start.state === 'unavailable' || start.state === 'retained' && (!resolved || resolved !== start.operationId)) return true;
  if (selected.some(group => {
    const report = input.schema[group.id];
    return group.workspace_id !== input.workspaceId || inspectSchemaOperation(input.workspaceId, group.id).state !== 'empty' ||
      !report || report.scope.workspaceId !== input.workspaceId || report.scope.workspaceRevision !== input.workspaceRevision ||
      report.scope.groupId !== group.id || report.scope.groupRevision !== group.revision ||
      report.scope.connectorId !== group.destination.connector_id || report.scope.connectorRevision !== group.destination.connector_revision ||
      (report.observation.state !== 'observed' && !(group.id === previewRetryGroupId && report.previewFailure === true &&
        report.observation.state === 'unavailable' && report.observation.reason === 'unreadable'));
  })) return true;
  if (!includeTestWrite) return false;
  const recovery = recordingTestWriteRecovery.read();
  return recovery.unavailable || recovery.references.some(ref => ref.workspace_id === input.workspaceId && selected.some(group => group.id === ref.group_id));
}

/** Retry-specific denial; the schema owner still validates its private failure receipt. */
export function recordingSchemaPreviewRetryBlocked(input: RecordingPreparationFenceInput, groupId: string): boolean {
  if (!groupId || input.groups.filter(group => group.id === groupId).length !== 1 || !input.schema[groupId]?.previewFailure) return true;
  return ownersBlocked(input, input.groups, true, groupId);
}
