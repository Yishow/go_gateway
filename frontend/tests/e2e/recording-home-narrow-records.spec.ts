import { expect, test, type Page } from '@playwright/test';
import { delivery } from '../unit/workbench-v2/recording-home-guidance-fixtures';
import { wide } from '../fixtures/cutoverFixtures';
import { snapshotLifecycle } from './fixtures/recording-home-snapshot-state';

// GET-only isolated mock: every write and unknown path is recorded and aborted, nothing reaches a real backend.
// Groups come from `wide()`, which the real write-group list parser accepts; the unit-only `group()`
// helper is rejected by that parser and leaves the home on its read-error alert.
const DEVICES = [{ id: 'dev-1', name: 'Boiler line PLC' }, { id: 'dev-2', name: 'Chiller plant meter' }];
const groups = [
  wide({ id: 'G1', revision: 'G1r', name: 'Boiler outlet temperature every 30 seconds', status: 'running' }),
  wide({ id: 'G2', revision: 'G2r', applied_revision: '', name: 'Chiller energy draft', status: 'draft' }),
  wide({ id: 'G3', revision: 'G3r', name: 'Line pressure', status: 'ready',
    members: wide().members.map(item => ({ ...item, device_id: 'dev-2' })), basic_managed_device_id: 'dev-2' }),
  wide({ id: 'G4', revision: 'G4r', name: 'Retired flow meter', status: 'disabled',
    destination: { ...wide().destination, storage_strategy: 'custom', table_name: 'legacy_flow_meter_readings_archive' } }),
];
const WS = groups[0].workspace_id;
// Full parser-valid lifecycle DTO from the home snapshot fixture; an incomplete one shows "設定生命週期未能讀取".
const lifecycle = (value: (typeof groups)[number]) => snapshotLifecycle({ workspace: { id: WS }, revision: 'W2' } as Parameters<typeof snapshotLifecycle>[0], value);

async function mockHome(page: Page) {
  const writes: string[] = [], unknown: string[] = [];
  await page.route('**/api/**', async route => {
    const request = route.request(), url = new URL(request.url()), path = url.pathname.replace('/api/v1/datalink', '');
    if (request.method() !== 'GET') { writes.push(`${request.method()} ${path}`); return route.abort(); }
    const match = path.match(/^\/studio-v2\/workspace\/write-groups\/([^/]+)\/(lifecycle|delivery)$/);
    const found = match && groups.find(item => item.id === decodeURIComponent(match[1]));
    let data: unknown;
    if (path === '/studio-v2/workspace') data = { id: WS, kind: 'single', status: 'ready', ordered_device_ids: DEVICES.map(d => d.id), created_at: '', updated_at: '' };
    else if (path === '/studio-v2/workspace/write-groups') data = { workspace_id: WS, workspace_revision: 'W2', groups };
    else if (path === '/studio-v2/recording-home/sources') data = { workspace_id: WS, state: 'available', source_digest: 'a'.repeat(64),
      devices: DEVICES.map((d, index) => ({ ...d, protocol: 'modbus_tcp', availability: index ? 'unknown' : 'available' })), candidates: [], excluded: [] };
    else if (path === '/studio-v2/workspace/database-operations/test-writes') data = { workspace_id: WS, complete: true, operations: [] };
    // Same shapes as recording-home-snapshot-api.ts: no active anomalies, runtime evidence explicitly unavailable.
    else if (path === '/acquisition-anomalies/summary') data = { items: DEVICES.map(d => ({ device_id: d.id, active_count: 0 })), total_active: 0 };
    else if (path === '/studio-v2/workspace/runtime-context' && DEVICES.some(d => d.id === url.searchParams.get('device_id'))) data = { workspace_id: WS,
      devices: [], default_device_id: null, setup: null, rule_evidence: { contract_version: 1, workspace_id: WS, workspace_revision: 'W2',
        device_id: url.searchParams.get('device_id'), observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null } };
    else if (found && match[2] === 'lifecycle') data = lifecycle(found);
    else if (found) data = delivery(found);
    else { unknown.push(path); return route.abort(); }
    return route.fulfill({ json: { success: true, data } });
  });
  return { writes, unknown };
}

async function openHome(page: Page, width: number, height: number) {
  await page.setViewportSize({ width, height });
  const ledger = await mockHome(page);
  await page.goto('/tests/e2e/fixtures/recording-home-source-host.html?home&light');
  const table = page.getByTestId('home-management-table');
  // A failed home read must fail here, not hide behind a partially rendered table.
  await expect(page.getByRole('row').filter({ hasText: groups[0].name }).or(page.getByText('無法讀取目前狀態'))).toBeVisible();
  await expect(page.getByText('無法讀取目前狀態')).toHaveCount(0);
  // Lifecycle, delivery and anomaly reads must all succeed: no read-failure wording anywhere on the home.
  await expect(page.getByText(/未能讀取|無法讀取/)).toHaveCount(0);
  await expect(table.getByRole('row')).toHaveCount(groups.length + 1);
  // Every row's lifecycle and delivery were read and parsed: no row falls back to the incomplete-status text.
  await expect(table.getByText('尚未取得完整狀態')).toHaveCount(0);
  return { ledger, table };
}

const NARROW = [{ width: 390, height: 844 }, { width: 768, height: 1024 }, { width: 820, height: 1180 }];
for (const viewport of NARROW) {
  test(`every recording's facts and its open action are reachable without sideways scrolling at ${viewport.width}px`, async ({ page }, info) => {
    const { ledger, table } = await openHome(page, viewport.width, viewport.height);
    await expect(table.getByRole('columnheader')).toHaveCount(6);
    const panel = page.getByTestId('home-collection-scroll');
    const overflow = await panel.evaluate(node => ({ panel: node.scrollWidth - node.clientWidth, page: document.documentElement.scrollWidth - window.innerWidth }));
    expect(overflow.panel, 'records panel scrolls sideways').toBeLessThanOrEqual(1);
    expect(overflow.page, 'page scrolls sideways').toBeLessThanOrEqual(0);
    const frame = (await panel.boundingBox())!;
    for (const item of groups) {
      const row = table.getByRole('row').filter({ hasText: item.name });
      const cells = row.getByRole('cell');
      await expect(cells).toHaveCount(6);
      // name, status, source, destination, writes, open record: every cell starts and ends inside the visible panel.
      for (let index = 0; index < 6; index += 1) {
        const box = (await cells.nth(index).boundingBox())!;
        expect(box.x, `${item.id} cell ${index} left`).toBeGreaterThanOrEqual(frame.x - 1);
        expect(box.x + box.width, `${item.id} cell ${index} right`).toBeLessThanOrEqual(frame.x + frame.width + 1);
      }
      await expect(row.getByText(item.name)).toBeVisible();
      await expect(row.locator('[data-status-badge]')).toBeVisible();
      await expect(row.locator('[data-flow=source]')).toBeVisible();
      await expect(row.locator('[data-flow=destination]')).toBeVisible();
      const open = row.getByRole('button', { name: '開啟記錄' });
      await expect(open).toBeVisible();
      const box = (await open.boundingBox())!;
      expect(box.height, 'touch target').toBeGreaterThanOrEqual(44);
    }
    await page.screenshot({ path: info.outputPath(`home-records-${viewport.width}.png`), fullPage: true });
    expect(ledger.writes).toEqual([]);
    expect(ledger.unknown).toEqual([]);
  });
}

for (const viewport of NARROW) {
  test(`every device's name, protocol, setting state and actions are reachable without sideways scrolling at ${viewport.width}px`, async ({ page }, info) => {
    const { ledger, table } = await openHome(page, viewport.width, viewport.height);
    await page.getByRole('navigation', { name: '管理記錄或設備' }).getByRole('button', { name: '設備' }).click();
    await expect(table).toHaveAttribute('data-mode', 'devices');
    await expect(table.getByRole('row')).toHaveCount(DEVICES.length + 1);
    await expect(table.getByRole('columnheader')).toHaveCount(4);
    const panel = page.getByTestId('home-collection-scroll');
    const overflow = await panel.evaluate(node => ({ panel: node.scrollWidth - node.clientWidth, page: document.documentElement.scrollWidth - window.innerWidth }));
    expect(overflow.panel, 'devices panel scrolls sideways').toBeLessThanOrEqual(1);
    expect(overflow.page, 'page scrolls sideways').toBeLessThanOrEqual(0);
    const frame = (await panel.boundingBox())!;
    for (const device of DEVICES) {
      const row = table.getByRole('row').filter({ hasText: device.name });
      await expect(row.getByRole('cell')).toHaveCount(4);
      const parts = [row.getByText(device.name), row.locator('[data-protocol]'), row.locator('[data-availability-badge]'),
        row.getByRole('button', { name: '連線設定' }), row.getByRole('button', { name: '讀取點位' }), row.getByTestId(`home-device-anomalies-${device.id}`)];
      for (const part of parts) {
        await expect(part).toBeVisible();
        const box = (await part.boundingBox())!;
        expect(box.x).toBeGreaterThanOrEqual(frame.x - 1);
        expect(box.x + box.width).toBeLessThanOrEqual(frame.x + frame.width + 1);
      }
    }
    await page.screenshot({ path: info.outputPath(`home-devices-${viewport.width}.png`), fullPage: true });
    expect(ledger.writes).toEqual([]);
    expect(ledger.unknown).toEqual([]);
  });
}

test('desktop keeps the six-column table, search, filters, paging and keyboard entry', async ({ page }, info) => {
  const { ledger, table } = await openHome(page, 1440, 900);
  await expect(table.getByRole('columnheader')).toHaveCount(6);
  const header = (await table.getByRole('columnheader').first().boundingBox())!, last = (await table.getByRole('columnheader').last().boundingBox())!;
  expect(Math.abs(header.y - last.y), 'headers stay on one row').toBeLessThanOrEqual(1);
  const panel = page.getByTestId('home-collection-scroll');
  expect(await panel.evaluate(node => node.scrollWidth - node.clientWidth)).toBeLessThanOrEqual(1);
  await page.screenshot({ path: info.outputPath('home-records-1440.png'), fullPage: true });
  const status = page.getByTestId('home-page-status');
  await expect(status).toContainText('共 4 筆');
  const filters = page.getByRole('search');
  await filters.getByRole('searchbox').fill('energy');
  await expect(table.getByRole('row')).toHaveCount(2);
  await filters.getByRole('searchbox').fill('');
  await filters.getByRole('combobox').nth(1).selectOption('disabled');
  await expect(table.getByRole('row')).toHaveCount(2);
  await expect(table.getByRole('row').nth(1)).toContainText('Retired flow meter');
  await filters.getByRole('button', { name: '清除篩選' }).click();
  await expect(table.getByRole('row')).toHaveCount(groups.length + 1);
  const open = table.getByRole('row').filter({ hasText: groups[0].name }).getByRole('button', { name: '開啟記錄' });
  await open.focus();
  await page.keyboard.press('Shift+Tab');
  await page.keyboard.press('Tab');
  await expect(open).toBeFocused();
  await page.getByRole('navigation', { name: '管理記錄或設備' }).getByRole('button', { name: '設備' }).click();
  await expect(table.getByRole('columnheader')).toHaveCount(4);
  expect(await panel.evaluate(node => node.scrollWidth - node.clientWidth)).toBeLessThanOrEqual(1);
  await page.screenshot({ path: info.outputPath('home-devices-1440.png'), fullPage: true });
  expect(ledger.writes).toEqual([]);
  expect(ledger.unknown).toEqual([]);
});
