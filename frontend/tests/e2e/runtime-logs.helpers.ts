import { expect, type Locator, type Page, type TestInfo } from '@playwright/test';
import type { RuntimeLogEvent, RuntimeLogMetadata, RuntimeLogRecord } from '../../src/types/runtimeLogs';
import { createRuntimeLogsMock } from './fixtures/runtime-logs-api';
import { RUNTIME_LOGS_HOST, RUNTIME_LOGS_PATH, RUNTIME_LOGS_TIME, runtimeLogsFrame, runtimeLogsSnapshot,
  type RuntimeLogsScenario } from './fixtures/runtime-logs-state';
import type { RuntimeLogsStreamController } from './fixtures/runtime-logs-stream';

export const LOG_VIEWPORTS = [{ width: 1024, height: 900 }, { width: 1600, height: 1000 }] as const;
export const LOG_LANGUAGES = ['en', 'zh-TW'] as const;
export type LogLanguage = typeof LOG_LANGUAGES[number];
type LogMock = ReturnType<typeof createRuntimeLogsMock>;
type FixtureWindow = Window & { __runtimeLogsStream: RuntimeLogsStreamController };
export const LOG_CASES = [
  ['RL-E01', 'Recording group fails and recovers', 'Recording group fails and recovers'],
  ['RL-E02', 'No proof of SQL commitment', 'No proof of SQL commitment'],
  ['RL-E03', 'Event capture fails', 'Event capture fails'],
  ['RL-E04', 'Hostile and secret-bearing fields', 'Hostile and secret-bearing fields'],
  ['RL-E05', 'Repeated failure and revision change', 'Repeated failure and revision change'],
  ['RL-E06', 'High-cardinality burst', 'High-cardinality burst'],
  ['RL-E07', 'Restart after a previous commit', 'Restart after a previous commit'],
  ['RL-T01', 'Remote viewing is denied without a proxy bypass', 'Remote viewing is denied without a proxy bypass'],
  ['RL-T02', 'Filters change with an old response in flight', 'Filters change with an old response in flight'],
  ['RL-T03', 'Filters change with an old response in flight', 'Filters change with an old response in flight'],
  ['RL-T04', 'Retention gap prevents complete fault history', 'Retention gap prevents complete fault history'],
  ['RL-N01', 'Operator opens the exact failed group', 'Operator opens the exact failed group'],
  ['RL-N02', 'Configuration was deleted or renamed', 'Configuration was deleted or renamed'],
  ['RL-N03', 'Configuration was deleted or renamed', 'Configuration was deleted or renamed'],
  ['RL-N04', 'Configuration was deleted or renamed', 'Configuration was deleted or renamed'],
  ['RL-N05', 'Operator opens the exact failed group', 'Operator opens the exact failed group'],
  ['RL-N06', 'Configuration was deleted or renamed', 'Configuration was deleted or renamed'],
] as const;
/** This catalog is preparation, never a browser result. Screenshots require real execution. */
export const FORMAL_LOG_CASE_CATALOG = Object.freeze(LOG_CASES.map(([id, scenario, example]) => ({ id, scenario, example,
  screenshotTemplate: `${id}-{variant}-{locale}-{width}-mock-data.png`,
  variants: id === 'RL-N04' ? ['foreign', 'unavailable'] : id === 'RL-N05' ? ['disclosure', 'native-focus'] : ['main'], requiresActualBrowser: true,
  exceptionScenario: 'Blocked browser execution is reported without a false pass',
  exceptionExample: 'Discovered case without a browser result' })));

export function logWords(language: LogLanguage) {
  return language === 'en' ? {
    title: 'Runtime logs', records: 'Retained runtime log records', scope: 'Scope filters (exact IDs)',
    apply: 'Apply filters', device: 'Device ID', rule: 'Source rule ID', level: 'Minimum severity', source: 'Source (exact ID)', stage: 'Evidence stage', search: 'Search safe text', workspace: 'Workspace ID', group: 'Recording group ID',
    actionable: 'Suggested action / unknown state', pause: 'Pause view', resume: 'Resume view', clear: 'Clear view',
    details: 'Event details', openGroup: 'Open group settings', revision: 'Original group revision', entityRevision: 'Original entity revision',
    reference: 'Reference ID', occurrences: 'Occurrences', first: 'First seen', last: 'Last seen', effect: 'Receipt effect key',
    receiptRevision: 'Receipt group revision', receipt: 'Runtime receipt (this effect only)',
    retrying: 'Delivery retrying', recovered: 'Delivery recovered', unknown: 'Delivery outcome unknown',
    unavailable: 'Event evidence unavailable', firstCommit: 'First SQL commitment observed',
    connected: 'Connected · log feed only', paused: 'Paused', reconnecting: 'Disconnected · reconnecting', denied: 'Access denied',
    historyUnknown: 'Current fault or recovery status is unknown where history is missing.',
    gap: 'Some capture history is missing.', reset: 'The gateway process restarted.', truncated: 'Some displayed content was truncated,',
    currentRevision: 'Current settings revision', historyHint: 'Requested revision R1 is informational',
    removed: 'This item was removed.', inaccessible: 'This workspace or selection is inaccessible.',
  } : {
    title: '執行日誌', records: '保留的執行日誌紀錄', scope: '範圍篩選（完整 ID）',
    apply: '套用篩選', device: '設備 ID', rule: '來源規則 ID', level: '最低嚴重程度', source: '來源（完整 ID）', stage: '證據階段', search: '搜尋安全文字', workspace: '工作區 ID', group: '寫入群組 ID',
    actionable: '有處理建議／狀態不明', pause: '暫停檢視', resume: '繼續檢視', clear: '清除畫面',
    details: '事件詳細資料', openGroup: '開啟群組設定', revision: '原始群組版本', entityRevision: '原始實體版本',
    reference: '參照 ID', occurrences: '發生次數', first: '首次觀察', last: '最近觀察', effect: '收據效果鍵',
    receiptRevision: '收據群組版本', receipt: '執行收據（僅此效果）',
    retrying: '資料交付重試中', recovered: '資料交付已恢復', unknown: '資料交付結果不明',
    unavailable: '事件證據無法取得', firstCommit: '首次觀察到 SQL 提交',
    connected: '已連線・僅日誌傳輸', paused: '已暫停', reconnecting: '已斷線・重新連線中', denied: '存取遭拒',
    historyUnknown: '歷史紀錄缺漏時，目前故障或恢復狀態不明。',
    gap: '部分接收歷史紀錄遺失', reset: '閘道程序已重新啟動', truncated: '部分內容已截斷',
    currentRevision: '目前設定版本', historyHint: '指定版本 R1 僅供參考', removed: '此項目已移除', inaccessible: '無法存取',
  };
}
export const logRows = (page: Page) => page.locator('tbody > tr');
export const fact = (scope: Locator, label: string) => scope.locator('dt').filter({ hasText: new RegExp(`^${label.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}$`) }).locator('..').locator('dd');
export const snapshotCalls = (mock: LogMock) => mock.ledger.filter(entry => entry.path === RUNTIME_LOGS_PATH);
export const locatorCalls = (mock: LogMock) => mock.ledger.filter(entry => entry.params.entry === 'locator');

export async function openLogs(page: Page, language: LogLanguage, scenario: RuntimeLogsScenario = 'normal', setup?: (mock: LogMock) => void) {
  const mock = createRuntimeLogsMock(scenario); setup?.(mock);
  const errors: string[] = []; page.on('pageerror', error => errors.push(error.message));
  const failedRequests: string[] = []; page.on('requestfailed', request => failedRequests.push(request.url()));
  const gates: Array<{ matches: (entry: LogMock['ledger'][number]) => boolean; wait: Promise<void>; arrived: () => void; finished: () => void }> = [];
  const holdNext = (matches: (entry: LogMock['ledger'][number]) => boolean) => {
    let release = () => {}; let arrived = () => {}; let finished = () => {};
    const wait = new Promise<void>(resolve => { release = resolve; });
    const reached = new Promise<void>(resolve => { arrived = resolve; });
    const settled = new Promise<void>(resolve => { finished = resolve; });
    gates.push({ matches, wait, arrived, finished }); return { release, reached, settled };
  };
  // Clone the resolved response BEFORE delaying it. Later scenario/filter changes
  // cannot turn an old A response into B data and accidentally hide cancellation bugs.
  await page.route('**/api/**', async route => {
    const request = route.request();
    const result = mock.handle({ method: request.method(), url: request.url(), body: request.postData() ?? undefined });
    const body = JSON.stringify(result.body), entry = mock.ledger.at(-1);
    const index = entry ? gates.findIndex(gate => gate.matches(entry)) : -1;
    const gate = index >= 0 ? gates.splice(index, 1)[0] : undefined;
    if (gate) { gate.arrived(); await gate.wait; }
    try { await route.fulfill({ status: result.status, contentType: 'application/json', body }); }
    finally { gate?.finished(); }
  });
  await page.clock.install({ time: new Date(RUNTIME_LOGS_TIME) });
  await page.goto(`${RUNTIME_LOGS_HOST}?lang=${language}`);
  await expect(page.getByTestId('runtime-logs-fixture-label')).toContainText('No backend, device or SQL verification.');
  await expect(page.getByRole('heading', { name: logWords(language).title, exact: true })).toBeVisible();
  await expect.poll(() => snapshotCalls(mock).length).toBe(1);
  return { mock, errors, failedRequests, holdNext, language };
}
export type LogJourney = Awaited<ReturnType<typeof openLogs>>;
export async function streams(page: Page) {
  return page.evaluate(() => (window as FixtureWindow).__runtimeLogsStream.requests.map(({ url, aborted, closed }) => ({ url, aborted, closed })));
}
export async function activeStream(page: Page, device?: string) {
  const find = async () => {
    const list = await streams(page);
    for (let i = list.length - 1; i >= 0; i--) {
      if (!list[i].closed && !list[i].aborted && (device === undefined || new URL(list[i].url, 'http://fixture.invalid').searchParams.get('device_id') === device)) return i;
    }
    return -1;
  };
  await expect.poll(find).toBeGreaterThanOrEqual(0); return find();
}
export async function sendFrame(page: Page, index: number, event: RuntimeLogEvent) {
  return page.evaluate(({ index, frame }) => (window as FixtureWindow).__runtimeLogsStream.send(index, frame), { index, frame: runtimeLogsFrame(event) });
}
export async function readyLogs(page: Page, language: LogLanguage, device?: string, metadata?: RuntimeLogMetadata) {
  const index = await activeStream(page, device);
  const request = (await streams(page))[index]; const after = new URL(request.url, 'http://fixture.invalid').searchParams.get('after') ?? '';
  const base = metadata ?? runtimeLogsSnapshot([], { latest_cursor: after });
  expect(await sendFrame(page, index, { type: 'handshake', metadata: base, cursor: after })).toBe(true);
  await expect(page.getByRole('status')).toHaveText(logWords(language).connected); return index;
}
export async function sendRecord(page: Page, index: number, record: RuntimeLogRecord) {
  return sendFrame(page, index, { type: 'log', record, cursor: `${record.instance_id}:${record.sequence}` });
}
export async function sendRecords(page: Page, index: number, records: RuntimeLogRecord[]) {
  const frames = records.map(record => runtimeLogsFrame({ type: 'log', record, cursor: `${record.instance_id}:${record.sequence}` }));
  return page.evaluate(({ index, frames }) => frames.filter(frame => (window as FixtureWindow).__runtimeLogsStream.send(index, frame)).length, { index, frames });
}
export async function failStream(page: Page, index: number) {
  await page.evaluate(index => (window as FixtureWindow).__runtimeLogsStream.fail(index), index);
}
export async function applyDevice(page: Page, language: LogLanguage, device: string, rule?: string) {
  const words = logWords(language);
  const prefix = new RegExp(`^${words.scope.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}`);
  const toggle = page.getByRole('button', { name: prefix });
  if (await toggle.getAttribute('aria-expanded') !== 'true') await toggle.click();
  await page.getByLabel(words.device, { exact: true }).fill(device);
  if (rule !== undefined) await page.getByLabel(words.rule, { exact: true }).fill(rule);
  await page.getByRole('button', { name: words.apply, exact: true }).click();
}
export async function disclose(row: Locator, language: LogLanguage) {
  const summary = row.locator('summary').filter({ hasText: logWords(language).details });
  const details = summary.locator('..');
  if (await details.getAttribute('open') === null) await summary.click();
  await expect(details).toHaveAttribute('open', ''); return details;
}
export async function tabTo(page: Page, target: Locator) {
  for (let i = 0; i < 240; i++) {
    await page.keyboard.press('Tab');
    if (await target.evaluate(element => element === document.activeElement)) break;
  }
  await expect(target).toBeFocused();
}
export async function captureLogs(page: Page, info: TestInfo, id: typeof LOG_CASES[number][0], language: LogLanguage, variant = 'main') {
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth), 'No horizontal viewport overflow').toBe(true);
  const width = page.viewportSize()?.width; if (!width) throw new Error('Explicit QA viewport required');
  const filename = `${id}-${variant}-${language}-${width}-mock-data.png`, path = info.outputPath(filename);
  await page.screenshot({ path, fullPage: true });
  await info.attach(`${filename} — real browser, mocked data; NOT backend/SQL proof`, { path, contentType: 'image/png' });
}
export async function finishLogs(info: TestInfo, journey: LogJourney) {
  expect(journey.errors).toEqual([]); expect(journey.mock.unexpected).toEqual([]);
  expect(journey.mock.ledger.filter(entry => entry.mutation)).toEqual([]);
  expect(journey.mock.ledger.filter(entry => entry.params.editor !== undefined || /retry|skip|activate|ddl|bootstrap|test-write/.test(entry.path))).toEqual([]);
  expect(journey.mock.ledger.filter(entry => /\/runtime\/logs|\/proxy\//.test(entry.path))).toEqual([]);
  await info.attach('strict-mock-http-ledger-not-backend-proof', { body: Buffer.from(JSON.stringify(journey.mock.ledger, null, 2)), contentType: 'application/json' });
}
