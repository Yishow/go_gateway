import { expect, test } from '@playwright/test';

test('cancel asks before discarding an unsent point name and pending selection', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto('/tests/e2e/fixtures/recording-point-name-host.html');
  await page.getByTestId('recording-points-open').click();
  const dialog = page.getByRole('dialog');
  const pending = dialog.getByRole('checkbox', { name: 'MOD_1 (40002)' });
  await pending.check();
  const name = dialog.getByRole('textbox', { name: '點位名稱 MOD_0 (40001)' });
  await name.click();
  await expect(name).toBeEditable();
  await expect(name).toHaveAttribute('aria-busy', 'false');
  await name.fill('尚未儲存的名稱');

  await dialog.getByRole('button', { name: '取消', exact: true }).click();
  const confirmation = page.getByRole('alertdialog');
  await expect(confirmation).toBeVisible({ timeout: 3_000 });
  await expect(confirmation).toHaveAccessibleName('捨棄這次的點位名稱修改？');
  await expect(confirmation).toContainText('已保存的名稱和記錄草稿會保留');
  await expect(page.getByTestId('fixture-saves')).toHaveText('[]');
  await confirmation.getByRole('button', { name: '繼續編輯' }).click();
  await expect(name).toHaveValue('尚未儲存的名稱');
  await expect(pending).toBeChecked();

  await dialog.getByRole('button', { name: '取消', exact: true }).click();
  await page.getByRole('alertdialog').getByRole('button', { name: '捨棄並關閉選取' }).click();
  await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.getByTestId('fixture-saves')).toHaveText('[]');
  await page.getByTestId('recording-points-open').click();
  await expect(page.getByRole('textbox', { name: '點位名稱 MOD_0 (40001)' })).toHaveValue('MOD_0 (40001)');
  await expect(page.getByRole('checkbox', { name: 'MOD_1 (40002)' })).not.toBeChecked();
});

test('discarding a picker name does not disable later name autosave in the recording', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  await page.getByTestId('recording-points-open').click();
  const picker = page.getByRole('dialog');
  const pickerName = picker.getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
  await pickerName.click();
  await expect(pickerName).toHaveAttribute('aria-busy', 'false');
  await pickerName.fill('捨棄的名稱');
  await picker.getByRole('button', { name: '取消', exact: true }).click();
  await page.getByRole('alertdialog').getByRole('button', { name: '捨棄並關閉選取' }).click();

  const row = page.getByTestId('group-member-row-P0');
  const mainName = row.getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
  await mainName.click();
  await expect(mainName).toHaveAttribute('aria-busy', 'false');
  await mainName.fill('主畫面的新名稱');
  await page.getByTestId('group-name').click();
  await expect(row.getByRole('textbox', { name: '點位名稱 主畫面的新名稱' })).toBeVisible({ timeout: 3_000 });
});

test.describe('touch cancellation', () => {
  test.use({ hasTouch: true, viewport: { width: 390, height: 844 } });

  test('discards only the unsent name and retains a previously saved name', async ({ page }, testInfo) => {
    await page.goto('/tests/e2e/fixtures/recording-point-name-host.html');
    await page.getByTestId('recording-points-open').tap();
    const picker = page.getByRole('dialog');
    const firstName = picker.getByRole('textbox', { name: '點位名稱 MOD_0 (40001)' });
    await firstName.tap();
    await expect(firstName).toHaveAttribute('aria-busy', 'false');
    await firstName.fill('已保存的名稱');
    await firstName.press('Enter');
    await expect(page.getByTestId('fixture-saves')).toHaveText('[{"point":"P32","name":"已保存的名稱"}]');

    await picker.getByRole('button', { name: '取消', exact: true }).tap({ timeout: 5_000 });
    await page.getByTestId('recording-points-open').tap();
    const reopened = page.getByRole('dialog');
    const savedName = reopened.getByRole('textbox', { name: '點位名稱 已保存的名稱' });
    await savedName.tap({ timeout: 5_000 });
    await expect(savedName).toBeEditable({ timeout: 5_000 });
    await expect(savedName).toHaveAttribute('aria-busy', 'false');
    await savedName.fill('尚未保存的名稱');
    await reopened.getByRole('button', { name: '取消', exact: true }).tap({ timeout: 5_000 });
    const confirmation = page.getByRole('alertdialog');
    await expect(confirmation).toHaveAccessibleName('捨棄這次的點位名稱修改？');
    await page.screenshot({ path: testInfo.outputPath('picker-cancel-confirmation-mobile.png') });
    await confirmation.getByRole('button', { name: '捨棄並關閉選取' }).tap({ timeout: 5_000 });
    await expect(page.getByTestId('fixture-saves')).toHaveText('[{"point":"P32","name":"已保存的名稱"}]');
    await page.getByTestId('recording-points-open').tap();
    await expect(page.getByRole('textbox', { name: '點位名稱 已保存的名稱' })).toHaveValue('已保存的名稱');
  });
});
