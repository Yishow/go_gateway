import { spawn, type ChildProcessWithoutNullStreams } from 'node:child_process';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { createServer, type AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';

import { expect, test, type Page, type Request } from '@playwright/test';

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
    await runProcess(process.platform === 'win32' ? 'npm.cmd' : 'npm', ['run', 'build'], { cwd: frontendRoot });
    await rm(join(staticRoot, 'assets'), { recursive: true, force: true });
    await cp(join(frontendRoot, 'dist'), staticRoot, { recursive: true, force: true });
    await runProcess('go', ['build', '-o', binary, './cmd/test_ui'], { cwd: repoRoot });
    server = spawn(binary, [], {
      cwd: workDir,
      env: {
        ...process.env,
        PORT: String(port),
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
