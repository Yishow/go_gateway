import { expect, test } from '@playwright/test';

for (const language of ['zh-TW', 'en']) {
  for (const width of [1440, 1024, 375]) {
    test.describe(`${language} point grid at ${width}px`, () => {
      test.use({ viewport: { width, height: 900 }, hasTouch: width === 375 });

      test('keeps each collection control inside its point and independently operable', async ({ page }, testInfo) => {
        // Real Studio shell and Step 2; API interception makes this a UI-only test.
        await page.route('**/api/**', route => route.fulfill({ json: { data: [] } }));
        await page.goto(`/tests/e2e/fixtures/studio-mapping-scroll.html?language=${language}`);
        await page.getByTestId('step-nav-button-2').click();
        const grid = page.getByTestId('point-grid');
        const cells = grid.locator('[data-testid^="point-cell-"]');
        await expect(cells).toHaveCount(48);
        await page.evaluate(() => document.fonts.ready);

        const geometry = await cells.evaluateAll(elements => elements.map(element => {
          const cell = element.getBoundingClientRect();
          const label = element.querySelector('label')!;
          const control = label.getBoundingClientRect();
          const address = element.querySelector(':scope > div:first-child span')!.getBoundingClientRect();
          const neighbours = elements.filter(other => other !== element).map(other => other.getBoundingClientRect());
          return {
            point: element.getAttribute('data-address'),
            contained: control.left >= cell.left && control.right <= cell.right
              && control.top >= cell.top && control.bottom <= cell.bottom,
            overlapsNeighbour: neighbours.some(other => control.left < other.right && control.right > other.left
              && control.top < other.bottom && control.bottom > other.top),
            overlapsAddress: control.left < address.right && control.right > address.left
              && control.top < address.bottom && control.bottom > address.top,
            fullyVisible: label.scrollWidth <= label.clientWidth && label.scrollHeight <= label.clientHeight,
            touchHeight: control.height,
            touchWidth: control.width,
          };
        }));
        await testInfo.attach('point-grid-geometry', { body: JSON.stringify(geometry, null, 2), contentType: 'application/json' });
        expect(geometry.filter(cell => !cell.contained || cell.overlapsNeighbour || cell.overlapsAddress || !cell.fullyVisible)).toEqual([]);
        if (width === 375) expect(geometry.every(cell => cell.touchHeight >= 44 && cell.touchWidth >= 44)).toBe(true);

        const firstCell = cells.first();
        const checkbox = firstCell.getByRole('checkbox', { name: language === 'en' ? 'Collect' : '採集' });
        await expect(checkbox).toBeChecked();
        await firstCell.focus();
        await firstCell.press('Space');
        await expect(firstCell).toHaveAttribute('data-selected', 'true');
        await expect(checkbox).toBeChecked();
        await checkbox.focus();
        await checkbox.press('Space');
        await expect(checkbox).not.toBeChecked();
        await expect(firstCell).toHaveAttribute('data-selected', 'true');
        await firstCell.locator('label').click();
        await expect(checkbox).toBeChecked();
        await expect(cells.nth(1).getByRole('checkbox')).toBeChecked();
        await grid.screenshot({ path: testInfo.outputPath('point-grid.png') });
      });
    });
  }
}
