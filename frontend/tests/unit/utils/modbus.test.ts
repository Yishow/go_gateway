import { describe, expect, it } from 'vitest';

import { normalizeRawDataFromPacket } from '@/utils/modbus';

describe('normalizeRawDataFromPacket', () => {
  it('在 raw_data 為陣列時直接回傳', () => {
    const raw = [0x00, 0x01, 0x02];
    const result = normalizeRawDataFromPacket(raw, '00 01 02');
    expect(result).toEqual([0, 1, 2]);
  });

  it('在 raw_data 無法使用時由 hex_data 解析', () => {
    const rawDataBase64 = 'AAECAwQ=';
    const result = normalizeRawDataFromPacket(rawDataBase64, '00 01 02 03 04');
    expect(result).toEqual([0, 1, 2, 3, 4]);
  });

  it('在 hex_data 格式不正確時回傳空陣列', () => {
    const result = normalizeRawDataFromPacket(undefined, 'ZZ ZZ');
    expect(result).toEqual([]);
  });
});
