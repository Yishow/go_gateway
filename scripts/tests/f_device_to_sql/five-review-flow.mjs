// Real embedded UI → MC 3E → disposable SQLite contract for the five review fixes.
import assert from 'node:assert/strict';
import { mkdirSync, mkdtempSync, writeFileSync } from 'node:fs';
import { join, resolve } from 'node:path';
import * as f from './fresh-ui.mjs';
import { readSQLiteRows } from './fresh-sqlite-observation.mjs';
import { startMCFixture } from './mc-fixture.mjs';

const evidence = process.env.F_EVIDENCE_DIR;
assert.ok(evidence && resolve(evidence) === evidence, 'explicit absolute F_EVIDENCE_DIR required');
mkdirSync(evidence, { recursive: true });
const work = mkdtempSync('/tmp/gateway-five-fixes-');
const port = Number(process.env.GW_PORT ?? 3599), mcPort = Number(process.env.F_MC_PORT ?? 15599);
const binary = process.env.F_GATEWAY_BINARY;
assert.ok(binary, 'explicit current gateway binary required');
const result = { work, binary: f.binaryIdentity(binary), checks: {}, failures: [] };
const children = [];
let ui, mc;
const q = value => `"${String(value).replaceAll('"', '""')}"`;
const list = data => Array.isArray(data) ? data : data.rules ?? data.mappings ?? [];
const sleep = ms => new Promise(resolveWait => setTimeout(resolveWait, ms));
async function until(read, predicate, timeout = 30000) {
  const deadline = Date.now() + timeout;
  let last;
  while (Date.now() < deadline) { last = await read(); if (predicate(last)) return last; await sleep(400); }
  throw new Error(`observation timeout: ${JSON.stringify(last).slice(0, 300)}`);
}
async function snapshot(name) {
  writeFileSync(join(evidence, `${name}.txt`), await ui.page.locator('body').innerText());
  await ui.page.screenshot({ path: join(evidence, `${name}.png`), fullPage: true });
}
try {
  await Promise.all([f.assertFreshPortFree(port), f.assertFreshPortFree(mcPort)]);
  mc = await startMCFixture(mcPort, { 0: 253, 1: 1234, 2: 55, 3: 34464, 4: 1, 5: 1, 6: 0, 7: 21 });
  children.push(f.startFreshGateway({ work, binary, port }));
  await f.waitForFreshGateway(`http://127.0.0.1:${port}`);
  ui = await f.openFreshUI({ port });
  const { page, base } = ui;
  await page.evaluate(() => localStorage.setItem('i18nextLng', 'zh-TW'));
  await page.reload();
  await page.getByTestId('btn-add-device').click();
  await page.getByTestId('editor-input-name').fill('MC mixed fixture');
  await page.getByTestId('protocol-card-mc_3e').click();
  await page.getByTestId('input-host').fill('127.0.0.1');
  await page.getByTestId('input-port').fill(String(mcPort));
  await page.waitForFunction(() => document.querySelector('[data-testid=device-save-banner]')?.getAttribute('data-save-state') === 'saved');
  await page.getByTestId('run-test-button').click();
  await page.getByTestId('success-readiness-card').waitFor();
  await page.getByTestId('btn-continue-step1').click();
  const specs = [
    ['D0', 'int16', '溫度', '°C', '0.1', 'float64'],
    ['D1', 'uint16', '壓力', 'bar', '0.01', 'float64'],
    ['D2', 'uint16', '瞬時流量', 'm3/h', '1', 'uint16'],
    ['D3', 'uint32', '累積流量', 'm3', '1', 'uint32'],
    ['D5', 'uint16', '狀態', '', '1', 'uint16'],
    ['D6', 'uint16', '警報', '', '1', 'int16'],
    ['D7', 'uint16', '累計產量', 'pcs', '1', 'uint16'],
  ];
  for (const [index, [address, type]] of specs.entries()) {
    const oldTabs = await page.locator('[data-testid^="rule-tab-rule-"]').count();
    await page.getByTestId('rule-add-btn').click();
    await until(() => page.locator('[data-testid^="rule-tab-rule-"]').count(), count => count === oldTabs + 1);
    const ruleId = (await page.locator('[data-testid^="rule-tab-rule-"]').last().getAttribute('data-testid')).replace('rule-tab-', '');
    await page.getByTestId('rule-count-input').fill('1');
    await page.getByTestId('rule-start-input').fill(address);
    await f.waitForObservation(base, '/studio-v2/workspace/source-rules', data => list(data).some(rule => rule.id === ruleId && rule.start_address === address && rule.count === 1));
    await page.getByTestId('rule-datatype-select').selectOption(type);
    await f.waitForObservation(base, '/studio-v2/workspace/source-rules', data => list(data).some(rule => rule.id === ruleId && rule.data_type === type));
    await page.getByTestId('rule-prefix-input').fill(`M${index}_`);
    await page.getByTestId('rule-prefix-input').press('Tab');
    await f.waitForObservation(base, '/studio-v2/workspace/source-rules', data => list(data).some(rule => rule.id === ruleId && rule.naming_prefix === `M${index}_`));
    await sleep(200);
  }
  result.checks.rules = await f.observeJSON(base, '/studio-v2/workspace/source-rules');
  assert.equal(list(result.checks.rules).length, 7);
  await page.getByTestId('continue-step3-btn').click();
  await page.getByTestId('mapping-advanced-toggle').click();
  const tags = [];
  for (const [address, , name, unit, scale, targetType] of specs) {
    const rows = page.locator('[data-testid^="mapping-row-"]');
    const row = rows.filter({ hasText: address }).first();
    const id = (await row.getAttribute('data-testid')).replace('mapping-row-', '');
    for (const [testId, field, expected] of [[`input-display-name-${id}`, 'display_name', name], [`input-unit-${id}`, 'unit', unit]]) {
      await page.getByTestId(testId).fill(expected); await page.getByTestId(testId).press('Tab');
      await f.waitForObservation(base, '/studio-v2/workspace/mappings', data => list(data).some(m => m.address === address && m[field] === expected));
      await sleep(150);
    }
    await page.getByTestId(`select-target-type-${id}`).selectOption(targetType);
    await f.waitForObservation(base, '/studio-v2/workspace/mappings', data => list(data).some(m => m.address === address && m.target_type === targetType));
    await page.getByTestId(`input-scale-${id}`).fill(scale); await page.getByTestId(`input-scale-${id}`).press('Tab');
    await f.waitForObservation(base, '/studio-v2/workspace/mappings', data => list(data).some(m => m.address === address && m.scale === Number(scale)));
    tags.push({ address, name, unit, scale, targetType });
  }
  await page.getByTestId('mapping-advanced-toggle').click();
  await snapshot('mc-mapping');
  // Real API preview stays unconfirmed; only the explicit UI button adds consent.
  await page.getByTestId('measurement-advanced').locator('summary').click();
  await page.getByTestId('template-select').selectOption('mixed_equipment_sensors_v1');
  const previewResponse = page.waitForResponse(r => r.url().endsWith('/templates/preview'));
  await page.getByTestId('preview-template-btn').click();
  const preview = (await (await previewResponse).json()).data;
  assert.equal(preview.confirmed, false);
  const rejected = await page.request.post(`${base}/api/v1/datalink/studio-v2/workspace/measurements/templates/apply`, { data: preview });
  assert.equal(rejected.status(), 422);
  const appliedResponse = page.waitForResponse(r => r.url().endsWith('/templates/apply') && r.request().method() === 'POST');
  await page.getByTestId('apply-template-btn').click();
  const applied = await appliedResponse;
  assert.equal(applied.request().postDataJSON().confirmed, true);
  assert.equal(applied.status(), 200, await applied.text());
  result.checks.template = { preview_confirmed: false, unconfirmed_status: rejected.status(), confirmed_status: applied.status() };
  await page.getByTestId('btn-continue').click();
  await f.configureFreshSQLite(page, base, { path: join(work, 'destination.db') });
  await page.getByTestId('basic-recording-interval').fill('2');
  const prepared = await f.prepareFreshBasicGroup(page, base);
  const ddlResponse = page.waitForResponse(r => r.url().endsWith('/schema-apply') && r.request().method() === 'POST');
  await page.getByTestId('group-schema-apply').click();
  const ddl = await ddlResponse;
  assert.equal(ddl.status(), 200);
  result.checks.ddl = (await ddl.json()).data;
  assert.equal(result.checks.ddl.status, 'succeeded');
  const startResponse = page.waitForResponse(r => r.url().endsWith('/recording-start') && r.request().method() === 'POST');
  await page.getByTestId('basic-recording-start').click();
  const started = await startResponse;
  assert.equal(started.status(), 200);
  const operation = (await started.json()).data;
  result.checks.start = await f.waitForObservation(base, `/studio-v2/workspace/recording-start/operations/${operation.operation_id}`, data => data.status === 'succeeded');
  const mappings = list(await f.observeJSON(base, '/studio-v2/workspace/mappings'));
  const group = prepared.group;
  const sql = () => readSQLiteRows(join(work, 'destination.db'), `SELECT * FROM ${q(group.destination.table_name)} ORDER BY bucket_start`);
  const rows = await until(sql, values => values.length >= 3);
  const expected = { D0: 25.3, D1: 12.34, D2: 55, D3: 100000, D5: 1, D6: 0, D7: 21 };
  const columns = {};
  for (const [address, value] of Object.entries(expected)) {
    const mapping = mappings.find(m => m.address === address);
    const column = group.members.find(m => m.point_id === mapping.point_id && m.tag_id === mapping.tag_id).target_column;
    columns[address] = column;
    for (const row of rows.slice(0, 3)) {
      assert.notEqual(row[column], null, `${address} must contain a value`);
      if (['D0', 'D1', 'D6'].includes(address)) {
        assert.equal(typeof row[column], 'number');
        assert.ok(Math.abs(row[column] - value) < 1e-9, `${address}: ${row[column]} != ${value}`);
      } else assert.equal(row[column], String(value), `${address} retains exact unsigned TEXT`);
    }
  }
  assert.ok(mc.requests.some(r => r.address === 3 && r.count === 2));
  assert.ok(!mc.requests.some(r => r.address === 3 && r.count === 4));
  result.checks.mc_three_buckets = { tags, rows: rows.slice(0, 3), columns, mappings, group };
  await snapshot('mc-sql-confirmed');
  // A valid uint16 device value beyond the deliberately selected int16 target
  // produces one typed mapping failure without breaking the other MC reads.
  mc.words.set(6, 65535);
  await sleep(7500);
  const requiredBuckets = readSQLiteRows(join(work, 'gateway.db'), `SELECT kind,reason,members FROM wg_delivery_buckets WHERE group_id='${group.id}' ORDER BY bucket_start DESC LIMIT 3`);
  assert.ok(requiredBuckets.some(bucket => bucket.kind === 'skipped' && bucket.reason === 'incomplete-required-member'));
  result.checks.required_gap = requiredBuckets;
  mc.words.set(6, 0);
  await sleep(3500);
  // Narrow Basic entry opens this group's existing required controls.
  await page.getByTestId('basic-recording-completeness').selectOption('partial');
  await page.getByTestId('basic-member-settings').click();
  await page.getByTestId('group-editor').waitFor();
  const alarmMapping = mappings.find(m => m.address === 'D6');
  const alarmRow = page.getByTestId(`group-member-row-${alarmMapping.point_id}`);
  await alarmRow.getByRole('checkbox', { name: /必要值|Required/i }).uncheck();
  await page.getByTestId('group-advanced').locator('summary').click();
  await page.getByTestId('group-incomplete').selectOption('partial');
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor();
  const editor = page.getByTestId('group-editor');
  await editor.getByTestId('group-schema-preview').click();
  await editor.getByTestId('group-schema-preview-result').waitFor();
  const updatedDDLResponse = page.waitForResponse(r => r.url().endsWith('/schema-apply') && r.request().method() === 'POST');
  await editor.getByTestId('group-schema-apply').click();
  const updatedDDL = await updatedDDLResponse;
  assert.equal(updatedDDL.status(), 200);
  assert.equal((await updatedDDL.json()).data.status, 'succeeded');
  await page.getByTestId('group-apply').click();
  await f.waitForObservation(base, '/studio-v2/workspace/write-groups', d => d.groups.some(g => g.id === group.id && g.applied_revision === g.revision && g.row_policy.incomplete_policy === 'partial' && g.members.some(m => m.point_id === alarmMapping.point_id && !m.required)));
  const before = sql().length;
  mc.words.set(6, 65535);
  const partialRows = await until(sql, values => values.slice(before).some(row => row[columns.D6] === null));
  const partial = partialRows.find(row => row[columns.D6] === null);
  assert.ok(JSON.parse(partial.provenance).some(m => m.status !== 'ok'));
  result.checks.optional_partial = partial;
  mc.words.set(6, 0);
  const recovered = await until(sql, values => values.some(row => row.bucket_start > partial.bucket_start && String(row[columns.D6]) === '0'));
  result.checks.recovery_rows = recovered.length;
  result.checks.recovered_row = recovered.find(row => row.bucket_start > partial.bucket_start && row[columns.D6] === 0);
  await snapshot('optional-recovered');
  await page.goto(`${base}/studio/runtime`);
  await page.getByTestId('runtime-dashboard-header').waitFor();
  assert.ok(!(await page.getByTestId('runtime-dashboard-header').innerText()).includes('即時寫入資料庫'));
  await page.getByTestId('runtime-header-recording-evidence').click();
  await page.getByTestId('basic-recording-panel').waitFor();
  result.checks.evidence_return_step4 = true;
  assert.deepEqual(ui.pageErrors, []);
  result.page_errors = ui.pageErrors;
} catch (error) {
  result.failures.push(error.stack); console.error(error.stack);
  if (ui) await snapshot('failure').catch(() => {});
} finally {
  if (ui) await f.closeFreshUI(ui.browser);
  result.cleanup = await f.stopProcesses(children);
  if (mc) { result.mc_requests = mc.requests; await mc.close(); }
  writeFileSync(join(evidence, 'result.json'), JSON.stringify(result, null, 2));
}
assert.equal(result.failures.length, 0, result.failures.join('\n'));
assert.ok(result.cleanup.every(f.isFreshCleanExit));
console.log(JSON.stringify({ passed: true, checks: Object.keys(result.checks) }));
