import axios from 'axios';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import { VITE_API_BASE_URL } from '../env';
import type { APIResponse } from '../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../types/sourceRuleCandidates';

const api = axios.create({
  baseURL: `${VITE_API_BASE_URL}/datalink`,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface SourceRuleCandidateScope {
  workspace_id: string;
  expected_workspace_revision: string;
  revision_id: string;
}

api.interceptors.request.use((config) => {
  const uiVersion = resolveGatewayUiVersion();
  if (uiVersion) {
    config.headers = config.headers ?? {};
    if (!config.headers['X-UI-Version']) {
      config.headers['X-UI-Version'] = uiVersion;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error),
);

export const sourceRuleCandidateAPI = {
  async get(ruleId: string, scope?: SourceRuleCandidateScope): Promise<SourceRuleCandidateSnapshotView> {
    const res = await api.get<APIResponse<SourceRuleCandidateSnapshotView>>(
      `/source-rules/${ruleId}/candidates`,
      scope ? { params: scope } : undefined,
    );
    return res.data.data!;
  },
};
