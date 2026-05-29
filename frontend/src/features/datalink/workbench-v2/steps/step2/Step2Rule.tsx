import * as React from 'react';
import { useState } from 'react';
import type { WorkbenchV2State } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import {
  useAllPoints,
  useShareLayout,
  useConflictAddrs,
  useRulePoints,
} from '../../state/selectors';
import { RuleTabRail } from './RuleTabRail';
import { RuleEditor } from './RuleEditor';
import { PointGridToolbar } from './PointGridToolbar';
import { PointGrid } from './PointGrid';
import { MergedPointTable } from './MergedPointTable';

export interface Step2RuleProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  onContinue: () => void;
}

/**
 * 接入規則與點位衍生容器元件 (Step 2)
 * 
 * 落地設計決策：「拆檔策略：容器與子元件分離」
 * 負責彙整 Selector 衍生狀態（所有點位、衝突地址、Share 佈局），
 * 並實作點位網格的 6 種批次操作邏輯與 state callbacks，
 * 採用 5-col / 7-col 橫向排版與底部合併總表組合出完整的 Step 2 畫面。
 */
export const Step2Rule: React.FC<Step2RuleProps> = ({
  state,
  dispatch,
  onContinue,
}) => {
  const { rules, devices, selectedRuleId, settings } = state;
  const [gridSelection, setGridSelection] = useState<Set<string>>(new Set());

  // 1. 取得備用設備 ID 並衍生全域點位與衝突
  const fallbackDeviceId = devices[0]?.id || 'dev-01';
  const allPoints = useAllPoints(rules, fallbackDeviceId);
  const conflictAddrs = useConflictAddrs(allPoints);

  // 2. 計算全域 Modbus Share 佈局
  const shareLayouts = useShareLayout(rules, settings.modbus_share.base_register);

  // 3. 取得當前選中規則及相關點位衍生狀態
  const currentRule = rules.find((r) => r.id === selectedRuleId) || rules[0];
  const skippedSet = new Set(currentRule?.skipped_addresses || []);
  const currentRulePoints = useRulePoints(
    currentRule,
    currentRule?.device_id || fallbackDeviceId,
    skippedSet
  );
  const currentShareLayout = shareLayouts[currentRule?.id] || null;

  // 4. 點位略過/啟用批次修改 Callback
  const handleBatchToggleSkip = (addresses: string[], shouldSkip: boolean) => {
    const currentSkipped = currentRule.skipped_addresses || [];
    let nextSkipped: string[];

    if (shouldSkip) {
      nextSkipped = Array.from(new Set([...currentSkipped, ...addresses]));
    } else {
      const removeSet = new Set(addresses);
      nextSkipped = currentSkipped.filter((a) => !removeSet.has(a));
    }

    dispatch({
      type: 'updateRuleSkipped',
      ruleId: currentRule.id,
      skippedAddresses: nextSkipped,
    });
  };

  // 5. 處理單一略過 Toggle
  const handleToggleSkip = (address: string) => {
    dispatch({
      type: 'toggleRuleSkippedAddress',
      ruleId: currentRule.id,
      address,
    });
  };

  // 6. 點位 Toolbar 批次按鈕事件處理
  const handleToolbarAction = (
    actionType:
      | 'all_enabled'
      | 'all_skipped'
      | 'reverse'
      | 'skip_selected'
      | 'enable_selected'
      | 'clear_selection'
  ) => {
    if (!currentRule) return;

    switch (actionType) {
      case 'all_enabled':
        dispatch({
          type: 'updateRuleSkipped',
          ruleId: currentRule.id,
          skippedAddresses: [],
        });
        break;
      case 'all_skipped': {
        const allAddrs = currentRulePoints.map((p) => p.address);
        dispatch({
          type: 'updateRuleSkipped',
          ruleId: currentRule.id,
          skippedAddresses: allAddrs,
        });
        break;
      }
      case 'reverse': {
        const skipped = currentRule.skipped_addresses || [];
        const skippedSet = new Set(skipped);
        const nextSkipped: string[] = [];

        currentRulePoints.forEach((p) => {
          if (!skippedSet.has(p.address)) {
            nextSkipped.push(p.address);
          }
        });

        dispatch({
          type: 'updateRuleSkipped',
          ruleId: currentRule.id,
          skippedAddresses: nextSkipped,
        });
        break;
      }
      case 'skip_selected': {
        const selectedAddrs = currentRulePoints
          .filter((p) => gridSelection.has(p.id))
          .map((p) => p.address);
        handleBatchToggleSkip(selectedAddrs, true);
        break;
      }
      case 'enable_selected': {
        const selectedAddrs = currentRulePoints
          .filter((p) => gridSelection.has(p.id))
          .map((p) => p.address);
        handleBatchToggleSkip(selectedAddrs, false);
        break;
      }
      case 'clear_selection':
        setGridSelection(new Set());
        break;
      default:
        break;
    }
  };

  if (!currentRule) {
    return (
      <div className="p-8 text-center text-slate-500">
        暫無配置規則，請先新增接入規則。
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6" data-testid="step2-rule-container">
      {/* 規則分頁切換列 */}
      <RuleTabRail
        rules={rules}
        devices={devices}
        selectedRuleId={currentRule.id}
        dispatch={dispatch}
      />

      {/* 雙欄核心配置區 */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* 左側：規則參數編輯器 */}
        <div className="lg:col-span-5">
          <RuleEditor
            rule={currentRule}
            devices={devices}
            globalShareEnabled={settings.modbus_share.enabled}
            dispatch={dispatch}
          />
        </div>

        {/* 右側：點位狀態網格與批次控制列 */}
        <div className="lg:col-span-7 flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <h4 className="text-xs font-semibold text-slate-300">
              當前規則點位網格
            </h4>
            <p className="text-[10px] text-slate-500 leading-normal">
              單擊網格以略過點位；配合 Shift + 點擊可批次選取區段；配合 Ctrl / ⌘ + 點擊可多選。被略過的點位將不參與後續映射轉發。
            </p>
          </div>

          <PointGridToolbar
            gridSelection={gridSelection}
            totalCount={currentRulePoints.length}
            stride={currentRulePoints[0]?.width || 1}
            onAction={handleToolbarAction}
          />

          <PointGrid
            ruleId={currentRule.id}
            points={currentRulePoints}
            conflictAddrs={conflictAddrs}
            shareLayout={currentShareLayout}
            gridSelection={gridSelection}
            setGridSelection={setGridSelection}
            onToggleSkipAddress={handleToggleSkip}
            onBatchToggleSkipAddresses={handleBatchToggleSkip}
          />
        </div>
      </div>

      {/* 全域點位合併表與工作流導航 */}
      <MergedPointTable
        points={allPoints}
        rules={rules}
        devices={devices}
        conflictAddrs={conflictAddrs}
        shareLayouts={shareLayouts}
        onContinue={onContinue}
      />
    </div>
  );
};
export default Step2Rule;
