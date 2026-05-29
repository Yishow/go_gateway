import type { WorkbenchV2State, SettingsConnector } from './types';
import type { WorkbenchV2Action } from './useWorkbenchV2State';
import { DEFAULT_SETTINGS, makeDefaultConnector } from './settingsDefaults';

/**
 * 處理系統設定與連接器池狀態變更的子 Reducer
 * 落地設計決策：「Connector 池：以陣列管理 + per-row status」
 * 
 * @param state 當前工作台狀態
 * @param action 觸發的動作
 * @returns 變更後的狀態
 */
export function settingsReducer(state: WorkbenchV2State, action: WorkbenchV2Action): WorkbenchV2State {
  switch (action.type) {
    case 'updateSettings':
      return {
        ...state,
        settings: {
          ...state.settings,
          ...action.patch
        }
      };

    case 'updateSettingsSection':
      return {
        ...state,
        settings: {
          ...state.settings,
          [action.section]: {
            ...state.settings[action.section],
            ...action.patch
          }
        }
      };

    case 'addConnector': {
      const idx = state.settings.connectors.length + 1;
      const newConn = makeDefaultConnector(idx);
      return {
        ...state,
        settings: {
          ...state.settings,
          connectors: [...state.settings.connectors, newConn]
        }
      };
    }

    case 'updateConnector':
      return {
        ...state,
        settings: {
          ...state.settings,
          connectors: state.settings.connectors.map((c: SettingsConnector) => {
            if (c.id !== action.id) return c;
            
            // 檢查是否修改了關鍵連線參數
            const hasKindChange = action.patch.kind !== undefined && action.patch.kind !== c.kind;
            const hasHostChange = action.patch.host !== undefined && action.patch.host !== c.host;
            const hasPortChange = action.patch.port !== undefined && action.patch.port !== c.port;

            const nextConnector = { ...c, ...action.patch };
            
            // 若修改關鍵參數，重置狀態為 unknown
            if (hasKindChange || hasHostChange || hasPortChange) {
              nextConnector.status = 'unknown';
            }
            return nextConnector;
          })
        }
      };

    case 'removeConnector':
      return {
        ...state,
        settings: {
          ...state.settings,
          connectors: state.settings.connectors.filter((c: SettingsConnector) => c.id !== action.id)
        }
      };

    case 'startConnectorTest':
      return {
        ...state,
        settings: {
          ...state.settings,
          connectors: state.settings.connectors.map((c: SettingsConnector) => 
            c.id === action.id ? { ...c, status: 'testing' } : c
          )
        }
      };

    case 'completeConnectorTest':
      return {
        ...state,
        settings: {
          ...state.settings,
          connectors: state.settings.connectors.map((c: SettingsConnector) => 
            c.id === action.id ? { ...c, ...action.result } : c
          )
        }
      };

    case 'resetSettingsToDefaults':
      return {
        ...state,
        settings: {
          ...DEFAULT_SETTINGS
        }
      };

    default:
      return state;
  }
}
