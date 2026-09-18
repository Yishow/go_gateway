import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { CommitLog } from '../../state/types';
import { getSafeErrorMessage } from '../../../../../utils/typedErrors';

/**
 * CommitProgress 元件屬性
 */
interface CommitProgressProps {
  logs: CommitLog[];
  status: 'idle' | 'committing' | 'success' | 'failed';
  onRetry?: (log: CommitLog) => Promise<void> | void;
}

/**
 * 提交進度動畫日誌列表元件
 * 落地需求：「Commit sequence and animation」之動畫日誌與進度狀態
 */
export function CommitProgress({ logs, status, onRetry }: CommitProgressProps) {
  const { t } = useTranslation('workbench-v2');

  const isCommitting = status === 'committing';
  const [retryingIndex, setRetryingIndex] = React.useState<number | null>(null);
  const retryingRef = React.useRef<number | null>(null);

  const handleRetry = async (log: CommitLog, index: number) => {
    if (!onRetry || isCommitting || retryingRef.current !== null) {
      return;
    }
    retryingRef.current = index;
    setRetryingIndex(index);
    try {
      await onRetry(log);
    } finally {
      retryingRef.current = null;
      setRetryingIndex(null);
    }
  };

  return (
    <div className="bg-gray-900/10 border border-gray-800 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm">
      <div className="space-y-4 flex-1">
        {/* 標題與動畫載入狀態 */}
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-white">
              {t('step4.progress_title')}
            </h3>
            <p className="text-xs text-gray-500 mt-1">
              {t('step4.progress_subtitle')}
            </p>
          </div>
          {isCommitting && (
            <div className="flex space-x-1" data-testid="pulse-loader">
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
              <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
            </div>
          )}
        </div>

        {/* Log 滾動區域 */}
        <div className="border border-gray-800/60 rounded-xl bg-gray-950/40 p-4 h-[240px] overflow-y-auto space-y-2.5 font-mono text-xs scrollbar-thin scrollbar-thumb-gray-800">
          {logs.map((log, index) => {
            const rowStatus = log.status ?? 'pending';
            const isFailed = rowStatus === 'failed';
            const isSkipped = rowStatus === 'skipped';
            const isPending = rowStatus === 'pending' || rowStatus === 'running';
            const isSuccess = rowStatus === 'success';
            const safeError = isFailed
              ? getSafeErrorMessage(
                  { code: log.code, request_id: log.request_id, retryable: log.retryable === true },
                  (key, options) => t(key, options),
                )
              : null;

            return (
              <div
                key={index}
                className={`flex items-start justify-between gap-3 animate-[slideIn_0.2s_ease-out] ${
                  isFailed ? 'text-rose-300' : isSkipped ? 'text-amber-300' : 'text-gray-300'
                }`}
                data-testid={`commit-log-row-${index}`}
                data-status={rowStatus}
              >
                <div className="flex items-start gap-2.5">
                  {isSuccess && (
                    <span className="text-emerald-500 font-bold select-none" aria-hidden="true">
                      ✓
                    </span>
                  )}
                  {isFailed && (
                    <span className="text-rose-500 font-bold select-none" aria-hidden="true">
                      ✗
                    </span>
                  )}
                  {isSkipped && (
                    <span className="text-amber-500 font-bold select-none" aria-hidden="true">
                      ↷
                    </span>
                  )}
                  {isPending && (
                    <span className="text-gray-500 font-bold select-none" aria-hidden="true">
                      ○
                    </span>
                  )}
                  <div>
                    <span className="text-gray-200 font-semibold block">{log.label}</span>
                    <span className="text-[10px] text-gray-500 block leading-tight mt-0.5">
                      {isFailed ? (
                        <>
                          {safeError?.message}
                          {safeError?.requestId && (
                            <span className="ml-1">({t('errors.request_id')}: {safeError.requestId})</span>
                          )}
                        </>
                      ) : t(`step4.progress_status.${rowStatus}`)}
                    </span>
                    {isFailed && safeError?.retryable && (
                      <button
                        type="button"
                        className="mt-1 rounded border border-rose-500/40 px-1.5 py-0.5 text-[10px] text-rose-300 font-sans hover:bg-rose-500/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-400 disabled:cursor-not-allowed disabled:opacity-50"
                        disabled={!onRetry || isCommitting || retryingIndex !== null}
                        onClick={() => void handleRetry(log, index)}
                        aria-label={retryingIndex === index ? t('step4.retrying') : t('step4.retry_action')}
                      >
                        {retryingIndex === index ? t('step4.retrying') : t('step4.retry_action')}
                      </button>
                    )}
                  </div>
                </div>
                {isSuccess && (
                  <span className="px-1.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 font-bold text-[10px] border border-emerald-500/20 select-none">
                    {t('step4.progress_confirmed')}
                  </span>
                )}
                {isFailed && (
                  <span className="px-1.5 py-0.5 rounded bg-rose-500/10 text-rose-400 font-bold text-[10px] border border-rose-500/20 select-none">
                    ERR
                  </span>
                )}
                {isSkipped && (
                  <span className="px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-400 font-bold text-[10px] border border-amber-500/20 select-none">
                    SKIP
                  </span>
                )}
                {isPending && (
                  <span className="px-1.5 py-0.5 rounded bg-gray-500/10 text-gray-400 font-bold text-[10px] border border-gray-500/20 select-none">
                    WAIT
                  </span>
                )}
              </div>
            );
          })}

          {logs.length === 0 && (
            <div className="h-full flex items-center justify-center text-gray-500 text-xs">
              {t('step4.waiting_logs')}
            </div>
          )}
        </div>
      </div>

      {/* 底部 pulse-dot */}
      {isCommitting && (
        <div className="mt-4 flex items-center gap-2 text-xs text-blue-400 select-none bg-blue-500/5 border border-blue-500/10 rounded-lg p-2.5 animate-pulse">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
          </span>
          <span>{t('step4.executing_next_command')}</span>
        </div>
      )}
    </div>
  );
}
