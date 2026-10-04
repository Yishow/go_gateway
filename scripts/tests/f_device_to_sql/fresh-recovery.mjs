import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';
import {
  assertAcquisitionProvenance,
  assertNeighborUnchanged,
  assertNoOwnedRows,
  assertOneEffect,
  assertOperationIdentity,
  assertTargetRowsUnchanged,
  assertTestWriteOutcome,
  sanitizeEvidence,
} from './fresh-recovery-assertions.mjs';
import { dropOneResponseAfterEffect, observeResponse } from './fresh-recovery-response.mjs';

export const FRESH_RECOVERY_BASELINE_CONTRACT = Object.freeze({
  createRun: 'createRun({ kind, runId, gatewayPort, simulatorPorts }) => run context',
  setupViaUI: 'run.setupViaUI({ intervalSeconds }) => { groupId, target, neighbor }',
  openGroup: 'run.openGroup(groupId) navigates the real UI to the saved group editor',
  produceNeighborViaUI: 'run.produceNeighborViaUI(groupId) creates the non-target owned neighbor through the real UI',
  observeAPI: 'run.observeAPI({ method: "GET", path }) performs read-only API observation',
  targetSnapshot: 'run.targetSnapshot() performs independent read-only SQL observation',
  neighborSnapshot: 'run.neighborSnapshot() performs independent read-only neighbor observation',
  gatewaySnapshot: 'run.gatewaySnapshot() performs independent config DB observation',
  exerciseDelayedResend: 'run.exerciseDelayedResend(groupId) uses only owned simulator/target controls and returns SQL facts',
  armStartPartial: 'run.armStartPartial() creates a real owned activation checkpoint failure after save/apply',
  releaseStartPartial: 'run.releaseStartPartial() removes the owned activation checkpoint fault before retry',
  restart: 'run.restart() restarts the same owned gateway/config DB',
  cleanup: 'run.cleanup() removes only this run resources and verifies absence',
});

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), '../../..');
const EVIDENCE_ROOT = resolve(process.env.F_EVIDENCE_ROOT ?? join(ROOT, 'docs/plans/studio-v2-flow-completion/evidence-f'));
const RECOVERY_CASES = Object.freeze(['schema_lost_response', 'start_lost_response', 'delayed_resend']);

function dataOf(response) {
  return response?.body?.data ?? response?.data ?? response;
}

function requireMethod(run, name) {
  assert.equal(typeof run?.[name], 'function', `fresh baseline adapter must export ${name}()`);
  return run[name].bind(run);
}

function writeEvidence({ run, kind, runId, groupId, results, cleanup, failure }) {
  const evidencePath = join(EVIDENCE_ROOT, `fresh-recovery-${kind}-${runId}.json`);
  mkdirSync(dirname(evidencePath), { recursive: true });
  let binaryHash;
  if (run?.gatewayBinary) {
    try { binaryHash = createHash('sha256').update(readFileSync(run.gatewayBinary)).digest('hex'); } catch { /* evidence still records the failure */ }
  }
  const witness = sanitizeEvidence({
    run_id: runId, kind, group_id: groupId, results, cleanup, failure, harness_started_at: run?.harness_started_at, harness_ended_at: new Date().toISOString(),
    binary_sha256: binaryHash,
    passed: !failure && Boolean(cleanup?.passed) && RECOVERY_CASES.every((name) => results[name]) && Boolean(results.main_ui),
  });
  writeFileSync(evidencePath, `${JSON.stringify(witness)}\n`);
}

async function openGroupUI(run, groupId) {
  await requireMethod(run, 'openGroup')(groupId);
}

async function reopenDatabaseStep(page) {
  await page.getByTestId('workbench-v2-root').waitFor({ timeout: 20_000 });
  const stepFour = page.getByTestId('step-nav-button-4');
  if (!(await stepFour.isEnabled().catch(() => false))) {
    await page.getByTestId('btn-continue-step1').waitFor({ timeout: 20_000 });
    await page.getByTestId('btn-continue-step1').click();
    await page.getByTestId('step2-rule-container').waitFor({ timeout: 20_000 });
    const stepTwoContinue = page.getByTestId('continue-step3-btn');
    await stepTwoContinue.waitFor({ timeout: 20_000 });
    await page.waitForFunction(() => !document.querySelector('[data-testid="continue-step3-btn"]')?.hasAttribute('disabled'), null, { timeout: 20_000 });
    await stepTwoContinue.click();
    await page.getByTestId('step3-mapping-container').waitFor({ timeout: 20_000 });
    await page.getByTestId('btn-continue').click();
  }
  await stepFour.waitFor({ timeout: 20_000 });
  await page.waitForFunction(() => !document.querySelector('[data-testid="step-nav-button-4"]')?.hasAttribute('disabled'), null, { timeout: 20_000 });
  await stepFour.click();
  await page.getByTestId('basic-recording-panel').waitFor({ timeout: 20_000 });
}

async function runSchemaLostReply(run, groupId) {
  const page = run.page;
  await openGroupUI(run, groupId);
  const previewResponse = await observeResponse(page, {
    path: `/workspace/write-groups/${groupId}/schema-preview`,
  }, () => page.getByTestId('group-schema-preview').click());
  const preview = dataOf(previewResponse);
  assert.ok(preview?.operation_id && preview?.token, 'schema preview must issue operation and token');
  const before = await run.targetSnapshot();
  const drop = await dropOneResponseAfterEffect(page, {
    path: `/workspace/write-groups/${groupId}/schema-apply`,
  });
  let result;
  try {
    await page.getByTestId('group-schema-apply').click();
    const dropped = await drop.done;
    assert.equal(dropped.status, 200, 'schema apply server effect did not return HTTP 200 before drop');
    const serverOperation = dataOf(dropped);
    assertOperationIdentity(preview, serverOperation, 'schema lost response');
    await page.getByTestId('group-schema-operation-pending').waitFor({ timeout: 15000 });
    const checkDrop = await dropOneResponseAfterEffect(page, {
      method: 'GET',
      path: `/workspace/database-operations/${preview.operation_id}`,
    });
    let errorUI;
    try {
      await page.getByTestId('group-schema-check').click();
      const droppedCheck = await checkDrop.done;
      assert.equal(droppedCheck.status, 200, 'schema operation check did not return HTTP 200 before drop');
      await page.getByTestId('group-schema-error').waitFor({ timeout: 15000 });
      const errorText = await page.getByTestId('group-schema-error').innerText();
      assert.ok(errorText.trim(), 'schema operation check did not render a safe error');
      errorUI = await requireMethod(run, 'captureResponsiveEvidence')('error');
    } finally {
      await checkDrop.close();
    }
    const operation = await run.observeAPI({ method: 'GET', path: `/workspace/database-operations/${preview.operation_id}` });
    const operationData = dataOf(operation);
    assertOperationIdentity(serverOperation, operationData, 'schema operation query');
    assert.equal(operationData.status, 'succeeded', 'schema operation query did not prove success');
    assert.ok(Number(operationData.executed_statements) > 0, 'schema operation did not report executed statements');
    const check = page.getByTestId('group-schema-check');
    if (await check.isVisible().catch(() => false)) await check.click();
    const after = await run.targetSnapshot();
    assert.ok(Number(after.managed_table_count) > Number(before.managed_table_count), 'schema response drop did not leave tables behind');
    result = { preview: { operation_id: preview.operation_id }, dropped, operation: operationData, target_before: before, target_after: after, error_ui: errorUI };
  } finally {
    await drop.close();
  }
  // The schema proof is persisted by the server, but the Basic panel can still
  // hold its pre-confirm group projection after the lost reply/check flow.
  // Reload the real UI before constructing the recording-start request so its
  // request carries the confirmed schema revision and digest.
  await page.reload();
  await reopenDatabaseStep(page);
  return result;
}

async function runStartLostReply(run, groupId) {
  const page = run.page;
  await page.getByTestId('basic-recording-panel').waitFor({ timeout: 15_000 });
  const armStartPartial = requireMethod(run, 'armStartPartial');
  await armStartPartial();
  const startDrop = await dropOneResponseAfterEffect(page, {
    path: '/workspace/recording-start',
  });
  try {
    await page.getByTestId('basic-recording-start').click();
    const dropped = await startDrop.done;
    assert.equal(dropped.status, 200, 'recording start server effect did not return HTTP 200 before drop');
    const first = dataOf(dropped);
    const initialRequest = dropped.request;
    assert.equal(first.action, 'recording_start', 'dropped response is not recording_start');
    assert.ok(first.operation_id, 'recording start response has no operation identity');
    assert.equal(first.status, 'partial', 'start recovery must begin from a real partial operation');
    const operationID = first.operation_id;
    const sourceIdentity = {
      request_id: initialRequest?.request_id,
      workspace_id: initialRequest?.workspace_id,
      expected_workspace_revision: initialRequest?.expected_workspace_revision,
      groups: initialRequest?.groups,
    };
    assert.ok(sourceIdentity.request_id && sourceIdentity.workspace_id && sourceIdentity.expected_workspace_revision && sourceIdentity.groups?.length,
      'recording start response drop lost the original source identity');
    await requireMethod(run, 'restart')();
    await page.reload();
    await reopenDatabaseStep(page);
    const queried = dataOf(await run.observeAPI({ method: 'GET', path: `/workspace/recording-start/operations/${operationID}` }));
    assertOperationIdentity(first, queried, 'start operation query after reload');
    assert.equal(queried.status, 'partial', 'reload must retain the partial operation');
    const appliedRevision = queried.groups?.[0]?.applied_revision;
    const partialResponse = await observeResponse(page, { path: '/workspace/recording-start' }, () => page.getByTestId('basic-recording-start').click());
    const partial = dataOf(partialResponse);
    assertOperationIdentity(first, partial, 'same-request partial retry');
    assert.equal(partialResponse.request?.request_id, initialRequest?.request_id, 'partial retry changed request identity');
    assert.deepEqual(partialResponse.request?.groups, sourceIdentity.groups, 'partial retry changed source group identity');
    assert.equal(partial.status, 'partial', 'armed activation fault did not return a real partial retry');
    await page.waitForFunction(() => {
      const text = document.querySelector('[data-testid="basic-recording-status"]')?.textContent ?? '';
      return /partial|partly|部分/i.test(text);
    }, null, { timeout: 15000 });
    await page.getByTestId('basic-recording-retry').waitFor({ timeout: 15000 });
    const partialText = await page.getByTestId('basic-recording-status').innerText();
    assert.match(partialText, /partial|partly|部分/i, 'partial status was not rendered in the Basic panel');
    const retryText = await page.getByTestId('basic-recording-retry').innerText();
    assert.match(retryText, /retry|重送/i, 'same-request retry control was not rendered');
    const partialUI = await requireMethod(run, 'captureResponsiveEvidence')('partial');
    await requireMethod(run, 'releaseStartPartial')();
    const retryResponse = await observeResponse(page, { path: '/workspace/recording-start' }, () => page.getByTestId('basic-recording-retry').click());
    assert.equal(retryResponse.request?.request_id, initialRequest?.request_id, 'lost-response retry changed request identity');
    assert.deepEqual(retryResponse.request?.groups, sourceIdentity.groups, 'lost-response retry changed source group identity');
    const replay = dataOf(retryResponse);
    assertOperationIdentity(first, replay, 'start retry');
    assert.equal(replay.operation_id, operationID, 'start retry created another operation');
    assert.equal(replay.status, 'succeeded', 'start retry did not finish the original operation');
    assert.equal(replay.groups?.[0]?.applied_revision, appliedRevision, 'start retry changed applied revision');
    const effects = (await run.gatewaySnapshot()).effects?.filter((effect) => effect.operation_id === operationID) ?? [];
    assertOneEffect(effects, 'start retry activation');
    return { dropped, first, queried, partial_response: partialResponse, replay, source_identity: sourceIdentity, effects, retry_request: retryResponse.request, partial_ui: partialUI };
  } finally {
    await startDrop.close();
  }
}

async function runDelayedResend(run, groupId) {
  const exercise = requireMethod(run, 'exerciseDelayedResend');
  const result = await exercise(groupId);
  assert.ok(result?.expected && result?.row, 'delayed resend adapter returned incomplete SQL evidence');
  assertAcquisitionProvenance(result.row, result.expected);
  assert.equal(result.row.record_id, result.durable_link?.record_id, 'delayed resend row was not linked to committed outbox record');
  assert.equal(result.row.group_id, result.durable_link?.group_id, 'delayed resend row changed durable group identity');
  assert.ok(result.durable_link?.group_revision, 'delayed resend has no durable group revision evidence');
  assertOneEffect(result.effects, 'delayed resend');
  const evidence = run.page.getByTestId('basic-recording-evidence');
  await evidence.waitFor({ timeout: 30_000 });
  await run.page.getByTestId('basic-recording-committed-effect').waitFor({ timeout: 30_000 });
  assert.equal(await evidence.locator('p.text-slate-500').count(), 0, 'current SQL committed target still shows the first-bucket pending hint');
  result.pending_hint = { present: false, evidence_text: await evidence.innerText() };
  return result;
}

async function runTestWriteCleanup(run, groupId) {
  const page = run.page;
  await requireMethod(run, 'produceNeighborViaUI')(groupId);
  const before = await run.neighborSnapshot();
  await openGroupUI(run, groupId);
  const targetBefore = await run.targetSnapshot();
  const previewResponse = await observeResponse(page, {
    path: `/workspace/write-groups/${groupId}/test-write-preview`,
  }, () => page.getByTestId('group-test-write-preview').click());
  const preview = dataOf(previewResponse);
  assert.ok(preview?.operation_id,
    `test-write preview must expose operation identity: status=${previewResponse.status} request=${JSON.stringify(previewResponse.request)} body=${JSON.stringify(previewResponse.body).slice(0, 4000)}`);
  assert.ok(preview?.owner_value, 'test-write preview must expose the operation-owned cleanup marker');
  await page.getByTestId('group-test-write-preview-card').waitFor({ timeout: 15000 });
  const confirmResponse = await observeResponse(page, {
    path: `/workspace/write-groups/${groupId}/test-write`,
  }, () => page.getByTestId('group-test-write-confirm').click());
  const operation = dataOf(confirmResponse);
  assertTestWriteOutcome(operation);
  await page.getByTestId('group-test-write-outcome').waitFor({ timeout: 30000 });
  const after = await run.neighborSnapshot();
  assertNeighborUnchanged(before, after);
  const targetAfter = await run.targetSnapshot();
  assertTargetRowsUnchanged(targetBefore.rows ?? targetBefore, targetAfter.rows ?? targetAfter, preview.owner_value, preview.owner_column);
  assertNoOwnedRows(targetAfter.rows ?? targetAfter, preview.owner_value, preview.owner_column);
  return {
    preview,
    operation,
    neighbor_before: before,
    neighbor_after: after,
    target_neighbor_baseline: run.targetNeighborBaseline,
    target_before: targetBefore,
    target_after: targetAfter,
  };
}

export async function runFreshRecovery({ createRun, kind, runId, gatewayPort = 3381, simulatorPorts = [15052, 15053] }) {
  assert.equal(typeof createRun, 'function', 'fresh baseline module must export createRun()');
  assert.match(String(kind), /^[A-Za-z0-9][A-Za-z0-9_-]{0,40}$/, 'fresh recovery kind is unsafe');
  assert.match(String(runId), /^[A-Za-z0-9][A-Za-z0-9_-]{0,80}$/, 'fresh recovery run id is unsafe');
  let run;
  try {
    run = await createRun({ kind, runId, gatewayPort, simulatorPorts });
  } catch (error) {
    const failure = { message: error instanceof Error ? error.message : String(error) };
    writeEvidence({ kind, runId, results: {}, cleanup: { passed: false }, failure });
    throw error;
  }
  let groupId;
  const results = {};
  let cleanup;
  let failure;
  try {
    const setup = await requireMethod(run, 'setupViaUI')({ intervalSeconds: 10 });
    groupId = setup.groupId;
    assert.ok(groupId, 'fresh UI setup did not return a saved group id');
    results.schema_lost_response = await runSchemaLostReply(run, groupId);
    results.start_lost_response = await runStartLostReply(run, groupId);
    results.delayed_resend = await runDelayedResend(run, groupId);
    results.main_ui = await requireMethod(run, 'captureResponsiveEvidence')('main');
    return { kind, run_id: runId, group_id: groupId, results };
  } catch (error) {
    run.keepWork = true;
    failure = { message: error instanceof Error ? error.message : String(error) };
    try { results.error_ui = await requireMethod(run, 'captureResponsiveEvidence')('error'); }
    catch (captureError) { failure.ui_capture = captureError instanceof Error ? captureError.message : String(captureError); }
    throw error;
  } finally {
    try {
      cleanup = await requireMethod(run, 'cleanup')();
    } catch (error) {
      cleanup = { passed: false, error: error instanceof Error ? error.message : String(error) };
    }
    writeEvidence({ run, kind, runId, groupId, results, cleanup, failure });
  }
}

async function main() {
  const moduleName = process.env.F_FRESH_BASELINE_MODULE ?? './scripts/tests/f_device_to_sql/fresh-recovery-adapter.mjs';
  const kind = process.env.F_SQL_KIND ?? 'sqlite';
  const runId = process.env.F_RUN_ID ?? `${kind}-recovery-${Date.now()}`;
  const moduleURL = moduleName.startsWith('file:') ? moduleName : pathToFileURL(resolve(ROOT, moduleName)).href;
  const baseline = await import(moduleURL);
  await runFreshRecovery({
    createRun: baseline.createRun,
    kind,
    runId,
    gatewayPort: Number(process.env.GW_PORT ?? 3381),
    simulatorPorts: [Number(process.env.F_SIM_PORT_A ?? 15052), Number(process.env.F_SIM_PORT_B ?? 15053)],
  });
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  main().catch((error) => {
    console.error(`fresh recovery failed: ${error.message}`);
    process.exitCode = 1;
  });
}
