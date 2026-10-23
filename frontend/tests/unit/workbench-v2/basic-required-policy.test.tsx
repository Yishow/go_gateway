import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, o?: Record<string, unknown>) => o && 'count' in o ? `${key}:${o.count}` : key }) }));
const state = savedGroupState();
const managed = savedGroup({ basic_managed_device_id: 'dev-1', destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_group_server' } });
const query = { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [managed] }, isLoading: false, isError: false };
const mutate = vi.fn();
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => query, useWriteGroupDeliveryQuery: () => ({}), useUpdateWriteGroupMutation: () => ({ mutateAsync: mutate }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: mutate }), useRecordingStartMutation: () => ({ mutateAsync: mutate }), useRecordingStartOperationQuery: () => ({}),
}));
beforeEach(() => { vi.clearAllMocks(); query.data.groups = [managed]; sessionStorage.clear(); });
function panel(onOpenMemberSettings = vi.fn()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><BasicRecordingPanel state={state} workspaceId="ws-1" readonly={false} {...{onOpenMemberSettings}} /></QueryClientProvider>);
  return onOpenMemberSettings;
}
it('explains required gating and opens only the selected saved member settings without mutating members', () => {
  const before = structuredClone(managed.members);
  const open = panel();
  fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
  expect(screen.getByTestId('basic-required-policy')).toHaveTextContent('basic.required_hint');
  expect(screen.getByTestId('basic-required-policy')).toHaveTextContent(`basic.required_count:${managed.members.filter(m => m.required).length}`);
  fireEvent.click(screen.getByTestId('basic-member-settings'));
  expect(open).toHaveBeenCalledWith(managed.id);
  expect(managed.members).toEqual(before);
  expect(mutate).not.toHaveBeenCalled();
});
it('shows required defaults before preparation and does not offer an unrelated group', () => {
  query.data.groups = [];
  panel();
  expect(screen.getByTestId('basic-required-policy')).toHaveTextContent('basic.prepare_first');
  expect(screen.queryByTestId('basic-member-settings')).not.toBeInTheDocument();
  expect(mutate).not.toHaveBeenCalled();
});
