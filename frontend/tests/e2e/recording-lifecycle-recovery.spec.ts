import { expect, test } from '@playwright/test';
import { VIEWPORTS, attachLedger, basicFor, captureGeometry, dialogFor, openChange,
  openFixture, responseGate, showAdvanced, toActivationReview, toSchemaPlan } from './recording-lifecycle.helpers';

for (const viewport of VIEWPORTS) test.describe(`recording lifecycle recovery ${viewport.width}x${viewport.height}`, () => {
  test.use({ viewport, timezoneId: 'UTC' });

  test('LC-R01 unknown schema stays status-only through mode switch, reload and remount with the same root', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    const root = mock.snapshot().envelope;
    mock.setSchemaOperationStatus('unknown');
    mock.loseNext('schema-confirm');
    await dialog.getByTestId('recording-change-confirm_schema').click();
    await expect(dialog.getByTestId('recording-change-check_status')).toBeEnabled();
    await expect(dialog).toContainText('不要建立替代操作或另一份資料表計畫');
    await dialog.getByTestId('recording-change-check_status').click();
    await expect(dialog.getByTestId('recording-change-repreview')).toHaveCount(0);
    await expect(dialog.getByRole('button', { name: '重送完全相同請求' })).toHaveCount(0);
    const unknownEffects = mock.ledger.filter(entry => entry.mutation);
    await captureGeometry(page, info, 'unknown-schema-status-only', dialog);
    await dialog.getByRole('button', { name: '關閉', exact: true }).click();
    await expect(page.getByTestId('recording-change-recovery')).toBeVisible();
    const advanced = await showAdvanced(page);
    await expect(advanced.getByTestId('saved-recording-transition')).toBeDisabled();
    await expect(basicFor(page).getByTestId('saved-recording-transition').first()).toBeDisabled();
    await page.reload();
    await expect(page.getByTestId('recording-change-recovery')).toBeVisible();
    await page.getByTestId('fixture-remount-host').click();
    await page.getByRole('button', { name: '查看原操作', exact: true }).click();
    await dialog.getByTestId('recording-change-check_status').click();
    await expect(dialog).toContainText('不要建立替代操作');
    expect(mock.snapshot().envelope).toEqual(root);
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.mutation)).toEqual(unknownEffects);
    const statusRequests = mock.ledger.filter(entry => entry.action === 'operation-status').map(entry => entry.body);
    expect(statusRequests.length).toBeGreaterThanOrEqual(2);
    for (const body of statusRequests) expect(body).toEqual({ action: 'operation-status', workspace_id: 'ws-1', lineage_id: 'L', operation_id: root?.request.operation_id });
    mock.setSchemaOperationStatus('succeeded');
    await dialog.getByTestId('recording-change-check_status').click();
    await expect(dialog.getByText('資料表準備已確認，尚未開始新記錄。', { exact: true })).toBeVisible();
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(1);
    await captureGeometry(page, info, 'same-root-schema-recovered', dialog);
    await attachLedger(info, mock);
  });

  test('LC-R02 only server-proven unsubmitted preview enables explicit same-root repreview', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await dialog.getByTestId('recording-change-review').click();
    mock.loseNext('schema-preview');
    await dialog.getByTestId('recording-change-reserve_preview').click();
    await expect(dialog.getByTestId('recording-change-check_status')).toBeEnabled();
    const root = mock.snapshot().envelope;
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(1);
    await expect(dialog.getByTestId('recording-change-repreview')).toHaveCount(0);
    await dialog.getByTestId('recording-change-check_status').click();
    await expect(dialog.getByTestId('recording-change-repreview')).toBeEnabled();
    await expect(dialog).toContainText('原操作確認尚未提交任何資料表準備');
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(1);
    await captureGeometry(page, info, 'proven-unsubmitted-explicit-repreview', dialog);
    await dialog.getByTestId('recording-change-repreview').click();
    await expect(dialog.getByTestId('recording-change-confirm_schema')).toBeEnabled();
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(2);
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(0);
    expect(mock.snapshot().envelope).toEqual(root);
    await dialog.getByText('先前預覽操作編號', { exact: true }).click();
    await expect(dialog).toContainText('僅預覽，尚未提交');
    await captureGeometry(page, info, 'retained-prior-preview-history', dialog);
    await attachLedger(info, mock);
  });

  test('LC-R03 failed status retains unknown root and does not silently replay any effect', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    mock.loseNext('schema-confirm');
    mock.setSchemaOperationStatus('unknown');
    await dialog.getByTestId('recording-change-confirm_schema').click();
    await expect(dialog.getByTestId('recording-change-check_status')).toBeEnabled();
    const state = mock.snapshot();
    const mutations = mock.ledger.filter(entry => entry.mutation);
    mock.rejectNext('operation-status');
    await dialog.getByTestId('recording-change-check_status').click();
    await expect(dialog.getByRole('alert')).toContainText('先前結果仍保留');
    await expect(dialog.getByTestId('recording-change-repreview')).toHaveCount(0);
    expect(mock.ledger.filter(entry => entry.mutation)).toEqual(mutations);
    expect(mock.snapshot()).toEqual(state);
    await captureGeometry(page, info, 'status-failure-retains-root', dialog);
    await attachLedger(info, mock);
  });

  test('LC-R04 lost activation ACK recovers terminal result; late error cannot regress it after remount', async ({ page }, info) => {
    const gate = responseGate('activation-confirm');
    const mock = await openFixture(page, 'applied-wide', gate);
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    await toActivationReview(page);
    const root = mock.snapshot().envelope;
    mock.loseNext('activation-confirm');
    const failed = page.waitForEvent('requestfailed', { predicate: request => request.postData()?.includes('"action":"activation-confirm"') ?? false });
    try {
    await dialog.getByTestId('recording-change-confirm_cutover').click();
    await expect.poll(() => mock.snapshot().status?.cutover.status).toBe('succeeded');
    await dialog.getByRole('button', { name: '關閉', exact: true }).click();
    await page.getByTestId('fixture-remount-host').click();
    await page.getByRole('button', { name: '查看原操作', exact: true }).click();
    await dialogFor(page).getByTestId('recording-change-check_status').click();
    await expect(dialogFor(page).getByText('已記錄確認時間的開始／切換排程。', { exact: true })).toBeVisible();
    gate.release();
    await failed;
    await expect(dialogFor(page).getByText('已記錄確認時間的開始／切換排程。', { exact: true })).toBeVisible();
    await expect(dialogFor(page).getByRole('alert')).toHaveCount(0);
    await expect(dialogFor(page).getByTestId('recording-change-repreview')).toHaveCount(0);
    await page.reload();
    await expect(page.getByTestId('recording-change-recovery')).toContainText('已記錄確認時間的開始／切換排程');
    await page.getByRole('button', { name: '查看原操作', exact: true }).click();
    await expect(dialogFor(page).getByRole('alert')).toHaveCount(0);
    expect(mock.snapshot().envelope).toEqual(root);
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(1);
    await captureGeometry(page, info, 'terminal-survives-late-error-reload', dialogFor(page));
    await attachLedger(info, mock);
    } finally { gate.release(); }
  });
});
