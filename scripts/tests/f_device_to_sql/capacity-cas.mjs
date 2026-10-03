// F2.3 capacity/CAS witness. This module is intentionally export-only: the
// caller owns the real gateway, loopback devices and disposable SQL target.
// Every mutation below goes through the production Studio V2 API or UI.
import assert from 'node:assert/strict';
import { BASE, boundedFetch, openBrowser } from './lib.mjs';
import { safePreview } from './test-write-lib.mjs';

const writeGroupPath = (id) => `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`;
const writeGroupsPath = '/studio-v2/workspace/write-groups';

function clone(value) {
  return value === undefined ? undefined : JSON.parse(JSON.stringify(value));
}

function responseErrorCode(reply) {
  return reply.body?.error?.code;
}

export function responseWitness(reply) {
  const request = clone(reply.request);
  const body = clone(reply.body);
  if (typeof request?.body?.token === 'string') request.body = safePreview(request.body);
  if (typeof body?.data?.token === 'string') body.data = safePreview(body.data);
  return {
    request,
    response: { status: reply.status, ok: reply.ok, body },
  };
}

export async function request(method, path, body) {
  const requestBody = body === undefined ? undefined : clone(body);
  const response = await boundedFetch(`${BASE}/api/v1/datalink${path}`, {
    method,
    headers: requestBody === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: requestBody === undefined ? undefined : JSON.stringify(requestBody),
  });
  const text = await response.text();
  let parsed;
  try {
    parsed = text === '' ? null : JSON.parse(text);
  } catch {
    parsed = { raw: text };
  }
  return {
    request: { method, path, body: requestBody ?? null },
    status: response.status,
    ok: response.ok,
    body: parsed,
  };
}

async function readGroups() {
  const reply = await request('GET', writeGroupsPath);
  assert.equal(reply.status, 200, `GET ${writeGroupsPath}: ${JSON.stringify(reply.body)}`);
  assert.equal(reply.body?.success, true, `GET ${writeGroupsPath}: ${JSON.stringify(reply.body)}`);
  assert.ok(reply.body.data?.workspace_id, 'write-group list must identify the workspace');
  return { reply, data: reply.body.data };
}

function findGroup(data, id) {
  const group = data.groups?.find((entry) => entry.id === id);
  assert.ok(group, `write group ${id} is present in the canonical workspace list`);
  return group;
}

function groupFromInput(groups, key, data) {
  const requested = groups?.[key];
  if (requested?.id) return findGroup(data, requested.id);
  const first = data.groups?.find((entry) => entry.status !== 'deleted');
  assert.ok(first, 'capacity/CAS requires one live write group');
  return first;
}

function mergeGroup(group, patch) {
  const next = clone(group);
  for (const [key, value] of Object.entries(patch)) {
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      next[key] = { ...(next[key] ?? {}), ...clone(value) };
    } else {
      next[key] = clone(value);
    }
  }
  return next;
}

function casUpdate(data, group, patch) {
  return {
    workspace_id: data.workspace_id,
    expected_workspace_revision: data.workspace_revision,
    expected_group_revision: group.revision,
    expected_connector_revision: group.destination.connector_revision,
    group: mergeGroup(group, patch),
  };
}

function applyRequest(data, group) {
  return {
    workspace_id: data.workspace_id,
    expected_workspace_revision: data.workspace_revision,
    expected_group_revision: group.revision,
    expected_connector_revision: group.destination.connector_revision,
  };
}

function record(e, name, details) {
  if (typeof e?.check === 'function') e.check(name, details);
  return { passed: true, ...details };
}

function sqlSnapshot(e) {
  assert.ok(e?.targets && typeof e.targets.rows === 'function', 'capacity/CAS requires owned SQL target readers');
  const targets = {};
  for (const name of e.targetNames ?? ['a', 'b']) targets[name] = clone(e.targets.rows(name));
  return {
    targets,
    local: {
      outbox: typeof e.outbox === 'function' ? clone(e.outbox()) : null,
      samples: typeof e.samples === 'function' ? clone(e.samples()) : null,
    },
  };
}

function deliveryContract(value) {
  if (!value) return null;
  return {
    intake: clone(value.intake),
    backlog: clone(value.backlog),
  };
}

async function readDelivery(groupID) {
  const reply = await request('GET', `${writeGroupPath(groupID)}/delivery`);
  assert.equal(reply.status, 200, `GET delivery ${groupID}: ${JSON.stringify(reply.body)}`);
  assert.equal(reply.body?.success, true, `GET delivery ${groupID}: ${JSON.stringify(reply.body)}`);
  return { reply, data: reply.body.data };
}

async function openGroupEditor(page, id) {
  assert.ok(page, 'capacity/CAS UI witness requires the root browser page');
  await page.reload();
  await page.waitForTimeout(1000);
  const step4 = page.getByTestId('step-nav-button-4');
  if (await step4.isVisible()) await step4.click();
  const open = page.getByTestId(`group-open-${id}`);
  await open.waitFor({ timeout: 20000 });
  await open.click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20000 });
}

async function runUIConflictWitness(e, groupID, baselineGroup, page = e.page) {
  await openGroupEditor(page, groupID);
  const localValue = `${baselineGroup.name} UI local conflict`;
  await page.getByTestId('group-name').fill(localValue);

  const externalValue = `${baselineGroup.name} server canonical`;
  let actualRead;
  let externalReply;
  let delayedActualRead = false;
  const listURL = '**/api/v1/datalink/studio-v2/workspace/write-groups';
  await page.route(listURL, async (route) => {
    if (!delayedActualRead && route.request().method() === 'GET') {
      delayedActualRead = true;
      const actualResponse = await route.fetch();
      const actualBody = await actualResponse.json();
      assert.equal(actualResponse.status(), 200, `actual UI GET write-groups: ${JSON.stringify(actualBody)}`);
      assert.equal(actualBody?.success, true, `actual UI GET write-groups: ${JSON.stringify(actualBody)}`);
      const actualData = actualBody.data;
      const actualGroup = findGroup(actualData, groupID);
      const externalRequest = casUpdate(actualData, actualGroup, { name: externalValue });
      externalReply = await request('PUT', writeGroupPath(groupID), externalRequest);
      assert.equal(externalReply.status, 200, `external canonical save: ${JSON.stringify(externalReply.body)}`);
      assert.equal(externalReply.body?.success, true, `external canonical save: ${JSON.stringify(externalReply.body)}`);
      actualRead = {
        request: { method: route.request().method(), path: new URL(route.request().url()).pathname.replace('/api/v1/datalink', ''), body: null },
        response: { status: actualResponse.status(), ok: actualResponse.ok(), body: clone(actualBody) },
      };
      await route.fulfill({ response: actualResponse });
      return;
    }
    await route.continue();
  });

  let uiConflictReply;
  try {
    const expectedURL = `${BASE}/api/v1/datalink${writeGroupPath(groupID)}`;
    const responsePromise = page.waitForResponse(
      (response) => response.url() === expectedURL && response.request().method() === 'PUT',
      { timeout: 20000 },
    );
    await page.getByTestId('group-save').click();
    const response = await responsePromise;
    const body = await response.json();
    uiConflictReply = {
      request: {
        method: response.request().method(),
        path: new URL(response.url()).pathname.replace('/api/v1/datalink', ''),
        body: response.request().postDataJSON(),
      },
      response: { status: response.status(), ok: response.ok(), body },
    };
  } finally {
    await page.unroute(listURL);
  }

  assert.equal(delayedActualRead, true, 'UI save fetched the actual backend list before the competing PUT');
  assert.ok(actualRead, 'actual UI list response is recorded');
  assert.ok(externalReply, 'competing canonical PUT is recorded');
  assert.equal(uiConflictReply.response.status, 409, JSON.stringify(uiConflictReply));
  assert.equal(uiConflictReply.response.body?.error?.code, 'revision_mismatch', JSON.stringify(uiConflictReply));
  await page.getByTestId('group-conflict').waitFor({ timeout: 10000 });
  assert.equal(await page.getByTestId('group-name').inputValue(), localValue, '409 keeps the local editor value');
  const localValueAfter409 = await page.getByTestId('group-name').inputValue();

  const reloadGET = page.waitForResponse(
    (response) => response.url() === `${BASE}/api/v1/datalink${writeGroupsPath}` && response.request().method() === 'GET',
    { timeout: 20000 },
  );
  await page.getByTestId('group-conflict-reload').click();
  const reloadResponse = await reloadGET;
  assert.equal(reloadResponse.status(), 200);
  const afterReload = await readGroups();
  const serverGroup = findGroup(afterReload.data, groupID);
  assert.equal(serverGroup.name, externalValue, 'reload reads the canonical value that won the external save');

  const repairValue = `${baselineGroup.name} UI repaired`;
  await page.getByTestId('group-name').fill(repairValue);
  const repairURL = `${BASE}/api/v1/datalink${writeGroupPath(groupID)}`;
  const repairResponsePromise = page.waitForResponse(
    (response) => response.url() === repairURL && response.request().method() === 'PUT',
    { timeout: 20000 },
  );
  await page.getByTestId('group-save').click();
  const repairResponse = await repairResponsePromise;
  const repairBody = await repairResponse.json();
  assert.equal(repairResponse.status(), 200, JSON.stringify(repairBody));
  assert.equal(repairBody?.success, true, JSON.stringify(repairBody));
  await page.getByTestId('group-saved-note').waitFor({ timeout: 15000 });

  return {
    actual_read: actualRead,
    external_save: responseWitness(externalReply),
    ui_save_409: uiConflictReply,
    local_value_after_409: localValueAfter409,
    reload: { request: { method: 'GET', path: writeGroupsPath, body: null }, response: { status: reloadResponse.status(), ok: reloadResponse.ok() } },
    repair_save: {
      request: { method: repairResponse.request().method(), path: new URL(repairResponse.url()).pathname.replace('/api/v1/datalink', ''), body: repairResponse.request().postDataJSON() },
      response: { status: repairResponse.status(), ok: repairResponse.ok(), body: repairBody },
    },
    repaired_value: repairValue,
  };
}

/**
 * Runs the F2.3 capacity/CAS cases against a running owned fixture.
 *
 * The caller must provide the `targets` readers and a group map such as
 * `{ a, b }`; an existing live `page` is reused when present, otherwise the
 * module opens and closes one UI client. No process, port, target schema or
 * SQL row is created here; the returned object is the raw witness for the
 * caller's compact evidence file.
 */
export async function runConcurrentEditCases(e, groups) {
  assert.ok(e, 'capacity/CAS context is required');
  const first = await readGroups();
  const baselineGroup = groupFromInput(groups, e.capacityGroupKey ?? 'a', first.data);
  const groupID = baselineGroup.id;
  const initialSQL = sqlSnapshot(e);
  const initialDelivery = await readDelivery(groupID);

  const racePatch = { name: `${baselineGroup.name} CAS winner` };
  const raceRequestA = casUpdate(first.data, baselineGroup, racePatch);
  const raceRequestB = clone(raceRequestA);
  const [raceA, raceB] = await Promise.all([
    request('PUT', writeGroupPath(groupID), raceRequestA),
    request('PUT', writeGroupPath(groupID), raceRequestB),
  ]);
  const raceReplies = [raceA, raceB];
  assert.deepEqual(raceReplies.map((reply) => reply.status).sort((a, b) => a - b), [200, 409], JSON.stringify(raceReplies.map(responseWitness)));
  assert.equal(raceReplies.filter((reply) => responseErrorCode(reply) === 'revision_mismatch').length, 1, JSON.stringify(raceReplies.map(responseWitness)));
  const casRace = record(e, 'CapacityCASConcurrentEdit', {
    requests: [raceRequestA, raceRequestB],
    responses: raceReplies.map(responseWitness),
  });

  const staleApplyRequest = applyRequest(first.data, baselineGroup);
  const staleApply = await request('POST', `${writeGroupPath(groupID)}/apply`, staleApplyRequest);
  assert.equal(staleApply.status, 409, JSON.stringify(responseWitness(staleApply)));
  assert.equal(responseErrorCode(staleApply), 'revision_mismatch', JSON.stringify(responseWitness(staleApply)));
  const staleApplyResult = record(e, 'CapacityCASStaleApply', {
    request: staleApplyRequest,
    response: responseWitness(staleApply),
  });

  const afterRace = await readGroups();
  const raceGroup = findGroup(afterRace.data, groupID);
  const beforePreviewSQL = sqlSnapshot(e);
  const beforePreviewDelivery = await readDelivery(groupID);
  const preview = await request('POST', `${writeGroupPath(groupID)}/test-write-preview`);
  assert.equal(preview.status, 200, JSON.stringify(responseWitness(preview)));
  assert.equal(preview.body?.success, true, JSON.stringify(responseWitness(preview)));
  const previewData = preview.body.data;
  assert.ok(previewData?.token && previewData?.operation_id, JSON.stringify(responseWitness(preview)));
  assert.equal(previewData.group_id, groupID, JSON.stringify(responseWitness(preview)));

  const savedPatch = { row_policy: { interval_seconds: Number(raceGroup.row_policy.interval_seconds) + 1 } };
  const canonicalSaveRequest = casUpdate(afterRace.data, raceGroup, savedPatch);
  const canonicalSave = await request('PUT', writeGroupPath(groupID), canonicalSaveRequest);
  assert.equal(canonicalSave.status, 200, JSON.stringify(responseWitness(canonicalSave)));
  assert.equal(canonicalSave.body?.success, true, JSON.stringify(responseWitness(canonicalSave)));
  const savedGroup = canonicalSave.body.data.group;
  assert.notEqual(savedGroup.revision, raceGroup.revision);
  assert.equal(savedGroup.applied_revision, raceGroup.applied_revision, 'saving a draft does not change applied_revision');

  const confirmRequest = { token: previewData.token, operation_id: previewData.operation_id };
  const staleConfirm = await request('POST', `${writeGroupPath(groupID)}/test-write`, confirmRequest);
  assert.equal(staleConfirm.status, 409, JSON.stringify(responseWitness(staleConfirm)));
  assert.equal(responseErrorCode(staleConfirm), 'WRITE_GROUP_TEST_WRITE_PREVIEW_STALE', JSON.stringify(responseWitness(staleConfirm)));
  const afterStaleConfirmSQL = sqlSnapshot(e);
  assert.deepEqual(afterStaleConfirmSQL.targets, beforePreviewSQL.targets, 'stale test-write confirmation creates no target row');
  const previewStale = record(e, 'CapacityCASPreviewStaleNoTargetRow', {
    preview: responseWitness(preview),
    canonical_save: responseWitness(canonicalSave),
    confirm_request: safePreview(confirmRequest),
    confirm_response: responseWitness(staleConfirm),
    sql_before_preview: beforePreviewSQL,
    sql_after_confirm: afterStaleConfirmSQL,
  });

  const afterDraftDelivery = await readDelivery(groupID);
  const afterDraftSQL = sqlSnapshot(e);
  assert.equal(savedGroup.applied_revision, raceGroup.applied_revision);
  assert.equal(savedGroup.destination.connector_id, raceGroup.destination.connector_id);
  assert.equal(savedGroup.destination.connector_revision, raceGroup.destination.connector_revision);
  assert.deepEqual(deliveryContract(afterDraftDelivery.data), deliveryContract(beforePreviewDelivery.data));
  if (beforePreviewSQL.local.outbox !== null) assert.deepEqual(afterDraftSQL.local.outbox, beforePreviewSQL.local.outbox);
  if (beforePreviewSQL.local.samples !== null) assert.deepEqual(afterDraftSQL.local.samples, beforePreviewSQL.local.samples);
  const draftHeld = record(e, 'CapacityCASUnappliedDraftHeld', {
    before: { group: raceGroup, delivery: deliveryContract(beforePreviewDelivery.data), sql: beforePreviewSQL },
    saved_draft: savedGroup,
    after: { group: savedGroup, delivery: deliveryContract(afterDraftDelivery.data), sql: afterDraftSQL },
  });

  let ownedBrowser;
  try {
    let uiPage = e.page;
    if (!uiPage) {
      const opened = await openBrowser();
      ownedBrowser = opened.browser;
      uiPage = opened.page;
    }
    const uiBaseline = await readGroups();
    const uiGroup = findGroup(uiBaseline.data, groupID);
    const ui = await runUIConflictWitness(e, groupID, uiGroup, uiPage);
    const finalSQL = sqlSnapshot(e);
    assert.deepEqual(finalSQL.targets, initialSQL.targets, 'CAS and stale test-write paths never write destination rows');
    const uiConflict = record(e, 'CapacityCASUI409KeepsDraftAndReloads', { group_id: groupID, ...ui });

    return {
      group_id: groupID,
      initial_sql: initialSQL,
      initial_delivery: initialDelivery.data,
      cas_race: casRace,
      stale_apply: staleApplyResult,
      preview_stale: previewStale,
      draft_held: draftHeld,
      ui_conflict: uiConflict,
      final_sql: finalSQL,
    };
  } finally {
    if (ownedBrowser) await ownedBrowser.close();
  }
}
