import { recordingHomeSourceContext, type HomeSourceContext } from '../../../../hooks/datalink/recordingHomeSourceCache';
import { initialCollectionState, type RecordingCollectionState } from './recordingHomeCollections';
export interface HomePresentation {
  state: RecordingCollectionState; workspaceRevision: string; scrollTop: number; focusId: string | null; observedIds?: readonly string[];
}
// Presentation only, scoped to this QueryClient/backend and workspace. No saved
// selection, receipt, permissions or command intent is restored from this cache.
const remembered = new WeakMap<object, Map<string, Map<string, HomePresentation>>>();
function entriesFor(context: object, sourceContext: HomeSourceContext): Map<string, HomePresentation> {
  const { apiOrigin, backendToken } = recordingHomeSourceContext(sourceContext);
  const key = JSON.stringify([backendToken, apiOrigin]);
  let partitions = remembered.get(context);
  if (!partitions) { partitions = new Map(); remembered.set(context, partitions); }
  let entries = partitions.get(key);
  if (!entries) { entries = new Map(); partitions.set(key, entries); }
  return entries;
}
export function readHomePresentation(context: object, workspaceId: string, workspaceRevision: string, sourceContext: HomeSourceContext = {}): HomePresentation {
  const previous = entriesFor(context, sourceContext).get(workspaceId);
  return previous ? { ...previous, state: { ...previous.state }, observedIds: previous.observedIds?.slice() } :
    { state: initialCollectionState(), workspaceRevision, scrollTop: 0, focusId: null };
}
export function rememberHomePresentation(context: object, workspaceId: string, value: HomePresentation, sourceContext: HomeSourceContext = {}): void {
  const entries = entriesFor(context, sourceContext);
  entries.delete(workspaceId); entries.set(workspaceId, { ...value, state: { ...value.state }, observedIds: value.observedIds?.slice() });
  while (entries.size > 8) entries.delete(entries.keys().next().value!);
}
