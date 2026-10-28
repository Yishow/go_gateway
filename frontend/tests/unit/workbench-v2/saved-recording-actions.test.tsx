import { SavedRecordingManagement, RemovedRecordingManagement } from '@/features/datalink/workbench-v2/steps/step4/SavedRecordingManagement';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

const mocks = vi.hoisted(() => ({ lifecycle: vi.fn(), refetch: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), list: vi.fn(), readiness: vi.fn(), delivery: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => options && 'count' in options ? `${key}:${options.count}` : key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({ useWriteGroupLifecycleQuery: mocks.lifecycle }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: mocks }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: [], status: 'exists', refetch: vi.fn() }) }));

function snapshot(over: Partial<WriteGroupLifecycle> = {}): WriteGroupLifecycle {
  return {
    workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'group-1', group_revision: 'rev-1', applied_revision: '', status: 'ready', connector_id: 'conn-1', connector_revision: 'crev-1',
    actions: ['edit', 'stop', 'remove', 'recreate', 'transition'].map((action) => ({ action: action as WriteGroupLifecycle['actions'][number]['action'], allowed: true, reason_code: '', next_action: '', requires_review: action === 'recreate' || action === 'transition' })),
    history: { canonical: true, current_group_id: 'group-1', generations: [] }, ...over,
  };
}
let current: WriteGroupLifecycle;
const wrap = (element: React.ReactNode) => render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>{element}</QueryClientProvider>);
const bar = (extra = {}) => <GroupLifecycleBar group={savedGroup()} workspaceId="ws-1" workspaceRevision="wrev-1" canApply readonly={false} onChanged={vi.fn()} onConflict={vi.fn()} {...extra} />;
beforeEach(() => {
  vi.clearAllMocks(); current = snapshot();
  mocks.lifecycle.mockImplementation(() => ({ data: current, isError: false, isLoading: false, isFetching: false, refetch: mocks.refetch }));
  mocks.refetch.mockImplementation(async () => ({ data: current }));
  mocks.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [savedGroup({ status: 'deleted' })] });
  mocks.readiness.mockResolvedValue({ group_revision: 'rev-1', ready: true, issues: [] });
  mocks.delivery.mockRejectedValue(new Error('unavailable'));
});

describe('saved recording lifecycle authority', () => {
  it('uses server stop/remove availability and actionable reasons instead of status guesses', () => {
    current.actions[1] = { action: 'stop', allowed: false, reason_code: 'successor_reserved', next_action: 'review_lifecycle', requires_review: false };
    wrap(bar());
    expect(screen.getByTestId('group-disable')).toBeDisabled();
    expect(screen.getByTestId('saved-recording-reason-stop')).toHaveTextContent('reason.successor_reserved');
    expect(screen.getByTestId('saved-recording-reason-stop')).toHaveTextContent('next.review_lifecycle');
  });
  it('fails closed for foreign or stale lifecycle scope', () => {
    current = snapshot({ workspace_id: 'foreign' }); wrap(bar());
    expect(screen.getByTestId('group-disable')).toBeDisabled();
    expect(screen.getByTestId('group-delete')).toBeDisabled();
  });
  it('refreshes authority before sending a stop; a newly blocked answer is not mutation permission', async () => {
    wrap(bar());
    mocks.refetch.mockResolvedValue({ data: snapshot({ actions: current.actions.map((a) => ({ ...a, allowed: false, reason_code: 'successor_reserved' })) }) });
    fireEvent.click(screen.getByTestId('group-disable'));
    await waitFor(() => expect(mocks.refetch).toHaveBeenCalled());
    expect(mocks.disable).not.toHaveBeenCalled();
  });
  it('does not treat retained cached facts as a fresh answer when refetch fails', async () => {
    mocks.refetch.mockResolvedValue({ data: current, isError: true, error: new Error('private transport detail') });
    wrap(bar()); fireEvent.click(screen.getByTestId('group-disable'));
    await waitFor(() => expect(mocks.refetch).toHaveBeenCalled());
    await waitFor(() => expect(screen.getByText('reason.stale_scope')).toBeInTheDocument());
    expect(mocks.disable).not.toHaveBeenCalled();
    expect(screen.queryByText('private transport detail')).not.toBeInTheDocument();
  });
  it('opens a deleted item requested by Basic without treating it as a new ordinary draft', async () => {
    wrap(<WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} openGroupRequest={{ id: 'group-1' }} />);
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toBeDisabled();
    expect(screen.getByTestId('group-save')).toBeDisabled();
  });
  it('shows separate samples and outbox units for each immutable revision when removing', async () => {
    current.history.generations = [{ group_id: 'group-1', generation: 1, current: true, reserved: false, status: 'ready', group_revision: 'rev-1', applied_revision: 'a1', saved_destination: { connector_id: 'conn-1', connector_revision: 'crev-1', destination_scope: 'scope' }, revisions: [{ group_revision: 'a1', effective_at: '2026-10-01T00:00:00Z', sources: [], destination: { connector_id: 'conn-1', connector_revision: 'crev-1', destination_scope: 'scope' }, pending: { samples: 11, outbox: { outstanding: 7, total: 9, queued: 2, retrying: 1, blocked: 3, unknown: 1, quarantined: 0, sql_committed: 2, skipped: 0 }, handoffs: [] } }] }];
    wrap(bar()); fireEvent.click(screen.getByTestId('group-delete'));
    const dialog = await screen.findByRole('alertdialog');
    expect(within(dialog).getByTestId('lifecycle-pending-samples')).toHaveTextContent('11');
    expect(within(dialog).getByTestId('lifecycle-pending-outbox')).toHaveTextContent('7');
    expect(within(dialog).getByTestId('lifecycle-pending-blocked')).toHaveTextContent('3');
    expect(within(dialog).getByTestId('lifecycle-pending-unknown')).toHaveTextContent('1');
  });
});


describe('saved action caller parity and interruption', () => {
  it('Basic and Advanced show the same allowed actions and server blockers', () => {
    current.actions[2] = { action: 'remove', allowed: false, reason_code: 'connector_identity_stale', next_action: 'stop', requires_review: false };
    const host = { onAction: vi.fn() };
    wrap(<><section data-testid="basic"><SavedRecordingManagement group={savedGroup()} workspaceId="ws-1" workspaceRevision="wrev-1" lifecycleHost={host} /></section><section data-testid="advanced">{bar({ lifecycleHost: host })}</section></>);
    const basic = within(screen.getByTestId('basic')); const advanced = within(screen.getByTestId('advanced'));
    for (const id of ['saved-recording-edit', 'group-disable', 'group-delete', 'saved-recording-recreate', 'saved-recording-transition']) {
      expect(basic.getByTestId(id).hasAttribute('disabled')).toBe(advanced.getByTestId(id).hasAttribute('disabled'));
      expect(basic.getByTestId(id).textContent).toBe(advanced.getByTestId(id).textContent);
    }
    expect(basic.getByTestId('saved-recording-reason-remove').textContent).toBe(advanced.getByTestId('saved-recording-reason-remove').textContent);
  });
  it('keeps reason IDs unique and described-by links local when both modes show the same group', () => {
    current.actions = current.actions.map((action) => ({ ...action, allowed: false, reason_code: 'successor_reserved', next_action: 'review_lifecycle' }));
    const host = { onAction: vi.fn() };
    wrap(<><section data-testid="basic"><SavedRecordingManagement group={savedGroup()} workspaceId="ws-1" workspaceRevision="wrev-1" lifecycleHost={host} /></section><section data-testid="advanced">{bar({ lifecycleHost: host })}</section></>);
    const reasons = screen.getAllByTestId(/^saved-recording-reason-/);
    expect(new Set(reasons.map((reason) => reason.id)).size).toBe(reasons.length);
    for (const surface of ['basic', 'advanced']) {
      const region = screen.getByTestId(surface);
      for (const button of within(region).getAllByRole('button').filter((button) => button.hasAttribute('aria-describedby'))) {
        const description = document.getElementById(button.getAttribute('aria-describedby')!);
        expect(description).not.toBeNull();
        expect(region).toContainElement(description);
        expect(description).toHaveTextContent('reason.successor_reserved');
      }
    }
  });
  it('the Basic tombstone entry emits recreate intent, not ordinary ensure/create, and stays inspectable readonly', () => {
    const onAction = vi.fn(); const onOpenAdvanced = vi.fn(); const group = savedGroup({ status: 'deleted' });
    const view = wrap(<RemovedRecordingManagement groups={[group]} deviceId="dev-1" workspaceId="ws-1" workspaceRevision="wrev-1" lifecycleHost={{ onAction }} onOpenAdvanced={onOpenAdvanced} />);
    fireEvent.click(screen.getByText('removed.title:1'));
    fireEvent.click(screen.getByTestId('saved-recording-recreate'));
    expect(onAction).toHaveBeenCalledWith({ action: 'recreate', workspaceId: 'ws-1', groupId: 'group-1', workspaceRevision: 'wrev-1', groupRevision: 'rev-1' });
    view.unmount();
    wrap(<RemovedRecordingManagement groups={[group]} deviceId="dev-1" workspaceId="ws-1" workspaceRevision="wrev-1" readonly lifecycleHost={{ onAction }} onOpenAdvanced={onOpenAdvanced} />);
    fireEvent.click(screen.getByText('removed.title:1'));
    expect(screen.getByTestId('saved-recording-recreate')).toBeDisabled();
    fireEvent.click(screen.getByTestId('basic-removed-open-group-1'));
    expect(onOpenAdvanced).toHaveBeenCalledWith('group-1');
  });
  it('Advanced tombstone readonly protects fields but does not block recreate review', async () => {
    const onAction = vi.fn();
    wrap(<WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} lifecycleHost={{ onAction }} openGroupRequest={{ id: 'group-1' }} />);
    await screen.findByTestId('group-editor');
    expect(screen.getByTestId('group-name')).toBeDisabled();
    expect(screen.getByTestId('saved-recording-recreate')).not.toBeDisabled();
    fireEvent.click(screen.getByTestId('saved-recording-recreate'));
    expect(onAction).toHaveBeenCalledWith(expect.objectContaining({ action: 'recreate', groupId: 'group-1' }));
  });
  it('routes a managed profile ordinary start for review without bypassing canApply', () => {
    const onReviewScopedStart = vi.fn();
    const group = savedGroup({ destination: { ...savedGroup().destination, storage_strategy: 'managed' }, layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 10, unit: 'second', sample_reduction: 'latest' } });
    const view = wrap(bar({ group, onReviewScopedStart }));
    fireEvent.click(screen.getByTestId('group-apply'));
    expect(onReviewScopedStart).toHaveBeenCalledWith(group); expect(mocks.apply).not.toHaveBeenCalled();
    view.unmount(); wrap(bar({ group, onReviewScopedStart, canApply: false }));
    expect(screen.getByTestId('group-apply')).toBeDisabled();
  });
  it('does not leave the same mounted bar busy when its revision changes during a read', async () => {
    let release: (value: { data: WriteGroupLifecycle }) => void = () => undefined;
    mocks.refetch.mockReturnValue(new Promise((resolve) => { release = resolve; }));
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const view = render(<QueryClientProvider client={client}>{bar()}</QueryClientProvider>);
    fireEvent.click(screen.getByTestId('group-disable'));
    current = snapshot({ group_revision: 'rev-2' });
    view.rerender(<QueryClientProvider client={client}>{bar({ group: savedGroup({ revision: 'rev-2' }) })}</QueryClientProvider>);
    await act(async () => { release({ data: snapshot() }); await Promise.resolve(); });
    expect(screen.getByTestId('group-disable')).not.toBeDisabled();
    expect(mocks.disable).not.toHaveBeenCalled();
  });
  it('cancel restores focus without stopping or removing any recording', async () => {
    wrap(bar()); const trigger = screen.getByTestId('group-delete'); trigger.focus(); fireEvent.click(trigger);
    const dialog = await screen.findByRole('alertdialog');
    fireEvent.click(within(dialog).getByTestId('group-delete-cancel'));
    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument());
    await waitFor(() => expect(trigger).toHaveFocus());
    expect(mocks.remove).not.toHaveBeenCalled(); expect(mocks.disable).not.toHaveBeenCalled();
  });
});
