import type { APIResponse } from '../types/datalink';
import type { HistoryQuery, HistoryReport } from '../types/historyReport';
import { studioV2DatalinkApi } from './studioV2Workspace';

export const studioV2WorkspaceHistoryAPI = {
  async queryHistory(query: HistoryQuery): Promise<HistoryReport> {
    const res = await studioV2DatalinkApi.post<APIResponse<HistoryReport>>(
      '/studio-v2/workspace/history/query',
      query
    );
    return res.data.data!;
  },

  async exportCSV(query: HistoryQuery): Promise<Blob> {
    const res = await studioV2DatalinkApi.post(
      '/studio-v2/workspace/history/export',
      query,
      {
        responseType: 'blob',
      }
    );
    return res.data as Blob;
  },
};
