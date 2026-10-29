import { expect, test, type Page } from '@playwright/test';
import {
  RUNTIME_VIEWPORTS, actionCalls, captureRuntime, fact, finishRuntime, groupCard,
  inspectSource, openRuntime, sourceDetail, tabTo,
} from './runtime-recording.helpers';

async function nativeSelection(page: Page, step: number) {
  if (step === 1) {
    await expect(page.getByTestId('device-editor')).toBeVisible();
    await expect(page.getByTestId('editor-input-name')).toHaveValue('Device A');
  } else if (step === 2) {
    await expect(page.getByTestId('rule-editor')).toBeVisible();
    await expect(page.getByTestId('rule-device-select')).toHaveValue('A');
    await expect(page.getByTestId('rule-name-input')).toHaveValue('Rule A1');
    await expect(page.getByTestId('rule-start-input')).toHaveValue('40001');
  } else if (step === 3) {
    const row = page.getByTestId('mapping-row-PA1');
    await expect(page.getByTestId('step3-mapping-container')).toBeVisible();
    await expect(row).toBeVisible(); await expect(row).toContainText('40001');
    await expect(row).toHaveClass(/outline-blue-500/);
  } else {
    const card = page.getByTestId('group-list-item-G1');
    await expect(card).toBeVisible(); await expect(card).toHaveAttribute('data-locator-selected', 'true');
    await expect(page.getByTestId('group-open-G1')).toBeVisible();
    await expect(page.getByTestId('group-editor')).toHaveCount(0); // A locator focuses, never auto-opens an editor.
  }
}

for (const viewport of RUNTIME_VIEWPORTS) {
  test.describe(`runtime recording navigation ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    for (const step of [1, 2, 3, 4] as const) {
      test(`RR-N0${step} exact Step${step} native selection survives real Back/Forward without writes`, async ({ page }, info) => {
        const journey = await openRuntime(page); await inspectSource(page);
        const link = step === 1 ? page.getByRole('link', { name: 'Open device settings', exact: true })
          : step === 4 ? groupCard(page).getByRole('link', { name: 'Open saved group settings', exact: true })
            : sourceDetail(page).getByRole('link', { name: step === 2 ? 'Open rule settings' : 'Open mapping settings', exact: true });
        const href = await link.getAttribute('href'); expect(href).not.toBeNull();
        if (!href) throw new Error('An exact native settings link is required.');
        const expected = new URL(href, 'http://fixture.invalid');
        expect(expected.pathname).toBe('/studio/v2');
        expect(Object.fromEntries(expected.searchParams)).toEqual({ entry: 'locator', workspace_id: 'W', step: String(step), device_id: 'A',
          ...(step === 2 ? { rule_id: 'A1' } : step === 3 ? { mapping_id: 'MA1' } : step === 4 ? { group_id: 'G1', group_revision: 'B' } : {}) });
        await link.click(); await expect(page).toHaveURL(url => url.pathname === '/studio/v2' && url.search === expected.search);
        await nativeSelection(page, step);
        const entry = actionCalls(journey.mock, 'locator').at(-1);
        expect(entry?.params).toMatchObject(Object.fromEntries(expected.searchParams));
        expect(entry?.params.editor).toBeUndefined(); expect(entry?.params.readiness).toBeUndefined();
        await captureRuntime(page, info, `rr-n0${step}-native-exact-selection`);
        await page.goBack(); // Actual browser history, not setting React/router state directly.
        await expect(page).toHaveURL(url => url.pathname === '/studio/runtime' && url.searchParams.get('device_id') === 'A');
        await expect(fact(sourceDetail(page), 'Rule')).toHaveText('A1');
        await expect(fact(sourceDetail(page), 'Mapping')).toHaveText('MA1');
        await expect(fact(sourceDetail(page), 'Point')).toHaveText('PA1');
        await page.goForward(); await nativeSelection(page, step);
        expect(actionCalls(journey.mock, 'locator').every(call => call.params.workspace_id === 'W')).toBe(true);
        await finishRuntime(info, journey);
      });
    }

    test('RR-N05 deleted G1 never selects same-name G2; only proven workspace links are offered', async ({ page }, info) => {
      const journey = await openRuntime(page, 'deleted', { settingsQuery: 'entry=locator&workspace_id=W&step=4&group_id=G1' });
      const notice = page.getByTestId('studio-locator-notice');
      await expect(notice).toContainText('This item was removed. A same-name replacement has not been selected.');
      await expect(page.getByTestId('group-open-G2')).toHaveCount(0); await expect(page.getByTestId('group-editor')).toHaveCount(0);
      const owned = notice.getByRole('link', { name: 'Open device A', exact: true });
      await expect(owned).toHaveAttribute('href', /workspace_id=W(?:&|$)/);
      expect(actionCalls(journey.mock, 'locator').every(call => call.params.group_id === 'G1')).toBe(true);
      expect(journey.mock.ledger.some(call => call.path.includes('/G2'))).toBe(false);
      await captureRuntime(page, info, 'rr-n05-deleted-not-namesake', notice); await finishRuntime(info, journey);
    });

    test('RR-N06 foreign workspace fails closed without another workspace or first-device fallback', async ({ page }, info) => {
      const journey = await openRuntime(page, 'foreign', { settingsQuery: 'entry=locator&workspace_id=FOREIGN&step=2&device_id=A&rule_id=A1' });
      const notice = page.getByTestId('studio-locator-notice');
      await expect(notice).toContainText('This workspace or selection is inaccessible. No other workspace has been selected.');
      await expect(notice.getByRole('link')).toHaveCount(0); await expect(page.getByTestId('rule-editor')).toHaveCount(0);
      expect(actionCalls(journey.mock, 'locator')).toHaveLength(1);
      expect(actionCalls(journey.mock, 'locator')[0].params.workspace_id).toBe('FOREIGN');
      await captureRuntime(page, info, 'rr-n06-foreign-fail-closed', notice);
      await finishRuntime(info, journey, [{ action: 'locator', status: 403 }]);
    });

    for (const language of ['en', 'zh-TW'] as const) {
      test(`RR-N07 focused order and keyboard report disclosure remain usable in ${language}`, async ({ page }, info) => {
        const journey = await openRuntime(page, 'normal', { language });
        const reportName = language === 'en' ? 'Historical reports' : '歷史報表';
        const show = language === 'en' ? 'Show historical reports' : '顯示歷史報表';
        const hide = language === 'en' ? 'Hide historical reports' : '隱藏歷史報表';
        const reports = page.getByRole('region', { name: reportName, exact: true });
        await expect(reports).toBeVisible();
        const order = ['runtime-dashboard-setup-context', 'runtime-recording-rules', 'runtime-dashboard-summary-panel', 'runtime-dashboard-live-points-table', 'runtime-dashboard-logs-panel'];
        for (let index = 0; index < order.length - 1; index++) {
          expect(await page.getByTestId(order[index]).evaluate((element, next) => {
            const after = document.querySelector(`[data-testid="${next}"]`);
            return Boolean(after && element.compareDocumentPosition(after) & Node.DOCUMENT_POSITION_FOLLOWING);
          }, order[index + 1]), `${order[index]} comes before ${order[index + 1]}`).toBe(true);
        }
        const lastSection = order.at(-1);
        if (!lastSection) throw new Error('Focused section order is required.');
        expect(await page.getByTestId(lastSection).evaluate((element, name) => {
          const report = [...document.querySelectorAll('section')].find(item => item.getAttribute('aria-label') === name);
          return Boolean(report && element.compareDocumentPosition(report) & Node.DOCUMENT_POSITION_FOLLOWING);
        }, reportName)).toBe(true);
        const toggle = reports.getByRole('button', { name: show, exact: true });
        await expect(toggle).toHaveAttribute('aria-expanded', 'false');
        const control = await toggle.getAttribute('aria-controls'); expect(control).toBeTruthy();
        const content = page.locator(`[id="${control}"]`); await expect(content).toBeHidden();
        await expect.poll(() => actionCalls(journey.mock, 'reports').length).toBe(1);
        const reportReads = actionCalls(journey.mock, 'reports').length;
        await tabTo(page, toggle); await page.keyboard.press('Enter');
        const collapse = reports.getByRole('button', { name: hide, exact: true });
        await expect(collapse).toBeFocused(); await expect(collapse).toHaveAttribute('aria-expanded', 'true');
        await expect(content).toBeVisible();
        await captureRuntime(page, info, `rr-n07-reports-${language}`, collapse);
        await page.keyboard.press('Space'); await expect(toggle).toBeFocused(); await expect(content).toBeHidden();
        expect(actionCalls(journey.mock, 'reports')).toHaveLength(reportReads);
        await finishRuntime(info, journey);
      });
    }
  });
}
