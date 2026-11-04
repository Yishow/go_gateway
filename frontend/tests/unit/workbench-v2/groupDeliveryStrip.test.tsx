import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { GroupDeliveryStrip } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip';
import { useWriteGroupDeliveryQuery } from '@/hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { savedGroup } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupDeliveryQuery: vi.fn(),
}));

const deliveryQuery = vi.mocked(useWriteGroupDeliveryQuery);

function delivery(overrides: Partial<WriteGroupDelivery> = {}): WriteGroupDelivery {
  return {
    group_id: 'group-1',
    intake: { state: 'active' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null,
    oldest_pending_seconds: 0,
    no_data_buckets: 3,
    skipped_buckets: 2,
    recent_bucket_issues: [],
    backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    ...overrides,
  };
}

beforeEach(() => {
  vi.clearAllMocks();
  deliveryQuery.mockReturnValue({ data: delivery(), isError: false } as never);
});

describe('GroupDeliveryStrip', () => {
  it('shows skipped and no-data bucket counts as separate delivery facts', () => {
    render(<GroupDeliveryStrip group={savedGroup()} />);

    expect(screen.getByTestId('delivery-skipped-buckets')).toHaveTextContent('2');
    expect(screen.getByTestId('delivery-no-data-buckets')).toHaveTextContent('3');
  });

  it('does not treat cached delivery as confirmed after a refresh error', () => {
    deliveryQuery.mockReturnValue({
      data: delivery({
        recent_bucket_issues: [{
          group_revision: 'rev-1', bucket_start: '2026-10-03T00:00:00Z', kind: 'skipped', causes: ['missing'],
        }],
      }),
      isError: true,
    } as never);
    render(<GroupDeliveryStrip group={savedGroup({ applied_revision: 'rev-applied' })} />);

    for (const testId of [
      'delivery-collecting', 'delivery-skipped-buckets', 'delivery-no-data-buckets',
      'delivery-local', 'delivery-attention', 'delivery-committed',
    ]) {
      expect(screen.getByTestId(testId)).toHaveTextContent('step4.group.delivery.unconfirmed');
    }
    expect(screen.queryByTestId('delivery-bucket-causes')).not.toBeInTheDocument();
    expect(screen.getByTestId('delivery-saved')).toHaveTextContent('step4.group.delivery.yes');
    expect(screen.getByTestId('delivery-applied')).toHaveTextContent('step4.group.delivery.yes');
    expect(screen.getByTestId('delivery-applied-revision')).toHaveTextContent('rev-applied');
    expect(screen.getByTestId('group-delivery-strip').querySelector('details')).toBeNull();
  });

  it('keeps delivery facts unconfirmed before the first response', () => {
    deliveryQuery.mockReturnValue({ data: undefined, isError: false, isLoading: true } as never);
    render(<GroupDeliveryStrip group={savedGroup()} />);

    for (const testId of [
      'delivery-collecting', 'delivery-skipped-buckets', 'delivery-no-data-buckets',
      'delivery-local', 'delivery-attention', 'delivery-committed',
    ]) {
      expect(screen.getByTestId(testId)).toHaveTextContent('step4.group.delivery.unconfirmed');
    }
    expect(screen.queryByTestId('delivery-bucket-causes')).not.toBeInTheDocument();
  });

  it('labels recent skipped causes with safe localized copy', () => {
    deliveryQuery.mockReturnValue({
      data: delivery({
        recent_bucket_issues: [{
          group_revision: 'rev-1', bucket_start: '2026-10-03T00:00:00Z', kind: 'skipped',
          causes: ['missing', 'bad', 'stale'],
        }],
      }),
      isError: false,
    } as never);
    render(<GroupDeliveryStrip group={savedGroup()} />);

    const causes = screen.getByTestId('delivery-bucket-causes');
    expect(causes).toBeInTheDocument();
    const issue = screen.getByTestId('delivery-recent-issue-0');
    expect(issue).toHaveTextContent('step4.group.delivery.issue.skipped');
    expect(issue).toHaveTextContent('step4.group.delivery.cause.missing');
    expect(issue).toHaveTextContent('step4.group.delivery.cause.bad');
    expect(issue).toHaveTextContent('step4.group.delivery.cause.stale');
    expect(issue).not.toHaveTextContent('private-dsn');
  });
});

it('shows compact delivery figures for an applied record without folding, while errors and bucket risks remain visible', () => {
  deliveryQuery.mockReturnValue({ data: delivery({ recent_bucket_issues: [{ group_revision: 'rev-1', bucket_start: '2026-10-03T00:00:00Z', kind: 'skipped', causes: ['missing'] }] }), isError: false } as never);
  const view = render(<GroupDeliveryStrip group={{ ...savedGroup(), applied_revision: 'rev-1' }} {...{ compact: true }} />);
  expect(screen.getByTestId('group-delivery-diagnostics').querySelector('dl')).toHaveClass('recording-kpi-strip');
  expect(screen.getByTestId('delivery-collecting')).toBeVisible();
  expect(screen.getByTestId('delivery-bucket-causes')).toBeVisible();
  deliveryQuery.mockReturnValue({ data: undefined, isError: true } as never);
  view.rerender(<GroupDeliveryStrip group={{ ...savedGroup(), applied_revision: 'rev-1' }} {...{ compact: true }} />);
  expect(screen.getByTestId('group-delivery-failed')).toBeVisible();
});


it('summarizes a confirmed not-running record without presenting eight zero cells as progress', () => {
  deliveryQuery.mockReturnValue({ data: delivery({ intake: { state: 'not_running' } }), isError: false } as never);
  render(<GroupDeliveryStrip group={savedGroup()} compact />);
  expect(screen.getByTestId('recording-delivery-summary')).toHaveTextContent('deliverySummary.notRunning');
  expect(screen.getByTestId('recording-delivery-summary')).toHaveTextContent('step4.group.delivery.yes');
  expect(screen.queryByTestId('delivery-collecting')).not.toBeInTheDocument();
  expect(screen.queryByTestId('recording-delivery-history')).not.toBeInTheDocument();
});

it('keeps committed history and pending rows visible even without an applied revision', () => {
  const data = delivery({ intake: { state: 'not_running' }, last_sql_committed_at: '2026-11-02T01:00:00Z' });
  data.stages.sql_committed = 9; data.stages.queued = 2;
  deliveryQuery.mockReturnValue({ data, isError: false } as never);
  render(<GroupDeliveryStrip group={savedGroup({ applied_revision: '' })} compact />);
  expect(screen.getByTestId('recording-delivery-history')).toBeVisible();
  expect(screen.getByTestId('recording-delivery-history')).toHaveTextContent('9');
  expect(screen.getByTestId('recording-delivery-pending')).toBeVisible();
  expect(screen.getByTestId('recording-delivery-pending')).toHaveTextContent('2');
});

it('does not confirm runtime or history from a foreign delivery response', () => {
  const data = delivery({ group_id: 'foreign-group' }); data.stages.sql_committed = 12;
  deliveryQuery.mockReturnValue({ data, isError: false } as never);
  render(<GroupDeliveryStrip group={savedGroup()} compact />);
  expect(screen.getByTestId('recording-delivery-summary')).toHaveTextContent('deliverySummary.unconfirmed');
  expect(screen.queryByTestId('recording-delivery-history')).not.toBeInTheDocument();
  // An unapplied record shows no figure cells at all, so a foreign 12 can never appear as committed rows.
  expect(screen.queryByTestId('delivery-committed')).not.toBeInTheDocument();
  expect(screen.queryByText('12')).not.toBeInTheDocument();
});
