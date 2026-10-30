import { buildStudioPresentationTarget } from '../../../../services/studioRecordingView';
import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
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
function ScopedRecordingHome({ workspaceId, workspaceRevision, location, onNavigate, reloadScope }: {
  workspaceId: string; workspaceRevision: string; location: RecordingHomeLocation; onNavigate?: (target: string) => void; reloadScope: () => Promise<unknown>;
}) {
  const [scope] = useState<RecordingHomeScope>(() => ({ workspaceId, workspaceRevision, epoch: ++generation }));
  const sources = useQuery({ queryKey: [...studioV2WorkspaceKeys.all, 'recording-home-sources', scope.workspaceId, scope.workspaceRevision, scope.epoch],
    queryFn: () => readRecordingHomeSources(scope), retry: false, refetchOnWindowFocus: false, gcTime: 0 });
  const home = useRecordingHomeEvidence({ scope, sources: sources.isFetching || sources.isPending
    ? { state: 'loading', scope } : sources.isError ? { state: 'error', scope } : sources.data });
  return <RecordingHome guidance={home.guidance} location={location} onNavigate={onNavigate}
    onRetry={() => { void Promise.all([reloadScope(), sources.refetch(), home.retry()]); }} />;
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
    workspaceId={workspace.data.id} workspaceRevision={groups.data.workspace_revision}
    location={location} onNavigate={onNavigate} reloadScope={reload} />;
}
