import { expect, test } from '@playwright/test';
import {
  RUNTIME_VIEWPORTS, actionCalls, captureRuntime, fact, finishRuntime, groupCard, inspectSource,
  openDisclosure, openRuntime, selectDevice, sourceDetail, sourceMappings, sourceRules,
} from './runtime-recording.helpers';

for (const viewport of RUNTIME_VIEWPORTS) {
  test.describe(`runtime recording rules ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('RR-R01 six rules/ten mappings remain inspectable; delayed A6 never replaces device B', async ({ page }, info) => {
      const journey = await openRuntime(page);
      await expect(sourceRules(page).getByRole('button', { name: /^Show mappings for / })).toHaveCount(6);
      await expect(sourceMappings(page).getByRole('button', { name: /^Inspect source / })).toHaveCount(10);
      expect(actionCalls(journey.mock, 'detail')).toHaveLength(0);
      const pending = journey.holdNext('detail');
      await sourceRules(page).getByRole('button', { name: 'Show mappings for A6', exact: true }).click();
      await sourceMappings(page).getByRole('button', { name: 'Inspect source 40051', exact: true }).click();
      await pending.reached;
      expect(actionCalls(journey.mock, 'detail').at(-1)?.params).toMatchObject({ device_id: 'A', rule_id: 'A6', mapping_id: 'MA6', point_id: 'PA6' });
      await selectDevice(page, 'B'); pending.release();
      await expect(sourceRules(page).getByRole('button', { name: 'Show mappings for B1', exact: true })).toBeVisible();
      await expect(sourceRules(page).getByRole('button', { name: /^Show mappings for / })).toHaveCount(1);
      await expect(sourceDetail(page)).toHaveCount(0);
      await expect(sourceMappings(page)).not.toContainText('PA6');
      expect(actionCalls(journey.mock, 'detail')).toHaveLength(1);
      await captureRuntime(page, info, 'rr-r01-late-a6-device-b', sourceRules(page)); await finishRuntime(info, journey);
    });

    test('RR-R02 polling1000ms and recording mean15s stay separate; saved R2 scale0.1 is not effective R1', async ({ page }, info) => {
      const journey = await openRuntime(page); await inspectSource(page);
      const details = sourceDetail(page);
      await expect(fact(details, 'Saved rule revision')).toHaveText('R2');
      await expect(fact(details, 'Register address')).toHaveText('40001');
      await expect(fact(details, 'Polling interval and source')).toContainText('1000 ms');
      await expect(details).toContainText('mean / 15 seconds');
      await expect(details).toContainText('× 0.1');
      await expect(details).toContainText('Saved settings do not prove the configuration currently executing.');
      await expect(details).toContainText('Effective source configuration unavailable');
      const request = actionCalls(journey.mock, 'detail').at(-1);
      expect(request?.params).toMatchObject({ workspace_id: 'W', workspace_revision: 'W1', device_id: 'A', rule_id: 'A1', rule_revision: 'R2', point_id: 'PA1', mapping_id: 'MA1' });
      await captureRuntime(page, info, 'rr-r02-source-saved-only', details);
      await details.getByRole('button', { name: 'Close source details', exact: true }).click();
      await expect(sourceMappings(page).getByRole('button', { name: 'Inspect source 40001', exact: true })).toBeFocused();
      await finishRuntime(info, journey);
    });

    test('RR-R03 draft B meter_new30 never relabels effective A meter_old15; saved C and requested B remain distinct', async ({ page }, info) => {
      const journey = await openRuntime(page); const card = groupCard(page);
      const saved = card.getByRole('region', { name: 'Saved configuration', exact: true });
      const effective = card.getByRole('region', { name: 'Effective recording snapshot', exact: true });
      await expect(fact(saved, 'Table')).toHaveText('meter_new'); await expect(fact(saved, 'Recording interval')).toHaveText('30 seconds');
      await expect(fact(effective, 'Table')).toHaveText('meter_old'); await expect(fact(effective, 'Recording interval')).toHaveText('15 seconds');
      journey.mock.setScenario('draft'); await page.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      await expect(fact(saved, 'Fact revision')).toHaveText('C'); await expect(fact(saved, 'Table')).toHaveText('meter_draft');
      await expect(fact(card.getByRole('region', { name: 'Latest requested revision', exact: true }), 'Requested revision')).toHaveText('B');
      await expect(fact(card.getByRole('region', { name: 'Latest requested revision', exact: true }), 'Requested effective boundary')).toHaveText('Unavailable');
      await expect(fact(effective, 'Fact revision')).toHaveText('A'); await expect(fact(effective, 'Table')).toHaveText('meter_old');
      journey.mock.setScenario('no-proof'); await page.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      await expect(fact(saved, 'Table')).toHaveText('meter_new');
      await expect(effective).toContainText('Effective recording facts unavailable');
      await expect(effective).not.toContainText('meter_old');
      await captureRuntime(page, info, 'rr-r03-saved-requested-effective', card); await finishRuntime(info, journey);
    });

    test('RR-R04 mixed W1/W2 evidence fails stale and only a read-only refresh recovers', async ({ page }, info) => {
      const journey = await openRuntime(page, 'stale');
      const panel = page.getByTestId('runtime-recording-rules');
      await expect(panel.getByRole('alert')).toContainText('Recording configuration changed');
      await expect(groupCard(page)).toHaveCount(0);
      journey.mock.setScenario('normal'); await panel.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      await expect(groupCard(page)).toContainText('meter_old');
      expect(actionCalls(journey.mock, 'evidence')).toHaveLength(2);
      await captureRuntime(page, info, 'rr-r04-stale-read-recovery', panel); await finishRuntime(info, journey);
    });

    test('RR-R05 disabled C2 draft retains C1 revision7 old_table three pending output rows', async ({ page }, info) => {
      const journey = await openRuntime(page, 'backlog'); const card = groupCard(page);
      const backlog = await openDisclosure(card, 'Backlog revision A0');
      await expect(fact(backlog, 'Connector ID')).toHaveText('C1'); await expect(fact(backlog, 'Connector revision')).toHaveText('7');
      await expect(fact(backlog, 'Table')).toHaveText('old_table'); await expect(fact(backlog, 'Pending output rows')).toHaveText('3 output rows');
      await expect(card.getByRole('region', { name: 'Saved configuration', exact: true })).toContainText('C2');
      await expect(backlog).not.toContainText('C2');
      await expect(fact(card.getByRole('region', { name: 'Current revision delivery', exact: true }), 'Observed intake state')).toHaveText('Retiring');
      await captureRuntime(page, info, 'rr-r05-original-destination-backlog', backlog); await finishRuntime(info, journey);
    });

    test('RR-R06 unknown delta, partial mean and missing source receipt never become zero or complete', async ({ page }, info) => {
      const journey = await openRuntime(page, 'uncertain'); const card = groupCard(page);
      const sources = card.getByRole('region', { name: 'Source evidence for selected device A', exact: true });
      const summaries = sources.locator('summary'); await expect(summaries).toHaveCount(2);
      for (const summary of await summaries.all()) await summary.click();
      await expect(sources).toContainText('Unknown; not zero'); await expect(sources).toContainText('Partial');
      await expect(sources.getByRole('region', { name: 'Latest observed source result', exact: true }).first()).toContainText('Unknown');
      await expect(sources).not.toContainText('Complete');
      const records = sources.locator('dt').filter({ hasText: /^Record ID$/ }).locator('..').locator('dd');
      await expect(records).toHaveCount(4);
      for (const record of await records.all()) await expect(record).toHaveText('Unavailable');
      const phases = card.getByRole('region', { name: 'Test-write evidence', exact: true });
      for (const label of ['Test-write operation ID', 'Test write recorded at', 'Independent readback verified at', 'Owned cleanup completed at']) await expect(fact(phases, label)).toHaveText('Unavailable');
      await captureRuntime(page, info, 'rr-r06-uncertain-source-evidence', sources); await finishRuntime(info, journey);
    });

    test('RR-R07 shared ten rows are not device commits; missing timing never becomes a15-second deadline', async ({ page }, info) => {
      const journey = await openRuntime(page, 'shared'); const card = groupCard(page);
      const delivery = card.getByRole('region', { name: 'Current revision delivery', exact: true });
      await expect(delivery).toContainText('Shared group totals across devices; not device or Tag success counts.');
      await expect(fact(delivery, 'SQL committed')).toHaveText('10 output rows');
      await expect(card.getByRole('region', { name: 'Source evidence for selected device A', exact: true })).toContainText('Selected-source evidence unavailable');
      journey.mock.setScenario('shared'); await page.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      const timing = card.getByRole('region', { name: 'Observed timing', exact: true });
      for (const label of ['Next due at', 'Next retry at', 'Attempt duration']) await expect(fact(timing, label)).toHaveText('Unavailable');
      await page.clock.setFixedTime(new Date('2026-10-29T10:00:01Z'));
      await expect(timing).not.toContainText('10:00:15'); await expect(timing).toContainText('not a live clock or countdown');
      await captureRuntime(page, info, 'rr-r07-shared-and-null-timing', timing); await finishRuntime(info, journey);
    });

    test('RR-R08 120 rules use bounded pages/search and disclosure creates no per-row request fan-out', async ({ page }, info) => {
      const journey = await openRuntime(page, 'many'); const rules = sourceRules(page);
      await expect(rules.getByRole('button', { name: /^Show mappings for / })).toHaveCount(50);
      await rules.getByRole('button', { name: 'Next source rules page', exact: true }).click();
      await expect(rules).toContainText('Page 2'); await expect(rules.getByRole('button', { name: /^Show mappings for / })).toHaveCount(50);
      await rules.getByRole('button', { name: 'Next source rules page', exact: true }).click();
      await expect(rules).toContainText('Page 3'); await expect(rules.getByRole('button', { name: /^Show mappings for / })).toHaveCount(20);
      expect(actionCalls(journey.mock, 'detail')).toHaveLength(0);
      const before = actionCalls(journey.mock, 'rules').length;
      await rules.getByRole('textbox', { name: 'Search source rules', exact: true }).fill('A120');
      expect(actionCalls(journey.mock, 'rules')).toHaveLength(before);
      await rules.getByRole('button', { name: 'Search rules', exact: true }).click();
      await expect(rules.getByRole('button', { name: /^Show mappings for / })).toHaveCount(1);
      await expect(rules.getByRole('button', { name: 'Show mappings for A120', exact: true })).toBeVisible();
      expect(actionCalls(journey.mock, 'rules')).toHaveLength(4);
      expect(actionCalls(journey.mock, 'rules').every(call => call.params.limit === '50')).toBe(true);
      expect(actionCalls(journey.mock, 'evidence')).toHaveLength(1);
      expect(actionCalls(journey.mock, 'detail')).toHaveLength(0);
      // Real SQLite Go tests, not this mock, prove limit500 clamps100.
      await captureRuntime(page, info, 'rr-r08-bounded-page-search', rules); await finishRuntime(info, journey);
    });

    test('RR-R10 quota retains group/global all-revision scope and overshoot; unknown storage is not healthy or repair permission', async ({ page }, info) => {
      const journey = await openRuntime(page, 'quota-group'); const card = groupCard(page);
      const quota = card.getByRole('region', { name: 'Quota and storage evidence', exact: true });
      const storage = quota.getByRole('region', { name: 'Frozen-destination storage history', exact: true });
      await expect(fact(quota, 'Quota scope')).toHaveText('This group — all revisions');
      await expect(fact(quota, 'Quota state')).toHaveText('Warning');
      await expect(fact(quota, 'Used durable bytes')).toHaveText('80');
      await expect(fact(quota, 'Configured maximum bytes')).toHaveText('100');
      await expect(fact(quota, 'New intake refused by quota')).toHaveText('No');
      await expect(quota).toContainText('Quota bytes are not selected-device or Tag totals and do not prove SQL delivery or intake permission.');
      const delivery = card.getByRole('region', { name: 'Current revision delivery', exact: true });
      await expect(delivery).toContainText('Shared group totals across devices; not device or Tag success counts.');
      await expect(fact(delivery, 'SQL committed')).toHaveText('0 output rows');
      journey.mock.setScenario('quota-global'); await page.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      await expect(fact(quota, 'Quota scope')).toHaveText('Global — all groups and revisions');
      await expect(fact(quota, 'Quota state')).toHaveText('Hard limit');
      await expect(fact(quota, 'Used durable bytes')).toHaveText('125'); // An observed overshoot is not clamped to100.
      await expect(fact(quota, 'Configured maximum bytes')).toHaveText('100');
      await expect(fact(quota, 'New intake refused by quota')).toHaveText('Yes');
      await expect(storage).toContainText('Frozen-destination storage history and events unavailable.');
      await expect(storage.locator('li')).toHaveCount(0); await expect(storage.getByRole('button')).toHaveCount(0);
      await captureRuntime(page, info, 'rr-r10-global-quota-overshoot-null-storage', quota);
      journey.mock.setScenario('shared'); await page.getByRole('button', { name: 'Refresh recording evidence', exact: true }).click();
      for (const label of ['Quota scope', 'Quota state', 'Used durable bytes', 'Configured maximum bytes', 'New intake refused by quota']) await expect(fact(quota, label)).toHaveText('Unavailable');
      await expect(storage).toContainText('Frozen-destination storage history and events unavailable.');
      await expect(storage.locator('li')).toHaveCount(0); await expect(storage.getByRole('button')).toHaveCount(0);
      await finishRuntime(info, journey);
    });

    test('RR-R09 hostile connector data stays out of operator DOM/URL; safe registers/table survive without log proxy', async ({ page }, info) => {
      const journey = await openRuntime(page, 'hostile'); await inspectSource(page);
      await expect(sourceDetail(page)).toContainText('40001'); await expect(groupCard(page)).toContainText('meter_old');
      const hostile = /fixture-secret|fixture-password|db\.private|\/private\/|Server=localhost|\[::1\]:502|driver-secret/i;
      await expect(page.locator('body')).not.toContainText(hostile);
      expect(await page.locator('body').innerHTML(), 'Attributes and hidden DOM must not become a secret side channel').not.toMatch(hostile);
      await expect(groupCard(page)).toContainText('安全資料表');
      for (const link of await page.locator('a[href]').all()) expect(await link.getAttribute('href')).not.toMatch(hostile);
      expect(page.url()).not.toMatch(hostile);
      const history = await page.evaluate(() => JSON.stringify(window.history.state)); expect(history).not.toMatch(hostile);
      expect(actionCalls(journey.mock, 'locator')).toHaveLength(0);
      await captureRuntime(page, info, 'rr-r09-safe-monitor', sourceDetail(page)); await finishRuntime(info, journey);
    });
  });
}
