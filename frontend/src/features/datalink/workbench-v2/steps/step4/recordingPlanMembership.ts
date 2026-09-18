import type { MeasurementDefinition } from '@/types/measurement';
import type { PlanMember, RecordingPlan } from '@/types/recordingPlan';
import type { DbConnector, WorkbenchV2State } from '../../state/types';

export interface RecordingPlanScope {
  workspaceId: string;
  members: PlanMember[];
  fingerprint: string;
}

/** Join saved measurements to the enabled points currently shown in Studio. */
export function resolveRecordingPlanScope(
  state: WorkbenchV2State,
  workspaceId: string,
  definitions: MeasurementDefinition[],
): RecordingPlanScope {
  const selected: MeasurementDefinition[] = [];
  const points = state.points.filter(point => point.enabled && !point.skipped);
  for (const point of points) {
    const device = state.devices.find(candidate => candidate.id === point.device_id);
    const mapping = state.mappings[point.id];
    const rule = state.rules.find(candidate => candidate.id === point.rule_id);
    if (rule?.enabled === false || mapping?.enabled === false) continue;
    if (!device?.persisted || device.save_state !== 'saved') return emptyScope(workspaceId);
    const pointId = mapping?.persisted_point_id ?? point.id;
    const matches = definitions.filter(definition => definition.workspace_id === workspaceId &&
      definition.device_id === point.device_id && definition.point_id === pointId &&
      definition.id && definition.equipment_id && definition.definition_revision && definition.source_binding_revision);
    if (matches.length !== 1) return emptyScope(workspaceId);
    selected.push(matches[0]);
  }
  if (new Set(selected.map(definition => definition.id)).size !== selected.length) return emptyScope(workspaceId);
  selected.sort((a, b) => a.id.localeCompare(b.id));
  return {
    workspaceId,
    members: selected.map(definition => ({
      member_id: `member-${definition.id}`, measurement_id: definition.id,
      equipment_id: definition.equipment_id, name: definition.name,
    })),
    fingerprint: JSON.stringify([workspaceId, selected.map(definition => [definition.id, definition.device_id,
      definition.point_id, definition.equipment_id, definition.definition_revision, definition.source_binding_revision, definition.series_epoch])]),
  };
}

function emptyScope(workspaceId: string): RecordingPlanScope {
  return { workspaceId, members: [], fingerprint: JSON.stringify([workspaceId, []]) };
}

/** A selected plan must describe the same saved members and destination. */
export function recordingPlanMatchesScope(plan: RecordingPlan, scope: RecordingPlanScope | undefined, connector: DbConnector | undefined): boolean {
  if (!scope?.members.length || plan.workspace_id !== scope.workspaceId || !connector?.connector_id ||
    !connector.identity_revision || !plan.members?.length || !plan.destinations?.length) return false;
  if (plan.destinations.some(destination => destination.connector_id !== connector.connector_id ||
    destination.connector_revision !== connector.identity_revision)) return false;
  if (plan.members.length !== scope.members.length || new Set(plan.members.map(member => member.measurement_id)).size !== plan.members.length) return false;
  return plan.members.every(member => scope.members.some(expected => expected.measurement_id === member.measurement_id &&
    expected.equipment_id === member.equipment_id)) && Boolean(plan.streams?.length) &&
    plan.streams.every(stream => plan.members.some(member => member.measurement_id === stream.measurement_id));
}
