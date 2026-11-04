import { expect, test } from '@playwright/test';

const catalog = { workspace_id: 'ws-1', workspace_revision: 'workspace-1', connector_id: 'conn-1', connector_revision: 'crev-1',
  kind: 'postgres', default_schema: 'public', schemas: [{ name: 'public', tables: ['customer_readings', 'gw_group_fixture', 'gw_other_owner'] }, { name: 'archive', tables: ['history'] }, { name: 'empty_schema', tables: [] }] };

for (const viewport of [{ width: 1360, height: 900 }, { width: 390, height: 844 }]) {
  test(`catalog choices at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize(viewport);
    let requested = false;
    await page.route('**/api/**', route => {
      const url = new URL(route.request().url());
      if (!url.pathname.endsWith('/database-catalog')) return route.fulfill({ status: 503, json: { success: false } });
      expect(Object.fromEntries(url.searchParams)).toEqual({ workspace_id: 'ws-1', expected_workspace_revision: 'workspace-1', connector_id: 'conn-1',
        expected_connector_revision: 'crev-1', group_id: 'group-1', expected_group_revision: 'rev-1' });
      requested = true;
      return route.fulfill({ json: { success: true, data: catalog } });
    });
    await page.goto('/tests/e2e/fixtures/recording-catalog-host.html');
    const destinationStep = page.getByRole('button', { name: '3 目的地與確認' });
    await destinationStep.click();
    await expect(destinationStep).toHaveAttribute('aria-current', 'step');
    const table = page.getByTestId('group-table'), schema = page.getByTestId('group-schema');
    await expect(table.locator('option', { hasText: 'gw_group_fixture' })).toHaveCount(1);
    expect(requested).toBe(true);
    await expect(table.locator('option', { hasText: 'customer_readings' })).toHaveCount(1);
    await expect(table.locator('option', { hasText: 'gw_other_owner' })).toHaveAttribute('disabled', '');
    await expect(table.locator('option', { hasText: 'pg_catalog' })).toHaveCount(0);
    await expect(schema.locator('option', { hasText: 'empty_schema' })).toHaveCount(1);
    const topEdgeIsTable = await table.evaluate(element => {
      const control = element.getBoundingClientRect();
      const help = document.querySelector('[data-parameter-help="record.table"]')!.getBoundingClientRect();
      return document.elementFromPoint((help.left + help.right) / 2, control.top + 2) === element;
    });
    expect(topEdgeIsTable).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`destination-viewport-${viewport.width}.png`) });
    await page.getByTestId('recording-editor-destination').screenshot({ path: testInfo.outputPath(`destination-${viewport.width}.png`) });
    await schema.selectOption({ label: 'archive' });
    const manualTable = page.getByTestId('group-table-manual');
    await expect(manualTable).toHaveValue('gw_group_fixture');
    const tableBounds = await table.boundingBox(), manualBounds = await manualTable.boundingBox();
    expect(tableBounds).not.toBeNull(); expect(manualBounds).not.toBeNull();
    expect(manualBounds!.y).toBeGreaterThanOrEqual(tableBounds!.y + tableBounds!.height);
    await table.selectOption({ label: 'history' });
    await expect(table.locator('option:checked')).toHaveText('history');
    await schema.selectOption({ label: 'empty_schema' });
    const hints = table.locator('..').locator('span');
    await expect(hints).toHaveCount(2);
    const firstHint = await hints.nth(0).boundingBox(), secondHint = await hints.nth(1).boundingBox();
    expect(firstHint).not.toBeNull(); expect(secondHint).not.toBeNull();
    expect(secondHint!.y).toBeGreaterThanOrEqual(firstHint!.y + firstHint!.height);
  });
}

test('failed and stale reads retain the draft and can be retried', async ({ page }) => {
  let attempts = 0;
  await page.route('**/api/**', route => {
    if (!new URL(route.request().url()).pathname.endsWith('/database-catalog')) return route.fulfill({ status: 503, json: { success: false } });
    attempts += 1;
    if (attempts === 1) return route.fulfill({ json: { success: true, data: { ...catalog, connector_revision: 'old' } } });
    return route.fulfill({ json: { success: true, data: catalog } });
  });
  await page.goto('/tests/e2e/fixtures/recording-catalog-host.html');
  await page.getByRole('button', { name: '3 目的地與確認' }).click();
  const table = page.getByTestId('group-table');
  await expect(page.getByRole('alert').filter({ hasText: '無法讀取資料表清單' })).toBeVisible();
  await table.fill('my_draft');
  await page.getByRole('button', { name: '重新讀取資料表' }).click();
  await expect(page.getByTestId('group-table-manual')).toHaveValue('my_draft');
  expect(attempts).toBe(2);
});

test('manual entry cannot bypass a known wide-table owner', async ({ page }) => {
  await page.route('**/api/**', route => new URL(route.request().url()).pathname.endsWith('/database-catalog')
    ? route.fulfill({ json: { success: true, data: catalog } }) : route.fulfill({ status: 503, json: { success: false } }));
  await page.goto('/tests/e2e/fixtures/recording-catalog-host.html');
  await page.getByRole('button', { name: '3 目的地與確認' }).click();
  const table = page.getByTestId('group-table');
  await expect(table.locator('option', { hasText: 'gw_other_owner' })).toHaveAttribute('disabled', '');
  await table.selectOption({ label: '自行輸入資料表名稱' });
  await page.getByTestId('group-table-manual').fill('gw_other_owner');
  await expect(page.getByRole('alert').filter({ hasText: '此寬表由原記錄擁有' })).toBeVisible();
  await expect(page.getByTestId('group-save')).toBeDisabled();
});
