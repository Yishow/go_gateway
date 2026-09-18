import { beforeEach, describe, expect, it, vi } from 'vitest';
import { dbTargetAPI } from '@/services/datalink';
import { api } from '@/services/datalinkClient';

// Only the transport is stubbed; the real list guard stays under test.
vi.mock('@/services/datalinkClient', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/services/datalinkClient')>()),
  api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), patch: vi.fn(), delete: vi.fn() },
  DATALINK_BASE: '/api/v1/datalink',
}));

describe('dbTargetAPI.listTables inspection results', () => {
  beforeEach(() => {
    vi.mocked(api.get).mockReset();
  });

  it('returns the tables reported by the saved connection', async () => {
    const tables = [{
      schema: 'main', name: 'sensor_values',
      columns: [{ name: 'recorded_at', data_type: 'TEXT', nullable: false, primary_key: true }],
    }];
    vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: tables } } as never);

    await expect(dbTargetAPI.listTables('db-1')).resolves.toEqual(tables);
  });

  it('rejects a reply without a table list instead of reporting that no table exists', async () => {
    vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: null } } as never);

    await expect(dbTargetAPI.listTables('db-1')).rejects.toThrow();
  });

  it('keeps a failed inspection request as a failure', async () => {
    vi.mocked(api.get).mockRejectedValueOnce(Object.assign(new Error('Request failed with status code 422'), { response: { status: 422 } }));

    await expect(dbTargetAPI.listTables('db-1')).rejects.toThrow('Request failed with status code 422');
  });
});
