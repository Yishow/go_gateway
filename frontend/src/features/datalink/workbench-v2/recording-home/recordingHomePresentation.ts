import { initialCollectionState, type RecordingCollectionState } from './recordingHomeCollections';
export interface HomePresentation {
  state: RecordingCollectionState; workspaceRevision: string; scrollTop: number; focusId: string | null; observedIds?: readonly string[];
}
// Presentation only, scoped to this QueryClient/backend and workspace. No saved
// selection, receipt, permissions or command intent is restored from this cache.
const remembered = new WeakMap<object, Map<string, HomePresentation>>();
export function readHomePresentation(context: object, workspaceId: string, workspaceRevision: string): HomePresentation {
  const previous = remembered.get(context)?.get(workspaceId);
  return previous ? { ...previous, state: { ...previous.state }, observedIds: previous.observedIds?.slice() } :
    { state: initialCollectionState(), workspaceRevision, scrollTop: 0, focusId: null };
}
export function rememberHomePresentation(context: object, workspaceId: string, value: HomePresentation): void {
  const entries = remembered.get(context) ?? new Map<string, HomePresentation>();
  entries.delete(workspaceId); entries.set(workspaceId, { ...value, state: { ...value.state }, observedIds: value.observedIds?.slice() });
  while (entries.size > 8) entries.delete(entries.keys().next().value!);
  remembered.set(context, entries);
}
