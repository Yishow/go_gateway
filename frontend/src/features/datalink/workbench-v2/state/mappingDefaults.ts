import type { Mapping } from './types';
import type { Point } from './types';

export const POINT_SEMANTIC = [
  { suffix: 'TEMP_IN', display: '入口溫度', unit: '°C', tag_key: 'line01.temp.inlet' },
  { suffix: 'TEMP_OUT', display: '出口溫度', unit: '°C', tag_key: 'line01.temp.outlet' },
  { suffix: 'PRES_01', display: '主管路壓力', unit: 'kPa', tag_key: 'line01.pressure.main' },
  { suffix: 'PRES_02', display: '次管路壓力', unit: 'kPa', tag_key: 'line01.pressure.sub' },
  { suffix: 'FLOW_01', display: '流量計', unit: 'L/min', tag_key: 'line01.flow.q1' },
  { suffix: 'HUM_01', display: '濕度', unit: '%', tag_key: 'line01.humidity.amb' },
  { suffix: 'VIB_01', display: '振動', unit: 'mm/s', tag_key: 'line01.vibration.motor' },
  { suffix: 'RPM_01', display: '馬達轉速', unit: 'rpm', tag_key: 'line01.motor.rpm' },
];

export const RAW_VALUE_SEEDS = [243, 251, 1024, 985, 67, 542, 18, 1450];

function slugifySegment(value: string): string {
  const normalized = value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '.')
    .replace(/^\.+|\.+$/g, '')
    .replace(/\.+/g, '.');

  return normalized || 'point';
}

function buildDefaultTagKey(point: Point): string {
  return [
    slugifySegment(point.device_id),
    slugifySegment(point.name),
    `r${point.address.replace(/[^0-9a-z]/gi, '').toLowerCase()}`,
  ].join('.');
}

/**
 * 根據點位與索引建立預設的 Mapping 物件
 * 
 * @param point 點位資料
 * @param idx 索引
 * @returns 預設 Mapping 物件
 */
export function buildDefaultMapping(point: Point, _idx: number): Mapping {
  const value = {
    tag_key: buildDefaultTagKey(point),
    display_name: point.display?.trim() || point.name,
    unit: point.unit?.trim() || '',
    target_type: 'float64' as const,
    scale: point._rule_scale,
    offset: point._rule_offset,
    enabled: true,
  };
  return {
    point_id: point.id,
    ...value,
    persisted: false,
    local_value: value,
    save_state: 'idle',
    save_error: null,
  };
}
