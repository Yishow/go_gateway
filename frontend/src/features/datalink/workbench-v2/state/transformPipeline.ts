import type { Mapping } from './types';
import type { Point } from './types';

/**
 * 執行線性轉換 (× Scale + Offset)
 * 
 * @param raw 原始解碼數值
 * @param scale 乘積係數
 * @param offset 偏移量
 * @returns 轉換後的數值
 */
export function runScale(raw: number, scale: number, offset: number): number {
  return raw * scale + offset;
}

/**
 * 轉換數值型別到目標型別 (Clamping 防止溢出)
 * 
 * @param scaled 縮放後的數值
 * @param targetType 目標資料型別
 * @returns 轉換後的數值或布林值或字串
 */
export function castValue(scaled: number, targetType: Mapping['target_type']): number | boolean | string {
  if (targetType === 'bool') {
    return scaled !== 0;
  }
  if (targetType === 'string') {
    return String(scaled);
  }

  // 取得四捨五入的整數值，供整數 Clamping 使用
  const rounded = Math.round(scaled);

  switch (targetType) {
    case 'int16':
      return Math.max(-32768, Math.min(32767, rounded));
    case 'int32':
      return Math.max(-2147483648, Math.min(2147483647, rounded));
    case 'int64':
      // JS 雙精度浮點數限制，使用 MIN_SAFE_INTEGER 和 MAX_SAFE_INTEGER
      return Math.max(Number.MIN_SAFE_INTEGER, Math.min(Number.MAX_SAFE_INTEGER, rounded));
    case 'uint16':
      return Math.max(0, Math.min(65535, rounded));
    case 'uint32':
      return Math.max(0, Math.min(4294967295, rounded));
    case 'uint64':
      return Math.max(0, Math.min(Number.MAX_SAFE_INTEGER, rounded));
    case 'float32':
    case 'float64':
    default:
      return scaled;
  }
}

/**
 * 格式化最終輸出值，供 UI 預覽顯示
 * 
 * @param value 已經轉型後的數值/布林/字串
 * @param targetType 目標資料型別
 * @returns 格式化後的字串
 */
export function formatFinal(value: number | boolean | string, targetType: Mapping['target_type']): string {
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false';
  }
  if (targetType === 'bool') {
    return value ? 'true' : 'false';
  }
  if (targetType === 'string') {
    return String(value);
  }
  if (typeof value === 'number') {
    if (targetType === 'float32' || targetType === 'float64') {
      return value.toFixed(2);
    }
    return Math.round(value).toString();
  }
  return String(value);
}

/**
 * 建立點位映射的 API 請求 JSON 結構
 * 
 * @param mapping 映射設定
 * @param point 原始點位
 * @returns 映射 API Payload 物件
 */
export function buildPayload(mapping: Mapping, point: Point): object {
  return {
    point_id: mapping.point_id,
    tag_id: mapping.tag_key,
    transform_pipeline: [
      {
        type: 'decode',
        order: 1,
        params: { data_type: point.data_type },
      },
      {
        type: 'scale',
        order: 2,
        params: { scale: mapping.scale, offset: mapping.offset },
      },
      {
        type: 'cast',
        order: 3,
        params: { target_type: mapping.target_type },
      },
    ],
    enabled: mapping.enabled,
  };
}
