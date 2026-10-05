import type { WorkbenchV2State } from '../../state/types';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';
import { recordingStartRevisionsMatch, recordingStartOperationOwnsCurrentSetupRevision, type BasicRecordingIntentRecord } from '../../state/basicRecording';
export function destinationIsSaved(state: WorkbenchV2State): boolean {
  const connector = state.db.connector;
  return connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);
}

export function optionalShareRequest(shareStatus: ModbusShareStatus | null | undefined): Pick<RecordingStartRequest, 'readiness_token' | 'settings_revision' | 'workspace_revision'> {
  if (!shareStatus?.readiness_token || !shareStatus.settings_revision || !shareStatus.workspace_revision) {
    return {};
  }
  return {
    readiness_token: shareStatus.readiness_token,
    settings_revision: shareStatus.settings_revision,
    workspace_revision: shareStatus.workspace_revision,
  };
}

export function groupMatchesRequest(record: BasicRecordingIntentRecord | undefined, group: WriteGroup | undefined, connectorRevision?: string, shareStatus?: ModbusShareStatus | null): boolean {
  const intent = record?.request.groups[0];
  if (!record || (record.operation_id && !recordingStartRevisionsMatch(record.request, shareStatus?.settings_revision, shareStatus?.workspace_revision))) return false;
  if (!group) return record.request.groups.length === 0;
  if (!intent) return false;
  return intent.group_id === group.id && intent.expected_group_revision === group.revision &&
    (!connectorRevision || intent.expected_connector_revision === connectorRevision);
}

export function operationScopeMatchesCurrent(
  record: BasicRecordingIntentRecord | undefined,
  operation: RecordingStartOperation | undefined,
  group: WriteGroup | undefined,
  workspaceId: string | undefined,
  deviceId: string,
  connectorRevision?: string,
): boolean {
  if (!record || !operation || !group || !connectorRevision || operation.workspace_id !== workspaceId ||
    !operation.device_ids.includes(deviceId)) return false;
  const intent = record.request.groups[0];
  const progress = operation.groups.find((item) => item.group_id === group.id);
  return Boolean(intent && progress && intent.group_id === group.id &&
    intent.expected_connector_revision === connectorRevision &&
    group.destination.connector_revision === connectorRevision &&
    progress.group_revision === group.revision &&
    (!progress.applied_revision || progress.applied_revision === group.applied_revision));
}

export function requestMatchesCurrentControls(
  record: BasicRecordingIntentRecord | undefined,
  operation: RecordingStartOperation | undefined,
  group: WriteGroup | undefined,
  workspaceId: string | undefined,
  deviceId: string,
  connectorRevision: string | undefined,
  intervalSeconds: number,
  incompletePolicy: string,
  shareStatus: ModbusShareStatus | null | undefined,
  workspaceRevision: string,
): boolean {
  if (!record || record.request.workspace_id !== workspaceId ||
    record.request.device_ids.length !== 1 || record.request.device_ids[0] !== deviceId) return false;
  const intent = record.request.groups[0];
  if (!group) {
    const unresolvedRequest = !record.operation_id && record.request.groups.length === 0;
    return record.request.groups.length === 0 &&
      (unresolvedRequest || recordingStartRevisionsMatch(record.request, shareStatus?.settings_revision, shareStatus?.workspace_revision));
  }
  const operationOwnsRevision = recordingStartOperationOwnsCurrentSetupRevision(operation, workspaceRevision, group);
  if (!intent || record.request.groups.length !== 1 || !connectorRevision ||
    intent.group_id !== group.id ||
    (intent.expected_group_revision !== group.revision && operation !== undefined && !operationOwnsRevision) ||
    intent.expected_connector_revision !== connectorRevision) return false;
  if (record.request.expected_workspace_revision !== workspaceRevision && operation !== undefined && !operationOwnsRevision) return false;
  const draftPolicy = intent.draft?.row_policy;
  const requestedInterval = draftPolicy?.interval_seconds ?? group.row_policy.interval_seconds;
  const requestedIncomplete = draftPolicy?.incomplete_policy ?? group.row_policy.incomplete_policy ?? 'skip_row';
  if (requestedInterval !== intervalSeconds || requestedIncomplete !== incompletePolicy) return false;
  return !record.operation_id || recordingStartRevisionsMatch(record.request, shareStatus?.settings_revision, shareStatus?.workspace_revision);
}

export function groupSnapshotMatchesRequest(
  record: BasicRecordingIntentRecord | undefined,
  groupId: string | undefined,
  groupRevision: string | undefined,
  connectorRevision?: string,
): boolean {
  if (!record) return false;
  const intent = record.request.groups[0];
  if (!groupId) return record.request.groups.length === 0;
  if (!intent) return false;
  return intent.group_id === groupId && intent.expected_group_revision === groupRevision &&
    (!connectorRevision || intent.expected_connector_revision === connectorRevision);
}

export function operationBelongsToGroupId(record: BasicRecordingIntentRecord | undefined, groupId: string | undefined): boolean {
  if (!record) return false;
  const intent = record.request.groups[0];
  return groupId ? intent?.group_id === groupId : record.request.groups.length === 0;
}

export function unresolvedRequestMatchesCurrentScope(
  record: BasicRecordingIntentRecord | undefined,
  groupId: string | undefined,
  connectorRevision: string | undefined,
): boolean {
  if (!record || record.operation_id) return false;
  const intent = record.request.groups[0];
  if (!groupId) return record.request.groups.length === 0;
  return Boolean(intent && record.request.groups.length === 1 && intent.group_id === groupId &&
    connectorRevision && intent.expected_connector_revision === connectorRevision);
}
