import type { MeasurementDefinition } from '../../../../types/measurement';

// Local explanation only: saved server readiness remains authoritative.
export function isConfirmedCounterDefinition(def: MeasurementDefinition): boolean {
  const policy = def.counter_policy;
  if (def.semantic_kind === 'counter') {
    if (policy?.allow_negative) return false;
  } else if (def.semantic_kind === 'signed_counter') {
    if (!policy?.allow_negative) return false;
  } else return false;
  if (![def.id, def.definition_revision, def.source_binding_revision, def.series_epoch, def.unit]
    .every((part) => typeof part === 'string' && part.trim().length > 0)) return false;
  const modulus = policy?.rollover_modulus ?? '';
  const rate = policy?.max_rate_per_second ?? '';
  if (!modulus && !rate) return true;
  // Match the backend decimal grammar without coercing exact evidence to Number.
  const positiveExact = (value: string) => value.length <= 64 && /^\+?(?:\d+(?:\.\d*)?|\.\d+)$/.test(value) && /[1-9]/.test(value);
  return positiveExact(modulus) && (!rate || positiveExact(rate)) && def.semantic_kind === 'counter' &&
    !def.numeric_encoding?.trim().toLowerCase().startsWith('float');
}
