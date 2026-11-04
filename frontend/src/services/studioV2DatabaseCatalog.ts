import { assertStudioLocatorWorkspaceBinding, captureStudioLocatorReadScope, studioLocatorReadParams, throwStudioLocatorReadError } from './studioEntityLocator';
import { studioV2DatalinkApi } from './studioV2Workspace';
import type { APIResponse, DatabaseConnectorKind } from '../types/datalink';
import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';

export interface RecordingCatalogScope {
  workspaceId: string; workspaceRevision: string; connectorId: string; connectorRevision: string; kind: DatabaseConnectorKind;
  groupId?: string; groupRevision?: string;
}
export interface RecordingDestinationCatalog {
  workspace_id: string; workspace_revision: string; connector_id: string; connector_revision: string;
  kind: DatabaseConnectorKind; default_schema: string; schemas: { name: string; tables: string[] }[];
}

/** Reject partial or differently scoped results: the catalog offers choices, never save authorization. */
export function parseRecordingDestinationCatalog(value: unknown, scope: RecordingCatalogScope): RecordingDestinationCatalog {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('Invalid catalog');
  const result = value as Record<string, unknown>;
  if (result.workspace_id !== scope.workspaceId || result.workspace_revision !== scope.workspaceRevision ||
      result.connector_id !== scope.connectorId || result.connector_revision !== scope.connectorRevision ||
      result.kind !== scope.kind || typeof result.default_schema !== 'string' || !Array.isArray(result.schemas)) throw new Error('Invalid catalog scope');
  const names = new Set<string>();
  const schemas = result.schemas.map((item: unknown) => {
    if (!item || typeof item !== 'object' || Array.isArray(item)) throw new Error('Invalid catalog schema');
    const schema = item as Record<string, unknown>;
    if (typeof schema.name !== 'string' || !schema.name.trim() || names.has(schema.name) || !Array.isArray(schema.tables) ||
        schema.tables.some((table: unknown) => typeof table !== 'string' || !table.trim()) || new Set(schema.tables).size !== schema.tables.length) throw new Error('Invalid catalog schema');
    names.add(schema.name);
    return { name: schema.name, tables: schema.tables as string[] };
  });
  return { workspace_id: scope.workspaceId, workspace_revision: scope.workspaceRevision, connector_id: scope.connectorId,
    connector_revision: scope.connectorRevision, kind: scope.kind, default_schema: result.default_schema, schemas };
}

export async function readRecordingDestinationCatalog(scope: RecordingCatalogScope, locatorScope?: StudioV2LocatorReadScope): Promise<RecordingDestinationCatalog> {
  const selected = captureStudioLocatorReadScope(locatorScope);
  const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/database-catalog', { params: {
    workspace_id: scope.workspaceId, expected_workspace_revision: scope.workspaceRevision, connector_id: scope.connectorId,
    expected_connector_revision: scope.connectorRevision, ...studioLocatorReadParams(selected),
    ...(scope.groupId && scope.groupRevision ? { group_id: scope.groupId, expected_group_revision: scope.groupRevision } : {}),
  } }).catch(error => throwStudioLocatorReadError(error, selected));
  if (res.data.success !== true) throw new Error('Catalog unavailable');
  assertStudioLocatorWorkspaceBinding(res.data.data, selected);
  return parseRecordingDestinationCatalog(res.data.data, scope);
}
