import { spawn, type ChildProcessWithoutNullStreams } from 'node:child_process';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { createServer, type AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';

import { expect, test, type Page, type Request } from '@playwright/test';
import type { RuntimeLogSnapshot } from '../../src/types/runtimeLogs';

const repoRoot = process.cwd().endsWith('frontend') ? resolve(process.cwd(), '..') : process.cwd();
const embeddedPort = async (): Promise<number> => {
  const probe = createServer();
  await new Promise<void>((resolveProbe, reject) => {
    probe.once('error', reject);
    probe.listen(0, '127.0.0.1', () => resolveProbe());
  });

  const address = probe.address() as AddressInfo;
  const port = address.port;
  await new Promise<void>((resolveClose, reject) => {
    probe.close((error) => (error ? reject(error) : resolveClose()));
  });
  return port;
};

async function runProcess(
  command: string,
  args: string[],
  options: { cwd: string; env?: Record<string, string | undefined> },
): Promise<string> {
  const child = spawn(command, args, {
    cwd: options.cwd,
    env: { ...process.env, ...options.env },
    stdio: ['ignore', 'pipe', 'pipe'],
    windowsHide: true,
  });
  const output: string[] = [];
  child.stdout.on('data', (chunk: Buffer) => output.push(chunk.toString()));
  child.stderr.on('data', (chunk: Buffer) => output.push(chunk.toString()));

  const exitCode = await new Promise<number>((resolveProcess, reject) => {
    child.once('error', reject);
    child.once('close', (code) => resolveProcess(code ?? 1));
  });
  if (exitCode !== 0) {
    throw new Error(`${command} ${args.join(' ')} failed (${exitCode})\n${output.join('').slice(-4000)}`);
  }
  return output.join('');
}

async function waitForServer(baseURL: string, child: ChildProcessWithoutNullStreams): Promise<void> {
  const deadline = Date.now() + 120_000;
  while (Date.now() < deadline) {
    if (child.exitCode !== null) {
      throw new Error(`embedded server exited before readiness with code ${child.exitCode}`);
    }
    try {
      const response = await fetch(`${baseURL}/`);
      if (response.status === 200) return;
    } catch {
      // The embedded process may still be compiling migrations or binding its port.
    }
    await new Promise((resolveDelay) => setTimeout(resolveDelay, 250));
  }
  throw new Error(`embedded server did not become ready at ${baseURL}`);
}

function isAssetRequest(request: Request, origin: string): boolean {
  const url = new URL(request.url());
  if (url.origin !== origin) return false;
  return (
    request.resourceType() === 'script' ||
    request.resourceType() === 'stylesheet' ||
    /\.(?:js|mjs|css)(?:$|[?#])/i.test(url.pathname)
  );
}

function observeAssetRequests(page: Page, origin: string) {
  const assets = new Set<string>();
  const missing = new Set<string>();
  const failed = new Set<string>();

  page.on('request', (request) => {
    if (isAssetRequest(request, origin)) assets.add(request.url());
  });
  page.on('response', (response) => {
    const request = response.request();
    if (isAssetRequest(request, origin) && response.status() === 404) {
      missing.add(response.url());
    }
  });
  page.on('requestfailed', (request) => {
    if (isAssetRequest(request, origin)) failed.add(request.url());
  });

  return {
    assetCount: () => assets.size,
    urls: () => [...assets],
    missing: () => [...missing],
    failed: () => [...failed],
  };
}

function assetPaths(observer: ReturnType<typeof observeAssetRequests>): string[] {
  return [...new Set(observer.urls().map((url) => new URL(url).pathname))].sort();
}

function expectNoLegacyAssets(paths: string[]): void {
  for (const path of paths) {
    expect(path).not.toMatch(/DatalinkWorkbenchPage|workbench-frame|workbench-experiment/i);
  }
}

test.describe.configure({ mode: 'serial', timeout: 180_000 });

test.describe('Embedded frontend delivery smoke', () => {
  let server: ChildProcessWithoutNullStreams;
  let workDir: string;
  let baseURL: string;

  test.beforeAll(async () => {
    workDir = await mkdtemp(join(tmpdir(), 'go-gateway-embedded-e2e-'));
    const port = await embeddedPort();
    baseURL = `http://127.0.0.1:${port}`;
    const binary = join(workDir, process.platform === 'win32' ? 'test-ui.exe' : 'test-ui');

    const frontendRoot = join(repoRoot, 'frontend');
    const staticRoot = join(repoRoot, 'cmd', 'test_ui', 'static');
    await runProcess(process.platform === 'win32' ? 'npm.cmd' : 'npm', ['run', 'build'], { cwd: frontendRoot, env: { VITE_API_BASE_URL: '/api/v1' } });
    await rm(join(staticRoot, 'assets'), { recursive: true, force: true });
    await cp(join(frontendRoot, 'dist'), staticRoot, { recursive: true, force: true });
    await runProcess('go', ['build', '-o', binary, './cmd/test_ui'], { cwd: repoRoot });
    server = spawn(binary, [], {
      cwd: workDir,
      env: {
        ...process.env,
        PORT: String(port),
        HOST: '127.0.0.1',
        API_BASE_PATH: '/api/v1',
        GATEWAY_DB_PATH: join(workDir, 'datalink.db'),
        LOG_FILE: join(workDir, 'gateway-runtime.jsonl'),
        LOG_OUTPUT: 'console',
        AUTO_OPEN_BROWSER: 'false',
      },
      stdio: ['ignore', 'pipe', 'pipe'],
      windowsHide: true,
    });
    await waitForServer(baseURL, server);
  });

  test.afterAll(async () => {
    if (server && server.exitCode === null) {
      const closed = new Promise<void>((resolveExit) => server.once('close', () => resolveExit()));
      server.kill();
      await closed;
    }
    if (workDir) await rm(workDir, { recursive: true, force: true });
  });

  async function openRoute(page: Page, path: string, readySelector: string) {
    const observer = observeAssetRequests(page, new URL(baseURL).origin);
    await page.goto(`${baseURL}${path}`, { waitUntil: 'domcontentloaded' });
    await expect(page.locator(readySelector)).toBeVisible();
    await expect.poll(() => observer.assetCount()).toBeGreaterThan(0);
    expect(observer.missing()).toEqual([]);
    expect(observer.failed()).toEqual([]);
    return observer;
  }

  test('/studio/v2 keeps the guided workbench route loadable', async ({ page }) => {
    await openRoute(page, '/studio/v2', '[data-workbench-v2="true"]');
    await expect(page).toHaveURL(/\/studio\/v2/);
  });

  test('/studio follows the generic V2 fallback after legacy deletion', async ({ page }) => {
    await openRoute(page, '/studio', '[data-workbench-v2="true"]');
    await expect(page).toHaveURL(/\/studio\/v2$/);
  });

  test('/studio and an arbitrary unknown path share the generic fallback', async ({ page }) => {
    await openRoute(page, '/unknown-route-for-retirement', '[data-workbench-v2="true"]');
    await expect(page).toHaveURL(/\/studio\/v2$/);
  });

  test('/studio and unknown route preserve navigation and request the same V2 assets', async ({ browser }) => {
    const studioContext = await browser.newContext();
    const unknownContext = await browser.newContext();

    try {
      const studioPage = await studioContext.newPage();
      const unknownPage = await unknownContext.newPage();
      const studioObserver = await openRoute(
        studioPage,
        '/studio?step=output&target=database#legacy',
        '[data-workbench-v2="true"]',
      );
      const unknownObserver = await openRoute(
        unknownPage,
        '/unknown-route-for-retirement?step=output&target=database#legacy',
        '[data-workbench-v2="true"]',
      );

      for (const url of [studioPage.url(), unknownPage.url()]) {
        const finalURL = new URL(url);
        expect(finalURL.pathname).toBe('/studio/v2');
        expect(finalURL.search).toBe('?step=output&target=database');
        expect(finalURL.hash).toBe('#legacy');
      }

      const studioAssets = assetPaths(studioObserver);
      const unknownAssets = assetPaths(unknownObserver);
      expect(studioAssets).toEqual(unknownAssets);
      expectNoLegacyAssets(studioAssets);
      expectNoLegacyAssets(unknownAssets);
    } finally {
      await Promise.all([studioContext.close(), unknownContext.close()]);
    }
  });

  test('/studio/runtime keeps the focused monitor route loadable', async ({ page }) => {
    await openRoute(page, '/studio/runtime?device_id=device-A', '[data-testid="runtime-dashboard-route-state"]');
    await expect(page).toHaveURL(/\/studio\/runtime\?device_id=device-A/);
  });

  test('/studio/logs connects locally and browser controls preserve logging', async ({ page }, testInfo) => {
    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
    const observer = await openRoute(page, '/studio/logs', 'h1');
    await expect(page.getByRole('heading', { name: 'Runtime logs' })).toBeVisible();
    await expect(page.getByRole('status')).toHaveText('Connected');
    expect(assetPaths(observer).some((path) => path.includes('RuntimeLogsPage-'))).toBe(true);
    await page.getByLabel('Minimum severity').selectOption('warn');
    await page.getByLabel('Source (exact ID)').selectOption('runtime');
    await page.getByLabel('Search safe text').fill('no-match-embedded-smoke');
    await page.getByRole('button', { name: 'Apply filters' }).click();
    await expect(page.getByRole('status')).toHaveText('Connected');
    await expect(page.getByText('No matching retained records')).toBeVisible();
    await page.getByRole('button', { name: 'Pause view' }).focus();
    await page.keyboard.press('Enter');
    await expect(page.getByRole('status')).toHaveText('Paused');
    await page.getByRole('button', { name: 'Clear view' }).click();
    const snapshot = await page.request.get(`${baseURL}/api/v1/system/logs`);
    expect(snapshot.ok()).toBe(true);
    expect(Array.isArray((await snapshot.json()).records)).toBe(true);
    await page.getByRole('button', { name: 'Resume view' }).click();
    await expect(page.getByRole('status')).toHaveText('Connected');
    await page.getByLabel('Follow tail').uncheck();
    await expect(page.getByLabel('Follow tail')).not.toBeChecked();
    await page.screenshot({ path: testInfo.outputPath('runtime-logs-embedded.png'), fullPage: true });
    await page.getByRole('link', { name: 'Open runtime monitor' }).click();
    await expect(page.getByRole('link', { name: 'Runtime logs' })).toBeVisible();
    await page.goBack();
    await expect(page.getByRole('heading', { name: 'Runtime logs' })).toBeVisible();
    await expect(page.getByRole('status')).toHaveText('Connected');
  });

  // These cases use the owned embedded process and its real retained records;
  // only the browser network is faulted, without mocking log success responses.
  test('logs reconnect after an injected transport abort and browser offline/online without replay duplicates', async ({ page, context }) => {
    await page.addInitScript(() => {
      localStorage.setItem('i18nextLng', 'en');
      const originalFetch = window.fetch.bind(window);
      window.fetch = (input, init) => {
        const url = input instanceof Request ? input.url : String(input);
        if (!new URL(url, location.href).pathname.endsWith('/system/logs/stream')) return originalFetch(input, init);
        const fault = new AbortController();
        (window as Window & { dropLogStream: () => void }).dropLogStream = () => fault.abort();
        const signals = init?.signal ? [init.signal, fault.signal] : [fault.signal];
        return originalFetch(input, { ...init, signal: AbortSignal.any(signals) });
      };
    });
    const streams: string[] = [];
    page.on('request', (request) => {
      if (new URL(request.url()).pathname === '/api/v1/system/logs/stream') streams.push(request.url());
    });
    await openRoute(page, '/studio/logs', 'h1');
    const status = page.getByRole('status');
    const rows = page.locator('tbody tr');
    await expect(status).toHaveText('Connected');
    await expect.poll(() => rows.count()).toBeGreaterThan(0);
    const before = await rows.allTextContents();
    const initialStreams = streams.length;
    try {
      await context.setOffline(true);
      // Chromium offline mode blocks new requests but can keep an existing SSE
      // socket alive; cancel that real fetch to inject the dropped connection.
      await page.evaluate(() => (window as Window & { dropLogStream: () => void }).dropLogStream());
      await expect(status).toHaveText('Disconnected · reconnecting');
      expect(await rows.allTextContents()).toEqual(before);
    } finally {
      await context.setOffline(false);
    }
    await expect(status).toHaveText('Connected');
    await expect.poll(() => streams.length).toBeGreaterThan(initialStreams);
    expect(new URL(streams.at(-1)!).searchParams.get('after')).toMatch(/^[A-Za-z0-9_-]+:[0-9]+$/);
    expect(await rows.allTextContents()).toEqual(before);
  });

  test('logs clear affects only browser rows and preserves the server retained identities', async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
    await openRoute(page, '/studio/logs', 'h1');
    const status = page.getByRole('status');
    await expect(status).toHaveText('Connected');
    await expect.poll(() => page.locator('tbody tr').count()).toBeGreaterThan(0);
    const beforeResponse = await page.request.get(`${baseURL}/api/v1/system/logs`);
    expect(beforeResponse.ok()).toBe(true);
    const before: RuntimeLogSnapshot = await beforeResponse.json();
    const identities = (snapshot: RuntimeLogSnapshot) => snapshot.records.map((record) => `${record.instance_id}:${record.sequence}`);
    expect(before.records.length).toBeGreaterThan(0);
    await page.getByRole('button', { name: 'Pause view' }).click();
    await expect(status).toHaveText('Paused');
    await page.getByRole('button', { name: 'Clear view' }).click();
    await expect(page.locator('tbody tr')).toHaveCount(0);
    await expect(status).toHaveText('Paused');
    const afterResponse = await page.request.get(`${baseURL}/api/v1/system/logs`);
    expect(afterResponse.ok()).toBe(true);
    const after: RuntimeLogSnapshot = await afterResponse.json();
    expect(after.instance_id).toBe(before.instance_id);
    expect(identities(after)).toEqual(expect.arrayContaining(identities(before)));
    expect(after.capture_dropped).toBe(before.capture_dropped);
    await page.getByRole('button', { name: 'Resume view' }).click();
    await expect(status).toHaveText('Connected');
    await expect(page.locator('tbody tr')).toHaveCount(0);
  });

  test('logs source and text filters select real retained records and distinguish an empty match', async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
    const response = await page.request.get(`${baseURL}/api/v1/system/logs`);
    expect(response.ok()).toBe(true);
    const snapshot: RuntimeLogSnapshot = await response.json();
    const retained = snapshot.records.find((record) => record.source === 'startup');
    expect(retained).toBeDefined();
    await openRoute(page, '/studio/logs', 'h1');
    await expect(page.getByRole('status')).toHaveText('Connected');
    await page.getByLabel('Source (exact ID)').selectOption(retained!.source);
    await page.getByLabel('Search safe text').fill(retained!.code);
    const filteredRequest = page.waitForResponse((candidate) => {
      const url = new URL(candidate.url());
      return url.pathname === '/api/v1/system/logs' && url.searchParams.get('q') === retained!.code;
    });
    await page.getByRole('button', { name: 'Apply filters' }).click();
    const filtered = await filteredRequest;
    expect(new URL(filtered.url()).searchParams.get('source')).toBe(retained!.source);
    expect(filtered.ok()).toBe(true);
    const matching: RuntimeLogSnapshot = await filtered.json();
    expect(matching.records.length).toBeGreaterThan(0);
    expect(matching.records.every((record) => record.source === retained!.source && record.code === retained!.code)).toBe(true);
    await expect(page.getByRole('status')).toHaveText('Connected');
    await expect(page.locator('tbody tr')).toHaveCount(matching.records.length);
    await expect(page.getByText(retained!.message, { exact: true })).toBeVisible();
    await page.getByLabel('Search safe text').fill('no-match-logs-retained-browser-contract');
    await page.getByRole('button', { name: 'Apply filters' }).click();
    await expect(page.getByRole('status')).toHaveText('Connected');
    await expect(page.getByText('No matching retained records')).toBeVisible();
    await expect(page.locator('tbody tr')).toHaveCount(0);
  });

  test('zh-TW keyboard path reaches /studio/logs from Studio v2 and runs real controls', async ({ page }, testInfo) => {
    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'zh-TW'));
    await openRoute(page, '/studio/v2', '[data-workbench-v2="true"]');
    await expect(page.locator('[data-testid="workbench-v2-bootstrap-loading"]')).toBeHidden();
    const studioEntry = page.getByRole('link', { name: '執行日誌', exact: true });
    await expect(studioEntry).toHaveAttribute('href', '/studio/logs');
    await studioEntry.focus();
    await expect(studioEntry).toBeFocused();
    await page.keyboard.press('Enter');
    await expect(page).toHaveURL(/\/studio\/logs$/);
    await expect(page.getByRole('heading', { name: '執行日誌' })).toBeVisible();
    const status = page.getByRole('status');
    await expect(status).toHaveText('已連線');
    await page.getByLabel('搜尋安全文字').focus();
    await page.keyboard.type('no-match-embedded-zh-tw');
    await page.keyboard.press('Tab');
    await expect(page.getByRole('button', { name: '套用篩選' })).toBeFocused();
    await page.keyboard.press('Enter');
    await expect(page.getByText('沒有符合條件的保留紀錄')).toBeVisible();
    await expect(status).toHaveText('已連線');
    const toggle = page.getByRole('button', { name: '暫停檢視' });
    await toggle.focus();
    await page.keyboard.press('Enter');
    await expect(status).toHaveText('已暫停');
    await expect(page.getByRole('button', { name: '繼續檢視' })).toBeFocused();
    await page.keyboard.press('Space');
    await expect(status).toHaveText('已連線');
    await page.screenshot({ path: testInfo.outputPath('runtime-logs-embedded-zh-TW.png'), fullPage: true });
    const runtimeEntry = page.getByRole('link', { name: '開啟執行監看', exact: true });
    await expect(runtimeEntry).toHaveAttribute('href', '/studio/runtime');
    await runtimeEntry.focus();
    await page.keyboard.press('Enter');
    await expect(page).toHaveURL(/\/studio\/runtime(?:[?#]|$)/);
    const logsEntry = page.getByRole('link', { name: '執行日誌', exact: true });
    await expect(logsEntry).toHaveAttribute('href', '/studio/logs');
    await logsEntry.focus();
    await page.keyboard.press('Enter');
    await expect(page).toHaveURL(/\/studio\/logs$/);
    await expect(page.getByRole('status')).toHaveText('已連線');
  });

  test('/test keeps the independent engineering tool loadable', async ({ page }) => {
    await openRoute(page, '/test', 'h1');
    await expect(page.locator('h1').first()).toBeVisible();
    await expect(page).toHaveURL(/\/test$/);
  });

  test('/gateway/quick-setup keeps the experimental route loadable', async ({ page }) => {
    await page.route('**/api/v1/datalink/settings', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [{ key: 'ENABLE_GATEWAY_DUAL_ENTRY', value: true }] }),
      });
    });
    await openRoute(page, '/gateway/quick-setup', '[data-testid="protocol-select"]');
    await expect(page).toHaveURL(/\/gateway\/quick-setup$/);
  });

  test('reports a deliberately missing same-origin asset as a 404', async ({ page }) => {
    const observer = observeAssetRequests(page, new URL(baseURL).origin);
    await page.goto(`${baseURL}/test`, { waitUntil: 'domcontentloaded' });
    await expect(page.locator('h1').first()).toBeVisible();

    const missingURL = `${baseURL}/assets/__embedded-delivery-missing__.js`;
    const status = await page.evaluate(async (url) => (await fetch(url, { cache: 'no-store' })).status, missingURL);
    expect(status).toBe(404);
    await expect.poll(() => observer.missing()).toContain(missingURL);
  });
});
