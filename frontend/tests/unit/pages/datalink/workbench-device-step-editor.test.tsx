import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchDeviceStep } from '../../../../src/pages/datalink/workbench/WorkbenchDeviceStep';
import {
  WorkbenchProvider,
  useWorkbench,
} from '../../../../src/pages/datalink/workbench/WorkbenchProvider';
import type { Device } from '../../../../src/types/datalink';

const {
  mockDeviceQuery,
  mockCreateDeviceMutation,
  mockTestDraftConnectionMutation,
  mockUpdateDeviceMutation,
} = vi.hoisted(() => ({
  mockDeviceQuery: {
    data: [] as Device[],
    isLoading: false,
    isSuccess: true,
    isFetching: false,
  },
  mockCreateDeviceMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockTestDraftConnectionMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockUpdateDeviceMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../src/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => mockDeviceQuery,
  useCreateDeviceMutation: () => mockCreateDeviceMutation,
  useTestDraftConnectionMutation: () => mockTestDraftConnectionMutation,
  useUpdateDeviceMutation: () => mockUpdateDeviceMutation,
}));

function OpenCreatePanelButton() {
  const { openCreateDevicePanel } = useWorkbench();

  return (
    <button onClick={openCreateDevicePanel} type="button">
      open-create
    </button>
  );
}

function OpenEditPanelButton() {
  const { openEditDevicePanel } = useWorkbench();

  return (
    <button onClick={() => openEditDevicePanel('device-1')} type="button">
      open-edit
    </button>
  );
}

function buildTree(queryClient: QueryClient): ReactNode {
  return (
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <OpenCreatePanelButton />
        <OpenEditPanelButton />
        <WorkbenchDeviceStep />
      </WorkbenchProvider>
    </QueryClientProvider>
  );
}

function renderDeviceStep() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return {
    queryClient,
    ...render(buildTree(queryClient)),
  };
}

describe('WorkbenchDeviceStep create editor', () => {
  beforeEach(() => {
    mockDeviceQuery.data = [];
    mockDeviceQuery.isLoading = false;
    mockDeviceQuery.isSuccess = true;
    mockDeviceQuery.isFetching = false;
    mockCreateDeviceMutation.mutateAsync.mockReset();
    mockTestDraftConnectionMutation.mutateAsync.mockReset();
    mockUpdateDeviceMutation.mutateAsync.mockReset();
  });

  it('preserves draft input and protocol selection after a device list refresh', () => {
    const view = renderDeviceStep();

    fireEvent.click(screen.getByRole('button', { name: 'open-create' }));
    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Editable PLC' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.fields.protocol'), {
      target: { value: 'modbus_rtu' },
    });

    expect(screen.getByLabelText('workbench.device.fields.name')).toHaveValue(
      'Editable PLC',
    );
    expect(screen.getByLabelText('workbench.device.fields.protocol')).toHaveValue(
      'modbus_rtu',
    );

    mockDeviceQuery.data = [];
    view.rerender(buildTree(view.queryClient));

    expect(screen.getByLabelText('workbench.device.fields.name')).toHaveValue(
      'Editable PLC',
    );
    expect(screen.getByLabelText('workbench.device.fields.protocol')).toHaveValue(
      'modbus_rtu',
    );
  });

  it('submits modbus tcp data_format from the create editor', async () => {
    mockCreateDeviceMutation.mutateAsync.mockResolvedValue({
      id: 'device-created',
    });

    renderDeviceStep();

    fireEvent.click(screen.getByRole('button', { name: 'open-create' }));
    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Modbus TCP Device' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.host'), {
      target: { value: '192.168.0.10' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.dataFormat'), {
      target: { value: 'BADC' },
    });

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(mockCreateDeviceMutation.mutateAsync).toHaveBeenCalledWith(
        expect.objectContaining({
          protocol: 'modbus_tcp',
          connection_config: expect.objectContaining({
            data_format: 'BADC',
          }),
        }),
      );
    });
  });

  it('loads and persists mc_3e data_format in the edit editor', async () => {
    mockDeviceQuery.data = [
      {
        id: 'device-1',
        name: 'MC Device',
        description: '',
        protocol: 'mc_3e',
        status: 'active',
        connection_config: JSON.stringify({
          host: '10.0.0.20',
          port: 5000,
          data_format: 'binary',
        }),
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    ];
    mockUpdateDeviceMutation.mutateAsync.mockResolvedValue(undefined);

    renderDeviceStep();

    fireEvent.click(screen.getByRole('button', { name: 'open-edit' }));

    const dataFormatField = screen.getByLabelText('workbench.device.connection.dataFormat');
    expect(dataFormatField).toHaveValue('CDAB');

    fireEvent.change(dataFormatField, { target: { value: 'DCBA' } });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(mockUpdateDeviceMutation.mutateAsync).toHaveBeenCalledWith({
        id: 'device-1',
        data: expect.objectContaining({
          connection_config: expect.objectContaining({
            data_format: 'DCBA',
          }),
        }),
      });
    });
  });
});
