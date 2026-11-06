import { createBasicRecordingRequestId } from '../state/basicRecording';
import { ParameterHelp } from '../components/ParameterHelp';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import { useWriteGroupDeliveryQuery } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { RecordingMemberResult } from './RecordingMemberResult';
import { useTranslation } from 'react-i18next';
import type { WriteGroupCalculationCapabilities } from '../../../../types/studioV2WriteGroupCalculation';
import type { EditorDraft } from '../state/writeGroup/draft';
import { mixedDraftIssues, newMixedDraft } from './recordingMixedDraft';
import { cadenceInputInvalid } from './recordingMemberCalculation';
import { RecordingIntentHelp } from './RecordingIntentHelp';
import { RecordingMethodGroups } from './RecordingMethodGroups';
export function mixedEditorAvailable(capabilities: WriteGroupCalculationCapabilities | undefined, dialect: string, layout?: string): boolean {
  return Boolean(capabilities?.editor_supported && capabilities.draft_supported &&
    capabilities.schema_dialects.some(value => value === dialect) && (!layout || capabilities.schema_layouts.some(value => value === layout)));
}
export function recordingMixedEditorIssues(draft: EditorDraft, capabilities: WriteGroupCalculationCapabilities | undefined, dialect: string, destinationSaved: boolean): string[] {
  const issues = capabilities?.editor_supported ? mixedDraftIssues(draft, capabilities, dialect) : ['mixed-capability-unavailable'];
  if (!destinationSaved) issues.push('destination-unsaved');
  if (draft.members.some(cadenceInputInvalid)) issues.push('member-cadence-invalid');
  return [...new Set(issues)];
}
export interface RecordingMixedEditorProps { section?: 'all' | 'method' | 'destination'; draft: EditorDraft; saved: boolean; disabled: boolean; capabilities?: WriteGroupCalculationCapabilities; dialect: string; onChange: (draft: EditorDraft) => void; onHiddenMember?: (pointId: string) => void }
/** Intent-only controls. The enclosing GroupEditor retains its single CAS/save/schema authority. */
export function RecordingMixedEditor({ section = 'all', draft, saved, disabled, capabilities, dialect, onChange, onHiddenMember }: RecordingMixedEditorProps) {
  const { t } = useTranslation('recording-editor');
  const mixed = draft.recording_profile?.kind === 'mixed';
  const available = mixedEditorAvailable(capabilities, dialect);
  return <section className="space-y-2 text-xs text-[var(--gw-muted)]" data-testid={section === 'destination' ? 'recording-mixed-destination' : 'recording-mixed-intent'}>
    {section !== 'destination' && !mixed && (saved ? <p>{t('calculation.legacyUnchanged')}</p> : <>
      <button type="button" className="inline-flex min-h-10 max-sm:min-h-11 items-center justify-center rounded-lg border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)] disabled:cursor-not-allowed disabled:opacity-50" disabled={disabled || !available || draft.members.length > 0} onClick={() => {
        if (disabled || !available || saved || draft.members.length) return;
        onChange({ ...newMixedDraft({ connector_id: draft.connector_id, connector_revision: draft.connector_revision, table_schema: draft.table_schema, table_name: '' }, createBasicRecordingRequestId()), name: draft.name });
      }}>{t('calculation.chooseMixed')}</button>
      <p>{t(!available ? 'calculation.capability' : draft.members.length ? 'calculation.chooseBeforePoints' : 'calculation.independent')}</p>
    </>)}
    {mixed && <>
      {section === 'all' && <p>{t('calculation.independent')}</p>}
      {section !== 'method' && <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block">{t('calculation.layoutLabel')}
        <select className="mt-1 block rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] p-1" value={draft.layout_version} disabled={disabled || saved} onChange={event => {
          const layout = event.target.value;
          if (!saved && !disabled && (layout === 'managed_long_v1' || layout === 'managed_wide_profile_v1') && capabilities?.schema_layouts.includes(layout)) onChange({ ...draft, layout_version: layout });
        }}>{(['managed_wide_profile_v1', 'managed_long_v1'] as const).map(layout => <option key={layout} value={layout} disabled={!capabilities?.schema_layouts.includes(layout)}>{t(`calculation.${layout}`)}</option>)}</select>
      </label><ParameterHelp help="record.layout" /></div>}
    </>}
    {section === 'method' && mixed && <RecordingMethodGroups draft={draft} onHiddenMember={onHiddenMember} />}
    {section !== 'destination' && <RecordingIntentHelp compact={section === 'method'} visibleSafety={section === 'method' && mixed} hasMembers={draft.members.length > 0} hasCounter={draft.members.some(member => member.calculation === 'delta')} hasRequired={draft.members.some(member => member.required)} />}
    {section === 'method' && mixed && !available && <p role="status">{t('calculation.capability')}</p>}
  </section>;
}

/** A record that was never applied has no result windows to show; the block appears once it is. */
export function RecordingMixedResults({ group, candidates }: { group: WriteGroup; candidates: GroupCandidate[] }) {
  return group.applied_revision ? <AppliedMixedResults group={group} candidates={candidates} /> : null;
}
function AppliedMixedResults({ group, candidates }: { group: WriteGroup; candidates: GroupCandidate[] }) {
  const query = useWriteGroupDeliveryQuery(group.id, true, 5000);
  return <RecordingMemberResult group={group} candidates={candidates} delivery={query.isError ? undefined : query.data} />;
}
