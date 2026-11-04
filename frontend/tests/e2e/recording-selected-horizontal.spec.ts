import { expect, test } from '@playwright/test';

for (const width of [320, 390, 640]) test(`phone users can scroll selected points and use their controls at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  await table.scrollIntoViewIfNeeded();
  await page.screenshot({ path: testInfo.outputPath(`selected-initial-${width}.png`) });
  const size = await region.evaluate(element => ({ client: element.clientWidth, scroll: element.scrollWidth, overflow: getComputedStyle(element).overflowX }));
  expect(size.overflow).toBe('auto');
  expect(size.scroll).toBeGreaterThan(size.client);
  await expect(table.getByText('左右滑動查看其餘欄位')).toBeVisible();
  const headerTops = await region.locator('thead th').evaluateAll(headers => headers.map(header => header.getBoundingClientRect().top));
  expect(Math.max(...headerTops) - Math.min(...headerTops)).toBeLessThanOrEqual(1);
  await region.focus();
  await page.keyboard.press('ArrowRight');
  await expect.poll(() => region.evaluate(element => element.scrollLeft)).toBeGreaterThan(0);
  const row = table.getByTestId('group-member-row-P0');
  await row.locator('td[data-column="actions"]').scrollIntoViewIfNeeded();
  await page.screenshot({ path: testInfo.outputPath(`selected-scrolled-${width}.png`) });
  const interval = row.getByTestId('group-member-cadence-P0').locator('input');
  await interval.fill('30'); await expect(interval).toHaveValue('30');
  await row.getByRole('button', { name: /移除/ }).click();
  await expect(page.getByTestId('group-member-scope')).toContainText('31 點');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await page.screenshot({ path: testInfo.outputPath(`selected-after-scroll-${width}.png`) });
});

test('the phone scroll hint stays hidden on tablet and desktop layouts', async ({ page }) => {
  await page.route('**/api/**', route => route.abort());
  for (const width of [641, 768, 1280]) {
    await page.setViewportSize({ width, height: 844 });
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
    await expect(page.getByTestId('group-member-table').getByText('左右滑動查看其餘欄位')).toBeHidden();
  }
});

test.describe('touch input', () => {
  test.use({ viewport: { width: 390, height: 844 }, isMobile: true, hasTouch: true });
  test('a horizontal swipe moves the selected points table', async ({ page, context }) => {
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
    const region = page.getByTestId('group-member-table').getByRole('region');
    await region.scrollIntoViewIfNeeded();
    const box = await region.boundingBox();
    expect(box).not.toBeNull();
    const cdp = await context.newCDPSession(page);
    await cdp.send('Input.synthesizeScrollGesture', { x: box!.x + box!.width * 0.8, y: box!.y + box!.height / 2,
      xDistance: -180, yDistance: 0, gestureSourceType: 'touch' });
    await expect.poll(() => region.evaluate(element => element.scrollLeft)).toBeGreaterThan(0);
  });
});
