import { expect, test } from '@playwright/test';

const kinds = ['form', 'memory', 'radix-confirm', 'legacy', 'confirm', 'slide'];
const sizes = [{ width: 390, height: 844 }, { width: 768, height: 1024 }, { width: 1024, height: 768 }, { width: 1366, height: 400 }];
for (const kind of kinds) for (const viewport of sizes) {
  test(`${kind} fills ${viewport.width}x${viewport.height} and preserves reachable controls`, async ({ page }) => {
    await page.setViewportSize(viewport);
    await page.goto(`/tests/e2e/fixtures/responsive-modal-host.html?kind=${kind}`);
    await page.getByRole('button', { name: '開啟測試' }).click();
    const dialog = page.locator('.gw-modal, [aria-modal=true]');
    await expect(dialog).toBeVisible();
    const box = (await dialog.boundingBox())!;
    expect(box.x).toBeCloseTo(0, 0); expect(box.y).toBeCloseTo(0, 0);
    expect(box.width).toBeCloseTo(viewport.width, 0); expect(box.height).toBeCloseTo(viewport.height, 0);
    expect(await dialog.evaluate(el => el.clientWidth >= el.scrollWidth)).toBe(true);
    const body = dialog.locator('.gw-modal-body, .gw-dialog-body');
    await body.evaluate(el => { el.scrollTop = el.scrollHeight; });
    await expect(dialog.locator('.gw-modal-header, .gw-dialog-header')).toBeInViewport();
    if (kind !== 'slide') await expect(dialog.locator('.gw-modal-footer, .gw-dialog-footer')).toBeInViewport();
    await page.keyboard.press('Escape');
    await expect(dialog).toHaveCount(0);
    expect((await page.getByTestId('non-modal').boundingBox())!.width).toBe(200);
  });
}
for (const hasTouch of [false, true]) test(`1366 tablet capability is explicit: touch=${hasTouch}`, async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 1366, height: 1024 }, hasTouch });
  const page = await context.newPage();
  try {
    await page.goto('http://127.0.0.1:4173/tests/e2e/fixtures/responsive-modal-host.html?kind=memory');
    await page.getByRole('button', { name: '開啟測試' }).click();
    expect((await page.locator('.gw-modal, [aria-modal=true]').boundingBox())!.width).toBeCloseTo(hasTouch ? 1366 : 1366 * 0.85, 0);
  } finally { await context.close(); }
});
test('rotating a tablet preserves an in-progress field and returns to desktop sizing', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 1024 });
  await page.goto('/tests/e2e/fixtures/responsive-modal-host.html?kind=form');
  await page.getByRole('button', { name: '開啟測試' }).click();
  await page.getByRole('textbox').fill('尚未儲存');
  await page.setViewportSize({ width: 1024, height: 768 });
  expect((await page.locator('.gw-modal, [aria-modal=true]').boundingBox())!.width).toBeCloseTo(1024, 0);
  await expect(page.getByRole('textbox')).toHaveValue('尚未儲存');
  await page.setViewportSize({ width: 1440, height: 900 });
  expect((await page.locator('.gw-modal, [aria-modal=true]').boundingBox())!.width).toBeCloseTo(1080, 0);
});

for (const viewport of [{ width: 390, height: 844 }, { width: 1024, height: 768 }]) test(`confirmation text scrolls from the keyboard at ${viewport.width}px`, async ({ page }) => {
  await page.setViewportSize(viewport);
  await page.goto('/tests/e2e/fixtures/responsive-modal-host.html?kind=confirm');
  await page.getByRole('button', { name: '開啟測試' }).click();
  const dialog = page.locator('.gw-modal, [aria-modal=true]'), body = dialog.locator('.gw-dialog-body');
  // Cancel is the deliberate initial focus. Shift+Tab must reach the scroll region.
  await expect(dialog.getByRole('button', { name: '取消' })).toBeFocused();
  await page.keyboard.press('Shift+Tab');
  await expect(body).toBeFocused();
  await page.keyboard.press('PageDown');
  await expect.poll(() => body.evaluate(el => el.scrollTop)).toBeGreaterThan(0);
  await page.keyboard.press('End');
  await expect.poll(() => body.evaluate(el => el.scrollHeight - el.clientHeight - el.scrollTop)).toBeLessThanOrEqual(2);
  await expect(dialog.locator('.gw-dialog-footer')).toBeInViewport();
  await page.keyboard.press('Tab');
  await expect(dialog.getByRole('button', { name: '取消' })).toBeFocused();
});

test('fullscreen text controls use at least 16px', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 1024 });
  await page.goto('/tests/e2e/fixtures/responsive-modal-host.html?kind=form');
  await page.getByRole('button', { name: '開啟測試' }).click();
  expect(await page.getByRole('textbox').evaluate(el => parseFloat(getComputedStyle(el).fontSize))).toBeGreaterThanOrEqual(16);
});
