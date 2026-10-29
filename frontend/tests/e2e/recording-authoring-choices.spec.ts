import { expect, test } from '@playwright/test';
import {
  AUTHORING_VIEWPORTS, actionCalls, attachAuthoringEvidence, authoring, captureAuthoring, ensureGroup,
  expandAuthoringSources, noAuthoringMutation, object, openAuthoring, selectDevice,
} from './recording-authoring.helpers';

for (const viewport of AUTHORING_VIEWPORTS) {
  test.describe(`authoring choices ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('AR-C01 numeric state authors latest role; gauge mean and nonnumeric roles remain distinct', async ({ page }, info) => {
      const { mock } = await openAuthoring(page);
      await expandAuthoringSources(page);
      const panel = authoring(page), role = panel.getByTestId('circuit-profile-role-tag-0');
      await expect(role).toBeEnabled();
      await expect(role.locator('option[value="state"]')).toContainText(/latest/i);
      await role.selectOption('state');
      await expect(role).toHaveValue('state');
      await expect(panel.getByTestId('circuit-profile-role-tag-1')).toHaveValue('gauge');
      await expect(panel.getByTestId('circuit-profile-role-tag-2').locator('option')).toHaveCount(1);
      await expect(panel.getByTestId('circuit-profile-circuit-name')).toHaveValue('');
      await expect(panel.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'false');
      await expect(panel.getByTestId('circuit-profile-prepare')).toBeEnabled();
      noAuthoringMutation(mock);
      await captureAuthoring(page, info, 'numeric-state-latest-and-gauge', role);
      await panel.getByTestId('circuit-profile-prepare').click();
      await expect.poll(() => actionCalls(mock, 'ensure').length).toBe(1);
      await expect(panel.getByTestId('group-schema-preview-result')).toBeVisible();
      const group = ensureGroup(mock), members = group.members as Array<Record<string, unknown>>;
      expect(members.find(member => member.tag_id === 'tag-0')?.value_role).toBe('state');
      expect(members.find(member => member.tag_id === 'tag-1')?.value_role).toBe('gauge');
      expect(object(group.recording_profile).sample_reduction).toBe('mean');
      noAuthoringMutation(mock, ['ensure', 'schema-preview']);
      // The real backend's 1→2 latest=2 / gauge=1.5 proof is separate; this mock never computes it.
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-C02 instant-only counter exposes required-name focus and never infers confirmed meaning from a name', async ({ page }, info) => {
      const { mock } = await openAuthoring(page);
      await expandAuthoringSources(page);
      const panel = authoring(page);
      await panel.getByTestId('circuit-profile-role-tag-0').selectOption('counter_raw');
      await expect(panel.getByTestId('circuit-profile-toggle-cumulative')).not.toBeChecked();
      const name = panel.getByTestId('circuit-profile-circuit-name');
      const confirm = panel.getByTestId('circuit-confirm-counter-submit');
      await expect(name).toHaveAttribute('aria-required', 'true');
      await expect(name).toHaveAccessibleDescription('Enter a circuit name before confirming these counter sources.');
      await expect(confirm).toBeDisabled();
      await expect(confirm).toHaveAccessibleDescription('Enter a circuit name before confirming these counter sources.');
      await panel.getByTestId('circuit-confirm-counter-name-focus').click();
      await expect(name).toBeFocused();
      await captureAuthoring(page, info, 'instant-counter-required-name', name);
      await name.fill('Counter A');
      const circuitId = await page.evaluate(() => JSON.parse(sessionStorage.getItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1') ?? '{}').circuit_id as string);
      expect(circuitId).toBeTruthy();
      await name.fill('Renamed Counter');
      expect(await page.evaluate(() => JSON.parse(sessionStorage.getItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1') ?? '{}').circuit_id as string)).toBe(circuitId);
      await expect(name).toHaveAttribute('aria-required', 'false');
      await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
      await expect(panel.getByTestId('circuit-counter-confirmation')).toBeVisible();
      noAuthoringMutation(mock);
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-C03 new cumulative picker begins empty and explicit selection survives toggle and reload', async ({ page }, info) => {
      const { mock } = await openAuthoring(page);
      const panel = authoring(page), toggle = panel.getByTestId('circuit-profile-toggle-cumulative');
      await toggle.check();
      await expect(panel.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'true');
      await expect(panel.getByTestId('circuit-profile-cumulative-selected')).toHaveText('Selected cumulative sources: 0');
      await expect(panel.getByTestId('circuit-profile-cumulative-sources-error')).toContainText('Choose the counters');
      for (const tag of ['tag-0', 'tag-1', 'tag-2']) await expect(panel.getByTestId(`circuit-profile-member-cum-${tag}`)).not.toBeChecked();
      await captureAuthoring(page, info, 'empty-cumulative-picker', panel.getByTestId('circuit-profile-cumulative-sources-error'));
      await panel.getByTestId('circuit-profile-member-cum-tag-0').check();
      await toggle.uncheck();
      await toggle.check();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-1')).not.toBeChecked();
      await page.reload();
      await expandAuthoringSources(page);
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-2')).not.toBeChecked();
      noAuthoringMutation(mock);
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-C04 instant-only measurement retry preserves drafts and late A response cannot overwrite B', async ({ page }, info) => {
      const { mock, holdNext } = await openAuthoring(page, 'new-sqlite', { setup: value => value.setMeasurementError('dev-1', true) });
      await expandAuthoringSources(page);
      const panel = authoring(page);
      await panel.getByTestId('circuit-profile-role-tag-0').selectOption('counter_raw');
      await panel.getByTestId('circuit-profile-circuit-name').fill('Counter A');
      const error = panel.getByTestId('circuit-profile-measurements-error');
      await expect(error).toContainText('Your draft is preserved');
      await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
      await captureAuthoring(page, info, 'instant-only-counter-retry', error);
      const before = actionCalls(mock, 'measurements:dev-1').length;
      mock.setMeasurementError('dev-1', false);
      const gate = holdNext('measurements:dev-1');
      try {
        await error.getByRole('button', { name: 'Retry', exact: true }).click();
        await expect.poll(() => actionCalls(mock, 'measurements:dev-1').length).toBe(before + 1);
        await panel.getByTestId('circuit-profile-circuit-name').fill('A edited while retrying');
        await selectDevice(page, 'dev-2');
        await panel.getByTestId('circuit-profile-circuit-name').fill('Draft B');
        gate.release();
        await expandAuthoringSources(page);
        await expect(panel.getByTestId('circuit-profile-role-b-tag-0')).toHaveValue('gauge');
        await expect(panel.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft B');
        await expect(error).toHaveCount(0);
        await selectDevice(page, 'dev-1');
        await expandAuthoringSources(page);
        await expect(panel.getByTestId('circuit-profile-role-tag-0')).toHaveValue('counter_raw');
        await expect(panel.getByTestId('circuit-profile-circuit-name')).toHaveValue('A edited while retrying');
      } finally { gate.release(); }
      noAuthoringMutation(mock);
      await attachAuthoringEvidence(info, mock, ['measurements:dev-1']);
    });

    for (const engine of ['sqlite', 'postgres', 'mysql', 'mariadb', 'sqlserver'] as const) {
      test(`AR-C05 ${engine} visible layout gate and outbound authored layout agree`, async ({ page }, info) => {
        const { mock } = await openAuthoring(page, `new-${engine}`);
        const panel = authoring(page);
        const wide = panel.getByTestId('circuit-profile-purpose-wide'), long = panel.getByTestId('circuit-profile-purpose-long');
        if (engine === 'sqlite' || engine === 'postgres') await expect(wide).toBeEnabled();
        else {
          await expect(wide).toBeDisabled();
          await expect(panel.getByTestId('circuit-profile-purpose-reason-wide')).toContainText('wide tables are unavailable');
        }
        await expect(long).toBeEnabled();
        await long.check();
        await expect(panel.getByTestId('circuit-profile-prepare')).toBeEnabled();
        await captureAuthoring(page, info, `layout-${engine}`, panel.getByTestId('circuit-profile-purpose'));
        await panel.getByTestId('circuit-profile-prepare').click();
        await expect.poll(() => actionCalls(mock, 'ensure').length).toBe(1);
        await expect(panel.getByTestId('group-schema-preview-result')).toBeVisible();
        expect(ensureGroup(mock).layout_version).toBe('managed_long_v1');
        noAuthoringMutation(mock, ['ensure', 'schema-preview']);
        await attachAuthoringEvidence(info, mock);
      });
    }

    for (const condition of ['stale', 'foreign', 'unknown', 'failed'] as const) {
      test(`AR-C06 ${condition} capability never enables preparation`, async ({ page }, info) => {
        const { mock } = await openAuthoring(page, 'new-mysql', { setup: value => {
          if (condition === 'failed') value.rejectNext('capability');
          else value.setCapability(condition === 'stale' ? { connector_revision: 'old-revision' }
            : condition === 'foreign' ? { connector_id: 'foreign' } : { metadata_status: 'unknown' });
        } });
        const panel = authoring(page);
        await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
        await expect(panel.getByTestId('circuit-profile-layout-blocked')).toContainText('unverified');
        await captureAuthoring(page, info, `capability-${condition}`, panel.getByTestId('circuit-profile-layout-blocked'));
        noAuthoringMutation(mock);
        await attachAuthoringEvidence(info, mock, condition === 'failed' ? ['capability'] : []);
      });
    }

    test('AR-C07 late capability R1 cannot authorize new saved connector revision R2', async ({ page }, info) => {
      const { mock, initialGate } = await openAuthoring(page, 'new-sqlserver', { holdInitial: 'capability' });
      const panel = authoring(page);
      try {
        await expect.poll(() => actionCalls(mock, 'capability').length).toBe(1);
        await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
        await page.getByTestId('fixture-refresh-connector-revision').click();
        await expect.poll(() => actionCalls(mock, 'capability').length).toBeGreaterThan(1);
        initialGate?.release();
        await expect(panel.getByTestId('circuit-profile-layout-blocked')).toContainText('unverified');
        await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
        noAuthoringMutation(mock);
        await captureAuthoring(page, info, 'capability-late-old-revision', panel.getByTestId('circuit-profile-layout-blocked'));
      } finally { initialGate?.release(); }
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-C08 unsupported saved wide remains readable rather than silently converted', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'saved-unsupported-wide');
      const before = mock.snapshot().groups;
      const panel = authoring(page);
      await expect(panel.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
      await expect(panel.getByTestId('circuit-profile-purpose-wide')).toBeDisabled();
      await expect(panel.getByTestId('circuit-profile-layout-blocked')).toContainText('wide tables are unavailable');
      await expect(panel.getByTestId('circuit-profile-prepare')).toBeDisabled();
      await page.reload();
      await expect(panel.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
      expect(mock.snapshot().groups).toEqual(before);
      noAuthoringMutation(mock, ['schema-preview']);
      await captureAuthoring(page, info, 'saved-unsupported-wide-readable', panel.getByTestId('circuit-profile-layout-blocked'));
      await attachAuthoringEvidence(info, mock);
    });
  });
}
