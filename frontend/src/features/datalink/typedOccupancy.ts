import type { DataType } from '../../types/datalink';

export type OccupancyStatus =
  | 'available'
  | 'planned'
  | 'used'
  | 'linked'
  | 'selected'
  | 'conflict';

export type ConflictSeverity = 'none' | 'soft' | 'hard';

export const SPAN_BY_DATA_TYPE: Record<DataType, number> = {
  bool: 1,
  int16: 1,
  uint16: 1,
  int32: 2,
  uint32: 2,
  float32: 2,
  int64: 4,
  uint64: 4,
  float64: 4,
  string: 1,
};

export interface TypedOccupancyValidationResult {
  valid: boolean;
  span: number;
  totalCells: number;
  errors: Array<'invalid_count' | 'invalid_span'>;
}

export function getSpanByDataType(dataType: DataType): number {
  return SPAN_BY_DATA_TYPE[dataType];
}

export function validateTypedOccupancyPlan(dataType: DataType, count: number): TypedOccupancyValidationResult {
  const span = getSpanByDataType(dataType);
  const errors: Array<'invalid_count' | 'invalid_span'> = [];

  const normalizedCount = Number.isFinite(count) ? Math.floor(count) : 0;
  if (normalizedCount <= 0 || normalizedCount > 200) errors.push('invalid_count');
  if (![1, 2, 4].includes(span)) errors.push('invalid_span');

  return {
    valid: errors.length === 0,
    span,
    totalCells: Math.max(0, normalizedCount) * span,
    errors,
  };
}

export function resolveConflictSeverity(input: {
  hasUsedPoint: boolean;
  hasLinkedAddress: boolean;
}): ConflictSeverity {
  if (input.hasUsedPoint) return 'hard';
  if (input.hasLinkedAddress) return 'soft';
  return 'none';
}
