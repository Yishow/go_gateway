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
import { useHomeMemoryValues } from './recordingHomeMemoryValues';
import { RecordingHomeRangeOptions } from './RecordingHomeRangeOptions';
import { useRuntimeStatus } from '../../runtime-dashboard/useRuntimeStatus';
import { runtimeIntervalReadout } from '../state/acquisitionSettings';
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
  const runtime = useRuntimeStatus({ deviceId: device?.id ?? null, pollingIntervalMs: 5000 });
  const started = useRef(false), innerClose = useRef<(() => void) | undefined>(undefined);
  const [, setConfirm] = useState(false);
  useEffect(() => { if (!started.current) { started.current = true; void buffer.start(); } }, [buffer]);
  useLayoutEffect(() => { controls.closeRef.current = () => { if (!buffer.locked) innerClose.current?.(); }; return () => { controls.closeRef.current = undefined; }; }, [buffer.locked, controls.closeRef]);
  useRecordingDraftGuard({ version: JSON.stringify([controls.target.epoch, buffer.buffer]), dirty: buffer.dirty && !buffer.locked, pending: buffer.locked, host: true, discard: buffer.reset });
  return <RecordingSourceDeviceContent cancelLabel={t('device.cancelEdit')} existingEdit={Boolean(device)} runtimeInterval={device ? runtimeIntervalReadout(device.id, runtime) : undefined} controller={{ ...buffer, disabled: buffer.disabled || buffer.preparing || buffer.blocked }} slot={{ active: true, closeRef: innerClose, onConfirmChange: setConfirm,
    onReturn: () => { if (!buffer.pending) controls.onClose(); }, onSaved: saved => { if (controls.target.kind === 'create') controls.onCreated(saved.id); else controls.onClose(); } }} />;
}
function NewRange({ device, source, suspended, controls, leave, options, hasRanges, refreshList, close }: { device: Device; source: DeviceSourcePort; suspended: boolean; controls: HomeSourceControls; leave: HomeSourceLeave; options: ReactNode; hasRanges: boolean; refreshList: () => void; close: () => void }) {
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
  // An open keep/discard question keeps its original destination; no other request may replace it.
  useLayoutEffect(() => { leave.current = proceed => { if (range.locked || confirm) return; if (range.dirty) setConfirm(() => proceed); else proceed(); }; return () => { leave.current = undefined; }; }, [range.locked, range.dirty, confirm, leave]);
  const { onLocked } = controls;
  useLayoutEffect(() => { onLocked(range.locked || Boolean(confirm)); }, [range.locked, confirm, onLocked]);
  useRecordingDraftGuard({ version: JSON.stringify([controls.target.epoch, range.buffer]), dirty: range.dirty && !range.locked, pending: range.locked, host: true, discard: range.reset });
  if (confirm) return <><div className="gw-modal-body">{t('source.discardDescription')}</div><footer className="gw-modal-footer">
    <button type="button" data-testid="home-source-keep" onClick={() => setConfirm(undefined)}>{t('source.keep')}</button>
    <button type="button" data-testid="home-source-discard" onClick={() => { const proceed = confirm; range.reset(); setConfirm(undefined); proceed(); }}>{t('source.discard')}</button></footer></>;
  const prefilled = startAddress && range.buffer.address === startAddress && !range.outcome ? <p role="status" data-testid="home-range-prefilled">{t('source.memory.prefilled', { address: startAddress })}</p> : undefined;
  // Without another range (or after this one is saved) the only way back is closing the frame.
  const closes = !hasRanges || range.outcome?.status === 'confirmed' || range.outcome?.status === 'saved-inactive';
  return <RecordingRangeContent controller={range} available={!suspended} beforeForm={<><fieldset className="gw-home-range-switch" disabled={range.locked || suspended}>{options}</fieldset>{prefilled}</>}
    onReview={ruleId => controls.onSelect(ruleId)} onBack={() => leave.current?.(closes ? close : () => controls.onSelect())} backLabel={t(closes ? 'source.close' : 'source.cancelNewRange')}
    onSubmit={() => { void range.submit(); }} onInspect={() => { void range.inspect(); }} onPrepare={() => { void range.prepareSaved(); }} />;
}
/** The original scoped Page supplies every command owner. This component selects
 * local form content only; list responses cannot install or authorize an editor. */
export function RecordingHomeSourceContent(props: RecordingHomeSourceContentProps) {
  const { target, state, dispatch, source, suspended, onDeviceName, onReady, closeRef, leaveRef, onClose, onSelect } = props;
  const { t } = useTranslation('recording-home');
  const leave = useRef<((proceed: () => void) => void) | undefined>(undefined);
  const list = useQuery({ queryKey: ['recording-home-device-ranges', target.workspaceId, target.deviceId, target.epoch],
    queryFn: () => { if (target.kind !== 'points') throw new Error('Range list requires an existing device'); return readHomeDeviceRanges(target.workspaceId, target.deviceId); }, enabled: target.kind === 'points' && !suspended,
    retry: false, refetchOnWindowFocus: false, gcTime: 0 });
  const resolved = useRef(false);
  const listed = list.data?.rules;
  const device = state.devices.find(item => item.id === target.deviceId && item.persisted === true);
  const unselected = target.kind === 'points' && !target.ruleId && !target.newRange && !target.view;
  // No separate range list: an unselected device opens its first range, or a new one when it has none.
  useEffect(() => {
    if (!unselected || resolved.current || suspended || !listed || list.isFetching || !device) return;
    if (!listed.every(rule => rule.workspace_id === target.workspaceId && rule.device_id === target.deviceId)) return;
    resolved.current = true;
    if (listed.length) onSelect(listed[0].id); else onSelect(undefined, true);
  }, [unselected, target.workspaceId, target.deviceId, suspended, list.isFetching, listed, onSelect, device]);
  const loadingRanges = target.kind === 'points' && (list.isPending || unselected && !list.isError);
  const ranges = useHomeRangeInputs(list.isFetching ? undefined : list.data, list.dataUpdatedAt, device?.protocol, target.workspaceId);
  const memoryValues = useHomeMemoryValues(device?.id, device?.protocol, target.view === 'memory');
  useLayoutEffect(() => { if (device) onDeviceName(device.name); }, [device, onDeviceName]);
  useLayoutEffect(() => { if (!suspended && (target.kind === 'connection' && device || target.kind === 'points' && (!device || !loadingRanges))) onReady(); }, [target.kind, suspended, device, loadingRanges, onReady]);
  useLayoutEffect(() => {
    if (target.kind === 'connection' || target.kind === 'create') return;
    closeRef.current = () => { if (leave.current) leave.current(onClose); else onClose(); };
    leaveRef.current = proceed => { if (leave.current) leave.current(proceed); else proceed(); };
    return () => { closeRef.current = undefined; leaveRef.current = undefined; };
  }, [target.kind, closeRef, leaveRef, onClose]);
  const select = (ruleId?: string, newRange?: boolean, options?: Parameters<typeof onSelect>[2]) => {
    const next = () => onSelect(ruleId, newRange, options);
    if (leave.current) leave.current(next); else next();
  };
  if (target.kind === 'create' && !suspended) return <SavedConnection source={source} suspended={suspended} controls={props} />;
  if (target.kind === 'create' || suspended) return <div className="gw-modal-body" role="status">{t('source.unavailable')}</div>;
  if (target.kind === 'connection') return device ? <SavedConnection device={device} source={source} suspended={suspended} controls={props} /> : <div className="gw-modal-body" role="status">{t('source.unavailable')}</div>;
  if (!device) return <div className="gw-modal-body" role="alert">{t('source.unavailable')}</div>;
  if (loadingRanges) return <RecordingHomeSourceLoading onClose={onClose} />;
  const retry = () => { void list.refetch(); };
  if (unselected) return <><div className="gw-modal-body" role="alert">{t('source.rangesUnavailable')}</div>
    <footer className="gw-modal-footer"><button type="button" onClick={onClose}>{t('source.close')}</button><button type="button" className="gw-primary" onClick={retry}>{t('source.memory.retry')}</button></footer></>;
  const options = <RecordingHomeRangeOptions protocol={device.protocol} listed={listed} inputs={ranges.inputs} state={state} target={target}
    ready={list.isSuccess && !list.isFetching} failed={list.isError} onRetry={retry} onSelect={(ruleId, newRange) => select(ruleId, newRange)} />;
  /** Unconfirmed new-range saves, unsaved local range edits, or a list that disagrees with
   * the hydrated revision make every free cell unknown (fail closed). */
  function memoryBlock(): HomeMemoryBlock | undefined {
    if (props.source.range?.observations().some(item => item.request.workspaceId === target.workspaceId && ['pending', 'unknown', 'persisted-unavailable'].includes(item.status))) return 'recovery';
    const local = state.rules.filter(item => item.device_id === target.deviceId && (!item.workspace_id || item.workspace_id === target.workspaceId));
    if (local.some(item => item.save_state !== undefined && item.save_state !== 'saved')) return 'local';
    if (local.some(item => listed?.some(value => value.id === item.id && item.revision_id && value.revision_id !== item.revision_id))) return 'stale';
  }
  const rule = state.rules.find(item => item.id === target.ruleId && item.workspace_id === target.workspaceId && item.device_id === target.deviceId && item.persisted);
  if (target.ruleId) return rule ? <RecordingHomeSavedRangeEditor shareEnabled={props.shareEnabled} onRetryShare={props.onRetryShare} onOpenShareSettings={props.onOpenShareSettings} target={target} state={state} rule={rule} mcEvidence={props.mcEvidenceByRule?.[rule.id]} dispatch={dispatch} leave={leave} suspended={suspended} list={options}
    headerSlot={props.headerSlot} onLocked={props.onLocked} onClose={() => leave.current?.(onClose)} onLeave={props.onClose} ruleObserverRef={props.ruleObserverRef} />
    : <div className="gw-modal-body" role="alert">{t('source.unavailable')}</div>;
  if (target.view === 'memory') return <RecordingHomeMemoryView protocol={device.protocol} inputs={ranges.inputs} checking={ranges.checking} blocked={memoryBlock()}
    listScopeMissing={ranges.listScopeMissing} onRetry={retry} values={memoryValues}
    status={list.isError ? 'error' : list.isPending || list.isFetching ? 'loading' : 'ready'} onBack={() => select(target.returnRuleId, target.returnNewRange || !target.returnRuleId && !listed?.length)}
    onClose={() => { if (closeRef.current) closeRef.current(); else onClose(); }}
    onOpen={ruleId => select(ruleId)} onAdd={address => select(undefined, true, { startAddress: address })} />;
  return <NewRange device={device} source={source} suspended={suspended} controls={props} leave={leave} options={options} hasRanges={Boolean(listed?.length)} refreshList={retry} close={onClose} />;
}
