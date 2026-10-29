import type { ReactNode } from 'react';
import type { EvidenceResult, RuleEvidenceView } from '../../../../../types/runtimeRuleEvidence';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { CircuitProfileLocalDraft } from '../../state/circuitProfileBasic';
import type { WorkbenchV2State } from '../../state/types';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

/** Viewing identity only. Reporting it never changes the existing editor's selection or draft. */
export interface RecordingViewedDevice { workspaceId: string; deviceId: string }
/** Transient display of the owner's current controls, never a save/start request or effective facts. */
export interface RecordingDraftSummary {
  profiles: readonly {
    kind: 'instantaneous' | 'cumulative'; interval: string; unit: 'second' | 'minute';
    sources: readonly GroupCandidate[]; reduction: 'latest' | 'mean' | 'cumulative_delta';
  }[];
  layout: 'wide' | 'long'; storageMode: 'plain' | 'timescale' | 'native_monthly';
  incompletePolicy: 'skip_row' | 'partial'; circuitName?: string;
}
export interface RecordingScopeView extends RecordingViewedDevice {
  workspaceRevision: string; circuitId?: string; groupIds: readonly string[]; draft: RecordingDraftSummary;
}
export type RenderRecordingScopeSummary = (scope: RecordingScopeView) => ReactNode;
export interface RecordingScopeSummaryProps {
  scope: RecordingScopeView;
  state: WorkbenchV2State;
  groups: readonly WriteGroup[];
  configurationState: 'ready' | 'loading' | 'error' | 'unknown';
  evidence: {
    data?: EvidenceResult<RuleEvidenceView>; isLoading: boolean; isFetching: boolean; isError: boolean;
    refetch: () => Promise<unknown>;
  };
}

/** Project controls without validating, defaulting or mutating their canonical owner. Invalid input remains visible. */
export function circuitRecordingScope(input: {
  workspaceId?: string; workspaceRevision: string; deviceId: string; circuitId: string;
  groups: readonly WriteGroup[]; draft: CircuitProfileLocalDraft;
  instantaneous: readonly GroupCandidate[]; cumulative: readonly GroupCandidate[];
  purpose: 'wide' | 'long'; storageMode: RecordingDraftSummary['storageMode']; reduction: 'latest' | 'mean';
}): RecordingScopeView {
  const profiles: RecordingDraftSummary['profiles'][number][] = [];
  if (input.draft.instantaneous_enabled) profiles.push({ kind: 'instantaneous',
    interval: input.draft.instantaneous_interval_input, unit: 'second', sources: input.instantaneous, reduction: input.reduction });
  if (input.draft.cumulative_enabled) profiles.push({ kind: 'cumulative',
    interval: input.draft.cumulative_interval_input, unit: 'minute', sources: input.cumulative, reduction: 'cumulative_delta' });
  return { workspaceId: input.workspaceId ?? '', workspaceRevision: input.workspaceRevision,
    deviceId: input.deviceId, ...(input.circuitId ? { circuitId: input.circuitId } : {}),
    groupIds: input.groups.map(group => group.id), draft: { profiles, layout: input.purpose,
      storageMode: input.storageMode, incompletePolicy: 'partial', circuitName: input.draft.circuit_name } };
}
