import test from 'node:test';
import assert from 'node:assert/strict';
import { verifyFreshDeliveryEvidence } from './fresh-recording-verification.mjs';

const revision = 'rev-current';
const table = 'gw_group_test';
const members = [
  { device_id: 'device-a', point_id: 'point-a-1', tag_id: 'tag-a-1', target_column: 'v_a1' },
  { device_id: 'device-a', point_id: 'point-a-2', tag_id: 'tag-a-2', target_column: 'v_a2' },
  { device_id: 'device-b', point_id: 'point-b-1', tag_id: 'tag-b-1', target_column: 'v_b1' },
  { device_id: 'device-b', point_id: 'point-b-2', tag_id: 'tag-b-2', target_column: 'v_b2' },
];
const mappings = members.map((member, index) => ({
  ...member, tag_key: `tag.${index + 1}`,
}));
const expectedByTag = Object.fromEntries(mappings.map((mapping, index) => [mapping.tag_key, String(index + 1)]));
const group = { id: 'group-1', applied_revision: revision, destination: { table_name: table }, members };

function provenanceFor(rowMembers, bucketStart, suffix = '') {
  return JSON.stringify(rowMembers.map((member, index) => ({
    member: JSON.stringify([member.device_id, member.point_id, member.tag_id]),
    status: 'ok', quality: 'good', sample_id: `sample-${suffix}-${index}`,
    observed_at: bucketStart.replace('00Z', '30Z'),
  })));
}

function rowsForBuckets() {
  const rows = [];
  for (const [bucketIndex, minute] of [0, 1, 2].entries()) {
    const bucketStart = `2026-10-04T00:0${minute}:00Z`;
    for (const device of ['device-a', 'device-b']) {
      const rowMembers = members.filter((member) => member.device_id === device);
      const recordId = `record-${bucketIndex}-${device}`;
      const row = {
        record_id: recordId, group_id: group.id, device_id: device, bucket_start: bucketStart,
        provenance: provenanceFor(rowMembers, bucketStart, recordId),
        v_a1: null, v_a2: null, v_b1: null, v_b2: null,
      };
      row[`v_${device === 'device-a' ? 'a1' : 'b1'}`] = device === 'device-a' ? '1' : '3';
      row[`v_${device === 'device-a' ? 'a2' : 'b2'}`] = device === 'device-a' ? '2' : '4';
      rows.push(row);
    }
  }
  return rows;
}

function observations(rows) {
  const outbox = rows.map((row) => ({
    effect_key: `effect-${row.record_id}`, record_id: row.record_id, group_id: group.id,
    group_revision: revision, bucket_start: row.bucket_start, table_name: table,
    state: 'sql_committed', payload_digest: `digest-${row.record_id}`,
  }));
  const receipts = outbox.map((effect) => ({
    effect_key: effect.effect_key, payload_digest: effect.payload_digest, committed_at: '2026-10-04T00:00:59Z',
  }));
  const closedBuckets = rows.map((row) => ({
    group_revision: revision, bucket_start: row.bucket_start, kind: 'row', record_id: row.record_id,
    members: row.provenance,
  }));
  return { outbox, receipts, closedBuckets };
}

test('verifies each multi-entity row against its own members and durable receipt', () => {
  const rows = rowsForBuckets();
  const { outbox, receipts, closedBuckets } = observations(rows);
  const evidence = verifyFreshDeliveryEvidence({
    group, table, rows, mappings, expectedByTag, closedBuckets, outbox, receipts,
  });
  assert.deepEqual(evidence.failures, []);
});

test('rejects changed source sample identity and even one nanosecond of acquisition time', () => {
  for (const change of ['sample', 'time']) {
    const rows = rowsForBuckets();
    const durable = observations(rows);
    const entries = JSON.parse(rows[0].provenance);
    if (change === 'sample') entries[0].sample_id = 'different-source-sample';
    else entries[0].observed_at = '2026-10-04T00:00:30.000000001Z';
    rows[0].provenance = JSON.stringify(entries);
    const evidence = verifyFreshDeliveryEvidence({ group, table, rows, mappings, expectedByTag, ...durable });
    assert.ok(evidence.failures.some((failure) => failure.includes('frozen acquisition provenance')), change);
  }
});

test('accepts UTC acquisition timestamps with equivalent fractional padding', () => {
  const rows = rowsForBuckets();
  const durable = observations(rows);
  const entries = JSON.parse(rows[0].provenance);
  entries[0].observed_at = '2026-10-04T00:00:30.000000000Z';
  rows[0].provenance = JSON.stringify(entries);
  const evidence = verifyFreshDeliveryEvidence({ group, table, rows, mappings, expectedByTag, ...durable });
  assert.deepEqual(evidence.failures, []);
});

test('rejects duplicate provenance tuple and non-UTC observed_at', () => {
  const rows = rowsForBuckets();
  const { outbox, receipts, closedBuckets } = observations(rows);
  const first = JSON.parse(rows[0].provenance);
  first[1] = { ...first[0], observed_at: '2026-10-04T00:00:30' };
  rows[0].provenance = JSON.stringify(first);
  const evidence = verifyFreshDeliveryEvidence({
    group, table, rows, mappings, expectedByTag, closedBuckets, outbox, receipts,
  });
  assert.ok(evidence.failures.some((failure) => failure.includes('quality/source provenance')));
});

test('uses the persisted entity key when the SQL row omits device_id', () => {
  const entityMembers = members.map((member) => ({ ...member, entity_key: member.device_id === 'device-a' ? 'A' : 'B' }));
  const entityGroup = { ...group, row_policy: { entity_key_column: 'entity_id' }, members: entityMembers };
  const rows = rowsForBuckets().map((row) => {
    const entity = row.device_id === 'device-a' ? 'A' : 'B';
    const copy = { ...row, entity_id: entity };
    delete copy.device_id;
    return copy;
  });
  const { outbox, receipts, closedBuckets } = observations(rows);
  const evidence = verifyFreshDeliveryEvidence({
    group: entityGroup, table, rows, mappings, expectedByTag, closedBuckets, outbox, receipts,
    entityColumn: 'entity_id',
  });
  assert.deepEqual(evidence.failures, []);
});

test('rejects entities whose individual three buckets do not share the same first three', () => {
  const rows = rowsForBuckets();
  for (const row of rows.filter((item) => item.device_id === 'device-b')) {
    row.bucket_start = new Date(Date.parse(row.bucket_start) + 60_000).toISOString();
    row.provenance = provenanceFor(members.filter((member) => member.device_id === 'device-b'), row.bucket_start, row.record_id);
  }
  rows.sort((a, b) => a.bucket_start.localeCompare(b.bucket_start));
  const evidence = verifyFreshDeliveryEvidence({ group, table, rows, mappings, expectedByTag, ...observations(rows) });
  assert.ok(evidence.failures.some((failure) => failure.includes('entity rows missing')));
});

test('rejects an unknown device with empty provenance and NULL business values', () => {
  const rows = rowsForBuckets();
  rows[0] = { ...rows[0], device_id: 'unknown-device', provenance: '[]', v_a1: null, v_a2: null };
  const evidence = verifyFreshDeliveryEvidence({ group, table, rows, mappings, expectedByTag, ...observations(rows) });
  assert.ok(evidence.failures.some((failure) => failure.includes('unknown entity')));
});
