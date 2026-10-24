import { expect, test } from '@playwright/test';

for (const advanced of [false, true]) {
  for (const viewport of [{ width: 1600, height: 1000 }, { width: 1024, height: 720 }]) {
    test(`${advanced ? 'advanced' : 'basic'} live values stay within the page at ${viewport.width}px`, async ({ page }, testInfo) => {
      await page.setViewportSize(viewport);
      await page.route('**/api/**', route => route.fulfill({ json: { data: [] } }));
      await page.goto('/tests/e2e/fixtures/mapping-live-overflow.html');
      const rows = page.getByTestId('mapping-table-wrapper').locator('tbody tr');
      const statuses = page.locator('[data-testid^="preview-state-"]');
      const scroller = page.getByTestId('mapping-table-wrapper').locator('table').locator('..');
      const pageOverflow = () => page.evaluate(() => {
        const bottom = document.querySelector('[data-testid=overflow-page]')!.getBoundingClientRect().bottom + scrollY;
        return document.documentElement.scrollHeight - Math.ceil(bottom);
      });
      await expect(rows).toHaveCount(48);
      await expect(statuses).toHaveCount(0);
      await expect.poll(pageOverflow).toBeLessThanOrEqual(1);
      if (advanced) await page.getByTestId('mapping-advanced-toggle').click();
      expect(await page.getByTestId('input-display-name-point-0').evaluate(element => element.getBoundingClientRect().width))
        .toBeGreaterThanOrEqual(160);

      // Non-empty values mount the real absolute sr-only live status in every row.
      await page.getByTestId('load-live-values').click();
      await expect(statuses).toHaveCount(48);
      await expect.poll(pageOverflow).toBeLessThanOrEqual(1);
      expect(await scroller.evaluate(element => element.clientHeight)).toBeGreaterThan(500);
      await scroller.scrollIntoViewIfNeeded();
      await scroller.evaluate(element => { element.scrollTop = element.scrollHeight; });
      await expect(rows.last()).toBeInViewport();
      await expect.poll(pageOverflow).toBeLessThanOrEqual(1);

      await page.evaluate(() => window.scrollTo(0, document.documentElement.scrollHeight));
      expect(await page.evaluate(() => {
        const bottom = document.querySelector('[data-testid=overflow-page]')!.getBoundingClientRect().bottom;
        return Math.max(0, innerHeight - bottom);
      })).toBeLessThanOrEqual(1);
      await page.screenshot({ path: testInfo.outputPath('live-values-no-blank.png') });
    });
  }
}
