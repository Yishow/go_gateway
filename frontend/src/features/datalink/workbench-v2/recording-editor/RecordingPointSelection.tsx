import { useRecordingPickerDiscovery, type RecordingPickerDiscovery } from './useRecordingPickerDiscovery';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { useRecordingDeviceBuffer } from './useRecordingDeviceBuffer';
import { RecordingSourceDeviceContent } from './RecordingSourceDeviceContent';
import { useEffect, useRef, useState, type ReactNode, type MutableRefObject } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import type { EditorMember } from '../state/writeGroup/draft';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '../../../../utils/returnFocus';
import { semanticTokens } from '../tokens';
import { RecordingPointPicker } from './RecordingPointPicker';
import { recordingPickerLabels } from './recordingPickerLabels';
import { addRecordingMembers, recordingPickerProjection } from './recordingMemberSelection';
import { cancelPicker, createPicker, removeMember, sourceKey, type CandidateEvidence, type MemberHistory } from './recordingPickerState';
import '../../../../styles/semantic-modal.css';

export interface RecordingSourceSlot {
  active: boolean;
  closeRef: MutableRefObject<(() => void) | undefined>;
  onConfirmChange: (confirming: boolean) => void;
  onReturn: () => void;
  onSaved: (device: StudioV2WorkspaceDeviceRecord) => void;
}
export interface RecordingPointSelectionProps {
  scope: string; workspaceId: string; candidates: readonly GroupCandidate[]; members: EditorMember[];
  onChange: (members: EditorMember[]) => void; disabled: boolean;
  devices: readonly { id: string; name: string }[];
  evidence?: Readonly<Record<string, CandidateEvidence | undefined>>;
  status?: 'loaded' | 'loading' | 'failed'; onRetry?: () => void;
  deviceSource?: DeviceSourcePort;
  discover?: () => Promise<RecordingPickerDiscovery>;
  onDiscovery?: (observed: { scope: string; value: RecordingPickerDiscovery }) => void;
  renderMembers?: (actions: { onRemove: (member: EditorMember) => void; undo?: { member: EditorMember; onUndo: () => void; disabled: boolean } }) => ReactNode;
}
/** Controlled intent only. GroupEditor remains the sole recording draft/save owner.
 * Source content shares this one Radix shell; never mounts a second dialog. */
export function RecordingPointSelection({ scope, workspaceId, candidates: suppliedCandidates, members, onChange, disabled, devices: suppliedDevices, evidence, status = 'loaded', onRetry,
  deviceSource, discover, onDiscovery, renderMembers }: RecordingPointSelectionProps) {
  const { t } = useTranslation('recording-editor');
  const deviceBuffer = useRecordingDeviceBuffer(workspaceId, deviceSource, disabled);
  const discovery = useRecordingPickerDiscovery(scope, discover);
  const candidates = discover ? discovery.result?.value?.candidates ?? [] : suppliedCandidates;
  const [confirmedDevices, setConfirmedDevices] = useState<{ device: StudioV2WorkspaceDeviceRecord; afterRequest: number }[]>([]);
  const observedDevices = discover ? discovery.result?.value?.devices ?? [] : suppliedDevices;
  const devices = [...observedDevices, ...confirmedDevices.filter(item => !(discover && discovery.result?.status === 'loaded' && discovery.result.request > item.afterRequest) && !observedDevices.some(observed => observed.id === item.device.id)).map(item => item.device)];
  useEffect(() => {
    setConfirmedDevices(current => {
      const next = current.filter(item => discover
        ? !(discovery.result?.status === 'loaded' && discovery.result.request > item.afterRequest)
        : !suppliedDevices.some(device => device.id === item.device.id));
      return next.length === current.length ? current : next;
    });
  }, [discover, discovery.result, suppliedDevices]);
  useEffect(() => { if (discovery.result?.status === 'loaded' && discovery.result.value) onDiscovery?.({ scope, value: discovery.result.value }); }, [discovery.result, scope, onDiscovery]);
  const pickerStatus = discover ? discovery.result?.status ?? 'loading' : status;
  const sourceIntent = useRef(0);
  const [picker, setPicker] = useState(() => createPicker(scope));
  const [history, setHistory] = useState<MemberHistory<EditorMember>>(() => ({ scope, removed: new Map() }));
  const [lastRemoved, setLastRemoved] = useState<EditorMember>();
  const [open, setOpen] = useState(false);
  const [view, setView] = useState<'picker' | 'device-choice' | 'device-buffer' | 'range-buffer'>('picker');
  const [choice, setChoice] = useState('');
  const [target, setTarget] = useState<{ id: string; name: string }>();
  const [range, setRange] = useState({ name: '', address: '', count: '' });
  const [rangeDirty, setRangeDirty] = useState(false);
  const [confirm, setConfirm] = useState(false);
  const sourceClose = useRef<(() => void) | undefined>(undefined);
  const [deviceConfirm, setDeviceConfirm] = useState(false);
  const [deviceMounted, setDeviceMounted] = useState(false);
  const trigger = useRef<HTMLButtonElement>(null), ticket = useRef<ReturnFocusTicket | null>(null);
  const owner = useRef({ scope });
  const active = !disabled && picker.scope === scope && history.scope === scope;
  const projection = recordingPickerProjection(scope, candidates, members, evidence);
  const back = () => { sourceIntent.current += 1; if (view === 'range-buffer' && rangeDirty) setConfirm(true); else setView('picker'); };
  const close = () => {
    sourceIntent.current += 1;
    if (view === 'device-buffer') { sourceClose.current?.(); return; }
    if (view !== 'picker') { back(); return; }
    setPicker(cancelPicker(picker)); setOpen(false);
  };
  const remove = (member: EditorMember) => {
    if (!active) return;
    const next = removeMember(members, history, sourceKey(member));
    setHistory(next.history); setLastRemoved(member); onChange([...next.members]);
  };
  const restored = lastRemoved ? addRecordingMembers(members, history, [lastRemoved], projection) : undefined;
  const targetAvailable = Boolean(target && devices.some(device => device.id === target.id));
  const startRange = (device: { id: string; name: string }) => {
    setTarget(device); setRange({ name: '', address: '', count: '' }); setRangeDirty(false); setView('range-buffer');
  };
  const source = (deviceId?: string) => {
    if (!active) return;
    const device = devices.find(item => item.id === deviceId);
    if (deviceId && !device) return;
    if (device) startRange(device);
    else { setChoice(''); setView('device-choice'); }
  };
  return <section className="space-y-3" data-testid="recording-point-selection">
    <button ref={trigger} type="button" className="gw-secondary" data-testid="recording-points-open" disabled={!active}
      onClick={() => { setView('picker'); setOpen(true); if (discover) void discovery.refresh(); }}>{t('picker.open')}</button>
    {deviceBuffer.locked && (!open || view !== 'device-buffer') && <p role="status">{t('device.unresolved')}</p>}
    {renderMembers?.({ onRemove: remove, undo: lastRemoved ? { member: lastRemoved, disabled: !active || !restored,
      onUndo: () => { if (active && restored) { onChange(restored); setLastRemoved(undefined); } } } : undefined })}
    <Dialog.Root open={open} onOpenChange={value => { if (!value) close(); }}>
      <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
        <Dialog.Content className="gw-modal" data-size={confirm || deviceConfirm ? 'confirm' : 'form'} role={confirm || deviceConfirm ? 'alertdialog' : 'dialog'}
          style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }}
          onOpenAutoFocus={() => { ticket.current = captureReturnFocus(owner.current, trigger.current); }}
          onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(ticket.current, owner.current, trigger.current); }}
          onPointerDownOutside={event => event.preventDefault()}
          onEscapeKeyDown={event => { event.preventDefault(); if (view === 'device-buffer') { sourceClose.current?.(); return; } if (confirm) setConfirm(false); else close(); }}>
          <header className="gw-modal-header"><Dialog.Title>{t(confirm || deviceConfirm ? 'device.discardTitle' : view === 'picker' ? 'picker.title' : view === 'range-buffer' ? 'range.title' : 'source.title')}</Dialog.Title>
            <Dialog.Description>{t(confirm || deviceConfirm ? 'device.discardDescription' : 'picker.description')}</Dialog.Description></header>
          {confirm ? <><div className="gw-modal-body" /><footer className="gw-modal-footer">
            <button type="button" className="gw-secondary" onClick={() => setConfirm(false)}>{t('device.keep')}</button>
            <button type="button" className="gw-primary" onClick={() => { setConfirm(false); setRangeDirty(false); setRange({ name: '', address: '', count: '' }); setView('picker'); }}>{t('device.discard')}</button>
          </footer></> : <>
            {view === 'picker' && <div className="gw-modal-body recording-picker-panel">{deviceBuffer.recoverable && <button type="button" data-testid="recording-source-resume" onClick={() => { if (deviceBuffer.resume()) { setDeviceMounted(true); setView('device-buffer'); } }}>{t('device.inspect')}</button>}<RecordingPointPicker state={picker} projection={projection} devices={devices}
              status={active ? pickerStatus : 'failed'} labels={recordingPickerLabels(t)} onChange={setPicker} onCancel={() => setOpen(false)} onRetry={discover ? () => { void discovery.refresh(); } : onRetry ?? (() => undefined)}
              onAdd={identities => { const next = addRecordingMembers(members, history, identities, projection); if (active && next) { onChange(next); setPicker(cancelPicker(picker)); setOpen(false); } }} onSource={source} /></div>}
            {view === 'device-choice' && <><div className="gw-modal-body space-y-4">
              <label>{t('source.chooseDevice')}<select value={choice} onChange={event => setChoice(event.target.value)}><option value="">{t('source.noDeviceSelected')}</option>
                {devices.map(device => <option key={device.id} value={device.id}>{device.name}</option>)}</select></label>
              <button type="button" disabled={!deviceSource || !active || deviceBuffer.preparing} onClick={() => { const intent = ++sourceIntent.current; void deviceBuffer.start().then(permitted => { if (permitted && sourceIntent.current === intent) { setDeviceMounted(true); setView('device-buffer'); } }); }}>{t('source.newDevice')}</button>
              {deviceBuffer.preparing && <p role="status">{t('device.checkingScope')}</p>}
              {deviceBuffer.blocked && <p role="alert">{t('device.scopeBlocked')}</p>}
            </div><footer className="gw-modal-footer"><button type="button" onClick={back}>{t('source.back')}</button>
              <button type="button" className="gw-primary" disabled={!devices.some(device => device.id === choice)} onClick={() => { const device = devices.find(item => item.id === choice); if (device) startRange(device); }}>{t('source.continue')}</button></footer></>}
            {view === 'range-buffer' && <><div className="gw-modal-body space-y-4"><p data-testid="recording-range-device">{target?.name}</p>
              {targetAvailable ? <p className="gw-muted" data-testid="recording-source-confirmed">{t('range.savedDeviceNote')}</p> : <p role="alert">{t('picker.unavailable')}</p>}
              {(['name', 'address', 'count'] as const).map(field => <label key={field} className="block">{t(`range.${field}`)}<input disabled={!targetAvailable || !active} value={range[field]}
                inputMode={field === 'count' ? 'numeric' : undefined} onChange={event => { setRange({ ...range, [field]: event.target.value }); setRangeDirty(true); }} /></label>)}
              <p role="status">{t('range.persistenceUnavailable')}</p>
            </div><footer className="gw-modal-footer"><button type="button" data-testid="recording-range-back" onClick={back}>{t('source.back')}</button>
              <button type="button" className="gw-primary" disabled>{t('range.save')}</button></footer></>}
          </>}
          {deviceMounted && <div className="recording-source-panel" hidden={view !== 'device-buffer' || confirm}>{<RecordingSourceDeviceContent controller={deviceBuffer} slot={{ active: open && view === 'device-buffer' && !confirm, closeRef: sourceClose, onConfirmChange: setDeviceConfirm,
            onReturn: () => { setDeviceConfirm(false); setView('picker'); }, onSaved: device => { setConfirmedDevices(current => [...current.filter(item => item.device.id !== device.id), { device, afterRequest: discovery.currentRequest() }]); startRange(device); } }} />}</div>}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </section>;
}
