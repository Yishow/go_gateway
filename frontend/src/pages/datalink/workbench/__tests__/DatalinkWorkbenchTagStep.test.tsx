import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '../../../../types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockCreateTagMutation,
  mockCreateMappingMutation,
  mockDeleteMappingMutation,
  mockBatchCreate,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockCreateTagMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockCreateMappingMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockDeleteMappingMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockBatchCreate: vi.fn(),
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
  }),
  useCreateDeviceMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useUpdateDeviceMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useTestConnectionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeletePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
    refetch: vi.fn().mockResolvedValue({ data: mockTags }),
  }),
  useCreateTagMutation: () => mockCreateTagMutation,
}));

vi.mock('../../../../services/datalink', () => ({
  tagAPI: {
    batchCreate: mockBatchCreate,
  },
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
  }),
  useCreateMappingMutation: () => mockCreateMappingMutation,
  useDeleteMappingMutation: () => mockDeleteMappingMutation,
}));

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

describe('DatalinkWorkbench tag step', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });

    mockPoints.splice(
      0,
      mockPoints.length,
      {
        id: 'point-1',
        device_id: 'device-1',
        name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        address: '40001',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
      {
        id: 'point-2',
        device_id: 'device-1',
        name: 'Pressure Sensor',
        description: '',
        data_type: 'int16',
        address: '40002',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    );

    mockTags.splice(0, mockTags.length);
    mockMappings.splice(0, mockMappings.length);
    mockBatchCreate.mockImplementation(async (tags: Array<{ key: string }>) => {
      // 模擬 batch create：所有 tags 都成功建立
      const created = tags.map((tag) => `tag-${tag.key}`);
      // 同時把 tags 加到 mockTags 以供 refetch 找到
      for (const tag of tags) {
        mockTags.push({
          id: `tag-${tag.key}`,
          key: tag.key,
          display_name: tag.key,
          description: '',
          data_type: 'int16',
          unit: '',
          labels: null,
          status: 'draft',
          created_at: '',
          updated_at: '',
        });
      }
      return { created, errors: [] };
    });
    mockCreateTagMutation.mutateAsync.mockImplementation(async (request: { key: string }) => ({
      id: `tag-${request.key}`,
      key: request.key,
      display_name: request.key,
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    }));
    mockCreateMappingMutation.mutateAsync.mockResolvedValue({
      id: 'mapping-created',
      point_id: 'point-1',
      tag_id: 'tag-created',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
  });

  it('shows an empty state when the selected device has no points', () => {
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByText('workbench.tag.empty.title')).toBeInTheDocument();
    expect(screen.getByTestId('tag-empty-eligible-spans')).toBeInTheDocument();
    expect(screen.getByText('workbench.tag.empty.eligibleSpansNone')).toBeInTheDocument();
  });

  it('shows a go-to-source action in the empty state that navigates to step 2', () => {
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const goToSourceButton = screen.getByTestId('tag-empty-goto-source');
    expect(goToSourceButton).toBeInTheDocument();

    fireEvent.click(goToSourceButton);

    // After clicking, active step should switch to source (step 2)
    const sourceStepButton = screen.getByRole('button', { name: 'workbench.steps.source' });
    expect(sourceStepButton).toHaveAttribute('aria-current', 'step');
  });

  it('previews tag keys for selected points and updates the preview when the prefix changes', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-preview-point-1')).toHaveTextContent('TAG_40001');

    fireEvent.click(screen.getByTestId('tag-template-toggle'));
    fireEvent.change(screen.getByLabelText('workbench.tag.template.prefix'), {
      target: { value: 'linea' },
    });

    expect(screen.getByTestId('tag-preview-point-1')).toHaveTextContent('LINEA_40001');
  });

  it('shows dense source metadata for each candidate row', () => {
    mockPoints[0].last_value = 123;

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-candidate-board')).toHaveAttribute('data-layout', 'row-board');
    expect(screen.getByTestId('tag-candidate-point-1')).toHaveAttribute('data-layout', 'row');
    expect(screen.getByTestId('tag-candidate-point-1')).toHaveTextContent('40001');
    expect(screen.getByTestId('tag-raw-point-1')).toHaveTextContent('123');
    expect(screen.getByTestId('tag-status-point-1')).toHaveTextContent(
      'workbench.tag.board.status.unbound',
    );
  });

  it('shows adaptive board guidance and flow hints as selection state changes', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByText('workbench.tag.board.selectionHint.none')).toBeInTheDocument();
    expect(screen.getByText('workbench.tag.board.flowModeHint.create')).toBeInTheDocument();

    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByText('workbench.tag.board.selectionHint.ready')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));

    expect(screen.getByText('workbench.tag.board.flowModeHint.existing')).toBeInTheDocument();
  });

  it('replaces preview keys with existing-tag selectors in existing flow rows', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'LINEA_FLOW',
      display_name: 'Line A Flow',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));

    const row = screen.getByTestId('tag-candidate-point-1');

    expect(within(row).queryByTestId('tag-preview-point-1')).not.toBeInTheDocument();
    expect(within(row).getByTestId('existing-tag-select-point-1')).toBeInTheDocument();
  });

  it('shows batch actions only after one or more candidate rows are selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.queryByTestId('batch-diff-preview')).not.toBeInTheDocument();
    expect(
      screen.queryByRole('button', { name: 'workbench.tag.actions.bind' }),
    ).not.toBeInTheDocument();

    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByTestId('batch-diff-preview')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'workbench.tag.actions.bind' })).toBeInTheDocument();
  });

  it('keeps conflict detail out of the row and surfaces it in the inspector', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'TAG_40001',
      display_name: 'Existing Flow Sensor',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    } as Tag);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const row = screen.getByTestId('tag-candidate-point-1');
    expect(within(row).queryByText('workbench.tag.badges.existingKey')).not.toBeInTheDocument();

    fireEvent.click(row);

    expect(screen.getByTestId('tag-inspector-panel')).toHaveTextContent(
      'workbench.tag.inspector.conflict.existing-key',
    );
  });

  it('shows tag inspector details when a candidate row is selected', () => {
    mockPoints[0].last_value = 123;

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.click(screen.getByTestId('tag-candidate-point-1'));

    expect(screen.getByTestId('tag-inspector-panel')).toBeInTheDocument();
    expect(screen.getByTestId('tag-inspector-tag-key')).toHaveTextContent('TAG_40001');
    expect(screen.getByTestId('tag-inspector-point-address')).toHaveTextContent('40001');
  });

  it('filters tag candidates by keyword and status', () => {
    mockMappings.push({
      id: 'mapping-2',
      point_id: 'point-2',
      tag_id: 'tag-2',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.tag.board.search'), {
      target: { value: 'Pressure' },
    });
    fireEvent.change(screen.getByLabelText('workbench.tag.board.statusFilter'), {
      target: { value: 'bound' },
    });

    expect(screen.queryByTestId('tag-candidate-point-1')).not.toBeInTheDocument();
    expect(screen.getByTestId('tag-candidate-point-2')).toBeInTheDocument();
  });

  it('supports binding selected points to an existing tag', async () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'LINEA_FLOW',
      display_name: 'Line A Flow',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.change(screen.getByTestId('existing-tag-select-point-1'), {
      target: { value: 'tag-existing' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockCreateMappingMutation.mutateAsync).toHaveBeenCalledWith({
        point_id: 'point-1',
        tag_id: 'tag-existing',
        enabled: true,
      });
    });

    expect(mockCreateTagMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('blocks batch binding when preview keys conflict with existing tags', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'TAG_40001',
      display_name: 'Existing Flow Sensor',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByTestId('tag-preview-point-1')).toHaveAttribute('data-conflict', 'true');
    expect(screen.getByRole('button', { name: 'workbench.tag.actions.bind' })).toBeDisabled();
    expect(screen.getByText('workbench.tag.conflicts.summary')).toBeInTheDocument();
  });

  it('shows a partial failure summary when part of the batch bind fails', async () => {
    // 模擬 batch create 部分失敗（default strategy=address, prefix=TAG）
    mockBatchCreate.mockImplementationOnce(async () => {
      // TAG_40001 成功，TAG_40002 失敗
      mockTags.push({
        id: 'tag-TAG_40001',
        key: 'TAG_40001',
        display_name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      });
      return {
        created: ['tag-TAG_40001'],
        errors: [{ key: 'TAG_40002', error: 'duplicate key' }],
      };
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockBatchCreate).toHaveBeenCalledTimes(1);
    });

    await waitFor(() => {
      expect(screen.getByText('workbench.tag.results.partialFailure')).toBeInTheDocument();
    });
    expect(screen.getByText('duplicate key')).toBeInTheDocument();
  });

  it('shows a batch diff preview before executing the bind', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));

    const diffPanel = screen.getByTestId('batch-diff-preview');
    expect(diffPanel).toBeInTheDocument();
    expect(diffPanel).toHaveTextContent('TAG_40001');
    expect(diffPanel).toHaveTextContent('TAG_40002');
    expect(screen.getByTestId('diff-to-create-count')).toHaveTextContent('2');
    expect(screen.getByTestId('diff-skipped-count')).toHaveTextContent('0');
  });

  it('shows skipped items in the diff preview when candidates have conflicts', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'TAG_40001',
      display_name: 'Existing Flow Sensor',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'draft',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));

    expect(screen.getByTestId('diff-to-create-count')).toHaveTextContent('1');
    expect(screen.getByTestId('diff-skipped-count')).toHaveTextContent('1');
  });

  it('switches diff preview to show toBind items in existing flow', () => {
    mockTags.push({
      id: 'tag-existing',
      key: 'LINEA_FLOW',
      display_name: 'Line A Flow',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByTestId('diff-to-bind-count')).toBeInTheDocument();
  });

  it('shows an enhanced result summary with created/linked/skipped/failed counts', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByLabelText('Flow Sensor'));
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockBatchCreate).toHaveBeenCalledTimes(1);
    });

    expect(screen.getByTestId('result-created-count')).toHaveTextContent('2');
    expect(screen.getByTestId('result-failed-count')).toHaveTextContent('0');
  });

  it('shows an unbind button for bound candidates and calls delete mutation after confirm', async () => {
    mockTags.push({
      id: 'tag-bound-1',
      key: 'TAG_40001',
      display_name: 'TAG 40001',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mockMappings.push({
      id: 'mapping-bound-1',
      point_id: 'point-1',
      tag_id: 'tag-bound-1',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
    mockDeleteMappingMutation.mutateAsync.mockResolvedValue(undefined);
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const candidate = screen.getByTestId('tag-candidate-point-1');
    const unbindButton = within(candidate).getByTestId('tag-unbind-point-1');
    expect(unbindButton).toBeInTheDocument();

    const boundLabel = within(candidate).getByTestId('tag-bound-label-point-1');
    expect(boundLabel).toBeInTheDocument();

    fireEvent.click(unbindButton);

    await waitFor(() => {
      expect(window.confirm).toHaveBeenCalled();
      expect(mockDeleteMappingMutation.mutateAsync).toHaveBeenCalledWith('mapping-bound-1');
    });

    vi.restoreAllMocks();
  });

  it('does not unbind when the user cancels the confirmation dialog', async () => {
    mockTags.push({
      id: 'tag-bound-2',
      key: 'TAG_40001',
      display_name: 'TAG 40001',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mockMappings.push({
      id: 'mapping-bound-2',
      point_id: 'point-1',
      tag_id: 'tag-bound-2',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
    vi.spyOn(window, 'confirm').mockReturnValue(false);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const candidate = screen.getByTestId('tag-candidate-point-1');
    const unbindButton = within(candidate).getByTestId('tag-unbind-point-1');
    fireEvent.click(unbindButton);

    expect(window.confirm).toHaveBeenCalled();
    expect(mockDeleteMappingMutation.mutateAsync).not.toHaveBeenCalled();

    vi.restoreAllMocks();
  });

  it('shows the template panel collapsed by default with a toggle button', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const toggle = screen.getByTestId('tag-template-toggle');
    expect(toggle).toBeInTheDocument();

    expect(screen.queryByLabelText('workbench.tag.template.prefix')).not.toBeInTheDocument();

    fireEvent.click(toggle);
    expect(screen.getByLabelText('workbench.tag.template.prefix')).toBeInTheDocument();
  });
});
