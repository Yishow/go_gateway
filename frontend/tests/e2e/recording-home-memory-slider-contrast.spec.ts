import { expect, test } from '@playwright/test';

// Native range painting is not exposed by getComputedStyle. Measure its actual
// inactive track against the immediately adjacent painted surface (DESIGN.md).
for (const width of [1440, 768, 390]) for (const light of [true, false]) {
  test(`native map slider boundary contrast ${width}px ${light ? 'light' : 'dark'}`, async ({ page }, info) => {
    await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=dense256&en${light ? '&light' : ''}`);
    await page.getByTestId('home-memory-layout-map').click();
    const slider = page.getByTestId('home-memory-map-select');
    await slider.focus();
    await page.keyboard.press('Home');
    for (let index = 0; index < 127; index++) await page.keyboard.press('ArrowRight');
    await expect(slider).toHaveAttribute('aria-valuetext', 'W7F');
    const png = await slider.screenshot();
    const measurement = await page.evaluate(async base64 => {
      const image = new Image(); image.src = `data:image/png;base64,${base64}`;
      await image.decode();
      const canvas = document.createElement('canvas');
      canvas.width = image.width; canvas.height = image.height;
      const context = canvas.getContext('2d')!; context.drawImage(image, 0, 0);
      const x = Math.floor(image.width * .75), center = Math.floor(image.height / 2);
      const pixel = (y: number) => [...context.getImageData(x, y, 1, 1).data].slice(0, 3);
      const luminance = (rgb: number[]) => rgb.map(value => {
        const channel = value / 255;
        return channel <= .04045 ? channel / 12.92 : ((channel + .055) / 1.055) ** 2.4;
      }).reduce((sum, value, index) => sum + value * [.2126, .7152, .0722][index], 0);
      // W7F leaves the right quarter inactive. Native track height is 8px;
      // this scan includes its antialiased edges but excludes the outer focus ring.
      const adjacent = pixel(center - 8), bg = luminance(adjacent);
      const samples = Array.from({ length: 13 }, (_, index) => {
        const y = center - 6 + index, color = pixel(y), lightness = luminance(color);
        return { y, color, ratio: (Math.max(lightness, bg) + .05) / (Math.min(lightness, bg) + .05) };
      });
      return { adjacent, x, samples, boundary: Math.max(...samples.map(sample => sample.ratio)) };
    }, png.toString('base64'));
    await info.attach('native-slider-pixels.json', { body: JSON.stringify(measurement, null, 2), contentType: 'application/json' });
    await info.attach('native-slider.png', { body: png, contentType: 'image/png' });
    expect(measurement.boundary, 'Enabled inactive track boundary must contrast at least 3:1 with its painted surroundings').toBeGreaterThanOrEqual(3);
    await expect(slider).toBeFocused();
    expect(await slider.evaluate(node => node.matches(':focus-visible'))).toBe(true);
    await page.keyboard.press('End');
    await expect(slider).toHaveAttribute('aria-valuetext', 'WFF');
    await expect(page.getByTestId('fixture-last-action')).toHaveText('');
  });
}
