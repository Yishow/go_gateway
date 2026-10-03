import { describe, expect, it } from 'vitest';
import { hydrateStudioV2Mapping } from '../../../src/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroupState } from '../../fixtures/writeGroupState';
import type { Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceMappingRecord } from '../../../src/types/datalink';

const point: Point = {
  id: 'rule-1-p-0', device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name: 'temperature',
  address: '40001', data_type: 'float64', function: 'holding_register', width: 1, enabled: true,
  skipped: false, _rule_scale: 1, _rule_offset: 0,
};

const record: StudioV2WorkspaceMappingRecord = {
  id: 'mapping-1', workspace_id: 'ws-1', point_id: 'pt-0', rule_id: 'rule-1', device_id: 'dev-1',
  address: '40001', tag_id: 'tag-1', tag_key: 'line.temperature', display_name: '溫度', unit: '°C',
  target_type: 'float64', scale: 0.1, offset: 0, enabled: true, created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z',
};

function draftFrom(fields: Partial<Mapping>): Mapping {
  return {
    point_id: point.id, tag_key: '', display_name: '', unit: '', target_type: 'float64',
    scale: 1, offset: 0, enabled: true, ...fields,
  };
}

describe('BasicPointTagToGroup', () => {
  it('a saved basic mapping carries the real point/tag/mapping IDs and invents no measurement', () => {
    const hydrated = hydrateStudioV2Mapping(point, record);
    expect(hydrated.persisted).toBe(true);
    expect(hydrated.mapping_id).toBe('mapping-1');
    expect(hydrated.persisted_point_id).toBe('pt-0');
    expect(hydrated.tag_id).toBe('tag-1');
    // Basic writing never fabricates a measurement identity.
    expect(JSON.stringify(hydrated)).not.toContain('measurement');
  });

  it('a Chinese composition draft survives a refetch instead of being clobbered', () => {
    const composing = draftFrom({
      tag_key: 'line.溫', save_state: 'saving',
      local_value: {
        tag_key: 'line.溫', display_name: '', unit: '', target_type: 'float64',
        scale: 1, offset: 0, enabled: true,
      },
    });
    const editing = hydrateStudioV2Mapping(point, record, composing);
    // The row is mid-save with an unfinished composition; hydration must keep the
    // local draft instead of overwriting the input with the stale saved value.
    expect(editing.tag_key).toBe('line.溫');
    expect(editing.save_state).toBe('saving');
  });

  it('once saved, the final composed value replaces the draft on reload', () => {
    const savedDraft = draftFrom({ tag_key: 'line.溫度', tag_id: 'tag-1', mapping_id: 'mapping-1' });
    const reloaded = hydrateStudioV2Mapping(point, record, savedDraft);
    expect(reloaded.tag_key).toBe('line.temperature');
    expect(reloaded.save_state).toBe('saved');
  });

  it('basic typed tags join a write group without any measurement semantics', () => {
    const { candidates, excluded } = buildGroupCandidates(savedGroupState());
    expect(excluded).toEqual([]);
    expect(candidates.map((candidate) => candidate.tag_key)).toEqual([
      'line.temperature', 'line.pressure', 'line.batch',
    ]);
    for (const candidate of candidates) {
      // Group membership uses the persisted point and tag IDs only.
      expect(candidate.point_id).toMatch(/^pt-\d$/);
      expect(candidate.tag_id).toMatch(/^tag-\d$/);
      expect(JSON.stringify(candidate)).not.toContain('measurement');
    }
  });
});
