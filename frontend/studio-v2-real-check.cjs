const { chromium } = require('@playwright/test');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1440, height: 1024 } });
  const base = 'http://127.0.0.1:4174';
  try {
    console.log('STEP1 goto studio/v2');
    await page.goto(`${base}/studio/v2`, { waitUntil: 'domcontentloaded', timeout: 30000 });
    console.log('URL@afterGoto=', page.url());

    const waitAny = await Promise.race([
      page.waitForSelector('[data-testid="step1-device-view"]', { timeout: 30000 }).then(() => 'step1'),
      page.waitForSelector('text=登入', { timeout: 30000 }).then(() => 'signin'),
      page.waitForSelector('text=404', { timeout: 30000 }).then(() => 'notfound'),
    ]);
    console.log('FIRST_MATCH=', waitAny);

    if (waitAny !== 'step1') {
      const html = await page.content();
      console.log('DOM_HEAD=', html.slice(0, 400).replace(/\s+/g, ' '));
      throw new Error('expected step1 view not found');
    }

    await page.fill('[data-testid="input-host"]', '127.0.0.1');
    await page.fill('[data-testid="input-port"]', '5020');
    await page.fill('[data-testid="input-slave-id"]', '1');
    await page.fill('[data-testid="input-timeout"]', '5');

    console.log('STEP2 run-test-button');
    await page.click('[data-testid="run-test-button"]');
    await page.waitForSelector('[data-testid="success-readiness-card"]', { timeout: 60000 });
    await page.waitForSelector('[data-testid="btn-continue-step1"]:not([disabled])', { timeout: 60000 });
    await page.click('[data-testid="btn-continue-step1"]');

    await page.waitForSelector('[data-testid="merged-point-table-container"]', { timeout: 60000 });
    await page.waitForSelector('[data-testid="continue-step3-btn"]:not([disabled])', { timeout: 60000 });
    await page.click('[data-testid="continue-step3-btn"]');

    await page.waitForSelector('[data-testid="step3-mapping-container"]', { timeout: 60000 });
    await page.waitForSelector('[data-testid="btn-continue"]:not([disabled])', { timeout: 60000 });
    await page.click('[data-testid="btn-continue"]');

    await page.waitForSelector('button:has-text("第一次啟動設備")', { timeout: 60000 });
    await page.click('button:has-text("第一次啟動設備")');
    const outcome = await Promise.race([
      page.waitForSelector('[data-testid^="activation-result-"]', { timeout: 120000 }).then(() => 'result'),
      page.waitForSelector('button:has-text("前往 Runtime Dashboard")', { timeout: 120000 }).then(() => 'goto'),
      page.waitForSelector('[data-testid="activation-empty-message"]', { timeout: 120000 }).then(() => 'empty'),
    ]);
    console.log('STEP4_OUTCOME=' + outcome);

    await page.waitForSelector('button:has-text("前往 Runtime Dashboard")', { timeout: 120000 });
    await page.click('button:has-text("前往 Runtime Dashboard")');

    await page.waitForSelector('[data-testid="runtime-dashboard-route"]', { timeout: 120000 });
    await page.waitForFunction(() => {
      const state = document.querySelector('[data-testid="runtime-dashboard-route-state"]');
      return Boolean(state && state.textContent && ['live', 'degraded'].includes(state.textContent.trim()));
    }, { timeout: 120000 });

    await page.waitForFunction(() => {
      const rows = document.querySelectorAll('[data-testid="runtime-dashboard-live-points-table"] tbody tr');
      return rows.length > 0;
    }, { timeout: 120000 });

    const routeState = await page.locator('[data-testid="runtime-dashboard-route-state"]').textContent();
    const rowCount = await page.locator('[data-testid="runtime-dashboard-live-points-table"] tbody tr').count();
    const selected = await page.locator('[data-testid="runtime-dashboard-selected-device"]').textContent();
    const stream = await page.locator('[data-testid="runtime-dashboard-stream-state"]').textContent();

    console.log('RUNTIME_ROUTE_STATE=' + (routeState || '').trim());
    console.log('RUNTIME_SELECTED_DEVICE=' + (selected || '').trim());
    console.log('RUNTIME_STREAM_STATE=' + (stream || '').trim());
    console.log('RUNTIME_ROWS=' + rowCount);
    await page.screenshot({ path: '/private/tmp/studio-v2-real-check-final.png', fullPage: true });

    const runtimeStatus = await page.evaluate(async () => {
      const r = await fetch('http://127.0.0.1:8080/api/v1/datalink/runtime/status');
      return await r.text();
    });
    console.log('RUNTIME_STATUS=' + runtimeStatus);
  } catch (err) {
    console.error('E2E_ERROR=' + (err && err.message ? err.message : String(err)));
    console.error('ERR_URL=' + page.url());
    const html = await page.content();
    console.error('DOM_HEAD=' + html.slice(0, 500).replace(/\s+/g, ' '));
    await page.screenshot({ path: '/private/tmp/studio-v2-real-check-error.png', fullPage: true });
    process.exitCode = 1;
  } finally {
    await browser.close();
  }
})();
