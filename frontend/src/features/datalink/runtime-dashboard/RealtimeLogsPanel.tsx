import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { RuntimeDashboardLog } from './useRuntimeStream';

interface RealtimeLogsPanelProps {
  logs?: RuntimeDashboardLog[];
}

export function RealtimeLogsPanel({ logs = [] }: RealtimeLogsPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // 實作自動滾動至底部效果以保持最新日誌可見
  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
    }
  }, [logs]);

  function getLevelColor(level: RuntimeDashboardLog['level']) {
    switch (level) {
      case 'error':
        return 'text-rose-400';
      case 'warn':
        return 'text-amber-400';
      default:
        return 'text-cyan-400';
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
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 flex flex-col h-[280px]"
      data-testid="runtime-dashboard-logs-panel"
    >
      <div className="mb-3 shrink-0">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('logs.title', 'Diagnostic logs')}
        </h2>
        <p className="mt-0.5 text-xs text-slate-400">
          {t('logs.description', 'Real-time diagnostic events and state transition history.')}
        </p>
      </div>

      <div
        ref={scrollContainerRef}
        className="flex-1 min-h-0 overflow-y-auto rounded-2xl bg-slate-950/80 p-4 font-mono text-xs text-slate-300 border border-slate-900 select-text"
      >
        {logs.length === 0 ? (
          <div className="text-slate-600 italic select-none">
            {t('logs.empty', 'No diagnostic logs emitted yet.')}
          </div>
        ) : (
          <div className="space-y-1.5">
            {logs.map((log, idx) => (
              <div key={idx} className="flex items-start gap-2 leading-relaxed break-all">
                <span className="text-slate-500 shrink-0 select-none">
                  [{formatTime(log.timestamp)}]
                </span>
                <span className={`font-semibold shrink-0 uppercase select-none ${getLevelColor(log.level)}`}>
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
