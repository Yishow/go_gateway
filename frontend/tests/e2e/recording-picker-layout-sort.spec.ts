import { expect, test } from '@playwright/test';

test('loading a point name keeps the table row height stable', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto('/tests/e2e/fixtures/recording-point-name-host.html');
  await page.getByTestId('recording-points-open').click();
  const dialog = page.getByRole('dialog');
  const input = dialog.getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
  const row = input.locator('xpath=ancestor::tr');
  const before = await row.boundingBox();
  await input.click();
  await expect(row.getByText('正在讀取點位來源…')).toBeVisible();
  const during = await row.boundingBox();
  expect(before).not.toBeNull(); expect(during).not.toBeNull();
  expect(during!.height).toBeCloseTo(before!.height, 0);
});

for (const viewport of [{ width: 1440, height: 900 }, { width: 390, height: 844 }]) {
  test(`sorting preserves pending points at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize(viewport);
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-point-name-host.html');
    await page.getByTestId('recording-points-open').click();
    const dialog = page.getByRole('dialog');
    const address = dialog.getByRole('columnheader', { name: '位址' });
    await expect(address).toHaveAttribute('aria-sort', 'ascending');
    const firstRow = dialog.locator('tbody tr').first();
    await expect(firstRow.locator('td').nth(2)).toHaveText('40001');
    const existing = dialog.getByRole('checkbox', { name: 'TAG_0 (D0)' });
    await expect(existing).toBeDisabled();
    await expect(existing).toHaveAccessibleDescription('已在此記錄');
    const pending = dialog.getByRole('checkbox', { name: 'MOD_0 (40001)' });
    await pending.check();
    await dialog.getByRole('button', { name: '位址' }).click();
    await expect(address).toHaveAttribute('aria-sort', 'descending');
    await expect(firstRow.locator('td').nth(2)).toHaveText('D31');
    await expect(dialog.getByText('待加入 1 點')).toBeVisible();
    await expect(dialog.getByText('另有 1 個已選點位不在目前頁面或篩選中')).toBeVisible();
    await dialog.getByRole('button', { name: '位址' }).click();
    await expect(pending).toBeChecked();
    const horizontalHint = dialog.getByText('左右滑動查看位址與資料型別');
    if (viewport.width === 390) await expect(horizontalHint).toBeVisible();
    else await expect(horizontalHint).not.toBeVisible();
    await expect(dialog.getByText('已在此記錄', { exact: true })).not.toBeVisible();
    await dialog.locator('.recording-point-picker__body').evaluate(element => { element.scrollLeft = 0; element.scrollTop = 0; });
    await page.screenshot({ path: testInfo.outputPath(`picker-sort-${viewport.width}.png`) });
  });
}
