import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { createRuntimeSetupFixture } from '../runtime-dashboard/runtimeSetupFixture';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const read = async (data: unknown) => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data } });
  return studioV2RuntimeContextAPI.get();
};
const fixture = () => ({ workspace_id: 'ws-1', default_device_id: 'device-A', devices: [{
  device_id: 'device-A', name: 'Pump A', protocol: 'modbus_tcp', running: false, availability_status: 'available',
}], setup: createRuntimeSetupFixture() });
describe('runtime context safe legacy projection', () => {
  beforeEach(() => vi.clearAllMocks());
  it('discards legacy credentials, endpoints, private database paths and nested payloads at every level', async () => {
    const data = fixture();
    const hostile = { host: '10.20.30.40', port: 5432, username: 'secret-user', password: 'secret-password',
      database: '/private/tenant/database.sqlite', dsn: 'postgres://secret-user:secret-password@10.20.30.40',
      payload: { nested: { error: 'SELECT secret FROM credentials' } } };
    Object.assign(data, hostile);
    Object.assign(data.setup.database_config!, hostile);
    Object.assign(data.setup.source_rules[0], hostile);
    Object.assign(data.setup.mappings[0], hostile);
    Object.assign(data.devices[0], hostile, { projection_message: 'driver: SELECT secret FROM credentials', availability_reason: '/private/tenant/database.sqlite' });
    const result = await read(data);
    const serialized = JSON.stringify(result);
    for (const secret of ['10.20.30.40', '5432', 'secret-user', 'secret-password', '/private/', 'SELECT secret', 'dsn', 'payload']) expect(serialized).not.toContain(secret);
    expect(result.setup?.source_rules[0].start_address).toBe('40001');
    expect(result.setup?.database_config).toMatchObject({ name: 'PostgreSQL Connector', table: 'sensor_readings', schema: 'public' });
  });
  it('degrades missing legacy fields safely without inventing effective evidence', async () => {
    expect(await read({ workspace_id: 'ws-1', devices: [] })).toEqual({ workspace_id: 'ws-1', devices: [], default_device_id: null, setup: null });
    const data = fixture();
    data.setup.database_config = undefined;
    const result = await read(data);
    expect(result.setup?.database_config).toBeNull();
    expect(result).not.toHaveProperty('effective');
  });
  it('preserves saved-only provenance and validated register notations without inventing unknown facts', async () => {
    const data = fixture();
    Object.assign(data.setup, { evidence_state: 'saved-only' });
    Object.assign(data.setup.source_rules[0], { start_address: 'D0100', data_type: 'future-type' });
    Object.assign(data.setup.database_config!, { write_mode: 'future-mode', kind: 'future-kind' });
    const result = await read(data);
    expect(result.setup).toHaveProperty('evidence_state', 'saved-only');
    expect(result.setup?.source_rules[0]).toMatchObject({ start_address: 'D0100', data_type: '' });
    expect(result.setup?.database_config).toMatchObject({ write_mode: 'unavailable', kind: 'unavailable' });
  });
  it.each(['D0100', '0xFFFF', '400001', 'R100.2'])('preserves validated source register %s', async (address) => {
    const data = fixture(); data.setup.source_rules[0].start_address = address;
    expect((await read(data)).setup?.source_rules[0].start_address).toBe(address);
  });
  it('does not permit a network address through identity fields', async () => {
    const data = fixture(); data.devices[0].device_id = '10.20.30.40';
    expect(JSON.stringify(await read(data))).not.toContain('10.20.30.40');
  });
  it.each(['line01.temp.inlet', 'factory/line_1.temp'])('retains semantic tag key %s rather than treating it as an opaque ID', async (key) => {
    const data = fixture(); data.setup.mappings[0].tag_key = key;
    expect((await read(data)).setup?.mappings[0].tag_key).toBe(key);
  });
  it.each(['file:private.sqlite', '/private/data.sqlite', 'uid=operator;pwd=secret', 'https://private.example', '../secret', 'line/../secret', './secret', 'line/./secret'])('rejects unsafe tag-key syntax %s', async (key) => {
    const data = fixture(); data.setup.mappings[0].tag_key = key;
    expect((await read(data)).setup?.mappings[0].tag_key).toBe('');
  });
  it.each(['m/s', 'µΩ', 'kWh/m³'])('retains safe engineering unit %s', async (unit) => {
    const data = fixture(); data.setup.mappings[0].unit = unit;
    expect((await read(data)).setup?.mappings[0].unit).toBe(unit);
  });
  it.each([['W7FFFFFFF', 'W7FFFFFFF'], ['W80000000', ''], ['WFFFFFFFF', '']])('validates MC register range %s', async (address, expected) => {
    const data = fixture(); data.setup.source_rules[0].start_address = address;
    expect((await read(data)).setup?.source_rules[0].start_address).toBe(expected);
  });
  it.each([{ availability_status: 'future-state' }, { protocol: 'future-protocol' }])('fails closed for unsupported saved device metadata %s', async (unknown) => {
    const data = fixture(); Object.assign(data.devices[0], { running: true, ...unknown });
    const result = await read(data);
    expect(result.devices[0].running).toBe(false);
    expect(result.default_device_id).toBeNull();
  });
  it.each(['deferred', 'degraded', 'restart_required'])('preserves known saved projection alignment %s', async (alignment) => {
    const data = fixture(); Object.assign(data.devices[0], { projection_alignment: alignment });
    expect((await read(data)).devices[0].projection_alignment).toBe(alignment);
  });
  it('does not expose a forward-slash Windows private path as destination metadata', async () => {
    const data = fixture(); data.setup.database_config!.table = 'C:/private/fixture.sqlite';
    expect((await read(data)).setup?.database_config?.table).toBe('');
  });
  it('retains only an available known default, never invents a replacement', async () => {
    expect((await read(fixture())).default_device_id).toBe('device-A');
    const data = fixture(); data.default_device_id = 'unknown-device';
    expect((await read(data)).default_device_id).toBeNull();
  });
  it.each(['path=/private/db', 'path=C:/private/db', '(/private/db)', '\\\\server\\share\\private', './private/db', '../private/db'])(
    'keeps private path %s out of service-projected destination metadata', async (path) => {
      const data = fixture(); Object.assign(data.setup.database_config!, { name: path, schema: path, table: path });
      expect((await read(data)).setup?.database_config).toMatchObject({ name: '', schema: '', table: '' });
    },
  );
  it('rejects DNS-like identity fields', async () => {
    const data = fixture(); data.devices[0].device_id = 'private.example';
    expect(JSON.stringify(await read(data))).not.toContain('private.example');
  });
  it('does not invent a readiness owner step for unknown metadata', async () => {
    const data = fixture(); Object.assign(data.setup.readiness_summary!.issues[0], { step: 'future-step' });
    expect((await read(data)).setup?.readiness_summary?.issues).not.toEqual(expect.arrayContaining([expect.objectContaining({ code: 'database-target-missing' })]));
  });
  it('never reflects hostile strings in nominally safe labels, identifiers, addresses or readiness messages', async () => {
    const data = fixture();
    data.devices[0].name = 'SELECT secret FROM credentials';
    data.setup.database_config!.name = '/private/tenant/database.sqlite';
    data.setup.database_config!.table = 'SELECT * FROM secrets';
    data.setup.source_rules[0].start_address = 'https://private-host';
    data.setup.mappings[0].display_name = 'password=secret';
    data.setup.readiness_summary!.issues[0].message = 'pq: SELECT * FROM secrets';
    const result = await read(data);
    expect(JSON.stringify(result)).not.toMatch(/private|SELECT|password=|postgres:\/\//);
    expect(result.setup?.source_rules[0].start_address).toBe('');
  });
});
