import { strict as assert } from 'node:assert';
import { createHash } from 'node:crypto';
import { createServer, type AddressInfo } from 'node:net';
import { join } from 'node:path';
import { performance } from 'node:perf_hooks';
import { startEmbedded, type EmbeddedBuild } from './embedded-server';
import { assertDevice, assertDisabledSources, assertRule, assertSeedOrigin, type SeedDevice, type SeedPoint, type SeedRule } from './recording-home-seed-contract';
import { assertEmptySources, assertNoConfigurationChanges, assertSeedTableCounts, changedTables, snapshotSQLite } from './recording-home-sqlite';
import type { RecordingHomeSourceAvailable } from '../../../src/types/recordingHomeSources';

export const compactSourcePath = '/api/v1/datalink/studio-v2/recording-home/sources';
const workspacePath = '/api/v1/datalink/studio-v2/workspace';
export interface SeedRequest { phase: 'setup' | 'compact'; method: string; path: string; status: number; bytes: number; sha256: string; elapsedMs: number }
export async function startProtocolSentinel() {
  const connections: { at: string; remoteAddress?: string }[] = [];
  const listener = createServer(socket => { connections.push({ at: new Date().toISOString(), remoteAddress: socket.remoteAddress }); socket.destroy(); });
  await new Promise<void>((done, reject) => { listener.once('error', reject); listener.listen(0, '127.0.0.1', done); });
  const port = (listener.address() as AddressInfo).port;
  return { port, connections, assertQuiet: () => { if (connections.length) throw new Error(`Disabled fixture made ${connections.length} protocol connections`); },
    stop: () => new Promise<void>((done, reject) => listener.close(error => error ? reject(error) : done())) };
}
/** This helper owns startup. There is deliberately no baseURL/server override or seed-existing API. */
export async function seedRecordingHomeFixture(build: EmbeddedBuild) {
  const server = await startEmbedded(build); // A's ownership/identity checks must succeed before any seed request.
  const ledger: SeedRequest[] = [];
  const shareObservations: unknown[] = [];
  let sentinel: Awaited<ReturnType<typeof startProtocolSentinel>> | undefined;
  try {
    assertSeedOrigin(server.baseURL, server.port);
    const databasePath = join(server.runtimeRoot, 'datalink.db');
    const snapshot = () => snapshotSQLite(databasePath);
    const startup = snapshot();
    assertEmptySources(startup);
    sentinel = await startProtocolSentinel();
    const ownedSentinel = sentinel;
    async function request<T>(path: string, body?: unknown, phase: SeedRequest['phase'] = 'setup') {
      assertSeedOrigin(server.baseURL, server.port);
      const method = body === undefined ? 'GET' : 'POST';
      if (method === 'POST' && path !== `${workspacePath}/devices` && path !== `${workspacePath}/source-rules`) throw new Error('Unapproved seed write');
      const started = performance.now();
      const response = await fetch(`${server.baseURL}${path}`, { method, redirect: 'manual', cache: 'no-store', signal: AbortSignal.timeout(30_000),
        ...(body === undefined ? {} : { headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) }) });
      const bytes = Buffer.from(await response.arrayBuffer());
      ledger.push({ phase, method, path, status: response.status, bytes: bytes.length,
        sha256: createHash('sha256').update(bytes).digest('hex'), elapsedMs: performance.now() - started });
      if (response.status !== (method === 'POST' ? 201 : 200)) throw new Error(`Seed ${method} ${path}: HTTP ${response.status}; redirects are never followed`);
      const envelope = JSON.parse(bytes.toString()) as { success?: unknown; data?: T };
      if (envelope.success !== true || envelope.data === undefined || envelope.data === null) throw new Error(`Malformed seed ACK for ${path}`);
      ownedSentinel.assertQuiet();
      return { data: envelope.data, bytes: bytes.length, sha256: ledger.at(-1)!.sha256 };
    }
    interface Workspace { id: string; ordered_device_ids: string[]; modbus_share?: { workspace_revision?: string; readiness?: boolean; status?: { enabled: boolean; configured_enabled: boolean; port: number; mapping_count: number; listener_state: string; bind_state: string } } }
    const assertShareDisabled = (workspace: Workspace) => {
      const status = workspace.modbus_share?.status;
      // Readiness is hydration evidence, not an enabled listener. Check the actual production Status DTO.
      shareObservations.push({ readiness: workspace.modbus_share?.readiness, status: status && { enabled: status.enabled, configured_enabled: status.configured_enabled, port: status.port, mapping_count: status.mapping_count, listener_state: status.listener_state, bind_state: status.bind_state } });
      if (!status || status.enabled !== false || status.configured_enabled !== false || status.port !== 0 || status.mapping_count !== 0 ||
        !['disabled', 'stopped'].includes(status.listener_state) || !['disabled', 'stopped'].includes(status.bind_state)) throw new Error('Unexpected active or unknown Share listener');
    };
    const initial = (await request<Workspace>(workspacePath)).data;
    if (!initial.id || !Array.isArray(initial.ordered_device_ids) || initial.ordered_device_ids.length) throw new Error('Refusing existing workspace membership');
    assertShareDisabled(initial);
    const bootstrap = snapshot();
    const acknowledgments: unknown[] = [];
    for (let index = 1; index <= 20; index++) {
      const device = (await request<SeedDevice>(`${workspacePath}/devices`, { id: `perf-device-${index}`, name: `Synthetic offline ${index}`, protocol: 'modbus_tcp',
        connection_config: { host: '127.0.0.1', port: ownedSentinel.port, slave_id: 1, timeout: 100 } })).data;
      assertDevice(device, index);
      const rule = (await request<SeedRule>(`${workspacePath}/source-rules`, { id: `perf-rule-${index}`, device_id: device.id, start_address: '40001', count: 100,
        data_type: 'int16', naming_prefix: `Perf_${index}_`, enabled: false, share_enabled: false })).data;
      assertRule(rule, initial.id, index);
      if (!rule.runtime_apply_status) throw new Error('Missing seed runtime projection ACK');
      acknowledgments.push({ deviceId: device.id, ruleId: rule.id, revision: rule.revision_id, saveState: rule.save_state, runtimeApplyStatus: rule.runtime_apply_status });
    }
    // Legacy GETs belong to setup only: device GET may update configuration timestamps.
    const workspace = (await request<Workspace>(workspacePath)).data;
    const devices = (await request<SeedDevice[]>(`${workspacePath}/devices`)).data;
    const rules = (await request<SeedRule[]>(`${workspacePath}/source-rules`)).data;
    const points = (await request<SeedPoint[]>('/api/v1/datalink/points')).data;
    assertDisabledSources(workspace.id, devices, rules, points);
    assert.equal(workspace.id, initial.id);
    assert.deepEqual(workspace.ordered_device_ids, Array.from({ length: 20 }, (_, i) => `perf-device-${i + 1}`));
    assertShareDisabled(workspace);
    const expectedDevices = devices.map(({ id, name, protocol, availability_status }) => ({ id, name, protocol, availability: availability_status }));
    const seeded = snapshot(); assertSeedTableCounts(seeded); ownedSentinel.assertQuiet();
    const setup = { workspaceId: workspace.id, workspaceRevision: workspace.modbus_share?.workspace_revision ?? null,
      counts: { devices: devices.length, rules: rules.length, points: points.length, mappings: 0, tags: 0, connectors: 0 },
      acknowledgments, expectedDevices, shareObservations, fixtureSHA256: createHash('sha256').update(JSON.stringify({ workspace, expectedDevices, rules, points })).digest('hex'),
      database: { startup, bootstrap, seeded, bootstrapChanges: changedTables(startup, bootstrap), seedChanges: changedTables(bootstrap, seeded) } };
    const compact = async () => {
      const response = await request<RecordingHomeSourceAvailable>(`${compactSourcePath}?workspace_id=${encodeURIComponent(workspace.id)}`, undefined, 'compact');
      assert.deepEqual(response.data.devices, expectedDevices);
      assert.deepEqual(response.data.candidates, []); assert.deepEqual(response.data.excluded, []);
      assert.equal(response.data.state, 'available'); assert.equal(response.data.workspace_id, workspace.id);
      assert.match(response.data.source_digest, /^[a-f0-9]{64}$/);
      assert.deepEqual(Object.keys(response.data).sort(), ['candidates', 'devices', 'excluded', 'source_digest', 'state', 'workspace_id']);
      for (const device of response.data.devices) assert.deepEqual(Object.keys(device).sort(), ['availability', 'id', 'name', 'protocol']);
      assert(response.bytes < 3_060_000, 'Compact JSON must be smaller than the documented approximate legacy transport, not a timing claim');
      return response;
    };
    return { server, sentinel: ownedSentinel, setup, ledger, snapshot, compact,
      stop: async () => { try { await server.stop(); ownedSentinel.assertQuiet(); } finally { await ownedSentinel.stop(); } } };
  } catch (error) {
    const evidence = { ledger, shareObservations, connections: sentinel?.connections ?? [], serverLog: server.logs() };
    try { await server.stop(); } finally { await sentinel?.stop(); }
    throw Object.assign(new Error(`Recording-home fixture setup failed: ${error instanceof Error ? error.message : String(error)}`, { cause: error }), { evidence });
  }
}
/** Non-browser real HTTP smoke. Configuration hashes, not the GET ledger, establish no persisted config changes. */
export async function measureCompactReads(fixture: Awaited<ReturnType<typeof seedRecordingHomeFixture>>) {
  const before = fixture.snapshot();
  const first = await fixture.compact(); const second = await fixture.compact();
  const after = fixture.snapshot();
  const evidence = { first, second, before, after, changes: changedTables(before, after), protocolConnections: fixture.sentinel.connections,
    limitation: 'Disabled unmapped fixture; canonical empty candidates/excluded. Table equality detects persisted differences, not transient or rolled-back writes.' };
  try {
    fixture.sentinel.assertQuiet(); assertNoConfigurationChanges(before, after);
    assert.equal(first.data.source_digest, second.data.source_digest);
  } catch (error) { throw Object.assign(error instanceof Error ? error : new Error(String(error)), { evidence }); }
  return evidence;
}
