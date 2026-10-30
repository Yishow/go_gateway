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
        type: 'word',
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
        type: 'word',
      });
    });

    it('should treat D100 same as D0100', () => {
      const result = addressParser.parse('D100', 'fatek_fbs');
      expect(result).toEqual({
        protocol: 'fatek_fbs',
        area: 'D',
        startNumber: 100,
        raw: 'D100', // raw is input
        type: 'word',
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
        type: 'word',
      });
    });

    it('should parse X bit contact address as octal', () => {
      const result = addressParser.parse('X10', 'mc_3e');
      expect(result).toEqual({
        protocol: 'mc_3e',
        area: 'X',
        startNumber: 8,
        raw: 'X10',
        type: 'bit',
      });
    });

    it('should parse MC areas with their protocol radix and type', () => {
      expect(addressParser.parse('Y7', 'mc_3e')).toEqual({
        protocol: 'mc_3e',
        area: 'Y',
        startNumber: 7,
        raw: 'Y7',
        type: 'bit',
      });
      expect(addressParser.parse('WA', 'mc_3e')).toEqual({
        protocol: 'mc_3e',
        area: 'W',
        startNumber: 10,
        raw: 'WA',
        type: 'word',
      });
      expect(addressParser.parse('BFF', 'mc_3e')).toEqual({
        protocol: 'mc_3e',
        area: 'B',
        startNumber: 255,
        raw: 'BFF',
        type: 'bit',
      });
      expect(addressParser.parse('M16', 'mc_3e')).toEqual({
        protocol: 'mc_3e',
        area: 'M',
        startNumber: 16,
        raw: 'M16',
        type: 'bit',
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

    it('should fail explicitly when offset receives an invalid protocol address', () => {
      expect(() => addressParser.offset('Z999', 1, 'mc_3e')).toThrow(/Invalid MC3E address/);
      expect(() => addressParser.offset('', 1, 'mc_3e')).toThrow(/Invalid MC3E address/);
    });

    it('should advance MC3E addresses using area-specific radix', () => {
      expect(addressParser.offset('X0', 16, 'mc_3e')).toBe('X20');
      expect(addressParser.offset('Y7', 1, 'mc_3e')).toBe('Y10');
      expect(addressParser.offset('W0', 10, 'mc_3e')).toBe('WA');
      expect(addressParser.offset('W0', 16, 'mc_3e')).toBe('W10');
      expect(addressParser.offset('WA', 1, 'mc_3e')).toBe('WB');
      expect(addressParser.offset('BFF', 1, 'mc_3e')).toBe('B100');
      expect(addressParser.offset('D0', 16, 'mc_3e')).toBe('D16');
      expect(addressParser.offset('M0', 16, 'mc_3e')).toBe('M16');
    });

    it('should expand MC addresses with the same radix as offset', () => {
      expect(addressParser.expand('W0', 17, 'mc_3e')).toEqual([
        'W0',
        'W1',
        'W2',
        'W3',
        'W4',
        'W5',
        'W6',
        'W7',
        'W8',
        'W9',
        'WA',
        'WB',
        'WC',
        'WD',
        'WE',
        'WF',
        'W10',
      ]);
      expect(addressParser.expand('X0', 3, 'mc_3e')).toEqual(['X0', 'X1', 'X2']);
      expect(addressParser.expand('Y7', 2, 'mc_3e')).toEqual(['Y7', 'Y10']);
    });

    it('should reject MC addresses outside each area radix', () => {
      for (const address of ['X8', 'X9', 'YA', 'WG', 'Z0']) {
        expect(addressParser.validate(address, 'mc_3e').valid).toBe(false);
        expect(addressParser.expand(address, 1, 'mc_3e')).toEqual([]);
        expect(() => addressParser.offset(address, 1, 'mc_3e')).toThrow(
          /Invalid MC3E address/,
        );
      }
    });

    it('should normalize MC address whitespace and case', () => {
      expect(addressParser.parse(' wa ', 'mc_3e')).toEqual({
        protocol: 'mc_3e',
        area: 'W',
        startNumber: 10,
        raw: 'WA',
        type: 'word',
      });
    });

    it('should satisfy the protocol radix boundary spec table', () => {
      expect(addressParser.offset('D0', 16, 'mc_3e')).toBe('D16');
      expect(addressParser.offset('R0', 16, 'fatek_fbs')).toBe('R16');
      expect(addressParser.offset('40001', 2, 'modbus_tcp')).toBe('40003');
    });
  });
});
