import { useTranslation } from 'react-i18next';
export function RuntimeLogsLink() {
  const { t } = useTranslation('runtime-logs');
  return <a href="/studio/logs" className="inline-flex rounded-lg border border-slate-600 px-3 py-2 text-xs font-medium text-cyan-200 hover:bg-slate-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-cyan-300">{t('title')}</a>;
}
