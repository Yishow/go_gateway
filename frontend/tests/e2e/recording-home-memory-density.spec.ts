import { expect, test } from '@playwright/test';

// MOCK SAMPLE only. Common controls must leave useful data visible on desktop opening.
for (const width of [1440, 768]) {
for (const light of [true, false]) {
  for (const layout of ['grid', 'navigator', 'map', 'chapters', 'detail']) {
    test(`${layout} initially shows two full address rows at ${width}px ${light ? 'light' : 'dark'}`, async ({ page }) => {
      await page.setViewportSize({ width, height: 1000 });
      await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en${light ? '&light' : ''}`);
      await page.getByTestId(`home-memory-layout-${layout}`).click();
      const view = page.getByTestId('home-memory-view');
      await view.evaluate(node => { node.scrollTop = 0; });
      // Density cannot be recovered by hiding the fixture's data-provenance warning.
      const notice = page.getByTestId('fixture-sample-notice');
      await expect(notice).toBeVisible();
      await expect(notice).toHaveText('MOCK SAMPLE readings · not connected to HSL or a device');
      const noticeBounds = await notice.boundingBox();
      expect(noticeBounds?.width).toBeGreaterThan(100);
      expect(noticeBounds?.height).toBeGreaterThan(15);
      await expect.poll(() => view.evaluate(node => {
        const viewport = node.getBoundingClientRect();
        const grid = node.querySelector('.gw-memory-grid')!;
        const columns = getComputedStyle(grid).gridTemplateColumns.split(' ').filter(Boolean).length;
        return [...grid.children].slice(0, columns * 2).every(cell => {
          const bounds = cell.getBoundingClientRect();
          return bounds.top >= viewport.top && bounds.bottom <= viewport.bottom;
        });
      })).toBe(true);
    });
  }
}
}
