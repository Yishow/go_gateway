import { expect, test } from '@playwright/test';

for (const advanced of [false, true]) {
  for (const width of [1024, 1600]) {
    test(`mapping save info ${advanced ? 'advanced' : 'basic'} at ${width}px`, async ({ page }, testInfo) => {
      await page.setViewportSize({ width, height: 720 });
      await page.route('**/api/**', route => route.fulfill({ json: { data: [] } }));
      await page.goto('/tests/e2e/fixtures/mapping-live-overflow.html');
      if (advanced) await page.getByTestId('mapping-advanced-toggle').click();
      const input = page.getByTestId('input-display-name-point-0');
      const info = page.getByTestId('mapping-save-info-point-0');
      const row = page.getByTestId('mapping-row-point-0');
      await expect(info).toBeVisible();
      await expect(page.getByRole('dialog')).toHaveCount(0);
      const inputBox = (await input.boundingBox())!, infoBox = (await info.boundingBox())!;
      expect(inputBox.width).toBeGreaterThanOrEqual(160);
      expect(infoBox.x).toBeGreaterThanOrEqual(inputBox.x + inputBox.width);
      expect(Math.abs(inputBox.y + inputBox.height / 2 - infoBox.y - infoBox.height / 2)).toBeLessThan(2);
      const rowHeight = (await row.boundingBox())!.height;
      await info.hover();
      await expect(page.getByRole('dialog')).toContainText('已儲存');
      expect((await row.boundingBox())!.height).toBe(rowHeight);
      await page.keyboard.press('Escape');
      await expect(page.getByRole('dialog')).toHaveCount(0);

      // Native keyboard activation, then an actionable safe error in the portal.
      const error = page.getByTestId('mapping-save-info-point-1');
      await error.focus();
      await page.keyboard.press('Space');
      await expect(page.getByRole('dialog')).toContainText('test-request');
      await expect(page.getByTestId('mapping-save-retry-point-1')).toBeVisible();
      await page.mouse.move(1, 1);
      await expect(page.getByRole('dialog')).toBeVisible();
      await page.screenshot({ path: testInfo.outputPath('save-error-info.png') });
      await page.keyboard.press('Escape');

      const scroller = page.getByTestId('mapping-table-wrapper').locator('table').locator('..');
      await scroller.evaluate(element => { element.scrollTop = element.scrollHeight; });
      await page.getByTestId('mapping-save-info-point-47').hover();
      const dialog = page.getByRole('dialog');
      await expect(dialog).toBeVisible();
      expect(await dialog.evaluate(element => element.parentElement === document.body)).toBe(true);
      const box = (await dialog.boundingBox())!;
      expect(box.y).toBeGreaterThanOrEqual(0);
      expect(box.y + box.height).toBeLessThanOrEqual(720);
      expect(box.x).toBeGreaterThanOrEqual(0);
      expect(box.x + box.width).toBeLessThanOrEqual(width);
      await scroller.evaluate(element => { element.scrollTop -= 20; });
      await expect(dialog).toHaveCount(0);
    });
  }
}
