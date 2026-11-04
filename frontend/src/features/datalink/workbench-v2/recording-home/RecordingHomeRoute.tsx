import { recordingHomeResultGuidance } from './recordingHomeResultGuidance';
import { useRecordingHomeFacts } from './useRecordingHomeFacts';
import { homeOwners } from './recordingHomeGuidanceEvidence';
import { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { cachedHomeReads } from './recordingHomeCachedReads';
import type { WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import { recordingHomeProjection } from './recordingHomeProjection';
import { RecordingHomeCollection } from './RecordingHomeCollection';
import { recordingCollection, deviceCollection } from './recordingHomeCollections';
import { readHomePresentation, rememberHomePresentation } from './recordingHomePresentation';
import { buildSimpleRecordingTarget, buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { useEffect, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceQuery } from '../../../../hooks/datalink/useStudioV2Workspace';
import { useWriteGroupsQuery } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';
import { readRecordingHomeSources } from './recordingHomeSources';
import { useRecordingHomeEvidence } from './useRecordingHomeEvidence';
import { RecordingHome, RecordingHomeTools, type RecordingHomeLocation } from './RecordingHome';
import { RecordingProductShell } from '../shell/RecordingProductShell';
import type { RecordingHomeScope } from './recordingHomeProjection';
import type { HomeSourceDevice } from './RecordingHomeSourceHost';

// Request generation only, not configuration or workspace authority. Remounting
// cannot reuse pending evidence from a previous visit, even A -> B -> A.
let generation = 0;
/** Upper bound for holding the first-load reading state; later reads render progressively as before. */
export const RECORDING_HOME_SETTLE_LIMIT_MS = 2500;
function ScopedRecordingHome({ workspaceId, workspaceRevision, list, location, onNavigate, reloadScope, onOpenDevice, onCreateDevice }: {
  workspaceId: string; workspaceRevision: string; list: WriteGroupListResponse; location: RecordingHomeLocation; onNavigate?: (target: string) => void; reloadScope: () => Promise<unknown>; onOpenDevice?: (workspaceId: string, deviceId: string, kind: 'connection' | 'points', devices?: readonly HomeSourceDevice[]) => void; onCreateDevice?: (workspaceId: string) => void;
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
  // The same last validated projection the collection shows (a failed refresh replaces it with an error state).
  // Every device chosen in the settings frame is re-validated by its own scoped locator read.
  const switchable = sources.data?.state === 'available' ? devices.map(({ id, name }) => ({ id, name })) : undefined;
  const cached = cachedHomeReads(context, groups.map(group => ({ ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision })));
  const collection = recordingCollection(groups, devices, state, new Set(filterObservedIds));
  const page = state.mode === 'devices' ? deviceCollection(devices, state).page : collection.page;
  // Capture at most one page after source names are available (or failed).
  // Same-intent render/refetch/cache updates never replenish this demand.
  if (!sources.isPending && demand?.state !== state) {
    setDemand({ state, ids: state.mode === 'devices' ? [] : collection.rows.map(group => group.id) });
  }
  const demandIds = demand?.state === state ? demand.ids : [];
  const facts = useRecordingHomeFacts(scope, groups, demandIds);
  const home = useRecordingHomeEvidence({ scope, owners: facts.owners, detailGroupIds: demandIds, sources: sources.isFetching || sources.isPending
    ? { state: 'loading', scope } : sources.isError ? { state: 'error', scope } : sources.data });
  const summaries = Object.fromEntries(home.guidance.recordingItems.map(item => [item.id, recordingHomeResultSummary(item.evidence, {
    runtime: facts.runtimeGroups[item.id], history: facts.histories[item.id], owner: homeOwners(scope, item.evidence, facts.owners[item.id]).inspection,
  })]));
  const guidance = recordingHomeResultGuidance(home.guidance, summaries);
  // Presentation only: the first visit shows one stable reading state until sources and the
  // captured page's reads settle (success or failure), instead of repainting per read batch.
  const [revealed, setRevealed] = useState(false);
  const settling = !revealed && (sources.isPending || sources.isFetching || demand?.state !== state || home.loading || facts.loading);
  if (!revealed && !settling) setRevealed(true);
  useEffect(() => {
    if (revealed) return;
    const timer = window.setTimeout(() => setRevealed(true), RECORDING_HOME_SETTLE_LIMIT_MS);
    return () => window.clearTimeout(timer);
  }, [revealed]);
  const observedIds = guidance.attentionItems.filter(item => item.reason !== 'evidence_unavailable').map(item => item.id);
  const unverified = (restored.observedIds ?? []).filter(id => groups.some(group => group.id === id) &&
    (cached.delivery[id]?.state !== 'available' || cached.lifecycle[id]?.state !== 'available' || !facts.runtimeGroups[id]));
  const displayedObservedIds = [...new Set([...observedIds, ...unverified])];
  return <RecordingHome settling={settling} mode={state.mode} onCreateDevice={() => onCreateDevice?.(workspaceId)} workspaceId={workspaceId} guidance={guidance} location={location} onNavigate={onNavigate}
    onConfigurePoints={id => { if (devices.some(device => device.id === id)) onOpenDevice?.(workspaceId, id, 'points', switchable); }}
    onRetry={() => { void Promise.all([reloadScope(), sources.refetch(), home.retry(), facts.retry()]); }}
    collection={<RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance}
      pendingDevices={facts.pendingDevices}
      summaries={summaries}
      observedIds={displayedObservedIds} filterObservedIds={filterObservedIds} unverifiedObserved={unverified.length > 0} onRecheck={() => { void Promise.all([home.retryDetails(), facts.retry()]); }}
      sourcesAvailable={sources.data?.state === 'available' && !sources.isFetching} presentationReady={!settling} restored={restored}
      onChange={next => {
        if (!state.observedOnly || !next.observedOnly) setFilterObservedIds(displayedObservedIds);
        setState(next); const presentation = { state: next, workspaceRevision, scrollTop: 0, focusId: null, observedIds: displayedObservedIds };
        setRestored(presentation); rememberHomePresentation(context, workspaceId, presentation);
      }}
      onRemember={(scrollTop, focusId, currentPage) => rememberHomePresentation(context, workspaceId,
        { state: { ...state, page: currentPage ?? page }, workspaceRevision, scrollTop, focusId, observedIds: displayedObservedIds })}
      onOpen={id => {
        const item = home.guidance.recordingItems.find(candidate => candidate.id === id);
        if (item) {
          const target = item.evidence.group.status === 'deleted'
            ? buildStudioPresentationTarget('?step=output&target=database', location.hash, 'advanced')
            : buildSimpleRecordingTarget({ kind: 'edit', workspaceId, groupId: id, groupRevision: item.evidence.group.revision }, location.hash);
          if (target) onNavigate?.(target);
        }
      }}
      onOpenDevice={(id, kind) => {
        if (devices.some(device => device.id === id)) onOpenDevice?.(workspaceId, id, kind, switchable);
      }} />} />;
}
/** Home mounts no autosave/activation owners. Bootstrap identity and DB setup
 * revision are read independently, with an explicit empty DB revision allowed. */
export function RecordingHomeRoute({ location, onNavigate, onOpenDevice, onWorkspaceChange, onCreateDevice }: { location: RecordingHomeLocation; onNavigate?: (target: string) => void;
  onOpenDevice?: (workspaceId: string, deviceId: string, kind: 'connection' | 'points', devices?: readonly HomeSourceDevice[]) => void; onCreateDevice?: (workspaceId: string) => void; onWorkspaceChange?: (id?: string) => void }) {
  const { t } = useTranslation('recording-home');
  const { t: editorText } = useTranslation('recording-editor');
  const workspace = useStudioV2WorkspaceQuery();
  useEffect(() => { onWorkspaceChange?.(workspace.data?.id); }, [workspace.data?.id, onWorkspaceChange]);
  const groups = useWriteGroupsQuery(workspace.isSuccess && Boolean(workspace.data?.id));
  const reload = () => Promise.all([workspace.refetch(), groups.refetch()]);
  const failed = workspace.isError || groups.isError;
  const loading = workspace.isFetching || workspace.isPending || groups.isFetching || groups.isPending;
  const known = !failed && !loading && workspace.data?.id && groups.data?.workspace_id === workspace.data.id &&
    typeof groups.data.workspace_revision === 'string';
  // Same frame as the scoped home (product header + page header) so the hand-over does not jump.
  const settling = loading && !failed;
  if (!known) return <><RecordingProductShell workspaceId={workspace.data?.id} scopePending={settling && !workspace.data?.id}
    tools={settling ? <RecordingHomeTools location={location} onNavigate={onNavigate} /> : undefined} />
  <main data-workbench-v2="true" data-testid="recording-home" className="gw-recording-home"
    aria-busy={settling || undefined} data-home-settling={settling ? 'true' : 'false'}>
    <header className="gw-home-header"><h1>{t('title')}</h1>
      {settling && <button type="button" className="gw-home-create" onClick={() => { const target = buildSimpleRecordingTarget({ kind: 'new' }, location.hash); if (target) onNavigate?.(target); }}>{editorText('simple.create')}</button>}</header>
    <section className="gw-home-guidance" aria-live="polite" data-guidance-state={failed ? 'error' : 'unknown'}>
      <p className={settling ? 'gw-guidance-heading gw-home-settling-status' : 'gw-guidance-heading'} role={failed ? 'alert' : 'status'}>{t(failed ? 'error' : loading ? 'loading' : 'unknown')}</p>
      <div className="gw-home-actions"><button type="button" className="gw-primary" onClick={() => { void reload(); }}>{t('retry')}</button>
      <button type="button" onClick={() => { const target = buildStudioPresentationTarget(location.search, location.hash, 'advanced'); if (target) onNavigate?.(target); }}>{t('advanced')}</button></div>
    </section>
  </main></>;
  return <ScopedRecordingHome key={JSON.stringify([workspace.data.id, groups.data.workspace_revision])}
    workspaceId={workspace.data.id} workspaceRevision={groups.data.workspace_revision} list={groups.data}
    location={location} onNavigate={onNavigate} reloadScope={reload} onOpenDevice={onOpenDevice} onCreateDevice={onCreateDevice} />;
}
