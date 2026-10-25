import fs from 'node:fs';
import path from 'node:path';
import { DatabaseSync, type SQLInputValue } from 'node:sqlite';
import { expect } from '@playwright/test';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';

export interface DeviceRecordingFixture {
  backend: string;
  frontend: string;
  simulator: string;
  config_db: string;
  target_db: string;
  workspace_id: string;
  devices: { a: string; b: string };
  mappings: { device_id: string; tag_id: string; target_type: string }[];
}

/** A preseeded, disposable normal gateway is required; never reset an existing workspace. */
export function loadFixture(): DeviceRecordingFixture | undefined {
  const file = process.env.DEVICE_RECORDING_E2E_FIXTURE;
  if (!file) return undefined;
  const fixture = JSON.parse(fs.readFileSync(file, 'utf8')) as DeviceRecordingFixture;
  for (const value of [fixture.backend, fixture.frontend]) {
    const url = new URL(value);
    if (url.protocol !== 'http:' || !['127.0.0.1', 'localhost'].includes(url.hostname)) {
      throw new Error('Only explicit loopback fixture URLs are allowed');
    }
  }
  for (const filePath of [fixture.config_db, fixture.target_db]) {
    if (path.dirname(path.resolve(filePath)) !== path.dirname(path.resolve(file))) {
      throw new Error('SQL databases must belong to the fixture directory');
    }
  }
  if (!fixture.workspace_id || !fixture.devices.a || !fixture.devices.b || fixture.devices.a === fixture.devices.b) {
    throw new Error('Fixture must name its workspace and two distinct devices');
  }
  return fixture;
}

export async function api<T>(fixture: DeviceRecordingFixture, endpoint: string): Promise<T> {
  const response = await fetch(fixture.backend + '/api/v1/datalink/studio-v2/workspace' + endpoint);
  expect(response.ok).toBe(true);
  const body = await response.json() as { success: boolean; data: T };
  expect(body.success).toBe(true);
  return body.data;
}

export async function groups(fixture: DeviceRecordingFixture): Promise<WriteGroup[]> {
  return (await api<{ groups: WriteGroup[] }>(fixture, '/write-groups')).groups;
}

/** SQL values remain strings, including exact uint64s and decimals. */
export function readSQL(file: string, sql: string, args: SQLInputValue[] = []) {
  const db = new DatabaseSync(file, { readOnly: true });
  try {
    return db.prepare(sql).all(...args) as Record<string, string | number | null>[];
  } finally {
    db.close();
  }
}

export function identifier(value: string) {
  if (!/^[a-zA-Z_][a-zA-Z0-9_]*$/.test(value)) throw new Error('Invalid managed SQL identifier');
  return '"' + value + '"';
}

export function catalog(fixture: DeviceRecordingFixture) {
  if (!fs.existsSync(fixture.target_db)) return [];
  return readSQL(fixture.target_db, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name");
}

/** Match the destination row and receipt to one committed outbox effect of the applied version. */
export function committedReadback(fixture: DeviceRecordingFixture, group: WriteGroup) {
  const effects = readSQL(fixture.config_db,
    "SELECT record_id, effect_key, payload_digest, group_revision, connector_revision, table_name FROM wg_delivery_outbox WHERE workspace_id=? AND group_id=? AND group_revision=? AND state='sql_committed' ORDER BY bucket_start LIMIT 1",
    [fixture.workspace_id, group.id, group.applied_revision]);
  if (effects.length === 0) return undefined;
  const effect = effects[0];
  expect(effect.table_name).toBe(group.destination.table_name);
  expect(effect.connector_revision).toBe(group.destination.connector_revision);
  const table = identifier(group.destination.table_name);
  const row = readSQL(fixture.target_db, 'SELECT * FROM ' + table + ' WHERE record_id=?', [effect.record_id])[0];
  expect(row).toBeDefined();
  const receipt = readSQL(fixture.target_db, 'SELECT payload_digest FROM gw_effect_receipts WHERE effect_key=?', [effect.effect_key]);
  expect(receipt).toHaveLength(1);
  expect(receipt[0].payload_digest).toBe(effect.payload_digest);
  const frozen = readSQL(fixture.config_db,
    'SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?',
    [fixture.workspace_id, group.id, group.applied_revision]);
  expect(frozen).toHaveLength(1);
  expect(JSON.parse(String(frozen[0].payload)).applied_revision).toBe(group.applied_revision);
  return { effect, row };
}

export function recordingOperations(fixture: DeviceRecordingFixture) {
  return readSQL(fixture.config_db,
    "SELECT operation_id, status, payload_digest, detail FROM managed_schema_operations WHERE workspace_id=? AND action='recording_start'",
    [fixture.workspace_id]);
}

export function runtimeVersions(fixture: DeviceRecordingFixture) {
  return readSQL(fixture.config_db,
    'SELECT group_id, group_revision, payload FROM wg_runtime_versions WHERE workspace_id=? ORDER BY group_id,group_revision',
    [fixture.workspace_id]);
}
