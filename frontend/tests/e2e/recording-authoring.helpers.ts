import { expect, type Locator, type Page, type TestInfo } from '@playwright/test';
import { createRecordingAuthoringMock } from './fixtures/recording-authoring-api';

export const AUTHORING_VIEWPORTS = [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }] as const;
export type AuthoringMock = ReturnType<typeof createRecordingAuthoringMock>;
export type AuthoringScenario = Parameters<typeof createRecordingAuthoringMock>[0];
export const basic = (page: Page) => page.getByTestId('step4-basic-recording');
export const authoring = (page: Page) => basic(page).getByTestId('circuit-profile-basic-panel');
export const advanced = (page: Page) => page.getByTestId('step4-advanced-recording');
export const mutationCalls = (mock: AuthoringMock) => mock.ledger.filter(entry => entry.mutation);
export const actionCalls = (mock: AuthoringMock, action: string) => mock.ledger.filter(entry => entry.action === action);

/** Test transport only. The real component, parser, query, persistence and controller remain installed. */
export async function openAuthoring(page: Page, scenario: AuthoringScenario = 'new-sqlite', options: {
  setup?: (mock: AuthoringMock) => void; language?: 'en' | 'zh-TW'; holdInitial?: string;
} = {}) {
  const mock = createRecordingAuthoringMock(scenario);
  const gates = new Map<string, { wait: Promise<void>; release: () => void }>();
  const holdNext = (action: string) => {
    let release = () => {};
    const wait = new Promise<void>(resolve => { release = resolve; });
    const gate = { wait, release };
    gates.set(action, gate);
    return gate;
  };
  const initialGate = options.holdInitial ? holdNext(options.holdInitial) : undefined;
  options.setup?.(mock);
  await page.route('**/api/**', async route => {
    const request = route.request();
    const body: unknown = request.postData() ? request.postDataJSON() : undefined;
    const response = mock.handle({ method: request.method(), url: request.url(), body });
    const action = mock.ledger.at(-1)?.action ?? '';
    const gate = gates.get(action);
    if (gate) { gates.delete(action); await gate.wait; }
    if (response.delayMs) await new Promise(resolve => setTimeout(resolve, response.delayMs));
    if (response.transport === 'abort') { await route.abort('failed'); return; }
    await route.fulfill({ status: response.status, contentType: 'application/json', body: JSON.stringify(response.body) });
  });
  await page.clock.setFixedTime(new Date('2030-01-01T00:00:00Z'));
  await page.goto(`/tests/e2e/fixtures/recording-authoring-host.html?scenario=${scenario}&lang=${options.language ?? 'en'}`);
  await expect(page.getByRole('heading', { name: 'Recording authoring — mock API / test data' })).toBeVisible();
  await expect(authoring(page)).toBeVisible();
  await expect.poll(() => actionCalls(mock, 'test-write-list').length).toBeGreaterThan(0);
  return { mock, holdNext, initialGate };
}

export async function expandAuthoringSources(page: Page) {
  const toggle = authoring(page).getByTestId('circuit-profile-sources-toggle');
  if (await toggle.getAttribute('aria-expanded') !== 'true') await toggle.click();
  await expect(authoring(page).getByTestId('circuit-profile-member-list')).toBeVisible();
}

export async function selectDevice(page: Page, id: 'dev-1' | 'dev-2') {
  const tile = basic(page).locator(`[data-testid="basic-recording-device-tile"][data-device-id="${id}"]`);
  await tile.click();
  await expect(tile).toHaveAttribute('aria-pressed', 'true');
}

export async function openAuthoringAdvanced(page: Page, groupId = 'G1') {
  const surface = advanced(page);
  if (await surface.getAttribute('open') === null) await surface.locator(':scope > summary').click();
  await surface.getByTestId(`group-open-${groupId}`).click();
  await expect(surface.getByTestId('group-editor')).toBeVisible();
  return surface;
}

export async function testWritePanel(page: Page) {
  const surface = await openAuthoringAdvanced(page);
  const panel = surface.getByTestId('group-test-write');
  await expect(panel).toBeVisible();
  return panel;
}

export function object(value: unknown): Record<string, unknown> {
  expect(value !== null && typeof value === 'object' && !Array.isArray(value)).toBe(true);
  return value as Record<string, unknown>;
}

export function ensureGroup(mock: AuthoringMock) {
  const calls = actionCalls(mock, 'ensure');
  expect(calls).toHaveLength(1);
  expect(calls[0].status, 'An attempted request rejected by the strict mock is not accepted authoring').toBe(200);
  expect(mock.snapshot().groups).toHaveLength(1);
  return object(object(calls[0].body).group);
}

/** Default zero-write assertions retain token preview attempts instead of hiding them. */
export function noAuthoringMutation(mock: AuthoringMock, allowed: string[] = []) {
  expect(mutationCalls(mock).filter(entry => !allowed.includes(entry.action))).toEqual([]);
}

export async function captureAuthoring(page: Page, info: TestInfo, name: string, focus?: Locator) {
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  if (focus) { await focus.scrollIntoViewIfNeeded(); await expect(focus).toBeVisible(); }
  const path = info.outputPath(`${name}-${page.viewportSize()!.width}-mock-api.png`);
  await page.screenshot({ path, fullPage: true });
  await info.attach(`${name} — actual browser, mock API/test data`, { path, contentType: 'image/png' });
}

export async function attachAuthoringEvidence(info: TestInfo, mock: AuthoringMock, expectedErrorActions: string[] = []) {
  expect(mock.ledger.filter(entry => entry.status === 404), 'Unexpected API paths are acceptance failures, including rejected mutations').toEqual([]);
  expect(mock.ledger.filter(entry => entry.status >= 400 && !expectedErrorActions.includes(entry.action)),
    'Only explicitly injected negative-case errors are allowed').toEqual([]);
  await info.attach('authoring-mock-api-ledger-not-SQL-or-production-wiring', {
    body: Buffer.from(JSON.stringify({ evidence_layer: 'real-component-browser-mock-api', ledger: mock.ledger, state: mock.snapshot() }, null, 2)),
    contentType: 'application/json',
  });
}
