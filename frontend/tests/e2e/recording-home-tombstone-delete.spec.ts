import { expect, test } from '@playwright/test';
import { wide } from '../fixtures/cutoverFixtures';
import { snapshotLifecycle } from './fixtures/recording-home-snapshot-state';

const group = wide({ id: 'G1', name: 'Test', status: 'deleted', revision: 'G1r' });
const workspaceId = group.workspace_id;
const state = { workspace: { id: workspaceId }, revision: 'W2' } as Parameters<typeof snapshotLifecycle>[0];

for (const width of [390, 1440]) {
  test(`unconfirmed tombstone opens the guarded delete dialog at ${width}px`, async ({ page }, info) => {
    const writes: string[] = [];
    await page.setViewportSize({ width, height: 900 });
    await page.route('**/api/**', async route => {
      const request = route.request();
      const url = new URL(request.url());
      const path = url.pathname.replace('/api/v1/datalink', '');
      if (request.method() !== 'GET') { writes.push(`${request.method()} ${path}`); return route.abort(); }
      let data: unknown;
      if (path === '/studio-v2/workspace') data = { id: workspaceId, kind: 'single', status: 'ready', ordered_device_ids: ['dev-1'], created_at: '', updated_at: '' };
      else if (path === '/studio-v2/workspace/write-groups') data = { workspace_id: workspaceId, workspace_revision: 'W2', groups: [group] };
      else if (path === '/studio-v2/recording-home/sources') data = { workspace_id: workspaceId, state: 'available', source_digest: 'a'.repeat(64),
        devices: [{ id: 'dev-1', name: 'MC', protocol: 'modbus_tcp', availability: 'available' }], candidates: [], excluded: [] };
      else if (path === '/studio-v2/workspace/database-operations/test-writes') data = { workspace_id: workspaceId, complete: true, operations: [] };
      else if (path === '/runtime/status' && !url.searchParams.has('device_id')) data = { running: false, uptime_seconds: 0, collectors: [],
        snapshot_state: { state: 'ready', empty: true, degraded: false, unavailable: false, stale: false } };
      else if (path === '/acquisition-anomalies/summary') data = { items: [], total_active: 0 };
      else if (path === '/studio-v2/workspace/runtime-context') data = { workspace_id: workspaceId, devices: [], default_device_id: null, setup: null,
        rule_evidence: { contract_version: 1, workspace_id: workspaceId, workspace_revision: 'W2', device_id: 'dev-1', observed_at: null,
          freshness: 'unavailable', source: null, groups: [], locator: null } };
      else if (path === '/studio-v2/workspace/write-groups/G1/lifecycle') {
        const facts = snapshotLifecycle(state, group);
        data = { ...facts, actions: facts.actions.map(action => action.action === 'remove'
          ? { ...action, allowed: true, reason_code: '', next_action: '' } : action) };
      } else if (path === '/studio-v2/workspace/write-groups/G1/delivery') {
        return route.fulfill({ status: 503, json: { success: false, error: { code: 'UNAVAILABLE' } } });
      } else return route.abort();
      return route.fulfill({ json: { success: true, data } });
    });
    await page.goto('/tests/e2e/fixtures/recording-home-source-host.html?home&light');
    const unconfirmed = page.getByTestId('recording-home-unconfirmed');
    await expect(unconfirmed).toBeVisible();
    await unconfirmed.locator('summary').first().click();
    await unconfirmed.getByRole('button', { name: '「Test」的操作' }).click();
    const remove = unconfirmed.getByRole('menuitem', { name: '刪除記錄' });
    await expect(remove).toBeEnabled();
    const bounds = await remove.boundingBox();
    expect(bounds).not.toBeNull();
    expect(bounds!.x).toBeGreaterThanOrEqual(0);
    expect(bounds!.x + bounds!.width).toBeLessThanOrEqual(width);
    expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBe(width);
    await page.screenshot({ path: info.outputPath(`tombstone-menu-${width}.png`), fullPage: true });
    await remove.click();
    const dialog = page.getByRole('alertdialog');
    await expect(dialog).toContainText('Test');
    await expect(page.getByTestId('group-delete-confirm-button')).toBeDisabled();
    await page.screenshot({ path: info.outputPath(`tombstone-dialog-${width}.png`) });
    expect(writes).toEqual([]);
  });
}
