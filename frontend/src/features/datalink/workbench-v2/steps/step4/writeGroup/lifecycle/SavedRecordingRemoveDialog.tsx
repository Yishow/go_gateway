import '@/styles/semantic-modal.css';
import { semanticTokens } from '../../../../tokens';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import type { WriteGroupLifecycleHistory } from '@/types/studioV2WriteGroupLifecycle';
import { RecordingLifecycleHistory } from './RecordingLifecycleHistory';

export function SavedRecordingRemoveDialog({ open, busy, history, onClose, onConfirm, returnFocus }: {
  open: boolean; busy: boolean; history?: WriteGroupLifecycleHistory; onClose: () => void; onConfirm: () => void; returnFocus?: HTMLElement | null;
}) {
  const { t } = useTranslation('recording-lifecycle');
  return <Dialog.Root open={open} onOpenChange={(next) => { if (!next && !busy) onClose(); }}>
    <Dialog.Portal>
      <Dialog.Overlay className="gw-modal-overlay" />
      <Dialog.Content role="alertdialog" onCloseAutoFocus={(event) => { if (returnFocus?.isConnected) { event.preventDefault(); returnFocus.focus(); } }}
        className="gw-modal" data-size="confirm" style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }} data-testid="group-delete-confirm">
        <header className="gw-modal-header"><Dialog.Title className="font-semibold">{t('remove.title')}</Dialog.Title></header>
        <div className="gw-modal-body space-y-4 gw-history"><Dialog.Description>{t('remove.keeps')}</Dialog.Description>
        <RecordingLifecycleHistory history={history} /></div>
        <footer className="gw-modal-footer">
          <button type="button" onClick={onConfirm} disabled={busy} className="gw-danger" data-testid="group-delete-confirm-button">{t('remove.confirm')}</button>
          <button type="button" onClick={onClose} disabled={busy} className="gw-secondary" data-testid="group-delete-cancel">{t('cancel')}</button>
        </footer>
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
}
