// Early D/E feedback only; this does not satisfy the seven-type, three-bucket or six-run matrix.
import assert from 'node:assert/strict';
import { existsSync, mkdirSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { createConnection } from 'node:net';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, binaryIdentity, configureFreshSQLite,
  openFreshAdvanced, openFreshUI, previewAndApplyFreshSchema,
  setupFreshDevice, setupFreshRules, startFreshBasic, startFreshGateway, removeOwnedWork, startFreshSimulator, stopProcesses,
} from './fresh-ui.mjs';
import { boundedFetch, sh } from './lib.mjs';
import { preflightSQLiteDestination } from './destination-preflight.mjs';

const runId = `early-root-${Date.now()}`;
const work = `/tmp/gw-f-${runId}`;
const destination = join(work, 'destination.db');
const gatewayPort = Number(process.env.F_EARLY_ROOT_GW_PORT ?? 3383);
assert.ok(Number.isInteger(gatewayPort) && gatewayPort >= 1024 && gatewayPort <= 65535);
const base = `http://127.0.0.1:${gatewayPort}`;
const simulatorPort = 15054;
const gatewayBinary = process.env.F_EARLY_ROOT_BINARY ?? DEFAULT_GATEWAY_BINARY;
const result = { run_id: runId, kind: 'early-single-sqlite-feedback', passed: false,
  linked_build_evidence: gatewayBinary.endsWith('sqlite-probe-v2')
    ? 'normal-build-sqlite-probe-v2-root-20261004.json' : 'normal-build-root-20261004.json',
  command: 'node scripts/tests/f_device_to_sql/fresh-early-root.mjs', environment: 'darwin/arm64; loopback simulator; disposable SQLite',
  expected_values: [215, 1013], limitations: ['Only two D/E acquisition points.', 'Not the seven-type/three-bucket/six-timing acceptance.'] };
const children = [];
let ownsWork = false;
let ui;
let stage = 'boot';

function sanitize(value) {
  if (typeof value === 'string') return value.replaceAll(work, '<owned-run>')
    .replace(/\u001b\[[0-9;]*m/g, '').replace(/\/Users\/[^/\s]+/g, '<user-home>')
    .replace(/\/(?:private\/)?var\/folders\/[^\s"']+/g, '<owned-temp>');
  if (Array.isArray(value)) return value.map(sanitize);
  if (value && typeof value === 'object') return Object.fromEntries(Object.entries(value).map(([key, item]) => [key, sanitize(item)]));
  return value;
}

async function observe(path) {
  const response = await boundedFetch(`${base}/api/v1/datalink${path}`, {}, 10_000);
  assert.equal(response.ok, true, `GET ${path} must succeed`);
  const body = await response.json();
  assert.ok(body?.data, `GET ${path} must contain data`);
  return body.data;
}

async function assertOwnedPortFree(port) {
  await new Promise((resolve, reject) => {
    const socket = createConnection({ host: '127.0.0.1', port });
    socket.setTimeout(1500, () => { socket.destroy(); reject(new Error('owned port availability could not be proven')); });
    socket.once('connect', () => { socket.destroy(); reject(new Error(`owned port ${port} is already occupied`)); });
    socket.once('error', (error) => {
      socket.destroy();
      if (error.code === 'ECONNREFUSED') resolve();
      else reject(error);
    });
  });
}

async function saveEarlyMappings() {
  const rows = ui.page.locator('[data-testid^="mapping-row-"]');
  await rows.first().waitFor();
  assert.equal(await rows.count(), 2);
  const tags = [];
  for (let index = 0; index < 2; index += 1) {
    const pointID = (await rows.nth(index).getAttribute('data-testid')).replace('mapping-row-', '');
    const tag = `early.${index + 1}`;
    for (const [key, testID, value] of [
      ['tag_key', 'input-tag-key', tag], ['display_name', 'input-display-name', tag],
      ['target_type', 'select-target-type', index === 0 ? 'int16' : 'uint16'],
    ]) {
      const field = ui.page.getByTestId(`${testID}-${pointID}`);
      if (key === 'target_type') await field.selectOption(value);
      else await field.fill(value);
      await field.press('Tab');
      const deadline = Date.now() + 20_000;
      while (true) {
        const data = await observe('/studio-v2/workspace/mappings');
        const saved = Array.isArray(data) ? data : data.mappings ?? [];
        const badge = await ui.page.getByTestId(`mapping-save-state-${pointID}`).innerText().catch(() => '');
        if (saved.some((item) => item.tag_key === tag && item[key] === value) &&
          await field.inputValue() === value && /saved|已儲存|已保存/i.test(badge)) break;
        if (Date.now() >= deadline) throw new Error(`early mapping ${index} ${key} did not persist exact desired value`);
        await new Promise((resolve) => setTimeout(resolve, 300));
      }
    }
    tags.push({ point_id: pointID, tag_key: tag });
  }
  await ui.page.getByTestId('btn-continue').click();
  await ui.page.getByTestId('step-nav-button-4').waitFor();
  return tags;
}

try {
  mkdirSync(DEFAULT_EVIDENCE_DIR, { recursive: true });
  assert.equal(existsSync(work), false, 'an early run must never reuse a workspace');
  mkdirSync(work);
  ownsWork = true;
  result.source_head = sh('git', ['rev-parse', 'HEAD']);
  result.binary_sha256 = binaryIdentity(gatewayBinary).sha256;
  result.harness_source_sha256 = Object.fromEntries(['fresh-early-root.mjs', 'fresh-ui.mjs', 'destination-preflight.mjs']
    .map((name) => [name, binaryIdentity(new URL(name, import.meta.url).pathname).sha256]));
  result.simulator = { binary_sha256: binaryIdentity(DEFAULT_SIMULATOR_BINARY).sha256,
    port: simulatorPort, owned: true, registers: { 0: 215, 1: 1013 } };
  result.preflight = preflightSQLiteDestination({ path: destination, work });
  assert.equal(result.preflight.exists, false);
  await assertOwnedPortFree(gatewayPort);
  await assertOwnedPortFree(simulatorPort);
  const simulator = startFreshSimulator({ work, name: 'early-root-simulator', port: simulatorPort,
    registers: result.simulator.registers });
  children.push(simulator);
  result.simulator.pid = simulator.pid;
  children.push(startFreshGateway({ work, binary: gatewayBinary, port: gatewayPort }));
  const bootDeadline = Date.now() + 30_000;
  while (true) {
    try { await observe('/studio-v2/workspace'); break; }
    catch (error) {
      if (Date.now() >= bootDeadline) throw error;
      await new Promise((resolve) => setTimeout(resolve, 300));
    }
  }
  ui = await openFreshUI({ port: gatewayPort });
  result.browser_version = ui.browser.version();
  result.t_open = ui.t_open;
  stage = 'device';
  await setupFreshDevice(ui.page, { name: 'F early A', port: simulatorPort });
  stage = 'rules';
  await setupFreshRules(ui.page, [
    { name: 'Early temperature', start: 40001, count: 1, dataType: 'int16', prefix: 'early_temperature' },
    { name: 'Early pressure', start: 40002, count: 1, dataType: 'uint16', prefix: 'early_pressure' },
  ]);
  stage = 'mapping';
  const mappings = await saveEarlyMappings();
  assert.equal(mappings.length, 2);
  stage = 'destination';
  await configureFreshSQLite(ui.page, base, { path: destination, table: 'readings' });
  result.destination_after_save = preflightSQLiteDestination({ path: destination, work });
  assert.equal(existsSync(destination), false);
  stage = 'first-start-without-proof';
  const initial = await startFreshBasic(ui.page);
  result.initial_status = initial.status_text;
  assert.equal(existsSync(destination), false, 'start without proof must perform no DDL');
  const saved = await observe('/studio-v2/workspace/write-groups');
  assert.equal(saved.groups.length, 1);
  const group = saved.groups[0];
  assert.equal(group.members.length, 2);
  result.group_id = group.id;
  stage = 'preview-confirm';
  await openFreshAdvanced(ui.page);
  await ui.page.getByTestId(`group-open-${group.id}`).click();
  await ui.page.getByTestId('group-schema-preview').click();
  await ui.page.getByTestId('group-schema-preview-result').waitFor({ timeout: 20_000 });
  assert.equal(existsSync(destination), false, 'preview leaves the missing destination absent');
  result.preview_left_destination_absent = true;
  // The shared helper performs a fresh preview followed by the explicit UI confirmation.
  result.schema = await previewAndApplyFreshSchema(ui.page);
  stage = 'start';
  await ui.page.reload();
  await ui.page.getByTestId('step-nav-button-4').click();
  const started = await startFreshBasic(ui.page);
  result.t_start = started.t_start;
  result.start_status = started.status_text;
  stage = 'first-independent-sql';
  const sqlDeadline = Date.now() + 180_000;
  while (true) {
    const closed = Number(sh('sqlite3', ['-readonly', `file:${join(work, 'gateway.db')}?mode=ro`,
      `SELECT COUNT(*) FROM wg_delivery_outbox WHERE group_id = '${group.id.replaceAll("'", "''")}'`]));
    if (closed > 0 && !result.t_closed) result.t_closed = process.hrtime.bigint().toString();
    const sql = `SELECT COUNT(*) FROM "${group.destination.table_name.replaceAll('"', '""')}"`;
    const rows = Number(sh('sqlite3', ['-readonly', `file:${destination}?mode=ro`, sql]));
    if (rows > 0) {
      result.t_sql = process.hrtime.bigint().toString();
      result.sql_rows = JSON.parse(sh('sqlite3', ['-readonly', '-json', `file:${destination}?mode=ro`,
        `SELECT * FROM "${group.destination.table_name.replaceAll('"', '""')}" ORDER BY rowid LIMIT 1`]));
      break;
    }
    if (Date.now() >= sqlDeadline) throw new Error('no first SQL row before the bounded deadline');
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  result.current_groups = await observe('/studio-v2/workspace/write-groups');
  result.delivery = await observe(`/studio-v2/workspace/write-groups/${group.id}/delivery`);
  result.source_mappings = mappings;
  const mappingData = await observe('/studio-v2/workspace/mappings');
  result.persisted_mappings = Array.isArray(mappingData) ? mappingData : mappingData.mappings ?? [];
  for (const member of result.current_groups.groups[0].members) {
    const row = result.sql_rows[0];
    const persisted = result.persisted_mappings.find((mapping) => mapping.point_id === member.point_id &&
      mapping.tag_id === member.tag_id && mapping.device_id === member.device_id);
    const index = mappings.findIndex((mapping) => mapping.tag_key === persisted?.tag_key);
    assert.ok(index >= 0, 'saved member must match an original UI mapping');
    assert.equal(Number(row[member.target_column]), result.expected_values[index]);
  }
  result.end_to_end_ms = Number(BigInt(result.t_sql) - BigInt(result.t_open)) / 1e6;
  assert.ok(result.t_closed, 'first row must have an independently observed durable closure');
  result.setup_ms = Number(BigInt(result.t_start) - BigInt(result.t_open)) / 1e6;
  result.system_wait_ms = Number(BigInt(result.t_closed) - BigInt(result.t_start)) / 1e6;
  result.delivery_observed_ms = Number(BigInt(result.t_sql) - BigInt(result.t_closed)) / 1e6;
  result.poll_interval_ms = 1000;
  result.polling_limitation = 'Closure is observed before target SELECT in each iteration; observation delay is up to one polling interval plus SQL command time.';
  result.passed = true;
} catch (error) {
  result.failed_stage = stage;
  result.error = sanitize(String(error.message).split('\n')[0]);
  const launchFailure = String(error.message).split('\n').find((line) => line.includes('bootstrap_check_in'));
  if (launchFailure) result.browser_launch_failure = sanitize(launchFailure);
  if (ui) {
    result.ui_text = (await ui.page.locator('body').innerText().catch(() => '')).replaceAll(work, '<owned-run>').slice(-12_000);
  }
} finally {
  try { await ui?.browser.close(); }
  catch (error) { result.browser_cleanup_error = error.message; result.passed = false; }
  const exits = await stopProcesses(children);
  result.child_exit = exits;
  result.child_exit_passed = exits.length === children.length && exits.every((exit) => exit.state === 'exited' && exit.exit_code === 0 && !exit.spawn_error && !exit.stop_error && !exit.kill_sent);
  if (!result.child_exit_passed) result.passed = false;
  const workCleanup = ownsWork ? removeOwnedWork(work) : null;
  result.owned_cleanup_passed = workCleanup ? workCleanup.removed : 'not-created';
  if (workCleanup?.error) result.owned_cleanup_error = workCleanup.error;
  if (!result.owned_cleanup_passed) result.passed = false;
  writeFileSync(join(DEFAULT_EVIDENCE_DIR, `${runId}.json`), `${JSON.stringify(sanitize(result), null, 2)}\n`);
  console.log(JSON.stringify({ run_id: runId, passed: result.passed, failed_stage: result.failed_stage, error: result.error }));
}
process.exitCode = result.passed ? 0 : 1;
