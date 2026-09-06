import { describe, it, expect } from 'vitest';
import { getColumnsFor, getDefaultConnector, getDbKindPatch } from '../../../src/features/datalink/workbench-v2/state/dbSchemas';

/**
 * @file dbSchemas.test.ts
 * @description 測試資料庫 Schema 與預設連接器產生邏輯，覆蓋 4 種資料庫類型。
 */

describe('dbSchemas', () => {
  describe('getColumnsFor', () => {
    it('應回傳 PostgreSQL 的欄位定義且順序正確，第一欄為主鍵 ts', () => {
      const columns = getColumnsFor('postgres');
      expect(columns).toHaveLength(9);
      expect(columns[0]).toEqual({
        name: 'ts',
        type: 'timestamptz',
        nullable: false,
        primary_key: true
      });
      expect(columns[1].name).toBe('temp_in_c');
      expect(columns[8].name).toBe('motor_rpm');
    });

    it('應支援 sqlite 的型別轉換', () => {
      const columns = getColumnsFor('sqlite');
      expect(columns[0].type).toBe('DATETIME');
      expect(columns[1].type).toBe('REAL');
      expect(columns[8].type).toBe('INTEGER');
    });

    it('應支援 mysql 的型別轉換', () => {
      const columns = getColumnsFor('mysql');
      expect(columns[0].type).toBe('TIMESTAMP');
      expect(columns[1].type).toBe('DOUBLE');
      expect(columns[8].type).toBe('INT');
    });

    it('應支援 sqlserver 的型別轉換', () => {
      const columns = getColumnsFor('sqlserver');
      expect(columns[0].type).toBe('DATETIMEOFFSET');
      expect(columns[1].type).toBe('FLOAT');
      expect(columns[8].type).toBe('INT');
    });
  });

  describe('getDefaultConnector', () => {
    it('應取得 postgres 的預設連線設定', () => {
      const connector = getDefaultConnector('postgres');
      expect(connector.kind).toBe('postgres');
      expect(connector.port).toBe(5432);
      expect(connector.host).toBe('127.0.0.1');
      expect(connector.database).toBe('gateway_metrics');
      expect(connector.schema).toBe('public');
      expect(connector.table).toBe('sensor_readings');
      expect(connector.write_mode).toBe('insert');
      expect(connector.write_interval_seconds).toBe(5);
    });

    it('應取得 sqlite 的預設連線設定，且 schema 為空字串', () => {
      const connector = getDefaultConnector('sqlite');
      expect(connector.kind).toBe('sqlite');
      expect(connector.host).toBe('');
      expect(connector.port).toBe(0);
      expect(connector.database).toBe('gateway.db');
      expect(connector.schema).toBe('');
    });

    it('應取得 mysql 的預設連線設定', () => {
      const connector = getDefaultConnector('mysql');
      expect(connector.kind).toBe('mysql');
      expect(connector.port).toBe(3306);
      expect(connector.host).toBe('127.0.0.1');
      expect(connector.schema).toBe('');
    });

    it('應取得 sqlserver 的預設連線設定，且 schema 為 dbo', () => {
      const connector = getDefaultConnector('sqlserver');
      expect(connector.kind).toBe('sqlserver');
      expect(connector.port).toBe(1433);
      expect(connector.host).toBe('127.0.0.1');
      expect(connector.schema).toBe('dbo');
    });
  });

  describe('getDbKindPatch', () => {
    it('從 postgres 切換至 mysql 時應更新 port 為 3306、username 為 root、schema 清空，並維持 loopback host', () => {
      const current = {
        kind: 'postgres' as const,
        host: '127.0.0.1',
        port: 5432,
        username: 'postgres',
        schema: 'public',
      };
      const patch = getDbKindPatch('mysql', current);
      expect(patch).toMatchObject({
        kind: 'mysql',
        host: '127.0.0.1',
        port: 3306,
        username: 'root',
        schema: '',
      });
    });

    it('從 postgres 切換至 sqlserver 時若為自訂 remote host 應予以保留，且 port 改為 1433、schema 改為 dbo', () => {
      const current = {
        kind: 'postgres' as const,
        host: '192.168.10.50',
        port: 5432,
        username: 'custom_user',
        schema: 'public',
      };
      const patch = getDbKindPatch('sqlserver', current);
      expect(patch).toMatchObject({
        kind: 'sqlserver',
        host: '192.168.10.50',
        port: 1433,
        username: 'custom_user',
        schema: 'dbo',
      });
    });

    it('切換至 sqlite 時應清空 host、port、username、password、schema', () => {
      const current = {
        kind: 'postgres' as const,
        host: '127.0.0.1',
        port: 5432,
        username: 'postgres',
        password: 'secret_password',
        schema: 'public',
        database: 'gateway_metrics',
      };
      const patch = getDbKindPatch('sqlite', current);
      expect(patch).toMatchObject({
        kind: 'sqlite',
        host: '',
        port: 0,
        username: '',
        schema: '',
        database: 'gateway.db',
      });
      expect(patch.password).toBeUndefined();

    });
  });
});

