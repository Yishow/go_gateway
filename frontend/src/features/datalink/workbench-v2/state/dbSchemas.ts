import type { DbConnector } from './types';

/**
 * 資料庫欄位定義介面
 */
export interface DbColumn {
  name: string;
  type: string;
  nullable: boolean;
  primary_key: boolean;
}

/**
 * 模擬的資料庫欄位資料（PostgreSQL 9 欄範例）
 * 落地設計決策：「Sample table schema」
 */
export const SAMPLE_DB_TABLES: DbColumn[] = [
  { name: 'ts', type: 'timestamptz', nullable: false, primary_key: true },
  { name: 'temp_in_c', type: 'double precision', nullable: true, primary_key: false },
  { name: 'temp_out_c', type: 'double precision', nullable: true, primary_key: false },
  { name: 'pressure_main_kpa', type: 'double precision', nullable: true, primary_key: false },
  { name: 'pressure_sub_kpa', type: 'double precision', nullable: true, primary_key: false },
  { name: 'flow_lpm', type: 'double precision', nullable: true, primary_key: false },
  { name: 'humidity_pct', type: 'double precision', nullable: true, primary_key: false },
  { name: 'vibration_mms', type: 'double precision', nullable: true, primary_key: false },
  { name: 'motor_rpm', type: 'integer', nullable: true, primary_key: false }
];

/**
 * 依據 connector 類型取得資料表欄位清單
 * @param kind 資料庫類型
 * @returns 欄位陣列
 */
export function getColumnsFor(kind: DbConnector['kind']): DbColumn[] {
  // 目前所有 kind 都共用同一份模擬 schema，後續可依 kind 調整格式
  switch (kind) {
    case 'sqlite':
      return SAMPLE_DB_TABLES.map(c => ({
        ...c,
        type: c.type === 'timestamptz' ? 'DATETIME' : c.type === 'double precision' ? 'REAL' : 'INTEGER'
      }));
    case 'mysql':
      return SAMPLE_DB_TABLES.map(c => ({
        ...c,
        type: c.type === 'timestamptz' ? 'TIMESTAMP' : c.type === 'double precision' ? 'DOUBLE' : 'INT'
      }));
    case 'sqlserver':
      return SAMPLE_DB_TABLES.map(c => ({
        ...c,
        type: c.type === 'timestamptz' ? 'DATETIMEOFFSET' : c.type === 'double precision' ? 'FLOAT' : 'INT'
      }));
    case 'postgres':
    default:
      return SAMPLE_DB_TABLES;
  }
}

/**
 * 依據 connector 類型取得預設的連線設定
 * @param kind 資料庫類型
 * @returns 預設的 DbConnector 設定
 */
export function getDefaultConnector(kind: DbConnector['kind']): DbConnector {
  switch (kind) {
    case 'sqlite':
      return {
        kind: 'sqlite',
        name: 'SQLite Connector',
        host: '',
        port: 0,
        database: 'gateway.db',
        username: '',
        schema: '',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      };
    case 'mysql':
      return {
        kind: 'mysql',
        name: 'MySQL Connector',
        host: 'localhost',
        port: 3306,
        database: 'gateway_metrics',
        username: 'root',
        schema: '',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      };
    case 'sqlserver':
      return {
        kind: 'sqlserver',
        name: 'SQL Server Connector',
        host: 'localhost',
        port: 1433,
        database: 'gateway_metrics',
        username: 'sa',
        schema: 'dbo',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      };
    case 'postgres':
    default:
      return {
        kind: 'postgres',
        name: 'PostgreSQL Connector',
        host: 'tsdb.internal',
        port: 5432,
        database: 'gateway_metrics',
        username: 'postgres',
        schema: 'public',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      };
  }
}
