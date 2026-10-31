import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CircuitProfileStatusView } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView';
import type { WriteGroupDelivery, WriteGroupSourceEvidence, WriteGroupLegacySourceRecord } from '@/types/studioV2WriteGroupDelivery';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandStatus } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

const APPLIED = 'applied-1';
const stage = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };
const record = (over: Partial<WriteGroupLegacySourceRecord> = {}): WriteGroupLegacySourceRecord => ({
  available: true, group_revision: APPLIED, connector_revision: 'crev-1', record_id: 'rec-1', effect_key: 'eff-1', payload_digest: 'a'.repeat(64),
  profile_kind: 'instantaneous', window_start: '2026-10-24T10:00:00Z', window_end: '2026-10-24T10:00:15Z', status: 'ok', raw_present: true,
  delivery: 'sql_committed', committed_at: '2026-10-24T10:00:20Z', ...over,
});
const delivery = (data: Partial<WriteGroupSourceEvidence> & { sources?: WriteGroupSourceEvidence['sources'] } = {}): WriteGroupDelivery => ({
  group_id: 'grp-i', intake: { state: 'active' }, stages: stage, revision_stages: { group_revision: APPLIED, stages: stage }, last_sql_committed_at: '2026-10-24T10:00:20Z',
  oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  source_evidence: { state: 'available', group_revision: APPLIED, connector_id: 'conn-1', connector_revision: 'crev-1',
    sources: [{ circuit_id: 'cir-a', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest: record(), last_sql_committed: null }], ...data },
});
const applied = () => ({ ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] }), applied_revision: APPLIED });
const withDraftDestination = (connector_id: string, connector_revision: string) => {
  const group = applied();
  return { ...group, revision: 'draft-2', destination: { ...group.destination, connector_id, connector_revision } };
};

describe('Receipts are judged against the frozen applied destination, not the draft', () => {
  it('keeps the applied receipt and says the new destination draft has not taken effect', () => {
    render(<CircuitProfileStatusView group={withDraftDestination('draft-connector', 'draft-crev')} delivery={delivery()} />);
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent('step4.basic.status.receipt_confirmed');
    expect(screen.getByTestId('circuit-status-applied-destination-note')).toHaveTextContent('step4.basic.status.applied_destination_note');
  });

  it('says nothing about a draft when the destination is unchanged', () => {
    render(<CircuitProfileStatusView group={applied()} delivery={delivery()} />);
    expect(screen.queryByTestId('circuit-status-applied-destination-note')).not.toBeInTheDocument();
  });

  it('never lets an old receipt count for the new destination draft', () => {
    const data = delivery({ sources: [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest: record({ connector_revision: 'draft-crev' }), last_sql_committed: null }] });
    render(<CircuitProfileStatusView group={withDraftDestination('conn-1', 'draft-crev')} delivery={data} />);
    expandStatus('grp-i');
    const latest = screen.getByTestId('circuit-status-record-latest-tag-0');
    expect(latest).toHaveTextContent('step4.basic.status.destination_changed');
    expect(latest).not.toHaveTextContent('receipt_confirmed');
  });

  it('gives no receipt credit when the answer does not name the frozen destination', () => {
    render(<CircuitProfileStatusView group={applied()} delivery={delivery({ connector_id: undefined, connector_revision: undefined })} />);
    expandStatus('grp-i');
    const latest = screen.getByTestId('circuit-status-record-latest-tag-0');
    expect(latest).not.toHaveTextContent('receipt_confirmed');
    expect(latest).toHaveTextContent('step4.basic.status.destination_changed');
  });
});

describe('Window times keep the precision the record has', () => {
  const windowed = (start: string, end: string) => delivery({ sources: [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest: record({ window_start: start, window_end: end }), last_sql_committed: null }] });

  it.each([
    ['whole seconds', '2026-10-24T10:00:00Z', '2026-10-24T10:00:15Z', '2026-10-24 10:00:00 → 2026-10-24 10:00:15 UTC'],
    ['a one microsecond cutover window', '2026-10-24T10:00:00Z', '2026-10-24T10:00:00.000001Z', '2026-10-24 10:00:00 → 2026-10-24 10:00:00.000001 UTC'],
    ['an unaligned cutover start', '2026-10-24T10:00:07.5Z', '2026-10-24T10:00:15Z', '2026-10-24 10:00:07.5 → 2026-10-24 10:00:15 UTC'],
  ])('shows %s', (_name, start, end, shown) => {
    render(<CircuitProfileStatusView group={applied()} delivery={windowed(start, end)} />);
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent(shown);
  });
});
