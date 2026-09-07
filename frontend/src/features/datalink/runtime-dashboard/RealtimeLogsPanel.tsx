import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Terminal } from 'lucide-react';
import type { RuntimeDashboardLog } from './useRuntimeStream';

interface RealtimeLogsPanelProps {
  logs?: RuntimeDashboardLog[];
}

export function RealtimeLogsPanel({ logs = [] }: RealtimeLogsPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // 自動滾動至底部效果以保持最新日誌可見
  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
    }
  }, [logs]);

  function getLevelColor(level: RuntimeDashboardLog['level']) {
    switch (level) {
      case 'error':
        return 'text-rose-400 bg-rose-500/10 border-rose-500/30';
      case 'warn':
        return 'text-amber-400 bg-amber-500/10 border-amber-500/30';
      default:
        return 'text-cyan-400 bg-cyan-500/10 border-cyan-500/30';
    }
  }

  function formatTime(isoString: string) {
    try {
      const d = new Date(isoString);
      return d.toLocaleTimeString(undefined, { hour12: false });
    } catch {
      return isoString;
    }
  }

  return (
    <section
      className="flex h-[360px] flex-col rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
      data-testid="runtime-dashboard-logs-panel"
    >
      <div className="mb-3 flex shrink-0 items-center justify-between border-b border-slate-800/80 pb-3">
        <div className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700 bg-slate-800/60 text-slate-300">
            <Terminal className="h-4 w-4" />
          </div>
          <div>
            <h2 className="text-sm font-semibold text-slate-100">
              {t('logs.title', 'Diagnostic logs')}
            </h2>
            <p className="text-[11px] text-slate-400">
              {t('logs.description', 'Real-time diagnostic events and state transition history.')}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-1.5 font-mono text-[11px] text-slate-400">
          <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
          <span>LIVE</span>
        </div>
      </div>

      <div
        ref={scrollContainerRef}
        className="flex-1 min-h-0 overflow-y-auto rounded-2xl border border-slate-900 bg-slate-950/90 p-3.5 font-mono text-xs text-slate-300 select-text"
      >
        {logs.length === 0 ? (
          <div className="flex h-full items-center justify-center text-xs text-slate-500 italic select-none">
            {t('logs.empty', 'No diagnostic logs emitted yet.')}
          </div>
        ) : (
          <div className="space-y-2">
            {logs.map((log, idx) => (
              <div key={idx} className="flex items-start gap-2 leading-relaxed break-all">
                <span className="text-slate-500 shrink-0 select-none text-[11px]">
                  [{formatTime(log.timestamp)}]
                </span>
                <span
                  className={`rounded border px-1.5 py-0.2 text-[10px] font-semibold shrink-0 uppercase select-none ${getLevelColor(
                    log.level,
                  )}`}
                >
                  {log.level}
                </span>
                <span className="text-slate-200">{log.message}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
