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
    expect(screen.getByText('rev-applied').closest('details')).not.toBeNull();
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
