import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CircuitProfileStatusView } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView';
import type { WriteGroupDelivery, WriteGroupMixedSourceRecord } from '@/types/studioV2WriteGroupDelivery';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandStatus } from '../../fixtures/sourceReview';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

function fixture(method: 'latest' | 'mean' | 'delta', notDue: boolean): WriteGroupDelivery {
  const start = '2026-10-31T00:00:00Z', end = '2026-10-31T00:00:15Z';
  const result: WriteGroupMixedSourceRecord = {
    available: true, profile_kind: 'mixed', group_revision: 'applied-1', connector_revision: 'crev-1', record_id: 'record', effect_key: 'effect', payload_digest: 'a'.repeat(64),
    window_start: notDue ? null : start, window_end: notDue ? null : end, status: notDue ? 'not_due' : 'ok', raw_present: !notDue,
    delivery: 'sql_committed', committed_at: '2026-10-31T00:00:16Z',
    calculation: { result_shape: 'member-calculation-v1', method, interval_seconds: notDue ? 60 : 15, alignment: 'utc_epoch', due_at: end,
      outcome_digest: 'b'.repeat(64), outcome_kind: 'row', status: notDue ? 'not_due' : 'ok', result_present: !notDue, raw_present: !notDue, row_emitted: true,
      source_type: 'int64', result_type: method === 'latest' ? 'int64' : 'decimal', ...(notDue ? { reason: 'not_due' } : { aligned_start: start, aligned_end: end, actual_start: start, actual_end: end }) },
  };
  return { last_sql_committed_at: '2026-10-31T00:00:16Z', group_id: 'grp-i', intake: { state: 'active' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 },
    oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    source_evidence: { state: 'available', group_revision: 'applied-1', connector_id: 'conn-1', connector_revision: 'crev-1',
      sources: [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest: result, last_sql_committed: result }] } };
}
describe('legacy raw status view does not reinterpret mixed calculation evidence', () => {
  it.each(['latest', 'mean', 'delta', 'not_due'] as const)('keeps %s on the existing unavailable/read-only path', scenario => {
    const group = { ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] }), applied_revision: 'applied-1' };
    render(<CircuitProfileStatusView group={group} delivery={fixture(scenario === 'not_due' ? 'latest' : scenario, scenario === 'not_due')} />);
    expandStatus('grp-i');
    const record = screen.getByTestId('circuit-status-record-latest-tag-0');
    expect(record).toHaveTextContent('step4.basic.status.record_unavailable');
    expect(record).not.toHaveTextContent('raw_recorded');
    expect(record).not.toHaveTextContent('status_value.ok');
    expect(record).not.toHaveTextContent('0001');
    expect(record).not.toHaveTextContent('receipt_confirmed');
  });
});
