import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { ReadinessStage } from '../../state/types';
import { Icon } from '../../components/Icon';

export interface ReadinessStagesProps {
  /** 診斷步驟狀態列表 */
  stages: ReadinessStage[];
  /** 整體連線測試狀態 */
  testStatus?: 'running' | 'success' | 'failed' | 'idle' | 'draft';
}

/**
 * 連線與診斷進度展示元件 (ReadinessStages)
 * 
 * 落地設計決策：「Readiness 群組：以 connect / probe 語義標籤」
 * 實作需求 **Readiness check with explicit connect / probe separation**。
 * 將診斷步驟區分為「連線通道（connect）」與「協議握手（probe）」兩大區段，並顯示合計延遲。
 */
export const ReadinessStages: React.FC<ReadinessStagesProps> = ({ stages }) => {
  const { t } = useTranslation('workbench-v2');

  // 分群組過濾
  const connectStages = stages.filter((s) => s.group === 'connect');
  const probeStages = stages.filter((s) => s.group === 'probe');

  // 計算特定群組的合計 latency
  const calculateGroupLatency = (groupStages: ReadinessStage[]) => {
    return groupStages.reduce((sum, s) => sum + (s.latency_ms ?? 0), 0);
  };

  const connectLatency = calculateGroupLatency(connectStages);
  const probeLatency = calculateGroupLatency(probeStages);

  /**
   * 渲染單一診斷步驟
   */
  const renderStageItem = (stage: ReadinessStage) => {
    let iconNode: React.ReactNode = null;
    let circleClass = 'border-slate-800 bg-slate-900/50 text-slate-600';
    let textClass = 'text-slate-500';

    if (stage.status === 'success') {
      circleClass = 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400';
      textClass = 'text-slate-300';
      iconNode = <Icon name="check" className="h-3 w-3 stroke-[2.5]" />;
    } else if (stage.status === 'skipped') {
      circleClass = 'border-slate-700 bg-slate-950/40 text-slate-500';
      textClass = 'text-slate-400';
      iconNode = <div className="h-px w-2.5 bg-slate-600" />;
    } else if (stage.status === 'running') {
      circleClass = 'border-blue-500/30 bg-blue-500/10 text-blue-400';
      textClass = 'text-slate-100 font-medium';
      iconNode = <Icon name="refresh" className="h-3 w-3 animate-spin stroke-[2]" />;
    } else if (stage.status === 'failed') {
      circleClass = 'border-rose-500/30 bg-rose-500/10 text-rose-400';
      textClass = 'text-rose-300';
      iconNode = <Icon name="close" className="h-3 w-3 stroke-[2.5]" />;
    } else {
      // pending
      circleClass = 'border-slate-800 bg-slate-950/40 text-slate-600';
      textClass = 'text-slate-500';
      iconNode = <div className="h-1.5 w-1.5 rounded-full bg-slate-700" />;
    }

    return (
      <div key={stage.id} className="flex items-center justify-between py-1.5" data-testid={`stage-item-${stage.id}`}>
        <div className="flex items-center gap-3">
          {/* 狀態圓圈 */}
          <div
            className={`flex h-6 w-6 items-center justify-center rounded-full border transition-all duration-300 ${circleClass}`}
            data-testid={`stage-circle-${stage.id}`}
          >
            {iconNode}
          </div>
          {/* 步驟名稱 */}
          <span className={`text-xs transition-colors duration-300 ${textClass}`}>
            {t(stage.label)}
          </span>
        </div>
        {/* 延遲時間 */}
        {stage.status === 'success' && stage.latency_ms !== undefined && (
          <span className="text-xs font-mono text-slate-500" data-testid={`stage-latency-${stage.id}`}>
            {stage.latency_ms} ms
          </span>
        )}
        {stage.status === 'failed' && stage.message && (
          <span className="text-[10px] text-rose-400 max-w-[150px] truncate" title={stage.message}>
            {stage.message}
          </span>
        )}
      </div>
    );
  };

  return (
    <div className="space-y-4">
      {/* 連線通道區段 */}
      {connectStages.length > 0 && (
        <div data-testid="group-connect">
          <div className="flex items-center justify-between border-b border-slate-800/80 pb-1.5 mb-2">
            <span className="text-xs font-semibold text-slate-400">
              {t('step1.groups.connect')}
            </span>
            {connectLatency > 0 && (
              <span className="text-[10px] font-mono text-slate-500" data-testid="connect-total-latency">
                {connectLatency} ms
              </span>
            )}
          </div>
          <div className="space-y-1.5">
            {connectStages.map(renderStageItem)}
          </div>
        </div>
      )}

      {/* 協議握手區段 */}
      {probeStages.length > 0 && (
        <div data-testid="group-probe">
          <div className="flex items-center justify-between border-b border-slate-800/80 pb-1.5 mb-2">
            <span className="text-xs font-semibold text-slate-400">
              {t('step1.groups.probe')}
            </span>
            {probeLatency > 0 && (
              <span className="text-[10px] font-mono text-slate-500" data-testid="probe-total-latency">
                {probeLatency} ms
              </span>
            )}
          </div>
          <div className="space-y-1.5">
            {probeStages.map(renderStageItem)}
          </div>
        </div>
      )}
    </div>
  );
};
