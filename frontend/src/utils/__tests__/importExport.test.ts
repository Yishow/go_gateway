import { describe, it, expect } from 'vitest';
import {
  parseCSV,
  generateCSVTemplate,
  exportToJSON,
  exportToCSV,
} from '../importExport';
import type { Point } from '../../types/datalink';

describe('CSV Import', () => {
  describe('parseCSV', () => {
    it('should parse valid CSV with header', () => {
      const csv = `name,address,data_type,description
Pump1,D100,uint16,幫浦1
Pump2,D101,uint16,幫浦2`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.points).toHaveLength(2);
      expect(result.points[0]).toEqual({
        device_id: 'device-1',
        name: 'Pump1',
        address: 'D100',
        data_type: 'uint16',
        description: '幫浦1',
        polling_group_id: undefined,
        enabled: true,
      });
    });

    it('should handle missing optional fields', () => {
      const csv = `name,address
Sensor1,D200
Sensor2,D201`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.points).toHaveLength(2);
      expect(result.points[0].data_type).toBe('uint16'); // default
      expect(result.points[0].description).toBe('');
    });

    it('should return error for empty CSV', () => {
      const result = parseCSV('', 'device-1');
      
      expect(result.success).toBe(false);
      expect(result.errors[0].message).toContain('為空');
    });

    it('should return error for missing required field name', () => {
      const csv = `address,data_type
D100,uint16`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(false);
      expect(result.errors.some(e => e.column === 'name')).toBe(true);
    });

    it('should return error for missing required field address', () => {
      const csv = `name,data_type
Pump1,uint16`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(false);
      expect(result.errors.some(e => e.column === 'address')).toBe(true);
    });

    it('should return error for empty name value', () => {
      const csv = `name,address,data_type
,D100,uint16`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(false);
      expect(result.errors.some(e => e.column === 'name' && e.message.includes('不能為空'))).toBe(true);
    });

    it('should return error for invalid data_type', () => {
      const csv = `name,address,data_type
Pump1,D100,invalid_type`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(false);
      expect(result.errors.some(e => e.column === 'data_type')).toBe(true);
    });

    it('should handle quoted values with commas', () => {
      const csv = `name,address,description
"Pump, Main",D100,"Description, with comma"`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.points[0].name).toBe('Pump, Main');
      expect(result.points[0].description).toBe('Description, with comma');
    });

    it('should handle quoted values with escaped quotes', () => {
      const csv = `name,address,description
"Pump ""A""",D100,Test`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.points[0].name).toBe('Pump "A"');
    });

    it('should use custom delimiter', () => {
      const csv = `name;address;data_type
Pump1;D100;uint16`;
      
      const result = parseCSV(csv, 'device-1', { delimiter: ';' });
      
      expect(result.success).toBe(true);
      expect(result.points[0].name).toBe('Pump1');
    });

    it('should handle CSV without header', () => {
      const csv = `Pump1,D100,uint16,Description`;
      
      const result = parseCSV(csv, 'device-1', { hasHeader: false });
      
      expect(result.success).toBe(true);
      expect(result.points[0].name).toBe('Pump1');
      expect(result.points[0].address).toBe('D100');
    });

    it('should add warning for insufficient columns', () => {
      const csv = `name,address,data_type,description
Pump1,D100`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.warnings.length).toBeGreaterThan(0);
    });

    it('should handle enabled field', () => {
      const csv = `name,address,enabled
Pump1,D100,true
Pump2,D101,false`;
      
      const result = parseCSV(csv, 'device-1');
      
      expect(result.success).toBe(true);
      expect(result.points[0].enabled).toBe(true);
      expect(result.points[1].enabled).toBe(false);
    });

    it('should handle all valid data types', () => {
      const dataTypes = ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'];
      
      for (const dt of dataTypes) {
        const csv = `name,address,data_type
Test,D100,${dt}`;
        
        const result = parseCSV(csv, 'device-1');
        expect(result.success).toBe(true);
        expect(result.points[0].data_type).toBe(dt);
      }
    });
  });

  describe('generateCSVTemplate', () => {
    it('should generate valid template', () => {
      const template = generateCSVTemplate();
      
      expect(template).toContain('name');
      expect(template).toContain('address');
      expect(template).toContain('data_type');
      expect(template.split('\n')).toHaveLength(2);
    });
  });
});

describe('JSON Export', () => {
  const mockPoints: Point[] = [
    {
      id: 'p1',
      device_id: 'd1',
      name: 'Pump1',
      description: '幫浦1',
      data_type: 'uint16',
      address: 'D100',
      enabled: true,
      polling_group_id: 'pg1',
      last_value: 100,
      last_read_at: '2024-01-01T00:00:00Z',
      last_error: '',
      error_count: 0,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
  ];

  describe('exportToJSON', () => {
    it('should export points to valid JSON', () => {
      const json = exportToJSON(mockPoints);
      const parsed = JSON.parse(json);
      
      expect(parsed.points).toHaveLength(1);
      expect(parsed.version).toBe('1.0.0');
      expect(parsed.exportedAt).toBeTruthy();
    });

    it('should include device info when provided', () => {
      const json = exportToJSON(mockPoints, { deviceId: 'd1', deviceName: 'PLC-001' });
      const parsed = JSON.parse(json);
      
      expect(parsed.deviceId).toBe('d1');
      expect(parsed.deviceName).toBe('PLC-001');
    });

    it('should be pretty printed', () => {
      const json = exportToJSON(mockPoints);
      expect(json).toContain('\n');
      expect(json).toContain('  ');
    });
  });

  describe('exportToCSV', () => {
    it('should export points to valid CSV', () => {
      const csv = exportToCSV(mockPoints);
      const lines = csv.split('\n');
      
      expect(lines[0]).toBe('name,address,data_type,description,polling_group_id,enabled');
      expect(lines[1]).toContain('Pump1');
      expect(lines[1]).toContain('D100');
    });

    it('should escape values with commas', () => {
      const pointsWithComma: Point[] = [{
        ...mockPoints[0],
        name: 'Pump, Main',
        description: 'Has, comma',
      }];
      
      const csv = exportToCSV(pointsWithComma);
      
      expect(csv).toContain('"Pump, Main"');
      expect(csv).toContain('"Has, comma"');
    });

    it('should escape values with quotes', () => {
      const pointsWithQuote: Point[] = [{
        ...mockPoints[0],
        name: 'Pump "A"',
      }];
      
      const csv = exportToCSV(pointsWithQuote);
      
      expect(csv).toContain('"Pump ""A"""');
    });

    it('should handle empty description', () => {
      const pointsNoDesc: Point[] = [{
        ...mockPoints[0],
        description: '',
      }];
      
      const csv = exportToCSV(pointsNoDesc);
      const lines = csv.split('\n');
      
      expect(lines[1]).toContain('Pump1,D100,uint16,,');
    });

    it('should handle boolean enabled field', () => {
      const csv = exportToCSV(mockPoints);
      
      expect(csv).toContain('true');
    });

    it('should export disabled points correctly', () => {
      const disabledPoint: Point[] = [{
        ...mockPoints[0],
        enabled: false,
      }];
      
      const csv = exportToCSV(disabledPoint);
      
      expect(csv).toContain('false');
    });
  });
});

describe('Round-trip Import/Export', () => {
  it('should preserve data through export and re-import', () => {
    const mockPoints: Point[] = [
      {
        id: 'p1',
        device_id: 'd1',
        name: 'TestPoint',
        description: 'Test description',
        data_type: 'float32',
        address: 'D500',
        enabled: true,
        polling_group_id: 'pg1',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    ];

    // Export to CSV
    const csv = exportToCSV(mockPoints);
    
    // Re-import
    const result = parseCSV(csv, 'd1');
    
    expect(result.success).toBe(true);
    expect(result.points[0].name).toBe('TestPoint');
    expect(result.points[0].address).toBe('D500');
    expect(result.points[0].data_type).toBe('float32');
    expect(result.points[0].description).toBe('Test description');
  });
});
