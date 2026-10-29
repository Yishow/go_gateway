import { studioV2DatalinkApi } from './studioV2Workspace';
import { parseRecordingTestWriteRecovery } from '../utils/recordingTestWriteRecoveryJson';
/** Workspace authority is the server session; reject rather than display another workspace. */
export async function listRecordingTestWriteRecovery(workspaceId: string) {
  const response = await studioV2DatalinkApi.get<unknown>('/studio-v2/workspace/database-operations/test-writes');
  return parseRecordingTestWriteRecovery(response.data, workspaceId);
}
