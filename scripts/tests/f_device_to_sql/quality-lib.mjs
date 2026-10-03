// Deterministic controls operate only the acceptance binary. All values come
// from real Modbus reads; releases use the production journal and SQL sender.
import assert from 'node:assert/strict';
import { api, BASE, boundedFetch, sh } from './lib.mjs';

export const FIXTURE_PORT = process.env.F_FIXTURE_PORT ?? '3355';
export const START = Date.parse('2026-01-01T00:00:00Z');
export const stamp = (seconds) => new Date(START + seconds * 1000).toISOString();
export const query = (path, sql) => JSON.parse(sh('sqlite3', ['-json', path, sql]) || '[]');

export async function control(path, body = {}) {
  const response = await boundedFetch(`http://127.0.0.1:${FIXTURE_PORT}${path}`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
  });
  const result = await response.json();
  assert.equal(response.ok, true, `fixture ${path}: ${JSON.stringify(result)}`);
  return result;
}

export async function eventually(label, read, predicate, timeout = 15000) {
  const deadline = Date.now() + timeout;
  let result;
  while (Date.now() < deadline) {
    result = await read();
    if (predicate(result)) return result;
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  assert.fail(`${label}: ${JSON.stringify(result)}`);
}

export async function groupMutation(group, patch, apply = true) {
  const list = await api('/studio-v2/workspace/write-groups');
  const latest = list.groups.find((item) => item.id === group.id);
  assert.ok(latest, 'fixture group must still exist');
  const request = {
    workspace_id: list.workspace_id, expected_workspace_revision: list.workspace_revision,
    expected_group_revision: latest.revision,
    expected_connector_revision: latest.destination.connector_revision,
    group: { ...latest, ...patch },
  };
  const saved = await mutate(`/studio-v2/workspace/write-groups/${group.id}`, 'PUT', request);
  if (!apply) return saved;
  const { group: ignored, ...revisions } = request;
  return mutate(`/studio-v2/workspace/write-groups/${group.id}/apply`, 'POST', {
    ...revisions, expected_workspace_revision: saved.workspace_revision, expected_group_revision: saved.group.revision,
  });
}

async function mutate(path, method, body) {
  const response = await boundedFetch(`${BASE}/api/v1/datalink${path}`, {
    method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body),
  });
  const result = await response.json();
  assert.equal(response.ok, true, `${method} ${path}: ${JSON.stringify(result)}`);
  return result.data;
}

export async function showDelivery(page, group, view, file) {
  await page.reload();
  await page.getByTestId('step-nav-button-4').click();
  await page.getByTestId(`group-open-${group.id}`).click();
  await page.getByTestId('group-delivery-strip').waitFor();
  await eventually('UI skipped count matches durable API', () => page.getByTestId('delivery-skipped-buckets').innerText(),
    (text) => text === String(view.skipped_buckets));
  await eventually('UI silent count matches durable API', () => page.getByTestId('delivery-no-data-buckets').innerText(),
    (text) => text === String(view.no_data_buckets));
  const causes = await page.getByTestId('delivery-bucket-causes').innerText();
  const labels = {
    missing: /missing data|缺少資料/, bad: /bad quality|品質不良/, stale: /stale data|資料過期/,
    invalid: /invalid data|資料無效/, no_data: /no data|沒有資料/, unavailable: /reason unavailable|原因未確認/,
  };
  for (const [index, issue] of view.recent_bucket_issues.entries()) {
    const shown = await page.getByTestId(`delivery-recent-issue-${index}`).innerText();
    assert.ok(shown.includes(issue.bucket_start), 'UI reason keeps its actual bucket time');
    for (const cause of issue.causes) assert.match(shown, labels[cause], `UI cause matches durable ${cause}`);
  }
  await page.screenshot({ path: file, fullPage: true });
  return { causes, skipped_buckets: view.skipped_buckets, no_data_buckets: view.no_data_buckets };
}
