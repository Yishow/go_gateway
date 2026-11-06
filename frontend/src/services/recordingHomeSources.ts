import type { AxiosInstance } from 'axios';
import { studioV2DatalinkApi } from './studioV2Workspace';
import type { RecordingHomeSourceObservation } from '../types/recordingHomeSources';
const object = (value: unknown): value is Record<string, unknown> => typeof value === 'object' && value !== null && !Array.isArray(value);
const strings = (value: Record<string, unknown>, keys: string[]) => keys.every(key => typeof value[key] === 'string');
// Go strings.TrimSpace uses Unicode White_Space (NEL included, BOM excluded).
// Check blankness only; canonical IDs must never be normalized.
const identity = (value: unknown): value is string => typeof value === 'string' && !/^\p{White_Space}*$/u.test(value);
function valid(value: unknown, workspaceId: string): value is RecordingHomeSourceObservation {
  if (!identity(workspaceId) || !object(value) || value.workspace_id !== workspaceId) return false;
  if (value.state === 'unknown') return ['source-inconsistent', 'source-changed'].includes(String(value.reason_code)) &&
    !['devices', 'candidates', 'excluded', 'source_digest'].some(key => key in value);
  if (value.state !== 'available' || typeof value.source_digest !== 'string' || !/^[a-f0-9]{64}$/.test(value.source_digest)) return false;
  if (!Array.isArray(value.devices) || !Array.isArray(value.candidates) || !Array.isArray(value.excluded)) return false;
  const devices = value.devices;
  return devices.every(item => object(item) && strings(item, ['id', 'name', 'protocol', 'availability']) && identity(item.id)) &&
    new Set(devices.map(item => item.id)).size === devices.length &&
    value.candidates.every(item => object(item) && strings(item, ['key', 'device_id', 'point_id', 'tag_id', 'tag_key', 'label', 'device_name', 'address', 'target_type']) &&
      ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'].includes(String(item.target_type)) &&
      identity(item.point_id) && Boolean(item.tag_id) && item.key === `${item.point_id}|${item.tag_id}` && devices.some(device => device.id === item.device_id)) &&
    new Set(value.candidates.map(item => item.key)).size === value.candidates.length &&
    new Set(value.candidates.map(item => item.point_id)).size === value.candidates.length &&
    value.excluded.every(item => object(item) && typeof item.label === 'string' && ['device-not-saved', 'mapping-not-saved', 'tag-missing'].includes(String(item.reason)));
}
/** One complete readonly observation; never a readiness or command token. */
export const recordingHomeSourcesAPI = {
  async read(workspaceId: string, signal?: AbortSignal, backend: AxiosInstance = studioV2DatalinkApi, apiOrigin?: string): Promise<RecordingHomeSourceObservation> {
    const response = await backend.get<{ success: boolean; data: unknown }>('/studio-v2/recording-home/sources', { params: { workspace_id: workspaceId }, signal, ...(apiOrigin ? { baseURL: apiOrigin } : {}) });
    if (response.data.success !== true || !valid(response.data.data, workspaceId)) throw new Error('Source observation unavailable');
    return response.data.data;
  },
};
