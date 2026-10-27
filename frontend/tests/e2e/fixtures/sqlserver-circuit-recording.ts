import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { DatabaseSync, type SQLInputValue } from 'node:sqlite';
import type { StudioV2SQLServerRecordingCapability } from '../../../src/types/studioV2SQLServerRecordingCapability';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';

export const SQLSERVER_FIXTURE_ENV = 'GW_SQLSERVER_TEST_FIXTURE';

export interface SQLServerFixtureManifest {
  schema_version: number; disposable: boolean; native_amd64: boolean;
  host_os: 'linux' | 'windows' | 'darwin'; host_arch: 'amd64' | 'arm64';
  browser_storage_mode: 'plain' | 'native_monthly';
  database: string; schema: string;
  server_version_prefix: string; server_edition?: string;
  engine_edition: number; compatibility_level: number;
  connection_config_file: string; image_digest?: string; provenance?: string;
  backend: string; frontend: string; config_db: string; workspace_id: string;
  devices: { a: string; b: string };
  mappings: { device_id: string; tag_id: string; target_type: string }[];
}

export interface SQLServerConnectionConfig {
  host: string; port?: number; user: string; password: string;
  database: string; schema: string;
  encrypt?: boolean; trust_server_certificate?: boolean;
}

const IDENTIFIER_REGEX = /^[a-zA-Z_][a-zA-Z0-9_]*$/;
const READONLY_SQL_REGEX = /^\s*SELECT\b/i;
const MUTATING_SQL_REGEX = /\b(INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|EXEC|EXECUTE|GRANT|REVOKE)\b/i;

/** Safely format schema and table names as SQL Server bracketed identifiers [schema].[table]. */
export function sqlserverIdentifier(schema: string, table: string): string {
  if (!IDENTIFIER_REGEX.test(schema) || !IDENTIFIER_REGEX.test(table)) {
    throw new Error(`Invalid SQL identifier for SQL Server target: "${schema}"."${table}"`);
  }
  return `[${schema}].[${table}]`;
}

/** Validate manifest contents against strict shared fixture contract requirements. */
export function validateFixtureManifest(
  manifest: SQLServerFixtureManifest,
  manifestPath?: string,
): SQLServerFixtureManifest {
  if (manifest.schema_version !== 1) throw new Error(`manifest schema_version must be 1, received: ${String(manifest.schema_version)}`);
  if (manifest.disposable !== true) throw new Error('manifest must attest disposable: true');
  if (typeof manifest.native_amd64 !== 'boolean') throw new Error('manifest must attest native_amd64: true or false');
  if (manifest.native_amd64 === true && manifest.host_arch !== 'amd64') {
    throw new Error(`manifest native_amd64: true contradicts host_arch: "${String(manifest.host_arch)}"`);
  }
  if (manifest.native_amd64 === false && manifest.host_arch !== 'arm64') {
    throw new Error(`manifest native_amd64: false contradicts host_arch: "${String(manifest.host_arch)}"`);
  }
  if (manifest.host_os !== 'linux' && manifest.host_os !== 'windows' && manifest.host_os !== 'darwin') {
    throw new Error(`manifest host_os must be "linux", "windows", or "darwin", received: "${String(manifest.host_os)}"`);
  }
  if (manifest.host_arch !== 'amd64' && manifest.host_arch !== 'arm64') {
    throw new Error(`manifest host_arch must be "amd64" or "arm64", received: "${String(manifest.host_arch)}"`);
  }
  if (manifest.browser_storage_mode !== 'plain' && manifest.browser_storage_mode !== 'native_monthly') {
    throw new Error(`manifest browser_storage_mode must be "plain" or "native_monthly", received: "${String(manifest.browser_storage_mode)}"`);
  }
  const is2022 = manifest.server_version_prefix?.startsWith('16.');
  const is2025 = manifest.server_version_prefix?.startsWith('17.');
  if (!is2022 && !is2025) {
    throw new Error('manifest server_version_prefix must target SQL Server 2022 ("16.*") or 2025 ("17.*")');
  }
  if (manifest.engine_edition !== 2 && manifest.engine_edition !== 3 && manifest.engine_edition !== 4) {
    throw new Error(`manifest engine_edition must be 2, 3, or 4, received: ${String(manifest.engine_edition)}`);
  }
  if (is2022 && manifest.compatibility_level !== 160) {
    throw new Error(`manifest compatibility_level must be 160 for SQL Server 2022, received: ${String(manifest.compatibility_level)}`);
  }
  if (is2025 && manifest.compatibility_level !== 170) {
    throw new Error(`manifest compatibility_level must be 170 for SQL Server 2025, received: ${String(manifest.compatibility_level)}`);
  }

  for (const urlStr of [manifest.backend, manifest.frontend]) {
    try {
      const url = new URL(urlStr);
      if (url.protocol !== 'http:' || !['127.0.0.1', 'localhost'].includes(url.hostname)) {
        throw new Error(`Only explicit loopback fixture URLs are allowed for safety. Rejected: ${urlStr}`);
      }
    } catch (err) {
      if (err instanceof Error && err.message.includes('loopback')) throw err;
      throw new Error(`Malformed URL in fixture: ${urlStr}`, { cause: err });
    }
  }

  if (manifest.config_db && !fs.existsSync(manifest.config_db)) {
    const resolved = manifestPath && !path.isAbsolute(manifest.config_db) ? path.resolve(path.dirname(manifestPath), manifest.config_db) : '';
    if (resolved && fs.existsSync(resolved)) {
      manifest.config_db = resolved;
    } else {
      throw new Error(`Config database does not exist: ${manifest.config_db}`);
    }
  }

  if (!manifest.workspace_id) throw new Error('Fixture manifest must provide a non-empty workspace_id.');
  if (!manifest.devices || !manifest.devices.a || !manifest.devices.b || manifest.devices.a === manifest.devices.b) {
    throw new Error('Fixture manifest must provide two distinct devices.');
  }
  if (!IDENTIFIER_REGEX.test(manifest.database) || !IDENTIFIER_REGEX.test(manifest.schema)) {
    throw new Error(`Invalid database or schema identifier: database="${manifest.database}", schema="${manifest.schema}"`);
  }
  if (!manifest.connection_config_file) throw new Error('manifest connection_config_file is required.');
  validateConnectionConfigFile(manifest.connection_config_file, manifest.database, manifest.schema);
  return manifest;
}

/** Validate saved SQL Server connection config file existence, permissions (0600), and fields. */
export function validateConnectionConfigFile(
  configPath: string,
  expectedDatabase: string,
  expectedSchema: string,
): SQLServerConnectionConfig {
  if (!path.isAbsolute(configPath)) throw new Error(`connection_config_file must be an absolute path: ${configPath}`);
  if (!fs.existsSync(configPath)) throw new Error(`connection_config_file does not exist: ${configPath}`);
  const stat = fs.statSync(configPath);
  if (!stat.isFile()) throw new Error(`connection_config_file is not a regular file: ${configPath}`);
  if (process.platform !== 'win32' && (stat.mode & 0o077) !== 0) {
    throw new Error(`connection_config_file permissions are too broad: ${configPath}. Expected mode 0600.`);
  }

  const raw = fs.readFileSync(configPath, 'utf8');
  let parsed: Record<string, unknown>;
  try {
    parsed = JSON.parse(raw) as Record<string, unknown>;
  } catch (err) {
    throw new Error(`connection_config_file is not valid JSON: ${configPath}`, { cause: err });
  }

  if (parsed.database !== expectedDatabase || parsed.schema !== expectedSchema) {
    throw new Error('connection_config_file scope (database/schema) differs from manifest.');
  }
  const host = typeof parsed.host === 'string' ? parsed.host.trim() : '';
  if (!host) throw new Error('connection_config_file lacks required host.');
  const user = typeof parsed.user === 'string' ? parsed.user.trim() : '';
  if (!user) throw new Error('connection_config_file lacks required SQL authentication "user" field.');
  if (typeof parsed.password !== 'string' || parsed.password.length === 0) {
    throw new Error('connection_config_file lacks required non-empty "password".');
  }

  // Reject insecure TLS defaults and aliases as backend does (no -C, verified TLS only)
  if (parsed.encrypt !== undefined) {
    const encStr = String(parsed.encrypt).trim().toLowerCase();
    if (encStr === 'false' || encStr === '0' || encStr === 'no' || encStr === 'disable' ||
        !['true', 'mandatory', 'yes', '1', 't', 'strict'].includes(encStr)) {
      throw new Error('connection_config_file cannot disable TLS: encrypt=false rejected.');
    }
  }

  for (const key of ['trustservercertificate', 'trust_server_certificate', 'trustServerCertificate', 'TrustServerCertificate']) {
    if (parsed[key] !== undefined) {
      const valStr = String(parsed[key]).trim().toLowerCase();
      if (valStr === 'true' || valStr === '1' || valStr === 'yes' || valStr === 't') {
        throw new Error(`connection_config_file cannot bypass certificate verification: ${key}=true rejected.`);
      }
    }
  }

  return {
    host,
    port: typeof parsed.port === 'number' ? parsed.port : undefined,
    user,
    password: parsed.password,
    database: String(parsed.database),
    schema: String(parsed.schema),
    encrypt: true,
    trust_server_certificate: false,
  };
}

/**
 * Load the SQL Server test fixture strictly from GW_SQLSERVER_TEST_FIXTURE.
 * Returns undefined if GW_SQLSERVER_TEST_FIXTURE is unset.
 */
export function loadSQLServerFixture(): SQLServerFixtureManifest | undefined {
  const filePath = process.env[SQLSERVER_FIXTURE_ENV];
  if (!filePath || filePath.trim() === '') return undefined;

  const resolved = path.resolve(filePath.trim());
  if (!fs.existsSync(resolved)) {
    throw new Error(`${SQLSERVER_FIXTURE_ENV} file does not exist: ${resolved}`);
  }

  const raw = fs.readFileSync(resolved, 'utf8');
  const parsed = JSON.parse(raw) as SQLServerFixtureManifest;
  return validateFixtureManifest(parsed, resolved);
}

/** Find a fixed sqlcmd binary in standard system locations or PATH. */
export function resolveSqlCmdBinary(): string | undefined {
  const candidates = [
    '/opt/homebrew/bin/sqlcmd', '/usr/local/bin/sqlcmd',
    '/opt/mssql-tools/bin/sqlcmd', '/opt/mssql-tools18/bin/sqlcmd',
  ];
  for (const c of candidates) if (fs.existsSync(c)) return c;
  const exe = process.platform === 'win32' ? 'sqlcmd.exe' : 'sqlcmd';
  for (const dir of (process.env.PATH ?? '').split(path.delimiter)) {
    if (dir && fs.existsSync(path.join(dir, exe))) return path.join(dir, exe);
  }
  return undefined;
}

/** Read-only SQL Server readback using fixed sqlcmd binary with argument arrays and env secrets. */
export function querySQLServerReadback(
  connectionConfig: SQLServerConnectionConfig,
  readonlySql: string,
): string {
  if (!READONLY_SQL_REGEX.test(readonlySql) || MUTATING_SQL_REGEX.test(readonlySql)) {
    throw new Error(`Refusing non-readonly SQL query for readback: "${readonlySql}"`);
  }

  const sqlcmd = resolveSqlCmdBinary();
  if (!sqlcmd) {
    throw new Error('sqlcmd binary unavailable on this host (NOT_RUN)');
  }

  const server = connectionConfig.port ? `${connectionConfig.host},${connectionConfig.port}` : connectionConfig.host;
  // Arguments array with no shell interpolation; never pass password via -P on command line.
  // Enforce explicit verified TLS (-N) and NEVER -C. Use -b so SQL errors abort and exit with code 1.
  // Enforce -y0 for untruncated NVARCHAR(MAX)/provenance; drop incompatible -W.
  const args = ['-b', '-S', server, '-d', connectionConfig.database, '-U', connectionConfig.user, '-N', '-y0', '-Q', readonlySql, '-h', '-1'];
  const env = { ...process.env, SQLCMDPASSWORD: connectionConfig.password };

  try {
    return execFileSync(sqlcmd, args, {
      env,
      encoding: 'utf8',
      timeout: 15_000,
      stdio: ['ignore', 'pipe', 'pipe'],
    });
  } catch {
    // Drop raw exec cause to prevent any leakage of command arguments or stderr
    throw new Error('sqlcmd execution failed during verified readback check');
  }
}

/** Execute a read-only query against the local SQLite config database. */
export function readConfigSQL(file: string, sql: string, args: SQLInputValue[] = []): Record<string, string | number | null>[] {
  const db = new DatabaseSync(file, { readOnly: true });
  try { return db.prepare(sql).all(...args) as Record<string, string | number | null>[]; } finally { db.close(); }
}

/** Query the workspace API on the running gateway backend. */
export async function api<T>(fixture: SQLServerFixtureManifest, endpoint: string): Promise<T> {
  let res: Response;
  try { res = await fetch(`${fixture.backend}/api/v1/datalink/studio-v2/workspace${endpoint}`); } catch { throw new Error(`API network error while reaching ${endpoint}`); }
  if (!res.ok) throw new Error(`API request failed with status: ${res.status}`);
  let body: unknown;
  try { body = await res.json(); } catch { throw new Error(`API returned invalid JSON for ${endpoint}`); }
  if (!body || typeof body !== 'object' || !(body as { success?: boolean }).success) throw new Error(`API returned unsuccessful response for ${endpoint}`);
  return (body as { data: T }).data;
}

const SQLSERVER_PERMISSION_FIELDS = [
  'verified', 'database_view_definition', 'schema_view_definition', 'create_table',
  'alter_schema', 'select_schema', 'insert_schema', 'delete_schema', 'alter_any_dataspace',
] as const;

/** Safely validate SQL Server capability projection matching backend and service contract. */
export function validateSQLServerCapability(
  value: unknown,
  expectedConnectorId: string,
  expectedConnectorRevision: string,
): StudioV2SQLServerRecordingCapability {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('SQL Server capability response is malformed');
  const record = value as Record<string, unknown>;
  const textFields = ['connector_id', 'connector_revision', 'database', 'schema', 'server_version', 'server_edition', 'capability_revision', 'metadata_status'] as const;
  for (const field of textFields) {
    if (typeof record[field] !== 'string' || (record[field] as string).length > 512) throw new Error(`SQL Server capability response field ${field} is malformed`);
  }
  const { connector_id: connectorId, connector_revision: connectorRevision, engine_edition: engineEdition, compatibility_level: compatibilityLevel } = record;
  if (typeof engineEdition !== 'number' || !Number.isInteger(engineEdition) || typeof compatibilityLevel !== 'number' || !Number.isInteger(compatibilityLevel)) {
    throw new Error('SQL Server capability response is malformed');
  }
  const booleanFields = ['connected', 'supports_managed_long', 'supports_native_partition', 'supports_group_test_write', 'supports_missing_repair'] as const;
  const permissions = record.permissions as Record<string, unknown> | undefined;
  if (!permissions || typeof permissions !== 'object' || Array.isArray(permissions)) throw new Error('SQL Server capability response is malformed');
  if (connectorId !== expectedConnectorId || connectorRevision !== expectedConnectorRevision || !connectorId || !connectorRevision || !record.capability_revision ||
    booleanFields.some((f) => typeof record[f] !== 'boolean') || SQLSERVER_PERMISSION_FIELDS.some((f) => typeof permissions[f] !== 'boolean') ||
    !Array.isArray(record.reasons) || record.reasons.length > 64 || record.reasons.some((r) => typeof r !== 'string' || !r || r.length > 128)) {
    throw new Error('SQL Server capability response does not match the saved connector');
  }
  return record as unknown as StudioV2SQLServerRecordingCapability;
}

export type WirePayloadCellKind = 'null' | 'bool' | 'int64' | 'float64' | 'string' | 'time';

export interface WirePayloadCell {
  column: string;
  kind: WirePayloadCellKind;
  value?: string;
}

export interface WireRowPayload {
  version: number;
  record_id: string;
  effect_key: string;
  entity_key?: string;
  bucket_start: string;
  partial?: boolean;
  cells: WirePayloadCell[];
}

/** Safely decode and validate delivery outbox wire payload strictly matching payload.go */
export function parseWireRowPayload(raw: string): WireRowPayload {
  let parsed: unknown;
  try { parsed = JSON.parse(raw); } catch { throw new Error('Outbox payload is not valid JSON'); }
  if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) throw new Error('Outbox payload must be a JSON object');
  const obj = parsed as Record<string, unknown>;
  if (obj.version !== 1) throw new Error(`Outbox payload unsupported version ${String(obj.version)}, expected 1`);
  if (typeof obj.record_id !== 'string' || !obj.record_id) throw new Error('Outbox payload lacks valid record_id');
  if (typeof obj.effect_key !== 'string' || !obj.effect_key) throw new Error('Outbox payload lacks valid effect_key');
  if (typeof obj.bucket_start !== 'string' || !obj.bucket_start) throw new Error('Outbox payload lacks valid bucket_start');
  if (!Array.isArray(obj.cells)) throw new Error('Outbox payload lacks valid cells array');

  const validKinds = new Set(['null', 'bool', 'int64', 'float64', 'string', 'time']);
  const validatedCells: WirePayloadCell[] = [];
  for (const item of obj.cells) {
    if (!item || typeof item !== 'object' || Array.isArray(item)) throw new Error('Outbox payload cell is not an object');
    const cell = item as Record<string, unknown>;
    if (typeof cell.column !== 'string' || !cell.column) throw new Error('Outbox payload cell missing column');
    if (typeof cell.kind !== 'string' || !validKinds.has(cell.kind)) throw new Error(`Outbox payload cell has unsupported kind: ${String(cell.kind)}`);
    if (cell.value !== undefined && typeof cell.value !== 'string') throw new Error('Outbox payload cell value must be a string if present');
    validatedCells.push({ column: cell.column, kind: cell.kind as WirePayloadCellKind, ...(cell.value !== undefined ? { value: cell.value } : {}) });
  }
  return {
    version: 1, record_id: obj.record_id, effect_key: obj.effect_key,
    ...(typeof obj.entity_key === 'string' ? { entity_key: obj.entity_key } : {}),
    bucket_start: obj.bucket_start, ...(typeof obj.partial === 'boolean' ? { partial: obj.partial } : {}),
    cells: validatedCells,
  };
}

/** Fetch current saved write groups. */
export async function groups(fixture: SQLServerFixtureManifest): Promise<WriteGroup[]> {
  return (await api<{ groups: WriteGroup[] }>(fixture, '/write-groups')).groups;
}

/** Read recording start operations only from the config ledger. */
export function recordingStartOperations(fixture: SQLServerFixtureManifest) {
  return readConfigSQL(fixture.config_db, "SELECT operation_id, action, status, payload_digest, detail FROM managed_schema_operations WHERE workspace_id=? AND action='recording_start' ORDER BY created_at", [fixture.workspace_id]);
}

/** Read schema apply operations from the config ledger. */
export function schemaApplyOperations(fixture: SQLServerFixtureManifest) {
  return readConfigSQL(fixture.config_db, "SELECT operation_id, action, status, payload_digest, detail FROM managed_schema_operations WHERE workspace_id=? AND action='schema_apply' ORDER BY created_at", [fixture.workspace_id]);
}

/** Read frozen applied runtime versions from the config database. */
export function runtimeVersions(fixture: SQLServerFixtureManifest) {
  return readConfigSQL(fixture.config_db, 'SELECT group_id, group_revision, payload FROM wg_runtime_versions WHERE workspace_id=? ORDER BY group_id, group_revision', [fixture.workspace_id]);
}

/** Read committed delivery effects from the outbox table (mandatory state: sql_committed). */
export function deliveryCommittedOutbox(fixture: SQLServerFixtureManifest, groupId: string, groupRevision: string) {
  return readConfigSQL(fixture.config_db, "SELECT record_id, effect_key, payload_digest, group_revision, connector_revision, table_name, state, bucket_start, payload FROM wg_delivery_outbox WHERE workspace_id=? AND group_id=? AND group_revision=? AND state='sql_committed' ORDER BY bucket_start", [fixture.workspace_id, groupId, groupRevision]);
}

/** Canonicalize decimal string by trimming trailing fraction zeros without converting to Number. */
export function canonicalizeDecimalText(val: string): string {
  const trimmed = val.trim();
  if (!trimmed.includes('.')) return trimmed;
  const [intPart, fracPart] = trimmed.split('.');
  const stripped = fracPart.replace(/0+$/, '');
  return stripped.length > 0 ? `${intPart}.${stripped}` : intPart;
}

/** Canonicalize UTC time string to 6-digit microsecond fraction (e.g. 2026-10-26T12:00:00.000000Z). */
export function canonicalizeUtcTimeText(val: string): string {
  let s = val.trim();
  if (!s.endsWith('Z')) s += 'Z';
  const match = /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d+))?Z$/.exec(s);
  if (!match || /[1-9]/.test((match[2] ?? '').slice(6))) {
    throw new Error('SQL timestamp is not an exact UTC microsecond value');
  }
  const frac = (match[2] ?? '').padEnd(6, '0').slice(0, 6);
  return `${match[1]}.${frac}Z`;
}

/** Build safely quoted SQL Server query returning a single JSON object for all target cells. */
export function buildSQLServerRowSelectSql(
  schema: string,
  table: string,
  recordId: string,
  cells: WirePayloadCell[],
): string {
  const targetTable = sqlserverIdentifier(schema, table);
  const colMap = new Map<string, WirePayloadCellKind>();
  for (const cell of cells) {
    if (!IDENTIFIER_REGEX.test(cell.column)) throw new Error(`Invalid column identifier: ${cell.column}`);
    colMap.set(cell.column, cell.kind);
  }
  if (!colMap.has('record_id')) colMap.set('record_id', 'string');

  const selectParts: string[] = [];
  for (const [col, kind] of colMap.entries()) {
    const quoted = `[${col}]`;
    if (kind === 'bool' || col === 'value_bool') {
      selectParts.push(`CASE WHEN ${quoted} = 1 THEN 'true' WHEN ${quoted} = 0 THEN 'false' ELSE NULL END AS ${quoted}`);
    } else if (kind === 'time' || ['observed_at', 'received_at', 'bucket_start', 'bucket_end', 'delta_from', 'delta_to'].includes(col)) {
      selectParts.push(`CASE WHEN ${quoted} IS NOT NULL THEN CONVERT(NVARCHAR(33), ${quoted}, 126) + 'Z' ELSE NULL END AS ${quoted}`);
    } else {
      selectParts.push(`CAST(${quoted} AS NVARCHAR(MAX)) AS ${quoted}`);
    }
  }

  const safeRecordId = recordId.replace(/'/g, "''");
  return `SELECT TOP 1 ${selectParts.join(', ')} FROM ${targetTable} WHERE [record_id] = '${safeRecordId}' FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER`;
}

function deepEqualJson(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  if (!a || !b || typeof a !== 'object' || typeof b !== 'object') return false;
  if (Array.isArray(a) !== Array.isArray(b)) return false;
  if (Array.isArray(a) && Array.isArray(b)) {
    return a.length === b.length && a.every((v, i) => deepEqualJson(v, b[i]));
  }
  const ka = Object.keys(a as Record<string, unknown>).sort();
  const kb = Object.keys(b as Record<string, unknown>).sort();
  return ka.length === kb.length && ka.every((k, i) => k === kb[i] && deepEqualJson((a as Record<string, unknown>)[k], (b as Record<string, unknown>)[k]));
}

/** Safely compare an outbox wire row against the row returned by SQL Server. */
export function verifyWirePayloadAgainstRow(
  wire: WireRowPayload,
  row: Record<string, string | null>,
): void {
  for (const cell of wire.cells) {
    const actual = row[cell.column];
    const expected = cell.value;
    if (cell.kind === 'null') {
      if (actual !== null) throw new Error(`Column "${cell.column}" expected null, received: "${String(actual)}"`);
      continue;
    }
    if (actual === null || actual === undefined) throw new Error(`Column "${cell.column}" expected "${String(expected)}", but actual value is null/missing`);

    const actualStr = String(actual);
    const expectedStr = String(expected ?? '');

    switch (cell.kind) {
      case 'bool':
        if (actualStr !== expectedStr) throw new Error(`Bool column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
        break;
      case 'int64':
        if (actualStr !== expectedStr) throw new Error(`Int64 column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
        break;
      case 'float64':
        if (!Number.isFinite(Number(actualStr)) || !Number.isFinite(Number(expectedStr)) ||
            Math.abs(Number(actualStr) - Number(expectedStr)) > 1e-5) throw new Error(`Float64 column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
        break;
      case 'time':
        if (canonicalizeUtcTimeText(actualStr) !== canonicalizeUtcTimeText(expectedStr)) {
          throw new Error(`Time column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
        }
        break;
      case 'string':
        if (cell.column === 'provenance') {
          if (actualStr !== expectedStr) {
            try {
              if (!deepEqualJson(JSON.parse(actualStr), JSON.parse(expectedStr))) throw new Error();
            } catch {
              throw new Error(`Provenance mismatch: expected "${expectedStr}", got "${actualStr}"`);
            }
          }
        } else if (cell.column === 'value_decimal' || cell.column === 'delta_value') {
          if (canonicalizeDecimalText(actualStr) !== canonicalizeDecimalText(expectedStr)) {
            throw new Error(`Decimal column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
          }
        } else if (actualStr !== expectedStr) {
          throw new Error(`String column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
        }
        break;
      default:
        if (actualStr !== expectedStr) throw new Error(`Column "${cell.column}" mismatch: expected "${expectedStr}", got "${actualStr}"`);
    }
  }
}
