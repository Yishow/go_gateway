import { expect, test, type Page } from '@playwright/test';
import { domainRecord, runtimeLogsSnapshot } from './fixtures/runtime-logs-state';
import { LOG_LANGUAGES, LOG_VIEWPORTS, captureLogs, disclose, fact, finishLogs, locatorCalls,
  logRows, logWords, openLogs, readyLogs, snapshotCalls, tabTo, type LogLanguage } from './runtime-logs.helpers';

async function noEditor(page: Page) {
  for (const id of ['group-editor', 'basic-recording-panel', 'workspace-test-write-recovery']) await expect(page.getByTestId(id)).toHaveCount(0);
}
async function selectedGroup(page: Page, group: string) {
  const card = page.getByTestId(`group-list-item-${group}`);
  await expect(card).toHaveAttribute('data-locator-selected', 'true'); await expect(card).toBeFocused();
  await noEditor(page); return card;
}
async function originalRevision(page: Page, language: LogLanguage) {
  await disclose(logRows(page).first(), language);
  await expect(fact(logRows(page).first(), logWords(language).revision)).toHaveText('R1');
}

for (const viewport of LOG_VIEWPORTS) for (const language of LOG_LANGUAGES) {
  test.describe(`runtime logs native navigation ${language} ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });
    const words = logWords(language);

    test('RL-N01 keyboard opens exact W/G/R1 and focuses native current G/R2 without hydration or commands', async ({ page }, info) => {
      const journey = await openLogs(page, language); await readyLogs(page, language); await originalRevision(page, language);
      expect(locatorCalls(journey.mock)).toHaveLength(0);
      const link = logRows(page).first().getByRole('link', { name: words.openGroup, exact: true });
      const href = await link.getAttribute('href'); if (!href) throw new Error('Exact settings target required');
      const target = new URL(href, 'http://fixture.invalid');
      expect(target.pathname).toBe('/studio/v2');
      expect(Object.fromEntries(target.searchParams)).toEqual({ entry: 'locator', workspace_id: 'W', step: '4', group_id: 'G', group_revision: 'R1' });
      await tabTo(page, link); await page.keyboard.press('Enter');
      await expect(page).toHaveURL(url => url.pathname === '/studio/v2' && url.search === target.search);
      const card = await selectedGroup(page, 'G'); await expect(card).toContainText('Shared recorder');
      await expect(page.getByText(words.historyHint, { exact: false })).toBeVisible();
      await expect(page.getByText(words.currentRevision, { exact: false })).toBeVisible();
      await expect(page.getByText('R2', { exact: true })).toBeVisible();
      await expect.poll(() => locatorCalls(journey.mock).length).toBe(1);
      expect(locatorCalls(journey.mock)[0].params).toMatchObject({ workspace_id: 'W', group_id: 'G', group_revision: 'R1' });
      await captureLogs(page, info, 'RL-N01', language); await finishLogs(info, journey);
    });

    test('RL-N02 real Back Forward rechecks renamed current R3 while event remains R1', async ({ page }, info) => {
      const journey = await openLogs(page, language); await readyLogs(page, language);
      await logRows(page).first().getByRole('link', { name: words.openGroup, exact: true }).click();
      await selectedGroup(page, 'G'); await expect(page.getByText('R2', { exact: true })).toBeVisible();
      journey.mock.setScenario('renamed'); await page.goBack(); await readyLogs(page, language);
      await originalRevision(page, language); await expect(page.locator('tbody')).not.toContainText('Renamed current recorder');
      await page.goForward(); const card = await selectedGroup(page, 'G');
      await expect.poll(() => locatorCalls(journey.mock).length).toBe(2);
      await expect(card).toContainText('Renamed current recorder'); await expect(page.getByText('R3', { exact: true })).toBeVisible();
      expect(locatorCalls(journey.mock).every(call => call.params.group_id === 'G' && call.params.group_revision === 'R1')).toBe(true);
      await captureLogs(page, info, 'RL-N02', language); await finishLogs(info, journey);
    });

    test('RL-N03 deleted G1 never opens same-name G2 and retains historical identity', async ({ page }, info) => {
      const journey = await openLogs(page, language, 'deleted'); await readyLogs(page, language);
      await originalRevision(page, language); await expect(fact(logRows(page).first(), words.group)).toHaveText('G1');
      await logRows(page).first().getByRole('link', { name: words.openGroup, exact: true }).click();
      const notice = page.getByTestId('studio-locator-notice'); await expect(notice).toContainText(words.removed);
      await expect(page.getByTestId('group-list-item-G1')).toHaveCount(0); await expect(page.getByTestId('group-list-item-G2')).toHaveCount(0);
      expect(locatorCalls(journey.mock)).toHaveLength(1); expect(locatorCalls(journey.mock)[0].params.group_id).toBe('G1');
      expect(journey.mock.ledger.some(call => call.path.includes('/G2') || call.params.group_id === 'G2')).toBe(false);
      const links = await notice.getByRole('link').evaluateAll(elements => elements.map(element => element.getAttribute('href')));
      expect(links.length).toBeGreaterThan(0);
      for (const href of links) {
        if (!href) throw new Error('Owned workspace fallback link required'); const url = new URL(href, 'http://fixture.invalid');
        expect(url.searchParams.get('workspace_id')).toBe('W'); expect(url.searchParams.get('step')).toBe('1');
      }
      await noEditor(page); await captureLogs(page, info, 'RL-N03', language);
      await page.goBack(); await readyLogs(page, language); await originalRevision(page, language);
      await expect(fact(logRows(page).first(), words.group)).toHaveText('G1'); await finishLogs(info, journey);
    });

    for (const scenario of ['foreign', 'unavailable'] as const) {
      test(`RL-N04 ${scenario} resolver fails closed without another workspace or name replacement`, async ({ page }, info) => {
        const journey = await openLogs(page, language, scenario); await readyLogs(page, language);
        await logRows(page).first().getByRole('link', { name: words.openGroup, exact: true }).click();
        const notice = page.getByTestId('studio-locator-notice');
        await expect(notice).toContainText(scenario === 'foreign' ? words.inaccessible : language === 'en' ? 'unavailable' : '目前無法取得');
        await expect(page.getByTestId('group-list-item-G')).toHaveCount(0); await noEditor(page);
        expect(locatorCalls(journey.mock)).toHaveLength(1); expect(locatorCalls(journey.mock)[0].params.workspace_id).toBe('W');
        expect(journey.mock.ledger.filter(call => call.path !== '/api/v1/system/logs')).toHaveLength(1);
        await expect(page.locator('body')).not.toContainText(/password_fixture|private\.example|OTHER/);
        await captureLogs(page, info, 'RL-N04', language, scenario);
        await page.goBack(); await readyLogs(page, language); await originalRevision(page, language);
        await finishLogs(info, journey);
      });
    }

    test('RL-N06 late G locator response cannot replace B after actual Back and another event action', async ({ page }, info) => {
      const records = [domainRecord('domain.delivery.retrying'), domainRecord('domain.delivery.blocked', { sequence: '9007199254740994', fields: { group_id: 'B' } })];
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot(records)));
      await readyLogs(page, language);
      const late = journey.holdNext(call => call.params.entry === 'locator' && call.params.group_id === 'G');
      await logRows(page).first().getByRole('link', { name: words.openGroup, exact: true }).click(); await late.reached;
      await page.goBack(); await readyLogs(page, language);
      await logRows(page).nth(1).getByRole('link', { name: words.openGroup, exact: true }).click();
      await selectedGroup(page, 'B'); late.release(); await late.settled;
      await expect.poll(() => journey.failedRequests.some(url => new URL(url).searchParams.get('group_id') === 'G')).toBe(true);
      await expect(page.getByTestId('group-list-item-G')).toHaveCount(0);
      await expect(page.getByTestId('group-list-item-B')).toHaveAttribute('data-locator-selected', 'true');
      expect(locatorCalls(journey.mock).map(call => call.params.group_id)).toEqual(['G', 'B']);
      await expect(page).toHaveURL(url => url.searchParams.get('group_id') === 'B' && url.searchParams.get('group_revision') === 'R1');
      await captureLogs(page, info, 'RL-N06', language); await finishLogs(info, journey);
    });
  });
}

for (const language of LOG_LANGUAGES) {
  test.describe(`runtime logs narrow keyboard ${language}`, () => {
    test.use({ viewport: { width: 390, height: 844 }, timezoneId: 'UTC' });
    test('RL-N05 narrow390 keyboard scope disclosure and exact navigation keep focus without viewport overflow', async ({ page }, info) => {
      const words = logWords(language), journey = await openLogs(page, language); await readyLogs(page, language);
      expect(await logRows(page).first().locator('td').last().evaluate(cell => {
        const table = cell.closest('table');
        return table ? cell.getBoundingClientRect().width / table.getBoundingClientRect().width : 0;
      }), 'Narrow event guidance must retain readable row width').toBeGreaterThan(0.8);
      const toggle = page.getByRole('button', { name: words.scope, exact: true });
      await expect(toggle).toHaveAttribute('aria-expanded', 'false'); await tabTo(page, toggle); await page.keyboard.press('Enter');
      await expect(toggle).toBeFocused(); await expect(toggle).toHaveAttribute('aria-expanded', 'true');
      const controls = await toggle.getAttribute('aria-controls'); if (!controls) throw new Error('Disclosure association required');
      await expect(page.locator(`[id="${controls}"]`)).toBeVisible();
      const workspace = page.getByLabel(words.workspace, { exact: true }); await tabTo(page, workspace); await workspace.fill('W');
      const group = page.getByLabel(words.group, { exact: true }); await tabTo(page, group); await group.fill('G');
      await page.getByRole('combobox', { name: words.level, exact: true }).selectOption('warn');
      await page.getByRole('combobox', { name: words.source, exact: true }).selectOption('runtime');
      await page.getByRole('combobox', { name: words.stage, exact: true }).selectOption('retrying');
      await page.getByLabel(words.search, { exact: true }).fill('REF1');
      const apply = page.getByRole('button', { name: words.apply, exact: true }); await tabTo(page, apply); await page.keyboard.press('Enter');
      await expect.poll(() => snapshotCalls(journey.mock).length).toBe(2);
      expect(snapshotCalls(journey.mock).at(-1)?.params).toMatchObject({ level: 'warn', source: 'runtime', stage: 'retrying', q: 'REF1', workspace_id: 'W', group_id: 'G' });
      await readyLogs(page, language); await expect(logRows(page)).toHaveCount(1);
      const summary = logRows(page).first().locator('summary'); await tabTo(page, summary); await page.keyboard.press('Enter');
      await expect(summary).toBeFocused(); await expect(summary.locator('..')).toHaveAttribute('open', '');
      await expect(fact(logRows(page).first(), words.revision)).toHaveText('R1');
      await captureLogs(page, info, 'RL-N05', language, 'disclosure');
      const link = logRows(page).first().getByRole('link', { name: words.openGroup, exact: true });
      await tabTo(page, link); await page.keyboard.press('Enter'); await selectedGroup(page, 'G');
      await captureLogs(page, info, 'RL-N05', language, 'native-focus'); await finishLogs(info, journey);
    });
  });
}
