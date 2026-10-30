import type { EvidenceResult, RecordingGroupView, RuleEvidenceView } from '../../../../types/runtimeRuleEvidence';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { RecordingHomeEvidence, RecordingHomeScope } from './recordingHomeProjection';
/** Existing runtime service validates snapshot/layout/boundary. This adapter only
 * fences its captured request and canonical current group; it creates no proof. */
export function recordingHomeRuntimeGroup(scope: RecordingHomeScope, group: RecordingReadonly<WriteGroup>,
  evidence?: RecordingHomeEvidence<EvidenceResult<RuleEvidenceView>>,
): RecordingReadonly<RecordingGroupView> | undefined {
  if (!evidence || evidence.state !== 'available' || !evidence.value.ok || evidence.scope.workspaceId !== scope.workspaceId ||
    evidence.scope.workspaceRevision !== scope.workspaceRevision || evidence.scope.epoch !== scope.epoch) return;
  const view = evidence.value.value;
  if (view.workspace_id !== scope.workspaceId || view.workspace_revision !== scope.workspaceRevision || view.freshness !== 'fresh' ||
    !view.observed_at || !group.members.some(member => member.device_id === view.device_id) && group.basic_managed_device_id !== view.device_id) return;
  const matches = view.groups.filter(item => item.group_id === group.id);
  if (matches.length !== 1) return;
  const result = matches[0];
  if (result.saved_revision !== group.revision || result.applied_revision !== group.applied_revision ||
    result.state !== 'effective' || result.proof !== 'applied-snapshot-and-layout' || !result.effective || !result.effective_revision ||
    result.effective.binding.entity_revision !== result.effective_revision || result.as_of !== view.observed_at) return;
  return result;
}
