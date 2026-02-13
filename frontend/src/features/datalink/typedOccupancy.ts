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

export function getSpanByDataType(dataType: DataType): number {
  return SPAN_BY_DATA_TYPE[dataType];
}

export function resolveConflictSeverity(input: {
  hasUsedPoint: boolean;
  hasLinkedAddress: boolean;
}): ConflictSeverity {
  if (input.hasUsedPoint) return 'hard';
  if (input.hasLinkedAddress) return 'soft';
  return 'none';
}
