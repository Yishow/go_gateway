import { useEffect, useRef, useState } from 'react';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '@/utils/returnFocus';
import '@/styles/semantic-modal.css';
import { semanticTokens } from '../../../../tokens';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import type { WriteGroupLifecycleHistory } from '@/types/studioV2WriteGroupLifecycle';
import { RecordingLifecycleHistory } from './RecordingLifecycleHistory';
import { RecordingDeleteDialogError } from '../../../../recording-home/RecordingDeleteDialogError';
import '../../../../recording-home/recording-delete-dialog.css';

export function SavedRecordingRemoveDialog({ action = 'remove', open, busy, completed = false, refreshFailed = false, group, history, error, onStop, onOpenRecording, onRefresh, onClose, onConfirm, returnFocus, focusOwner, fallbackFocus, objectLabel }: {
  action?: 'stop' | 'remove'; open: boolean; busy: boolean; completed?: boolean; refreshFailed?: boolean; history?: WriteGroupLifecycleHistory; group?: { readonly id?: string; readonly revision?: string; readonly destination: { readonly table_schema: string; readonly table_name: string } }; error?: unknown; onStop?: () => void; onOpenRecording?: () => void; onRefresh?: () => void; onClose: () => void; onConfirm: () => void; returnFocus?: HTMLElement | null; focusOwner?: object; fallbackFocus?: HTMLElement | null; objectLabel?: string;
}) {
  const { t } = useTranslation('recording-lifecycle');
  const { t: deletion } = useTranslation('recording-home');
  const scope = JSON.stringify([group?.id, group?.revision]);
  const [acknowledgement, setAcknowledgement] = useState<{ scope: string; checked: boolean }>();
  const acknowledged = acknowledgement?.scope === scope && acknowledgement.checked;
  useEffect(() => { if (!open) setAcknowledgement(undefined); }, [open]);
  const testId = action === 'stop' ? 'group-stop' : 'group-delete';
  const ownScope = useRef({}); const currentOwner = useRef(focusOwner ?? ownScope.current); currentOwner.current = focusOwner ?? ownScope.current;
  const ticket = useRef<ReturnFocusTicket | null>(null);
  return <Dialog.Root open={open} onOpenChange={(next) => { if (!next && !busy) onClose(); }}>
    <Dialog.Portal>
      <Dialog.Overlay className="gw-modal-overlay" />
      <Dialog.Content role="alertdialog" onOpenAutoFocus={() => { ticket.current = captureReturnFocus(currentOwner.current, returnFocus ?? (document.activeElement instanceof HTMLElement ? document.activeElement : null)); }}
        onCloseAutoFocus={(event) => { event.preventDefault(); restoreReturnFocus(ticket.current, currentOwner.current, fallbackFocus); }}
        className="gw-modal gw-delete-dialog" data-size="confirm" style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }} data-testid={`${testId}-confirm`}>
        <header className="gw-modal-header"><Dialog.Title className="font-semibold">{action === 'stop' ? t('stop.title') : deletion('deleteDialogs.recording.title')}</Dialog.Title>{objectLabel && <p>{objectLabel}</p>}</header>
        <div className="gw-modal-body space-y-4 gw-history">{completed ? <Dialog.Description role="status">{deletion('deleteDialogs.recording.success')}</Dialog.Description> : action === 'stop' ? <Dialog.Description>{t('stop.keeps')}</Dialog.Description> : <>
          <Dialog.Description>{deletion('deleteDialogs.recording.deleted')}</Dialog.Description>
          <p>{deletion('deleteDialogs.recording.retained')} {group && <strong>{[group.destination.table_schema, group.destination.table_name].filter(Boolean).join('.')}</strong>}</p>
          <p>{deletion('deleteDialogs.recording.cannotUndo')}</p>
          {history && <RecordingLifecycleHistory history={history} deleting />}
          <RecordingDeleteDialogError error={error} onStop={onStop} onOpenRecording={onOpenRecording} onRefresh={onRefresh} />
          {refreshFailed && <p role="alert" className="gw-warning">{deletion('deleteDialogs.refreshFailed')}</p>}
          <label><input type="checkbox" checked={acknowledged} onChange={event => setAcknowledgement({ scope, checked: event.target.checked })} /> {deletion('deleteDialogs.recording.acknowledge')}</label>
        </>}</div>
        <footer className="gw-modal-footer">
          {!completed && <button type="button" onClick={onConfirm} disabled={busy || (action === 'remove' && (!acknowledged || Boolean(error)))} className="gw-danger" data-testid={`${testId}-confirm-button`}>{action === 'stop' ? t('stop.confirm') : deletion('deleteDialogs.recording.confirm')}</button>}
          <button type="button" onClick={onClose} disabled={busy} className="gw-secondary" data-testid={`${testId}-cancel`}>{completed ? deletion('deleteDialogs.recording.close') : t('cancel')}</button>
        </footer>
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
}
