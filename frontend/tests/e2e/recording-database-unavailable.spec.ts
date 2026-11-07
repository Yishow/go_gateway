import { expect, test } from '@playwright/test';
import { quietDelivery } from '../fixtures/cutoverFixtures';
import { captureGuided, openGuided, openGuidedAdvanced } from './recording-guided.helpers';

// Real components with an isolated API fixture; this does not claim a real database outage.
test('shows and clears the database destination warning as delivery evidence changes', async ({ page }, info) => {
  let unavailable = true;
  await openGuided(page, 'saved-wide', {
    language: 'zh-TW',
    setup: mock => {
      const handle = mock.handle;
      mock.handle = request => {
        const response = handle(request);
        if (request.method !== 'GET' || !new URL(request.url).pathname.endsWith('/write-groups/G1/delivery') || response.status !== 200 || !unavailable) return response;
        return { ...response, body: { success: true, data: {
          ...quietDelivery, group_id: 'G1', stages: { ...quietDelivery.stages, retrying: 1 },
          backlog: [{ group_revision: 'R1', connector_id: 'conn-1', connector_revision: 'crev-1',
            table_schema: 'main', table_name: 'readings', pending: 1, error_codes: ['target-unavailable'] }],
        } } };
      };
    },
  });
  const surface = await openGuidedAdvanced(page, 'G1');
  const warning = surface.getByTestId('group-destination-unavailable');
  await expect(warning).toContainText('資料庫目的地曾無法使用');
  await expect(surface.getByTestId('delivery-local')).toHaveText('1');
  await expect(surface.getByTestId('delivery-committed')).toHaveText('0');
  await captureGuided(page, info, 'database-destination-unavailable', warning);
  const panelPath = info.outputPath('database-warning-panel.png');
  await surface.getByTestId('group-delivery-strip').screenshot({ path: panelPath });
  await info.attach('database warning panel — mock API/test data', { path: panelPath, contentType: 'image/png' });
  unavailable = false;
  await expect(warning).toHaveCount(0, { timeout: 15_000 });
});
