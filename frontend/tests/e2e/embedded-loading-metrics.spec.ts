import { expect, test } from '@playwright/test';
import { buildEmbedded, startEmbedded, type EmbeddedBuild } from './helpers/embedded-server';
import { browserLedger, performanceLedger, summarize } from './helpers/embedded-browser-ledger';

test.describe('Embedded loading measurements: real fresh-empty backend, not live reading', () => {
  test.describe.configure({ mode: 'serial', timeout: 600_000 });
  let build: EmbeddedBuild;
  let server: Awaited<ReturnType<typeof startEmbedded>>;
  test.beforeAll(async () => { build = await buildEmbedded(); server = await startEmbedded(build); });
  test.afterAll(async () => { if (server) await server.stop(); if (build) await build.dispose(); });

  test('five cold contexts and five same-context warm reloads retain raw browser ledgers', async ({ browser }, testInfo) => {
    const samples: unknown[] = [];
    const values = { cold: { frame: [] as number[], settled: [] as number[], actionable: [] as number[], transferBytes: [] as number[], requestCount: [] as number[] },
      warm: { frame: [] as number[], settled: [] as number[], actionable: [] as number[], transferBytes: [] as number[], requestCount: [] as number[] } };
    try {
      for (let run = 1; run <= 5; run++) {
        const context = await browser.newContext({ serviceWorkers: 'block' });
        try {
          await context.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
          const page = await context.newPage();
          const ledger = await browserLedger(page);
          for (const mode of ['cold', 'warm'] as const) {
            ledger.events.length = 0;
            if (mode === 'cold') await page.goto(`${server.baseURL}/studio/v2?view=recordings`, { waitUntil: 'domcontentloaded' });
            else await page.reload({ waitUntil: 'domcontentloaded' });
            const home = page.getByTestId('recording-home');
            await expect(home).toBeVisible();
            const frameObservedMs = await page.evaluate(() => performance.now());
            await expect(home).toHaveAttribute('data-home-settling', 'false');
            await expect(page.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh');
            const settledObservedMs = await page.evaluate(() => performance.now());
            // Trial action checks hit target/actionability without opening editor or mutating settings.
            await home.getByRole('button').filter({ visible: true }).first().click({ trial: true });
            const actionableObservedMs = await page.evaluate(() => performance.now());
            const measurements = await performanceLedger(page);
            values[mode].frame.push(frameObservedMs);
            values[mode].settled.push(settledObservedMs);
            values[mode].actionable.push(actionableObservedMs);
            values[mode].transferBytes.push(measurements.totals.transferBytes);
            values[mode].requestCount.push(measurements.totals.requestCount);
            samples.push({ run, mode, frameObservedMs, settledObservedMs, actionableObservedMs,
              observation: 'automation-observed upper bounds from navigation timeOrigin, not paint or live-reading latency',
              performance: measurements, network: structuredClone(ledger.events) });
            await page.screenshot({ path: testInfo.outputPath(`${mode}-${run}.png`), fullPage: true });
          }
          await ledger.dispose();
        } finally { await context.close(); }
      }
    } finally {
      await testInfo.attach('embedded-loading-metrics.json', { contentType: 'application/json', body: Buffer.from(JSON.stringify({
        build: build.identity, browser: browser.version(), fixture: { kind: 'real fresh-empty SQLite', devices: 0, points: 0, externalEndpoints: 0 },
        methodology: '5 independent empty-cache browser contexts, each followed by one same-context warm reload; no request interception; same binary and server; no before/after A-only speed claim',
        summaries: Object.fromEntries(Object.entries(values).map(([mode, metrics]) => [mode, Object.fromEntries(Object.entries(metrics).map(([name, raw]) => [name, summarize(raw)]))])), samples,
      }, null, 2)) });
    }
    expect(samples).toHaveLength(10);
  });
});
