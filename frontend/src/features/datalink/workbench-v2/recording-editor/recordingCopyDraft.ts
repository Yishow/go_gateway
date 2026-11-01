import type { WriteGroup, WriteGroupMember } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupCalculationCapabilities } from '../../../../types/studioV2WriteGroupCalculation';
import { memberKey, type EditorDraft } from '../state/writeGroup/draft';
import { newMixedDraft } from './recordingMixedDraft';

export type RecordingCopyMemberIntent = Pick<WriteGroupMember, 'device_id' | 'point_id' | 'tag_id' | 'calculation' | 'cadence' | 'required' | 'max_age_seconds' | 'source_value_type' | 'value_role' | 'measurement_id'>;
export interface RecordingCopyPrefill {
  sourceName: string;
  layout: 'managed_wide_profile_v1' | 'managed_long_v1';
  allowedLatenessSeconds: number;
  incompletePolicy: 'partial' | 'skip_row';
  members: RecordingCopyMemberIntent[];
}
type Destination = Pick<EditorDraft, 'connector_id' | 'connector_revision' | 'table_schema' | 'table_name'>;

/** A new configuration intent only: no source revision, SQL binding, confirmation or execution state. */
function memberIntent(member: RecordingCopyMemberIntent): RecordingCopyMemberIntent {
  return {
    device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id, required: member.required,
    ...(member.calculation !== undefined ? { calculation: member.calculation } : {}),
    ...(member.cadence ? { cadence: { interval: member.cadence.interval, unit: member.cadence.unit, alignment: member.cadence.alignment } } : {}),
    ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
    ...(member.source_value_type !== undefined ? { source_value_type: member.source_value_type } : {}),
    ...(member.value_role !== undefined ? { value_role: member.value_role } : {}),
    ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
  };
}
/** Canonical freshness belongs to the caller; these exclusions never grant save or start authority. */
export function recordingCopyPrefill(group: WriteGroup, capabilities: WriteGroupCalculationCapabilities | undefined, dialect: string): RecordingCopyPrefill | undefined {
  const layout = group.layout_version;
  if (!['draft', 'ready', 'running', 'disabled'].includes(group.status) || !group.id || !group.workspace_id || !group.revision || !group.record_id ||
      group.recording_profile?.kind !== 'mixed' || group.recording_profile.result_shape !== 'member-calculation-v1' ||
      group.destination.storage_strategy !== 'managed' || group.storage_mode !== undefined && group.storage_mode !== 'plain' ||
      layout !== 'managed_wide_profile_v1' && layout !== 'managed_long_v1' || !capabilities?.editor_supported || !capabilities.draft_supported ||
      !capabilities.schema_dialects.some(value => value === dialect) || !capabilities.schema_layouts.includes(layout)) return undefined;
  if (!group.members.length || group.members.some(member => !member.calculation || !capabilities.methods.includes(member.calculation) ||
      !member.cadence || !capabilities.cadence.units.includes(member.cadence.unit) || member.cadence.alignment !== capabilities.cadence.alignment ||
      !Number.isSafeInteger(member.cadence.interval) || member.cadence.interval < capabilities.cadence.minimum ||
      member.cadence.interval > (member.cadence.unit === 'minute' ? capabilities.cadence.max_minutes : capabilities.cadence.max_seconds) ||
      member.exact_recording_source && capabilities.unsupported_combinations.some(rule => rule.method === member.calculation &&
        rule.source_transform === 'exact_counter_scale' && (!rule.layout || rule.layout === layout)))) return undefined;
  const incompletePolicy = group.row_policy.incomplete_policy;
  if (incompletePolicy !== 'partial' && incompletePolicy !== 'skip_row') return undefined;
  return { sourceName: group.name, layout, allowedLatenessSeconds: group.row_policy.allowed_lateness_seconds,
    incompletePolicy, members: group.members.map(memberIntent) };
}
/** Caller supplies a newly generated ID and current connector; managed table assignment stays server-owned. */
export function createRecordingCopyDraft(prefill: RecordingCopyPrefill, destination: Destination, recordId: string, name: string): EditorDraft {
  return { ...newMixedDraft({ ...destination, table_name: '' }, recordId, prefill.layout), name,
    allowed_lateness_seconds: prefill.allowedLatenessSeconds, incomplete_policy: prefill.incompletePolicy,
    members: prefill.members.map(member => ({ ...memberIntent(member), key: memberKey(member.point_id, member.tag_id), entity_key: '', target_column: '' })) };
}
