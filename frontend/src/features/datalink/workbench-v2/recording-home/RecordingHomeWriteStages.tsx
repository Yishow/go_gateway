import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import { AlertTriangle, Check, CircleHelp } from 'lucide-react';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { recordingHomeWriteStageMapping } from './recordingHomeWriteStageMapping';
import './recordingHomeRows.css';

type Summary = ReturnType<typeof recordingHomeResultSummary>;
const stages = ['configuration', 'activation', 'sql', 'readback'] as const;

export function RecordingHomeWriteStages({ name, summary }: { name: string; summary?: Summary }) {
  const { t } = useTranslation('recording-home');
  const id = useId();
  const states = summary ? recordingHomeWriteStageMapping(summary) : { configuration: 'unconfirmed', activation: 'unconfirmed', sql: 'unconfirmed', readback: 'unconfirmed' } as const;
  return <div role="group" className="gw-write-stages" aria-label={t('rowEvidence.stagesFor', { name })}>
    {stages.map(stage => {
      const state = states[stage], Icon = state === 'complete' ? Check : state === 'attention' ? AlertTriangle : CircleHelp;
      const status = stage === 'configuration' && summary?.configuration === 'not_requested' ? 'rowEvidence.draftNotRequested' :
        stage === 'sql' && (summary?.sql === 'historical_commit' || summary?.sql === 'previous_revision_receipt') && state !== 'attention' ? 'rowEvidence.historical' :
          stage === 'readback' && summary?.readback === 'test_write_verified' ? 'rowEvidence.testWriteReadback' : `rowEvidence.state.${state}`;
      const tooltip = `${id}-${stage}`;
      return <div key={stage} tabIndex={0} className="gw-write-stage" data-write-stage={stage} data-stage-state={state} aria-describedby={tooltip}>
        <span className="gw-write-stage-title"><Icon aria-hidden="true" size={14} />{t(`rowEvidence.stage.${stage}`)}</span>
        <span>{t(status)}</span>
        <div role="tooltip" id={tooltip} className="gw-write-stage-tooltip">
          {!summary && <p>{t('rowEvidence.evidenceUnavailable')}</p>}
          {summary && stage === 'configuration' && <><p>{t(`result.configuration.${summary.configuration}`)}</p>{summary.requestedBoundary && <p>{t('result.requestedBoundary')} <time>{summary.requestedBoundary}</time></p>}</>}
          {summary && stage === 'activation' && <><p>{t(`result.effective.${summary.effective}`)}{summary.effectiveRevision && ` · ${summary.effectiveRevision}`}</p>{summary.observedAt && <p>{t('result.observedAt')} <time>{summary.observedAt}</time></p>}</>}
          {summary && stage === 'sql' && <>{summary.lastCommittedAt && <p>{t('result.lastCommit')} <time>{summary.lastCommittedAt}</time>{summary.historicalDestination && ` · ${summary.historicalDestination}`}</p>}<p>{t(`result.owners.${summary.owners}`)}</p></>}
          {summary && stage === 'readback' && <><p>{t(summary.readback === 'test_write_verified' ? 'result.testWriteReadback' : 'result.readbackUnknown')}</p><p>{t(summary.cleanup === 'test_write_completed' ? 'result.testWriteCleanup' : 'rowEvidence.cleanupUnknown')}</p></>}
        </div>
      </div>;
    })}
  </div>;
}
