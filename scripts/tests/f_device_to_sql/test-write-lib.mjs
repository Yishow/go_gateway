import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { join } from 'node:path';
import { api, BASE, boundedFetch, sh } from './lib.mjs';
import { mutation } from './recovery-lib.mjs';

export const operationSQL = 'SELECT operation_id,status,detail,updated_at FROM managed_schema_operations WHERE operation_id=?';
export const safePreview = ({ token, ...preview }) => ({ ...preview, token_sha256: createHash('sha256').update(token).digest('hex') });

export async function showGroup(page, group) {
  await page.reload(); await page.getByTestId('step-nav-button-4').click();
  await page.getByTestId(`group-open-${group.id}`).click();
  await page.getByTestId('group-editor').waitFor();
}

export async function uiPreview(page, group, rows) {
  await showGroup(page, group);
  const before = rows();
  const button = page.getByTestId('group-test-write-preview');
  assert.equal(await button.isEnabled(), true, `real UI permits preview: ${await page.getByTestId('group-test-write').innerText()}`);
  const [actual] = await Promise.all([
    page.waitForResponse((r) => r.url().endsWith(`/write-groups/${group.id}/test-write-preview`) && r.request().method() === 'POST'),
    button.click(),
  ]);
  assert.equal(actual.status(), 200);
  const preview = (await actual.json()).data;
  await page.getByTestId('group-test-write-preview-card').waitFor();
  assert.deepEqual(rows(), before, 'actual UI preview does not mutate target rows');
  assert.match(preview.owner_value, /^gw-test-[a-f0-9-]+$/);
  return preview;
}

export async function uiConfirm(page, group) {
  const [actual] = await Promise.all([
    page.waitForResponse((r) => r.url().endsWith(`/write-groups/${group.id}/test-write`) && r.request().method() === 'POST'),
    page.getByTestId('group-test-write-confirm').click(),
  ]);
  assert.equal(actual.status(), 200);
  await page.getByTestId('group-test-write-outcome').waitFor();
  return (await actual.json()).data;
}

export async function confirm(group, preview) {
  const response = await boundedFetch(`${BASE}/api/v1/datalink/studio-v2/workspace/write-groups/${group.id}/test-write`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token: preview.token, operation_id: preview.operation_id }),
  });
  return { status: response.status, body: await response.json() };
}

export function assertResult(op, write, cleanup) {
  assert.equal(op.write_outcome, write); assert.equal(op.cleanup_status, cleanup);
}

export async function retained(group, preview, original, rows) {
  const before = rows(); const replay = await confirm(group, preview);
  assert.equal(replay.status, 200); assert.deepEqual(replay.body.data, original);
  assert.deepEqual(rows(), before, 'same-operation retry never rewrites or repeats cleanup');
  const read = await api(`/studio-v2/workspace/database-operations/${preview.operation_id}`);
  assert.deepEqual(read, original);
  return { confirm_status: replay.status, operation_get: read, unchanged_rows: true };
}

export async function savedClone(group, connector, name) {
  const list = await api('/studio-v2/workspace/write-groups');
  const { id, revision, applied_revision, status, created_at, updated_at, ...copy } = group;
  const result = await mutation('/studio-v2/workspace/write-groups', 'POST', {
    workspace_id: list.workspace_id, expected_workspace_revision: list.workspace_revision,
    expected_group_revision: '', expected_connector_revision: connector.identity_revision,
    group: {
      ...copy, name, write_policy: { dedupe_capability: 'none' },
      destination: { ...copy.destination, connector_id: connector.id, connector_revision: connector.identity_revision },
    },
  });
  return result.group;
}

export function permissionTargets(targets, runNumber) {
  const schema = targets.schema('a'); assert.match(schema, /^gw_f_recovery_a_\d+$/);
  const role = `gw_f_permissions_${runNumber}`;
  const psql = (sql) => sh('docker', ['exec', '-i', 'gw-wg-pg-test', 'psql', '-v', 'ON_ERROR_STOP=1', '-U', 'postgres', '-d', 'gwtest', '-At'], { input: sql });
  const password = createHash('sha256').update(`${role}:disposable-fixture`).digest('hex');
  psql(`CREATE ROLE "${role}" LOGIN PASSWORD '${password}'`);
  try {
    for (const name of ['a', 'b']) {
      const ownedSchema = targets.schema(name); assert.match(ownedSchema, /^gw_f_recovery_[ab]_\d+$/);
      psql(`GRANT USAGE ON SCHEMA "${ownedSchema}" TO "${role}"; GRANT INSERT,SELECT,DELETE ON TABLE "${ownedSchema}".readings,"${ownedSchema}".gw_effect_receipts TO "${role}";`);
    }
  } catch (error) {
    psql(`DROP OWNED BY "${role}"; DROP ROLE "${role}"`);
    throw error;
  }
  return {
    role,
    destination: { ...targets.destination, user: role, password },
    restrict: (kind) => {
      assert.ok(['no_select', 'no_delete'].includes(kind));
      psql(`GRANT INSERT,SELECT,DELETE ON TABLE "${schema}".readings TO "${role}"; REVOKE ${kind === 'no_select' ? 'SELECT' : 'DELETE'} ON TABLE "${schema}".readings FROM "${role}";`);
    },
    removeOwned: (owner) => {
      assert.match(owner, /^gw-test-[a-f0-9-]+$/);
      psql(`DELETE FROM "${schema}".readings WHERE line='${owner}'`);
    },
    cleanup: () => psql(`DROP OWNED BY "${role}"; DROP ROLE "${role}"`),
  };
}

export async function screenshotOutcome(page, evidence, kind, suffix) {
  const write = await page.getByTestId('group-test-write-outcome').innerText();
  const cleanup = await page.getByTestId('group-test-write-cleanup').innerText();
  await page.screenshot({ path: join(evidence, `test-write-${kind}-${suffix}.png`), fullPage: true });
  return { write, cleanup };
}
