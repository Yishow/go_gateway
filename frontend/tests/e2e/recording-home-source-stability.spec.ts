import { expect, test } from '@playwright/test';

for (const viewport of [{ width: 1440, height: 900 }, { width: 768, height: 1024 }, { width: 390, height: 844 },
  { width: 320, height: 568 }, { width: 900, height: 430 }]) {
  test(`home source dialog stays in place while a saved range loads at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize(viewport);
    const mutations: string[] = [];
    let releaseList!: () => void;
    const listGate = new Promise<void>(resolve => { releaseList = resolve; });
    await page.route('**/api/**', async route => {
      const request = route.request(), path = new URL(request.url()).pathname;
      if (request.method() !== 'GET') { mutations.push(request.method()); return route.abort(); }
      if (path === '/api/v1/datalink/studio-v2/workspace') return route.fulfill({ json: { success: true, data: { id: 'W', ordered_device_ids: ['D'] } } });
      if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') {
        await listGate;
        return route.fulfill({ json: { success: true, data: [{ id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'R2',
          start_address: '40001', count: 16, data_type: 'int16', naming_prefix: 'BLOCK_', enabled: true, skipped_addresses: [],
          locked: false, origin: 'manual', created_at: '2026-10-29T00:00:00Z', updated_at: '2026-10-29T00:00:00Z' }] } });
      }
      return route.abort();
    });
    await page.goto('/tests/e2e/fixtures/recording-home-source-host.html');
    await page.getByRole('button', { name: '讀取點位' }).click();
    const dialog = page.getByRole('dialog');
    await expect(dialog.getByRole('status')).toBeVisible();
    const loading = await dialog.boundingBox();
    await dialog.evaluate(element => {
      const sizes: number[][] = [];
      const observer = new ResizeObserver(() => { const { x, y, width, height } = element.getBoundingClientRect(); sizes.push([x, y, width, height]); });
      observer.observe(element);
      (window as unknown as { sourceDialogSizes: number[][] }).sourceDialogSizes = sizes;
    });
    await page.screenshot({ path: testInfo.outputPath(`loading-${viewport.width}.png`) });
    releaseList();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('BLOCK_');
    const loaded = await dialog.boundingBox();
    await page.screenshot({ path: testInfo.outputPath(`loaded-${viewport.width}.png`) });
    expect(loading).not.toBeNull(); expect(loaded).not.toBeNull();
    expect(Math.abs(loaded!.x - loading!.x)).toBeLessThanOrEqual(1);
    expect(Math.abs(loaded!.width - loading!.width)).toBeLessThanOrEqual(1);
    expect(Math.abs(loaded!.y - loading!.y)).toBeLessThanOrEqual(1);
    expect(Math.abs(loaded!.height - loading!.height)).toBeLessThanOrEqual(1);
    const sizes = await page.evaluate(() => (window as unknown as { sourceDialogSizes: number[][] }).sourceDialogSizes);
    expect(sizes.length).toBeGreaterThan(0);
    for (const [x, y, width, height] of sizes) {
      expect(Math.abs(x - loading!.x)).toBeLessThanOrEqual(1);
      expect(Math.abs(y - loading!.y)).toBeLessThanOrEqual(1);
      expect(Math.abs(width - loading!.width)).toBeLessThanOrEqual(1);
      expect(Math.abs(height - loading!.height)).toBeLessThanOrEqual(1);
    }
    await expect(dialog.getByTestId('home-range-done')).toBeInViewport({ ratio: 1 });
    await expect(dialog.getByTestId('home-range-save')).toBeInViewport({ ratio: 1 });
    await expect(dialog.getByTestId('home-range-save')).toBeDisabled();
    await dialog.getByTestId('rule-prefix-input').fill('EDIT_');
    await expect(dialog.getByTestId('home-range-save')).toBeEnabled();
    await dialog.getByTestId('home-range-done').click();
    await expect(dialog.getByTestId('home-source-keep')).toBeVisible();
    await dialog.getByTestId('home-source-keep').click();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('EDIT_');
    expect(mutations).toEqual([]);
  });
}

test('a failed range read keeps the dialog stable and allows closing', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  let releaseList!: () => void;
  const listGate = new Promise<void>(resolve => { releaseList = resolve; });
  await page.route('**/api/**', async route => {
    const path = new URL(route.request().url()).pathname;
    if (route.request().method() !== 'GET') return route.abort();
    if (path === '/api/v1/datalink/studio-v2/workspace') return route.fulfill({ json: { success: true, data: { id: 'W', ordered_device_ids: ['D'] } } });
    if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') {
      await listGate;
      return route.fulfill({ status: 503, json: { success: false } });
    }
    return route.abort();
  });
  await page.goto('/tests/e2e/fixtures/recording-home-source-host.html');
  const trigger = page.getByRole('button', { name: '讀取點位' });
  await trigger.click();
  const dialog = page.getByRole('dialog');
  await expect(dialog.getByRole('status')).toBeVisible();
  const loading = await dialog.boundingBox();
  releaseList();
  await expect(dialog.getByRole('alert')).toBeVisible();
  const failed = await dialog.boundingBox();
  expect(loading).not.toBeNull(); expect(failed).not.toBeNull();
  expect(Math.abs(failed!.height - loading!.height)).toBeLessThanOrEqual(1);
  expect(Math.abs(failed!.y - loading!.y)).toBeLessThanOrEqual(1);
  await dialog.getByTestId('home-source-close').click();
  await expect(dialog).toHaveCount(0);
  await expect(trigger).toBeFocused();
});
