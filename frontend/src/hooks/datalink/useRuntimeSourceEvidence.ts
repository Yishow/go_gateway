import { useQuery } from '@tanstack/react-query';
import type { EvidenceResult, EvidenceScope, RuleEvidenceView } from '../../types/runtimeRuleEvidence';
import type { RuntimeSourceSelection } from '../../types/studioV2RuntimeContext';
import { studioV2RuntimeContextAPI } from '../../services/studioV2RuntimeContext';
import { studioV2WorkspaceKeys } from './keys';

/** Single selected source, scoped independently from recording-only evidence. */
export function useRuntimeSourceEvidence(scope: EvidenceScope | null, selection: RuntimeSourceSelection | null) {
  const observed = scope === null ? null : { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id };
  const selected = selection === null ? null : {
    rule_id: selection.rule_id, rule_revision: selection.rule_revision, point_id: selection.point_id,
    ...(selection.mapping_id === undefined ? {} : { mapping_id: selection.mapping_id }),
  };
  return useQuery<EvidenceResult<RuleEvidenceView>>({
    queryKey: [...studioV2WorkspaceKeys.runtimeContext(), 'source-evidence', observed?.workspace_id ?? null,
      observed?.workspace_revision ?? null, observed?.device_id ?? null, selected?.rule_id ?? null,
      selected?.rule_revision ?? null, selected?.point_id ?? null, selected?.mapping_id ?? null],
    queryFn: ({ signal }) => observed === null || selected === null
      ? Promise.resolve<EvidenceResult<RuleEvidenceView>>({ ok: false, reason: 'unavailable', action: 'refresh' })
      : studioV2RuntimeContextAPI.getSourceEvidence(observed, selected, signal),
    enabled: observed !== null && selected !== null,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
