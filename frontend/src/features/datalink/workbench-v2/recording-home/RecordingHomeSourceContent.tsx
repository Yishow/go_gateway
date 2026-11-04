import { useEffect, useLayoutEffect, useRef, useState, type Dispatch, type ReactNode } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { WorkbenchV2State, Device } from '../state/types';
import type { PersistedMCRuleState } from '../state/mcPersistedPoints';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { useRecordingDeviceBuffer } from '../recording-editor/useRecordingDeviceBuffer';
import { RecordingSourceDeviceContent } from '../recording-editor/RecordingSourceDeviceContent';
import { useRecordingRangeBuffer } from '../recording-editor/useRecordingRangeBuffer';
import { RecordingRangeContent } from '../recording-editor/RecordingRangeContent';
import { useRecordingDraftGuard } from '../recording-editor/RecordingDraftNavigation';
import { RecordingHomeSavedRangeEditor, type HomeSourceLeave } from './RecordingHomeSavedRangeEditor';
import { readHomeDeviceRanges } from './recordingHomeSourceList';
import type { HomeSourceControls } from './RecordingHomeSourceHost';
import { RecordingHomeMemoryView, useHomeRangeInputs, type HomeMemoryBlock } from './RecordingHomeMemoryView';
import { describeMemoryRule } from './recordingHomeMemoryMap';
export interface RecordingHomeSourceContentProps extends HomeSourceControls {
  shareEnabled?: boolean; onRetryShare: () => void; onOpenShareSettings: () => void;
  state: WorkbenchV2State; dispatch: Dispatch<WorkbenchV2Action>; source: DeviceSourcePort; suspended: boolean; mcEvidenceByRule?: Readonly<Record<string, PersistedMCRuleState>>;
}
export function RecordingHomeSourceLoading({ onClose }: { onClose: () => void }) {
  const { t } = useTranslation('recording-home');
  return <><div className="gw-modal-body gw-home-source-loading" role="status">{t('loading')}</div>
    <footer className="gw-modal-footer"><button type="button" onClick={onClose}>{t('source.close')}</button></footer></>;
}
function SavedConnection({ device, source, suspended, controls }: { device?: Device; source: DeviceSourcePort; suspended: boolean; controls: HomeSourceControls }) {
  const { t } = useTranslation('recording-editor');
  const buffer = useRecordingDeviceBuffer(controls.target.workspaceId, source, suspended, device);
  const started = useRef(false), innerClose = useRef<(() => void) | undefined>(undefined);
  const [, setConfirm] = useState(false);
  useEffect(() => { if (!started.current) { started.current = true; void buffer.start(); } }, [buffer]);
  useLayoutEffect(() => { controls.closeRef.current = () => { if (!buffer.locked) innerClose.current?.(); }; return () => { controls.closeRef.current = undefined; }; }, [buffer.locked, controls.closeRef]);
  useRecordingDraftGuard({ version: JSON.stringify([controls.target.epoch, buffer.buffer]), dirty: buffer.dirty && !buffer.locked, pending: buffer.locked, host: true, discard: buffer.reset });
  return <RecordingSourceDeviceContent cancelLabel={t('device.cancelEdit')} existingEdit={Boolean(device)} controller={{ ...buffer, disabled: buffer.disabled || buffer.preparing || buffer.blocked }} slot={{ active: true, closeRef: innerClose, onConfirmChange: setConfirm,
    onReturn: () => { if (!buffer.pending) controls.onClose(); }, onSaved: saved => { if (controls.target.kind === 'create') controls.onCreated(saved.id); else controls.onClose(); } }} />;
}
function NewRange({ device, source, suspended, controls, leave, list, refreshList }: { device: Device; source: DeviceSourcePort; suspended: boolean; controls: HomeSourceControls; leave: HomeSourceLeave; list: ReactNode; refreshList: () => void }) {
  const { t } = useTranslation('recording-home');
  const startAddress = controls.target.startAddress;
  const range = useRecordingRangeBuffer(controls.target.workspaceId, source.range, suspended);
  const started = useRef(false), refreshed = useRef<string | undefined>(undefined), [confirm, setConfirm] = useState<(() => void)>();
  useEffect(() => {
    const saved = range.outcome;
    if (!saved || !['confirmed', 'saved-inactive'].includes(saved.status) || !saved.rule?.id || saved.request.workspaceId !== controls.target.workspaceId ||
        saved.request.deviceId !== device.id || saved.rule.workspace_id !== controls.target.workspaceId || saved.rule.device_id !== device.id || refreshed.current === saved.request.requestId) return;
    refreshed.current = saved.request.requestId; refreshList();
  }, [range.outcome, controls.target.workspaceId, device.id, refreshList]);
  // A memory-view address is a non-dirty prefill only; recovery of an unresolved save keeps priority.
  useEffect(() => { if (!started.current) { started.current = true; void range.start(device, startAddress ? { address: startAddress } : undefined); } }, [range, device, startAddress]);
  useLayoutEffect(() => { leave.current = proceed => { if (range.locked) return; if (range.dirty) setConfirm(() => proceed); else proceed(); }; return () => { leave.current = undefined; }; }, [range.locked, range.dirty, leave]);
  useRecordingDraftGuard({ version: JSON.stringify([controls.target.epoch, range.buffer]), dirty: range.dirty && !range.locked, pending: range.locked, host: true, discard: range.reset });
  if (confirm) return <><div className="gw-modal-body">{t('source.discardDescription')}</div><footer className="gw-modal-footer">
    <button type="button" data-testid="home-source-keep" onClick={() => setConfirm(undefined)}>{t('source.keep')}</button>
    <button type="button" data-testid="home-source-discard" onClick={() => { const proceed = confirm; range.reset(); setConfirm(undefined); proceed(); }}>{t('source.discard')}</button></footer></>;
  const prefilled = startAddress && range.buffer.address === startAddress && !range.outcome ? <p role="status" data-testid="home-range-prefilled">{t('source.memory.prefilled', { address: startAddress })}</p> : undefined;
  return <RecordingRangeContent controller={range} available={!suspended} beforeForm={list ?? prefilled}
    onReview={ruleId => controls.onSelect(ruleId)} onBack={() => leave.current?.(() => controls.onSelect())} backLabel={t('source.backToList')}
    onSubmit={() => { void range.submit(); }} onInspect={() => { void range.inspect(); }} onPrepare={() => { void range.prepareSaved(); }} />;
}
/** The original scoped Page supplies every command owner. This component selects
 * local form content only; list responses cannot install or authorize an editor. */
export function RecordingHomeSourceContent(props: RecordingHomeSourceContentProps) {
  const { target, state, dispatch, source, suspended, onDeviceName, onReady, closeRef, onClose, onSelect } = props;
  const { t } = useTranslation('recording-home');
  const leave = useRef<((proceed: () => void) => void) | undefined>(undefined);
  const list = useQuery({ queryKey: ['recording-home-device-ranges', target.workspaceId, target.deviceId, target.epoch],
    queryFn: () => { if (target.kind !== 'points') throw new Error('Range list requires an existing device'); return readHomeDeviceRanges(target.workspaceId, target.deviceId); }, enabled: target.kind === 'points' && !suspended,
    retry: false, refetchOnWindowFocus: false, gcTime: 0 });
  const initialSelection = useRef(false);
  const listed = list.data?.rules;
  useEffect(() => {
    if (!target.initialSelection || initialSelection.current || target.kind !== 'points' || target.ruleId || target.newRange || target.view || suspended || !listed || list.isFetching) return;
    if (!state.devices.some(device => device.id === target.deviceId && device.persisted)) return;
    if (!listed.every(rule => rule.workspace_id === target.workspaceId && rule.device_id === target.deviceId)) return;
    initialSelection.current = true;
    if (listed.length === 1) onSelect(listed[0].id);
    else if (listed.length === 0) onSelect(undefined, true);
  }, [target, suspended, list.isFetching, listed, onSelect, state.devices]);
  const device = state.devices.find(item => item.id === target.deviceId && item.persisted === true);
  const autoSelecting = target.kind === 'points' && target.initialSelection && !target.ruleId && !target.newRange && device && listed &&
    listed.length <= 1 && listed.every(rule => rule.workspace_id === target.workspaceId && rule.device_id === target.deviceId);
  const loadingRanges = target.kind === 'points' && !target.ruleId && !target.newRange && !target.view && (list.isPending || list.isFetching || autoSelecting);
  const ranges = useHomeRangeInputs(list.isFetching ? undefined : list.data, list.dataUpdatedAt, device?.protocol, target.workspaceId);
  useLayoutEffect(() => { if (device) onDeviceName(device.name); }, [device, onDeviceName]);
  useLayoutEffect(() => { if (!suspended && (target.kind === 'connection' && device || target.kind === 'points' && (!device || !loadingRanges))) onReady(); }, [target.kind, suspended, device, loadingRanges, onReady]);
  useLayoutEffect(() => {
    if (target.kind === 'connection' || target.kind === 'create') return;
    closeRef.current = () => { if (leave.current) leave.current(onClose); else onClose(); };
    return () => { closeRef.current = undefined; };
  }, [target.kind, closeRef, onClose]);
  const select = (ruleId?: string, newRange?: boolean, options?: Parameters<typeof onSelect>[2]) => {
    const next = () => onSelect(ruleId, newRange, options);
    if (leave.current) leave.current(next); else next();
  };
  if (target.kind === 'create' && !suspended) return <SavedConnection source={source} suspended={suspended} controls={props} />;
  if (target.kind === 'create' || !device || suspended) return <div className="gw-modal-body" role="status">{t('source.unavailable')}</div>;
  if (target.kind === 'connection') return <SavedConnection device={device} source={source} suspended={suspended} controls={props} />;
  if (loadingRanges) return <RecordingHomeSourceLoading onClose={onClose} />;
  const listContent = <section aria-label={t('source.ranges')} className="gw-home-range-list">
    <header className="gw-range-list-heading"><h3>{t('source.ranges')}</h3>{list.isSuccess && !list.isFetching && <button type="button" data-testid="home-range-new" onClick={() => select(undefined, true)}>{t('source.newRange')}</button>}</header>
    {list.isPending || list.isFetching ? <p role="status">{t('loading')}</p> : list.isError ? <p role="alert">{t('source.unavailable')}</p> : <>
      <ul>{listed?.map(item => {
        const rule = state.rules.find(saved => saved.id === item.id && saved.workspace_id === target.workspaceId && saved.device_id === target.deviceId && saved.save_state === 'saved') ?? item;
        const input = ranges.inputs.find(value => value.id === item.id), span = input && device ? describeMemoryRule(device.protocol, { ...input, label: rule.naming_prefix }) : undefined;
        const footprint = span?.cells !== undefined ? t(`source.memory.${span.unit}`, { count: span.cells }) : t('source.memory.footprintUnknown');
        return <li key={rule.id}><div><strong>{rule.naming_prefix}</strong>{!rule.enabled && <span className="gw-range-flag">{t('source.memory.disabled')}</span>}
          <p className="gw-range-facts"><span className="gw-range-address">{span?.end && span.end !== span.start ? `${span.start} – ${span.end}` : span?.start ?? rule.start_address}</span>
            <span className="gw-range-fact"> · {t('source.memory.items', { count: rule.count })}</span><span className="gw-range-fact"> · <span className="gw-range-type">{rule.data_type}</span></span><span className="gw-range-fact"> · {footprint}</span>
            {rule.skipped_addresses.length > 0 && <span className="gw-range-fact"> · {t('source.memory.skipped', { count: rule.skipped_addresses.length })}</span>}</p></div>
        <button type="button" data-testid={`home-range-edit-${rule.id}`} onClick={() => select(rule.id)}>{t('source.edit')}</button></li>; })}</ul>
      {!listed?.length && <p>{t('source.emptyRanges')}</p>}
    </>}
  </section>;
  /** Unconfirmed new-range saves, unsaved local range edits, or a list that disagrees with
   * the hydrated revision make every free cell unknown (fail closed). */
  function memoryBlock(): HomeMemoryBlock | undefined {
    if (props.source.range?.observations().some(item => item.request.workspaceId === target.workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(item.status))) return 'recovery';
    const local = state.rules.filter(item => item.device_id === target.deviceId && (!item.workspace_id || item.workspace_id === target.workspaceId));
    if (local.some(item => item.save_state !== undefined && item.save_state !== 'saved')) return 'local';
    if (local.some(item => listed?.some(value => value.id === item.id && item.revision_id && value.revision_id !== item.revision_id))) return 'stale';
  }
  const rule = state.rules.find(item => item.id === target.ruleId && item.workspace_id === target.workspaceId && item.device_id === target.deviceId && item.persisted);
  if (target.ruleId) return rule ? <RecordingHomeSavedRangeEditor shareEnabled={props.shareEnabled} onRetryShare={props.onRetryShare} onOpenShareSettings={props.onOpenShareSettings} target={target} state={state} rule={rule} mcEvidence={props.mcEvidenceByRule?.[rule.id]} dispatch={dispatch} leave={leave} suspended={suspended} list={undefined} onDone={() => select()} onLeave={props.onClose} ruleObserverRef={props.ruleObserverRef} />
    : <div className="gw-modal-body" role="alert">{t('source.unavailable')}</div>;
  if (target.view === 'memory') return <RecordingHomeMemoryView protocol={device.protocol} inputs={ranges.inputs} checking={ranges.checking} blocked={memoryBlock()}
    listScopeMissing={ranges.listScopeMissing} onRetry={() => { void list.refetch(); }}
    status={list.isError ? 'error' : list.isPending || list.isFetching ? 'loading' : 'ready'} onBack={() => select()}
    onOpen={ruleId => select(ruleId)} onAdd={address => select(undefined, true, { startAddress: address })} />;
  if (target.newRange) return <NewRange device={device} source={source} suspended={suspended} controls={props} leave={leave} list={undefined} refreshList={() => { void list.refetch(); }} />;
  return <div className="gw-modal-body">{listContent}</div>;
}
