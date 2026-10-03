import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { CommitSummary } from '@/features/datalink/workbench-v2/steps/step4/CommitSummary';
import { proposeColumns } from '@/features/datalink/workbench-v2/state/writeGroup/proposal';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && typeof options === 'object' && 'count' in options ? `${key}:${options.count}` : key) }),
}));
const columns = [{ name: 'id', type: 'integer', nullable: false, primary_key: true }, { name: 'note', type: 'text', nullable: true, primary_key: false }];
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn() },
}));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] });
});

describe('Fewer columns than points / new table proposal', () => {
  it('names proposals as proposals and assigns nothing when the table has too few columns', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={client}><WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-create'));
    fireEvent.click(await screen.findByTestId('group-bulk-include'));
    const panel = await screen.findByTestId('group-column-proposal');
    expect(panel).toHaveTextContent('proposal.label');
    expect(screen.getByTestId('group-column-proposal-list').children).toHaveLength(3);
    for (const id of ['pt-0', 'pt-1', 'pt-2']) expect(screen.getByTestId(`group-member-column-${id}`)).toHaveValue('');
    expect(screen.getByTestId('group-save')).toBeDisabled();
  });

  it('proposes distinct, valid names that avoid existing columns', () => {
    const cand = (key: string, tagKey: string, t: 'float64' | 'string') => ({ key, device_id: 'd', point_id: key, tag_id: key, tag_key: tagKey, label: tagKey, device_name: 'p', address: '1', target_type: t });
    const result = proposeColumns([cand('a', 'line.Temp-In', 'float64'), cand('b', 'other.temp_in', 'float64'), cand('c', 'x.1st', 'string')], ['temp_in']);
    expect(result.map((p) => p.column)).toEqual(['temp_in_2', 'temp_in_3', 'v_1st']);
    expect(result.map((p) => p.sql_type)).toEqual(['DOUBLE PRECISION', 'DOUBLE PRECISION', 'TEXT']);
  });
});

describe('Share only', () => {
  const props = { deviceCount: 1, ruleCount: 1, pointCount: 2, mappingCount: 2, connector: savedGroupState().db.connector, onActivate: vi.fn() };
  it('activates with no database output at all when nothing blocks', () => {
    render(<CommitSummary {...props} groupSummary={{ state: 'ready', total: 0, applied: 0 }} readinessSummary={{ ready: true, blocking_count: 0, warning_count: 0, issues: [] }} />);
    expect(screen.getByRole('button', { name: /step4.activate_btn/ })).not.toBeDisabled();
    expect(screen.getByTestId('step4-group-summary-note')).toHaveTextContent('summary.none');
  });
  it('still stops on a Share readiness blocker', () => {
    render(<CommitSummary {...props} groupSummary={{ state: 'ready', total: 0, applied: 0 }}
      readinessSummary={{ ready: false, blocking_count: 1, warning_count: 0, issues: [{ code: 'modbus-share-hydration-required', severity: 'blocking', step: 'Step 4', scope: 'share', message: 'Share is not hydrated' }] }} />);
    expect(screen.getByRole('button', { name: /step4.activate_btn/ })).toBeDisabled();
  });
});
