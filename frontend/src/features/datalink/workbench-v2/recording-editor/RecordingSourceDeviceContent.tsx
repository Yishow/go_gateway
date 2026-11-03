import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import type { RecordingSourceSlot } from './RecordingPointSelection';
import type { RecordingDeviceBufferController } from './useRecordingDeviceBuffer';
import { RecordingDeviceSourceForm } from './RecordingDeviceSourceForm';
/** Content of the host's sole dialog, never a dialog or mutation owner itself. */
export function RecordingSourceDeviceContent({ controller: device, slot, existingEdit = false }: { controller: RecordingDeviceBufferController; slot: RecordingSourceSlot; existingEdit?: boolean }) {
  const { t } = useTranslation('recording-editor');
  const { onConfirmChange } = slot;
  const [confirm, setConfirm] = useState<'back' | 'clear' | DeviceSourceBuffer>();
  const close = () => {
    if (confirm) { setConfirm(undefined); return; }
    if (device.locked) return;
    if (device.dirty) setConfirm('back'); else slot.onReturn();
  };
  useEffect(() => { slot.closeRef.current = close; return () => { slot.closeRef.current = undefined; }; });
  useEffect(() => { onConfirmChange(Boolean(confirm)); }, [confirm, onConfirmChange]);
  useEffect(() => {
    if (!confirm && slot.active && device.outcome?.status === 'confirmed' && !device.dirty && device.outcome.device) slot.onSaved(device.outcome.device);
  }, [slot, device.outcome, device.dirty, confirm]);
  const discard = () => {
    if (confirm === 'back') { device.reset(); slot.onReturn(); }
    else if (confirm === 'clear') device.clearRecent();
    else if (confirm) device.selectConnection(confirm);
    setConfirm(undefined);
  };
  return <>
    <div className="gw-modal-body space-y-4">
      {!confirm && <>
        <RecordingDeviceSourceForm buffer={device.buffer} lockedProtocol={Boolean(device.deviceId) || device.locked} onChange={device.change}
          onSubtype={!device.deviceId ? next => { if (device.dirty) setConfirm(next); else device.selectConnection(next); } : undefined}
          onProtocol={next => { if (device.dirty) setConfirm(next); else device.selectConnection(next); }} />
        {device.reused && !device.deviceId && <div data-testid="recording-source-reuse" className="space-y-2">
          <p className="gw-muted">{t('recentConnection.reused')}</p>
          <p className="gw-muted">{t('recentConnection.verifyEndpoint')}</p>
          <button type="button" className="gw-secondary" data-testid="recording-source-clear-recent"
            onClick={() => { if (device.connectionEdited) setConfirm('clear'); else device.clearRecent(); }}>{t('recentConnection.clear')}</button>
        </div>}
        {device.pending && <p role="status">{t('device.pending')}</p>}
        {device.outcome?.status === 'confirmed' && device.dirty && <p role="status" data-testid="recording-source-newer">{t('device.newer')}</p>}
        {device.outcome && device.outcome.status !== 'confirmed' && device.outcome.status !== 'pending' && <p role="alert" data-testid={`recording-source-${device.outcome.status}`}>{t(`device.${device.outcome.status}`)}</p>}
        {device.dirty && !device.parsed && <p role="status">{t('device.invalid')}</p>}
      </>}
    </div>
    <footer className="gw-modal-footer" data-edit-actions={existingEdit || undefined} style={existingEdit ? { justifyContent: 'flex-end' } : undefined}>
      {confirm ? <><button type="button" className="gw-secondary" onClick={() => setConfirm(undefined)}>{t('device.keep')}</button>
        <button type="button" className="gw-primary" data-testid="recording-source-discard" onClick={discard}>{t('device.discard')}</button></> : <>
        <button type="button" className="gw-secondary" data-testid="recording-source-close" disabled={device.locked} onClick={close}>{t(existingEdit ? 'device.cancelEdit' : device.locked ? 'device.close' : 'source.back')}</button>
        {device.locked && !device.pending && <button type="button" data-testid="recording-source-inspect" disabled={device.checking} onClick={() => { void device.inspect(); }}>{t('device.inspect')}</button>}
        <button type="button" className="gw-primary" data-testid="recording-source-save" disabled={device.disabled || device.locked || !device.parsed || !device.dirty}
          onClick={() => { void device.submit(); }}>{t(existingEdit ? 'device.saveAndReturn' : device.deviceId ? 'device.saveChanges' : 'device.save')}</button>
      </>}
    </footer>
  </>;
}
