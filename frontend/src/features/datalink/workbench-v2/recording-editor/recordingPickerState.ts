import type { SemanticKind } from '../../../../types/measurement';
import type { WriteGroupMemberDraft } from '../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../state/writeGroup/candidates';

export type SourceIdentity = Pick<WriteGroupMemberDraft, 'device_id' | 'point_id' | 'tag_id'>;
export type PointFilter = 'all' | 'gauge' | 'counter' | 'state';
export interface CandidateEvidence {
  /** Caller-owned validation; candidate IDs alone are NOT revision/CAS proof. */
  availability: 'eligible' | 'invalid' | 'stale' | 'unknown';
  semantic?: SemanticKind;
  reason?: string;
}
export interface PickerState {
  scope: string;
  query: string;
  deviceId: string;
  pointType: PointFilter;
  page: number;
  pageSize: 25 | 50;
  pending: readonly SourceIdentity[];
  scrollTop: number;
  focusKey: string;
}
export interface PickerProjection {
  scope: string;
  candidates: readonly GroupCandidate[];
  evidence: Readonly<Record<string, CandidateEvidence | undefined>>;
  members: readonly SourceIdentity[];
}
export function sourceKey(source: SourceIdentity): string {
  return JSON.stringify([source.device_id, source.point_id, source.tag_id]);
}
export function createPicker(scope: string): PickerState {
  return { scope, query: '', deviceId: '', pointType: 'all', page: 0, pageSize: 25, pending: [], scrollTop: 0, focusKey: '' };
}
export function filterPicker(state: PickerState, patch: Partial<Pick<PickerState, 'query' | 'deviceId' | 'pointType' | 'pageSize'>>): PickerState {
  return { ...state, ...patch, page: 0, scrollTop: 0 };
}
function eligible(source: SourceIdentity, projection: PickerProjection): boolean {
  const key = sourceKey(source);
  return Boolean(source.device_id && source.point_id && source.tag_id &&
    projection.candidates.filter(row => sourceKey(row) === key).length === 1 &&
    projection.evidence[key]?.availability === 'eligible');
}
export function togglePending(state: PickerState, source: SourceIdentity, projection: PickerProjection): PickerState {
  if (state.scope !== projection.scope) return state;
  const key = sourceKey(source);
  if (projection.members.some(member => sourceKey(member) === key)) return state;
  if (state.pending.some(member => sourceKey(member) === key)) {
    return { ...state, pending: state.pending.filter(member => sourceKey(member) !== key) };
  }
  if (!eligible(source, projection)) return state;
  const { device_id, point_id, tag_id } = source;
  return { ...state, pending: [...state.pending, { device_id, point_id, tag_id }] };
}
export function pickerResults(state: PickerState, projection: PickerProjection) {
  const matchingScope = state.scope === projection.scope;
  const query = state.query.trim().toLocaleLowerCase();
  const members = new Set(projection.members.map(sourceKey));
  const filtered = matchingScope ? projection.candidates.filter(row => {
    const semantic = projection.evidence[sourceKey(row)]?.semantic;
    const typeMatches = state.pointType === 'all' || semantic === state.pointType ||
      state.pointType === 'counter' && semantic === 'signed_counter';
    return (!state.deviceId || row.device_id === state.deviceId) && typeMatches &&
      (!query || [row.label, row.device_name, row.address, row.tag_key].some(text => text.toLocaleLowerCase().includes(query)));
  }).sort((a, b) => sourceKey(a) < sourceKey(b) ? -1 : sourceKey(a) > sourceKey(b) ? 1 : 0) : [];
  const pages = Math.max(1, Math.ceil(filtered.length / state.pageSize));
  const page = Math.max(0, Math.min(Math.floor(state.page) || 0, pages - 1));
  const rows = filtered.slice(page * state.pageSize, (page + 1) * state.pageSize);
  const visible = new Set(rows.map(sourceKey));
  const invalidPending = state.pending.filter(source => !matchingScope || !eligible(source, projection) || members.has(sourceKey(source)));
  return { rows, total: filtered.length, added: filtered.filter(row => members.has(sourceKey(row))).length,
    page, pages, hiddenPending: state.pending.filter(source => !visible.has(sourceKey(source))).length,
    invalidPending, canAdd: state.pending.length > 0 && invalidPending.length === 0 };
}
/** An addition intent only: the integration owner must revalidate revisions before saving. */
export function pendingAddition(state: PickerState, projection: PickerProjection): readonly SourceIdentity[] | undefined {
  return pickerResults(state, projection).canAdd ? state.pending : undefined;
}
export function cancelPicker(state: PickerState): PickerState { return { ...state, pending: [] }; }
/** Session-only UI history. Never persisted, never interpreted as source validity. */
export interface MemberHistory<T extends SourceIdentity> { scope: string; removed: ReadonlyMap<string, T> }
export function removeMember<T extends SourceIdentity>(members: readonly T[], history: MemberHistory<T>, key: string) {
  const member = members.find(item => sourceKey(item) === key);
  if (!member) return { members, history };
  const removed = new Map(history.removed);
  removed.set(key, member);
  return { members: members.filter(item => sourceKey(item) !== key), history: { ...history, removed } };
}
export function restoreMember<T extends SourceIdentity>(members: readonly T[], history: MemberHistory<T>, source: SourceIdentity, projection: PickerProjection, initialize: (source: SourceIdentity) => T): readonly T[] {
  if (history.scope !== projection.scope || !eligible(source, projection)) return members;
  const key = sourceKey(source);
  if (members.some(member => sourceKey(member) === key)) return members;
  const restored = history.removed.get(key) ?? initialize(source);
  // Defaults may supply fields, but may never silently substitute another source.
  if (sourceKey(restored) !== key) return members;
  return [...members, restored];
}
