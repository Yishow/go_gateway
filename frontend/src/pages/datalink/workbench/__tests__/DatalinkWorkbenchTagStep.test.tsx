import { fireEvent, render, screen, waitFor } from '@testing-library/react';
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
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
  }),
  useCreateTagMutation: () => mockCreateTagMutation,
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
  }),
  useCreateMappingMutation: () => mockCreateMappingMutation,
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
  });

  it('previews tag keys for selected points and updates the preview when the prefix changes', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-preview-point-1')).toHaveTextContent('TAG_40001');

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

    expect(screen.getByTestId('tag-candidate-point-1')).toHaveTextContent('40001');
    expect(screen.getByTestId('tag-raw-point-1')).toHaveTextContent('123');
    expect(screen.getByTestId('tag-status-point-1')).toHaveTextContent(
      'workbench.tag.board.status.unbound',
    );
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

    expect(screen.queryByText('Flow Sensor')).not.toBeInTheDocument();
    expect(screen.getByText('Pressure Sensor')).toBeInTheDocument();
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
    fireEvent.click(screen.getByLabelText('Pressure Sensor'));
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

    expect(screen.getByTestId('tag-preview-point-1')).toHaveAttribute('data-conflict', 'true');
    expect(screen.getByRole('button', { name: 'workbench.tag.actions.bind' })).toBeDisabled();
    expect(screen.getByText('workbench.tag.conflicts.summary')).toBeInTheDocument();
  });

  it('shows a partial failure summary when part of the batch bind fails', async () => {
    mockCreateTagMutation.mutateAsync
      .mockResolvedValueOnce({
        id: 'tag-1',
        key: 'TAG_40001',
        display_name: 'TAG_40001',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      })
      .mockRejectedValueOnce(new Error('duplicate key'));

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.actions.bind' }));

    await waitFor(() => {
      expect(mockCreateTagMutation.mutateAsync).toHaveBeenCalledTimes(2);
    });

    expect(mockCreateMappingMutation.mutateAsync).toHaveBeenCalledTimes(1);
    expect(screen.getByText('workbench.tag.results.partialFailure')).toBeInTheDocument();
    expect(screen.getByText('duplicate key')).toBeInTheDocument();
  });
});
