// @vitest-environment node
import { afterEach, describe, expect, it, vi } from 'vitest';
import { createServer, type Server } from 'node:http';
import { connect, type AddressInfo } from 'node:net';
import { DatabaseSync } from 'node:sqlite';
import { mkdtempSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { startEmbedded, type EmbeddedBuild } from '../../e2e/helpers/embedded-server';
import { measureCompactReads, seedRecordingHomeFixture, startProtocolSentinel } from '../../e2e/helpers/recording-home-seed';
vi.mock('../../e2e/helpers/embedded-server', () => ({ startEmbedded: vi.fn() }));
const tables = ['devices', 'source_rules', 'source_rule_links', 'points', 'tags', 'mappings', 'database_connectors', 'database_target_mappings', 'write_groups', 'recording_plans'];
const roots: string[] = []; const servers: Server[] = []; const databases: DatabaseSync[] = [];
afterEach(async () => {
  vi.clearAllMocks();
  for (const server of servers.splice(0)) await new Promise<void>(done => server.close(() => done()));
  for (const db of databases.splice(0)) db.close();
  for (const root of roots.splice(0)) rmSync(root, { recursive: true, force: true });
});
async function stub(mode = 'valid') {
  const root = mkdtempSync(join(tmpdir(), 'home-seed-unit-')); roots.push(root);
  const db = new DatabaseSync(join(root, 'datalink.db')); databases.push(db);
  for (const table of tables) db.exec(`CREATE TABLE ${table}(id TEXT)`);
  if (mode === 'existing') db.exec("INSERT INTO devices VALUES ('user-device')");
  const calls: string[] = []; const devices: Record<string, unknown>[] = []; const rules: Record<string, unknown>[] = []; const points: Record<string, unknown>[] = [];
  let redirectRequests = 0;
  const server = createServer(async (request, response) => {
    calls.push(`${request.method} ${request.url}`);
    if (request.url === '/redirected') { redirectRequests++; response.end('{}'); return; }
    const chunks: Buffer[] = []; for await (const chunk of request) chunks.push(Buffer.from(chunk));
    const body = chunks.length ? JSON.parse(Buffer.concat(chunks).toString()) as Record<string, unknown> : {};
    let data: unknown;
    if (request.method === 'POST' && mode === 'redirect') { response.writeHead(307, { Location: '/redirected' }); response.end(); return; }
    if (request.url?.endsWith('/workspace')) data = { id: 'workspace', ordered_device_ids: devices.map(d => d.id), modbus_share: mode === 'missing-share' ? undefined : { readiness: true, status: { enabled: mode === 'active-share', configured_enabled: false, port: 0, mapping_count: 0, listener_state: 'disabled', bind_state: 'disabled' } } };
    else if (request.url?.endsWith('/devices') && request.method === 'POST') {
      data = { ...body, status: mode === 'active-device' ? 'active' : 'draft', running: false, availability_status: 'available', last_test_at: null, last_test_success: null };
      devices.push(data as Record<string, unknown>); db.prepare('INSERT INTO devices VALUES (?)').run(String(body.id));
    } else if (request.url?.endsWith('/source-rules') && request.method === 'POST') {
      data = { ...body, workspace_id: 'workspace', revision_id: `rev-${rules.length}`, save_state: mode === 'bad-ack' ? 'pending' : 'saved', runtime_apply_status: mode === 'missing-runtime' ? undefined : 'not_running' };
      rules.push(data as Record<string, unknown>); db.prepare('INSERT INTO source_rules VALUES (?)').run(String(body.id));
      for (let i = 0; i < 100; i++) {
        const point = { id: `${String(body.id)}-${i}`, device_id: body.device_id, enabled: false, last_read_at: null, last_value: null };
        points.push(point); db.prepare('INSERT INTO points VALUES (?)').run(point.id); db.prepare('INSERT INTO source_rule_links VALUES (?)').run(point.id);
      }
    } else if (request.url?.endsWith('/devices')) data = devices;
    else if (request.url?.endsWith('/source-rules')) data = rules;
    else if (request.url?.endsWith('/points')) data = mode === 'bad-count' ? points.slice(1) : points;
    else { response.writeHead(404); response.end(); return; }
    response.writeHead(request.method === 'POST' ? 201 : 200, { 'Content-Type': 'application/json' });
    response.end(JSON.stringify(mode === 'malformed' ? { success: false, data } : { success: true, data }));
  });
  servers.push(server); await new Promise<void>(done => server.listen(0, '127.0.0.1', done));
  const port = (server.address() as AddressInfo).port;
  const stop = vi.fn(async () => {});
  const handle = { baseURL: `http://127.0.0.1:${port}`, port, runtimeRoot: root, stop, logs: () => 'test-only stub, not production evidence' };
  vi.mocked(startEmbedded).mockResolvedValue(handle);
  return { handle, calls, redirectRequests: () => redirectRequests };
}
describe('owned fixture seeding', () => {
  it('starts through A owner, uses only 40 supported disabled writes, and keeps the sentinel alive through measurement', async () => {
    const fake = await stub(); const build = {} as EmbeddedBuild;
    const fixture = await seedRecordingHomeFixture(build);
    expect(startEmbedded).toHaveBeenCalledWith(build);
    expect(fixture.setup.counts).toEqual({ devices: 20, rules: 20, points: 2000, mappings: 0, tags: 0, connectors: 0 });
    expect(fixture.ledger.filter(item => item.method === 'POST')).toHaveLength(40);
    expect(fake.calls.some(call => /activate|poll|test|readiness/.test(call))).toBe(false);
    fixture.sentinel.assertQuiet(); await fixture.stop(); expect(fake.handle.stop).toHaveBeenCalledOnce();
  });
  it.each(['existing', 'malformed', 'active-device', 'bad-ack', 'bad-count', 'redirect', 'active-share', 'missing-share', 'missing-runtime'])('stops owned process on %s without fallback writes', async mode => {
    const fake = await stub(mode);
    await expect(seedRecordingHomeFixture({} as EmbeddedBuild)).rejects.toThrow('fixture setup failed');
    expect(fake.handle.stop).toHaveBeenCalledOnce();
    expect(fake.redirectRequests()).toBe(0);
    if (mode === 'existing') expect(fake.calls).toEqual([]);
    if (mode === 'malformed') expect(fake.calls.filter(call => call.startsWith('POST'))).toEqual([]);
    if (mode === 'redirect' || mode === 'active-device') expect(fake.calls.filter(call => call.startsWith('POST'))).toHaveLength(1);
  });
  it('rejects arbitrary origin even if a caller forges a structural server in a mocked owner', async () => {
    const fake = await stub(); fake.handle.baseURL = 'http://example.com';
    await expect(seedRecordingHomeFixture({} as EmbeddedBuild)).rejects.toThrow('owned loopback');
    expect(fake.calls).toEqual([]); expect(fake.handle.stop).toHaveBeenCalledOnce();
  });
  it('does not seed if the A ownership helper rejects startup', async () => {
    vi.mocked(startEmbedded).mockRejectedValue(new Error('identity mismatch'));
    await expect(seedRecordingHomeFixture({} as EmbeddedBuild)).rejects.toThrow('identity mismatch');
  });
  it('retains both SQLite snapshots when a measured compact read changes configuration', async () => {
    const before = { reader: 'test', tables: { devices: { rows: 20, sha256: 'before', kind: 'configuration' } } };
    const after = { reader: 'test', tables: { devices: { rows: 20, sha256: 'after', kind: 'configuration' } } };
    const fixture = { snapshot: vi.fn().mockReturnValueOnce(before).mockReturnValueOnce(after),
      compact: async () => ({ data: { source_digest: 'same' }, bytes: 1, sha256: 'body' }),
      sentinel: { assertQuiet: () => {}, connections: [] } } as unknown as Parameters<typeof measureCompactReads>[0];
    await expect(measureCompactReads(fixture)).rejects.toMatchObject({ evidence: { before, after, changes: [{ name: 'devices' }] } });
  });
  it('sentinel observes an actual local TCP attempt and never supplies protocol data', async () => {
    const sentinel = await startProtocolSentinel();
    try {
      const received: Buffer[] = [];
      await new Promise<void>((done, reject) => { const socket = connect(sentinel.port, '127.0.0.1'); socket.on('data', data => received.push(Buffer.from(data))); socket.once('error', reject); socket.once('close', () => done()); });
      expect(received).toEqual([]); expect(sentinel.connections).toHaveLength(1); expect(() => sentinel.assertQuiet()).toThrow('protocol connections');
    } finally { await sentinel.stop(); }
  });
});
