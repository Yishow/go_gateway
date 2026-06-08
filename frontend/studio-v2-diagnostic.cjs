const { chromium } = require('@playwright/test');

(async () => {
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();

  page.on('console', (msg) => {
    console.log('C', msg.type(), msg.text());
  });
  page.on('pageerror', (err) => {
    console.log('PE', err.message);
  });

  await page.goto('http://127.0.0.1:4174/studio/v2', { waitUntil: 'domcontentloaded', timeout: 20000 });
  await page.waitForTimeout(3000);

  const hasStep1 = await page.locator('[data-testid="step1-device-view"]').count();
  const hasInputHost = await page.locator('[data-testid="input-host"]').count();
  const hasRunTest = await page.locator('[data-testid="run-test-button"]').count();
  const hasLoginText = await page.getByText('登入').count();
  const has404 = await page.locator('text=404').count();

  const testids = await page.$$eval('[data-testid]', (els) =>
    els.slice(0, 120).map((el) => el.getAttribute('data-testid')).filter(Boolean),
  );
  const bodyText = await page.locator('body').innerText().catch(() => '');

  console.log('TITLE=', await page.title());
  console.log('URL=', page.url());
  console.log('HAS_step1=', hasStep1);
  console.log('HAS_input_host=', hasInputHost);
  console.log('HAS_run_test=', hasRunTest);
  console.log('HAS_signin=', hasLoginText);
  console.log('HAS_404=', has404);
  console.log('TESTID_HEAD=', testids.slice(0, 40).join(','));
  console.log('BODY=', bodyText.replace(/\s+/g, ' ').slice(0, 1800));

  await page.screenshot({ path: '/private/tmp/studio-v2-diagnostic.png', fullPage: true });
  await browser.close();
})();
