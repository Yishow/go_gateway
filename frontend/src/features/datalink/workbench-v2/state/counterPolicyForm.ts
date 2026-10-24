import type { MeasurementDefinition } from '../../../../types/measurement';
import { isConfirmedCounterDefinition } from './circuitCounterConfirmation';

export interface CounterPolicyForm {
  /** Explicit choice: a signed counter may report negative differences. */
  signed: boolean;
  /** Exact decimal text; never passed through a number. */
  modulus: string;
  rate: string;
}

export type CounterPolicyIssue = 'modulus_invalid' | 'rate_invalid' | 'rate_needs_modulus' | 'rollover_with_signed' | 'rollover_with_float';

/** Positive exact decimal in the grammar the backend accepts. */
const POSITIVE_EXACT = /^\+?(?:\d+(?:\.\d*)?|\.\d+)$/;
const MAX_DIGITS_LENGTH = 64;

const positiveExact = (value: string) => value.length <= MAX_DIGITS_LENGTH && POSITIVE_EXACT.test(value) && /[1-9]/.test(value);

export const isCounterDefinition = (def: MeasurementDefinition) => def.semantic_kind === 'counter' || def.semantic_kind === 'signed_counter';

export const isFloatEncoding = (def: MeasurementDefinition) => (def.numeric_encoding ?? '').trim().toLowerCase().startsWith('float');

/** A source as a saved group holds it, including the exact measurement definition it is bound to. */
export interface PolicyMember {
  device_id: string;
  point_id: string;
  tag_id: string;
  measurement_id?: string;
}

export type PolicySourceState =
  /** The bound definition, whether its saved policy already counts as confirmed. */
  | { kind: 'editable'; def: MeasurementDefinition; confirmed: boolean }
  /** No usable bound definition: it is confirmed in Basic, never guessed or created here. */
  | { kind: 'repair' };

const filled = (value: string | undefined) => typeof value === 'string' && value.trim() !== '';

/**
 * The definition a saved member is bound to, by its persisted measurement ID and only if that definition belongs to
 * the same workspace and source. A missing, foreign or incomplete definition (no revisions, epoch or unit) is never
 * replaced by another definition of the same source: the runtime counts from exactly the bound one.
 */
export function resolvePolicyDefinition(workspaceId: string, member: PolicyMember, definitions: MeasurementDefinition[]): PolicySourceState {
  const bound = member.measurement_id?.trim();
  const def = bound ? definitions.find((item) => item.id === bound) : undefined;
  if (!def || def.workspace_id !== workspaceId || def.device_id !== member.device_id || def.point_id !== member.point_id ||
    def.tag_id !== member.tag_id || !isCounterDefinition(def)) return { kind: 'repair' };
  if (![def.id, def.definition_revision, def.source_binding_revision, def.series_epoch, def.unit].every(filled)) return { kind: 'repair' };
  return { kind: 'editable', def, confirmed: isConfirmedCounterDefinition(def) };
}

export function formFromDefinition(def: MeasurementDefinition): CounterPolicyForm {
  return { signed: def.semantic_kind === 'signed_counter', modulus: def.counter_policy?.rollover_modulus ?? '', rate: def.counter_policy?.max_rate_per_second ?? '' };
}

/** Local explanation only: the backend stays the authority and may refuse more. */
export function validateCounterPolicy(def: MeasurementDefinition, form: CounterPolicyForm): CounterPolicyIssue[] {
  const modulus = form.modulus.trim();
  const rate = form.rate.trim();
  const issues: CounterPolicyIssue[] = [];
  if (modulus && !positiveExact(modulus)) issues.push('modulus_invalid');
  if (rate && !positiveExact(rate)) issues.push('rate_invalid');
  if (rate && !modulus) issues.push('rate_needs_modulus');
  if ((modulus || rate) && form.signed) issues.push('rollover_with_signed');
  if (isFloatEncoding(def)) issues.push('rollover_with_float');
  return issues;
}

/** A float counter's note is always listed but only blocks a save when a wrap policy is actually set. */
export function blockingIssues(issues: CounterPolicyIssue[], form: CounterPolicyForm): CounterPolicyIssue[] {
  const hasRollover = form.modulus.trim() !== '' || form.rate.trim() !== '';
  return issues.filter((issue) => issue !== 'rollover_with_float' || hasRollover);
}

export function counterPolicyChanged(def: MeasurementDefinition, form: CounterPolicyForm): boolean {
  const saved = formFromDefinition(def);
  return saved.signed !== form.signed || saved.modulus !== form.modulus.trim() || saved.rate !== form.rate.trim();
}

/**
 * The full saved definition with only the reviewed policy changed. Identity, binding, numeric encoding, unit and the
 * legacy float policy fields are carried over untouched; empty exact fields are removed rather than sent as zero.
 */
export function buildCounterPolicyDefinition(def: MeasurementDefinition, form: CounterPolicyForm): MeasurementDefinition {
  const policy = { ...def.counter_policy, allow_negative: form.signed };
  const modulus = form.modulus.trim();
  const rate = form.rate.trim();
  if (modulus) policy.rollover_modulus = modulus; else delete policy.rollover_modulus;
  if (rate) policy.max_rate_per_second = rate; else delete policy.max_rate_per_second;
  return { ...def, semantic_kind: form.signed ? 'signed_counter' : 'counter', counter_policy: policy };
}
