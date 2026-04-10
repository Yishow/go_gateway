import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type {
  DatabaseConnector,
  DatabaseTableInfo,
  DatabaseTargetMapping,
  DatabaseTargetValidationResult,
  SourceRuleRecord,
} from '../../../../src/types/datalink';
import { dbTargetAPI } from '../../../../src/services/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../../src/types/sourceRuleCandidates';
import { SourceRuleDatabaseTargetBoard } from '../../../../src/pages/datalink/workbench/SourceRuleDatabaseTargetBoard';
import {
  WorkbenchProvider,
  useWorkbench,
} from '../../../../src/pages/datalink/workbench/WorkbenchProvider';

const { mockCandidateSnapshot, mockDBTargetAPI, mockPersistedRules, mockCandidateRuleId } = vi.hoisted(() => ({
  mockCandidateSnapshot: {
    value: null as SourceRuleCandidateSnapshotView | null,
  },
  mockPersistedRules: {
    value: [] as SourceRuleRecord[],
  },
  mockCandidateRuleId: {
    value: null as string | null,
  },
  mockDBTargetAPI: {
    listConnectors: vi.fn(),
    getConnector: vi.fn(),
    createConnector: vi.fn(),
    updateConnector: vi.fn(),
    deleteConnector: vi.fn(),
    testConnector: vi.fn(),
    generateSchema: vi.fn(),
    dryRunMappings: vi.fn(),
    listTables: vi.fn(),
    validateConnector: vi.fn(),
    listMappings: vi.fn(),
    getMapping: vi.fn(),
    createMapping: vi.fn(),
    updateMapping: vi.fn(),
    deleteMapping: vi.fn(),
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: Record<string, string | number>) =>
      params ? `${key}:${JSON.stringify(params)}` : key,
  }),
}));

vi.mock('../../../../src/services/datalink', () => ({
  dbTargetAPI: mockDBTargetAPI,
}));

vi.mock('../../../../src/hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => {
    mockCandidateRuleId.value = ruleId ?? null;
    return {
      data: mockCandidateSnapshot.value,
      isLoading: false,
      isFetching: false,
    };
  },
}));

vi.mock('../../../../src/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: mockPersistedRules.value,
    isLoading: false,
    isSuccess: true,
  }),
}));

function FocusedRuleBootstrap({ ruleId }: { ruleId: string }) {
  const { setFocusedRuleId, setSelectedDeviceId } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    setFocusedRuleId(ruleId);
  }, [ruleId, setFocusedRuleId, setSelectedDeviceId]);

  return null;
}

function SourcePlanningRuleBootstrap({ ruleId }: { ruleId: string }) {
  const { setSelectedDeviceId, setSourcePlanningState } = useWorkbench();

  useEffect(() => {
    setSelectedDeviceId('device-1');
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedRuleId: ruleId,
    }));
  }, [ruleId, setSelectedDeviceId, setSourcePlanningState]);

  return null;
}

function buildConnector(id: string, name: string): DatabaseConnector {
  return {
    id,
    name,
    kind: 'sqlite',
    connection_config: { dsn: `/tmp/${id}.db` },
    status: 'ready',
    last_check_at: null,
    last_check_error: '',
    enabled: true,
    default_write_interval_seconds: 15,
    created_at: '',
    updated_at: '',
  };
}

function buildTable(schema: string, name: string): DatabaseTableInfo {
  return {
    schema,
    name,
    columns: [
      {
        name: 'value',
        data_type: 'REAL',
        nullable: false,
        primary_key: false,
      },
    ],
  };
}

function buildValidation(): DatabaseTargetValidationResult {
  return {
    ready: true,
    issues: [],
  };
}

function renderBoard(options?: { focusedRuleId?: string | null; sourcePlanningRuleId?: string | null }) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        {options?.focusedRuleId === null ? null : (
          <FocusedRuleBootstrap ruleId={options?.focusedRuleId ?? 'rule-1'} />
        )}
        {options?.sourcePlanningRuleId ? (
          <SourcePlanningRuleBootstrap ruleId={options.sourcePlanningRuleId} />
        ) : null}
        <SourceRuleDatabaseTargetBoard
          candidates={[
            {
              tagId: 'tag-1',
              tagKey: 'TAG_FLOW',
              pointName: 'Flow',
              pointAddress: '40001',
              dataType: 'int16',
              lastValue: 12,
            },
            {
              tagId: 'tag-2',
              tagKey: 'TAG_PRESSURE',
              pointName: 'Pressure',
              pointAddress: '40002',
              dataType: 'int16',
              lastValue: 18,
            },
          ]}
          selectedTagId="tag-1"
        />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

describe('SourceRuleDatabaseTargetBoard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCandidateRuleId.value = null;
    mockPersistedRules.value = [{
      id: 'rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 1,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: [],
      created_at: '',
      updated_at: '',
    }];
    mockCandidateSnapshot.value = {
      source_rule_id: 'rule-1',
      revision_id: 'rev-1',
      tags: {
        status: 'ready',
        candidates: [],
      },
      database_outputs: {
        status: 'blocked',
        reason: 'Connector schema drift detected',
        candidates: [
          {
            id: 'db-candidate-unbound',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'database_output',
              candidate_kind: 'database_output',
              derived_from_rule_address: '40003',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-0',
            address: '40003',
            point_id: 'point-3',
            tag_key: 'TAG_TEMPERATURE',
            display_name: 'Temperature',
            data_type: 'int16',
            status: 'ready',
          },
          {
            id: 'db-candidate-1',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'database_output',
              candidate_kind: 'database_output',
              derived_from_rule_address: '40001',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-1',
            address: '40001',
            point_id: 'point-1',
            tag_id: 'tag-1',
            mapping_id: 'mapping-1',
            tag_key: 'TAG_FLOW',
            display_name: 'Flow',
            data_type: 'int16',
            status: 'ready',
            connector_id: 'connector-1',
            table_schema: 'public',
            table_name: 'flow_metrics',
            column_name: 'value',
            write_mode: 'insert',
          },
          {
            id: 'db-candidate-2',
            identity: {
              source_rule_id: 'rule-1',
              candidate_type: 'database_output',
              candidate_kind: 'database_output',
              derived_from_rule_address: '40002',
              target_binding_scope: [],
            },
            proposed_signature: 'sig-2',
            address: '40002',
            point_id: 'point-2',
            tag_id: 'tag-2',
            mapping_id: 'mapping-2',
            tag_key: 'TAG_PRESSURE',
            display_name: 'Pressure',
            data_type: 'int16',
            status: 'out_of_sync',
            blocking_reason: 'Column value missing',
            connector_id: 'connector-2',
            table_schema: 'public',
            table_name: 'pressure_metrics',
            column_name: 'value',
            write_mode: 'upsert',
            timestamp_column: 'ts',
          },
        ],
      },
      local_modbus_outputs: {
        status: 'deferred',
        candidates: [],
      },
    };

    vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([
      buildConnector('connector-1', 'Primary DB'),
      buildConnector('connector-2', 'Backup DB'),
    ]);
    mockDBTargetAPI.generateSchema.mockResolvedValue({
      connector_id: 'connector-1',
      dry_run: false,
      statements: ['CREATE TABLE flow_metrics (value REAL NOT NULL)'],
      executed: 1,
    });
    mockDBTargetAPI.dryRunMappings.mockResolvedValue({
      connector_id: 'connector-1',
      results: [
        {
          candidate_id: 'db-mapping-1',
          status: 'ready',
          mapping_id: 'db-mapping-1',
          tag_id: 'tag-1',
        },
      ],
    });
    vi.mocked(dbTargetAPI.listMappings).mockResolvedValue([]);
    vi.mocked(dbTargetAPI.listTables).mockImplementation(async (connectorId) => {
      if (connectorId === 'connector-2') {
        return [buildTable('public', 'pressure_metrics')];
      }
      return [buildTable('public', 'flow_metrics')];
    });
    vi.mocked(dbTargetAPI.validateConnector).mockResolvedValue(buildValidation());
  });

  it('filters the review surface to the currently selected connector context', async () => {
    renderBoard();

    await screen.findByTestId('database-output-review-surface');

    expect(screen.getByTestId('database-review-candidate-db-candidate-1')).toBeInTheDocument();
    expect(
      screen.getByTestId('database-review-candidate-db-candidate-unbound'),
    ).toBeInTheDocument();
    expect(
      screen.queryByTestId('database-review-candidate-db-candidate-2'),
    ).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: /Backup DB/ }));

    await waitFor(() => {
      expect(
        screen.getByTestId('database-review-candidate-db-candidate-2'),
      ).toBeInTheDocument();
    });
    expect(
      screen.queryByTestId('database-review-candidate-db-candidate-1'),
    ).not.toBeInTheDocument();
    expect(
      screen.getByTestId('database-review-candidate-db-candidate-unbound'),
    ).toBeInTheDocument();
  });

  it('hides connector-bound candidates when no connector is selected', async () => {
    vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
    vi.mocked(dbTargetAPI.listTables).mockResolvedValue([]);

    renderBoard();

    await screen.findByTestId('database-output-review-surface');

    expect(
      screen.getByTestId('database-review-candidate-db-candidate-unbound'),
    ).toBeInTheDocument();
    expect(
      screen.queryByTestId('database-review-candidate-db-candidate-1'),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByTestId('database-review-candidate-db-candidate-2'),
    ).not.toBeInTheDocument();
  });

  it('shows revision and blocking evidence from the rule-scoped database candidate set', async () => {
    renderBoard();

    await screen.findByTestId('database-output-review-surface');

    expect(screen.getByTestId('database-review-revision')).toHaveTextContent('rev-1');
    expect(screen.getByTestId('database-review-set-status')).toHaveTextContent(
      'workbench.output.database.reviewSurface.status.blocked',
    );
    expect(
      screen.getByText((content) => content.includes('Connector schema drift detected')),
    ).toBeInTheDocument();
  });

  it('falls back to the persisted source-planning rule when cross-step focus is empty', async () => {
    renderBoard({
      focusedRuleId: null,
      sourcePlanningRuleId: 'rule-1',
    });

    await screen.findByTestId('database-output-review-surface');

    expect(mockCandidateRuleId.value).toBe('rule-1');
    expect(screen.getByTestId('database-review-revision')).toHaveTextContent('rev-1');
  });

  it('keeps rendering when validation refresh returns ready with null issues', async () => {
    vi.mocked(dbTargetAPI.validateConnector).mockReset();
    vi.mocked(dbTargetAPI.validateConnector)
      .mockResolvedValueOnce(buildValidation())
      .mockResolvedValueOnce({
        ready: true,
        issues: null,
      } as unknown as DatabaseTargetValidationResult);

    renderBoard();

    await screen.findByTestId('database-output-review-surface');

    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.database.actions.refreshValidation',
      }),
    );

    await waitFor(() => {
      expect(dbTargetAPI.validateConnector).toHaveBeenCalledTimes(2);
    });
    expect(screen.getByTestId('database-output-review-surface')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent(
      'workbench.output.database.results.validationRefreshed',
    );
  });

  it('skips review queries for a draft-focused rule id', async () => {
    mockPersistedRules.value = [];
    mockCandidateSnapshot.value = null;

    renderBoard();

    await screen.findByText('workbench.output.database.reviewSurface.noRule');
    expect(mockCandidateRuleId.value).toBeNull();
  });

  it('can generate schema when the selected connector has no tables yet', async () => {
    let tables: DatabaseTableInfo[] = [];

    vi.mocked(dbTargetAPI.listTables).mockImplementation(async () => tables);
    mockDBTargetAPI.generateSchema.mockImplementation(async () => {
      tables = [buildTable('public', 'flow_metrics')];
      return {
        connector_id: 'connector-1',
        dry_run: false,
        statements: ['CREATE TABLE flow_metrics (value REAL NOT NULL)'],
        executed: 1,
      };
    });

    renderBoard();

    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.database.actions.generateSchema',
      }),
    );

    await waitFor(() => {
      expect(mockDBTargetAPI.generateSchema).toHaveBeenCalledWith('connector-1', {
        dry_run: false,
      });
    });
    await waitFor(() => {
      expect(screen.getByText('public.flow_metrics')).toBeInTheDocument();
    });
  });

  it('reports when schema generation produces no changes for the current scope', async () => {
    vi.mocked(dbTargetAPI.listTables).mockResolvedValue([]);
    mockDBTargetAPI.generateSchema.mockResolvedValue({
      connector_id: 'connector-1',
      dry_run: false,
      statements: [],
      executed: 0,
    });

    renderBoard();

    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.database.actions.generateSchema',
      }),
    );

    await waitFor(() => {
      expect(mockDBTargetAPI.generateSchema).toHaveBeenCalledWith('connector-1', {
        dry_run: false,
      });
    });
    await screen.findByText('workbench.output.database.results.schemaGenerateNoChanges');
  });

  it('runs a dry-run check after binding a tag to a database column', async () => {
    let mappings: DatabaseTargetMapping[] = [];

    vi.mocked(dbTargetAPI.listMappings).mockImplementation(async () => mappings);
    vi.mocked(dbTargetAPI.createMapping).mockImplementation(async () => {
      mappings = [{
        id: 'db-mapping-1',
        tag_id: 'tag-1',
        connector_id: 'connector-1',
        table_schema: 'public',
        table_name: 'flow_metrics',
        column_name: 'value',
        write_mode: 'insert',
        timestamp_column: null,
        enabled: true,
        created_at: '',
        updated_at: '',
      }];
      return mappings[0];
    });

    renderBoard();

    await screen.findByTestId('schema-column-value');

    fireEvent.change(screen.getByLabelText('workbench.output.database.mapping.table'), {
      target: { value: 'public.flow_metrics' },
    });
    fireEvent.click(screen.getByTestId('schema-column-value'));

    await waitFor(() => {
      expect(mockDBTargetAPI.dryRunMappings).toHaveBeenCalledWith('connector-1', {
        candidate_ids: ['db-mapping-1'],
      });
    });
    expect(screen.getByTestId('database-dry-run-results')).toBeInTheDocument();
  });
});
