import { expect, type Locator, type Page, type TestInfo } from '@playwright/test';
import { createRuntimeRecordingMock } from './fixtures/runtime-recording-api';
import { installRuntimeRecordingStreamFixture } from './fixtures/runtime-recording-stream';

// Formal browser journeys use real route/components and mock API/test data.
// They are not SQL, backend clamp, deployed browser or screenshot evidence until executed.
export const RUNTIME_VIEWPORTS = [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }] as const;
export type RuntimeMock = ReturnType<typeof createRuntimeRecordingMock>;
export type RuntimeScenario = Parameters<typeof createRuntimeRecordingMock>[0];
export const actionCalls = (mock: RuntimeMock, action: string) => mock.ledger.filter(entry => entry.action === action);
export const sourceRules = (page: Page) => page.getByRole('region', { name: 'Source rules', exact: true });
export const sourceMappings = (page: Page) => page.getByRole('region', { name: 'Source mappings', exact: true });
export const sourceDetail = (page: Page) => page.getByRole('region', { name: 'Source rule details', exact: true });
export const groupCard = (page: Page) => page.getByRole('article', { name: 'Recording group G1', exact: true });
export const fact = (region: Locator, label: string) => region.locator('dt').filter({ hasText: new RegExp(`^${label.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}$`) }).locator('..').locator('dd');

type StreamControls = {
  list(): Array<{ url: string; closed: boolean }>;
  emit(index: number, event: string, data: unknown): void;
  fail(index: number): void;
};
type FixtureWindow = Window & { __runtimeRecordingStreams: StreamControls };

export async function openRuntime(page: Page, scenario: RuntimeScenario = 'normal', options: {
  language?: 'en' | 'zh-TW'; settingsQuery?: string; setup?: (mock: RuntimeMock) => void;
} = {}) {
  const mock = createRuntimeRecordingMock(scenario);
  const pageErrors: string[] = [];
  page.on('pageerror', error => pageErrors.push(error.message));
  const gates = new Map<string, { wait: Promise<void>; held: () => void }>();
  const holdNext = (action: string) => {
    if (gates.has(action)) throw new Error(`Already holding ${action}`);
    let release = () => {}; let held = () => {};
    const wait = new Promise<void>(resolve => { release = resolve; });
    const reached = new Promise<void>(resolve => { held = resolve; });
    gates.set(action, { wait, held });
    return { release, reached };
  };
  options.setup?.(mock);
  await page.route('**/api/**', async route => {
    const request = route.request();
    const response = mock.handle({ method: request.method(), url: request.url(), body: request.postData() ?? undefined });
    const action = mock.ledger.at(-1)?.action ?? 'unknown';
    const gate = gates.get(action);
    if (gate) { gates.delete(action); gate.held(); await gate.wait; }
    await route.fulfill({ status: response.status, contentType: 'application/json', body: JSON.stringify(response.body) });
  });
  await page.addInitScript(installRuntimeRecordingStreamFixture);
  await page.clock.setFixedTime(new Date('2026-10-29T10:00:00Z'));
  const query = options.settingsQuery
    ? `route=settings&${options.settingsQuery}` : 'device_id=A';
  await page.goto(`/tests/e2e/fixtures/runtime-recording-host.html?lang=${options.language ?? 'en'}&${query}`);
  if (!options.settingsQuery) {
    await expect(page.getByTestId('runtime-dashboard-header')).toBeVisible();
    await expect.poll(() => actionCalls(mock, 'context').length).toBeGreaterThan(0);
  }
  return { mock, holdNext, pageErrors };
}
export type RuntimeJourney = Awaited<ReturnType<typeof openRuntime>>;

export async function activeStream(page: Page, device = 'A') {
  await expect.poll(() => page.evaluate(id => (window as FixtureWindow).__runtimeRecordingStreams.list()
    .some(stream => { const url = new URL(stream.url, location.origin); return !stream.closed && url.pathname.endsWith('/runtime/stream') && url.searchParams.get('device_id') === id && url.searchParams.get('point_ids')?.split(',').includes(`P${id}1`); }), device)).toBe(true);
  return page.evaluate(id => {
    const streams = (window as FixtureWindow).__runtimeRecordingStreams.list();
    for (let index = streams.length - 1; index >= 0; index--) {
      const stream = streams[index]; const url = new URL(stream.url, location.origin);
      if (!stream.closed && url.pathname.endsWith('/runtime/stream') && url.searchParams.get('device_id') === id && url.searchParams.get('point_ids')?.split(',').includes(`P${id}1`)) return index;
    }
    return -1;
  }, device);
}
export async function emitStream(page: Page, index: number, event: string, data: unknown) {
  await page.evaluate(value => (window as FixtureWindow).__runtimeRecordingStreams.emit(value.index, value.event, value.data), { index, event, data });
}
export async function readyStream(page: Page, device = 'A') {
  const index = await activeStream(page, device);
  await emitStream(page, index, 'stream_state', { device_id: device, timestamp: '2026-10-29T10:00:00Z',
    stream_state: { state: 'ready', empty: false, stale: false, degraded: false, unavailable: false } });
  return index;
}
export async function failStream(page: Page, index: number) {
  await page.evaluate(value => (window as FixtureWindow).__runtimeRecordingStreams.fail(value), index);
}
export async function selectDevice(page: Page, device: 'A' | 'B') {
  await page.getByTestId('runtime-dashboard-device-switcher').getByRole('button', { name: `Device ${device}`, exact: true }).click();
  await expect(page.getByTestId('runtime-dashboard-header').getByRole('heading', { name: `Device ${device}`, exact: true })).toBeVisible();
}
export async function inspectSource(page: Page, rule = 'A1', address = '40001') {
  await sourceRules(page).getByRole('button', { name: `Show mappings for ${rule}`, exact: true }).click();
  await sourceMappings(page).getByRole('button', { name: `Inspect source ${address}`, exact: true }).click();
  await expect(sourceDetail(page)).toContainText('Saved-only configuration');
}
export async function openDisclosure(region: Locator, name: string) {
  const summary = region.locator('summary').filter({ hasText: new RegExp(`^${name}$`) });
  const detail = summary.locator('..');
  if (await detail.getAttribute('open') === null) await summary.click();
  await expect(detail).toHaveAttribute('open', '');
  return detail;
}
export async function tabTo(page: Page, target: Locator) {
  for (let index = 0; index < 240; index++) {
    await page.keyboard.press('Tab');
    if (await target.evaluate(element => element === document.activeElement)) break;
  }
  await expect(target).toBeFocused();
}
export async function captureRuntime(page: Page, info: TestInfo, name: string, focus?: Locator) {
  if (focus) { await focus.scrollIntoViewIfNeeded(); await expect(focus).toBeVisible(); }
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), 'No horizontal viewport overflow').toBe(true);
  const viewport = page.viewportSize();
  if (!viewport) throw new Error('The formal screenshot requires the declared viewport.');
  const path = info.outputPath(`${name}-${viewport.width}-mock-api.png`);
  await page.screenshot({ path, fullPage: true });
  await info.attach(`${name} — actual browser only; mocked API/test data, not SQL`, { path, contentType: 'image/png' });
}
export async function finishRuntime(info: TestInfo, journey: RuntimeJourney, expectedErrors: Array<{ action: string; status: number }> = []) {
  expect(journey.pageErrors, 'Actual page script errors are never silently accepted').toEqual([]);
  expect(journey.mock.ledger.filter(entry => entry.action === 'unknown'), 'No unhandled API passthrough').toEqual([]);
  expect(journey.mock.ledger.filter(entry => entry.mutation), 'All journeys are read-only').toEqual([]);
  expect(journey.mock.ledger.filter(entry => /\/system\/logs|\/runtime\/logs/.test(entry.path)), 'No process-log proxy').toEqual([]);
  expect(journey.mock.ledger.filter(entry => entry.status >= 400 && !expectedErrors.some(allowed => allowed.action === entry.action && allowed.status === entry.status)), 'Only named injected failures are allowed').toEqual([]);
  await info.attach('runtime-mock-api-ledger-NOT-SQL-or-production-proof', {
    body: Buffer.from(JSON.stringify({ evidence_layer: 'formal-browser-real-components-mocked-api', ledger: journey.mock.ledger }, null, 2)), contentType: 'application/json',
  });
}
