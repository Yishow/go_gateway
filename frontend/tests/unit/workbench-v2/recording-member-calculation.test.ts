import { describe, expect, it } from 'vitest';
import { memberCalculationReason, changeMemberCalculation, canConfirmMemberGauge, confirmMemberGauge, type MemberCalculationContext } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingMemberCalculation';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { MeasurementDefinition } from '../../../src/types/measurement';
const member: EditorMember = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', target_column: '', entity_key: '', required: true, source_value_type: 'float32' };
const gauge: MeasurementDefinition = { id: 'M', workspace_id: 'W', device_id: 'D', point_id: 'P', tag_id: 'T', equipment_id: 'E', name: 'Temperature', quantity: 'temperature', unit: 'C', semantic_kind: 'gauge', numeric_encoding: 'float32', definition_revision: 'd1', source_binding_revision: 'b1', series_epoch: 'e1' };
const context: MemberCalculationContext = { workspaceId: 'W', recordId: 'record-1', definitions: [], definitionsAvailable: true, layout: 'managed_wide_profile_v1', dialect: 'sqlite', capabilities: { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] } };
describe('per-member choices use canonical source semantics, not raw type guesses', () => {
 it('keeps raw latest available without mandatory measurement or activation readiness', () => { expect(memberCalculationReason(member, 'latest', context)).toBeUndefined(); });
 it('does not replace editor support with draft support or a known dialect guess', () => {
  expect(memberCalculationReason(member, 'latest', { ...context, capabilities: { ...context.capabilities!, editor_supported: false } })).toBe('capability');
  expect(memberCalculationReason(member, 'latest', { ...context, dialect: 'mysql' })).toBe('dialect');
  expect(memberCalculationReason(member, 'latest', { ...context, capabilities: undefined })).toBe('capability');
 });
 it('offers mean only for a unique complete numeric gauge definition of this exact source', () => {
  expect(memberCalculationReason(member, 'mean', { ...context, definitions: [gauge] })).toBeUndefined();
  expect(memberCalculationReason({ ...member, source_value_type: undefined }, 'mean', { ...context, definitions: [gauge] })).toBe('gauge_required');
  for (const definitions of [[], [gauge, { ...gauge, id: 'other' }], [{ ...gauge, workspace_id: 'foreign' }], [{ ...gauge, semantic_kind: 'state' as const }], [{ ...gauge, definition_revision: '' }]]) expect(memberCalculationReason(member, 'mean', { ...context, definitions })).toBe('gauge_required');
  expect(memberCalculationReason({ ...member, measurement_id: 'missing' }, 'mean', { ...context, definitions: [gauge] })).toBe('gauge_required');
 });
 it('requires confirmed counter semantics but preserves signed or rollover policies as eligible', () => {
  for (const definition of [{ ...gauge, semantic_kind: 'counter' as const }, { ...gauge, semantic_kind: 'signed_counter' as const, counter_policy: { allow_negative: true } }, { ...gauge, semantic_kind: 'counter' as const, numeric_encoding: 'uint32', counter_policy: { allow_negative: false, rollover_modulus: '4294967296' } }]) expect(memberCalculationReason(member, 'delta', { ...context, definitions: [definition] })).toBeUndefined();
  expect(memberCalculationReason(member, 'delta', { ...context, definitions: [gauge] })).toBe('counter_required');
 });
 it('honors every-layout exact source restrictions without treating a capability as record readiness', () => {
  const exact = { ...member, exact_recording_source: { measurement_id: 'M', definition_revision: 'd1', source_binding_revision: 'b1', series_epoch: 'e1', original_data_type: 'uint64' as const, scale: '0.1', unit: 'kWh', semantic_kind: 'counter' as const } };
  const scoped = { ...context, capabilities: { ...context.capabilities!, unsupported_combinations: [{ method: 'latest' as const, source_transform: 'exact_counter_scale' as const, phase: 'runtime' as const, reason: 'exact_latest_unsupported' }] } };
  expect(memberCalculationReason(exact, 'latest', scoped)).toBe('exact_latest_unsupported');
 });
 it('changes only this member intent and binds canonical semantics without carrying a stale result type', () => {
  const before = { ...member, source_revision: 's1', mapping_revision: 'm1', source_value_type: 'float32' as const, result_value_type: 'int64' as const, calculation: 'latest' as const, cadence: { interval: 15, unit: 'second' as const, alignment: 'utc_epoch' as const } };
  const after = changeMemberCalculation(before, 'mean', { ...context, definitions: [gauge] });
  expect(after).toMatchObject({ calculation: 'mean', measurement_id: 'M', value_role: 'gauge', source_value_type: 'float32', source_revision: 's1', mapping_revision: 'm1', cadence: before.cadence });
  expect(after).not.toHaveProperty('result_value_type'); expect(before.calculation).toBe('latest');
  expect(changeMemberCalculation(before, 'delta', { ...context, definitions: [gauge] })).toBe(before);
 });
});
it('does not invent a unit requirement for a canonical gauge, while counter confirmation still requires its unit', () => {
 const definition = { ...gauge, unit: undefined };
 expect(memberCalculationReason(member, 'mean', { ...context, definitions: [definition] })).toBeUndefined();
 expect(memberCalculationReason(member, 'delta', { ...context, definitions: [{ ...definition, semantic_kind: 'counter' }] })).toBe('counter_required');
});

describe('explicit per-record gauge authoring without global measurement creation', () => {
 it('requires explicit confirmation for a fully read unbound numeric source, not a default Latest choice', () => {
  expect(memberCalculationReason(member, 'mean', context)).toBe('gauge_required');
  expect(canConfirmMemberGauge(member, context)).toBe(true);
  const confirmed = confirmMemberGauge(member, context);
  expect(confirmed).toMatchObject({ value_role: 'gauge' }); expect(confirmed).not.toHaveProperty('measurement_id');
  expect(memberCalculationReason(confirmed, 'mean', context)).toBeUndefined(); expect(member.value_role).toBeUndefined();
 });
 it('does not let counter/state/conflicting or partial definitions be replaced by a gauge declaration', () => {
  for (const definitions of [[{ ...gauge, semantic_kind: 'counter' as const }], [{ ...gauge, semantic_kind: 'state' as const }], [gauge, { ...gauge, id: 'other' }], [{ ...gauge, definition_revision: '' }]]) {
   const read = { ...context, definitions }; expect(canConfirmMemberGauge(member, read)).toBe(false); expect(confirmMemberGauge(member, read)).toBe(member);
  }
  for (const value_role of ['counter_raw', 'state'] as const) expect(canConfirmMemberGauge({ ...member, value_role }, context)).toBe(false);
 });
 it('cannot clear a foreign, missing or non-gauge binding or an exact-source proof to gain mean eligibility', () => {
  const bound = { ...member, measurement_id: 'M' };
  for (const definitions of [[], [{ ...gauge, workspace_id: 'foreign' }], [{ ...gauge, semantic_kind: 'counter' as const }]]) {
   const read = { ...context, definitions }; expect(canConfirmMemberGauge(bound, read)).toBe(false); expect(confirmMemberGauge(bound, read)).toBe(bound);
  }
  const exact = { ...member, exact_recording_source: { measurement_id: 'M', definition_revision: 'd1', source_binding_revision: 'b1', series_epoch: 'e1', original_data_type: 'uint64' as const, scale: '0.1', unit: 'kWh', semantic_kind: 'counter' as const } };
  expect(canConfirmMemberGauge(exact, context)).toBe(false);
 });
 it('never confirms from unavailable reads or nonnumeric types and clears authority on a source change', () => {
  expect(canConfirmMemberGauge(member, { ...context, definitionsAvailable: false })).toBe(false);
  expect(canConfirmMemberGauge({ ...member, source_value_type: 'bool' }, context)).toBe(false);
  expect(canConfirmMemberGauge({ ...member, point_id: '' }, context)).toBe(false);
  const confirmed = confirmMemberGauge(member, context);
  expect(memberCalculationReason({ ...confirmed, point_id: 'other-point' }, 'mean', context)).toBe('gauge_required');
  expect(memberCalculationReason(confirmed, 'mean', { ...context, workspaceId: 'other-workspace' })).toBe('gauge_required');
  expect(memberCalculationReason(confirmed, 'mean', { ...context, recordId: 'other-record' })).toBe('gauge_required');
  expect(memberCalculationReason({ ...confirmed, mapping_revision: 'new-mapping' }, 'mean', context)).toBe('gauge_required');
 });
 it('preserves a canonical saved mean or explicitly saved gauge role without making ordinary saved Latest a gauge', () => {
  const saved = { ...member, source_revision: 'source-1', mapping_revision: 'mapping-1', calculation: 'mean' as const };
  expect(memberCalculationReason(saved, 'mean', context)).toBeUndefined();
  expect(memberCalculationReason(changeMemberCalculation(saved, 'latest', context), 'mean', context)).toBeUndefined();
  expect(memberCalculationReason({ ...saved, calculation: 'latest' }, 'mean', context)).toBe('gauge_required');
  expect(memberCalculationReason({ ...saved, calculation: 'latest', value_role: 'gauge' }, 'mean', context)).toBeUndefined();
 });
});
