import { useTranslation } from 'react-i18next';
/** Failed reads retain only the mounted view; retry uses the original query owner. */
export function LegacyRecordingReadError({ retry }: { retry: () => Promise<unknown> }) {
  const { t } = useTranslation('workbench-v2');
  return <div role="alert" data-testid="circuit-profile-groups-error" className="space-y-2 rounded border border-rose-800/70 p-4 text-xs text-rose-700 dark:text-rose-200">
    <p>{t('step4.basic.circuit.groups_error')}</p>
    <button type="button" className="gw-secondary" onClick={() => { void retry(); }}>{t('step4.group.retry')}</button>
  </div>;
}
