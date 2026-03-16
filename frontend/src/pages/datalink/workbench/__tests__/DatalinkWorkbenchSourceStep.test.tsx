import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Point } from '../../../../types/datalink';

const { mockDevices, mockPoints, mockCreatePointMutation } = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockCreatePointMutation: {
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
  useCreatePointMutation: () => mockCreatePointMutation,
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

describe('DatalinkWorkbench source step', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: 'Line A mixer',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });

    mockPoints.splice(0, mockPoints.length, {
      id: 'point-1',
      device_id: 'device-1',
      name: 'Existing Pressure',
      description: '',
      data_type: 'int16',
      address: '40005',
      enabled: true,
      polling_group_id: '',
      last_value: 12,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });
    mockCreatePointMutation.mutateAsync.mockResolvedValue(undefined);
  });

  it('gates source planning behind device selection', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));

    expect(screen.getByText('workbench.source.empty.title')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Mixer PLC' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.source.planner.count')).toBeInTheDocument();
  });

  it('renders planned addresses in grid and ledger views', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.apply' }));

    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40005')).toHaveAttribute('data-status', 'used');

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.table' }));

    expect(screen.getByTestId('address-ledger-row-40001')).toHaveTextContent('40001');
    expect(screen.getByTestId('address-ledger-row-40005')).toHaveTextContent('Existing Pressure');
  });

  it('batch creates points from the planned address range', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.apply' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.batchCreate' }));

    await waitFor(() => {
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(2);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenNthCalledWith(1, {
      device_id: 'device-1',
      address: '40001',
      data_type: 'float32',
      name: 'SRC_40001',
    });
    expect(mockCreatePointMutation.mutateAsync).toHaveBeenNthCalledWith(2, {
      device_id: 'device-1',
      address: '40003',
      data_type: 'float32',
      name: 'SRC_40003',
    });
  });

  it('blocks batch create when planned cells overlap existing points', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.apply' }));

    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'conflict');
    expect(
      screen.getByRole('button', { name: 'workbench.source.planner.batchCreate' }),
    ).toBeDisabled();
  });
});
