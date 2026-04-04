/**
 * 點位匯入/匯出工具
 * 
 * 支援:
 * - CSV 匯入點位
 * - JSON 匯出設定
 */

import type { Point, CreatePointRequest } from '../types/datalink';
import { DATALINK_DATA_TYPES } from '../types/datalink';

// =============================================================================
// CSV 匯入
// =============================================================================

export interface CSVParseResult {
  success: boolean;
  points: CreatePointRequest[];
  errors: CSVParseError[];
  warnings: string[];
}

export interface CSVParseError {
  row: number;
  column?: string;
  message: string;
}

/**
 * 解析 CSV 內容為點位請求
 */
export function parseCSV(
  csvContent: string,
  deviceId: string,
  options: { delimiter?: string; hasHeader?: boolean } = {}
): CSVParseResult {
  const { delimiter = ',', hasHeader = true } = options;
  const result: CSVParseResult = {
    success: true,
    points: [],
    errors: [],
    warnings: [],
  };

  const lines = csvContent.trim().split('\n').map(line => line.trim()).filter(line => line);
  
  if (lines.length === 0) {
    result.success = false;
    result.errors.push({ row: 0, message: 'CSV 檔案為空' });
    return result;
  }

  // 解析標題行
  let headers: string[] = ['name', 'address', 'data_type', 'description', 'polling_group_id'];
  let startRow = 0;

  if (hasHeader) {
    headers = lines[0].split(delimiter).map(h => h.trim().toLowerCase().replace(/\s+/g, '_'));
    startRow = 1;
  }

  // 驗證必要欄位
  const requiredFields = ['name', 'address'];
  for (const field of requiredFields) {
    if (!headers.includes(field)) {
      result.success = false;
      result.errors.push({ row: 0, column: field, message: `缺少必要欄位: ${field}` });
    }
  }

  if (!result.success) {
    return result;
  }

  // 解析資料行
  for (let i = startRow; i < lines.length; i++) {
    const rowNum = i + 1;
    const values = parseCSVLine(lines[i], delimiter);
    
    if (values.length < headers.length) {
      result.warnings.push(`第 ${rowNum} 行欄位數不足，已補空值`);
    }

    const row: Record<string, string> = {};
    headers.forEach((header, index) => {
      row[header] = values[index]?.trim() || '';
    });

    // 驗證必要欄位值
    if (!row.name) {
      result.errors.push({ row: rowNum, column: 'name', message: '名稱不能為空' });
      continue;
    }
    if (!row.address) {
      result.errors.push({ row: rowNum, column: 'address', message: '位址不能為空' });
      continue;
    }

    // 驗證資料類型（與 DATALINK_DATA_TYPES 一致）
    const dataType = row.data_type || 'uint16';
    if (!(DATALINK_DATA_TYPES as readonly string[]).includes(dataType)) {
      result.errors.push({ row: rowNum, column: 'data_type', message: `無效的資料類型: ${dataType}` });
      continue;
    }

    result.points.push({
      device_id: deviceId,
      name: row.name,
      address: row.address,
      data_type: dataType as CreatePointRequest['data_type'],
      description: row.description || '',
      polling_group_id: row.polling_group_id || undefined,
      enabled: row.enabled !== 'false',
    });
  }

  if (result.errors.length > 0) {
    result.success = false;
  }

  return result;
}

/**
 * 解析 CSV 行（處理引號內的逗號）
 */
function parseCSVLine(line: string, delimiter: string): string[] {
  const result: string[] = [];
  let current = '';
  let inQuotes = false;

  for (let i = 0; i < line.length; i++) {
    const char = line[i];
    
    if (char === '"') {
      if (inQuotes && line[i + 1] === '"') {
        current += '"';
        i++;
      } else {
        inQuotes = !inQuotes;
      }
    } else if (char === delimiter && !inQuotes) {
      result.push(current);
      current = '';
    } else {
      current += char;
    }
  }
  
  result.push(current);
  return result;
}

/**
 * 產生 CSV 範本
 */
export function generateCSVTemplate(): string {
  const headers = ['name', 'address', 'data_type', 'description', 'polling_group_id', 'enabled'];
  const example = ['Pump_Status', 'D100', 'uint16', '幫浦狀態', '', 'true'];
  
  return [headers.join(','), example.join(',')].join('\n');
}

// =============================================================================
// JSON 匯出
// =============================================================================

export interface ExportConfig {
  points: Point[];
  exportedAt: string;
  version: string;
  deviceId?: string;
  deviceName?: string;
}

/**
 * 匯出點位為 JSON
 */
export function exportToJSON(
  points: Point[],
  options: { deviceId?: string; deviceName?: string } = {}
): string {
  const config: ExportConfig = {
    points,
    exportedAt: new Date().toISOString(),
    version: '1.0.0',
    deviceId: options.deviceId,
    deviceName: options.deviceName,
  };

  return JSON.stringify(config, null, 2);
}

/**
 * 匯出點位為 CSV
 */
export function exportToCSV(points: Point[]): string {
  const headers = ['name', 'address', 'data_type', 'description', 'polling_group_id', 'enabled'];
  
  const rows = points.map(point => {
    return [
      escapeCSVValue(point.name),
      escapeCSVValue(point.address),
      point.data_type,
      escapeCSVValue(point.description || ''),
      point.polling_group_id || '',
      point.enabled ? 'true' : 'false',
    ].join(',');
  });

  return [headers.join(','), ...rows].join('\n');
}

/**
 * 轉義 CSV 值
 */
function escapeCSVValue(value: string): string {
  if (value.includes(',') || value.includes('"') || value.includes('\n')) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

// =============================================================================
// 檔案操作
// =============================================================================

/**
 * 下載檔案
 */
export function downloadFile(content: string, filename: string, mimeType: string): void {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  URL.revokeObjectURL(url);
}

/**
 * 讀取檔案內容
 */
export function readFile(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = () => reject(new Error('讀取檔案失敗'));
    reader.readAsText(file);
  });
}
