import axios, { type AxiosRequestConfig } from 'axios';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import { VITE_API_BASE_URL } from '../env';
import type { APIResponse, StudioV2Workspace } from '../types/datalink';

export const studioV2DatalinkApi = axios.create({
  baseURL: `${VITE_API_BASE_URL}/datalink`,
  headers: {
    'Content-Type': 'application/json',
  },
});

studioV2DatalinkApi.interceptors.request.use((config) => {
  const uiVersion = resolveGatewayUiVersion();
  if (uiVersion) {
    config.headers = config.headers ?? {};
    if (!config.headers['X-UI-Version']) {
      config.headers['X-UI-Version'] = uiVersion;
    }
  }
  return config;
});

studioV2DatalinkApi.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error),
);

export const studioV2WorkspaceAPI = {
  async get(options?: AxiosRequestConfig): Promise<StudioV2Workspace> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2Workspace>>('/studio-v2/workspace', options);
    return res.data.data!;
  },
};
