import type { ReactNode } from 'react';
import { isValidRecordingProfile } from '../../../../utils/studioV2RecordingProfile';
import { useTranslation } from 'react-i18next';
import type { EditorDraft } from '../state/writeGroup/draft';
import { recordingMethodGroups } from './recordingMethodGrouping';

/** Read-only local intent summary. No readiness, schedule, SQL or running evidence is inferred. */
export function RecordingEditorSummary({ draft, destinationName, saved, dirty, diagnostics }: {
  draft: EditorDraft; destinationName?: string; saved: boolean; dirty: boolean; diagnostics?: ReactNode;
}) {
  const { t } = useTranslation('recording-editor');
  const { t: constraints } = useTranslation('recording-editor-constraints');
  const devices = new Set(draft.members.map(member => member.device_id)).size;
  const methods = recordingMethodGroups(draft, t);
  const profile = draft.recording_profile;
  const knownProfile = profile && profile.kind !== 'mixed' && isValidRecordingProfile(profile) ? profile : undefined;
  const knownInterval = knownProfile?.interval ?? (!profile && Number.isFinite(draft.interval_seconds) && Number(draft.interval_seconds) > 0 ? draft.interval_seconds : undefined);
  const wholeMethod = knownProfile ? t(knownProfile.kind === 'cumulative' ? 'simple.summary.cumulative' : `calculation.${knownProfile.sample_reduction ?? 'latest'}`) : !profile ? constraints('legacyMethod.label') : t('simple.summary.methodUnknown');
  // Always shown: the rail is the record's identity, not an optional disclosure.
  return <section data-testid="recording-editor-summary" aria-labelledby="recording-editor-summary-title">
    <h3 id="recording-editor-summary-title" className="recording-editor-summary-title">{t('simple.summary.title')}</h3>
    <p className="recording-editor-summary-state">{t(!saved || dirty ? 'compact.unsavedSettings' : 'compact.savedSettings')}</p>
    <dl>
      <div><dt>{t('simple.summary.name')}</dt><dd data-testid="recording-summary-name">{draft.name || t('simple.unset')}</dd></div>
      <div><dt>{t('simple.section.sources')}</dt><dd>{t('simple.summary.sources', { devices, points: draft.members.length })}</dd></div>
      <div><dt>{t(draft.recording_profile?.kind === 'mixed' ? 'simple.summary.results' : 'simple.summary.groupResults')}</dt><dd>{draft.recording_profile?.kind === 'mixed'
        ? methods.length ? <ul>{methods.map(value => <li key={value.key}>{value.method} · {value.interval} {value.unit} · {t('simple.summary.points', { count: value.count })}</li>)}</ul> : t('simple.summary.noPoints')
        : <>{wholeMethod} · {knownInterval === undefined ? t('simple.summary.intervalUnknown') : <>{knownInterval} {!profile ? constraints('legacyMethod.secondsUnit') : t(`calculation.${knownProfile?.unit ?? 'second'}`)}</>}</>}</dd></div>
      <div><dt>{t('simple.destination')}</dt><dd>{destinationName || t('simple.unset')}</dd></div>
      <div><dt>{t('simple.summary.table')}</dt><dd data-testid="recording-summary-table">{draft.table_schema && <>{draft.table_schema} · </>}{draft.table_name || t(destinationName && draft.storage_strategy === 'managed' ? 'simple.pendingColumn' : 'simple.unset')}</dd></div>
    </dl>
    {diagnostics}
  </section>;
}
