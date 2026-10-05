// F2/F3/F6 first real run: one fresh Modbus device, nine managed members,
// UI-only setup/DDL confirmation/Start, and independent SQLite observations.
import { join } from 'node:path';
import { existsSync, rmSync } from 'node:fs';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY,
  binaryIdentity, closeFreshUI, configureFreshSQLite,
  makeFreshRun, monotonicNs, observeJSON,
  openFreshUI, persistFreshMappings, previewAndApplyFreshSchema, preflightFreshSQLite,
  prepareFreshBasicGroup, setupFreshDevice, setupFreshRules, startFreshGateway, startFreshSimulator,
  startFreshBasic, stopProcesses, waitForFreshGateway, waitForObservation, assertFreshPortFree, isFreshCleanExit,
} from './fresh-ui.mjs';
import { sh } from './lib.mjs';
import {
  FRESH_EXPECTED_VALUES, FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES,
  expectedValuesByTag, verifyFreshDeliveryEvidence,
} from './fresh-recording-verification.mjs';
import { assertFreshTimingWitness } from './fresh-timing.mjs';
import { writeFreshResult } from './fresh-result.mjs';

const runId = process.env.F_RUN_ID ?? `single-sqlite-${Date.now()}`;
const port = Number(process.env.GW_PORT ?? 3380);
const simPort = Number(process.env.F_SIM_PORT_A ?? 15050);
const run = makeFreshRun(runId);
const base = `http://127.0.0.1:${port}`;
const target = run.destination_db;
const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
const evidencePath = join(DEFAULT_EVIDENCE_DIR, `${runId}.json`);
const processes = [];
const lines = [];
const log = (message) => { lines.push(message); console.log(message); };
const regs = [...FRESH_REGISTERS];
const rules = FRESH_RULES.map((rule) => ({ ...rule }));
const targetTypes = [...FRESH_TARGET_TYPES];
const targetTypeByAddress = new Map(FRESH_RULES.map((rule, index) => [String(rule.start), FRESH_TARGET_TYPES[index]]));
const expectedValues = [...FRESH_EXPECTED_VALUES];
let browser;
let page;
let ui;
let result;
let failure;

function jsonRows(path, sql) {
  const raw = sh('sqlite3', ['-readonly', '-json', '-cmd', '.timeout 15000', `file:${path}?mode=ro`, sql]);
  return raw ? JSON.parse(raw) : [];
}

function quoteSQL(value) {
  return `'${String(value).replaceAll("'", "''")}'`;
}

async function waitForClosedBucket(groupId, groupRevision) {
  const deadline = Date.now() + 270_000;
  let rows = [];
  let firstClosedAt;
  const query = `SELECT group_id,group_revision,bucket_start,kind,reason,record_id,members
    FROM wg_delivery_buckets WHERE group_id = ${quoteSQL(groupId)}
    AND group_revision = ${quoteSQL(groupRevision)} AND kind = 'row'
    ORDER BY bucket_start, entity_key`;
  while (Date.now() < deadline) {
    rows = jsonRows(run.gateway_db, query);
    if (rows.length > 0 && !firstClosedAt) firstClosedAt = monotonicNs();
    if (rows.length >= 3) return { bucket: rows[0], rows, t_closed: firstClosedAt, query };
    await new Promise((resolveWait) => setTimeout(resolveWait, 2000));
  }
  throw new Error(`gateway did not close three complete buckets; got ${rows.length}`);
}

async function waitForManagedRows(table, columns, minimum = 3) {
  const names = columns.map((column) => column.name);
  const select = names.map((name) => {
    const escaped = name.replaceAll('"', '""');
    return `CASE WHEN "${escaped}" IS NULL THEN NULL ELSE printf('%s', "${escaped}") END AS "${escaped}"`;
  }).join(', ');
  const query = `SELECT ${select} FROM "${table.replaceAll('"', '""')}" ORDER BY bucket_start, record_id`;
  const polls = [];
  const deadline = Date.now() + 270_000;
  let rows = [];
  let firstRowAt;
  while (Date.now() < deadline) {
    const opened = monotonicNs();
    rows = jsonRows(target, query);
    const completed = monotonicNs();
    polls.push({ t_poll_open: opened, t_poll_complete: completed, row_count: rows.length, interval_ms: 2000 });
    if (rows.length > 0 && !firstRowAt) firstRowAt = completed;
    if (rows.length >= minimum) return { rows, polls, query, t_sql: firstRowAt };
    await new Promise((resolveWait) => setTimeout(resolveWait, 2000));
  }
  throw new Error(`managed table did not reach ${minimum} rows; got ${rows.length}`);
}

async function captureWidths() {
  for (const [width, height] of [[390, 1000], [768, 1100], [1440, 1400]]) {
    await page.setViewportSize({ width, height });
    await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, `${runId}-${width}.png`), fullPage: true });
  }
}

try {
  const initial = preflightFreshSQLite(target, run.work);
  if (initial.recording_table_present !== false || existsSync(target)) throw new Error('fresh SQLite target was present before UI');
  if (port === simPort) throw new Error('gateway and simulator ports must be distinct');
  await Promise.all([assertFreshPortFree(port), assertFreshPortFree(simPort)]);
  const sim = startFreshSimulator({ work: run.work, port: simPort, registers: Object.fromEntries(regs.map((value, index) => [index, value])) });
  const gateway = startFreshGateway({ work: run.work, binary: gatewayBinary, port });
  processes.push(sim, gateway);
  await waitForFreshGateway(base);
  (ui = await openFreshUI({ port, width: 1440, height: 1400 }));
  ({ browser, page } = ui);
  await setupFreshDevice(page, { name: 'Line A', port: simPort });
  await setupFreshRules(page, rules);
  const tags = await persistFreshMappings(page, base, {
    tagPrefix: 'line_a', targetTypes,
    targetTypeFor: ({ address }) => targetTypeByAddress.get(address),
  });
  const destination = await configureFreshSQLite(page, base, { path: target, table: 'f_connector_placeholder' });
  const afterConnector = preflightFreshSQLite(target, run.work);
  const firstStart = await prepareFreshBasicGroup(page, base);
  const groupDraft = firstStart.group;
  const beforePreview = preflightFreshSQLite(target, run.work);
  if (beforePreview.recording_table_present !== false) throw new Error('schema preview exposed a recording table');
  const schema = await previewAndApplyFreshSchema(page);
  const afterApply = preflightFreshSQLite(target, run.work);
  if (afterApply.recording_table_present !== false || !existsSync(target)) throw new Error('schema confirmation did not create SQLite target');
  if (await page.getByTestId('group-editor-close').count()) await page.getByTestId('group-editor-close').click();
  const start = await startFreshBasic(page);
  const tableRows = jsonRows(target, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name");
  const tables = tableRows.map((row) => row.name).filter((name) => /^gw_group_[A-Za-z0-9_]+$/.test(name));
  if (tables.length !== 1) throw new Error(`expected one managed table, found ${JSON.stringify(tables)}`);
  const table = tables[0];
  const columns = jsonRows(target, `PRAGMA table_info("${table.replaceAll('"', '""')}")`);
  const groups = await waitForObservation(base, '/studio-v2/workspace/write-groups', (data) => {
    const matches = (data.groups ?? []).filter((candidate) => candidate.id === groupDraft.id);
    return matches.length === 1 && Boolean(matches[0].applied_revision || matches[0].revision);
  }, 30_000);
  const groupMatches = groups.groups?.filter((candidate) => candidate.id === groupDraft.id) ?? [];
  if (groupMatches.length !== 1) throw new Error(`expected one exact current Basic group, found ${groupMatches.length}`);
  const group = groupMatches[0];
  const groupRevision = group.applied_revision || group.revision;
  const [closed, first] = await Promise.all([
    waitForClosedBucket(group.id, groupRevision),
    waitForManagedRows(table, columns, 3),
  ]);
  const timing = assertFreshTimingWitness({
    t_open: ui.t_open,
    t_start: start.t_start,
    t_closed: closed.t_closed,
    t_sql: first.t_sql,
    polls: first.polls,
  });
  const mappingData = await observeJSON(base, '/studio-v2/workspace/mappings');
  const mappings = Array.isArray(mappingData) ? mappingData : mappingData.mappings ?? [];
  const memberColumns = (group.members ?? []).map((member) => ({
    ...member,
    mapping: mappings.find((mapping) => (mapping.point_id ?? mapping.persisted_point_id) === member.point_id && mapping.tag_id === member.tag_id),
  }));
  const valuesByMember = memberColumns.map((member) => ({
    point_id: member.point_id, tag_id: member.tag_id, target_column: member.target_column,
    tag_key: member.mapping?.tag_key, target_type: member.mapping?.target_type,
    values: first.rows.map((row) => row[member.target_column]),
  }));
  const expectedByTag = expectedValuesByTag(tags, expectedValues);
  const observed = valuesByMember.map((member) => ({
    tag_key: member.tag_key, value: member.values[0], expected: expectedByTag[member.tag_key],
  }));
  const outbox = jsonRows(run.gateway_db, `SELECT effect_key,record_id,group_id,group_revision,bucket_start,table_name,state,payload_digest
    FROM wg_delivery_outbox WHERE group_id = ${quoteSQL(group.id)} ORDER BY bucket_start,record_id`);
  const receipts = jsonRows(run.gateway_db, `SELECT effect_key,payload_digest,committed_at FROM wg_delivery_receipts
    WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE group_id = ${quoteSQL(group.id)})`);
  const closedRows = jsonRows(run.gateway_db, `SELECT group_revision,bucket_start,kind,record_id,members
    FROM wg_delivery_buckets WHERE group_id = ${quoteSQL(group.id)} AND group_revision = ${quoteSQL(groupRevision)} AND kind = 'row'
    ORDER BY bucket_start,entity_key`);
  const failures = [];
  if (afterConnector.exists) failures.push('SQLite target was created by connector save before schema preview');
  if ((group.members ?? []).length !== 9) failures.push(`managed member count ${group.members?.length}`);
  if (observed.length !== expectedValues.length || observed.some((item) => String(item.value) !== item.expected)) {
    failures.push(`first bucket values differ: ${JSON.stringify(observed)}`);
  }
  if (first.rows.length < 3) failures.push(`only ${first.rows.length} buckets`);
  if (first.rows.some((row) => !row.bucket_start || !row.provenance || !row.device_id || !row.group_id)) failures.push('row identity/provenance fields incomplete');
  const delivery = verifyFreshDeliveryEvidence({
    group, table, rows: first.rows, mappings, expectedByTag, closedBuckets: closedRows, outbox, receipts,
  });
  failures.push(...delivery.failures);
  await captureWidths();
  result = {
    run_id: runId, passed: failures.length === 0, failures, destination, table, ...timing,
    target_preflight: {
      before_ui: initial, after_connector_save: afterConnector, before_preview: beforePreview, after_apply: afterApply,
      strict_missing_until_preview: !afterConnector.exists && !beforePreview.exists,
    },
    binary: { gateway: binaryIdentity(gatewayBinary), simulator: binaryIdentity(DEFAULT_SIMULATOR_BINARY) },
    registers: regs, rules, tags, group, member_columns: memberColumns, values_by_member: valuesByMember,
    columns, rows: first.rows, sql_query: first.query, polls: first.polls, delivery,
    clocks: { t_open: ui.t_open, t_start: start.t_start, t_closed: closed.t_closed, t_sql: first.t_sql, t_basic_prepare_start: firstStart.t_start },
    timing,
    ui_status: { prepare: firstStart.first_status_text, committed: start.status_text }, schema, page_errors: ui.pageErrors, api_traffic: ui.apiTraffic, log: lines,
    retained_work: run.work, screenshots: [390, 768, 1440].map((width) => `${runId}-${width}.png`),
  };
  log(result.passed ? `PASS ${runId}: ${first.rows.length} SQLite buckets` : `FAIL ${runId}: ${failures.join('; ')}`);
  process.exitCode = result.passed ? 0 : 1;
} catch (error) {
  failure = error;
  process.exitCode = 1;
  log(`ERROR ${runId}: ${error?.message?.split('\n')[0] ?? String(error)}`);
} finally {
  const cleanup = await stopProcesses(processes).catch((error) => [{ state: 'error', error: error.message }]);
  let browserCleanup;
  try {
    await closeFreshUI(browser);
    browserCleanup = { state: 'closed' };
  } catch (error) {
    browserCleanup = { state: 'error', error: error.message.split('\n')[0] };
  }
  const payload = result ?? { run_id: runId, passed: false, failures: [failure?.message ?? 'fresh run failed'], retained_work: run.work, log: lines };
  const processCleanupOk = cleanup.every(isFreshCleanExit);
  const namespaceCleanup = { state: 'pending', path: run.work, absent: false };
  if (processCleanupOk && browserCleanup.state === 'closed') {
    try {
      rmSync(run.work, { recursive: true, force: false });
      namespaceCleanup.state = 'removed';
      namespaceCleanup.absent = !existsSync(run.work);
    } catch (error) {
      namespaceCleanup.state = 'error';
      namespaceCleanup.error = error.message.split('\n')[0];
    }
  } else {
    namespaceCleanup.state = 'blocked';
  }
  payload.cleanup = { processes: cleanup, browser: browserCleanup, namespace: namespaceCleanup };
  if (!processCleanupOk || browserCleanup.state !== 'closed' || namespaceCleanup.absent !== true) {
    payload.passed = false;
    payload.failures = [...new Set([...(payload.failures ?? []), 'cleanup did not prove all children/browser/namespace closed'])];
  }
  writeFreshResult(evidencePath, payload);
}

if (failure) throw failure;
