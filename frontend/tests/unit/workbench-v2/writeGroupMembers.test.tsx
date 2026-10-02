import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && 'count' in options ? `${key}:${options.count}` : options && 'column' in options ? `${key}:${options.column}` : key) }),
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

async function openNewGroup() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  fireEvent.click(await screen.findByTestId('group-create'));
  await screen.findByTestId('group-editor');
}

beforeEach(() => {
  vi.clearAllMocks();
  metadata = { status: 'exists', columns: TABLE_COLUMNS };
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] });
});

describe('RealMetadataAndReviewableAssignment (editor)', () => {
  it('proposes columns from the tag name as reviewable suggestions and assigns nothing on its own', async () => {
    await openNewGroup();
    expect(screen.getByTestId('group-member-suggested-pt-0')).toHaveTextContent('temperature');
    expect(screen.getByTestId('group-member-suggested-pt-1')).toHaveTextContent('pressure');
    expect(screen.getByTestId('group-member-suggested-pt-2')).toHaveTextContent('batch');
    for (const id of ['pt-0', 'pt-1', 'pt-2']) expect(screen.getByTestId(`group-member-include-${id}`)).not.toBeChecked();
    // The proposed column names a real column of the table, never a sample column.
    expect(screen.queryByText('temp_in_c')).not.toBeInTheDocument();
  });

  it('accepts only the suggestions of the rows shown and says how many that is', async () => {
    await openNewGroup();
    expect(screen.getByTestId('group-bulk-accept')).toHaveTextContent(':3');
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: 'pressure' } });
    expect(screen.getByTestId('group-bulk-accept')).toHaveTextContent(':1');
    expect(screen.getByTestId('group-member-scope')).toHaveTextContent('step4.group.members.scope');
    fireEvent.click(screen.getByTestId('group-bulk-accept'));
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: '' } });
    expect(screen.getByTestId('group-member-include-pt-1')).toBeChecked();
    expect(screen.getByTestId('group-member-column-pt-1')).toHaveValue('pressure');
    expect(screen.getByTestId('group-member-include-pt-0')).not.toBeChecked();
    expect(screen.getByTestId('group-member-include-pt-2')).not.toBeChecked();
  });

  it('bulk add and remove touch only the rows the filter shows', async () => {
    await openNewGroup();
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: 'line.b' } });
    expect(screen.getByTestId('group-bulk-include')).toHaveTextContent(':1');
    fireEvent.click(screen.getByTestId('group-bulk-include'));
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: '' } });
    expect(screen.getByTestId('group-member-include-pt-2')).toBeChecked();
    expect(screen.getByTestId('group-member-include-pt-0')).not.toBeChecked();
    fireEvent.change(screen.getByTestId('group-member-search'), { target: { value: 'line.t' } });
    expect(screen.getByTestId('group-bulk-exclude')).toHaveTextContent(':0');
    expect(screen.getByTestId('group-bulk-exclude')).toBeDisabled();
  });

  it('flags an incompatible column and a shared column on the affected rows, with text and not only color', async () => {
    await openNewGroup();
    for (const id of ['pt-0', 'pt-1']) fireEvent.click(screen.getByTestId(`group-member-include-${id}`));
    fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'note' } });
    expect(screen.getByTestId('group-member-problem-pt-0')).toHaveTextContent('problem_column_incompatible');
    fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
    fireEvent.change(screen.getByTestId('group-member-column-pt-1'), { target: { value: 'temperature' } });
    expect(screen.getByTestId('group-member-problem-pt-0')).toHaveTextContent('problem_column_conflict');
    expect(screen.getByTestId('group-member-problem-pt-1')).toHaveTextContent('problem_column_conflict');
    expect(screen.getByTestId('group-issues')).toHaveTextContent('column_conflict');
    expect(screen.getByTestId('group-save')).toBeDisabled();
    // Problems-only keeps just those rows visible.
    fireEvent.click(screen.getByTestId('group-member-problems-only'));
    expect(screen.queryByTestId('group-member-row-pt-2')).not.toBeInTheDocument();
  });

  it('allows a shared column when each Tag carries its own entity key', async () => {
    await openNewGroup();
    fireEvent.click(screen.getByText('step4.group.editor.advanced'));
    fireEvent.change(screen.getByTestId('group-entity-column'), { target: { value: 'note' } });
    for (const id of ['pt-0', 'pt-1']) {
      fireEvent.click(screen.getByTestId(`group-member-include-${id}`));
      fireEvent.change(screen.getByTestId(`group-member-column-${id}`), { target: { value: 'temperature' } });
    }
    expect(screen.getByTestId('group-issues')).toHaveTextContent('column_conflict');
    const entityInputs = screen.getAllByLabelText(/step4.group.members.entity_label/);
    fireEvent.change(entityInputs[0], { target: { value: 'line-1' } });
    fireEvent.change(entityInputs[1], { target: { value: 'line-2' } });
    expect(screen.getByTestId('group-issues')).not.toHaveTextContent('column_conflict');
  });

  it('never fills in sample columns when the table cannot be read, and says why', async () => {
    metadata = { status: 'forbidden', columns: [] };
    await openNewGroup();
    expect(screen.getByTestId('group-member-metadata-note')).toHaveTextContent('metadata_forbidden');
    expect(screen.queryByTestId('group-member-suggested-pt-0')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
    // Without metadata the column is typed by hand and nothing is claimed verified.
    expect(screen.getByTestId('group-member-column-pt-0').tagName).toBe('INPUT');
    expect(screen.getByTestId('group-member-column-pt-0')).toHaveValue('');
  });

  it('separates not-checked, checking, missing and failed table states', async () => {
    for (const status of ['not_checked', 'checking', 'missing', 'failed']) {
      metadata = { status, columns: [] };
      document.body.innerHTML = '';
      await openNewGroup();
      expect(screen.getByTestId('group-member-metadata-note')).toHaveTextContent(`metadata_${status}`);
    }
  });

  it('marks a saved assignment for repair when its column disappeared instead of rebinding it', async () => {
    api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [savedGroup()] });
    metadata = { status: 'exists', columns: TABLE_COLUMNS.filter((column) => column.name !== 'pressure') };
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={client}><WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    expect(await screen.findByTestId('group-member-problem-pt-1')).toHaveTextContent('problem_column_missing');
    expect(screen.getByTestId('group-member-column-pt-1')).toHaveValue('pressure');
    await waitFor(() => expect(screen.getByTestId('group-saved-note')).toBeInTheDocument());
  });

  it('lists points that cannot join a group and why', async () => {
    const state = savedGroupState();
    state.mappings['rule-1-p-2'] = { ...state.mappings['rule-1-p-2'], persisted: false, save_state: 'idle' };
    api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={client}><WriteGroupSection state={state} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-create'));
    expect(await screen.findByTestId('group-member-excluded')).toHaveTextContent('excluded_mapping_not_saved');
  });
});
