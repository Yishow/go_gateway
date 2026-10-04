// F PostgreSQL matrix: two Basic devices plus one UI-created multi-entity group.
// Browser actions perform every setup mutation; API and SQL calls below observe.
import assert from 'node:assert/strict';
import { Buffer } from 'node:buffer';
import { writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  DEFAULT_EVIDENCE_DIR, DEFAULT_GATEWAY_BINARY, DEFAULT_SIMULATOR_BINARY, assertFreshPortFree,
  binaryIdentity, makeFreshRun, monotonicNs, observeJSON, openFreshUI, persistFreshMappings,
  previewAndApplyFreshSchema, setupFreshMultiDeviceSources, startFreshGateway,
  removeOwnedWork, startFreshSimulator, stopProcesses, waitForFreshGateway,
} from './fresh-ui.mjs';
import {
  FRESH_EXPECTED_VALUES, FRESH_REGISTERS, FRESH_RULES, FRESH_TARGET_TYPES,
  verifyFreshDeliveryEvidence,
} from './fresh-recording-verification.mjs';
import { createOwnedPostgresSchema } from './destination-preflight.mjs';
import { configureFreshPostgres } from './fresh-postgres-setup.mjs';
import { restoreUniqueDraftColumn, runFreshMatrixTestWrite } from './fresh-matrix-test-write.mjs';
import {
  localDurableRows, ownedPGTarget, pgColumnMetadata, pgCommand, pgRecordingRelations,
  pgRecordingRows, quotePGIdentifier,
} from './fresh-postgres-observation-root.mjs';
import { sh } from './lib.mjs';
import { assertFreshTimingWitness } from './fresh-timing.mjs';
const runID = process.env.F_RUN_ID || ('pg-matrix-' + Date.now());
const gatewayBinary = process.env.F_GATEWAY_BINARY || DEFAULT_GATEWAY_BINARY;
const schema = 'gw_f_pg_matrix_' + Date.now();
const gatewayPort = Number(process.env.F_MATRIX_GATEWAY_PORT || 3385);
const simulatorPortA = Number(process.env.F_MATRIX_SIM_PORT_A || 15057);
const simulatorPortB = Number(process.env.F_MATRIX_SIM_PORT_B || 15058);
const base = 'http://127.0.0.1:' + gatewayPort;
const result = { run_id: runID, kind: 'postgres-dual-basic-multi-entity', passed: false,
  command: 'node scripts/tests/f_device_to_sql/fresh-postgres-matrix-root.mjs',
  environment: 'darwin/arm64; owned loopback simulators and disposable PostgreSQL schema',
  poll_interval_ms: 1000, limitations: [
    'No real PLC, LAN, production database, or deployment acceptance.',
    'Basic A/B are UI-created and activated prerequisites; this run asserts three buckets per entity only for the third multi-entity group.',
  ] };
let run;
let target;
let createdSchema = false;
let ui;
let stage = 'preflight';
const children = [];
function sanitize(value) {
  if (Array.isArray(value)) return value.map(sanitize);
  if (value && typeof value === 'object') return Object.fromEntries(Object.entries(value).map(([key, item]) =>
    [key, /password|credential|authorization|cookie|dsn|token/i.test(key) ? '[redacted]' : sanitize(item)]));
  if (typeof value !== 'string') return value;
  return value.replaceAll(run?.work || '<never-used>', '<owned-run>')
    .replace(/\/(?:Users|home)\/[^/\s]+/g, '<user-home>')
    .replace(/\/(?:private\/)?var\/folders\/[^\s"']+/g, '<owned-temp>')
    .replace(/password=[^\s]+/gi, 'password=[redacted]');
}
function literal(value) { return "'" + String(value).replaceAll("'", "''") + "'"; }
function durableSnapshot(group) {
  const scope = 'group_id=' + literal(group.id) + ' AND group_revision=' + literal(group.applied_revision);
  return {
    buckets: localDurableRows(run.gateway_db, 'SELECT group_id,group_revision,entity_key,bucket_start,kind,reason,record_id,members FROM wg_delivery_buckets WHERE ' + scope + " AND kind='row' ORDER BY bucket_start,entity_key"),
    outbox: localDurableRows(run.gateway_db, 'SELECT effect_key,record_id,group_id,group_revision,entity_key,bucket_start,table_name,state,payload_digest,committed_at FROM wg_delivery_outbox WHERE ' + scope + ' ORDER BY bucket_start,record_id'),
    receipts: localDurableRows(run.gateway_db, 'SELECT effect_key,payload_digest,committed_at FROM wg_delivery_receipts WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE ' + scope + ')'),
    checkpoints: localDurableRows(run.gateway_db, 'SELECT group_id,group_revision,next_close FROM wg_delivery_checkpoints WHERE ' + scope),
  };
}
function bRegisters() {
  const values = [...FRESH_REGISTERS];
  values.splice(2, 4, 0, 0, 0, 2);
  values[0] = 187; values[1] = 777; values[6] = 0;
  values[7] = 0xc010; values[8] = 0;
  const encoded = Buffer.alloc(8);
  encoded.writeDoubleBE(-654.321);
  for (let index = 0; index < 4; index += 1) values[9 + index] = encoded.readUInt16BE(index * 2);
  values.splice(14, 4, 0x8000, 0, 0, 0);
  values.splice(18, 4, 0xffff, 0xffff, 0xffff, 0xffff);
  values[22] = 0x4232;
  for (let index = 23; index < values.length; index += 1) values[index] = 0;
  return values;
}
function registerMap(values) {
  return Object.fromEntries(values.map((value, index) => [index, value]));
}
async function groupsAt(baseURL) {
  return observeJSON(baseURL, '/studio-v2/workspace/write-groups');
}
async function waitGroup(baseURL, id, predicate, timeoutMs = 20_000) {
  const deadline = Date.now() + timeoutMs;
  let current;
  while (Date.now() < deadline) {
    const data = await groupsAt(baseURL);
    current = data.groups?.find((group) => group.id === id);
    if (current && predicate(current)) return current;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error('write group ' + id + ' did not reach expected state: ' + JSON.stringify(current));
}
async function waitBasicGroup(baseURL, deviceID, timeoutMs = 30_000) {
  const deadline = Date.now() + timeoutMs;
  let current;
  while (Date.now() < deadline) {
    const data = await groupsAt(baseURL);
    current = data.groups?.find((group) => group.basic_managed_device_id === deviceID);
    if (current) return current;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  throw new Error('Basic group was not persisted for device ' + deviceID);
}
async function selectBasicDevice(page, deviceID) {
  const selector = page.getByTestId('basic-recording-device');
  await selector.selectOption(deviceID);
  assert.equal(await selector.inputValue(), deviceID);
  await page.waitForTimeout(400);
}
async function openGroup(page, id) {
  const details = page.getByTestId('step4-advanced-recording');
  if ((await details.getAttribute('open')) === null) await details.locator('summary').click();
  await page.getByTestId('group-open-' + id).click();
  await page.getByTestId('group-editor').waitFor({ timeout: 20_000 });
}
async function ensureBasicDevice(page, baseURL, deviceID) {
  const relationsBeforeSave = pgRecordingRelations(schema);
  await selectBasicDevice(page, deviceID);
  await page.getByTestId('basic-recording-start').click();
  await page.getByTestId('basic-recording-status').waitFor({ timeout: 30_000 });
  const first = await waitBasicGroup(baseURL, deviceID);
  if (!first.destination.schema_revision) {
    const relationsBeforeSchema = pgRecordingRelations(schema);
    assert.deepEqual(relationsBeforeSchema, relationsBeforeSave, 'saving a Basic group must not create its recording table');
    await openGroup(page, first.id);
    const schemaResult = await previewAndApplyFreshSchema(page);
    await page.getByTestId('group-editor-close').click();
    await page.getByTestId('basic-recording-panel').waitFor({ timeout: 20_000 });
    return { group: first, schema: schemaResult, relations_before_schema: relationsBeforeSchema };
  }
  return { group: first, schema: null, relations_before_schema: null };
}
async function activateBasicDevice(page, baseURL, deviceID, label) {
  await selectBasicDevice(page, deviceID);
  const tStart = monotonicNs();
  await page.getByTestId('basic-recording-start').click();
  await page.getByTestId('basic-recording-status').waitFor({ timeout: 30_000 });
  await page.waitForFunction(() => /succeeded|成功|partial|部分|failed|失敗/.test(
    document.querySelector('[data-testid="basic-recording-status"]')?.textContent || ''), null, { timeout: 120_000 });
  const status = await page.getByTestId('basic-recording-status').innerText();
  assert.match(status, /succeeded|成功/, label + ' Basic activation did not succeed');
  const data = await groupsAt(baseURL);
  const group = data.groups?.find((item) => item.basic_managed_device_id === deviceID);
  assert.ok(group?.applied_revision, label + ' Basic group was not applied');
  return { t_start: tStart, status, group };
}
async function createSeedGroup(page, tagsA, tagsB, baseURL) {
  const relationsBeforeSave = pgRecordingRelations(schema);
  const details = page.getByTestId('step4-advanced-recording');
  if ((await details.getAttribute('open')) === null) await details.locator('summary').click();
  await page.getByTestId('group-create').click();
  await page.getByTestId('group-name').fill('F PG matrix A+B');
  await page.getByTestId('group-storage-managed').check();
  const advanced = page.getByTestId('group-advanced');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-entity-column').fill('entity_id');
  await page.getByTestId('group-bulk-include').click();
  for (const tag of [...tagsA, ...tagsB]) {
    const row = page.getByTestId('group-member-row-' + tag.persisted_point_id);
    await row.locator('input[type="text"]').first().fill(tag.device_id === tagsA[0].device_id ? 'A' : 'B');
  }
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const data = await groupsAt(baseURL);
  const matches = data.groups?.filter((group) => group.name === 'F PG matrix A+B') || [];
  assert.equal(matches.length, 1);
  const group = matches[0];
  assert.equal(group.members.length, 18, 'advanced matrix group must contain both devices');
  assert.deepEqual(pgRecordingRelations(schema), relationsBeforeSave, 'advanced group save must not create a recording table');
  const schemaResult = await previewAndApplyFreshSchema(page);
  return { group, schema: schemaResult };
}
async function convertSeedToCustom(page, baseURL, group, tagsA, tagsB) {
  await page.getByTestId('group-storage-custom').check();
  const aInt = page.getByTestId('group-member-column-' + tagsA[0].persisted_point_id);
  const bInt = page.getByTestId('group-member-column-' + tagsB[0].persisted_point_id);
  await aInt.waitFor({ timeout: 20_000 });
  const sharedColumn = await aInt.inputValue();
  assert.ok(sharedColumn);
  await bInt.selectOption(sharedColumn);
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const data = await groupsAt(baseURL);
  const saved = data.groups?.find((item) => item.id === group.id);
  assert.ok(saved);
  assert.equal(saved.destination.storage_strategy, 'custom');
  assert.equal(saved.row_policy.entity_key_column, 'entity_id');
  const shared = saved.members.filter((member) => member.target_column === sharedColumn);
  assert.equal(shared.length, 2, 'exactly the A/B int16 members must share one column');
  assert.notEqual(shared[0].entity_key, shared[1].entity_key);
  const targetCounts = new Map();
  for (const member of saved.members) targetCounts.set(member.target_column, (targetCounts.get(member.target_column) || 0) + 1);
  assert.deepEqual([...targetCounts.entries()].filter(([, count]) => count > 1), [[sharedColumn, 2]]);
  return { group: saved, sharedColumn };
}
async function disableGroup(page, baseURL, group) {
  await openGroup(page, group.id);
  await page.getByTestId('group-lifecycle').waitFor({ timeout: 20_000 });
  const current = (await groupsAt(baseURL)).groups?.find((item) => item.id === group.id);
  if (current?.status !== 'disabled') {
    await page.getByTestId('group-disable').click();
    return waitGroup(baseURL, group.id, (item) => item.status === 'disabled').then(async (disabled) => {
      await page.getByTestId('group-editor-close').click();
      return disabled;
    });
  }
  await page.getByTestId('group-editor-close').click();
  return current;
}
function memberIdentity(member) { return JSON.stringify([member.device_id, member.point_id, member.tag_id]); }
function utcMillis(value) {
  const text = String(value ?? '');
  if (!/(?:Z|\+00:00)$/i.test(text)) return NaN;
  return Date.parse(text);
}
function sameBucket(left, right) {
  const first = utcMillis(left);
  const second = utcMillis(right);
  return Number.isFinite(first) && first === second;
}
function consecutiveBucketStarts(rows) {
  const values = [...new Set(rows.map((row) => row.bucket_start))].sort();
  const times = values.map((value) => utcMillis(value));
  return { values, times, valid: values.length === 3 && times.every(Number.isFinite) &&
    times[1] - times[0] === 60_000 && times[2] - times[1] === 60_000 };
}
function verifyMatrix({ group, rows, durable, mappings, sharedColumn, expectedByIdentity }) {
  const failures = [];
  const members = group.members || [];
  const entities = ['A', 'B'];
  if (members.length !== 18) failures.push('member count ' + members.length);
  if (group.destination.storage_strategy !== 'custom') failures.push('advanced group did not remain custom');
  if (group.row_policy.entity_key_column !== 'entity_id') failures.push('entity key column was not persisted');
  if (new Set(members.map((member) => member.entity_key)).size !== 2) failures.push('entity keys are not A/B');
  const mappingIdentities = new Set(mappings.filter((mapping) => mapping.point_id && mapping.tag_id && mapping.device_id)
    .map((mapping) => JSON.stringify([mapping.device_id, mapping.point_id, mapping.tag_id])));
  if (mappingIdentities.size !== mappings.length) failures.push('mapping identity incomplete');
  if (!['A', 'B'].every((entity) => consecutiveBucketStarts(rows.filter((row) => row.entity_id === entity)).values
    .every((start) => rows.some((row) => row.entity_id !== entity && sameBucket(row.bucket_start, start))))) {
    failures.push('A/B entities do not share the same three buckets');
  }
  for (const entity of entities) {
    const entityRows = rows.filter((row) => row.entity_id === entity);
    const buckets = consecutiveBucketStarts(entityRows);
    if (entityRows.length !== 3 || !buckets.valid) failures.push(entity + ' does not have three consecutive buckets');
    const entityMembers = members.filter((member) => member.entity_key === entity);
    const otherMembers = members.filter((member) => member.entity_key !== entity);
    const entityMappings = mappings.filter((mapping) => entityMembers.some((member) =>
      memberIdentity(member) === JSON.stringify([mapping.device_id, mapping.point_id, mapping.tag_id])));
    const entityDevices = new Set(entityMembers.map((member) => member.device_id));
    if (entityDevices.size !== 1) failures.push(entity + ' members do not resolve to one device');
    if (entityDevices.size === 1) {
      // Cross-device layouts intentionally omit device_id from the SQL row. The
      // entity's persisted member identities prove the one device used here;
      // normalize only the shared verifier input while retaining raw rows below.
      const deviceID = [...entityDevices][0];
      const expectedByTag = Object.fromEntries(entityMembers.map((member) => {
        const mapping = entityMappings.find((item) => memberIdentity(member) ===
          JSON.stringify([item.device_id, item.point_id, item.tag_id]));
        return mapping?.tag_key ? [mapping.tag_key, expectedByIdentity.get(memberIdentity(member))] : null;
      }).filter(Boolean));
      const sharedVerification = verifyFreshDeliveryEvidence({
        group: { ...group, members: entityMembers },
        table: group.destination.table_name,
        rows: entityRows.map((row) => ({ ...row, device_id: deviceID })),
        mappings: entityMappings,
        expectedByTag,
        closedBuckets: durable.buckets.filter((bucket) => bucket.entity_key === entity),
        outbox: durable.outbox.filter((item) => item.entity_key === entity),
        receipts: durable.receipts,
        sharedTargetColumns: [sharedColumn],
      });
      failures.push(...sharedVerification.failures.map((failure) => entity + ': ' + failure));
    }
    for (const row of entityRows) {
      if (row.group_id !== group.id || row.entity_id !== entity) failures.push(entity + ' row identity mismatch');
      const outbox = durable.outbox.find((item) => item.record_id === row.record_id && sameBucket(item.bucket_start, row.bucket_start));
      const receipt = outbox && durable.receipts.find((item) => item.effect_key === outbox.effect_key);
      if (!outbox || outbox.state !== 'sql_committed' || outbox.group_revision !== group.applied_revision ||
        outbox.table_name !== group.destination.table_name) failures.push(entity + ' row lacks current committed outbox');
      if (!receipt || receipt.payload_digest !== outbox.payload_digest || !receipt.committed_at) failures.push(entity + ' row lacks matching receipt');
      for (const member of entityMembers) {
        const expected = expectedByIdentity.get(memberIdentity(member));
        if (expected === undefined || row[member.target_column] === null || String(row[member.target_column]) !== expected) failures.push(entity + ' value mismatch for ' + member.target_column);
      }
      for (const member of otherMembers) {
        if (member.target_column !== sharedColumn && row[member.target_column] !== null && row[member.target_column] !== '') failures.push(entity + ' row leaked ' + member.target_column);
      }
      let provenance;
      try { provenance = JSON.parse(row.provenance); } catch { provenance = null; }
      const expectedProvenance = new Set(entityMembers.map(memberIdentity));
      const actualProvenance = new Set((provenance || []).map((entry) => entry?.member));
      const bucketTime = utcMillis(row.bucket_start);
      if (!Array.isArray(provenance) || actualProvenance.size !== expectedProvenance.size ||
        actualProvenance.size !== provenance.length ||
        [...expectedProvenance].some((identity) => !actualProvenance.has(identity)) ||
        provenance.some((entry) => entry?.status !== 'ok' || entry?.quality !== 'good' || !entry?.sample_id ||
          !Number.isFinite(utcMillis(entry?.observed_at)) || utcMillis(entry.observed_at) < bucketTime ||
          utcMillis(entry.observed_at) >= bucketTime + 60_000)) {
        failures.push(entity + ' provenance or quality mismatch');
      }
    }
  }
  if (rows.length !== 6) failures.push('expected six SQL rows, got ' + rows.length);
  if (durable.buckets.length !== 6) failures.push('expected six durable row closures, got ' + durable.buckets.length);
  if (durable.checkpoints.length !== 1) failures.push('expected one current checkpoint, got ' + durable.checkpoints.length);
  for (const member of members) {
    if (!member.source_revision || !member.mapping_revision || !member.target_column) failures.push('member revision/column missing: ' + member.point_id);
    if (!mappingIdentities.has(memberIdentity(member))) failures.push('member identity missing: ' + member.point_id);
  }
  return { failures, entities, member_count: members.length, sql_row_count: rows.length,
    durable_row_count: durable.buckets.length, mapping_count: mappings.length };
}
try {
  target = ownedPGTarget(schema);
  run = makeFreshRun(runID);
  result.destination = { host: target.host, port: target.port, database: target.database,
    username: target.username, schema: target.schema, table: target.table };
  result.source_head = sh('git', ['rev-parse', 'HEAD']);
  result.binary_sha256 = binaryIdentity(gatewayBinary).sha256;
  result.simulator_sha256 = binaryIdentity(DEFAULT_SIMULATOR_BINARY).sha256;
  const dependencies = ['fresh-postgres-matrix-root.mjs', 'fresh-postgres-observation-root.mjs',
    'fresh-postgres-setup.mjs', 'fresh-matrix-test-write.mjs', 'fresh-recording-verification.mjs',
    'fresh-ui.mjs', 'destination-preflight.mjs', 'lib.mjs'];
  result.harness_source_sha256 = Object.fromEntries(dependencies.map((name) =>
    [name, binaryIdentity(new URL(name, import.meta.url).pathname).sha256]));
  result.target_preflight = createOwnedPostgresSchema({ psql: pgCommand, schema, table: target.table });
  createdSchema = true;
  result.relations_before_ui = pgRecordingRelations(schema);
  assert.deepEqual(result.relations_before_ui, [], 'fresh schema must contain no recording objects');
  await assertFreshPortFree(gatewayPort); await assertFreshPortFree(simulatorPortA); await assertFreshPortFree(simulatorPortB);
  children.push(startFreshSimulator({ work: run.work, name: 'pg-matrix-a', port: simulatorPortA, registers: registerMap(FRESH_REGISTERS) }));
  children.push(startFreshSimulator({ work: run.work, name: 'pg-matrix-b', port: simulatorPortB, registers: registerMap(bRegisters()) }));
  children.push(startFreshGateway({ work: run.work, binary: gatewayBinary, port: gatewayPort }));
  await waitForFreshGateway(base);
  ui = await openFreshUI({ port: gatewayPort });
  result.browser_version = ui.browser.version(); result.t_open = ui.t_open;
  stage = 'devices-and-mappings';
  await setupFreshMultiDeviceSources(ui.page, [
    { name: 'F PostgreSQL A', port: simulatorPortA, prefix: 'A' },
    { name: 'F PostgreSQL B', port: simulatorPortB, prefix: 'B' },
  ], FRESH_RULES);
  const targetTypeByAddress = new Map(FRESH_RULES.map((rule, index) => [String(rule.start), FRESH_TARGET_TYPES[index]]));
  result.tags = await persistFreshMappings(ui.page, base, { tagPrefix: 'pg_matrix',
    targetTypeFor: ({ address }) => targetTypeByAddress.get(address) });
  assert.equal(result.tags.length, 18);
  result.devices = await observeJSON(base, '/studio-v2/workspace/devices');
  const devices = Array.isArray(result.devices) ? result.devices : result.devices.devices;
  const deviceA = devices.find((device) => device.name === 'F PostgreSQL A');
  const deviceB = devices.find((device) => device.name === 'F PostgreSQL B');
  assert.ok(deviceA?.id && deviceB?.id && deviceA.id !== deviceB.id);
  const tagsA = result.tags.filter((tag) => tag.device_id === deviceA.id).sort((left, right) => Number(left.address) - Number(right.address));
  const tagsB = result.tags.filter((tag) => tag.device_id === deviceB.id).sort((left, right) => Number(left.address) - Number(right.address));
  assert.equal(tagsA.length, 9); assert.equal(tagsB.length, 9);
  assert.deepEqual(tagsA.map((tag) => tag.address), tagsB.map((tag) => tag.address), 'A/B source rules must share the nine addresses');
  result.tag_identities = { A: tagsA, B: tagsB };
  const savedMappings = await observeJSON(base, '/studio-v2/workspace/mappings');
  const mappingRows = Array.isArray(savedMappings) ? savedMappings : savedMappings.mappings || [];
  result.mapping_target_types = {};
  for (const [entity, tags] of [['A', tagsA], ['B', tagsB]]) {
    result.mapping_target_types[entity] = tags.map((tag, index) => {
      const mapping = mappingRows.find((item) => item.point_id === tag.persisted_point_id && item.tag_id === tag.tag_id && item.device_id === tag.device_id);
      assert.equal(mapping?.target_type, FRESH_TARGET_TYPES[index], entity + ' mapping target type was not persisted at address ' + tag.address);
      return { address: tag.address, target_type: mapping.target_type };
    });
  }
  stage = 'destination';
  result.connector_setup = await configureFreshPostgres(ui.page, base, target);
  assert.equal(await ui.page.getByTestId('basic-recording-interval').inputValue(), '60');
  stage = 'dual-basic';
  result.relations_before_schema = pgRecordingRelations(schema);
  assert.deepEqual(result.relations_before_schema, [], 'connector and saved groups must not pre-create recording tables');
  const basicA = await ensureBasicDevice(ui.page, base, deviceA.id);
  const basicAStarted = await activateBasicDevice(ui.page, base, deviceA.id, 'A');
  const basicB = await ensureBasicDevice(ui.page, base, deviceB.id);
  const basicBStarted = await activateBasicDevice(ui.page, base, deviceB.id, 'B');
  result.basic = { A: { ...basicAStarted, schema: basicA.schema }, B: { ...basicBStarted, schema: basicB.schema } };
  assert.ok(result.basic.A.group.applied_revision && result.basic.B.group.applied_revision);
  stage = 'advanced-entity-group';
  const seed = await createSeedGroup(ui.page, tagsA, tagsB, base);
  const originalBIntColumn = seed.group.members.find((member) => member.point_id === tagsB[0].persisted_point_id)?.target_column;
  assert.ok(originalBIntColumn, 'seed must retain B int16 original column');
  const converted = await convertSeedToCustom(ui.page, base, seed.group, tagsA, tagsB);
  const groupList = await groupsAt(base);
  const advancedGroups = (groupList.groups || []).filter((item) => !item.basic_managed_device_id && item.status !== 'deleted');
  assert.equal(advancedGroups.length, 1, 'fresh matrix must have one non-Basic third group');
  assert.equal(advancedGroups[0].id, converted.group.id);
  result.group_counts = { total: (groupList.groups || []).length, basic: 2, advanced: advancedGroups.length };
  result.advanced_schema = seed.schema;
  result.shared_column = converted.sharedColumn;
  await ui.page.getByTestId('group-editor-close').click();
  const disabledA = await disableGroup(ui.page, base, basicAStarted.group);
  const disabledB = await disableGroup(ui.page, base, basicBStarted.group);
  result.basic_disabled = { A: disabledA, B: disabledB };
  await openGroup(ui.page, converted.group.id);
  await ui.page.getByTestId('group-apply').waitFor({ state: 'visible', timeout: 20_000 });
  await ui.page.waitForFunction(() => !document.querySelector('[data-testid="group-apply"]')?.hasAttribute('disabled'), null, { timeout: 30_000 });
  const tStart = monotonicNs();
  await ui.page.getByTestId('group-apply').click();
  const applied = await waitGroup(base, converted.group.id, (group) => Boolean(group.applied_revision) && group.status !== 'draft', 30_000);
  result.t_start = tStart; result.group = applied;
  await ui.page.getByTestId('group-editor-close').click();
  result.columns = pgColumnMetadata(schema, applied.destination.table_name);
  assert.ok(result.columns.some((column) => column.name === 'entity_id'));
  stage = 'three-formal-buckets-per-entity';
  const expectedByIdentity = new Map();
  // PostgreSQL's SELECT ...::text renders the boolean target as true/false;
  // keep the exact textual values used by the SQL witness rather than the
  // simulator's 1/0 register representation.
  const expectedB = ['187', '777', '2', 'false', '-2.25', '-654.321', '9223372036854775808', '18446744073709551615', 'B2'];
  const expectedA = [...FRESH_EXPECTED_VALUES]; expectedA[3] = 'true';
  for (const pair of [[tagsA, expectedA], [tagsB, expectedB]]) {
    pair[0].forEach((tag, index) => expectedByIdentity.set(JSON.stringify([tag.device_id, tag.persisted_point_id, tag.tag_id]), pair[1][index]));
  }
  const deadline = Date.now() + 270_000; result.polls = []; result.poll_count = 0;
  while (true) {
    result.durable = durableSnapshot(applied);
    if (!result.t_closed && result.durable.buckets.length > 0) result.t_closed = monotonicNs();
    result.rows = pgRecordingRows(schema, applied.destination.table_name, result.columns);
    const observed = monotonicNs();
    if (!result.t_sql && result.rows.length > 0) result.t_sql = observed;
    result.poll_count += 1;
    const last = result.polls.at(-1);
    if (!last || last.local_closed_rows !== result.durable.buckets.length || last.sql_rows !== result.rows.length) {
      result.polls.push({ observed_at_ns: observed, local_closed_rows: result.durable.buckets.length, sql_rows: result.rows.length });
    }
    const threeEach = ['A', 'B'].every((entity) => consecutiveBucketStarts(result.rows.filter((row) => row.entity_id === entity)).valid);
    const settled = result.rows.length >= 6 && result.rows.every((row) => {
      const outbox = result.durable.outbox.find((item) => item.record_id === row.record_id && sameBucket(item.bucket_start, row.bucket_start)); const receipt = outbox && result.durable.receipts.find((item) => item.effect_key === outbox.effect_key);
      return outbox?.state === 'sql_committed' && outbox.group_revision === applied.applied_revision && outbox.table_name === applied.destination.table_name &&
        receipt?.payload_digest === outbox.payload_digest && Boolean(receipt?.committed_at);
    });
    if (result.durable.buckets.length >= 6 && threeEach && settled) break;
    if (Date.now() >= deadline) throw new Error('dual-entity group did not produce three consecutive buckets per entity');
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  result.mappings = await observeJSON(base, '/studio-v2/workspace/mappings');
  result.verification = verifyMatrix({ group: applied, rows: result.rows, durable: result.durable,
    mappings: Array.isArray(result.mappings) ? result.mappings : result.mappings.mappings || [],
    sharedColumn: result.shared_column, expectedByIdentity });
  assert.deepEqual(result.verification.failures, []);
  assert.equal(result.rows.length, 6); assert.equal(result.durable.buckets.length, 6);
  stage = 'test-write-draft';
  const draftGroup = await restoreUniqueDraftColumn({ page: ui.page, baseURL: base, group: applied,
    pointID: tagsB[0].persisted_point_id, originalColumn: originalBIntColumn,
    readRows: () => pgRecordingRows(schema, applied.destination.table_name, result.columns) });
  result.draft_group = draftGroup;
  result.test_write = await runFreshMatrixTestWrite({ page: ui.page, baseURL: base, group: draftGroup,
    readRows: () => pgRecordingRows(schema, applied.destination.table_name, result.columns),
    readDelivery: () => durableSnapshot(applied),
    readLedger: (operationID) => localDurableRows(run.gateway_db, 'SELECT operation_id,action,detail FROM managed_schema_operations WHERE operation_id=' + literal(operationID)) });
  result.delivery = await observeJSON(base, '/studio-v2/workspace/write-groups/' + applied.id + '/delivery');
  result.setup_ms = Number(BigInt(result.t_start) - BigInt(result.t_open)) / 1e6;
  result.system_wait_ms = Number(BigInt(result.t_closed) - BigInt(result.t_start)) / 1e6;
  result.delivery_observed_ms = Number(BigInt(result.t_sql) - BigInt(result.t_closed)) / 1e6;
  result.end_to_end_ms = Number(BigInt(result.t_sql) - BigInt(result.t_open)) / 1e6;
  Object.assign(result, assertFreshTimingWitness(result));
  assert.deepEqual(ui.pageErrors, []);
  result.page_errors = ui.pageErrors;
  result.passed = true;
} catch (error) {
  result.failed_stage = stage; result.error = String(error.message).split('\n')[0];
  if (ui) result.ui_text = (await ui.page.locator('body').innerText().catch(() => '')).slice(-10_000);
} finally {
  try { await ui?.browser.close(); result.browser_cleanup = 'closed'; }
  catch { result.browser_cleanup = 'failed'; result.passed = false; }
  result.child_exit = await stopProcesses(children);
  result.child_exit_passed = result.child_exit.every((child) => child.state === 'exited' && child.exit_code === 0 &&
    !child.signal && !child.spawn_error && !child.stop_error && !child.kill_sent);
  if (!result.child_exit_passed) result.passed = false;
  try {
    if (createdSchema) pgCommand('DROP SCHEMA ' + quotePGIdentifier(schema) + ' CASCADE');
    result.schema_cleanup = createdSchema ? !pgCommand("SELECT nspname FROM pg_namespace WHERE nspname='" + schema + "'") : 'not-created';
    if (result.schema_cleanup === false) result.passed = false;
  } catch { result.schema_cleanup = false; result.passed = false; }
  if (run) {
    const workCleanup = removeOwnedWork(run.work);
    result.owned_work_cleanup = workCleanup.removed;
    if (workCleanup.error) result.owned_work_cleanup_error = workCleanup.error;
    if (!result.owned_work_cleanup) result.passed = false;
  }
  const fields = Object.entries(sanitize(result)).map(([key, value]) => '  ' + JSON.stringify(key) + ': ' + JSON.stringify(value));
  writeFileSync(join(DEFAULT_EVIDENCE_DIR, runID + '.json'), '{\n' + fields.join(',\n') + '\n}\n');
  console.log(JSON.stringify({ run_id: runID, passed: result.passed, failed_stage: result.failed_stage, error: sanitize(result.error) }));
}
process.exitCode = result.passed ? 0 : 1;
