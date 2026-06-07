import * as React from 'react';
import type { WorkbenchV2State } from '../state/types';
import { Icon } from '../components';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../types/studioV2WorkspaceReadiness';
import { WorkspaceReadinessPanel } from '../components/WorkspaceReadinessPanel';

export interface SummaryRailProps {
  state: WorkbenchV2State;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
}

/**
 * Workbench V2 右側即時摘要列元件
 * 
 * 落地設計決策：「Shell layout regions」
 * 根據 state 即時推導並顯示設備、接入規則、點位映射與資料庫摘要，
 * 並在視窗寬度 < 1280px 時完全不渲染至 DOM 中。
 */
export const SummaryRail: React.FC<SummaryRailProps> = ({ state, workspaceReadiness }) => {
  const [isWide, setIsWide] = React.useState<boolean>(() => {
    if (typeof window !== 'undefined') {
      return window.innerWidth >= 1280;
    }
    return true;
  });

  React.useEffect(() => {
    if (typeof window === 'undefined') return;

    const handleResize = () => {
      setIsWide(window.innerWidth >= 1280);
    };

    window.addEventListener('resize', handleResize);
    // 初始化呼叫一次，以防測試中手動修改 width 後未觸發 resize
    handleResize();

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  if (!isWide) return null;

  const { devices = [], rules = [], points = [], mappings = {}, db } = state;
  const enabledPoints = points.filter((p) => !p.skipped);
  const connector = db?.connector;

  // 輔助列渲染
  const Row = ({
    label,
    value,
    mono = false,
    dim = false,
  }: {
    label: string;
    value: string | number;
    mono?: boolean;
    dim?: boolean;
  }) => (
    <div className="flex items-baseline justify-between gap-3 py-1">
      <span className="text-[10px] uppercase tracking-wider text-slate-500 flex-shrink-0">
        {label}
      </span>
      <span
        className={`text-xs text-right truncate ${mono ? 'font-mono' : ''} ${
          dim ? 'text-slate-500' : 'text-slate-200'
        }`}
      >
        {value}
      </span>
    </div>
  );

  const totalRanges = rules.filter((r) => r.enabled).length;
  // 推導點位數量：使用 rules 的 planned 點位數作為備用
  const totalPlanned = rules
    .filter((r) => r.enabled)
    .reduce((sum, r) => sum + r.count - (r.skipped_addresses?.length || 0), 0);
  const testedDevs = devices.filter((d) => d.test?.status === 'success').length;

  return (
    <div className="space-y-3 w-[260px] flex-shrink-0" data-testid="summary-rail">
      <WorkspaceReadinessPanel summary={workspaceReadiness} dataTestId="summary-rail-readiness" maxIssues={3} />

      {/* 設備卡片 */}
      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="device" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            設備
          </span>
        </div>
        <Row label="總數" value={`${devices.length} 個`} mono />
        <Row
          label="已測試"
          value={`${testedDevs} / ${devices.length}`}
          mono
          dim={testedDevs !== devices.length}
        />
        <div className="mt-1.5 space-y-0.5">
          {devices.slice(0, 4).map((d, i) => {
            const colors = ['bg-blue-500', 'bg-emerald-500', 'bg-amber-500', 'bg-fuchsia-500'];
            return (
              <div key={d.id} className="flex items-center gap-1.5 text-[10px] font-mono text-slate-500">
                <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${colors[i % colors.length]}`} />
                <span className="truncate text-slate-400">{d.name}</span>
                <span className="ml-auto text-slate-500">{d.protocol}</span>
              </div>
            );
          })}
          {devices.length > 4 && (
            <div className="text-[10px] text-slate-600 pl-3">+ {devices.length - 4} 更多</div>
          )}
        </div>
      </div>

      {/* 規則卡片 */}
      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="rule" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            接入規則
          </span>
        </div>
        <Row label="條數" value={`${totalRanges} 條`} mono />
        <Row label="點位" value={`${enabledPoints.length || totalPlanned} 個`} mono />
      </div>

      {/* 映射卡片 */}
      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="map" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            映射
          </span>
        </div>
        <Row
          label="Tags"
          value={Object.values(mappings).filter((m) => m.enabled).length || '—'}
          mono
        />
        <Row
          label="個別轉換"
          value={Object.values(mappings).filter((m) => m.scale !== 1 || m.offset !== 0).length}
          mono
        />
      </div>

      {/* 資料庫卡片 */}
      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="db" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            資料庫寫入
          </span>
        </div>
        {connector && connector.name ? (
          <>
            <Row label="類型" value={connector.kind} mono />
            <Row label="表" value={`${connector.schema}.${connector.table}`} mono />
            <Row label="間隔" value={`${connector.write_interval_seconds}s`} mono />
            <Row
              label="目標"
              value={`${db?.targets ? Object.values(db.targets).filter((t) => t.enabled).length : 0} 欄位`}
              mono
            />
          </>
        ) : (
          <div className="text-[11px] text-slate-500 italic">尚未設定</div>
        )}
      </div>

      {/* 部署成功提示 */}
      {state.devices.some((device) => device.running) && (
        <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/[0.08] p-3 sweep-in">
          <div className="flex items-center gap-2 text-emerald-200 text-xs font-semibold">
            <Icon name="check" className="w-3.5 h-3.5" />
            已有設備在 Runtime 執行
          </div>
          <div className="text-[10px] text-emerald-300/70 mt-1">可前往 Runtime 觀察已成功啟動的設備。</div>
        </div>
      )}
    </div>
  );
};
