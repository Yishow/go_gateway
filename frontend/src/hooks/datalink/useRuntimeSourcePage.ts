import { useQuery } from '@tanstack/react-query';
import type { EvidenceResult, EvidenceScope } from '../../types/runtimeRuleEvidence';
import type { RuntimeSourcePage, RuntimeSourcePageRequest } from '../../types/studioV2RuntimeContext';
import { studioV2RuntimeContextAPI } from '../../services/studioV2RuntimeContext';
import { studioV2WorkspaceKeys } from './keys';

/** One requested page only; no automatic detail/page fan-out or previous-scope data. */
export function useRuntimeSourcePage(scope: EvidenceScope | null, request: RuntimeSourcePageRequest | null) {
  const observed = scope === null ? null : { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id };
  const selected = request === null ? null : { kind: request.kind, query: request.query, limit: request.limit,
    cursor: request.cursor, rule_id: request.rule_id, rule_revision: request.rule_revision };
  return useQuery<EvidenceResult<RuntimeSourcePage>>({
    queryKey: [...studioV2WorkspaceKeys.runtimeContext(), 'source-page', observed?.workspace_id ?? null,
      observed?.workspace_revision ?? null, observed?.device_id ?? null, selected?.kind ?? null, selected?.query ?? null,
      selected?.limit ?? null, selected?.cursor ?? null, selected?.rule_id ?? null, selected?.rule_revision ?? null],
    queryFn: ({ signal }) => observed === null || selected === null
      ? Promise.resolve<EvidenceResult<RuntimeSourcePage>>({ ok: false, reason: 'unavailable', action: 'refresh' })
      : studioV2RuntimeContextAPI.getSourcePage(observed, selected, signal),
    enabled: observed !== null && selected !== null,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
