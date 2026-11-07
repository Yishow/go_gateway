import { useRef, useState } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import { useDeleteStudioV2WorkspaceDeviceMutation } from '@/hooks/datalink/useStudioV2WorkspaceDevices';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';
import { RecordingDeleteDialogError } from './RecordingDeleteDialogError';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '@/utils/returnFocus';
import '@/styles/semantic-modal.css';
import './recording-delete-dialog.css';

interface DeviceRecording {
  readonly id: string;
  readonly name: string;
  readonly status: string;
  readonly basic_managed_device_id?: string;
  readonly members: readonly { readonly device_id: string; readonly point_id: string }[];
}
export interface DeviceDeleteDialogProps {
  open: boolean;
  device: { readonly id: string; readonly name: string };
  workspaceId: string;
  recordings: readonly DeviceRecording[];
  scopeValid?: boolean;
  onClose: () => void;
  onDeleted: () => void;
  onRefresh: () => void;
  refreshFailed?: boolean;
}

/** An open session owns its typed confirmation and cannot reuse it for another device. */
export function DeviceDeleteDialog(props: DeviceDeleteDialogProps) {
  return props.open ? <DeviceDeleteSession key={props.device.id} {...props} /> : null;
}

function DeviceDeleteSession({ device, workspaceId, recordings, onClose, onDeleted, onRefresh, refreshFailed, scopeValid = true }: DeviceDeleteDialogProps) {
  const { t } = useTranslation('recording-home');
  const mutation = useDeleteStudioV2WorkspaceDeviceMutation(workspaceId);
  const [name, setName] = useState('');
  const [error, setError] = useState<unknown>(null);
  const inFlight = useRef(false);
  const owner = useRef({});
  const ticket = useRef<ReturnFocusTicket | null>(null);
  const local = recordings.filter(recording => recording.status !== 'deleted' &&
    (recording.basic_managed_device_id === device.id || recording.members.some(member => member.device_id === device.id)))
    .map(recording => ({ id: recording.id, name: recording.name,
      points: new Set(recording.members.filter(member => member.device_id === device.id).map(member => member.point_id)).size }));
  const references = error instanceof DeleteResponseError && error.code === 'DEVICE_REFERENCED_BY_RECORDINGS' && error.details &&
    'recordings' in error.details ? error.details : undefined;
  const recordingBlock = error instanceof DeleteResponseError && error.outcome === 'failed' && error.code?.startsWith('WRITE_GROUP_DELETE_');
  const deviceErrorMessage = error instanceof DeleteResponseError && error.outcome === 'deleted'
    ? error.code === 'DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED' ? 'shareRefreshFailed' : 'connectionCloseFailed'
    : error instanceof DeleteResponseError && error.outcome === 'failed' ? 'refused' : 'unknown';
  const blocked = references ? references.total > 0 : local.length > 0;
  const submit = async () => {
    if (inFlight.current || !scopeValid || blocked || name !== device.name || error) return;
    inFlight.current = true;
    try { await mutation.mutateAsync(device.id); onDeleted(); }
    catch (cause) { setError(cause); }
    finally { inFlight.current = false; }
  };
  return <Dialog.Root open onOpenChange={next => { if (!next && !mutation.isPending) onClose(); }}>
    <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
      <Dialog.Content role="alertdialog" className="gw-modal gw-delete-dialog" data-size="confirm"
        onOpenAutoFocus={() => { ticket.current = captureReturnFocus(owner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null); }}
        onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(ticket.current, owner.current); }}>
        <header className="gw-modal-header"><Dialog.Title>{t('deleteDialogs.device.title')}</Dialog.Title><p>{device.name}</p></header>
        <div className="gw-modal-body space-y-4">
          <Dialog.Description>{t('deleteDialogs.device.description')}</Dialog.Description>
          {!scopeValid && <div className="gw-warning" role="alert"><p>{t('deleteDialogs.device.stale')}</p><button type="button" onClick={onRefresh}>{t('deleteDialogs.device.refresh')}</button></div>}
          {blocked && <div className="gw-warning" role="alert"><p>{t('deleteDialogs.device.referenced')}</p>
            <ul>{(references?.recordings ?? local).map(recording => <li key={recording.id}>
              <span>{recording.name}</span>{'points' in recording && typeof recording.points === 'number' && <small> · {t('deleteDialogs.device.points')}: {recording.points}</small>}
            </li>)}</ul>
            {references && references.total > references.recordings.length && <p>{t('deleteDialogs.device.total')} {references.total}</p>}
          </div>}
          {recordingBlock && <><p>{t('deleteDialogs.device.blockedByRecording')}</p><RecordingDeleteDialogError error={error} onRefresh={onRefresh} /></>}
          {Boolean(error) && !references && !recordingBlock && <div role="alert" className="gw-warning"><p>{t(`deleteDialogs.device.${deviceErrorMessage}`)}</p>
            <button type="button" onClick={onRefresh}>{t('deleteDialogs.device.refresh')}</button></div>}
          {refreshFailed && <p role="alert" className="gw-warning">{t('deleteDialogs.refreshFailed')}</p>}
          {!blocked && <label>{t('deleteDialogs.device.typeName')}<input value={name} onChange={event => setName(event.target.value)} autoComplete="off" /></label>}
        </div>
        <footer className="gw-modal-footer"><button type="button" className="gw-danger" data-testid="device-delete-confirm-button"
          disabled={!scopeValid || blocked || name !== device.name || mutation.isPending || Boolean(error)} onClick={() => { void submit(); }}>{t('deleteDialogs.device.confirm')}</button>
          <button type="button" className="gw-secondary" disabled={mutation.isPending} onClick={onClose}>{t('deleteDialogs.device.cancel')}</button></footer>
      </Dialog.Content></Dialog.Portal>
  </Dialog.Root>;
}
