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
  devicesRefetch,
  pointsRefetch,
  tagsRefetch,
  mappingsRefetch,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  devicesRefetch: vi.fn().mockResolvedValue(undefined),
  pointsRefetch: vi.fn().mockResolvedValue(undefined),
  tagsRefetch: vi.fn().mockResolvedValue(undefined),
  mappingsRefetch: vi.fn().mockResolvedValue(undefined),
}));

const shellTranslations: Record<string, string> = {
  'workbench.shell.eyebrow': 'Incident context',
  'workbench.shell.labels.blocker': 'Current blocker',
  'workbench.shell.labels.refresh': 'Diagnostics status',
  'workbench.shell.blockers.outputRuleRequired':
    'Output needs an active rule handoff. Return to Tag review to reopen the source rule context, then continue in Output.',
  'workbench.shell.blockers.outputPending':
    'Output still needs review. Check Local Modbus or Database, then apply the mapping.',
  'workbench.shell.blockers.sourceRequired':
    'Source is not ready. Add or recover source points before advancing.',
  'workbench.shell.refresh.success': 'Diagnostics refreshed',
  'workbench.shell.actions.returnTag': 'Return to tag review',
  'workbench.shell.actions.returnOutput': 'Open output workspace',
  'workbench.shell.actions.returnSource': 'Return to source',
};

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, unknown>) => {
      if (key === 'workbench.shell.readinessSummary') {
        return `${vars?.readyCount}/${vars?.stepCount} steps ready · ${vars?.attentionCount} need attention`;
      }
      return shellTranslations[key] ?? key;
    },
  }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
    refetch: devicesRefetch,
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
  useTestDraftConnectionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useToggleDeviceStatusMutation: () => ({
    mutateAsync: vi.fn().mockResolvedValue(undefined),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
    refetch: pointsRefetch,
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
    refetch: tagsRefetch,
  }),
  useCreateTagMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeleteTagMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
    refetch: mappingsRefetch,
  }),
  useCreateMappingMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeleteMappingMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  modbusShareAPI: {
    status: vi.fn().mockResolvedValue({
      enabled: false,
      port: 5020,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    }),
    listMappings: vi.fn().mockResolvedValue([]),
    start: vi.fn(),
    stop: vi.fn(),
    upsertMapping: vi.fn(),
    deleteMapping: vi.fn(),
    writeTagValue: vi.fn(),
    sync: vi.fn(),
  },
  tagAPI: {
    batchCreate: vi.fn().mockResolvedValue({ created: [], errors: [] }),
  },
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

function seedHappyPath() {
  mockDevices.splice(0, mockDevices.length, {
    id: 'device-1',
    name: 'Mixer PLC',
    description: '',
    protocol: 'modbus_tcp',
    status: 'active',
    connection_config: '{"host":"192.168.1.10","port":502,"slave_id":7,"timeout":5}',
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

  mockTags.splice(0, mockTags.length, {
    id: 'tag-1',
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

  mockMappings.splice(0, mockMappings.length, {
    id: 'mapping-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    enabled: true,
    transform_pipeline: '',
    created_at: '',
    updated_at: '',
  });
}

describe('MuiWorkbench shell incident desk (reopened Phase 5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    devicesRefetch.mockResolvedValue(undefined);
    pointsRefetch.mockResolvedValue(undefined);
    tagsRefetch.mockResolvedValue(undefined);
    mappingsRefetch.mockResolvedValue(undefined);
    seedHappyPath();
  });

  it('surfaces the output blocker and shortest return action in the shell', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('shell-incident-strip')).toBeInTheDocument();
    expect(screen.getByText('Incident context')).toBeInTheDocument();
    expect(screen.getByText('Current blocker')).toBeInTheDocument();
    expect(screen.getByText('Diagnostics status')).toBeInTheDocument();
    expect(screen.getByTestId('shell-incident-blocker')).toHaveTextContent(
      'Output needs an active rule handoff. Return to Tag review to reopen the source rule context, then continue in Output.',
    );
    expect(screen.getByTestId('shell-return-action')).toHaveTextContent(
      'Return to tag review',
    );

    fireEvent.click(screen.getByTestId('shell-return-action'));

    expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'tag');
  });

  it('routes back to source when upstream readiness is missing', () => {
    mockPoints.splice(0, mockPoints.length);
    mockTags.splice(0, mockTags.length);
    mockMappings.splice(0, mockMappings.length);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('shell-incident-blocker')).toHaveTextContent(
      'Source is not ready. Add or recover source points before advancing.',
    );
    expect(screen.getByTestId('shell-return-action')).toHaveTextContent(
      'Return to source',
    );

    fireEvent.click(screen.getByTestId('shell-return-action'));

    expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'source');
  });

  it('refreshes shell diagnostics and reports success', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('shell-refresh-action'));

    await waitFor(() => {
      expect(devicesRefetch).toHaveBeenCalledTimes(1);
      expect(pointsRefetch).toHaveBeenCalledTimes(1);
      expect(tagsRefetch).toHaveBeenCalledTimes(1);
      expect(mappingsRefetch).toHaveBeenCalledTimes(1);
    });

    await waitFor(() => {
      expect(screen.getByTestId('shell-refresh-status')).toHaveTextContent(
        'Diagnostics refreshed',
      );
    });
  });

  it('exposes retry state when shell diagnostics refresh fails', async () => {
    devicesRefetch.mockRejectedValueOnce(new Error('boom'));

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('shell-refresh-action'));

    await waitFor(() => {
      expect(screen.getByTestId('shell-refresh-status')).toHaveTextContent(
        'workbench.shell.refresh.error',
      );
    });

    expect(screen.getByTestId('shell-refresh-action')).toHaveTextContent(
      'workbench.shell.refresh.retry',
    );
  });
});
