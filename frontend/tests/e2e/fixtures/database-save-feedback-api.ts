import type { Page, Route } from '@playwright/test';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '../../../src/types/datalink';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from '../../unit/workbench-v2/studio-locator-fixtures';

// Real App/router, simulated HTTP only. No request falls through to a gateway, PLC or database.
export const databaseSavePaths = {
  config: '/api/v1/datalink/studio-v2/workspace/database-config',
  targets: '/api/v1/datalink/studio-v2/workspace/database-targets',
  bootstrap: '/api/v1/datalink/studio-v2/workspace',
} as const;
export type ReadKind = keyof typeof databaseSavePaths;
export const readKinds: ReadKind[] = ['config', 'targets', 'bootstrap'];
export type AckMode = 'valid' | 'missing-id' | 'foreign-workspace' | 'missing-workspace' | 'missing-identity' | 'missing-setup' | 'malformed' | 'conflict' | 'lost' | 'timeout';
type LedgerEntry = { sequence: number; elapsedMs: number; event: string; path?: string; method?: string; detail?: unknown };
function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>(done => { resolve = done; });
  return { promise, resolve };
}
export function heldRead() {
  const arrival = deferred<void>(), gate = deferred<'ok' | 'fail'>(), finish = deferred<void>();
  return { started: arrival.promise, finished: finish.promise, arrive: () => arrival.resolve(),
    outcome: gate.promise, release: () => gate.resolve('ok'), fail: () => gate.resolve('fail'), finish: () => finish.resolve() };
}
function savePlan(mode: AckMode) {
  const received = deferred<Record<string, unknown>>(), gate = deferred<void>(), finish = deferred<void>();
  return { mode, received: received.promise, arrive: received.resolve, gate: gate.promise,
    release: () => gate.resolve(), finished: finish.promise, finish: () => finish.resolve() };
}
const at = '2026-11-06T00:00:00Z';
export function databaseConfig(): StudioV2WorkspaceDatabaseConfigRecord {
  return { id: 'db-W', workspace_id: 'W', identity_revision: 'identity-0', setup_revision: 'setup-0',
    kind: 'sqlite', name: 'Fixture SQLite', host: '', port: 0, database: '/tmp/fixture-only.db', username: '',
    schema: 'main', table: 'readings', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
    status: 'ready', save_state: 'saved', created_at: at, updated_at: at };
}
function safeBody(body: Record<string, unknown>) {
  // Ledger is an attachment: never persist credentials, even if a future test fills one.
  const { password: _password, ...safe } = body;
  return safe;
}
export class DatabaseSaveFeedbackApi {
  readonly ledger: LedgerEntry[] = [];
  readonly unexpected: string[] = [];
  readonly pageErrors: string[] = [];
  readonly writes: Record<string, unknown>[] = [];
  config = databaseConfig();
  private readonly start = Date.now();
  private readonly plans: ReturnType<typeof savePlan>[] = [];
  private readonly saves: ReturnType<typeof savePlan>[] = [];
  private readonly holds: ReturnType<typeof heldRead>[] = [];
  private readonly reads: Record<ReadKind, ReturnType<typeof heldRead>[]> = { config: [], targets: [], bootstrap: [] };
  private disposed = false;
  private revision = 0;
  private readonly device = { ...nativeDevice(), name: 'Fixture Modbus', connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1}' };
  private readonly rule = nativeRule();
  private readonly mapping = { ...nativeMapping(), address: '40001', target_type: 'int16' as const, scale: 1 };
  selectWorkspace(id: string) {
    this.config = { ...databaseConfig(), id: `db-${id}`, workspace_id: id, name: `Fixture ${id} SQLite`,
      identity_revision: `${id}-identity-0`, setup_revision: `${id}-setup-0` };
    this.revision = 0; this.mark('backend-workspace-selected', { workspaceId: id });
  }
  private workspace() {
    return { ...locatorRead().workspace, id: this.config.workspace_id, ordered_device_ids: ['D'],
      modbus_share: { workspace_id: this.config.workspace_id, workspace_revision: `${this.config.workspace_id}2`, settings_revision: 'settings-1',
        hydration_state: 'ready', readiness: false, configured_enabled: false, listener_state: 'stopped' } };
  }
  mark(event: string, detail?: unknown) {
    this.ledger.push({ sequence: this.ledger.length + 1, elapsedMs: Date.now() - this.start, event, detail });
  }
  prepareSave(mode: AckMode = 'valid') {
    const plan = savePlan(mode); this.plans.push(plan); this.saves.push(plan); return plan;
  }
  holdRefreshes() {
    const round = { config: heldRead(), targets: heldRead(), bootstrap: heldRead() };
    for (const kind of readKinds) { this.reads[kind].push(round[kind]); this.holds.push(round[kind]); }
    return round;
  }
  private responseData(path: string, url: URL): { data: unknown; status?: number } | undefined {
    const root = databaseSavePaths.bootstrap;
    if (path === root) {
      // Reuse complete locator fixtures; do not invent a reduced editor-hydration envelope.
      if (url.searchParams.get('entry') === 'locator') {
        const step = Number(url.searchParams.get('step'));
        if (step !== 1 && step !== 2 && step !== 3) return undefined;
        const read = locatorRead(step); read.entry.devices = [this.device];
        return { data: { ...this.workspace(), locator_entry: read.entry } };
      }
      return { data: this.workspace() };
    }
    if (path === databaseSavePaths.config) return { data: this.config };
    if (path === databaseSavePaths.targets) return { data: [] };
    if (path === `${root}/devices`) return { data: [this.device] };
    if (path === `${root}/source-rules`) return { data: [{ ...this.rule, workspace_id: this.config.workspace_id }] };
    if (path === `${root}/mappings`) return { data: [{ ...this.mapping, workspace_id: this.config.workspace_id }] };
    if (path === `${root}/write-groups`) return { data: { workspace_id: this.config.workspace_id, workspace_revision: this.config.setup_revision, groups: [] } };
    if (path === `${root}/audit-history`) return { data: { entries: [] } };
    if (path === `${root}/database-operations/test-writes`) return { data: { workspace_id: this.config.workspace_id, complete: true, operations: [] } };
    if (path === `${root}/measurements` || path === `${root}/measurements/templates`) return { data: [] };
    if (path === `${root}/runtime-context`) return { status: 503, data: null }; // No fabricated runtime readiness.
    if (path === `${root}/database-metadata`) return { data: { workspace_id: this.config.workspace_id,
      connector_id: this.config.id, connector_revision: this.config.identity_revision, database: this.config.database,
      schema: this.config.schema, table: this.config.table, inspection_status: 'missing', columns: [] } };
    if (path === `${root}/database-catalog`) return { data: { workspace_id: this.config.workspace_id,
      workspace_revision: this.config.setup_revision, connector_id: this.config.id, connector_revision: this.config.identity_revision,
      kind: 'sqlite', default_schema: 'main', schemas: [{ name: 'main', tables: [] }] } };
    if (path === '/api/v1/datalink/settings' || path === '/api/v1/datalink/points') return { data: [] };
    if (path === '/api/v1/datalink/modbus-share/status') return { data: { ...this.workspace().modbus_share, enabled: false, running: false, bind_state: 'disabled', port: 5020, address: '127.0.0.1', mapping_count: 0 } };
    if (path === '/api/v1/datalink/db-targets/connectors') return { data: [{ id: this.config.id, identity_revision: this.config.identity_revision,
      kind: 'sqlite', name: this.config.name, enabled: true, status: 'ready', connection_config: { database: this.config.database, schema: 'main' }, created_at: at, updated_at: at }] };
    if (path === '/api/v1/datalink/acquisition-anomalies/summary') return { data: { items: [{ device_id: 'D', active_count: 0 }], total_active: 0 } };
    if (path === '/api/v1/datalink/studio-v2/recording-home/sources') return { data: { workspace_id: this.config.workspace_id,
      state: 'available', source_digest: 'a'.repeat(64), devices: [{ id: 'D', name: this.device.name, protocol: 'modbus_tcp', availability: 'unknown' }], candidates: [], excluded: [] } };
    return undefined;
  }
  async handle(route: Route) {
    const request = route.request(), url = new URL(request.url()), path = url.pathname, method = request.method();
    this.ledger.push({ sequence: this.ledger.length + 1, elapsedMs: Date.now() - this.start, event: 'request', path, method });
    if (method === 'PUT' && path === databaseSavePaths.config) return this.write(route);
    const response = method === 'GET' ? this.responseData(path, url) : undefined;
    if (!response) {
      this.unexpected.push(`${method} ${path}`); this.mark('unmatched-aborted', `${method} ${path}`);
      return route.abort('blockedbyclient');
    }
    // Capture before waiting: a late read MUST carry the old snapshot, not mutable live state.
    const snapshot = structuredClone(response.data);
    const kind = readKinds.find(value => databaseSavePaths[value] === path && !url.searchParams.has('entry'));
    const held = kind ? this.reads[kind].shift() : undefined;
    if (held) { this.mark(`read-start:${kind}`, snapshot); held.arrive(); }
    const outcome = held ? await held.outcome : 'ok';
    try {
      if (this.disposed) return await route.abort('aborted');
      const status = outcome === 'fail' ? 503 : response.status ?? 200;
      await route.fulfill({ status, json: status === 200 ? { success: true, data: snapshot }
        : { success: false, error: { code: 'fixture_read_unavailable', message: 'Observation unavailable' } } });
      this.mark(held ? `read-end:${kind}:${outcome}` : 'read-response', { path, status, snapshot });
    } catch (error) {
      // C deliberately cancels superseded reads. Retain the attempted old release as evidence.
      if (!request.failure() && !this.disposed) throw error;
      this.mark(`read-cancelled:${kind ?? path}`);
    } finally { held?.finish(); }
  }
  private async write(route: Route) {
    const plan = this.plans.shift(), body = route.request().postDataJSON() as Record<string, unknown>;
    this.writes.push(safeBody(body)); this.mark('config-submit', safeBody(body));
    if (!plan) { this.unexpected.push('unplanned config PUT'); return route.abort('blockedbyclient'); }
    plan.arrive(safeBody(body));
    await plan.gate;
    try {
      if (this.disposed) return await route.abort('aborted');
      if (plan.mode === 'conflict') {
        await route.fulfill({ status: 409, json: { success: false, error: { code: 'workspace_revision_conflict', message: 'Configuration changed' } } });
        this.mark('config-rejected:409'); return;
      }
      // Lost/timeout ACK simulates a committed write with no reply. Only GET can observe it later.
      this.revision++;
      this.config = { ...this.config, ...safeBody(body), identity_revision: `identity-${this.revision}`,
        setup_revision: `setup-${this.revision}`, updated_at: `2026-11-06T00:00:0${this.revision}Z` } as StudioV2WorkspaceDatabaseConfigRecord;
      if (plan.mode === 'lost' || plan.mode === 'timeout') {
        await route.abort(plan.mode === 'timeout' ? 'timedout' : 'connectionclosed'); this.mark(`config-ack-${plan.mode}`); return;
      }
      const ack: Record<string, unknown> = { ...this.config, runtime_apply_status: 'apply_failed', runtime_apply_code: 'fixture_runtime_not_applied' };
      if (plan.mode === 'missing-id') delete ack.id;
      if (plan.mode === 'foreign-workspace') ack.workspace_id = 'foreign-workspace';
      if (plan.mode === 'missing-workspace') delete ack.workspace_id;
      if (plan.mode === 'missing-identity') delete ack.identity_revision;
      if (plan.mode === 'missing-setup') delete ack.setup_revision;
      await route.fulfill({ json: { success: true, data: plan.mode === 'malformed' ? null : ack } });
      this.mark('config-ack', plan.mode === 'malformed' ? null : ack);
    } finally { plan.finish(); }
  }
  async install(page: Page) {
    await page.addInitScript(() => { localStorage.setItem('i18nextLng', 'zh-TW'); });
    page.on('pageerror', error => this.pageErrors.push(error.message));
    // Catch all HTTP requests: same-origin app assets may continue; every other transport aborts.
    await page.route('**/*', route => {
      const url = new URL(route.request().url());
      if (url.pathname.startsWith('/api/')) return this.handle(route);
      const request = route.request();
      const asset = ['script', 'stylesheet', 'font', 'image', 'manifest'].includes(request.resourceType());
      const document = request.resourceType() === 'document' && url.pathname === '/studio/v2';
      if (url.origin === 'http://127.0.0.1:4173' && request.method() === 'GET' && (asset || document)) return route.continue();
      this.unexpected.push(`non-fixture transport ${request.method()} ${url.origin}${url.pathname}`);
      return route.abort('blockedbyclient');
    });
  }
  dispose() {
    this.disposed = true;
    this.holds.forEach(held => held.release()); this.saves.forEach(save => save.release());
  }
}
