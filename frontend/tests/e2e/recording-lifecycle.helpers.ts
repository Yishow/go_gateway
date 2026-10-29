import { expect, type Locator, type Page, type TestInfo } from '@playwright/test';
import { createRecordingLifecycleMock } from './fixtures/recording-lifecycle-api';

export const VIEWPORTS = [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }] as const;
export const FUTURE_BOUNDARY = '2030-01-01T01:00:00.000Z';
export type LifecycleMock = ReturnType<typeof createRecordingLifecycleMock> & {
  initialMutations: ReturnType<typeof createRecordingLifecycleMock>['ledger'];
};
export type LifecycleScenario = Parameters<typeof createRecordingLifecycleMock>[0];
export const dialogFor = (page: Page) => page.getByRole('dialog', { name: '變更已保存的記錄' });
export const basicFor = (page: Page) => page.getByTestId('step4-basic-recording');
export const advancedFor = (page: Page) => page.getByTestId('step4-advanced-recording');

/** Holds only a mock response after the mock ledger has accepted it; no timing race. */
export function responseGate(action: string) {
  let release = () => {};
  const wait = new Promise<void>((resolve) => { release = resolve; });
  return { action, wait, release };
}

/** Real Step4Database DOM with a deliberately fake, fail-closed API; not DB/backend acceptance. */
export async function openFixture(page: Page, scenario: LifecycleScenario, gate?: ReturnType<typeof responseGate>) {
  const mock = createRecordingLifecycleMock(scenario);
  await page.route('**/api/v1/**', async (route) => {
    const request = route.request();
    const body: unknown = request.postData() ? request.postDataJSON() : undefined;
    const response = mock.handle({ method: request.method(), url: request.url(), body });
    if (gate && body && typeof body === 'object' && 'action' in body && body.action === gate.action) await gate.wait;
    if (response.delayMs) await new Promise((resolve) => setTimeout(resolve, response.delayMs));
    if (response.transport === 'abort') { await route.abort('failed'); return; }
    await route.fulfill({ status: response.status, contentType: 'application/json', body: JSON.stringify(response.body) });
  });
  await page.clock.setFixedTime(new Date('2030-01-01T00:00:00Z'));
  await page.goto(`/tests/e2e/fixtures/recording-lifecycle-host.html?scenario=${scenario}&lang=zh-TW`);
  await expect(page.getByRole('heading', { name: 'Recording lifecycle — mock API / test data' })).toBeVisible();
  await expect(basicFor(page)).toBeVisible();
  await expect(basicFor(page).getByTestId('saved-recording-actions').first()).toBeAttached();
  const prepared = basicFor(page).getByTestId('circuit-prepared-group-instantaneous');
  if (await prepared.count()) {
    // Guided mount is navigation; only an explicit preview command may create a ledger/token.
    await expect(prepared.getByTestId('group-schema-preview-result')).toHaveCount(0);
    await expect(prepared.getByTestId('group-schema-error')).toHaveCount(0);
  }
  expect(mock.ledger.filter(entry => entry.mutation)).toEqual([]);
  return Object.assign(mock, { initialMutations: structuredClone(mock.ledger.filter(entry => entry.mutation)) });
}

export function assertNoNewEffects(mock: LifecycleMock) {
  expect(mock.ledger.filter(entry => entry.mutation)).toEqual(mock.initialMutations);
  expect(mock.initialMutations).toEqual([]);
}

export async function showBasicActions(page: Page, removed = false) {
  const basic = basicFor(page);
  if (removed) {
    const details = basic.getByTestId('basic-removed-recordings');
    if (await details.getAttribute('open') === null) await details.locator('summary').click();
  }
  const actions = basic.getByTestId('saved-recording-actions').first();
  await expect(actions).toBeVisible();
  return actions;
}

export async function showAdvanced(page: Page, removed = false) {
  const advanced = advancedFor(page);
  if (await advanced.getAttribute('open') === null) await advanced.locator(':scope > summary').click();
  if (removed) {
    const details = advanced.getByTestId('group-removed');
    if (await details.getAttribute('open') === null) await details.locator('summary').click();
  }
  await advanced.getByTestId('group-open-G1').click();
  await expect(advanced.getByTestId('group-editor')).toBeVisible();
  return advanced;
}

export async function openChange(page: Page, removed = false, mode: 'basic' | 'advanced' = 'basic') {
  const surface = mode === 'basic' ? await showBasicActions(page, removed) : await showAdvanced(page, removed);
  const trigger = surface.getByTestId(removed ? 'saved-recording-recreate' : 'saved-recording-transition');
  await expect(trigger).toBeEnabled();
  await trigger.focus();
  await page.keyboard.press('Enter');
  const dialog = dialogFor(page);
  await expect(dialog).toBeVisible();
  await expect(dialog.getByLabel('記錄名稱', { exact: true })).toBeEnabled();
  return { dialog, trigger };
}

export async function toSchemaPlan(page: Page) {
  const dialog = dialogFor(page);
  await dialog.getByRole('button', { name: '檢視變更影響', exact: true }).click();
  await dialog.getByRole('button', { name: '保留設定並預覽資料表', exact: true }).click();
  await expect(dialog.getByRole('region', { name: '伺服器回傳的資料表準備計畫' })).toBeVisible();
}

/** Convert an absolute instant through the browser's actual timezone, never an ISO-string slice. */
export async function fillBoundary(page: Page, instant = FUTURE_BOUNDARY) {
  const localValue = await page.evaluate((iso) => {
    const date = new Date(iso);
    const pad = (value: number) => String(value).padStart(2, '0');
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
  }, instant);
  const input = dialogFor(page).getByLabel('未來開始時間（您的當地時間）', { exact: true });
  await input.fill(localValue);
  expect(await input.evaluate(element => new Date((element as HTMLInputElement).value).toISOString())).toBe(new Date(instant).toISOString());
}

export async function toActivationReview(page: Page, doubleSchemaConfirm = false) {
  const dialog = dialogFor(page);
  const confirm = dialog.getByRole('button', { name: '確認此資料表計畫', exact: true });
  if (doubleSchemaConfirm) await confirm.dblclick({ delay: 30 });
  else await confirm.click();
  await expect(dialog.getByText('資料表準備已確認，尚未開始新記錄。', { exact: true })).toBeVisible();
  await expect(dialog.getByTestId('recording-change-review_activation')).toBeDisabled();
  await fillBoundary(page);
  await dialog.getByRole('button', { name: '審核開始時間與就緒狀態', exact: true }).click();
  await expect(dialog.locator('time')).toHaveText(FUTURE_BOUNDARY);
}

export async function captureGeometry(page: Page, info: TestInfo, name: string, dialog?: Locator) {
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  if (dialog) {
    const geometry = await dialog.evaluate((element) => {
      const box = element.getBoundingClientRect();
      return { x: box.x, y: box.y, right: box.right, bottom: box.bottom, width: element.clientWidth,
        scrollWidth: element.scrollWidth, viewportWidth: innerWidth, viewportHeight: innerHeight };
    });
    expect(geometry.x).toBeGreaterThanOrEqual(0);
    expect(geometry.y).toBeGreaterThanOrEqual(0);
    expect(geometry.right).toBeLessThanOrEqual(geometry.viewportWidth);
    expect(geometry.bottom).toBeLessThanOrEqual(geometry.viewportHeight);
    expect(geometry.scrollWidth).toBeLessThanOrEqual(geometry.width);
    await expect(dialog.getByRole('heading', { name: '變更已保存的記錄' })).toBeInViewport();
    await dialog.locator('footer').scrollIntoViewIfNeeded();
    await expect(dialog.locator('footer')).toBeInViewport({ ratio: 1 });
    const body = dialog.locator(':scope > div[aria-busy]');
    const last = body.locator(':scope > *').last();
    await last.scrollIntoViewIfNeeded();
    await expect(last).toBeInViewport();
    if (page.viewportSize()?.height === 720 && await body.getByTestId('lifecycle-pending-samples').count() > 0) {
      const scroll = await body.evaluate(element => ({ height: element.clientHeight, content: element.scrollHeight, top: element.scrollTop }));
      expect(scroll.content).toBeGreaterThan(scroll.height);
      expect(scroll.top).toBeGreaterThan(0);
    }
    const alert = dialog.getByRole('alert').first();
    if (await alert.count()) { await alert.scrollIntoViewIfNeeded(); await expect(alert).toBeInViewport(); }
  }
  const path = info.outputPath(`${name}-${page.viewportSize()!.width}-mock-api.png`);
  await page.screenshot({ path, fullPage: true });
  await info.attach(`${name} (actual browser; mock API/test data)`, { path, contentType: 'image/png' });
}

export async function assertPendingUnits(scope: Locator) {
  const history = scope.getByTestId('recording-lifecycle-history');
  await expect(history).toContainText('以下單位不同，請勿加總');
  await expect(history.getByTestId('lifecycle-pending-samples').first()).toHaveText('3');
  await expect(history.getByTestId('lifecycle-pending-outbox').first()).toHaveText('4');
  await expect(history.getByTestId('lifecycle-pending-blocked').first()).toHaveText('2');
  await expect(history.getByTestId('lifecycle-pending-unknown').first()).toHaveText('1');
  await expect(history).toContainText('原始連線');
  await expect(history).not.toContainText('7 筆');
}

export async function attachLedger(info: TestInfo, mock: LifecycleMock) {
  expect(mock.ledger.filter(entry => entry.status === 404), 'Unexpected API routes must not hide behind fail-closed mock responses').toEqual([]);
  await info.attach('mock-api-ledger-not-production-http-or-db', {
    body: Buffer.from(JSON.stringify({ evidence_layer: 'mock-api-browser', host_initial_mutations: mock.initialMutations, ledger: mock.ledger }, null, 2)),
    contentType: 'application/json',
  });
}
