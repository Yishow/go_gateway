import { useTranslation } from 'react-i18next';
import type { EditorDraft } from '../../../state/writeGroup/draft';
import { destinationLock } from './GroupEditorCapabilities';

export interface RecordingDestinationLockNoticeProps {
  draft: EditorDraft;
  /** The recording has a persisted identity. */
  saved: boolean;
  /** Existing copy-as-new action; omitted when this recording cannot be copied. */
  onCopy?: () => void;
  copyDisabled?: boolean;
}

/**
 * Explains why a saved recording's destination is fixed. The reason follows the
 * same `destinationLock` result that disables the fields, so text and locks agree.
 */
export function RecordingDestinationLockNotice({ draft, saved, onCopy, copyDisabled = false }: RecordingDestinationLockNoticeProps) {
  const { t } = useTranslation('recording-editor-constraints');
  const { t: editorText } = useTranslation('recording-editor');
  if (!saved || !destinationLock(draft, saved).locked) return null;
  const reason = draft.layout_version === 'managed_long_v1' ? 'longTable' : 'boundTable';
  const secondaryBtn = 'gw-secondary rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-1.5 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-ground)] disabled:opacity-50';
  return <div className="recording-compact-members__metadata" data-testid="recording-destination-lock-notice">
    <p data-testid="recording-destination-lock-reason">{t(`destinationLock.${reason}`)}</p>
    {onCopy && <button type="button" className={secondaryBtn} disabled={copyDisabled} onClick={onCopy} data-testid="recording-destination-lock-copy">{editorText('copy.action')}</button>}
  </div>;
}
