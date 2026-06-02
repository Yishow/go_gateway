import { fireEvent, render, screen, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Mapping, Point, SourceRuleRecord, Tag } from '../../../../types/datalink';
import type { SourceRuleCandidateSnapshotView, SourceRuleDatabaseOutputCandidateView } from '../../../../types/sourceRuleCandidates';

const {
  mockDevices,
  mockPoints,
  mockMappings,
  mockTags,
  mockSourceRules,
  mockCandidateViews,
  mockRuntimeStatus,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCandidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
  mockRuntimeStatus: vi.fn(),
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestDraftConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useToggleDeviceStatusMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeletePointMutation: () => ({ mutate: vi.fn(), mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockSourceRules.filter((rule) => rule.device_id === filters.device_id)
      : [],
    isLoading: false,
    isSuccess: true,
  }),
  useCreateSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useEnableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDisableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), mutate: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false, refetch: vi.fn() }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useRuntimeStream', () => ({
  useRuntimeStream: () => ({ liveValues: {} }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockCandidateViews[ruleId] ?? null : null,
    isLoading: false,
    isError: false,
    error: null,
    isRefetching: false,
    refetch: vi.fn(),
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  runtimeAPI: { getStatus: mockRuntimeStatus },
}));

function makeDatabaseCandidate(
  ruleId: string,
  pointId: string,
  address: string,
  tagKey: string,
  groupKey: string,
  columnName: string,
): SourceRuleDatabaseOutputCandidateView {
  return {
    id: `db-${pointId}`,
    identity: { source_rule_id: ruleId, candidate_type: 'database_outputs', candidate_kind: 'database_output', derived_from_rule_address: address },
    proposed_signature: `db-signature-${pointId}`,
    address,
    point_id: pointId,
    tag_key: tagKey,
    display_name: `${pointId}-database`,
    data_type: 'float64',
    status: 'ready',
    connector_id: 'connector-1',
    table_schema: 'public',
    table_name: 'line_metrics',
    column_name: columnName,
    group_key: groupKey,
    write_mode: 'insert',
    write_interval_seconds: 5,
  };
}

function makeCandidateSnapshot(
  ruleId: string,
  revisionId: string,
  pointId: string,
  address: string,
  tagKey: string,
  groupKey: string,
  columnName: string,
): SourceRuleCandidateSnapshotView {
  return {
    source_rule_id: ruleId,
    revision_id: revisionId,
    tags: {
      status: 'ready',
      candidates: [
        {
          id: `candidate-${pointId}`,
          identity: { source_rule_id: ruleId, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: address },
          proposed_signature: `signature-${revisionId}-${pointId}`,
          address,
          point_id: pointId,
          tag_key: tagKey,
          display_name: `${pointId}-tag`,
          data_type: 'int16',
          status: 'draft',
        },
      ],
    },
    database_outputs: { status: 'ready', candidates: [makeDatabaseCandidate(ruleId, pointId, address, tagKey, groupKey, columnName)] },
    local_modbus_outputs: { status: 'deferred', candidates: [] },
  };
}

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchPage />
    </QueryClientProvider>,
  );
}

async function openSourceStep() {
  renderPage();
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
  fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
  fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));
  await screen.findByTestId('source-active-rule-summary');
}

async function switchToDatabaseTarget() {
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
  fireEvent.click(screen.getByRole('tab', { name: 'workbench.output.targetSwitcher.database' }));
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
  fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));
  await screen.findByTestId('source-active-rule-summary');
}

describe('DatalinkWorkbench source preview continuity', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockDevices.splice(0, mockDevices.length, { id: 'device-1', name: 'Mixer PLC', description: 'Line A mixer', protocol: 'modbus_tcp', status: 'active', connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' });
    mockPoints.splice(0, mockPoints.length,
      { id: 'point-1', device_id: 'device-1', name: 'Flow Sensor', description: '', data_type: 'int16', address: '40001', enabled: true, polling_group_id: '', last_value: null, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
      { id: 'point-2', device_id: 'device-1', name: 'Pressure Sensor', description: '', data_type: 'int16', address: '40002', enabled: true, polling_group_id: '', last_value: null, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
    );
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length,
      { id: 'rule-1', device_id: 'device-1', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'SRC_A', enabled: true, locked: false, origin: 'manual', skipped_addresses: [], revision_id: 'rev-1', created_at: '', updated_at: '' },
      { id: 'rule-2', device_id: 'device-1', start_address: '40002', count: 1, data_type: 'int16', naming_prefix: 'SRC_B', enabled: true, locked: false, origin: 'manual', skipped_addresses: [], revision_id: 'rev-2', created_at: '', updated_at: '' },
    );
    mockCandidateViews['rule-1'] = makeCandidateSnapshot('rule-1', 'rev-1', 'point-1', '40001', 'TAG_40001', 'line-a', 'flow_value');
    mockCandidateViews['rule-2'] = makeCandidateSnapshot('rule-2', 'rev-2', 'point-2', '40002', 'TAG_40002', 'line-b', 'pressure_value');
    mockRuntimeStatus.mockResolvedValue({ collectors: [] });
  });

  it('temporarily previews hovered rules and restores the selected rule on mouse leave', async () => {
    await openSourceStep();

    const preview = screen.getByTestId('source-tag-preview');
    expect(within(preview).getByText('TAG_40001')).toBeInTheDocument();

    const rule2Chip = screen.getByTestId('source-preview-rule-rule-2');
    fireEvent.mouseEnter(rule2Chip);
    expect(within(preview).getByText('workbench.source.preview.hoverTitle')).toBeInTheDocument();
    expect(within(preview).getByText('TAG_40002')).toBeInTheDocument();

    fireEvent.mouseLeave(rule2Chip);
    expect(within(preview).getByText('workbench.source.preview.title')).toBeInTheDocument();
    expect(within(preview).getByText('TAG_40001')).toBeInTheDocument();

    fireEvent.click(rule2Chip);
    expect(within(preview).getByText('TAG_40002')).toBeInTheDocument();

    const rule1Chip = screen.getByTestId('source-preview-rule-rule-1');
    fireEvent.mouseEnter(rule1Chip);
    expect(within(preview).getByText('TAG_40001')).toBeInTheDocument();

    fireEvent.mouseLeave(rule1Chip);
    expect(within(preview).getByText('TAG_40002')).toBeInTheDocument();
  });

  it('keeps source preview, tag scope, and database grouping hints aligned for the selected rule', async () => {
    await openSourceStep();
    await switchToDatabaseTarget();

    fireEvent.click(screen.getByTestId('source-preview-rule-rule-2'));

    const preview = screen.getByTestId('source-tag-preview');
    expect(within(preview).getByText('TAG_40002')).toBeInTheDocument();

    expect(within(preview).getByText('line-b')).toBeInTheDocument();
    expect(within(preview).getByText('pressure_value')).toBeInTheDocument();
    expect(screen.getByTestId('source-database-planning-hints')).toHaveTextContent('MBT_ROW');
    expect(screen.getByTestId('source-diagnostic-toolbar')).toHaveTextContent(
      'workbench.source.diagnostics.description',
    );

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    expect(await screen.findByLabelText('Pressure Sensor')).toBeInTheDocument();
    expect(screen.queryByLabelText('Flow Sensor')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));
    expect(await screen.findByTestId('source-tag-preview')).toHaveTextContent('TAG_40002');
    expect(screen.getByTestId('source-active-rule-address')).toHaveTextContent('40002');
  });
});
