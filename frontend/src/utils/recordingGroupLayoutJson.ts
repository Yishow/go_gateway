import type { GroupSchemaColumn, GroupSchemaIndex, GroupSchemaLayout } from '../types/recordingPlan';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';
import { isRecord } from './studioV2WriteGroupJsonPrimitives';

function parseGroupSchemaIndexes(value: unknown): GroupSchemaIndex[] | undefined | null {
  if (value === undefined || value === null) return undefined;
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const indexes: GroupSchemaIndex[] = [];
  for (const item of value) {
    if (!isRecord(item) || !Array.isArray(item.columns) || item.columns.length === 0 ||
      item.columns.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
    const name = boundedString(item.name);
    const columns = item.columns.map((column) => boundedString(column));
    if (!name || columns.some((column) => !column)) return null;
    indexes.push({ name, columns: columns as string[] });
  }
  return indexes;
}

/** Rebuilds the server-resolved group table layout; a malformed layout invalidates the preview. */
export function parseGroupSchemaLayout(value: unknown): GroupSchemaLayout | undefined | null {
  if (value === undefined || value === null) return undefined;
  if (!isRecord(value)) return null;
  const tableName = boundedString(value.table_name);
  const ownerColumn = boundedString(value.owner_column);
  if (!tableName || !ownerColumn || !Array.isArray(value.columns) || value.columns.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const columns: GroupSchemaColumn[] = [];
  for (const item of value.columns) {
    if (!isRecord(item)) return null;
    const name = boundedString(item.name);
    const sqlType = boundedString(item.sql_type);
    if (!name || !sqlType || typeof item.nullable !== 'boolean' || typeof item.primary_key !== 'boolean') return null;
    columns.push({ name, sql_type: sqlType, nullable: item.nullable, primary_key: item.primary_key });
  }
  const indexes = parseGroupSchemaIndexes(value.indexes);
  if (indexes === null) return null;
  return {
    table_name: tableName, columns, owner_column: ownerColumn,
    ...(indexes !== undefined ? { indexes } : {}),
  };
}
