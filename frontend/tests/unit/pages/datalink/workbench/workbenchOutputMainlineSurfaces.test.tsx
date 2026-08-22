import { useEffect } from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { Device, Mapping, Point, Tag } from '@/types/datalink';
import type { SourceRuleCandidateSnapshotView } from '@/types/sourceRuleCandidates';
import { MuiWorkbenchBottomSummaryBar } from '@/pages/datalink/workbench/MuiWorkbenchBottomSummaryBar';
import { MuiWorkbenchIncidentStrip } from '@/pages/datalink/workbench/MuiWorkbenchIncidentStrip';
import { MuiWorkbenchStepRail } from '@/pages/datalink/workbench/MuiWorkbenchStepRail';
import { useWorkbench, WorkbenchProvider } from '@/pages/datalink/workbench/WorkbenchProvider';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockSnapshotState,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockSnapshotState: { value: null as SourceRuleCandidateSnapshotView | null },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
  }),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
}));

vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
  }),
}));

vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
  }),
}));

vi.mock('@/hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: () => ({
    data: mockSnapshotState.value,
    isLoading: false,
  }),
}));

vi.mock('@/pages/datalink/workbench/useWorkbenchShellDiagnostics', () => ({
  useWorkbenchShellDiagnostics: () => ({
    isRefreshing: false,
    statusKey: 'workbench.shell.refresh.idle',
    actionKey: 'workbench.shell.refresh.action',
    refreshDiagnostics: vi.fn().mockResolvedValue(undefined),
  }),
}));

function SeedOutputReadyState() {
  const {
    setActiveStep,
    setFocusedRuleId,
    setSelectedDeviceId,
    setSourcePlanningState,
  } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    setActiveStep('output');
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedRuleId: 'rule-1',
    }));
    setFocusedRuleId('rule-1');
  }, [setActiveStep, setFocusedRuleId, setSelectedDeviceId, setSourcePlanningState]);

  return null;
}

describe('workbench Output mainline surfaces', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'UI 4.3 Modbus TCP',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: true,
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
        name: 'MBT_40001',
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
        name: 'MBT_40002',
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

    mockTags.splice(
      0,
      mockTags.length,
      {
        id: 'tag-1',
        key: 'meter/A1',
        display_name: 'Meter A1',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'active',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'tag-2',
        key: 'meter/kw',
        display_name: 'Meter KW',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'active',
        created_at: '',
        updated_at: '',
      },
    );

    mockMappings.splice(
      0,
      mockMappings.length,
      {
        id: 'mapping-1',
        point_id: 'point-1',
        tag_id: 'tag-1',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'mapping-2',
        point_id: 'point-2',
        tag_id: 'tag-2',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
    );

    mockSnapshotState.value = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [],
      },
      database_outputs: {
        status: 'ready',
        candidates: [
          {
            id: 'db-candidate-1',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'database_outputs',
              candidate_kind: 'database_output',
              derived_from_rule_address: '40001',
            },
            proposed_signature: 'db-1',
            address: '40001',
            point_id: 'point-1',
            tag_id: 'tag-1',
            mapping_id: 'db-map-1',
            tag_key: 'meter/A1',
            display_name: 'Meter A1',
            data_type: 'int16',
            status: 'ready',
            connector_id: 'connector-1',
            table_schema: 'main',
            table_name: 'meter_rows',
            column_name: 'a1',
          },
        ],
      },
      local_modbus_outputs: {
        status: 'deferred',
        reason: 'not-configured',
        candidates: [],
      },
    };
  });

  it('marks Output as mainline-ready across the rail, shell, and bottom summary when Database is applied', async () => {
    render(
      <WorkbenchProvider>
        <SeedOutputReadyState />
        <MuiWorkbenchStepRail />
        <MuiWorkbenchIncidentStrip />
        <MuiWorkbenchBottomSummaryBar />
      </WorkbenchProvider>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('shell-incident-blocker')).toHaveTextContent(
        'workbench.shell.blockers.stable',
      );
    });

    expect(
      screen.getByRole('tab', { name: /workbench\.steps\.output - workbench\.readiness\.ready/ }),
    ).toBeInTheDocument();
    expect(screen.getByTestId('shell-return-action')).toHaveTextContent(
      'workbench.shell.actions.focusCurrent',
    );
    expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-readiness', 'ready');
  });
});
