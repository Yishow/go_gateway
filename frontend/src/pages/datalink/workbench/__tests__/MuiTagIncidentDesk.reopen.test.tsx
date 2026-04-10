import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { MuiTagIncidentDesk } from '../MuiTagIncidentDesk';
import { WorkbenchProvider, useWorkbench } from '../WorkbenchProvider';
import type { Device, Mapping, Point, SourceRuleRecord, Tag } from '../../../../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../../types/sourceRuleCandidates';
import type { SourceRuleTagReviewDecision } from '../../../../types/sourceRuleTagReviewDecisions';

const mocks = vi.hoisted(() => ({
  devices: [] as Device[],
  points: [] as Point[],
  tags: [] as Tag[],
  mappings: [] as Mapping[],
  sourceRules: [] as SourceRuleRecord[],
  candidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
  decisions: {} as Record<string, SourceRuleTagReviewDecision[]>,
  candidateState: {
    isLoading: false,
    isError: false,
    error: null as Error | null,
    refetch: vi.fn(),
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mocks.devices, isLoading: false }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({ data: mocks.points, isLoading: false }),
  useDeletePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mocks.tags,
    isLoading: false,
    refetch: vi.fn().mockResolvedValue({ data: mocks.tags }),
  }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mocks.mappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mocks.sourceRules,
    isLoading: false,
    isSuccess: true,
    isRefetching: false,
    refetch: vi.fn().mockResolvedValue({ data: mocks.sourceRules }),
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mocks.candidateViews[ruleId] ?? null : null,
    isLoading: mocks.candidateState.isLoading,
    isError: mocks.candidateState.isError,
    error: mocks.candidateState.error,
    isRefetching: false,
    refetch: mocks.candidateState.refetch,
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleTagReviewDecisions', () => ({
  useSourceRuleTagReviewDecisionsQuery: (ruleId?: string | null) => ({
    data: ruleId ? mocks.decisions[ruleId] ?? [] : [],
    isLoading: false,
    isError: false,
    error: null,
  }),
  useUpsertSourceRuleTagReviewDecisionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  tagAPI: {
    batchCreate: vi.fn(),
  },
}));

function Bootstrap({ focusedRuleId }: { focusedRuleId?: string }) {
  const { setFocusedRuleId, setSelectedDeviceId } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    if (focusedRuleId) {
      setFocusedRuleId(focusedRuleId);
    }
  }, [focusedRuleId, setFocusedRuleId, setSelectedDeviceId]);

  return null;
}

function renderDesk(focusedRuleId = 'rule-1') {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <Bootstrap focusedRuleId={focusedRuleId} />
        <MuiTagIncidentDesk />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

function seedReadyTagFlow() {
  mocks.devices.splice(0, mocks.devices.length, {
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
  mocks.points.splice(0, mocks.points.length, {
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
  mocks.sourceRules.splice(0, mocks.sourceRules.length, {
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
    revision_id: 'rev-1',
    created_at: '',
    updated_at: '',
  });
  mocks.candidateViews['rule-1'] = {
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
    database_outputs: { status: 'deferred', candidates: [] },
    local_modbus_outputs: { status: 'deferred', candidates: [] },
  };
}

describe('MuiTagIncidentDesk reopened v2 Tag surface', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mocks.devices.splice(0);
    mocks.points.splice(0);
    mocks.tags.splice(0);
    mocks.mappings.splice(0);
    mocks.sourceRules.splice(0);
    Object.keys(mocks.candidateViews).forEach((key) => delete mocks.candidateViews[key]);
    Object.keys(mocks.decisions).forEach((key) => delete mocks.decisions[key]);
    mocks.candidateState.isLoading = false;
    mocks.candidateState.isError = false;
    mocks.candidateState.error = null;
    mocks.candidateState.refetch.mockResolvedValue({ data: null });
  });

  it('renders incident-desk command surfaces while preserving shared Tag workboard panels', async () => {
    seedReadyTagFlow();

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('tag-incident-command-panel')).toBeInTheDocument();
    });
    expect(screen.getByTestId('tag-incident-priority-card')).toBeInTheDocument();
    expect(screen.getByTestId('tag-incident-summary-strip')).toBeInTheDocument();
    expect(screen.getByTestId('tag-incident-handoff-panel')).toBeInTheDocument();
    expect(screen.getByTestId('tag-board-surface')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-tag-review-surface')).toBeInTheDocument();
    expect(screen.getByTestId('tag-candidate-board')).toBeInTheDocument();
    expect(screen.getByTestId('tag-master-surface')).toBeInTheDocument();
  });

  it('surfaces blocker-first recovery states above the shared workboard', async () => {
    seedReadyTagFlow();
    mocks.candidateState.isError = true;
    mocks.candidateState.error = new Error('network timeout');

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('tag-incident-retry-banner')).toBeInTheDocument();
    });
    expect(screen.getByTestId('tag-incident-retry-banner')).toHaveTextContent('network timeout');
    expect(screen.getByTestId('tag-incident-retry-btn')).toBeInTheDocument();
  });

  it('surfaces stale review recovery narrative when source revision drifts', async () => {
    seedReadyTagFlow();
    mocks.sourceRules[0]!.revision_id = 'rev-2';

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('tag-incident-stale-banner')).toBeInTheDocument();
    });
    expect(screen.getByTestId('tag-incident-stale-refresh')).toBeInTheDocument();
  });

  it('keeps output handoff visible and explains blocked readiness', async () => {
    seedReadyTagFlow();
    mocks.tags.splice(0, mocks.tags.length, {
      id: 'tag-1',
      key: 'SRC_40001',
      display_name: 'Flow Sensor',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mocks.mappings.splice(0, mocks.mappings.length, {
      id: 'mapping-1',
      point_id: 'point-1',
      tag_id: 'tag-1',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });

    renderDesk();

    await waitFor(() => {
      expect(screen.getByTestId('tag-incident-handoff-panel')).toBeInTheDocument();
    });
    expect(screen.getByTestId('tag-incident-goto-output')).toBeEnabled();
    expect(screen.getByTestId('tag-incident-handoff-status')).toBeInTheDocument();
  });
});
