import { describe, expect, it } from 'vitest';
import type { Mapping, Point, Tag } from '../../../../types/datalink';
import {
  buildTagBindingCandidates,
  buildTagBindingRequests,
  buildTagKey,
} from '../tagBindingModel';

function createPoint(overrides: Partial<Point>): Point {
  return {
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
    ...overrides,
  };
}

function createTag(overrides: Partial<Tag>): Tag {
  return {
    id: 'tag-1',
    key: 'TAG_40001',
    display_name: 'Flow Sensor',
    description: '',
    data_type: 'int16',
    unit: '',
    labels: null,
    status: 'draft',
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

function createMapping(overrides: Partial<Mapping>): Mapping {
  return {
    id: 'mapping-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    enabled: true,
    transform_pipeline: '',
    created_at: '',
    updated_at: '',
    ...overrides,
  };
}

describe('tagBindingModel', () => {
  it('builds address-based tag keys with normalized prefix', () => {
    expect(
      buildTagKey(createPoint({ address: 'D100' }), {
        prefix: 'line a',
        strategy: 'address',
      }),
    ).toBe('LINE-A_D100');
  });

  it('flags existing-key conflicts and already-linked points', () => {
    const candidates = buildTagBindingCandidates({
      points: [createPoint({ id: 'point-1', address: '40001' })],
      tags: [createTag({ key: 'TAG_40001' })],
      mappings: [createMapping({ point_id: 'point-1' })],
      template: {
        prefix: 'tag',
        strategy: 'address',
      },
    });

    expect(candidates[0]).toMatchObject({
      previewKey: 'TAG_40001',
      conflict: true,
      alreadyLinked: true,
    });
  });

  it('builds create requests only for selected, bindable points', () => {
    const candidates = buildTagBindingCandidates({
      points: [
        createPoint({ id: 'point-1', address: '40001' }),
        createPoint({ id: 'point-2', address: '40002', name: 'Pressure' }),
      ],
      tags: [createTag({ key: 'TAG_40002' })],
      mappings: [] as Mapping[],
      template: {
        prefix: 'tag',
        strategy: 'address',
      },
    });

    expect(
      buildTagBindingRequests(candidates, ['point-1', 'point-2']),
    ).toEqual([
      {
        pointId: 'point-1',
        tagKey: 'TAG_40001',
        tagRequest: {
          key: 'TAG_40001',
          data_type: 'int16',
          display_name: 'Flow Sensor',
        },
      },
    ]);
  });
});
