import { useRecordingPointNameBuffer } from './useRecordingPointNameBuffer';
import { RecordingPointNameContent } from './RecordingPointNameContent';
import { useRecordingDraftGuard } from './RecordingDraftNavigation';
import { useRecordingRangeBuffer } from './useRecordingRangeBuffer';
import { RecordingRangeContent } from './RecordingRangeContent';
import { useRecordingPickerDiscovery, type RecordingPickerDiscovery } from './useRecordingPickerDiscovery';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { useRecordingDeviceBuffer } from './useRecordingDeviceBuffer';
import { RecordingSourceDeviceContent } from './RecordingSourceDeviceContent';
import { useEffect, useLayoutEffect, useMemo, useRef, useState, type ReactNode, type MutableRefObject } from 'react';
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
import { cancelPicker, comparePickerRows, createPicker, removeMember, sourceKey, type CandidateEvidence, type MemberHistory } from './recordingPickerState';
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
  renderMembers?: (actions: { onRemove: (member: EditorMember) => void; renderName?: (candidate: GroupCandidate) => ReactNode; nameLocked?: boolean; undo?: { member: EditorMember; onUndo: () => void; disabled: boolean } }) => ReactNode;
}
/** Controlled intent only. GroupEditor remains the sole recording draft/save owner.
 * Source content shares this one Radix shell; never mounts a second dialog. */
export function RecordingPointSelection({ scope, workspaceId, candidates: suppliedCandidates, members, onChange, disabled, devices: suppliedDevices, evidence, status = 'loaded', onRetry,
  deviceSource, discover, onDiscovery, renderMembers }: RecordingPointSelectionProps) {
  const { t } = useTranslation('recording-editor');
  const deviceBuffer = useRecordingDeviceBuffer(workspaceId, deviceSource, disabled);
  const nameBuffer = useRecordingPointNameBuffer(scope, deviceSource?.name, disabled);
  const rangeBuffer = useRecordingRangeBuffer(workspaceId, deviceSource?.range, disabled);
  const rangeSaved = useRef('');
  const rangeNavigation = useRef<{ requestId: string; sessionId: string; generation: number; scope: string; intent: number } | undefined>(undefined);
  const [locatedCount, setLocatedCount] = useState<number>();
  const discovery = useRecordingPickerDiscovery(scope, discover);
  const [confirmedCandidates, setConfirmedCandidates] = useState<{ rows: readonly GroupCandidate[]; afterRequest: number }>();
  const [renamedCandidates, setRenamedCandidates] = useState<{ scope: string; rows: GroupCandidate[]; afterRequest: number; supplied: readonly GroupCandidate[] }>();
  const candidates = useMemo(() => {
    const observed = discover ? discovery.result?.value?.candidates ?? [] : suppliedCandidates;
    return [...observed.map(row => renamedCandidates?.scope === scope && !(discover ? discovery.result?.status === 'loaded' && discovery.result.request > renamedCandidates.afterRequest : suppliedCandidates !== renamedCandidates.supplied) ? renamedCandidates.rows.find(named => sourceKey(named) === sourceKey(row)) ?? row : row), ...(confirmedCandidates && !(discover && discovery.result?.status === 'loaded' && discovery.result.request > confirmedCandidates.afterRequest) ? confirmedCandidates.rows.filter(row => !observed.some(item => sourceKey(item) === sourceKey(row))) : [])];
  }, [discover, discovery.result, suppliedCandidates, confirmedCandidates, renamedCandidates, scope]);
  const [confirmedDevices, setConfirmedDevices] = useState<{ device: StudioV2WorkspaceDeviceRecord; afterRequest: number }[]>([]);
  const devices = useMemo(() => {
    const observed = discover ? discovery.result?.value?.devices ?? [] : suppliedDevices;
    return [...observed, ...confirmedDevices.filter(item => !(discover && discovery.result?.status === 'loaded' && discovery.result.request > item.afterRequest) && !observed.some(row => row.id === item.device.id)).map(item => item.device)];
  }, [discover, discovery.result, suppliedDevices, confirmedDevices]);
  useEffect(() => {
    setConfirmedDevices(current => {
      const next = current.filter(item => discover
        ? !(discovery.result?.status === 'loaded' && discovery.result.request > item.afterRequest)
        : !suppliedDevices.some(device => device.id === item.device.id));
      return next.length === current.length ? current : next;
    });
  }, [discover, discovery.result, suppliedDevices]);
  useEffect(() => { if (discovery.result?.status === 'loaded' && discovery.result.value) onDiscovery?.({ scope, value: discovery.result.value }); }, [discovery.result, scope, onDiscovery]);
  useEffect(() => {
    if (confirmedCandidates && (discover ? discovery.result?.status === 'loaded' && discovery.result.request > confirmedCandidates.afterRequest : confirmedCandidates.rows.every(row => suppliedCandidates.some(observed => sourceKey(observed) === sourceKey(row))))) setConfirmedCandidates(undefined);
  }, [confirmedCandidates, discover, discovery.result, suppliedCandidates]);
  const pickerStatus = discover ? discovery.result?.status ?? 'loading' : status;
  const sourceIntent = useRef(0);
  const [picker, setPicker] = useState(() => createPicker(scope));
  const [history, setHistory] = useState<MemberHistory<EditorMember>>(() => ({ scope, removed: new Map() }));
  const [lastRemoved, setLastRemoved] = useState<EditorMember>();
  const [open, setOpen] = useState(false);
  const [view, setView] = useState<'picker' | 'device-choice' | 'device-buffer' | 'range-buffer'>('picker');
  const [choice, setChoice] = useState('');

  const [confirm, setConfirm] = useState(false);
  const [nameConfirm, setNameConfirm] = useState(false);
  const nameConfirmFocus = useRef<HTMLElement | null>(null), nameKeepButton = useRef<HTMLButtonElement>(null);
  const nameConfirmSnapshot = useRef<{ scope: string; key: string; name: string } | undefined>(undefined);
  const sourceClose = useRef<(() => void) | undefined>(undefined);
  const [deviceConfirm, setDeviceConfirm] = useState(false);
  const [deviceMounted, setDeviceMounted] = useState(false);
  const suppressNameBlur = useRef(false);
  // The selected list shares this name controller after the picker closes.
  useEffect(() => { if (!open) suppressNameBlur.current = false; }, [open]);
  const navigation = useRecordingDraftGuard({
    version: JSON.stringify([scope, picker.pending, deviceBuffer.buffer, rangeBuffer.buffer, nameBuffer.value?.name]),
    dirty: picker.pending.length > 0 || deviceBuffer.dirty || rangeBuffer.dirty || nameBuffer.dirty,
    pending: deviceBuffer.locked || rangeBuffer.locked || nameBuffer.pending, host: open,
    discard: () => { suppressNameBlur.current = true; nameConfirmSnapshot.current = undefined; setNameConfirm(false); sourceIntent.current += 1; deviceBuffer.reset(); rangeBuffer.reset(); nameBuffer.reset(); setPicker(cancelPicker(picker)); setOpen(false); },
  });
  useLayoutEffect(() => { if (nameConfirm && !navigation.confirming) nameKeepButton.current?.focus(); }, [nameConfirm, navigation.confirming]);
  const trigger = useRef<HTMLButtonElement>(null), ticket = useRef<ReturnFocusTicket | null>(null);
  const owner = useRef({ scope });
  const active = !disabled && picker.scope === scope && history.scope === scope;
  const projection = recordingPickerProjection(scope, candidates, members, evidence);
  for (const row of candidates) if (deviceSource?.range?.canSelect?.(row.device_id, row.point_id) === false) projection.evidence = { ...projection.evidence, [sourceKey(row)]: { ...projection.evidence[sourceKey(row)], availability: 'unknown' } };
  const nameIntent = useRef<GroupCandidate | undefined>(undefined);
  const nameFlight = useRef(false);
  const finishName = () => {
    const next = nameIntent.current; nameIntent.current = undefined;
    nameBuffer.reset();
    if (next && document.activeElement instanceof HTMLElement && document.activeElement.dataset.pickerFocus === `name:${sourceKey(next)}`) void nameBuffer.begin(next);
  };
  const saveName = () => {
    if (!active || nameFlight.current) return;
    nameFlight.current = true;
    void nameBuffer.save().then(saved => {
      if (!saved || !active) { nameIntent.current = undefined; return; }
      const rows = candidates.map(value => sourceKey(value) === sourceKey(saved) ? saved : value);
      setRenamedCandidates(current => ({ scope, afterRequest: discovery.currentRequest(), supplied: suppliedCandidates, rows: [...(current?.scope === scope ? current.rows.filter(value => sourceKey(value) !== sourceKey(saved)) : []), saved] }));
      onDiscovery?.({ scope, value: { candidates: rows, devices } }); finishName(); if (discover) void discovery.refresh();
    }).finally(() => { nameFlight.current = false; });
  };
  const back = () => { sourceIntent.current += 1; if (view === 'range-buffer' && rangeBuffer.dirty && !rangeBuffer.locked) setConfirm(true); else setView('picker'); };
  const keepName = () => {
    setNameConfirm(false); nameConfirmSnapshot.current = undefined; suppressNameBlur.current = false;
    requestAnimationFrame(() => { if (nameConfirmFocus.current?.isConnected) nameConfirmFocus.current.focus({ preventScroll: true }); });
  };
  const close = () => {
    if (nameBuffer.pending || nameBuffer.value?.outcome?.status === 'unknown') return;
    if (nameBuffer.dirty) {
      const value = nameBuffer.value!;
      nameConfirmSnapshot.current = { scope, key: sourceKey(value.candidate), name: value.name };
      nameConfirmFocus.current = document.activeElement instanceof HTMLElement ? document.activeElement : null;
      suppressNameBlur.current = true; setNameConfirm(true); return;
    }
    nameBuffer.reset();
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
  const targetAvailable = Boolean(rangeBuffer.target && devices.some(device => device.id === rangeBuffer.target?.id));
  const startRange = (device: { id: string; name: string }) => {
    if (!rangeBuffer.canStart(device.id)) { setChoice(device.id); setView('device-choice'); return; }
    setView('range-buffer'); void rangeBuffer.start(device);
  };
  const captureRangeNavigation = () => {
    const request = rangeBuffer.submission();
    if (request) rangeNavigation.current = { requestId: request.requestId, sessionId: request.sessionId, generation: request.generation, scope, intent: sourceIntent.current };
  };
  useEffect(() => {
    const result = rangeBuffer.outcome;
    if (result?.status !== 'confirmed' || !result.candidates || !result.rule || rangeSaved.current === result.request.requestId) return;
    rangeSaved.current = result.request.requestId;
    const rows = result.candidates.candidates;
    setConfirmedCandidates({ rows, afterRequest: discovery.currentRequest() });
    const union = [...candidates.filter(row => !rows.some(next => sourceKey(next) === sourceKey(row))), ...rows];
    onDiscovery?.({ scope, value: { candidates: union, devices } });
    const navigation = rangeNavigation.current;
    // Completion is a source fact, not authority to navigate a later user intent.
    if (!open || view !== 'range-buffer' || rangeBuffer.dirty || navigation?.scope !== scope || navigation.intent !== sourceIntent.current ||
      navigation.requestId !== result.request.requestId || navigation.sessionId !== result.request.sessionId || navigation.generation !== result.request.generation) return;
    setLocatedCount(rows.length);
    const added = new Set(rows.map(sourceKey));
    setPicker(current => {
      const sorted = union.filter(row => row.device_id === result.request.deviceId).sort((a, b) => comparePickerRows(a, b, current.sort));
      const first = sorted.findIndex(row => added.has(sourceKey(row)));
      return { ...current, query: '', pointType: 'all', deviceId: result.request.deviceId, page: Math.max(0, Math.floor(first / current.pageSize)), scrollTop: 0, focusKey: first >= 0 ? sourceKey(sorted[first]) : 'search' };
    });
    if (!rangeBuffer.dirty) setView('picker');
  }, [rangeBuffer.outcome, rangeBuffer.dirty, candidates, devices, discovery, onDiscovery, scope, open, view]);
  const source = (deviceId?: string) => {
    if (!active) return;
    const device = devices.find(item => item.id === deviceId);
    if (deviceId && !device) return;
    if (device) startRange(device);
    else { setChoice(''); setView('device-choice'); }
  };
  const renderName = deviceSource?.name ? (row: GroupCandidate) => <RecordingPointNameContent candidate={row} disabled={!active || (nameBuffer.value?.outcome?.status === 'unknown' || nameBuffer.dirty && !nameBuffer.value?.name.trim()) && sourceKey(nameBuffer.value!.candidate) !== sourceKey(row)}
    controller={nameBuffer} onBegin={() => {
      if (!active || nameBuffer.value?.outcome?.status === 'unknown') return;
      if (nameBuffer.pending || nameBuffer.dirty) { nameIntent.current = row; if (!nameBuffer.pending && nameBuffer.value?.name.trim()) saveName(); return; }
      nameBuffer.reset(); void nameBuffer.begin(row);
    }} onBack={finishName} onSave={saveName} canAutoSave={() => !suppressNameBlur.current && !navigation.confirming} /> : undefined;
  const nameLocked = Boolean(nameBuffer.pending || nameBuffer.dirty || nameBuffer.value?.outcome?.status === 'unknown');
  return <section className="recording-point-selection space-y-3" data-testid="recording-point-selection">
    <button ref={trigger} type="button" className="gw-secondary" data-testid="recording-points-open" disabled={!active}
      onClick={() => { suppressNameBlur.current = false; setView('picker'); setLocatedCount(undefined); setOpen(true); if (discover) void discovery.refresh(); }}>{t('picker.open')}</button>
    {deviceBuffer.locked && (!open || view !== 'device-buffer') && <p role="status">{t('device.unresolved')}</p>}
    {renderMembers?.({ onRemove: remove, renderName, nameLocked, undo: lastRemoved ? { member: lastRemoved, disabled: !active || !restored,
      onUndo: () => { if (active && restored) { onChange(restored); setLastRemoved(undefined); } } } : undefined })}
    <Dialog.Root open={open} onOpenChange={value => { if (!value) { if (navigation.confirming) navigation.keep(); else if (nameConfirm) keepName(); else close(); } }}>
      <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
        <Dialog.Content className="gw-modal" data-size={navigation.confirming || nameConfirm || confirm || deviceConfirm ? 'confirm' : 'form'} role={navigation.confirming || nameConfirm || confirm || deviceConfirm ? 'alertdialog' : 'dialog'}
          style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }}
          onOpenAutoFocus={() => { ticket.current = captureReturnFocus(owner.current, trigger.current); }}
          onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(ticket.current, owner.current, trigger.current); }}
          onPointerDownOutside={event => event.preventDefault()}
          onEscapeKeyDown={event => { event.preventDefault(); if (event.isComposing) return; if (navigation.confirming) { navigation.keep(); return; } if (nameConfirm) { keepName(); return; }
            const name = nameBuffer.value, focused = document.activeElement;
            if (view === 'picker' && name && focused instanceof HTMLElement && focused.dataset.pickerFocus === `name:${sourceKey(name.candidate)}` &&
              (name.outcome?.status === 'ready' || name.preparing && !name.outcome)) { nameIntent.current = undefined; nameBuffer.reset(); return; }
            if (view === 'device-buffer') { sourceClose.current?.(); return; } if (confirm) setConfirm(false); else close(); }}>
          {navigation.confirmation ?? <header className="gw-modal-header"><Dialog.Title>{t(nameConfirm ? 'name.cancelTitle' : confirm || deviceConfirm ? 'device.discardTitle' : view === 'picker' ? 'picker.title' : view === 'range-buffer' ? 'range.title' : 'source.title')}</Dialog.Title>
            <Dialog.Description>{t(nameConfirm ? 'name.cancelDescription' : confirm || deviceConfirm ? 'device.discardDescription' : view === 'range-buffer' ? 'range.description' : 'picker.description')}</Dialog.Description></header>}
          {nameConfirm && !navigation.confirming && <><div className="gw-modal-body" /><footer className="gw-modal-footer">
            <button ref={nameKeepButton} type="button" className="gw-secondary" onClick={keepName}>{t('name.cancelKeep')}</button>
            <button type="button" className="gw-primary" onClick={() => {
              const value = nameBuffer.value, snapshot = nameConfirmSnapshot.current;
              if (!active || !value || !nameBuffer.dirty || nameBuffer.pending || snapshot?.scope !== scope || snapshot.key !== sourceKey(value.candidate) || snapshot.name !== value.name) { keepName(); return; }
              nameConfirmSnapshot.current = undefined; nameIntent.current = undefined; nameBuffer.reset(); sourceIntent.current += 1; setPicker(cancelPicker(picker)); setNameConfirm(false); setOpen(false);
            }}>{t('name.cancelDiscard')}</button>
          </footer></>}
          <div style={{ display: navigation.confirming || nameConfirm ? 'none' : 'contents' }}>
          {confirm ? <><div className="gw-modal-body" /><footer className="gw-modal-footer">
            <button type="button" className="gw-secondary" onClick={() => setConfirm(false)}>{t('device.keep')}</button>
            <button type="button" className="gw-primary" onClick={() => { setConfirm(false); rangeBuffer.reset(); setView('picker'); }}>{t('device.discard')}</button>
          </footer></> : <>
            {view === 'picker' && <div className="gw-modal-body recording-picker-panel">{locatedCount !== undefined && <p role="status" data-testid="recording-range-located">{t('range.located', { count: locatedCount })}</p>}{rangeBuffer.recoverable && rangeBuffer.canResume() && <button type="button" data-testid="recording-range-resume" onClick={() => { if (rangeBuffer.resume()) setView('range-buffer'); }}>{t('range.resume')}</button>}{rangeBuffer.inactiveSources.map(item => <button key={item.request.requestId} type="button" data-testid={`recording-range-recover-${item.request.requestId}`} disabled={!rangeBuffer.canResume(item.request.requestId)} onClick={() => { sourceIntent.current += 1; if (rangeBuffer.resume(item.request.requestId)) setView('range-buffer'); }}>{t('range.resume')} · {devices.find(device => device.id === item.request.deviceId)?.name ?? item.request.deviceId} · {item.request.buffer.name} · {item.request.buffer.start_address}</button>)}{deviceBuffer.recoverable && <button type="button" data-testid="recording-source-resume" onClick={() => { if (deviceBuffer.resume()) { setDeviceMounted(true); setView('device-buffer'); } }}>{t('device.inspect')}</button>}<RecordingPointPicker state={picker} projection={projection} devices={devices}
              status={active ? pickerStatus : 'failed'} labels={recordingPickerLabels(t)} onChange={setPicker} onCancel={close} onRetry={discover ? () => { void discovery.refresh(); } : onRetry ?? (() => undefined)}
              nameLocked={nameLocked} nameNotice={deviceSource?.name ? t('name.shared') : undefined} renderName={renderName}
              onAdd={identities => { const next = addRecordingMembers(members, history, identities, projection); if (active && next) { onChange(next); setPicker(cancelPicker(picker)); setOpen(false); } }} onSource={source} /></div>}
            {view === 'device-choice' && <><div className="gw-modal-body space-y-4">
              {choice && !rangeBuffer.canStart(choice) && <><p role="status" data-testid="recording-range-other-pending">{t('range.otherPending')}</p>
                <button type="button" data-testid="recording-range-resume-choice" onClick={() => { if (rangeBuffer.resume()) setView('range-buffer'); }}>{t('range.resume')}</button></>}
              <label>{t('source.chooseDevice')}<select value={choice} onChange={event => setChoice(event.target.value)}><option value="">{t('source.noDeviceSelected')}</option>
                {devices.map(device => <option key={device.id} value={device.id}>{device.name}</option>)}</select></label>
              <button type="button" disabled={!deviceSource || !active || deviceBuffer.preparing} onClick={() => { const intent = ++sourceIntent.current; void deviceBuffer.start().then(permitted => { if (permitted && sourceIntent.current === intent) { setDeviceMounted(true); setView('device-buffer'); } }); }}>{t('source.newDevice')}</button>
              {deviceBuffer.preparing && <p role="status">{t('device.checkingScope')}</p>}
              {deviceBuffer.blocked && <p role="alert">{t('device.scopeBlocked')}</p>}
            </div><footer className="gw-modal-footer"><button type="button" onClick={back}>{t('source.back')}</button>
              <button type="button" className="gw-primary" disabled={!devices.some(device => device.id === choice) || !rangeBuffer.canStart(choice)} onClick={() => { const device = devices.find(item => item.id === choice); if (device) startRange(device); }}>{t('source.continue')}</button></footer></>}
            {view === 'range-buffer' && <RecordingRangeContent controller={rangeBuffer} available={targetAvailable && active} onBack={back}
              onSubmit={() => { void rangeBuffer.submit(); captureRangeNavigation(); }} onInspect={() => { captureRangeNavigation(); void rangeBuffer.inspect(); }} onPrepare={() => { void rangeBuffer.prepareSaved(); captureRangeNavigation(); }} />}
          </>}
          {deviceMounted && <div className="recording-source-panel" hidden={view !== 'device-buffer' || confirm}>{<RecordingSourceDeviceContent controller={deviceBuffer} slot={{ active: open && view === 'device-buffer' && !confirm, closeRef: sourceClose, onConfirmChange: setDeviceConfirm,
            onReturn: () => { setDeviceConfirm(false); setView('picker'); }, onSaved: device => { setConfirmedDevices(current => [...current.filter(item => item.device.id !== device.id), { device, afterRequest: discovery.currentRequest() }]); startRange(device); } }} />}</div>}
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </section>;
}
