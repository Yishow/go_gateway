import { expect, test } from '@playwright/test';

const fixtureUrl = (query = '') => `/tests/e2e/fixtures/recording-selected-host.html${query}`;

for (const width of [320, 390]) test(`shows the first selected point and reaches its actions at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(fixtureUrl());

  const row = page.getByTestId('group-member-row-P0');
  const region = page.getByTestId('group-member-table').getByRole('region');
  if (width === 320) await region.scrollIntoViewIfNeeded();
  const rowBox = await row.boundingBox();
  const nameBox = await row.getByRole('textbox').boundingBox(), regionBox = await region.boundingBox();
  const action = row.locator('td[data-column="actions"] button').first();
  expect(nameBox!.x).toBeGreaterThanOrEqual(regionBox!.x);
  expect(nameBox!.x + nameBox!.width, JSON.stringify({ nameBox, regionBox, device: await row.locator('td[data-column="device"]').boundingBox() })).toBeLessThanOrEqual(regionBox!.x + regionBox!.width);
  await action.scrollIntoViewIfNeeded();
  const actionBox = await action.boundingBox();
  expect(rowBox).not.toBeNull();
  expect(actionBox).not.toBeNull();
  expect(rowBox!.y).toBeLessThan(844 - 120);
  expect(actionBox!.y + actionBox!.height).toBeLessThanOrEqual(844);
  expect(await region.evaluate(element => element.scrollLeft)).toBeGreaterThan(0);
  await page.screenshot({ path: testInfo.outputPath(`selected-mobile-${width}.png`) });
});

test('shows the first selected point in a mixed recording on the first phone view', async ({ page }, testInfo) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(fixtureUrl('?mixed'));
  const point = page.getByTestId('group-member-row-P0').getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
  const box = await point.boundingBox();
  expect(box).not.toBeNull();
  expect(box!.y + box!.height).toBeLessThanOrEqual(844);
  await page.screenshot({ path: testInfo.outputPath('selected-mobile-mixed.png') });
});

test('keeps point help before the table in visual and keyboard order', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(fixtureUrl());
  const help = page.locator('[data-parameter-help="member.required"]');
  const row = page.getByTestId('group-member-row-P0');
  const helpBox = await help.boundingBox(), rowBox = await row.boundingBox();
  expect(helpBox).not.toBeNull(); expect(rowBox).not.toBeNull();
  expect(helpBox!.y + helpBox!.height).toBeLessThanOrEqual(rowBox!.y);
  await help.focus();
  await page.keyboard.press('Tab');
  await expect(page.locator('[data-parameter-help="point.order"]')).toBeFocused();
});

for (const width of [390, 768]) {
  test(`keeps a long device name readable at ${width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width, height: 844 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(fixtureUrl('?longdevice'));

    const device = page.getByTestId('group-member-row-P0').locator('td[data-column="device"]');
    await expect(device).toContainText('MC Building North');
    const metrics = await device.evaluate(element => {
      const style = getComputedStyle(element);
      return { clientWidth: element.clientWidth, scrollWidth: element.scrollWidth, whiteSpace: style.whiteSpace };
    });
    expect(metrics.scrollWidth, JSON.stringify(metrics)).toBeLessThanOrEqual(metrics.clientWidth);
    expect(metrics.whiteSpace).not.toBe('nowrap');
    await page.screenshot({ path: testInfo.outputPath(`selected-long-name-${width}.png`) });
  });
}
