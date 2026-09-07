import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Layers, CheckCircle2, AlertCircle, Database } from 'lucide-react';
import type { RuntimeValueEvent } from '../../../types/datalink';
import type { StudioV2RuntimeSetupContext } from '../../../types/studioV2RuntimeContext';

interface LivePointsTableProps {
  liveValues: Record<string, RuntimeValueEvent>;
  setupContext?: StudioV2RuntimeSetupContext | null;
}

function formatValue(value: unknown) {
  if (value === null || value === undefined) {
    return '-';
  }
  if (typeof value === 'object') {
    return JSON.stringify(value);
  }
  return String(value);
}

function formatRelativeTime(isoString: string): string {
  try {
    const diff = Math.floor((Date.now() - new Date(isoString).getTime()) / 1000);
    if (diff < 3) return '剛剛 (Just now)';
    if (diff < 60) return `${diff} 秒前`;
    return new Date(isoString).toLocaleTimeString();
  } catch {
    return isoString;
  }
}

export function LivePointsTable({ liveValues, setupContext }: LivePointsTableProps) {
  const { t } = useTranslation('runtime-dashboard');
  const rows = Object.values(liveValues).sort((left, right) =>
    left.address.localeCompare(right.address),
  );

  const pointMetaMap = useMemo(() => {
    const targetMap = new Map<string, string>();
    for (const tgt of setupContext?.database_targets ?? []) {
      targetMap.set(tgt.tag_id, tgt.column_name);
    }
    const map = new Map<string, { name: string; column?: string; unit?: string }>();
    for (const m of setupContext?.mappings ?? []) {
      const col = m.tag_id ? targetMap.get(m.tag_id) : undefined;
      const info = {
        name: m.display_name || m.tag_key || m.address,
        column: col,
        unit: m.unit,
      };
      map.set(m.point_id, info);
      map.set(m.address, info);
    }
    return map;
  }, [setupContext]);

  return (
    <section
      className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
      data-testid="runtime-dashboard-live-points-table"
    >
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-cyan-500/30 bg-cyan-500/10 text-cyan-300">
            <Layers className="h-5 w-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="text-base font-semibold text-slate-50">
                {t('points.title', 'Live points')} • PLC 點位即時數值表
              </h2>
              {rows.length > 0 && (
                <span className="rounded-md border border-cyan-500/30 bg-cyan-500/10 px-2 py-0.5 text-xs font-mono font-medium text-cyan-300">
                  {rows.length} 點即時更新中
                </span>
              )}
            </div>
            <p className="text-xs text-slate-400">
              {t('points.description', 'Live point rows update after the first runtime value event arrives.')} (即時顯示暫存器原始讀值與工程轉換值)
            </p>
          </div>
        </div>

        {rows.length > 0 && (
          <div className="flex items-center gap-2 text-xs text-slate-400">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-cyan-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-cyan-500" />
            </span>
            <span className="font-mono text-cyan-300">即時推播已同步</span>
          </div>
        )}
      </div>

      {rows.length === 0 ? (
        <div
          className="rounded-2xl border border-dashed border-slate-800 bg-slate-950/40 p-8 flex flex-col items-center justify-center space-y-4"
          data-testid="runtime-dashboard-live-points-placeholder"
        >
          <div className="flex items-center gap-2 text-cyan-400/90 text-sm font-medium">
            <span className="relative flex h-2.5 w-2.5">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75" />
              <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-cyan-500" />
            </span>
            {t('points.placeholder', 'Waiting for the first live point update...')}
          </div>
          <div className="w-full max-w-md space-y-2.5 animate-pulse">
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-full" />
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-11/12 mx-auto" />
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-10/12 mx-auto" />
          </div>
        </div>
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-slate-800/90 bg-slate-950/70">
          <table className="min-w-full text-left text-sm text-slate-200">
            <thead className="border-b border-slate-800 bg-slate-900/80 text-[11px] font-semibold uppercase tracking-wider text-slate-400">
              <tr>
                <th className="px-4 py-3">點位與標籤 ({t('points.address', 'Address')})</th>
                <th className="px-4 py-3">{t('points.rawValue', 'Raw value')}</th>
                <th className="px-4 py-3">{t('points.transformedValue', 'Transformed value')}</th>
                <th className="px-4 py-3">{t('points.quality', 'Quality')}</th>
                <th className="px-4 py-3">{t('points.timestamp', 'Timestamp')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60 font-mono text-xs">
              {rows.map((row) => {
                const isGood = row.quality === 'good' && !row.stale;
                const meta = pointMetaMap.get(row.point_id) || pointMetaMap.get(row.address);
                return (
                  <tr key={row.point_id} className="transition hover:bg-slate-900/60">
                    <td className="px-4 py-3 font-sans">
                      <div className="flex flex-wrap items-center gap-1.5 font-semibold text-slate-100">
                        <span>{meta?.name ?? row.address}</span>
                        {meta?.column && (
                          <span className="inline-flex items-center gap-1 rounded bg-indigo-500/15 border border-indigo-500/30 px-1.5 py-0.2 font-mono text-[10px] text-indigo-300">
                            <Database className="h-2.5 w-2.5" />
                            <span>{meta.column}</span>
                          </span>
                        )}
                      </div>
                      <div className="mt-0.5 flex items-center gap-2 font-mono text-[11px] text-cyan-300">
                        <span>暫存器: {row.address}</span>
                        <span className="sr-only">{row.address}</span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-slate-300 font-mono">
                      {formatValue(row.raw_value)}
                    </td>
                    <td className="px-4 py-3 font-bold text-emerald-300 font-mono">
                      {formatValue(row.transformed_value)}
                      {meta?.unit && <span className="ml-1 text-[10px] text-slate-400 font-normal">{meta.unit}</span>}
                    </td>
                    <td className="px-4 py-3 font-sans">
                      {row.stale ? (
                        <span className="inline-flex items-center gap-1 rounded-md border border-amber-500/30 bg-amber-500/10 px-2 py-0.5 text-[11px] font-medium text-amber-300">
                          <AlertCircle className="h-3 w-3" />
                          <span>{t('points.stale', 'stale')}</span>
                        </span>
                      ) : isGood ? (
                        <span className="inline-flex items-center gap-1 rounded-md border border-emerald-500/30 bg-emerald-500/10 px-2 py-0.5 text-[11px] font-medium text-emerald-300">
                          <CheckCircle2 className="h-3 w-3" />
                          <span>{row.quality}</span>
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 rounded-md border border-slate-700 bg-slate-800 px-2 py-0.5 text-[11px] font-medium text-slate-300">
                          <span>{row.quality}</span>
                        </span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-slate-400 font-mono text-[11px]">
                      <span title={row.timestamp}>{formatRelativeTime(row.timestamp)}</span>
                      <span className="sr-only">{row.timestamp}</span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
