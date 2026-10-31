import { expect, it } from 'vitest';
import { parseWriteGroupCalculationCapabilities } from '@/utils/studioV2WriteGroupCalculationJson';

const restriction = { method: 'delta', source_transform: 'exact_counter_scale', phase: 'draft', reason: 'mixed_exact_source_draft_unsupported' };
function capabilities(unsupported: unknown[]) {
  return {
    version: 1, result_shape: 'member-calculation-v1', draft_supported: true,
    activation_supported: false, editor_supported: false, reason: 'production_acceptance_pending',
    methods: ['latest', 'mean', 'delta'],
    cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1,
      max_seconds: 9223372036, max_minutes: 153722867, window_bounds_checked: true },
    schema_dialects: ['sqlite', 'postgres'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'],
    unsupported_combinations: unsupported,
  };
}

it('retains an omitted layout as all layouts without inventing a default', () => {
  const result = parseWriteGroupCalculationCapabilities(capabilities([restriction]));
  expect(result?.unsupported_combinations[0]).toEqual(restriction);
  expect(result?.unsupported_combinations[0]).not.toHaveProperty('layout');
});

it.each(['managed_wide_profile_v1', 'managed_long_v1'])('preserves the explicit %s restriction detached from input', layout => {
  const input = { ...restriction, layout };
  const result = parseWriteGroupCalculationCapabilities(capabilities([input]));
  expect(result?.unsupported_combinations[0]).toEqual(input);
  expect(result?.activation_supported).toBe(false);
  expect(result?.editor_supported).toBe(false);
  input.layout = 'unknown';
  expect(result?.unsupported_combinations[0]).toHaveProperty('layout', layout);
});

it.each([null, undefined, '', 'unknown', 'managed_wide_v1', 1, ['managed_long_v1']])('rejects an invalid present layout %#', layout => {
  expect(parseWriteGroupCalculationCapabilities(capabilities([{ ...restriction, layout }]))).toBeNull();
});

it('rejects extra keys even on a layout-qualified restriction', () => {
  expect(parseWriteGroupCalculationCapabilities(capabilities([{ ...restriction, layout: 'managed_long_v1', allowed: true }]))).toBeNull();
});

it('permits distinct non-overlapping layout restrictions', () => {
  const restrictions = ['managed_wide_profile_v1', 'managed_long_v1'].map(layout => ({ ...restriction, layout }));
  expect(parseWriteGroupCalculationCapabilities(capabilities(restrictions))?.unsupported_combinations).toEqual(restrictions);
});

it.each([
  [restriction, { ...restriction, layout: 'managed_long_v1' }],
  [{ ...restriction, layout: 'managed_long_v1' }, restriction],
  [{ ...restriction, layout: 'managed_long_v1' }, { ...restriction, layout: 'managed_long_v1' }],
])('rejects overlapping restrictions instead of silently dropping one %#', (...restrictions) => {
  expect(parseWriteGroupCalculationCapabilities(capabilities(restrictions))).toBeNull();
});

it('does not reinterpret omitted layout using an inherited property', () => {
  const inherited = Object.assign(Object.create({ layout: 'managed_long_v1' }) as Record<string, unknown>, restriction);
  expect(parseWriteGroupCalculationCapabilities(capabilities([inherited]))?.unsupported_combinations[0]).toEqual(restriction);
});
