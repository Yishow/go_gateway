import { expect, test } from '@playwright/test';
import { buildEmbedded, startEmbedded, type EmbeddedBuild } from './helpers/embedded-server';
import { browserLedger, performanceLedger } from './helpers/embedded-browser-ledger';

test.describe('Embedded controlled two-build cache protocol (not a performance baseline)', () => {
  test.describe.configure({ mode: 'serial', timeout: 900_000 });
  let first: EmbeddedBuild;
  let second: EmbeddedBuild;
  let server: Awaited<ReturnType<typeof startEmbedded>>;
  test.beforeAll(async () => {
    first = await buildEmbedded({ cacheProtocolVariant: 'version-one' });
    second = await buildEmbedded({ cacheProtocolVariant: 'version-two' });
    server = await startEmbedded(first);
  });
  test.afterAll(async () => {
    if (server) await server.stop();
    if (first) await first.dispose();
    if (second) await second.dispose();
  });

  test('HTML revalidates across same-origin restart; unchanged hash is cached and missing old hash is 404', async ({ browser }, testInfo) => {
    const oldAssets = first.assets.filter(path => !second.assets.includes(path) && path.endsWith('.js'));
    const newAssets = second.assets.filter(path => !first.assets.includes(path) && path.endsWith('.js'));
    const unchanged = first.assets.find(path => /react-vendor.*\.js$/.test(path) && second.identity.assetSHA256[path] === first.identity.assetSHA256[path]);
    expect(oldAssets.length).toBeGreaterThan(0);
    expect(newAssets.length).toBeGreaterThan(0);
    expect(unchanged).toBeDefined();
    const context = await browser.newContext({ serviceWorkers: 'block' });
    const captures: unknown[] = [];
    try {
      const page = await context.newPage();
      const ledger = await browserLedger(page);
      const url = `${server.baseURL}/studio/v2?view=recordings`;
      const html1 = await page.goto(url, { waitUntil: 'domcontentloaded' });
      expect(html1?.headers()['cache-control']).toBe('no-cache');
      await expect(page.getByTestId('recording-home')).toBeVisible();
      const probe = async (path: string, cache: 'default' | 'no-store' = 'default') => page.evaluate(async ({ path, cache }) => {
        const response = await fetch(`/${path}`, { cache });
        return { status: response.status, cacheControl: response.headers.get('cache-control'), type: response.headers.get('content-type'), etag: response.headers.get('etag'), text: await response.text() };
      }, { path, cache });
      const unchanged1 = await probe(unchanged!);
      expect(unchanged1.cacheControl).toContain('immutable');
      const old1 = await probe(oldAssets[0]);
      expect(old1.status).toBe(200);
      captures.push({ phase: 'first-build', performance: await performanceLedger(page), network: structuredClone(ledger.events) });
      const port = server.port;
      await server.stop();
      server = await startEmbedded(second, port);
      ledger.events.length = 0;
      const html2 = await page.reload({ waitUntil: 'domcontentloaded' });
      expect(html2?.headers()['cache-control']).toBe('no-cache');
      await expect(page.getByTestId('recording-home')).toBeVisible();
      expect(await page.evaluate(() => Reflect.get(globalThis, '__embeddedCacheProtocolVersion'))).toBe('version-two');
      const unchanged2 = await probe(unchanged!);
      expect(unchanged2.text).toBe(unchanged1.text);
      expect(unchanged2.etag).toBe(unchanged1.etag);
      const newAsset = await probe(newAssets[0]);
      expect(newAsset.status).toBe(200);
      expect(newAsset.text).not.toBe(old1.text);
      expect(newAsset.etag).not.toBe(old1.etag);
      // Browser may correctly retain immutable old bytes. Force network for explicit missing-old-hash check.
      const missing = await probe(oldAssets[0], 'no-store');
      expect(missing.status).toBe(404);
      expect(missing.type ?? '').not.toContain('text/html');
      const resources = await page.evaluate((unchanged) => performance.getEntriesByType('resource').filter(entry => new URL(entry.name).pathname === `/${unchanged}`).map(entry => ({ transferSize: (entry as PerformanceResourceTiming).transferSize, decodedBodySize: (entry as PerformanceResourceTiming).decodedBodySize })), unchanged!);
      expect(resources.length).toBeGreaterThan(0);
      expect(resources.every(entry => entry.transferSize === 0 && entry.decodedBodySize > 0)).toBe(true);
      captures.push({ phase: 'second-build-warm', performance: await performanceLedger(page), network: structuredClone(ledger.events) });
      await ledger.dispose();
      const fresh = await browser.newContext({ serviceWorkers: 'block' });
      try {
        const freshPage = await fresh.newPage();
        await freshPage.goto(url, { waitUntil: 'domcontentloaded' });
        await expect(freshPage.getByTestId('recording-home')).toBeVisible();
        expect(await freshPage.evaluate(() => Reflect.get(globalThis, '__embeddedCacheProtocolVersion'))).toBe('version-two');
        captures.push({ phase: 'second-build-cold', performance: await performanceLedger(freshPage) });
      } finally { await fresh.close(); }
    } finally {
      await testInfo.attach('embedded-cache-upgrade.json', { contentType: 'application/json', body: Buffer.from(JSON.stringify({ first: first.identity, second: second.identity, oldAssets, newAssets, unchanged, captures }, null, 2)) });
      await context.close();
    }
  });
});
