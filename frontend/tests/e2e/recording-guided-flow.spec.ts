import { expect, test } from '@playwright/test';
import zh from '../../src/i18n/locales/zh-TW/workbench-v2-recording-guided.json' with { type: 'json' };
import {
  GUIDED_VIEWPORTS, actionCalls, advanced, attachGuidedEvidence, authoring, basic, captureGuided,
  expandGuidedSources, mutationCalls, noGuidedMutation, openGuided, prepareGuided, startGuided, tabTo,
} from './recording-guided.helpers';

// Formal browser cases only. Discovery/typechecking is not browser, SQL or runtime verification.
for (const viewport of GUIDED_VIEWPORTS) {
  test.describe(`guided recording flow ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('GD-F01 ordinary non-counter flow has separate save, preview, DDL consent and start consent', async ({ page }, info) => {
      const { mock } = await openGuided(page);
      const panel = authoring(page), interval = panel.getByTestId('circuit-profile-interval-instantaneous');
      const setupURL = page.url();
      await expect(panel.getByTestId('circuit-counter-confirmation')).toHaveCount(0);
      await expect(advanced(page)).not.toHaveAttribute('open');
      await expect(panel.getByTestId('recording-scope-summary')).toContainText('Mock PLC A');
      await expect(panel.getByTestId('recording-scope-summary')).toContainText('Line DB');
      await expect(panel.getByTestId('recording-scope-summary')).toContainText('Machine state');
      await expect(panel.getByTestId('recording-scope-summary')).not.toContainText('/mock-private/');
      await interval.fill('0');
      await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
      await panel.getByRole('button', { name: 'Review blocking controls', exact: true }).click();
      await expect(interval).toBeFocused();
      noGuidedMutation(mock);
      await interval.fill('15');
      const save = panel.getByTestId('circuit-profile-prepare');
      await expect(save).toBeEnabled(); await tabTo(page, save);
      await expect(panel.locator('button[data-primary="true"]')).toHaveCount(1);
      await captureGuided(page, info, 'new-recording-keyboard-primary', save);
      await page.keyboard.press('Enter');
      await expect.poll(() => actionCalls(mock, 'ensure').length).toBe(1);
      expect(actionCalls(mock, 'ensure')[0].status).toBe(200);
      expect(actionCalls(mock, 'schema-preview')).toHaveLength(0);
      expect(actionCalls(mock, 'schema-apply')).toHaveLength(0);
      expect(actionCalls(mock, 'start')).toHaveLength(0);
      const preview = panel.getByTestId('group-schema-preview');
      await expect(preview).toBeEnabled(); await preview.click();
      await expect(panel.getByTestId('group-schema-preview-result')).toBeVisible();
      expect(actionCalls(mock, 'schema-preview')).toHaveLength(1);
      expect(actionCalls(mock, 'schema-apply')).toHaveLength(0);
      await captureGuided(page, info, 'explicit-table-review-before-DDL', panel.getByTestId('group-schema-preview-result'));
      await panel.getByTestId('group-schema-apply').click();
      await expect.poll(() => actionCalls(mock, 'schema-apply').length).toBe(1);
      expect(actionCalls(mock, 'schema-apply')[0].status).toBe(200);
      const review = panel.getByTestId('circuit-profile-start');
      await expect(review).toBeEnabled(); await review.click();
      await expect(page.getByRole('dialog')).toBeVisible();
      await expect(page.getByRole('dialog')).toContainText('not SQL success');
      expect(actionCalls(mock, 'start')).toHaveLength(0);
      await page.getByTestId('recording-start-cancel').click();
      await expect(page.getByRole('dialog')).toHaveCount(0); await expect(review).toBeFocused();
      expect(actionCalls(mock, 'start')).toHaveLength(0);
      await startGuided(page, mock);
      await expect(panel.getByTestId('recording-completion-summary')).toContainText('Confirmed active devices');
      await expect(page).toHaveURL(setupURL);
      expect(actionCalls(mock, 'start')[0].status).toBe(200);
      expect(actionCalls(mock, 'start')[0].body).toMatchObject({ workspace_id: 'ws-1', device_ids: ['dev-1'],
        groups: [{ group_id: mock.snapshot().groups[0].id, expected_group_revision: mock.snapshot().groups[0].revision, expected_connector_revision: 'crev-1' }] });
      expect(mutationCalls(mock).map(call => call.action)).toEqual(['ensure', 'schema-preview', 'schema-apply', 'start']);
      await expect(advanced(page)).not.toHaveAttribute('open');
      await captureGuided(page, info, 'started-not-SQL-confirmed', panel.getByTestId('recording-completion-summary'));
      await attachGuidedEvidence(info, mock);
    });

    test('GD-F02 a failed schema preview exposes a keyboard-reachable retry without DDL or start', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'new', { setup: value => value.rejectNext('schema-preview') });
      const panel = authoring(page);
      await panel.getByTestId('circuit-profile-prepare').click();
      await expect.poll(() => actionCalls(mock, 'ensure').length).toBe(1);
      const preview = panel.getByTestId('group-schema-preview'); await expect(preview).toBeEnabled(); await preview.click();
      await expect(panel.getByTestId('group-schema-error')).toBeVisible();
      expect(actionCalls(mock, 'schema-preview')).toHaveLength(1);
      expect(actionCalls(mock, 'schema-preview')[0].status).toBeGreaterThanOrEqual(400);
      expect(actionCalls(mock, 'schema-apply')).toHaveLength(0); expect(actionCalls(mock, 'start')).toHaveLength(0);
      await panel.getByRole('button', { name: 'Inspect this scope', exact: true }).click();
      await expect(panel.getByTestId('recording-guided-preparation')).toBeFocused();
      await tabTo(page, preview); await captureGuided(page, info, 'schema-failure-actionable-retry', preview);
      await page.keyboard.press('Enter');
      await expect(panel.getByTestId('group-schema-preview-result')).toBeVisible();
      expect(actionCalls(mock, 'schema-preview')).toHaveLength(2);
      expect(actionCalls(mock, 'schema-preview')[1].status).toBe(200);
      noGuidedMutation(mock, ['ensure', 'schema-preview']);
      await attachGuidedEvidence(info, mock, ['schema-preview']);
    });

    test('GD-F03 failed original-operation lookup stays unknown through reload and never resends', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'unknown');
      const setupURL = page.url(); await startGuided(page, mock);
      const original = structuredClone(actionCalls(mock, 'start')[0].body);
      const check = basic(page).getByTestId('basic-recording-check'); await expect(check).toBeEnabled();
      mock.rejectNext('start-operation'); const before = actionCalls(mock, 'start-operation').length;
      await tabTo(page, check); await page.keyboard.press('Enter');
      await expect.poll(() => actionCalls(mock, 'start-operation').length).toBeGreaterThan(before);
      await expect(basic(page)).toContainText('The saved operation could not be read. Check again before retrying.');
      await expect(basic(page).getByTestId('basic-recording-continue-runtime')).toHaveCount(0);
      await expect(basic(page).getByTestId('circuit-profile-start')).toHaveCount(0);
      await captureGuided(page, info, 'lookup-failure-original-operation-only', check);
      await page.reload(); await expect(check).toBeEnabled(); await check.click();
      expect(actionCalls(mock, 'start-operation').every(call => new URL(call.url).pathname.endsWith('/O1'))).toBe(true);
      expect(actionCalls(mock, 'start')).toHaveLength(1); expect(actionCalls(mock, 'start')[0].body).toEqual(original);
      expect(mock.snapshot().originalStart).toEqual(original);
      await expect(page).toHaveURL(setupURL); noGuidedMutation(mock, ['start']);
      await attachGuidedEvidence(info, mock, ['start-operation']);
    });

    test('GD-F04 partial result names incomplete groups and only explicitly continues to confirmed A', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'partial');
      const setupURL = page.url(); await startGuided(page, mock);
      const completion = basic(page).getByTestId('recording-completion-summary');
      await expect(completion).toContainText('Partial');
      await expect(completion.getByRole('list', { name: 'Confirmed active devices', exact: true })).toHaveText('Mock PLC A');
      await expect(completion.getByRole('list', { name: 'Recording groups not completed', exact: true })).toContainText('Mock saved recording');
      await expect(completion).not.toContainText('Mock PLC B'); await expect(completion).toContainText('not proof of SQL commit');
      await expect(page).toHaveURL(setupURL);
      await expect(basic(page).getByTestId('basic-recording-check')).toHaveAttribute('data-primary', 'true');
      const proceed = completion.getByTestId('basic-recording-continue-runtime'); await expect(proceed).not.toHaveAttribute('data-primary');
      await captureGuided(page, info, 'partial-A-confirmed-group-incomplete', completion);
      await proceed.dblclick({ delay: 20 });
      await expect(page).toHaveURL(/\/studio\/runtime\?device_id=dev-1$/);
      await expect(page.getByTestId('fixture-runtime-destination')).toBeVisible();
      await expect(page.getByTestId('fixture-runtime-scope')).toHaveText('dev-1');
      expect(actionCalls(mock, 'start')).toHaveLength(1);
      await page.goBack(); await expect(page).toHaveURL(setupURL);
      await expect(basic(page).locator('[data-device-id="dev-1"]')).toHaveAttribute('aria-pressed', 'true');
      await expect(completion).toContainText('Recording groups not completed');
      expect(mock.snapshot().originalStart?.device_ids).toEqual(['dev-1']); noGuidedMutation(mock, ['start']);
      await captureGuided(page, info, 'returned-to-original-partial-scope', completion);
      await attachGuidedEvidence(info, mock);
    });

    test('GD-F05 same saved scope keeps missing effective evidence explicit and drawer restores focus', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'saved-long');
      const summary = basic(page).getByTestId('recording-scope-summary');
      await expect(summary).toContainText('Mock PLC A'); await expect(summary).toContainText('Saved circuit');
      await expect(summary.getByRole('region', { name: 'Saved configuration', exact: true })).toContainText('Unavailable');
      await expect(summary.getByRole('region', { name: 'Verified effective configuration', exact: true })).toContainText('Unavailable');
      await expect(summary).not.toContainText('/mock-private/'); await expect(summary).not.toContainText('recordings.sqlite');
      const details = summary.getByRole('button', { name: 'View complete rules and delivery evidence', exact: true });
      await tabTo(page, details); await page.keyboard.press('Enter');
      const drawer = page.getByRole('dialog', { name: 'Rules and evidence for this scope' });
      await expect(drawer).toBeVisible(); await expect(drawer).toContainText('G1'); await expect(drawer).toContainText('R1');
      await expect(drawer.getByRole('region', { name: 'Saved configuration details', exact: true })).toContainText('Saved schema references are not a current schema verification.');
      await expect(drawer.getByRole('region', { name: 'Verified effective details', exact: true })).toContainText('Unavailable');
      await captureGuided(page, info, 'saved-versus-unavailable-effective-drawer', drawer);
      await page.keyboard.press('Escape'); await expect(drawer).toHaveCount(0); await expect(details).toBeFocused();
      noGuidedMutation(mock); await attachGuidedEvidence(info, mock);
    });

    test('GD-F06 long source lists remain inspectable without document overflow or implicit saves', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'long-sources');
      const summary = basic(page).getByTestId('recording-scope-summary');
      const sources = summary.getByRole('list', { name: 'Sources and types', exact: true });
      await expect(sources.locator('li')).toHaveCount(40);
      expect(await sources.evaluate(element => element.scrollHeight > element.clientHeight)).toBe(true);
      await expandGuidedSources(page);
      await expect(authoring(page).getByTestId('circuit-profile-role-tag-39')).toBeVisible();
      await expect(summary).toContainText('Source 40'); await expect(summary).not.toContainText('/mock-private/');
      await captureGuided(page, info, 'long-source-list-and-single-primary', sources);
      await expect(authoring(page).locator('button[data-primary="true"]')).toHaveCount(1);
      noGuidedMutation(mock); await attachGuidedEvidence(info, mock);
    });

    test('GD-F07 Traditional Chinese normal flow preserves translated consent and counted stages', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'new', { language: 'zh-TW' });
      const panel = authoring(page), setupURL = page.url();
      await expect(panel.getByTestId('recording-scope-summary')).toContainText(zh.unsaved);
      await expect(panel.getByTestId('circuit-counter-confirmation')).toHaveCount(0);
      noGuidedMutation(mock); await prepareGuided(page, mock);
      expect(mutationCalls(mock).map(call => call.action)).toEqual(['ensure', 'schema-preview', 'schema-apply']);
      expect(mutationCalls(mock).every(call => call.status === 200)).toBe(true);
      const review = panel.getByTestId('circuit-profile-start'); await expect(review).toHaveText(zh.next.actions.review_scoped_start);
      await review.click(); const dialog = page.getByRole('dialog', { name: zh.next.startTitle });
      await expect(dialog).toBeVisible(); await expect(dialog).toContainText(zh.next.deliverySeparate);
      expect(actionCalls(mock, 'start')).toHaveLength(0);
      await page.getByTestId('recording-start-cancel').click(); await expect(review).toBeFocused();
      await startGuided(page, mock);
      const completion = panel.getByTestId('recording-completion-summary');
      await expect(completion).toContainText(zh.completion.heading); await expect(completion).toContainText(zh.completion.deliverySeparate);
      await expect(completion.getByRole('list', { name: zh.completion.activated, exact: true })).toHaveText('Mock PLC A');
      await expect(completion.getByTestId('basic-recording-continue-runtime')).toHaveText(zh.completion.continue);
      await expect(page).toHaveURL(setupURL);
      expect(mutationCalls(mock).map(call => call.action)).toEqual(['ensure', 'schema-preview', 'schema-apply', 'start']);
      expect(mock.snapshot().originalStart?.device_ids).toEqual(['dev-1']);
      await captureGuided(page, info, 'zh-TW-confirmed-start-separate-delivery', completion);
      await attachGuidedEvidence(info, mock);
    });

    test('GD-F08 simulated runtime evidence separates saved R2 at 30 seconds from effective R1 at 15 seconds', async ({ page }, info) => {
      info.annotations.push({ type: 'evidence-layer', description: 'Simulated safe runtime DTO; not actual runtime, device or SQL verification.' });
      const { mock } = await openGuided(page, 'saved-different');
      await expect(page.getByText('No database, device, SQL or effective runtime verification. Scenario: saved-different', { exact: true })).toBeVisible();
      const summary = basic(page).getByTestId('recording-scope-summary');
      const saved = summary.getByRole('region', { name: 'Saved configuration', exact: true });
      const effective = summary.getByRole('region', { name: 'Verified effective configuration', exact: true });
      await expect(summary).toContainText('Mock PLC A');
      await expect(saved.getByText('30 seconds', { exact: true })).toBeVisible();
      await expect(effective.getByText('15 seconds', { exact: true })).toBeVisible();
      await expect(saved).not.toContainText('15 seconds'); await expect(effective).not.toContainText('30 seconds');
      await expect(saved).toContainText('Saved configuration is not activation or SQL delivery.');
      await expect(effective).toContainText('retained history does not prove active intake or SQL success.');
      await expect(summary).not.toContainText('/mock-private/'); await expect(summary).not.toContainText('recordings.sqlite');
      expect(actionCalls(mock, 'runtime-context').some(call => call.status === 200)).toBe(true);
      noGuidedMutation(mock);
      await captureGuided(page, info, 'SIMULATED-saved-30-effective-15-distinct', summary);
      const details = summary.getByRole('button', { name: 'View complete rules and delivery evidence', exact: true });
      await tabTo(page, details); await page.keyboard.press('Enter');
      const drawer = page.getByRole('dialog', { name: 'Rules and evidence for this scope' });
      for (const [label, value] of [['Saved revision', 'R2'], ['Latest requested revision', 'R1'],
        ['Effective revision at observation', 'R1'], ['Server observation', '2030-01-01T00:00:00.000Z']] as const) {
        await expect(drawer.locator('dt').filter({ hasText: new RegExp(`^${label}$`) }).locator('..').locator('dd')).toHaveText(value);
      }
      await page.keyboard.press('Escape'); await expect(drawer).toHaveCount(0); await expect(details).toBeFocused();
      noGuidedMutation(mock); await attachGuidedEvidence(info, mock);
    });
  });
}
