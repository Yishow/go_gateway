import { expect, test } from '@playwright/test';

for (const viewport of [{ width: 1440, height: 900 }, { width: 390, height: 844 }]) {
  test(`point names save from the picker at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize(viewport);
    await page.goto('/tests/e2e/fixtures/recording-point-name-host.html');
    await page.getByTestId('recording-points-open').click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await page.screenshot({ path: testInfo.outputPath(`point-picker-${viewport.width}.png`) });
    const first = dialog.getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
    const second = dialog.getByRole('textbox', { name: '點位名稱 TAG_1 (D1)' });
    await first.click();
    await first.fill('主馬達');
    await page.screenshot({ path: testInfo.outputPath(`point-picker-editing-${viewport.width}.png`) });
    await second.click();
    await expect(page.getByTestId('fixture-saves')).toContainText('主馬達');
    await second.fill('進水泵');
    await second.press('Enter');
    await expect(page.getByTestId('fixture-saves')).toContainText('進水泵');
    await expect(dialog.getByRole('button', { name: '保存名稱' })).toHaveCount(0);
    const renamed = dialog.getByRole('textbox', { name: '點位名稱 主馬達' });
    await renamed.click();
    await renamed.fill('不保存');
    await renamed.press('Escape');
    await expect(renamed).toHaveValue('主馬達');
    await expect(page.getByTestId('fixture-saves')).not.toContainText('不保存');
    await dialog.getByLabel('來源設備').selectOption('modbus');
    const modbus = dialog.getByRole('textbox', { name: '點位名稱 MOD_0 (40001)' });
    await modbus.click();
    await modbus.fill('Modbus 馬達');
    await modbus.press('Enter');
    await expect(page.getByTestId('fixture-saves')).toContainText('Modbus 馬達');
  });
}
