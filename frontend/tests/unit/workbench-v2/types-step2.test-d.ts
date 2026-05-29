import { assertType } from 'vitest';
import type { Rule, Point, ShareLayout } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * 測試 Rule 結構的靜態斷言
 */
const dummyRule: Rule = {
  id: 'rule-01',
  device_id: 'dev-01',
  name: 'Telemetry Registers',
  start_address: '40001',
  count: 10,
  data_type: 'float32',
  naming_prefix: 'TEMP_',
  enabled: true,
  scale_multiplier: 0.1,
  scale_offset: 10.0,
  data_format: 'ABCD',
  skipped_addresses: ['40003', '40004'],
  share_enabled: true,
  share_start_register: 40010,
  share_stride: 2,
};
assertType<Rule>(dummyRule);

/**
 * 測試 Point 結構的靜態斷言
 */
const dummyPoint: Point = {
  id: 'p-01',
  device_id: 'dev-01',
  rule_id: 'rule-01',
  rule_name: 'Telemetry Registers',
  name: 'TEMP_0',
  address: '40001',
  data_type: 'float32',
  function: 'holding_register',
  width: 2,
  enabled: true,
  skipped: false,
  display: 'Temperature 1',
  unit: '°C',
  tag_key_suggest: 'temp.1',
  _rule_scale: 0.1,
  _rule_offset: 10.0,
};
assertType<Point>(dummyPoint);

/**
 * 測試 ShareLayout 結構的靜態斷言
 */
const dummyShareLayout: ShareLayout = {
  start: 40010,
  stride: 2,
  end: 40030,
  auto: false,
};
assertType<ShareLayout>(dummyShareLayout);
