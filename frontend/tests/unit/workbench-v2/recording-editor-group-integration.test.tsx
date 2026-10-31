import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { DeviceSourcePort } from '@/pages/datalink/workbench-v2/deviceSourceSubmission';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && typeof options === 'object' && 'count' in options ? `${key}:${options.count}` : key) }),
}));

let metadata: { status: string; columns: typeof TABLE_COLUMNS } = { status: 'exists', columns: TABLE_COLUMNS };
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: metadata.columns, status: metadata.status, refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

function list(groups: WriteGroup[], revision = 'wrev-1') {
  return { workspace_id: 'ws-1', workspace_revision: revision, groups };
}

function renderSection(props: Partial<React.ComponentProps<typeof WriteGroupSection>> & { deviceSource?: DeviceSourcePort } = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const utils = render(
    <QueryClientProvider client={client}>
      <WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} />
    </QueryClientProvider>,
  );
  return { client, ...utils };
}

beforeEach(() => {
  vi.clearAllMocks();
  metadata = { status: 'exists', columns: TABLE_COLUMNS };
  api.list.mockResolvedValue(list([]));
  api.delivery.mockResolvedValue({
    group_id: 'group-1', intake: { state: 'not_running' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });
});


describe('actual GroupEditor additive picker integration', () => {
  it('offers one entry, commits only Add N, and restores exact column intent after remove/undo', async () => {
    renderSection({ deviceSource: { submit: vi.fn(), inspect: vi.fn(), observations: () => [] } });
    fireEvent.click(await screen.findByTestId('group-create'));
    expect(screen.queryByTestId('recording-source-open')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-member-include-pt-0')).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep name' } });
    fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
    fireEvent.click(screen.getByRole('button', { name: 'picker.add:1' }));
    expect(screen.getByTestId('group-member-row-pt-0')).toBeVisible();
    fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
    fireEvent.click(screen.getByRole('button', { name: 'picker.removeMember' }));
    expect(screen.queryByTestId('group-member-row-pt-0')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'picker.undo' }));
    expect(screen.getByTestId('group-member-column-pt-0')).toHaveValue('temperature');
    fireEvent.click(screen.getByTestId('recording-points-open'));
    expect(screen.getByRole('checkbox', { name: 'line.temperature' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'line.temperature' })).toBeDisabled();
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.pressure' }));
    fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
    expect(screen.queryByTestId('group-member-row-pt-1')).not.toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Keep name');
    expect(api.create).not.toHaveBeenCalled();
  });
});

it.each(['pending', 'unknown'] as const)('blocks Save only for a selected source with an unsettled %s write and preserves draft input', async status => {
  let unsettled: 'pending' | 'unknown' | undefined;
  const source: DeviceSourcePort = { submit: vi.fn(), inspect: vi.fn(), observations: () => [], range: {
    prepare: vi.fn(), prepareSaved: vi.fn(), submit: vi.fn(), inspect: vi.fn(), observations: () => [],
    sourceWriteState: (_device, point) => point === 'pt-0' ? unsettled : undefined,
  } };
  renderSection({ deviceSource: source });
  fireEvent.click(await screen.findByTestId('group-create'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep draft' } });
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
  fireEvent.click(screen.getByRole('button', { name: 'picker.add:1' }));
  fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
  expect(screen.getByTestId('group-save')).toBeEnabled();
  unsettled = status;
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep latest draft' } });
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.getByTestId('group-source-unsettled')).toHaveTextContent('line.temperature');
  expect(screen.getByTestId('group-member-column-pt-0')).toHaveValue('temperature');
  expect(api.create).not.toHaveBeenCalled();
  unsettled = undefined;
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Now resolved' } });
  expect(screen.getByTestId('group-save')).toBeEnabled();
});
it('rechecks pending source writes after the Save revision read, not only at click time', async () => {
  let pending = false;
  const source: DeviceSourcePort = { submit: vi.fn(), inspect: vi.fn(), observations: () => [], range: {
    prepare: vi.fn(), prepareSaved: vi.fn(), submit: vi.fn(), inspect: vi.fn(), observations: () => [],
    sourceWriteState: (_device, point) => pending && point === 'pt-0' ? 'pending' : undefined,
  } };
  renderSection({ deviceSource: source });
  fireEvent.click(await screen.findByTestId('group-create'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep draft' } });
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
  fireEvent.click(screen.getByRole('button', { name: 'picker.add:1' }));
  fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
  let resolve!: (value: ReturnType<typeof list>) => void;
  api.list.mockImplementationOnce(() => new Promise(yes => { resolve = yes; }));
  const reads = api.list.mock.calls.length;
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(reads + 1));
  pending = true;
  await act(async () => { resolve(list([])); });
  expect(api.create).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-name')).toHaveValue('Keep draft');
});

it('does not treat a known unavailable or inactive source as an ambiguous write when saving an offline draft', async () => {
  let selectable = true;
  const source: DeviceSourcePort = { submit: vi.fn(), inspect: vi.fn(), observations: () => [], range: {
    prepare: vi.fn(), prepareSaved: vi.fn(), submit: vi.fn(), inspect: vi.fn(), observations: () => [],
    canSelect: () => selectable, sourceWriteState: () => undefined,
  } };
  renderSection({ deviceSource: source });
  fireEvent.click(await screen.findByTestId('group-create'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Offline draft' } });
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
  fireEvent.click(screen.getByRole('button', { name: 'picker.add:1' }));
  fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
  selectable = false;
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep offline draft' } });
  expect(screen.getByTestId('group-save')).toBeEnabled();
  expect(screen.queryByTestId('group-source-unsettled')).not.toBeInTheDocument();
});
