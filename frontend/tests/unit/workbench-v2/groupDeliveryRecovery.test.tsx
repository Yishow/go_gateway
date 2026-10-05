import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupDeliveryStrip } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip';
import { savedGroup } from '../../fixtures/writeGroupState';
const api = vi.hoisted(() => ({ resolveDelivery: vi.fn() }));
const query = vi.hoisted(() => ({ data: undefined as unknown, isError: false, refetch: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: api }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useWriteGroupDeliveryQuery: () => query }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const item = { effect_key: 'effect-old', state: 'blocked', state_revision: 1, payload_digest: 'digest-old', group_revision: 'revision-old', connector_id: 'connector-old', connector_revision: 'connector-revision-old', table_schema: 'main', table_name: 'original_readings', error_code: 'missing-table' };
beforeEach(() => {
 vi.clearAllMocks(); query.isError = false;
 query.data = { stages: { collecting: 0, queued: 1, retrying: 0, blocked: 1, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, no_data_buckets: 0, skipped_buckets: 0, attention: [item] };
 api.resolveDelivery.mockResolvedValue({ decision_id: 'decision', effect_key: item.effect_key, state: 'pending', duplicate: false });
});
it('requires a repair note and explicit skip confirmation, preserving frozen row identity', async () => {
 render(<GroupDeliveryStrip group={savedGroup()} />);
 expect(screen.getByText('main.original_readings')).toBeInTheDocument();
 const retry = screen.getByTestId('delivery-retry-effect-old');
 const skip = screen.getByTestId('delivery-skip-effect-old');
 expect(retry).toBeDisabled(); expect(skip).toBeDisabled();
 fireEvent.change(screen.getByTestId('delivery-reason-effect-old'), { target: { value: 'Repaired the missing table' } });
 expect(skip).toBeDisabled();
 fireEvent.click(retry);
 await waitFor(() => expect(api.resolveDelivery).toHaveBeenCalledTimes(1));
 expect(api.resolveDelivery.mock.calls[0]).toEqual([savedGroup().id, expect.objectContaining({ effect_key: 'effect-old', expected_state: 'blocked', expected_state_revision: 1, payload_digest: 'digest-old', resolution: 'retry', reason: 'Repaired the missing table', confirm_skip: false })]);
 await waitFor(() => expect(query.refetch).toHaveBeenCalled());
});
it('does not offer retry or skip for unknown or cached evidence after refresh error', () => {
 query.data = { ...(query.data as object), attention: [{ ...item, state: 'unknown' }] };
 const view = render(<GroupDeliveryStrip group={savedGroup()} />);
 expect(screen.queryByTestId('delivery-retry-effect-old')).not.toBeInTheDocument();
 expect(screen.getByText('step4.group.delivery.recovery.unknown')).toBeInTheDocument();
 query.isError = true; view.rerender(<GroupDeliveryStrip group={savedGroup()} />);
 expect(screen.queryByTestId('delivery-recovery')).not.toBeInTheDocument();
});
it('an uncertain reply reuses its decision ID and shows a safe error', async () => {
 api.resolveDelivery.mockRejectedValueOnce(new Error('postgres://secret')).mockResolvedValueOnce({ state: 'pending' });
 render(<GroupDeliveryStrip group={savedGroup()} />);
 fireEvent.change(screen.getByTestId('delivery-reason-effect-old'), { target: { value: 'Repaired permission' } });
 fireEvent.click(screen.getByTestId('delivery-retry-effect-old'));
 await screen.findByText('step4.group.delivery.recovery.failed');
 expect(screen.queryByText(/postgres:\/\/secret/)).not.toBeInTheDocument();
 fireEvent.click(screen.getByTestId('delivery-retry-effect-old'));
 await waitFor(() => expect(api.resolveDelivery).toHaveBeenCalledTimes(2));
 expect(api.resolveDelivery.mock.calls[1][1].decision_id).toBe(api.resolveDelivery.mock.calls[0][1].decision_id);
});

it('a newly blocked attempt permits a fresh decision instead of remaining resolved forever', async () => {
 const view = render(<GroupDeliveryStrip group={savedGroup()} />);
 fireEvent.change(screen.getByTestId('delivery-reason-effect-old'), { target: { value: 'First repair' } });
 fireEvent.click(screen.getByTestId('delivery-retry-effect-old'));
 await screen.findByText('step4.group.delivery.recovery.saved');
 query.data = { ...(query.data as object), attention: [{ ...item, state_revision: 2 }] };
 view.rerender(<GroupDeliveryStrip group={savedGroup()} />);
 expect(screen.getByTestId('delivery-reason-effect-old')).not.toBeDisabled();
 expect(screen.queryByText('step4.group.delivery.recovery.saved')).not.toBeInTheDocument();
});
