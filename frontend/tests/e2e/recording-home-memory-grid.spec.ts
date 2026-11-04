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
