// Browser transport-fault injection against real template endpoints in an owned fixture.
import assert from 'node:assert/strict';
import { mkdirSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import * as f from './fresh-ui.mjs';
const work = process.env.F_EXISTING_WORK, evidence = process.env.F_EVIDENCE_DIR;
assert.ok(work?.startsWith('/tmp/gateway-five-fixes-') && evidence, 'owned fixture work and evidence required');
mkdirSync(evidence, { recursive: true });
const children = [], result = { checks: [], failures: [] };
let ui;
try {
  await f.assertFreshPortFree(3600);
  children.push(f.startFreshGateway({ work, binary: process.env.F_GATEWAY_BINARY, port: 3600 }));
  await f.waitForFreshGateway('http://127.0.0.1:3600');
  ui = await f.openFreshUI({ port: 3600 });
  const { page } = ui;
  const templatePattern = '**/measurements/templates';
  const fail = status => route => route.fulfill({ status, contentType: 'application/json', body: JSON.stringify({ success: false, error: { code: 'TEST_FAILURE', message: 'private password SHOULD_NOT_RENDER' } }) });
  await page.route(templatePattern, fail(503));
  await page.getByTestId('step-nav-button-3').click();
  await page.getByTestId('measurement-advanced').locator('summary').click();
  await page.getByTestId('template-load-error').waitFor();
  assert.ok(!(await page.locator('body').innerText()).includes('SHOULD_NOT_RENDER'));
  await page.unroute(templatePattern);
  await page.getByTestId('template-load-retry').click();
  await page.getByTestId('template-select').selectOption('three_phase_power_meter_v1');
  result.checks.push('load_503_safe_reload');
  for (const action of ['preview', 'apply']) {
    for (const kind of [409, 500, 'offline']) {
      if (action === 'apply') {
        await page.getByTestId('preview-template-btn').click();
        await page.getByTestId('template-preview-summary').waitFor();
      }
      const pattern = `**/templates/${action}`;
      if (kind === 'offline') await page.context().setOffline(true);
      else await page.route(pattern, fail(kind));
      await page.getByTestId(action === 'preview' ? 'preview-template-btn' : 'apply-template-btn').click();
      await page.getByTestId('template-error').waitFor({ timeout: 20000 });
      assert.ok(!(await page.locator('body').innerText()).includes('SHOULD_NOT_RENDER'));
      if (kind === 'offline') await page.context().setOffline(false);
      else await page.unroute(pattern);
      await page.getByTestId('preview-template-btn').click();
      await page.getByTestId('template-preview-summary').waitFor();
      result.checks.push(`${action}_${kind}_safe_repreview`);
    }
  }
  await page.screenshot({ path: join(evidence, 'template-recovered.png'), fullPage: true });
  assert.deepEqual(ui.pageErrors, []);
  result.page_errors = ui.pageErrors;
} catch (error) { result.failures.push(error.stack); console.error(error.stack); }
finally {
  if (ui) await f.closeFreshUI(ui.browser);
  result.cleanup = await f.stopProcesses(children);
  writeFileSync(join(evidence, 'result.json'), JSON.stringify(result, null, 2));
}
assert.equal(result.failures.length, 0);
assert.ok(result.cleanup.every(f.isFreshCleanExit));
console.log(JSON.stringify(result));
