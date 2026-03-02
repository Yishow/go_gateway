import { expect, test } from '@playwright/test';

async function mockDualEntryEnabled(page: import('@playwright/test').Page) {
  await page.route('**/api/v1/datalink/settings', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        data: [{ key: 'ENABLE_GATEWAY_DUAL_ENTRY', value: true, description: '' }],
      }),
    });
  });
}

test.describe('Gateway Dual Entry E2E', () => {
  test('Quick Setup flow from entry renders payload preview', async ({ page }) => {
    await mockDualEntryEnabled(page);
    await page.goto('/gateway/entry');

    await expect(page).toHaveURL(/\/gateway\/entry/);
    await page.locator('a[href="/gateway/quick-setup"]').first().click();

    await expect(page).toHaveURL(/\/gateway\/quick-setup/);
    await page.getByTestId('host-input').fill('10.10.10.10');
    await page.getByTestId('port-input').fill('1502');

    const preview = page.getByTestId('payload-preview');
    await expect(preview).toContainText('10.10.10.10');
    await expect(preview).toContainText('1502');
  });

  test('Expert flow from entry renders workbench and payload', async ({ page }) => {
    await mockDualEntryEnabled(page);
    await page.goto('/gateway/entry');

    await page.locator('a[href="/gateway/expert-workbench"]').first().click();
    await expect(page).toHaveURL(/\/gateway\/expert-workbench/);

    await expect(page.getByTestId('expert-route-table')).toBeVisible();
    await expect(page.getByTestId('expert-plugin-chain')).toBeVisible();

    const preview = page.getByTestId('expert-payload-preview');
    await expect(preview).toContainText('modbus-tcp');
    await expect(preview).toContainText('routes');
  });

  test('Expert invalid manifest blocks submit', async ({ page }) => {
    await mockDualEntryEnabled(page);
    await page.goto('/gateway/expert-workbench');

    const textarea = page.getByTestId('expert-raw-manifest');
    await textarea.fill('{ invalid }');

    await expect(page.getByTestId('expert-manifest-error')).toBeVisible();
    await expect(page.getByTestId('expert-submit-btn')).toBeDisabled();
  });
});
