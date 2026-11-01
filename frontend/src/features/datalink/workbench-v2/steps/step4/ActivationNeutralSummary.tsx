import * as React from 'react';
import type { TFunction } from 'i18next';
import { useTranslation } from 'react-i18next';
import type {
  StudioV2ActivationRecovery,
  StudioV2ActivationResponse,
} from '../../../../../types/studioV2Activation';

function resultStatusLabel(
  status: StudioV2ActivationResponse['results'][number]['status'],
  t: TFunction<'workbench-v2'>,
) {
  return t(`step4.result_status.${status}`, {
    defaultValue: status === 'success'
      ? 'Confirmed'
      : status === 'failed'
        ? 'Rejected'
        : status === 'pending'
          ? 'Pending'
          : 'Skipped',
  });
}

export function ActivationRecoveryNotice({
  recovery,
  recoveryUnavailable = false,
}: {
  recovery?: StudioV2ActivationRecovery | null;
  recoveryUnavailable?: boolean;
}) {
  const { t } = useTranslation('workbench-v2');
  if (!recovery && !recoveryUnavailable) {
    return null;
  }

  const runningDevices = recovery?.devices.filter((device) => device.running) ?? [];

  return (
    <div className="space-y-3">
      {recovery && (
        <div
          data-testid="activation-recovery-status"
          className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] px-3 py-2 text-xs text-[var(--gw-text)]"
        >
          <div>
            {runningDevices.length > 0
              ? t('step4.activation_recovery_running', {
                count: runningDevices.length,
              })
              : t('step4.activation_recovery_read')}
          </div>
          {recovery.operation_id && (
            <div className="mt-1 font-mono text-xs text-[var(--gw-muted)]">
              {t('step4.activation_operation_id')}: {recovery.operation_id}
            </div>
          )}
        </div>
      )}
      {recoveryUnavailable && (
        <div
          data-testid="activation-recovery-unavailable"
          className="rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-3 py-2 text-xs text-[var(--gw-warning-fg)]"
        >
          {t('step4.activation_recovery_unavailable')}
        </div>
      )}
    </div>
  );
}

export function ActivationNeutralSummary({
  onReset,
  canContinue,
  onCommit,
  response,
  recovery,
  recoveryUnavailable = false,
}: {
  onReset: () => void;
  canContinue?: boolean;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  response?: StudioV2ActivationResponse | null;
  recovery?: StudioV2ActivationRecovery | null;
  recoveryUnavailable?: boolean;
}) {
  const { t } = useTranslation('workbench-v2');
  const handoffTriggeredRef = React.useRef(false);
  const hasUnresolvedResult = Boolean(
    response?.results.some((result) => result.status === 'pending' || result.status === 'skipped'),
  );
  const isUnconfirmed = response
    ? response.outcome === 'unconfirmed' ||
      (response.outcome !== 'failed' && (response.results.length === 0 || hasUnresolvedResult))
    : recoveryUnavailable;
  const handleCommit = () => {
    if (!onCommit || handoffTriggeredRef.current) {
      return;
    }
    handoffTriggeredRef.current = true;
    onCommit();
  };

  return (
    <div
      data-testid="activation-empty-message"
      data-outcome={isUnconfirmed ? 'unconfirmed' : 'failed'}
      className="rounded-2xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-5 py-6 text-sm text-[var(--gw-warning-fg)]"
    >
      <p>
        {isUnconfirmed
          ? t('step4.activation_unconfirmed')
          : t('step4.activation_no_success')}
      </p>
      {(response?.operation_id || response?.request_id) && (
        <div
          data-testid="activation-response-identifiers"
          className="mt-3 space-y-1 font-mono text-xs text-[var(--gw-warning-fg)]"
        >
          {response.operation_id && (
            <div>
              {t('step4.activation_operation_id')}: {response.operation_id}
            </div>
          )}
          {response.request_id && (
            <div>
              {t('step4.activation_request_id')}: {response.request_id}
            </div>
          )}
        </div>
      )}
      {response && response.results.length > 0 && (
        <div
          data-testid="activation-neutral-results"
          className="mt-3 space-y-2"
        >
          {response.results.map((result) => (
            <div
              key={result.device_id}
              data-testid={`activation-neutral-result-${result.device_id}`}
              className="flex items-center justify-between gap-3 rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] px-3 py-2 text-xs"
            >
              <span className="font-mono text-[var(--gw-text)]">{result.device_id}</span>
              <span className="text-[var(--gw-warning-fg)]">{resultStatusLabel(result.status, t)}</span>
            </div>
          ))}
        </div>
      )}
      {response && (
        <div
          data-testid="activation-delivery-status"
          className="mt-3 rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-4 py-3 text-xs text-[var(--gw-warning-fg)]"
        >
          <span className="font-semibold">
            {t('step4.activation_delivery_status_label')}:
          </span>{' '}
          {t('step4.activation_delivery_unconfirmed')}
        </div>
      )}
      <ActivationRecoveryNotice
        recovery={recovery}
        recoveryUnavailable={recoveryUnavailable && (!response || isUnconfirmed)}
      />
      <div className="mt-4 flex flex-wrap items-center gap-3">
        {canContinue && onCommit && (
          <button
            type="button"
            className="rounded-xl bg-[var(--gw-primary)] px-4 py-2 font-medium text-[var(--gw-on-primary)] transition-all hover:bg-[var(--gw-primary-hover)]"
            onClick={handleCommit}
          >
            {t('step4.go_to_dashboard_btn')}
          </button>
        )}
        <button
          type="button"
          className="rounded-xl border border-[var(--gw-control)] px-4 py-2 hover:bg-[var(--gw-surface)]"
          onClick={onReset}
        >
          {t('step4.reset_activation_btn')}
        </button>
      </div>
    </div>
  );
}
