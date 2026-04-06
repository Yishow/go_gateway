import axios from 'axios';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import { VITE_API_BASE_URL } from '../env';
import type { APIResponse } from '../types/datalink';
import type {
  SourceRuleTagReviewDecision,
  UpsertSourceRuleTagReviewDecisionRequest,
} from '../types/sourceRuleTagReviewDecisions';

const api = axios.create({
  baseURL: `${VITE_API_BASE_URL}/datalink`,
  headers: {
    'Content-Type': 'application/json',
  },
});

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
  (error) => {
    if (error.response?.data?.error?.message) {
      error.message = error.response.data.error.message;
    } else if (error.response?.data?.message) {
      error.message = error.response.data.message;
    } else if (!error.message) {
      error.message = 'Request failed with status code ' + (error.response?.status || 'unknown');
    }
    return Promise.reject(error);
  },
);

export const sourceRuleTagReviewDecisionAPI = {
  async list(ruleId: string): Promise<SourceRuleTagReviewDecision[]> {
    const res = await api.get<APIResponse<SourceRuleTagReviewDecision[]>>(
      `/source-rules/${ruleId}/tag-review-decisions`,
    );
    return res.data.data ?? [];
  },

  async upsert(
    ruleId: string,
    request: UpsertSourceRuleTagReviewDecisionRequest,
  ): Promise<SourceRuleTagReviewDecision> {
    const res = await api.post<APIResponse<SourceRuleTagReviewDecision>>(
      `/source-rules/${ruleId}/tag-review-decisions`,
      request,
    );
    return res.data.data!;
  },
};
