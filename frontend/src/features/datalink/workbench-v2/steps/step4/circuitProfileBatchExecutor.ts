import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupMutationResponse,
  WriteGroupProfileKind,
  WriteGroupUpdateRequest,
} from '../../../../../types/studioV2WriteGroup';
import {
  BATCH_PROFILE_ORDER,
  type BatchProfileResults,
  type BatchTargetExecutionResult,
  classifyBatchMutationError,
  groupProfileInterval,
  summarizeBatchTarget,
} from '../../state/circuitProfileBatchOutcome';
import { circuitProfileLayoutBlockReason, type CircuitLayoutBlockReason, type CircuitProfileCapabilityContext } from '../../state/circuitProfileCapabilities';
import { existingProfileUpdateGroup, newProfileDraft } from '../../state/circuitProfileBatchIntent';

export interface FrozenBatchTarget {
  circuit_id: string;
  circuit_name: string;
  instantaneousGroup?: WriteGroup;
  cumulativeGroup?: WriteGroup;
  anyGroup: WriteGroup;
}

/** Everything the operator reviewed; the executor sends nothing that is not in here. */
export interface FrozenBatchPlan {
  workspaceId: string;
  workspaceRevision: string;
  connectorRevision?: string;
  /** Identity of the Basic scope (workspace, device, circuit, destination) that reviewed the plan. */
  scopeKey?: string;
  applyInst: boolean;
  instInterval: number;
  applyCum: boolean;
  cumInterval: number;
  targets: FrozenBatchTarget[];
}

export interface ExecuteBatchPlanInput {
  plan: FrozenBatchPlan;
  capabilityContext?: CircuitProfileCapabilityContext;
  connectorRevision?: string;
  updateWriteGroup: (vars: { id: string; request: WriteGroupUpdateRequest }) => Promise<WriteGroupMutationResponse>;
  ensureCircuitProfile?: (vars: { request: WriteGroupCreateRequest }) => Promise<WriteGroupMutationResponse>;
}

export interface ExecuteBatchPlanResult {
  updatedGroups: WriteGroup[];
  workspaceRevision: string;
  targetResults: Record<string, BatchTargetExecutionResult>;
}

/** Thrown when the batch stopped early; carries what was saved and the per profile evidence. */
export class BatchPlanError extends Error {
  constructor(cause: unknown, readonly updatedGroups: WriteGroup[], readonly workspaceRevision: string,
    readonly targetResults: Record<string, BatchTargetExecutionResult>) {
    super('Circuit profile batch stopped before every reviewed change was confirmed.', { cause });
    this.name = 'BatchPlanError';
  }
}

const UNAVAILABLE_WRITER = new Error('Creating a circuit profile is not available for this batch.');

function requestedKinds(plan: Pick<FrozenBatchPlan, 'applyInst' | 'applyCum'>): WriteGroupProfileKind[] {
  return BATCH_PROFILE_ORDER.filter((kind) => (kind === 'instantaneous' ? plan.applyInst : plan.applyCum));
}

export type BatchLayoutBlockReason = CircuitLayoutBlockReason | 'batch_layout_unsupported';

export class BatchLayoutError extends Error {
  constructor(readonly reason: BatchLayoutBlockReason) {
    super('Recording layout is unavailable for this batch.');
    this.name = 'BatchLayoutError';
  }
}

/** The current authority is separate from the frozen review; every selected profile must match it. */
export function batchLayoutBlockReason(
  plan: Pick<FrozenBatchPlan, 'targets' | 'applyInst' | 'applyCum' | 'connectorRevision'>,
  context?: CircuitProfileCapabilityContext, connectorRevision?: string,
): BatchLayoutBlockReason | undefined {
  const destination = context?.destination;
  if (!context || !destination?.connector_id || !destination.connector_revision) return 'capability_unavailable';
  for (const target of plan.targets) {
    for (const kind of requestedKinds(plan)) {
      const existing = kind === 'instantaneous' ? target.instantaneousGroup : target.cumulativeGroup;
      const source = existing ?? target.anyGroup;
      const revision = plan.connectorRevision ?? connectorRevision ?? source.destination.connector_revision;
      if (source.destination.connector_id !== destination.connector_id ||
        source.destination.connector_revision !== destination.connector_revision || revision !== destination.connector_revision) {
        return 'capability_unavailable';
      }
      if (source.layout_version !== 'managed_long_v1') return 'batch_layout_unsupported';
      const reason = circuitProfileLayoutBlockReason(context, 'long', source.storage_mode ?? 'plain');
      if (reason) return reason;
    }
  }
  return undefined;
}

export async function executeBatchPlan(input: ExecuteBatchPlanInput): Promise<ExecuteBatchPlanResult> {
  const { plan, connectorRevision, updateWriteGroup, ensureCircuitProfile } = input;
  const blockedReason = batchLayoutBlockReason(plan, input.capabilityContext, connectorRevision);
  if (blockedReason) throw new BatchLayoutError(blockedReason);
  const kinds = requestedKinds(plan);
  const intervals = { instantaneous: plan.instInterval, cumulative: plan.cumInterval };
  let currentWsRev = plan.workspaceRevision;
  const updatedGroups: WriteGroup[] = [];
  const targetResults: Record<string, BatchTargetExecutionResult> = {};
  for (const target of plan.targets) {
    const profiles: BatchProfileResults = {};
    for (const kind of kinds) profiles[kind] = { kind, status: 'not-attempted' };
    targetResults[target.circuit_id] = { circuit_id: target.circuit_id, circuit_name: target.circuit_name, status: 'not-attempted', profiles };
  }

  const write = async (target: FrozenBatchTarget, kind: WriteGroupProfileKind, existing: WriteGroup | undefined) => {
    const revision = plan.connectorRevision ?? connectorRevision ?? (existing ?? target.anyGroup).destination.connector_revision ?? '';
    if (existing) {
      return updateWriteGroup({
        id: existing.id,
        request: {
          workspace_id: plan.workspaceId, expected_workspace_revision: currentWsRev,
          expected_group_revision: existing.revision, expected_connector_revision: revision,
          group: existingProfileUpdateGroup(existing, kind, intervals[kind]),
        },
      });
    }
    if (!ensureCircuitProfile) throw UNAVAILABLE_WRITER;
    return ensureCircuitProfile({
      request: {
        workspace_id: plan.workspaceId, expected_workspace_revision: currentWsRev, expected_connector_revision: revision,
        group: newProfileDraft(plan.workspaceId, target, kind, intervals[kind]),
      },
    });
  };

  let stoppedBy: { cause: unknown } | undefined;
  for (const target of plan.targets) {
    for (const kind of kinds) {
      if (stoppedBy) break;
      const result = targetResults[target.circuit_id].profiles[kind];
      if (!result) continue;
      const existing = kind === 'instantaneous' ? target.instantaneousGroup : target.cumulativeGroup;
      if (existing && groupProfileInterval(existing, kind) === intervals[kind]) {
        result.status = 'unchanged';
        continue;
      }
      try {
        const saved = await write(target, kind, existing);
        if (saved.workspace_revision) currentWsRev = saved.workspace_revision;
        updatedGroups.push(saved.group);
        result.status = 'saved';
        result.group = saved.group;
      } catch (cause) {
        result.status = classifyBatchMutationError(cause);
        result.cause = cause;
        stoppedBy = { cause };
      }
    }
    targetResults[target.circuit_id].status = summarizeBatchTarget(targetResults[target.circuit_id].profiles);
  }
  if (stoppedBy) throw new BatchPlanError(stoppedBy.cause, updatedGroups, currentWsRev, targetResults);
  return { updatedGroups, workspaceRevision: currentWsRev, targetResults };
}
