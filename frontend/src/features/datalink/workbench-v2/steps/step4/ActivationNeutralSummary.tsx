import { useTranslation } from 'react-i18next';

export function ActivationNeutralSummary({ onReset }: { onReset: () => void }) {
  const { t } = useTranslation('workbench-v2');
  return (
    <div
      data-testid="activation-empty-message"
      className="rounded-2xl border border-amber-500/20 bg-amber-500/5 px-5 py-6 text-sm text-amber-100"
    >
      <p>{t('step4.activation_no_success')}</p>
      <button type="button" className="mt-4 rounded-xl border border-slate-700 px-4 py-2" onClick={onReset}>
        {t('step4.reset_activation_btn')}
      </button>
    </div>
  );
}
