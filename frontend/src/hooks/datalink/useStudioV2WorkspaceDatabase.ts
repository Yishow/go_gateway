import type { StudioV2LocatorReadScope } from '../../types/studioV2LocatorEntry';
import { studioLocatorReadKey } from '../../services/studioEntityLocator';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceDatabaseAPI,
  type StudioV2WorkspaceDatabaseConfigRequest,
  type StudioV2WorkspaceDatabaseTargetRequest,
} from '../../services/studioV2WorkspaceDatabase';

export function useStudioV2DatabaseConfigQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.databaseConfig(),
    queryFn: () => studioV2WorkspaceDatabaseAPI.getConfig(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}

export function useStudioV2MySQLRecordingCapabilityQuery(connectorId: string, connectorRevision: string, enabled: boolean, locatorScope?: StudioV2LocatorReadScope) {
  const selectedLocatorScope = locatorScope === undefined ? undefined : { workspaceId: locatorScope.workspaceId, workspaceRevision: locatorScope.workspaceRevision };
  return useQuery({
    queryKey: [...studioV2WorkspaceKeys.mysqlRecordingCapability(connectorId, connectorRevision), ...studioLocatorReadKey(selectedLocatorScope)],
    queryFn: () => selectedLocatorScope ? studioV2WorkspaceDatabaseAPI.getMySQLRecordingCapability(connectorId, connectorRevision, selectedLocatorScope) : studioV2WorkspaceDatabaseAPI.getMySQLRecordingCapability(connectorId, connectorRevision),
    enabled: enabled && Boolean(connectorId && connectorRevision),
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 0,
  });
}

export function useStudioV2SQLServerRecordingCapabilityQuery(connectorId: string, connectorRevision: string, enabled: boolean, locatorScope?: StudioV2LocatorReadScope) {
  const selectedLocatorScope = locatorScope === undefined ? undefined : { workspaceId: locatorScope.workspaceId, workspaceRevision: locatorScope.workspaceRevision };
  return useQuery({
    queryKey: [...studioV2WorkspaceKeys.sqlserverRecordingCapability(connectorId, connectorRevision), ...studioLocatorReadKey(selectedLocatorScope)],
    queryFn: () => selectedLocatorScope ? studioV2WorkspaceDatabaseAPI.getSQLServerRecordingCapability(connectorId, connectorRevision, selectedLocatorScope) : studioV2WorkspaceDatabaseAPI.getSQLServerRecordingCapability(connectorId, connectorRevision),
    enabled: enabled && Boolean(connectorId && connectorRevision),
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 0,
  });
}

export function useStudioV2DatabaseTargetsQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.databaseTargets(),
    queryFn: () => studioV2WorkspaceDatabaseAPI.listTargets(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}

export interface StudioV2DatabaseMetadataScope {
  connectorId?: string;
  connectorRevision?: string;
  database: string;
  schema: string;
  table: string;
  groupId?: string;
  groupRevision?: string;
}

/** 以完整範圍作為 key 查詢實際欄位；任一範圍改變即換成新的查詢，舊回覆不會沿用。 */
export function useStudioV2DatabaseMetadataQuery(scope: StudioV2DatabaseMetadataScope, enabled: boolean, locatorScope?: StudioV2LocatorReadScope) {
  const selectedLocatorScope = locatorScope === undefined ? undefined : { workspaceId: locatorScope.workspaceId, workspaceRevision: locatorScope.workspaceRevision };
  const connectorId = scope.connectorId ?? '';
  const connectorRevision = scope.connectorRevision ?? '';
  return useQuery({
    queryKey: [...studioV2WorkspaceKeys.databaseMetadata(connectorId, connectorRevision, scope.database, scope.schema, scope.table, scope.groupId, scope.groupRevision), ...studioLocatorReadKey(selectedLocatorScope)],
    queryFn: () => {
      const group = scope.groupId && scope.groupRevision ? { groupId: scope.groupId, groupRevision: scope.groupRevision } : undefined;
      return selectedLocatorScope ? studioV2WorkspaceDatabaseAPI.getMetadata(connectorRevision, group, selectedLocatorScope) : studioV2WorkspaceDatabaseAPI.getMetadata(connectorRevision, group);
    },
    enabled: enabled && Boolean(connectorId && connectorRevision),
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useUpdateStudioV2DatabaseConfigMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: StudioV2WorkspaceDatabaseConfigRequest) => studioV2WorkspaceDatabaseAPI.updateConfig(request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseConfig() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpsertStudioV2DatabaseTargetMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ pointId, request }: { pointId: string; request: StudioV2WorkspaceDatabaseTargetRequest }) =>
      studioV2WorkspaceDatabaseAPI.upsertTarget(pointId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
