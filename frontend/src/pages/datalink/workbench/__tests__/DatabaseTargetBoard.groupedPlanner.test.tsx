import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { WorkbenchProvider } from '../WorkbenchProvider';
import { DatabaseTargetBoard } from '../DatabaseTargetBoard';
import type { WorkbenchOutputCandidate } from '../workbenchOutputTypes';
import type { DatabaseConnector, DatabaseTargetMapping, DatabaseTableInfo } from '../../../../types/datalink';
import type {
  SourceRuleCandidateSetView,
  SourceRuleDatabaseOutputCandidateView,
} from '../../../../types/sourceRuleCandidates';

const mocks = vi.hoisted(() => ({
  connectors: [] as DatabaseConnector[],
  mappings: [] as DatabaseTargetMapping[],
  tables: [] as DatabaseTableInfo[],
  createMapping: vi.fn(),
  updateMapping: vi.fn(),
  deleteMapping: vi.fn(),
  listConnectors: vi.fn(),
  listMappings: vi.fn(),
  listTables: vi.fn(),
  validateConnector: vi.fn(),
  dryRunMappings: vi.fn(),
  generateSchema: vi.fn(),
  createConnector: vi.fn(),
  updateConnector: vi.fn(),
  deleteConnector: vi.fn(),
  testConnector: vi.fn(),
  refreshReview: vi.fn(),
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: Record<string, string | number>) => {
      if (key === 'workbench.output.database.rowPlanner.issue.interval_mismatch') {
        return 'interval mismatch';
      }
      if (params && typeof params.defaultValue === 'string') {
        return params.defaultValue;
      }
      return key;
    },
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  dbTargetAPI: {
    listConnectors: mocks.listConnectors,
    getConnector: vi.fn(),
    createConnector: mocks.createConnector,
    updateConnector: mocks.updateConnector,
    deleteConnector: mocks.deleteConnector,
    testConnector: mocks.testConnector,
    generateSchema: mocks.generateSchema,
    dryRunMappings: mocks.dryRunMappings,
    listTables: mocks.listTables,
    validateConnector: mocks.validateConnector,
    listMappings: mocks.listMappings,
    getMapping: vi.fn(),
    createMapping: mocks.createMapping,
    updateMapping: mocks.updateMapping,
    deleteMapping: mocks.deleteMapping,
  },
}));

vi.mock('../useRefreshSourceRuleCandidates', () => ({
  useRefreshSourceRuleCandidates: () => mocks.refreshReview,
}));

function createReviewCandidate(
  overrides: Partial<SourceRuleDatabaseOutputCandidateView>,
): SourceRuleDatabaseOutputCandidateView {
  return {
    id: 'db-candidate-1',
    identity: {
      source_rule_id: 'rule-1',
      candidate_type: 'database_outputs',
      candidate_kind: 'database_output',
      derived_from_rule_address: '40001',
    },
    proposed_signature: 'sig-1',
    address: '40001',
    point_id: 'point-1',
    tag_id: 'tag-1',
    tag_key: 'meter/A1',
    display_name: 'Flow A1',
    data_type: 'int16',
    status: 'ready',
    group_key: 'meter',
    column_name: 'a1',
    write_mode: 'insert',
    timestamp_column: '',
    write_interval_seconds: 15,
    ...overrides,
  };
}

const candidates: WorkbenchOutputCandidate[] = [
  {
    tagId: 'tag-1',
    tagKey: 'meter/A1',
    pointName: 'Flow A1',
    pointAddress: '40001',
    dataType: 'int16',
    lastValue: 12,
  },
  {
    tagId: 'tag-2',
    tagKey: 'meter/kw',
    pointName: 'Flow KW',
    pointAddress: '40002',
    dataType: 'int16',
    lastValue: 28,
  },
];

const reviewSet: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> = {
  status: 'ready',
  candidates: [
    createReviewCandidate({
      id: 'db-candidate-1',
      point_id: 'point-1',
      tag_id: 'tag-1',
      column_name: 'a1',
      tag_key: 'meter/A1',
    }),
    createReviewCandidate({
      id: 'db-candidate-2',
      point_id: 'point-2',
      tag_id: 'tag-2',
      column_name: 'kw',
      tag_key: 'meter/kw',
      display_name: 'Flow KW',
      address: '40002',
    }),
  ],
};

describe('DatabaseTargetBoard grouped planner', () => {
  const groupedRowId = 'group-db-candidate-1-db-candidate-2';

  beforeEach(() => {
    vi.clearAllMocks();
    mocks.connectors.splice(0, mocks.connectors.length, {
      id: 'connector-1',
      name: 'Main SQLite',
      kind: 'sqlite',
      connection_config: { dsn: '/tmp/target.db' },
      status: 'ready',
      last_check_at: '',
      last_check_error: '',
      enabled: true,
      default_write_interval_seconds: 15,
      created_at: '',
      updated_at: '',
    });
    mocks.mappings.splice(0, mocks.mappings.length);
    mocks.tables.splice(0, mocks.tables.length, {
      schema: 'main',
      name: 'sensor_values',
      columns: [
        { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true },
        { name: 'value', data_type: 'real', nullable: false, primary_key: false },
        { name: 'kw', data_type: 'real', nullable: false, primary_key: false },
      ],
    });
    mocks.listConnectors.mockResolvedValue(mocks.connectors);
    mocks.listMappings.mockResolvedValue(mocks.mappings);
    mocks.listTables.mockResolvedValue(mocks.tables);
    mocks.validateConnector.mockResolvedValue({ ready: true, issues: [] });
    mocks.createMapping.mockResolvedValueOnce({
      id: 'mapping-1',
      tag_id: 'tag-1',
      connector_id: 'connector-1',
      table_schema: 'main',
      table_name: 'sensor_values',
      column_name: 'flow_kw',
      write_mode: 'insert',
      timestamp_column: null,
      group_key: 'line',
      write_interval_seconds: 30,
      enabled: true,
      created_at: '',
      updated_at: '',
    });
    mocks.createMapping.mockResolvedValueOnce({
      id: 'mapping-2',
      tag_id: 'tag-2',
      connector_id: 'connector-1',
      table_schema: 'main',
      table_name: 'sensor_values',
      column_name: 'kw',
      write_mode: 'insert',
      timestamp_column: null,
      group_key: 'line',
      write_interval_seconds: 30,
      enabled: true,
      created_at: '',
      updated_at: '',
    });
    mocks.dryRunMappings.mockResolvedValue({
      connector_id: 'connector-1',
      results: [
        { candidate_id: 'mapping-1', status: 'ready' },
        { candidate_id: 'mapping-2', status: 'ready' },
      ],
    });
    mocks.refreshReview.mockResolvedValue(undefined);
  });

  it('applies a grouped row with overridden group, columns, and interval', async () => {
    render(
      <WorkbenchProvider>
        <DatabaseTargetBoard
          candidates={candidates}
          selectedTagId="tag-1"
          reviewRuleId="rule-1"
          reviewRevisionId="rev-1"
          reviewSet={reviewSet}
          reviewLoading={false}
        />
      </WorkbenchProvider>,
    );

    await waitFor(() =>
      expect(screen.getByTestId(`database-row-plan-${groupedRowId}`)).toBeInTheDocument(),
    );
    await waitFor(() =>
      expect(screen.getByTestId(`database-row-apply-${groupedRowId}`)).not.toBeDisabled(),
    );

    fireEvent.change(screen.getByTestId(`database-row-interval-${groupedRowId}`), {
      target: { value: '30' },
    });
    fireEvent.change(screen.getByTestId(`database-row-group-${groupedRowId}`), {
      target: { value: 'line' },
    });
    fireEvent.change(screen.getByTestId('database-row-column-db-candidate-1'), {
      target: { value: 'flow_kw' },
    });
    fireEvent.click(screen.getByTestId(`database-row-apply-${groupedRowId}`));

    await waitFor(() => expect(mocks.createMapping).toHaveBeenCalledTimes(2));

    expect(mocks.createMapping).toHaveBeenNthCalledWith(
      1,
      expect.objectContaining({
        tag_id: 'tag-1',
        connector_id: 'connector-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        group_key: 'line',
        column_name: 'flow_kw',
        write_interval_seconds: 30,
      }),
    );
    expect(mocks.createMapping).toHaveBeenNthCalledWith(
      2,
      expect.objectContaining({
        tag_id: 'tag-2',
        connector_id: 'connector-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        group_key: 'line',
        column_name: 'kw',
        write_interval_seconds: 30,
      }),
    );
  });
});
