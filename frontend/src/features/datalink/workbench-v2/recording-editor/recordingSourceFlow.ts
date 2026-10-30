import type { Device } from '../state/types';
import type { PickerState } from './recordingPickerState';

export type SourceView = 'picker' | 'device-choice' | 'device-buffer' | 'range-buffer';
export interface SourceRequest<T> {
  scope: string;
  session: string;
  requestId: string;
  generation: number;
  kind: 'device' | 'range';
  deviceId?: Device['id'];
  buffer: T;
}
export interface SourceFlow<T> {
  scope: string;
  session: string;
  view: SourceView;
  picker: PickerState;
  targetDeviceId?: Device['id'];
  buffer?: T;
  generation: number;
  dirty: boolean;
  operation?: { request: SourceRequest<T>; status: 'pending' | 'unknown' | 'failed' | 'confirmed' };
  /** Caller-confirmed identities only; retaining these never implies activation or SQL. */
  confirmed: readonly { kind: 'device' | 'range'; deviceId: Device['id']; ruleId?: string }[];
}
export type SourceSettlement = {
  scope: string; session: string; requestId: string; generation: number;
} & ({ status: 'unknown' | 'failed' } | { status: 'confirmed'; deviceId: Device['id']; ruleId?: string });
export function createSourceFlow<T>(picker: PickerState, session: string): SourceFlow<T> {
  return { scope: picker.scope, session, picker, view: 'picker', generation: 0, dirty: false, confirmed: [] };
}
export function beginSource<T>(flow: SourceFlow<T>, devices: readonly Pick<Device, 'id'>[]): SourceFlow<T> {
  if (flow.view !== 'picker') return flow;
  const targetDeviceId = flow.picker.deviceId || undefined;
  if (targetDeviceId && !devices.some(device => device.id === targetDeviceId)) return flow;
  return { ...flow, view: targetDeviceId ? 'range-buffer' : 'device-choice', targetDeviceId };
}
export function chooseSourceDevice<T>(flow: SourceFlow<T>, deviceId: string | undefined, devices: readonly Pick<Device, 'id'>[], buffer: T): SourceFlow<T> {
  if (flow.view !== 'device-choice' || deviceId !== undefined && !devices.some(device => device.id === deviceId)) return flow;
  return { ...flow, view: deviceId ? 'range-buffer' : 'device-buffer', targetDeviceId: deviceId, buffer };
}
export function editSourceBuffer<T>(flow: SourceFlow<T>, buffer: T): SourceFlow<T> {
  if (flow.view !== 'device-buffer' && flow.view !== 'range-buffer') return flow;
  return { ...flow, buffer, generation: flow.generation + 1, dirty: true };
}
export function submitSource<T>(flow: SourceFlow<T>, requestId: string): { flow: SourceFlow<T>; request?: SourceRequest<T> } {
  if (!requestId || flow.buffer === undefined || !flow.dirty ||
    (flow.view !== 'device-buffer' && flow.view !== 'range-buffer') ||
    (flow.view === 'range-buffer' && !flow.targetDeviceId) ||
    flow.operation?.status === 'pending' || flow.operation?.status === 'unknown') return { flow };
  const request: SourceRequest<T> = { scope: flow.scope, session: flow.session, requestId, generation: flow.generation,
    kind: flow.view === 'device-buffer' ? 'device' : 'range', deviceId: flow.targetDeviceId, buffer: flow.buffer };
  return { flow: { ...flow, operation: { request, status: 'pending' } }, request };
}
/** This consumes the integration owner's verified settlement; it does no saving or identity generation. */
export function settleSource<T>(flow: SourceFlow<T>, result: SourceSettlement): SourceFlow<T> {
  const operation = flow.operation;
  if (!operation || operation.status === 'confirmed' || operation.status === 'failed') return flow;
  const request = operation.request;
  if (result.scope !== flow.scope || result.session !== flow.session || result.requestId !== request.requestId || result.generation !== request.generation) return flow;
  if (result.status !== 'confirmed') return { ...flow, operation: { ...operation, status: result.status } };
  if (!result.deviceId || request.kind === 'range' && (!result.ruleId || result.deviceId !== request.deviceId)) return flow;
  const confirmed = [...flow.confirmed, { kind: request.kind, deviceId: result.deviceId, ...(result.ruleId ? { ruleId: result.ruleId } : {}) }];
  const settled = { ...flow, confirmed, operation: { ...operation, status: 'confirmed' as const } };
  if (flow.generation !== request.generation) return settled;
  return { ...settled, dirty: false, buffer: undefined,
    view: request.kind === 'device' ? 'range-buffer' : 'picker',
    targetDeviceId: result.deviceId };
}
export function leaveSource<T>(flow: SourceFlow<T>, discardUnsent = false): { flow: SourceFlow<T>; outcome: 'returned' | 'needs-confirmation' | 'needs-recovery' } {
  if (flow.operation?.status === 'pending' || flow.operation?.status === 'unknown') return { flow, outcome: 'needs-recovery' };
  if (flow.dirty && !discardUnsent) return { flow, outcome: 'needs-confirmation' };
  return { flow: { ...flow, view: 'picker', dirty: false, buffer: undefined, targetDeviceId: undefined }, outcome: 'returned' };
}
