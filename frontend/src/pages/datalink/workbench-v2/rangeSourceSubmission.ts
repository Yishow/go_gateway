import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import type { Rule } from '../../../features/datalink/workbench-v2/state/types';
import type { GroupCandidates } from '../../../features/datalink/workbench-v2/state/writeGroup/candidates';
import type { CreateSourceRuleRequest, SourceRuleRecord } from '../../../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../../../types/studioV2RuntimeApply';
import type { StudioV2MappingSaveResult, StudioV2WorkspaceMappingRequest } from '../../../services/studioV2Mappings';

export type RangeSourceBuffer = Omit<Rule, 'id' | 'device_id' | 'workspace_id' | 'revision_id' | 'persisted' | 'save_state' | 'save_error' | 'mc_persisted_points'>;
export interface RangeSourceRequest {
  workspaceId: string; sessionId: string; requestId: string; generation: number;
  deviceId: string; buffer: RangeSourceBuffer;
}
export interface RangeRuleReceipt {
  submitted: CreateSourceRuleRequest;
  acknowledged: StudioV2RuntimeAppliedRecord<SourceRuleRecord>;
}
export interface RangeMappingReceipt {
  submitted: StudioV2WorkspaceMappingRequest;
  acknowledged: StudioV2MappingSaveResult;
}
export interface SourceIdentity {
  device_id: string; point_id: string; tag_id: string;
}
export interface RangeSourcePreparation {
  device: StudioV2WorkspaceDeviceRecord;
  activationAllowed: boolean;
}
export interface RangeCompletion extends GroupCandidates {
  inactive: boolean;
  savedPoints: SourceIdentity[];
}
export interface RangeSourceOutcome {
  request: RangeSourceRequest;
  status: 'pending' | 'confirmed' | 'unknown' | 'blocked' | 'persisted-unavailable' | 'saved-inactive';
  rule?: SourceRuleRecord;
  candidates?: GroupCandidates;
  savedPoints?: SourceIdentity[];
  preparationError?: 'diagnostic-failed' | 'source-changed';
  /** A current-state read is never a receipt for an uncertain request. */
  inspection?: { status: 'not-checked' | 'current-state' | 'unavailable'; ruleId?: string; currentRule?: SourceRuleRecord };
}
export interface RangeSourcePort {
  /** Prevent partially observed new source points entering a recording draft. */
  canSelect?(deviceId: string, pointId: string): boolean;
  sourceWriteState?(deviceId: string, pointId: string): 'pending' | 'unknown' | undefined;
  prepare(deviceId: string): Promise<RangeSourcePreparation | undefined>;
  prepareSaved(sourceRequestId: string, request: RangeSourceRequest): Promise<RangeSourceOutcome>;
  submit(request: RangeSourceRequest): Promise<RangeSourceOutcome>;
  inspect(requestId: string): Promise<RangeSourceOutcome | undefined>;
  observations(): readonly RangeSourceOutcome[];
}
