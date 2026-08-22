import * as React from 'react';
import { useState, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { WorkbenchV2State } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { useAllPoints, useMappingValidation } from '../../state/selectors';
import { MappingTable } from './MappingTable';
import { useStep3LiveValues } from './useStep3LiveValues';

export interface Step3MappingProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  onContinue: () => void;
  onBack: () => void;
}

/**
 * 點位映射步驟容器元件 (Step 3)
 * 
 * 落地設計決策：「自動 mapping 初始化：reducer action 而非 component setState」
 * 組合映射編輯表與表格內預覽欄位，管理選取列狀態與底部繼續閥。
 */
export const Step3Mapping: React.FC<Step3MappingProps> = ({
  state,
  dispatch,
  onContinue,
  onBack,
}) => {
  const { t } = useTranslation('workbench-v2');

  // 1. 取得點位與啟用點位列表
  const allPoints = useAllPoints(state.rules, state.devices[0]?.id || 'dev-01');
  const enabledPoints = useMemo(() => allPoints.filter((p) => p.enabled && !p.skipped), [allPoints]);

  const enabledPointsIdStr = enabledPoints.map((p) => p.id).join(',');

  // 2. 自動 mapping 初始化效應
  useEffect(() => {
    dispatch({ type: 'initMappingsForPoints', points: enabledPoints });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabledPointsIdStr, dispatch]);

  // 3. 選取列管理
  const [selectedIdx, setSelectedIdx] = useState<number | null>(
    enabledPoints.length > 0 ? 0 : null
  );

  // 當點位列表發生變動時，重置選取索引
  useEffect(() => {
    setSelectedIdx(enabledPoints.length > 0 ? 0 : null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabledPointsIdStr]);

  // 4. 衍生 validation 與表格即時值
  const { canContinue, emptyTagCount, enabledCount, totalCount } = useMappingValidation(
    state.mappings
  );

  const liveValues = useStep3LiveValues(enabledPoints, state.mappings);

  return (
    <div className="space-y-6" data-testid="step3-mapping-container">
      {/* 標頭 */}
      <div>
        <h2 className="text-lg font-semibold text-slate-200">
          {t('step3.title', { defaultValue: '點位 → Tag 映射' })}
        </h2>
        <p className="text-xs text-slate-400 mt-1">
          {t('step3.subtitle', { defaultValue: '為啟用的 PLC 暫存器點位綁定對應的 Tag，並設定 Scale/Offset 線性轉換。' })}
        </p>
      </div>

      <MappingTable
        points={enabledPoints}
        mappings={state.mappings}
        selectedIdx={selectedIdx}
        setSelectedIdx={setSelectedIdx}
        devices={state.devices}
        rawValues={liveValues.rawValues}
        connectionByDevice={liveValues.connectionByDevice}
        dispatch={dispatch}
      />

      {/* 底部導覽 Footer */}
      <div className="flex items-center justify-between pt-6 border-t border-slate-900">
        {/* 左側 aside chip */}
        <div className="flex items-center">
          {enabledCount === 0 ? (
            <div
              className="bg-amber-500/10 text-amber-400 border border-amber-500/20 px-3 py-1.5 rounded-lg text-xs flex items-center gap-1.5 font-medium"
              data-testid="aside-chip-no-enabled"
            >
              <span>⚠</span>
              {t('step3.footer.noEnabled', { defaultValue: '無啟用點位' })}
            </div>
          ) : emptyTagCount > 0 ? (
            <div
              className="bg-amber-500/10 text-amber-400 border border-amber-500/20 px-3 py-1.5 rounded-lg text-xs flex items-center gap-1.5 font-medium"
              data-testid="aside-chip-empty-tag"
            >
              <span>⚠</span>
              {t('step3.footer.emptyTagWarning', {
                defaultValue: '{{count}} 個 Tag Key 留空',
                count: emptyTagCount,
              })}
            </div>
          ) : (
            <div
              className="bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-3 py-1.5 rounded-lg text-xs flex items-center gap-1.5 font-medium"
              data-testid="aside-chip-success"
            >
              <span>✓</span>
              {t('step3.footer.success', {
                defaultValue: '{{enabled}}/{{total}} 個點位已啟用',
                enabled: enabledCount,
                total: totalCount,
              })}
            </div>
          )}
        </div>

        {/* 右側按鈕組 */}
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={onBack}
            className="px-4 py-2 text-xs font-semibold text-slate-400 hover:text-slate-200 transition-colors focus:outline-none"
            data-testid="btn-back"
          >
            {t('step3.footer.back', { defaultValue: '上一步' })}
          </button>
          <button
            type="button"
            disabled={!canContinue}
            onClick={onContinue}
            className={`px-5 py-2 text-xs font-semibold rounded-lg transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-950 ${
              canContinue
                ? 'bg-blue-600 hover:bg-blue-500 text-white cursor-pointer shadow-lg shadow-blue-900/20'
                : 'bg-slate-900 text-slate-500 cursor-not-allowed border border-slate-850'
            }`}
            data-testid="btn-continue"
          >
            {t('step3.footer.continue', { defaultValue: '設定資料庫寫入' })}
          </button>
        </div>
      </div>
    </div>
  );
};
