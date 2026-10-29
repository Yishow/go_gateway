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
      <Dialog.Overlay className="fixed inset-0 z-50 bg-black/50" />
      <Dialog.Content role="alertdialog" onCloseAutoFocus={(event) => { if (returnFocus?.isConnected) { event.preventDefault(); returnFocus.focus(); } }}
        className="fixed left-1/2 top-1/2 z-50 max-h-[85vh] w-[min(92vw,36rem)] -translate-x-1/2 -translate-y-1/2 space-y-4 overflow-y-auto rounded-xl border border-slate-600 bg-slate-900 p-5 text-sm text-slate-100" data-testid="group-delete-confirm">
        <Dialog.Title className="font-semibold">{t('remove.title')}</Dialog.Title>
        <Dialog.Description>{t('remove.keeps')}</Dialog.Description>
        <RecordingLifecycleHistory history={history} />
        <div className="flex gap-2">
          <button type="button" onClick={onConfirm} disabled={busy} className="rounded bg-red-700 px-3 py-2 text-white disabled:opacity-40" data-testid="group-delete-confirm-button">{t('remove.confirm')}</button>
          <button type="button" onClick={onClose} disabled={busy} className="rounded border border-slate-600 px-3 py-2" data-testid="group-delete-cancel">{t('cancel')}</button>
        </div>
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
}
