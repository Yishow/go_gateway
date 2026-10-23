import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { BasicMemberUpdate } from '@/features/datalink/workbench-v2/steps/step4/BasicMemberUpdate';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
const mutateAsync = vi.hoisted(() => vi.fn());
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useUpdateWriteGroupMutation: () => ({ mutateAsync }) }));
const group = savedGroup({ basic_managed_device_id: 'dev-1', members: savedGroup().members.map(m => ({ ...m, target_column: '' })), destination: { ...savedGroup().destination, storage_strategy: 'managed' } });
const props = { group, deviceId: 'dev-1', candidates: buildGroupCandidates(savedGroupState()).candidates, workspaceRevision: 'wrev-1', disabled: false, onUpdated: vi.fn() };
beforeEach(() => { vi.clearAllMocks(); mutateAsync.mockResolvedValue({ group, workspace_revision: 'wrev-2' }); });
it('requires a review and sends one revision-guarded membership update despite double-click', async () => {
  render(<BasicMemberUpdate {...props} />);
  expect(mutateAsync).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  fireEvent.click(screen.getByTestId('basic-members-confirm'));
  fireEvent.click(screen.getByTestId('basic-members-confirm'));
  await waitFor(() => expect(mutateAsync).toHaveBeenCalledTimes(1));
  expect(mutateAsync.mock.calls[0][0]).toMatchObject({ id: group.id, request: { expected_group_revision: group.revision, expected_workspace_revision: 'wrev-1', expected_connector_revision: 'crev-1', group: { members: expect.arrayContaining([{ device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', target_column: '', required: true }]) } } });
});
it('revokes a review after the canonical revision changes and cancels without saving', () => {
  const { rerender } = render(<BasicMemberUpdate {...props} />);
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  rerender(<BasicMemberUpdate {...props} workspaceRevision="newer" />);
  expect(screen.queryByTestId('basic-members-confirm')).not.toBeInTheDocument();
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  fireEvent.click(screen.getByText('step3.bulk.cancel'));
  expect(screen.queryByTestId('basic-members-confirm')).not.toBeInTheDocument();
  expect(mutateAsync).not.toHaveBeenCalled();
});
it('keeps readonly/unknown callers from saving', () => {
  render(<BasicMemberUpdate {...props} disabled />);
  expect(screen.getByTestId('basic-recording-members-review')).toBeDisabled();
  expect(mutateAsync).not.toHaveBeenCalled();
});
