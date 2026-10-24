import type { WriteGroup } from '../../../../types/studioV2WriteGroup';

/** Causes of a failed schema preparation that belong to the destination itself, not to one group or source. */
export type DestinationBlockerCode =
  | 'timescale_unavailable' | 'storage_owner_mismatch' | 'schema_repair_blocked'
  | 'permission_denied' | 'target_unavailable' | 'schema_incompatible';

/** Backend codes, including their older spellings, mapped to one stable blocker cause. */
const BLOCKER_BY_CODE: Record<string, DestinationBlockerCode> = {
  timescale_unavailable: 'timescale_unavailable',
  WRITE_GROUP_TIMESCALE_UNAVAILABLE: 'timescale_unavailable',
  RECORDING_SCHEMA_TIMESCALE_UNAVAILABLE: 'timescale_unavailable',
  storage_owner_mismatch: 'storage_owner_mismatch',
  schema_repair_blocked: 'schema_repair_blocked',
  RECORDING_SCHEMA_PERMISSION_DENIED: 'permission_denied',
  RECORDING_SCHEMA_TARGET_UNAVAILABLE: 'target_unavailable',
  RECORDING_SCHEMA_INCOMPATIBLE: 'schema_incompatible',
};

/** Where the operator fixes it: the destination settings, or the one group whose table needs review. */
export const BLOCKER_REPAIR: Record<DestinationBlockerCode, 'destination' | 'group'> = {
  timescale_unavailable: 'destination', storage_owner_mismatch: 'destination', permission_denied: 'destination',
  target_unavailable: 'destination', schema_incompatible: 'destination', schema_repair_blocked: 'group',
};

export function destinationBlockerCode(code: string | undefined): DestinationBlockerCode | undefined {
  return code !== undefined && Object.prototype.hasOwnProperty.call(BLOCKER_BY_CODE, code) ? BLOCKER_BY_CODE[code] : undefined;
}

export interface DestinationBlocker {
  /** Stable cause plus the destination identity and revision the affected groups are frozen to. */
  key: string;
  code: DestinationBlockerCode;
  groups: WriteGroup[];
}

const SEPARATOR = '␟';

/**
 * One blocker per stable cause and owning destination (connector, its revision, schema, table and storage mode).
 * The same cause on a different destination or revision stays separate, and nothing here confirms or clears a
 * group: each group keeps its own schema confirmation.
 */
export function aggregateDestinationBlockers(
  groups: WriteGroup[],
  reports: Readonly<Record<string, DestinationBlockerCode | undefined>>,
): DestinationBlocker[] {
  const found = new Map<string, DestinationBlocker>();
  for (const group of groups) {
    const code = reports[group.id];
    if (!code) continue;
    const { connector_id, connector_revision, table_schema, table_name } = group.destination;
    const key = [code, connector_id, connector_revision, table_schema, table_name, group.storage_mode ?? ''].join(SEPARATOR);
    const existing = found.get(key);
    if (existing) existing.groups.push(group);
    else found.set(key, { key, code, groups: [group] });
  }
  return Array.from(found.values());
}
