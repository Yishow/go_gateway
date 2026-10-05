// Narrow F acceptance for the changed Basic flow. Configuration mutations use
// the real UI; fault DDL is confined to the run-owned SQLite target.
import assert from 'node:assert/strict';
import { join } from 'node:path';
import { sh } from './lib.mjs';
import { readSQLiteRows, quoteSQLite } from './fresh-sqlite-observation.mjs';
import { observeJSON, openFreshAdvanced, startFreshGateway, stopProcesses, waitForFreshGateway, isFreshCleanExit } from './fresh-ui.mjs';

const identifier = (value) => `"${String(value).replaceAll('"', '""')}"`;
async function until(read, predicate, timeout = 100_000) {
  const deadline = Date.now() + timeout;
  let last;
  while (Date.now() < deadline) {
    last = await read();
    if (predicate(last)) return last;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  throw new Error(`acceptance timeout: ${JSON.stringify(last).slice(0, 500)}`);
}

async function observeDeviceSwitch(page, base, devices, tags, evidenceDir, runId) {
  await page.goto(`${base}/studio/runtime?device_id=${encodeURIComponent(devices[0].id)}`);
  await page.getByTestId('runtime-dashboard-device-switcher').waitFor({ timeout: 30_000 });
  const observations = [];
  for (const [device, expected] of [[devices[0], '215'], [devices[1], '187'], [devices[0], '215']]) {
    await page.getByTestId('runtime-dashboard-device-switcher').getByRole('button', { name: device.name, exact: true }).click();
    const own = tags.find((tag) => tag.device_id === device.id && tag.address === '40001');
    const foreign = tags.find((tag) => tag.device_id !== device.id && tag.address === '40001');
    assert.ok(own && foreign, 'same-address mappings missing');
    await page.waitForFunction(({ tag, expectedValue }) => {
      const rows = document.querySelectorAll('[data-testid="runtime-dashboard-live-points-table"] tbody tr');
      return [...rows].some((row) => row.cells[0]?.querySelector('div > span')?.textContent === tag && row.cells[1]?.textContent.trim() === expectedValue && row.cells[2]?.textContent.trim() === expectedValue);
    }, { tag: own.tag_key, expectedValue: expected }, { timeout: 25_000 });
    const text = await page.getByTestId('runtime-dashboard-live-points-table').innerText();
    const names = await page.getByTestId('runtime-dashboard-live-points-table').locator('tbody tr td:first-child div:first-child > span:first-child').allTextContents();
    assert.ok(!names.includes(foreign.tag_key), 'selected device table contains foreign mapping');
    assert.equal(names.length, 8, 'selected device must have exactly eight own measurement rows');
    if (device.id === devices[0].id) assert.ok(text.includes('9007199254740993'), 'live uint64 lost precision');
    observations.push({ device_id: device.id, own_tag: own.tag_key, foreign_tag_absent: foreign.tag_key, expected, text });
  }
  await page.screenshot({ path: join(evidenceDir, `${runId}-runtime.png`), fullPage: true });
  return observations;
}

async function openBasic(page, base, deviceId) {
  await page.goto(`${base}/studio/v2`);
  await page.getByTestId('step-nav-button-4').click();
  await page.getByTestId('basic-recording-device').selectOption(deviceId);
  await page.getByTestId('basic-recording-panel').waitFor();
}

function effectSnapshot(run, effectKey) {
  const rows = readSQLiteRows(run.gateway_db, `SELECT effect_key,record_id,bucket_start,payload,payload_digest,group_revision,connector_id,connector_revision,table_name,state,claim_epoch FROM wg_delivery_outbox WHERE effect_key=${quoteSQLite(effectKey)}`);
  assert.equal(rows.length, 1);
  return rows[0];
}
function frozen(snapshot) {
  const { state, claim_epoch, ...identity } = snapshot;
  return identity;
}

async function recoverUI({ page, base, run, target, group, fault, repair, resolution }) {
  const targetWrite = (sql) => sh('sqlite3', ['-cmd', '.timeout 15000', target, sql]);
  targetWrite(fault);
  const view = await until(() => observeJSON(base, `/studio-v2/workspace/write-groups/${group.id}/delivery`),
    (value) => value.attention?.some((item) => item.state === (resolution === 'retry' ? 'blocked' : 'quarantined')));
  const attention = view.attention.find((item) => item.state === (resolution === 'retry' ? 'blocked' : 'quarantined'));
  const before = effectSnapshot(run, attention.effect_key);
  targetWrite(repair);
  await openBasic(page, base, group.basic_managed_device_id);
  const recovery = page.getByTestId('delivery-recovery');
  await recovery.waitFor({ timeout: 20_000 });
  const note = page.getByTestId(`delivery-reason-${attention.effect_key}`);
  await note.fill(`Owned ${resolution === 'retry' ? 'missing table' : 'poison row'} repaired; ${resolution} reviewed`);
  if (resolution === 'skip') {
    assert.equal(await page.getByTestId(`delivery-skip-${attention.effect_key}`).isDisabled(), true, 'skip needs explicit confirmation');
    await note.locator('xpath=ancestor::li').getByRole('checkbox').check();
  }
  const response = page.waitForResponse((reply) => reply.url().endsWith(`/write-groups/${group.id}/delivery/resolve`) && reply.request().method() === 'POST');
  await page.getByTestId(`delivery-${resolution}-${attention.effect_key}`).click();
  const reply = await response;
  assert.equal(reply.status(), 200);
  const request = reply.request().postDataJSON();
  const result = await reply.json();
  const after = await until(() => effectSnapshot(run, attention.effect_key), (item) => item.state === (resolution === 'retry' ? 'sql_committed' : 'operator_skipped'), 40_000);
  assert.deepEqual(frozen(after), frozen(before), 'operator resolution retargeted accepted payload');
  const audits = readSQLiteRows(run.gateway_db, `SELECT details FROM workspace_audit_history WHERE id=${quoteSQLite(request.decision_id)}`);
  assert.equal(audits.length, 1, 'resolution audit must persist exactly once');
  const sqlRows = readSQLiteRows(target, `SELECT record_id FROM ${identifier(group.destination.table_name)} WHERE record_id=${quoteSQLite(before.record_id)}`);
  assert.equal(sqlRows.length, resolution === 'retry' ? 1 : 0, 'skip must not pretend the row was delivered');
  await page.screenshot({ path: join(run.evidenceDir, `${run.run_id}-recovery-${resolution}.png`), fullPage: true });
  return { before, after, request, result, audit: JSON.parse(audits[0].details), destination_row_count: sqlRows.length };
}

export async function verifyBasicIntegrityLifecycle({ page, base, port, run, target, groups, savedDevices, tags, children, gatewayBinary, evidenceDir }) {
  run.evidenceDir = evidenceDir;
  const devices = ['Line A', 'Line B'].map((name) => savedDevices.find((device) => device.name === name));
  console.log('F integrity: actual runtime A215/B187/A215 switching');
  const runtime = await observeDeviceSwitch(page, base, devices, tags, evidenceDir, run.run_id);
  const group = groups.find((item) => item.basic_managed_device_id === devices[0].id);
  const table = identifier(group.destination.table_name);
  const parked = identifier(`${group.destination.table_name}_fault_saved`);
  console.log('F integrity: owned missing-table fault and Basic retry');
  const retry = await recoverUI({ page, base, run, target, group, resolution: 'retry',
    fault: `ALTER TABLE ${table} RENAME TO ${parked};`, repair: `ALTER TABLE ${parked} RENAME TO ${table};` });
  console.log('F integrity: owned poison fault and explicit Basic skip');
  const skip = await recoverUI({ page, base, run, target, group, resolution: 'skip',
    fault: `CREATE TRIGGER f_owned_poison BEFORE INSERT ON ${table} BEGIN SELECT RAISE(ABORT,'rejected row'); END;`, repair: 'DROP TRIGGER f_owned_poison;' });
  // The next complete bucket must still be delivered after the skipped head.
  const subsequent = await until(() => readSQLiteRows(target, `SELECT record_id,bucket_start FROM ${table} WHERE bucket_start > ${quoteSQLite(skip.before.bucket_start)} ORDER BY bucket_start`), (rows) => rows.length > 0);
  console.log('F integrity: disable, drain, normal gateway restart');
  await openFreshAdvanced(page);
  await page.getByTestId(`group-open-${group.id}`).click();
  await page.getByTestId('group-disable').click();
  await until(() => observeJSON(base, '/studio-v2/workspace/write-groups'), (value) => value.groups.some((item) => item.id === group.id && item.status === 'disabled'), 20_000);
  // Let the accepted final bucket drain, then stop only this runner's gateway.
  await until(() => observeJSON(base, `/studio-v2/workspace/write-groups/${group.id}/delivery`), (value) => value.intake.state === 'not_running' && value.stages.collecting === 0 && value.stages.queued === 0 && value.stages.retrying === 0);
  const beforeRestart = readSQLiteRows(target, `SELECT record_id,bucket_start FROM ${table} ORDER BY bucket_start`);
  const oldGateway = children.find((child) => child.fHarnessName === 'gateway');
  const stopped = await stopProcesses([oldGateway]);
  assert.ok(stopped.every(isFreshCleanExit), 'owned gateway did not shut down normally');
  const gateway = startFreshGateway({ work: run.work, binary: gatewayBinary, port });
  children.push(gateway);
  await waitForFreshGateway(base);
  await openBasic(page, base, devices[0].id);
  const afterRestart = await observeJSON(base, '/studio-v2/workspace/write-groups');
  assert.equal(afterRestart.groups.find((item) => item.id === group.id).status, 'disabled');
  const enabledGroup = groups.find((item) => item.id !== group.id);
  const healthyCount = readSQLiteRows(target, `SELECT COUNT(*) AS n FROM ${identifier(enabledGroup.destination.table_name)}`)[0].n;
  await until(() => readSQLiteRows(target, `SELECT COUNT(*) AS n FROM ${identifier(enabledGroup.destination.table_name)}`)[0].n, (n) => n > healthyCount, 90_000);
  const disabledRows = readSQLiteRows(target, `SELECT record_id,bucket_start FROM ${table} ORDER BY bucket_start`);
  assert.deepEqual(disabledRows, beforeRestart, 'disabled group wrote again after restart while another group delivered');
  await page.screenshot({ path: join(evidenceDir, `${run.run_id}-disabled-restart.png`), fullPage: true });
  return { runtime, missing_table_retry: retry, poison_explicit_skip: skip, subsequent_rows: subsequent,
    disable_restart: { stop: stopped, disabled_group_id: group.id, before: beforeRestart, after: disabledRows, healthy_group_progressed: true } };
}
