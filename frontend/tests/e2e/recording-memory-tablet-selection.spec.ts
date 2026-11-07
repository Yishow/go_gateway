import { expect, test } from '@playwright/test';

// MOCK SAMPLE UI-only regression: no external device or database calls.
for (const width of [390, 768, 1024, 1440]) {
  for (const light of [true, false]) {
    test(`selected memory row stays compact ${width}px ${light ? 'light' : 'dark'}`, async ({ page }) => {
      await page.setViewportSize({ width, height: width === 390 ? 844 : 1024 });
      await page.route('**/api/**', route => route.abort());
      await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=long${light ? '&light' : ''}`);
      const first = page.getByTestId('home-memory-cell-40001');
      const inspection = page.getByTestId('home-memory-inspection');
      const controls = page.locator('.gw-memory-toolbar');
      const initial = await first.boundingBox();
      const footer = await page.locator('.gw-modal-footer').boundingBox();
      const toolbar = await controls.boundingBox();
      await first.focus();
      await expect(inspection).toContainText('40001');
      const selected = await first.boundingBox();
      expect(selected!.height).toBeCloseTo(initial!.height, 0);
      expect(selected!.height).toBeLessThan(110);
      expect(await first.locator('.gw-memory-focus-information').count()).toBe(0);
      expect((await page.locator('.gw-modal-footer').boundingBox())!.y).toBeCloseTo(footer!.y, 0);
      expect((await controls.boundingBox())!.y).toBeCloseTo(toolbar!.y, 0);
      await page.keyboard.press('ArrowRight');
      await expect(page.getByTestId('home-memory-cell-40002')).toBeFocused();
      await expect(inspection).toContainText('40002');
      await inspection.focus();
      await page.keyboard.press('End');
      expect(await inspection.evaluate(node => node.scrollWidth <= node.clientWidth + 1)).toBe(true);
      for (const layout of ['navigator', 'map', 'chapters', 'detail', 'grid']) {
        await page.getByTestId(`home-memory-layout-${layout}`).click();
        await expect(page.getByTestId('home-memory-cell-40002')).toHaveAttribute('tabindex', '0');
        expect((await page.getByTestId('home-memory-cell-40002').boundingBox())!.height).toBeLessThan(110);
      }
      await expect(page.getByTestId('home-memory-close')).toBeInViewport();
      await page.screenshot({ path: `test-results/tablet-selection-${width}-${light ? 'light' : 'dark'}.png` });
    });
  }
}
