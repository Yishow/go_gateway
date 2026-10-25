import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { parseWriteGroupSourceEvidence } from '@/utils/studioV2WriteGroupSourceEvidenceJson';
import { CircuitProfileStatusView } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandStatus } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

const APPLIED = 'applied-1';
const stage = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };
function row(start = '2026-10-24T10:00:00Z', end = '2026-10-24T10:00:15Z') {
  return { available: true as const, group_revision: APPLIED, connector_revision: 'crev-1', record_id: 'record-1', effect_key: 'effect-1',
    payload_digest: 'a'.repeat(64), profile_kind: 'instantaneous' as const, window_start: start, window_end: end,
    status: 'partial' as const, raw_present: true, delivery: 'sql_committed' as const, committed_at: '2026-10-24T10:00:20Z' };
}
function evidence(start?: string, end?: string) {
  return { state: 'available' as const, group_revision: APPLIED, connector_id: 'conn-1', connector_revision: 'crev-1',
    sources: [{ circuit_id: 'cir-a', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest: row(start, end), last_sql_committed: null }] };
}
function delivery(sourceEvidence = evidence()): WriteGroupDelivery {
  return { group_id: 'grp-i', intake: { state: 'active' }, stages: stage, revision_stages: { group_revision: APPLIED, stages: stage },
    last_sql_committed_at: '2026-10-24T10:00:20Z', oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0,
    backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    source_evidence: sourceEvidence };
}
function group() {
  return { ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'Circuit A', interval: 15, points: [0, 1] }), applied_revision: APPLIED };
}

describe('Source evidence preserves the applied destination and exact partial window', () => {
  it('accepts a positive one-microsecond end-truncated window instead of rounding it to zero', () => {
    const parsed = parseWriteGroupSourceEvidence(evidence('2026-10-24T10:00:00Z', '2026-10-24T10:00:00.000001Z'));
    expect(parsed.ok).toBe(true);
  });

  it('keeps the frozen applied connector identity in the wire contract', () => {
    const parsed = parseWriteGroupSourceEvidence(evidence());
    expect(parsed).toEqual(expect.objectContaining({ ok: true, value: expect.objectContaining({ connector_id: 'conn-1', connector_revision: 'crev-1' }) }));
  });

  it('shows the true partial window precision instead of two identical second-rounded times', () => {
    render(<CircuitProfileStatusView group={group()} delivery={delivery(evidence('2026-10-24T10:00:00Z', '2026-10-24T10:00:00.000001Z'))} />);
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent('10:00:00.000001');
  });

  it('keeps the applied receipt visible after an unapplied destination draft changes', () => {
    const applied = group();
    const draft = { ...applied, revision: 'draft-2', destination: { ...applied.destination, connector_id: 'draft-connector', connector_revision: 'draft-crev' } };
    render(<CircuitProfileStatusView group={draft} delivery={delivery()} />);
    expandStatus('grp-i');
    const latest = screen.getByTestId('circuit-status-record-latest-tag-0');
    expect(latest).toHaveTextContent('step4.basic.status.receipt_confirmed');
    expect(latest).toHaveTextContent('step4.basic.status.raw_recorded');
  });
});
