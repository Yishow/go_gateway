import { expect, test } from '@playwright/test';
import {
  GUIDED_VIEWPORTS, actionCalls, attachGuidedEvidence, authoring, basic, captureGuided, expandGuidedSources,
  noGuidedMutation, openGuided, openGuidedAdvanced, selectGuidedDevice, startGuided, tabTo,
} from './recording-guided.helpers';

for (const viewport of GUIDED_VIEWPORTS) {
  test.describe(`guided navigation ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('GD-N01 Cancel, 30-second draft, B → A, Advanced and reload preserve scope without mutations', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'saved-long');
      const panel = authoring(page), review = panel.getByTestId('circuit-profile-start');
      await expect(review).toBeEnabled(); await review.click(); await expect(page.getByRole('dialog')).toBeVisible();
      await page.getByTestId('recording-start-cancel').click(); await expect(review).toBeFocused();
      const interval = panel.getByTestId('circuit-profile-interval-instantaneous'); await interval.fill('30');
      await selectGuidedDevice(page, 'dev-2'); await selectGuidedDevice(page, 'dev-1'); await expect(interval).toHaveValue('30');
      const advanced = await openGuidedAdvanced(page); await advanced.locator(':scope > summary').click();
      await expect(interval).toHaveValue('30'); noGuidedMutation(mock);
      await page.reload(); await expect(authoring(page).getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30');
      await expect(basic(page).locator('[data-device-id="dev-1"]')).toHaveAttribute('aria-pressed', 'true');
      noGuidedMutation(mock); await captureGuided(page, info, 'cancel-draft-device-reload', interval); await attachGuidedEvidence(info, mock);
    });

    test('GD-N02 keyboard Tab reaches primary; Cancel/Escape and details restore current scope focus', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'saved-long');
      const review = authoring(page).getByTestId('circuit-profile-start'); await expect(review).toBeEnabled();
      await tabTo(page, review); await page.keyboard.press('Enter');
      await expect(page.getByRole('dialog', { name: 'Confirm this recording scope' })).toBeVisible();
      await page.keyboard.press('Escape'); await expect(review).toBeFocused();
      const details = page.getByRole('button', { name: 'View complete rules and delivery evidence' });
      await details.click(); await expect(page.getByRole('dialog', { name: 'Rules and evidence for this scope' })).toBeVisible();
      await page.keyboard.press('Escape'); await expect(details).toBeFocused();
      await captureGuided(page, info, 'keyboard-dialog-details-focus', details); noGuidedMutation(mock); await attachGuidedEvidence(info, mock);
    });

    test('GD-N03 explicit empty device stays empty through disclosure and reload, never borrows A sources', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'empty-device'); await selectGuidedDevice(page, 'dev-empty');
      await expect(page.getByTestId('recording-scope-summary')).toContainText('Empty device');
      await expandGuidedSources(page); await expect(authoring(page).getByTestId('circuit-profile-member-list').getByRole('checkbox')).toHaveCount(0);
      await openGuidedAdvanced(page); await page.reload();
      await expect(basic(page).locator('[data-device-id="dev-empty"]')).toHaveAttribute('aria-pressed', 'true');
      await expect(page.getByTestId('recording-scope-summary')).toContainText('Empty device');
      await expect(authoring(page).getByTestId('circuit-profile-prepare')).toHaveCount(0);
      noGuidedMutation(mock); await captureGuided(page, info, 'empty-selected-device'); await attachGuidedEvidence(info, mock);
    });

    test('GD-N04 unknown original O1 stays query-only through reload and A → B → A', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'unknown'); await startGuided(page, mock);
      const original = structuredClone(actionCalls(mock, 'start')[0].body);
      await expect(authoring(page).getByTestId('basic-recording-check')).toBeEnabled();
      await page.reload(); await selectGuidedDevice(page, 'dev-2'); await selectGuidedDevice(page, 'dev-1');
      const query = authoring(page).getByTestId('basic-recording-check'); await expect(query).toBeEnabled();
      await expect(query).toHaveAttribute('data-primary', 'true');
      await expect(authoring(page).getByTestId('circuit-profile-start')).toHaveCount(0);
      const reads = actionCalls(mock, 'start-operation').length; await query.click();
      await expect.poll(() => actionCalls(mock, 'start-operation').length).toBe(reads + 1);
      expect(actionCalls(mock, 'start-operation').every(call => call.url.endsWith('/O1'))).toBe(true);
      expect(actionCalls(mock, 'start')).toHaveLength(1); expect(actionCalls(mock, 'start')[0].body).toEqual(original);
      noGuidedMutation(mock, ['start']); await captureGuided(page, info, 'unknown-original-query', query); await attachGuidedEvidence(info, mock);
    });

    test('GD-N05 lost ACK without operation ID keeps manual recovery and never fabricates an ID or resends', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'saved-long', { setup: value => value.loseNext('start') });
      await startGuided(page, mock); await expect(authoring(page).getByRole('button', { name: 'Inspect operation recovery' })).toBeVisible();
      await page.reload(); await expect(authoring(page).getByRole('button', { name: 'Inspect operation recovery' })).toBeVisible();
      await expect(authoring(page).getByTestId('circuit-profile-start')).toHaveCount(0);
      const recovery = authoring(page).getByRole('button', { name: 'Inspect operation recovery' }); await recovery.click();
      expect(actionCalls(mock, 'start')).toHaveLength(1); expect(actionCalls(mock, 'start-operation')).toHaveLength(0);
      noGuidedMutation(mock, ['start']); await captureGuided(page, info, 'lost-ack-manual-guidance', recovery); await attachGuidedEvidence(info, mock);
    });

    test('GD-N06 late accepted POST after A → B → A exposes original O1 without another start', async ({ page }, info) => {
      const { mock, holdNext } = await openGuided(page, 'unknown'); const pending = holdNext('start');
      await startGuided(page, mock); await selectGuidedDevice(page, 'dev-2'); await selectGuidedDevice(page, 'dev-1');
      expect(actionCalls(mock, 'start')).toHaveLength(1); pending.release();
      const query = authoring(page).getByTestId('basic-recording-check'); await expect(query).toBeEnabled();
      await expect.poll(() => actionCalls(mock, 'start-operation').length).toBeGreaterThan(0);
      expect(actionCalls(mock, 'start-operation').every(call => call.url.endsWith('/O1'))).toBe(true);
      expect(actionCalls(mock, 'start')).toHaveLength(1); noGuidedMutation(mock, ['start']);
      await captureGuided(page, info, 'late-ack-original-operation', query); await attachGuidedEvidence(info, mock);
    });

    test('GD-N07 source repair focuses native controls; long lists scroll without losing selection', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'long-sources'); await expandGuidedSources(page);
      const panel = authoring(page), list = panel.getByTestId('circuit-profile-member-list');
      expect(await list.evaluate(element => element.scrollHeight > element.clientHeight)).toBe(true);
      const lastSource = list.getByRole('checkbox', { name: 'Source 40', exact: true }); await lastSource.scrollIntoViewIfNeeded();
      expect(await list.evaluate(element => element.scrollTop)).toBeGreaterThan(0);
      await lastSource.uncheck(); await panel.getByTestId('circuit-profile-sources-toggle').click(); await expandGuidedSources(page);
      await expect(lastSource).not.toBeChecked();
      await panel.getByTestId('circuit-profile-toggle-cumulative').check();
      await panel.getByTestId('circuit-profile-sources-toggle').click();
      const repair = panel.getByTestId('circuit-profile-cumulative-choose-sources'); await repair.click();
      await expect(panel.getByTestId('circuit-profile-sources-toggle')).toBeFocused();
      await expect(panel.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'true');
      noGuidedMutation(mock); await captureGuided(page, info, 'source-scroll-and-repair-focus', repair); await attachGuidedEvidence(info, mock);
    });

    test('GD-N08 URL-independent typed selection distinguishes current, historical hint, stale and foreign without mutation', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'saved-long');
      const harness = page.getByTestId('fixture-selection-contract'); await harness.locator(':scope > summary').click();
      await page.getByTestId('fixture-selection-owned').click(); await expect(page.getByTestId('fixture-selection-result')).toHaveText('selected:dev-1:G1:R1');
      await page.getByTestId('fixture-selection-historical').click(); await expect(page.getByTestId('fixture-selection-result')).toHaveText('selected:dev-1:G1:R1:history-different');
      await page.getByTestId('fixture-selection-stale').click(); await expect(page.getByTestId('fixture-selection-result')).toHaveText('unavailable:evidence_stale');
      await page.getByTestId('fixture-selection-foreign').click(); await expect(page.getByTestId('fixture-selection-result')).toContainText('unavailable:');
      await expect(basic(page).locator('[data-device-id="dev-1"]')).toHaveAttribute('aria-pressed', 'true');
      noGuidedMutation(mock); await captureGuided(page, info, 'selection-contract-not-runtime-url', harness); await attachGuidedEvidence(info, mock);
    });

    test('GD-N09 offline preserves local drafts/navigation but disables network commands', async ({ page, context }, info) => {
      const { mock } = await openGuided(page, 'saved-long'); await context.setOffline(true);
      try {
        const panel = authoring(page), review = panel.getByTestId('circuit-profile-start'); await expect(review).toBeDisabled();
        await panel.getByTestId('circuit-profile-interval-instantaneous').fill('30');
        await openGuidedAdvanced(page); await selectGuidedDevice(page, 'dev-2'); await selectGuidedDevice(page, 'dev-1');
        await expect(panel.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30');
        noGuidedMutation(mock); await captureGuided(page, info, 'offline-local-draft');
      } finally { await context.setOffline(false); }
      noGuidedMutation(mock); await attachGuidedEvidence(info, mock);
    });
  });
}
