import { expect, test } from '@playwright/test';

for (const width of [390, 1353]) test(`aligns the real recording editor destination help with each label at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 820 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-editor-real-host.html');
  await expect(page.getByTestId('group-editor')).toBeVisible();
  for (const key of ['record.strategy', 'record.table', 'record.schema', 'policy.interval', 'policy.lateness', 'policy.incomplete', 'policy.dedupe', 'policy.entityColumn', 'policy.provenanceColumn']) {
    const help = page.locator(`[data-parameter-help="${key}"]`);
    await expect(help).toBeVisible();
    const geometry = await help.evaluate(button => {
      const sibling = button.previousSibling;
      const text = sibling instanceof HTMLElement && sibling.tagName === 'LABEL' ? [...sibling.childNodes].find(node => node.nodeType === Node.TEXT_NODE && node.textContent?.trim()) : sibling;
      const label = document.createRange(), glyph = document.createRange(); label.selectNodeContents(text!); glyph.selectNodeContents(button);
      const a = label.getBoundingClientRect(), b = glyph.getBoundingClientRect(), hit = button.getBoundingClientRect();
      return { gap: b.left - a.right, verticalOffset: Math.abs(a.top + a.height / 2 - b.top - b.height / 2), width: hit.width, height: hit.height };
    });
    expect.soft(geometry.gap, `${key}: ${JSON.stringify(geometry)}`).toBeGreaterThanOrEqual(0);
    expect.soft(geometry.gap, `${key}: ${JSON.stringify(geometry)}`).toBeLessThanOrEqual(12);
    expect.soft(geometry.verticalOffset, `${key}: ${JSON.stringify(geometry)}`).toBeLessThanOrEqual(5);
    expect.soft(geometry.width, key).toBeGreaterThanOrEqual(44); expect.soft(geometry.height, key).toBeGreaterThanOrEqual(44);
  }
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await page.getByTestId('group-advanced').scrollIntoViewIfNeeded();
  await page.screenshot({ path: testInfo.outputPath('destination.png') });
  const help = page.locator('[data-parameter-help="policy.incomplete"]');
  await help.click(); await expect(page.getByRole('note')).toBeVisible();
  await page.keyboard.press('Escape'); await expect(help).toHaveAttribute('aria-expanded', 'false');
  await expect(page.getByTestId('group-editor')).toBeVisible();
});

for (const width of [390, 1280]) test(`keeps editor help glyphs close to their labels at ${width}px`, async ({ page }) => {
  await page.setViewportSize({ width, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  for (const key of ['member.required', 'point.order', 'record.strategy', 'record.table', 'record.schema', 'policy.interval', 'policy.lateness', 'policy.incomplete', 'policy.dedupe', 'policy.entityColumn', 'policy.provenanceColumn']) {
    const help = page.locator(`[data-parameter-help="${key}"]`);
    const geometry = await help.evaluate(button => {
      const sibling = button.previousSibling;
      const text = sibling instanceof HTMLElement && sibling.tagName === 'LABEL' ? [...sibling.childNodes].find(node => node.nodeType === Node.TEXT_NODE && node.textContent?.trim()) : sibling;
      const label = document.createRange(), glyph = document.createRange(); label.selectNodeContents(text!); glyph.selectNodeContents(button);
      const hit = button.getBoundingClientRect();
      return { gap: glyph.getBoundingClientRect().left - label.getBoundingClientRect().right, width: hit.width, height: hit.height };
    });
    expect.soft(geometry.gap, `${key}: ${JSON.stringify(geometry)}`).toBeLessThanOrEqual(12);
    expect.soft(geometry.gap, `${key}: ${JSON.stringify(geometry)}`).toBeGreaterThanOrEqual(0);
    expect.soft(geometry.width, key).toBeGreaterThanOrEqual(44); expect.soft(geometry.height, key).toBeGreaterThanOrEqual(44);
  }
});

test('keeps destination and policy fields labelled, editable and their help usable', async ({ page }, testInfo) => {
  await page.setViewportSize({ width: 390, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const tableName = page.getByRole('textbox', { name: '資料表', exact: true });
  await tableName.fill('fixture_table'); await expect(tableName).toHaveValue('fixture_table');
  const interval = page.getByRole('spinbutton', { name: '成列間隔（秒）' });
  await interval.fill('30'); await expect(interval).toHaveValue('30');
  const help = page.locator('[data-parameter-help="policy.interval"]');
  await help.click(); await expect(page.getByRole('note')).toBeVisible();
  await page.keyboard.press('Escape'); await expect(help).toHaveAttribute('aria-expanded', 'false');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
  await page.screenshot({ path: testInfo.outputPath('policy.png') });
});

test('keeps mixed destination layout help close to its label', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
  const help = page.locator('[data-parameter-help="record.layout"]');
  const gap = await help.evaluate(button => {
    const text = [...button.previousElementSibling!.childNodes].find(node => node.nodeType === Node.TEXT_NODE && node.textContent?.trim());
    const label = document.createRange(), glyph = document.createRange(); label.selectNodeContents(text!); glyph.selectNodeContents(button);
    return glyph.getBoundingClientRect().left - label.getBoundingClientRect().right;
  });
  expect(gap).toBeGreaterThanOrEqual(0); expect(gap).toBeLessThanOrEqual(12);
  await help.click(); await expect(page.getByRole('note')).toBeVisible();
});

for (const width of [1280, 1440]) test(`keeps desktop point names balanced within the table at ${width}px`, async ({ page }) => {
  await page.setViewportSize({ width, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const table = page.getByTestId('group-member-table');
  const tableWidth = (await table.getByRole('table').boundingBox())!.width;
  const pointWidth = (await table.locator('th[data-column="pointOnly"]').boundingBox())!.width;
  expect(pointWidth).toBeGreaterThanOrEqual(200);
  expect(pointWidth / tableWidth).toBeLessThanOrEqual(0.5);
});
