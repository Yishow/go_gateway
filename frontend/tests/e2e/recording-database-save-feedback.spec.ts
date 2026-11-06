import { expect, test as base, type Page, type TestInfo } from '@playwright/test';
import { DatabaseSaveFeedbackApi, databaseSavePaths, readKinds, type AckMode, type ReadKind } from './fixtures/database-save-feedback-api';

// AUTHORING ONLY in cloud. Execute in the joint Mac browser pass (C 3.1).
// Real /studio/v2 App/router with simulated HTTP, not a direct Step4 fixture host.
// Config autosave has no standalone navigation guard. Navigation cases explicitly
// exercise the EXISTING recording-draft guard; they do not claim database-only protection.
const advanced = '/studio/v2?view=advanced&step=output&target=database#save-feedback';
const home = '/studio/v2?view=recordings#save-feedback';
const saved = (page: Page) => page.getByTestId('database-connector-save-status');
const refreshing = (page: Page) => page.getByTestId('database-connector-refresh-status');
const name = (page: Page) => page.getByRole('textbox', { name: '連線名稱', exact: true });
const test = base.extend<{ api: DatabaseSaveFeedbackApi }>({
  api: async ({ page }, runFixture, info) => {
    const api = new DatabaseSaveFeedbackApi();
    await api.install(page);
    page.on('response', response => {
      if (response.request().method() === 'PUT' && new URL(response.url()).pathname === databaseSavePaths.config)
        api.mark('browser-ack-received', { status: response.status() });
    });
    try { await runFixture(api); }
    finally {
      api.dispose();
      await info.attach('configuration-http-ledger', { body: JSON.stringify(api.ledger, null, 2), contentType: 'application/json' });
      expect(api.unexpected, 'Unexpected APIs/external transports abort and fail acceptance').toEqual([]);
      expect(api.pageErrors, 'No unhandled rejection or React failure').toEqual([]);
    }
  },
});
async function openAdvanced(page: Page) {
  await page.goto(advanced);
  await expect(page.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'W');
  await expect(saved(page)).toContainText('已儲存');
  const toggle = page.getByTestId('database-connection-card-toggle');
  if (await toggle.getAttribute('aria-expanded') === 'false') await toggle.click();
  await expect(name(page)).toHaveValue('Fixture SQLite');
  await page.getByTestId('step4-advanced-recording').locator(':scope > summary').click();
  await expect(page.getByTestId('group-create')).toBeEnabled();
}
async function capture(page: Page, info: TestInfo, label: string) {
  await info.attach(label, { body: await page.screenshot({ fullPage: true }), contentType: 'image/png' });
}
async function beginSave(page: Page, api: DatabaseSaveFeedbackApi, value: string, mode: AckMode = 'valid') {
  const ack = api.prepareSave(mode);
  await name(page).fill(value);
  expect(await ack.received).toMatchObject({ connector_id: 'db-W', name: value });
  await expect(saved(page)).toContainText('儲存中');
  return ack;
}
async function expectPending(page: Page, api: DatabaseSaveFeedbackApi, round: ReturnType<DatabaseSaveFeedbackApi['holdRefreshes']>) {
  await Promise.all(readKinds.map(kind => round[kind].started));
  await expect(saved(page)).toContainText('已儲存');
  await expect(refreshing(page)).toContainText('正在更新目前狀態');
  await expect(page.getByTestId('group-create')).toBeDisabled();
  api.mark('ui-saved-with-three-reads-held');
}
async function release(round: ReturnType<DatabaseSaveFeedbackApi['holdRefreshes']>, order: ReadKind[]) {
  for (const kind of order) { round[kind].release(); await round[kind].finished; }
}
async function cancelLeave(page: Page) {
  const dialog = page.getByRole('alertdialog');
  await expect(dialog).toBeVisible();
  await dialog.getByRole('button', { name: '繼續編輯', exact: true }).click();
  await expect(dialog).toBeHidden();
}
// An event fence, not a latency threshold or replacement for independently held reads.
async function rendered(page: Page) {
  await page.evaluate(() => new Promise<void>(done => requestAnimationFrame(() => requestAnimationFrame(() => done()))));
}

for (const order of [['config', 'targets', 'bootstrap'], ['bootstrap', 'targets', 'config']] as ReadKind[][]) {
  test(`scoped ACK saves before independently held reads (${order.join(' → ')})`, async ({ page, api }, info) => {
    await openAdvanced(page);
    const round = api.holdRefreshes(), ack = await beginSave(page, api, 'ACK confirmed');
    expect(api.writes).toHaveLength(1);
    ack.release(); await expectPending(page, api, round);
    // Failed runtime apply is separate from configuration saved, never SQL/running proof.
    expect(api.ledger.find(entry => entry.event === 'config-ack')?.detail).toMatchObject({ runtime_apply_status: 'apply_failed' });
    expect(api.ledger.some(entry => entry.event.startsWith('read-end:'))).toBe(false);
    await page.getByTestId('database-connection-card-toggle').click();
    await expect(name(page)).toBeHidden();
    await expect(saved(page)).toBeVisible(); await expect(refreshing(page)).toBeVisible();
    await capture(page, info, 'saved-while-three-reads-pending');
    for (const [index, kind] of order.entries()) {
      round[kind].release(); await round[kind].finished;
      if (index < 2) await expect(refreshing(page)).toContainText('正在更新目前狀態');
    }
    await expect(refreshing(page)).toHaveCount(0);
    await expect(saved(page)).toContainText('已儲存');
    await expect(page.getByTestId('group-create')).toBeEnabled();
    expect(api.writes).toHaveLength(1);
  });
}
for (const failed of readKinds) {
  test(`${failed} refresh failure keeps collapsed saved summary; retry is GET only`, async ({ page, api }, info) => {
    await openAdvanced(page);
    const round = api.holdRefreshes(), ack = await beginSave(page, api, 'Saved despite stale read');
    ack.release(); await expectPending(page, api, round);
    for (const kind of readKinds) { if (kind === failed) round[kind].fail(); else round[kind].release(); }
    await Promise.all(readKinds.map(kind => round[kind].finished));
    await page.getByTestId('database-connection-card-toggle').click();
    await expect(refreshing(page)).toHaveAttribute('role', 'alert');
    await expect(refreshing(page)).toContainText('設定已儲存，但狀態更新失敗');
    await expect(saved(page)).toContainText('已儲存');
    await expect(page.getByTestId('workbench-v2-bootstrap-error')).toHaveCount(0);
    await expect(page.getByTestId('group-create')).toBeDisabled();
    await capture(page, info, `saved-with-${failed}-refresh-error`);
    const retry = api.holdRefreshes(), before = api.writes.length;
    await refreshing(page).getByRole('button', { name: '重新讀取目前狀態', exact: true }).click();
    await expectPending(page, api, retry);
    await release(retry, ['targets', 'bootstrap', 'config']);
    await expect(refreshing(page)).toHaveCount(0);
    expect(api.writes).toHaveLength(before);
    await page.getByTestId('database-connection-card-toggle').click();
    await expect(name(page)).toHaveValue('Saved despite stale read');
  });
}
for (const mode of ['missing-id', 'foreign-workspace', 'missing-workspace', 'missing-identity', 'missing-setup', 'malformed', 'conflict', 'timeout'] as AckMode[]) {
  test(`${mode} ACK cannot claim saved or trigger the successful refresh path`, async ({ page, api }, info) => {
    await openAdvanced(page);
    const ack = await beginSave(page, api, `Unconfirmed ${mode}`, mode);
    ack.release(); await ack.finished;
    await expect(saved(page)).toContainText('儲存失敗');
    await expect(saved(page)).not.toContainText('已儲存');
    await expect(refreshing(page)).toHaveCount(0);
    await expect(name(page)).toHaveValue(`Unconfirmed ${mode}`);
    await expect(page.getByTestId('group-create')).toBeDisabled();
    expect(api.writes).toHaveLength(1);
    if (mode !== 'conflict') {
      await name(page).fill(`Newer local ${mode}`); await rendered(page);
      await expect(saved(page)).toContainText('儲存失敗'); expect(api.writes).toHaveLength(1);
    }
    await capture(page, info, `unconfirmed-${mode}`);
    if (mode === 'conflict') {
      const retry = await beginSave(page, api, 'Explicit edit after rejected conflict');
      retry.release(); await expect(saved(page)).toContainText('已儲存');
      expect(api.writes).toHaveLength(2);
    }
  });
}

test('late old refreshes cannot overwrite newer edits or newer ACK revisions', async ({ page, api }, info) => {
  await openAdvanced(page);
  const old = api.holdRefreshes(), first = await beginSave(page, api, 'First accepted');
  first.release(); await expectPending(page, api, old);
  const latest = api.holdRefreshes(), second = await beginSave(page, api, 'Newer text');
  expect(await second.received).toMatchObject({ expected_connector_revision: 'identity-1', expected_setup_revision: 'setup-1' });
  await expect(name(page)).toHaveValue('Newer text');
  second.release(); await expectPending(page, api, latest);
  await release(latest, ['targets', 'config', 'bootstrap']);
  await expect(refreshing(page)).toHaveCount(0);
  await release(old, ['bootstrap', 'config', 'targets']); await rendered(page);
  await expect(name(page)).toHaveValue('Newer text'); await expect(saved(page)).toContainText('已儲存');
  const configReads = api.ledger.filter(entry => entry.event === 'read-start:config');
  expect(configReads.map(entry => (entry.detail as { name: string }).name)).toEqual(['First accepted', 'Newer text']);
  const third = await beginSave(page, api, 'Third edit');
  expect(await third.received).toMatchObject({ expected_connector_revision: 'identity-2', expected_setup_revision: 'setup-2' });
  third.release(); await expect(saved(page)).toContainText('已儲存');
  expect(api.writes).toHaveLength(3); await capture(page, info, 'newer-edit-retained');
});

test('lost ACK never replays on repeated edits or document reopen; GET observes saved config', async ({ page, api }, info) => {
  await openAdvanced(page);
  const ack = await beginSave(page, api, 'Persisted but reply lost', 'lost');
  ack.release(); await ack.finished; await expect(saved(page)).toContainText('儲存失敗');
  await name(page).fill('Newer unsent intent'); await name(page).fill('Newest unsent intent');
  await rendered(page); expect(api.writes).toHaveLength(1);
  await expect(name(page)).toHaveValue('Newest unsent intent');
  await capture(page, info, 'lost-ack-local-intent');
  // No config-only route guard exists. Reopen is GET observation, not replay/recovery PUT.
  await page.goto(home); await expect(page.getByTestId('recording-home')).toBeVisible();
  await page.goBack(); await expect(page.getByTestId('database-connection-card')).toBeVisible();
  if (await page.getByTestId('database-connection-card-toggle').getAttribute('aria-expanded') === 'false')
    await page.getByTestId('database-connection-card-toggle').click();
  await expect(name(page)).toHaveValue('Persisted but reply lost');
  await expect(saved(page)).toContainText('已儲存'); expect(api.writes).toHaveLength(1);
});

test('late workspace A reads cannot replace a newly opened workspace B', async ({ page, api }, info) => {
  await openAdvanced(page);
  const old = api.holdRefreshes(), ack = await beginSave(page, api, 'Workspace A accepted');
  ack.release(); await expectPending(page, api, old);
  // Simulated backend session change, not a fictitious UI workspace selector.
  api.selectWorkspace('B');
  await page.goto(home); await expect(page.getByTestId('recording-home')).toBeVisible();
  await page.goto(advanced);
  await expect(page.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'B');
  if (await page.getByTestId('database-connection-card-toggle').getAttribute('aria-expanded') === 'false')
    await page.getByTestId('database-connection-card-toggle').click();
  await expect(name(page)).toHaveValue('Fixture B SQLite');
  await release(old, ['targets', 'bootstrap', 'config']); await rendered(page);
  await expect(page.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'B');
  await expect(name(page)).toHaveValue('Fixture B SQLite');
  expect(api.writes).toHaveLength(1); await capture(page, info, 'old-workspace-read-fenced');
});

for (const phase of ['unsent', 'sent', 'unknown'] as const) {
  test(`existing recording-draft Close/Escape guard survives ${phase} config state`, async ({ page, api }, info) => {
    await openAdvanced(page);
    await page.getByTestId('group-create').click();
    await page.getByTestId('group-name').fill('Independent unsaved recording');
    if (phase === 'unsent') { await name(page).fill(''); await expect(saved(page)).toContainText('草稿需要修正'); }
    const ack = phase === 'unsent' ? undefined : await beginSave(page, api, `Config ${phase}`, 'lost');
    if (ack && phase === 'unknown') { ack.release(); await expect(saved(page)).toContainText('儲存失敗'); }
    await page.getByTestId('group-editor-close').click();
    await capture(page, info, `existing-draft-guard-${phase}`); await cancelLeave(page);
    await expect(page.getByTestId('group-name')).toHaveValue('Independent unsaved recording');
    await page.getByTestId('group-editor-close').focus(); await page.keyboard.press('Escape'); await cancelLeave(page);
    await expect(page.getByTestId('group-name')).toHaveValue('Independent unsaved recording');
    expect(api.writes).toHaveLength(phase === 'unsent' ? 0 : 1);
    if (ack && phase === 'sent') { ack.release(); await expect(saved(page)).toContainText('儲存失敗'); }
    await page.getByTestId('group-editor-close').click();
    await page.getByRole('alertdialog').getByRole('button', { name: '捨棄未保存欄位並離開', exact: true }).click();
    await expect(page.getByTestId('group-editor')).toHaveCount(0);
    expect(api.writes).toHaveLength(phase === 'unsent' ? 0 : 1);
  });
}

test('genuine browser Back and Forward keep the existing recording draft on cancelled leave', async ({ page, api }, info) => {
  // Three real router entries: home -> new recording -> home. No pushState or injected Link.
  await page.goto(home);
  await page.getByTestId('recording-home').getByRole('button', { name: '新增記錄', exact: true }).click();
  await expect(page.getByTestId('group-name')).toBeVisible();
  await page.locator('.gw-simple-recording-title').getByRole('button').click();
  await page.getByRole('alertdialog').getByRole('button', { name: '捨棄未保存欄位並離開', exact: true }).click();
  await expect(page.getByTestId('recording-home')).toBeVisible(); await page.goBack();
  await expect(page.getByTestId('group-name')).toBeVisible();
  await page.getByTestId('group-name').fill('History draft remains');
  const editorURL = page.url();
  await page.goBack(); await cancelLeave(page);
  await expect(page).toHaveURL(editorURL); await expect(page.getByTestId('group-name')).toHaveValue('History draft remains');
  await page.goForward(); await cancelLeave(page);
  await expect(page).toHaveURL(editorURL); await expect(page.getByTestId('group-name')).toHaveValue('History draft remains');
  await capture(page, info, 'cancelled-back-and-forward-retain-draft'); expect(api.writes).toEqual([]);
});
