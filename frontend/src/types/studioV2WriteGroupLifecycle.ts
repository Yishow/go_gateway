import type { WriteGroup, WriteGroupStatus } from './studioV2WriteGroup';
import type { CutoverMemberRef } from './studioV2WriteGroupCutover';

export type WriteGroupLifecycleActionName = 'edit' | 'stop' | 'remove' | 'recreate' | 'transition';
export interface SavedRecordingActionIntent {
  action: WriteGroupLifecycleActionName;
  workspaceId: string;
  groupId: string;
  workspaceRevision: string;
  groupRevision: string;
}
export interface WriteGroupLifecycleAction {
  action: WriteGroupLifecycleActionName;
  allowed: boolean;
  reason_code: string;
  next_action: string;
  requires_review: boolean;
}
export interface WriteGroupLifecycleDestination {
  connector_id: string;
  connector_revision: string;
  destination_scope: string;
}
export interface WriteGroupLifecycleGroupKey {
  workspace_id: string;
  group_id: string;
  group_revision: string;
}
export interface WriteGroupLifecyclePending {
  samples: number;
  outbox: {
    outstanding: number;
    total: number;
    queued: number;
    retrying: number;
    blocked: number;
    quarantined: number;
    unknown: number;
    sql_committed: number;
    skipped: number;
  };
  handoffs: Array<{
    origin: WriteGroupLifecycleGroupKey;
    source: WriteGroupLifecycleGroupKey;
    target: WriteGroupLifecycleGroupKey;
    state: string;
    samples: number;
  }>;
}
export interface WriteGroupLifecycleRevision {
  group_revision: string;
  effective_at: string;
  ends_at?: string;
  sources: CutoverMemberRef[];
  destination: WriteGroupLifecycleDestination;
  pending: WriteGroupLifecyclePending;
}
export interface WriteGroupLifecycleGeneration {
  group_id: string;
  generation?: number;
  predecessor_group_id?: string;
  current: boolean;
  reserved: boolean;
  status: WriteGroupStatus;
  group_revision: string;
  applied_revision: string;
  saved_destination: WriteGroupLifecycleDestination;
  revisions: WriteGroupLifecycleRevision[];
}
export interface WriteGroupLifecycleHistory {
  canonical: boolean;
  lineage_id?: string;
  lineage_revision?: number;
  current_group_id: string;
  reserved_group_id?: string;
  generations: WriteGroupLifecycleGeneration[];
}
/** Advisory source semantics only; a seed is never authority for reservation or DDL. */
export interface WriteGroupSuccessorSeed {
  basis_kind: 'applied' | 'saved';
  basis_revision: string;
  basis_digest: string;
  group: WriteGroup;
}
export interface WriteGroupLifecycle {
  workspace_id: string;
  workspace_revision: string;
  group_id: string;
  group_revision: string;
  applied_revision: string;
  status: WriteGroupStatus;
  connector_id: string;
  connector_revision: string;
  actions: WriteGroupLifecycleAction[];
  history: WriteGroupLifecycleHistory;
  successor_seed?: WriteGroupSuccessorSeed;
}
