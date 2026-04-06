import { fireEvent, render, screen } from '@testing-library/react';
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

function buildTree(queryClient: QueryClient): ReactNode {
  return (
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <OpenCreatePanelButton />
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
});
