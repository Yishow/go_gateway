import { cachedHomeReads } from './recordingHomeCachedReads';
import { recordingHomeGuidance } from './recordingHomeGuidance';
import type { WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import { recordingHomeProjection } from './recordingHomeProjection';
import { RecordingHomeCollection } from './RecordingHomeCollection';
import { recordingCollection, deviceCollection } from './recordingHomeCollections';
import { readHomePresentation, rememberHomePresentation } from './recordingHomePresentation';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import { buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceQuery } from '../../../../hooks/datalink/useStudioV2Workspace';
import { useWriteGroupsQuery } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';
import { readRecordingHomeSources } from './recordingHomeSources';
import { useRecordingHomeEvidence } from './useRecordingHomeEvidence';
import { RecordingHome, type RecordingHomeLocation } from './RecordingHome';
import type { RecordingHomeScope } from './recordingHomeProjection';

// Request generation only, not configuration or workspace authority. Remounting
// cannot reuse pending evidence from a previous visit, even A -> B -> A.
let generation = 0;
function ScopedRecordingHome({ workspaceId, workspaceRevision, list, location, onNavigate, reloadScope }: {
  workspaceId: string; workspaceRevision: string; list: WriteGroupListResponse; location: RecordingHomeLocation; onNavigate?: (target: string) => void; reloadScope: () => Promise<unknown>;
}) {
  const context = useQueryClient();
  const [restored, setRestored] = useState(() => readHomePresentation(context, workspaceId, workspaceRevision));
  const [state, setState] = useState(restored.state);
  // Freeze membership for this user-selected page. Recovered rows stay visible
  // until another filter/page intent; replies cannot pull the next batch in.
  const [filterObservedIds, setFilterObservedIds] = useState<readonly string[]>(restored.observedIds ?? []);
  const [demand, setDemand] = useState<{ state: typeof state; ids: string[] } | null>(null);
  const [scope] = useState<RecordingHomeScope>(() => ({ workspaceId, workspaceRevision, epoch: ++generation }));
  const sources = useQuery({ queryKey: [...studioV2WorkspaceKeys.all, 'recording-home-sources', scope.workspaceId, scope.workspaceRevision, scope.epoch],
    queryFn: () => readRecordingHomeSources(scope), retry: false, refetchOnWindowFocus: false, gcTime: 0 });
  const groupsRead = recordingHomeProjection({ scope, groups: { state: 'available', scope, value: list } }).groups;
  const groups = groupsRead.state === 'available' ? groupsRead.value.map(item => item.group) : [];
  const devices = sources.data?.state === 'available' ? sources.data.value.devices : [];
  const cached = cachedHomeReads(context, groups.map(group => ({ ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision })));
  const observed = recordingHomeGuidance({ evidence: { scope, groups: { state: 'available', scope, value: list }, ...cached } })
    .attentionItems.filter(item => item.reason !== 'evidence_unavailable').map(item => item.id);
  const unverified = (restored.observedIds ?? []).filter(id => groups.some(group => group.id === id) &&
    (cached.delivery[id]?.state !== 'available' || cached.lifecycle[id]?.state !== 'available'));
  const observedIds = [...new Set([...observed, ...unverified])];
  const collection = recordingCollection(groups, devices, state, new Set(filterObservedIds));
  const page = state.mode === 'devices' ? deviceCollection(devices, state).page : collection.page;
  // Capture at most one page after source names are available (or failed).
  // Same-intent render/refetch/cache updates never replenish this demand.
  if (!sources.isPending && demand?.state !== state) {
    setDemand({ state, ids: state.mode === 'devices' ? [] : collection.rows.map(group => group.id) });
  }
  const home = useRecordingHomeEvidence({ scope, detailGroupIds: demand?.state === state ? demand.ids : [], sources: sources.isFetching || sources.isPending
    ? { state: 'loading', scope } : sources.isError ? { state: 'error', scope } : sources.data });
  return <RecordingHome guidance={home.guidance} location={location} onNavigate={onNavigate}
    onRetry={() => { void Promise.all([reloadScope(), sources.refetch(), home.retry()]); }}
    collection={<RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={home.guidance}
      observedIds={observedIds} filterObservedIds={filterObservedIds} unverifiedObserved={unverified.length > 0} onRecheck={() => { void home.retry(); }}
      sourcesAvailable={sources.data?.state === 'available' && !sources.isFetching} restored={restored}
      onChange={next => {
        if (!state.observedOnly || !next.observedOnly) setFilterObservedIds(observedIds);
        setState(next); const presentation = { state: next, workspaceRevision, scrollTop: 0, focusId: null, observedIds };
        setRestored(presentation); rememberHomePresentation(context, workspaceId, presentation);
      }}
      onRemember={(scrollTop, focusId, currentPage) => rememberHomePresentation(context, workspaceId,
        { state: { ...state, page: currentPage ?? page }, workspaceRevision, scrollTop, focusId, observedIds })}
      onOpen={id => {
        const item = home.guidance.recordingItems.find(candidate => candidate.id === id);
        if (item) {
          const target = buildStudioEntityLocatorTarget({ step: 4, workspace_id: workspaceId, device_id: null, rule_id: null,
            mapping_id: null, group_id: id, group_revision: item.evidence.group.revision });
          if (target) onNavigate?.(target);
        }
      }}
      onOpenDevice={id => {
        if (!devices.some(device => device.id === id)) return;
        const target = buildStudioEntityLocatorTarget({ step: 1, workspace_id: workspaceId, device_id: id, rule_id: null, mapping_id: null, group_id: null, group_revision: null });
        if (target) onNavigate?.(target);
      }} />} />;
}
/** Home mounts no autosave/activation owners. Bootstrap identity and DB setup
 * revision are read independently, with an explicit empty DB revision allowed. */
export function RecordingHomeRoute({ location, onNavigate }: { location: RecordingHomeLocation; onNavigate?: (target: string) => void }) {
  const { t } = useTranslation('recording-home');
  const workspace = useStudioV2WorkspaceQuery();
  const groups = useWriteGroupsQuery(workspace.isSuccess && Boolean(workspace.data?.id));
  const reload = () => Promise.all([workspace.refetch(), groups.refetch()]);
  const failed = workspace.isError || groups.isError;
  const loading = workspace.isFetching || workspace.isPending || groups.isFetching || groups.isPending;
  const known = !failed && !loading && workspace.data?.id && groups.data?.workspace_id === workspace.data.id &&
    typeof groups.data.workspace_revision === 'string';
  if (!known) return <main data-workbench-v2="true" data-testid="recording-home" className="min-h-screen p-6">
    <h1 className="text-2xl font-semibold">{t('title')}</h1>
    <p role={failed ? 'alert' : 'status'}>{t(failed ? 'error' : loading ? 'loading' : 'unknown')}</p>
    <button type="button" onClick={() => { void reload(); }}>{t('retry')}</button>
    <button type="button" className="ml-4" onClick={() => { const target = buildStudioPresentationTarget(location.search, location.hash, 'advanced'); if (target) onNavigate?.(target); }}>{t('advanced')}</button>
  </main>;
  return <ScopedRecordingHome key={JSON.stringify([workspace.data.id, groups.data.workspace_revision])}
    workspaceId={workspace.data.id} workspaceRevision={groups.data.workspace_revision} list={groups.data}
    location={location} onNavigate={onNavigate} reloadScope={reload} />;
}
