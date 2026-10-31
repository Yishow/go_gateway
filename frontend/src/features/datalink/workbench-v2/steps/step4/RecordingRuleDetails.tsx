import type { ReactNode } from 'react';
import type { RecordingReadonly, RecordingWorkflowView } from '../../state/recordingWorkflowView';
import type { WriteGroupMember } from '../../../../../types/studioV2WriteGroup';
import type { RecordingMember } from '../../../../../types/runtimeRuleEvidence';
import { isSafeRuntimeSensorText, safeRuntimeIdentity } from '../../../../../utils/runtimeOperatorSafety';
import type { RecordingGuidedLabels } from './recordingGuidedLabels';

type Labels = RecordingGuidedLabels;
type Row = readonly [string, ReactNode];
const enumText = (value: unknown, allowed: readonly string[], labels: Labels) =>
  typeof value === 'string' && allowed.includes(value) ? value : labels.unavailable;
const mapped = (value: unknown, values: Readonly<Record<string, string>>, fallback: string) =>
  typeof value === 'string' && Object.hasOwn(values, value) ? values[value] : fallback;
const knownState = (value: unknown, labels: Labels) => typeof value === 'string' && Object.hasOwn(labels.states, value)
  ? labels.states[value as keyof Labels['states']] : labels.unavailable;
const identity = (value: unknown, labels: Labels) => safeRuntimeIdentity(value) || labels.unavailable;
const scalar = (value: unknown, labels: Labels) => typeof value === 'number' && Number.isFinite(value) ? String(value)
  : typeof value === 'string' && value.length > 0 && isSafeRuntimeSensorText(value) ? value : labels.unavailable;
const decimal = (value: unknown, labels: Labels) => typeof value === 'string' && /^[+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?$/.test(value)
  ? value : labels.unavailable;
const yesNo = (value: unknown, labels: Labels) => typeof value === 'boolean'
  ? labels.detailFields[value ? 'true' : 'false'] : labels.unavailable;
const timestamp = (value: string | null | undefined, labels: Labels) => value && /^\d{4}-\d{2}-\d{2}T[\d:.+-]+Z?$/.test(value)
  && Number.isFinite(Date.parse(value)) ? value : labels.unavailable;

function Fields({ rows }: { rows: readonly Row[] }) {
  return <dl className="grid min-w-0 gap-3 sm:grid-cols-2">{rows.map(([label, value], index) =>
    <div key={`${label}-${index}`} className="min-w-0"><dt className="text-xs text-slate-400">{label}</dt>
      <dd className="break-words text-sm text-slate-100">{value}</dd></div>)}</dl>;
}
function memberRows(member: RecordingReadonly<RecordingMember>, labels: Labels): Row[] {
  const f = labels.detailFields;
  return [[labels.fields.device, identity(member.device_id, labels)], [f.point, identity(member.point_id, labels)],
    [f.tag, identity(member.tag_id, labels)], [labels.fields.circuit, identity(member.circuit_id, labels)],
    [f.column, identity(member.column, labels)], [labels.fields.sources, mapped(member.source_value_type, labels.sourceTypes, labels.unavailable)]];
}
function SavedMember({ member, labels }: { member: RecordingReadonly<WriteGroupMember>; labels: Labels }) {
  const f = labels.detailFields;
  const proof = member.exact_recording_source;
  const policy = proof?.counter_policy;
  const rows: Row[] = [
    ...memberRows({ device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      column: member.target_column, circuit_id: member.circuit_id ?? null, source_value_type: member.source_value_type ?? null }, labels),
    [f.sourceRevision, identity(member.source_revision, labels)], [f.mappingRevision, identity(member.mapping_revision, labels)],
    [f.required, yesNo(member.required, labels)], [f.maxAge, scalar(member.max_age_seconds, labels)],
    [f.valueRole, member.value_role && ['gauge', 'counter_raw', 'state'].includes(member.value_role) ? f[member.value_role] : labels.unavailable],
    [f.resultType, enumText(member.result_value_type, ['bool', 'text', 'int64', 'uint64', 'decimal', 'float64'], labels)],
    [f.measurement, identity(proof?.measurement_id ?? member.measurement_id, labels)],
  ];
  if (proof) rows.push([f.originalType, mapped(proof.original_data_type, labels.sourceTypes, labels.unavailable)], [f.definitionRevision, identity(proof.definition_revision, labels)],
    [f.sourceBindingRevision, identity(proof.source_binding_revision, labels)], [f.epoch, identity(proof.series_epoch, labels)],
    [f.scale, decimal(proof.scale, labels)], [f.physicalUnit, scalar(proof.unit, labels)],
    [f.semanticKind, proof.semantic_kind === 'counter' || proof.semantic_kind === 'signed_counter' ? f[proof.semantic_kind] : labels.unavailable], [f.allowNegative, yesNo(policy?.allow_negative, labels)],
    [f.resetThreshold, scalar(policy?.reset_threshold, labels)], [f.maxRollover, scalar(policy?.max_roll_over, labels)],
    [f.rolloverModulus, decimal(policy?.rollover_modulus, labels)], [f.maxRate, decimal(policy?.max_rate_per_second, labels)]);
  return <li className="min-w-0 rounded border border-slate-700 p-3"><Fields rows={rows} /></li>;
}

/** Allowlisted projections of existing owners only. No query, editor, raw token, SQL, or mutation. */
export function RecordingRuleDetails({ view, labels }: { view: RecordingWorkflowView; labels: Labels }) {
  if (view.state !== 'available') return <p role="status">{labels.reasons[view.reason]}</p>;
  const group = view.saved.group;
  const effective = view.effective.state === 'available' ? view.effective.value : null;
  const runtime = view.runtime.state === 'available' ? view.runtime.value : null;
  const delivery = runtime?.delivery;
  const current = delivery?.current;
  const operation = view.operation.state === 'available' && view.operation.value.kind !== 'none' ? view.operation.value.operation : null;
  const f = labels.detailFields;
  const stages = current?.stages;
  const counterRelevant = group.recording_profile?.kind === 'cumulative' || group.members.some(member =>
    member.value_role === 'counter_raw' || member.exact_recording_source !== undefined) || effective?.facts.recording_profile?.kind === 'cumulative';
  const columns = (values: readonly string[] | undefined) => values?.length ? values.map(value => identity(value, labels)).join(', ') : labels.unavailable;
  return <div className="space-y-5">
    <section aria-label={labels.savedDetails} className="min-w-0 space-y-3">
      <h4 className="font-medium">{labels.savedDetails}</h4>
      <p className="text-sm text-slate-400">{labels.savedDescription}</p>
      <Fields rows={[
        [labels.fields.interval, `${scalar(group.row_policy.interval_seconds, labels)} ${labels.units.second}`],
        [labels.fields.lateness, `${scalar(group.row_policy.allowed_lateness_seconds, labels)} ${labels.units.second}`],
        [labels.fields.incomplete, mapped(group.row_policy.incomplete_policy, labels.incompletePolicies, labels.unavailable)],
        [f.connector, identity(group.destination.connector_id, labels)], [f.connectorRevision, identity(group.destination.connector_revision, labels)],
        [f.schemaRevision, identity(group.destination.schema_revision, labels)], [f.schemaDigest, identity(group.destination.schema_digest, labels)],
        [f.writeMode, group.write_policy.mode === 'append' || group.write_policy.mode === 'latest' ? f[group.write_policy.mode] : labels.unavailable],
        [f.recreateMissing, yesNo(group.write_policy.recreate_missing, labels)], [f.dedupeCapability, scalar(group.write_policy.dedupe_capability, labels)],
        [f.entityKeyColumn, identity(group.row_policy.entity_key_column, labels)], [f.groupKeyColumns, columns(group.row_policy.group_key_columns)],
        [f.uniqueKeyColumns, columns(group.row_policy.unique_key_columns)], [f.valueColumn, identity(group.row_policy.value_column, labels)],
        [f.qualityColumn, identity(group.row_policy.quality_column, labels)], [f.provenanceColumn, identity(group.row_policy.provenance_column, labels)],
        [f.recordKeyColumn, identity(group.row_policy.record_key_column, labels)], [f.bucketStartColumn, identity(group.row_policy.bucket_start_column, labels)],
        [f.groupIdColumn, identity(group.row_policy.group_id_column, labels)], [f.deviceIdColumn, identity(group.row_policy.device_id_column, labels)],
      ]} />
      <p className="text-xs text-slate-400">{labels.savedReferences}</p>
      <ul aria-label={labels.savedDetails} className="space-y-2">{group.members.map((member) =>
        <SavedMember key={JSON.stringify([member.device_id, member.point_id, member.tag_id])} member={member} labels={labels} />)}</ul>
    </section>
    <section aria-label={labels.effectiveDetails} className="min-w-0 space-y-3">
      <h4 className="font-medium">{labels.effectiveDetails}</h4>
      <p className="text-sm text-slate-400">{labels.effectiveDescription}</p>
      {effective ? <>
        <Fields rows={[[f.connector, identity(effective.facts.destination.connector_id, labels)],
          [f.connectorRevision, identity(effective.facts.destination.connector_revision, labels)]]} />
        <ul aria-label={labels.effectiveDetails} className="space-y-2">{effective.facts.members.map(member =>
          <li className="min-w-0 rounded border border-slate-700 p-3" key={JSON.stringify([member.device_id, member.point_id, member.tag_id])}>
            <Fields rows={memberRows(member, labels)} /></li>)}</ul>
      </> : <p role="status">{labels.unavailable}</p>}
      {counterRelevant && <p role="status" className="text-sm text-slate-400">{labels.effectiveCounterUnavailable}</p>}
    </section>
    <section aria-label={f.schemaEvidence} className="space-y-3">
      <h4 className="font-medium">{f.schemaEvidence}</h4>
      <Fields rows={[
        [f.readiness, view.readiness.state === 'available' ? yesNo(view.readiness.value.ready, labels) : labels.unavailable],
        [f.schemaEvidence, view.schema.state === 'available' ? identity(view.schema.value.schema_digest ?? view.schema.value.digest, labels) : labels.unavailable],
        [f.operation, identity(operation?.operation_id, labels)],
        [f.operationStatus, knownState(operation?.status, labels)],
        [f.writeOutcome, view.operation.state === 'available' && view.operation.value.kind === 'test_write' ? knownState(view.operation.value.operation.write_outcome, labels) : labels.unavailable],
        [f.cleanupStatus, view.operation.state === 'available' && view.operation.value.kind === 'test_write' ? knownState(view.operation.value.operation.cleanup_status, labels) : labels.unavailable],
        [f.requestedAt, timestamp(runtime?.applied_effective_at, labels)],
      ]} />
    </section>
    <section aria-label={f.delivery} className="space-y-3">
      <h4 className="font-medium">{f.delivery}</h4>
      <p className="text-sm text-slate-400">{labels.deliveryDescription}</p>
      <Fields rows={[
        [f.intake, knownState(current?.intake?.state, labels)],
        [f.collecting, scalar(stages?.collecting, labels)], [f.queued, scalar(stages?.queued, labels)],
        [f.retrying, scalar(stages?.retrying, labels)], [f.blocked, scalar(stages?.blocked, labels)],
        [f.quarantined, scalar(stages?.quarantined, labels)], [f.unknown, scalar(stages?.unknown, labels)],
        [f.sqlCommitted, scalar(stages?.sql_committed, labels)], [f.skipped, scalar(stages?.skipped, labels)],
        [f.lastCommit, timestamp(current?.last_commit?.committed_at, labels)],
        [f.operation, identity(delivery?.test_write?.operation_id, labels)], [f.writtenAt, timestamp(delivery?.test_write?.written_at, labels)],
        [f.readback, timestamp(delivery?.test_write?.readback_at, labels)], [f.cleanup, timestamp(delivery?.test_write?.cleanup_at, labels)],
      ]} />
      {delivery && <ul aria-label={f.backlog} className="space-y-2">{delivery.backlog.map((backlog, index) =>
        <li key={`${backlog.binding.entity_revision}-${index}`}><Fields rows={[
          [labels.fields.savedRevision, identity(backlog.binding.entity_revision, labels)], [f.backlog, scalar(backlog.pending, labels)],
          [f.connector, identity(backlog.destination.connector_id, labels)], [f.connectorRevision, identity(backlog.destination.connector_revision, labels)],
          [labels.fields.destination, [backlog.destination.schema, backlog.destination.table].filter(Boolean).join('.')],
        ]} /></li>)}</ul>}
    </section>
  </div>;
}
