import type { DataType } from '../types/datalink';
import type { CounterPolicy } from '../types/measurement';
import type { WriteGroupExactRecordingSource } from '../types/studioV2WriteGroup';
import { isRecord, readRequiredText } from './studioV2WriteGroupJsonPrimitives';

const EXACT_SOURCE_REQUIRED_KEYS = new Set([
  'measurement_id', 'definition_revision', 'source_binding_revision', 'series_epoch',
  'original_data_type', 'scale', 'unit', 'semantic_kind',
]);
const EXACT_SOURCE_KEYS = new Set([...EXACT_SOURCE_REQUIRED_KEYS, 'counter_policy']);
const EXACT_COUNTER_DATA_TYPES = new Set<DataType>([
  'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64',
]);
const EXACT_COUNTER_SEMANTICS = new Set(['counter', 'signed_counter']);
const COUNTER_POLICY_KEYS = new Set([
  'reset_threshold', 'max_roll_over', 'allow_negative', 'rollover_modulus', 'max_rate_per_second',
]);
const PLAIN_DECIMAL_TEXT = /^\+?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)$/;

function isPositiveDecimalText(value: string): boolean {
  return value.length <= 64 && PLAIN_DECIMAL_TEXT.test(value) && /[1-9]/.test(value);
}

function hasOwn(record: Record<string, unknown>, key: string): boolean {
  return Object.prototype.hasOwnProperty.call(record, key);
}

function optionalPolicyNumber(record: Record<string, unknown>, key: string): number | undefined | null {
  if (!hasOwn(record, key)) return undefined;
  const value = record[key];
  return typeof value === 'number' && Number.isFinite(value) ? value : null;
}

function optionalPolicyDecimal(record: Record<string, unknown>, key: string): string | undefined | null {
  if (!hasOwn(record, key)) return undefined;
  const value = record[key];
  return typeof value === 'string' && isPositiveDecimalText(value) ? value : null;
}

function parseExactCounterPolicy(value: unknown): CounterPolicy | null {
  if (!isRecord(value) || Object.keys(value).some((key) => !COUNTER_POLICY_KEYS.has(key)) ||
    typeof value.allow_negative !== 'boolean') return null;

  const resetThreshold = optionalPolicyNumber(value, 'reset_threshold');
  const maxRollOver = optionalPolicyNumber(value, 'max_roll_over');
  const rolloverModulus = optionalPolicyDecimal(value, 'rollover_modulus');
  const maxRatePerSecond = optionalPolicyDecimal(value, 'max_rate_per_second');
  if (resetThreshold === null || maxRollOver === null || rolloverModulus === null || maxRatePerSecond === null ||
    (maxRatePerSecond !== undefined && rolloverModulus === undefined)) return null;

  return {
    ...(resetThreshold !== undefined ? { reset_threshold: resetThreshold } : {}),
    ...(maxRollOver !== undefined ? { max_roll_over: maxRollOver } : {}),
    allow_negative: value.allow_negative,
    ...(rolloverModulus !== undefined ? { rollover_modulus: rolloverModulus } : {}),
    ...(maxRatePerSecond !== undefined ? { max_rate_per_second: maxRatePerSecond } : {}),
  };
}

/** Parses the complete backend-issued proof; partial, changed, or unknown fields fail closed. */
export function parseWriteGroupExactRecordingSource(value: unknown): WriteGroupExactRecordingSource | null {
  if (!isRecord(value) || Object.keys(value).some((key) => !EXACT_SOURCE_KEYS.has(key)) ||
    [...EXACT_SOURCE_REQUIRED_KEYS].some((key) => !hasOwn(value, key))) return null;

  const measurementId = readRequiredText(value, 'measurement_id');
  const definitionRevision = readRequiredText(value, 'definition_revision');
  const sourceBindingRevision = readRequiredText(value, 'source_binding_revision');
  const seriesEpoch = readRequiredText(value, 'series_epoch');
  const originalDataType = readRequiredText(value, 'original_data_type');
  const scale = value.scale;
  const unit = readRequiredText(value, 'unit');
  const semanticKind = readRequiredText(value, 'semantic_kind');
  const policyPresent = hasOwn(value, 'counter_policy');
  const counterPolicy = policyPresent ? parseExactCounterPolicy(value.counter_policy) : undefined;
  if (!measurementId || !definitionRevision || !sourceBindingRevision || !seriesEpoch ||
    !originalDataType || !EXACT_COUNTER_DATA_TYPES.has(originalDataType as DataType) ||
    typeof scale !== 'string' || !isPositiveDecimalText(scale) || !unit ||
    !semanticKind || !EXACT_COUNTER_SEMANTICS.has(semanticKind) ||
    (policyPresent && !counterPolicy) ||
    (semanticKind === 'counter' && counterPolicy?.allow_negative === true) ||
    (semanticKind === 'signed_counter' && (counterPolicy?.allow_negative !== true ||
      counterPolicy.rollover_modulus !== undefined || counterPolicy.max_rate_per_second !== undefined))) {
    return null;
  }

  return {
    measurement_id: measurementId,
    definition_revision: definitionRevision,
    source_binding_revision: sourceBindingRevision,
    series_epoch: seriesEpoch,
    original_data_type: originalDataType as DataType,
    scale,
    unit,
    semantic_kind: semanticKind as WriteGroupExactRecordingSource['semantic_kind'],
    ...(counterPolicy ? { counter_policy: counterPolicy } : {}),
  };
}

/** Clones the backend proof without changing an omitted optional policy into `{}`. */
export function cloneWriteGroupExactRecordingSource(
  source: WriteGroupExactRecordingSource,
): WriteGroupExactRecordingSource {
  const { counter_policy: counterPolicy, ...proof } = source;
  return {
    ...proof,
    ...(counterPolicy !== undefined ? { counter_policy: { ...counterPolicy } } : {}),
  };
}
