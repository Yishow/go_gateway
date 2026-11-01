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
    <div className="bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-6 flex flex-col justify-between h-full">
      <div className="space-y-4 flex-1">
        {/* 標題與動畫載入狀態 */}
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold text-[var(--gw-text)]">
              {t('step4.progress_title')}
            </h3>
            <p className="text-xs text-[var(--gw-muted)] mt-1">
              {t('step4.progress_subtitle')}
            </p>
          </div>
          {isCommitting && (
            <div className="flex space-x-1" data-testid="pulse-loader">
              <div className="w-2 h-2 bg-[var(--gw-active)] rounded-full" />
              <div className="w-2 h-2 bg-[var(--gw-active)] rounded-full" />
              <div className="w-2 h-2 bg-[var(--gw-active)] rounded-full" />
            </div>
          )}
        </div>

        {/* Log 滾動區域 */}
        <div className="border border-[var(--gw-border)] rounded-xl bg-[var(--gw-surface)] p-4 h-[240px] overflow-y-auto space-y-2.5 font-mono text-xs">
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
                className={`flex items-start justify-between gap-3 ${
                  isFailed ? 'text-[var(--gw-danger-fg)]' : isSkipped ? 'text-[var(--gw-warning-fg)]' : 'text-[var(--gw-text)]'
                }`}
                data-testid={`commit-log-row-${index}`}
                data-status={rowStatus}
              >
                <div className="flex items-start gap-2.5">
                  {isSuccess && (
                    <span className="text-[var(--gw-success-fg)] font-bold select-none" aria-hidden="true">
                      ✓
                    </span>
                  )}
                  {isFailed && (
                    <span className="text-[var(--gw-danger-fg)] font-bold select-none" aria-hidden="true">
                      ✗
                    </span>
                  )}
                  {isSkipped && (
                    <span className="text-[var(--gw-warning-fg)] font-bold select-none" aria-hidden="true">
                      ↷
                    </span>
                  )}
                  {isPending && (
                    <span className="text-[var(--gw-muted)] font-bold select-none" aria-hidden="true">
                      ○
                    </span>
                  )}
                  <div>
                    <span className="text-[var(--gw-text)] font-semibold block">{log.label}</span>
                    <span className="text-xs text-[var(--gw-muted)] block leading-tight mt-0.5">
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
                        className="mt-1 rounded border border-[var(--gw-danger-border)] px-1.5 py-0.5 text-xs text-[var(--gw-danger-fg)] font-sans hover:bg-[var(--gw-danger-bg)] focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--gw-danger-border)] disabled:cursor-not-allowed disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]"
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
                  <span className="px-1.5 py-0.5 rounded bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)] font-bold text-xs border border-[var(--gw-success-border)] select-none">
                    {t('step4.progress_confirmed')}
                  </span>
                )}
                {isFailed && (
                  <span className="px-1.5 py-0.5 rounded bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)] font-bold text-xs border border-[var(--gw-danger-border)] select-none">
                    ERR
                  </span>
                )}
                {isSkipped && (
                  <span className="px-1.5 py-0.5 rounded bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)] font-bold text-xs border border-[var(--gw-warning-border)] select-none">
                    SKIP
                  </span>
                )}
                {isPending && (
                  <span className="px-1.5 py-0.5 rounded bg-[var(--gw-surface)] text-[var(--gw-muted)] font-bold text-xs border border-[var(--gw-border)] select-none">
                    WAIT
                  </span>
                )}
              </div>
            );
          })}

          {logs.length === 0 && (
            <div className="h-full flex items-center justify-center text-[var(--gw-muted)] text-xs">
              {t('step4.waiting_logs')}
            </div>
          )}
        </div>
      </div>

      {/* 底部 */}
      {isCommitting && (
        <div className="mt-4 flex items-center gap-2 text-xs text-[var(--gw-accent)] select-none bg-[var(--gw-active)] border border-[var(--gw-focus)] rounded-lg p-2.5">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[var(--gw-active)] opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-[var(--gw-active)]"></span>
          </span>
          <span>{t('step4.executing_next_command')}</span>
        </div>
      )}
    </div>
  );
}
