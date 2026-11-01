import type { EditorDraft } from '../state/writeGroup/draft';
import type { WriteGroupDraft, WriteGroupLayoutVersion } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupCalculationCapabilities } from '../../../../types/studioV2WriteGroupCalculation';
import { parseWriteGroupCalculationCapabilities } from '../../../../utils/studioV2WriteGroupCalculationJson';
import { cloneWriteGroupExactRecordingSource, parseWriteGroupExactRecordingSource } from '../../../../utils/studioV2ExactRecordingSourceJson';

type Destination = Pick<EditorDraft, 'connector_id' | 'connector_revision' | 'table_schema' | 'table_name'>;
const validRecordId = (id: string | undefined): boolean => typeof id === 'string' &&
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/.test(id) && id !== '00000000-0000-0000-0000-000000000000';
const numericTypes = new Set(['int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64']);

/** New independent intent only. No legacy conversion, copied baseline or server-owned layout bindings. */
export function newMixedDraft(destination: Destination, recordId: string, layout: WriteGroupLayoutVersion = 'managed_wide_profile_v1'): EditorDraft {
  if (!validRecordId(recordId)) throw new Error('mixed-record-id-invalid');
  if (layout !== 'managed_wide_profile_v1' && layout !== 'managed_long_v1') throw new Error('mixed-layout-unsupported');
  return {
    name: '', connector_id: destination.connector_id, connector_revision: destination.connector_revision,
    table_schema: destination.table_schema, table_name: destination.table_name,
    record_id: recordId, storage_strategy: 'managed', layout_version: layout,
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
    allowed_lateness_seconds: 0, incomplete_policy: 'partial', entity_key_column: '', provenance_column: '', dedupe_capability: '', members: [],
  };
}

/** Local exclusions only: backend CAS, counter confirmation, source typing and schema remain authoritative. */
function issuesFor(draft: EditorDraft, capabilities: WriteGroupCalculationCapabilities, dialect?: string): string[] {
  const issues: string[] = [];
  const capability = parseWriteGroupCalculationCapabilities(capabilities);
  if (!capability?.draft_supported) issues.push('mixed-capability-unavailable');
  const profile = draft.recording_profile;
  if (profile?.kind !== 'mixed' || profile.result_shape !== 'member-calculation-v1' ||
      Object.keys(profile).some(key => key !== 'kind' && key !== 'result_shape')) issues.push('mixed-profile-required');
  if (!validRecordId(draft.record_id)) issues.push('mixed-record-id-invalid');
  if (draft.interval_seconds !== undefined) issues.push('mixed-group-interval-forbidden');
  if (draft.storage_strategy !== 'managed' || (draft.storage_mode !== undefined && draft.storage_mode !== 'plain')) issues.push('mixed-storage-unsupported');
  if (!draft.layout_version || !['managed_wide_profile_v1', 'managed_long_v1'].includes(draft.layout_version) ||
      !capability?.schema_layouts.includes(draft.layout_version)) issues.push('mixed-layout-unsupported');
  if (dialect !== undefined && !capability?.schema_dialects.some(supported => supported === dialect)) issues.push('mixed-dialect-unsupported');
  if (!draft.name.trim()) issues.push('name-required');
  if (!draft.connector_id.trim() || !draft.connector_revision.trim()) issues.push('destination-unsaved');
  if (!draft.members.length) issues.push('members-required');
  if (!Number.isFinite(draft.allowed_lateness_seconds) || draft.allowed_lateness_seconds < 0) issues.push('lateness-invalid');
  if (draft.layout_version === 'managed_wide_profile_v1' && draft.incomplete_policy !== 'partial') issues.push('mixed-row-policy-invalid');
  const seen = new Set<string>();
  for (const member of draft.members) {
    const identity = JSON.stringify([member.device_id, member.point_id, member.tag_id]);
    if (!member.device_id.trim() || !member.point_id.trim() || !member.tag_id.trim()) issues.push('member-source-required');
    if (seen.has(identity)) issues.push('member-source-duplicate');
    seen.add(identity);
    if (member.entity_key.trim()) issues.push('member-entity-unsupported');
    if (!member.calculation || !capability?.methods.includes(member.calculation)) issues.push('member-method-unsupported');
    const cadence = member.cadence;
    if (!cadence || !capability || cadence.alignment !== capability.cadence.alignment || !capability.cadence.units.includes(cadence.unit) ||
        !Number.isSafeInteger(cadence.interval) || cadence.interval < capability.cadence.minimum ||
        cadence.interval > (cadence.unit === 'minute' ? capability.cadence.max_minutes : capability.cadence.max_seconds)) issues.push('member-cadence-invalid');
    const numeric = member.source_value_type !== undefined && numericTypes.has(member.source_value_type);
    if (member.calculation === 'mean' && (!numeric || (member.value_role !== undefined && member.value_role !== 'gauge'))) issues.push('member-mean-inapplicable');
    if ((member.calculation === 'delta' && (!numeric || (member.value_role !== undefined && member.value_role !== 'counter_raw') || !member.measurement_id?.trim())) ||
        (member.value_role === 'counter_raw' && (!numeric || !member.measurement_id?.trim()))) issues.push('member-delta-unconfirmed');
    if (member.exact_recording_source !== undefined) {
      const exact = parseWriteGroupExactRecordingSource(member.exact_recording_source);
      if (!exact || exact.measurement_id !== member.measurement_id || exact.original_data_type !== member.source_value_type) issues.push('member-exact-source-invalid');
      for (const restriction of capability?.unsupported_combinations ?? []) {
        if (restriction.method === member.calculation && (restriction.layout === undefined || restriction.layout === draft.layout_version)) issues.push(restriction.reason);
      }
    }
  }
  return [...new Set(issues)];
}

export function mixedDraftIssues(draft: EditorDraft, capabilities: WriteGroupCalculationCapabilities, dialect: string): string[] {
  return issuesFor(draft, capabilities, dialect);
}

/** Explicit canonical whitelist; UI input strings, leases, readiness and baseline state never enter the request. */
export function mixedDraftToRequestGroup(draft: EditorDraft, workspaceId: string, capabilities: WriteGroupCalculationCapabilities): WriteGroupDraft {
  const issues = issuesFor(draft, capabilities);
  if (!workspaceId.trim()) issues.push('workspace-required');
  if (issues.length) throw new Error(issues.join(','));
  return {
    workspace_id: workspaceId, record_id: draft.record_id, name: draft.name.trim(),
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, layout_version: draft.layout_version,
    ...(draft.storage_mode !== undefined ? { storage_mode: draft.storage_mode } : {}),
    destination: { connector_id: draft.connector_id, connector_revision: draft.connector_revision,
      table_schema: draft.table_schema, table_name: draft.table_name.trim(), storage_strategy: draft.storage_strategy },
    members: draft.members.map(member => ({
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      calculation: member.calculation, cadence: member.cadence ? { interval: member.cadence.interval, unit: member.cadence.unit, alignment: member.cadence.alignment } : undefined,
      target_column: member.target_column.trim(), required: member.required,
      ...(member.source_revision !== undefined ? { source_revision: member.source_revision } : {}),
      ...(member.mapping_revision !== undefined ? { mapping_revision: member.mapping_revision } : {}),
      ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
      ...(member.entity_key.trim() ? { entity_key: member.entity_key.trim() } : {}),
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
      ...(member.circuit_id !== undefined ? { circuit_id: member.circuit_id } : {}),
      ...(member.circuit_name !== undefined ? { circuit_name: member.circuit_name } : {}),
      ...(member.value_role !== undefined ? { value_role: member.value_role } : {}),
      ...(member.source_value_type !== undefined ? { source_value_type: member.source_value_type } : {}),
      ...(member.result_value_type !== undefined ? { result_value_type: member.result_value_type } : {}),
      ...(member.exact_recording_source !== undefined ? { exact_recording_source: cloneWriteGroupExactRecordingSource(member.exact_recording_source) } : {}),
    })),
    row_policy: {
      allowed_lateness_seconds: draft.allowed_lateness_seconds, incomplete_policy: draft.incomplete_policy,
      ...(draft.entity_key_column.trim() ? { entity_key_column: draft.entity_key_column.trim() } : {}),
      ...(draft.provenance_column.trim() ? { provenance_column: draft.provenance_column.trim() } : {}),
      ...(draft.group_key_columns !== undefined ? { group_key_columns: [...draft.group_key_columns] } : {}),
      ...(draft.unique_key_columns !== undefined ? { unique_key_columns: [...draft.unique_key_columns] } : {}),
      ...(draft.value_column !== undefined ? { value_column: draft.value_column } : {}),
      ...(draft.quality_column !== undefined ? { quality_column: draft.quality_column } : {}),
      ...(draft.record_key_column !== undefined ? { record_key_column: draft.record_key_column } : {}),
      ...(draft.bucket_start_column !== undefined ? { bucket_start_column: draft.bucket_start_column } : {}),
      ...(draft.due_at_column !== undefined ? { due_at_column: draft.due_at_column } : {}),
      ...(draft.group_id_column !== undefined ? { group_id_column: draft.group_id_column } : {}),
      ...(draft.device_id_column !== undefined ? { device_id_column: draft.device_id_column } : {}),
    },
    write_policy: {
      ...(draft.dedupe_capability ? { dedupe_capability: draft.dedupe_capability } : {}),
      ...(draft.recreate_missing !== undefined ? { recreate_missing: draft.recreate_missing } : {}),
    },
  };
}
