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
      className="rounded-2xl border border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] px-5 py-4 text-sm text-[var(--gw-danger-fg)]"
      data-testid="settings-operation-error"
    >
      <p>{operationError.message}</p>
      {operationError.requestId ? (
        <p className="mt-1 text-xs text-[var(--gw-danger-fg)]" data-testid="settings-operation-request-id">
          {t('errors.request_id', 'Request ID')}: {operationError.requestId}
        </p>
      ) : null}
      {operationError.retryable && onRetry ? (
        <button
          type="button"
          className="mt-3 rounded-lg border border-[var(--gw-danger-border)] px-3 py-1.5 text-xs font-semibold text-[var(--gw-danger-fg)]"
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
        className="rounded-2xl border border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] px-5 py-4 text-sm text-[var(--gw-danger-fg)]"
        data-testid="settings-backend-error"
      >
        <p>{safeError?.message ?? t('settings.backend_load_error')}</p>
        {safeError?.requestId ? (
          <p className="mt-1 text-xs text-[var(--gw-danger-fg)]">
            {t('errors.request_id', 'Request ID')}: {safeError.requestId}
          </p>
        ) : null}
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl border border-[var(--gw-border)] bg-[var(--gw-surface)] px-5 py-4 text-sm text-[var(--gw-text)]"
      data-testid="settings-backend-loading"
    >
      {t('settings.backend_loading')}
    </div>
  );
}
