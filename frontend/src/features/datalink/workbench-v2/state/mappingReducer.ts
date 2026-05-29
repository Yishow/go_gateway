import type { WorkbenchV2State, Mapping } from './types';
import { buildDefaultMapping } from './mappingDefaults';

/**
 * 處理點位映射相關狀態轉移的 Sub-Reducer
 * 
 * 落地設計決策：「Reducer 拆檔與模組化」
 * 封裝 4 個點位映射管理 Action，防止主 Reducer 行數爆滿。
 */
export function mappingReducer(state: WorkbenchV2State, action: any): WorkbenchV2State {
  switch (action.type) {
    case 'initMappingsForPoints': {
      const nextMappings: Record<string, Mapping> = {};
      action.points.forEach((p: any, idx: number) => {
        if (state.mappings[p.id]) {
          nextMappings[p.id] = state.mappings[p.id];
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
      return {
        ...state,
        mappings: {
          ...state.mappings,
          [action.pointId]: {
            ...state.mappings[action.pointId],
            ...action.patch,
          },
        },
      };
    }
    case 'toggleMappingEnabled': {
      if (!state.mappings[action.pointId]) return state;
      return {
        ...state,
        mappings: {
          ...state.mappings,
          [action.pointId]: {
            ...state.mappings[action.pointId],
            enabled: !state.mappings[action.pointId].enabled,
          },
        },
      };
    }
    case 'bulkApplyTransform': {
      const source = state.mappings[action.fromPointId];
      if (!source) return state;

      const nextMappings = { ...state.mappings };
      Object.keys(nextMappings).forEach((pointId) => {
        const m = nextMappings[pointId];
        const updated = { ...m };
        action.fields.forEach((field: 'scale' | 'offset' | 'target_type') => {
          if (field === 'scale') updated.scale = source.scale;
          if (field === 'offset') updated.offset = source.offset;
          if (field === 'target_type') updated.target_type = source.target_type;
        });
        nextMappings[pointId] = updated;
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
