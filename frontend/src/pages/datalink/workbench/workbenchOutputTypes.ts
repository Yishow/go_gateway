import type { DataType } from '../../../types/datalink';

export interface WorkbenchOutputCandidate {
  tagId: string;
  tagKey: string;
  pointName: string;
  pointAddress: string;
  dataType: DataType;
  lastValue: unknown;
}

export type AutoMapStrategy = 'sequential' | 'gapAware' | 'aligned';

export type DryRunResult = {
  tagId: string;
  tagKey: string;
  register: number;
  valid: boolean;
  reason?: string;
};

export type OutputReadiness = 'ready' | 'partial' | 'unmapped';

export function computeOutputReadiness(
  hasModbus: boolean,
  hasDatabase: boolean,
): OutputReadiness {
  if (hasModbus && hasDatabase) {
    return 'ready';
  }
  if (hasModbus || hasDatabase) {
    return 'partial';
  }
  return 'unmapped';
}
