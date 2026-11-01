import { useTranslation } from 'react-i18next';
/** Explanations describe stored outcomes; neither this disclosure nor its examples starts collection. */
export function RecordingIntentHelp() {
  const { t } = useTranslation('recording-editor');
  return <details className="text-xs text-slate-400" data-testid="recording-intent-help">
    <summary>{t('calculation.help')}</summary>
    <dl className="mt-2 space-y-2">
      {(['latest', 'mean', 'delta'] as const).map(method => <div key={method}><dt className="font-medium text-slate-200">{t(`calculation.${method}`)}</dt><dd>{t(`calculation.${method}Example`)}</dd></div>)}
    </dl>
    <p>{t('calculation.mixedExample')}</p>
    <p>{t('calculation.firstBaseline')}</p><p>{t('calculation.counterPolicy')}</p><p>{t('calculation.requiredPolicy')}</p>
  </details>;
}
