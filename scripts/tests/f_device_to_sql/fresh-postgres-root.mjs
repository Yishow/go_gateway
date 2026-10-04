// F first-use PostgreSQL: normal embedded UI, no pre-created recording table.
import assert from 'node:assert/strict';
import { writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, assertFreshPortFree,
  binaryIdentity, makeFreshRun, monotonicNs, observeJSON, openFreshUI, persistFreshMappings,
  prepareFreshBasicGroup, previewAndApplyFreshSchema, setupFreshDevice, setupFreshRules,
  startFreshBasic, startFreshGateway, removeOwnedWork, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import {
  FRESH_EXPECTED_VALUES, FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES,
  expectedValuesByTag, verifyFreshDeliveryEvidence,
} from './fresh-recording-verification.mjs';
import { createOwnedPostgresSchema } from './destination-preflight.mjs';
import { configureFreshPostgres } from './fresh-postgres-setup.mjs';
import {
  localDurableRows, ownedPGTarget, pgColumnMetadata, pgCommand, pgRecordingRelations,
  pgRecordingRows, quotePGIdentifier,
} from './fresh-postgres-observation-root.mjs';
import { sh } from './lib.mjs';

const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
const runID = `pg-root-${Date.now()}`;
const schema = `gw_f_pg_root_${Date.now()}`;
const gatewayPort = 3384;
const simulatorPort = 15055;
const base = `http://127.0.0.1:${gatewayPort}`;
const result = { run_id: runID, kind: 'postgres-single-seven-types', passed: false,
  command: 'node scripts/tests/f_device_to_sql/fresh-postgres-root.mjs',
  environment: 'darwin/arm64; owned loopback simulator and disposable PostgreSQL schema',
  linked_build_evidence: 'normal-build-ui-evidence-v3-root-20261004.json',
  poll_interval_ms: 1000, limitations: ['Single device; multi-device and same-group entity runs are separate.'] };
let run;
let target;
let createdSchema = false;
let ui;
let stage = 'preflight';
const children = [];

function sanitize(value) {
  if (Array.isArray(value)) return value.map(sanitize);
  if (value && typeof value === 'object') return Object.fromEntries(Object.entries(value).map(([key, item]) =>
    [key, /password|credential|authorization|cookie|dsn|token/i.test(key) ? '[redacted]' : sanitize(item)]));
  if (typeof value !== 'string') return value;
  return value.replaceAll(run?.work ?? '<never-used>', '<owned-run>')
    .replace(/\/(?:Users|home)\/[^/\s]+/g, '<user-home>')
    .replace(/\/(?:private\/)?var\/folders\/[^\s"']+/g, '<owned-temp>')
    .replace(/password=[^\s]+/gi, 'password=[redacted]');
}

function literal(value) { return `'${String(value).replaceAll("'", "''")}'`; }

async function captureUIBounds(suffix = 'diagnostic') {
  const toggle = ui.page.getByTitle(/^收合側邊欄/);
  if (await toggle.isVisible()) {
    await toggle.focus();
    await toggle.press('Control+b');
    await ui.page.getByTitle(/^展開側邊欄/).waitFor({ timeout: 5000 });
  }
  const connector = ui.page.getByRole('button', { name: /Collapse connector setup|收合連線設定/i });
  if (await connector.isVisible()) await connector.click();
  const advanced = ui.page.getByTestId('step4-advanced-recording');
  if ((await advanced.getAttribute('open')) !== null) await advanced.locator('summary').press('Enter');
  const observations = [];
  for (const width of [390, 768, 1440]) {
    await ui.page.setViewportSize({ width, height: 1000 });
    await ui.page.getByTestId('basic-recording-panel').scrollIntoViewIfNeeded();
    await ui.page.waitForTimeout(300);
    const bounds = await ui.page.evaluate(() => {
      const basic = document.querySelector('[data-testid="basic-recording-panel"]');
      const rect = basic.getBoundingClientRect();
      const outside = [...document.querySelectorAll('body *')].filter((element) => {
        if (typeof element.checkVisibility === 'function' && !element.checkVisibility()) return false;
        const box = element.getBoundingClientRect();
        return box.width > 0 && box.height > 0 && box.right > innerWidth + 1;
      }).map((element) => ({ tag: element.tagName, test_id: element.dataset.testid,
        classes: element.className, text: element.textContent?.slice(0, 80),
        right: element.getBoundingClientRect().right, width: element.getBoundingClientRect().width,
        in_basic_recording: basic.contains(element) }));
      return { viewport_width: innerWidth, document_width: document.documentElement.scrollWidth,
        basic_bounds: { left: rect.left, right: rect.right, width: rect.width },
        basic_outside: outside.filter((element) => element.in_basic_recording), outside: outside.slice(-30) };
    });
    const file = `${runID}-${suffix}-${width}.png`;
    const height = await ui.page.evaluate(() => { window.scrollTo(0, 0); return document.documentElement.scrollHeight; });
    await ui.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, file), clip: { x: 0, y: 0, width, height } });
    await ui.page.getByTestId('basic-recording-start').scrollIntoViewIfNeeded();
    const controlsFile = `${runID}-${suffix}-controls-${width}.png`;
    await ui.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, controlsFile) });
    observations.push({ width, file, controls_file: controlsFile, ...bounds });
  }
  return observations;
}

function durableSnapshot(group) {
  const scope = `group_id=${literal(group.id)} AND group_revision=${literal(group.applied_revision)}`;
  return {
    buckets: localDurableRows(run.gateway_db, `SELECT group_id,group_revision,entity_key,bucket_start,kind,reason,record_id,members
      FROM wg_delivery_buckets WHERE ${scope} AND kind='row' ORDER BY bucket_start,entity_key`),
    outbox: localDurableRows(run.gateway_db, `SELECT effect_key,record_id,group_id,group_revision,entity_key,bucket_start,
      table_name,state,payload_digest,committed_at FROM wg_delivery_outbox WHERE ${scope} ORDER BY bucket_start,record_id`),
    receipts: localDurableRows(run.gateway_db, `SELECT effect_key,payload_digest,committed_at FROM wg_delivery_receipts
      WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE ${scope})`),
    checkpoints: localDurableRows(run.gateway_db, `SELECT group_id,group_revision,next_close FROM wg_delivery_checkpoints WHERE ${scope}`),
  };
}

try {
  target = ownedPGTarget(schema);
  run = makeFreshRun(runID);
  result.source_head = sh('git', ['rev-parse', 'HEAD']);
  result.binary_sha256 = binaryIdentity(gatewayBinary).sha256;
  result.simulator_sha256 = binaryIdentity(DEFAULT_SIMULATOR_BINARY).sha256;
  const dependencies = ['fresh-postgres-root.mjs', 'fresh-postgres-observation-root.mjs', 'fresh-postgres-setup.mjs',
    'fresh-recording-verification.mjs', 'fresh-ui.mjs', 'destination-preflight.mjs', 'lib.mjs'];
  result.harness_source_sha256 = Object.fromEntries(dependencies.map((name) =>
    [name, binaryIdentity(new URL(name, import.meta.url).pathname).sha256]));
  result.target_preflight = createOwnedPostgresSchema({ psql: pgCommand, schema, table: target.table });
  createdSchema = true;
  result.relations_before_ui = pgRecordingRelations(schema);
  assert.deepEqual(result.relations_before_ui, [], 'fresh schema must contain no recording objects');
  await assertFreshPortFree(gatewayPort);
  await assertFreshPortFree(simulatorPort);
  const registers = Object.fromEntries(FRESH_REGISTERS.map((value, index) => [index, value]));
  children.push(startFreshSimulator({ work: run.work, name: 'pg-root-simulator', port: simulatorPort, registers }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port: gatewayPort }));
  await waitForFreshGateway(base);
  ui = await openFreshUI({ port: gatewayPort });
  result.browser_version = ui.browser.version();
  result.t_open = ui.t_open;
  stage = 'device';
  await setupFreshDevice(ui.page, { name: 'F PostgreSQL A', port: simulatorPort });
  stage = 'rules';
  await setupFreshRules(ui.page, FRESH_RULES.map((rule) => ({ ...rule })));
  stage = 'mappings';
  result.tags = await persistFreshMappings(ui.page, base, { tagPrefix: 'pg_a', targetTypes: FRESH_TARGET_TYPES });
  assert.equal(result.tags.length, 9);
  stage = 'destination';
  result.connector_setup = await configureFreshPostgres(ui.page, base, target);
  assert.equal(await ui.page.getByTestId('basic-recording-interval').inputValue(), '60', 'new draft default remains 60 seconds');
  stage = 'schema';
  const prepared = await prepareFreshBasicGroup(ui.page, base);
  result.first_prepare_status = prepared.first_status_text;
  assert.deepEqual(pgRecordingRelations(schema), [], 'saving the Basic group must not create recording tables');
  result.schema = await previewAndApplyFreshSchema(ui.page);
  await ui.page.getByTestId('group-editor-close').click();
  stage = 'start';
  const started = await startFreshBasic(ui.page);
  result.t_start = started.t_start;
  result.start_status = started.status_text;
  assert.match(result.start_status, /succeeded|成功/);
  const groups = await observeJSON(base, '/studio-v2/workspace/write-groups');
  result.group = groups.groups.find((group) => group.id === prepared.group.id);
  assert.ok(result.group?.applied_revision, 'canonical group must actually be applied');
  assert.equal(result.group.row_policy.interval_seconds, 60);
  assert.equal(result.group.destination.table_schema, schema);
  const table = result.group.destination.table_name;
  result.columns = pgColumnMetadata(schema, table);
  assert.ok(result.columns.length > 9);
  if (process.env.F_UI_DIAGNOSTIC === '1') {
    stage = 'ui-diagnostic-only';
    result.ui_diagnostic = await captureUIBounds();
    throw new Error('UI diagnostic only; no three-bucket acceptance verdict');
  }
  stage = 'three-formal-buckets';
  const deadline = Date.now() + 270_000;
  result.polls = [];
  result.poll_count = 0;
  while (true) {
    result.durable = durableSnapshot(result.group);
    if (!result.t_closed && result.durable.buckets.length > 0) result.t_closed = monotonicNs();
    result.rows = pgRecordingRows(schema, table, result.columns);
    const observed = monotonicNs();
    if (!result.t_sql && result.rows.length > 0) result.t_sql = observed;
    result.poll_count += 1;
    const previousPoll = result.polls.at(-1);
    if (!previousPoll || previousPoll.local_closed_rows !== result.durable.buckets.length || previousPoll.sql_rows !== result.rows.length) {
      result.polls.push({ observed_at_ns: observed, local_closed_rows: result.durable.buckets.length, sql_rows: result.rows.length });
    }
    const settled = result.rows.slice(0, 3).every((row) => {
      const outbox = result.durable.outbox.find((item) => item.record_id === row.record_id && item.state === 'sql_committed');
      return outbox && result.durable.receipts.some((receipt) => receipt.effect_key === outbox.effect_key &&
        receipt.payload_digest === outbox.payload_digest && receipt.committed_at);
    });
    if (result.rows.length >= 3 && result.durable.buckets.length >= 3 && settled) break;
    if (Date.now() >= deadline) throw new Error('three consecutive formal buckets were not observed within the bounded window');
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  const data = await observeJSON(base, '/studio-v2/workspace/mappings');
  result.mappings = Array.isArray(data) ? data : data.mappings;
  const expected = [...FRESH_EXPECTED_VALUES];
  expected[3] = 'true';
  result.verification = verifyFreshDeliveryEvidence({ group: result.group, table, rows: result.rows,
    mappings: result.mappings, expectedByTag: expectedValuesByTag(result.tags, expected),
    closedBuckets: result.durable.buckets, outbox: result.durable.outbox, receipts: result.durable.receipts });
  assert.deepEqual(result.verification.failures, [], 'independent SQL/durable/source assertions must all pass');
  assert.equal(result.durable.checkpoints.length, 1);
  assert.ok(Date.parse(result.durable.checkpoints[0].next_close) >= Date.parse(result.rows[2].bucket_start) + 60_000);
  result.delivery = await observeJSON(base, `/studio-v2/workspace/write-groups/${result.group.id}/delivery`);
  result.setup_ms = Number(BigInt(result.t_start) - BigInt(result.t_open)) / 1e6;
  result.system_wait_ms = Number(BigInt(result.t_closed) - BigInt(result.t_start)) / 1e6;
  result.delivery_observed_ms = Number(BigInt(result.t_sql) - BigInt(result.t_closed)) / 1e6;
  result.end_to_end_ms = Number(BigInt(result.t_sql) - BigInt(result.t_open)) / 1e6;
  assert.ok([result.setup_ms, result.system_wait_ms, result.delivery_observed_ms].every((value) => value >= 0),
    'monotonic observation segments must not be negative');
  assert.ok(result.end_to_end_ms <= 300_000, 'controlled first independent SQL must be at most 300 seconds');
  result.polling_limitation = 'Local closure SELECT precedes target SELECT in each loop; observation delay includes up to one poll interval and actual SQL command durations.';
  assert.deepEqual(ui.pageErrors, [], 'normal UI must have no uncaught browser errors');
  result.page_errors = ui.pageErrors;
  stage = 'responsive-main-ui';
  await ui.page.getByTestId('basic-recording-committed-effect').waitFor({ timeout: 20_000 });
  assert.equal(await ui.page.getByTestId('basic-recording-panel').getByText(/first complete UTC snapshot is still pending|第一個完整.*等待/i).count(), 0,
    'current SQL receipt must remove the first-bucket pending claim');
  result.screenshots = await captureUIBounds('main');
  result.main_view = { sidebar_control_b: true, connector_collapsed: true, advanced_collapsed: true,
    limitation: 'Basic recording controls are bounded below. Existing long connector-summary labels can overflow the document on mobile, as preserved E layout context; full-page responsive redesign is outside this change.' };
  for (const observation of result.screenshots) {
    assert.ok(observation.basic_bounds.left >= 0 && observation.basic_bounds.right <= observation.width + 1);
    assert.deepEqual(observation.basic_outside, [], `Basic controls must fit viewport ${observation.width}`);
  }
  const advanced = ui.page.getByTestId('step4-advanced-recording');
  const summary = advanced.locator('summary');
  const beforeKeyboard = (await advanced.getAttribute('open')) !== null;
  await summary.focus();
  await summary.press('Enter');
  const afterKeyboard = (await advanced.getAttribute('open')) !== null;
  assert.notEqual(afterKeyboard, beforeKeyboard, 'advanced details must toggle with keyboard');
  await summary.press('Enter');
  result.keyboard = { advanced_summary_enter_toggled: true, restored: ((await advanced.getAttribute('open')) !== null) === beforeKeyboard };
  result.passed = true;
} catch (error) {
  result.failed_stage = stage;
  result.error = String(error.message).split('\n')[0];
  if (ui) result.ui_text = (await ui.page.locator('body').innerText().catch(() => '')).slice(-10_000);
} finally {
  try { await ui?.browser.close(); result.browser_cleanup = 'closed'; }
  catch { result.browser_cleanup = 'failed'; result.passed = false; }
  result.child_exit = await stopProcesses(children);
  result.child_exit_passed = result.child_exit.every((child) => child.state === 'exited' && child.exit_code === 0 &&
    !child.signal && !child.spawn_error && !child.stop_error && !child.kill_sent);
  if (!result.child_exit_passed) result.passed = false;
  try {
    if (createdSchema) pgCommand(`DROP SCHEMA ${quotePGIdentifier(schema)} CASCADE`);
    result.schema_cleanup = createdSchema ? !pgCommand(`SELECT nspname FROM pg_namespace WHERE nspname='${schema}'`) : 'not-created';
    if (result.schema_cleanup === false) result.passed = false;
  } catch { result.schema_cleanup = false; result.passed = false; }
  if (run) {
    const workCleanup = removeOwnedWork(run.work);
    result.owned_work_cleanup = workCleanup.removed;
    if (workCleanup.error) result.owned_work_cleanup_error = workCleanup.error;
    if (!result.owned_work_cleanup) result.passed = false;
  }
  // Keep independent observations intact while bounding generated evidence line count.
  const evidence = sanitize(result);
  const fields = Object.entries(evidence).map(([key, value]) => `  ${JSON.stringify(key)}: ${JSON.stringify(value)}`);
  writeFileSync(join(DEFAULT_EVIDENCE_DIR, `${runID}.json`), `{\n${fields.join(',\n')}\n}\n`);
  console.log(JSON.stringify({ run_id: runID, passed: result.passed, failed_stage: result.failed_stage, error: sanitize(result.error) }));
}
process.exitCode = result.passed ? 0 : 1;
