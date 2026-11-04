import { useTranslation } from 'react-i18next';
/** Selected intent determines which risks are immediately relevant; all explanations remain reachable. */
export function RecordingIntentHelp({ visibleSafety = false, compact = false, hasMembers = true, hasCounter = true, hasRequired = true }: {
  visibleSafety?: boolean; compact?: boolean; hasMembers?: boolean; hasCounter?: boolean; hasRequired?: boolean;
}) {
  const { t } = useTranslation('recording-editor');
  // Compact editor: only the data-loss risks that apply stay visible; the method explanations
  // live in the method section's help button (record.method).
  if (compact) return visibleSafety && (hasCounter || hasRequired) ? <div className="space-y-1 text-sm text-[var(--gw-muted)]" data-testid="recording-intent-risks">
    {hasCounter && <p>{t('calculation.firstBaseline')}</p>}
    {hasRequired && <p>{t('simple.requiredRisk')}</p>}
  </div> : null;
  return <div className="space-y-2 text-sm text-[var(--gw-muted)]">
    {visibleSafety && <p>{t(hasMembers ? 'calculation.dueUnion' : 'simple.methodEmpty')}</p>}
    {visibleSafety && hasCounter && <p>{t('calculation.firstBaseline')}</p>}
    {visibleSafety && hasRequired && <p>{t('simple.requiredRisk')}</p>}
    <details data-testid="recording-intent-help">
      <summary className="cursor-pointer">{t('calculation.help')}</summary>
      <dl className="mt-2 space-y-2">
        {(['latest', 'mean', 'delta'] as const).map(method => <div key={method}><dt className="font-medium text-[var(--gw-text)]">{t(`calculation.${method}`)}</dt><dd>{t(`calculation.${method}Example`)}</dd></div>)}
      </dl>
      <p>{t('calculation.mixedExample')}</p>
      {(!visibleSafety || !hasCounter) && <p>{t('calculation.firstBaseline')}</p>}
      <p>{t('calculation.counterPolicy')}</p><p>{t('calculation.requiredPolicy')}</p>
    </details>
  </div>;
}
