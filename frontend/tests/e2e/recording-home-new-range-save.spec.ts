import { expect, test, type Page } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

// Regression: a new range opened from the device settings frame (directly or from a free
// memory cell) must finish preparing under React StrictMode and submit exactly one create.
// Every request is simulated: the rule create gets a disposable ACK, every other write aborts.
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: true } };
const device = { ...nativeDevice(), name: 'Line 1 Modbus', protocol: 'modbus_tcp', connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1}' };
const rules = [{ ...nativeRule('I'), revision_id: 'I1', start_address: '40001', count: 16, data_type: 'int16', naming_prefix: 'TEMP_' }] as const;
const createPath = '/api/v1/datalink/studio-v2/workspace/source-rules';
type Seen = { creates: Record<string, unknown>[]; writes: string[] };

async function simulate(page: Page, readiness: 'ok' | 'fail', seen: Seen) {
  await page.route('**/api/**', async route => {
    const request = route.request(), url = new URL(request.url()), path = url.pathname;
    if (request.method() === 'POST' && path === '/api/v1/datalink/devices/D/readiness') {
      if (readiness === 'fail') return route.fulfill({ status: 503, json: { success: false, error: 'unavailable' } });
      return route.fulfill({ json: { success: true, data: { device_id: 'D', status: 'not_ready', connect_status: 'unknown', probe_status: 'unknown',
        activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] } } });
    }
    if (request.method() === 'POST' && path === createPath) {
      const body = request.postDataJSON() as Record<string, unknown>; seen.creates.push(body);
      return route.fulfill({ json: { success: true, data: { ...body, workspace_id: 'W', revision_id: 'N1', locked: false, origin: 'manual', created_at: '', updated_at: '', runtime_apply_status: 'not_running' } } });
    }
    if (request.method() !== 'GET') { seen.writes.push(`${request.method()} ${path}`); return route.abort(); }
    let data: unknown;
    const candidates = /^\/api\/v1\/datalink\/source-rules\/([^/]+)\/candidates$/.exec(path);
    if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
      // The home frame hydrates through the real locator entry (step 1 device, step 2 rule head + scope).
      const step = url.searchParams.get('step') === '2' ? 2 : 1, read = locatorRead(step), selected = rules.find(rule => rule.id === url.searchParams.get('rule_id'));
      read.workspace = workspace; read.entry.devices = [device as never];
      if (step === 2 && selected) { read.locator.rule_id = selected.id; read.entry.source_rules = [selected as never]; read.entry.rule_heads = { [selected.id]: selected.revision_id }; read.entry.hydration_scope.rule_ids = [selected.id]; }
      const { modbus_share: _share, ...plain } = workspace;
      data = { ...plain, locator_entry: read.entry };
    } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
    else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
    else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = [device];
    else if (path === createPath) data = rules;
    else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
    else if (candidates && candidates[1] === 'I') {
      const rule = rules[0];
      data = { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: url.searchParams.get('expected_workspace_revision'), revision_id: rule.revision_id,
        tags: { status: 'ready', candidates: Array.from({ length: rule.count }, (_, index) => ({ id: `I${index}`, identity: { source_rule_id: rule.id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: rule.start_address },
          proposed_signature: `I${index}`, address: String(40001 + index), point_id: `IP${index}`, tag_key: `I${index}`, display_name: `I${index}`, data_type: 'int16', status: 'draft' })) } };
    } else return route.abort();
    return route.fulfill({ json: { success: true, data } });
  });
}

async function open(page: Page, entry: 'direct' | 'memory') {
  await page.goto('/tests/e2e/fixtures/recording-home-source-host.html?home&strict&light');
  await page.getByTestId('home-device-points-D').click();
  const dialog = page.getByRole('dialog');
  await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('TEMP_');
  if (entry === 'direct') {
    await dialog.getByTestId('home-range-option-new').click();
    await dialog.getByRole('textbox', { name: '起始位址', exact: true }).fill('40017');
  } else {
    await dialog.getByTestId('home-range-memory').click();
    await dialog.getByTestId('home-memory-cell-40017').click();
    await expect(dialog.getByRole('textbox', { name: '起始位址', exact: true })).toHaveValue('40017');
  }
  await dialog.getByRole('textbox', { name: '資料項目數', exact: true }).fill('8');
  await dialog.getByRole('combobox', { name: '資料型別', exact: true }).selectOption('float32');
  await dialog.getByRole('textbox', { name: '點位名稱前綴', exact: true }).fill('FLOW_');
  return dialog;
}

for (const entry of ['direct', 'memory'] as const) {
  test(`new range from ${entry} entry prepares under StrictMode and submits one create`, async ({ page }) => {
    const seen: Seen = { creates: [], writes: [] };
    await simulate(page, 'ok', seen);
    const dialog = await open(page, entry);
    // Preparation is transient: it must resolve to the saved-readiness note, not hang.
    await expect(dialog.getByText('正在核對目前工作區…')).toHaveCount(0);
    await expect(dialog.getByText('設備尚未就緒，將先保存離線來源；不會測試連線或開始採集。')).toBeVisible();
    await expect(dialog.getByText('8 筆資料 · 占用 16 個位址 · 結束於 40032')).toBeVisible();
    const save = dialog.getByRole('button', { name: '保存範圍' });
    await expect(save).toBeEnabled();
    await save.dblclick();
    await expect.poll(() => seen.creates.length).toBe(1);
    expect(seen.creates[0]).toMatchObject({ device_id: 'D', start_address: '40017', count: 8, data_type: 'float32', naming_prefix: 'FLOW_', enabled: false });
    // The submitted buffer leaves the editable state; a second click cannot resubmit it.
    await expect(save).toBeDisabled();
    await page.waitForTimeout(500);
    expect(seen.creates).toHaveLength(1);
  });
}

test('readiness 503 blocks the new range instead of checking forever and never writes', async ({ page }) => {
  const seen: Seen = { creates: [], writes: [] };
  await simulate(page, 'fail', seen);
  const dialog = await open(page, 'direct');
  await expect(dialog.getByText('正在核對目前工作區…')).toHaveCount(0);
  await expect(dialog.getByText('無法核對目前工作區或版本。記錄草稿未變更；請先重新讀取目前範圍再試。')).toBeVisible();
  await expect(dialog.getByRole('button', { name: '保存範圍' })).toBeDisabled();
  expect(seen.creates).toEqual([]);
  expect(seen.writes).toEqual([]);
});
