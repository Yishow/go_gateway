import { render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { CircuitProfileStatusView } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const group = { ...longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] }), applied_revision: 'applied-1' };
const delivery: WriteGroupDelivery = {
  group_id: group.id, intake: { state: 'active' }, last_sql_committed_at: null,
  stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 },
  oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  storage: { generation: 2, history_may_be_lost: true, reconstructions: 1, recreate_missing: true, events: [] },
};

it('keeps the lost-history warning visible after automatic reconstruction and later SQL commits', () => {
  render(<CircuitProfileStatusView group={group} delivery={delivery} />);
  expect(screen.getByRole('alert')).toHaveTextContent('step4.basic.status.storage_history_uncertain');
  expect(screen.getByTestId('circuit-status-storage-generation')).toHaveTextContent('2');
});

it('does not attach a different groups reconstruction warning to this group', () => {
  render(<CircuitProfileStatusView group={group} delivery={{ ...delivery, group_id: 'other' }} />);
  expect(screen.queryByTestId('circuit-status-storage-history')).not.toBeInTheDocument();
});

it('does not warn that history was lost for a new table or legacy answer without storage evidence', () => {
  render(<CircuitProfileStatusView group={group} delivery={{ ...delivery, storage: undefined }} />);
  expect(screen.queryByTestId('circuit-status-storage-history')).not.toBeInTheDocument();
});
