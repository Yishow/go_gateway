import { mkdirSync } from 'node:fs';
import { join } from 'node:path';
import { expect, test, type Page } from '@playwright/test';
import { delivery } from '../unit/workbench-v2/recording-home-guidance-fixtures';
import { wide } from '../fixtures/cutoverFixtures';
import { snapshotLifecycle } from './fixtures/recording-home-snapshot-state';

// Offline, GET-only fixture. This is production /studio/v2 UI with synthetic observations.
const NOW = '2026-11-08T00:00:00.000Z';
const DEVICES = [
  { id: 'dev-1', name: 'Boiler line PLC', state: 'collecting' },
  { id: 'dev-2', name: 'Chiller plant meter', state: 'degraded' },
  { id: 'dev-3', name: 'Idle pump controller', state: 'idle' },
] as const;
const groups = [
  wide({ id: 'G1', revision: 'G1r', applied_revision: 'G1r', name: 'Boiler outlet temperature', status: 'running' }),
  wide({ id: 'G2', revision: 'G2r', applied_revision: 'G2r', name: 'Chiller energy', status: 'ready',
    members: wide().members.map(item => ({ ...item, device_id: 'dev-2' })), basic_managed_device_id: 'dev-2' }),
  wide({ id: 'G3', revision: 'G3r', applied_revision: '', name: 'Pump pressure draft', status: 'draft',
    members: wide().members.map(item => ({ ...item, device_id: 'dev-3' })), basic_managed_device_id: 'dev-3' }),
];
const WS = groups[0].workspace_id;
type Theme = 'light' | 'dark';

async function mockHome(page: Page, theme: Theme, baseURL: string) {
  const unexpected: string[] = [];
  const origin = new URL(baseURL).origin;
  await page.addInitScript(selectedTheme => {
    localStorage.setItem('gateway.theme.preference.v1:/api/v1/datalink', JSON.stringify({ version: 1, preference: selectedTheme }));
    localStorage.setItem('i18nextLng', 'zh-TW');
    Reflect.set(window, '__REACT_GRAB_DISABLED__', true);
  }, theme);
  await page.clock.setFixedTime(new Date(NOW));
  await page.route('**/*', async route => {
    const request = route.request(), url = new URL(request.url()), path = url.pathname.replace('/api/v1/datalink', '');
    if (url.origin !== origin) { unexpected.push(`${request.method()} ${url.origin}${url.pathname}`); return route.abort(); }
    if (!url.pathname.startsWith('/api/')) {
      if (request.method() === 'GET') return route.continue();
      unexpected.push(`${request.method()} ${url.pathname}`); return route.abort();
    }
    if (request.method() !== 'GET') { unexpected.push(`${request.method()} ${url.pathname}`); return route.abort(); }
    const match = path.match(/^\/studio-v2\/workspace\/write-groups\/([^/]+)\/(lifecycle|delivery)$/);
    const found = match && groups.find(item => item.id === decodeURIComponent(match[1]));
    let data: unknown;
    if (path === '/settings') data = [{ key: 'theme', value: theme }];
    else if (path === '/studio-v2/workspace') data = { id: WS, kind: 'single', status: 'ready', ordered_device_ids: DEVICES.map(d => d.id), created_at: '', updated_at: '' };
    else if (path === '/studio-v2/workspace/write-groups') data = { workspace_id: WS, workspace_revision: 'W2', groups };
    else if (path === '/studio-v2/recording-home/sources') data = { workspace_id: WS, state: 'available', source_digest: 'a'.repeat(64),
      devices: DEVICES.map(({ id, name }) => ({ id, name, protocol: 'modbus_tcp', availability: 'available' })), candidates: [], excluded: [] };
    else if (path === '/studio-v2/workspace/database-operations/test-writes') data = { workspace_id: WS, complete: true, operations: [] };
    else if (path === '/runtime/status' && !url.searchParams.has('device_id')) data = {
      running: true, uptime_seconds: 30, metrics: { write_success_total: 10, write_error_total: 1 },
      snapshot_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false },
      collectors: DEVICES.map(device => ({ device_id: device.id, device_name: device.name, protocol: 'modbus_tcp',
        status: device.state === 'idle' ? 'idle' : 'running', availability_status: 'available', running: device.state !== 'idle',
        points_total: device.state === 'idle' ? 0 : 4, points_healthy: device.state === 'collecting' ? 4 : device.state === 'degraded' ? 3 : 0,
        points_stale: 0, points_error: device.state === 'degraded' ? 1 : 0, last_read_at: device.state === 'idle' ? null : NOW,
        last_error: null, breaker_state: 'closed', effective_interval_ms: 1000 })),
    };
    else if (path === '/acquisition-anomalies/summary') data = { items: DEVICES.map(d => ({ device_id: d.id, active_count: 0 })), total_active: 0 };
    else if (path === '/studio-v2/workspace/runtime-context' && DEVICES.some(d => d.id === url.searchParams.get('device_id'))) data = { workspace_id: WS,
      devices: [], default_device_id: null, setup: null, rule_evidence: { contract_version: 1, workspace_id: WS, workspace_revision: 'W2',
        device_id: url.searchParams.get('device_id'), observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null } };
    else if (found && match[2] === 'lifecycle') data = snapshotLifecycle({ workspace: { id: WS }, revision: 'W2' } as Parameters<typeof snapshotLifecycle>[0], found);
    else if (found) data = { ...delivery(found), ...(found.id === 'G1' ? { last_sql_committed_effect: {
      group_revision: found.applied_revision, connector_revision: found.destination.connector_revision,
      record_id: 'receipt-1', effect_key: 'effect-1', payload_digest: 'a'.repeat(64), committed_at: NOW,
    } } : {}) };
    else { unexpected.push(`${request.method()} ${url.pathname}`); return route.abort(); }
    return route.fulfill({ json: { success: true, data } });
  });
  return unexpected;
}

const viewports = [{ width: 1440, height: 900 }, { width: 1024, height: 768 }, { width: 768, height: 1024 }, { width: 390, height: 844 }];
for (const viewport of viewports) for (const theme of ['light', 'dark'] as const) {
  test(`recording home status layout at ${viewport.width}x${viewport.height} in ${theme}`, async ({ page, baseURL }, info) => {
    if (!baseURL) throw new Error('A local baseURL is required');
    await page.setViewportSize(viewport);
    const unexpected = await mockHome(page, theme, baseURL);
    await page.goto('/studio/v2?view=recordings');
    await expect(page.locator('html')).toHaveClass(new RegExp(`(^|\\s)${theme}(\\s|$)`));
    await expect(page.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false');
    await expect(page.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh');
    const rows = page.locator('table[data-mode="recordings"] tbody tr');
    await expect(rows).toHaveCount(groups.length);
    await expect(rows.locator('details')).toHaveCount(0);
    for (const device of DEVICES) await expect(page.getByTestId(`device-status-${device.id}`).locator('[data-connection-state]')).toHaveAttribute('data-connection-state', device.state);
    await expect(rows.first().locator('[data-write-stage]')).toHaveCount(4);
    await expect(rows.first().locator('[data-write-stage="configuration"]')).toHaveAttribute('data-stage-state', 'complete');
    await expect(rows.first().locator('[data-write-stage="sql"]')).toHaveAttribute('data-stage-state', 'complete');
    expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBe(await page.evaluate(() => window.innerWidth));
    const shots = process.env.RECORDING_HOME_LAYOUT_SHOTS ?? info.outputDir;
    mkdirSync(shots, { recursive: true });
    await page.screenshot({ path: join(shots, `home-${viewport.width}-${theme}.png`) });
    // Keep a full-page companion: at short viewports the record table begins below the fold.
    await page.screenshot({ path: join(shots, `home-${viewport.width}-${theme}-full.png`), fullPage: true });
    const strip = page.getByRole('region', { name: '設備連線狀態' });
    for (const device of DEVICES) {
      const card = page.getByTestId(`device-status-${device.id}`);
      await card.scrollIntoViewIfNeeded();
      const frame = (await strip.boundingBox())!, box = (await card.boundingBox())!;
      expect(box.x).toBeGreaterThanOrEqual(frame.x - 1);
      expect(box.x + box.width).toBeLessThanOrEqual(frame.x + frame.width + 1);
      const action = strip.getByRole('button', { name: `刪除設備 ${device.name}` });
      await action.scrollIntoViewIfNeeded();
      const actionBox = (await action.boundingBox())!;
      expect(actionBox.x).toBeGreaterThanOrEqual(frame.x - 1);
      expect(actionBox.x + actionBox.width).toBeLessThanOrEqual(frame.x + frame.width + 1);
      expect(actionBox.y + actionBox.height).toBeLessThanOrEqual(frame.y + frame.height + 1);
    }
    for (let index = 0; index < groups.length; index += 1) {
      const open = rows.nth(index).getByRole('button', { name: '開啟記錄' });
      await open.scrollIntoViewIfNeeded();
      await expect(open).toBeInViewport({ ratio: 1 });
    }
    expect(unexpected, 'API fixture refuses unlisted reads and every write').toEqual([]);
  });
}
