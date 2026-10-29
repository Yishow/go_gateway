import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SavedRecordingActionHost } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
const mocks = vi.hoisted(() => ({ lifecycle: vi.fn(), list: vi.fn(), basic: vi.fn(), advanced: vi.fn(), support: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: mocks.lifecycle } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: mocks.list, get: vi.fn() } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel', () => ({ BasicRecordingPanel: (props: { readonly: boolean; lifecycleHost: SavedRecordingActionHost; onOpenAdvanced: () => void }) => {
  mocks.basic(props);
  return <><button disabled={props.readonly} onClick={() => props.lifecycleHost.onAction({ action: 'recreate', workspaceId: 'ws-1', groupId: 'group-1', workspaceRevision: 'wrev-1', groupRevision: 'rev-1' })}>Basic recreate</button><button onClick={props.onOpenAdvanced}>Advanced mode</button></>;
} }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection', () => ({ WriteGroupSection: (props: { readonly: boolean; lifecycleHost: SavedRecordingActionHost }) => {
  mocks.advanced(props);
  return <button disabled={props.readonly} onClick={() => props.lifecycleHost.onAction({ action: 'transition', workspaceId: 'ws-1', groupId: 'group-1', workspaceRevision: 'wrev-1', groupRevision: 'rev-1' })}>Advanced transition</button>;
} }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels', () => ({ Step4SupportPanels: (props: unknown) => { mocks.support(props); return null; } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/CommitSummary', () => ({ CommitSummary: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/ShareOutputSummary', () => ({ ShareOutputSummary: () => null }));
const group = savedGroup({ status: 'deleted', layout_version: 'managed_wide_v1', destination: { ...savedGroup().destination, storage_strategy: 'managed' } });
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear();
  mocks.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [group] });
  mocks.lifecycle.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: group.id, group_revision: group.revision, connector_id: 'conn-1', connector_revision: 'crev-1', status: group.status, applied_revision: '', actions: [{ action: 'recreate', allowed: true, reason_code: '', next_action: '', requires_review: true }, { action: 'transition', allowed: true, reason_code: '', next_action: '', requires_review: true }], history: { canonical: true, lineage_id: 'lineage', lineage_revision: 1, current_group_id: group.id, generations: [] }, successor_seed: { basis_kind: 'saved', basis_revision: group.revision, basis_digest: 'basis', group } });
});
function mount() {
  return render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}><Step4Database state={savedGroupState()} dispatch={vi.fn()} workspaceId="ws-1" /></QueryClientProvider>);
}
describe('one lifecycle host across Basic and Advanced', () => {
  it('keeps the original selection, locks both modes, and does not open another operation when switching modes', async () => {
    mount(); fireEvent.click(screen.getByText('Basic recreate'));
    expect(await screen.findByRole('dialog')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByDisplayValue('Line A')).toBeInTheDocument());
    fireEvent.click(screen.getByText('Advanced mode'));
    expect(screen.getAllByRole('dialog')).toHaveLength(1);
    expect(screen.getByText('Basic recreate')).toBeDisabled();
    expect(screen.getByText('Advanced transition')).toBeDisabled();
    fireEvent.click(screen.getByText('Advanced transition'));
    expect(mocks.lifecycle).toHaveBeenCalledTimes(1);
    expect(mocks.support.mock.lastCall?.[0]).toMatchObject({ disabled: true });
  });
  it('closing before submission restores the host without performing a mutation', async () => {
    mount(); fireEvent.click(screen.getByText('Basic recreate'));
    await screen.findByDisplayValue('Line A');
    fireEvent.click(screen.getByRole('button', { name: 'dialog.cancel' }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(screen.getByText('Basic recreate')).not.toBeDisabled();
    expect(screen.getByText('Advanced transition')).not.toBeDisabled();
    expect(mocks.lifecycle).toHaveBeenCalledTimes(1);
  });
});
