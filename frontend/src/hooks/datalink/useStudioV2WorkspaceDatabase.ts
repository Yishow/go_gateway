import { beginConfigSave, confirmConfigSave, failConfigSave, validateConfigAck, configReadContext, configReadKey, fencedConfigRead } from './recordingConfigRefresh';
import type { StudioV2LocatorReadScope } from '../../types/studioV2LocatorEntry';
import { studioLocatorReadKey } from '../../services/studioEntityLocator';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceDatabaseAPI,
  type StudioV2WorkspaceDatabaseConfigRequest,
  type StudioV2WorkspaceDatabaseTargetRequest,
} from '../../services/studioV2WorkspaceDatabase';

export function useStudioV2DatabaseConfigQuery(enabled: boolean, workspaceId?: string) {
  const context = configReadContext(workspaceId);
  return useQuery({
    queryKey: configReadKey(studioV2WorkspaceKeys.databaseConfig(), context),
    queryFn: ({ signal, client }) => fencedConfigRead(client, context, signal, async () => {
      const record = await studioV2WorkspaceDatabaseAPI.getConfig({ signal, baseURL: context.apiOrigin });
      if (record && workspaceId && record.workspace_id !== workspaceId) throw new Error('Configuration observation unavailable');
      return record;
    }),
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

export function useStudioV2DatabaseTargetsQuery(enabled: boolean, workspaceId?: string) {
  const context = configReadContext(workspaceId);
  return useQuery({
    queryKey: configReadKey(studioV2WorkspaceKeys.databaseTargets(), context),
    queryFn: ({ signal, client }) => fencedConfigRead(client, context, signal, async () => {
      const records = await studioV2WorkspaceDatabaseAPI.listTargets({ signal, baseURL: context.apiOrigin });
      if (workspaceId && records.some(record => record.workspace_id !== workspaceId)) throw new Error('Target observation unavailable');
      return records;
    }),
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

export function useUpdateStudioV2DatabaseConfigMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: async (request: StudioV2WorkspaceDatabaseConfigRequest) => {
      const captured = beginConfigSave(queryClient, workspaceId);
      try {
        const record = await studioV2WorkspaceDatabaseAPI.updateConfig(request, { baseURL: captured.context.apiOrigin });
        validateConfigAck(record, request, captured.context.workspaceId);
        confirmConfigSave(queryClient, captured);
        return record;
      } catch (error) { failConfigSave(queryClient, captured); throw error; }
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
