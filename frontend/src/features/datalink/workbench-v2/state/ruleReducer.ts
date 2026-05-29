import type { WorkbenchV2State } from './types';
import type { WorkbenchV2Action } from './useWorkbenchV2State';

/**
 * 處理接入規則相關狀態轉移的 Sub-Reducer
 * 
 * 落地設計決策：「Reducer 拆檔與模組化」
 * 將 11 個規則管理 action 移至獨立檔案，規避 useWorkbenchV2State.ts 行數超限問題，
 * 同時支援自動/手動 share reset、count 限制、skipped addresses 重設等細節邏輯。
 */
export function ruleReducer(state: WorkbenchV2State, action: WorkbenchV2Action): WorkbenchV2State {
  switch (action.type) {
    case 'selectRule':
      return {
        ...state,
        selectedRuleId: action.ruleId,
      };
    case 'addRule':
      return {
        ...state,
        rules: [...state.rules, action.rule],
        selectedRuleId: action.rule.id,
      };
    case 'removeRule': {
      const nextRules = state.rules.filter((r) => r.id !== action.ruleId);
      let nextSelectedRuleId = state.selectedRuleId;
      if (nextSelectedRuleId === action.ruleId) {
        nextSelectedRuleId = nextRules.length > 0 ? nextRules[0].id : null;
      }
      return {
        ...state,
        rules: nextRules,
        selectedRuleId: nextSelectedRuleId,
      };
    }
    case 'updateRule':
      return {
        ...state,
        rules: state.rules.map((r) => {
          if (r.id !== action.ruleId) return r;
          const patch = { ...action.patch };
          if (patch.count !== undefined) {
            patch.count = Math.max(1, Math.min(64, patch.count));
          }
          const hasLinkedResetKeys =
            patch.start_address !== undefined ||
            patch.count !== undefined ||
            patch.data_type !== undefined;

          return {
            ...r,
            ...patch,
            skipped_addresses: hasLinkedResetKeys ? [] : r.skipped_addresses,
          };
        }),
      };
    case 'renameRule':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, name: action.name } : r
        ),
      };
    case 'toggleRuleEnabled':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, enabled: !r.enabled } : r
        ),
      };
    case 'updateRuleSkipped':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, skipped_addresses: action.skippedAddresses } : r
        ),
      };
    case 'toggleRuleSkippedAddress':
      return {
        ...state,
        rules: state.rules.map((r) => {
          if (r.id !== action.ruleId) return r;
          const skipped = r.skipped_addresses || [];
          const exists = skipped.includes(action.address);
          const nextSkipped = exists
            ? skipped.filter((a) => a !== action.address)
            : [...skipped, action.address];
          return {
            ...r,
            skipped_addresses: nextSkipped,
          };
        }),
      };
    case 'toggleRuleShareEnabled':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, share_enabled: !r.share_enabled } : r
        ),
      };
    case 'updateRuleShareStart':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, share_start_register: action.shareStart } : r
        ),
      };
    case 'updateRuleShareStride':
      return {
        ...state,
        rules: state.rules.map((r) =>
          r.id === action.ruleId ? { ...r, share_stride: action.shareStride } : r
        ),
      };
    default:
      return state;
  }
}
