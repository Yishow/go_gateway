// Fresh SQLite third matrix: one UI-created A/B entity group with a shared
// value column and distinct member columns, followed by explicit test-write.
// Browser actions perform setup mutations; API/SQL below are observations.
import assert from 'node:assert/strict';
import { existsSync, rmSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, assertFreshPortFree,
  binaryIdentity, closeFreshUI, configureFreshSQLite, ensureFreshSidebarCollapsed, isFreshCleanExit, makeFreshRun,
  monotonicNs, observeJSON, openFreshAdvanced, openFreshUI, persistFreshMappings,
  prepareFreshBasicGroup, previewAndApplyFreshSchema, setupFreshMultiDeviceSources,
  startFreshBasic, startFreshGateway, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import {
  FRESH_EXPECTED_VALUES, FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES,
  verifyFreshDeliveryEvidence,
} from './fresh-recording-verification.mjs';
import { preflightSQLiteDestination } from './destination-preflight.mjs';
import { readSQLiteDelivery, readSQLiteRows } from './fresh-sqlite-observation.mjs';
import { observePreservedRuntimeRows } from './fresh-matrix-test-write.mjs';
import { assertFreshTimingWitness } from './fresh-timing.mjs';

const runId = process.env.F_RUN_ID ?? `multi-sqlite-${Date.now()}`;
const port = Number(process.env.GW_PORT ?? 3380);
const simPortA = Number(process.env.F_SIM_PORT_A ?? 15050);
const simPortB = Number(process.env.F_SIM_PORT_B ?? 15051);
const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
const run = makeFreshRun(runId);
const base = `http://127.0.0.1:${port}`;
const target = run.destination_db;
const evidencePath = join(DEFAULT_EVIDENCE_DIR, `${runId}.json`);
const children = [];
let browser; let page; let ui; let result; let failure; let stage = 'preflight';

const valuesByAddress = (values) => Object.fromEntries(FRESH_RULES.map((rule, i) => [String(rule.start), values[i]]));
const valuesA = valuesByAddress(FRESH_EXPECTED_VALUES);
const registersB = [...FRESH_REGISTERS];
registersB.splice(0, 14, 187, 777, 0, 0, 0, 2, 0, 0xc010, 0, 0xc084, 0x7291, 0x6872, 0xb021, 0x4232);
// The string rule reads register index 22 (40023), independently of the
// float words above; 0x4232 is the simulator's B2 literal.
registersB[22] = 0x4232;
const valuesB = valuesByAddress(['187', '777', '2', '0', '-2.25', '-654.321', FRESH_EXPECTED_VALUES[6], FRESH_EXPECTED_VALUES[7], 'B2']);
const log = (message) => console.log(message);
const qid = (value) => `"${String(value).replaceAll('"', '""')}"`;
const qlit = (value) => `'${String(value).replaceAll("'", "''")}'`;

function tableColumns(table) {
  return readSQLiteRows(target, `PRAGMA table_info(${qid(table)})`).map((column) => ({ name: column.name }));
}

function tableSelect(columns) {
  return columns.map(({ name }) => {
    const quoted = qid(name);
    return `CASE WHEN ${quoted} IS NULL THEN NULL ELSE printf('%s', ${quoted}) END AS ${quoted}`;
  }).join(', ');
}

async function groupsAt() { return observeJSON(base, '/studio-v2/workspace/write-groups'); }

async function waitGroup(id, predicate, timeoutMs = 30_000) {
  const deadline = Date.now() + timeoutMs; let current;
  while (Date.now() < deadline) {
    current = (await groupsAt()).groups?.find((group) => group.id === id);
    if (current && predicate(current)) return current;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error(`group ${id} did not reach expected state: ${JSON.stringify(current)}`);
}

async function openGroup(id) {
  const advanced = page.getByTestId('step4-advanced-recording');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId(`group-open-${id}`).click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
}

async function waitBasicGroup(deviceId) {
  const deadline = Date.now() + 30_000; let group;
  while (Date.now() < deadline) {
    group = (await groupsAt()).groups?.find((item) => item.basic_managed_device_id === deviceId);
    if (group) return group;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error(`Basic group not persisted for ${deviceId}`);
}

async function ensureBasic(deviceId) {
  const prepared = await prepareFreshBasicGroup(page, base, deviceId);
  const schema = await previewAndApplyFreshSchema(page);
  await page.getByTestId('group-editor-close').click();
  const started = await startFreshBasic(page);
  return { prepared, schema, started, group: await waitBasicGroup(deviceId) };
}

async function disableGroup(group) {
  await openGroup(group.id);
  const current = (await groupsAt()).groups?.find((item) => item.id === group.id);
  if (current?.status !== 'disabled') {
    await page.getByTestId('group-disable').click();
    await waitGroup(group.id, (item) => item.status === 'disabled');
  }
  await page.getByTestId('group-editor-close').click();
}

async function createEntityGroup(tagsA, tagsB) {
  await openFreshAdvanced(page);
  await page.getByTestId('group-create').click();
  await page.getByTestId('group-name').fill('F SQLite entity A+B');
  await page.getByTestId('group-storage-managed').check();
  const advanced = page.getByTestId('group-advanced');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-entity-column').fill('entity_id');
  await page.getByTestId('group-provenance-column').fill('provenance');
  await page.getByTestId('group-bulk-include').click();
  for (const tag of [...tagsA, ...tagsB]) {
    const row = page.getByTestId(`group-member-row-${tag.persisted_point_id}`);
    await row.locator('input[type="text"]').first().fill(tag.device_id === tagsA[0].device_id ? 'A' : 'B');
  }
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  let groups = await groupsAt();
  const matches = groups.groups?.filter((group) => group.name === 'F SQLite entity A+B') ?? [];
  assert.equal(matches.length, 1);
  const managed = matches[0];
  assert.equal(managed.members.length, 18);
  const testWriteDistinctColumn = managed.members.find((member) => member.device_id === tagsB[0].device_id && member.point_id === tagsB[0].persisted_point_id)?.target_column;
  assert.ok(testWriteDistinctColumn);
  const schema = await previewAndApplyFreshSchema(page);
  await page.getByTestId('group-storage-custom').check();
  const aColumn = page.getByTestId(`group-member-column-${tagsA[0].persisted_point_id}`);
  const bColumn = page.getByTestId(`group-member-column-${tagsB[0].persisted_point_id}`);
  await aColumn.waitFor({ timeout: 20_000 });
  const sharedColumn = await aColumn.inputValue();
  assert.ok(sharedColumn);
  await bColumn.selectOption(sharedColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  groups = await groupsAt();
  const saved = groups.groups?.find((group) => group.id === managed.id);
  assert.ok(saved);
  assert.equal(saved.destination.storage_strategy, 'custom');
  assert.equal(saved.row_policy.entity_key_column, 'entity_id');
  const shared = saved.members.filter((member) => member.target_column === sharedColumn);
  assert.equal(shared.length, 2, 'exactly the A/B first members must share one column');
  assert.ok(saved.members.some((member) => member.target_column !== sharedColumn), 'remaining members must stay distinct');
  return { group: saved, schema, sharedColumn, testWriteDistinctColumn };
}

function consecutive(rows) {
  const starts = [...new Set(rows.map((row) => row.bucket_start))].sort();
  const times = starts.map((value) => Date.parse(value));
  return starts.length === 3 && times.every(Number.isFinite) && times[1] - times[0] === 60_000 && times[2] - times[1] === 60_000;
}

async function waitMultiRows(group, table, columns) {
  const select = tableSelect(columns);
  const deadline = Date.now() + 270_000; const polls = [];
  let rows = []; let closed = []; let tClosed; let tSql;
  while (Date.now() < deadline) {
    const before = monotonicNs();
    rows = readSQLiteRows(target, `SELECT ${select} FROM ${qid(table)} ORDER BY bucket_start, record_id`);
    const gatewayRows = readSQLiteRows(run.gateway_db, `SELECT group_id,group_revision,bucket_start,kind,reason,record_id,members,entity_key FROM wg_delivery_buckets WHERE group_id=${qlit(group.id)} AND group_revision=${qlit(group.applied_revision)} AND kind='row' ORDER BY bucket_start,entity_key`);
    const delivery = readSQLiteDelivery({ path: run.gateway_db, groupId: group.id, groupRevision: group.applied_revision });
    if (gatewayRows.length && !tClosed) tClosed = before;
    if (rows.length && !tSql) tSql = monotonicNs();
    const changed = polls.at(-1)?.rows !== rows.length || polls.at(-1)?.closed !== gatewayRows.length;
    if (changed) polls.push({ rows: rows.length, closed: gatewayRows.length, observed_ns: monotonicNs() });
    const committedEffects = delivery.outbox.filter((item) => item.state === 'sql_committed');
    const receiptKeys = new Set(delivery.receipts.map((item) => item.effect_key));
    const durableComplete = committedEffects.length >= 6 && committedEffects.every((item) => receiptKeys.has(item.effect_key));
    if (rows.length >= 6 && gatewayRows.length >= 6 && durableComplete && ['A', 'B'].every((entity) => consecutive(rows.filter((row) => row.entity_id === entity)))) {
      closed = gatewayRows; return { rows, closed, delivery, polls, t_closed: tClosed, t_sql: tSql };
    }
    await new Promise((resolve) => setTimeout(resolve, 1_000));
  }
  throw new Error(`multi-entity group did not produce six rows/three buckets: rows=${rows.length} closed=${closed.length}`);
}

async function prepareTestWriteGroup(entity, tagsB, appliedRevision) {
  // The running revision deliberately keeps the shared A/B column.  The
  // test-write contract needs a unique member value for its safe readback, so
  // save only the B int16 draft column and leave the applied revision alone.
  const selector = page.getByTestId(`group-member-column-${tagsB[0].persisted_point_id}`);
  await selector.selectOption(entity.testWriteDistinctColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const current = (await groupsAt()).groups?.find((group) => group.id === entity.group.id);
  assert.ok(current);
  assert.notEqual(current.revision, entity.group.revision, 'test-write draft change did not create a saved revision');
  assert.equal(current.applied_revision, appliedRevision, 'test-write draft unexpectedly re-applied the shared revision');
  assert.equal(current.members.find((member) => member.point_id === tagsB[0].persisted_point_id)?.target_column, entity.testWriteDistinctColumn);
  await page.getByTestId('group-editor-close').click();
  return { group: current, draft_saved_without_apply: true };
}

async function runTestWrite(group, table, columns) {
  await openGroup(group.id);
  const select = tableSelect(columns);
  const before = readSQLiteRows(target, `SELECT ${select} FROM ${qid(table)} ORDER BY bucket_start, record_id`);
  const previewResponse = await Promise.all([
    page.waitForResponse((response) => response.url().endsWith(`/write-groups/${group.id}/test-write-preview`) && response.request().method() === 'POST'),
    page.getByTestId('group-test-write-preview').click(),
  ]).then(([response]) => response);
  assert.equal(previewResponse.status(), 200);
  const preview = (await previewResponse.json()).data;
  assert.ok(preview?.operation_id && preview?.token);
  assert.equal(preview.group_id, group.id);
  assert.equal(preview.group_revision, group.revision);
  assert.equal(preview.target.table, group.destination.table_name);
  assert.equal(preview.owner_column, group.row_policy.entity_key_column);
  assert.match(preview.owner_value, /^gw-test-[a-f0-9-]+$/);
  try {
    await page.getByTestId('group-test-write-preview-card').waitFor({ timeout: 20_000 });
  } catch (error) {
    const uiError = await page.getByTestId('group-test-write-error').innerText().catch(() => '');
    const panelText = await page.getByTestId('group-test-write').innerText().catch(() => '');
    throw new Error(`${error.message.split('\n')[0]}; preview-card-error=${uiError || '<none>'}; panel=${panelText.replace(/\s+/g, ' ').slice(-600)}`);
  }
  const afterPreview = readSQLiteRows(target, `SELECT ${select} FROM ${qid(table)} ORDER BY bucket_start, record_id`);
  const afterPreviewByRecord = new Map(afterPreview.map((row) => [row.record_id, row]));
  for (const row of before) assert.deepEqual(afterPreviewByRecord.get(row.record_id), row, `preview changed runtime row ${row.record_id}`);
  const confirmResponse = await Promise.all([
    page.waitForResponse((response) => response.url().endsWith(`/write-groups/${group.id}/test-write`) && response.request().method() === 'POST'),
    page.getByTestId('group-test-write-confirm').click(),
  ]).then(([response]) => response);
  assert.equal(confirmResponse.status(), 200);
  const operation = (await confirmResponse.json()).data;
  assert.equal(operation.operation_id, preview.operation_id);
  assert.equal(operation.action, 'test_write');
  await page.getByTestId('group-test-write-outcome').waitFor({ timeout: 30_000 });
  assert.equal(operation.write_outcome, 'written_verified');
  assert.equal(operation.cleanup_status, 'cleaned');
  const after = readSQLiteRows(target, `SELECT ${select} FROM ${qid(table)} ORDER BY bucket_start, record_id`);
  const afterByRecord = new Map(after.map((row) => [row.record_id, row]));
  for (const row of before) assert.deepEqual(afterByRecord.get(row.record_id), row, `runtime row ${row.record_id} changed during test write`);
  assert.equal(after.filter((row) => String(row[preview.owner_column] ?? '') === preview.owner_value).length, 0, 'operation-owned test row remains');
  const ledgerRows = readSQLiteRows(run.gateway_db, `SELECT operation_id,status,action,scope_key,write_outcome,cleanup_status,payload_digest,detail FROM managed_schema_operations WHERE operation_id=${qlit(operation.operation_id)}`);
  assert.equal(ledgerRows.length, 1);
  const ledger = ledgerRows[0];
  assert.equal(ledger.operation_id, preview.operation_id);
  assert.equal(ledger.action, 'test_write');
  assert.ok(ledger.scope_key, 'test-write ledger scope is missing');
  const detail = JSON.parse(ledger.detail);
  assert.equal(detail.owner_column, preview.owner_column);
  assert.equal(detail.owner_value, preview.owner_value);
  assert.equal(detail.table, preview.target.table);
  const { after: durableAfter, delivery: runtimeNeighborDelivery } = await observePreservedRuntimeRows({
    before,
    group,
    table,
    readRows: () => readSQLiteRows(target, `SELECT ${select} FROM ${qid(table)} ORDER BY bucket_start, record_id`),
    readDelivery: () => readSQLiteDelivery({ path: run.gateway_db, groupId: group.id, groupRevision: group.applied_revision }),
  });
  assert.equal(durableAfter.filter((row) => String(row[preview.owner_column] ?? '') === preview.owner_value).length, 0, 'operation-owned test row remains after durable observation');
  return { preview: { operation_id: preview.operation_id, group_id: preview.group_id, owner_column: preview.owner_column, owner_value: preview.owner_value, table: preview.target.table }, operation, ledger, before, after: durableAfter, runtime_neighbor_delivery: runtimeNeighborDelivery };
}

async function captureWidths() {
  await page.getByTestId('group-editor-close').click();
  const advanced = page.getByTestId('step4-advanced-recording');
  if ((await advanced.getAttribute('open')) !== null) await advanced.locator('summary').click();
  for (const button of await page.getByRole('button', { name: /Collapse|收合/i }).all()) {
    if (await button.isVisible().catch(() => false)) await button.click();
  }
  await ensureFreshSidebarCollapsed(page);
  const screenshots = [];
  for (const [width, height] of [[390, 1100], [768, 1200], [1440, 1400]]) {
    await page.setViewportSize({ width, height }); await page.waitForTimeout(300);
    const file = `${runId}-${width}.png`;
    screenshots.push({ width, file, overflow: await page.evaluate(() => document.documentElement.scrollWidth > innerWidth),
      basic_visible: await page.getByTestId('basic-recording-panel').isVisible() });
    await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, file), fullPage: true });
  }
  return screenshots;
}

try {
  assert.ok(new Set([port, simPortA, simPortB]).size === 3, 'fresh ports must be distinct');
  const initial = preflightSQLiteDestination({ path: target, work: run.work });
  assert.equal(initial.recording_table_present, false); assert.equal(existsSync(target), false);
  await Promise.all([assertFreshPortFree(port), assertFreshPortFree(simPortA), assertFreshPortFree(simPortB)]);
  children.push(startFreshSimulator({ work: run.work, name: 'sim-a', port: simPortA, registers: Object.fromEntries(FRESH_REGISTERS.map((value, i) => [i, value])) }));
  children.push(startFreshSimulator({ work: run.work, name: 'sim-b', port: simPortB, registers: Object.fromEntries(registersB.map((value, i) => [i, value])) }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port }));
  await waitForFreshGateway(base); ui = await openFreshUI({ port }); ({ browser, page } = ui);
  stage = 'devices-and-mappings';
  await setupFreshMultiDeviceSources(page, [{ name: 'Line A', port: simPortA, prefix: 'A' }, { name: 'Line B', port: simPortB, prefix: 'B' }], FRESH_RULES);
  const byAddress = new Map(FRESH_RULES.map((rule, i) => [String(rule.start), FRESH_TARGET_TYPES[i]]));
  const tags = await persistFreshMappings(page, base, { tagPrefix: 'entity', targetTypes: [...FRESH_TARGET_TYPES, ...FRESH_TARGET_TYPES], targetTypeFor: ({ address }) => byAddress.get(address) });
  const deviceData = await observeJSON(base, '/studio-v2/workspace/devices');
  const savedDevices = Array.isArray(deviceData) ? deviceData : deviceData.devices ?? [];
  const deviceA = savedDevices.find((device) => device.name === 'Line A'); const deviceB = savedDevices.find((device) => device.name === 'Line B');
  assert.ok(deviceA?.id && deviceB?.id && deviceA.id !== deviceB.id);
  const tagsA = tags.filter((tag) => tag.device_id === deviceA.id); const tagsB = tags.filter((tag) => tag.device_id === deviceB.id);
  assert.equal(tagsA.length, 9); assert.equal(tagsB.length, 9);
  stage = 'destination';
  const destination = await configureFreshSQLite(page, base, { path: target, table: 'f_connector_placeholder' });
  assert.equal(preflightSQLiteDestination({ path: target, work: run.work }).exists, false);
  stage = 'basic-prerequisites';
  const basicA = await ensureBasic(deviceA.id); const basicB = await ensureBasic(deviceB.id);
  await disableGroup(basicA.group); await disableGroup(basicB.group);
  stage = 'same-group-save-schema-apply';
  const entity = await createEntityGroup(tagsA, tagsB);
  await page.getByTestId('group-editor-close').click();
  const entityBeforeApply = (await groupsAt()).groups.find((group) => group.id === entity.group.id);
  await openGroup(entity.group.id);
  await page.getByTestId('group-lifecycle').waitFor({ timeout: 20_000 });
  await page.waitForFunction(() => !document.querySelector('[data-testid="group-apply"]')?.hasAttribute('disabled'), null, { timeout: 30_000 });
  const tStart = monotonicNs(); await page.getByTestId('group-apply').click();
  const applied = await waitGroup(entity.group.id, (group) => Boolean(group.applied_revision) && group.status !== 'draft');
  const groupRevision = applied.applied_revision; const table = applied.destination.table_name; const columns = tableColumns(table);
  assert.ok(columns.some((column) => column.name === 'entity_id'));
  stage = 'three-buckets';
  const expectedByTag = Object.fromEntries(tags.map((tag) => [tag.tag_key, (tag.device_id === deviceA.id ? valuesA : valuesB)[String(tag.address)]]));
  const observed = await waitMultiRows(applied, table, columns);
  const mappingsData = await observeJSON(base, '/studio-v2/workspace/mappings');
  const mappings = Array.isArray(mappingsData) ? mappingsData : mappingsData.mappings ?? [];
  const delivery = readSQLiteDelivery({ path: run.gateway_db, groupId: applied.id, groupRevision });
  const verification = verifyFreshDeliveryEvidence({ group: applied, table, rows: observed.rows, mappings, expectedByTag,
    closedBuckets: delivery.closedBuckets, outbox: delivery.outbox, receipts: delivery.receipts,
    sharedTargetColumns: [entity.sharedColumn], entityColumn: 'entity_id' });
  assert.deepEqual(verification.failures, []);
  assert.equal(observed.rows.length, 6); assert.equal(delivery.closedBuckets.length, 6);
  stage = 'explicit-test-write';
  const savedForTestWrite = await prepareTestWriteGroup(entity, tagsB, groupRevision);
  const testWriteEvidence = await runTestWrite(savedForTestWrite.group, table, columns);
  const screenshots = await captureWidths();
  const timings = assertFreshTimingWitness({ t_open: ui.t_open, t_start: tStart, t_closed: observed.t_closed, t_sql: observed.t_sql, polls: observed.polls });
  result = { run_id: runId, passed: true, kind: 'sqlite-single-group-multi-entity', destination, target_preflight: { before_ui: initial },
    binary: { gateway: binaryIdentity(gatewayBinary), simulator: binaryIdentity(DEFAULT_SIMULATOR_BINARY) },
    source_registers: { B_string_rule: { address: 40023, register_index: 22, word: '0x4232', value: 'B2' } }, devices: [deviceA, deviceB],
    tags, groups: { basic: [basicA.group, basicB.group], entity_before_apply: entityBeforeApply, applied },
    shared_column: entity.sharedColumn, columns, expected_by_tag: expectedByTag, rows: observed.rows, delivery, verification,
    test_write: { ...testWriteEvidence, saved_group: savedForTestWrite.group, applied_revision_before_save: groupRevision }, timings,
    screenshots, page_errors: ui.pageErrors, api_traffic: ui.apiTraffic, retained_work: run.work };
  log(`PASS ${runId}: one A/B entity group, three buckets each, explicit test-write verified and cleaned`);
} catch (error) {
  failure = error; result = { run_id: runId, passed: false, failed_stage: stage, error: String(error?.message ?? error).split('\n')[0], retained_work: run.work };
  log(`ERROR ${runId} ${stage}: ${result.error}`);
} finally {
  const cleanup = await stopProcesses(children).catch((error) => [{ state: 'error', error: String(error) }]);
  let browserCleanup = { state: 'closed' }; try { await closeFreshUI(browser); } catch (error) { browserCleanup = { state: 'error', error: error.message.split('\n')[0] }; }
  const processOK = cleanup.every(isFreshCleanExit);
  const payload = result ?? { run_id: runId, passed: false, error: 'no result', retained_work: run.work };
  const namespace = { path: run.work, absent: false, state: 'blocked' };
  if (processOK && browserCleanup.state === 'closed') { try { rmSync(run.work, { recursive: true, force: false }); namespace.state = 'removed'; namespace.absent = !existsSync(run.work); } catch (error) { namespace.error = error.message.split('\n')[0]; } }
  payload.cleanup = { processes: cleanup, browser: browserCleanup, namespace };
  if (!processOK || browserCleanup.state !== 'closed' || !namespace.absent) { payload.passed = false; payload.error = `${payload.error ?? ''}; cleanup not proven`; }
  writeFileSync(evidencePath, `${JSON.stringify(payload, null, 2)}\n`);
}
if (failure) throw failure;
