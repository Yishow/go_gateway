import assert from 'node:assert/strict';
import { join } from 'node:path';
import { control, eventually, query } from './quality-lib.mjs';
import { sh } from './lib.mjs';
import { holdTargetWrites } from './capacity-target-lock.mjs';

const usageSQL = `SELECT
 (SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))),0) FROM wg_delivery_samples WHERE consumed=0 AND group_id=?) +
 (SELECT COALESCE(SUM(LENGTH(CAST(payload AS BLOB))),0) FROM wg_delivery_outbox WHERE state!='sql_committed' AND group_id=?) AS bytes`;

export async function runCapacityCases(e, groups, work, gatewayEnv) {
  const { a, b } = groups;
  const gatewayDB = join(work, 'gateway.db');
  const usage = (group) => {
    // IDs are actual UUIDs from canonical UI saves, never arbitrary SQL input.
    assert.match(group.id, /^[a-f0-9-]{36}$/);
    return query(gatewayDB, usageSQL.replaceAll('?', `'${group.id}'`))[0].bytes;
  };
  const attempted = async (second, group, label, pointIDs) => {
    await e.clock(second);
    const captured = await control('/poll', { group_id: group.id, acquisition_id: label, ...(pointIDs ? { point_ids: pointIDs } : {}) });
    assert.ok(captured.indices.length > 0, `actual production poll captured samples: ${JSON.stringify(captured)}`);
    return { captured, released: await control('/release', { indices: captured.indices }) };
  };
  await e.targets.disconnect();
  await e.feed(18, a, 'quota-first-owned-row'); await e.close(20);
  await eventually('first accepted row stays pending in outage', () => e.item(a, 10), (row) => row?.state === 'retrying');
  await e.feed(28, a, 'quota-second-owned-row'); await e.close(30);
  assert.equal(e.outbox().filter((row) => row.group_id === a.id).length, 2);
  const acceptedBefore = e.samples();
  // A disconnected target cannot supply restart schema metadata, so keep its
  // existing file readable and block only real writes for the quota cases.
  const writeLock = await holdTargetWrites(`${e.targets.database('a')}.away`);
  await e.targets.reconnect();
  try {
  const groupLimit = usage(a) + 1;
  gatewayEnv.F_FIXTURE_GROUP_QUOTA_BYTES = String(groupLimit);
  await e.restart(30);
  const beforeRefusal = e.samples();
  const groupRefusal = await attempted(38, a, 'quota-group-refused');
  assert.ok(groupRefusal.released.results.length > 0);
  assert.ok(groupRefusal.released.results.every((row) => !row.accepted && row.reason === 'quota-hard-limit'));
  assert.deepEqual(e.samples(), beforeRefusal);
  const hardGroup = await e.delivery(a);
  assert.equal(hardGroup.quota.state, 'hard_limit'); assert.equal(hardGroup.quota.scope, 'group');
  assert.equal(hardGroup.quota.intake_refused, true); assert.ok(hardGroup.quota.loss_risk_notice);
  await e.feed(38, b, 'quota-independent-healthy'); await e.close(40); await e.committed(b, 30);
  assert.equal(e.targets.rows('a').length, 0); assert.equal(e.targets.rows('b').length, 1);
  e.check('QuotaGroupRefusesNewAckAndHealthyScopeContinues', {
    configured_bytes: groupLimit, accepted_before: acceptedBefore, refusal: groupRefusal,
    actual_target_lock: { pid: writeLock.pid, sql: writeLock.sql },
    delivery: hardGroup, ui: await e.ui(a, 'quota-group-pending'), healthy_effect: e.proof(b, 30, 'b'),
  });

  delete gatewayEnv.F_FIXTURE_GROUP_QUOTA_BYTES;
  gatewayEnv.F_FIXTURE_GLOBAL_QUOTA_BYTES = String(usage(a) + 1);
  await e.restart(40);
  const beforeGlobal = e.samples();
  const globalRefusal = await attempted(48, b, 'quota-global-refused');
  assert.ok(globalRefusal.released.results.every((row) => !row.accepted && row.reason === 'quota-hard-limit'));
  assert.deepEqual(e.samples(), beforeGlobal);
  const hardGlobal = await e.delivery(b);
  assert.equal(hardGlobal.quota.scope, 'global'); assert.equal(hardGlobal.quota.state, 'hard_limit');
  assert.equal(hardGlobal.quota.intake_refused, true); assert.ok(hardGlobal.quota.loss_risk_notice);
  assert.equal(e.targets.rows('b').length, 1);
  e.check('QuotaGlobalRefusesNewAckWithoutDeletingAcceptedData', { refusal: globalRefusal, delivery: hardGlobal, ui: await e.ui(b, 'quota-global-refused') });
  } finally {
    await writeLock.close();
  }

  delete gatewayEnv.F_FIXTURE_GLOBAL_QUOTA_BYTES;
  await e.restart(50); await e.targets.reconnect(); await e.close(50);
  await e.committed(a, 10); await e.committed(a, 20);
  assert.equal(e.targets.rows('a').length, 2);
  // Successful closure may prune its consumed journal. The actual target
  // provenance and durable receipts, rather than retained raw journal rows,
  // prove both previously accepted buckets reached SQL.
  const deliveredIDs = e.targets.rows('a').flatMap((row) => JSON.parse(row.prov).map((member) => member.sample_id));
  assert.ok(acceptedBefore.every((sample) => deliveredIDs.includes(sample.sample_id)));
  e.check('CapacityRestoredDeliversOriginallyAcceptedRows', { effects: [e.proof(a, 10, 'a'), e.proof(a, 20, 'a')], rows: e.targets.rows('a') });

  await e.feed(58, a, 'disk-original-accepted-row');
  const pageGate = await control('/capacity', { kind: 'disk_full', enabled: true });
  assert.equal(pageGate.enabled, true); assert.equal(pageGate.max_page_count, pageGate.page_count);
  let diskRefusal;
  let beforeDiskRefusal;
  const acceptedAttempts = [];
  for (let index = 0; index < 256; index++) {
    beforeDiskRefusal = e.samples();
    const attempt = await attempted(58, a, `actual-disk-full-${index}`, [a.members[0].point_id]);
    assert.equal(attempt.released.results.length, 1);
    if (!attempt.released.results[0].accepted) { diskRefusal = attempt; break; }
    acceptedAttempts.push(attempt);
  }
  assert.ok(diskRefusal, 'real file page growth eventually reaches SQLITE_FULL');
  assert.equal(diskRefusal.released.results[0].reason, 'disk-full');
  assert.deepEqual(e.samples(), beforeDiskRefusal);
  assert.equal(e.targets.rows('a').length, 2);
  await control('/capacity', { kind: 'disk_full', enabled: false });
  await e.close(60); await e.committed(a, 50);
  assert.equal(e.targets.rows('a').length, 3);
  // The snapshot contract selects the latest accepted value per member;
  // repeated temperature samples intentionally become one row member.
  const diskRow = e.targets.rows('a').find((row) => JSON.parse(row.prov).every((member) => Date.parse(member.observed_at) === Date.parse('2026-01-01T00:00:58Z')));
  assert.ok(diskRow); const diskMembers = JSON.parse(diskRow.prov);
  assert.equal(diskMembers.length, a.members.length);
  assert.ok(diskMembers.every((member) => beforeDiskRefusal.some((sample) => sample.sample_id === member.sample_id)));
  e.check('ActualDiskFullRefusesAckAndKeepsAcceptedData', {
    page_gate: pageGate, accepted_attempts: acceptedAttempts, refusal: diskRefusal,
    before_refusal: beforeDiskRefusal, effect_after_restore: e.proof(a, 50, 'a'), ui: await e.ui(a, 'disk-restored'),
  });

  const targetDB = e.targets.database('a');
  const trigger = "CREATE TRIGGER f_capacity_poison BEFORE INSERT ON readings WHEN NEW.line='A' BEGIN SELECT RAISE(ABORT, 'row_rejected'); END";
  sh('sqlite3', [targetDB, trigger]);
  await e.feed(68, a, 'poison-first-a'); await e.feed(68, b, 'poison-healthy-b'); await e.close(70);
  const poisoned = await eventually('actual destination row rejection is quarantined', () => e.item(a, 60), (row) => row?.state === 'quarantined');
  assert.equal(poisoned.last_error_code, 'destination-rejected-row');
  await e.committed(b, 60);
  await e.feed(78, a, 'poison-held-successor'); await e.feed(78, b, 'poison-second-healthy'); await e.close(80);
  await e.committed(b, 70);
  assert.equal(e.item(a, 70).state, 'pending'); assert.equal(e.targets.rows('a').length, 3);
  assert.equal(e.targets.rows('b').length, 3);
  e.check('PoisonRetainsPayloadAndOnlyBlocksItsOwnPartition', {
    fixture_ddl: trigger, quarantined: poisoned, held_successor: e.item(a, 70),
    healthy_effects: [e.proof(b, 60, 'b'), e.proof(b, 70, 'b')], ui: await e.ui(a, 'poison-attention'),
  });
  sh('sqlite3', [targetDB, 'DROP TRIGGER f_capacity_poison']);
  return { usage_query: usageSQL };
}
