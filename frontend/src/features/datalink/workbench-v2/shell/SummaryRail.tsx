import * as React from 'react';
import type { WorkbenchV2State } from '../state/types';
import { Icon } from '../components';
import type { StudioV2WorkspaceAuditEntry } from '../../../../types/studioV2WorkspaceAudit';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../types/studioV2WorkspaceReadiness';
import { WorkspaceReadinessPanel } from '../components/WorkspaceReadinessPanel';
import { WorkspaceAuditHistoryPanel } from './WorkspaceAuditHistoryPanel';

export interface SummaryRailProps {
  state: WorkbenchV2State;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
  workspaceAuditHistory?: StudioV2WorkspaceAuditEntry[];
  workspaceAuditUnavailable?: boolean;
  onOpenRuntime?: () => void;
}

/**
 * Workbench V2 右側即時摘要列元件
 * 
 * 落地設計決策：「Shell layout regions」
 * 根據 state 即時推導並顯示設備、接入規則、點位映射與資料庫摘要，
 * 並在視窗寬度 < 1280px 時完全不渲染至 DOM 中。
 */
export const SummaryRail: React.FC<SummaryRailProps> = ({
  state,
  workspaceReadiness,
  workspaceAuditHistory,
  workspaceAuditUnavailable,
  onOpenRuntime,
}) => {
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
      <span className="text-xs text-[var(--gw-muted)] flex-shrink-0">
        {label}
      </span>
      <span
        className={`text-xs text-right truncate ${mono ? 'font-mono' : ''} ${
          dim ? 'text-[var(--gw-muted)]' : 'text-[var(--gw-text)]'
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
  const hasRunningDevice = state.devices.some((device) => device.running);

  return (
    <div className="space-y-3 w-[260px] flex-shrink-0" data-testid="summary-rail">
      <WorkspaceReadinessPanel
        summary={workspaceReadiness}
        dataTestId="summary-rail-readiness"
        compact
      />
      <WorkspaceAuditHistoryPanel
        entries={workspaceAuditHistory}
        unavailable={workspaceAuditUnavailable}
      />

      {/* 設備卡片 */}
      <div className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="device" className="w-3.5 h-3.5 text-[var(--gw-muted)]" />
          <span className="text-xs font-semibold text-[var(--gw-muted)]">
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
          {devices.slice(0, 4).map((d) => {
            return (
              <div key={d.id} className="flex items-center gap-1.5 text-xs font-mono text-[var(--gw-muted)]">
                <span className="w-1.5 h-1.5 rounded-full flex-shrink-0 bg-[var(--gw-muted)]" />
                <span className="truncate text-[var(--gw-muted)]">{d.name}</span>
                <span className="ml-auto text-[var(--gw-muted)]">{d.protocol}</span>
              </div>
            );
          })}
          {devices.length > 4 && (
            <div className="text-xs text-[var(--gw-muted)] pl-3">+ {devices.length - 4} 更多</div>
          )}
        </div>
      </div>

      {/* 規則卡片 */}
      <div className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="rule" className="w-3.5 h-3.5 text-[var(--gw-muted)]" />
          <span className="text-xs font-semibold text-[var(--gw-muted)]">
            接入規則
          </span>
        </div>
        <Row label="條數" value={`${totalRanges} 條`} mono />
        <Row label="點位" value={`${enabledPoints.length || totalPlanned} 個`} mono />
      </div>

      {/* 映射卡片 */}
      <div className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="map" className="w-3.5 h-3.5 text-[var(--gw-muted)]" />
          <span className="text-xs font-semibold text-[var(--gw-muted)]">
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
      <div className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="db" className="w-3.5 h-3.5 text-[var(--gw-muted)]" />
          <span className="text-xs font-semibold text-[var(--gw-muted)]">
            資料庫連線
          </span>
        </div>
        {connector && connector.name ? (
          <>
            <Row label="名稱" value={connector.name} />
            <Row label="類型" value={connector.kind} mono />
            <p className="mt-2 text-xs text-[var(--gw-muted)]">表格、記錄間隔與量測點以 Basic／Advanced 群組設定為準。</p>
          </>
        ) : (
          <div className="text-xs text-[var(--gw-muted)] italic">尚未設定</div>
        )}
      </div>

      {/* 部署成功提示 */}
      {hasRunningDevice && (
        <div className="rounded-xl border border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] p-3">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <div className="flex items-center gap-2 text-[var(--gw-success-fg)] text-xs font-semibold">
                <Icon name="check" className="w-3.5 h-3.5" />
                已有設備在 Runtime 執行
              </div>
              <div className="text-xs text-[var(--gw-success-fg)] mt-1">
                可直接進入 Runtime 觀察已成功啟動的設備。
              </div>
            </div>
            <button
              type="button"
              data-testid="summary-rail-runtime-link"
              onClick={onOpenRuntime}
              className="shrink-0 rounded-md border border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] px-2 py-1 text-xs font-semibold text-[var(--gw-success-fg)] transition-colors hover:bg-[var(--gw-success-bg)]"
            >
              前往 Runtime
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
