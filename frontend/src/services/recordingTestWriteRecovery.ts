import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
import { captureStudioLocatorReadScope, studioLocatorReadParams, throwStudioLocatorReadError } from './studioEntityLocator';
import { studioV2DatalinkApi } from './studioV2Workspace';
import { parseRecordingTestWriteRecovery } from '../utils/recordingTestWriteRecoveryJson';
/** Workspace authority is the server session; reject rather than display another workspace. */
export async function listRecordingTestWriteRecovery(workspaceId: string, locatorScope?: StudioV2LocatorReadScope) {
  const selected = captureStudioLocatorReadScope(locatorScope);
  if (selected && selected.workspaceId !== workspaceId) throw new Error('Scoped settings read unavailable.');
  const url = '/studio-v2/workspace/database-operations/test-writes';
  const response = await (selected ? studioV2DatalinkApi.get<unknown>(url, { params: studioLocatorReadParams(selected) }) : studioV2DatalinkApi.get<unknown>(url))
    .catch(error => throwStudioLocatorReadError(error, selected));
  return parseRecordingTestWriteRecovery(response.data, workspaceId);
}
