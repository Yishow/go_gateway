import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SourceRuleCandidateSnapshotView } from '../../../../types/sourceRuleCandidates';
import type { SourceRuleRecord } from '../../../../types/datalink';
import { WorkbenchProvider, useWorkbench } from '../WorkbenchProvider';
import { SourceRuleDatabaseTargetBoard } from '../SourceRuleDatabaseTargetBoard';

const { mockSourceRules, mockCandidateViews } = vi.hoisted(() => ({
  mockSourceRules: [] as SourceRuleRecord[],
  mockCandidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mockSourceRules,
    isLoading: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mockCandidateViews[ruleId] ?? null : null,
    isLoading: false,
    isFetching: false,
  }),
}));

vi.mock('../DatabaseTargetBoard', () => ({
  DatabaseTargetBoard: ({ reviewSet }: { reviewSet?: SourceRuleCandidateSnapshotView['database_outputs'] | null }) => (
    <pre data-testid="database-target-board-review-set">
      {JSON.stringify(reviewSet?.candidates ?? [])}
    </pre>
  ),
}));

function OutputBootstrap() {
  const { setFocusedRuleId, setSelectedDeviceId, setTagGroupingOverride } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    setFocusedRuleId('rule-1');
    setTagGroupingOverride('rule-1', 'point-1', {
      groupKey: 'line',
      columnName: 'flow_kw',
    });
  }, [setFocusedRuleId, setSelectedDeviceId, setTagGroupingOverride]);

  return <SourceRuleDatabaseTargetBoard candidates={[]} selectedTagId="" />;
}

describe('SourceRuleDatabaseTargetBoard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockSourceRules.splice(0, mockSourceRules.length, {
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
    mockCandidateViews['rule-1'] = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: { status: 'ready', candidates: [] },
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
            proposed_signature: 'db-sig-1',
            address: '40001',
            point_id: 'point-1',
            tag_key: 'meter/A1',
            display_name: 'Flow Sensor',
            data_type: 'int16',
            status: 'ready',
            group_key: 'meter',
            column_name: 'a1',
            write_interval_seconds: 15,
          },
        ],
      },
      local_modbus_outputs: { status: 'ready', candidates: [] },
    };
  });

  it('passes tag-review grouping overrides into the Output review set', async () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <WorkbenchProvider>
          <OutputBootstrap />
        </WorkbenchProvider>
      </QueryClientProvider>,
    );

    await waitFor(() =>
      expect(screen.getByTestId('database-target-board-review-set')).toHaveTextContent('flow_kw'),
    );

    expect(screen.getByTestId('database-target-board-review-set')).toHaveTextContent('"group_key":"line"');
    expect(screen.getByTestId('database-target-board-review-set')).toHaveTextContent('"column_name":"flow_kw"');
  });
});
