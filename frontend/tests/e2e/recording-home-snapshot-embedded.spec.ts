import { createHash } from 'node:crypto';
import { expect, test, type Request } from '@playwright/test';
import { buildEmbedded, type EmbeddedBuild } from './helpers/embedded-server';
import { browserLedger, performanceLedger } from './helpers/embedded-browser-ledger';
import { compactSourcePath, measureCompactReads, seedRecordingHomeFixture } from './helpers/recording-home-seed';
import { assertNoConfigurationChanges, changedTables } from './helpers/recording-home-sqlite';

// Production frontend + production API; passive observers only, never page.route/context.route.
// Requires Node with node:sqlite DatabaseSync/readOnly (22.13+; cloud smoke uses 24).
// The empty mappings fixture exercises device parity; mapped identity guards remain separate unit evidence.
test.use({ trace: 'on' });
test.describe('Recording home real embedded 20-device / 2000-point snapshots', () => {
  test.describe.configure({ mode: 'serial', timeout: 600_000 });
  let build: EmbeddedBuild;
  let fixture: Awaited<ReturnType<typeof seedRecordingHomeFixture>>;
  test.beforeAll(async () => { build = await buildEmbedded(); fixture = await seedRecordingHomeFixture(build); });
  test.afterAll(async () => { try { await fixture?.stop(); } finally { await build?.dispose(); } });

  test('compact parity and real home returns retain separate HTTP, browser and table-hash evidence', async ({ page, browser }, testInfo) => {
    let endpoint: Awaited<ReturnType<typeof measureCompactReads>>;
    try { endpoint = await measureCompactReads(fixture); } catch (error) {
      await testInfo.attach('recording-home-compact-failure.json', { contentType: 'application/json', body: Buffer.from(JSON.stringify({
        message: String(error), evidence: error instanceof Error && 'evidence' in error ? error.evidence : null,
      }, null, 2)) });
      throw error;
    }
    const network = await browserLedger(page);
    let phase = 'cold-home';
    const requests: { phase: string; method: string; path: string; url: string }[] = [];
    const pending = new Set<Request>();
    const captures: Promise<void>[] = [];
    const compactReplies: unknown[] = [];
    const captureErrors: string[] = [];
    const segments: unknown[] = [];
    page.on('request', request => {
      const url = new URL(request.url());
      if (url.origin !== fixture.server.baseURL || !url.pathname.startsWith('/api/')) return;
      requests.push({ phase, method: request.method(), path: url.pathname, url: request.url() });
      if (!url.pathname.endsWith('/stream')) pending.add(request);
    });
    page.on('requestfinished', request => pending.delete(request));
    page.on('requestfailed', request => pending.delete(request));
    page.on('response', response => {
      if (new URL(response.url()).pathname !== compactSourcePath) return;
      const observedPhase = phase;
      captures.push((async () => {
        const bytes = await response.body();
        const body = JSON.parse(bytes.toString());
        expect(response.status()).toBe(200);
        expect(body).toMatchObject({ success: true, data: { state: 'available', workspace_id: fixture.setup.workspaceId,
          devices: fixture.setup.expectedDevices, candidates: [], excluded: [] } });
        expect(body.data.source_digest).toBe(endpoint.first.data.source_digest);
        compactReplies.push({ phase: observedPhase, status: response.status(), bytes: bytes.length,
          sha256: createHash('sha256').update(bytes).digest('hex'), sourceDigest: body.data.source_digest });
      })().catch(error => { captureErrors.push(String(error)); }));
    });
    const drain = async () => { await expect.poll(() => pending.size, { message: 'Wait for observed non-stream API requests to finish' }).toBe(0); await Promise.all(captures); };
    const settleHome = async () => {
      await expect(page.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false');
      await expect(page.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh');
      await drain();
    };
    const checkpoint = () => ({ database: fixture.snapshot(), requests: requests.length, compact: compactReplies.length });
    const finishSegment = async (name: string, before: ReturnType<typeof checkpoint>, compactCount: number, enforceReadonly = true) => {
      await drain(); const after = fixture.snapshot();
      const selected = requests.slice(before.requests);
      segments.push({ name, url: page.url(), before: before.database, after, changes: changedTables(before.database, after),
        requests: selected, compactReplies: compactReplies.slice(before.compact), performance: await performanceLedger(page),
        interpretation: enforceReadonly ? 'Configuration hashes compared; operational differences retained separately' : 'Editor/modal legacy source GET phase; configuration writes are reported, not counted as compact/home writes' });
      expect(selected.filter(item => item.path === compactSourcePath)).toHaveLength(compactCount);
      expect(selected.filter(item => !['GET', 'HEAD', 'OPTIONS'].includes(item.method))).toEqual([]);
      if (enforceReadonly) {
        assertNoConfigurationChanges(before.database, after);
        expect(selected.filter(item => /\/workspace\/(devices|source-rules|mappings)$|\/datalink\/points$/.test(item.path))).toEqual([]);
      }
      fixture.sentinel.assertQuiet();
    };
    try {
      await page.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
      let before = checkpoint();
      await page.goto(`${fixture.server.baseURL}/studio/v2?view=recordings#snapshot`, { waitUntil: 'domcontentloaded' });
      await settleHome();
      await finishSegment(phase, before, 1);
      await page.screenshot({ path: testInfo.outputPath('cold-home.png'), fullPage: true });

      // New unsaved recording follows the real editor route, but never saves or starts it.
      // All editor discovery GETs are outside the home-read interval, because legacy device reads mutate timestamps.
      phase = 'editor-discovery'; before = checkpoint();
      await page.getByRole('button', { name: 'New recording', exact: true }).click();
      await expect(page.getByTestId('group-editor')).toBeVisible();
      await expect(page.getByRole('button', { name: 'Back to overview', exact: true })).toBeEnabled();
      await drain(); await finishSegment(phase, before, 0, false);

      phase = 'warm-editor-return'; before = checkpoint();
      await page.getByRole('button', { name: 'Back to overview', exact: true }).click();
      await settleHome(); await finishSegment(phase, before, 1);
      await expect(page).toHaveURL(/\/studio\/v2\?view=recordings#snapshot$/);
      await page.screenshot({ path: testInfo.outputPath('warm-return.png'), fullPage: true });

      phase = 'history-editor'; before = checkpoint();
      await page.goBack(); await expect(page.getByTestId('group-editor')).toBeVisible();
      await drain(); await finishSegment(phase, before, 0, false);
      phase = 'warm-history-forward'; before = checkpoint();
      await page.goForward(); await settleHome(); await finishSegment(phase, before, 1);

      // Home remains mounted while a source modal opens. Close is not a route return.
      // Source discovery is isolated; only closing the settled modal is a zero-write interval.
      phase = 'source-modal-discovery'; before = checkpoint();
      const source = page.getByTestId('home-device-points-perf-device-1');
      await source.click();
      await expect(page.getByRole('dialog')).toHaveAttribute('data-initial-loading', 'false');
      await drain(); await finishSegment(phase, before, 0, false);
      phase = 'source-modal-close'; before = checkpoint();
      await page.getByTestId('home-range-done').click();
      await expect(page.getByRole('dialog')).toHaveCount(0); await expect(source).toBeFocused();
      await finishSegment(phase, before, 0);
      expect(captureErrors).toEqual([]);
      expect(compactReplies).toHaveLength(3);
    } finally {
      await Promise.all(captures);
      await testInfo.attach('recording-home-embedded-evidence.json', { contentType: 'application/json', body: Buffer.from(JSON.stringify({
        build: build.identity, browser: browser.version(), fixture: fixture.setup, seedLedger: fixture.ledger, endpoint,
        methodology: 'One cold navigation and two warm same-context route returns; no API interception. Raw timing/bytes only, no speed claim. Setup, bootstrap, editor discovery and modal reads are separate.',
        persistedWriteLimit: 'Hashes detect persisted table differences including configuration timestamps. Equal hashes do not prove zero transient, rolled-back, or write-then-restore operations. The client method ledger alone proves no server-write property.',
        segments, requests, compactReplies, captureErrors, network: network.events, protocolConnections: fixture.sentinel.connections,
      }, null, 2)) });
      await testInfo.attach('recording-home-embedded-server.log', { contentType: 'text/plain', body: Buffer.from(fixture.server.logs()) });
      await network.dispose();
    }
  });
});
