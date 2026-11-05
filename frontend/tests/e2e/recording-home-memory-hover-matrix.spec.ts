import { expect, test } from '@playwright/test';

// Deterministic MOCK SAMPLE values only; this is not evidence of a gateway/HSL live read.
for (const width of [1440, 768, 390]) {
  for (const light of [true, false]) {
    for (const layout of ['map', 'detail']) {
      test(`${layout} hover previews without moving keyboard focus at ${width}px ${light ? 'light' : 'dark'}`, async ({ page }) => {
        await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
        await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en${light ? '&light' : ''}`);
        await expect(page.getByTestId('fixture-sample-notice')).toContainText('MOCK SAMPLE');
        await page.getByTestId(`home-memory-layout-${layout}`).click();
        const first = page.getByTestId('home-memory-cell-40001'), continuation = page.getByTestId('home-memory-cell-40018');
        const detail = page.getByTestId('home-memory-detail'), action = page.getByTestId('fixture-last-action');
        await first.focus();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40001');
        await continuation.hover();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40018');
        await expect(detail.locator('.gw-memory-detail-reading')).toHaveText('12.5');
        await expect(detail).toContainText('40017');
        await expect(first).toBeFocused();
        await expect(first).toHaveAttribute('tabindex', '0');
        await expect(continuation).toHaveAttribute('tabindex', '-1');
        await expect(action).toHaveText('');

        await page.keyboard.press('ArrowRight');
        await expect(page.getByTestId('home-memory-cell-40002')).toBeFocused();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40002');
        await expect(detail.locator('.gw-memory-detail-reading')).toHaveText('-42');

        // Move the pointer off the cell before entering again; the keyboard owns selection meanwhile.
        await page.getByTestId(`home-memory-layout-${layout}`).hover();
        await continuation.hover();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40018');
        await detail.getByTestId('home-memory-detail-previous').click();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40017');
        await expect(page.getByTestId('home-memory-cell-40017')).toHaveAttribute('tabindex', '0');
        await detail.getByTestId('home-memory-detail-next').click();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40018');
        await expect(page.getByTestId('home-memory-cell-40018')).toHaveAttribute('tabindex', '0');
        await expect(detail.locator('.gw-memory-detail-reading')).toHaveText('12.5');
        await expect(action).toHaveText('');

        await first.focus();
        await expect(detail.locator('.gw-memory-detail-address')).toHaveText('40001');
        await expect(detail.locator('.gw-memory-detail-reading')).toHaveText('0');
        await expect(action).toHaveText('');
      });
    }
  }
}
