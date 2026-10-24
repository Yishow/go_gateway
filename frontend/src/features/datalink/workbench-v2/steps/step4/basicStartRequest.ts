import type { RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupDraft } from '../../../../../types/studioV2WriteGroup';
import type { ModbusShareStatus } from '../../../../../types/modbusShare';
import { optionalShareRequest } from './basicRecordingPanelHelpers';

export interface BasicStartScope {
  workspaceId?: string;
  deviceId: string;
  connectorRevision?: string;
  shareStatus?: ModbusShareStatus | null;
}

/** The recording-start request of the legacy per-device wide group, built from exactly the revisions the operator sees. */
export function buildBasicStartRequest(scope: BasicStartScope, group: WriteGroup | undefined, expectedWorkspaceRevision: string, requestId: string, draft?: WriteGroupDraft): RecordingStartRequest {
  return {
    request_id: requestId,
    workspace_id: scope.workspaceId ?? group?.workspace_id ?? '',
    expected_workspace_revision: expectedWorkspaceRevision,
    device_ids: [scope.deviceId],
    groups: group ? [{
      group_id: group.id,
      expected_group_revision: group.revision,
      expected_connector_revision: scope.connectorRevision ?? group.destination.connector_revision,
      ...(draft ? { draft } : {}),
    }] : [],
    ...optionalShareRequest(scope.shareStatus),
  };
}
