import { expect, test } from '@playwright/test';
import { FUTURE_BOUNDARY, VIEWPORTS, assertNoNewEffects, assertPendingUnits, attachLedger, captureGeometry, fillBoundary,
  openChange, openFixture, toActivationReview, toSchemaPlan } from './recording-lifecycle.helpers';

for (const viewport of VIEWPORTS) test.describe(`recording lifecycle transitions ${viewport.width}x${viewport.height}`, () => {
  test.use({ viewport, timezoneId: 'Asia/Taipei' });

  for (const variant of [
    { name: 'wide-to-long', scenario: 'applied-wide', layout: 'managed_long_v1', destination: 'conn-1' },
    { name: 'long-to-wide', scenario: 'applied-long', layout: 'managed_wide_profile_v1', destination: 'conn-1' },
    { name: 'same-mode-destination', scenario: 'applied-long', layout: 'managed_long_v1', destination: 'conn-2' },
  ] as const) test(`LC-T01 ${variant.name}: schema plan and future cutover require distinct confirmation`, async ({ page }, info) => {
    const mock = await openFixture(page, variant.scenario);
    const original = mock.snapshot().source;
    const { dialog } = await openChange(page, false, variant.name === 'long-to-wide' ? 'advanced' : 'basic');
    await dialog.getByLabel('資料表布局', { exact: true }).selectOption(variant.layout);
    await dialog.getByLabel('目的地', { exact: true }).selectOption(variant.destination);
    await captureGeometry(page, info, `${variant.name}-choose`, dialog);
    await dialog.getByTestId('recording-change-review').click();
    await assertPendingUnits(dialog);
    await expect(dialog).toContainText('保留原資料表、歷史與收據');
    assertNoNewEffects(mock);
    await captureGeometry(page, info, `${variant.name}-impact`, dialog);
    mock.delayNext('reserve', 250);
    await dialog.getByTestId('recording-change-reserve_preview').dblclick({ delay: 30 });
    await expect(dialog.getByTestId('recording-change-confirm_schema')).toBeEnabled();
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(0);
    const prepared = mock.snapshot();
    expect(prepared.source).toEqual(original);
    expect(prepared.target?.id).toBe('G2');
    expect(prepared.target?.applied_revision).toBe('');
    expect(prepared.target?.layout_version).toBe(variant.layout);
    expect(prepared.target?.storage_mode ?? 'plain').toBe('plain');
    expect(prepared.target?.destination.connector_id).toBe(variant.destination);
    expect(prepared.target?.destination.table_name).not.toBe(original.destination.table_name);
    await expect(dialog).toContainText('設定已保留，尚未開始新記錄');
    await expect(dialog).not.toContainText('mock-only.db');
    await captureGeometry(page, info, `${variant.name}-schema-plan`, dialog);
    mock.delayNext('schema-confirm', 250);
    await toActivationReview(page, true);
    await expect(dialog.getByLabel('未來開始時間（您的當地時間）', { exact: true })).toHaveValue(/^2030-01-01T09:00(?::00)?$/);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(0);
    await expect(dialog.getByText('依來源觀測時間決定歸屬：邊界前屬舊世代；等於或晚於邊界屬新世代。')).toBeVisible();
    await expect(dialog.getByText('已審核的切換保留舊歷史，僅讓符合條件的新觀測歸入新世代。')).toBeVisible();
    await expect(dialog.getByTestId('recording-change-confirm_first_start')).toHaveCount(0);
    await captureGeometry(page, info, `${variant.name}-future-cutover`, dialog);
    mock.delayNext('activation-confirm', 250);
    await dialog.getByTestId('recording-change-confirm_cutover').dblclick({ delay: 30 });
    await expect(dialog.getByText('已記錄確認時間的開始／切換排程。', { exact: true })).toBeVisible();
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'first-start-confirm')).toHaveLength(0);
    expect(mock.snapshot().status?.cutover.result?.boundary).toBe(FUTURE_BOUNDARY);
    await expect(dialog).toContainText('這是設定結果，不代表設備正在執行，也不代表 SQL 資料已提交');
    await captureGeometry(page, info, `${variant.name}-terminal`, dialog);
    await attachLedger(info, mock);
  });

  test('LC-T02 first start has an explicit future boundary and no invented historical cutover', async ({ page }, info) => {
    const mock = await openFixture(page, 'deleted-unapplied');
    const { dialog } = await openChange(page, true);
    await toSchemaPlan(page);
    await toActivationReview(page);
    await expect(dialog.getByText('這是首次啟用，不會結束不存在的歷史啟用世代。')).toBeVisible();
    await expect(dialog.getByTestId('recording-change-confirm_cutover')).toHaveCount(0);
    await captureGeometry(page, info, 'first-start-future-boundary', dialog);
    mock.delayNext('first-start-confirm', 250);
    await dialog.getByTestId('recording-change-confirm_first_start').dblclick({ delay: 30 });
    await expect(dialog.getByText('已記錄確認時間的開始／切換排程。', { exact: true })).toBeVisible();
    expect(mock.ledger.filter(entry => entry.action === 'first-start-confirm')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(0);
    expect(mock.snapshot().status?.first_start.result?.boundary).toBe(FUTURE_BOUNDARY);
    expect(mock.snapshot().source.status).toBe('deleted');
    await expect(dialog).toContainText('不代表 SQL 資料已提交');
    await captureGeometry(page, info, 'first-start-terminal', dialog);
    await attachLedger(info, mock);
  });

  test('LC-T03 Basic-only tombstone recreation locks shape and creates only one distinct reserved generation', async ({ page }, info) => {
    const mock = await openFixture(page, 'deleted-basic');
    const original = mock.snapshot().source;
    const { dialog } = await openChange(page, true);
    await expect(dialog.getByLabel('資料表布局', { exact: true })).toBeDisabled();
    await expect(dialog.getByLabel('目的地', { exact: true })).toBeDisabled();
    await expect(dialog).toContainText('基本記錄沿用原布局、目的地與完整來源重新建立');
    await toSchemaPlan(page);
    const state = mock.snapshot();
    expect(state.source).toEqual(original);
    expect(state.target?.id).toBe('G2');
    expect(state.target?.members).toEqual(original.members);
    expect(state.target?.layout_version).toBe('managed_wide_v1');
    expect(state.target?.recording_profile).toBeUndefined();
    expect(state.target?.applied_revision).toBe('');
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(0);
    await assertPendingUnits(dialog);
    await captureGeometry(page, info, 'basic-same-shape-recreate', dialog);
    await attachLedger(info, mock);
  });

  test('LC-T04 initial blocker and stale scope explain repair without reservation or DDL', async ({ page }, info) => {
    const mock = await openFixture(page, 'review-blocked-wide');
    const { dialog } = await openChange(page);
    await dialog.getByTestId('recording-change-review').click();
    await expect(dialog.getByRole('alert')).toContainText('目前證據不足');
    await expect(dialog.getByTestId('recording-change-reserve_preview')).toBeDisabled();
    assertNoNewEffects(mock);
    await captureGeometry(page, info, 'impact-blocked', dialog);
    await dialog.getByRole('button', { name: '取消', exact: true }).click();
    await openChange(page);
    mock.staleNextLifecycle();
    await dialog.getByTestId('recording-change-review').click();
    await expect(dialog.getByRole('alert')).toContainText('保存版本已變更');
    expect(mock.ledger.filter(entry => entry.action === 'reserve')).toHaveLength(0);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(0);
    await captureGeometry(page, info, 'stale-scope-repair', dialog);
    await attachLedger(info, mock);
  });

  test('LC-T05 past or changed boundary cannot silently reuse final activation review', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    await dialog.getByTestId('recording-change-confirm_schema').click();
    await fillBoundary(page, '2029-12-31T23:00:00.000Z');
    await dialog.getByTestId('recording-change-review_activation').click();
    await expect(dialog.getByRole('alert')).toContainText('請選擇未來時間');
    expect(mock.ledger.filter(entry => entry.action === 'activation-review')).toHaveLength(0);
    await fillBoundary(page);
    await dialog.getByTestId('recording-change-review_activation').click();
    await expect(dialog.getByTestId('recording-change-confirm_cutover')).toBeEnabled();
    await fillBoundary(page, '2030-01-01T02:00:00.000Z');
    await expect(dialog.getByTestId('recording-change-confirm_cutover')).toHaveCount(0);
    await expect(dialog.getByTestId('recording-change-review_activation')).toBeEnabled();
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(0);
    await captureGeometry(page, info, 'changed-boundary-requires-review', dialog);
    await attachLedger(info, mock);
  });

  test('LC-T06 final stale CAS preserves the reserved root but sends no activation confirmation', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    await toActivationReview(page);
    const root = mock.snapshot().envelope;
    mock.staleNextLifecycle();
    await dialog.getByTestId('recording-change-confirm_cutover').click();
    await expect(dialog.getByRole('alert')).toContainText('保存版本已變更');
    await expect(dialog.getByTestId('recording-change-confirm_cutover')).toHaveCount(0);
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(0);
    expect(mock.ledger.filter(entry => entry.action === 'first-start-confirm')).toHaveLength(0);
    expect(mock.snapshot().envelope).toEqual(root);
    expect(mock.snapshot().target?.applied_revision).toBe('');
    await captureGeometry(page, info, 'final-stale-cas', dialog);
    await attachLedger(info, mock);
  });

  test('LC-T07 unready activation retains prepared schema but never offers final start', async ({ page }, info) => {
    const mock = await openFixture(page, 'applied-wide');
    const { dialog } = await openChange(page);
    await toSchemaPlan(page);
    await dialog.getByTestId('recording-change-confirm_schema').click();
    mock.setActivationReady(false);
    await fillBoundary(page);
    await dialog.getByTestId('recording-change-review_activation').click();
    await expect(dialog.getByRole('alert')).toContainText('此操作目前受阻');
    await expect(dialog.getByTestId('recording-change-confirm_cutover')).toHaveCount(0);
    await expect(dialog.getByTestId('recording-change-confirm_first_start')).toHaveCount(0);
    expect(mock.ledger.filter(entry => entry.action === 'activation-readiness')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.action === 'activation-confirm')).toHaveLength(0);
    expect(mock.snapshot().target?.applied_revision).toBe('');
    await captureGeometry(page, info, 'activation-readiness-blocked', dialog);
    await attachLedger(info, mock);
  });
});
