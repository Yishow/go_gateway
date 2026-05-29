import * as React from 'react';
import type { Rule, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { Field } from '../../components/Field';
import { Input, Select } from '../../components/inputs';
import { Toggle } from '../../components/Toggle';
import { RangeSummary } from './RangeSummary';
import { ScaleSection } from './ScaleSection';
import { ShareSection } from './ShareSection';

export interface RuleEditorProps {
  rule: Rule;
  devices: Device[];
  globalShareEnabled: boolean;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 接入規則編輯器面板元件
 * 
 * 落地設計決策：「Rule editor with linked reset」
 * 提供設備選擇、名稱、前綴、起始暫存器與數量的 2x2 緊湊 Grid 編輯區。
 * 底下整合 RangeSummary、ScaleSection、ShareSection 展開區，以及控制該規則啟用的 Toggle。
 */
export const RuleEditor: React.FC<RuleEditorProps> = ({
  rule,
  devices,
  globalShareEnabled,
  dispatch,
}) => {
  const handleTextChange = (field: keyof Rule, val: string) => {
    dispatch({
      type: 'updateRule',
      ruleId: rule.id,
      patch: { [field]: val },
    });
  };

  const handleNumberChange = (field: keyof Rule, valStr: string, fallback = 1) => {
    const val = parseInt(valStr, 10);
    dispatch({
      type: 'updateRule',
      ruleId: rule.id,
      patch: { [field]: isNaN(val) ? fallback : val },
    });
  };

  return (
    <div
      className="flex flex-col gap-5 p-5 rounded-xl border border-slate-800 bg-slate-900/20"
      data-testid="rule-editor"
    >
      <h3 className="text-sm font-semibold text-slate-200">接入規則參數配置</h3>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* 裝置選擇 (僅在有多個裝置時啟用) */}
        <Field label="所屬設備來源">
          <Select
            value={rule.device_id}
            disabled={devices.length <= 1}
            onChange={(e) => handleTextChange('device_id', e.target.value)}
            data-testid="rule-device-select"
          >
            {devices.map((d) => (
              <option key={d.id} value={d.id}>
                {d.name}
              </option>
            ))}
          </Select>
        </Field>

        {/* 規則名稱 */}
        <Field label="規則名稱">
          <Input
            type="text"
            value={rule.name}
            onChange={(e) => handleTextChange('name', e.target.value)}
            data-testid="rule-name-input"
          />
        </Field>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {/* 起始位址 */}
        <Field label="起始暫存器位址" hint="例: 40001">
          <Input
            type="text"
            value={rule.start_address}
            onChange={(e) => handleTextChange('start_address', e.target.value)}
            data-testid="rule-start-input"
          />
        </Field>

        {/* 點位數量 */}
        <Field label="點位數量" hint="範圍 1 至 64">
          <Input
            type="number"
            min={1}
            max={64}
            value={rule.count}
            onChange={(e) => handleNumberChange('count', e.target.value, 1)}
            data-testid="rule-count-input"
          />
        </Field>

        {/* 資料型別 */}
        <Field label="資料型別">
          <Select
            value={rule.data_type}
            onChange={(e) => handleTextChange('data_type', e.target.value)}
            data-testid="rule-datatype-select"
          >
            <option value="bool">bool</option>
            <option value="int16">int16</option>
            <option value="uint16">uint16</option>
            <option value="int32">int32</option>
            <option value="uint32">uint32</option>
            <option value="int64">int64</option>
            <option value="uint64">uint64</option>
            <option value="float32">float32</option>
            <option value="float64">float64</option>
            <option value="string">string</option>
          </Select>
        </Field>

        {/* 命名前綴 */}
        <Field label="點位名稱前綴" hint="例: TEMP_">
          <Input
            type="text"
            value={rule.naming_prefix}
            onChange={(e) => handleTextChange('naming_prefix', e.target.value)}
            data-testid="rule-prefix-input"
          />
        </Field>
      </div>

      {/* 暫存器區段摘要 */}
      <RangeSummary
        startAddress={rule.start_address}
        count={rule.count}
        dataType={rule.data_type}
      />

      {/* 線性縮放 Details */}
      <ScaleSection
        ruleId={rule.id}
        multiplier={rule.scale_multiplier}
        offset={rule.scale_offset}
        dataFormat={rule.data_format}
        dispatch={dispatch}
      />

      {/* Modbus Share Details */}
      <ShareSection
        ruleId={rule.id}
        shareEnabled={rule.share_enabled}
        shareStartRegister={rule.share_start_register}
        shareStride={rule.share_stride}
        globalShareEnabled={globalShareEnabled}
        dispatch={dispatch}
      />

      {/* 啟用此規則 */}
      <div className="border-t border-slate-800/60 pt-4 flex items-center justify-between">
        <Toggle
          checked={rule.enabled}
          onChange={() => dispatch({ type: 'toggleRuleEnabled', ruleId: rule.id })}
          label="啟用此規則 (停用時將不衍生其所屬點位)"
          data-testid="rule-editor-enable-toggle"
        />
      </div>
    </div>
  );
};
