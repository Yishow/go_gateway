// Real UI regression for selection, reviewed bulk changes, lost DDL replies,
// refresh recovery and narrow Basic membership updates. Only loopback fixtures.
import assert from 'node:assert/strict';
import { existsSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, assertFreshPortFree, binaryIdentity, closeFreshUI,
  configureFreshSQLite, isFreshCleanExit, makeFreshRun, observeJSON, openFreshUI,
  persistFreshMappings, prepareFreshBasicGroup, previewAndApplyFreshSchema,
  setupFreshDevice, startFreshBasic, startFreshGateway,
  startFreshSimulator, stopProcesses, waitForFreshGateway, waitForObservation,
} from './fresh-ui.mjs';
import { readSQLiteRows } from './fresh-sqlite-observation.mjs';
import { sh } from './lib.mjs';
import { writeFreshResult } from './fresh-result.mjs';

const runId = process.env.F_RUN_ID ?? `simplify-flow-${Date.now()}`;
const port = Number(process.env.GW_PORT ?? 3484), simPort = Number(process.env.F_SIM_PORT_A ?? 15484);
const run = makeFreshRun(runId), children = [];
const gatewayBinary = process.env.F_GATEWAY_BINARY;
assert.ok(gatewayBinary, 'explicit final binary required');
let ui;
const result = { run_id: runId, passed: false, checks: {}, failures: [] };
const q = value => `"${String(value).replaceAll('"', '""')}"`;
async function sqlUntil(sql, predicate, timeout = 45_000) {
  const deadline = Date.now() + timeout;
  while (Date.now() < deadline) {
    const rows = readSQLiteRows(run.destination_db, sql);
    if (predicate(rows)) return rows;
    await new Promise(resolve => setTimeout(resolve, 500));
  }
  throw new Error('SQL observation timeout');
}
try {
  await Promise.all([assertFreshPortFree(port), assertFreshPortFree(simPort)]);
  children.push(startFreshSimulator({ work: run.work, port: simPort, registers: { 0: 32767, 1: 65535, 2: 65535, 3: 65535, 4: 32768, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 777 } }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port }));
  const base = `http://127.0.0.1:${port}`;
  await waitForFreshGateway(base);
  ui = await openFreshUI({ port });
  const { page } = ui;
  await setupFreshDevice(page, { name: 'Simplification fixture', port: simPort });
  // Grid click and Shift range cannot mutate collection settings.
  await page.getByTestId('rule-add-btn').click();
  await page.getByTestId('rule-name-input').fill('Read a range');
  await page.getByTestId('rule-count-input').fill('2');
  await page.getByTestId('rule-start-input').fill('40001');
  await page.getByTestId('rule-datatype-select').selectOption('int64');
  await page.waitForTimeout(1500);
  const cells = page.locator('[data-testid^="point-cell-"]');
  await cells.first().click(); await cells.nth(1).click({ modifiers: ['Shift'] });
  assert.equal(await cells.locator('input:checked').count(), 2);
  assert.equal(await page.locator('[data-selected="true"]').count(), 2);
  result.checks.grid_selection_preserves_collection = true;
  await page.getByTestId('continue-step3-btn').click();
  assert.equal(await page.getByRole('columnheader').count(), 5);
  await page.getByTestId('mapping-advanced-toggle').click();
  const rows = page.locator('[data-testid^="mapping-row-"]');
  const ids = await rows.evaluateAll(nodes => nodes.map(n => n.dataset.testid.replace('mapping-row-', '')));
  await page.getByTestId(`input-unit-${ids[0]}`).fill('C');
  await page.getByTestId(`input-unit-${ids[1]}`).fill('bar');
  await page.getByTestId(`input-unit-${ids[0]}`).focus();
  await page.getByTestId(`bulk-target-${ids[1]}`).check();
  await page.getByTestId('btn-bulk-apply-unit').click();
  assert.ok((await page.getByTestId('bulk-preview').innerText()).includes('bar → C'));
  await page.getByTestId('bulk-preview').getByRole('button', { name: /Cancel|取消/ }).click();
  assert.equal(await page.getByTestId(`input-unit-${ids[1]}`).inputValue(), 'bar');
  await page.getByTestId('mapping-advanced-toggle').click();
  await page.getByTestId('mapping-advanced-toggle').click();
  assert.equal(await page.getByTestId(`input-unit-${ids[1]}`).inputValue(), 'bar');
  result.checks.review_cancel_and_view_switch_preserve_values = true;
  await page.getByTestId('mapping-advanced-toggle').click();
  await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, `${runId}-mapping.png`), fullPage: true });
  await persistFreshMappings(page, base, { tagPrefix: 'simplify' });
  await configureFreshSQLite(page, base, { path: run.destination_db });
  assert.equal(await page.getByTestId('database-connection-card-toggle').getAttribute('aria-expanded'), 'false');
  await page.getByTestId('basic-recording-interval').fill('2');
  const prepared = await prepareFreshBasicGroup(page, base);
  assert.equal(existsSync(run.destination_db), false, 'read-only preview created target DB');
  let lostOperation;
  const applyPosts = [];
  page.on('request', request => { if (request.method() === 'POST' && request.url().endsWith('/schema-apply')) applyPosts.push(request.postDataJSON().operation_id); });
  await page.route('**/write-groups/*/schema-apply', async route => {
    lostOperation = route.request().postDataJSON().operation_id;
    const response = await route.fetch();
    assert.equal(response.status(), 200);
    await route.abort('failed'); // Lose the real reply after the gateway executes it.
  }, { times: 1 });
  await page.getByTestId('group-schema-apply').click();
  await page.getByTestId('group-schema-check').waitFor();
  await page.reload(); await page.getByTestId('step-nav-button-4').click();
  await page.getByTestId('basic-recording-prepare').click();
  await page.getByTestId('group-schema-check').waitFor();
  assert.equal(await page.getByTestId('group-schema-preview').isDisabled(), true);
  await page.getByTestId('group-schema-check').click();
  await page.waitForFunction(() => /succeeded|成功/.test(document.querySelector('[data-testid="group-schema-operation"]')?.textContent ?? ''));
  assert.deepEqual(applyPosts, [lostOperation]);
  result.checks.lost_ddl_reply_refresh_queries_original_operation = { operation_id: lostOperation, apply_count: applyPosts.length };
  await startFreshBasic(page);
  const groupPath = '/studio-v2/workspace/write-groups';
  let saved = (await observeJSON(base, groupPath)).groups.find(g => g.id === prepared.group.id);
  const beforeRows = await sqlUntil(`SELECT record_id, CAST(${q(saved.members[0].target_column)} AS TEXT) AS signed_max, CAST(${q(saved.members[1].target_column)} AS TEXT) AS signed_min FROM ${q(saved.destination.table_name)}`, rows => rows.length >= 2);
  assert.equal(beforeRows[0].signed_max, '9223372036854775807');
  assert.equal(beforeRows[0].signed_min, '-9223372036854775808');
  result.checks.signed_int64_sql_precision = { max: beforeRows[0].signed_max, min: beforeRows[0].signed_min };
  await page.getByTestId('step-nav-button-2').click();
  await page.getByTestId('rule-count-input').fill('3');
  await page.waitForTimeout(1500);
  await page.getByTestId('continue-step3-btn').click();
  await persistFreshMappings(page, base, { tagPrefix: 'simplify' });
  await page.getByTestId('basic-recording-members-review').click();
  await page.getByTestId('basic-members-confirm').waitFor();
  assert.equal(await page.getByTestId('basic-recording-start').isDisabled(), true);
  const updateResponse = page.waitForResponse(r => r.request().method() === 'PUT' && r.url().endsWith(`/write-groups/${saved.id}`));
  await page.getByTestId('basic-members-confirm').click();
  const updated = await updateResponse;
  assert.equal(updated.status(), 200, await updated.text());
  const groupBefore = saved;
  saved = await waitForObservation(base, groupPath, data => data.groups.some(g => g.id === groupBefore.id && g.members.length === 3)).then(data => data.groups.find(g => g.id === groupBefore.id));
  assert.deepEqual(saved.row_policy, groupBefore.row_policy);
  assert.deepEqual(saved.write_policy, groupBefore.write_policy);
  for (const member of groupBefore.members) {
    const retained = saved.members.find(m => m.point_id === member.point_id && m.tag_id === member.tag_id);
    assert.equal(retained.target_column, member.target_column); assert.equal(retained.required, member.required);
  }
  // A second explicit preparation must preview the newly added column.
  await prepareFreshBasicGroup(page, base);
  await previewAndApplyFreshSchema(page);
  await startFreshBasic(page);
  const added = saved.members.find(m => !groupBefore.members.some(old => old.point_id === m.point_id));
  const newRows = await sqlUntil(`SELECT record_id, CAST(${q(added.target_column)} AS TEXT) AS added FROM ${q(saved.destination.table_name)}`, rows => rows.some(r => r.added === '777'));
  assert.ok(beforeRows.every(old => newRows.some(row => row.record_id === old.record_id)), 'membership update removed history');
  result.checks.new_point_cas_schema_and_sql = { group_id: saved.id, before_members: 2, after_members: 3, old_rows_retained: beforeRows.length, added_value: '777' };

  // Fault DDL is restricted to this runner's owned SQLite target. No data is
  // inserted by the harness; both decisions operate on real gateway rows.
  const decisionEvidence = [];
  for (const resolution of ['retry', 'skip']) {
    const trigger = `simplify_owned_${resolution}`;
    sh('sqlite3', [run.destination_db, `CREATE TRIGGER ${trigger} BEFORE INSERT ON ${q(saved.destination.table_name)} BEGIN SELECT RAISE(ABORT,'owned rejected row'); END;`]);
    const attention = await waitForObservation(base, `/studio-v2/workspace/write-groups/${saved.id}/delivery`, data => data.attention?.some(item => item.state === 'quarantined'), 30_000).then(data => data.attention.find(item => item.state === 'quarantined'));
    sh('sqlite3', [run.destination_db, `DROP TRIGGER ${trigger};`]);
    const literal = value => `'${String(value).replaceAll("'", "''")}'`;
    const effectSQL = `SELECT effect_key,record_id,payload_digest,group_revision,connector_revision,state FROM wg_delivery_outbox WHERE effect_key=${literal(attention.effect_key)}`;
    const before = readSQLiteRows(run.gateway_db, effectSQL)[0];
    const note = page.getByTestId(`delivery-reason-${attention.effect_key}`);
    await note.fill(`Owned destination repaired; reviewed ${resolution}`);
    if (resolution === 'skip') {
      const skipButton = page.getByTestId(`delivery-skip-${attention.effect_key}`);
      assert.equal(await skipButton.isVisible(), false);
      await page.getByTestId(`delivery-skip-details-${attention.effect_key}`).click();
      await note.locator('xpath=ancestor::li').getByRole('checkbox').check();
    }
    const replyPromise = page.waitForResponse(reply => reply.request().method() === 'POST' && reply.url().endsWith(`/write-groups/${saved.id}/delivery/resolve`));
    await page.getByTestId(`delivery-${resolution}-${attention.effect_key}`).click();
    const reply = await replyPromise; assert.equal(reply.status(), 200);
    const request = reply.request().postDataJSON();
    const expected = resolution === 'retry' ? 'sql_committed' : 'operator_skipped';
    const deadline = Date.now() + 30_000;
    let after;
    do {
      after = readSQLiteRows(run.gateway_db, effectSQL)[0];
      if (after.state === expected) break;
      await new Promise(resolve => setTimeout(resolve, 300));
    } while (Date.now() < deadline);
    assert.equal(after.state, expected);
    assert.deepEqual({ ...after, state: before.state }, before);
    const audits = readSQLiteRows(run.gateway_db, `SELECT id FROM workspace_audit_history WHERE id=${literal(request.decision_id)}`);
    assert.equal(audits.length, 1);
    const delivered = readSQLiteRows(run.destination_db, `SELECT record_id FROM ${q(saved.destination.table_name)} WHERE record_id=${literal(before.record_id)}`);
    assert.equal(delivered.length, resolution === 'retry' ? 1 : 0);
    decisionEvidence.push({ resolution, decision_id: request.decision_id, before, after, audit_count: 1, destination_row_count: delivered.length });
  }
  result.checks.real_retry_and_explicit_skip = decisionEvidence;
  const screenshot = `${runId}.png`;
  await page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, screenshot), fullPage: true });
  result.screenshot = screenshot;
  result.binary = binaryIdentity(gatewayBinary);
  result.page_errors = ui.pageErrors;
  assert.deepEqual(ui.pageErrors, []);
  result.passed = true;
} catch (error) {
  result.failures.push(error?.stack ?? String(error));
  if (ui) { result.screen_text = await ui.page.locator('body').innerText().catch(() => ''); await ui.page.screenshot({ path: join(DEFAULT_EVIDENCE_DIR, `${runId}-error.png`), fullPage: true }).catch(() => {}); }
} finally {
  await closeFreshUI(ui?.browser);
  const processes = await stopProcesses(children);
  if (processes.every(isFreshCleanExit)) rmSync(run.work, { recursive: true });
  result.cleanup = { processes, namespace_absent: !existsSync(run.work) };
  if (!result.cleanup.namespace_absent) { result.passed = false; result.failures.push('owned cleanup incomplete'); }
  writeFreshResult(join(DEFAULT_EVIDENCE_DIR, `${runId}.json`), result);
  console.log(JSON.stringify({ passed: result.passed, checks: result.checks, failures: result.failures }));
}
