import { describe, expect, it } from 'vitest';
import { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { group, delivery } from './recording-home-guidance-fixtures';
import type { RecordingHomeGroupProjection } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
const receipt = { group_revision: 'G2', connector_revision: 'DB2', record_id: 'R', effect_key: 'E', payload_digest: 'a'.repeat(64), committed_at: '2026-10-30T10:32:00Z' };
function item(): { -readonly [K in keyof RecordingHomeGroupProjection]: RecordingHomeGroupProjection[K] } {
  const saved = { ...group('G'), applied_revision: 'G2' };
  return { group: saved, lifecycle: { state: 'unavailable', reason: 'unknown' }, delivery: { state: 'available', value: {
    ...delivery(saved), last_sql_committed_effect: receipt, revision_stages: { group_revision: 'G2', stages: delivery(saved).stages },
  } } };
}
describe('home result facts never imply present health or command permission', () => {
  it('labels matching applied SQL independently from unproven owners/readback/current health', () => {
    expect(recordingHomeResultSummary(item())).toMatchObject({ saved: true, configuration: 'requested', sql: 'revision_receipt',
      lastCommittedAt: receipt.committed_at, current: 'unknown', readback: 'unknown', owners: 'unknown' });
  });
  it('receipt matching old applied revision never blesses a newer desired revision', () => {
    const value = item(); value.group = { ...value.group, revision: 'G3' };
    expect(recordingHomeResultSummary(value)).toMatchObject({ configuration: 'not_requested', sql: 'previous_revision_receipt' });
  });
  it('start/intake acknowledgement without revision-scoped commits waits rather than claims success', () => {
    const value = item(); if (value.delivery.state !== 'available') throw new Error('fixture');
    value.delivery = { state: 'available', value: { ...value.delivery.value, last_sql_committed_at: null, last_sql_committed_effect: null,
      stages: { ...value.delivery.value.stages, sql_committed: 0 }, revision_stages: { group_revision: 'G2', stages: { ...value.delivery.value.stages, sql_committed: 0 } } } };
    expect(recordingHomeResultSummary(value)).toMatchObject({ sql: 'unknown', collection: 'unknown' });
    const runtime = effectiveGroup(); runtime.delivery!.current!.last_commit = null; runtime.delivery!.current!.stages.sql_committed = 0;
    expect(recordingHomeResultSummary(value, { runtime })).toMatchObject({ sql: 'awaiting_first_write', collection: 'accepting', current: 'unknown' });
  });
  it('unscoped historical timestamp is retained but is not a matching SQL receipt', () => {
    const value = item(); if (value.delivery.state !== 'available') throw new Error('fixture');
    value.delivery = { state: 'available', value: { ...value.delivery.value, last_sql_committed_effect: undefined, revision_stages: undefined } };
    expect(recordingHomeResultSummary(value)).toMatchObject({ sql: 'historical_commit', lastCommittedAt: receipt.committed_at });
  });
  it('failed current read may retain only separately identified history, not the prior live receipt', () => {
    const value = item(); value.delivery = { state: 'unavailable', reason: 'error' };
    expect(recordingHomeResultSummary(value, { history: { workspaceId: 'W1', groupId: 'G', committedAt: receipt.committed_at } }))
      .toMatchObject({ sql: 'historical_commit', current: 'unknown', lastCommittedAt: receipt.committed_at });
    expect(recordingHomeResultSummary(value, { history: { workspaceId: 'FOREIGN', groupId: 'G', committedAt: receipt.committed_at } }).lastCommittedAt).toBeNull();
  });
  it('rejects a mismatched connector receipt and never fabricates destination values', () => {
    const value = item(); if (value.delivery.state !== 'available') throw new Error('fixture');
    value.delivery = { state: 'available', value: { ...value.delivery.value, last_sql_committed_effect: { ...receipt, connector_revision: 'OLD' } } };
    expect(recordingHomeResultSummary(value).sql).toBe('historical_commit');
    expect(JSON.stringify(recordingHomeResultSummary(value))).not.toContain('payload_digest');
  });
});

import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
function effectiveGroup() {
  const runtime = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' }).groups[0];
  runtime.as_of = '2026-10-30T10:34:00Z';
  runtime.delivery = { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 0, counter_replay_quarantines: null,
    backlog: [], test_write: null, current: { binding: runtime.effective!.binding, totals_scope: 'group', shared: true,
      intake: { group_revision: 'A', state: 'active', reason: null }, stages: delivery(group('G')).stages,
      last_commit: { ...receipt, group_revision: 'A', connector_revision: '7', provenance: 'runtime' }, sources: [], source_state: 'unavailable',
      next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null } };
  return runtime;
}
describe('requested and effective boundaries are independent SQL facts', () => {
  it('requested B may still have effective A; a receipt for A is never B success', () => {
    const value = item(); value.group = { ...value.group, revision: 'B', applied_revision: 'B' };
    const result = recordingHomeResultSummary(value, { runtime: effectiveGroup() });
    expect(result).toMatchObject({ configuration: 'requested', effective: 'requested_not_effective', effectiveRevision: 'A',
      sql: 'effective_revision_receipt', historicalDestination: 'meter_old', current: 'unknown', owners: 'unknown' });
  });
  it('a later canonical boundary may establish B effective without claiming live device health', () => {
    const value = item(); value.group = { ...value.group, revision: 'B', applied_revision: 'B' };
    const runtime = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' }, 'B').groups[0];
    expect(recordingHomeResultSummary(value, { runtime })).toMatchObject({ effective: 'requested_effective', effectiveRevision: 'B', current: 'unknown' });
  });
  it('independent test-write readback and cleanup do not establish normal-row readback', () => {
    const runtime = effectiveGroup(); runtime.delivery!.test_write = { operation_id: 'TEST', provenance: 'test-write',
      written_at: receipt.committed_at, readback_at: receipt.committed_at, cleanup_at: null };
    expect(recordingHomeResultSummary(item(), { runtime })).toMatchObject({ readback: 'test_write_verified', cleanup: 'unknown' });
    runtime.delivery!.test_write.cleanup_at = receipt.committed_at;
    expect(recordingHomeResultSummary(item(), { runtime }).cleanup).toBe('test_write_completed');
  });
});

it('newer same-effective-source missing window preserves the previous SQL receipt and destination', () => {
  const runtime = effectiveGroup();
  runtime.delivery!.current!.source_state = 'available';
  runtime.delivery!.current!.sources = [{ device_id: 'D1', point_id: 'P', tag_id: 'T', last_committed: null,
    latest: { state: 'missing', window_start: '2026-10-30T10:32:30Z', window_end: '2026-10-30T10:33:00Z', delta_known: null, coverage: null, receipt: null } }];
  expect(recordingHomeResultSummary(item(), { runtime })).toMatchObject({ sql: 'effective_revision_receipt', current: 'source_attention',
    lastCommittedAt: '2026-10-30T10:32:00Z', historicalDestination: 'meter_old' });
  runtime.delivery!.current!.sources[0].last_committed = { ...runtime.delivery!.current!.sources[0].latest!, state: 'available', window_end: '2026-10-30T10:33:30Z' };
  expect(recordingHomeResultSummary(item(), { runtime }).current).toBe('unknown');
  runtime.delivery!.current!.sources[0].last_committed = null;
  runtime.delivery!.current!.sources[0].latest!.window_end = '2026-10-30T11:00:00Z';
  expect(recordingHomeResultSummary(item(), { runtime }).current).toBe('unknown');
});
it('historical aggregate error counts without effective revision proof never classify current group health', () => {
  const value = item(); if (value.delivery.state !== 'available') throw new Error('fixture');
  value.delivery = { state: 'available', value: { ...value.delivery.value, stages: { ...value.delivery.value.stages, blocked: 20 } } };
  expect(recordingHomeResultSummary(value).current).toBe('unknown');
});
