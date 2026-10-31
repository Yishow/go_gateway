import type { WriteGroupCalculationCapabilities } from '../types/studioV2WriteGroupCalculation';
import { isRecord, readRequiredText } from './studioV2WriteGroupJsonPrimitives';

const METHODS = ['latest', 'mean', 'delta'] as const;
const UNITS = ['second', 'minute'] as const;
const DIALECTS = ['sqlite', 'postgres'] as const;
const LAYOUTS = ['managed_wide_profile_v1', 'managed_long_v1'] as const;
const TOP_KEYS = ['version', 'result_shape', 'draft_supported', 'activation_supported', 'editor_supported', 'reason',
  'methods', 'cadence', 'schema_dialects', 'schema_layouts', 'unsupported_combinations'];
const CADENCE_KEYS = ['alignment', 'units', 'minimum', 'max_seconds', 'max_minutes', 'window_bounds_checked'];
const RESTRICTION_KEYS = ['method', 'source_transform', 'phase', 'reason'];

function exactKeys(value: Record<string, unknown>, keys: readonly string[]): boolean {
  const actual = Object.keys(value);
  return actual.length === keys.length && actual.every(key => keys.includes(key));
}

function enumArray<T extends string>(value: unknown, allowed: readonly T[]): T[] | null {
  if (!Array.isArray(value) || value.length === 0 || value.length > allowed.length ||
    value.some(item => typeof item !== 'string' || !allowed.includes(item as T)) || new Set(value).size !== value.length) return null;
  return [...value] as T[];
}

function positiveInteger(value: unknown): value is number {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= 1;
}

/** Missing means unknown older backend. Present malformed data cannot become support.
 * This only preserves backend capabilities: it grants no record readiness or operation authority.
 */
export function parseWriteGroupCalculationCapabilities(value: unknown): WriteGroupCalculationCapabilities | undefined | null {
  if (value === undefined) return undefined;
  if (!isRecord(value) || !exactKeys(value, TOP_KEYS) || value.version !== 1 || value.result_shape !== 'member-calculation-v1' ||
    typeof value.draft_supported !== 'boolean' || typeof value.activation_supported !== 'boolean' || typeof value.editor_supported !== 'boolean') return null;
  const reason = readRequiredText(value, 'reason');
  const methods = enumArray(value.methods, METHODS);
  const dialects = enumArray(value.schema_dialects, DIALECTS);
  const layouts = enumArray(value.schema_layouts, LAYOUTS);
  const cadence = value.cadence;
  if (!reason || !methods || !dialects || !layouts || !isRecord(cadence) || !exactKeys(cadence, CADENCE_KEYS) ||
    cadence.alignment !== 'utc_epoch' || cadence.window_bounds_checked !== true ||
    !positiveInteger(cadence.minimum) || !positiveInteger(cadence.max_seconds) || !positiveInteger(cadence.max_minutes) ||
    cadence.max_seconds > 9_223_372_036 || cadence.max_minutes > Math.floor(cadence.max_seconds / 60) ||
    cadence.minimum > cadence.max_seconds || cadence.minimum > cadence.max_minutes) return null;
  const units = enumArray(cadence.units, UNITS);
  const rawRestrictions = value.unsupported_combinations;
  if (!units || !Array.isArray(rawRestrictions) || rawRestrictions.length > METHODS.length * 2 * LAYOUTS.length) return null;
  const restrictions: Array<WriteGroupCalculationCapabilities['unsupported_combinations'][number]> = [];
  const seen = new Map<string, Set<string>>();
  for (const raw of rawRestrictions) {
    if (!isRecord(raw) || !exactKeys(raw, Object.hasOwn(raw, 'layout') ? [...RESTRICTION_KEYS, 'layout'] : RESTRICTION_KEYS) ||
      (raw.method !== 'latest' && raw.method !== 'mean' && raw.method !== 'delta') ||
      raw.source_transform !== 'exact_counter_scale' || (raw.phase !== 'draft' && raw.phase !== 'runtime')) return null;
    const hasLayout = Object.hasOwn(raw, 'layout');
    const layout = hasLayout ? raw.layout : undefined;
    if (hasLayout && layout !== 'managed_wide_profile_v1' && layout !== 'managed_long_v1') return null;
    const restrictionReason = readRequiredText(raw, 'reason');
    const identity = `${raw.method}/${raw.phase}/${raw.source_transform}`;
    const scope = typeof layout === 'string' ? layout : '*';
    const scopes = seen.get(identity) ?? new Set<string>();
    if (!restrictionReason || scopes.has(scope) || scopes.has('*') || (scope === '*' && scopes.size > 0)) return null;
    scopes.add(scope);
    seen.set(identity, scopes);
    restrictions.push({ method: raw.method, source_transform: raw.source_transform, phase: raw.phase, reason: restrictionReason,
      ...(layout === 'managed_wide_profile_v1' || layout === 'managed_long_v1' ? { layout } : {}) });
  }
  return {
    version: 1, result_shape: 'member-calculation-v1', draft_supported: value.draft_supported,
    activation_supported: value.activation_supported, editor_supported: value.editor_supported, reason,
    methods, schema_dialects: dialects, schema_layouts: layouts,
    cadence: { alignment: 'utc_epoch', units, minimum: cadence.minimum, max_seconds: cadence.max_seconds,
      max_minutes: cadence.max_minutes, window_bounds_checked: true },
    unsupported_combinations: restrictions,
  };
}
