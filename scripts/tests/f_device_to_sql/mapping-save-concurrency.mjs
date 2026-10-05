// User-case regression through the real isolated embedded UI and SQL repositories.
// Configuration writes are browser actions; fault injection only touches this owned DB.
import assert from 'node:assert/strict';
import { existsSync, readFileSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, assertFreshPortFree, binaryIdentity,
  closeFreshUI, ensureFreshSidebarCollapsed, isFreshCleanExit, makeFreshRun, observeJSON,
  openFreshUI, persistFreshMappings, setupFreshDevice, setupFreshRules,
  startFreshGateway, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import { sh } from './lib.mjs';
import { quoteSQLite, readSQLiteRows } from './fresh-sqlite-observation.mjs';
import { writeFreshResult } from './fresh-result.mjs';

const runId = process.env.F_RUN_ID ?? 'mapping-concurrency-20261005a';
const port = Number(process.env.GW_PORT ?? 3492);
const simPort = Number(process.env.F_SIM_PORT_A ?? 15193);
assert.notEqual(port, 3222);
assert.notEqual(port, 5173);
assert.notEqual(simPort, port);
const run = makeFreshRun(runId);
const children = [];
let ui;
const detailedTraffic = [];
const replies = new Set();
const result = { run_id: runId, passed: false, failures: [], phases: [] };
const sql = `SELECT l.address,l.point_id,l.mapping_id,l.tag_id,m.status,m.enabled,
 m.transform_pipeline,m.proposed_signature,m.last_applied_signature,m.rule_candidate_id,
 t.key,t.display_name,t.data_type FROM source_rule_links l JOIN mappings m ON m.id=l.mapping_id
 JOIN tags t ON t.id=l.tag_id ORDER BY l.address`;

async function until(read, predicate, timeoutMs = 30_000) {
  const deadline = Date.now() + timeoutMs;
  let last;
  while (Date.now() < deadline) {
    last = await read();
    if (predicate(last)) return last;
    await new Promise((resolve) => setTimeout(resolve, 250));
  }
  throw new Error(`UI mapping state did not settle: ${JSON.stringify(last).slice(0, 500)}`);
}
async function saved(page, predicate) {
  await until(() => observeJSON(ui.base, '/studio-v2/workspace/mappings'), predicate);
  await page.waitForFunction(() => {
    const states = [...document.querySelectorAll('[data-testid^="mapping-save-state-"]')];
    return states.length === 8 && states.every((row) => /Saved|已儲存/.test(row.textContent ?? ''));
  }, null, { timeout: 30_000 });
}
function snapshot(name) {
  const rows = readSQLiteRows(run.gateway_db, sql);
  assert.equal(rows.length, 8);
  assert.equal(new Set(rows.map((row) => row.point_id)).size, 8);
  assert.equal(new Set(rows.map((row) => row.mapping_id)).size, 8);
  assert.equal(new Set(rows.map((row) => row.tag_id)).size, 8);
  assert.ok(rows.every((row) => row.rule_candidate_id && row.last_applied_signature && row.proposed_signature));
  result.phases.push({ name, rows });
  return rows;
}

try {
  await assertFreshPortFree(port);
  await assertFreshPortFree(simPort);
  children.push(startFreshSimulator({ work: run.work, port: simPort,
    registers: { 0: 243, 1: 2, 2: 3, 3: 4, 4: 5, 5: 6, 6: 7, 7: 8 } }));
  children.push(startFreshGateway({ work: run.work, port }));
  await waitForFreshGateway(`http://127.0.0.1:${port}`);
  ui = await openFreshUI({ port, width: 1440, height: 1100 });
  const { page } = ui;
  page.on('response', (reply) => {
    const path = new URL(reply.url()).pathname;
    const method = reply.request().method();
    if (!/workspace\/mappings|mappings\/preview/.test(path) || !['PUT', 'POST'].includes(method)) return;
    const task = (async () => {
      detailedTraffic.push({ path, method, status: reply.status(),
        request: reply.request().postDataJSON(), response: await reply.json().catch(() => null) });
    })();
    replies.add(task);
    task.finally(() => replies.delete(task));
  });
  await setupFreshDevice(page, { name: 'Mapping regression', port: simPort });
  await setupFreshRules(page, [{ name: 'Eight same-rule points', start: 40001, count: 8,
    dataType: 'int16', prefix: 'BLOCK' }]);
  await persistFreshMappings(page, ui.base, { tagPrefix: 'race' });
  await page.getByTestId('step-nav-button-3').click();
  await page.getByTestId('step3-mapping-container').waitFor();
  await saved(page, (rows) => rows.length === 8 && rows.every((row) => row.target_type === 'int16'));
  const before = snapshot('default-int16');
  const identities = before.map(({ address, point_id, mapping_id, tag_id }) => ({ address, point_id, mapping_id, tag_id }));
  for (let batch = 0; batch < 3; batch += 1) {
    await page.getByTestId('btn-enable-all-mappings').click();
    await page.getByTestId('btn-enable-all-mappings').click();
    await saved(page, (rows) => rows.length === 8 && rows.every((row) => row.enabled));
    const after = snapshot(`enable-batch-${batch + 1}`);
    assert.deepEqual(after.map(({ address, point_id, mapping_id, tag_id }) => ({ address, point_id, mapping_id, tag_id })), identities);
    assert.ok(after.every((row) => row.transform_pipeline === '[]'));
  }
  const rows = page.locator('[data-testid^="mapping-row-"]');
  const firstId = (await rows.first().getAttribute('data-testid')).replace('mapping-row-', '');
  await rows.first().locator('td').first().click();
  await page.getByTestId(`select-target-type-${firstId}`).selectOption('float64');
  await saved(page, (records) => records.some((row) => row.address === '40001' && row.target_type === 'float64'));
  await page.getByTestId('btn-bulk-apply-target-type').click();
  await saved(page, (records) => records.length === 8 && records.every((row) => row.target_type === 'float64'));
  await page.getByTestId(`input-scale-${firstId}`).fill('0.5');
  await page.getByTestId(`input-offset-${firstId}`).fill('10');
  await page.getByTestId('btn-bulk-apply-all').click();
  await saved(page, (records) => records.length === 8 && records.every((row) => row.target_type === 'float64' && row.scale === 0.5 && row.offset === 10));
  const name = page.getByTestId(`input-display-name-${firstId}`);
  for (const text of ['First edit', 'Second edit', 'Final rapid edit']) await name.fill(text);
  await saved(page, (records) => records.some((row) => row.address === '40001' && row.display_name === 'Final rapid edit'));
  assert.equal(await name.inputValue(), 'Final rapid edit');
  const previewText = await until(() => page.getByTestId(`preview-final-${firstId}`).innerText(),
    (text) => Number(text.split('\n', 1)[0].trim()) === 131.5);
  result.scaled_preview = { display_text: previewText, value: Number(previewText.split('\n', 1)[0].trim()) };
  const scaled = snapshot('bulk-scaled-and-latest-draft');
  assert.ok(scaled.every((row) => row.last_applied_signature !== row.proposed_signature));
  const persistedFirst = scaled.find((row) => row.address === '40001');
  assert.equal(persistedFirst.display_name, 'Final rapid edit');

  // Failed UPDATE occurs before tag/pipeline mutation; late rollback has separate SQL tests.
  sh('sqlite3', ['-cmd', '.timeout 15000', run.gateway_db,
    `CREATE TRIGGER owned_mapping_failure BEFORE UPDATE ON mappings WHEN OLD.id=${quoteSQLite(persistedFirst.mapping_id)} BEGIN SELECT RAISE(ABORT,'owned failure'); END;`]);
  await name.fill('Recovered draft');
  await page.getByTestId(`mapping-save-state-${firstId}`).waitFor();
  await until(() => page.getByTestId(`mapping-save-state-${firstId}`).innerText(),
    (text) => /500/.test(text) && /workspace_mapping_save_failed/.test(text));
  const failedText = await page.getByTestId(`mapping-save-state-${firstId}`).innerText();
  assert.ok(!/constraint|owned failure|SQLITE|RAISE|trigger/i.test(failedText));
  assert.ok(!/owned failure|owned_mapping_failure|RAISE\(ABORT|constraint failed/i.test(await page.locator('body').innerText()));
  assert.equal(await name.inputValue(), 'Recovered draft');
  result.error_screenshot = `${runId}-error-state-1440.png`;
  await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, result.error_screenshot), fullPage: true });
  const failuresBeforeRetry = detailedTraffic.filter((item) => item.status >= 400).length;
  await page.waitForTimeout(500);
  assert.equal(detailedTraffic.filter((item) => item.status >= 400).length, failuresBeforeRetry, 'failed save replayed without operator retry');
  sh('sqlite3', ['-cmd', '.timeout 15000', run.gateway_db, 'DROP TRIGGER owned_mapping_failure;']);
  await page.getByTestId(`mapping-save-retry-${firstId}`).click();
  await saved(page, (records) => records.some((row) => row.address === '40001' && row.display_name === 'Recovered draft'));
  const recovered = snapshot('explicit-retry');
  assert.equal(recovered.find((row) => row.address === '40001').display_name, 'Recovered draft');
  await ensureFreshSidebarCollapsed(page);
  const screenshot = `${runId}-1440.png`;
  await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, screenshot), fullPage: true });
  await Promise.all([...replies]);
  const errors = detailedTraffic.filter((item) => item.status >= 400);
  assert.equal(errors.length, 1, JSON.stringify(errors));
  assert.equal(errors[0].status, 500);
  assert.equal(errors[0].response.error.code, 'workspace_mapping_save_failed');
  assert.deepEqual(ui.pageErrors, []);
  Object.assign(result, { passed: true, screenshot, binary: binaryIdentity(DEFAULT_GATEWAY_BINARY),
    page_errors: ui.pageErrors, failed_safe_text: failedText, injected_error_count: 1,
    mapping_sql_observation: 'Owned gateway configuration SQLite; destination recording is unchanged and covered separately',
    api_traffic: detailedTraffic, source: JSON.parse(readFileSync(process.env.F_SOURCE_MANIFEST, 'utf8')) });
} catch (error) {
  result.failures.push(error?.message ?? String(error));
  try { await ui?.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, `${runId}-error.png`), fullPage: true }); } catch { /* best effort failure evidence */ }
} finally {
  let browser = { state: 'closed' };
  try { await closeFreshUI(ui?.browser); } catch (error) { browser = { state: 'error', error: error.message }; }
  const processes = await stopProcesses(children).catch((error) => [{ state: 'error', error: error.message }]);
  const namespace = { path: run.work, absent: false };
  if (browser.state === 'closed' && processes.every(isFreshCleanExit)) {
    try { rmSync(run.work, { recursive: true }); namespace.absent = !existsSync(run.work); }
    catch (error) { namespace.error = error.message; }
  }
  result.cleanup = { browser, processes, namespace };
  if (!namespace.absent) { result.passed = false; result.failures.push('owned cleanup not proved'); }
  await Promise.allSettled([...replies]);
  result.api_traffic = detailedTraffic;
  writeFreshResult(join(DEFAULT_EVIDENCE_DIR, `${runId}.json`), result);
  console.log(`${result.passed ? 'PASS' : 'FAIL'} ${runId}: ${result.failures.join('; ')}`);
  if (!result.passed) process.exitCode = 1;
}
