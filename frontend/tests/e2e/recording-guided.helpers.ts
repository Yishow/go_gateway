import { expect, type Locator, type Page, type TestInfo } from '@playwright/test';
import { createGuidedRecordingMock } from './fixtures/recording-guided-api';

export const GUIDED_VIEWPORTS = [{ width: 390, height: 844 }, { width: 768, height: 1024 }, { width: 1280, height: 900 }] as const;
export type GuidedMock = ReturnType<typeof createGuidedRecordingMock>;
export type GuidedScenario = Parameters<typeof createGuidedRecordingMock>[0];
export const basic = (page: Page) => page.getByTestId('step4-basic-recording');
export const authoring = (page: Page) => basic(page).getByTestId('circuit-profile-basic-panel');
export const advanced = (page: Page) => page.getByTestId('step4-advanced-recording');
export const mutationCalls = (mock: GuidedMock) => mock.ledger.filter(entry => entry.mutation);
export const actionCalls = (mock: GuidedMock, action: string) => mock.ledger.filter(entry => entry.action === action);
export async function openGuided(page: Page, scenario: GuidedScenario = 'new', options: {
  setup?: (mock: GuidedMock) => void; language?: 'en' | 'zh-TW'; holdInitial?: string;
} = {}) {
  const mock = createGuidedRecordingMock(scenario);
  const gates = new Map<string, { wait: Promise<void>; release: () => void }>();
  const holdNext = (action: string) => {
    let release = () => {};
    const wait = new Promise<void>(resolve => { release = resolve; });
    const gate = { wait, release }; gates.set(action, gate); return gate;
  };
  const initialGate = options.holdInitial ? holdNext(options.holdInitial) : undefined;
  options.setup?.(mock);
  await page.route('**/api/**', async route => {
    const request = route.request();
    const response = mock.handle({ method: request.method(), url: request.url(), body: request.postData() ? request.postDataJSON() : undefined });
    const action = mock.ledger.at(-1)?.action ?? '';
    const gate = gates.get(action); if (gate) { gates.delete(action); await gate.wait; }
    if (response.delayMs) await new Promise(resolve => setTimeout(resolve, response.delayMs));
    if (response.transport === 'abort') { await route.abort('failed'); return; }
    await route.fulfill({ status: response.status, contentType: 'application/json', body: JSON.stringify(response.body) });
  });
  await page.clock.setFixedTime(new Date('2030-01-01T00:00:00Z'));
  await page.goto(`/tests/e2e/fixtures/recording-guided-host.html?scenario=${scenario}&lang=${options.language ?? 'en'}`);
  await expect(page.getByRole('heading', { name: 'Guided recording — real components / mock API / test data' })).toBeVisible();
  await expect(page.getByTestId('workbench-v2-shell')).toBeVisible();
  await expect.poll(() => actionCalls(mock, 'list').length).toBeGreaterThan(0);
  return { mock, holdNext, initialGate };
}
export async function expandGuidedSources(page: Page) {
  const toggle = authoring(page).getByTestId('circuit-profile-sources-toggle');
  if (await toggle.getAttribute('aria-expanded') !== 'true') await toggle.click();
  await expect(authoring(page).getByTestId('circuit-profile-member-list')).toBeVisible();
}
export async function selectGuidedDevice(page: Page, id: 'dev-1' | 'dev-2' | 'dev-empty') {
  const tile = basic(page).locator(`[data-testid="basic-recording-device-tile"][data-device-id="${id}"]`);
  await tile.click(); await expect(tile).toHaveAttribute('aria-pressed', 'true');
}
export async function openGuidedAdvanced(page: Page, groupId?: string) {
  const surface = advanced(page);
  if (await surface.getAttribute('open') === null) await surface.locator(':scope > summary').click();
  if (groupId) { await surface.getByTestId(`group-open-${groupId}`).click(); await expect(surface.getByTestId('group-editor')).toBeVisible(); }
  return surface;
}
/** Existing explicit commands only; never injects a token/readiness or bypasses review. */
export async function prepareGuided(page: Page, mock: GuidedMock) {
  const panel = authoring(page), save = panel.getByTestId('circuit-profile-prepare');
  await expect(save).toBeEnabled(); await save.click();
  await expect.poll(() => actionCalls(mock, 'ensure').filter(call => call.status === 200).length).toBe(1);
  expect(actionCalls(mock, 'schema-preview')).toHaveLength(0);
  const preview = panel.getByTestId('group-schema-preview'); await expect(preview).toBeEnabled(); await preview.click();
  await expect(panel.getByTestId('group-schema-preview-result')).toBeVisible();
  expect(actionCalls(mock, 'schema-apply')).toHaveLength(0); expect(actionCalls(mock, 'start')).toHaveLength(0);
  const confirm = panel.getByTestId('group-schema-apply'); await expect(confirm).toBeEnabled(); await confirm.click();
  await expect.poll(() => actionCalls(mock, 'schema-apply').filter(call => call.status === 200).length).toBe(1);
  await expect(panel.getByTestId('circuit-profile-start')).toBeEnabled();
}
export async function startGuided(page: Page, mock: GuidedMock) {
  const before = actionCalls(mock, 'start').length;
  const review = authoring(page).getByTestId('circuit-profile-start'); await expect(review).toBeEnabled(); await review.click();
  await expect(page.getByRole('dialog')).toBeVisible(); expect(actionCalls(mock, 'start')).toHaveLength(before);
  await page.getByTestId('recording-start-confirm').click();
  await expect.poll(() => actionCalls(mock, 'start').length).toBe(before + 1);
}
export function noGuidedMutation(mock: GuidedMock, allowed: string[] = []) {
  expect(mutationCalls(mock).filter(entry => !allowed.includes(entry.action))).toEqual([]);
}
export async function tabTo(page: Page, target: Locator) {
  for (let index = 0; index < 120; index++) {
    await page.keyboard.press('Tab');
    if (await target.evaluate(element => element === document.activeElement)) break;
  }
  await expect(target).toBeFocused();
}
export async function captureGuided(page: Page, info: TestInfo, name: string, focus?: Locator) {
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  if (focus) { await focus.scrollIntoViewIfNeeded(); await expect(focus).toBeVisible(); }
  const path = info.outputPath(`${name}-${page.viewportSize()!.width}-mock-api.png`);
  await page.screenshot({ path, fullPage: true });
  await info.attach(`${name} — actual browser, mock API/test data`, { path, contentType: 'image/png' });
}
export async function attachGuidedEvidence(info: TestInfo, mock: GuidedMock, expectedErrorActions: string[] = []) {
  expect(mock.ledger.filter(entry => entry.status === 404), 'Unexpected paths, including rejected writes, are failures').toEqual([]);
  expect(mock.ledger.filter(entry => entry.status >= 400 && !expectedErrorActions.includes(entry.action)), 'Only named injected errors are allowed').toEqual([]);
  await info.attach('guided-mock-api-ledger-NOT-SQL-runtime-or-production-wiring', {
    body: Buffer.from(JSON.stringify({ evidence_layer: 'actual-components-browser-mock-api', ledger: mock.ledger, snapshot: mock.snapshot() }, null, 2)),
    contentType: 'application/json',
  });
}
