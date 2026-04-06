import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchProvider, useWorkbench } from '../WorkbenchProvider';
import { TagBindingStudio } from '../TagBindingStudio';
import type { Device, Point, SourceRuleRecord } from '../../../../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../../types/sourceRuleCandidates';

const {
  mockDevices,
  mockPoints,
  mockSourceRules,
  mockCandidateViews,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCandidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
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
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({
    data: mockPoints,
    isLoading: false,
  }),
  useDeletePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: [],
    isLoading: false,
    refetch: vi.fn().mockResolvedValue({ data: [] }),
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
    data: [],
    isLoading: false,
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

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mockSourceRules,
    isLoading: false,
    isSuccess: true,
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockCandidateViews[ruleId] ?? null : null,
    isLoading: false,
    isError: false,
    error: null,
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  tagAPI: {
    batchCreate: vi.fn(),
  },
}));

function TagStepBootstrap({ focusedRuleId }: { focusedRuleId?: string }) {
  const { setSelectedDeviceId, setFocusedRuleId } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    if (focusedRuleId) {
      setFocusedRuleId(focusedRuleId);
    }
  }, [focusedRuleId, setFocusedRuleId, setSelectedDeviceId]);

  return <TagBindingStudio />;
}

function renderTagStep(focusedRuleId?: string) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <TagStepBootstrap focusedRuleId={focusedRuleId} />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbench tag review surface', () => {
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

    mockPoints.splice(0, mockPoints.length, {
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
    });

    mockSourceRules.splice(0, mockSourceRules.length);
    for (const key of Object.keys(mockCandidateViews)) {
      delete mockCandidateViews[key];
    }
  });

  it('loads generated tag candidates from the active source-rule revision', async () => {
    mockSourceRules.push({
      id: 'rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 1,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      skipped_addresses: [],
      created_at: '',
      updated_at: '',
    });
    mockCandidateViews['rule-1'] = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [
          {
            id: 'candidate-1',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'tags',
              candidate_kind: 'tag',
              derived_from_rule_address: '40001',
            },
            proposed_signature: 'sig-1',
            address: '40001',
            point_id: 'point-1',
            tag_key: 'SRC_40001',
            display_name: 'Flow Sensor',
            data_type: 'int16',
            status: 'draft',
          },
        ],
      },
      database_outputs: {
        status: 'deferred',
        candidates: [],
      },
      local_modbus_outputs: {
        status: 'deferred',
        candidates: [],
      },
    };

    renderTagStep('rule-1');

    await waitFor(() =>
      expect(screen.getByTestId('source-rule-tag-review-surface')).toBeInTheDocument(),
    );

    expect(screen.getByTestId('source-rule-tag-review-revision')).toHaveTextContent('rev-1');
    expect(screen.getByTestId('source-rule-tag-review-row-candidate-1')).toHaveTextContent(
      'SRC_40001',
    );
    expect(screen.getByTestId('source-rule-tag-review-status-candidate-1')).toHaveTextContent(
      'workbench.tag.reviewSurface.status.generated',
    );
    expect(
      screen.getByTestId('source-rule-tag-review-mapping-intent-candidate-1'),
    ).toHaveTextContent('workbench.tag.reviewSurface.mappingIntent.pending');
  });

  it('switches Step 3 to the newly selected source-rule revision without mixing old candidates', async () => {
    mockSourceRules.push(
      {
        id: 'rule-1',
        device_id: 'device-1',
        start_address: '40001',
        count: 1,
        data_type: 'int16',
        naming_prefix: 'SRC',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        created_at: '',
        updated_at: '',
      },
      {
        id: 'rule-2',
        device_id: 'device-1',
        start_address: '40101',
        count: 1,
        data_type: 'float32',
        naming_prefix: 'LINE',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        created_at: '',
        updated_at: '',
      },
    );
    mockCandidateViews['rule-1'] = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [
          {
            id: 'candidate-1',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'tags',
              candidate_kind: 'tag',
              derived_from_rule_address: '40001',
            },
            proposed_signature: 'sig-1',
            address: '40001',
            point_id: 'point-1',
            tag_key: 'SRC_40001',
            display_name: 'Flow Sensor',
            data_type: 'int16',
            status: 'draft',
          },
        ],
      },
      database_outputs: {
        status: 'deferred',
        candidates: [],
      },
      local_modbus_outputs: {
        status: 'deferred',
        candidates: [],
      },
    };
    mockCandidateViews['rule-2'] = {
      source_rule_id: 'rule-2',
      revision_id: 'rev-2',
      tags: {
        status: 'ready',
        candidates: [
          {
            id: 'candidate-2',
            identity: {
              source_rule_id: 'rule-2',
              candidate_type: 'tags',
              candidate_kind: 'tag',
              derived_from_rule_address: '40101',
            },
            proposed_signature: 'sig-2',
            address: '40101',
            point_id: 'point-1',
            tag_key: 'LINE_40101',
            display_name: 'Pressure Sensor',
            data_type: 'float32',
            status: 'draft',
          },
        ],
      },
      database_outputs: {
        status: 'deferred',
        candidates: [],
      },
      local_modbus_outputs: {
        status: 'deferred',
        candidates: [],
      },
    };

    renderTagStep('rule-1');

    await waitFor(() =>
      expect(screen.getByTestId('source-rule-tag-review-row-candidate-1')).toBeInTheDocument(),
    );

    fireEvent.change(screen.getByTestId('source-rule-tag-review-rule-select'), {
      target: { value: 'rule-2' },
    });

    await waitFor(() =>
      expect(screen.getByTestId('source-rule-tag-review-row-candidate-2')).toBeInTheDocument(),
    );

    expect(screen.queryByTestId('source-rule-tag-review-row-candidate-1')).not.toBeInTheDocument();
    expect(screen.getByTestId('source-rule-tag-review-revision')).toHaveTextContent('rev-2');
  });
});
