import { useQuery, useQueryClient } from '@tanstack/react-query';
import type { EvidenceResult, EvidenceScope, RuleEvidenceView } from '../../types/runtimeRuleEvidence';
import { studioV2RuntimeContextAPI } from '../../services/studioV2RuntimeContext';
import { studioV2WorkspaceKeys } from './keys';

/** Shared read-only evidence; never promotes editable setup data to applied facts. */
export function useRuntimeRuleEvidence(scope: EvidenceScope | null) {
  const client = useQueryClient();
  const selected = scope === null ? null : { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id };
  return useQuery<EvidenceResult<RuleEvidenceView>>({
    queryKey: [...studioV2WorkspaceKeys.runtimeContext(), 'rule-evidence', selected?.workspace_id ?? null, selected?.workspace_revision ?? null, selected?.device_id ?? null],
    queryFn: async ({ signal }) => {
      const result: EvidenceResult<RuleEvidenceView> = selected === null
        ? { ok: false, reason: 'unavailable', action: 'refresh' }
        : await studioV2RuntimeContextAPI.getRuleEvidence(selected, signal);
      if (!signal.aborted && !result.ok && result.reason === 'stale-reference') {
        // Refresh only the scope authority. Its new revision keys the next evidence
        // read; publish the stale result while that refresh is still pending.
        void client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.runtimeContext(), exact: true });
      }
      return result;
    },
    enabled: selected !== null,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
