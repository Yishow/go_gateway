import axios from 'axios';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import { VITE_API_BASE_URL } from '../env';

const API_BASE = VITE_API_BASE_URL;
export const DATALINK_BASE = `${API_BASE}/datalink`;

export const api = axios.create({
  baseURL: DATALINK_BASE,
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
  (error) => Promise.reject(error),
);

/** A reply without a list must not be read as "nothing exists". */
export function requireList<T>(value: T[] | null | undefined, what: string): T[] {
  if (!Array.isArray(value)) {
    throw new Error(`${what} returned no list`);
  }
  return value;
}
