import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { UseQueryResult } from '@tanstack/react-query';
import { useRuntimeRuleEvidence } from '../../../../../hooks/datalink/useRuntimeRuleEvidence';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import type { EvidenceScope } from '../../../../../types/runtimeRuleEvidence';
import type { WriteGroupListResponse } from '../../../../../types/studioV2WriteGroup';
import type { RecordingScopeSummaryProps } from './recordingScopeModel';

const groupsScopeAuthority = studioV2WorkspaceKeys.writeGroups();
/** Step4's authority is the groups read, not the dashboard's runtime-context read.
 * Refresh only observes scope; it cannot acknowledge editable draft/CAS state. */
export function useStep4RecordingEvidence(scope: EvidenceScope | null,
  groups: Pick<UseQueryResult<WriteGroupListResponse>, 'refetch'>): RecordingScopeSummaryProps['evidence'] {
  const evidence = useRuntimeRuleEvidence(scope, groupsScopeAuthority);
  const workspace = scope?.workspace_id, device = scope?.device_id;
  const owner = useMemo(() => ({ workspace, device }), [workspace, device]);
  const current = useRef({ owner, scope, evidence }); current.current = { owner, scope, evidence };
  const mounted = useRef(true); const generation = useRef(0);
  const [refresh, setRefresh] = useState<{ owner: object; generation: number; pending: boolean; failed: boolean }>();
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; generation.current += 1; }; }, []);
  const refetch = useCallback(async () => {
    if (!scope) return;
    const request = ++generation.current;
    const stillOwned = () => mounted.current && generation.current === request && current.current.owner === owner;
    setRefresh({ owner, generation: request, pending: true, failed: false });
    let failed = false;
    try {
      const latest = await groups.refetch();
      if (!stillOwned()) return;
      if (latest.isError || !latest.data || latest.data.workspace_id !== scope.workspace_id || !latest.data.workspace_revision) {
        failed = true; return;
      }
      // A changed head already updates the caller's groups observer and therefore
      // keys a new exact evidence read. Never refetch the captured old scope.
      if (latest.data.workspace_revision === scope.workspace_revision &&
          current.current.scope?.workspace_revision === scope.workspace_revision) await current.current.evidence.refetch();
    } catch {
      failed = true;
    } finally {
      if (stillOwned()) setRefresh({ owner, generation: request, pending: false, failed });
    }
  }, [groups, owner, scope]);
  const ownRefresh = refresh?.owner === owner && refresh.generation === generation.current ? refresh : undefined;
  const blocked = ownRefresh?.pending || ownRefresh?.failed;
  return { data: blocked ? undefined : evidence.data, isLoading: evidence.isLoading,
    isFetching: evidence.isFetching || Boolean(ownRefresh?.pending), isError: evidence.isError || Boolean(ownRefresh?.failed), refetch };
}
