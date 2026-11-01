import { render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingMemberResult } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingMemberResult';
import { parseMixedSourceRecord } from '../../../src/utils/studioV2MixedSourceEvidenceJson';
import { savedGroup } from '../../fixtures/writeGroupState';
import type { WriteGroupDelivery, WriteGroupMixedSourceRecord } from '../../../src/types/studioV2WriteGroupDelivery';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const group = savedGroup({ applied_revision: 'applied-1', record_id: '11111111-1111-4111-8111-111111111111', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } });
const record: WriteGroupMixedSourceRecord = { available: true, group_revision: 'applied-1', connector_revision: 'frozen-connector', record_id: '22222222-2222-4222-8222-222222222222', effect_key: 'E', payload_digest: 'a'.repeat(64), profile_kind: 'mixed', window_start: null, window_end: null, status: 'not_due', raw_present: false, delivery: 'queued', delta_known: false, calculation: { result_shape: 'member-calculation-v1', method: 'delta', interval_seconds: 60, alignment: 'utc_epoch', due_at: '2026-10-31T00:00:15Z', next_due: '2026-10-31T00:01:00Z', outcome_digest: 'b'.repeat(64), outcome_kind: 'row', status: 'not_due', reason: 'not_due', result_present: false, row_emitted: true, raw_present: false, source_type: 'uint32', result_type: 'decimal' } };
function dueResult(present: boolean): WriteGroupMixedSourceRecord {
 const value = structuredClone(record); value.status = present ? 'ok' : 'baseline'; value.raw_present = true; value.delta_known = present;
 value.window_start = '2026-10-31T00:00:05Z'; value.window_end = '2026-10-31T00:00:55Z';
 Object.assign(value.calculation, { status: present ? 'ok' : 'missing', reason: present ? undefined : 'first-reading', result_present: present, raw_present: true, due_at: '2026-10-31T00:01:00Z', next_due: '2026-10-31T00:02:00Z', aligned_start: '2026-10-31T00:00:00Z', aligned_end: '2026-10-31T00:01:00Z', actual_start: value.window_start, actual_end: value.window_end });
 return value;
}
function evidence(value = record): WriteGroupDelivery {
 const canonical = parseMixedSourceRecord({ ...value, window_start: value.window_start ?? '0001-01-01T00:00:00Z', window_end: value.window_end ?? '0001-01-01T00:00:00Z' }, () => null);
 expect(canonical, 'fixture must satisfy the real canonical producer parser').not.toBeNull();
 return { group_id: group.id, intake: { state: 'not_running' }, stages: { collecting: 0, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false, loss_risk_notice: "" }, source_evidence: { state: 'available', group_revision: 'applied-1', connector_revision: 'frozen-connector', sources: [{ ...group.members[0], latest: value, last_sql_committed: null }] } }; }
it('separates not-due NULL from raw provenance and source health; never fabricates a short window or committed result', () => {
 render(<RecordingMemberResult group={group} delivery={evidence()} candidates={[{ key: 'pt-0|tag-0', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', tag_key: 'temperature', label: 'Temperature', device_name: 'PLC A', address: '40001', target_type: 'float64' }]} />);
 expect(screen.getByRole('heading', { name: 'PLC A · Temperature · 40001' })).toBeInTheDocument();
 expect(screen.getByTestId('member-result-pt-0-latest')).toHaveTextContent('result.not_due');
 expect(screen.getByTestId('member-result-pt-0-latest')).toHaveTextContent('result.null');
 expect(screen.getByText('2026-10-31T00:01:00Z')).toBeInTheDocument();
 expect(screen.queryByText('result.committed')).not.toBeInTheDocument(); expect(screen.queryByTestId('member-result-window-pt-0-latest')).not.toBeInTheDocument();
});
it('renders only canonical aligned and actual bounds with first-baseline reason, not a raw or zero substitute', () => {
 const value = dueResult(false);
 render(<RecordingMemberResult group={group} delivery={evidence(value)} />);
 expect(screen.getByTestId('member-result-window-pt-0-latest')).toHaveTextContent('2026-10-31T00:00:00Z');
 expect(screen.getByText('result.raw')).toBeInTheDocument(); expect(screen.getByTestId('member-result-pt-0-latest')).toHaveTextContent('result.null');
 expect(screen.getByText('result.firstBaseline')).toBeInTheDocument(); expect(screen.getByText('2026-10-31T00:00:05Z')).toBeInTheDocument();
});
it('fails closed on foreign applied revision instead of showing stale success', () => {
 render(<RecordingMemberResult group={group} delivery={evidence({ ...record, group_revision: 'foreign-applied-revision' })} />);
 expect(screen.getByText('result.unavailable')).toBeInTheDocument(); expect(screen.queryByTestId('member-result-pt-0-latest')).not.toBeInTheDocument();
});
it('uses canonical result_present for a valid zero-capable result, never raw presence or value truthiness', () => {
 const value = dueResult(true);
 render(<RecordingMemberResult group={group} delivery={evidence(value)} />);
 expect(screen.getByTestId('member-result-pt-0-latest')).toHaveTextContent('result.present'); expect(screen.getByTestId('member-result-pt-0-latest')).not.toHaveTextContent('result.null');
 expect(screen.getByText('result.raw')).toBeInTheDocument();
});
