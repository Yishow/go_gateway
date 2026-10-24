import type { MeasurementDefinition } from '../../../../../types/measurement';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupDestinationDraft,
  WriteGroupMutationResponse,
  WriteGroupProfileKind,
  WriteGroupStorageMode,
  WriteGroupUpdateRequest,
} from '../../../../../types/studioV2WriteGroup';
import {
  buildCircuitProfileGroupDraft,
  canonicalProfileGroupMatchesDraft,
  groupCircuitId,
  groupProfileKind,
} from '../../state/circuitProfileBasic';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

export interface EnsureOrUpdateCircuitProfilesInput {
  workspaceId: string;
  initialWorkspaceRevision: string;
  deviceName: string;
  enabledProfiles: Array<{
    profile: WriteGroupProfileKind;
    intervalValue: number;
    intervalSeconds: number;
    candidates?: GroupCandidate[];
  }>;
  storageMode: WriteGroupStorageMode;
  recreateMissing: boolean;
  circuitId: string;
  circuitName: string;
  destination: WriteGroupDestinationDraft;
  connectorKind?: string;
  selectedCandidates?: GroupCandidate[];
  measurementsByPointKey: Map<string, MeasurementDefinition>;
  savedLongGroups: WriteGroup[];
  preparedGroups: WriteGroup[];
  ensureCircuitProfile: (variables: { request: WriteGroupCreateRequest }) => Promise<WriteGroupMutationResponse>;
  updateWriteGroup: (variables: { id: string; request: WriteGroupUpdateRequest }) => Promise<WriteGroupMutationResponse>;
}

export interface EnsureOrUpdateCircuitProfilesResult {
  groups: WriteGroup[];
  workspaceRevision: string;
}

export async function ensureOrUpdateCircuitProfileGroups(
  input: EnsureOrUpdateCircuitProfilesInput,
): Promise<EnsureOrUpdateCircuitProfilesResult> {
  const {
    workspaceId,
    initialWorkspaceRevision,
    deviceName,
    enabledProfiles,
    storageMode,
    recreateMissing,
    circuitId,
    circuitName,
    destination,
    selectedCandidates,
    measurementsByPointKey,
    savedLongGroups,
    preparedGroups,
    ensureCircuitProfile,
    updateWriteGroup,
  } = input;

  let currentWorkspaceRevision = initialWorkspaceRevision;
  const resolvedGroups: WriteGroup[] = [];

  for (const { profile, intervalValue, intervalSeconds, candidates } of enabledProfiles) {
    const profileCandidates = candidates ?? selectedCandidates ?? [];
    const existing = [...preparedGroups, ...savedLongGroups].find((group) =>
      group.workspace_id === workspaceId && groupProfileKind(group) === profile &&
      groupCircuitId(group) === circuitId &&
      group.members.every((member) => profileCandidates.some((candidate) => candidate.device_id === member.device_id)),
    );

    const effectiveStorageMode: WriteGroupStorageMode = existing?.storage_mode ?? storageMode;

    const draft = buildCircuitProfileGroupDraft({
      workspaceId,
      deviceName,
      profile,
      intervalValue,
      intervalSeconds,
      storageMode: effectiveStorageMode,
      recreateMissing,
      circuitId,
      circuitName,
      destination,
      selectedCandidates: profileCandidates,
      measurementsByPointKey,
    });

    if (existing) {
      if (canonicalProfileGroupMatchesDraft(existing, draft, input.connectorKind)) {
        resolvedGroups.push(existing);
        continue;
      }
      const updated = await updateWriteGroup({
        id: existing.id,
        request: {
          workspace_id: workspaceId,
          expected_workspace_revision: currentWorkspaceRevision,
          expected_group_revision: existing.revision,
          expected_connector_revision: destination.connector_revision,
          group: draft,
        },
      });
      if (updated.workspace_revision) {
        currentWorkspaceRevision = updated.workspace_revision;
      }
      resolvedGroups.push(updated.group);
    } else {
      const ensured = await ensureCircuitProfile({
        request: {
          workspace_id: workspaceId,
          expected_workspace_revision: currentWorkspaceRevision,
          expected_connector_revision: destination.connector_revision,
          group: draft,
        },
      });
      if (ensured.workspace_revision) {
        currentWorkspaceRevision = ensured.workspace_revision;
      }
      resolvedGroups.push(ensured.group);
    }
  }

  return {
    groups: resolvedGroups,
    workspaceRevision: currentWorkspaceRevision,
  };
}

export { BatchPlanError, executeBatchPlan } from './circuitProfileBatchExecutor';
export type {
  ExecuteBatchPlanInput,
  ExecuteBatchPlanResult,
  FrozenBatchPlan,
  FrozenBatchTarget,
} from './circuitProfileBatchExecutor';
export type {
  BatchProfileResult,
  BatchProfileStatus,
  BatchTargetExecutionResult,
  BatchTargetStatus,
} from '../../state/circuitProfileBatchOutcome';
