import { test as base, expect, type Page } from '@playwright/test';
import { createRecordingHomeSnapshotAPI, type RecordingHomeSnapshotAPI, type SnapshotCall } from '../fixtures/recording-home-snapshot-api';
import { installSnapshotTransport } from '../fixtures/recording-home-snapshot-transport';
import { assertSnapshotHomeDemand } from './recording-home-snapshot-ledger';
export const HOME_URL = '/studio/v2?view=recordings#snapshot-proof';
export const sourceCount = (api: RecordingHomeSnapshotAPI) => api.ledger.filter(call => call.action === 'sources').length;
export const observation = (page: Page) => page.getByTestId('home-source-observation');
export const home = (page: Page) => page.getByTestId('recording-home');
export const callSlice = (api: RecordingHomeSnapshotAPI, checkpoint: number) => api.ledger.slice(checkpoint);
export async function rendered(page: Page) {
  await page.evaluate(() => new Promise<void>(resolve => requestAnimationFrame(() => requestAnimationFrame(() => resolve()))));
}
export async function settledRequests(api: RecordingHomeSnapshotAPI, ignored: readonly SnapshotCall[] = []) {
  await expect.poll(() => api.ledger.filter(call => !ignored.includes(call) && ['pending', 'fulfilled'].includes(call.delivery)).length).toBe(0);
}
export async function readyHome(page: Page, api: RecordingHomeSnapshotAPI) {
  await expect(observation(page)).toHaveAttribute('data-source-status', 'fresh');
  await expect(home(page)).toHaveAttribute('data-home-settling', 'false');
  await settledRequests(api); await rendered(page);
}
export async function visibleGroupIds(page: Page) {
  return page.locator('[data-testid="home-management-table"] [data-home-focus^="record-"]').evaluateAll(elements =>
    elements.map(element => element.getAttribute('data-home-focus')!.slice('record-'.length)));
}
export async function assertPageDemand(page: Page, api: RecordingHomeSnapshotAPI, checkpoint: number, requireFullPage = false) {
  await rendered(page);
  if (requireFullPage) for (const action of ['lifecycle', 'delivery']) {
    await expect.poll(() => callSlice(api, checkpoint).filter(call => call.phase === 'home' && call.action === action).length).toBe(25);
  }
  await settledRequests(api); await rendered(page);
  assertSnapshotHomeDemand(callSlice(api, checkpoint), await visibleGroupIds(page));
}
export async function openEditor(page: Page, api: RecordingHomeSnapshotAPI, groupId: string) {
  api.setPhase('editor');
  await page.locator(`[data-home-focus="record-${groupId}"]`).click();
  await expect(page).toHaveURL(/view=recording(?:&|#)/);
  await expect(page.getByTestId('simple-recording-shell')).toBeVisible();
  await expect(page.getByTestId('group-editor')).toBeVisible();
}
export async function returnHome(page: Page, api: RecordingHomeSnapshotAPI) {
  api.setPhase('home');
  await page.getByTestId('simple-recording-shell').getByRole('button', { name: 'Back to overview', exact: true }).click();
  await expect(page).toHaveURL(/view=recordings(?:&|#)/);
}
export async function filterPageTwo(page: Page, api: RecordingHomeSnapshotAPI, count: number) {
  const filters = page.getByRole('search', { name: 'Search and filters' });
  const changes = [() => filters.getByRole('searchbox').fill('Line A'),
    () => filters.getByRole('combobox', { name: 'Device', exact: true }).selectOption('D-A'),
    () => filters.getByRole('combobox', { name: 'Saved status', exact: true }).selectOption('ready'),
    () => filters.getByRole('combobox', { name: 'Destination', exact: true }).selectOption('conn-1')];
  for (const change of changes) { const checkpoint = api.ledger.length; await change(); await assertPageDemand(page, api, checkpoint); }
  const checkpoint = api.ledger.length;
  await home(page).getByRole('button', { name: 'Next', exact: true }).click();
  await expect(page.getByTestId('home-page-status')).toHaveText(`Page 2 of ${Math.ceil(count / 25)} · ${count} matches`);
  await assertPageDemand(page, api, checkpoint, true);
  expect(await visibleGroupIds(page)).toHaveLength(25);
}
export async function expectFiltersRestored(page: Page, count = 60) {
  const filters = page.getByRole('search', { name: 'Search and filters' });
  await expect(filters.getByRole('searchbox')).toHaveValue('Line A');
  await expect(filters.getByRole('combobox', { name: 'Device', exact: true })).toHaveValue('D-A');
  await expect(filters.getByRole('combobox', { name: 'Saved status', exact: true })).toHaveValue('ready');
  await expect(filters.getByRole('combobox', { name: 'Destination', exact: true })).toHaveValue('conn-1');
  await expect(page.getByTestId('home-page-status')).toHaveText(`Page 2 of ${Math.ceil(count / 25)} · ${count} matches`);
}
export async function openSourceRange(page: Page, api: RecordingHomeSnapshotAPI) {
  api.setPhase('source-modal'); await page.getByTestId('home-device-points-D-A').click();
  await expect(page.getByRole('dialog').getByTestId('rule-prefix-input')).toHaveValue(api.snapshot().rules[0].naming_prefix);
  await settledRequests(api); await rendered(page);
}

export const test = base.extend<{ snapshot: RecordingHomeSnapshotAPI; groupCount: number }>({
  groupCount: [60, { option: true }],
  snapshot: async ({ page, groupCount }, runCase, info) => {
    const api = createRecordingHomeSnapshotAPI({ groupCount });
    const unexpectedTransport = await installSnapshotTransport(page, api, String(info.project.use.baseURL));
    await page.setViewportSize({ width: 1280, height: 720 });
    await page.clock.setFixedTime(new Date('2030-01-01T00:00:00Z'));
    try { await runCase(api); }
    finally {
      api.releaseAll();
      await info.attach('routed-snapshot-ledger-MOCK-API-NOT-production-backend', { contentType: 'application/json',
        body: Buffer.from(JSON.stringify({ layer: 'production-/studio/v2-frontend-on-Vite-with-mocked-API', browserExecuted: true,
          route: page.url(), groupCount, lateResponseMeaning: 'finished = browser transport completed; fulfilled alone does not prove delivery; cancelled/not-delivered never counts as an ignored delivered reply', ledger: api.ledger }, null, 2)) });
      expect(unexpectedTransport, 'No external or non-fixture transport may escape').toEqual([]);
      expect(api.ledger.filter(call => call.unexpected), 'Every unlisted API/method/body must fail acceptance').toEqual([]);
      expect(api.ledger.filter(call => call.mutation && !(call.action === 'save-rule' && call.status === 200)), 'Only the separately approved source save may mutate').toEqual([]);
    }
  },
});
export { expect };
