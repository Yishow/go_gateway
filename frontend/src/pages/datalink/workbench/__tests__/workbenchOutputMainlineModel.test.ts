import { describe, expect, it } from 'vitest';
import type { SourceRuleCandidateSnapshotView } from '../../../../types/sourceRuleCandidates';
import { buildWorkbenchOutputMainlineState } from '../workbenchOutputMainlineModel';

function createSnapshot(
  overrides?: Partial<SourceRuleCandidateSnapshotView>,
): SourceRuleCandidateSnapshotView {
  return {
    source_rule_id: 'rule-1',
    revision_id: 'revision-1',
    tags: {
      status: 'ready',
      candidates: [],
    },
    database_outputs: {
      status: 'ready',
      candidates: [],
    },
    local_modbus_outputs: {
      status: 'ready',
      candidates: [],
    },
    ...overrides,
  };
}

describe('buildWorkbenchOutputMainlineState', () => {
  it('requires an active rule handoff before Output can claim the recovery surface', () => {
    expect(
      buildWorkbenchOutputMainlineState({
        hasSelectedDevice: true,
        sourceReady: true,
        tagReady: true,
        activeRuleId: null,
        snapshot: null,
      }),
    ).toMatchObject({
      outputReady: false,
      outputReadiness: { status: 'partial', reason: 'output-rule-required' },
      databaseConfigured: false,
      localModbusConfigured: false,
    });
  });

  it('marks Output mainline ready when a database mapping has been applied', () => {
    expect(
      buildWorkbenchOutputMainlineState({
        hasSelectedDevice: true,
        sourceReady: true,
        tagReady: true,
        activeRuleId: 'rule-1',
        snapshot: createSnapshot({
          database_outputs: {
            status: 'ready',
            candidates: [
              {
                id: 'database-candidate-1',
                identity: {
                  source_rule_id: 'rule-1',
                  candidate_type: 'database_output',
                  candidate_kind: 'database',
                  derived_from_rule_address: '40001',
                },
                proposed_signature: 'db-1',
                address: '40001',
                point_id: 'point-1',
                tag_id: 'tag-1',
                mapping_id: 'mapping-1',
                tag_key: 'meter/kw',
                display_name: 'Power',
                data_type: 'int16',
                status: 'ready',
              },
            ],
          },
        }),
      }),
    ).toMatchObject({
      outputReady: true,
      outputReadiness: { status: 'ready' },
      databaseConfigured: true,
      localModbusConfigured: false,
    });
  });

  it('treats a Local Modbus register at HR0 as a valid mainline-complete binding', () => {
    expect(
      buildWorkbenchOutputMainlineState({
        hasSelectedDevice: true,
        sourceReady: true,
        tagReady: true,
        activeRuleId: 'rule-1',
        snapshot: createSnapshot({
          local_modbus_outputs: {
            status: 'ready',
            candidates: [
              {
                id: 'modbus-candidate-1',
                identity: {
                  source_rule_id: 'rule-1',
                  candidate_type: 'local_modbus_output',
                  candidate_kind: 'local_modbus',
                  derived_from_rule_address: '40001',
                },
                proposed_signature: 'modbus-1',
                address: '40001',
                point_id: 'point-1',
                tag_id: 'tag-1',
                tag_key: 'meter/kw',
                display_name: 'Power',
                data_type: 'int16',
                register: 0,
                register_count: 1,
                status: 'ready',
              },
            ],
          },
        }),
      }),
    ).toMatchObject({
      outputReady: true,
      outputReadiness: { status: 'ready' },
      databaseConfigured: false,
      localModbusConfigured: true,
    });
  });
});
