import assert from 'node:assert/strict';
import { api } from './lib.mjs';
import { eventually, stamp } from './quality-lib.mjs';
import { createNoneGroup, disable, mutation, saveApply } from './recovery-lib.mjs';

const ownOpen = (e, group, second) => e.samples().filter((row) => row.group_id === group.id && !row.consumed && Date.parse(row.bucket_start) === Date.parse(stamp(second)));
const checkpoint = (e, group) => e.checkpoints().find((row) => row.group_id === group.id && row.group_revision === group.applied_revision)?.next_close;

export function verifyRecoverySQL(e, groups) {
  const keys = new Set();
  for (const [name, count] of [['a', 6], ['b', 5], ['other', 0]]) {
    const rows = e.targets.rows(name);
    assert.equal(rows.length, count);
    for (const row of rows) {
      const group = groups[row.line.toLowerCase()];
      assert.ok(group, 'target row belongs to an actual persisted group entity');
      assert.equal(row.temperature, row.line === 'A' ? 215 : 187);
      assert.equal(Number(row.pressure), row.line === 'A' ? 1013 : 777);
      assert.equal(Number(row.running), row.line === 'A' ? 1 : 0);
      assert.equal(row.batch, row.line === 'A' ? '9007199254740993' : '2');
      const provenance = JSON.parse(row.prov);
      assert.equal(provenance.length, group.members.length);
      const expectedMembers = group.members.map((member) => JSON.stringify([member.device_id, member.point_id, member.tag_id])).sort();
      assert.deepEqual(provenance.map((member) => member.member).sort(), expectedMembers);
      for (const member of provenance) {
        assert.equal(member.status, 'ok'); assert.equal(member.quality, 'good');
        assert.ok(member.sample_id); assert.match(member.observed_at, /Z$/);
        assert.ok(Number.isFinite(Date.parse(member.observed_at)));
      }
      const bucket = Math.floor(Date.parse(provenance[0].observed_at) / 10000);
      assert.ok(provenance.every((member) => Math.floor(Date.parse(member.observed_at) / 10000) === bucket));
      const key = `${row.line}|${bucket}`;
      assert.equal(keys.has(key), false, 'actual SQL contains one row for each entity/bucket'); keys.add(key);
    }
  }
  assert.equal(e.targets.receipts('a').length, 6);
  assert.equal(e.targets.receipts('b').length, 2);
}

export async function runRecoveryCases(e, groups) {
  let { a, b } = groups;
  await e.feed(18, a, 'ack-eighteen');
  assert.equal(ownOpen(e, a, 10).length, 4);
  assert.equal(e.targets.rows('a').length, 0);
  const ackBefore = { samples: e.samples(), checkpoints: e.checkpoints() };
  await e.restart(18); await e.close(20); await e.committed(a, 10);
  assert.equal(e.targets.rows('a').length, 1);
  e.check('AtomicAckThenKillRecovery', { before: ackBefore, effect: e.proof(a, 10, 'a') });

  await e.feed(28, a, 'closure-twenty-eight');
  await e.fault('closure_hold', a, true);
  const heldTick = e.close(30).catch((error) => ({ interrupted: error.message }));
  const held = await e.reached('closure_hold', a);
  assert.equal(ownOpen(e, a, 20).length, 4);
  assert.equal(e.item(a, 20), undefined);
  assert.equal(Date.parse(checkpoint(e, a)), Date.parse(stamp(20)));
  assert.equal(e.targets.rows('a').length, 1);
  const uncommitted = { samples: e.samples(), checkpoints: e.checkpoints(), outbox: e.outbox(), held };
  e.kill('gateway', 'SIGKILL'); await heldTick;
  await e.restart(28); await e.fault('closure_hold', a, false); await e.close(30); await e.committed(a, 20);
  assert.equal(e.targets.rows('a').length, 2);
  e.check('AtomicClosureTransactionKillRollbackAndRecovery', { before_kill: uncommitted, effect: e.proof(a, 20, 'a') });

  await e.targets.disconnect();
  await e.feed(38, a, 'outage-thirty-eight-a'); await e.feed(38, b, 'outage-thirty-eight-b'); await e.close(40);
  const offline = await eventually('real target outage keeps a frozen backlog', () => e.item(a, 30), (row) => row?.state === 'retrying');
  await e.committed(b, 30);
  assert.equal(e.targets.rows('a').length, 2); assert.equal(e.targets.rows('b').length, 1);
  const offlineUI = await e.ui(a, 'outage-pending');
  assert.equal(offlineUI.view.stages.sql_committed, 2);
  await e.restart(40);
  await e.feed(48, b, 'healthy-after-restart'); await e.close(50); await e.committed(b, 40);
  assert.equal(e.targets.rows('b').length, 2); assert.equal(e.targets.rows('a').length, 2);
  await e.targets.reconnect(); await e.close(50); await e.committed(a, 30);
  assert.equal(e.targets.rows('a').length, 3);
  assert.equal(e.item(a, 30).payload_digest, offline.payload_digest);
  e.check('ProductionGroupOutageRecoveryAndHealthyConnector', { frozen: offline, ui_pending: offlineUI, effect: e.proof(a, 30, 'a'), healthy_effect: e.proof(b, 40, 'b') });

  await e.fault('target_commit_response_lost', a, true);
  await e.feed(58, a, 'lost-response-fifty-eight'); await e.close(60);
  const returnedLost = await eventually('ambiguous actual commit enters safe retry', () => e.item(a, 50), (row) => ['retrying', 'sql_committed'].includes(row?.state) && row.retry_count >= 1);
  const lostFault = (await e.state()).faults.find((fault) => fault.kind === 'target_commit_response_lost' && fault.group_id === a.id);
  assert.ok(lostFault.enabled && lostFault.reached >= 1 && lostFault.commits >= 1);
  assert.equal(e.targets.rows('a').length, 4);
  await e.restart(60); await e.committed(a, 50);
  assert.equal(e.targets.rows('a').length, 4);
  e.check('TargetReceiptIdentityResponseLostOneEffect', { fault: lostFault, after_response_lost: returnedLost, effect: e.proof(a, 50, 'a') });

  await e.fault('target_commit_hold', a, true);
  await e.feed(68, a, 'commit-then-kill-sixty-eight'); await e.close(70);
  const commitHeld = await e.reached('target_commit_hold', a);
  const sending = e.item(a, 60);
  assert.equal(sending.state, 'sending'); assert.equal(e.targets.rows('a').length, 5);
  assert.equal(e.receipts().some((row) => row.effect_key === sending.effect_key), false);
  const unconfirmedUI = await e.ui(a, 'real-commit-unconfirmed');
  assert.equal(unconfirmedUI.view.stages.sql_committed, 4);
  await e.restart(70); await e.committed(a, 60);
  assert.equal(e.targets.rows('a').length, 5);
  e.check('TargetCommitThenKillReceiptRecoveryOneEffect', { barrier: commitHeld, before_kill: sending, ui_unconfirmed: unconfirmedUI, effect: e.proof(a, 60, 'a') });

  await e.fault('local_receipt_failure', a, true);
  await e.feed(78, a, 'local-receipt-seventy-eight'); await e.close(80);
  const receiptFailed = await eventually('target commit with failed local receipt', () => e.item(a, 70), (row) => row?.state === 'sending' && e.targets.rows('a').length === 6);
  assert.equal(e.receipts().some((row) => row.effect_key === receiptFailed.effect_key), false);
  await e.restart(80); await e.fault('local_receipt_failure', a, false);
  // Startup may retry before the persisted receipt trigger is removed. The
  // production lease remains valid for that attempt; observe its real expiry.
  await e.committed(a, 70, 100000);
  assert.equal(e.targets.rows('a').length, 6);
  e.check('LocalReceiptFailureWithReceiptConvergesOnce', { before_restart: receiptFailed, effect: e.proof(a, 70, 'a') });

  await e.targets.disconnect(); await e.feed(88, a, 'frozen-endpoint-eighty-eight'); await e.close(90);
  const frozen = await eventually('endpoint-edit fixture has real pending data', () => e.item(a, 80), (row) => row?.state === 'retrying');
  const connector = await api(`/db-targets/connectors/${a.destination.connector_id}`);
  const updated = await mutation(`/db-targets/connectors/${connector.id}`, 'PUT', { expected_identity_revision: connector.identity_revision, connection_config: e.targets.otherConfig });
  assert.notEqual(updated.identity_revision, connector.identity_revision);
  const blocked = await eventually('frozen destination rejects retargeting', () => e.item(a, 80), (row) => row?.state === 'blocked' && row.last_error_code === 'target-blocked', 40000);
  assert.equal(blocked.connector_revision, frozen.connector_revision); assert.equal(blocked.payload_digest, frozen.payload_digest);
  assert.equal(e.targets.rows('other').length, 0);
  await disable(a); await e.close(100);
  assert.equal(e.item(a, 80).state, 'blocked'); assert.equal(e.targets.rows('other').length, 0);
  e.check('RevisionBoundBacklogEndpointEditAndDisableNeverRetarget', { original: frozen, blocked, ui: await e.ui(a, 'edited-endpoint-blocked') });

  // Independent group/entity uses the same actual persisted B device members.
  await e.clock(100);
  const c = await createNoneGroup(b);
  b = await saveApply(b, { write_policy: { dedupe_capability: 'none' } });
  await e.close(110);
  await e.fault('local_receipt_failure', c, true);
  await e.feed(118, c, 'none-local-receipt-one-eighteen'); await e.close(120);
  const noneReceipt = await eventually('none target committed before local receipt failed', () => e.item(c, 110), (row) => row?.state === 'sending' && e.targets.rows('b').some((target) => target.line === 'C'));
  await e.restart(120); await e.fault('local_receipt_failure', c, false);
  const unknownReceipt = await eventually('none interrupted attempt becomes unknown', () => e.item(c, 110), (row) => row?.state === 'unknown');
  assert.equal(e.targets.rows('b').filter((row) => row.line === 'C').length, 1);
  assert.equal(e.receipts().some((row) => row.effect_key === unknownReceipt.effect_key), false);
  // C and B own the same persisted source points. One real read legitimately
  // reaches both groups; verify B's independent effect rather than ignoring it.
  await e.committed(b, 110);
  await disable(c); await e.close(130);
  e.check('LocalReceiptFailureWithoutDedupeNeverResends', { before_kill: noneReceipt, unknown: unknownReceipt, effect: e.proof(c, 110, 'b'), healthy_shared_source_effect: e.proof(b, 110, 'b') });

  await e.fault('target_commit_response_lost', b, true);
  await e.feed(138, b, 'none-response-one-thirty-eight'); await e.close(140);
  const noneUnknown = await eventually('no dedupe ambiguous commit is unknown', () => e.item(b, 130), (row) => row?.state === 'unknown');
  const noneFault = (await e.state()).faults.find((fault) => fault.kind === 'target_commit_response_lost' && fault.group_id === b.id);
  assert.ok(noneFault.enabled && noneFault.reached >= 1 && noneFault.commits >= 1);
  const targetCount = e.targets.rows('b').length;
  await e.restart(140);
  await e.feed(148, b, 'successor-blocked-by-unknown'); await e.close(150);
  await eventually('later accepted row stays behind unknown effect', () => e.item(b, 140), (row) => row?.state === 'pending');
  assert.equal(e.targets.rows('b').length, targetCount);
  assert.equal(e.item(b, 130).state, 'unknown');
  assert.equal(e.receipts().some((row) => row.effect_key === noneUnknown.effect_key), false);
  e.check('ResponseLostWithoutDedupeUnknownAndNoBlindInsert', { fault: noneFault, unknown: noneUnknown, successor: e.item(b, 140), target_rows: e.targets.rows('b'), ui: await e.ui(b, 'none-unknown') });
  return { a, b, c };
}
