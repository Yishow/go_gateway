import { expect, test } from '@playwright/test';

// MOCK SAMPLE only. Exercise the real modal CSS without connecting to a device.
test('memory modal follows 85% of the desktop viewport across resizes', async ({ page }) => {
  await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en');
  const dialog = page.getByRole('dialog');
  for (const viewport of [{ width: 1920, height: 1080 }, { width: 1440, height: 1000 }, { width: 768, height: 1000 }]) {
    await page.setViewportSize(viewport);
    await expect.poll(async () => {
      const box = await dialog.boundingBox();
      return box && Math.abs(box.width - viewport.width * 0.85) < 1
        && Math.abs(box.height - viewport.height * 0.85) < 1
        && Math.abs(box.x - viewport.width * 0.075) < 1
        && Math.abs(box.y - viewport.height * 0.075) < 1;
    }).toBe(true);
  }
});

for (const viewport of [{ width: 390, height: 844 }, { width: 320, height: 568 }, { width: 844, height: 390 }]) {
  test(`memory modal stays inside small viewport ${viewport.width}x${viewport.height}`, async ({ page }) => {
    await page.setViewportSize(viewport);
    await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en');
    const dialog = page.getByRole('dialog');
    const box = await dialog.boundingBox();
    expect(box).not.toBeNull();
    expect(box!.x).toBeGreaterThanOrEqual(8);
    expect(box!.y).toBeGreaterThanOrEqual(8);
    expect(box!.x + box!.width).toBeLessThanOrEqual(viewport.width - 8);
    expect(box!.y + box!.height).toBeLessThanOrEqual(viewport.height - 8);
    expect(box!.width).toBeCloseTo(viewport.width - 16, 0);
    expect(box!.height).toBeCloseTo(viewport.height - 16, 0);
    await expect(dialog.locator('.gw-modal-footer')).toBeInViewport();
    const body = page.getByTestId('home-memory-view');
    expect(await body.evaluate(node => node.clientWidth >= node.scrollWidth)).toBe(true);
    await dialog.getByTestId('home-memory-close').click();
    await expect(dialog).toHaveCount(0);
    await page.getByTestId('fixture-reopen').click();
    await expect(page.getByRole('dialog')).toBeVisible();
  });
}

test('memory sizing does not change other source forms or confirmation dialogs', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.goto('/tests/e2e/fixtures/recording-home-memory-host.html?scenario=reference&en');
  const dialog = page.getByRole('dialog');
  for (const kind of ['points', 'connection']) {
    await dialog.evaluate((node, sourceKind) => {
      node.removeAttribute('data-source-view');
      node.setAttribute('data-source-kind', sourceKind);
    }, kind);
    const box = await dialog.boundingBox();
    expect(box!.width).toBeCloseTo(1000, 0);
    expect(box!.height).toBeCloseTo(750, 0);
  }
  await dialog.evaluate(node => {
    node.setAttribute('data-source-view', 'memory');
    node.setAttribute('data-size', 'confirm');
  });
  expect((await dialog.boundingBox())!.width).toBeCloseTo(420, 0);
});
