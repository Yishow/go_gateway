import { describe, expect, it, vi } from 'vitest';
import { runtimeAPI } from '@/services/datalink';
import { api } from '@/services/datalinkClient';
vi.mock('@/services/datalinkClient', async (original) => ({
  ...(await original<typeof import('@/services/datalinkClient')>()), api: { get: vi.fn() },
}));
const read = async (data: unknown) => {
  vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data } });
  return runtimeAPI.getStatus();
};
describe('runtime status safe legacy projection', () => {
  it('removes every raw diagnostic text and arbitrary nested field while preserving health and delivery stages', async () => {
    const secret = 'driver: postgres://user:pass@10.1.1.1:5432 SELECT * FROM private /private/db.sqlite';
    const hostile = { error: secret, last_error: secret, reason: secret, failure_reason: secret,
      last_failure_reason: secret, projection_message: secret, availability_reason: secret, nested: { password: secret } };
    const result = await read({ running: true, uptime_seconds: 20, ...hostile,
      metrics: { collected_total: 3, write_success_total: 1, ...hostile },
      snapshot_state: { state: 'degraded', degraded: true, ...hostile },
      collectors: [{ device_id: 'device-A', device_name: 'SELECT secret FROM credentials', protocol: 'modbus_tcp', status: 'error', points_error: 1, ...hostile }],
      diagnostics: [{ scope: 'point', point_id: 'point-1', stages: [{ stage: 'mapping', status: 'failed', ...hostile }], ...hostile }],
      database_delivery: [{ point_id: 'point-1', tag_id: 'tag-1', status: 'failed', stages: ['mapped'], ...hostile }],
      modbus_share_delivery: [{ point_id: 'point-1', tag_id: 'tag-1', status: 'failed', stage: 'mapping', ...hostile }],
    });
    expect(JSON.stringify(result)).not.toMatch(/driver:|private|postgres:\/\/|SELECT|nested|password/);
    expect(result.collectors[0]).toMatchObject({ device_id: 'device-A', points_error: 1 });
    expect(result.database_delivery?.[0]).toMatchObject({ point_id: 'point-1', status: 'failed' });
  });
  it('preserves the existing database stage vocabulary and restart-required alignment', async () => {
    const result = await read({ collectors: [{ device_id: 'device-A', projection_alignment: 'restart-required' }],
      database_delivery: [{ stages: ['collected', 'mapped', 'db_write', 'db_write_failed', 'db_write_succeeded'] }] });
    expect(result.database_delivery?.[0].stages).toEqual(['collected', 'mapped', 'db_write', 'db_write_failed', 'db_write_succeeded']);
    expect(result.collectors[0]).toHaveProperty('projection_alignment', 'restart-required');
  });
  it('does not permit a network address through identity fields', async () => {
    expect(JSON.stringify(await read({ collectors: [{ device_id: '10.20.30.40' }] }))).not.toContain('10.20.30.40');
  });
  it.each(['modbus_udp', 'fatek_fbs', 'mc_3e', 'mqtt'])('preserves supported protocol %s', async (protocol) => {
    expect((await read({ collectors: [{ protocol }] })).collectors[0].protocol).toBe(protocol);
  });
  it('rejects DNS-like identity fields without copying unknown failure codes', async () => {
    const result = await read({ collectors: [{ device_id: 'private.example' }], diagnostics: [{ scope: 'device:private.example', failure_code: 'uid=operator;pwd=secret' }] });
    expect(JSON.stringify(result)).not.toMatch(/private\.example|uid=|pwd=/);
  });
  it('distinguishes an absent failure stage from an unsupported provided stage', async () => {
    const result = await read({ diagnostics: [{ scope: 'runtime', stages: [] }, { scope: 'runtime', failure_stage: 'future-stage', stages: [] }] });
    const diagnostics = (result as import('@/types/runtimeDiagnostics').RuntimeStatusWithDiagnostics).diagnostics!;
    expect(diagnostics[0]).not.toHaveProperty('failure_stage');
    expect(diagnostics[1].failure_stage).toBe('unknown');
  });
  it.each([{ status: 'future-status' }, { availability_status: 'future-state' }, { protocol: 'future-protocol' }])('fails closed for unsupported collector metadata %s', async (unknown) => {
    const result = await read({ collectors: [{ device_id: 'device-A', status: 'running', protocol: 'modbus_tcp', availability_status: 'available', running: true, ...unknown }] });
    expect(result.collectors[0].running).toBe(false);
  });
  it.each(['deferred', 'degraded', 'restart_required'])('preserves known projection alignment %s', async (alignment) => {
    expect((await read({ collectors: [{ projection_alignment: alignment }] })).collectors[0]).toHaveProperty('projection_alignment', alignment);
  });
  it('retains valid running evidence while rejecting missing collector identity', async () => {
    const valid = { device_id: 'device-A', status: 'running', protocol: 'modbus_tcp', availability_status: 'available', running: true };
    expect((await read({ collectors: [valid] })).collectors[0].running).toBe(true);
    expect((await read({ collectors: [{ ...valid, device_id: '' }] })).collectors[0].running).toBe(false);
  });
  it('marks missing status as unavailable rather than a fabricated healthy snapshot', async () => {
    expect(await read(undefined)).toMatchObject({ snapshot_state: { state: 'unavailable', unavailable: true }, collectors: [] });
  });
});
