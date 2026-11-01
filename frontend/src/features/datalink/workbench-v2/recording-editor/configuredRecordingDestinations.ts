import type { WriteGroup, WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import type { DbConnector } from '../state/types';
import type { EditorDraft } from '../state/writeGroup/draft';
export interface ConfiguredDestinationContext { workspaceId: string; connector: DbConnector; draft: EditorDraft; groupId?: string }
/** Configuration identity only. The backend remains the sole canonical table-name authority. */
export function configuredDestinationIdentity(group: WriteGroup): string {
  const d = group.destination;
  return JSON.stringify([group.workspace_id, d.connector_id, d.database, d.table_schema, d.table_name,
    d.storage_strategy, group.layout_version ?? 'managed_wide_v1', group.storage_mode ?? 'plain']);
}
export function configuredDestinationReason(group: WriteGroup, context: ConfiguredDestinationContext): string | undefined {
  const { connector, draft } = context, d = group.destination;
  if (group.workspace_id !== context.workspaceId) return 'scope';
  if (group.status === 'deleted' || !group.revision || !d.table_name || !d.connector_revision) return 'unavailable';
  if (d.storage_strategy === 'managed' && group.layout_version !== 'managed_long_v1' && group.id !== context.groupId) return 'wideOwner';
  if (!connector.persisted || connector.save_state !== 'saved' || !connector.identity_revision) return 'connector';
  if (d.connector_id !== connector.connector_id || d.database !== connector.database || d.table_schema !== connector.schema) return 'destination';
  if (d.connector_revision !== connector.identity_revision || draft.connector_revision !== connector.identity_revision) return 'revision';
  if (draft.connector_id !== d.connector_id || draft.table_schema !== d.table_schema || draft.storage_strategy !== d.storage_strategy ||
    (draft.layout_version ?? 'managed_wide_v1') !== (group.layout_version ?? 'managed_wide_v1') ||
    (draft.storage_mode ?? 'plain') !== (group.storage_mode ?? 'plain')) return 'layout';
  if (group.id !== context.groupId && (d.storage_strategy !== 'managed' || group.layout_version !== 'managed_long_v1')) return 'layout';
  return undefined;
}
/** Explicit allowlist: never spread a destination or a saved group into a draft. */
export function configuredDestinationPatch(group: WriteGroup): Partial<EditorDraft> {
  const d = group.destination;
  return { connector_id: d.connector_id, connector_revision: d.connector_revision, table_schema: d.table_schema, table_name: d.table_name };
}
export function configuredDestinationStillCurrent(group: WriteGroup, list: Pick<WriteGroupListResponse, 'groups'> & Partial<Pick<WriteGroupListResponse, 'workspace_id'>> | undefined): boolean {
  if (!list || list.workspace_id !== undefined && list.workspace_id !== group.workspace_id) return false;
  return list.groups.some(current => current.id === group.id && current.status !== 'deleted' && current.revision === group.revision &&
    current.destination.connector_revision === group.destination.connector_revision && configuredDestinationIdentity(current) === configuredDestinationIdentity(group));
}
