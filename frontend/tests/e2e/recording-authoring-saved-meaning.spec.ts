import { expect, test } from '@playwright/test';
import {
  AUTHORING_VIEWPORTS, actionCalls, attachAuthoringEvidence, authoring, captureAuthoring,
  expandAuthoringSources, noAuthoringMutation, object, openAuthoring, openAuthoringAdvanced, selectDevice,
} from './recording-authoring.helpers';

const lockedReason = 'Saved source roles cannot currently be changed. Layout or destination review preserves this role and its frozen binding.';
for (const viewport of AUTHORING_VIEWPORTS) {
  test.describe(`saved meaning ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    for (const role of ['gauge', 'state', 'counter_raw'] as const) {
      test(`AR-S01 saved ${role} stays locked through real Basic/Advanced, review Back/Cancel and reload`, async ({ page }, info) => {
        const { mock } = await openAuthoring(page, role === 'counter_raw' ? 'saved-counter' : `saved-${role}`);
        const before = structuredClone(mock.snapshot().groups);
        await expandAuthoringSources(page);
        const control = authoring(page).getByTestId('circuit-profile-role-tag-0');
        await expect(control).toHaveValue(role);
        await expect(control).toBeDisabled();
        await expect(control).toHaveAccessibleDescription(lockedReason);
        await captureAuthoring(page, info, `saved-${role}-locked`, control);
        const surface = await openAuthoringAdvanced(page);
        await surface.getByTestId('saved-recording-transition').click();
        const dialog = page.getByRole('dialog', { name: 'Change saved recording' });
        await expect(dialog).toBeVisible();
        await expect(dialog.getByText(/Sources and recording meaning stay unchanged/)).toBeVisible();
        await expect(dialog.getByTestId('circuit-profile-role-tag-0')).toHaveCount(0);
        await dialog.getByRole('button', { name: 'Review impact', exact: true }).click();
        await expect(dialog.getByRole('button', { name: 'Back', exact: true })).toBeEnabled();
        const review = object(actionCalls(mock, 'successor-review').at(-1)?.body);
        const members = object(review.group).members as Array<Record<string, unknown>>;
        expect(members.find(member => member.tag_id === 'tag-0')?.value_role).toBe(role);
        await captureAuthoring(page, info, `saved-${role}-supported-layout-review`, dialog);
        await dialog.getByRole('button', { name: 'Back', exact: true }).click();
        await dialog.getByRole('button', { name: 'Cancel', exact: true }).click();
        await expect(dialog).toHaveCount(0);
        await selectDevice(page, 'dev-2');
        await selectDevice(page, 'dev-1');
        await page.reload();
        await expandAuthoringSources(page);
        await expect(control).toHaveValue(role);
        await expect(control).toBeDisabled();
        await expect(control).toHaveAccessibleDescription(lockedReason);
        expect(mock.snapshot().groups).toEqual(before);
        expect(mock.snapshot().frozenGroups).toEqual(before);
        noAuthoringMutation(mock, ['schema-preview']);
        expect(actionCalls(mock, 'ensure')).toHaveLength(0);
        // The unchanged mock snapshot is a browser oracle, not real durable/frozen SQL proof.
        await attachAuthoringEvidence(info, mock);
      });
    }

    test('AR-S02 saved explicit cumulative membership survives toggles, device navigation and reload', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'saved-cumulative');
      const before = structuredClone(mock.snapshot().groups);
      const panel = authoring(page);
      await expandAuthoringSources(page);
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-1')).not.toBeChecked();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-2')).not.toBeChecked();
      const toggle = panel.getByTestId('circuit-profile-toggle-cumulative');
      await toggle.uncheck();
      await toggle.check();
      await selectDevice(page, 'dev-2');
      await selectDevice(page, 'dev-1');
      await page.reload();
      await expandAuthoringSources(page);
      await expect(panel.getByTestId('circuit-profile-cumulative-selected')).toHaveText('Selected cumulative sources: 1');
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
      await expect(panel.getByTestId('circuit-profile-member-cum-tag-1')).not.toBeChecked();
      expect(mock.snapshot().groups).toEqual(before);
      noAuthoringMutation(mock, ['schema-preview']);
      await captureAuthoring(page, info, 'saved-cumulative-explicit-membership', panel.getByTestId('circuit-profile-cumulative-selected'));
      await attachAuthoringEvidence(info, mock);
    });
  });
}
