import { expect, test } from '@playwright/test';
import { VIEWPORTS, advancedFor, assertNoNewEffects, assertPendingUnits, attachLedger, basicFor, captureGeometry,
  dialogFor, openChange, openFixture, showAdvanced, showBasicActions } from './recording-lifecycle.helpers';

// Formal browser source. Discovery/typecheck do not execute these assertions or produce screenshots.
for (const viewport of VIEWPORTS) test.describe(`recording lifecycle actions ${viewport.width}x${viewport.height}`, () => {
  test.use({ viewport, timezoneId: 'UTC' });

  for (const [scenario, enabled] of [
    ['draft-wide', [true, true, true, false, false]],
    ['applied-wide', [true, true, true, false, true]],
    ['disabled-wide', [true, false, true, false, true]],
    ['deleted-wide', [false, false, false, true, false]],
  ] as const) test(`LC-A01 ${scenario}: Basic/Advanced server-action parity and saved revisions`, async ({ page }, info) => {
    const mock = await openFixture(page, scenario);
    const removed = scenario === 'deleted-wide';
    const basic = await showBasicActions(page, removed);
    const ids = ['saved-recording-edit', 'group-disable', 'group-delete', 'saved-recording-recreate', 'saved-recording-transition'];
    for (const [index, id] of ids.entries()) {
      if (enabled[index]) await expect(basic.getByTestId(id)).toBeEnabled();
      else await expect(basic.getByTestId(id)).toBeDisabled();
    }
    await captureGeometry(page, info, `${scenario}-basic`);
    const advanced = await showAdvanced(page, removed);
    for (const [index, id] of ids.entries()) {
      if (enabled[index]) await expect(advanced.getByTestId(id)).toBeEnabled();
      else await expect(advanced.getByTestId(id)).toBeDisabled();
    }
    await expect(advanced.getByTestId('group-revisions')).toContainText('R1');
    if (scenario !== 'draft-wide') await expect(advanced.getByTestId('group-applied-revision')).toHaveText('APPLIED1');
    else await expect(advanced.getByTestId('group-applied-revision')).not.toHaveText('APPLIED1');
    assertNoNewEffects(mock);
    await captureGeometry(page, info, `${scenario}-advanced`);
    await attachLedger(info, mock);
  });

  test('LC-A02 after host setup, initial Cancel, Escape and Back have zero new effects; keyboard focus restores', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const before = mock.snapshot();
    const { dialog, trigger } = await openChange(page);
    const steps = dialog.getByRole('list', { name: '記錄變更步驟' });
    await expect(steps.getByRole('listitem')).toHaveCount(3);
    const first = dialog.getByLabel('記錄名稱', { exact: true });
    await first.focus();
    await page.keyboard.press('Shift+Tab');
    await expect(dialog.getByTestId('recording-change-review')).toBeFocused();
    await page.keyboard.press('Tab');
    await expect(first).toBeFocused();
    await captureGeometry(page, info, 'keyboard-choose', dialog);
    await page.keyboard.press('Escape');
    await expect(dialog).toHaveCount(0);
    await expect(trigger).toBeFocused();
    await trigger.press('Enter');
    await dialog.getByRole('button', { name: '檢視變更影響', exact: true }).click();
    await assertPendingUnits(dialog);
    await captureGeometry(page, info, 'initial-impact', dialog);
    await dialog.getByRole('button', { name: '返回', exact: true }).click();
    await expect(dialog.getByLabel('資料表布局', { exact: true })).toBeVisible();
    await dialog.getByRole('button', { name: '取消', exact: true }).click();
    await expect(dialog).toHaveCount(0);
    await expect(trigger).toBeFocused();
    expect(mock.snapshot().source).toEqual(before.source);
    expect(mock.snapshot().workspaceRevision).toBe(before.workspaceRevision);
    expect(mock.snapshot().lineageRevision).toBe(before.lineageRevision);
    expect(mock.snapshot().preparation).toBeUndefined();
    expect(mock.snapshot().target).toBeUndefined();
    expect(mock.snapshot().status).toBeUndefined();
    assertNoNewEffects(mock);
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(0);
    await attachLedger(info, mock);
  });

  test('LC-A03 stop, remove confirmation and tombstone preserve distinct historical units', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const original = mock.snapshot().source;
    let actions = await showBasicActions(page);
    await actions.getByRole('button', { name: '停止新記錄', exact: true }).click();
    await expect.poll(() => mock.ledger.filter(entry => entry.action === 'disable').length).toBe(1);
    expect(mock.snapshot().source.status).toBe('disabled');
    expect(mock.snapshot().source.destination).toEqual(original.destination);
    expect(mock.snapshot().source.applied_revision).toBe(original.applied_revision);
    actions = await showBasicActions(page);
    await expect(actions.getByTestId('group-disable')).toBeDisabled();
    await actions.getByRole('button', { name: '移除設定', exact: true }).click();
    const remove = page.getByRole('alertdialog', { name: '要移除這份已保存的記錄設定嗎？' });
    await assertPendingUnits(remove);
    await expect(remove).toContainText('SQL 資料表、歷史資料及交付回執都會保留');
    expect(mock.ledger.filter(entry => entry.action === 'remove')).toHaveLength(0);
    await remove.getByRole('button', { name: '取消', exact: true }).click();
    expect(mock.snapshot().source.status).toBe('disabled');
    await actions.getByTestId('group-delete').click();
    await remove.getByRole('button', { name: '移除設定並保留歷史', exact: true }).click();
    await expect.poll(() => mock.snapshot().source.status).toBe('deleted');
    const tombstone = await showBasicActions(page, true);
    await expect(tombstone.getByTestId('saved-recording-recreate')).toBeEnabled();
    const advanced = await showAdvanced(page, true);
    await expect(advanced.getByTestId('group-applied-revision')).toHaveText('APPLIED1');
    expect(mock.snapshot().source.destination).toEqual(original.destination);
    expect(mock.ledger.filter(entry => entry.action === 'remove')).toHaveLength(1);
    expect(mock.snapshot().target).toBeUndefined();
    await captureGeometry(page, info, 'removed-readable-tombstone');
    await attachLedger(info, mock);
  });

  test('LC-A04 blocked removal explains stale connector and preserves saved-revision stop route', async ({ page }, info) => {
    const mock = await openFixture(page, 'blocked-wide');
    const actions = await showBasicActions(page);
    await expect(actions.getByTestId('group-delete')).toBeDisabled();
    await expect(actions.getByTestId('saved-recording-reason-remove')).toContainText('原先保存的目的地連線已變更');
    await expect(actions.getByTestId('saved-recording-reason-remove')).toContainText('仍可先停止新記錄');
    await expect(actions.getByTestId('group-disable')).toBeEnabled();
    await actions.getByTestId('group-disable').click();
    await expect.poll(() => mock.snapshot().source.status).toBe('disabled');
    expect(mock.ledger.find(entry => entry.action === 'disable')?.body).toMatchObject({ expected_connector_revision: 'crev-1', expected_group_revision: 'R1' });
    expect(mock.ledger.filter(entry => entry.action === 'remove')).toHaveLength(0);
    await captureGeometry(page, info, 'blocked-repair-route');
    await attachLedger(info, mock);
  });

  test('LC-A05 managed scoped start remains a review route, not direct apply or workspace start', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const advanced = await showAdvanced(page);
    const start = advanced.getByRole('button', { name: '到基本模式審核並開始', exact: true });
    await expect(start).toBeEnabled();
    await start.click();
    await expect(basicFor(page)).toBeFocused();
    await expect(basicFor(page).locator('[data-testid="basic-recording-device-tile"][data-device-id="dev-1"]')).toHaveAttribute('aria-pressed', 'true');
    await expect(basicFor(page).getByTestId('circuit-status-instantaneous')).toBeVisible();
    assertNoNewEffects(mock);
    await expect(dialogFor(page)).toHaveCount(0);
    await expect(advancedFor(page)).toHaveAttribute('open', '');
    await captureGeometry(page, info, 'managed-scoped-start-review');
    await attachLedger(info, mock);
  });
});
