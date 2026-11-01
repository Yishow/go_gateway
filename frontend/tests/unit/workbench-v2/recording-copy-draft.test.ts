import { describe, expect, it } from 'vitest';
import { recordingCopyPrefill, createRecordingCopyDraft } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingCopyDraft';
import { memberCalculationReason } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingMemberCalculation';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
import { savedGroup } from '../../fixtures/writeGroupState';

const capability: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
const destination = { connector_id: 'current-connector', connector_revision: 'current-revision', table_schema: 'current-schema', table_name: 'must-not-reuse' };
const freshId = '22222222-2222-4222-8222-222222222222';
function source(overrides: Partial<WriteGroup> = {}): WriteGroup {
  const base = savedGroup();
  const member = { ...base.members[0], calculation: 'mean' as const, cadence: { interval: 15, unit: 'second' as const, alignment: 'utc_epoch' as const },
    source_value_type: 'float64' as const, value_role: 'gauge' as const, result_value_type: 'float64' as const,
    gauge_confirmation: 'old-gauge-proof', lineage: 'old-lineage', baseline: 'old-baseline' };
  return { ...base, record_id: '11111111-1111-4111-8111-111111111111', status: 'running', applied_revision: 'old-applied',
    destination: { ...base.destination, storage_strategy: 'managed', schema_revision: 'old-schema', schema_digest: 'old-digest' },
    layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
    row_policy: { allowed_lateness_seconds: 3, incomplete_policy: 'partial', record_key_column: 'record_id', group_key_columns: ['old-group'], value_column: 'old-value' },
    write_policy: { recreate_missing: true, dedupe_capability: 'old-dedupe' }, migration: { source_ids: ['old-lineage'] }, members: [member], ...overrides };
}
function prefill(group = source()) {
  const value = recordingCopyPrefill(group, capability, 'sqlite');
  expect(value).toBeDefined();
  if (!value) throw new Error('Expected a supported canonical copy source');
  return value;
}
it('projects only configuration intent and drops all source authority and proof fields', () => {
  const value = prefill();
  expect(Object.keys(value).sort()).toEqual(['allowedLatenessSeconds', 'incompletePolicy', 'layout', 'members', 'sourceName']);
  expect(value.members).toEqual([{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', calculation: 'mean',
    cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, required: true, source_value_type: 'float64', value_role: 'gauge' }]);
});
it('creates a new identity with current destination and empty server-owned bindings without changing the original', () => {
  const original = source(); const snapshot = structuredClone(original);
  const value = prefill(original); const draft = createRecordingCopyDraft(value, destination, freshId, 'Line A copy');
  expect(draft).toMatchObject({ name: 'Line A copy', record_id: freshId, connector_id: 'current-connector', connector_revision: 'current-revision',
    table_schema: 'current-schema', table_name: '', allowed_lateness_seconds: 3, incomplete_policy: 'partial', members: [{ key: 'pt-0|tag-0', target_column: '', entity_key: '', calculation: 'mean' }] });
  for (const key of ['id', 'revision', 'applied_revision', 'schema_digest', 'migration', 'group_key_columns', 'record_key_column', 'value_column', 'recreate_missing']) expect(draft).not.toHaveProperty(key);
  for (const key of ['source_revision', 'mapping_revision', 'result_value_type', 'gauge_confirmation', 'baseline', 'lineage']) expect(draft.members[0]).not.toHaveProperty(key);
  expect(draft.dedupe_capability).toBe(''); expect(original).toEqual(snapshot);
  draft.members[0].cadence!.interval = 30;
  expect(value.members[0].cadence?.interval).toBe(15); expect(original.members[0].cadence?.interval).toBe(15);
});
it('preserves long-layout intent without copying its shared table or managed columns', () => {
  const draft = createRecordingCopyDraft(prefill(source({ layout_version: 'managed_long_v1', storage_mode: 'plain' })), destination, freshId, 'Long copy');
  expect(draft.layout_version).toBe('managed_long_v1'); expect(draft.table_name).toBe(''); expect(draft.record_id).toBe(freshId);
  expect(draft).not.toHaveProperty('record_key_column');
});
it('preserves delta intent but drops exact proof and requires fresh counter validation', () => {
  const base = source(); const group = source({ members: [{ ...base.members[0], calculation: 'delta', source_value_type: 'uint64', value_role: 'counter_raw', measurement_id: 'counter-1',
    exact_recording_source: { measurement_id: 'counter-1', definition_revision: 'old-definition', source_binding_revision: 'old-binding', series_epoch: 'old-epoch', original_data_type: 'uint64', scale: '1', unit: 'kWh', semantic_kind: 'counter' } }] });
  const draft = createRecordingCopyDraft(prefill(group), destination, freshId, 'Counter copy');
  expect(draft.members[0]).toMatchObject({ calculation: 'delta', measurement_id: 'counter-1' });
  expect(draft.members[0]).not.toHaveProperty('exact_recording_source');
  expect(memberCalculationReason(draft.members[0], 'delta', { workspaceId: 'ws-1', recordId: freshId, definitions: [], definitionsAvailable: true,
    capabilities: capability, layout: draft.layout_version!, dialect: 'sqlite' })).toBe('counter_required');
});
describe('unsupported sources never silently become a different recording kind', () => {
  it.each(['deleted', 'custom', 'legacy', 'timescale', 'native_monthly'] as const)('rejects %s', kind => {
    const group = source();
    if (kind === 'deleted') group.status = kind;
    else if (kind === 'custom') group.destination.storage_strategy = kind;
    else if (kind === 'legacy') group.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second' };
    else group.storage_mode = kind;
    expect(recordingCopyPrefill(group, capability, 'sqlite')).toBeUndefined();
  });
  it.each(['missing', 'editor', 'draft', 'dialect', 'layout', 'method'] as const)('rejects %s capability mismatch', kind => {
    const cap = structuredClone(capability);
    if (kind === 'editor') cap.editor_supported = false;
    if (kind === 'draft') cap.draft_supported = false;
    if (kind === 'layout') cap.schema_layouts = ['managed_long_v1'];
    if (kind === 'method') cap.methods = ['latest'];
    expect(recordingCopyPrefill(source(), kind === 'missing' ? undefined : cap, kind === 'dialect' ? 'postgres' : 'sqlite')).toBeUndefined();
  });
});
