import { getSafeErrorMessage, type SafeErrorMessage } from '../../../../utils/typedErrors';
import { useTranslation } from 'react-i18next';

interface SettingsStatusProps {
  operationError: SafeErrorMessage | null;
  onRetry?: () => void;
}

export function SettingsStatus({ operationError, onRetry }: SettingsStatusProps) {
  const { t } = useTranslation('workbench-v2');
  if (!operationError) {
    return null;
  }

  return (
    <div
      className="rounded-2xl border border-red-500/20 bg-red-500/5 px-5 py-4 text-sm text-red-200"
      data-testid="settings-operation-error"
    >
      <p>{operationError.message}</p>
      {operationError.requestId ? (
        <p className="mt-1 text-xs text-red-200/80" data-testid="settings-operation-request-id">
          {t('errors.request_id', 'Request ID')}: {operationError.requestId}
        </p>
      ) : null}
      {operationError.retryable && onRetry ? (
        <button
          type="button"
          className="mt-3 rounded-lg border border-red-300/30 px-3 py-1.5 text-xs font-semibold text-red-100"
          onClick={onRetry}
        >
          {operationError.action ?? t('errors.retry', 'Retry')}
        </button>
      ) : null}
    </div>
  );
}

interface SettingsBackendStatusProps {
  hasError: boolean;
  error?: unknown;
}

export function SettingsBackendStatus({ hasError, error }: SettingsBackendStatusProps) {
  const { t } = useTranslation('workbench-v2');
  if (hasError) {
    const safeError = error
      ? getSafeErrorMessage(error, (key, options) => t(key, options))
      : null;
    return (
      <div
        className="rounded-2xl border border-red-500/20 bg-red-500/5 px-5 py-4 text-sm text-red-200"
        data-testid="settings-backend-error"
      >
        <p>{safeError?.message ?? t('settings.backend_load_error')}</p>
        {safeError?.requestId ? (
          <p className="mt-1 text-xs text-red-200/80">
            {t('errors.request_id', 'Request ID')}: {safeError.requestId}
          </p>
        ) : null}
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl border border-slate-700/60 bg-slate-900/40 px-5 py-4 text-sm text-slate-300"
      data-testid="settings-backend-loading"
    >
      {t('settings.backend_loading')}
    </div>
  );
}
