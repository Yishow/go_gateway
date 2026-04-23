import { getDefaultNamingPrefixForProtocol } from '../../../features/datalink/sourcePlannerContract';
import type { DataType, ProtocolType } from '../../../types/datalink';
import type { AddressCanvasItem } from './sourceCanvasModel';

export type SourceRuleCoverageSummary = {
  planned: number;
  used: number;
  conflict: number;
};

export type SourcePlanningDatabaseAdvisory = {
  namingPrefix: string;
  targetDataType: DataType;
  scaleMultiplier: string;
  scaleOffset: string;
};

export function buildSourceRuleCoverageSummary(
  items: ReadonlyArray<AddressCanvasItem>,
  ruleId: string | null,
): SourceRuleCoverageSummary {
  if (!ruleId) {
    return { planned: 0, used: 0, conflict: 0 };
  }

  return items.reduce<SourceRuleCoverageSummary>(
    (summary, item) => {
      if (item.mergeOffset > 0 || !item.ruleIds.includes(ruleId)) {
        return summary;
      }

      if (item.status === 'planned') {
        summary.planned += 1;
      } else if (item.status === 'used') {
        summary.used += 1;
      } else if (item.status === 'conflict') {
        summary.conflict += 1;
      }
      return summary;
    },
    { planned: 0, used: 0, conflict: 0 },
  );
}

export function getDatabaseAdvisoryTargetDataType(dataType: DataType): DataType {
  switch (dataType) {
    case 'bool':
    case 'string':
      return dataType;
    default:
      return 'float64';
  }
}

export function buildSourcePlanningDatabaseAdvisory(
  protocol: ProtocolType,
  dataType: DataType,
): SourcePlanningDatabaseAdvisory {
  return {
    namingPrefix: `${getDefaultNamingPrefixForProtocol(protocol)}_ROW`,
    targetDataType: getDatabaseAdvisoryTargetDataType(dataType),
    scaleMultiplier: '1',
    scaleOffset: '0',
  };
}
