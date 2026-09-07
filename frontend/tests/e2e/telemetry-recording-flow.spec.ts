import { expect, test } from '@playwright/test';

const mockRuntimeContext = {
  workspace_id: 'ws-default',
  default_device_id: 'dev-meter-1',
  devices: [
    {
      device_id: 'dev-meter-1',
      name: '三相多功能電表 A',
      protocol: 'modbus_tcp',
      running: true,
      availability_status: 'available',
    },
  ],
  setup: null,
};

const mockRuntimeStatus = {
  running: true,
  uptime_seconds: 120,
  metrics: {
    collected_total: 128,
    write_success_total: 128,
    write_error_total: 0,
    mapping_error_total: 0,
    point_state_error_total: 0,
  },
  collectors: [
    {
      device_id: 'dev-meter-1',
      device_name: '三相多功能電表 A',
      protocol: 'modbus_tcp',
      status: 'running',
      availability_status: 'available',
      availability_reason: null,
      running: true,
      points_total: 1,
      points_healthy: 1,
      points_stale: 0,
      points_error: 0,
      last_read_at: '2026-09-07T10:00:00Z',
      last_error: null,
      breaker_state: 'closed',
    },
  ],
  diagnostics: [],
  database_delivery: [],
  modbus_share_delivery: [],
};

test.describe('Telemetry Recording & History Reports E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Mock runtime context
    await page.route('**/*runtime-context*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: mockRuntimeContext }),
      });
    });

    // Mock runtime status
    await page.route('**/*runtime/status*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: mockRuntimeStatus }),
      });
    });

    // Mock runtime stream SSE
    await page.route('**/*runtime/stream*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'text/event-stream',
        body: ': keepalive\n\n',
      });
    });

    // Mock points query
    await page.route('**/*points*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [] }),
      });
    });
  });

  test('Runtime Dashboard renders history reports panel with empty state when no plan', async ({
    page,
  }) => {
    // Mock no recording plans
    await page.route('**/*recording-plans*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [] }),
      });
    });

    await page.goto('/studio/runtime?device_id=dev-meter-1');
    const emptyPrompt = page.getByTestId('history-reports-empty');
    await expect(emptyPrompt).toBeVisible({ timeout: 15_000 });
    await expect(emptyPrompt).toContainText('尚未建立本設備之記錄方案');
  });

  test('Runtime Dashboard displays summary cards and exports CSV when plan exists', async ({
    page,
  }) => {
    // Mock recording plans
    await page.route('**/*recording-plans*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            {
              id: 'plan-power-1',
              name: '高壓迴路電力記錄方案',
              status: 'running',
              members: [{ member_id: 'm1', equipment_id: 'dev-meter-1' }],
            },
          ],
        }),
      });
    });

    // Mock history query
    await page.route('**/*history/query*', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            workspace_id: 'ws-default',
            plan_id: 'plan-power-1',
            points: [
              {
                measurement_id: 'meas-power-active',
                observed_at: new Date().toISOString(),
                value_numeric: 220.5,
                time_weighted_mean: 220.5,
                sampled_min: 218.0,
                sampled_max: 222.0,
                usage_delta: 5.4,
                quality: 'good',
                coverage_ratio: 0.99,
                is_estimated: false,
                is_provisional: false,
              },
            ],
            total_count: 1,
            generated_at: new Date().toISOString(),
          },
        }),
      });
    });

    // Mock CSV export
    let csvExportTriggered = false;
    await page.route('**/*history/export*', async (route) => {
      csvExportTriggered = true;
      await route.fulfill({
        status: 200,
        contentType: 'text/csv',
        body: 'observed_at,measurement_id,value\n2026-09-07T10:00:00Z,meas-power-active,220.5\n',
      });
    });

    await page.goto('/studio/runtime?device_id=dev-meter-1');

    const panel = page.getByTestId('history-reports-panel');
    await expect(panel).toBeVisible({ timeout: 15_000 });
    await expect(panel).toContainText('高壓迴路電力記錄方案');

    // Verify summary statistics cards
    await expect(page.getByTestId('summary-avg')).toContainText('220.50');
    await expect(page.getByTestId('summary-usage')).toContainText('5.40');
    await expect(page.getByTestId('summary-coverage')).toContainText('99.0%');

    // Click Export CSV button
    const exportBtn = page.getByTestId('history-export-csv-btn');
    await expect(exportBtn).toBeVisible();
    await exportBtn.click();

    await expect.poll(() => csvExportTriggered).toBeTruthy();
  });
});
