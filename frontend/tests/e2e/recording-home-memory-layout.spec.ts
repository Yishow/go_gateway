import { expect, test, type Locator } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

// Narrow-screen layout of the range settings and memory view; every request is mocked, nothing reaches a device.
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: true } };
const device = { ...nativeDevice(), name: 'Packaging line 3 north PLC', protocol: 'mc_3e', connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' };
const rules = [
  { ...nativeRule('I'), revision_id: 'I1', start_address: 'D0', count: 32, data_type: 'int16', naming_prefix: 'BOILER_OUTLET_TEMPERATURE_SENSOR_' },
  { ...nativeRule('F'), revision_id: 'F1', start_address: 'D39', count: 2, data_type: 'float32', naming_prefix: 'FLOW_' },
] as const;

/** Elements inside the dialog that visibly stick out past its edges (a scrolling row clips its own content). */
async function overflowing(dialog: Locator) {
  return dialog.evaluate(root => {
    const box = root.getBoundingClientRect();
    return [...root.querySelectorAll<HTMLElement>('*')].filter(node => {
      let { left, right } = node.getBoundingClientRect();
      for (let parent = node.parentElement; parent && parent !== root; parent = parent.parentElement) {
        if (getComputedStyle(parent).overflowX === 'visible') continue;
        const clip = parent.getBoundingClientRect(); left = Math.max(left, clip.left); right = Math.min(right, clip.right);
      }
      const rect = node.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0 && right > left && (right > box.right + 0.5 || left < box.left - 0.5);
    }).map(node => `${node.tagName}.${node.className}:${node.textContent?.slice(0, 30)}`);
  });
}
/** On a narrow screen the full title stays unclipped; with two header actions they move below it. */
async function titleRow(dialog: Locator) {
  return dialog.locator('.gw-modal-header').evaluate(header => {
    const title = header.querySelector('h2')!, box = title.getBoundingClientRect();
    return { text: title.textContent, clipped: title.scrollWidth > title.clientWidth + 1,
      beside: [...header.querySelectorAll('button')].filter(button => button.getBoundingClientRect().top < box.bottom).length };
  });
}
/** Rendered lines a phrase of the dialog title spans (a view name split across lines leaves an orphan character). */
async function titlePhraseLines(dialog: Locator, phrase: string) {
  return dialog.locator('.gw-modal-header h2').evaluate((title, wanted) => {
    const walker = document.createTreeWalker(title, NodeFilter.SHOW_TEXT);
    for (let node = walker.nextNode(); node; node = walker.nextNode()) {
      const at = node.textContent?.indexOf(wanted) ?? -1;
      if (at < 0) continue;
      const range = document.createRange(); range.setStart(node, at); range.setEnd(node, at + wanted.length);
      return new Set([...range.getClientRects()].map(rect => Math.round(rect.top))).size;
    }
    return -1;
  }, phrase);
}
/** Header text lines must not collide with header buttons. */
async function headerCollisions(dialog: Locator) {
  return dialog.locator('.gw-modal-header').evaluate(header => {
    const buttons = [...header.querySelectorAll('button')].map(node => node.getBoundingClientRect());
    const texts = [...header.querySelectorAll('h2, p')].map(node => node.getBoundingClientRect());
    return texts.flatMap(text => buttons.filter(button => text.left < button.right && button.left < text.right && text.top < button.bottom && button.top < text.bottom).map(() => 'overlap'));
  });
}

test('range settings and memory view fit a 320px screen with a long device and range name', async ({ page }, testInfo) => {
  await page.setViewportSize({ width: 320, height: 568 });
  await page.route('**/api/**', async route => {
    const request = route.request(), url = new URL(request.url()), path = url.pathname;
    if (request.method() === 'POST' && path.endsWith('/readiness')) return route.fulfill({ json: { success: true, data: { device_id: 'D', activation_allowed: false } } });
    if (request.method() !== 'GET') return route.abort();
    let data: unknown;
    const candidates = /^\/api\/v1\/datalink\/source-rules\/([^/]+)\/candidates$/.exec(path);
    if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
      const step = url.searchParams.get('step') === '2' ? 2 : 1, selected = rules.find(rule => rule.id === url.searchParams.get('rule_id'));
      const read = locatorRead(step); read.workspace = workspace; read.entry.devices = [device as never];
      if (step === 2 && selected) {
        read.locator.rule_id = selected.id; read.entry.source_rules = [selected as never]; read.entry.rule_heads = { [selected.id]: selected.revision_id };
        read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [selected.id] };
      }
      const { modbus_share: _share, ...plain } = workspace; data = { ...plain, locator_entry: read.entry };
    } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
    else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
    else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = [device];
    else if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') data = rules;
    else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
    else if (candidates) {
      const rule = rules.find(item => item.id === candidates[1])!, start = Number(rule.start_address.slice(1)), step = rule.data_type === 'float32' ? 2 : 1;
      data = { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: 'W2', revision_id: rule.revision_id, tags: { status: 'ready', candidates: Array.from({ length: rule.count }, (_, index) => ({
        id: `${rule.id}${index}`, identity: { source_rule_id: rule.id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: rule.start_address }, proposed_signature: `${index}`,
        address: `D${start + index * step}`, point_id: `${rule.id}${index}`, tag_key: `${rule.id}${index}`, display_name: `${rule.id}${index}`, data_type: rule.data_type, status: 'draft' })) } };
    } else return route.abort();
    return route.fulfill({ json: { success: true, data } });
  });
  await page.goto('/tests/e2e/fixtures/recording-home-source-host.html?home&light');
  await page.getByTestId('home-device-points-D').click();
  const dialog = page.getByRole('dialog');
  // No separate list: the first range opens with both ranges as inline choices under the device row.
  await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('BOILER_OUTLET_TEMPERATURE_SENSOR_');
  await expect(dialog.getByTestId('home-range-option-I')).toBeVisible(); await expect(dialog.getByTestId('home-range-option-F')).toContainText('暫存器：D39～D42');
  await expect(dialog.getByRole('switch', { name: '啟用此規則' })).toBeVisible();
  await page.screenshot({ path: testInfo.outputPath('settings-320.png') });
  expect.soft(await overflowing(dialog)).toEqual([]);
  expect.soft(await headerCollisions(dialog)).toEqual([]);
  expect.soft(await titleRow(dialog)).toEqual({ text: 'Packaging line 3 north PLC · 設定讀取點位', clipped: false, beside: 0 });
  expect.soft(await titlePhraseLines(dialog, '設定讀取點位')).toBe(1);
  await expect(dialog.getByTestId('home-range-memory')).toBeInViewport(); await expect(dialog.getByTestId('home-range-done')).toBeInViewport({ ratio: 1 });
  await dialog.getByTestId('home-range-memory').click();
  await expect(dialog.getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'used');
  await page.screenshot({ path: testInfo.outputPath('memory-320.png') });
  expect.soft(await overflowing(dialog)).toEqual([]);
  expect.soft(await headerCollisions(dialog)).toEqual([]);
  expect.soft(await titleRow(dialog)).toMatchObject({ text: 'Packaging line 3 north PLC · 記憶體視圖', clipped: false });
  expect.soft(await titlePhraseLines(dialog, '記憶體視圖')).toBe(1);
  expect.soft(await titlePhraseLines(dialog, 'PLC ·')).toBe(1);
  await expect(dialog.getByText('前往位址', { exact: true })).toBeVisible();
  await expect(dialog.getByText('此區範圍', { exact: false })).toBeInViewport({ ratio: 1 });
  await expect(dialog.getByRole('button', { name: /^D0 · BOILER_OUTLET_TEMPERATURE_SENSOR_$/ })).toBeVisible();
  // The first screen after opening shows memory cells, not only navigation.
  await expect.soft(dialog.getByTestId('home-memory-cell-D0')).toBeInViewport();
  const close = dialog.getByTestId('home-memory-close');
  await expect(close).toBeInViewport();
  await close.click();
  await expect(dialog).toBeHidden();
});
