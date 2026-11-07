import { recordingHomeResultGuidance } from './recordingHomeResultGuidance';
import { useRecordingHomeFacts } from './useRecordingHomeFacts';
import { homeOwners } from './recordingHomeGuidanceEvidence';
import { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { cachedHomeReads } from './recordingHomeCachedReads';
import type { WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import { recordingHomeProjection } from './recordingHomeProjection';
import { RecordingHomeCollection } from './RecordingHomeCollection';
import { recordingCollection, type HomeDeviceSummary, type RecordingCollectionState } from './recordingHomeCollections';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { ReactNode } from 'react';
import { readHomePresentation, rememberHomePresentation } from './recordingHomePresentation';
import { buildSimpleRecordingTarget, buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { useEffect, useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useStudioV2WorkspaceQuery } from '../../../../hooks/datalink/useStudioV2Workspace';
import { useWriteGroupsQuery } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { useAcquisitionAnomalySummaryQuery } from '../../../../hooks/datalink/useAcquisitionAnomalies';
import { anomalySummaryStatus } from './AcquisitionAnomalyButton';
import { AcquisitionAnomalyDialog } from './AcquisitionAnomalyDialog';
import { useRecordingHomeSources, recordingHomeSourceContext } from '../../../../hooks/datalink/useRecordingHomeSources';
import { useRecordingHomeEvidence } from './useRecordingHomeEvidence';
import { useRecordingHomeDeviceStatus } from './useRecordingHomeDeviceStatus';
import { RecordingHome, RecordingHomeTools, type RecordingHomeLocation } from './RecordingHome';
import { deriveDeviceConnectionState } from './deviceConnectionState';
import { RecordingHomeDeviceStrip } from './RecordingHomeDeviceStrip';
import { RecordingHomeStatusSummary } from './RecordingHomeStatusSummary';
import { RecordingHomeDeviceActions } from './RecordingHomeDeviceActions';
import { RecordingDeleteDialogController } from './RecordingDeleteDialogController';
import type { DeviceCollector, DeviceConnectionObservation } from '../../../../types/deviceConnection';
import { RecordingProductShell } from '../shell/RecordingProductShell';
import type { RecordingHomeScope } from './recordingHomeProjection';
import type { HomeSourceDevice } from './RecordingHomeSourceHost';

// Request generation only, not configuration or workspace authority. Remounting
// cannot reuse pending evidence from a previous visit, even A -> B -> A.
let generation = 0;
const clientTokens = new WeakMap<object, number>();
let clientGeneration = 0;
/** Upper bound for holding the first-load reading state; later reads render progressively as before. */
export const RECORDING_HOME_SETTLE_LIMIT_MS = 2500;
function ScopedRecordingHome({ workspaceId, workspaceRevision, list, location, onNavigate, reloadScope, onOpenDevice, onCreateDevice, onDelete, onDeleteDevice, notice }: {
  workspaceId: string; workspaceRevision: string; list: WriteGroupListResponse; location: RecordingHomeLocation; onNavigate?: (target: string) => void; reloadScope: () => Promise<unknown>; onOpenDevice?: (workspaceId: string, deviceId: string, kind: 'connection' | 'points', devices?: readonly HomeSourceDevice[]) => void; onCreateDevice?: (workspaceId: string) => void;
  onDelete: (group: RecordingReadonly<WriteGroup>) => void; onDeleteDevice: (device: HomeDeviceSummary) => void; notice?: ReactNode;
}) {
  const context = useQueryClient();
  const homeRef = useRef<HTMLElement>(null);
  const { t } = useTranslation('recording-home');
  const [presentationContext] = useState(() => recordingHomeSourceContext());
  const [restored, setRestored] = useState(() => readHomePresentation(context, workspaceId, workspaceRevision, presentationContext));
  // Devices are managed in the strip; an older devices-mode presentation reopens on the record list.
  const [state, setState] = useState<RecordingCollectionState>(() => ({ ...restored.state, mode: 'recordings',
    ...(restored.state.mode === 'devices' ? { keyword: '', page: 1 } : {}) }));
  const [anomalyDevice, setAnomalyDevice] = useState<{ id: string; name: string } | null>(null);
  // Freeze membership for this user-selected page. Recovered rows stay visible
  // until another filter/page intent; replies cannot pull the next batch in.
  const [filterObservedIds, setFilterObservedIds] = useState<readonly string[]>(restored.observedIds ?? []);
  const [demand, setDemand] = useState<{ state: typeof state; ids: string[] } | null>(null);
  const [scope] = useState<RecordingHomeScope>(() => ({ workspaceId, workspaceRevision, epoch: ++generation }));
  const status = useRecordingHomeDeviceStatus(workspaceId);
  const sources = useRecordingHomeSources(scope);
  const groupsRead = recordingHomeProjection({ scope, groups: { state: 'available', scope, value: list } }).groups;
  const groups = groupsRead.state === 'available' ? groupsRead.value.map(item => item.group) : [];
  const devices = sources.observed?.devices ?? [];
  const collectors = (status.snapshot?.collectors ?? []) as DeviceCollector[];
  const observations = new Map<string, DeviceConnectionObservation>(devices.map(device => {
    const collector = collectors.find(item => item.device_id === device.id);
    return [device.id, { state: deriveDeviceConnectionState(collector, status.now), collector }];
  }));
  // Last observed presentation remains visible; only the current visit can provide fresh evidence.
  // Every device chosen in the settings frame is re-validated by its own scoped locator read.
  const switchable = sources.observed ? devices.map(({ id, name }) => ({ id, name })) : undefined;
  const cached = cachedHomeReads(context, groups.map(group => ({ ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision })));
  const collection = recordingCollection(groups, devices, state, new Set(filterObservedIds));
  const page = collection.page;
  // Capture at most one page after source names are available (or failed).
  // Same-intent render/refetch/cache updates never replenish this demand.
  if ((sources.observed || !sources.isPending) && demand?.state !== state) {
    setDemand({ state, ids: collection.rows.map(group => group.id) });
  }
  const demandIds = demand?.state === state ? demand.ids : [];
  const facts = useRecordingHomeFacts(scope, groups, demandIds);
  const anomalySummary = useAcquisitionAnomalySummaryQuery();
  const home = useRecordingHomeEvidence({ scope, owners: facts.owners, detailGroupIds: demandIds, sources: sources.evidence });
  const summaries = Object.fromEntries(home.guidance.recordingItems.map(item => [item.id, recordingHomeResultSummary(item.evidence, {
    runtime: facts.runtimeGroups[item.id], history: facts.histories[item.id], owner: homeOwners(scope, item.evidence, facts.owners[item.id]).inspection,
  })]));
  const guidance = recordingHomeResultGuidance(home.guidance, summaries);
  // Presentation only: the first visit shows one stable reading state until sources and the
  // captured page's reads settle (success or failure), instead of repainting per read batch.
  const [revealed, setRevealed] = useState(Boolean(sources.observed));
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
  const changeState = (next: RecordingCollectionState) => {
    if (!state.observedOnly || !next.observedOnly) setFilterObservedIds(displayedObservedIds);
    setState(next); const presentation = { state: next, workspaceRevision, scrollTop: 0, focusId: null, observedIds: displayedObservedIds };
    setRestored(presentation); rememberHomePresentation(context, workspaceId, presentation, presentationContext);
  };
  const openDevice = (id: string, kind: 'connection' | 'points') => {
    if (devices.some(device => device.id === id)) onOpenDevice?.(workspaceId, id, kind, switchable);
  };
  const remember = (scrollTop: number, focusId: string | null, currentPage = page) => rememberHomePresentation(context, workspaceId,
    { state: { ...state, page: currentPage }, workspaceRevision, scrollTop, focusId, observedIds: displayedObservedIds }, presentationContext);
  const renderDeviceActions = (device: HomeDeviceSummary) => <RecordingHomeDeviceActions device={device} available={sources.fresh && groupsRead.state === 'available'} onDelete={onDeleteDevice} />;
  return <RecordingHome homeRef={homeRef} settling={settling} onDelete={onDelete} workspaceId={workspaceId} guidance={guidance} location={location} onNavigate={onNavigate} updatedAt={sources.dataUpdatedAt || undefined}
    statusContent={<><RecordingHomeStatusSummary devices={devices} observations={observations} groups={groups} metrics={status.snapshot?.metrics} outdated={status.outdated} />
      <RecordingHomeDeviceStrip devices={devices} observations={observations} selected={state.deviceId} now={status.displayNow}
        renderDeviceActions={renderDeviceActions} onCreateDevice={onCreateDevice ? () => onCreateDevice(workspaceId) : undefined}
        anomalySummary={anomalySummaryStatus(anomalySummary)} onOpenAnomalies={device => setAnomalyDevice({ id: device.id, name: device.name })}
        onOpenDevice={onOpenDevice ? openDevice : undefined} returnFocusId={restored.focusId}
        onRemember={focusId => remember(homeRef.current?.querySelector<HTMLDivElement>('[data-testid="home-collection-scroll"]')?.scrollTop ?? 0, focusId)}
        onSelect={id => changeState({ ...state, mode: 'recordings', deviceId: state.deviceId === id ? '' : id, page: 1 })} /></>}
    onConfigurePoints={id => openDevice(id, 'points')}
    onRetry={() => { void Promise.all([reloadScope(), sources.refetch(), home.retry(), facts.retry(), status.refetch()]); }}
    collection={<>{notice}<p data-testid="home-source-observation" data-source-status={sources.status} role={sources.status === 'error' ? 'alert' : 'status'}>{t(`sources.${sources.status}`)}</p><RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance}
      pendingDevices={facts.pendingDevices}
      summaries={summaries}
      returnLayoutPending={sources.isFetching || home.loading || facts.loading || anomalySummary.isFetching}
      observedIds={displayedObservedIds} filterObservedIds={filterObservedIds} unverifiedObserved={unverified.length > 0} onRecheck={() => { void Promise.all([home.retryDetails(), facts.retry()]); }}
      sourcesAvailable={sources.fresh} observedSourcesAvailable={Boolean(sources.observed)} presentationReady={!settling} restored={restored}
      onDelete={onDelete} onChange={changeState}
      onRemember={remember}
      onOpen={id => {
        const item = home.guidance.recordingItems.find(candidate => candidate.id === id);
        if (item) {
          const target = item.evidence.group.status === 'deleted'
            ? buildStudioPresentationTarget('?step=output&target=database', location.hash, 'advanced')
            : buildSimpleRecordingTarget({ kind: 'edit', workspaceId, groupId: id, groupRevision: item.evidence.group.revision }, location.hash);
          if (target) onNavigate?.(target);
        }
      }}
      onOpenDevice={openDevice} />
      {anomalyDevice && <AcquisitionAnomalyDialog device={anomalyDevice} onClose={() => { setAnomalyDevice(null); void anomalySummary.refetch(); }} />}</>} />;
}
/** Home mounts no autosave/activation owners. Bootstrap identity and DB setup
 * revision are read independently, with an explicit empty DB revision allowed. */
export function RecordingHomeRoute({ location, onNavigate, onOpenDevice, onWorkspaceChange, onCreateDevice }: { location: RecordingHomeLocation; onNavigate?: (target: string) => void;
  onOpenDevice?: (workspaceId: string, deviceId: string, kind: 'connection' | 'points', devices?: readonly HomeSourceDevice[]) => void; onCreateDevice?: (workspaceId: string) => void; onWorkspaceChange?: (id?: string) => void }) {
  const { t } = useTranslation('recording-home');
  const { t: editorText } = useTranslation('recording-editor');
  const queryClient = useQueryClient();
  if (!clientTokens.has(queryClient)) clientTokens.set(queryClient, ++clientGeneration);
  const workspace = useStudioV2WorkspaceQuery(true, undefined, { refetchOnMount: 'always' });
  useEffect(() => { onWorkspaceChange?.(workspace.data?.id); }, [workspace.data?.id, onWorkspaceChange]);
  // The singleton workspace can change while away, even inside either cache stale window.
  // Revalidate both owners on each home visit; the identity check still rejects mismatched scopes.
  const groups = useWriteGroupsQuery(workspace.isSuccess && Boolean(workspace.data?.id), undefined, { refetchOnMount: 'always' });
  const reload = () => Promise.all([workspace.refetch(), groups.refetch()]);
  const failed = workspace.isError || groups.isError;
  const loading = workspace.isFetching || workspace.isPending || groups.isFetching || groups.isPending;
  // Keep the validated home mounted while refreshing the same scope. Its reads
  // already expose pending evidence; refetching is not a new workspace visit.
  const known = !failed && workspace.data?.id && groups.data?.workspace_id === workspace.data.id &&
    typeof groups.data.workspace_revision === 'string';
  // Same frame as the scoped home (product header + page header) so the hand-over does not jump.
  const settling = loading && !failed;
  const fallback = <><RecordingProductShell workspaceId={workspace.data?.id} scopePending={settling && !workspace.data?.id}
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
  if (!workspace.data?.id || !groups.data || groups.data.workspace_id !== workspace.data.id) return fallback;
  const backendContext = recordingHomeSourceContext();
  return <RecordingDeleteDialogController key={JSON.stringify([workspace.data.id, clientTokens.get(queryClient), backendContext.backendToken, backendContext.apiOrigin])}
    workspaceId={workspace.data.id} workspaceRevision={groups.data.workspace_revision} groups={groups.data.groups} location={location} onNavigate={onNavigate}
    refresh={async () => { const [currentWorkspace, currentGroups] = await reload(); return !currentWorkspace.isError && !currentGroups.isError; }}>
    {actions => !known ? fallback : <ScopedRecordingHome key={JSON.stringify([workspace.data!.id, groups.data!.workspace_revision, clientTokens.get(queryClient), backendContext.backendToken, backendContext.apiOrigin])}
      workspaceId={workspace.data!.id} workspaceRevision={groups.data!.workspace_revision} list={groups.data!}
      location={location} onNavigate={onNavigate} reloadScope={reload} onOpenDevice={onOpenDevice} onCreateDevice={onCreateDevice} {...actions} />}
  </RecordingDeleteDialogController>;
}
