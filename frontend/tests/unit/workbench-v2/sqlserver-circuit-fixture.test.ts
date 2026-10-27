import fs from 'node:fs';
import path from 'node:path';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import {
  buildSQLServerRowSelectSql,
  canonicalizeDecimalText,
  canonicalizeUtcTimeText,
  loadSQLServerFixture,
  parseWireRowPayload,
  querySQLServerReadback,
  SQLSERVER_FIXTURE_ENV,
  sqlserverIdentifier,
  validateConnectionConfigFile,
  validateFixtureManifest,
  verifyWirePayloadAgainstRow,
  type SQLServerConnectionConfig,
  type SQLServerFixtureManifest,
  type WireRowPayload,
} from '../../e2e/fixtures/sqlserver-circuit-recording';

describe('SQL Server browser acceptance fixture contract and guards', () => {
  const tmpDir = path.join('/private/tmp/circuit-sql-work', 'fixture-test-' + Date.now());

  beforeEach(() => {
    fs.mkdirSync(tmpDir, { recursive: true });
    delete process.env[SQLSERVER_FIXTURE_ENV];
  });

  afterEach(() => {
    delete process.env[SQLSERVER_FIXTURE_ENV];
    fs.rmSync(tmpDir, { recursive: true, force: true });
  });

  function validManifest(overrides: Partial<SQLServerFixtureManifest> = {}): SQLServerFixtureManifest {
    const configDb = path.join(tmpDir, 'test_config.db');
    fs.writeFileSync(configDb, '');
    const connConfigPath = path.join(tmpDir, 'connection.json');
    fs.writeFileSync(
      connConfigPath,
      JSON.stringify({ host: '127.0.0.1', user: 'fixture', password: ' ', database: 'gw_disposable_db', schema: 'dbo' }),
    );
    if (process.platform !== 'win32') fs.chmodSync(connConfigPath, 0o600);

    return {
      schema_version: 1,
      disposable: true,
      native_amd64: true,
      host_os: 'linux',
      host_arch: 'amd64',
      browser_storage_mode: 'plain',
      database: 'gw_disposable_db',
      schema: 'dbo',
      server_version_prefix: '16.',
      engine_edition: 3,
      compatibility_level: 160,
      connection_config_file: connConfigPath,
      backend: 'http://127.0.0.1:8080',
      frontend: 'http://127.0.0.1:4173',
      config_db: configDb,
      workspace_id: 'ws-sql-test-1',
      devices: { a: 'meter-a', b: 'meter-b' },
      mappings: [
        { device_id: 'meter-a', tag_id: 'tag-v1', target_type: 'float32' },
        { device_id: 'meter-a', tag_id: 'tag-c1', target_type: 'uint64' },
      ],
      ...overrides,
    };
  }

  it('returns undefined when GW_SQLSERVER_TEST_FIXTURE is unset', () => {
    expect(loadSQLServerFixture()).toBeUndefined();
  });

  it('throws an error when GW_SQLSERVER_TEST_FIXTURE points to a non-existent file', () => {
    process.env[SQLSERVER_FIXTURE_ENV] = path.join(tmpDir, 'missing.json');
    expect(() => loadSQLServerFixture()).toThrow(/not found|does not exist/i);
  });

  it('strictly rejects manifests missing required fields or having invalid editions', () => {
    const m1 = validManifest({ schema_version: 2 });
    expect(() => validateFixtureManifest(m1)).toThrow(/schema_version must be 1/i);

    const m2 = validManifest({ disposable: false });
    expect(() => validateFixtureManifest(m2)).toThrow(/disposable: true/i);

    const m3 = validManifest({ native_amd64: false });
    expect(() => validateFixtureManifest(m3)).toThrow(/native_amd64/i);

    const m4 = validManifest({ server_version_prefix: '15.' });
    expect(() => validateFixtureManifest(m4)).toThrow(/SQL Server 2022/i);

    const m5 = validManifest({ engine_edition: 1 });
    expect(() => validateFixtureManifest(m5)).toThrow(/engine_edition must be 2, 3, or 4/i);

    const m6 = validManifest({ compatibility_level: 150 });
    expect(() => validateFixtureManifest(m6)).toThrow(/compatibility_level must be 160/i);

    const m7 = validManifest({ server_version_prefix: '17.', compatibility_level: 160 });
    expect(() => validateFixtureManifest(m7)).toThrow(/compatibility_level must be 170/i);
  });

  it('rejects non-loopback backend or frontend endpoints', () => {
    const manifest1 = validManifest({ backend: 'http://192.168.1.100:8080' });
    expect(() => validateFixtureManifest(manifest1)).toThrow(/loopback/i);

    const manifest2 = validManifest({ frontend: 'https://remote.example.com' });
    expect(() => validateFixtureManifest(manifest2)).toThrow(/loopback/i);
  });

  it('rejects missing or identical devices', () => {
    const manifestSame = validManifest({ devices: { a: 'dev-1', b: 'dev-1' } });
    expect(() => validateFixtureManifest(manifestSame)).toThrow(/two distinct devices/i);
  });

  it('validates SQL Server identifier escaping against malicious inputs', () => {
    expect(sqlserverIdentifier('dbo', 'gw_readings')).toBe('[dbo].[gw_readings]');
    expect(sqlserverIdentifier('my_schema', 'table_123')).toBe('[my_schema].[table_123]');
    expect(() => sqlserverIdentifier('dbo', 'table; DROP TABLE users;--')).toThrow(/Invalid SQL identifier/i);
    expect(() => sqlserverIdentifier('dbo;--', 'table')).toThrow(/Invalid SQL identifier/i);
  });

  it('loads valid manifest correctly from file', () => {
    const manifest = validManifest();
    const manifestPath = path.join(tmpDir, 'valid_manifest.json');
    fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2));

    process.env[SQLSERVER_FIXTURE_ENV] = manifestPath;
    const loaded = loadSQLServerFixture();
    expect(loaded).toBeDefined();
    expect(loaded?.workspace_id).toBe('ws-sql-test-1');
    expect(loaded?.database).toBe('gw_disposable_db');
    expect(loaded?.schema).toBe('dbo');
    expect(loaded?.devices).toEqual({ a: 'meter-a', b: 'meter-b' });
  });

  it('validates connection_config_file with strict user key, preserved password literals, and TLS safety', () => {
    const configPath = path.join(tmpDir, 'conn_check.json');
    const validConfig = {
      host: '127.0.0.1',
      user: 'sql_agent',
      password: ' ',
      database: 'gw_db',
      schema: 'dbo',
    };
    fs.writeFileSync(configPath, JSON.stringify(validConfig));
    if (process.platform !== 'win32') fs.chmodSync(configPath, 0o600);

    const loaded = validateConnectionConfigFile(configPath, 'gw_db', 'dbo');
    expect(loaded.user).toBe('sql_agent');
    expect(loaded.password).toBe(' ');
    expect(loaded.encrypt).toBe(true);
    expect(loaded.trust_server_certificate).toBe(false);

    // Rejects config with encrypt=false, string "false", "0", "no"
    for (const badEncrypt of [false, 'false', '0', 'no', 'disable']) {
      fs.writeFileSync(configPath, JSON.stringify({ ...validConfig, encrypt: badEncrypt }));
      expect(() => validateConnectionConfigFile(configPath, 'gw_db', 'dbo')).toThrow(/cannot disable TLS/i);
    }

    // Rejects trustservercertificate aliases and string "true", "1", "yes"
    for (const key of ['trustservercertificate', 'trust_server_certificate', 'trustServerCertificate', 'TrustServerCertificate']) {
      for (const badTrust of [true, 'true', '1', 'yes']) {
        fs.writeFileSync(configPath, JSON.stringify({ ...validConfig, [key]: badTrust }));
        expect(() => validateConnectionConfigFile(configPath, 'gw_db', 'dbo')).toThrow(/cannot bypass certificate verification/i);
      }
    }

    // Rejects config using 'username' instead of backend 'user' key
    const invalidAliasConfig = {
      host: '127.0.0.1',
      username: 'sql_agent',
      password: 'secret',
      database: 'gw_db',
      schema: 'dbo',
    };
    fs.writeFileSync(configPath, JSON.stringify(invalidAliasConfig));
    if (process.platform !== 'win32') fs.chmodSync(configPath, 0o600);
    expect(() => validateConnectionConfigFile(configPath, 'gw_db', 'dbo')).toThrow(
      /lacks required SQL authentication "user" field/i,
    );

    // Rejects if mode is broader than 0600 on POSIX
    if (process.platform !== 'win32') {
      fs.writeFileSync(configPath, JSON.stringify(validConfig));
      fs.chmodSync(configPath, 0o644);
      expect(() => validateConnectionConfigFile(configPath, 'gw_db', 'dbo')).toThrow(/permissions are too broad/i);
    }
  });

  it('validates browser_storage_mode strictly requires plain or native_monthly', () => {
    const manifestPlain = validManifest({ browser_storage_mode: 'plain' });
    expect(() => validateFixtureManifest(manifestPlain)).not.toThrow();

    const manifestNative = validManifest({ browser_storage_mode: 'native_monthly' });
    expect(() => validateFixtureManifest(manifestNative)).not.toThrow();

    const manifestInvalid = validManifest({ browser_storage_mode: 'timescale' as 'plain' });
    expect(() => validateFixtureManifest(manifestInvalid)).toThrow(/manifest browser_storage_mode/i);
  });

  it('accepts both native windows and linux amd64 SQL Server hosts', () => {
    const manifestLinux = validManifest({ host_os: 'linux', host_arch: 'amd64' });
    expect(() => validateFixtureManifest(manifestLinux)).not.toThrow();

    const manifestWindows = validManifest({ host_os: 'windows', host_arch: 'amd64' });
    expect(() => validateFixtureManifest(manifestWindows)).not.toThrow();

    const manifestInvalid = validManifest({ host_os: 'openbsd' as 'linux', host_arch: 'amd64' });
    expect(() => validateFixtureManifest(manifestInvalid)).toThrow(/manifest host_os/i);
  });

  it('accepts truthful SQL Server 2025 on emulated arm64 darwin host with exact 17.x and compat 170', () => {
    const manifest2025 = validManifest({
      native_amd64: false,
      host_os: 'darwin' as 'linux',
      host_arch: 'arm64' as 'amd64',
      server_version_prefix: '17.',
      compatibility_level: 170,
    });
    expect(() => validateFixtureManifest(manifest2025)).not.toThrow();
  });

  it('rejects contradictory manifest claiming native_amd64: true on arm64 host', () => {
    const contradictory = validManifest({
      native_amd64: true,
      host_arch: 'arm64' as 'amd64',
    });
    expect(() => validateFixtureManifest(contradictory)).toThrow(/native_amd64: true contradicts host_arch/i);
  });

  it('enforces readonly SQL only in readback helper and rejects mutating statements', () => {
    const config: SQLServerConnectionConfig = {
      host: '127.0.0.1',
      user: 'reader',
      password: 'pwd',
      database: 'gw_db',
      schema: 'dbo',
    };

    expect(() => querySQLServerReadback(config, 'DROP TABLE gw_readings')).toThrow(/Refusing non-readonly SQL/i);
    expect(() => querySQLServerReadback(config, 'INSERT INTO gw_readings VALUES (1)')).toThrow(/Refusing non-readonly SQL/i);
    expect(() => querySQLServerReadback(config, 'DELETE FROM gw_readings')).toThrow(/Refusing non-readonly SQL/i);
    expect(() => querySQLServerReadback(config, 'TRUNCATE TABLE gw_readings')).toThrow(/Refusing non-readonly SQL/i);
  });

  it('safely decodes and validates delivery outbox wire payload strictly matching payload.go', () => {
    const validPayload = JSON.stringify({
      version: 1,
      record_id: 'rec-001',
      effect_key: 'eff-abc-123',
      entity_key: 'ent-456',
      bucket_start: '2026-10-26T12:00:00Z',
      partial: false,
      cells: [
        { column: 'val_i64', kind: 'int64', value: '1234567890123' },
        { column: 'val_f64', kind: 'float64', value: '123.456' },
        { column: 'val_str', kind: 'string', value: 'meter_ok' },
        { column: 'val_b', kind: 'bool', value: 'true' },
        { column: 'val_t', kind: 'time', value: '2026-10-26T12:00:00Z' },
        { column: 'val_null', kind: 'null' },
      ],
    });

    const decoded = parseWireRowPayload(validPayload);
    expect(decoded.version).toBe(1);
    expect(decoded.record_id).toBe('rec-001');
    expect(decoded.effect_key).toBe('eff-abc-123');
    expect(decoded.cells).toHaveLength(6);
    expect(decoded.cells[0]).toEqual({ column: 'val_i64', kind: 'int64', value: '1234567890123' });
    expect(decoded.cells[5]).toEqual({ column: 'val_null', kind: 'null' });

    // Rejects unsupported version
    const badVersion = JSON.stringify({ ...JSON.parse(validPayload), version: 2 });
    expect(() => parseWireRowPayload(badVersion)).toThrow(/unsupported version 2/i);

    // Rejects missing record_id or effect_key
    const noRecordId = JSON.stringify({ ...JSON.parse(validPayload), record_id: '' });
    expect(() => parseWireRowPayload(noRecordId)).toThrow(/lacks valid record_id/i);

    const noEffectKey = JSON.stringify({ ...JSON.parse(validPayload), effect_key: '' });
    expect(() => parseWireRowPayload(noEffectKey)).toThrow(/lacks valid effect_key/i);

    // Rejects invalid cell kind
    const badCellKind = JSON.stringify({
      version: 1,
      record_id: 'rec-1',
      effect_key: 'eff-1',
      bucket_start: '2026-10-26T00:00:00Z',
      cells: [{ column: 'col1', kind: 'unknown_kind', value: '1' }],
    });
    expect(() => parseWireRowPayload(badCellKind)).toThrow(/unsupported kind/i);

    // Rejects non-string cell value (ensures no numeric rounding in wire format)
    const badCellValue = JSON.stringify({
      version: 1,
      record_id: 'rec-1',
      effect_key: 'eff-1',
      bucket_start: '2026-10-26T00:00:00Z',
      cells: [{ column: 'col1', kind: 'int64', value: 123 }],
    });
    expect(() => parseWireRowPayload(badCellValue)).toThrow(/value must be a string/i);
  });

  it('canonicalizes decimal text by trimming trailing fraction zeros without converting to Number', () => {
    expect(canonicalizeDecimalText('0.400000')).toBe('0.4');
    expect(canonicalizeDecimalText('0.4')).toBe('0.4');
    expect(canonicalizeDecimalText('123.0000')).toBe('123');
    expect(canonicalizeDecimalText('0.00')).toBe('0');
    expect(canonicalizeDecimalText('-45.6700')).toBe('-45.67');
    expect(canonicalizeDecimalText('100')).toBe('100');
    // Preserves huge uint64/bigint precision beyond JS Number.MAX_SAFE_INTEGER
    expect(canonicalizeDecimalText('18446744073709551615.000000')).toBe('18446744073709551615');
  });

  it('canonicalizes UTC ISO time text by padding fraction to 6 microsecond digits', () => {
    expect(canonicalizeUtcTimeText('2026-10-26T12:00:00Z')).toBe('2026-10-26T12:00:00.000000Z');
    expect(canonicalizeUtcTimeText('2026-10-26T12:00:00.123Z')).toBe('2026-10-26T12:00:00.123000Z');
    expect(canonicalizeUtcTimeText('2026-10-26T12:00:00.123456Z')).toBe('2026-10-26T12:00:00.123456Z');
    expect(canonicalizeUtcTimeText('2026-10-26T12:00:00.123456000Z')).toBe('2026-10-26T12:00:00.123456Z');
    expect(() => canonicalizeUtcTimeText('2026-10-26T12:00:00.123456789Z')).toThrow();
  });

  it('builds safely quoted SQL Server row query using typed scalar NVARCHAR and FOR JSON PATH', () => {
    const sql = buildSQLServerRowSelectSql('dbo', 'gw_readings', 'rec-001', [
      { column: 'val_b', kind: 'bool', value: 'true' },
      { column: 'bucket_start', kind: 'time', value: '2026-10-26T12:00:00Z' },
      { column: 'val_dec', kind: 'string', value: '0.4' },
      { column: 'val_i', kind: 'int64', value: '18446744073709551615' },
      { column: 'delta_value', kind: 'string', value: '0.4' },
    ]);

    expect(sql).toContain('SELECT TOP 1');
    expect(sql).toContain('[dbo].[gw_readings]');
    expect(sql).toContain("CASE WHEN [val_b] = 1 THEN 'true' WHEN [val_b] = 0 THEN 'false' ELSE NULL END AS [val_b]");
    expect(sql).toContain("CASE WHEN [bucket_start] IS NOT NULL THEN CONVERT(NVARCHAR(33), [bucket_start], 126) + 'Z' ELSE NULL END AS [bucket_start]");
    expect(sql).toContain('CAST([val_dec] AS NVARCHAR(MAX)) AS [val_dec]');
    expect(sql).toContain('CAST([val_i] AS NVARCHAR(MAX)) AS [val_i]');
    expect(sql).toContain('CAST([delta_value] AS NVARCHAR(MAX)) AS [delta_value]');
    expect(sql).toContain("WHERE [record_id] = 'rec-001'");
    expect(sql).toContain('FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER');
  });

  it('verifies outbox wire payload against SQL Server JSON row for bool, null, uint64max, decimal 0.4, UTC microseconds, Unicode, and provenance', () => {
    const wire: WireRowPayload = {
      version: 1,
      record_id: 'rec-001',
      effect_key: 'eff-123',
      bucket_start: '2026-10-26T12:00:00Z',
      cells: [
        { column: 'val_bool', kind: 'bool', value: 'true' },
        { column: 'val_null', kind: 'null' },
        { column: 'val_uint64', kind: 'int64', value: '18446744073709551615' },
        { column: 'value_decimal', kind: 'string', value: '0.4' },
        { column: 'delta_value', kind: 'string', value: '0.4' },
        { column: 'bucket_start', kind: 'time', value: '2026-10-26T12:00:00Z' },
        { column: 'circuit_name', kind: 'string', value: '台電 1600kVA 高壓側' },
        { column: 'provenance', kind: 'string', value: '{"source":"modbus_tcp","device_id":"d1"}' },
      ],
    };

    // SQL Server returns padded decimal, microsecond time, matching bool, null, unicode, and large int64
    const validRow: Record<string, string | null> = {
      record_id: 'rec-001',
      val_bool: 'true',
      val_null: null,
      val_uint64: '18446744073709551615',
      value_decimal: '0.400000',
      delta_value: '0.400000',
      bucket_start: '2026-10-26T12:00:00.000000Z',
      circuit_name: '台電 1600kVA 高壓側',
      provenance: '{"device_id":"d1","source":"modbus_tcp"}', // JSON key reordered
    };

    expect(() => verifyWirePayloadAgainstRow(wire, validRow)).not.toThrow();

    // Rejects bool mismatch
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, val_bool: 'false' })).toThrow(/Bool column.*mismatch/i);

    // Rejects null mismatch
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, val_null: 'not_null' })).toThrow(/expected null/i);
    const missingNull = { ...validRow };
    delete missingNull.val_null;
    expect(() => verifyWirePayloadAgainstRow(wire, missingNull)).toThrow(/expected null/i);

    // Rejects uint64 mismatch (e.g. if floating point rounding altered the exact string)
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, val_uint64: '18446744073709552000' })).toThrow(/Int64 column.*mismatch/i);

    // Rejects decimal mismatch
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, value_decimal: '0.500000' })).toThrow(/Decimal column.*mismatch/i);

    // Rejects time mismatch
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, bucket_start: '2026-10-26T12:00:01.000000Z' })).toThrow(/Time column.*mismatch/i);

    // Rejects Unicode mismatch
    expect(() => verifyWirePayloadAgainstRow(wire, { ...validRow, circuit_name: '台電 3200kVA' })).toThrow(/String column.*mismatch/i);
  });

  it('rejects non-finite float readback instead of treating NaN as equal', () => {
    const wire: WireRowPayload = {
      version: 1, record_id: 'record', effect_key: 'effect', bucket_start: '2026-10-26T00:00:00Z',
      cells: [{ column: 'value_float', kind: 'float64', value: '1.25' }],
    };
    expect(() => verifyWirePayloadAgainstRow(wire, { value_float: 'NaN' })).toThrow();
  });
});
