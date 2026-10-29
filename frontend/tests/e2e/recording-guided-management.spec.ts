import { expect, test } from '@playwright/test';
import {
  GUIDED_VIEWPORTS, actionCalls, advanced, attachGuidedEvidence, captureGuided,
  noGuidedMutation, openGuided, openGuidedAdvanced, tabTo,
} from './recording-guided.helpers';

// Existing controls and real dialogs are exercised; mock responses establish no SQL/history durability.
for (const viewport of GUIDED_VIEWPORTS) {
  test.describe(`guided saved recording management ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    for (const scenario of ['saved-wide', 'saved-long', 'saved-custom'] as const) {
      test(`GD-M01 ${scenario} exposes the same saved settings without conversion or another start`, async ({ page }, info) => {
        const { mock } = await openGuided(page, scenario);
        const before = structuredClone(mock.snapshot().groups), original = before[0];
        const surface = await openGuidedAdvanced(page);
        await expect(surface.locator(':scope > summary')).toHaveText('Advanced settings and workspace status');
        await expect(surface.getByTestId('recording-workspace-status-note')).toContainText('does not start or activate them');
        const card = surface.getByTestId('group-list-item-G1');
        await expect(card).toContainText(original.name);
        await expect(card.getByTestId('saved-recording-edit')).toBeVisible();
        await expect(card.getByTestId('group-disable')).toHaveText('Stop new recording');
        await expect(card.getByTestId('group-delete')).toHaveText('Remove settings');
        const open = card.getByTestId('group-open-G1'); await tabTo(page, open); await page.keyboard.press('Enter');
        const editor = surface.getByTestId('group-editor'); await expect(editor).toBeVisible();
        await expect(editor.getByTestId('group-name')).toHaveValue(original.name);
        await expect(editor.getByTestId(`group-storage-${original.destination.storage_strategy}`)).toBeChecked();
        await expect(editor.getByTestId('group-applied-revision')).toHaveText(original.applied_revision!);
        if (scenario !== 'saved-custom') {
          await expect(editor.getByTestId('group-profile-layout')).toContainText(scenario === 'saved-wide' ? 'wide table' : 'long table');
          await expect(editor.getByTestId('group-profile-interval')).toContainText('15');
        }
        await captureGuided(page, info, `${scenario}-original-settings-inspection`, editor);
        await editor.getByTestId('group-editor-close').click(); await expect(card).toBeVisible();
        expect(mock.snapshot().groups).toEqual(before); noGuidedMutation(mock);
        await attachGuidedEvidence(info, mock);
      });
    }

    test('GD-M02 explicit Stop new recording preserves the original destination and historical backlog', async ({ page }, info) => {
      const { mock, holdNext } = await openGuided(page, 'saved-long');
      const surface = await openGuidedAdvanced(page);
      const card = surface.getByTestId('group-list-item-G1');
      await expect(card).toContainText('Already accepted backlog still goes to its original destination');
      const stop = card.getByTestId('group-disable'); await expect(stop).toBeEnabled(); await tabTo(page, stop);
      const before = structuredClone(mock.snapshot());
      noGuidedMutation(mock);
      const gate = holdNext('stop');
      try {
        await stop.dblclick({ delay: 20 });
        await expect.poll(() => actionCalls(mock, 'stop').length).toBe(1);
        expect(actionCalls(mock, 'stop')[0].body).toEqual({ workspace_id: 'ws-1', expected_workspace_revision: before.workspaceRevision,
          expected_group_revision: before.groups[0].revision, expected_connector_revision: 'crev-1' });
        gate.release();
        await expect(stop).toBeDisabled();
        expect(actionCalls(mock, 'stop')[0].status).toBe(200);
        expect(mock.snapshot().groups[0].status).toBe('disabled');
        expect(mock.snapshot().groups[0].destination).toEqual(before.groups[0].destination);
        expect(mock.snapshot().history).toEqual(before.history);
        noGuidedMutation(mock, ['stop']);
        await captureGuided(page, info, 'stopped-intake-history-preserved', card);
      } finally { gate.release(); }
      await attachGuidedEvidence(info, mock);
    });

    test('GD-M03 Remove settings requires separate consent, cancels safely, and retains history after confirmation', async ({ page }, info) => {
      const { mock, holdNext } = await openGuided(page, 'saved-long');
      const surface = await openGuidedAdvanced(page);
      const card = surface.getByTestId('group-list-item-G1'), remove = card.getByTestId('group-delete');
      await expect(remove).toBeEnabled(); const before = structuredClone(mock.snapshot()); await remove.click();
      const dialog = page.getByRole('alertdialog', { name: 'Remove these saved recording settings?' });
      await expect(dialog).toBeVisible(); await expect(dialog).toContainText('SQL tables, recorded history and delivery receipts are kept.');
      await expect(dialog.getByTestId('lifecycle-pending-samples')).toHaveText('3');
      await expect(dialog.getByTestId('lifecycle-pending-outbox')).toHaveText('4');
      await expect(dialog.getByTestId('lifecycle-pending-blocked')).toHaveText('2');
      await expect(dialog.getByTestId('lifecycle-pending-unknown')).toHaveText('1');
      noGuidedMutation(mock);
      await captureGuided(page, info, 'remove-review-distinct-backlog-units', dialog);
      await dialog.getByTestId('group-delete-cancel').click(); await expect(dialog).toHaveCount(0); await expect(remove).toBeFocused();
      expect(mock.snapshot()).toEqual(before); noGuidedMutation(mock);
      await remove.click(); await expect(dialog).toBeVisible();
      const confirm = dialog.getByTestId('group-delete-confirm-button'); await tabTo(page, confirm);
      const gate = holdNext('remove');
      try {
        await confirm.dblclick({ delay: 20 });
        await expect.poll(() => actionCalls(mock, 'remove').length).toBe(1);
        expect(actionCalls(mock, 'remove')[0].body).toEqual({ workspace_id: 'ws-1', expected_workspace_revision: before.workspaceRevision,
          expected_group_revision: before.groups[0].revision, expected_connector_revision: 'crev-1' });
        gate.release(); await expect(dialog).toHaveCount(0);
        expect(actionCalls(mock, 'remove')[0].status).toBe(200);
        expect(mock.snapshot().groups[0].status).toBe('deleted');
        expect(mock.snapshot().groups[0].destination).toEqual(before.groups[0].destination);
        expect(mock.snapshot().history).toEqual(before.history);
        const removed = surface.getByTestId('group-removed'); await removed.locator(':scope > summary').click();
        await expect(removed.getByTestId('group-list-item-G1')).toBeVisible();
        await removed.getByTestId('group-open-G1').click();
        await expect(surface.getByTestId('group-editor')).toBeVisible(); await expect(surface.getByTestId('group-save')).toBeDisabled();
        noGuidedMutation(mock, ['remove']);
        await captureGuided(page, info, 'removed-settings-readonly-history-not-deleted', surface.getByTestId('group-editor'));
      } finally { gate.release(); }
      await attachGuidedEvidence(info, mock);
    });

    test('GD-M04 server-blocked saved wide actions keep their reason and allow read-only inspection', async ({ page }, info) => {
      const { mock } = await openGuided(page, 'blocked-wide');
      const before = structuredClone(mock.snapshot().groups), surface = await openGuidedAdvanced(page);
      const card = surface.getByTestId('group-list-item-G1');
      await expect(card.getByTestId('group-disable')).toBeEnabled(); // Server permits safe stop despite stale destination identity.
      for (const [action, testId] of [['edit', 'saved-recording-edit'], ['remove', 'group-delete']] as const) {
        const button = card.getByTestId(testId), reason = card.getByTestId(`saved-recording-reason-${action}`);
        await expect(button).toBeDisabled(); await expect(reason).toBeVisible();
        await expect(button).toHaveAttribute('aria-describedby', await reason.getAttribute('id') ?? 'missing-reason');
        await expect(reason).toContainText('The saved destination connection has changed.');
        await expect(reason).toContainText(action === 'remove' ? 'You can stop new recording while investigating.' : 'Check the saved connection in destination settings.');
      }
      await captureGuided(page, info, 'blocked-wide-server-reasons-visible', card);
      await card.getByTestId('group-open-G1').click(); await expect(advanced(page).getByTestId('group-editor')).toBeVisible();
      await expect(surface.getByTestId('group-profile-layout')).toContainText('wide table');
      expect(mock.snapshot().groups).toEqual(before); noGuidedMutation(mock);
      await attachGuidedEvidence(info, mock);
    });
  });
}
