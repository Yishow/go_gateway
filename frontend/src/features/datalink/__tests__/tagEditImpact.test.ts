import { describe, expect, it } from 'vitest';
import type { Mapping, Tag } from '../../../types/datalink';
import {
  buildGlobalTagEditDraft,
  getAffectedMappingCountForTag,
  getAffectedMappingsForTag,
  hasGlobalTagEditChanges,
  toTagUpdateRequest,
} from '../tagEditImpact';

const sampleTag: Tag = {
  id: 'tag-1',
  key: 'line_a_temp',
  display_name: 'Line A Temp',
  description: 'temperature',
  data_type: 'float32',
  unit: 'C',
  labels: null,
  status: 'active',
  created_at: '',
  updated_at: '',
};

const mappings: Mapping[] = [
  { id: 'm-1', point_id: 'p-1', tag_id: 'tag-1', enabled: true, transform_pipeline: '[]', created_at: '', updated_at: '' },
  { id: 'm-2', point_id: 'p-2', tag_id: 'tag-2', enabled: true, transform_pipeline: '[]', created_at: '', updated_at: '' },
  { id: 'm-3', point_id: 'p-3', tag_id: 'tag-1', enabled: false, transform_pipeline: '[]', created_at: '', updated_at: '' },
];

describe('tagEditImpact', () => {
  it('finds all mappings affected by the same global tag', () => {
    const affected = getAffectedMappingsForTag(mappings, 'tag-1');
    expect(affected.map((mapping) => mapping.id)).toEqual(['m-1', 'm-3']);
    expect(getAffectedMappingCountForTag(mappings, 'tag-1')).toBe(2);
  });

  it('builds trimmed draft and detects changes correctly', () => {
    const draft = buildGlobalTagEditDraft({
      display_name: '  Line A Temperature  ',
      unit: ' C ',
      description: ' updated ',
    });

    expect(draft).toEqual({
      display_name: 'Line A Temperature',
      unit: 'C',
      description: 'updated',
    });
    expect(hasGlobalTagEditChanges(sampleTag, draft)).toBe(true);
  });

  it('returns no change when draft equals original tag values', () => {
    const draft = buildGlobalTagEditDraft({
      display_name: sampleTag.display_name,
      unit: sampleTag.unit,
      description: sampleTag.description,
    });

    expect(hasGlobalTagEditChanges(sampleTag, draft)).toBe(false);
  });

  it('converts draft to update payload for API mutation', () => {
    const draft = buildGlobalTagEditDraft({
      display_name: 'Name',
      unit: 'kPa',
      description: 'desc',
    });

    expect(toTagUpdateRequest(draft)).toEqual({
      display_name: 'Name',
      unit: 'kPa',
      description: 'desc',
    });
  });
});
