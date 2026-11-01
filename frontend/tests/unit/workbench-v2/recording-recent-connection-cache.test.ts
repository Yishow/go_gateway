import { describe, expect, it } from 'vitest';
import { connectionWhitelist, createRecentConnectionCache } from '../../../src/features/datalink/workbench-v2/recording-editor/recentConnectionDefaults';
import type { DeviceSourceOutcome, DeviceSourceRequest } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';
const request = (id = 'one', config: Record<string, unknown> = { host: '192.0.2.20', port: 502, slave_id: 3, timeout: 5 }): DeviceSourceRequest => ({
  workspaceId: 'workspace-1', sessionId: 'form-1', requestId: id, generation: 7,
  buffer: { name: 'Private PLC name', description: 'Private note', protocol: 'modbus_tcp', config },
});
function confirmed(input: DeviceSourceRequest): DeviceSourceOutcome {
  return { request: structuredClone(input), status: 'confirmed', device: {
    id: 'new-device', protocol: input.buffer.protocol, runtime_apply_status: 'not_running',
  } as StudioV2WorkspaceDeviceRecord };
}
function remember(cache: ReturnType<typeof createRecentConnectionCache>, input: DeviceSourceRequest) {
  cache.capture(input)(confirmed(input));
}
const unwanted = { name: 'old', id: 'old', password: 'secret', token: 'secret', credentials: { password: 'secret' },
  enabled: true, activation: 'running', readiness: 'ready', receipt: 'receipt', operation: 'operation', ownership: 'owned',
  profile: 'fx5', readable_coverage: ['D0-D100'], capability: { random: true }, points: ['D0'], mappings: ['M'], vendor_key: 'unknown' };
describe('protocol-schema connection whitelist', () => {
  it.each(['modbus_tcp', 'modbus_udp'] as const)('only reuses non-sensitive %s connection parameters', protocol => {
    expect(connectionWhitelist({ protocol, config: { ...unwanted, host: '192.0.2.5', port: 503, slave_id: 4, timeout: 8, data_format: 'CDAB' } })).toEqual({
      protocol, config: { host: '192.0.2.5', port: 503, slave_id: 4, timeout: 8, data_format: 'CDAB' },
    });
  });
  it('normalizes RTU aliases and copies all serial settings but no identities or evidence', () => {
    expect(connectionWhitelist({ protocol: 'modbus_rtu', config: { ...unwanted, port: '/dev/ttyS4', baud: 19200, parity: 'E', data_bits: 7, stop_bits: 2, slave_id: 5, timeout: 9, data_format: 'BADC' } })).toEqual({
      protocol: 'modbus_rtu', config: { serial_port: '/dev/ttyS4', baud_rate: 19200, parity: 'even', data_bits: 7, stop_bits: 2, slave_id: 5, timeout: 9, data_format: 'BADC' },
    });
  });
  it('normalizes MC aliases with canonical priority and excludes FX5 profiles and coverage', () => {
    expect(connectionWhitelist({ protocol: 'mc_3e', config: { ...unwanted, host: '192.0.2.6', port: 6000, network_no: 2, pc_no: 254, station: 6, module_station_no: 7, station_no: 8, module_io_no: 1023, io_no: 1000, timeout: 12, data_format: 'DCBA' } })).toEqual({
      protocol: 'mc_3e', config: { host: '192.0.2.6', port: 6000, network_no: 2, pc_no: 254, station_no: 8, io_no: 1000, timeout: 12, data_format: 'DCBA' },
    });
  });
  it('isolates FATEK mode fields, retaining station aliases but never inactive fields', () => {
    const config = { ...unwanted, host: '192.0.2.7', port: 500, serial_port: '/dev/ttyS5', baud_rate: 4800, parity: 'odd', data_bits: 7, stop_bits: 2, station: 12, timeout: 11 };
    expect(connectionWhitelist({ protocol: 'fatek_fbs', config: { ...config, mode: 'tcp' } })?.config).toEqual({ mode: 'tcp', host: '192.0.2.7', port: 500, station_no: 12, timeout: 11 });
    expect(connectionWhitelist({ protocol: 'fatek_fbs', config: { ...config, mode: 'serial' } })?.config).toEqual({ mode: 'serial', serial_port: '/dev/ttyS5', baud_rate: 4800, parity: 'odd', data_bits: 7, stop_bits: 2, station_no: 12, timeout: 11 });
    expect(connectionWhitelist({ protocol: 'fatek_fbs', config: { mode: 'unsupported' } })).toBeUndefined();
  });
  it('cannot clone credentials through nested allowlisted values or enable MQTT', () => {
    expect(connectionWhitelist({ protocol: 'modbus_tcp', config: { host: { password: 'secret' }, port: Number.NaN, data_format: ['secret'] } })?.config).toEqual({});
    expect(connectionWhitelist({ protocol: 'mqtt', config: { broker: 'mqtt://localhost', username: 'user', password: 'secret' } })).toBeUndefined();
  });
});
describe('page-session causal connection cache', () => {
  it('keeps an immutable request snapshot, not current form values, and returns defensive copies', () => {
    const cache = createRecentConnectionCache('workspace-1'), input = request();
    const saved = confirmed(input), publish = cache.capture(input);
    input.buffer.config.host = 'new unsaved value'; publish(saved);
    const value = cache.port.read('workspace-1');
    expect(value).toEqual({ protocol: 'modbus_tcp', config: { host: '192.0.2.20', port: 502, slave_id: 3, timeout: 5 } });
    value!.config.host = 'mutated copy';
    expect(cache.port.read('workspace-1')?.config.host).toBe('192.0.2.20');
  });
  it.each(['pending', 'failed', 'unknown', 'blocked', 'persisted-unavailable'] as const)('never replaces a confirmed cache with %s', status => {
    const cache = createRecentConnectionCache('workspace-1'); remember(cache, request());
    const input = request('two', { host: 'bad newer', port: 502 }); cache.capture(input)({ ...confirmed(input), status });
    expect(cache.port.read('workspace-1')?.config.host).toBe('192.0.2.20');
  });
  it('rejects failed runtime apply and existing device updates, without populating an empty cache', () => {
    const cache = createRecentConnectionCache('workspace-1'), input = request(), outcome = confirmed(input);
    outcome.device!.runtime_apply_status = 'apply_failed'; cache.capture(input)(outcome);
    remember(cache, { ...input, deviceId: 'existing' });
    expect(cache.port.read('workspace-1')).toBeUndefined();
  });
  it.each(['workspaceId', 'sessionId', 'requestId', 'generation'] as const)('requires the captured %s', field => {
    const cache = createRecentConnectionCache('workspace-1'), input = request(), outcome = confirmed(input);
    outcome.request = { ...outcome.request, [field]: field === 'generation' ? 99 : 'different' };
    cache.capture(input)(outcome); expect(cache.port.read('workspace-1')).toBeUndefined();
  });
  it('rejects changed subtype or snapshot even when the request identity fields match', () => {
    const cache = createRecentConnectionCache('workspace-1'), input = request(), outcome = confirmed(input);
    outcome.request.buffer.config.host = 'not saved'; cache.capture(input)(outcome);
    expect(cache.port.read('workspace-1')).toBeUndefined();
  });
  it('clears all FATEK subtypes but not another protocol and rejects pre-clear callbacks', () => {
    const cache = createRecentConnectionCache('workspace-1'); remember(cache, request());
    const tcp = { ...request('tcp'), buffer: { ...request().buffer, protocol: 'fatek_fbs' as const, config: { mode: 'tcp', host: '192.0.2.9', port: 500, station: 3, timeout: 5 } } };
    const serial = { ...tcp, requestId: 'serial', buffer: { ...tcp.buffer, config: { mode: 'serial', serial_port: '/dev/ttyS6', baud_rate: 19200, station: 4 } } };
    remember(cache, tcp); remember(cache, serial);
    const late = cache.capture({ ...tcp, requestId: 'late' });
    expect(cache.port.read('workspace-1', 'fatek_fbs', 'serial')?.config.serial_port).toBe('/dev/ttyS6');
    cache.port.clear('workspace-1', 'fatek_fbs'); late(confirmed({ ...tcp, requestId: 'late' }));
    expect(cache.port.read('workspace-1', 'fatek_fbs', 'tcp')).toBeUndefined();
    expect(cache.port.read('workspace-1', 'fatek_fbs', 'serial')).toBeUndefined();
    expect(cache.port.read('workspace-1')?.protocol).toBe('modbus_tcp');
    remember(cache, { ...serial, requestId: 'after-clear' });
    expect(cache.port.read('workspace-1')?.config.mode).toBe('serial');
  });
  it('does not let an older request overwrite a newer confirmed connection or last selection', () => {
    const cache = createRecentConnectionCache('workspace-1'), old = request(), publish = cache.capture(old);
    remember(cache, request('newer', { host: '192.0.2.80' })); publish(confirmed(old));
    expect(cache.port.read('workspace-1')?.config.host).toBe('192.0.2.80');
  });
  it('invalidates scope A→B→A and page teardown callbacks, with no persistence or TTL', () => {
    const cache = createRecentConnectionCache('workspace-1'), input = request(), late = cache.capture(input);
    remember(cache, input); cache.resetScope('workspace-2');
    expect(cache.port.read('workspace-1')).toBeUndefined(); expect(cache.port.read('workspace-2')).toBeUndefined();
    cache.resetScope('workspace-1'); late(confirmed(input)); expect(cache.port.read('workspace-1')).toBeUndefined();
    const unmounted = cache.capture(input); remember(cache, input); cache.dispose(); unmounted(confirmed(input));
    expect(cache.port.read('workspace-1')).toBeUndefined();
    expect(createRecentConnectionCache('workspace-1').port.read('workspace-1')).toBeUndefined();
  });
  it('foreign workspace clear/read cannot change the current workspace cache', () => {
    const cache = createRecentConnectionCache('workspace-1'); remember(cache, request());
    cache.port.clear('workspace-2', 'modbus_tcp');
    expect(cache.port.read('workspace-2')).toBeUndefined();
    expect(cache.port.read('workspace-1')?.config.host).toBe('192.0.2.20');
  });
});
