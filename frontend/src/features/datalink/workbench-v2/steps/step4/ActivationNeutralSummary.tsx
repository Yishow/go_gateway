import { useTranslation } from 'react-i18next';

export function ActivationNeutralSummary({
  onReset,
  canContinue,
  onCommit,
}: {
  onReset: () => void;
  canContinue?: boolean;
  onCommit?: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  return (
    <div
      data-testid="activation-empty-message"
      className="rounded-2xl border border-amber-500/20 bg-amber-500/5 px-5 py-6 text-sm text-amber-100"
    >
      <p>{t('step4.activation_no_success')}</p>
      <div className="mt-4 flex flex-wrap items-center gap-3">
        {canContinue && onCommit && (
          <button
            type="button"
            className="rounded-xl bg-emerald-600 px-4 py-2 font-medium text-white transition-all hover:bg-emerald-500 active:scale-[0.98]"
            onClick={onCommit}
          >
            {t('step4.go_to_dashboard_btn', '前往 Runtime Dashboard')}
          </button>
        )}
        <button
          type="button"
          className="rounded-xl border border-slate-700 px-4 py-2 hover:bg-slate-800/60"
          onClick={onReset}
        >
          {t('step4.reset_activation_btn')}
        </button>
      </div>
    </div>
  );
}
