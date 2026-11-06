import { useTranslation } from 'react-i18next';
import type { WriteGroupCalculationCapabilities } from '../../../../../../types/studioV2WriteGroupCalculation';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import { RecordingMixedEditor } from '../../../recording-editor/RecordingMixedEditor';
import { initializeMixedMembers } from '../../../recording-editor/recordingMemberCalculation';
import { newMixedDraft } from '../../../recording-editor/recordingMixedDraft';
import { RecordingInstantaneousMethod } from '../../../recording-editor/RecordingInstantaneousMethod';
import { isInstantaneousProfileDraft } from '../../../recording-editor/instantaneousReduction';
import { createBasicRecordingRequestId } from '../../../state/basicRecording';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import { memberKey, newDraft, type EditorDraft, type EditorMember } from '../../../state/writeGroup/draft';
import { GroupProfileSummary } from './GroupProfileSummary';

export interface GroupEditorMethodSectionProps {
  presentation: 'guided' | 'simple';
  draft: EditorDraft;
  saved: boolean;
  disabled: boolean;
  capabilities?: WriteGroupCalculationCapabilities;
  dialect: string;
  candidates?: GroupCandidate[];
  /** The saved record, compared with the draft to tell the operator a method change needs a fresh schema check. */
  savedGroup?: WriteGroup;
  /** The editor's last save failure, shown at the instantaneous method fields when it rejects the configuration. */
  saveError?: unknown;
  /** Replaces the whole draft; the editor keeps the save guard. */
  onReplaceDraft: (next: EditorDraft) => void;
  /** Opens a new unsaved recording prefilled from a legacy recording without mutating the saved group. */
  onStartNew?: (next: EditorDraft) => void;
  onHiddenMember?: (pointId: string) => void;
}

export function createNewRecordingFromLegacy(
  draft: EditorDraft,
  candidates: GroupCandidate[],
  capabilities: WriteGroupCalculationCapabilities | undefined,
  presentation: 'guided' | 'simple',
  copyName: string,
): EditorDraft {
  const destination = {
    connector_id: draft.connector_id,
    connector_revision: draft.connector_revision,
    table_schema: draft.table_schema,
    table_name: '',
  };
  const baseMembers: EditorMember[] = draft.members.map(member => ({
    key: memberKey(member.point_id, member.tag_id),
    device_id: member.device_id,
    point_id: member.point_id,
    tag_id: member.tag_id,
    entity_key: '',
    target_column: '',
    required: member.required,
    ...(member.source_value_type ? { source_value_type: member.source_value_type } : {}),
  }));
  if (presentation === 'simple') {
    const base = newMixedDraft(destination, createBasicRecordingRequestId());
    return {
      ...base,
      name: copyName,
      members: initializeMixedMembers(baseMembers, [], candidates, capabilities, {
        interval: 15,
        unit: 'second',
        alignment: 'utc_epoch',
      }),
    };
  }
  return { ...newDraft(destination), name: copyName, members: baseMembers };
}

/**
 * The "recording method" section of the record editor. It only assembles the
 * existing method views; the editor remains the single draft and save owner.
 */
export function GroupEditorMethodSection({ presentation, draft, saved, disabled, capabilities, dialect, candidates = [], savedGroup, saveError, onReplaceDraft, onStartNew, onHiddenMember }: GroupEditorMethodSectionProps) {
  const { t } = useTranslation('recording-editor-constraints');
  const { t: editorText } = useTranslation('recording-editor');
  // Legacy profiles keep their canonical group policy; mixed members own their individual cadence.
  const profile = (draft.layout_version === 'managed_long_v1' || draft.layout_version === 'managed_wide_profile_v1') && draft.recording_profile !== undefined;
  const mixed = draft.recording_profile?.kind === 'mixed';
  const legacyWide = draft.recording_profile === undefined;
  const editableInstantaneous = presentation === 'simple' && saved && isInstantaneousProfileDraft(draft);
  return <>
    {legacyWide && (
      <div className="space-y-2 text-xs text-[var(--gw-text)]" data-testid="recording-legacy-method-summary">
        <p data-testid="recording-legacy-method-text">{t('legacyMethod.summary', { seconds: draft.interval_seconds })}</p>
        <button
          type="button"
          className="gw-secondary rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-1.5 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-ground)] disabled:opacity-50"
          disabled={disabled}
          data-testid="recording-legacy-create-new"
          onClick={() => {
            if (disabled) return;
            const next = createNewRecordingFromLegacy(draft, candidates, capabilities, presentation, editorText('copy.name', { name: draft.name }));
            (onStartNew ?? onReplaceDraft)(next);
          }}
        >
          {t('legacyMethod.createWithPoints')}
        </button>
      </div>
    )}
    <RecordingMixedEditor section={presentation === 'simple' ? 'method' : 'all'} draft={draft} saved={saved} disabled={disabled} capabilities={capabilities} dialect={dialect} onChange={onReplaceDraft} onHiddenMember={onHiddenMember} />
    {profile && !mixed && <GroupProfileSummary draft={draft} editableInstantaneous={editableInstantaneous} />}
    {/* Only an instantaneous profile may change its reduction here; a cumulative record never offers mean. */}
    {editableInstantaneous && (
      <RecordingInstantaneousMethod draft={draft} savedGroup={savedGroup} candidates={candidates} disabled={disabled} saveError={saveError} onReplaceDraft={onReplaceDraft} />
    )}
  </>;
}
