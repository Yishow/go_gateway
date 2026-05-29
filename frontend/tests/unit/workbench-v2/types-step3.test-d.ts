import { assertType } from 'vitest';
import type { Mapping, TargetType } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * 測試 Mapping 結構與 TargetType 的靜態斷言
 * 
 * 落地設計決策：「拆檔策略：6 個元件 + 2 個 state module」型別契約
 */
const dummyMapping: Mapping = {
  point_id: 'p-01',
  tag_key: 'line01.temp.inlet',
  display_name: '進水溫度',
  unit: '°C',
  target_type: 'float64',
  scale: 0.1,
  offset: 0,
  enabled: true,
};

assertType<Mapping>(dummyMapping);

// 驗證 target_type 的各種可能值
const validTypes: TargetType[] = [
  'bool',
  'int16',
  'int32',
  'int64',
  'uint16',
  'uint32',
  'uint64',
  'float32',
  'float64',
  'string',
];
assertType<TargetType[]>(validTypes);
