import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
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
