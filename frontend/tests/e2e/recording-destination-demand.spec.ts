import { expect, test } from '@playwright/test';

for (const width of [390, 1360]) test(`editing sources defers the destination catalog at ${width}px`, async ({ page }, testInfo) => {
  await page.setViewportSize({ width, height: 600 });
  let catalogReads = 0;
  await page.route('**/api/**', route => {
    const pathname = new URL(route.request().url()).pathname;
    if (pathname.endsWith('/database-metadata')) return route.fulfill({ json: { success: true, data: {
      workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: 'test_database',
      schema: 'public', table: 'gw_group_fixture', inspection_status: 'exists', columns: [],
    } } });
    if (!pathname.endsWith('/database-catalog')) {
      return route.fulfill({ status: 503, json: { success: false } });
    }
    catalogReads += 1;
    if (catalogReads === 1) return route.fulfill({ status: 503, json: { success: false } });
    return route.fulfill({ json: { success: true, data: {
      workspace_id: 'ws-1', workspace_revision: 'workspace-1', connector_id: 'conn-1', connector_revision: 'crev-1',
      kind: 'postgres', default_schema: 'public', schemas: [{ name: 'public', tables: ['gw_group_fixture'] }],
    } } });
  });

  await page.goto('/tests/e2e/fixtures/recording-catalog-host.html', { waitUntil: 'networkidle' });
  const destination = await page.getByTestId('recording-editor-destination').boundingBox();
  expect(destination?.y).toBeGreaterThan(600);
  await page.getByTestId('group-name').fill('來源編輯中');
  await page.getByRole('button', { name: '1 來源' }).click();
  await page.waitForLoadState('networkidle');
  expect(catalogReads).toBe(0);

  await page.getByRole('button', { name: '3 目的地與確認' }).click();
  await expect(page.getByRole('alert').filter({ hasText: '無法讀取資料表清單' })).toBeVisible();
  expect(catalogReads).toBe(1);
  await page.screenshot({ path: testInfo.outputPath('destination-unavailable.png') });
  await page.getByRole('button', { name: '重新讀取資料表' }).click();
  await expect(page.getByTestId('group-table').locator('option', { hasText: 'gw_group_fixture' })).toHaveCount(1);
  expect(catalogReads).toBe(2);
  await expect(page.getByTestId('group-name')).toHaveValue('來源編輯中');
  await page.screenshot({ path: testInfo.outputPath('destination-after-retry.png') });
});
