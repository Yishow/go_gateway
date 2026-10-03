import type { WriteGroupDraft } from './studioV2WriteGroup';
import type { SchemaOperationStatus } from './recordingPlan';

export interface RecordingStartGroupIntent {
  group_id: string;
  expected_group_revision: string;
  expected_connector_revision: string;
  draft?: WriteGroupDraft;
}

export interface RecordingStartRequest {
  request_id: string;
  workspace_id: string;
  expected_workspace_revision: string;
  device_ids: string[];
  groups: RecordingStartGroupIntent[];
  readiness_token?: string;
  settings_revision?: string;
  workspace_revision?: string;
}

export interface RecordingStartGroupProgress {
  group_id: string;
  group_revision: string;
  applied_revision?: string;
  saved: boolean;
  ready: boolean;
  applied: boolean;
}

export interface RecordingStartDeviceProgress {
  device_id: string;
  activated: boolean;
  reason?: string;
}

export interface RecordingStartOperation {
  operation_id: string;
  action: 'recording_start';
  status: SchemaOperationStatus;
  intent_digest: string;
  reason?: string;
  next_action?: string;
  workspace_id: string;
  device_ids: string[];
  setup_revision: string;
  groups: RecordingStartGroupProgress[];
  devices: RecordingStartDeviceProgress[];
  stage: 'save' | 'readiness' | 'apply' | 'activation' | 'complete';
  created_at: string;
  updated_at: string;
}
