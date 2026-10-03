import type { RecordingStartOperation, RecordingStartRequest } from '../types/studioV2RecordingStart';
import type { WriteGroupCreateRequest, WriteGroupMutationResponse } from '../types/studioV2WriteGroup';
import { parseRecordingStartOperation } from '../utils/studioV2RecordingStartJson';
import { parseWriteGroupAPIData, parseWriteGroupMutationData } from '../utils/studioV2WriteGroupJson';
import { studioV2DatalinkApi } from './studioV2Workspace';

function sameIds(left: string[], right: string[]): boolean {
  const orderedRight = [...right].sort();
  return left.length === right.length && [...left].sort().every((id, index) => id === orderedRight[index]);
}

export const studioV2RecordingStartAPI = {
  async ensureBasic(deviceId: string, request: WriteGroupCreateRequest): Promise<WriteGroupMutationResponse> {
    const result = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/basic/${encodeURIComponent(deviceId)}`, request,
    );
    return parseWriteGroupAPIData(result.data, 'basic write group', parseWriteGroupMutationData);
  },

  async start(request: RecordingStartRequest): Promise<RecordingStartOperation> {
    const result = await studioV2DatalinkApi.post<unknown>('/studio-v2/workspace/recording-start', request);
    return parseWriteGroupAPIData(result.data, 'recording start', (data) => {
      const operation = parseRecordingStartOperation(data);
      return operation && operation.workspace_id === request.workspace_id &&
        sameIds(operation.device_ids, request.device_ids) &&
        sameIds(operation.groups.map((group) => group.group_id), request.groups.map((group) => group.group_id))
        ? operation : null;
    });
  },

  async operation(operationId: string): Promise<RecordingStartOperation> {
    const result = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/recording-start/operations/${encodeURIComponent(operationId)}`,
    );
    return parseWriteGroupAPIData(result.data, 'recording start operation', (data) => {
      const operation = parseRecordingStartOperation(data);
      return operation?.operation_id === operationId ? operation : null;
    });
  },
};
