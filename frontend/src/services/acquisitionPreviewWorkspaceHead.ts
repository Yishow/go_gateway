import { studioV2DatalinkApi } from './studioV2Workspace';
import { parseWriteGroupAPIData, parseWriteGroupListData } from '../utils/studioV2WriteGroupJson';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';

export type AcquisitionWorkspaceHeadResult = { ok: true; revision: string } | { ok: false; reason: 'foreign-scope' | 'unavailable' };
/** A detached canonical read receipt, never an editable setup CAS acknowledgement. */
export async function readAcquisitionPreviewWorkspaceHead(workspace: string, signal: AbortSignal): Promise<AcquisitionWorkspaceHeadResult> {
  try {
    const response = await studioV2DatalinkApi.get<unknown>('/studio-v2/workspace/write-groups', { signal });
    const value = parseWriteGroupAPIData(response.data, 'write-group list', parseWriteGroupListData, 8);
    if (signal.aborted) return { ok: false, reason: 'unavailable' };
    if (value.workspace_id !== workspace) return { ok: false, reason: 'foreign-scope' };
    if (value.workspace_revision !== '' && safeRuntimeIdentity(value.workspace_revision) === '') return { ok: false, reason: 'unavailable' };
    return { ok: true, revision: value.workspace_revision };
  } catch {
    return { ok: false, reason: 'unavailable' };
  }
}
