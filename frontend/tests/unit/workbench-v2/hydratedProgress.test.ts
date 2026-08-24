import { describe, expect, it } from 'vitest';
import { inferHydratedProgress } from '../../../src/pages/datalink/workbench-v2/hydratedProgress';
import type { Device, Mapping, Rule } from '../../../src/features/datalink/workbench-v2/state/types';

const device: Device = {
  id: 'dev-valid',
  name: 'Valid PLC',
  description: '',
  protocol: 'modbus_tcp',
  config: {},
  status: 'active',
  test: null,
};

const rule: Rule = {
  id: 'rule-valid',
  device_id: 'dev-valid',
  name: 'Valid Rule',
  start_address: '40001',
  count: 1,
  data_type: 'int16',
  naming_prefix: 'TAG_',
  enabled: true,
  scale_multiplier: 1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: [],
  share_enabled: false,
  share_start_register: null,
  share_stride: null,
  persisted: true,
};

const mapping: Mapping = {
  point_id: 'rule-valid-p-0',
  tag_key: 'line.valid',
  display_name: 'Valid tag',
  unit: '',
  target_type: 'int16',
  scale: 1,
  offset: 0,
  enabled: true,
  persisted: true,
  tag_id: 'tag-valid',
};

describe('Studio V2 hydrated progress', () => {
  it('keeps an orphan persisted rule at Step 2 without completing Step 2 or 3', () => {
    const orphanRule = { ...rule, device_id: 'deleted-device' };
    const progress = inferHydratedProgress(
      [device],
      [orphanRule],
      { [mapping.point_id]: mapping },
    );

    expect(progress.current).toBe(2);
    expect(progress.completed).toEqual(new Set([1]));
  });

  it('allows valid hydrated Step 2 and persisted mapping to advance to Step 4', () => {
    const progress = inferHydratedProgress([device], [rule], { [mapping.point_id]: mapping });

    expect(progress.current).toBe(4);
    expect(progress.completed).toEqual(new Set([1, 2, 3]));
  });
});
