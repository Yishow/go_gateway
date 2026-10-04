import assert from 'node:assert/strict';
import { observeJSON } from './fresh-ui.mjs';

function rowsByRecordID(rows, label) {
  const result = new Map();
  for (const row of rows) {
    assert.ok(row.record_id, label + ' row has no record_id');
    assert.equal(result.has(row.record_id), false, label + ' has duplicate record_id ' + row.record_id);
    result.set(row.record_id, row);
  }
  return result;
}

export function assertRowsPreserved(before, after, label, runtime) {
  const afterByID = rowsByRecordID(after, label + ' after');
  for (const [recordID, row] of rowsByRecordID(before, label + ' before')) {
    assert.deepEqual(afterByID.get(recordID), row, label + ' changed runtime record ' + recordID);
  }
  if (!runtime) return;
  const { group, table, outbox, receipts } = runtime;
  const entityColumn = group.row_policy.entity_key_column;
  for (const row of after) {
    const effect = outbox.find((entry) => entry.record_id === row.record_id);
    const receipt = effect && receipts.find((entry) => entry.effect_key === effect.effect_key);
    assert.ok(effect && effect.state === 'sql_committed' && effect.group_id === group.id &&
      effect.group_revision === group.applied_revision && effect.table_name === table &&
      row.group_id === group.id && row[entityColumn] && effect.entity_key === row[entityColumn] &&
      Number.isFinite(Date.parse(row.bucket_start)) && Date.parse(effect.bucket_start) === Date.parse(row.bucket_start) &&
      receipt?.payload_digest === effect.payload_digest && receipt?.committed_at,
    label + ' row lacks current durable runtime evidence: ' + row.record_id);
  }
}

/** Allow a just-committed runtime row to finish its local receipt before checking. */
export async function observePreservedRuntimeRows({ before, group, table, readRows, readDelivery }) {
  const deadline = Date.now() + 15_000;
  while (true) {
    const after = readRows();
    const delivery = readDelivery();
    try {
      assertRowsPreserved(before, after, 'test-write', { group, table, ...delivery });
      return { after, delivery };
    } catch (error) {
      if (Date.now() >= deadline) throw error;
      await new Promise((resolve) => setTimeout(resolve, 250));
    }
  }
}

export async function restoreUniqueDraftColumn({ page, baseURL, group, pointID, originalColumn, readRows }) {
  const currentMember = group.members.find((member) => member.point_id === pointID);
  assert.ok(currentMember && currentMember.target_column !== originalColumn, 'B int16 must start on the shared column');
  const before = readRows();
  await page.reload();
  await page.getByTestId('step-nav-button-4').click();
  const advanced = page.getByTestId('step4-advanced-recording');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-open-' + group.id).click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
  const columnControl = page.getByTestId('group-member-column-' + pointID);
  if (await columnControl.evaluate((element) => element.tagName) === 'SELECT') {
    await columnControl.selectOption(originalColumn);
  } else {
    await columnControl.fill(originalColumn);
  }
  assert.equal(await columnControl.inputValue(), originalColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  const data = await observeJSON(baseURL, '/studio-v2/workspace/write-groups');
  const saved = data.groups?.find((item) => item.id === group.id);
  assert.ok(saved, 'saved matrix group disappeared');
  assert.notEqual(saved.revision, group.revision, 'draft save must create a new revision');
  assert.equal(saved.applied_revision, group.applied_revision, 'draft save must retain the applied revision');
  assert.equal(saved.destination.storage_strategy, 'custom');
  assert.equal(saved.members.find((member) => member.point_id === pointID)?.target_column, originalColumn);
  assertRowsPreserved(before, readRows(), 'draft save');
  await page.getByTestId('group-editor-close').click();
  return saved;
}

/**
 * The operation detail (row identity, ownership marker) is service-owned and
 * never exposed by the API, so it is read from the local ledger row only.
 */
export function readLedgerDetail({ apiOperation, ledgerRows, preview }) {
  assert.equal(Object.hasOwn(apiOperation, 'detail'), false, 'API operation exposes detail');
  assert.equal(ledgerRows.length, 1, 'expected exactly one local ledger row for ' + preview.operation_id);
  const [ledger] = ledgerRows;
  assert.equal(ledger.operation_id, preview.operation_id);
  assert.equal(ledger.action, 'test_write', 'ledger action is not test_write');
  const detail = JSON.parse(ledger.detail);
  assert.equal(detail.owner_column, preview.owner_column, 'ledger owner column mismatch');
  assert.equal(detail.owner_value, preview.owner_value, 'ledger owner value mismatch');
  assert.equal(detail.table, preview.target.table);
  return detail;
}

export async function runFreshMatrixTestWrite({ page, baseURL, group, readRows, readDelivery, readLedger }) {
  assert.equal(typeof readRows, 'function');
  assert.equal(typeof readDelivery, 'function');
  assert.equal(typeof readLedger, 'function');
  const before = readRows();
  await page.reload();
  await page.getByTestId('step-nav-button-4').click();
  const advanced = page.getByTestId('step4-advanced-recording');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-open-' + group.id).click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
  const previewButton = page.getByTestId('group-test-write-preview');
  assert.equal(await previewButton.isEnabled(), true, 'matrix group test-write preview is not enabled');
  const [previewResponse] = await Promise.all([
    page.waitForResponse((response) => response.url().endsWith(`/write-groups/${group.id}/test-write-preview`) && response.request().method() === 'POST'),
    previewButton.click(),
  ]);
  const previewBody = await previewResponse.json();
  assert.equal(previewResponse.status(), 200, 'matrix test-write preview was not accepted: ' + JSON.stringify(previewBody));
  const preview = previewBody.data;
  assert.ok(preview?.operation_id && preview?.token, 'matrix test-write preview lacks operation identity');
  assert.equal(preview.group_id, group.id);
  assert.equal(preview.group_revision, group.revision);
  assert.equal(preview.target.table, group.destination.table_name);
  assert.equal(preview.owner_column, group.row_policy.entity_key_column);
  assert.match(preview.owner_value, /^gw-test-[a-f0-9-]+$/);
  assert.equal(preview.values.find((value) => value.column === preview.owner_column)?.value, preview.owner_value);
  await page.getByTestId('group-test-write-preview-card').waitFor({ timeout: 20_000 });
  assertRowsPreserved(before, readRows(), 'preview');
  const [confirmResponse] = await Promise.all([
    page.waitForResponse((response) => response.url().endsWith(`/write-groups/${group.id}/test-write`) && response.request().method() === 'POST'),
    page.getByTestId('group-test-write-confirm').click(),
  ]);
  const confirmBody = await confirmResponse.json();
  assert.equal(confirmResponse.status(), 200, 'matrix test-write confirmation was not terminal: ' + JSON.stringify(confirmBody));
  const operation = confirmBody.data;
  assert.equal(operation.operation_id, preview.operation_id);
  assert.equal(operation.action, 'test_write');
  assert.equal(operation.write_outcome, 'written_verified');
  assert.equal(operation.cleanup_status, 'cleaned');
  await page.getByTestId('group-test-write-outcome').waitFor({ timeout: 20_000 });
  const ledger = await observeJSON(baseURL, '/studio-v2/workspace/database-operations/' + encodeURIComponent(preview.operation_id));
  assert.equal(ledger.operation_id, preview.operation_id);
  const detail = readLedgerDetail({ apiOperation: ledger, ledgerRows: readLedger(preview.operation_id), preview });
  assert.equal(detail.table, group.destination.table_name);
  const { after, delivery } = await observePreservedRuntimeRows({ before, group,
    table: group.destination.table_name, readRows, readDelivery });
  assert.equal(after.some((row) => String(row[preview.owner_column]) === preview.owner_value), false,
    'operation-owned test row remained after cleanup');
  return {
    operation_id: operation.operation_id, owner_column: preview.owner_column, owner_value: preview.owner_value,
    write_outcome: operation.write_outcome, cleanup_status: operation.cleanup_status,
    before_record_ids: [...rowsByRecordID(before, 'test-write baseline').keys()],
    after_record_ids: [...rowsByRecordID(after, 'test-write final').keys()],
    preview,
    operation,
    ledger_detail: detail,
    runtime_neighbor_delivery: delivery,
  };
}
