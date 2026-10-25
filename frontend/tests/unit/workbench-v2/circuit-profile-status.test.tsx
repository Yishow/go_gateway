import { render, screen, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileStatus } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus';
import { CircuitProfileStatusView } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView';
import type {
  WriteGroupDelivery,
  WriteGroupSourceEvidenceItem,
  WriteGroupSourceRecord,
  WriteGroupSourceRecordUnavailable,
} from '@/types/studioV2WriteGroupDelivery';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandStatus } from '../../fixtures/sourceReview';

const { deliveryQuery, useDelivery } = vi.hoisted(() => {
  const deliveryQuery = { data: undefined as unknown, isError: false, isLoading: false };
  return { deliveryQuery, useDelivery: vi.fn(() => deliveryQuery) };
});
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useWriteGroupDeliveryQuery: useDelivery }));

const APPLIED = 'applied-1';
const instGroup = () => ({ ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }), applied_revision: APPLIED });
const cumGroup = () => ({ ...longProfileGroup({ id: 'grp-c', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 10, points: [0, 1], kind: 'cumulative' }), applied_revision: APPLIED });

function record(over: Partial<WriteGroupSourceRecord> = {}): WriteGroupSourceRecord {
  return {
    available: true, group_revision: APPLIED, connector_revision: 'crev-1', record_id: 'rec-1', effect_key: 'eff-1',
    payload_digest: 'a'.repeat(64), profile_kind: 'instantaneous', window_start: '2026-10-24T10:00:00Z', window_end: '2026-10-24T10:00:15Z',
    status: 'ok', raw_present: true, delivery: 'sql_committed', committed_at: '2026-10-24T10:00:20Z', ...over,
  };
}
function source(tag: string, latest: WriteGroupSourceEvidenceItem['latest'], confirmed: WriteGroupSourceEvidenceItem['last_sql_committed'] = latest): WriteGroupSourceEvidenceItem {
  return { circuit_id: 'cir-a', device_id: 'dev-1', point_id: `pt-${tag.slice(4)}`, tag_id: tag, latest, last_sql_committed: confirmed };
}
function delivery(sources: WriteGroupSourceEvidenceItem[], over: Partial<WriteGroupDelivery> = {}, revision = APPLIED): WriteGroupDelivery {
  const stages = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };
  return {
    group_id: 'grp-i', intake: { state: 'active' }, stages, last_sql_committed_at: '2026-10-24T10:00:20Z', oldest_pending_seconds: 0,
    no_data_buckets: 0, skipped_buckets: 0, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    revision_stages: { group_revision: revision, stages }, source_evidence: { state: 'available', group_revision: revision, connector_id: 'conn-1', connector_revision: 'crev-1', sources }, ...over,
  };
}
const view = (group = instGroup(), data?: WriteGroupDelivery) =>
  render(<CircuitProfileStatusView group={group} delivery={data} sourceLabels={{ 'tag-0': 'Power', 'tag-1': 'Energy' }} />);

describe('Per source status is read from the frozen record and stays separate from SQL evidence', () => {
  it('records the first cumulative raw reading as recorded with an unknown difference, never zero or a failure', () => {
    view(cumGroup(), delivery([source('tag-0', record({ profile_kind: 'cumulative', status: 'baseline', reason: 'first-reading', delta_known: false }))], { group_id: 'grp-c' }));
    expandStatus('grp-c');
    const row = screen.getByTestId('circuit-status-source-tag-0');
    expect(row).toHaveTextContent('Power');
    expect(row).toHaveTextContent('step4.basic.status.raw_recorded');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).toHaveTextContent('step4.basic.status.difference_unknown');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).toHaveTextContent('step4.basic.status.reason.first-reading');
    expect(row.textContent).not.toMatch(/failed|error|write_failed/i);
    expect(screen.getByTestId('circuit-status-delta-tag-0').textContent).not.toMatch(/(^|[^\w.])0([^\w.]|$)/);
  });

  it('shows a known difference only when the record says one was stored', () => {
    view(cumGroup(), delivery([source('tag-0', record({ profile_kind: 'cumulative', delta_known: true }))], { group_id: 'grp-c' }));
    expandStatus('grp-c');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).toHaveTextContent('step4.basic.status.difference_recorded');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).not.toHaveTextContent('difference_unknown');
  });

  it('keeps a missing source as a stored NULL with its reason beside an independent healthy source', () => {
    view(instGroup(), delivery([
      source('tag-0', record({ status: 'missing', reason: 'no-samples', raw_present: false })),
      source('tag-1', record({ status: 'ok' })),
    ]));
    expandStatus('grp-i');
    const missing = screen.getByTestId('circuit-status-source-tag-0');
    const healthy = screen.getByTestId('circuit-status-source-tag-1');
    expect(missing).toHaveTextContent('step4.basic.status.null_stored');
    expect(missing).toHaveTextContent('step4.basic.status.reason.no-samples');
    expect(healthy).toHaveTextContent('step4.basic.status.raw_recorded');
    expect(healthy).not.toHaveTextContent('null_stored');
    expect(document.body.textContent).not.toMatch(/skip_row|partial_row|completeness|incomplete/i);
  });

  it('treats a zero or false reading as present because only the record flag decides', () => {
    view(instGroup(), delivery([source('tag-0', record({ raw_present: true })), source('tag-1', record({ raw_present: true, status: 'partial' }))]));
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-source-tag-0')).toHaveTextContent('step4.basic.status.raw_recorded');
    expect(screen.getByTestId('circuit-status-source-tag-1')).toHaveTextContent('step4.basic.status.raw_recorded');
  });

  it('separates a newer record still waiting locally from the earlier SQL confirmed one', () => {
    const queued = record({ record_id: 'rec-2', effect_key: 'eff-2', window_start: '2026-10-24T10:00:15Z', window_end: '2026-10-24T10:00:30Z', delivery: 'queued', committed_at: undefined });
    view(instGroup(), delivery([source('tag-0', queued, record())]));
    expandStatus('grp-i');
    const latest = screen.getByTestId('circuit-status-record-latest-tag-0');
    const confirmed = screen.getByTestId('circuit-status-record-confirmed-tag-0');
    expect(latest).toHaveTextContent('2026-10-24 10:00:15');
    expect(latest).toHaveTextContent('step4.basic.status.delivery.queued');
    expect(latest).not.toHaveTextContent('receipt_confirmed');
    expect(confirmed).toHaveTextContent('2026-10-24 10:00:00');
    expect(confirmed).toHaveTextContent('step4.basic.status.receipt_confirmed');
  });

  it('never presents a record that could not be verified, or an unknown local state, as written', () => {
    const broken: WriteGroupSourceRecordUnavailable = { available: false, cause: 'digest_mismatch', delivery: 'sql_committed' };
    view(instGroup(), delivery([source('tag-0', broken, null), source('tag-1', record({ delivery: 'unknown', committed_at: undefined }), null)]));
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent('step4.basic.status.record_unavailable');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).not.toHaveTextContent('receipt_confirmed');
    expect(screen.getByTestId('circuit-status-record-latest-tag-1')).toHaveTextContent('step4.basic.status.delivery.unknown');
    expect(screen.queryByTestId('circuit-status-record-confirmed-tag-1')).not.toBeInTheDocument();
  });

  it('shows only the stable reason copy and never renders an unrecognised reason text', () => {
    view(instGroup(), delivery([source('tag-0', record({ status: 'bad', reason: 'driver says <b>secret</b> password=hunter2', raw_present: false })), source('tag-1', record({ status: 'bad', reason: 'reason_unavailable', raw_present: false }))]));
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-source-tag-0')).toHaveTextContent('step4.basic.status.reason_unavailable');
    expect(document.body.textContent).not.toMatch(/hunter2|secret/);
    expect(screen.getByTestId('circuit-status-source-tag-1')).toHaveTextContent('step4.basic.status.reason_unavailable');
  });

  it('says nothing was recorded yet, without baseline, zero or failure, when a source has no record', () => {
    view(instGroup(), delivery([source('tag-0', null, null)]));
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-source-tag-0')).toHaveTextContent('step4.basic.status.source_none');
    expect(screen.getByTestId('circuit-status-source-tag-0').textContent).not.toMatch(/baseline|failed|null_stored/);
  });
});

describe('Stages are separate and evidence of another revision or destination is not shown as current', () => {
  it('keeps saved, applied, collecting, locally durable and SQL committed apart and never claims readback', () => {
    const stages = { collecting: 2, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 };
    view(instGroup(), delivery([source('tag-0', record({ delivery: 'queued', committed_at: undefined }), null)], { stages, last_sql_committed_at: null, revision_stages: { group_revision: APPLIED, stages } }));
    expect(screen.getByTestId('circuit-status-stage-saved')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-status-stage-applied')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-status-stage-collecting')).toHaveTextContent('2');
    expect(screen.getByTestId('circuit-status-stage-local')).toHaveTextContent('1');
    expect(screen.getByTestId('circuit-status-stage-sql')).toHaveTextContent('0');
    expect(screen.queryByTestId('circuit-status-sql-confirmed')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-status-readback')).toHaveTextContent('step4.basic.status.readback_not_checked');
    expect(document.body.textContent).not.toMatch(/readback_verified/);
  });

  it('does not use evidence frozen to another applied revision', () => {
    view(instGroup(), delivery([source('tag-0', record({ group_revision: 'old-applied' }))], {}, 'old-applied'));
    expect(screen.queryByTestId('circuit-status-source-tag-0')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-status-sources-grp-i-toggle')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-status-evidence-note')).toHaveTextContent('step4.basic.status.evidence_not_current');
  });

  it('does not credit a record frozen to a different destination revision', () => {
    view(instGroup(), delivery([source('tag-0', record({ connector_revision: 'other-crev' }))]));
    expandStatus('grp-i');
    expect(screen.getByTestId('circuit-status-source-tag-0')).toHaveTextContent('step4.basic.status.destination_changed');
    expect(screen.getByTestId('circuit-status-source-tag-0')).not.toHaveTextContent('receipt_confirmed');
  });

  it('ignores a delivery answer that belongs to another group', () => {
    view(instGroup(), delivery([source('tag-0', record())], { group_id: 'grp-other' }));
    expect(screen.getByTestId('circuit-status-stage-saved')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-status-source-tag-0')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-status-sources-grp-i-toggle')).not.toBeInTheDocument();
  });

  it('reports unavailable evidence honestly for a non long layout or a missing snapshot', () => {
    view(instGroup(), delivery([], { source_evidence: { state: 'unavailable', unavailable_reason: 'applied_snapshot_unavailable', group_revision: APPLIED, sources: [] } }));
    expect(screen.getByTestId('circuit-status-evidence-note')).toHaveTextContent('step4.basic.status.evidence_unavailable');
  });
});

describe('The container reads the real delivery query of its own group', () => {
  beforeEach(() => { vi.clearAllMocks(); deliveryQuery.data = undefined; deliveryQuery.isError = false; });

  it('queries only an applied group and renders the answer for that group', () => {
    deliveryQuery.data = delivery([source('tag-0', record())]);
    render(<CircuitProfileStatus group={instGroup()} sourceLabels={{ 'tag-0': 'Power' }} />);
    expect(useDelivery).toHaveBeenCalledWith('grp-i', true, expect.any(Number));
    expandStatus('grp-i');
    expect(within(screen.getByTestId('circuit-status-instantaneous')).getByTestId('circuit-status-source-tag-0')).toBeInTheDocument();
  });

  it('does not query a group that was never applied and says so', () => {
    render(<CircuitProfileStatus group={{ ...instGroup(), applied_revision: '' }} />);
    expect(useDelivery).toHaveBeenCalledWith('grp-i', false, expect.any(Number));
    expect(screen.getByTestId('circuit-status-instantaneous')).toHaveTextContent('step4.basic.status.not_applied');
  });

  it('shows a visible query failure instead of an empty or successful state', () => {
    deliveryQuery.isError = true;
    render(<CircuitProfileStatus group={instGroup()} />);
    expect(screen.getByTestId('circuit-status-instantaneous')).toHaveTextContent('step4.basic.status.query_failed');
  });
});
