import type { WriteGroupTestWriteOperation } from './studioV2WriteGroupTestWrite';
export interface RecordingTestWriteRecoveryOperation extends WriteGroupTestWriteOperation {
  workspace_id: string;
  workspace_revision: string;
  group_id: string;
  group_revision: string;
  connector_id: string;
  connector_revision: string;
  target_fingerprint: string;
}
export interface RecordingTestWriteRecoveryList {
  workspace_id: string;
  complete: true;
  operations: RecordingTestWriteRecoveryOperation[];
}
