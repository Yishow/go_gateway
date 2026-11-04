import { expect, test } from '@playwright/test';

const pickerFixture = '/tests/e2e/fixtures/recording-point-name-host.html';

test('keeps point identity available when the mobile picker scrolls horizontally', async ({ page }, testInfo) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(pickerFixture);
  await page.getByTestId('recording-points-open').click();

  const dialog = page.getByRole('dialog');
  const body = dialog.locator('.recording-point-picker__body');
  const firstRow = dialog.locator('tbody tr').first();
  const checkbox = firstRow.getByRole('checkbox');
  const name = firstRow.getByRole('textbox');

  await body.evaluate(element => { element.scrollLeft = element.scrollWidth; });
  const bounds = await checkbox.evaluate(element => {
    const input = element.getBoundingClientRect();
    const viewport = element.closest('.recording-point-picker__body')!.getBoundingClientRect();
    return { inputLeft: input.left, inputRight: input.right, viewportLeft: viewport.left, viewportRight: viewport.right };
  });
  const notes = await dialog.locator('.recording-point-picker__meaning-note').evaluateAll(elements => elements.map(element => {
    const note = element.getBoundingClientRect();
    const viewport = element.closest('.recording-point-picker__body')!.getBoundingClientRect();
    return { noteLeft: note.left, noteRight: note.right, viewportLeft: viewport.left, viewportRight: viewport.right };
  }));

  expect(bounds.inputLeft).toBeGreaterThanOrEqual(bounds.viewportLeft);
  expect(bounds.inputRight).toBeLessThanOrEqual(bounds.viewportRight);
  for (const note of notes) {
    expect(note.noteLeft).toBeGreaterThanOrEqual(note.viewportLeft);
    expect(note.noteRight).toBeLessThanOrEqual(note.viewportRight);
  }
  await expect(name).toBeVisible();
  await page.screenshot({ path: testInfo.outputPath('picker-mobile-scrolled.png') });
});

test('gives inline point selection its own 44px touch target', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(pickerFixture);
  await page.getByTestId('recording-points-open').click();

  const firstRow = page.getByRole('dialog').locator('tbody tr').first();
  const checkbox = firstRow.getByRole('checkbox');
  const label = checkbox.locator('xpath=ancestor::label');
  const target = await label.boundingBox();

  expect(target).not.toBeNull();
  expect(target!.width).toBeGreaterThanOrEqual(44);
  expect(target!.height).toBeGreaterThanOrEqual(44);
  await label.click({ position: { x: 3, y: 22 } });
  await expect(checkbox).toBeChecked();
  await label.click({ position: { x: 3, y: 22 } });
  await expect(checkbox).not.toBeChecked();

  const name = firstRow.getByRole('textbox');
  await name.click();
  await expect(name).toBeEditable();
  await name.fill('名稱編輯不應選取');
  await expect(checkbox).not.toBeChecked();
});
