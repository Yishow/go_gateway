import { describe, it, expect } from 'vitest';
import { runScale, castValue, formatFinal, buildPayload } from '../../../src/features/datalink/workbench-v2/state/transformPipeline';
import type { Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

describe('transformPipeline', () => {
  describe('runScale', () => {
    it('should scale values correctly', () => {
      expect(runScale(243, 0.1, 0)).toBeCloseTo(24.3);
      expect(runScale(100, 1.5, -10)).toBe(140);
    });
  });

  describe('castValue', () => {
    it('should cast bool correctly', () => {
      expect(castValue(0, 'bool')).toBe(false);
      expect(castValue(1, 'bool')).toBe(true);
      expect(castValue(-5, 'bool')).toBe(true);
    });

    it('should cast int16 with clamping', () => {
      expect(castValue(24.3, 'int16')).toBe(24);
      expect(castValue(24.7, 'int16')).toBe(25);
      expect(castValue(40000, 'int16')).toBe(32767);
      expect(castValue(-40000, 'int16')).toBe(-32768);
    });

    it('should cast uint16 with clamping', () => {
      expect(castValue(100.2, 'uint16')).toBe(100);
      expect(castValue(70000, 'uint16')).toBe(65535);
      expect(castValue(-5, 'uint16')).toBe(0);
    });

    it('should cast float32 and float64 without clamping', () => {
      expect(castValue(123456.78, 'float32')).toBe(123456.78);
      expect(castValue(-123.456, 'float64')).toBe(-123.456);
    });

    it('should cast to string', () => {
      expect(castValue(12.34, 'string')).toBe('12.34');
    });
  });

  describe('formatFinal', () => {
    it('should format bool values to string', () => {
      expect(formatFinal(true, 'bool')).toBe('true');
      expect(formatFinal(false, 'bool')).toBe('false');
    });

    it('should format float values with 2 decimals', () => {
      expect(formatFinal(24.3456, 'float64')).toBe('24.35');
      expect(formatFinal(12, 'float32')).toBe('12.00');
    });

    it('should format int values to string', () => {
      expect(formatFinal(24.3, 'int16')).toBe('24');
      expect(formatFinal(25, 'uint32')).toBe('25');
    });

    it('should format string values directly', () => {
      expect(formatFinal('hello', 'string')).toBe('hello');
    });
  });

  describe('buildPayload', () => {
    it('should generate valid mapping payload', () => {
      const point: Point = {
        id: 'p-01',
        device_id: 'dev-01',
        rule_id: 'rule-01',
        rule_name: 'Holding Registers',
        name: 'SENSOR_1',
        address: '40001',
        data_type: 'int16',
        function: 'holding_register',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 0.1,
        _rule_offset: 0,
      };

      const mapping: Mapping = {
        point_id: 'p-01',
        tag_key: 'line01.temp.inlet',
        display_name: '進水溫度',
        unit: '°C',
        target_type: 'float64',
        scale: 0.1,
        offset: 0,
        enabled: true,
      };

      const payload = buildPayload(mapping, point) as {
        point_id: string;
        tag_id: string;
        enabled: boolean;
        transform_pipeline: Array<{ type: string; order: number; params: Record<string, unknown> }>;
      };

      expect(payload.point_id).toBe('p-01');
      expect(payload.tag_id).toBe('line01.temp.inlet');
      expect(payload.enabled).toBe(true);
      expect(payload.transform_pipeline).toHaveLength(3);

      expect(payload.transform_pipeline[0]).toEqual({
        type: 'decode',
        order: 1,
        params: { data_type: 'int16' },
      });

      expect(payload.transform_pipeline[1]).toEqual({
        type: 'scale',
        order: 2,
        params: { scale: 0.1, offset: 0 },
      });

      expect(payload.transform_pipeline[2]).toEqual({
        type: 'cast',
        order: 3,
        params: { target_type: 'float64' },
      });
    });
  });
});
