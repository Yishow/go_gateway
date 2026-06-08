const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1600, height: 1100 } });

  const take = async (name) => {
    const path = `/tmp/${name}.png`;
    await page.screenshot({ path, fullPage: true });
    console.log('screenshot=', path);
  };

  const waitMs = (ms) => new Promise((r) => setTimeout(r, ms));

  try {
    console.log('goto studio/v2');
    await page.goto('http://localhost:4180/studio/v2', { waitUntil: 'domcontentloaded', timeout: 20000 });
    await page.getByTestId('step1-device-view').waitFor({ timeout: 20000 });

    const deviceLabel = await page.locator('[data-testid^="status-tested-"]').first().textContent().catch(() => '');
    console.log('initial-step1-status-testchip=', deviceLabel);

    const runTest = page.getByTestId('run-test-button');
    await runTest.waitFor({ timeout: 20000 });
    await runTest.click();

    const successCard = page.getByTestId('success-readiness-card');
    await successCard.waitFor({ timeout: 30000 });

    const continueStep1 = page.getByTestId('btn-continue-step1');
    await continueStep1.waitFor({ timeout: 20000 });
    if (await continueStep1.isDisabled()) {
      await take('step1_disabled_after_test');
      throw new Error('step1 continue still disabled after run test');
    }
    await continueStep1.click();

    await page.getByTestId('step2-rule-container').waitFor({ timeout: 20000 });
    const continueStep2 = page.getByTestId('continue-step3-btn');
    await continueStep2.waitFor({ timeout: 20000 });
    if (await continueStep2.isDisabled()) {
      const statusText = await page.getByTestId('step2-rule-container').textContent().catch(() => '');
      await take('step2_disabled_after_page');
      throw new Error(`step2 continue disabled: ${statusText?.slice(0, 200)}`);
    }
    await continueStep2.click();

    await page.getByTestId('step3-mapping-container').waitFor({ timeout: 20000 });
    const continueStep3 = page.getByTestId('btn-continue');
    await continueStep3.waitFor({ timeout: 20000 });
    if (await continueStep3.isDisabled()) {
      const rowText = await page.locator('[data-testid^="mapping-row-"]').first().textContent().catch(() => '');
      await take('step3_disabled_after_page');
      throw new Error(`step3 continue disabled: ${rowText?.slice(0, 200)}`);
    }
    await continueStep3.click();

    await page.getByText('第一次啟動設備').waitFor({ timeout: 20000 });
    const activateBtn = page.getByRole('button', { name: '第一次啟動設備' });
    await activateBtn.waitFor({ timeout: 20000 });
    await activateBtn.click();

    const toRuntimeBtn = page.getByRole('button', { name: '前往 Runtime Dashboard' });
    await toRuntimeBtn.waitFor({ timeout: 30000 });
    await take('step4_pending');
    await toRuntimeBtn.click();

    await page.waitForURL('**/runtime-dashboard', { timeout: 30000 });
    await page.getByTestId('runtime-dashboard-route').waitFor({ timeout: 30000 });

    await page.waitForFunction(() => {
      const streamEl = document.querySelector('[data-testid="runtime-dashboard-stream-state"]');
      if (!streamEl) return false;
      const state = (streamEl.textContent || '').trim();
      return ['live', 'ready', 'degraded', 'empty', 'loading'].includes(state);
    }, { timeout: 30000 });

    const streamState = await page.locator('[data-testid="runtime-dashboard-stream-state"]').textContent().catch(() => '');
    const routeState = await page.locator('[data-testid="runtime-dashboard-route-state"]').textContent().catch(() => '');
    const hasLiveTable = await page.getByTestId('runtime-dashboard-live-points-table').isVisible().catch(() => false);
    const hasPointRows = await page.locator('[data-testid="runtime-dashboard-live-points-table"] tbody tr').first().isVisible().catch(() => false);

    await take('runtime-dashboard');

    const rawValue = await page.locator('[data-testid="runtime-dashboard-live-points-table"] tbody tr .px-3').first().textContent().catch(() => null);

    console.log('RESULT: PASS');
    console.log('streamState=', (streamState || '').trim());
    console.log('routeState=', (routeState || '').trim());
    console.log('liveTable=', hasLiveTable);
    console.log('pointRow=', hasPointRows);
    console.log('firstCell=', rawValue);

    await browser.close();
  } catch (err) {
    try { await take('failure'); } catch (_) {}
    console.log('RESULT: FAIL');
    console.error(String(err));
    await browser.close();
    process.exit(1);
  }
})();
