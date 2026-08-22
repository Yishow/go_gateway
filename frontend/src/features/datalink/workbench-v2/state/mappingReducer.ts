import type { WorkbenchV2State, Mapping, MappingValue } from './types';
import { buildDefaultMapping } from './mappingDefaults';
import type { Point } from './types';
import type { WorkbenchV2Action } from './useWorkbenchV2State';

/**
 * 處理點位映射相關狀態轉移的 Sub-Reducer
 * 
 * 落地設計決策：「Reducer 拆檔與模組化」
 * 封裝 4 個點位映射管理 Action，防止主 Reducer 行數爆滿。
 */
function editableValueOf(mapping: Mapping): MappingValue {
  return {
    tag_key: mapping.tag_key,
    display_name: mapping.display_name,
    unit: mapping.unit,
    target_type: mapping.target_type,
    scale: mapping.scale,
    offset: mapping.offset,
    enabled: mapping.enabled,
  };
}

function withLocalValue(mapping: Mapping): Mapping {
  return {
    ...mapping,
    local_value: editableValueOf(mapping),
  };
}

function isSamePointRow(mapping: Mapping, point: Point): boolean {
  if (mapping.rule_id && mapping.rule_id !== point.rule_id) {
    return false;
  }
  if (mapping.device_id && mapping.device_id !== point.device_id) {
    return false;
  }
  if (mapping.address && mapping.address.trim().toUpperCase() !== point.address.trim().toUpperCase()) {
    return false;
  }
  return true;
}

export function mappingReducer(state: WorkbenchV2State, action: WorkbenchV2Action): WorkbenchV2State {
  switch (action.type) {
    case 'initMappingsForPoints': {
      const nextMappings: Record<string, Mapping> = {};
      action.points.forEach((p: Point, idx: number) => {
        const current = state.mappings[p.id];
        if (current && isSamePointRow(current, p)) {
          nextMappings[p.id] = current;
        } else {
          nextMappings[p.id] = buildDefaultMapping(p, idx);
        }
      });
      return {
        ...state,
        points: action.points,
        mappings: nextMappings,
      };
    }
    case 'updateMapping': {
      if (!state.mappings[action.pointId]) return state;
      const nextMapping = {
        ...state.mappings[action.pointId],
        ...action.patch,
      } satisfies Mapping;
      return {
        ...state,
        mappings: {
          ...state.mappings,
          [action.pointId]: withLocalValue(nextMapping),
        },
      };
    }
    case 'toggleMappingEnabled': {
      if (!state.mappings[action.pointId]) return state;
      const nextMapping = {
        ...state.mappings[action.pointId],
        enabled: !state.mappings[action.pointId].enabled,
      } satisfies Mapping;
      return {
        ...state,
        mappings: {
          ...state.mappings,
          [action.pointId]: withLocalValue(nextMapping),
        },
      };
    }
    case 'setAllMappingsEnabled': {
      const nextMappings = Object.fromEntries(
        Object.entries(state.mappings).map(([pointId, mapping]) => [
          pointId,
          withLocalValue({
            ...mapping,
            enabled: action.enabled,
          }),
        ]),
      );

      return {
        ...state,
        mappings: nextMappings,
      };
    }
    case 'bulkApplyTransform': {
      const source = state.mappings[action.fromPointId];
      if (!source) return state;

      const nextMappings = { ...state.mappings };
      Object.keys(nextMappings).forEach((pointId) => {
        const m = nextMappings[pointId];
        if (!m.enabled) {
          return;
        }
        const updated = { ...m };
        action.fields.forEach((field: 'scale' | 'offset' | 'target_type' | 'unit') => {
          if (field === 'scale') updated.scale = source.scale;
          if (field === 'offset') updated.offset = source.offset;
          if (field === 'target_type') updated.target_type = source.target_type;
          if (field === 'unit') updated.unit = source.unit;
        });
        nextMappings[pointId] = withLocalValue(updated);
      });

      return {
        ...state,
        mappings: nextMappings,
      };
    }
    default:
      return state;
  }
}
