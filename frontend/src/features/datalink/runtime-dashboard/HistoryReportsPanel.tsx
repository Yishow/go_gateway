import React, { useState, useMemo } from 'react';
import { useStudioV2WorkspaceRecordingPlansQuery } from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import {
  useExportHistoryCSVMutation,
  useStudioV2WorkspaceHistoryQuery,
} from '@/hooks/datalink/useStudioV2WorkspaceHistory';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { HistoryQuery, Resolution } from '@/types/historyReport';

export interface HistoryReportsPanelProps {
  selectedDeviceId?: string | null;
}

const HistoryReportsPanelContent: React.FC<HistoryReportsPanelProps> = ({ selectedDeviceId }) => {
  const [timeRange, setTimeRange] = useState<'1h' | '24h' | '7d'>('1h');
  const [resolution, setResolution] = useState<Resolution>('1m');

  const { data: plans = [], isLoading: isPlansLoading } = useStudioV2WorkspaceRecordingPlansQuery(
    true,
    selectedDeviceId || undefined
  );
  const activePlan = plans[0];

  const exportMutation = useExportHistoryCSVMutation();

  const queryParams: HistoryQuery = useMemo(() => {
    const now = new Date();
    let startTime = new Date(now.getTime() - 60 * 60 * 1000);
    if (timeRange === '24h') {
      startTime = new Date(now.getTime() - 24 * 60 * 60 * 1000);
    } else if (timeRange === '7d') {
      startTime = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
    }

    return {
      plan_id: activePlan?.id || '',
      start_time: startTime.toISOString(),
      end_time: now.toISOString(),
      resolution,
      limit: 100,
    };
  }, [activePlan?.id, timeRange, resolution]);

  const {
    data: historyReport,
    isLoading: isHistoryLoading,
    refetch,
  } = useStudioV2WorkspaceHistoryQuery(queryParams, !!activePlan?.id);

  const points = useMemo(() => historyReport?.points || [], [historyReport?.points]);

  const summary = useMemo(() => {
    if (points.length === 0) {
      return { count: 0, avg: 0, min: 0, max: 0, totalUsage: 0, avgCoverage: 0 };
    }
    let sum = 0;
    let min = Infinity;
    let max = -Infinity;
    let totalUsage = 0;
    let coverageSum = 0;

    for (const p of points) {
      const val = p.value_numeric ?? p.time_weighted_mean ?? 0;
      sum += val;
      if (val < min) min = val;
      if (val > max) max = val;
      if (p.usage_delta) totalUsage += p.usage_delta;
      coverageSum += p.coverage_ratio;
    }

    return {
      count: points.length,
      avg: sum / points.length,
      min: min === Infinity ? 0 : min,
      max: max === -Infinity ? 0 : max,
      totalUsage,
      avgCoverage: coverageSum / points.length,
    };
  }, [points]);

  const [exportError, setExportError] = useState<string | null>(null);

  const handleExportCSV = async () => {
    if (!activePlan?.id) return;
    try {
      setExportError(null);
      const blob = await exportMutation.mutateAsync(queryParams);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `telemetry-report-${activePlan.id}-${timeRange}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setExportError(err instanceof Error ? err.message : '匯出 CSV 失敗，請重試');
    }
  };

  if (isPlansLoading) {
    return (
      <div className="rounded-3xl border border-slate-800 bg-slate-900/60 p-6 text-sm text-slate-400">
        正在載入記錄方案與歷史查詢模組…
      </div>
    );
  }

  if (!activePlan) {
    return (
      <div
        className="rounded-3xl border border-dashed border-slate-800 bg-slate-900/40 p-6 text-center"
        data-testid="history-reports-empty"
      >
        <div className="text-slate-300 font-medium mb-1">尚未建立本設備之記錄方案</div>
        <p className="text-xs text-slate-500 max-w-md mx-auto">
          若需檢視歷史曲線、計算用量差分或匯出資料庫報表，請前往 Studio V2 步驟四完成方案建立。
        </p>
      </div>
    );
  }

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/80 p-6 space-y-5"
      data-testid="history-reports-panel"
    >
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div>
          <div className="flex items-center gap-2">
            <h3 className="text-base font-semibold text-slate-100">歷史遙測與用量報表</h3>
            <span className="rounded bg-indigo-500/10 border border-indigo-500/30 px-2 py-0.5 text-[11px] font-medium text-indigo-300">
              方案: {activePlan.name}
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            支援時序統計加權平均、累積用量差分與防注入安全 CSV 匯出
          </p>
        </div>

        <div className="flex items-center gap-2">
          <select
            className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-slate-200"
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as '1h' | '24h' | '7d')}
            data-testid="history-time-range-select"
          >
            <option value="1h">最近 1 小時</option>
            <option value="24h">最近 24 小時</option>
            <option value="7d">最近 7 天</option>
          </select>

          <select
            className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-slate-200"
            value={resolution}
            onChange={(e) => setResolution(e.target.value as Resolution)}
            data-testid="history-resolution-select"
          >
            <option value="raw">原始樣本 (Raw)</option>
            <option value="1m">1 分鐘彙總 (1m)</option>
            <option value="1h">1 小時彙總 (1h)</option>
            <option value="1d">1 天彙總 (1d)</option>
          </select>

          <button
            type="button"
            className="rounded-lg border border-slate-700 bg-slate-800 px-3 py-1.5 text-xs font-medium text-slate-200 hover:bg-slate-700 transition"
            onClick={() => void refetch()}
            disabled={isHistoryLoading}
            data-testid="history-refresh-btn"
          >
            {isHistoryLoading ? '查詢中…' : '刷新'}
          </button>

          <button
            type="button"
            className="rounded-lg border border-indigo-500/40 bg-indigo-600/30 px-3 py-1.5 text-xs font-medium text-indigo-200 hover:bg-indigo-600/50 transition"
            onClick={handleExportCSV}
            disabled={exportMutation.isPending}
            data-testid="history-export-csv-btn"
          >
            {exportMutation.isPending ? '匯出中…' : '匯出 CSV'}
          </button>
        </div>
      </div>

      {exportError && (
        <div className="rounded-lg border border-rose-800 bg-rose-950/40 p-3 text-xs text-rose-300" role="alert" data-testid="history-export-error">
          {exportError}
        </div>
      )}

      <div className="grid grid-cols-2 sm:grid-cols-5 gap-3">
        <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <div className="text-[11px] text-slate-400">數據筆數</div>
          <div className="text-lg font-semibold text-slate-100 font-mono mt-0.5" data-testid="summary-count">
            {summary.count}
          </div>
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <div className="text-[11px] text-slate-400">時序平均值</div>
          <div className="text-lg font-semibold text-sky-400 font-mono mt-0.5" data-testid="summary-avg">
            {summary.avg.toFixed(2)}
          </div>
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <div className="text-[11px] text-slate-400">極值區間 (Min~Max)</div>
          <div className="text-lg font-semibold text-amber-400 font-mono mt-0.5" data-testid="summary-range">
            {summary.min.toFixed(1)} ~ {summary.max.toFixed(1)}
          </div>
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <div className="text-[11px] text-slate-400">累計區間用量</div>
          <div className="text-lg font-semibold text-emerald-400 font-mono mt-0.5" data-testid="summary-usage">
            {summary.totalUsage.toFixed(2)}
          </div>
        </div>
        <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
          <div className="text-[11px] text-slate-400">平均有效覆蓋率</div>
          <div className="text-lg font-semibold text-purple-400 font-mono mt-0.5" data-testid="summary-coverage">
            {(summary.avgCoverage * 100).toFixed(1)}%
          </div>
        </div>
      </div>

      <div className="overflow-x-auto rounded-xl border border-slate-800 bg-slate-950/40">
        <table className="w-full text-left text-xs" data-testid="history-points-table">
          <thead className="border-b border-slate-800 bg-slate-900/80 text-[11px] text-slate-400">
            <tr>
              <th className="px-3 py-2 font-medium">觀測時間</th>
              <th className="px-3 py-2 font-medium">量測編號</th>
              <th className="px-3 py-2 font-medium">觀測數值 / 平均值</th>
              <th className="px-3 py-2 font-medium">用量差分</th>
              <th className="px-3 py-2 font-medium">品質標籤</th>
              <th className="px-3 py-2 font-medium">覆蓋率</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 text-slate-300 font-mono">
            {points.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-3 py-6 text-center text-slate-500 font-sans">
                  {isHistoryLoading ? '載入歷史資料中…' : '該區間內尚無觀測紀錄'}
                </td>
              </tr>
            ) : (
              points.slice(0, 10).map((p, idx) => (
                <tr key={`${p.measurement_id}-${p.observed_at}-${idx}`} className="hover:bg-slate-900/40">
                  <td className="px-3 py-2 text-slate-400">
                    {new Date(p.observed_at).toLocaleTimeString()}
                  </td>
                  <td className="px-3 py-2 text-slate-300">{p.measurement_id}</td>
                  <td className="px-3 py-2 text-sky-300">
                    {p.time_weighted_mean !== undefined
                      ? p.time_weighted_mean.toFixed(2)
                      : p.value_numeric !== undefined
                      ? p.value_numeric.toFixed(2)
                      : p.value_string || '-'}
                  </td>
                  <td className="px-3 py-2 text-emerald-400">
                    {p.usage_delta !== undefined ? `+${p.usage_delta.toFixed(2)}` : '-'}
                  </td>
                  <td className="px-3 py-2">
                    <span
                      className={`inline-block px-1.5 py-0.5 rounded text-[10px] ${
                        p.quality === 'good'
                          ? 'bg-emerald-500/10 text-emerald-300 border border-emerald-500/20'
                          : 'bg-amber-500/10 text-amber-300 border border-amber-500/20'
                      }`}
                    >
                      {p.quality}
                    </span>
                  </td>
                  <td className="px-3 py-2 text-slate-400">
                    {(p.coverage_ratio * 100).toFixed(0)}%
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
};

export const HistoryReportsPanel: React.FC<HistoryReportsPanelProps> = (props) => {
  return (
    <SafeQueryBoundary>
      <HistoryReportsPanelContent {...props} />
    </SafeQueryBoundary>
  );
};
