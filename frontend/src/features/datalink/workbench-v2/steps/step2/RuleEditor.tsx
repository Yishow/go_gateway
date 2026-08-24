import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Rule, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { Field } from '../../components/Field';
import { Input, Select } from '../../components/inputs';
import { Toggle } from '../../components/Toggle';
import { RangeSummary } from './RangeSummary';
import { ScaleSection } from './ScaleSection';
import { ShareSection } from './ShareSection';
import { getRuleReadinessReason } from '../../state/sourceRule';

export interface RuleEditorProps {
  rule: Rule;
  devices: Device[];
  globalShareEnabled: boolean;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 接入規則編輯器面板元件
 * 
 * 落地設計決策：「局部 State vs Reducer：即時數值同步」
 * 提供設備選擇、名稱、起始地址、點位數量、資料類型、
 * 縮放比例/偏移、以及 Modbus Share 獨立覆蓋等配置，即時連動全域 Reducer。
 */
export const RuleEditor: React.FC<RuleEditorProps> = ({
  rule,
  devices,
  globalShareEnabled,
  dispatch,
}) => {
  const { t } = useTranslation('workbench-v2');

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

  const owningDevice = devices.find((d) => d.id === rule.device_id);
  const readinessReason = getRuleReadinessReason(
    rule,
    owningDevice ? { [owningDevice.id]: owningDevice.protocol } : {},
  );
  const hasMissingDevice = readinessReason === 'unknown_device' || readinessReason === 'deleted_device';
  const missingDeviceKey = readinessReason === 'deleted_device' ? 'deleted_device' : 'unknown_device';
  const hasInvalidAddress = readinessReason === 'invalid_address';
  const hasValidAddress = readinessReason === null;
  const missingDeviceFallback = readinessReason === 'deleted_device'
    ? {
      title: '此規則所屬設備已被刪除。',
      action: '請選擇替代設備或先還原設備後再保存。',
      option: '設備已刪除，請選擇替代設備',
    }
    : {
      title: '此規則沒有目前工作區的設備歸屬。',
      action: '請先選擇目前設備，再保存此規則。',
      option: '沒有目前設備，請選擇設備',
    };
  const protocol = owningDevice?.protocol;
  const isModbus = protocol?.startsWith('modbus') ?? false;
  const addressErrorId = `rule-start-error-${rule.id}`;
  const addrHint = isModbus
    ? t('step2.editor.addr_hint_modbus', '例: 40001')
    : t('step2.editor.addr_hint_plc', '例: D0, M0, W0');
  const addrPlaceholder = owningDevice ? (isModbus ? '40001' : 'D0') : '';

  return (
    <div
      className="flex flex-col gap-5 p-5 rounded-xl border border-slate-800 bg-slate-900/20"
      data-testid="rule-editor"
    >
      <h3 className="text-sm font-semibold text-slate-200">接入規則參數配置</h3>

      {hasMissingDevice && (
        <div
          className="rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-200"
          role="alert"
          data-testid="rule-missing-device"
          data-readiness-reason={readinessReason}
        >
          <div className="font-semibold">
            {t(`step2.editor.${missingDeviceKey}`, missingDeviceFallback.title)}
          </div>
          <div className="mt-1">
            {t(`step2.editor.${missingDeviceKey}_action`, missingDeviceFallback.action)}
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* 裝置選擇 (僅在有多個裝置時啟用) */}
        <Field label="所屬設備來源">
          <Select
            value={rule.device_id}
            disabled={devices.length <= 1 && !hasMissingDevice}
            aria-invalid={hasMissingDevice}
            onChange={(e) => handleTextChange('device_id', e.target.value)}
            data-testid="rule-device-select"
          >
            {hasMissingDevice && (
              <option value={rule.device_id} disabled>
                {t(`step2.editor.${missingDeviceKey}_option`, missingDeviceFallback.option)}
              </option>
            )}
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
        <Field label="起始暫存器位址" hint={hasValidAddress ? addrHint : undefined}>
          <Input
            type="text"
            value={rule.start_address}
            placeholder={addrPlaceholder}
            disabled={hasMissingDevice}
            onChange={(e) => handleTextChange('start_address', e.target.value)}
            aria-invalid={!hasValidAddress}
            aria-describedby={!hasValidAddress ? addressErrorId : undefined}
            data-testid="rule-start-input"
          />
          {!hasValidAddress && (
            <span
              id={addressErrorId}
              data-testid="rule-start-error"
              className="mt-1 flex items-center gap-1 text-xs text-red-300"
            >
              {hasMissingDevice
                ? t('step2.editor.missing_device_address', '請先重新選擇或修復設備，再檢查位址。')
                : hasInvalidAddress
                ? t('step2.editor.addr_invalid', '位址不符合所選協議格式，請輸入有效位址。')
                : null}
            </span>
          )}
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
      {owningDevice ? (
        <RangeSummary
          startAddress={rule.start_address}
          count={rule.count}
          dataType={rule.data_type}
          protocol={owningDevice.protocol}
        />
      ) : (
        <div
          className="rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-200"
          data-testid="rule-missing-device-summary"
        >
          {t('step2.editor.missing_device_summary', '無法計算位址範圍，請先修復設備來源。')}
        </div>
      )}

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
