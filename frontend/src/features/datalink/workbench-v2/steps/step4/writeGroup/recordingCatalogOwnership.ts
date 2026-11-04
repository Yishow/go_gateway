import type { ConfiguredDestinationObservation } from '../../../recording-editor/ConfiguredDestinationChooser';
import type { DbConnector } from '../../../state/types';

/** A catalog is not ownership evidence. Only current configured groups can identify a known wide-table owner. */
export function ownedWideCatalogTables(observation: ConfiguredDestinationObservation | undefined, scope: {
  workspaceId: string; workspaceRevision: string; connector: DbConnector; groupId?: string;
}, defaultSchema: string): { schema: string; name: string }[] {
  const data = observation?.data;
  if (!data || observation?.loading || observation?.failed || data.workspace_id !== scope.workspaceId || data.workspace_revision !== scope.workspaceRevision) return [];
  return data.groups.flatMap(group => {
    const destination = group.destination;
    if (group.id === scope.groupId || group.workspace_id !== scope.workspaceId || group.status === 'deleted' ||
        destination.storage_strategy !== 'managed' || group.layout_version === 'managed_long_v1' ||
        destination.connector_id !== scope.connector.connector_id || destination.database !== scope.connector.database || !destination.table_name) return [];
    return [{ schema: destination.table_schema || defaultSchema, name: destination.table_name }];
  });
}

/** Refuse manually entered known wide owners as well as catalog selections. */
export function hasOwnedWideDestination(observation: ConfiguredDestinationObservation | undefined, scope: {
  workspaceId: string; workspaceRevision: string; connector: DbConnector; groupId?: string;
}, schema: string, table: string): boolean {
  const defaultSchema = scope.connector.schema || (scope.connector.kind === 'postgres' ? 'public' : scope.connector.kind === 'sqlite' ? 'main' : scope.connector.kind === 'mysql' ? scope.connector.database : '');
  const targetSchema = schema.trim() || defaultSchema;
  return ownedWideCatalogTables(observation, scope, defaultSchema).some(owner => owner.name === table.trim() &&
    (owner.schema === targetSchema || scope.connector.kind === 'sqlserver' && !targetSchema));
}
