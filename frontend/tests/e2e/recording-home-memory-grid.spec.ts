import { expect, test } from '@playwright/test';

// Isolated fixture: in-memory ranges only, no API. Keyboard and wrap cues follow the rendered column count.
for (const viewport of [{ name: 'tablet', width: 768, height: 1024, columns: 10 }, { name: 'desktop', width: 1440, height: 900, columns: 16 }]) {
  test(`W area keyboard and multi-register wrap follow the ${viewport.columns} rendered columns at ${viewport.name}`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width: viewport.width, height: viewport.height });
    await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=w&light&en');
    const grid = page.getByRole('group', { name: /^W addresses/ });
    await expect(grid).toBeVisible();
    const rendered = await grid.evaluate(node => getComputedStyle(node).gridTemplateColumns.split(' ').filter(Boolean).length);
    expect(rendered).toBe(viewport.columns);
    await page.getByTestId('home-memory-cell-W0').focus();
    await page.keyboard.press('ArrowDown');
    await expect(page.getByTestId(viewport.columns === 10 ? 'home-memory-cell-WA' : 'home-memory-cell-W10')).toBeFocused();
    await page.keyboard.press('ArrowUp');
    await expect(page.getByTestId('home-memory-cell-W0')).toBeFocused();
    const w9 = page.getByTestId('home-memory-cell-W9'), wa = page.getByTestId('home-memory-cell-WA');
    if (viewport.columns === 10) {
      await expect(w9).toContainText('next row'); await expect(w9).not.toHaveAttribute('data-join-next', 'true');
      await expect(wa).toContainText('cont. W9');
    } else {
      await expect(w9).toHaveAttribute('data-join-next', 'true'); await expect(w9).not.toContainText('next row');
      await expect(wa).not.toContainText('cont. W9');
    }
    await page.screenshot({ path: testInfo.outputPath(`w-area-${viewport.name}.png`) });
  });
}

test('area tabs stay on one scrollable row at tablet width and keyboard selection scrolls into view', async ({ page }, testInfo) => {
  await page.setViewportSize({ width: 768, height: 1024 });
  await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=mc&light');
  const tablist = page.getByRole('tablist');
  const tabs = tablist.getByRole('tab');
  await expect(tabs).toHaveCount(6);
  const rows = await tabs.evaluateAll(nodes => new Set(nodes.map(node => Math.round(node.getBoundingClientRect().top))).size);
  expect(rows).toBe(1);
  // More tabs to the right are announced visually before any scrolling.
  await expect(tablist).toHaveAttribute('data-more-after', 'true');
  await page.screenshot({ path: testInfo.outputPath('tabs-tablet-start.png') });
  await tabs.first().focus();
  await page.keyboard.press('End');
  const last = tablist.getByRole('tab', { name: /^Y/ });
  await expect(last).toBeFocused(); await expect(last).toHaveAttribute('aria-selected', 'true');
  await expect(last).toBeInViewport({ ratio: 0.99 });
  await expect(tablist).toHaveAttribute('data-more-before', 'true');
  await page.screenshot({ path: testInfo.outputPath('tabs-tablet-end.png') });
  await page.keyboard.press('ArrowRight');
  const first = tablist.getByRole('tab', { name: /^D/ });
  await expect(first).toBeFocused(); await expect(first).toBeInViewport({ ratio: 0.99 });
  await page.screenshot({ path: testInfo.outputPath('tabs-tablet-after-keys.png') });
});

test('pointer selection of a half-visible area tab brings the whole tab into view', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 1024 });
  await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=mc&light');
  const tablist = page.getByRole('tablist'), last = tablist.getByRole('tab', { name: /^Y/ });
  await expect(tablist).toHaveAttribute('data-more-after', 'true');
  // Click the visible sliver of the clipped tab by coordinates, so Playwright does not scroll it first.
  const [list, tab] = await Promise.all([tablist.boundingBox(), last.boundingBox()]);
  expect(tab!.x + tab!.width).toBeGreaterThan(list!.x + list!.width);
  const visibleRight = list!.x + list!.width;
  await page.mouse.click((tab!.x + visibleRight) / 2, tab!.y + tab!.height / 2);
  await expect(last).toHaveAttribute('aria-selected', 'true');
  await expect.poll(() => tablist.evaluate(node => node.scrollLeft)).toBeGreaterThan(0);
  await expect(last).toBeInViewport({ ratio: 0.99 });
  const after = await last.boundingBox(), box = await tablist.boundingBox();
  expect(after!.x + after!.width).toBeLessThanOrEqual(box!.x + box!.width + 0.5);
});

test('pointer re-selection of the already selected, half-visible area tab brings it back into view', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 1024 });
  await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=mc&light');
  const tablist = page.getByRole('tablist'), first = tablist.getByRole('tab', { name: /^D/ });
  await expect(first).toHaveAttribute('aria-selected', 'true');
  await tablist.evaluate(node => { node.scrollLeft = node.scrollWidth; });
  await expect.poll(() => tablist.evaluate(node => node.scrollLeft)).toBeGreaterThan(0);
  // D is selected but clipped on the left; click its visible sliver by coordinates.
  const [list, tab] = await Promise.all([tablist.boundingBox(), first.boundingBox()]);
  expect(tab!.x).toBeLessThan(list!.x);
  await page.mouse.click((list!.x + tab!.x + tab!.width) / 2, tab!.y + tab!.height / 2);
  await expect(first).toHaveAttribute('aria-selected', 'true');
  await expect.poll(() => tablist.evaluate(node => node.scrollLeft)).toBe(0);
  await expect(first).toBeInViewport({ ratio: 0.99 });
  const after = await first.boundingBox(), box = await tablist.boundingBox();
  expect(after!.x).toBeGreaterThanOrEqual(box!.x - 0.5);
});
