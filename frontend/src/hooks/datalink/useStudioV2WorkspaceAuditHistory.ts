import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceAuditAPI } from '../../services/studioV2WorkspaceAudit';

export function useStudioV2WorkspaceAuditHistoryQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.auditHistory(),
    queryFn: () => studioV2WorkspaceAuditAPI.list(10),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
