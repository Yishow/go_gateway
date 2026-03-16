import type { DataType } from '../../../types/datalink';

export interface WorkbenchOutputCandidate {
  tagId: string;
  tagKey: string;
  pointName: string;
  pointAddress: string;
  dataType: DataType;
  lastValue: unknown;
}
