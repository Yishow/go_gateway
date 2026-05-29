import * as React from 'react';
import { useState, useRef, useEffect } from 'react';
import type { Rule, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { useDeviceColor } from '../../state/deviceColors';
import { Icon } from '../../components';

export interface RuleTabRailProps {
  rules: Rule[];
  devices: Device[];
  selectedRuleId: string | null;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

interface RuleTabItemProps {
  rule: Rule;
  isSelected: boolean;
  editingRuleId: string | null;
  editingName: string;
  setEditingName: (name: string) => void;
  setEditingRuleId: (id: string | null) => void;
  devices: Device[];
  rules: Rule[];
  dispatch: React.Dispatch<WorkbenchV2Action>;
  handleRenameSubmit: (ruleId: string) => void;
  inputRef: React.RefObject<HTMLInputElement | null>;
}

/**
 * 個別規則 Tab 項目元件
 * 
 * 落地設計決策：規避迴圈呼叫 Hook，封裝單一 Tab 渲染與 useDeviceColor 呼叫。
 */
const RuleTabItem: React.FC<RuleTabItemProps> = ({
  rule,
  isSelected,
  editingRuleId,
  editingName,
  setEditingName,
  setEditingRuleId,
  devices,
  rules,
  dispatch,
  handleRenameSubmit,
  inputRef,
}) => {
  const colorTheme = useDeviceColor(rule.device_id);

  return (
    <div
      onClick={() => {
        if (editingRuleId !== rule.id) {
          dispatch({ type: 'selectRule', ruleId: rule.id });
        }
      }}
      className={`group relative flex flex-col justify-between p-3.5 rounded-xl border transition-all duration-200 cursor-pointer select-none min-w-[180px] h-[104px] ${
        isSelected
          ? 'bg-slate-800/80 border-slate-600 shadow-md ring-1 ring-slate-700/50'
          : 'bg-slate-900/40 border-slate-800/80 hover:bg-slate-800/40 hover:border-slate-700/60'
      } ${!rule.enabled ? 'opacity-60' : ''}`}
      data-testid={`rule-tab-${rule.id}`}
    >
      {/* Top Row: Name and Action Icons */}
      <div className="flex items-center justify-between gap-2">
        {editingRuleId === rule.id ? (
          <input
            ref={inputRef}
            type="text"
            value={editingName}
            onChange={(e) => setEditingName(e.target.value)}
            onBlur={() => handleRenameSubmit(rule.id)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleRenameSubmit(rule.id);
              if (e.key === 'Escape') setEditingRuleId(null);
            }}
            className="bg-slate-950 text-slate-100 text-xs px-1.5 py-0.5 rounded border border-slate-700 focus:outline-none focus:border-blue-500 w-full"
            onClick={(e) => e.stopPropagation()}
          />
        ) : (
          <span
            onDoubleClick={(e) => {
              e.stopPropagation();
              setEditingRuleId(rule.id);
              setEditingName(rule.name);
            }}
            className="font-medium text-xs text-slate-200 truncate pr-6"
            title="雙擊可改名"
          >
            {rule.name}
          </span>
        )}

        {/* Hover Actions: Toggle enabled & Close */}
        <div
          className="absolute top-2.5 right-2.5 flex items-center gap-1.5 opacity-0 group-hover:opacity-100 transition-opacity duration-150"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Enable / Disable toggle button */}
          <button
            onClick={() => dispatch({ type: 'toggleRuleEnabled', ruleId: rule.id })}
            className="p-1 rounded hover:bg-slate-700 text-slate-400 hover:text-slate-200"
            title={rule.enabled ? '停用規則' : '啟用規則'}
            data-testid={`rule-toggle-enabled-${rule.id}`}
          >
            <Icon
              name={rule.enabled ? 'eye' : 'alert'}
              className={`w-3.5 h-3.5 ${rule.enabled ? 'text-blue-400' : 'text-slate-500'}`}
            />
          </button>

          {/* Close/Remove rule button (only allow deletion if there is more than 1 rule) */}
          {rules.length > 1 && (
            <button
              onClick={() => dispatch({ type: 'removeRule', ruleId: rule.id })}
              className="p-1 rounded hover:bg-red-500/20 text-slate-400 hover:text-red-400"
              title="刪除規則"
              data-testid={`rule-delete-${rule.id}`}
            >
              <Icon name="close" className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      </div>

      {/* Middle Row: Register Info */}
      <div className="text-[10px] text-slate-400 font-mono mt-1">
        {rule.start_address} ({rule.data_type}) · {rule.count} Pts
      </div>

      {/* Bottom Row: Device Badge (shown if 2+ devices exist) */}
      {devices.length >= 2 ? (
        <div className="flex items-center gap-1.5 mt-2 text-[10px] text-slate-500">
          <span className={`w-1.5 h-1.5 rounded-full ${colorTheme.solid}`} />
          <span className="truncate max-w-[130px]">
            {devices.find((d) => d.id === rule.device_id)?.name || '未知裝置'}
          </span>
        </div>
      ) : (
        <div className="h-4" />
      )}
    </div>
  );
};

/**
 * 規則分頁切換列元件
 * 
 * 落地設計決策：「Multi-rule tab management」與「簡潔改名互動」
 */
export const RuleTabRail: React.FC<RuleTabRailProps> = ({
  rules,
  devices,
  selectedRuleId,
  dispatch,
}) => {
  const [editingRuleId, setEditingRuleId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (editingRuleId && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [editingRuleId]);

  const handleAddRule = () => {
    const defaultDevId = devices[0]?.id || 'dev-01';
    const newId = `rule-${Date.now()}`;
    const newRule: Rule = {
      id: newId,
      device_id: defaultDevId,
      name: `規則 ${rules.length + 1}`,
      start_address: '40001',
      count: 8,
      data_type: 'int16',
      naming_prefix: 'TAG_',
      enabled: true,
      scale_multiplier: 1,
      scale_offset: 0,
      data_format: '',
      skipped_addresses: [],
      share_enabled: false,
      share_start_register: null,
      share_stride: null,
    };
    dispatch({ type: 'addRule', rule: newRule });
  };

  const handleRenameSubmit = (ruleId: string) => {
    const trimmed = editingName.trim();
    if (trimmed) {
      dispatch({ type: 'renameRule', ruleId, name: trimmed });
    }
    setEditingRuleId(null);
  };

  return (
    <div
      className="flex items-center gap-3 overflow-x-auto pb-2 scrollbar-thin"
      data-testid="rule-tab-rail"
    >
      {rules.map((rule) => (
        <RuleTabItem
          key={rule.id}
          rule={rule}
          isSelected={rule.id === selectedRuleId}
          editingRuleId={editingRuleId}
          editingName={editingName}
          setEditingName={setEditingName}
          setEditingRuleId={setEditingRuleId}
          devices={devices}
          rules={rules}
          dispatch={dispatch}
          handleRenameSubmit={handleRenameSubmit}
          inputRef={inputRef}
        />
      ))}

      {/* Add Rule Button */}
      <button
        onClick={handleAddRule}
        className="flex flex-col items-center justify-center gap-2 p-3.5 rounded-xl border border-dashed border-slate-700/60 hover:border-slate-500 hover:bg-slate-900/30 text-slate-400 hover:text-slate-200 transition-all duration-200 min-w-[120px] h-[104px]"
        data-testid="rule-add-btn"
      >
        <Icon name="plus" className="w-5 h-5" />
        <span className="text-[11px] font-medium">新增規則</span>
      </button>
    </div>
  );
};
