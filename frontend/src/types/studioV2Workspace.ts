import type { StudioV2WorkspaceReadinessSummary } from './studioV2WorkspaceReadiness';

export type StudioV2WorkspaceKind = 'single';
export type StudioV2WorkspaceStatus = 'empty' | 'ready';

export interface StudioV2Workspace {
  id: string;
  kind: StudioV2WorkspaceKind;
  status: StudioV2WorkspaceStatus;
  ordered_device_ids: string[];
  readiness_summary?: StudioV2WorkspaceReadinessSummary | null;
  created_at: string;
  updated_at: string;
}
