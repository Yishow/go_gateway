import { expect, test, type Page } from '@playwright/test';
import { buildEmbedded, startEmbedded } from './helpers/embedded-server';

// Production assets and real disposable embedded API. Only named JS/compact-read
// requests are faulted. No fixture HTML, substitute App, simulated API success, or PLC access.
test.describe('Embedded route loading and navigation', () => {
  test.describe.configure({ mode: 'serial', timeout: 180_000 });
  let build: Awaited<ReturnType<typeof buildEmbedded>>;
  let server: Awaited<ReturnType<typeof startEmbedded>>;
  let advancedURL: string;

  test.beforeAll(async () => {
    test.setTimeout(900_000);
    build = await buildEmbedded();
    server = await startEmbedded(build);
    const chunks = build.assets.filter(asset => /(?:^|\/)AdvancedWorkbenchV2Page-[^/]+\.js$/.test(asset));
    expect(chunks, 'Locate exactly one shipped advanced entry, never a guessed hash').toHaveLength(1);
    advancedURL = new URL(chunks[0].replace(/^\//, ''), `${server.baseURL}/`).href;
  });
  test.afterAll(async () => {
    try { await server?.stop(); } finally { await build?.dispose(); }
  });
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'en'));
  });

  async function home(page: Page, suffix = '') {
    await page.goto(`${server.baseURL}/studio/v2${suffix}`, { waitUntil: 'domcontentloaded' });
    await expect(page.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false');
    await expect(page.getByRole('button', { name: 'Devices', exact: true })).toBeVisible();
  }
  async function devices(page: Page) {
    await page.getByRole('button', { name: 'Devices', exact: true }).click();
    await expect(page.getByRole('button', { name: 'Add device', exact: true })).toBeEnabled();
  }
  function advancedRequests(page: Page) {
    const requests: string[] = [];
    page.on('request', request => { if (request.url() === advancedURL) requests.push(request.url()); });
    return requests;
  }

  test('cold home stays usable without fetching the advanced entry', async ({ page }) => {
    const requests = advancedRequests(page);
    await home(page, '?view=recordings#home');
    await devices(page);
    const search = page.getByRole('searchbox', { name: 'Search device names or protocols' });
    await search.fill('unsaved search');
    await expect(search).toHaveValue('unsaved search');
    await page.getByRole('button', { name: 'Recordings', exact: true }).click();
    await expect(page.getByRole('heading', { name: 'Recordings', exact: true })).toBeVisible();
    expect(requests).toEqual([]);
    await expect(page).toHaveURL(`${server.baseURL}/studio/v2?view=recordings#home`);
  });

  test('fault injection / cache disabled: delayed compact GET keeps the reading frame, then reveals real empty evidence', async ({ page }) => {
    let release!: () => void;
    const gate = new Promise<void>(resolve => { release = resolve; });
    const compactPath = '/api/v1/datalink/studio-v2/recording-home/sources';
    let held = 0;
    // Playwright routing disables HTTP cache. This is a response-delay fault
    // assertion, not the successful-empty latency measurement or a mock API.
    await page.route(url => url.origin === server.baseURL && url.pathname === compactPath, async route => {
      expect(route.request().method()).toBe('GET');
      held += 1;
      await gate;
      await route.continue();
    });
    try {
      await page.goto(`${server.baseURL}/studio/v2`, { waitUntil: 'domcontentloaded' });
      // Capture all initial-frame properties together, before the production
      // 2.5-second settle bound. We do not claim it waits forever for evidence.
      await expect.poll(async () => page.evaluate(() => {
        const visible = (element: Element | null) => Boolean(element && element.getClientRects().length &&
          getComputedStyle(element).visibility === 'visible' && getComputedStyle(element).display !== 'none');
        const frame = document.querySelector('[data-testid="recording-home"]');
        const emptySuccess = [...(frame?.querySelectorAll('p') ?? [])].some(element =>
          ['No items have been saved.', 'Connect your first device'].includes(element.textContent ?? '') && visible(element));
        return {
          frame: visible(frame), header: visible(document.querySelector('[data-testid="recording-product-header"]')),
          title: visible(frame?.querySelector('h1') ?? null), settling: frame?.getAttribute('data-home-settling'),
          reading: visible(frame?.querySelector('.gw-home-settling-status') ?? null),
          source: frame?.querySelector('[data-testid="home-source-observation"]')?.getAttribute('data-source-status'),
          emptySuccess,
        };
      })).toEqual({ frame: true, header: true, title: true, settling: 'true', reading: true, source: 'loading', emptySuccess: false });
      expect(held).toBe(1);
      // Past the bounded reading frame, pending evidence is allowed to reveal
      // an unknown state, but must never become successful-empty evidence.
      await expect(page.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false');
      await expect(page.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'loading');
      await expect(page.getByText('No items have been saved.', { exact: true })).toBeHidden();
      await expect(page.getByText('Connect your first device', { exact: true })).toBeHidden();
      const responsePromise = page.waitForResponse(response => new URL(response.url()).pathname === compactPath);
      release();
      const response = await responsePromise;
      expect(response.ok()).toBe(true);
      expect(await response.json()).toMatchObject({ success: true, data: { state: 'available', devices: [], candidates: [], excluded: [] } });
      await expect(page.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh');
      await expect(page.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false');
      await expect(page.getByText('No items have been saved.', { exact: true })).toBeVisible();
      await expect(page.getByText('Connect your first device', { exact: true })).toBeVisible();
    } finally { release(); }
  });

  test('delayed shipped chunk can close and reopen without replacing the home', async ({ page }) => {
    let release!: () => void;
    const gate = new Promise<void>(resolve => { release = resolve; });
    const requests = advancedRequests(page);
    await page.route(advancedURL, async route => { await gate; await route.continue(); });
    try {
      await home(page);
      await devices(page);
      const trigger = page.getByRole('button', { name: 'Add device', exact: true });
      await trigger.click();
      await expect(page.getByTestId('advanced-load-pending')).toBeVisible();
      await expect.poll(() => requests.length).toBe(1);
      await expect(page.getByTestId('recording-home')).toHaveCount(1);
      await page.getByTestId('advanced-load-pending').getByRole('button', { name: 'Close', exact: true }).click();
      await expect(page.getByRole('dialog')).toHaveCount(0);
      await expect(trigger).toBeFocused();
      await page.getByRole('searchbox').fill('still usable');
      await page.getByRole('searchbox').fill('');
      await trigger.click();
      await expect(page.getByTestId('advanced-load-pending')).toBeVisible();
      release();
      await expect(page.getByTestId('recording-source-name')).toBeEditable();
      await expect(page.getByTestId('advanced-load-pending')).toHaveCount(0);
      await page.getByTestId('recording-source-close').click();
      await expect(page.getByRole('dialog')).toHaveCount(0);
      await trigger.click();
      await expect(page.getByTestId('recording-source-name')).toBeEditable();
      expect(requests).toHaveLength(1);
    } finally { release(); }
  });

  test('rejected shipped chunk leaves home usable and explicit reload recovers', async ({ page }) => {
    const requests = advancedRequests(page);
    await page.route(advancedURL, route => route.abort('failed'));
    await home(page);
    await devices(page);
    const trigger = page.getByRole('button', { name: 'Add device', exact: true });
    await trigger.click();
    const error = page.getByTestId('advanced-load-error');
    await expect(error).toBeVisible();
    await expect(error).not.toContainText(/Failed to fetch|https?:\/\/|\.js/);
    await error.getByRole('button', { name: 'Close', exact: true }).click();
    await expect(trigger).toBeFocused();
    await page.getByRole('searchbox').fill('home after rejection');
    await page.getByRole('searchbox').fill('');
    await trigger.click();
    await expect(error).toBeVisible();
    expect(requests).toHaveLength(1); // React.lazy caches its rejection until reload.
    await page.unroute(advancedURL);
    await error.getByRole('button', { name: 'Refresh', exact: true }).click();
    await expect(page.getByRole('button', { name: 'Devices', exact: true })).toBeVisible();
    await devices(page);
    await trigger.click();
    await expect(page.getByTestId('recording-source-name')).toBeEditable();
    await expect(error).toHaveCount(0);
    expect(requests).toHaveLength(2);
  });

  test('advanced direct link retains query and hash through chunk delay and reload', async ({ page }) => {
    let release!: () => void;
    const gate = new Promise<void>(resolve => { release = resolve; });
    await page.route(advancedURL, async route => { await gate; await route.continue(); });
    const target = `${server.baseURL}/studio/v2?view=advanced&step=output&target=database#delivery`;
    try {
      await page.goto(target, { waitUntil: 'domcontentloaded' });
      await expect(page.getByTestId('advanced-load-pending')).toBeVisible();
      await expect(page).toHaveURL(target);
      release();
      await expect(page.getByTestId('workbench-v2-root')).toBeVisible();
      await expect(page.getByTestId('workbench-v2-bootstrap-error')).toHaveCount(0);
      await expect(page).toHaveURL(target);
      await page.reload();
      await expect(page.getByTestId('workbench-v2-root')).toBeVisible();
      await expect(page).toHaveURL(target);
    } finally { release(); }
  });

  test('repeated native history navigation preserves route intent and does not reload the chunk', async ({ page }) => {
    const requests = advancedRequests(page);
    const original = `${server.baseURL}/studio/v2?view=recordings#history`;
    const advanced = `${server.baseURL}/studio/v2?view=advanced#history`;
    await home(page, '?view=recordings#history');
    await page.locator('.gw-product-tools > details > summary').click();
    await page.getByRole('button', { name: 'Advanced setup', exact: true }).click();
    await expect(page.getByTestId('workbench-v2-root')).toBeVisible();
    await expect(page).toHaveURL(advanced);
    for (let visit = 0; visit < 2; visit++) {
      await page.goBack();
      await expect(page.getByTestId('recording-home')).toBeVisible();
      await expect(page).toHaveURL(original);
      await page.goForward();
      await expect(page.getByTestId('workbench-v2-root')).toBeVisible();
      await expect(page).toHaveURL(advanced);
    }
    expect(requests).toHaveLength(1);
  });

  test('dirty source draft survives cancelled native Forward and cancelled Close', async ({ page }) => {
    const mutations: string[] = [];
    // Safety tripwire: this case edits local fields only. Never reach a device,
    // test connection, or permit an accidental save against even the disposable DB.
    await page.route('**/api/**', async route => {
      if (['GET', 'HEAD', 'OPTIONS'].includes(route.request().method())) return route.continue();
      mutations.push(`${route.request().method()} ${new URL(route.request().url()).pathname}`);
      return route.abort('blockedbyclient');
    });
    await home(page, '?view=recordings#draft');
    await page.locator('.gw-product-tools > details > summary').click();
    await page.getByRole('button', { name: 'Advanced setup', exact: true }).click();
    await expect(page.getByTestId('workbench-v2-root')).toBeVisible();
    await page.goBack();
    await expect(page.getByTestId('recording-home')).toBeVisible();
    await devices(page);
    await page.getByRole('button', { name: 'Add device', exact: true }).click();
    const name = page.getByTestId('recording-source-name');
    await expect(name).toBeEditable();
    await name.fill('Unsent route draft');
    const original = page.url();
    // Native browser history, not a test-only navigate callback or pushState.
    await page.goForward();
    const confirmation = page.getByRole('alertdialog');
    await expect(confirmation).toBeVisible();
    await confirmation.getByRole('button', { name: 'Keep editing', exact: true }).click();
    await expect(page).toHaveURL(original);
    await expect(name).toHaveValue('Unsent route draft');
    await page.getByTestId('recording-source-close').click();
    await page.getByRole('dialog').getByRole('button', { name: 'Keep editing', exact: true }).click();
    await expect(name).toHaveValue('Unsent route draft');
    expect(mutations).toEqual([]);
  });
});
