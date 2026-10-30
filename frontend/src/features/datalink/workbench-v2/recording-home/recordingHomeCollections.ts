import type { WriteGroup, WriteGroupStatus } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
export interface HomeDeviceSummary { readonly id: string; readonly name: string; readonly protocol: string; readonly availability: string }
export interface RecordingCollectionState {
  mode: 'recordings' | 'devices'; keyword: string; deviceId: string; destination: string;
  /** Saved configuration status, never actual runtime health. */
  status: WriteGroupStatus | ''; observedOnly: boolean; page: number; pageSize: 25 | 50; descending: boolean;
}
export function initialCollectionState(): RecordingCollectionState {
  return { mode: 'recordings', keyword: '', deviceId: '', destination: '', status: '', observedOnly: false, page: 1, pageSize: 25, descending: false };
}
const identity = (a: string, b: string) => a < b ? -1 : a > b ? 1 : 0;
const searchText = (value: string) => value.trim().toLocaleLowerCase();
function byName<T extends { readonly id: string; readonly name: string }>(descending: boolean) {
  return (a: T, b: T) => (identity(searchText(a.name), searchText(b.name)) || identity(a.id, b.id)) * (descending ? -1 : 1);
}
function pageRows<T>(all: readonly T[], filtered: readonly T[], state: RecordingCollectionState) {
  const pageSize = state.pageSize === 50 ? 50 : 25;
  const pageCount = Math.max(1, Math.ceil(filtered.length / pageSize));
  const page = Math.min(pageCount, Math.max(1, Number.isSafeInteger(state.page) ? state.page : 1));
  return { rows: filtered.slice((page - 1) * pageSize, page * pageSize), page, pageCount, pageSize, total: all.length, filteredCount: filtered.length,
    empty: !all.length ? 'empty' as const : !filtered.length ? 'no-results' as const : null };
}
/** Canonical data stays untouched. This is not server-side search or a health index. */
export function recordingCollection(groups: readonly RecordingReadonly<WriteGroup>[], devices: readonly HomeDeviceSummary[], state: RecordingCollectionState, observedIds: ReadonlySet<string> = new Set()) {
  const names = new Map(devices.map(device => [device.id, device.name]));
  const keyword = searchText(state.keyword);
  const filtered = groups.filter(group => (!state.observedOnly || observedIds.has(group.id)) && (!state.status || group.status === state.status) &&
    (!state.destination || group.destination.connector_id === state.destination) &&
    (!state.deviceId || group.members.some(member => member.device_id === state.deviceId) || group.basic_managed_device_id === state.deviceId) &&
    (!keyword || searchText([group.name, group.destination.database, group.destination.table_schema, group.destination.table_name,
      ...group.members.map(member => names.get(member.device_id) ?? member.device_id), names.get(group.basic_managed_device_id ?? '') ?? ''].join(' ')).includes(keyword)))
    .sort(byName(state.descending));
  return pageRows(groups, filtered, state);
}
export function deviceCollection(devices: readonly HomeDeviceSummary[], state: RecordingCollectionState) {
  const keyword = searchText(state.keyword);
  const sorted = [...devices].sort(byName(state.descending));
  const filtered = sorted.filter(device => !keyword || searchText(`${device.name} ${device.id} ${device.protocol}`).includes(keyword));
  return { ...pageRows(devices, filtered, state), summary: sorted.slice(0, 4) };
}
