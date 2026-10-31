import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingRuleDetails } from '@/features/datalink/workbench-v2/steps/step4/RecordingRuleDetails';
import { RecordingRuleSummary } from '@/features/datalink/workbench-v2/steps/step4/RecordingRuleSummary';
import { recordingGuidedLabels } from '@/features/datalink/workbench-v2/steps/step4/recordingGuidedLabels';
import { recordingWorkflowView, type RecordingWorkflowView } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import en from '@/i18n/locales/en/workbench-v2-recording-guided.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2-recording-guided.json';

const labels = recordingGuidedLabels(key => key.split('.').reduce<unknown>((node, part) =>
  typeof node === 'object' && node !== null && part in node ? (node as Record<string, unknown>)[part] : undefined, en) as string);
function source(): WriteGroup {
  return { id: 'G1', workspace_id: 'W1', revision: 'R2', applied_revision: 'R1', name: 'Private group name', status: 'ready',
    members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_revision: 'S2', mapping_revision: 'M2',
      target_column: 'reading', required: true, max_age_seconds: 45, source_value_type: 'uint64', result_value_type: 'decimal', value_role: 'counter_raw',
      exact_recording_source: { measurement_id: 'measurement-1', definition_revision: 'definition-2', source_binding_revision: 'binding-3',
        series_epoch: 'epoch-4', original_data_type: 'uint64', scale: '0.0001', unit: 'kWh', semantic_kind: 'counter',
        counter_policy: { allow_negative: false, reset_threshold: 0, max_roll_over: 99, rollover_modulus: '18446744073709551616', max_rate_per_second: '100000.001' } } }],
    destination: { connector_id: 'DB1', connector_revision: 'D1', database: '/private/archive.sqlite', table_schema: 'public', table_name: 'readings',
      storage_strategy: 'managed', schema_revision: 'schema-1', schema_digest: 'digest-1' },
    row_policy: { interval_seconds: 30, allowed_lateness_seconds: 5, incomplete_policy: 'partial', group_key_columns: ['device_id', 'circuit_id'],
      unique_key_columns: ['record_key'], entity_key_column: 'entity_key', value_column: 'reading', quality_column: 'quality', provenance_column: 'origin',
      record_key_column: 'record_key', bucket_start_column: 'bucket_start', group_id_column: 'group_id', device_id_column: 'device_id' },
    write_policy: { mode: 'append', recreate_missing: false }, migration: {}, created_at: '', updated_at: '' };
}
function view(group = source(), evidence?: RuleEvidenceView): RecordingWorkflowView {
  const scope = { workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', groupId: 'G1',
    workspaceRevision: 'W2', groupRevision: group.revision, appliedRevision: group.applied_revision };
  return recordingWorkflowView({ scope, savedGroup: group, runtimeEvidence: evidence ? { state: 'available', scope, value: evidence } : undefined });
}
function field(region: HTMLElement, name: string) { return within(region).getByText(name, { selector: 'dt' }).nextElementSibling; }
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

describe('read-only saved configuration and verified effective technical details', () => {
  it('projects canonical mixed member identity without mutating its calculation intent', () => {
    const group = source();
    group.recording_profile = { kind: 'mixed', result_shape: 'member-calculation-v1' };
    group.record_id = '12345678-1234-4234-8234-123456789abc';
    delete group.row_policy.interval_seconds;
    group.members[0].calculation = 'delta';
    group.members[0].cadence = { interval: 15, unit: 'second', alignment: 'utc_epoch' };
    const before = structuredClone(group);
    render(<RecordingRuleDetails view={view(group)} labels={labels} />);
    const saved = screen.getByRole('region', { name: labels.savedDetails });
    expect(saved).toHaveTextContent('P1'); expect(saved).toHaveTextContent('reading');
    expect(group).toEqual(before);
  });
  it('labels the opaque Tag ID without conflating it with a semantic Tag key', () => {
    expect(en.detailFields.tag).toBe('Tag ID'); expect(zh.detailFields.tag).toBe('Tag ID');
  });
  it('retains exact saved source identities, types, bindings, counter policy and epoch without representing them as effective proof', () => {
    render(<RecordingRuleDetails view={view()} labels={labels} />);
    const saved = screen.getByRole('region', { name: labels.savedDetails });
    for (const value of ['T1', 'P1', 'S2', 'M2', 'reading', 'UInt64', 'decimal', 'measurement-1', 'definition-2', 'binding-3',
      'epoch-4', '0.0001', 'kWh', '18446744073709551616', '100000.001', 'schema-1', 'digest-1']) {
      expect(saved.textContent?.toLowerCase()).toContain(value.toLowerCase());
    }
    expect(field(saved, labels.detailFields.required)).toHaveTextContent('Yes');
    expect(field(saved, labels.detailFields.maxAge)).toHaveTextContent('45');
    expect(field(saved, labels.detailFields.resetThreshold)).toHaveTextContent('0');
    const effective = screen.getByRole('region', { name: labels.effectiveDetails });
    expect(effective).toHaveTextContent(labels.effectiveCounterUnavailable);
    expect(effective).not.toHaveTextContent(/epoch-4|binding-3|definition-2|18446744073709551616/);
    expect(document.body).not.toHaveTextContent(/private|archive.sqlite|Private group name/);
  });
  it('retains saved row key and column mapping policies rather than hiding them in a second editor', () => {
    render(<RecordingRuleDetails view={view()} labels={labels} />);
    const saved = screen.getByRole('region', { name: labels.savedDetails });
    for (const name of ['entity_key', 'record_key', 'bucket_start', 'group_id', 'quality', 'origin']) expect(saved).toHaveTextContent(name);
    expect(saved.querySelector('input, select, textarea, button')).toBeNull();
  });
  it('does not add counter semantics to ordinary non-counter inspection', () => {
    const group = source(); delete group.members[0].exact_recording_source; group.members[0].value_role = 'gauge';
    render(<RecordingRuleDetails view={view(group)} labels={labels} />);
    expect(document.body).not.toHaveTextContent(labels.effectiveCounterUnavailable);
    expect(document.body).not.toHaveTextContent(labels.detailFields.epoch);
  });
  it('allowlists editable details without exposing raw diagnostics or private path values', () => {
    const group = source(); group.members[0].source_revision = '/private/source.key'; group.members[0].target_column = 'password=hidden';
    group.members[0].exact_recording_source!.unit = 'postgres://private.internal';
    group.members[0].exact_recording_source!.counter_policy!.rollover_modulus = 'secret=hidden';
    group.destination.schema_digest = 'SELECT private FROM secret';
    render(<RecordingRuleDetails view={view(group)} labels={labels} />);
    expect(field(screen.getByRole('region', { name: labels.savedDetails }), labels.detailFields.tag)).toHaveTextContent(/^T1$/);
    expect(document.body).not.toHaveTextContent(/private|source.key|password|hidden|SELECT|postgres:/);
  });
  it.each(['gateway.internal', 'gateway.internal/metrics', '10.23.0.1', 'service/path'])('does not treat opaque tag_id %s as a semantic Tag key', (tagId) => {
    const group = source(); group.members[0].tag_id = tagId;
    render(<RecordingRuleDetails view={view(group)} labels={labels} />);
    const saved = screen.getByRole('region', { name: labels.savedDetails });
    expect(field(saved, labels.detailFields.tag)).toHaveTextContent(labels.unavailable);
    expect(document.body).not.toHaveTextContent(tagId);
  });
  it('also suppresses endpoint-like opaque tag_id in verified effective technical details', () => {
    const evidence = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'A' });
    const runtime = evidence.groups[0]; runtime.group_id = 'G1'; runtime.saved_revision = 'R2'; runtime.applied_revision = 'R1'; runtime.effective_revision = 'R1';
    runtime.saved!.binding.entity_revision = 'R2'; runtime.effective!.binding.entity_revision = 'R1';
    runtime.effective!.members[0].tag_id = 'gateway.internal';
    render(<RecordingRuleDetails view={view(source(), evidence)} labels={labels} />);
    const effective = screen.getByRole('region', { name: labels.effectiveDetails });
    expect(field(effective, labels.detailFields.tag)).toHaveTextContent(labels.unavailable);
    expect(document.body).not.toHaveTextContent('gateway.internal');
    expect(field(screen.getByRole('region', { name: labels.savedDetails }), labels.detailFields.tag)).toHaveTextContent(/^T1$/);
  });
  it('keeps unknown evidence unavailable rather than inventing successful schema, delivery or operation facts', () => {
    render(<RecordingRuleDetails view={view()} labels={labels} />);
    const delivery = screen.getByRole('region', { name: labels.detailFields.delivery });
    expect(field(delivery, labels.detailFields.sqlCommitted)).toHaveTextContent(labels.unavailable);
    expect(field(delivery, labels.detailFields.readback)).toHaveTextContent(labels.unavailable);
    expect(field(delivery, labels.detailFields.cleanup)).toHaveTextContent(labels.unavailable);
    const schema = screen.getByRole('region', { name: labels.detailFields.schemaEvidence });
    expect(field(schema, labels.detailFields.operationStatus)).toHaveTextContent(labels.unavailable);
    expect(field(schema, labels.detailFields.readiness)).toHaveTextContent(labels.unavailable);
  });
  it('uses each verified effective source and separate durable/test-write stage evidence, never the newer saved counter proof', () => {
    const evidence = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'A' });
    const runtime = evidence.groups[0]; runtime.group_id = 'G1'; runtime.saved_revision = 'R2'; runtime.applied_revision = 'R1'; runtime.effective_revision = 'R1';
    runtime.saved!.binding.entity_revision = 'R2'; runtime.effective!.binding.entity_revision = 'R1';
    runtime.saved!.members[0] = { device_id: 'A', point_id: 'P1', tag_id: 'T1', column: 'reading', source_value_type: 'uint64', circuit_id: 'C1' };
    runtime.effective!.members[0] = { device_id: 'A', point_id: 'P-old', tag_id: 'T-old', column: 'old_reading', source_value_type: 'int16', circuit_id: 'C1' };
    runtime.delivery = { counter_replay_quarantines: null, oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 8,
      backlog: [], current: { binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R1' }, totals_scope: 'group', shared: false,
        intake: { group_revision: 'R1', state: 'active', reason: null },
        stages: { collecting: 0, queued: 2, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 },
        last_commit: { group_revision: 'R1', connector_revision: '7', record_id: 'record-1', effect_key: 'effect-1', payload_digest: 'digest-1', committed_at: '2026-10-27T09:00:01Z', provenance: 'runtime' }, sources: [], source_state: 'available', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
      test_write: { operation_id: 'test-1', provenance: 'test-write', written_at: '2026-10-27T09:00:00Z', readback_at: null, cleanup_at: '2026-10-27T09:00:03Z' } };
    render(<RecordingRuleDetails view={view(source(), evidence)} labels={labels} />);
    const effective = screen.getByRole('region', { name: labels.effectiveDetails });
    expect(effective).toHaveTextContent('T-old'); expect(effective).toHaveTextContent('old_reading'); expect(effective).toHaveTextContent('Int16');
    expect(effective).not.toHaveTextContent(/readingSources.*Uint64|binding-3|epoch-4/);
    const delivery = screen.getByRole('region', { name: labels.detailFields.delivery });
    expect(field(delivery, labels.detailFields.queued)).toHaveTextContent('2'); expect(field(delivery, labels.detailFields.sqlCommitted)).toHaveTextContent('1');
    expect(field(delivery, labels.detailFields.readback)).toHaveTextContent(labels.unavailable);
    expect(field(delivery, labels.detailFields.cleanup)).toHaveTextContent('2026-10-27T09:00:03Z');
    expect(field(delivery, labels.detailFields.lastCommit)).toHaveTextContent('2026-10-27T09:00:01Z');
  });
  it('performs no requests or storage writes on render, drawer open, close, rerender or unmount', () => {
    const request = vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Forbidden request'));
    const storage = vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => undefined);
    const current = view(); if (current.state !== 'available') throw new Error('Expected owned view');
    const content = <RecordingRuleDetails view={current} labels={labels} />;
    const props = { view: current, labels, scopeLabel: 'Panel A', details: { scope: current.scope, asOf: null, effectiveRevision: null, content } };
    const mounted = render(<RecordingRuleSummary {...props} />);
    fireEvent.click(screen.getByRole('button', { name: labels.details }));
    expect(screen.getByRole('dialog')).toHaveTextContent('binding-3');
    fireEvent.click(screen.getByRole('button', { name: labels.close }));
    mounted.rerender(<RecordingRuleSummary {...props} />); mounted.unmount();
    expect(request).not.toHaveBeenCalled(); expect(storage).not.toHaveBeenCalled();
  });
});
