import { expect, test } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

// Every request is answered here; any write or unknown path is recorded and aborted.
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: true } };
const device = { ...nativeDevice(), name: 'Line 3 PLC', protocol: 'mc_3e', connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' };
const rules = [
  { ...nativeRule('I'), revision_id: 'I1', start_address: 'D0', count: 32, data_type: 'int16', naming_prefix: 'TEMP_' },
  { ...nativeRule('F'), revision_id: 'F1', start_address: 'D39', count: 2, data_type: 'float32', naming_prefix: 'FLOW_' },
  { ...nativeRule('M'), revision_id: 'M1', start_address: 'M0', count: 8, data_type: 'bool', naming_prefix: 'RUN_', enabled: false },
] as const;
const width = (type: string) => type === 'float32' ? 2 : 1;

for (const viewport of [{ name: 'desktop', width: 1440, height: 900 }, { name: 'tablet', width: 768, height: 1024 }, { name: 'phone', width: 390, height: 844 },
  { name: 'desktop-dark', width: 1440, height: 900, dark: true }, { name: 'phone-dark', width: 390, height: 844, dark: true }]) {
  test(`memory view separates MC areas and prefills an exact free address at ${viewport.name}`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width: viewport.width, height: viewport.height });
    const writes: string[] = [];
    await page.route('**/api/**', async route => {
      const request = route.request(), url = new URL(request.url()), path = url.pathname;
      // The original new-range form prepares with saved readiness diagnostics only (no activation, no device I/O).
      if (request.method() === 'POST' && path === '/api/v1/datalink/devices/D/readiness') return route.fulfill({ json: { success: true, data: { device_id: 'D', status: 'not_ready',
        connect_status: 'unknown', probe_status: 'unknown', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] } } });
      if (request.method() !== 'GET') { writes.push(`${request.method()} ${path}`); return route.abort(); }
      let data: unknown;
      const candidates = /^\/api\/v1\/datalink\/source-rules\/([^/]+)\/candidates$/.exec(path);
      if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
        const step = url.searchParams.get('step') === '2' ? 2 : 1, read = locatorRead(step), selected = rules.find(rule => rule.id === url.searchParams.get('rule_id'));
        read.workspace = workspace; read.entry.devices = [device as never];
        if (step === 2 && selected) { read.locator.rule_id = selected.id; read.entry.source_rules = [selected as never]; read.entry.rule_heads = { [selected.id]: selected.revision_id }; read.entry.hydration_scope.rule_ids = [selected.id]; }
        const { modbus_share: _share, ...plain } = workspace;
        data = { ...plain, locator_entry: read.entry };
      } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
      else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
      else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = [device];
      else if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') data = rules;
      else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
      else if (candidates) {
        const rule = rules.find(item => item.id === candidates[1])!, area = rule.start_address[0], start = Number(rule.start_address.slice(1));
        data = { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: url.searchParams.get('expected_workspace_revision'), revision_id: rule.revision_id,
          tags: { status: 'ready', candidates: Array.from({ length: rule.count }, (_, index) => ({ id: `${rule.id}${index}`, identity: { source_rule_id: rule.id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: rule.start_address },
            proposed_signature: `${rule.id}${index}`, address: `${area}${start + index * width(rule.data_type)}`, point_id: `${rule.id}P${index}`, tag_key: `${rule.id}${index}`, display_name: `${rule.id}${index}`, data_type: rule.data_type, status: 'draft' })) } };
      } else return route.abort();
      return route.fulfill({ json: { success: true, data } });
    });
    await page.goto(`/tests/e2e/fixtures/recording-home-source-host.html?home${'dark' in viewport ? '' : '&light'}`);
    await page.getByTestId('home-device-points-D').click();
    const dialog = page.getByRole('dialog');
    // The first range opens directly; every range is an inline choice with its address span.
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('TEMP_');
    await expect(dialog.getByTestId('home-range-option-F')).toContainText('D39 – D42');
    await page.screenshot({ path: testInfo.outputPath(`settings-${viewport.name}.png`) });
    await dialog.getByTestId('home-range-memory').click();
    await expect(dialog.getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'used');
    await expect(dialog.getByTestId('home-memory-cell-D40')).toHaveAttribute('aria-label', 'D40，FLOW_，float32，第 2/2 個暫存器，起始於 D39');
    await expect(dialog.getByTestId('home-memory-cell-D39')).toContainText('float32 1/2');
    await expect(dialog.getByTestId('home-memory-cell-D0')).toContainText('int16');
    // D39 ends a row at 10 and at 5 columns: the value continues on the next row instead of joining across the wrap.
    await expect(dialog.getByTestId('home-memory-cell-D39')).toContainText('續下一列');
    await expect(dialog.getByTestId('home-memory-cell-D40')).toContainText('接續 D39');
    await expect(dialog.getByTestId('home-memory-cell-D41')).toHaveAttribute('data-join-next', 'true');
    await expect(dialog.getByTestId('home-memory-cell-D42')).not.toHaveAttribute('data-join-next', 'true');
    await expect(dialog.getByTestId('home-memory-cell-D33')).toHaveAttribute('data-state', 'empty');
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth);
    expect(overflow).toBeLessThanOrEqual(0);
    await page.screenshot({ path: testInfo.outputPath(`memory-${viewport.name}.png`) });
    await dialog.getByTestId('home-memory-cell-D41').scrollIntoViewIfNeeded();
    await page.screenshot({ path: testInfo.outputPath(`memory-float-${viewport.name}.png`) });
    await dialog.getByRole('tab', { name: /^M/ }).click();
    await expect(dialog.getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'used');
    await dialog.getByRole('tab', { name: /^D/ }).click();
    await dialog.getByTestId('home-memory-cell-D33').focus();
    await page.keyboard.press('ArrowRight');
    await expect(dialog.getByTestId('home-memory-cell-D34')).toBeFocused();
    await dialog.getByTestId('home-memory-cell-D33').click();
    await expect(dialog.getByRole('textbox', { name: '起始位址', exact: true })).toHaveValue('D33');
    const footer = dialog.locator('.gw-modal-footer');
    const [back, save, edge] = await Promise.all([
      footer.getByTestId('recording-range-back').boundingBox(), footer.getByRole('button', { name: '保存範圍' }).boundingBox(), footer.boundingBox(),
    ]);
    expect(back).not.toBeNull(); expect(save).not.toBeNull(); expect(edge).not.toBeNull();
    expect(save!.x - (back!.x + back!.width)).toBeGreaterThanOrEqual(0);
    expect(save!.x - (back!.x + back!.width)).toBeLessThanOrEqual(24);
    expect(edge!.x + edge!.width - (save!.x + save!.width)).toBeLessThanOrEqual(32);
    await page.screenshot({ path: testInfo.outputPath(`prefill-${viewport.name}.png`) });
    await page.keyboard.press('Escape');
    await expect(dialog).toBeHidden();
    expect(writes).toEqual([]);
  });
}
