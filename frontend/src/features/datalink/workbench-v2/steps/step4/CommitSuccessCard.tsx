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
    card: 'bg-amber-950/20 border-amber-500/30 shadow-[0_0_24px_rgba(245,158,11,0.05)]',
    icon: 'border-amber-500/30 bg-amber-500/10 text-amber-300',
    title: 'text-amber-200',
    subtitle: 'text-amber-300/80',
  } : {
    card: 'bg-emerald-950/20 border-emerald-500/30 shadow-[0_0_24px_rgba(16,185,129,0.05)]',
    icon: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400',
    title: 'text-emerald-300',
    subtitle: 'text-emerald-400/80',
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
      className={`${cardTone.card} border rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm`}
    >
      <div className="space-y-6 flex-1">
        {/* 結果 Icon */}
        <div className="space-y-2 text-center">
          <div className={`mx-auto flex h-16 w-16 items-center justify-center rounded-full border text-3xl shadow-[0_0_16px_rgba(16,185,129,0.15)] ${cardTone.icon}`}>
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
            className="rounded-xl border border-amber-500/20 bg-amber-500/10 px-4 py-3 text-sm text-amber-100"
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
              className="rounded-xl border border-slate-800 bg-slate-950/40 px-4 py-3"
            >
              <div className="flex items-center justify-between gap-3">
                <span className="font-mono text-xs text-slate-200">{result.device_id}</span>
                <span
                  className={`rounded-full px-2 py-0.5 text-[11px] font-semibold ${
                    result.status === 'success'
                      ? 'bg-emerald-500/10 text-emerald-300'
                      : result.status === 'failed'
                        ? 'bg-rose-500/10 text-rose-300'
                        : 'bg-amber-500/10 text-amber-300'
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
          className="rounded-xl border border-amber-500/20 bg-amber-500/5 px-4 py-3 text-xs text-amber-100"
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
            className="w-full py-3 px-4 rounded-xl font-medium text-sm text-center bg-emerald-600 hover:bg-emerald-500 text-white cursor-pointer transition-all shadow-lg shadow-emerald-600/10 hover:shadow-emerald-500/20 active:scale-[0.98]"
          >
            {t('step4.go_to_dashboard_btn')}
          </button>
        )}

        <button
          type="button"
          onClick={onReset}
          className="w-full py-3 px-4 rounded-xl border border-slate-700 bg-slate-900/40 text-sm font-medium text-slate-200 transition-all hover:bg-slate-900/70"
        >
          {t('step4.reset_activation_btn')}
        </button>

        <p className="text-[10px] text-gray-500 text-center select-none">
          {canContinue
            ? t('step4.success_info')
            : t('step4.empty_info')}
        </p>
      </div>
    </div>
  );
}
