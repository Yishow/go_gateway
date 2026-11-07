import { expect, test } from '@playwright/test';

// Mock-only presentation checks. No API, device, persisted settings or runtime writes.
for (const light of [true, false]) {
  test(`phone map keeps painted type and empty regions distinct ${light ? 'light' : 'dark'}`, async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en${light ? '&light' : ''}`);
    await page.getByTestId('home-memory-layout-map').click();
    const strip = page.locator('.gw-memory-map-strip');
    await expect(strip.locator('span')).toHaveCount(220);
    // Check painted pixels, not only DOM backgrounds: a 1px inset on both sides
    // previously covered every 1.58px slot even though its computed fill was correct.
    const png = (await strip.screenshot()).toString('base64');
    const colors = await page.evaluate(async base64 => {
      const image = new Image();
      image.src = `data:image/png;base64,${base64}`;
      await image.decode();
      const canvas = document.createElement('canvas');
      canvas.width = image.width; canvas.height = image.height;
      const context = canvas.getContext('2d')!;
      context.drawImage(image, 0, 0);
      // Known fixture addresses: signed 40009, float 40025, empty 40061.
      return [8, 24, 60].map(index => [...context.getImageData(
        Math.floor((index + 0.5) * image.width / 220), Math.floor(image.height / 2), 1, 1,
      ).data]);
    }, png);
    expect(colors[0], 'Signed and empty slots must not paint the same color').not.toEqual(colors[2]);
    expect(colors[1], 'Float and empty slots must not paint the same color').not.toEqual(colors[2]);
    expect(colors[0], 'Signed and float regions remain distinguishable').not.toEqual(colors[1]);
  });

  test(`sample notice remains readable ${light ? 'light' : 'dark'}`, async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en${light ? '&light' : ''}`);
    const notice = page.getByTestId('fixture-sample-notice');
    await expect(notice).toBeVisible();
    // The production header hides explanatory paragraphs in a 1px clipped box.
    // Text presence (and even visibility alone) cannot prove this sample warning is readable.
    const bounds = await notice.boundingBox();
    expect(bounds?.width).toBeGreaterThan(100);
    expect(bounds?.height).toBeGreaterThan(15);
  });
}
