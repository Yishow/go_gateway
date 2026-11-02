import { useTranslation } from 'react-i18next';
import type { EditorDraft } from '../state/writeGroup/draft';
import type { CalculationEditorMember } from './recordingMemberCalculation';

/** Read-only local intent summary. No readiness, schedule, SQL or running evidence is inferred. */
export function RecordingEditorSummary({ draft, destinationName, saved, dirty }: {
  draft: EditorDraft; destinationName?: string; saved: boolean; dirty: boolean;
}) {
  const { t } = useTranslation('recording-editor');
  const devices = new Set(draft.members.map(member => member.device_id)).size;
  const methods = new Map<string, { method: string; interval: string; unit: string; count: number }>();
  const members: readonly CalculationEditorMember[] = draft.members;
  for (const member of members) {
    const method = member.calculation ? t(`calculation.${member.calculation}`) : t('simple.unset');
    const interval = member.cadence_input ?? (Number.isFinite(member.cadence?.interval) ? String(member.cadence?.interval) : t('simple.unset'));
    const unit = member.cadence?.unit ? t(`calculation.${member.cadence.unit}`) : '';
    const key = JSON.stringify([method, interval, unit]);
    methods.set(key, { method, interval, unit, count: (methods.get(key)?.count ?? 0) + 1 });
  }
  return <details open data-testid="recording-editor-summary">
    <summary>{t('simple.summary.title')}</summary>
    <p className="recording-editor-summary-state">{t(!saved || dirty ? 'compact.unsavedSettings' : 'compact.savedSettings')}</p>
    <dl>
      <div><dt>{t('simple.section.sources')}</dt><dd>{t('simple.summary.sources', { devices, points: draft.members.length })}</dd></div>
      <div><dt>{t('simple.summary.results')}</dt><dd>{draft.recording_profile?.kind === 'mixed'
        ? methods.size ? <ul>{[...methods].map(([key, value]) => <li key={key}>{value.method} · {value.interval} {value.unit} · {t('simple.summary.points', { count: value.count })}</li>)}</ul> : t('simple.summary.noPoints')
        : t('simple.summary.legacy')}</dd></div>
      <div><dt>{t('simple.destination')}</dt><dd>{destinationName || t('simple.unset')}</dd></div>
    </dl>
    <details className="recording-editor-summary-identities"><summary>{t('compact.identities')}</summary>
      <p>{t(!saved || dirty ? 'simple.summary.localIntent' : 'simple.summary.savedIntent')}</p>
      <dl><div><dt>{t('simple.summary.name')}</dt><dd data-testid="recording-summary-name">{draft.name || t('simple.unset')}</dd></div>
        <div><dt>{t('simple.destination')}</dt><dd>{draft.table_schema && <>{draft.table_schema} · </>}{draft.table_name || t(destinationName && draft.storage_strategy === 'managed' ? 'simple.pendingColumn' : 'simple.unset')}</dd></div></dl>
    </details>
  </details>;
}
