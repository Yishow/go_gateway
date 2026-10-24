import type { WriteGroup, WriteGroupDraft, WriteGroupProfileKind } from '../../../../types/studioV2WriteGroup';

/** The circuit a batch writes to: its identity plus the saved group whose members and destination a new profile reuses. */
export interface BatchIntentTarget {
  circuit_id: string;
  circuit_name: string;
  anyGroup: WriteGroup;
}

export function profileSettings(kind: WriteGroupProfileKind, interval: number) {
  return kind === 'instantaneous'
    ? { recording_profile: { kind, interval, unit: 'second' as const }, intervalSeconds: interval, label: 'Instantaneous' }
    : { recording_profile: { kind, interval, unit: 'minute' as const }, intervalSeconds: interval * 60, label: 'Cumulative' };
}

/** Body of the PUT for a saved profile: only the reviewed period changes, everything else is carried over. */
export function existingProfileUpdateGroup(existing: WriteGroup, kind: WriteGroupProfileKind, interval: number): WriteGroup {
  const settings = profileSettings(kind, interval);
  return { ...existing, recording_profile: settings.recording_profile, row_policy: { ...existing.row_policy, interval_seconds: settings.intervalSeconds } };
}

/** Body of the ensure request for a profile the circuit does not have yet. */
export function newProfileDraft(workspaceId: string, target: BatchIntentTarget, kind: WriteGroupProfileKind, interval: number): WriteGroupDraft {
  const settings = profileSettings(kind, interval);
  const base = target.anyGroup;
  return {
    workspace_id: workspaceId,
    name: `${target.circuit_name || target.circuit_id} (${settings.label})`,
    layout_version: base.layout_version || 'managed_long_v1',
    storage_mode: base.storage_mode || 'plain',
    recording_profile: settings.recording_profile,
    destination: base.destination,
    members: base.members.map((member) => ({ ...member, circuit_id: target.circuit_id, circuit_name: target.circuit_name })),
    row_policy: { interval_seconds: settings.intervalSeconds, allowed_lateness_seconds: base.row_policy.allowed_lateness_seconds ?? 0, incomplete_policy: base.row_policy.incomplete_policy ?? 'partial' },
    write_policy: base.write_policy,
  };
}

/**
 * The fields that make a profile "the same configuration" as the one that was reviewed. Revisions, status,
 * timestamps, schema evidence, server-owned column bindings and display names are deliberately left out.
 */
function comparableIntent(group: WriteGroup | WriteGroupDraft, withTableName: boolean) {
  return {
    workspace_id: group.workspace_id,
    layout_version: group.layout_version ?? 'managed_wide_v1',
    storage_mode: group.storage_mode ?? 'plain',
    recording_profile: group.recording_profile ?? null,
    destination: {
      connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision,
      table_schema: group.destination.table_schema, storage_strategy: group.destination.storage_strategy,
      table_name: withTableName ? group.destination.table_name : '',
    },
    row_policy: {
      interval_seconds: group.row_policy.interval_seconds, allowed_lateness_seconds: group.row_policy.allowed_lateness_seconds ?? 0,
      incomplete_policy: group.row_policy.incomplete_policy ?? '',
    },
    write_policy: { mode: group.write_policy.mode ?? '', dedupe_capability: group.write_policy.dedupe_capability ?? '', recreate_missing: group.write_policy.recreate_missing ?? false },
    members: group.members.map((member) => ({
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id, measurement_id: member.measurement_id ?? '',
      circuit_id: member.circuit_id ?? '', circuit_name: member.circuit_name ?? '', required: member.required,
    })).sort((a, b) => `${a.device_id}\0${a.point_id}\0${a.tag_id}`.localeCompare(`${b.device_id}\0${b.point_id}\0${b.tag_id}`)),
  };
}

/** Whether a saved group carries the reviewed configuration. An intent without a table name lets the server choose it. */
export function profileMatchesIntent(live: WriteGroup, intent: WriteGroup | WriteGroupDraft): boolean {
  const withTableName = intent.destination.table_name !== '';
  return JSON.stringify(comparableIntent(live, withTableName)) === JSON.stringify(comparableIntent(intent, withTableName));
}
