import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { parseWriteGroupDeliveryData } from '../../../src/utils/studioV2WriteGroupDeliveryJson';
import { GroupDeliveryStrip } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip';
import type { WriteGroupDelivery } from '../../../src/types/studioV2WriteGroupDelivery';
import { savedGroup } from '../../fixtures/writeGroupState';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
function fixture(): WriteGroupDelivery {
  return { group_id: 'group-1', intake: { state: 'blocked', reason: 'target-blocked' },
    stages: { collecting: 2, queued: 3, retrying: 1, blocked: 1, quarantined: 0, unknown: 1, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 42.5, no_data_buckets: 1, skipped_buckets: 2,
    backlog: [{ group_revision: 'rev-1', connector_id: 'C', connector_revision: 'C1', table_schema: 'main', table_name: 'meter_readings', pending: 4, error_codes: ['insert-failed', 'target-blocked'] }],
    attention: [{ effect_key: 'effect-1', state: 'blocked', state_revision: 7, payload_digest: 'digest-1', group_revision: 'rev-1', connector_id: 'C', connector_revision: 'C1', table_schema: 'main', table_name: 'meter_readings', error_code: 'insert-failed' }],
    quota: { configured: true, state: 'ok', scope: 'global', used_bytes: 10, max_bytes: 100, intake_refused: false } };
}
const timeFields = ['last-commit', 'receipt', 'bucket', 'storage-rebuilt', 'storage-event', 'source-start', 'source-end', 'source-commit'] as const;
function timestampFixture(field: typeof timeFields[number], timestamp: string): WriteGroupDelivery {
  const value = fixture(); const base = '2026-10-29T00:00:00Z';
  if (field === 'last-commit' || field === 'receipt') {
    value.stages.sql_committed = 1; value.last_sql_committed_at = field === 'last-commit' ? timestamp : base;
    if (field === 'receipt') value.last_sql_committed_effect = { group_revision: 'rev-1', connector_revision: 'C1', record_id: 'record-1', effect_key: 'effect-1', payload_digest: 'a'.repeat(64), committed_at: timestamp };
  }
  if (field === 'bucket') value.recent_bucket_issues = [{ group_revision: 'rev-1', bucket_start: timestamp, kind: 'skipped', causes: ['missing'] }];
  if (field === 'storage-rebuilt' || field === 'storage-event') value.storage = { generation: 2, history_may_be_lost: true, reconstructions: 1,
    last_reconstructed_at: field === 'storage-rebuilt' ? timestamp : base, events: [{ generation: 2, kind: 'reconstructed', occurrences: 1, occurred_at: field === 'storage-event' ? timestamp : base }] };
  if (field.startsWith('source-')) value.source_evidence = { state: 'available', group_revision: 'rev-1', connector_id: 'C', connector_revision: 'C1', sources: [{ device_id: 'D', point_id: 'P', tag_id: 'T', last_sql_committed: null,
    latest: { available: true, group_revision: 'rev-1', connector_revision: 'C1', record_id: 'record-1', effect_key: 'effect-1', payload_digest: 'a'.repeat(64), profile_kind: 'instantaneous', status: 'ok', raw_present: true, delivery: 'sql_committed',
      window_start: field === 'source-start' ? timestamp : base, window_end: field === 'source-end' ? timestamp : '2026-10-29T00:00:15Z', committed_at: field === 'source-commit' ? timestamp : base } }] };
  return value;
}
function mount(value: unknown) {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: value } });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  return render(<QueryClientProvider client={client}><GroupDeliveryStrip group={savedGroup({ id: 'group-1' })} readonly /></QueryClientProvider>);
}
beforeEach(async () => { await i18n.changeLanguage('en'); vi.clearAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });
describe('legacy delivery is a detached safe read DTO, not editable setup', () => {
  it('detaches named attention fields and does not retain hidden connector payload or mutate input', () => {
    const value = fixture(); const item = value.attention?.[0]; if (!item) throw new Error('Missing attention fixture');
    Object.assign(item, { connection_config: { password: 'fixture-password', database: '/private/fixture.sqlite' } }); Object.freeze(item);
    const parsed = parseWriteGroupDeliveryData(value); expect(parsed).not.toBeNull();
    expect(parsed?.attention?.[0]).toEqual({ effect_key: 'effect-1', state: 'blocked', state_revision: 7, payload_digest: 'digest-1', group_revision: 'rev-1', connector_id: 'C', connector_revision: 'C1', table_schema: 'main', table_name: 'meter_readings', error_code: 'insert-failed' });
    expect(parsed?.attention?.[0]).not.toBe(item); expect(JSON.stringify(parsed)).not.toContain('fixture-password');
    expect(item).toHaveProperty('connection_config');
  });
  it.each(['effect_key', 'group_revision', 'connector_id', 'connector_revision', 'payload_digest', 'table_name', 'table_schema'])('rejects unsafe attention %s rather than rewriting CAS or exposing it', (key) => {
    const value = fixture(); const item = value.attention?.[0]; if (!item) throw new Error('Missing attention');
    Object.assign(item, { [key]: '/private/fixture-secret' }); expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it.each(['group_id', 'group_revision', 'connector_id', 'connector_revision', 'table_name', 'table_schema'])('rejects unsafe group/backlog %s', (key) => {
    const value = fixture(); if (key === 'group_id') value.group_id = 'db.private.internal'; else Object.assign(value.backlog[0], { [key]: 'db.private.internal' });
    expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it.each(['group', 'backlog', 'receipt', 'revision', 'source'] as const)('does not trim or substitute the %s proof identity', (target) => {
    const value = fixture();
    if (target === 'group') value.group_id = ' group-1 ';
    if (target === 'backlog') value.backlog[0].connector_revision = ' C1 ';
    if (target === 'revision') value.revision_stages = { group_revision: ' rev-1 ', stages: { ...value.stages } };
    if (target === 'receipt') {
      value.stages.sql_committed = 1; value.last_sql_committed_at = '2026-10-29T00:00:00Z';
      value.last_sql_committed_effect = { group_revision: 'rev-1', connector_revision: 'C1', record_id: ' record-1 ', effect_key: 'effect-1', payload_digest: 'a'.repeat(64), committed_at: value.last_sql_committed_at };
    }
    if (target === 'source') value.source_evidence = { state: 'unavailable', unavailable_reason: 'not_long_layout', group_revision: ' rev-1 ', sources: [] };
    expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it('normalizes arbitrary diagnostics without changing stages, nullable facts or valid legacy codes', () => {
    const value = fixture(); value.intake.reason = 'driver exception password=fixture-password'; value.backlog[0].error_codes.push('driver exploded', 'future-unknown-code');
    const item = value.attention?.[0]; if (!item) throw new Error('Missing attention'); item.error_code = '/private/driver';
    value.quota.loss_risk_notice = 'Server=localhost;Database=metrics';
    const parsed = parseWriteGroupDeliveryData(value); expect(parsed).not.toBeNull();
    expect(parsed).toMatchObject({ stages: value.stages, last_sql_committed_at: null, oldest_pending_seconds: 42.5, intake: { reason: 'unavailable' }, quota: { loss_risk_notice: 'unavailable' } });
    expect(parsed?.attention?.[0].error_code).toBe('unavailable'); expect(parsed?.backlog[0].error_codes).toEqual(['insert-failed', 'target-blocked', 'unavailable', 'unavailable']);
    expect(value.intake.reason).toBe('driver exception password=fixture-password');
  });
  it('applies the read-only postgate to already parsed optional storage and source proof without helper rewrites', () => {
    const value = fixture(); value.storage = { generation: 1, history_may_be_lost: false, reconstructions: 0, repair_state: 'blocked', repair_reason: 'raw driver private diagnostic', events: [{ generation: 1, kind: 'repair_blocked', occurrences: 1, occurred_at: '2026-10-29T00:00:00Z', code: 'raw error', reason: '/private/fixture' }] };
    const parsed = parseWriteGroupDeliveryData(value); expect(parsed?.storage).toMatchObject({ repair_reason: 'unavailable', events: [{ code: 'unavailable', reason: 'unavailable' }] });
    value.source_evidence = { state: 'unavailable', unavailable_reason: 'not_long_layout', group_revision: '/private/proof', sources: [] };
    expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it('preserves published sender/storage/ownership/replay codes and the two producer quota notices', () => {
    const codes = ['insert-failed', 'storage-restore-in-progress', 'storage-recovery-unavailable', 'commit-ambiguous',
      'destination-identity-conflict', 'destination-unsupported', 'target-blocked', 'destination-rejected-row', 'target-unusable', 'target-unavailable',
      'payload-digest-mismatch', 'payload-undecodable', 'schema_repair_blocked', 'storage_owner_mismatch', 'repair_attempt_abandoned',
      'counter-replay-proof-conflict', 'recording-ownership-conflict', 'recovery-policy-unconfirmed'];
    for (const notice of ['new readings for this scope are not recorded until space is available; already accepted data is kept', 'capacity is running low; intake will be refused at the hard limit']) {
      const value = fixture(); value.backlog[0].error_codes = codes; value.quota.loss_risk_notice = notice;
      expect(parseWriteGroupDeliveryData(value)?.backlog[0].error_codes).toEqual(codes);
      expect(parseWriteGroupDeliveryData(value)?.quota.loss_risk_notice).toBe(notice);
    }
  });
  it('preserves published activation boundary diagnostics instead of degrading them to unknown', () => {
    const value = fixture(); value.backlog[0].error_codes = ['group-missing', 'sink-missing', 'group-not-applied', 'member-identity-incomplete',
      'tag-type-unsupported', 'layout-blocked', 'destination-scope-invalid', 'restore-failed', 'dedupe-unsupported', 'table-rebuilt', 'table-created-again'];
    expect(parseWriteGroupDeliveryData(value)?.backlog[0].error_codes).toEqual(value.backlog[0].error_codes);
  });
  it.each(['not_long_layout', 'applied_snapshot_unavailable', 'revision_mismatch'] as const)('preserves published unavailable source envelope %s and rejects unknown enums', (reason) => {
    const value = fixture(); value.source_evidence = { state: 'unavailable', unavailable_reason: reason, group_revision: 'rev-1', sources: [] };
    expect(parseWriteGroupDeliveryData(value)?.source_evidence).toEqual(value.source_evidence);
    Object.assign(value.source_evidence, { unavailable_reason: 'future-reason' }); expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it.each(['device_id', 'point_id', 'tag_id', 'circuit_id', 'record_id', 'effect_key', 'connector_revision', 'group_revision'])('rejects unsafe nested source proof %s without rewriting it', (key) => {
    const value = fixture();
    const record = { available: true as const, group_revision: 'rev-1', connector_revision: 'C1', record_id: 'record-1', effect_key: 'effect-1', payload_digest: 'a'.repeat(64), profile_kind: 'instantaneous' as const,
      window_start: '2026-10-29T00:00:00Z', window_end: '2026-10-29T00:00:15Z', status: 'partial' as const, raw_present: true, delivery: 'queued' as const,
      reason: 'future-source-code', aggregation: { method: 'arithmetic-mean' as const, version: 'arithmetic-mean-v1', basis: 'exact' as const, source_type: 'int16', result_type: 'decimal', valid_count: 1, rejected_count: 1, coverage: 'partial' as const, reason: 'future-aggregation-code' } };
    const source = { circuit_id: 'circuit-1', device_id: 'D', point_id: 'P', tag_id: 'T', latest: record, last_sql_committed: null };
    value.source_evidence = { state: 'available', group_revision: 'rev-1', connector_id: 'C', connector_revision: 'C1', sources: [source] };
    const parsed = parseWriteGroupDeliveryData(value); expect(parsed).not.toBeNull();
    const latest = parsed?.source_evidence?.sources[0].latest;
    expect(latest).toMatchObject({ reason: 'reason_unavailable', aggregation: { reason: 'reason_unavailable', valid_count: 1, rejected_count: 1 } });
    expect(record.reason).toBe('future-source-code');
    Object.assign(['device_id', 'point_id', 'tag_id', 'circuit_id'].includes(key) ? source : record, { [key]: '/private/proof' });
    expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it.each(['group_id', 'operation_id'])('rejects unsafe detached storage %s', (key) => {
    const value = fixture(); value.storage = { generation: 1, history_may_be_lost: false, reconstructions: 0, events: [{ generation: 1, kind: 'confirmed', occurrences: 1, occurred_at: '2026-10-29T00:00:00Z', [key]: '/private/proof' }] };
    expect(parseWriteGroupDeliveryData(value)).toBeNull();
  });
  it.each(timeFields)('rejects permissive Date.parse comments in %s without retaining private strings', (field) => {
    expect(parseWriteGroupDeliveryData(timestampFixture(field, '2026-10-29T08:00:00.123456789+08:00'))).not.toBeNull();
    const timestamp = '2026-10-29 (password=fixture-password /private/fixture.sqlite)';
    expect(Number.isFinite(Date.parse(timestamp))).toBe(true);
    expect(parseWriteGroupDeliveryData(timestampFixture(field, timestamp))).toBeNull();
  });
  it.each(timeFields)('preserves valid nanosecond and offset %s bytes exactly', (field) => {
    const timestamp = '2026-10-29T08:00:00.123456789+08:00'; const value = timestampFixture(field, timestamp);
    const parsed = parseWriteGroupDeliveryData(value); expect(parsed).not.toBeNull(); expect(JSON.stringify(parsed)).toContain(timestamp);
  });
  it('actual service refuses a hostile bucket time before text and datetime attributes render it', async () => {
    const mounted = mount(timestampFixture('bucket', '2026-10-29 (password=fixture-password /private/fixture.sqlite)'));
    await screen.findByTestId('group-delivery-failed'); expect(mounted.container.innerHTML).not.toMatch(/fixture-password|private\/fixture/);
    expect(screen.queryByTestId('delivery-bucket-causes')).not.toBeInTheDocument();
  });
  it('actual delivery service and native readonly strip preserve safe table and CAS values without mutation', async () => {
    const value = fixture(); mount(value); await screen.findByText('main.meter_readings');
    expect(screen.getByTestId('delivery-retry-effect-1')).toBeDisabled(); expect(screen.getByTestId('delivery-skip-effect-1')).toBeDisabled();
    const parsed = await studioV2WorkspaceWriteGroupsAPI.delivery('group-1'); expect(parsed.attention?.[0]).toEqual(value.attention?.[0]);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it.each(['table_name', 'table_schema', 'effect_key'])('actual service→native strip refuses hostile %s from old servers', async (key) => {
    const value = fixture(); const item = value.attention?.[0]; if (!item) throw new Error('Missing attention'); Object.assign(item, { [key]: '/private/fixture-secret' });
    const mounted = mount(value); await screen.findByTestId('group-delivery-failed');
    expect(mounted.container.innerHTML).not.toMatch(/fixture-secret|private\/|delivery-retry-/);
  });
});
