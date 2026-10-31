import { useEffect, useMemo, useRef, useState } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import type { DeviceSourceBuffer, DeviceSourceOutcome, DeviceSourcePort, DeviceSourceRequest } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '../../../../utils/returnFocus';
import { semanticTokens } from '../tokens';
import { newDeviceSourceBuffer, parsedDeviceSourceBuffer } from './deviceSourceBuffer';
import { RecordingDeviceSourceForm } from './RecordingDeviceSourceForm';
import '../../../../styles/semantic-modal.css';

/** One local form, backed only by the already-mounted workspace mutation owner. */
export function RecordingDeviceSourceSetup({ workspaceId, port, disabled }: {
  workspaceId: string; port?: DeviceSourcePort; disabled: boolean;
}) {
  const { t } = useTranslation('recording-editor');
  const [preparing, setPreparing] = useState(false);
  const [scopeBlocked, setScopeBlocked] = useState(false);
  const prepareFlight = useRef(false);
  const [open, setOpen] = useState(false);
  const [buffer, setBuffer] = useState(newDeviceSourceBuffer);
  const [generation, setGeneration] = useState(0);
  const [dirty, setDirty] = useState(false);
  const [confirm, setConfirm] = useState<'close' | DeviceSourceBuffer>();
  const [outcome, setOutcome] = useState<DeviceSourceOutcome>();
  const [pending, setPending] = useState(false);
  const [checking, setChecking] = useState(false);
  const [savedDeviceId, setSavedDeviceId] = useState<string>();
  const [sessionId] = useState(() => crypto.randomUUID());
  const owner = useMemo(() => ({ workspaceId, sessionId }), [workspaceId, sessionId]);
  const ownerRef = useRef(owner); ownerRef.current = owner;
  const generationRef = useRef(generation); generationRef.current = generation;
  const submitting = useRef(false);
  const requestRef = useRef<DeviceSourceRequest | undefined>(undefined);
  const mounted = useRef(true);
  const focusTicket = useRef<ReturnFocusTicket | null>(null);
  const trigger = useRef<HTMLButtonElement>(null);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const locked = pending || outcome?.status === 'unknown' || outcome?.status === 'pending';
  const parsed = parsedDeviceSourceBuffer(buffer);
  const change = (next: DeviceSourceBuffer) => { setBuffer(next); setGeneration(value => value + 1); generationRef.current += 1; setDirty(true); };
  const accept = (result: DeviceSourceOutcome | undefined, requestOwner: typeof owner) => {
    const request = requestRef.current;
    if (!mounted.current || ownerRef.current !== requestOwner || !request || !result ||
      result.request.workspaceId !== request.workspaceId || result.request.sessionId !== request.sessionId ||
      result.request.requestId !== request.requestId || result.request.generation !== request.generation) return;
    const verified = result.status !== 'confirmed' || Boolean(result.device?.id);
    setOutcome(verified ? result : { ...result, status: 'unknown' });
    setPending(false); submitting.current = false;
    if (result.device?.id && (result.status === 'confirmed' || result.status === 'persisted-unavailable')) setSavedDeviceId(result.device.id);
    if (verified && result.status === 'confirmed' && generationRef.current === request.generation) setDirty(false);
  };
  const submit = async () => {
    if (!port || disabled || locked || !parsed || !dirty || pending || submitting.current) return;
    const request: DeviceSourceRequest = { workspaceId, sessionId, requestId: crypto.randomUUID(), generation,
      ...(savedDeviceId ? { deviceId: savedDeviceId } : {}), buffer: parsed };
    requestRef.current = request; submitting.current = true; setPending(true); setOutcome(undefined);
    const requestOwner = owner;
    try { accept(await port.submit(request), requestOwner); }
    catch { accept({ request, status: 'unknown' }, requestOwner); }
  };
  const inspect = async () => {
    const request = requestRef.current;
    if (!port || !request || checking) return;
    setChecking(true); const requestOwner = owner;
    try { accept(await port.inspect(request.requestId), requestOwner); }
    catch { accept({ request, status: 'unknown' }, requestOwner); }
    finally { if (mounted.current && ownerRef.current === requestOwner) setChecking(false); }
  };
  const close = () => {
    if (locked) { setOpen(false); return; } // Close the view, never cancel a submitted operation.
    if (dirty) { setConfirm('close'); return; }
    setOpen(false);
  };
  const discard = () => {
    if (confirm === 'close') { setOpen(false); setBuffer(newDeviceSourceBuffer()); setDirty(false); setGeneration(0); setOutcome(undefined); setSavedDeviceId(undefined); }
    else if (confirm) change(confirm);
    setConfirm(undefined);
  };
  const unresolved = (port?.observations() ?? []).filter(item => item.request.workspaceId === workspaceId &&
    ['pending', 'unknown', 'persisted-unavailable'].includes(item.status));
  const resume = (item: DeviceSourceOutcome) => {
    requestRef.current = item.request; setBuffer(item.request.buffer); setGeneration(item.request.generation);
    generationRef.current = item.request.generation; setDirty(false); setOutcome(item); setPending(false);
    setSavedDeviceId(item.device?.id); setOpen(true);
  };
  const openSource = async () => {
    if (prepareFlight.current || disabled || !port) return;
    if (locked) { setOpen(true); return; }
    prepareFlight.current = true; setPreparing(true); setScopeBlocked(false);
    const capturedOwner = owner;
    let permitted: boolean;
    try { permitted = port.prepare ? await port.prepare(workspaceId) : true; } catch { permitted = false; }
    prepareFlight.current = false;
    if (!mounted.current || ownerRef.current !== capturedOwner) return;
    setPreparing(false);
    if (!permitted) { setScopeBlocked(true); return; }
    if (outcome?.status === 'confirmed' && !dirty) {
      setBuffer(newDeviceSourceBuffer()); setGeneration(0); generationRef.current = 0;
      setOutcome(undefined); setSavedDeviceId(undefined); requestRef.current = undefined;
    }
    setOpen(true);
  };
  const newer = outcome?.status === 'confirmed' && dirty;
  return <div className="space-y-2">
    <button ref={trigger} type="button" data-testid="recording-source-open" disabled={disabled || preparing || !port || unresolved.some(item => item.request.requestId !== requestRef.current?.requestId)} onClick={() => { void openSource(); }}>{t(locked ? 'device.resume' : 'device.open')}</button>
    {preparing && <p role="status">{t('device.checkingScope')}</p>}
    {scopeBlocked && <p role="alert">{t('device.scopeBlocked')}</p>}
    {disabled && <p className="gw-muted">{t('device.blocked')}</p>}
    {!open && unresolved.filter(item => item.request.requestId !== requestRef.current?.requestId).map(item =>
      <button key={item.request.requestId} type="button" data-testid={`recording-source-recover-${item.request.requestId}`}
        disabled={disabled} onClick={() => resume(item)}>{t('device.recover', { name: item.request.buffer.name })}</button>)}
    {!open && locked && <p role="status">{t('device.unresolved')}</p>}
    <Dialog.Root open={open} onOpenChange={value => { if (!value) close(); }}>
      <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
        <Dialog.Content className="gw-modal" data-size={confirm ? 'confirm' : 'form'} role={confirm ? 'alertdialog' : 'dialog'}
          style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }}
          onOpenAutoFocus={() => { focusTicket.current = captureReturnFocus(owner, trigger.current); }}
          onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(focusTicket.current, ownerRef.current, trigger.current); }}
          onPointerDownOutside={event => event.preventDefault()}
          onEscapeKeyDown={event => { event.preventDefault(); if (confirm) setConfirm(undefined); else close(); }}>
          <header className="gw-modal-header"><Dialog.Title>{t(confirm ? 'device.discardTitle' : 'device.title')}</Dialog.Title>
            <Dialog.Description>{t(confirm ? 'device.discardDescription' : 'device.description')}</Dialog.Description></header>
          <div className="gw-modal-body">
            {!confirm && <>
              {outcome?.status === 'confirmed' && !dirty ? <p role="status" data-testid="recording-source-confirmed">{t('device.confirmed', { name: outcome.device?.name })}</p> :
                <RecordingDeviceSourceForm buffer={buffer} lockedProtocol={Boolean(savedDeviceId) || locked}
                  onChange={change} onProtocol={next => { if (dirty) setConfirm(next); else change(next); }} />}
              {newer && <p role="status" data-testid="recording-source-newer">{t('device.newer')}</p>}
              {pending && <p role="status">{t('device.pending')}</p>}
              {outcome && outcome.status !== 'confirmed' && outcome.status !== 'pending' && <p role="alert" data-testid={`recording-source-${outcome.status}`}>{t(`device.${outcome.status}`)}</p>}
              {dirty && !parsed && <p role="status">{t('device.invalid')}</p>}
            </>}
          </div>
          <footer className="gw-modal-footer">
            {confirm ? <><button type="button" className="gw-secondary" onClick={() => setConfirm(undefined)}>{t('device.keep')}</button>
              <button type="button" className="gw-primary" data-testid="recording-source-discard" onClick={discard}>{t('device.discard')}</button></> : <>
              <button type="button" className="gw-secondary" data-testid="recording-source-close" onClick={close}>{t(locked || outcome?.status === 'confirmed' ? 'device.close' : 'device.cancel')}</button>
              {locked && !pending && <button type="button" data-testid="recording-source-inspect" disabled={checking} onClick={() => { void inspect(); }}>{t('device.inspect')}</button>}
              {(!outcome || outcome.status !== 'confirmed' || dirty) && <button type="button" className="gw-primary" data-testid="recording-source-save" disabled={disabled || !port || locked || !parsed || !dirty} onClick={() => { void submit(); }}>{t(savedDeviceId ? 'device.saveChanges' : 'device.save')}</button>}
            </>}
          </footer>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </div>;
}
