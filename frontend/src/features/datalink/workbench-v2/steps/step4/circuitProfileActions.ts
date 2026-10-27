import type { MeasurementDefinition } from '../../../../../types/measurement';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupDestinationDraft,
  WriteGroupDraft,
  WriteGroupMutationResponse,
  WriteGroupProfileKind,
  WriteGroupSampleReduction,
  WriteGroupStorageMode,
  WriteGroupUpdateRequest,
  WriteGroupValueRole,
} from '../../../../../types/studioV2WriteGroup';
import {
  buildCircuitProfileGroupDraft,
  canonicalProfileGroupMatchesDraft,
  groupCircuitId,
  groupProfileKind,
  purposeOfLayout,
  type CircuitPurpose,
} from '../../state/circuitProfileBasic';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { cloneWriteGroupExactRecordingSource } from '../../../../../utils/studioV2ExactRecordingSourceJson';

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
  /** The layout a new circuit is created with (long when omitted); a saved circuit always keeps its own. */
  purpose?: CircuitPurpose;
  /** The roles the operator picked for instantaneous sources, by source key. */
  memberRoles?: Record<string, WriteGroupValueRole>;
  /** The reduction for a new instantaneous profile; saved profiles retain their persisted reduction. */
  newReduction?: WriteGroupSampleReduction;
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

/**
 * Carries the reduction the saved profile already holds into the update draft. The
 * builder only knows the reviewed kind and period, so a period-only save would
 * otherwise drop an explicit mean back to the omitted latest default. A profile
 * saved without a reduction gains none.
 */
function carrySavedReduction(draft: WriteGroupDraft, existing: WriteGroup): WriteGroupDraft {
  const saved = existing.recording_profile;
  if (saved?.sample_reduction === undefined || saved.kind !== draft.recording_profile?.kind) return draft;
  return { ...draft, recording_profile: { ...draft.recording_profile, sample_reduction: saved.sample_reduction } };
}

/**
 * Carries the persisted role and value types of the same saved source into the update
 * draft, matched by `(device_id, point_id, tag_id)` only. A newly added source gets no
 * guessed role or type, and a legacy member keeps exactly what it already had. A
 * confirmed binding is restored only for `counter_raw`: that role stores its latest raw
 * value, and the binding is what confirms its counter meaning. The instantaneous
 * builder never resolves one.
 */
function carrySavedMemberMetadata(draft: WriteGroupDraft, existing: WriteGroup): WriteGroupDraft {
  const savedBySource = new Map(
    existing.members.map((member) => [`${member.device_id}\0${member.point_id}\0${member.tag_id}`, member]),
  );
  const instantaneous = draft.recording_profile?.kind === 'instantaneous';
  return {
    ...draft,
    members: draft.members.map((member) => {
      const saved = savedBySource.get(`${member.device_id}\0${member.point_id}\0${member.tag_id}`);
      if (!saved) return member;
      if (saved.value_role === undefined && saved.source_value_type === undefined && saved.result_value_type === undefined &&
        saved.exact_recording_source === undefined) {
        return member;
      }
      const savedMeasurementId = saved.exact_recording_source?.measurement_id ??
        (instantaneous && saved.value_role === 'counter_raw' ? saved.measurement_id : undefined);
      return {
        ...member,
        ...(saved.value_role !== undefined ? { value_role: saved.value_role } : {}),
        ...(saved.source_value_type !== undefined ? { source_value_type: saved.source_value_type } : {}),
        ...(saved.result_value_type !== undefined ? { result_value_type: saved.result_value_type } : {}),
        ...(savedMeasurementId !== undefined ? { measurement_id: savedMeasurementId } : {}),
        ...(saved.exact_recording_source !== undefined ? { exact_recording_source: cloneWriteGroupExactRecordingSource(saved.exact_recording_source) } : {}),
      };
    }),
  };
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
    // A saved profile keeps its layout and what it already says about its values: a latest profile stays latest, and
    // only a profile created now is an explicit mean. The layout is never changed here, so no second owner appears.
    const reduction = existing ? existing.recording_profile?.sample_reduction ?? 'latest' : input.newReduction ?? 'mean';

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
      purpose: (existing && purposeOfLayout(existing.layout_version)) || input.purpose || 'long',
      reduction,
      memberRoles: input.memberRoles,
    });

    if (existing) {
      const reviewedDraft = carrySavedMemberMetadata(carrySavedReduction(draft, existing), existing);
      if (canonicalProfileGroupMatchesDraft(existing, reviewedDraft, input.connectorKind)) {
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
          group: reviewedDraft,
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
