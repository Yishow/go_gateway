import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider, initReactI18next } from 'react-i18next';
import { afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import { SavedRecordingCard } from '@/features/datalink/workbench-v2/steps/step4/SavedRecordingCard';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import workbench from '@/i18n/locales/en/workbench-v2.json';
import lifecycleLabels from '@/i18n/locales/en/recording-lifecycle.json';
import guided from '@/i18n/locales/en/workbench-v2-recording-guided.json';
import { savedGroup } from '../../fixtures/writeGroupState';
const api = vi.hoisted(() => ({ lifecycle: vi.fn(), list: vi.fn(), delivery: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), ensureCircuitProfile: vi.fn(), ensureBasic: vi.fn(), create: vi.fn(), update: vi.fn(), start: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: api }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: api.lifecycle } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: api.start, ensureBasic: api.ensureBasic } }));
const i18n = createInstance();
const clients: QueryClient[] = [];
const original = savedGroup({ id: 'G1', name: 'Line A recording', layout_version: 'managed_wide_v1',
  destination: { ...savedGroup().destination, table_name: 'line_readings', storage_strategy: 'managed' } });
function snapshot(group = original): WriteGroupLifecycle {
  const destination = { connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision, destination_scope: 'original-scope' };
  return { workspace_id: 'ws-1', workspace_revision: 'W1', group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, status: group.status, connector_id: destination.connector_id, connector_revision: destination.connector_revision,
    actions: (['edit', 'stop', 'remove', 'recreate', 'transition'] as const).map(action => ({ action, allowed: true, reason_code: '', next_action: '', requires_review: action === 'recreate' || action === 'transition' })),
    history: { canonical: true, lineage_id: 'lineage-G1', lineage_revision: 1, current_group_id: group.id,
      generations: [{ group_id: group.id, generation: 1, current: true, reserved: false, status: group.status, group_revision: group.revision,
        applied_revision: group.applied_revision, saved_destination: destination, revisions: [{ group_revision: 'original-applied', effective_at: '2026-10-28T00:00:00Z', sources: [], destination,
          pending: { samples: 11, outbox: { outstanding: 1, total: 4, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 3, skipped: 0 }, handoffs: [] } }] }] } };
}
function delivery(): WriteGroupDelivery {
  return { group_id: 'G1', intake: { state: 'active' }, stages: { collecting: 11, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 3, skipped: 0 },
    last_sql_committed_at: '2026-10-28T00:00:00Z', oldest_pending_seconds: 1, no_data_buckets: 0, skipped_buckets: 0,
    backlog: [{ group_revision: 'original-applied', connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: 'line_readings', pending: 1, error_codes: [] }],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false } };
}
let current: WriteGroupLifecycle;
beforeAll(async () => { await i18n.use(initReactI18next).init({ lng: 'en', fallbackLng: 'en', interpolation: { escapeValue: false }, resources: {
  en: { 'workbench-v2': workbench, 'recording-lifecycle': lifecycleLabels, 'workbench-v2-recording-guided': guided } } }); });
beforeEach(() => { vi.clearAllMocks(); current = snapshot();
  api.lifecycle.mockImplementation(async () => current); api.delivery.mockResolvedValue(delivery());
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: [original] });
  api.disable.mockResolvedValue({ workspace_revision: 'W2', group: { ...original, status: 'disabled' } });
  api.remove.mockResolvedValue({ workspace_revision: 'W2', group: { ...original, status: 'deleted' } }); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
function host(group = original, managed: boolean | 'without-host' = true, readonly = false) {
  const onOpen = vi.fn(), onAction = vi.fn(), onReload = vi.fn(), tryAcquire = vi.fn(() => true);
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const management = { workspaceId: 'ws-1', workspaceRevision: 'W1', readonly, lifecycleHost: managed === 'without-host' ? undefined : { onAction, tryAcquire }, onReload };
  const tree = (next: WriteGroup) => <QueryClientProvider client={client}><I18nextProvider i18n={i18n}><ul>
    <SavedRecordingCard group={next} current={false} onOpen={onOpen} management={managed ? management : undefined} />
  </ul></I18nextProvider></QueryClientProvider>;
  const view = render(tree(group));
  return { ...view, onOpen, onAction, onReload, tryAcquire, client, rerenderGroup: (next: WriteGroup) => view.rerender(tree(next)) };
}
function noNewRecording() {
  for (const mutation of [api.apply, api.ensureCircuitProfile, api.ensureBasic, api.create, api.update, api.start]) expect(mutation).not.toHaveBeenCalled();
  expect(screen.queryByTestId('group-apply')).not.toBeInTheDocument();
}
const ready = async () => { await screen.findByTestId('group-list-delivery-G1'); };

describe('discoverable saved-recording card through the existing lifecycle owner', () => {
  it('exposes the server-authorized lifecycle actions directly on saved G1 without a second apply/start', async () => {
    host(); await ready();
    for (const id of ['saved-recording-edit', 'group-disable', 'group-delete', 'saved-recording-recreate', 'saved-recording-transition']) expect(screen.getByTestId(id)).toBeEnabled();
    expect(screen.getByText(/Already accepted backlog still goes to its original destination/)).toBeInTheDocument(); noNewRecording();
  });
  it.each([false, 'without-host'] as const)('keeps Open pure navigation and does not synthesize a management host when management is %s', async managed => {
    const view = host(original, managed); await ready(); fireEvent.click(screen.getByTestId('group-open-G1'));
    expect(view.onOpen).toHaveBeenCalledOnce(); expect(api.lifecycle).not.toHaveBeenCalled(); expect(api.list).not.toHaveBeenCalled();
    expect(screen.queryByTestId('saved-recording-actions')).not.toBeInTheDocument(); noNewRecording();
  });
  it('uses safe collapsed name/table fallbacks rather than exposing credentials or private database paths', async () => {
    const raw = savedGroup({ ...original, name: 'password=secret', destination: { ...original.destination, table_name: '/private/customer.db', database: '/private/database.db' } });
    host(raw, false); await ready();
    expect(document.body).not.toHaveTextContent(/password|secret|private|customer\.db|database\.db/);
    expect(document.body).toHaveTextContent('Saved recording'); expect(document.body).toHaveTextContent('Destination table unavailable');
  });
  it('keeps the original mismatch and delivery evidence without confusing queued Q1 with committed rows', async () => {
    host(original, false); await ready(); expect(screen.getByTestId('group-list-mismatch-G1')).toBeInTheDocument();
    const summary = screen.getByTestId('group-list-delivery-G1'); expect(summary).toHaveTextContent('3'); expect(summary).toHaveTextContent('1');
    expect(api.delivery).toHaveBeenCalledExactlyOnceWith('G1');
  });
});

describe('saved-card lifecycle integration and negative authority paths', () => {
  it.each(['wide', 'long', 'custom'] as const)('keeps %s saved settings unchanged and delegates only exact server-scoped management intent', async kind => {
    const group = kind === 'wide' ? original : kind === 'long' ? savedGroup({ ...original, layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 10, unit: 'second' } })
      : savedGroup({ ...original, layout_version: undefined, destination: { ...original.destination, storage_strategy: 'custom' } });
    const before = structuredClone(group); current = snapshot(group); const view = host(group); await ready();
    await waitFor(() => expect(screen.getByTestId('saved-recording-edit')).toBeEnabled());
    for (const [id, action] of [['saved-recording-edit', 'edit'], ['saved-recording-recreate', 'recreate'], ['saved-recording-transition', 'transition']]) {
      fireEvent.click(screen.getByTestId(id)); expect(view.onAction).toHaveBeenLastCalledWith({ action, workspaceId: 'ws-1', groupId: 'G1', workspaceRevision: 'W1', groupRevision: 'rev-1' });
    }
    fireEvent.click(screen.getByTestId('group-open-G1')); expect(view.onOpen).toHaveBeenCalledOnce(); expect(group).toEqual(before);
    expect(api.list).not.toHaveBeenCalled(); noNewRecording();
  });
  it.each(['foreign workspace', 'foreign group', 'stale workspace revision', 'stale group revision', 'missing availability', 'query error'] as const)('fails closed with actionable reasons for %s', async invalid => {
    if (invalid === 'foreign workspace') current.workspace_id = 'foreign';
    if (invalid === 'foreign group') current.group_id = 'foreign';
    if (invalid === 'stale workspace revision') current.workspace_revision = 'old';
    if (invalid === 'stale group revision') current.group_revision = 'old';
    if (invalid === 'missing availability') current.actions = [];
    if (invalid === 'query error') api.lifecycle.mockRejectedValue(new Error('postgres://private:secret@host/database'));
    host(); await ready(); await waitFor(() => expect(api.lifecycle).toHaveBeenCalledOnce());
    expect(screen.getByTestId('group-disable')).toBeDisabled(); expect(screen.getByTestId('group-delete')).toBeDisabled();
    expect(screen.getByTestId('saved-recording-reason-stop')).toHaveTextContent('Current action availability is not confirmed');
    expect(document.body).not.toHaveTextContent(/postgres:\/\/|private:secret/); noNewRecording();
  });
  it('renders server blockers with local accessible descriptions rather than guessing from the saved status', async () => {
    current.actions[1] = { action: 'stop', allowed: false, reason_code: 'successor_reserved', next_action: 'review_lifecycle', requires_review: false };
    host(); await ready(); const stop = screen.getByTestId('group-disable'); expect(stop).toBeDisabled();
    const reason = document.getElementById(stop.getAttribute('aria-describedby') ?? '');
    expect(reason).toHaveTextContent('A recording change is already reserved'); expect(reason).toHaveTextContent('Review the existing recording change');
    expect(screen.getByTestId('saved-recording-edit')).toBeEnabled();
  });
  it('discloses G1 queued Q1 original destination/history and separate units, then cancel returns focus without mutation', async () => {
    host(); await ready(); const trigger = screen.getByTestId('group-delete'); trigger.focus(); fireEvent.click(trigger);
    const dialog = await screen.findByRole('alertdialog');
    expect(dialog).toHaveTextContent('Accepted backlog continues to its original destination'); expect(dialog).toHaveTextContent('SQL tables, recorded history and delivery receipts are kept');
    expect(within(dialog).getByTestId('lifecycle-pending-samples')).toHaveTextContent('11');
    expect(within(dialog).getByTestId('lifecycle-pending-outbox')).toHaveTextContent('1');
    expect(dialog).toHaveTextContent('Original connection'); expect(dialog).toHaveTextContent('conn-1');
    expect(api.remove).not.toHaveBeenCalled(); fireEvent.click(within(dialog).getByTestId('group-delete-cancel'));
    await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()); await waitFor(() => expect(trigger).toHaveFocus());
    expect(api.remove).not.toHaveBeenCalled(); expect(api.disable).not.toHaveBeenCalled(); noNewRecording();
  });
  it('removes only after explicit confirmation and freshly matching lifecycle CAS, using the parent list reload', async () => {
    const view = host(); await ready(); fireEvent.click(screen.getByTestId('group-delete'));
    const dialog = await screen.findByRole('alertdialog'); expect(api.remove).not.toHaveBeenCalled();
    expect(api.lifecycle).toHaveBeenCalledTimes(2);
    api.remove.mockImplementationOnce(async () => {
      expect(api.lifecycle).toHaveBeenCalledTimes(3); // Initial, review and explicit confirm each inspect the existing owner.
      return { workspace_revision: 'W2', group: { ...original, status: 'deleted' } };
    });
    fireEvent.click(within(dialog).getByTestId('group-delete-confirm-button'));
    await waitFor(() => expect(view.onReload).toHaveBeenCalledOnce());
    expect(api.lifecycle).toHaveBeenCalledTimes(4); // Unchanged native mutation invalidation refreshes the fourth read.
    expect(api.remove).toHaveBeenCalledExactlyOnceWith('G1', { workspace_id: 'ws-1', expected_workspace_revision: 'W1', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' });
    expect(api.list).not.toHaveBeenCalled(); expect(current.history.generations[0].revisions[0].pending.outbox.queued).toBe(1); noNewRecording();
  });
  it('refreshes stop authority before CAS and repeated clicks dispatch only one original stop', async () => {
    const view = host(); await ready(); await waitFor(() => expect(screen.getByTestId('group-disable')).toBeEnabled());
    let release!: (snapshot: WriteGroupLifecycle) => void; const fresh = new Promise<WriteGroupLifecycle>(resolve => { release = resolve; });
    api.lifecycle.mockReturnValueOnce(fresh); fireEvent.click(screen.getByTestId('group-disable')); fireEvent.click(screen.getByTestId('group-disable'));
    expect(api.disable).not.toHaveBeenCalled(); await waitFor(() => expect(api.lifecycle).toHaveBeenCalledTimes(2));
    await act(async () => { release(current); }); await waitFor(() => expect(view.onReload).toHaveBeenCalledOnce());
    expect(api.disable).toHaveBeenCalledExactlyOnceWith('G1', { workspace_id: 'ws-1', expected_workspace_revision: 'W1', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' });
    expect(view.tryAcquire).toHaveBeenCalledOnce(); expect(api.list).not.toHaveBeenCalled(); noNewRecording();
  });
  it.each(['new blocker', 'refetch error'] as const)('does not use cached allowed stop facts after %s', async failure => {
    const view = host(); await ready(); await waitFor(() => expect(screen.getByTestId('group-disable')).toBeEnabled());
    if (failure === 'new blocker') api.lifecycle.mockResolvedValueOnce({ ...current, actions: current.actions.map(action => ({ ...action, allowed: false })) });
    else api.lifecycle.mockRejectedValueOnce(new Error('password=private-secret'));
    fireEvent.click(screen.getByTestId('group-disable')); await waitFor(() => expect(view.onReload).toHaveBeenCalledOnce());
    expect(api.disable).not.toHaveBeenCalled(); expect(screen.getByRole('alert')).toHaveTextContent('Refresh and review again');
    expect(document.body).not.toHaveTextContent('private-secret'); noNewRecording();
  });
  it('does not dispatch an old lifecycle response after the displayed group revision changes', async () => {
    const view = host(); await ready(); let release!: (snapshot: WriteGroupLifecycle) => void;
    api.lifecycle.mockReturnValueOnce(new Promise<WriteGroupLifecycle>(resolve => { release = resolve; }));
    fireEvent.click(screen.getByTestId('group-disable')); await waitFor(() => expect(api.lifecycle).toHaveBeenCalledTimes(2));
    view.rerenderGroup({ ...original, revision: 'rev-2' }); await act(async () => { release(current); });
    expect(api.disable).not.toHaveBeenCalled(); expect(view.onReload).not.toHaveBeenCalled(); noNewRecording();
  });
  it('keeps readonly settings inspectable without enabling lifecycle mutation or conversion', async () => {
    const view = host(original, true, true); await ready(); fireEvent.click(screen.getByTestId('group-open-G1'));
    expect(view.onOpen).toHaveBeenCalledOnce(); expect(screen.getByTestId('group-disable')).toBeDisabled(); expect(screen.getByTestId('group-delete')).toBeDisabled();
    expect(screen.getByTestId('saved-recording-management-unavailable-G1')).toHaveTextContent('You can still view saved settings and delivery evidence');
    expect(api.disable).not.toHaveBeenCalled(); expect(api.remove).not.toHaveBeenCalled(); noNewRecording();
  });
});
