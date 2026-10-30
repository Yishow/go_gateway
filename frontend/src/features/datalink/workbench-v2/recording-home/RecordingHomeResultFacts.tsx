import { useTranslation } from 'react-i18next';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
/** Compact result facts; secondary timestamps/owner/readback details are opt-in. */
export function RecordingHomeResultFacts({ id, summary }: { id: string; summary: ReturnType<typeof recordingHomeResultSummary> }) {
  const { t } = useTranslation('recording-home');
  return <div data-testid={`home-result-${id}`} className="space-y-1 text-xs">
    <p>{t(`result.sql.${summary.sql}`)}</p>
    <p>{t(`result.current.${summary.current}`)}</p>
    {summary.configuration === 'not_requested' && <p>{t('result.configuration.not_requested')}</p>}
    {summary.effective === 'requested_not_effective' && <p>{t('result.effective.requested_not_effective')}</p>}
    <details><summary>{t('result.details')}</summary><div className="space-y-1 py-2">
      <p>{t(`result.configuration.${summary.configuration}`)}</p>
      {summary.lastCommittedAt && <p>{t('result.lastCommit')} <time>{summary.lastCommittedAt}</time>{summary.historicalDestination && ` · ${summary.historicalDestination}`}</p>}
      <p>{t(`result.effective.${summary.effective}`)}{summary.effectiveRevision && ` · ${summary.effectiveRevision}`}</p>
      {summary.requestedBoundary && <p>{t('result.requestedBoundary')} <time>{summary.requestedBoundary}</time></p>}
      <p>{t(`result.owners.${summary.owners}`)}</p>
      <p>{t(summary.readback === 'test_write_verified' ? 'result.testWriteReadback' : 'result.readbackUnknown')}</p>
      {summary.cleanup === 'test_write_completed' && <p>{t('result.testWriteCleanup')}</p>}
      {summary.observedAt && <p>{t('result.observedAt')} <time>{summary.observedAt}</time></p>}
    </div></details>
  </div>;
}
