import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchProvider, useWorkbench } from '../WorkbenchProvider';
import { TagBindingStudio } from '../TagBindingStudio';
import type { Device, Point, SourceRuleRecord, Tag } from '../../../../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../../types/sourceRuleCandidates';
import type { SourceRuleTagReviewDecision } from '../../../../types/sourceRuleTagReviewDecisions';

const {
  mockDevices,
  mockPoints,
  mockSourceRules,
  mockTags,
  mockCandidateViews,
  mockReviewDecisions,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockTags: [] as Tag[],
  mockCandidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
  mockReviewDecisions: {} as Record<string, SourceRuleTagReviewDecision[]>,
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
    data: mockTags,
    isLoading: false,
    refetch: vi.fn().mockResolvedValue({ data: mockTags }),
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
    isRefetching: false,
    refetch: vi.fn().mockResolvedValue({ data: mockSourceRules }),
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockCandidateViews[ruleId] ?? null : null,
    isLoading: false,
    isError: false,
    error: null,
    isRefetching: false,
    refetch: vi.fn().mockResolvedValue({ data: null }),
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleTagReviewDecisions', () => ({
  useSourceRuleTagReviewDecisionsQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockReviewDecisions[ruleId] ?? [] : [],
    isLoading: false,
    isError: false,
    error: null,
  }),
  useUpsertSourceRuleTagReviewDecisionMutation: (ruleId?: string | null) => ({
    mutateAsync: async (request: {
      candidate_id: string;
      action: 'rename' | 'skip' | 'override';
      tag_key?: string;
      override_tag_id?: string;
    }) => {
      const nextDecision: SourceRuleTagReviewDecision = {
        source_rule_id: ruleId ?? '',
        candidate_id: request.candidate_id,
        action: request.action,
        tag_key:
          request.action === 'rename'
            ? request.tag_key
            : request.action === 'override'
              ? mockTags.find((tag) => tag.id === request.override_tag_id)?.key
              : undefined,
        override_tag_id:
          request.action === 'override' ? request.override_tag_id : undefined,
        stale: false,
        created_at: '',
        updated_at: '',
      };
      const current = mockReviewDecisions[ruleId ?? ''] ?? [];
      mockReviewDecisions[ruleId ?? ''] = [
        ...current.filter((decision) => decision.candidate_id !== request.candidate_id),
        nextDecision,
      ];
      return nextDecision;
    },
    isPending: false,
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

function buildCandidate(id: string, address: string, tagKey: string) {
  return {
    id,
    identity: {
      source_rule_id: 'rule-1',
      candidate_type: 'tags',
      candidate_kind: 'tag',
      derived_from_rule_address: address,
    },
    proposed_signature: `${id}-sig`,
    address,
    point_id: `point-${address}`,
    tag_key: tagKey,
    display_name: '',
    data_type: 'int16' as const,
    status: 'draft' as const,
  };
}

describe('DatalinkWorkbench tag review actions', () => {
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
        id: 'point-40001',
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
        id: 'point-40002',
        device_id: 'device-1',
        name: 'Level Sensor',
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

    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      skipped_addresses: [],
      revision_id: 'rev-1',
      created_at: '',
      updated_at: '',
    });

    mockTags.splice(0, mockTags.length);

    for (const key of Object.keys(mockCandidateViews)) {
      delete mockCandidateViews[key];
    }
    for (const key of Object.keys(mockReviewDecisions)) {
      delete mockReviewDecisions[key];
    }

    mockCandidateViews['rule-1'] = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [
          buildCandidate('candidate-1', '40001', 'SRC_40001'),
          buildCandidate('candidate-2', '40002', 'SRC_40002'),
        ],
      },
      database_outputs: {
        status: 'ready',
        candidates: [],
      },
      local_modbus_outputs: {
        status: 'ready',
        candidates: [],
      },
    };
  });

  it('saves rename feedback and updates the effective tag key', async () => {
    renderTagStep('rule-1');

    fireEvent.change(screen.getByTestId('source-rule-tag-review-rename-input-candidate-1'), {
      target: { value: 'RENAMED_40001' },
    });
    fireEvent.click(screen.getByTestId('source-rule-tag-review-rename-save-candidate-1'));

    await waitFor(() => {
      expect(screen.getByTestId('source-rule-tag-review-feedback')).toBeInTheDocument();
      expect(screen.getAllByText('RENAMED_40001').length).toBeGreaterThan(0);
      expect(
        screen.getByTestId('source-rule-tag-review-decision-candidate-1'),
      ).toBeInTheDocument();
    });
  });

  it('persists skip and override decisions inline', async () => {
    mockTags.push({
      id: 'tag-override-1',
      key: 'EXISTING_40002',
      display_name: 'Existing tag',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderTagStep('rule-1');

    fireEvent.click(screen.getByTestId('source-rule-tag-review-skip-candidate-1'));

    await waitFor(() => {
      expect(
        screen.getByTestId('source-rule-tag-review-decision-candidate-1'),
      ).toBeInTheDocument();
    });

    fireEvent.change(screen.getByTestId('source-rule-tag-review-override-select-candidate-2'), {
      target: { value: 'tag-override-1' },
    });
    fireEvent.click(screen.getByTestId('source-rule-tag-review-override-save-candidate-2'));

    await waitFor(() => {
      expect(screen.getAllByText('EXISTING_40002').length).toBeGreaterThan(0);
      expect(
        screen.getByTestId('source-rule-tag-review-decision-candidate-2'),
      ).toBeInTheDocument();
    });
  });

  it('disables Step 3 review actions when the open candidate revision is stale', async () => {
    mockSourceRules[0].revision_id = 'rev-2';
    mockCandidateViews['rule-1'] = {
      ...mockCandidateViews['rule-1'],
      revision_id: 'rev-1',
    };
    mockTags.push({
      id: 'tag-override-1',
      key: 'EXISTING_40002',
      display_name: 'Existing tag',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });

    renderTagStep('rule-1');

    await waitFor(() =>
      expect(screen.getByTestId('source-rule-tag-review-stale')).toBeInTheDocument(),
    );

    expect(screen.getByTestId('source-rule-tag-review-rename-input-candidate-1')).toBeDisabled();
    expect(screen.getByTestId('source-rule-tag-review-rename-save-candidate-1')).toBeDisabled();
    expect(screen.getByTestId('source-rule-tag-review-skip-candidate-1')).toBeDisabled();
    expect(screen.getByTestId('source-rule-tag-review-override-select-candidate-2')).toBeDisabled();
    expect(screen.getByTestId('source-rule-tag-review-override-save-candidate-2')).toBeDisabled();
  });
});
