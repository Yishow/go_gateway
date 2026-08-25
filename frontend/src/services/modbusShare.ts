import type { APIResponse } from '../types/datalink';
import type {
  ModbusShareMapping,
  ModbusShareReconcileOutcome,
  ModbusShareReconcileRequest,
  ModbusShareStatus,
  ModbusShareLifecycleRequest,
} from '../types/modbusShare';
import { api } from './datalinkClient';

export const modbusShareAPI = {
  async status(): Promise<ModbusShareStatus> {
    const res = await api.get<APIResponse<ModbusShareStatus>>('/modbus-share/status');
    return res.data.data!;
  },

  async reconcile(request: ModbusShareReconcileRequest): Promise<ModbusShareReconcileOutcome> {
    const res = await api.post<APIResponse<ModbusShareReconcileOutcome>>(
      '/modbus-share/reconcile',
      request,
    );
    return res.data.data!;
  },

  async start(
    request?: ModbusShareLifecycleRequest | number,
    expectedSettingsRevision?: string,
  ): Promise<ModbusShareStatus> {
    const normalized = typeof request === 'number'
      ? { port: request, expected_settings_revision: expectedSettingsRevision }
      : request;
    const res = await api.post<APIResponse<ModbusShareStatus>>('/modbus-share/start', {
      ...(normalized?.port !== undefined ? { port: normalized.port } : {}),
      ...(normalized?.expected_settings_revision
        ? { expected_settings_revision: normalized.expected_settings_revision }
        : {}),
    });
    return res.data.data!;
  },

  async stop(expectedSettingsRevision?: string): Promise<ModbusShareStatus> {
    const res = await api.post<APIResponse<ModbusShareStatus>>('/modbus-share/stop', {
      ...(expectedSettingsRevision ? { expected_settings_revision: expectedSettingsRevision } : {}),
    });
    return res.data.data!;
  },

  async listMappings(): Promise<ModbusShareMapping[]> {
    const res = await api.get<APIResponse<ModbusShareMapping[]>>('/modbus-share/mappings');
    return res.data.data ?? [];
  },

  async upsertMapping(tagId: string, register: number): Promise<ModbusShareMapping> {
    const res = await api.put<APIResponse<ModbusShareMapping>>(
      `/modbus-share/mappings/${tagId}`,
      { register },
    );
    return res.data.data!;
  },

  async deleteMapping(tagId: string): Promise<void> {
    await api.delete(`/modbus-share/mappings/${tagId}`);
  },

  async writeTagValue(tagId: string, value: unknown): Promise<void> {
    await api.post('/modbus-share/write-tag-value', { tag_id: tagId, value });
  },

  async sync(): Promise<{ updated: number; skipped: number; errors: string[] }> {
    const res = await api.post<APIResponse<{
      updated: number;
      skipped: number;
      errors: string[];
    }>>('/modbus-share/sync');
    return res.data.data ?? { updated: 0, skipped: 0, errors: [] };
  },
};
