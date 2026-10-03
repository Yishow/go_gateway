import assert from 'node:assert/strict';
import { join } from 'node:path';
import { api, BASE, boundedFetch, GW_PORT, startProcess, waitForGateway, waitForPortFree } from './lib.mjs';
import { control, eventually, FIXTURE_PORT, query, stamp } from './quality-lib.mjs';

export const localOutboxSQL = 'SELECT effect_key, record_id, group_id, group_revision, bucket_start, connector_id, connector_revision, database_name, dedupe_capability, payload_digest, state, retry_count, next_retry_at, claim_owner, claim_expires_at, last_error_code FROM wg_delivery_outbox ORDER BY bucket_start,group_id';
export const localSamplesSQL = 'SELECT group_id,group_revision,sample_id,bucket_start,consumed FROM wg_delivery_samples ORDER BY seq';
export const localReceiptsSQL = 'SELECT effect_key,payload_digest FROM wg_delivery_receipts';
export const localCheckpointSQL = 'SELECT group_id,group_revision,next_close FROM wg_delivery_checkpoints ORDER BY group_id,group_revision';

export async function mutation(path, method, body) {
  const response = await boundedFetch(`${BASE}/api/v1/datalink${path}`, {
    method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
  });
  const result = await response.json();
  assert.equal(response.ok, true, `${method} ${path}: ${JSON.stringify(result)}`);
  return result.data;
}

export async function saveApply(group, patch) {
  const list = await api('/studio-v2/workspace/write-groups');
  const latest = list.groups.find((item) => item.id === group.id);
  const connectorRevision = patch.destination?.connector_revision ?? latest.destination.connector_revision;
  const revisions = {
    workspace_id: list.workspace_id, expected_workspace_revision: list.workspace_revision,
    expected_group_revision: latest.revision, expected_connector_revision: connectorRevision,
  };
  const saved = await mutation(`/studio-v2/workspace/write-groups/${group.id}`, 'PUT', { ...revisions, group: { ...latest, ...patch } });
  const readiness = await api(`/studio-v2/workspace/write-groups/${group.id}/readiness`);
  assert.equal(readiness.ready, true, `actual saved group readiness: ${JSON.stringify(readiness)}`);
  return (await mutation(`/studio-v2/workspace/write-groups/${group.id}/apply`, 'POST', {
    ...revisions, expected_workspace_revision: saved.workspace_revision, expected_group_revision: saved.group.revision,
  })).group;
}

export async function disable(group) {
  const list = await api('/studio-v2/workspace/write-groups');
  const latest = list.groups.find((item) => item.id === group.id);
  return mutation(`/studio-v2/workspace/write-groups/${group.id}/disable`, 'POST', {
    workspace_id: list.workspace_id, expected_workspace_revision: list.workspace_revision,
    expected_group_revision: latest.revision, expected_connector_revision: latest.destination.connector_revision,
  });
}

export async function createNoneGroup(group) {
  const list = await api('/studio-v2/workspace/write-groups');
  const { id, revision, applied_revision, status, created_at, updated_at, ...copy } = group;
  const saved = await mutation('/studio-v2/workspace/write-groups', 'POST', {
    workspace_id: list.workspace_id, expected_workspace_revision: list.workspace_revision,
    expected_group_revision: '', expected_connector_revision: group.destination.connector_revision,
    group: { ...copy, name: 'Local receipt without dedupe', write_policy: { dedupe_capability: 'none' }, members: group.members.map((member) => ({ ...member, entity_key: 'C' })) },
  });
  return (await mutation(`/studio-v2/workspace/write-groups/${saved.group.id}/apply`, 'POST', {
    workspace_id: list.workspace_id, expected_workspace_revision: saved.workspace_revision,
    expected_group_revision: saved.group.revision, expected_connector_revision: group.destination.connector_revision,
  })).group;
}

export function recoveryContext(work, children, targets, page, evidence, results, gatewayEnv = {}) {
  const gatewayDB = join(work, 'gateway.db');
  let incarnation = children.gateway ? 1 : 0;
  const startGateway = (second = -10) => {
    children.gateway = startProcess(work, `gateway-${incarnation++}`, join(work, 'gw-ui'), [], {
      ...gatewayEnv,
      GATEWAY_DB_PATH: gatewayDB, PORT: GW_PORT, F_FIXTURE_PORT: FIXTURE_PORT, F_FIXTURE_START_AT: stamp(second),
    });
  };
  const kill = (name, signal = 'SIGTERM') => { try { process.kill(-children[name].pid, signal); } catch { /* owned child already exited */ } };
  const restart = async (second) => {
    kill('gateway', 'SIGKILL');
    await waitForPortFree();
    await eventually('fixture listener closed after SIGKILL', async () => {
      try { await fetch(`http://127.0.0.1:${FIXTURE_PORT}/faults`, { signal: AbortSignal.timeout(2000) }); return false; } catch { return true; }
    }, Boolean);
    startGateway(second); await waitForGateway(); await control('/pause');
  };
  const clock = (second) => control('/clock', { at: stamp(second) });
  const feed = async (second, group, label) => {
    await clock(second);
    const captured = await control('/poll', { group_id: group.id, acquisition_id: label });
    assert.equal(captured.indices.length, group.members.length);
    const released = await control('/release', { indices: captured.indices });
    assert.equal(released.results.every((result) => result.accepted), true, JSON.stringify(released));
    return { captured, released };
  };
  const close = async (second) => { await clock(second); return control('/tick'); };
  const state = async () => {
    const response = await fetch(`http://127.0.0.1:${FIXTURE_PORT}/faults`, { signal: AbortSignal.timeout(3000) });
    assert.equal(response.ok, true); return response.json();
  };
  const fault = (kind, group, enabled) => control('/fault', { kind, group_id: group.id, enabled });
  const reached = (kind, group, previous = 0) => eventually(`${kind} actual barrier reached`, state,
    (value) => value.faults?.some((item) => item.kind === kind && item.group_id === group.id && item.reached > previous && item.holding));
  const outbox = () => query(gatewayDB, localOutboxSQL);
  const samples = () => query(gatewayDB, localSamplesSQL);
  const receipts = () => query(gatewayDB, localReceiptsSQL);
  const checkpoints = () => query(gatewayDB, localCheckpointSQL);
  const item = (group, second) => outbox().find((row) => row.group_id === group.id && Date.parse(row.bucket_start) === Date.parse(stamp(second)));
  const committed = (group, second, timeout = 40000) => eventually('actual local confirmed effect', () => item(group, second), (row) => row?.state === 'sql_committed', timeout);
  const delivery = (group) => api(`/studio-v2/workspace/write-groups/${group.id}/delivery`);
  const proof = (group, second, targetName) => {
    const row = item(group, second);
    assert.ok(row);
    const matches = targets.receipts(targetName).filter((entry) => entry.effect_key === row.effect_key);
    assert.equal(matches.length, row.dedupe_capability === 'receipt' ? 1 : 0);
    if (matches.length) assert.equal(matches[0].payload_digest, row.payload_digest);
    return row;
  };
  const ui = async (group, suffix) => {
    await page.reload(); await page.getByTestId('step-nav-button-4').click();
    await page.getByTestId(`group-open-${group.id}`).click();
    await page.getByTestId('group-delivery-strip').waitFor();
    const observed = await eventually('UI database confirmation follows durable receipt', async () => ({
      view: await delivery(group), text: await page.getByTestId('delivery-committed').innerText(),
    }), ({ view, text }) => text === String(view.stages.sql_committed));
    await page.screenshot({ path: join(evidence, `recovery-${targets.destination.kind}-${suffix}.png`), fullPage: true });
    return { view: observed.view, text: await page.getByTestId('group-delivery-strip').innerText() };
  };
  const check = (name, details) => { results[name] = { passed: true, ...details }; console.log(`PASS ${name}`); };
  return { startGateway, kill, restart, clock, feed, close, state, fault, reached, outbox, samples, receipts, checkpoints, item, committed, delivery, proof, ui, check, targets, page };
}
