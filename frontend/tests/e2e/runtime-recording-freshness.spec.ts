import { expect, test } from '@playwright/test';
import {
  RUNTIME_VIEWPORTS, actionCalls, captureRuntime, emitStream, failStream, finishRuntime,
  openRuntime, readyStream, selectDevice,
} from './runtime-recording.helpers';

for (const viewport of RUNTIME_VIEWPORTS) {
  test.describe(`runtime recording freshness ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('RR-F01 10:00 receipt stays stale after10:00:05 disconnect; reconnect alone cannot claim LIVE', async ({ page }, info) => {
      const journey = await openRuntime(page);
      const receipt = page.getByTestId('runtime-snapshot-received-at');
      const banner = page.getByTestId('runtime-dashboard-live-state-banner');
      const route = page.getByTestId('runtime-dashboard-route-state');
      await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
      const stream = await readyStream(page); await expect(route).toHaveText('live');
      await expect(banner).toContainText('does not prove every point is fresh');
      const nextSnapshot = journey.holdNext('status');
      await page.clock.setFixedTime(new Date('2026-10-29T10:00:05Z')); await failStream(page, stream);
      await expect(route).toHaveText('degraded'); await expect(banner).toContainText('Stale snapshot');
      await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
      await expect(page.getByTestId('runtime-dashboard-summary-panel')).toBeVisible();
      const logs = page.getByTestId('runtime-dashboard-logs-panel');
      await expect(logs.getByRole('heading', { name: 'Connection events', exact: true })).toBeVisible();
      await expect(logs).toContainText('not SQL delivery receipts or process logs');
      await expect(logs.getByText('LIVE', { exact: true })).toHaveCount(0);
      await page.getByRole('button', { name: 'Reconnect live stream', exact: true }).click();
      const reconnected = await readyStream(page); expect(reconnected).toBeGreaterThan(stream);
      await expect(route).toHaveText('degraded'); await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
      // Hold the next existing polling read rather than cancelling it with a second manual refetch.
      await nextSnapshot.reached;
      await expect(route).toHaveText('degraded');
      await captureRuntime(page, info, 'rr-f01-reconnected-but-retained-stale', banner);
      await page.clock.setFixedTime(new Date('2026-10-29T10:00:06Z')); nextSnapshot.release();
      await expect(route).toHaveText('live'); await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:06.000Z');
      await expect(banner).not.toContainText('Stale snapshot'); await finishRuntime(info, journey);
    });

    test('RR-F02 failed snapshot preserves old receipt and existing polling recovers without a write', async ({ page }, info) => {
      const journey = await openRuntime(page); const receipt = page.getByTestId('runtime-snapshot-received-at');
      await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
      const stream = await readyStream(page); await expect(page.getByTestId('runtime-dashboard-route-state')).toHaveText('live');
      await page.clock.setFixedTime(new Date('2026-10-29T10:00:05Z')); await failStream(page, stream);
      await page.getByRole('button', { name: 'Reconnect live stream', exact: true }).click(); await readyStream(page);
      journey.mock.failNext('status', 503);
      const before = actionCalls(journey.mock, 'status').length;
      await page.getByRole('button', { name: 'Retry snapshot', exact: true }).click();
      await expect.poll(() => actionCalls(journey.mock, 'status').length).toBeGreaterThan(before);
      await expect(page.getByTestId('runtime-dashboard-live-state-banner')).toContainText('Stale snapshot');
      await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:00.000Z');
      await page.clock.setFixedTime(new Date('2026-10-29T10:00:07Z'));
      // Real production 5-second degraded polling runs; no custom timer or state injection.
      await expect.poll(() => actionCalls(journey.mock, 'status').length, { timeout: 12_000 }).toBeGreaterThan(before + 1);
      await expect(page.getByTestId('runtime-dashboard-route-state')).toHaveText('live');
      await expect(receipt).toHaveAttribute('datetime', '2026-10-29T10:00:07.000Z');
      await captureRuntime(page, info, 'rr-f02-polling-recovery', page.getByTestId('runtime-dashboard-live-state-banner'));
      await finishRuntime(info, journey, [{ action: 'status', status: 503 }]);
    });

    test('RR-F03 delayed A values/state on active B stream cannot populate or degrade B', async ({ page }, info) => {
      const journey = await openRuntime(page); await readyStream(page);
      await selectDevice(page, 'B'); const stream = await readyStream(page, 'B');
      const points = page.getByTestId('runtime-dashboard-live-points-table');
      const value = { device_id: 'B', point_id: 'PB1', address: '41001', raw_value: 123, transformed_value: 14.3,
        quality: 'good', stale: false, timestamp: '2026-10-29T10:00:00Z' };
      await emitStream(page, stream, 'value', value); await expect(points.getByRole('cell', { name: '123', exact: true })).toBeVisible();
      // Deliberately use the current OPEN B instance. A closed-instance fixture would swallow the event and prove nothing.
      await emitStream(page, stream, 'value', { ...value, device_id: 'A', point_id: 'PA1', address: '40001', raw_value: 999999, transformed_value: 999999 });
      await emitStream(page, stream, 'value', { ...value, device_id: 'A', raw_value: 888888, transformed_value: 888888 });
      await emitStream(page, stream, 'stream_state', { device_id: 'A', timestamp: value.timestamp,
        stream_state: { state: 'unavailable', empty: false, stale: true, degraded: true, unavailable: true } });
      // A subsequent valid event crosses the real 200ms display batch, making the negative DOM check observable.
      await emitStream(page, stream, 'value', { ...value, raw_value: 124 });
      await expect(points.getByRole('cell', { name: '124', exact: true })).toBeVisible();
      await expect(points).not.toContainText(/999999|888888|40001/);
      await expect(page.getByTestId('runtime-dashboard-route-state')).toHaveText('live');
      await expect(page.getByTestId('runtime-dashboard-header').getByRole('heading', { name: 'Device B', exact: true })).toBeVisible();
      await captureRuntime(page, info, 'rr-f03-active-b-rejects-late-a', points); await finishRuntime(info, journey);
    });
  });
}
