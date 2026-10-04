import assert from 'node:assert/strict';
import { spawn } from 'node:child_process';
import { existsSync, mkdirSync, rmSync, statSync } from 'node:fs';
import { createHash } from 'node:crypto';
import { join } from 'node:path';
import {
  assertFreshPortFree, DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, closeFreshUI, configureFreshSQLite,
  makeFreshRun, monotonicNs, observeJSON, openFreshAdvanced, openFreshUI,
  prepareFreshBasicGroup, setupFreshDevice, setupFreshRules, startFreshGateway, startFreshSimulator,
  stopProcesses,
} from './fresh-ui.mjs';
import { assertBasicLayout, inspectBasicLayout, screenshotAnchor } from './fresh-recovery-assertions.mjs';
import {
  makeOwnedPostgresProxy, parseOwnedPostgresDSN, postgresExec, postgresJSON, quoteIdentifier,
  readPostgresReceipt, readPostgresTable, readSQLiteReceipt, readSQLiteTable, sqliteExec, sqliteJSON, sqliteOwnedNeighborExec,
} from './fresh-recovery-sql.mjs';
import { persistRecoveryMappings } from './fresh-recovery-mappings.mjs';
import { configureFreshPostgres } from './fresh-postgres-setup.mjs';

const RULES = [
  { name: 'A int16', start: 40001, count: 1, dataType: 'int16', prefix: 'A_INT16_' },
  { name: 'A uint16', start: 40002, count: 1, dataType: 'uint16', prefix: 'A_UINT16_' },
  { name: 'A uint64', start: 40003, count: 1, dataType: 'uint64', prefix: 'A_U64_' },
  { name: 'A bool', start: 40007, count: 1, dataType: 'bool', prefix: 'A_BOOL_' },
];
const TARGET_TYPES = ['int16', 'uint16', 'uint64', 'bool'];
const REGISTERS = { 0: 215, 1: 1013, 2: 32, 6: 1 };
const CONNECTOR_TABLE = 'f_connector_placeholder';

function sleep(ms) { return new Promise((resolve) => setTimeout(resolve, ms)); }

async function holdSQLiteTargetWrites(path) {
  assert.match(path, /^\/tmp\/gw-f-[A-Za-z0-9_-]+\/destination\.db$/);
  const script = `import sqlite3,sys,urllib.parse
db=sqlite3.connect('file:'+urllib.parse.quote(sys.argv[1],safe='/')+'?mode=rw',uri=True,timeout=2)
db.execute('BEGIN IMMEDIATE')
print('READY',flush=True)
sys.stdin.buffer.read(1)
db.rollback()
db.close()
`;
  const child = spawn('python3', ['-u', '-c', script, path], { stdio: ['pipe', 'pipe', 'pipe'] });
  child.fHarnessName = 'sqlite-target-write-lock';
  let stderr = '';
  child.stderr.on('data', (value) => { stderr += value.toString(); });
  const exited = new Promise((resolve) => child.once('exit', resolve));
  const ready = new Promise((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error(`owned SQLite lock was not ready: ${stderr}`)), 5000);
    child.once('error', (error) => { clearTimeout(timer); reject(error); });
    child.once('exit', (code) => { clearTimeout(timer); reject(new Error(`owned SQLite lock exited ${code}: ${stderr}`)); });
    child.stdout.once('data', (value) => {
      clearTimeout(timer);
      if (value.toString().trim() === 'READY') resolve();
      else reject(new Error('owned SQLite lock returned an unexpected readiness message'));
    });
  });
  await ready;
  let closed = false;
  return {
    child,
    pid: child.pid,
    sql: 'BEGIN IMMEDIATE; ROLLBACK',
    close: async () => {
      if (closed) return;
      closed = true;
      child.stdin.end();
      const [code] = await Promise.all([exited]);
      if (code !== 0) throw new Error(`owned SQLite lock exited ${code}: ${stderr}`);
    },
  };
}

async function eventually(label, read, accept, timeoutMs = 90_000, intervalMs = 1000) {
  const deadline = Date.now() + timeoutMs;
  let last;
  while (Date.now() < deadline) {
    last = await read();
    if (accept(last)) return last;
    await sleep(intervalMs);
  }
  throw new Error(`${label} did not reach expected state: ${JSON.stringify(last).slice(0, 800)}`);
}

async function waitForGateway(base, timeoutMs = 60_000) {
  const deadline = Date.now() + timeoutMs;
  let last = 'no response';
  while (Date.now() < deadline) {
    try {
      const timeout = Math.max(1, Math.min(3000, deadline - Date.now()));
      const signal = AbortSignal.timeout(timeout);
      const [apiResponse, uiResponse] = await Promise.all([
        fetch(`${base}/api/v1/datalink/studio-v2/workspace`, { signal }),
        fetch(`${base}/studio/v2`, { signal }),
      ]);
      const body = await apiResponse.json();
      last = `api=${apiResponse.status} ui=${uiResponse.status} id=${body?.data?.id ?? ''} readiness=${body?.data?.readiness_summary != null}`;
      if (apiResponse.status === 200 && uiResponse.status === 200 && body?.data?.id && body?.data?.readiness_summary != null) return body.data;
    } catch (error) {
      last = error instanceof Error ? error.message.split('\n')[0] : String(error);
    }
    await sleep(300);
  }
  throw new Error(`fresh gateway did not start: ${base} (${last})`);
}

function tableName(group) { return group?.destination?.table_name ?? ''; }

function parseDSN() { return parseOwnedPostgresDSN(process.env.POSTGRES_DSN); }

function safeOutbox(path, groupID) {
  try {
    const scope = groupID ? `WHERE group_id = '${groupID.replaceAll("'", "''")}'` : '';
    return sqliteJSON(path, `SELECT effect_key,record_id,group_id,group_revision,bucket_start,state,retry_count,next_retry_at,payload,payload_digest
      FROM wg_delivery_outbox ${scope} ORDER BY bucket_start`);
  } catch (error) {
    if (/no such table/i.test(error.message)) return [];
    throw error;
  }
}

function sameTimestamp(left, right) {
  const leftMillis = Date.parse(String(left ?? ''));
  const rightMillis = Date.parse(String(right ?? ''));
  if (!Number.isFinite(leftMillis) || !Number.isFinite(rightMillis)) return String(left ?? '') === String(right ?? '');
  return leftMillis === rightMillis;
}

function frozenOutboxFacts(outbox) {
  let payload;
  try { payload = typeof outbox?.payload === 'string' ? JSON.parse(outbox.payload) : outbox?.payload; } catch { payload = null; }
  if (!payload || payload.record_id !== outbox.record_id || payload.effect_key !== outbox.effect_key || !sameTimestamp(payload.bucket_start, outbox.bucket_start)) {
    throw new Error('delayed outbox payload did not match its durable identity');
  }
  const provenanceCell = (payload.cells ?? []).find((cell) => cell.column === 'provenance');
  let provenance;
  try { provenance = typeof provenanceCell?.value === 'string' ? JSON.parse(provenanceCell.value) : provenanceCell?.value; } catch { provenance = null; }
  if (!Array.isArray(provenance) || provenance.length === 0) throw new Error('delayed outbox payload has no provenance evidence');
  return {
    effect_key: payload.effect_key,
    record_id: payload.record_id,
    group_id: outbox.group_id,
    group_revision: outbox.group_revision,
    bucket_start: outbox.bucket_start,
    bucket_start_raw: outbox.bucket_start,
    payload_bucket_start_raw: payload.bucket_start,
    payload_raw: outbox.payload,
    provenance,
  };
}

function operationEffects(path) {
  const operations = sqliteJSON(path, `SELECT operation_id,action,status,detail,updated_at
    FROM managed_schema_operations WHERE action = 'recording_start' ORDER BY updated_at`);
  const effects = [];
  for (const operation of operations) {
    let detail;
    try { detail = JSON.parse(operation.detail ?? '{}'); } catch { continue; }
    for (const group of detail.view?.groups ?? []) {
      if (group.applied && group.applied_revision) {
        effects.push({ operation_id: operation.operation_id, effect_key: `${group.group_id}:${group.applied_revision}`, group_id: group.group_id, group_revision: group.applied_revision });
      }
    }
  }
  return { operations, effects, outbox: safeOutbox(path, '') };
}

export async function createRun({ kind, runId, gatewayPort, simulatorPorts }) {
  assert.ok(['sqlite', 'postgres'].includes(kind), `unsupported fresh destination ${kind}`);
  const run = makeFreshRun(runId);
  const base = `http://127.0.0.1:${gatewayPort}`; const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
  const simPort = simulatorPorts[0];
  const hash = createHash('sha256').update(runId).digest('hex').slice(0, 20);
  let dsn;
  let target;
  let neighbor;
  let state;
  try {
    dsn = kind === 'postgres' ? parseDSN() : null;
    target = kind === 'sqlite'
      ? { kind, path: run.destination_db, table: CONNECTOR_TABLE }
      : { kind, host: '127.0.0.1', port: 55434, database: 'gwtest', username: 'postgres', password: dsn.password, schema: `gw_f_recovery_${hash}`, table: CONNECTOR_TABLE };
    neighbor = kind === 'sqlite'
      ? { kind, path: join(run.work, 'neighbor.db'), table: 'neighbor_rows' }
      : { kind, database: 'gwtest', schema: `gw_f_neighbor_${hash}`, table: 'neighbor_rows' };
    state = {
      ...run, kind, base, gatewayPort, simPort, target, neighbor, dsn, gateway: null,
      harness_started_at: new Date().toISOString(),
      gatewayBinary,
      simulator: null, browser: null, page: null, proxy: null, group: null,
      neighborReady: false, targetNeighborBaseline: null, cleaned: false, cleanupResult: null, activationFaultArmed: false, activationTrigger: `f_recovery_activation_${hash}`,
      targetSchemaCreated: false, neighborSchemaCreated: false,
      children: [], stopReports: [],
    };
    if (kind === 'postgres') {
      await assertFreshPortFree(55434);
      state.proxy = await makeOwnedPostgresProxy({ listenPort: 55434, upstreamPort: 55432 });
      postgresExec({ database: 'gwtest' }, `CREATE SCHEMA ${quoteIdentifier(target.schema)}`);
      state.targetSchemaCreated = true;
    }
  } catch (error) {
    await state?.proxy?.stop().catch(() => {});
    if (state?.targetSchemaCreated) {
      try { postgresExec({ database: 'gwtest' }, `DROP SCHEMA ${quoteIdentifier(target.schema)} CASCADE`); } catch { /* preserve the original initialization failure */ }
    }
    try { rmSync(run.work, { recursive: true, force: true }); } catch { /* preserve the original initialization failure */ }
    throw error;
  }

  const badStop = (status) => !status || status.term_state !== 'exited' || status.kill_sent || status.stop_error || status.spawn_error || status.kill_error || (status.exit_code !== null && status.exit_code !== 0);
  // cmd/test_ui bounds its runtime close at 5 seconds, so a harness timeout of the
  // same length races it whenever a fault fixture still holds the target; wait
  // longer than that bound. An exit that needed SIGKILL still fails the run.
  const OWNED_STOP_TIMEOUT_MS = 8000;
  const stopOwned = async (label, child) => {
    if (!child) return;
    const status = (await stopProcesses([child], OWNED_STOP_TIMEOUT_MS))[0];
    state.stopReports.push({ label, ...status });
    if (badStop(status)) throw new Error(`${label} did not exit cleanly`);
  };

  async function currentGroup() {
    const groups = await state.observeAPI({ method: 'GET', path: '/workspace/write-groups' });
    const group = (groups.groups ?? []).find((candidate) => candidate.id === state.group?.id);
    if (!group) throw new Error('fresh canonical managed group is no longer visible');
    state.group = group;
    return group;
  }

  state.observeAPI = async ({ method = 'GET', path }) => {
    assert.equal(method, 'GET', 'fresh adapter API observer is read-only');
    const requestPath = path.startsWith('/api/v1/datalink') ? path : `/api/v1/datalink/studio-v2${path.replace(/^\/studio-v2/, '')}`;
    return observeJSON(state.base, requestPath);
  };

  state.setupViaUI = async ({ intervalSeconds = 10 } = {}) => {
    await Promise.all([...new Set([gatewayPort, ...simulatorPorts])].map((port) => assertFreshPortFree(port)));
    if (kind === 'sqlite' && existsSync(target.path)) throw new Error('fresh SQLite target existed before UI setup');
    if (kind === 'postgres') {
      const tables = postgresJSON({ database: target.database }, `SELECT c.relname FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace WHERE n.nspname='${target.schema}'`);
      if (tables.length > 0) throw new Error('fresh PostgreSQL target schema was not empty');
    }
    state.simulator = startFreshSimulator({ work: state.work, name: 'sim-a', binary: DEFAULT_SIMULATOR_BINARY, port: simPort, registers: REGISTERS });
    state.children.push(state.simulator);
    state.gateway = startFreshGateway({ work: state.work, binary: gatewayBinary, port: gatewayPort });
    state.children.push(state.gateway);
    await waitForGateway(state.base);
    const ui = await openFreshUI({ port: gatewayPort, width: 1440, height: 1400 });
    state.browser = ui.browser; state.page = ui.page; state.ui = ui;
    await setupFreshDevice(state.page, { name: 'Line A', port: simPort });
    await setupFreshRules(state.page, RULES);
    await persistRecoveryMappings(state.page, state.base, { tagPrefix: `f_${runId}`, targetTypes: TARGET_TYPES, expectedCount: RULES.length });
    if (kind === 'sqlite') await configureFreshSQLite(state.page, state.base, { path: target.path, table: target.table });
    else await configureFreshPostgres(state.page, state.base, target);
    await state.page.getByTestId('basic-recording-interval').fill(String(intervalSeconds));
    const prepared = await prepareFreshBasicGroup(state.page, state.base);
    state.group = prepared.group;
    return { groupId: state.group.id, target, neighbor };
  };

  state.captureResponsiveEvidence = async (stage) => {
    if (!state.page) return null;
    const toggle = state.page.getByTitle(/^收合側邊欄/);
    const rail = state.page.getByTestId('sidebar-rail-aside');
    if ((await rail.getAttribute('data-collapsed')) !== 'true') {
      await toggle.focus();
      const activeTag = await state.page.evaluate(() => document.activeElement?.tagName ?? '');
      if (['INPUT', 'TEXTAREA', 'SELECT'].includes(activeTag)) throw new Error('sidebar keyboard capture retained form focus');
      await toggle.press('Control+b');
      await state.page.waitForTimeout(150);
    }
    if ((await rail.getAttribute('data-collapsed')) !== 'true') throw new Error('sidebar did not collapse with focused Ctrl+B');
    const screenshots = [];
    const focusIds = stage === 'partial' ? ['basic-recording-status', 'basic-recording-retry', 'basic-recording-panel'] : stage === 'error' ? ['group-schema-error', 'group-schema-operation-pending', 'group-schema-panel'] : ['database-delivery-truth', 'basic-recording-evidence', 'basic-recording-panel'];
    const controlsIds = stage === 'partial' ? ['basic-recording-retry', 'basic-recording-start', 'basic-recording-panel'] : stage === 'error' ? ['group-schema-check', 'group-schema-error', 'group-schema-panel'] : ['basic-recording-start', 'basic-recording-panel'];
    for (const [width, height] of [[390, 1000], [768, 1100], [1440, 1400]]) {
      await state.page.setViewportSize({ width, height });
      await state.page.waitForTimeout(250);
      const layout = await inspectBasicLayout(state.page); assertBasicLayout(layout, width);
      const file = `fresh-recovery-${kind}-${runId}-${stage}-${width}.png`;
      await state.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, file), fullPage: true });
      const focus = await screenshotAnchor(state.page, focusIds); const focusTestId = focus.testId;
      await focus.locator.scrollIntoViewIfNeeded();
      const focusFile = `fresh-recovery-${kind}-${runId}-${stage}-focus-${width}.png`;
      await state.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, focusFile) });
      const start = (await screenshotAnchor(state.page, controlsIds)).locator;
      await start.scrollIntoViewIfNeeded();
      const controlsFile = `fresh-recovery-${kind}-${runId}-${stage}-controls-${width}.png`;
      await state.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, controlsFile) });
      screenshots.push({ width, file, focus_file: focusFile, controls_file: controlsFile, focus_test_id: focusTestId, ...layout, document_horizontal_overflow: layout.document_width > width });
    }
    return { keyboard: { focused_control_b_collapsed: true }, screenshots };
  };
  state.openGroup = async (groupID) => {
    await state.page.getByTestId('step-nav-button-4').click();
    if (await state.page.getByTestId('group-editor').isVisible().catch(() => false)) return;
    await openFreshAdvanced(state.page);
    await state.page.getByTestId(`group-open-${groupID}`).click();
    await state.page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
  };

  state.targetSnapshot = async () => {
    const group = await currentGroup();
    const table = tableName(group);
    if (kind === 'sqlite') {
      if (!existsSync(target.path)) return { kind, exists: false, managed_table_count: 0, table, rows: [] };
      const stat = statSync(target.path);
      const sqliteMaster = sqliteJSON(target.path, "SELECT type,name,tbl_name FROM sqlite_master WHERE type IN ('table','index') ORDER BY name");
      const names = sqliteMaster.filter((candidate) => candidate.type === 'table' && candidate.name.startsWith('gw_group_'));
      const tableExists = Boolean(table) && names.some((candidate) => candidate.name === table);
      return {
        kind, exists: true, size_bytes: stat.size, mode: (stat.mode & 0o777).toString(8), sqlite_master: sqliteMaster,
        managed_table_count: names.length, table, table_exists: tableExists, rows: tableExists ? readSQLiteTable(target.path, table) : [],
      };
    }
    const names = postgresJSON({ database: target.database }, `SELECT c.relname AS name FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace WHERE n.nspname='${target.schema}' AND c.relname LIKE 'gw_group_%' AND c.relkind IN ('r','p')`);
    const tableExists = Boolean(table) && names.some((candidate) => candidate.name === table);
    return { kind, exists: true, managed_table_count: names.length, schema: target.schema, table, table_exists: tableExists, rows: tableExists ? readPostgresTable({ database: target.database }, target.schema, table) : [] };
  };

  state.targetReceiptSnapshot = async (effectKey) => kind === 'sqlite'
    ? readSQLiteReceipt(target.path, effectKey)
    : readPostgresReceipt({ database: target.database }, target.schema, effectKey);

  state.produceNeighborViaUI = async () => {
    if (state.neighborReady) return;
    // The managed destination starts empty and must never be pre-created by the
    // harness. Wait for a real runtime row so test-write cleanup proves that it
    // removes only its owned row from this exact target table.
    state.targetNeighborBaseline = await eventually(
      'managed target neighbor row',
      () => state.targetSnapshot(),
      (snapshot) => snapshot.table_exists && snapshot.rows.length > 0,
      90_000,
      1000,
    );
    if (kind === 'sqlite') {
      sqliteOwnedNeighborExec(neighbor.path, `CREATE TABLE IF NOT EXISTS neighbor_rows (fixture_id TEXT PRIMARY KEY, value INTEGER NOT NULL); DELETE FROM neighbor_rows; INSERT INTO neighbor_rows VALUES ('${runId}', 17)`);
    } else {
      postgresExec({ database: neighbor.database }, `CREATE SCHEMA ${quoteIdentifier(neighbor.schema)}; CREATE TABLE ${quoteIdentifier(neighbor.schema)}.${quoteIdentifier(neighbor.table)} (fixture_id TEXT PRIMARY KEY, value INTEGER NOT NULL); INSERT INTO ${quoteIdentifier(neighbor.schema)}.${quoteIdentifier(neighbor.table)} VALUES ('${runId}', 17)`);
      state.neighborSchemaCreated = true;
    }
    state.neighborReady = true;
  };

  state.neighborSnapshot = async () => {
    if (!state.neighborReady) return [];
    const rows = kind === 'sqlite'
      ? sqliteJSON(neighbor.path, `SELECT fixture_id,value FROM ${quoteIdentifier(neighbor.table)} ORDER BY fixture_id`)
      : postgresJSON({ database: neighbor.database }, `SELECT fixture_id,value FROM ${quoteIdentifier(neighbor.schema)}.${quoteIdentifier(neighbor.table)} ORDER BY fixture_id`);
    return rows;
  };

  state.gatewaySnapshot = async () => {
    const operations = sqliteJSON(state.gateway_db, `SELECT operation_id,action,status,detail,updated_at FROM managed_schema_operations ORDER BY updated_at`);
    const effects = operationEffects(state.gateway_db).effects;
    return { operations, effects, outbox: safeOutbox(state.gateway_db, state.group?.id ?? '') };
  };

  state.armStartPartial = async () => {
    if (state.activationFaultArmed) return;
    sqliteExec(state.gateway_db, `CREATE TRIGGER ${quoteIdentifier(state.activationTrigger)}
      BEFORE UPDATE OF detail ON managed_schema_operations
      WHEN NEW.action = 'recording_start' AND instr(NEW.detail, '"activated":true') > 0
      BEGIN SELECT RAISE(ABORT, 'fresh recovery activation checkpoint hold'); END`);
    state.activationFaultArmed = true;
  };

  state.releaseStartPartial = async () => {
    if (!state.activationFaultArmed) return;
    sqliteExec(state.gateway_db, `DROP TRIGGER ${quoteIdentifier(state.activationTrigger)}`);
    state.activationFaultArmed = false;
  };

  state.restart = async () => {
    if (state.gateway) await stopOwned('gateway restart', state.gateway);
    await assertFreshPortFree(gatewayPort);
    state.gateway = startFreshGateway({ work: state.work, binary: gatewayBinary, port: gatewayPort });
    state.children.push(state.gateway);
    await waitForGateway(state.base);
  };

  state.exerciseDelayedResend = async (groupID) => {
    const baselineRows = await eventually('first normal delivery', () => safeOutbox(state.gateway_db, groupID), (rows) => rows.some((row) => row.state === 'sql_committed'));
    const baselineCommitted = baselineRows.filter((row) => row.state === 'sql_committed');
    const baselineTime = Math.max(...baselineCommitted.map((row) => Date.parse(row.bucket_start)));
    let delayed;
    let sqliteLock;
    try {
      if (kind === 'sqlite') {
        sqliteLock = await holdSQLiteTargetWrites(target.path);
        state.children.push(sqliteLock.child);
      } else {
        await state.proxy.stop();
      }
      const delayedRows = await eventually('delayed delivery retry', () => safeOutbox(state.gateway_db, groupID), (rows) => rows.some((row) => Date.parse(row.bucket_start) > baselineTime && row.state !== 'sql_committed' && Number(row.retry_count) > 0), 60_000, 1000);
      delayed = delayedRows.find((row) => Date.parse(row.bucket_start) > baselineTime && row.state !== 'sql_committed' && Number(row.retry_count) > 0);
      if (!delayed) throw new Error('delayed delivery retry row disappeared before it could be frozen');
    } finally {
      if (kind === 'sqlite' && sqliteLock) await sqliteLock.close();
      if (kind === 'postgres') await state.proxy.start();
    }
    const frozen = frozenOutboxFacts(delayed);
    const committedRows = await eventually('delayed delivery commit', () => safeOutbox(state.gateway_db, groupID), (rows) => rows.some((row) => row.effect_key === frozen.effect_key && row.state === 'sql_committed'), 90_000, 1000);
    const committed = committedRows.find((row) => row.effect_key === frozen.effect_key && row.state === 'sql_committed');
    if (!committed) throw new Error('delayed outbox row disappeared before committed state could be observed');
    if (committed.record_id !== frozen.record_id || committed.payload_digest !== delayed.payload_digest) throw new Error('delayed outbox identity or payload digest changed during retry');
    const targetSnapshot = await eventually('delayed target row', () => state.targetSnapshot(), (snapshot) => snapshot.rows.some((row) => row.record_id === frozen.record_id), 30_000, 1000);
    const row = targetSnapshot.rows.find((candidate) => candidate.record_id === frozen.record_id);
    assert.ok(row, 'delayed target row was not observed');
    if (!sameTimestamp(row.bucket_start, frozen.bucket_start)) throw new Error('delayed resend changed original acquisition time');
    const receipt = await eventually('delayed destination receipt', () => state.targetReceiptSnapshot(frozen.effect_key), (rows) => rows.some((candidate) => candidate.effect_key === frozen.effect_key && candidate.payload_digest === delayed.payload_digest), 30_000, 1000);
    const current = await currentGroup();
    if (row.group_id !== committed.group_id || current.id !== committed.group_id || current.revision !== committed.group_revision) {
      throw new Error('delayed target row was not linked to the frozen group identity');
    }
    return {
      expected: { observed_at: frozen.bucket_start, provenance: frozen.provenance },
      durable_link: {
        record_id: committed.record_id,
        effect_key: committed.effect_key,
        group_id: committed.group_id,
        group_revision: committed.group_revision,
      },
      row,
      frozen_outbox: {
        effect_key: frozen.effect_key,
        record_id: frozen.record_id,
        group_id: frozen.group_id,
        group_revision: frozen.group_revision,
        bucket_start: frozen.bucket_start,
        bucket_start_raw: frozen.bucket_start_raw,
        payload_bucket_start_raw: frozen.payload_bucket_start_raw,
        payload_raw: frozen.payload_raw,
        provenance: frozen.provenance,
        payload_digest: delayed.payload_digest,
      },
      receipt,
      effects: [{ effect_key: frozen.effect_key, record_id: frozen.record_id, payload_digest: delayed.payload_digest }],
    };
  };

  state.cleanup = async () => {
    if (state.cleaned) return state.cleanupResult;
    const errors = [];
    let stopFailed = false;
    await closeFreshUI(state.browser).catch((error) => { stopFailed = true; errors.push(`browser: ${error.message}`); });
    for (const child of state.children.filter((candidate, index, all) => all.indexOf(candidate) === index)) {
      if (child.exitCode !== null || child.signalCode !== null) {
        const status = {
          name: child.fHarnessName,
          term_state: 'exited',
          kill_sent: false,
          exit_code: child.exitCode,
          signal: child.signalCode,
        };
        state.stopReports.push(status);
        if (badStop(status)) {
          stopFailed = true;
          errors.push(`${child.fHarnessName ?? 'owned child'} exited unsuccessfully`);
        }
        continue;
      }
      await stopOwned(child.fHarnessName ?? 'owned child', child).catch((error) => { stopFailed = true; errors.push(error.message); });
    }
    if (state.proxy) await state.proxy.stop().catch((error) => { stopFailed = true; errors.push(`proxy: ${error.message}`); });
    if (stopFailed) state.keepWork = true;
    if (state.activationFaultArmed) {
      try {
        sqliteExec(state.gateway_db, `DROP TRIGGER ${quoteIdentifier(state.activationTrigger)}`);
        state.activationFaultArmed = false;
      } catch (error) {
        errors.push(`activation trigger cleanup failed: ${error.message}`);
        state.keepWork = true;
      }
    }
    if (kind === 'postgres') {
      if (state.targetSchemaCreated) {
        try {
          postgresExec({ database: 'gwtest' }, `DROP SCHEMA ${quoteIdentifier(target.schema)} CASCADE`);
          if (postgresJSON({ database: 'gwtest' }, `SELECT nspname FROM pg_namespace WHERE nspname='${target.schema}'`).length !== 0) errors.push('target schema remained after cleanup');
          state.targetSchemaCreated = false;
        } catch (error) { errors.push('target schema cleanup failed'); }
      }
      if (state.neighborSchemaCreated) {
        try {
          postgresExec({ database: neighbor.database }, `DROP SCHEMA ${quoteIdentifier(neighbor.schema)} CASCADE`);
          if (postgresJSON({ database: neighbor.database }, `SELECT nspname FROM pg_namespace WHERE nspname='${neighbor.schema}'`).length !== 0) errors.push('neighbor schema remained after cleanup');
          state.neighborSchemaCreated = false;
        } catch (error) { errors.push('neighbor schema cleanup failed'); }
      }
    }
    let workAbsent = null;
    if (!state.keepWork) {
      try { rmSync(state.work, { recursive: true, force: true }); workAbsent = !existsSync(state.work); if (!workAbsent) errors.push('owned work directory remained after cleanup'); } catch { errors.push('owned work cleanup failed'); }
    } else {
      workAbsent = false;
    }
    state.cleanupResult = { passed: errors.length === 0, errors, stop_reports: state.stopReports, work_absent: workAbsent, work_retained: Boolean(state.keepWork) };
    state.cleaned = true;
    return state.cleanupResult;
  };
  return state;
}
