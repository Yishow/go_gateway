import { describe, expect, it } from 'vitest';
import { mixedDraftIssues, mixedDraftToRequestGroup, newMixedDraft } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { draftToRequestGroup, type EditorDraft, type EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
const destination = { connector_id: 'db', connector_revision: 'db-r1', table_schema: '', table_name: '' };
const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', draft_supported: true, activation_supported: false, editor_supported: false, reason: 'not_readiness', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 3600, max_minutes: 60, window_bounds_checked: true }, schema_dialects: ['sqlite', 'postgres'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [{ method: 'delta', source_transform: 'exact_counter_scale', phase: 'draft', layout: 'managed_wide_profile_v1', reason: 'exact_wide_unsupported' }] };
function member(id: string, method: EditorMember['calculation'] = 'latest'): EditorMember {
  return { key: `p-${id}|t-${id}`, device_id: 'D', point_id: `p-${id}`, tag_id: `t-${id}`, entity_key: '', target_column: `v_${id}`, required: true, calculation: method, cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, source_revision: 'source-r1', mapping_revision: 'map-r1', value_role: 'gauge', source_value_type: 'float32', result_value_type: 'float64' };
}
function draft(): EditorDraft { return { ...destination, name: 'Mixed', record_id: '11111111-1111-4111-8111-111111111111', storage_strategy: 'managed', layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, allowed_lateness_seconds: 0, incomplete_policy: 'partial', entity_key_column: '', provenance_column: '', dedupe_capability: '', members: [member('a', 'mean'), member('b', 'mean'), { ...member('c', 'delta'), value_role: 'counter_raw', measurement_id: 'confirmed-counter', cadence: { interval: 1, unit: 'minute', alignment: 'utc_epoch' } }] }; }
describe('recording mixed draft canonical adapter', () => {
  it('creates only a fresh independent mixed shell with no legacy group interval, borrowed bindings or members', () => {
    const next = newMixedDraft(destination, '22222222-2222-4222-8222-222222222222');
    expect(next).toMatchObject({ record_id: '22222222-2222-4222-8222-222222222222', members: [], storage_strategy: 'managed', layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, incomplete_policy: 'partial' });
    expect(next).not.toHaveProperty('interval_seconds'); expect(next).not.toHaveProperty('record_key_column'); expect(next).not.toHaveProperty('due_at_column');
    expect(newMixedDraft(destination, '33333333-3333-4333-8333-333333333333', 'managed_long_v1').record_id).toBe('33333333-3333-4333-8333-333333333333');
  });
  it.each(['', 'not-a-uuid', '00000000-0000-0000-0000-000000000000'])('rejects noncanonical independent record identity %s', id => {
    expect(() => newMixedDraft(destination, id)).toThrow();
    const value = draft(); value.record_id = id; expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('mixed-record-id-invalid');
  });
  it('maps independent mean/mean/delta windows without synthetic group cadence or activation authority', () => {
    const value = draft(); expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual([]);
    const wire = mixedDraftToRequestGroup(value, 'W', capabilities);
    expect(wire.record_id).toBe('11111111-1111-4111-8111-111111111111'); expect(wire.workspace_id).toBe('W');
    expect(wire.members.map(m => m.calculation)).toEqual(['mean', 'mean', 'delta']);
    expect(wire.members[2].cadence).toEqual({ interval: 1, unit: 'minute', alignment: 'utc_epoch' });
    expect(wire.row_policy).not.toHaveProperty('interval_seconds'); expect(wire).not.toHaveProperty('status'); expect(wire).not.toHaveProperty('readiness');
    expect(() => draftToRequestGroup(value, 'W')).toThrow();
  });
  it('retains revisions, measurement, roles and complete saved layout intent with detached nested values', () => {
    const value = draft(); Object.assign(value, { due_at_column: 'due_at', record_key_column: 'record_id', group_id_column: 'group_id', device_id_column: 'device_id', provenance_column: 'provenance', group_key_columns: ['group_id'], unique_key_columns: ['record_id'], value_column: 'value', quality_column: 'quality', dedupe_capability: 'receipt', recreate_missing: false });
    const wire = mixedDraftToRequestGroup(value, 'W', capabilities);
    expect(wire.members[2]).toMatchObject({ source_revision: 'source-r1', mapping_revision: 'map-r1', measurement_id: 'confirmed-counter', value_role: 'counter_raw', source_value_type: 'float32', result_value_type: 'float64' });
    expect(wire.row_policy).toMatchObject({ due_at_column: 'due_at', record_key_column: 'record_id', group_id_column: 'group_id', device_id_column: 'device_id', value_column: 'value', quality_column: 'quality', unique_key_columns: ['record_id'] });
    expect(wire.write_policy).toEqual({ dedupe_capability: 'receipt', recreate_missing: false });
    value.members[2].cadence!.interval = 50; value.unique_key_columns!.push('mutated');
    expect(wire.members[2].cadence!.interval).toBe(1); expect(wire.row_policy.unique_key_columns).toEqual(['record_id']);
  });
  it.each([0, -1, 1.5, NaN, Infinity, 3601])('rejects invalid or capability-exceeding cadence %s without coercion', interval => {
    const value = draft(); value.members[0].cadence!.interval = interval;
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('member-cadence-invalid');
    expect(() => mixedDraftToRequestGroup(value, 'W', capabilities)).toThrow();
  });
  it('rejects missing methods, unsupported units and absent source identity', () => {
    const value = draft(); delete value.members[0].calculation; value.members[1].cadence!.unit = 'minute'; value.members[1].cadence!.interval = 61; value.members[2].point_id = '';
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual(expect.arrayContaining(['member-method-unsupported', 'member-cadence-invalid', 'member-source-required']));
  });
  it('refuses legacy conversion, fake group interval, unsupported layout/storage/dialect and absent capabilities', () => {
    const legacy = draft(); legacy.recording_profile = { kind: 'instantaneous', interval: 60, unit: 'second' };
    expect(() => mixedDraftToRequestGroup(legacy, 'W', capabilities)).toThrow();
    const value = draft(); value.interval_seconds = 60; value.storage_strategy = 'custom'; value.layout_version = 'managed_wide_v1';
    expect(mixedDraftIssues(value, capabilities, 'mysql')).toEqual(expect.arrayContaining(['mixed-group-interval-forbidden', 'mixed-storage-unsupported', 'mixed-layout-unsupported', 'mixed-dialect-unsupported']));
    expect(mixedDraftIssues(draft(), { ...capabilities, draft_supported: false }, 'sqlite')).toContain('mixed-capability-unavailable');
  });
  it('does not copy invalid cadence input text or other UI-only member state into the canonical request', () => {
    const value = draft(); Object.assign(value.members[0], { cadence_input: '1e', baseline: 99, readiness: true });
    const wire = mixedDraftToRequestGroup(value, 'W', capabilities);
    expect(wire.members[0].cadence).toEqual({ interval: 15, unit: 'second', alignment: 'utc_epoch' });
    expect(wire.members[0]).not.toHaveProperty('cadence_input'); expect(wire.members[0]).not.toHaveProperty('baseline');
    expect(wire.members[0]).not.toHaveProperty('readiness'); expect(wire.members[0]).not.toHaveProperty('key');
  });
  it('rejects empty name/members and nonfinite or negative lateness without mutating the draft', () => {
    const value = draft(); value.name = ' '; value.members = []; value.allowed_lateness_seconds = NaN;
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual(expect.arrayContaining(['name-required', 'members-required', 'lateness-invalid']));
    value.allowed_lateness_seconds = -1; expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('lateness-invalid');
    value.allowed_lateness_seconds = Infinity; expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('lateness-invalid');
    expect(value.name).toBe(' '); expect(Number.isFinite(value.allowed_lateness_seconds)).toBe(false);
  });
  it('preserves distinct canonical device sources and rejects duplicate canonical source membership', () => {
    const value = draft(); value.members[1].device_id = 'other-device';
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual([]);
    expect(mixedDraftToRequestGroup(value, 'W', capabilities).members[1].device_id).toBe('other-device');
    value.members.push({ ...value.members[0], key: 'misleading-display-key' });
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('member-source-duplicate');
  });
  it('prevents nonnumeric mean and unbound or noncounter delta without inventing measurement policies', () => {
    const value = draft(); value.members[0].source_value_type = 'bool'; value.members[1].value_role = 'state'; delete value.members[2].measurement_id;
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual(expect.arrayContaining(['member-mean-inapplicable', 'member-delta-unconfirmed']));
  });
  it('applies exact-source restrictions by method/layout and preserves exact proof unchanged on supported long delta', () => {
    const value = draft(); const m = value.members[2]; m.source_value_type = 'uint64'; m.result_value_type = 'decimal'; m.exact_recording_source = { measurement_id: 'confirmed-counter', definition_revision: 'definition-r', source_binding_revision: 'binding-r', series_epoch: 'epoch', original_data_type: 'uint64', scale: '0.1', unit: 'kWh', semantic_kind: 'counter', counter_policy: { allow_negative: false, rollover_modulus: '18446744073709551616' } };
    expect(mixedDraftIssues(value, capabilities, 'sqlite')).toContain('exact_wide_unsupported');
    value.layout_version = 'managed_long_v1'; expect(mixedDraftIssues(value, capabilities, 'sqlite')).toEqual([]);
    const wire = mixedDraftToRequestGroup(value, 'W', capabilities); expect(wire.members[2].exact_recording_source).toEqual(m.exact_recording_source);
    m.exact_recording_source.counter_policy!.rollover_modulus = '10'; expect(wire.members[2].exact_recording_source!.counter_policy!.rollover_modulus).toBe('18446744073709551616');
    const allLayouts = { ...capabilities, unsupported_combinations: capabilities.unsupported_combinations.map(({ layout: _layout, ...restriction }) => restriction) };
    expect(mixedDraftIssues(value, allLayouts, 'sqlite')).toContain('exact_wide_unsupported');
  });
});
