import { describe, it, expect } from 'vitest';
import { addressParser, getDefaultPlannerStartAddress } from '@/utils/addressParser';

describe('AddressParser', () => {
  describe('Modbus', () => {
    it('should parse holding register address (4xxxx)', () => {
      const result = addressParser.parse('40001', 'modbus_tcp');
      expect(result).toEqual({
        protocol: 'modbus_tcp',
        area: 'HR',
        startNumber: 1,
        raw: '40001',
      });
    });

    it('should expand address range', () => {
      const result = addressParser.expand('40001', 3, 'modbus_tcp');
      expect(result).toEqual(['40001', '40002', '40003']);
    });
  });

  describe('FATEK', () => {
    it('should parse D register address', () => {
      const result = addressParser.parse('D0100', 'fatek_fbs');
      expect(result).toEqual({
        protocol: 'fatek_fbs',
        area: 'D',
        startNumber: 100,
        raw: 'D0100',
      });
    });

    it('should treat D100 same as D0100', () => {
      const result = addressParser.parse('D100', 'fatek_fbs');
      expect(result).toEqual({
        protocol: 'fatek_fbs',
        area: 'D',
        startNumber: 100,
        raw: 'D100', // raw is input
      });
    });
  });

  describe('MC3E', () => {
    it('should parse D register address', () => {
      const result = addressParser.parse('D100', 'mc_3e');
      expect(result).toEqual({
        protocol: 'mc_3e',
        area: 'D',
        startNumber: 100,
        raw: 'D100',
      });
    });
  });

  describe('Planner defaults', () => {
    it('should return protocol-specific default start addresses', () => {
      expect(getDefaultPlannerStartAddress('modbus_tcp')).toBe('40001');
      expect(getDefaultPlannerStartAddress('modbus_rtu')).toBe('40001');
      expect(getDefaultPlannerStartAddress('modbus_udp')).toBe('40001');
      expect(getDefaultPlannerStartAddress('fatek_fbs')).toBe('D0');
      expect(getDefaultPlannerStartAddress('mc_3e')).toBe('D0');
    });
  });

  describe('Validation', () => {
    it('should return valid for correct address', () => {
      expect(addressParser.validate('40001', 'modbus_tcp').valid).toBe(true);
    });

    it('should return invalid for incorrect format', () => {
      expect(addressParser.validate('INVALID', 'modbus_tcp').valid).toBe(false);
    });
  });
});
