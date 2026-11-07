import { expect, test } from '@playwright/test';

for (const width of [768, 1000, 1024]) test(`tablet users can scroll an aligned selected-points table at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed&theme=light');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  await table.scrollIntoViewIfNeeded();
  await page.screenshot({ path: testInfo.outputPath(`tablet-initial-${width}.png`) });
  const headers = await region.locator('thead th').evaluateAll(cells => cells.map(cell => cell.getBoundingClientRect().top));
  expect(Math.max(...headers) - Math.min(...headers)).toBeLessThanOrEqual(1);
  expect(await region.evaluate(element => element.scrollWidth)).toBeGreaterThan(await region.evaluate(element => element.clientWidth));
  await region.focus();
  await page.keyboard.press('ArrowRight');
  await expect.poll(() => region.evaluate(element => element.scrollLeft)).toBeGreaterThan(0);
  await page.keyboard.press('ArrowLeft');
  await expect.poll(() => region.evaluate(element => element.scrollLeft)).toBe(0);
  const row = table.getByTestId('group-member-row-P0');
  const interval = row.getByTestId('group-member-cadence-P0').locator('input');
  await interval.fill('30');
  await expect(interval).toHaveValue('30');
  await interval.press('End');
  await interval.press('ArrowLeft');
  expect(await interval.evaluate(element => (element as HTMLInputElement).selectionStart)).toBe(1);
  await row.getByRole('button', { name: /移除/ }).click();
  await expect(page.getByTestId('group-member-scope')).toContainText('31 點');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await page.screenshot({ path: testInfo.outputPath(`tablet-scrolled-${width}.png`) });
});

for (const width of [641, 768, 1000]) test(`keeps each point identifiable beside its actions at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed&theme=light&samepoint');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  await region.scrollIntoViewIfNeeded();
  await region.evaluate(element => { element.scrollLeft = element.scrollWidth; });
  const regionBox = await region.boundingBox();
  expect(regionBox).not.toBeNull();
  for (const [id, device] of [['P0', 'MC A'], ['P1', 'MC B']] as const) {
    const row = table.getByTestId(`group-member-row-${id}`);
    const point = row.locator('td[data-column="pointOnly"]').getByRole('textbox');
    const deviceCell = row.locator('td[data-column="device"]');
    const action = row.getByRole('button', { name: /移除/ });
    await expect(point).toHaveValue('Power (D0)');
    await expect(deviceCell).toContainText(device);
    await expect(deviceCell.getByText('D0', { exact: true })).toBeVisible({ timeout: 2000 });
    const pointBox = await point.boundingBox(), actionBox = await action.boundingBox(), deviceBox = await deviceCell.boundingBox();
    expect(pointBox).not.toBeNull(); expect(actionBox).not.toBeNull(); expect(deviceBox).not.toBeNull();
    expect(deviceBox!.x).toBeGreaterThanOrEqual(regionBox!.x);
    expect(deviceBox!.x + deviceBox!.width).toBeLessThanOrEqual(pointBox!.x);
    expect(pointBox!.x).toBeGreaterThanOrEqual(regionBox!.x);
    expect(pointBox!.x + pointBox!.width).toBeLessThanOrEqual(actionBox!.x);
    expect(actionBox!.x + actionBox!.width).toBeLessThanOrEqual(regionBox!.x + regionBox!.width);
  }
  await page.screenshot({ path: testInfo.outputPath(`tablet-identities-${width}.png`) });
  const secondRow = table.getByTestId('group-member-row-P1');
  const interval = secondRow.getByTestId('group-member-cadence-P1').locator('input');
  const pointCell = secondRow.locator('td[data-column="pointOnly"]');
  const pointBox = await pointCell.boundingBox(), intervalBox = await interval.boundingBox();
  expect(pointBox).not.toBeNull(); expect(intervalBox).not.toBeNull();
  const adjustment = Math.min(0, intervalBox!.x - pointBox!.x - pointBox!.width - 8);
  await region.evaluate((element, distance) => { element.scrollLeft += distance; }, adjustment);
  await interval.click();
  await interval.fill('30');
  await expect(interval).toHaveValue('30');
  const visiblePoint = await pointCell.boundingBox(), visibleInterval = await interval.boundingBox();
  expect(visiblePoint!.x).toBeGreaterThanOrEqual(regionBox!.x);
  expect(visiblePoint!.x + visiblePoint!.width).toBeLessThanOrEqual(visibleInterval!.x);
  expect(visibleInterval!.x + visibleInterval!.width).toBeLessThanOrEqual(regionBox!.x + regionBox!.width);
  await page.screenshot({ path: testInfo.outputPath(`tablet-interval-${width}.png`) });
  await region.evaluate(element => { element.scrollLeft = element.scrollWidth; });
  await table.getByTestId('group-member-row-P1').getByRole('button', { name: /移除/ }).click();
  await expect(table.getByTestId('group-member-row-P1')).toHaveCount(0);
  await expect(table.getByTestId('group-member-row-P0')).toBeVisible();
});

for (const width of [1025, 1194]) test(`keeps point names readable on a wide tablet at ${width}px`, async ({ page }) => {
  await page.setViewportSize({ width, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  const point = table.getByTestId('group-member-row-P0').locator('td[data-column="pointOnly"]').getByRole('textbox');
  expect((await point.boundingBox())!.width).toBeGreaterThanOrEqual(120);
  const headers = await region.locator('thead th').evaluateAll(cells => cells.map(cell => cell.getBoundingClientRect().top));
  expect(Math.max(...headers) - Math.min(...headers)).toBeLessThanOrEqual(1);
  await table.getByTestId('group-member-row-P0').getByRole('button', { name: /移除/ }).click();
  await expect(page.getByTestId('group-member-scope')).toContainText('31 點');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test('a horizontal wheel gesture moves the tablet table', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
  const region = page.getByTestId('group-member-table').getByRole('region');
  await region.scrollIntoViewIfNeeded();
  const box = await region.boundingBox();
  expect(box).not.toBeNull();
  await page.mouse.move(box!.x + box!.width / 2, box!.y + box!.height / 2);
  await page.mouse.wheel(250, 0);
  await expect.poll(() => region.evaluate(element => element.scrollLeft)).toBeGreaterThan(0);
});

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

test('the scroll hint follows whether the table has hidden columns', async ({ page }) => {
  await page.route('**/api/**', route => route.abort());
  await page.setViewportSize({ width: 1280, height: 844 });
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  const hint = table.getByText('左右滑動查看其餘欄位');
  for (const [width, scrollable] of [[1280, false], [640, true], [1280, false]] as const) {
    await page.setViewportSize({ width, height: 844 });
    await expect.poll(() => region.evaluate(element => element.scrollWidth > element.clientWidth)).toBe(scrollable);
    if (scrollable) await expect(hint).toBeVisible();
    else await expect(hint).toBeHidden();
  }
});

test('updates the scroll hint after search results return without ResizeObserver', async ({ page }) => {
  await page.addInitScript(() => { Object.defineProperty(window, 'ResizeObserver', { value: undefined, configurable: true }); });
  await page.route('**/api/**', route => route.abort());
  await page.setViewportSize({ width: 1280, height: 844 });
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
  expect(await page.evaluate(() => typeof window.ResizeObserver)).toBe('undefined');
  const table = page.getByTestId('group-member-table');
  const region = table.getByRole('region');
  const search = table.getByTestId('group-member-search');
  const hint = table.getByText('左右滑動查看其餘欄位');
  await search.fill('no-such-point');
  await expect(table.getByTestId('group-member-empty')).toBeVisible();
  await page.setViewportSize({ width: 768, height: 844 });
  await expect.poll(() => region.evaluate(element => element.scrollWidth > element.clientWidth)).toBe(false);
  await page.evaluate(() => new Promise<void>(resolve => requestAnimationFrame(() => requestAnimationFrame(() => resolve()))));
  await page.evaluate(() => window.dispatchEvent(new Event('resize')));
  await page.evaluate(() => new Promise<void>(resolve => requestAnimationFrame(() => requestAnimationFrame(() => resolve()))));
  await search.clear();
  await expect(table.getByTestId('group-member-row-P0')).toBeVisible();
  await expect.poll(() => region.evaluate(element => element.scrollWidth > element.clientWidth)).toBe(true);
  await expect(hint).toBeVisible({ timeout: 2000 });
});

test.describe('touch input', () => {
  test.use({ isMobile: true, hasTouch: true });
  for (const width of [390, 768, 1024]) test(`a horizontal swipe moves the selected points table at ${width}px`, async ({ page, context }) => {
    await page.setViewportSize({ width, height: 844 });
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
