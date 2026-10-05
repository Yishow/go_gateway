// Final rendered UI check for the sidebar-only review fix. It deliberately
// does not repeat the unchanged eight-point bucket/recovery acceptance.
import assert from 'node:assert/strict';
import { existsSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, assertFreshPortFree, binaryIdentity,
  closeFreshUI, configureFreshSQLite, ensureFreshSidebarCollapsed, isFreshCleanExit,
  makeFreshRun, openFreshUI, persistFreshMappings, setupFreshDevice, setupFreshRules,
  startFreshGateway, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import { writeFreshResult } from './fresh-result.mjs';

const runId = process.env.F_RUN_ID ?? 'basic-integrity-20261005-ui-final';
const port = Number(process.env.GW_PORT ?? 3491);
const simPort = Number(process.env.F_SIM_PORT_A ?? 15192);
const run = makeFreshRun(runId);
const children = [];
let ui;
let result = { run_id: runId, passed: false, failures: [], revalidation_boundary: 'SummaryRail and connector copy only; prior c proof covers unchanged default mapping, schema/start, SQL buckets, runtime, recovery and disable/restart' };
try {
  assert.notEqual(port, simPort);
  await assertFreshPortFree(port);
  await assertFreshPortFree(simPort);
  children.push(startFreshSimulator({ work: run.work, port: simPort, registers: { 0: 215 } }));
  children.push(startFreshGateway({ work: run.work, port }));
  await waitForFreshGateway(`http://127.0.0.1:${port}`);
  ui = await openFreshUI({ port });
  const { page, base } = ui;
  await setupFreshDevice(page, { name: 'UI summary device', port: simPort });
  await setupFreshRules(page, [{ name: 'UI summary point', start: 40001, count: 1, dataType: 'int16', prefix: 'UI' }]);
  await persistFreshMappings(page, base, { tagPrefix: 'summary' });
  const config = await configureFreshSQLite(page, base, { path: run.destination_db });
  await ensureFreshSidebarCollapsed(page);
  const rail = page.getByTestId('summary-rail');
  await rail.waitFor();
  const text = await rail.innerText();
  assert.ok(text.includes('F fresh SQLite'), 'actual sidebar lacks connection name');
  assert.ok(text.includes('sqlite') && text.includes('Basic') && text.includes('Advanced'));
  assert.ok(!text.includes('sensor_readings') && !text.includes('5s') && !text.includes('資料庫寫入'), 'legacy effective recording summary remains');
  assert.equal(await page.getByTestId('basic-recording-interval').inputValue(), '60');
  assert.equal(await page.getByTestId('input-table-name').count(), 0);
  // Read the actual current-locale subtitle; no configuration endpoint is stubbed.
  const subtitle = page.getByText(/^(Configure database connection parameters|設定資料庫連線參數)$/);
  await subtitle.waitFor();
  const connectorText = await subtitle.innerText();
  assert.ok(/Configure database connection parameters|設定資料庫連線參數/.test(connectorText));
  assert.ok(!/write strategy parameters|連線參數與策略/.test(connectorText));
  assert.equal(existsSync(run.destination_db), false, 'connection-only UI check implicitly created destination');
  const screenshot = `${runId}-1440.png`;
  await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, screenshot), fullPage: true });
  assert.deepEqual(ui.pageErrors, []);
  result = { ...result, passed: true, binary: binaryIdentity(DEFAULT_GATEWAY_BINARY), sidebar_text: text,
    connector_text: connectorText, config, basic_interval_seconds: 60, target_created: false,
    screenshot, page_errors: ui.pageErrors, api_traffic: ui.apiTraffic };
} catch (error) {
  result.failures.push(error?.message ?? String(error));
  try { await ui?.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, `${runId}-error.png`), fullPage: true }); } catch { /* failure evidence is best effort */ }
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
  writeFreshResult(join(DEFAULT_EVIDENCE_DIR, `${runId}.json`), result);
  console.log(`${result.passed ? 'PASS' : 'FAIL'} ${runId}: ${result.failures.join('; ')}`);
}
