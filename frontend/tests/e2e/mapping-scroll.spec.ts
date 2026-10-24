import { expect, test } from '@playwright/test';

for (const advanced of [false, true]) {
  test(`${advanced ? 'advanced' : 'basic'} mapping table boundary scrolling does not move the surrounding page`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.goto('/tests/e2e/fixtures/mapping-scroll.html');
    const table = page.getByTestId('mapping-table-wrapper').locator('table');
    const scroller = table.locator('..');
    await expect(page.getByTestId('mapping-row-point-47')).toBeAttached();
    if (advanced) await page.getByTestId('mapping-advanced-toggle').click();
    const box = await scroller.boundingBox();
    expect(box).not.toBeNull();
    await page.mouse.move(box!.x + box!.width / 2, box!.y + box!.height / 2);

    // Exercise real wheel scrolling, then keep scrolling past the table's lower edge.
    await page.mouse.wheel(0, 200);
    await expect.poll(() => scroller.evaluate(element => element.scrollTop)).toBeGreaterThan(0);
    await expect.poll(() => page.evaluate(() => window.scrollY)).toBe(0);
    await scroller.evaluate(element => { element.scrollTop = element.scrollHeight; });
    await page.mouse.wheel(0, 500);
    await page.waitForTimeout(300); // Allow the browser's asynchronous wheel scroll to settle.
    expect(await page.evaluate(() => window.scrollY)).toBe(0);
    await page.screenshot({ path: testInfo.outputPath('scroll-bottom.png') });

    // The table's top edge also stays contained, while scrolling outside remains available.
    await page.evaluate(() => window.scrollTo(0, 100));
    await scroller.evaluate(element => { element.scrollTop = 0; });
    const movedBox = await scroller.boundingBox();
    await page.mouse.move(movedBox!.x + movedBox!.width / 2, movedBox!.y + 30);
    await page.mouse.wheel(0, -500);
    await page.waitForTimeout(300);
    expect(await page.evaluate(() => window.scrollY)).toBe(100);
    await page.mouse.move(5, 100);
    await page.mouse.wheel(0, 150);
    await expect.poll(() => page.evaluate(() => window.scrollY)).toBeGreaterThan(100);
  });
}
