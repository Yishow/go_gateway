import { expect, test } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

// Every request is answered here from fixtures; any write or unknown path is recorded and aborted.
const ids = ['D', 'B', 'C', 'E', 'F'];
const workspace = { ...locatorRead().workspace, ordered_device_ids: ids };
const devices = ids.map(id => ({ ...nativeDevice(id), name: `Packaging line ${id} controller`, connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1,"timeout":1500}' }));
const rules = [nativeRule('R', 'D'), nativeRule('S', 'B'), nativeRule('T', 'C'), nativeRule('U', 'E'),
  { ...nativeRule('V1', 'F'), naming_prefix: 'FLOW_A_' }, { ...nativeRule('V2', 'F'), start_address: '40021', naming_prefix: 'FLOW_B_' }];

for (const viewport of [{ width: 320, height: 568 }, { width: 390, height: 844, light: true }]) {
  test(`a narrow frame scrolls its device row to a later device and switches that device's ranges at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width: viewport.width, height: viewport.height });
    const writes: string[] = [], locators: string[] = [];
    await page.route('**/api/**', async route => {
      const request = route.request(), url = new URL(request.url()), path = url.pathname;
      if (request.method() === 'POST' && /^\/api\/v1\/datalink\/devices\/[^/]+\/readiness$/.test(path)) return route.fulfill({ json: { success: true, data: { device_id: 'D', activation_allowed: false } } });
      if (request.method() !== 'GET') { writes.push(`${request.method()} ${path}`); return route.abort(); }
      let data: unknown;
      if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
        const step = url.searchParams.get('step') === '2' ? 2 : 1, deviceId = url.searchParams.get('device_id')!, ruleId = url.searchParams.get('rule_id');
        locators.push(`${step}:${deviceId}${ruleId ? `:${ruleId}` : ''}`);
        const read = locatorRead(step, deviceId); read.workspace = workspace; read.entry.devices = [devices.find(item => item.id === deviceId)! as never];
        const selected = rules.find(rule => rule.id === ruleId && rule.device_id === deviceId);
        if (step === 2 && selected) {
          read.locator.rule_id = selected.id; read.entry.source_rules = [selected as never]; read.entry.rule_heads = { [selected.id]: selected.revision_id! };
          read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [selected.id] };
          read.selection = { ...read.selection, target: { ...(read.selection as { target: object }).target, ruleId: selected.id, ruleRevision: selected.revision_id } } as typeof read.selection;
        }
        data = { ...read.workspace, locator_entry: read.entry };
      } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
      else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
      else if (path === '/api/v1/datalink/studio-v2/recording-home/sources') data = { workspace_id: 'W', state: 'available', source_digest: 'a'.repeat(64),
        devices: devices.map(({ id, name, protocol }) => ({ id, name, protocol, availability: 'unknown' })), candidates: [], excluded: [] };
      else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = devices;
      else if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') data = rules;
      else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
      else return route.abort();
      return route.fulfill({ json: { success: true, data } });
    });
    await page.goto(`/tests/e2e/fixtures/recording-home-source-host.html?home${'light' in viewport ? '&light' : ''}`);
    await page.getByTestId('home-device-points-D').click();
    const dialog = page.getByRole('dialog');
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    const row = dialog.getByRole('group', { name: '設備' });
    await expect(row.getByRole('button')).toHaveCount(ids.length);
    await expect(row.getByTestId('home-source-device-D')).toHaveAttribute('aria-pressed', 'true');
    // The row really overflows sideways and the last device starts outside it.
    expect(await row.evaluate(node => node.scrollWidth > node.clientWidth)).toBe(true);
    const last = row.getByTestId('home-source-device-F');
    expect(await last.evaluate((button, parent) => button.getBoundingClientRect().left >= parent!.getBoundingClientRect().right - 1, await row.elementHandle())).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`row-start-${viewport.width}.png`) });
    // A real horizontal wheel scroll over the row (not scrollIntoView) brings the later device into view.
    const box = (await row.boundingBox())!;
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2);
    await page.mouse.wheel(2000, 0);
    await expect.poll(() => row.evaluate(node => node.scrollLeft)).toBeGreaterThan(0);
    await expect(last).toBeInViewport({ ratio: 0.95 });
    expect(await last.evaluate((button, parent) => button.getBoundingClientRect().right <= parent!.getBoundingClientRect().right + 1, await row.elementHandle())).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`row-scrolled-${viewport.width}.png`) });
    await last.click();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('FLOW_A_');
    await expect(dialog.getByRole('heading')).toContainText('Packaging line F controller');
    await expect(dialog.getByRole('group', { name: '設備' }).getByTestId('home-source-device-F')).toHaveAttribute('aria-pressed', 'true');
    // Then the other range of that same device, inside the same frame.
    await dialog.getByTestId('home-range-option-V2').click();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('FLOW_B_');
    await expect(dialog.getByTestId('home-range-option-V2')).toHaveAttribute('aria-pressed', 'true');
    await page.screenshot({ path: testInfo.outputPath(`switched-${viewport.width}.png`) });
    expect(locators).toEqual(expect.arrayContaining(['1:F', '2:F:V1', '2:F:V2']));
    expect(await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)).toBeLessThanOrEqual(0);
    await expect(page.getByRole('dialog')).toHaveCount(1);
    expect(writes).toEqual([]);
  });
}
