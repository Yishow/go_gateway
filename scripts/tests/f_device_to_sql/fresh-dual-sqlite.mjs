// Fresh two-device Basic groups: all source setup and schema confirmation are
// performed through the embedded UI; SQLite is observed only through read-only
// queries after the UI has explicitly confirmed each managed schema.
import { join } from 'node:path';
import { existsSync } from 'node:fs';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, assertFreshPortFree,
  binaryIdentity, closeFreshUI, configureFreshSQLite, ensureFreshSidebarCollapsed, isFreshCleanExit,
  makeFreshRun, observeJSON, openFreshUI, persistFreshMappings, previewAndApplyFreshSchema,
  prepareFreshBasicGroup, setupFreshMultiDeviceSources, startFreshBasic, startFreshGateway,
  startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import {
  FRESH_EXPECTED_VALUES, FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES,
  expectedValuesByTag, verifyFreshDeliveryEvidence,
} from './fresh-recording-verification.mjs';
import {
  readSQLiteDelivery, readSQLiteRows, waitForClosedSQLiteBuckets, waitForSQLiteRows,
} from './fresh-sqlite-observation.mjs';
import { assertFreshTimingWitness } from './fresh-timing.mjs';
import { writeFreshResult } from './fresh-result.mjs';
import { addDevice, addRule } from './lib.mjs';
import { preflightSQLiteDestination } from './destination-preflight.mjs';

const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
const observedTimings = [];
const runId = process.env.F_RUN_ID ?? `dual-sqlite-${Date.now()}`;
const port = Number(process.env.GW_PORT ?? 3380);
const simPortA = Number(process.env.F_SIM_PORT_A ?? 15050);
const simPortB = Number(process.env.F_SIM_PORT_B ?? 15051);
const run = makeFreshRun(runId);
const base = `http://127.0.0.1:${port}`;
const target = run.destination_db;
const evidencePath = join(DEFAULT_EVIDENCE_DIR, `${runId}.json`);
const children = [];
const lines = [];
let browser;
let page;
let ui;
let result;
let failure;

const valuesByAddress = (values) => Object.fromEntries(FRESH_RULES.map((rule, index) => [String(rule.start), values[index]]));
const valuesA = valuesByAddress(FRESH_EXPECTED_VALUES);
const registersB = [...FRESH_REGISTERS];
registersB.splice(0, 14, 187, 777, 0, 0, 0, 2, 0, 0xc010, 0, 0xc084, 0x7291, 0x6872, 0xb021, 0x4232);
// The string source uses the same two-byte fixture register on both loopback
// devices; device identity is proven by persisted IDs/provenance, not by a
// display value that the simulator does not vary.
const valuesB = valuesByAddress(['187', '777', '2', '0', '-2.25', '-654.321', FRESH_EXPECTED_VALUES[6], FRESH_EXPECTED_VALUES[7], 'A1']);
const log = (message) => { lines.push(message); console.log(message); };

function quoteIdentifier(value) { return `"${String(value).replaceAll('"', '""')}"`; }

function tableColumns(table) {
  return readSQLiteRows(target, `PRAGMA table_info(${quoteIdentifier(table)})`)
    .map((column) => ({ name: column.name }));
}

function tableSelect(table, columns) {
  return columns.map(({ name }) => {
    const quoted = quoteIdentifier(name);
    return `CASE WHEN ${quoted} IS NULL THEN NULL ELSE printf('%s', ${quoted}) END AS ${quoted}`;
  }).join(', ');
}

async function currentGroups(deviceIds) {
  const data = await observeJSON(base, '/studio-v2/workspace/write-groups');
  const groups = (data.groups ?? []).filter((group) => deviceIds.includes(group.basic_managed_device_id));
  if (groups.length !== deviceIds.length) throw new Error(`expected one Basic group per device; got ${groups.length}`);
  for (const deviceId of deviceIds) {
    if (groups.filter((group) => group.basic_managed_device_id === deviceId).length !== 1) {
      throw new Error(`Basic group ownership is not unique for device ${deviceId}`);
    }
  }
  return groups;
}

function expectedForGroup(group, mappings, expectedByTag) {
  return Object.fromEntries(group.members.map((member) => {
    const mapping = mappings.find((item) => item.device_id === member.device_id &&
      item.point_id === member.point_id && item.tag_id === member.tag_id);
    if (!mapping) throw new Error(`missing persisted mapping for ${member.device_id}/${member.point_id}/${member.tag_id}`);
    return [mapping.tag_key, expectedByTag[mapping.tag_key]];
  }));
}

async function captureWidths() {
  await ensureFreshSidebarCollapsed(page);
  const screenshots = [];
  for (const [width, height] of [[390, 1000], [768, 1100], [1440, 1400]]) {
    await page.setViewportSize({ width, height });
    await page.waitForTimeout(300);
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth > innerWidth);
    const basicReadable = await page.getByTestId('basic-recording-panel').isVisible();
    const file = `${runId}-${width}.png`;
    await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, file), fullPage: true });
    screenshots.push({ width, file, document_horizontal_overflow: overflow, basic_panel_visible: basicReadable });
  }
  return screenshots;
}

async function runGroup(group, mappings, expectedByTag) {
  if (!group.applied_revision) throw new Error(`Basic group ${group.id} has no applied revision`);
  const table = group.destination.table_name;
  const columns = tableColumns(table);
  const select = tableSelect(table, columns);
  const [closed, destination] = await Promise.all([
    waitForClosedSQLiteBuckets({ path: run.gateway_db, groupId: group.id, groupRevision: group.applied_revision }),
    waitForSQLiteRows({ path: target, sql: `SELECT ${select} FROM ${quoteIdentifier(table)} ORDER BY bucket_start, record_id`, minimum: 3 }),
  ]);
  const delivery = readSQLiteDelivery({ path: run.gateway_db, groupId: group.id, groupRevision: group.applied_revision });
  const verification = verifyFreshDeliveryEvidence({
    group, table, rows: destination.rows, mappings, expectedByTag: expectedForGroup(group, mappings, expectedByTag),
    closedBuckets: delivery.closedBuckets, outbox: delivery.outbox, receipts: delivery.receipts,
  });
  return { group_id: group.id, device_id: group.basic_managed_device_id, table, columns,
    rows: destination.rows, polls: destination.polls, t_closed: closed.t_closed, t_sql: destination.t_sql,
    delivery, verification };
}

try {
  if (port === simPortA || port === simPortB || simPortA === simPortB) throw new Error('fresh gateway and simulator ports must be distinct');
  const initial = preflightSQLiteDestination({ path: target, work: run.work });
  if (initial.recording_table_present !== false || existsSync(target)) throw new Error('fresh SQLite target was present before UI');
  await Promise.all([assertFreshPortFree(port), assertFreshPortFree(simPortA), assertFreshPortFree(simPortB)]);
  children.push(startFreshSimulator({ work: run.work, name: 'sim-a', port: simPortA,
    registers: Object.fromEntries(FRESH_REGISTERS.map((value, index) => [index, value])) }));
  children.push(startFreshSimulator({ work: run.work, name: 'sim-b', port: simPortB,
    registers: Object.fromEntries(registersB.map((value, index) => [index, value])) }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port }));
  await waitForFreshGateway(base);
  ui = await openFreshUI({ port, width: 1440, height: 1400 });
  ({ browser, page } = ui);
  const devices = [{ name: 'Line A', port: simPortA, prefix: 'A' }, { name: 'Line B', port: simPortB, prefix: 'B' }];
  await setupFreshMultiDeviceSources(page, devices, FRESH_RULES.map((rule) => ({ ...rule })));
  const targetTypes = [...FRESH_TARGET_TYPES, ...FRESH_TARGET_TYPES];
  const targetTypeByAddress = new Map(FRESH_RULES.map((rule, index) => [String(rule.start), FRESH_TARGET_TYPES[index]]));
  const tags = await persistFreshMappings(page, base, {
    tagPrefix: 'dual', targetTypes,
    targetTypeFor: ({ address }) => targetTypeByAddress.get(address),
  });
  const deviceData = await observeJSON(base, '/studio-v2/workspace/devices');
  const savedDevices = Array.isArray(deviceData) ? deviceData : deviceData.devices ?? [];
  const deviceIds = devices.map((device) => savedDevices.find((saved) => saved.name === device.name)?.id);
  if (deviceIds.some((id) => !id)) throw new Error('fresh dual device IDs were not persisted');
  const valuesByDevice = new Map([[deviceIds[0], valuesA], [deviceIds[1], valuesB]]);
  const expectedByTag = Object.fromEntries(tags.map((tag) => [tag.tag_key, valuesByDevice.get(tag.device_id)?.[tag.address]]));
  if (Object.values(expectedByTag).some((value) => value === undefined)) throw new Error('fresh dual mapping identity did not expose device values');
  const destination = await configureFreshSQLite(page, base, { path: target, table: 'f_connector_placeholder' });
  const afterConnector = preflightSQLiteDestination({ path: target, work: run.work });
  if (afterConnector.exists) throw new Error('SQLite connector save created the destination before schema confirmation');
  const prepared = [];
  for (const deviceId of deviceIds) {
    const first = await prepareFreshBasicGroup(page, base, deviceId);
    const schema = await previewAndApplyFreshSchema(page);
    await page.getByTestId('group-editor-close').click();
    const start = await startFreshBasic(page);
    prepared.push({ device_id: deviceId, group: first.group, schema, start });
  }
  const groups = await currentGroups(deviceIds);
  if (groups.some((group) => group.members.length !== FRESH_RULES.length)) throw new Error('dual Basic groups do not each contain nine members');
  const observations = [];
  // Observe both groups at once: a serial loop would stamp the second group's first
  // SQL row only after the first group's three buckets, inflating its t_sql.
  const mappings = await observeJSON(base, '/studio-v2/workspace/mappings').then((data) => Array.isArray(data) ? data : data.mappings ?? []);
  observations.push(...await Promise.all(groups.map((group) => runGroup(group, mappings, expectedByTag))));
  const timings = observations.map((observation) => {
    const preparedGroup = prepared.find((candidate) => candidate.device_id === observation.device_id);
    if (!preparedGroup) throw new Error(`missing prepared timing source for device ${observation.device_id}`);
    // Keep the raw clocks so a failed timing criterion still leaves its measurements.
    observedTimings.push({ device_id: observation.device_id, t_open: ui.t_open, t_start: preparedGroup.start.t_start, t_closed: observation.t_closed, t_sql: observation.t_sql, poll_count: observation.polls?.length });
    return {
      device_id: observation.device_id,
      ...assertFreshTimingWitness({
        t_open: ui.t_open,
        t_start: preparedGroup.start.t_start,
        t_closed: observation.t_closed,
        t_sql: observation.t_sql,
        polls: observation.polls,
      }),
    };
  });
  const failures = observations.flatMap((item) => item.verification.failures.map((failureText) => `${item.device_id}: ${failureText}`));
  const screenshots = await captureWidths();
  result = {
    run_id: runId, passed: failures.length === 0, failures, kind: 'sqlite-dual-basic-same-address', destination,
    target_preflight: { before_ui: initial, after_connector_save: afterConnector }, binary: {
      gateway: binaryIdentity(gatewayBinary), simulator: binaryIdentity(DEFAULT_SIMULATOR_BINARY),
    }, devices: savedDevices.filter((device) => deviceIds.includes(device.id)), rules: FRESH_RULES, tags,
    expected_by_tag: expectedByTag, groups, observations, prepared, screenshots, page_errors: ui.pageErrors,
    api_traffic: ui.apiTraffic, clocks: { t_open: ui.t_open }, timings, log: lines, retained_work: run.work,
  };
  log(result.passed ? `PASS ${runId}: dual Basic groups with three buckets each` : `FAIL ${runId}: ${failures.join('; ')}`);
  process.exitCode = result.passed ? 0 : 1;
} catch (error) {
  failure = error;
  process.exitCode = 1;
  log(`ERROR ${runId}: ${error?.message?.split('\n')[0] ?? String(error)}`);
} finally {
  const cleanup = await stopProcesses(children).catch((error) => [{ state: 'error', error: error.message }]);
  let browserCleanup = { state: 'closed' };
  try { await closeFreshUI(browser); } catch (error) { browserCleanup = { state: 'error', error: error.message.split('\n')[0] }; }
  const processCleanupOk = cleanup.every(isFreshCleanExit);
  const payload = result ?? { run_id: runId, passed: false, failures: [failure?.message ?? 'fresh dual run failed'], timings_observed: observedTimings, retained_work: run.work, log: lines };
  const namespace = { path: run.work, absent: false, state: 'pending' };
  if (processCleanupOk && browserCleanup.state === 'closed') {
    try {
      const { rmSync } = await import('node:fs');
      rmSync(run.work, { recursive: true, force: false });
      namespace.state = 'removed'; namespace.absent = !existsSync(run.work);
    } catch (error) { namespace.state = 'error'; namespace.error = error.message.split('\n')[0]; }
  } else namespace.state = 'blocked';
  payload.cleanup = { processes: cleanup, browser: browserCleanup, namespace };
  if (!processCleanupOk || browserCleanup.state !== 'closed' || !namespace.absent) {
    payload.passed = false;
    payload.failures = [...new Set([...(payload.failures ?? []), 'cleanup did not prove all children/browser/namespace closed'])];
  }
  writeFreshResult(evidencePath, payload);
}

if (failure) throw failure;
