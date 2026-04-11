import { describe, expect, it } from 'vitest';
import type { SourceRuleDatabaseOutputCandidateView } from '../../../../types/sourceRuleCandidates';
import { buildDatabaseRowPlans } from '../databaseRowPlannerModel';

function createCandidate(
  overrides: Partial<SourceRuleDatabaseOutputCandidateView>,
): SourceRuleDatabaseOutputCandidateView {
  return {
    id: 'candidate-1',
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
    connector_id: 'connector-1',
    table_schema: 'main',
    table_name: 'sensor_values',
    column_name: 'a1',
    group_key: 'meter',
    write_mode: 'insert',
    timestamp_column: '',
    write_interval_seconds: 15,
    ...overrides,
  };
}

describe('buildDatabaseRowPlans', () => {
  it('groups compatible members into one row plan', () => {
    const rowPlans = buildDatabaseRowPlans({
      candidates: [
        createCandidate({ id: 'candidate-1', point_id: 'point-1', tag_id: 'tag-1', tag_key: 'meter/A1', column_name: 'a1' }),
        createCandidate({ id: 'candidate-2', point_id: 'point-2', tag_id: 'tag-2', tag_key: 'meter/kw', column_name: 'kw', display_name: 'Flow KW', address: '40002' }),
      ],
      selectedConnectorId: '',
      selectedTableKey: '',
      selectedWriteMode: 'insert',
      selectedTimestampColumn: '',
      connectorDefaultWriteIntervalSeconds: 15,
    });

    expect(rowPlans).toHaveLength(1);
    expect(rowPlans[0]).toMatchObject({
      id: 'group-candidate-1-candidate-2',
      kind: 'grouped',
      groupKey: 'meter',
      intervalSeconds: 15,
      status: 'ready',
      issueCodes: [],
    });
    expect(rowPlans[0].members.map((member) => member.columnName)).toEqual(['a1', 'kw']);
  });

  it('blocks a grouped row when member intervals disagree', () => {
    const rowPlans = buildDatabaseRowPlans({
      candidates: [
        createCandidate({ id: 'candidate-1', point_id: 'point-1', write_interval_seconds: 15 }),
        createCandidate({
          id: 'candidate-2',
          point_id: 'point-2',
          tag_id: 'tag-2',
          tag_key: 'meter/kw',
          column_name: 'kw',
          display_name: 'Flow KW',
          address: '40002',
          write_interval_seconds: 30,
        }),
      ],
      selectedConnectorId: '',
      selectedTableKey: '',
      selectedWriteMode: 'insert',
      selectedTimestampColumn: '',
      connectorDefaultWriteIntervalSeconds: 15,
    });

    expect(rowPlans).toHaveLength(1);
    expect(rowPlans[0].status).toBe('blocked');
    expect(rowPlans[0].issueCodes).toContain('interval_mismatch');
  });

  it('keeps slashless or legacy mappings as single-member rows', () => {
    const rowPlans = buildDatabaseRowPlans({
      candidates: [
        createCandidate({
          id: 'candidate-3',
          point_id: 'point-3',
          tag_id: 'tag-3',
          tag_key: 'FLOW',
          group_key: null,
          column_name: 'flow',
          display_name: 'Flow Legacy',
          address: '40003',
        }),
      ],
      selectedConnectorId: '',
      selectedTableKey: '',
      selectedWriteMode: 'insert',
      selectedTimestampColumn: '',
      connectorDefaultWriteIntervalSeconds: 15,
    });

    expect(rowPlans).toHaveLength(1);
    expect(rowPlans[0]).toMatchObject({
      id: 'single-candidate-3',
      kind: 'single',
      groupKey: null,
      status: 'ready',
    });
    expect(rowPlans[0].members).toHaveLength(1);
  });
});
