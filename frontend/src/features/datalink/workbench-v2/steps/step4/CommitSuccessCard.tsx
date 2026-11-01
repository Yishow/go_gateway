import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioV2ActivationResponse } from '../../../../../types/studioV2Activation';
import { normalizeTypedEnvelope } from '../../../../../utils/safeJson';

/**
 * CommitSuccessCard 元件屬性
 */
interface CommitSuccessCardProps {
  response: StudioV2ActivationResponse;
  canContinue: boolean;
  onCommit?: (confirmedDeviceIds?: string[]) => void;
  onReset: () => void;
}

/**
 * 提交完成成功展示卡片元件
 * 落地需求：「Commit completion card」之綠色完成卡片與導向按鈕
 */
export function CommitSuccessCard({
  response,
  canContinue,
  onCommit,
  onReset,
}: CommitSuccessCardProps) {
  const { t } = useTranslation('workbench-v2');
  const responseEnvelope = normalizeTypedEnvelope(response);
  const handoffTriggeredRef = React.useRef(false);
  const confirmedDeviceIds = response.results
    .filter((result) => result.status === 'success')
    .map((result) => result.device_id);
  const isPartial = response.results.length === 0 ||
    response.results.some((result) => result.status !== 'success');
  const cardTone = isPartial ? {
    card: 'bg-[var(--gw-surface)] border-[var(--gw-warning-border)]',
    icon: 'border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)]',
    title: 'text-[var(--gw-warning-fg)]',
    subtitle: 'text-[var(--gw-warning-fg)]',
  } : {
    card: 'bg-[var(--gw-surface)] border-[var(--gw-success-border)]',
    icon: 'border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]',
    title: 'text-[var(--gw-success-fg)]',
    subtitle: 'text-[var(--gw-success-fg)]',
  };

  const handleCommit = () => {
    if (!onCommit || handoffTriggeredRef.current) {
      return;
    }
    handoffTriggeredRef.current = true;
    onCommit(confirmedDeviceIds);
  };

  const resultStatusLabel = (status: StudioV2ActivationResponse['results'][number]['status']) =>
    t(`step4.result_status.${status}`, {
      defaultValue: status === 'success'
        ? 'Confirmed'
        : status === 'failed'
          ? 'Rejected'
          : status === 'pending'
            ? 'Pending'
            : 'Skipped',
    });

  return (
    <div
      data-testid="activation-results-card"
      data-outcome={isPartial ? 'partial' : 'confirmed'}
      className={`${cardTone.card} border rounded-2xl p-6 flex flex-col justify-between h-full`}
    >
      <div className="space-y-6 flex-1">
        {/* 結果 Icon */}
        <div className="space-y-2 text-center">
          <div className={`mx-auto flex h-16 w-16 items-center justify-center rounded-full border text-3xl ${cardTone.icon}`}>
            {isPartial ? '!' : '✓'}
          </div>
          <h3 className={`text-xl font-bold ${cardTone.title}`}>
            {t('step4.results_title')}
          </h3>
          <p className={`text-sm ${cardTone.subtitle}`}>
            {t('step4.results_subtitle')}
          </p>
        </div>

        {responseEnvelope.code && (
          <div
            data-testid="activation-empty-message"
            className="rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-4 py-3 text-sm text-[var(--gw-warning-fg)]"
          >
            {t(`errors.${responseEnvelope.code}`, {
              defaultValue: t('errors.activation_failed'),
            })}
          </div>
        )}

        <div className="space-y-2">
          {response.results.map((result) => (
            <div
              key={result.device_id}
              data-testid={`activation-result-${result.device_id}`}
              className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] px-4 py-3"
            >
              <div className="flex items-center justify-between gap-3">
                <span className="font-mono text-xs text-[var(--gw-text)]">{result.device_id}</span>
                <span
                  className={`rounded-full px-2 py-0.5 text-xs font-semibold ${
                    result.status === 'success'
                      ? 'bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]'
                      : result.status === 'failed'
                        ? 'bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)]'
                        : 'bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)]'
                  }`}
                >
                  {resultStatusLabel(result.status)}
                </span>
              </div>
            </div>
          ))}
        </div>

        <div
          data-testid="activation-delivery-status"
          className="rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] px-4 py-3 text-xs text-[var(--gw-warning-fg)]"
        >
          <span className="font-semibold">
            {t('step4.activation_delivery_status_label')}:
          </span>{' '}
          {t('step4.activation_delivery_unconfirmed')}
        </div>
      </div>

      <div className="mt-8 space-y-3">
        {canContinue && onCommit && (
          <button
            type="button"
            onClick={handleCommit}
            className="w-full py-3 px-4 rounded-xl font-medium text-sm text-center bg-[var(--gw-primary)] hover:bg-[var(--gw-primary-hover)] text-[var(--gw-on-primary)] cursor-pointer transition-colors"
          >
            {t('step4.go_to_dashboard_btn')}
          </button>
        )}

        <button
          type="button"
          onClick={onReset}
          className="w-full py-3 px-4 rounded-xl border border-[var(--gw-control)] bg-[var(--gw-surface)] text-sm font-medium text-[var(--gw-text)] transition-all hover:bg-[var(--gw-surface)]"
        >
          {t('step4.reset_activation_btn')}
        </button>

        <p className="text-xs text-[var(--gw-muted)] text-center select-none">
          {canContinue
            ? t('step4.success_info')
            : t('step4.empty_info')}
        </p>
      </div>
    </div>
  );
}
