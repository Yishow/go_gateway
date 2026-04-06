import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { SourceRuleLocalModbusReviewSurface } from '../../../../src/pages/datalink/workbench/SourceRuleLocalModbusReviewSurface';
import {
  WorkbenchProvider,
  useWorkbench,
} from '../../../../src/pages/datalink/workbench/WorkbenchProvider';
import type { ModbusShareStatus } from '../../../../src/types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../../src/types/sourceRuleCandidates';

const { mockCandidateSnapshot } = vi.hoisted(() => ({
  mockCandidateSnapshot: {
    value: null as SourceRuleCandidateSnapshotView | null,
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: Record<string, string | number>) =>
      params ? `${key}:${JSON.stringify(params)}` : key,
  }),
}));

vi.mock('../../../../src/hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: () => ({
    data: mockCandidateSnapshot.value,
    isLoading: false,
    isFetching: false,
  }),
}));

function FocusedRuleBootstrap({ ruleId }: { ruleId: string }) {
  const { setFocusedRuleId } = useWorkbench();

  useEffect(() => {
    setFocusedRuleId(ruleId);
  }, [ruleId, setFocusedRuleId]);

  return null;
}

function buildStatus(enabled: boolean): ModbusShareStatus {
  return {
    enabled,
    port: 1502,
    address: '127.0.0.1:1502',
    bind_state: 'pass',
    mapping_count: 3,
  };
}

function renderSurface(options?: {
  focusedRuleId?: string | null;
  selectedTagId?: string;
  status?: ModbusShareStatus | null;
  conflictCount?: number;
}) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        {options?.focusedRuleId ? <FocusedRuleBootstrap ruleId={options.focusedRuleId} /> : null}
        <SourceRuleLocalModbusReviewSurface
          selectedTagId={options?.selectedTagId ?? 'tag-1'}
          status={options?.status ?? buildStatus(false)}
          conflictCount={options?.conflictCount ?? 2}
        />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

describe('SourceRuleLocalModbusReviewSurface', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCandidateSnapshot.value = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [],
      },
      database_outputs: {
        status: 'deferred',
        candidates: [],
      },
      local_modbus_outputs: {
        status: 'deferred',
        candidates: [
          {
            id: 'lm-candidate-1',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'local_modbus_output',
              candidate_kind: 'local_modbus_output',
              derived_from_rule_address: '40001',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-1',
            address: '40001',
            point_id: 'point-1',
            tag_id: 'tag-1',
            tag_key: 'TAG_FLOW',
            display_name: 'Flow',
            data_type: 'int32',
            register: 10,
            register_count: 2,
            status: 'blocked_conflict',
            blocking_reason: 'Overlaps HR40012 with rule-2',
          },
          {
            id: 'lm-candidate-2',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'local_modbus_output',
              candidate_kind: 'local_modbus_output',
              derived_from_rule_address: '40002',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-2',
            address: '40002',
            point_id: 'point-2',
            tag_id: 'tag-2',
            tag_key: 'TAG_PRESSURE',
            display_name: 'Pressure',
            data_type: 'int16',
            register: 20,
            register_count: 1,
            status: 'out_of_sync',
            blocking_reason: 'Override tag missing',
          },
          {
            id: 'lm-candidate-3',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'local_modbus_output',
              candidate_kind: 'local_modbus_output',
              derived_from_rule_address: '40003',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-3',
            address: '40003',
            point_id: 'point-3',
            tag_key: 'TAG_RENAMED',
            display_name: 'Temperature',
            data_type: 'int16',
            register: 30,
            register_count: 1,
            status: 'ready',
          },
        ],
      },
    };
  });

  it('shows conflict reasons, verification state, and runtime health inline', () => {
    renderSurface({ focusedRuleId: 'rule-1' });

    expect(screen.getByTestId('local-modbus-review-surface')).toBeTruthy();
    expect(screen.getByText('TAG_FLOW')).toBeTruthy();
    expect(screen.getByText('TAG_PRESSURE')).toBeTruthy();
    expect(screen.getByText('Overlaps HR40012 with rule-2')).toBeTruthy();
    expect(screen.getByText('Override tag missing')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.serverStopped')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.summary.blocked:{"count":2}')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.currentSelection')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.status.blocked_conflict')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.status.out_of_sync')).toBeTruthy();
    expect(screen.getByText('workbench.output.modbusReviewSurface.verification.pending')).toBeTruthy();
  });

  it('shows the focus hint when no rule is selected', () => {
    renderSurface({
      focusedRuleId: null,
      selectedTagId: '',
      status: buildStatus(true),
      conflictCount: 0,
    });

    expect(screen.getByText('workbench.output.modbusReviewSurface.noRule')).toBeTruthy();
  });
});
