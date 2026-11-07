import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { RecordingHomeGuidanceItem } from './recordingHomeGuidance';

/** Lifecycle facts must belong to the displayed workspace, recording and destination revisions. */
export function recordingHomeLifecycleMatches(group: RecordingReadonly<WriteGroup>, evidence?: RecordingHomeGuidanceItem): boolean {
  const lifecycle = evidence?.evidence.lifecycle;
  const fact = lifecycle?.state === 'available' ? lifecycle.value : undefined;
  const scope = evidence?.action.scope;
  return Boolean(fact && scope && fact.workspace_id === group.workspace_id && fact.workspace_id === scope.workspaceId &&
    fact.workspace_revision === scope.workspaceRevision && fact.group_id === group.id && fact.group_revision === group.revision &&
    fact.applied_revision === group.applied_revision && fact.connector_id === group.destination.connector_id &&
    fact.connector_revision === group.destination.connector_revision);
}
