// Tight fresh UI diagnosis for a test-write preview that returns 200 while
// the embedded UI shows a generic error. It stops before runtime collection.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { existsSync, rmSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, assertFreshPortFree,
  binaryIdentity, closeFreshUI, configureFreshSQLite, makeFreshRun, observeJSON, openFreshAdvanced,
  openFreshUI, persistFreshMappings, previewAndApplyFreshSchema, setupFreshMultiDeviceSources,
  startFreshGateway, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import { FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES } from './fresh-recording-verification.mjs';
import { preflightSQLiteDestination } from './destination-preflight.mjs';

const runId = process.env.F_RUN_ID ?? `preview-diagnosis-${Date.now()}`;
const port = Number(process.env.GW_PORT ?? 3380);
const simPortA = Number(process.env.F_SIM_PORT_A ?? 15050);
const simPortB = Number(process.env.F_SIM_PORT_B ?? 15051);
const gatewayBinary = process.env.F_GATEWAY_BINARY ?? DEFAULT_GATEWAY_BINARY;
const run = makeFreshRun(runId);
const base = `http://127.0.0.1:${port}`;
const target = run.destination_db;
const evidencePath = join(DEFAULT_EVIDENCE_DIR, `${runId}.json`);
const children = [];
let browser; let page; let ui; let result; let failure; let stage = 'preflight';

const qid = (value) => `"${String(value).replaceAll('"', '""')}"`;

async function groupsAt() { return observeJSON(base, '/studio-v2/workspace/write-groups'); }

async function waitGroup(id, predicate, timeoutMs = 30_000) {
  const deadline = Date.now() + timeoutMs; let current;
  while (Date.now() < deadline) {
    current = (await groupsAt()).groups?.find((group) => group.id === id);
    if (current && predicate(current)) return current;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error(`group ${id} did not reach expected state`);
}

async function createEntityGroup(tagsA, tagsB) {
  await openFreshAdvanced(page);
  await page.getByTestId('group-create').click();
  await page.getByTestId('group-name').fill('F preview diagnosis A+B');
  await page.getByTestId('group-storage-managed').check();
  const advanced = page.getByTestId('group-advanced');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-entity-column').fill('entity_id');
  await page.getByTestId('group-provenance-column').fill('provenance');
  await page.getByTestId('group-bulk-include').click();
  for (const tag of [...tagsA, ...tagsB]) {
    const row = page.getByTestId(`group-member-row-${tag.persisted_point_id}`);
    await row.locator('input[type="text"]').first().fill(tag.device_id === tagsA[0].device_id ? 'A' : 'B');
  }
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const groups = await groupsAt();
  const matches = groups.groups?.filter((group) => group.name === 'F preview diagnosis A+B') ?? [];
  assert.equal(matches.length, 1);
  const managed = matches[0];
  assert.equal(managed.members.length, 18);
  const distinctColumn = managed.members.find((member) => member.device_id === tagsB[0].device_id && member.point_id === tagsB[0].persisted_point_id)?.target_column;
  assert.ok(distinctColumn);
  await previewAndApplyFreshSchema(page);
  await page.getByTestId('group-storage-custom').check();
  const aColumn = page.getByTestId(`group-member-column-${tagsA[0].persisted_point_id}`);
  const bColumn = page.getByTestId(`group-member-column-${tagsB[0].persisted_point_id}`);
  await aColumn.waitFor({ timeout: 20_000 });
  const sharedColumn = await aColumn.inputValue();
  assert.ok(sharedColumn);
  await bColumn.selectOption(sharedColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const savedGroups = await groupsAt();
  const saved = savedGroups.groups?.find((group) => group.id === managed.id);
  assert.ok(saved);
  assert.equal(saved.destination.storage_strategy, 'custom');
  return { group: saved, sharedColumn, distinctColumn };
}

function safePreview(data) {
  return {
    keys: data && typeof data === 'object' ? Object.keys(data).sort() : [],
    action: data?.action, operation_id: data?.operation_id, group_id: data?.group_id,
    group_revision: data?.group_revision, owner_column: data?.owner_column,
    owner_value_prefix: typeof data?.owner_value === 'string' ? data.owner_value.slice(0, 9) : undefined,
    cleanup: data?.cleanup, cleanup_length: typeof data?.cleanup === 'string' ? data.cleanup.length : undefined,
    field_lengths: Object.fromEntries(['token', 'operation_id', 'group_id', 'group_revision', 'owner_column', 'owner_value', 'cleanup'].map((key) => [key, typeof data?.[key] === 'string' ? data[key].length : null])),
    token_sha256: typeof data?.token === 'string' ? createHash('sha256').update(data.token).digest('hex') : undefined,
    target: data?.target && {
      keys: Object.keys(data.target).sort(), connector_id: data.target.connector_id,
      dialect: data.target.dialect, table: data.target.table,
      database_present: typeof data.target.database === 'string', schema_present: typeof data.target.schema === 'string',
    },
    values: Array.isArray(data?.values) ? data.values.map((entry) => ({
      column: entry?.column, type: entry?.type, value: entry?.value,
    })) : undefined,
  };
}

try {
  assert.equal(new Set([port, simPortA, simPortB]).size, 3, 'fresh ports must be distinct');
  const initial = preflightSQLiteDestination({ path: target, work: run.work });
  assert.equal(initial.recording_table_present, false); assert.equal(existsSync(target), false);
  await Promise.all([assertFreshPortFree(port), assertFreshPortFree(simPortA), assertFreshPortFree(simPortB)]);
  children.push(startFreshSimulator({ work: run.work, name: 'sim-a', port: simPortA,
    registers: Object.fromEntries(FRESH_REGISTERS.map((value, index) => [index, value])) }));
  children.push(startFreshSimulator({ work: run.work, name: 'sim-b', port: simPortB,
    registers: Object.fromEntries(FRESH_REGISTERS.map((value, index) => [index, value])) }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port }));
  await waitForFreshGateway(base); ui = await openFreshUI({ port }); ({ browser, page } = ui);
  stage = 'devices-and-mappings';
  await setupFreshMultiDeviceSources(page, [
    { name: 'Line A', port: simPortA, prefix: 'A' }, { name: 'Line B', port: simPortB, prefix: 'B' },
  ], FRESH_RULES);
  const targetTypes = new Map(FRESH_RULES.map((rule, index) => [String(rule.start), FRESH_TARGET_TYPES[index]]));
  const tags = await persistFreshMappings(page, base, {
    tagPrefix: 'preview', targetTypes: [...FRESH_TARGET_TYPES, ...FRESH_TARGET_TYPES],
    targetTypeFor: ({ address }) => targetTypes.get(address),
  });
  const devicesData = await observeJSON(base, '/studio-v2/workspace/devices');
  const devices = Array.isArray(devicesData) ? devicesData : devicesData.devices ?? [];
  const deviceA = devices.find((device) => device.name === 'Line A');
  const deviceB = devices.find((device) => device.name === 'Line B');
  assert.ok(deviceA?.id && deviceB?.id && deviceA.id !== deviceB.id);
  const tagsA = tags.filter((tag) => tag.device_id === deviceA.id);
  const tagsB = tags.filter((tag) => tag.device_id === deviceB.id);
  assert.equal(tagsA.length, 9); assert.equal(tagsB.length, 9);
  stage = 'destination';
  const destination = await configureFreshSQLite(page, base, { path: target, table: 'f_connector_placeholder' });
  assert.equal(preflightSQLiteDestination({ path: target, work: run.work }).exists, false);
  stage = 'schema-and-preview';
  const entity = await createEntityGroup(tagsA, tagsB);
  const applyButton = page.getByTestId('group-apply');
  await page.waitForFunction(() => !document.querySelector('[data-testid="group-apply"]')?.hasAttribute('disabled'), null, { timeout: 30_000 });
  await applyButton.click();
  const applied = await waitGroup(entity.group.id, (group) => Boolean(group.applied_revision) && group.status !== 'draft');
  const uniqueColumn = page.getByTestId(`group-member-column-${tagsB[0].persisted_point_id}`);
  await uniqueColumn.selectOption(entity.distinctColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  const draft = await waitGroup(entity.group.id, (group) => group.revision !== applied.revision && group.applied_revision === applied.applied_revision);
  await page.getByTestId('group-editor-close').click();
  await openFreshAdvanced(page);
  await page.getByTestId(`group-open-${draft.id}`).click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
  const previewButton = page.getByTestId('group-test-write-preview');
  const enabled = await previewButton.isEnabled();
  const previewResponse = await Promise.all([
    page.waitForResponse((response) => response.url().endsWith(`/write-groups/${entity.group.id}/test-write-preview`) && response.request().method() === 'POST'),
    previewButton.click(),
  ]).then(([response]) => response);
  const body = await previewResponse.json();
  const data = body?.data;
  await page.waitForTimeout(500);
  const card = page.getByTestId('group-test-write-preview-card');
  const cardVisible = await card.isVisible().catch(() => false);
  const uiError = await page.getByTestId('group-test-write-error').innerText().catch(() => '');
  const panelText = await page.getByTestId('group-test-write').innerText().catch(() => '');
  result = {
    run_id: runId, kind: 'sqlite-test-write-preview-diagnosis', passed: cardVisible,
    stage, destination, initial_preflight: initial, binary: { gateway: binaryIdentity(gatewayBinary), simulator: binaryIdentity(DEFAULT_SIMULATOR_BINARY) },
    preview_http_status: previewResponse.status(), response_success: body?.success === true, preview: safePreview(data),
    ui: { preview_enabled: enabled, card_visible: cardVisible, error: uiError, panel: panelText.replace(/\s+/g, ' ').slice(-1200) },
    page_errors: ui.pageErrors, api_traffic: ui.apiTraffic, retained_work: run.work,
  };
  console.log(`${cardVisible ? 'PASS' : 'DIAG'} ${runId}: preview status ${previewResponse.status()}, card=${cardVisible}`);
} catch (error) {
  failure = error;
  result = { run_id: runId, passed: false, failed_stage: stage, error: String(error?.message ?? error).split('\n')[0], retained_work: run.work };
  console.log(`ERROR ${runId} ${stage}: ${result.error}`);
} finally {
  const cleanup = await stopProcesses(children).catch((error) => [{ state: 'error', error: String(error) }]);
  let browserCleanup = { state: 'closed' };
  try { await closeFreshUI(browser); } catch (error) { browserCleanup = { state: 'error', error: error.message.split('\n')[0] }; }
  const processOK = cleanup.every((status) => status?.state === 'exited' && status.exit_code === 0 && !status.signal && status.kill_sent !== true && !status.spawn_error && !status.stop_error);
  const payload = result ?? { run_id: runId, passed: false, error: 'no result', retained_work: run.work };
  const namespace = { path: run.work, absent: false, state: 'blocked' };
  if (processOK && browserCleanup.state === 'closed') {
    try { rmSync(run.work, { recursive: true, force: false }); namespace.state = 'removed'; namespace.absent = !existsSync(run.work); }
    catch (error) { namespace.error = error.message.split('\n')[0]; }
  }
  payload.cleanup = { processes: cleanup, browser: browserCleanup, namespace };
  if (!processOK || browserCleanup.state !== 'closed' || !namespace.absent) { payload.passed = false; payload.error = `${payload.error ?? ''}; cleanup not proven`; }
  writeFileSync(evidencePath, `${JSON.stringify(payload, null, 2)}\n`);
}
if (failure) throw failure;
