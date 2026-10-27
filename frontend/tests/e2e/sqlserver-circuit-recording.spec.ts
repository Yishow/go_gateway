import { expect, test } from '@playwright/test';
import fs from 'node:fs';
import { createHash } from 'node:crypto';
import type { RecordingStartRequest, RecordingStartOperation } from '../../src/types/studioV2RecordingStart';
import type { WriteGroupTestWriteOperation, WriteGroupTestWritePreview } from '../../src/types/studioV2WriteGroupTestWrite';
import {
  api, buildSQLServerRowSelectSql, deliveryCommittedOutbox, groups, loadSQLServerFixture,
  parseWireRowPayload, querySQLServerReadback, readConfigSQL, recordingStartOperations, resolveSqlCmdBinary,
  runtimeVersions, schemaApplyOperations, sqlserverIdentifier, validateConnectionConfigFile,
  validateSQLServerCapability, verifyWirePayloadAgainstRow, type SQLServerFixtureManifest,
} from './fixtures/sqlserver-circuit-recording';

const fixture = loadSQLServerFixture();

for (const storageMode of ['plain', 'native_monthly'] as const) {
  test(`SQL Server circuit recording acceptance: ${storageMode} mode (dual profiles, schema, write/readback, screenshots)`, async ({ page }, info) => {
    // Truthful skip when no disposable fixture is configured
    test.skip(!fixture, 'NOT_RUN: GW_SQLSERVER_TEST_FIXTURE required (no disposable SQL Server instance configured)');
    if (!fixture) throw new Error('Fixture required');

    // Only matching storage mode runs; the other is NOT_RUN (separate fresh manifest required per mode)
    test.skip(
      fixture.browser_storage_mode !== storageMode,
      `NOT_RUN: fixture is configured for ${fixture.browser_storage_mode}, not ${storageMode} (separate fresh manifest required per mode)`,
    );

    // Mandatory SQL readback requires sqlcmd binary on host
    const sqlcmd = resolveSqlCmdBinary();
    test.skip(!sqlcmd, 'NOT_RUN: fixed sqlcmd binary is required for mandatory host-side readback');

    test.setTimeout(300_000);
    page.setDefaultTimeout(15_000);
    page.setDefaultNavigationTimeout(15_000);
    const f: SQLServerFixtureManifest = fixture;
    const ui = (id: string) => page.getByTestId(id);
    const tile = (id: string) => page.locator(`[data-testid="basic-recording-device-tile"][data-device-id="${id}"]`);

    const open = async () => {
      await page.goto(`${f.frontend}/studio/v2?step=4`);
      await expect(ui('circuit-profile-basic-panel')).toBeVisible();
    };

    const screenshot = async (name: string) => {
      const file = info.outputPath(`${name}.png`);
      await page.screenshot({ path: file, fullPage: true });
      await info.attach(name, { path: file, contentType: 'image/png' });
    };

    const attachJSON = async (name: string, value: unknown) => {
      const file = info.outputPath(`${name}.json`);
      fs.writeFileSync(file, JSON.stringify(value, null, 2));
      await info.attach(name, { path: file, contentType: 'application/json' });
    };

    const noOverflow = async () => {
      expect(await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)).toBe(false);
    };

    // Pre-condition validations: fresh, unused workspace required for each mode
    expect((await api<{ id: string }>(f, '')).id).toBe(f.workspace_id);
    const devices = await api<{ id: string; name: string }[]>(f, '/devices');
    expect(devices.map((d) => d.id).sort()).toEqual([f.devices.a, f.devices.b].sort());

    const initialGroups = await groups(f);
    if (initialGroups.length > 0) {
      throw new Error(`Workspace ${f.workspace_id} already contains write groups; fresh workspace required for ${storageMode} test`);
    }
    const initialOps = recordingStartOperations(f);
    if (initialOps.length > 0) {
      throw new Error(`Workspace ${f.workspace_id} has already executed recording operations; fresh workspace required for ${storageMode} test`);
    }

    // Check actual capability facts BEFORE any UI mutations; closed gates => NOT_RUN
    const dbConfig = await api<{ id?: string; connector_id?: string; identity_revision?: string; connector_revision?: string }>(f, '/database-config');
    const connectorId = dbConfig.connector_id ?? dbConfig.id;
    const connectorRevision = dbConfig.connector_revision ?? dbConfig.identity_revision;
    if (!connectorId || !connectorRevision) {
      test.skip(true, 'NOT_RUN: database-config has no saved connector');
      return;
    }
    const capData = await api<unknown>(
      f,
      `/recording-plans/capabilities?kind=sqlserver&connector_id=${encodeURIComponent(connectorId)}&expected_connector_revision=${encodeURIComponent(connectorRevision)}`,
    );
    const cap = validateSQLServerCapability(capData, connectorId, connectorRevision);
    if (
      !cap.connected ||
      cap.metadata_status !== 'verified' ||
      !cap.supports_managed_long ||
      !cap.supports_group_test_write ||
      !cap.supports_missing_repair ||
      (storageMode === 'native_monthly' && !cap.supports_native_partition)
    ) {
      test.skip(true, 'NOT_RUN: SQL Server recording operation gates are closed on this instance');
    }

    const mutations: string[] = [];
    page.on('request', (request) => {
      if (['POST', 'PUT', 'PATCH', 'DELETE'].includes(request.method()) && request.url().includes('/api/')) {
        mutations.push(`${request.method()} ${new URL(request.url()).pathname}`);
      }
    });

    await page.addInitScript(() => localStorage.setItem('i18nextLng', 'zh-TW'));
    await open();

    // 1. Dual resolution viewport verification
    for (const viewport of [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }]) {
      await page.setViewportSize(viewport);
      await expect(ui('circuit-profile-sqlserver-capability')).toBeVisible();

      // Verify exact SQL Server label without guessing IDs
      const nativeOption = ui('circuit-profile-storage-native-monthly');
      await expect(nativeOption).toBeVisible();
      const nativeLabel = nativeOption.locator('xpath=..');
      await expect(nativeLabel).toContainText('SQL Server');
      await expect(nativeLabel).not.toContainText('MySQL');

      await tile(f.devices.a).click();
      await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue(viewport.width === 1024 ? '15' : '47');
      await ui('circuit-profile-interval-instantaneous').fill('47');

      await tile(f.devices.a).focus();
      await page.keyboard.press('Tab');
      await expect(tile(f.devices.b)).toBeFocused();
      await page.keyboard.press('Space');
      await expect(tile(f.devices.b)).toHaveAttribute('aria-pressed', 'true');
      await ui('circuit-profile-interval-instantaneous').fill('33');

      await tile(f.devices.b).focus();
      await page.keyboard.press('Shift+Tab');
      await expect(tile(f.devices.a)).toBeFocused();
      await page.keyboard.press('Enter');
      await expect(tile(f.devices.a)).toHaveAttribute('aria-pressed', 'true');
      await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('47');

      // Reload preservation check
      await page.reload();
      await expect(tile(f.devices.a)).toHaveAttribute('aria-pressed', 'true');
      await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('47');
      await tile(f.devices.b).click();
      await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('33');
      await tile(f.devices.a).click();
      await noOverflow();
      await screenshot(`sqlserver-${storageMode}-drafts-${viewport.width}`);
    }

    expect(mutations).toEqual([]);
    expect(await groups(f)).toHaveLength(0);

    // 2. Select configured storage mode & missing repair option
    await ui(storageMode === 'plain' ? 'circuit-profile-storage-plain' : 'circuit-profile-storage-native-monthly').check();
    await ui('circuit-profile-recreate-missing').check();

    // 3. Configure Dual Profiles (Instantaneous + Cumulative with counter confirmation)
    await ui('circuit-profile-circuit-name').fill(`電錶 ${storageMode.toUpperCase()}`);
    await ui('circuit-profile-interval-instantaneous').fill('15');

    await ui('circuit-profile-sources-toggle').click();
    const counter = f.mappings.find((m) => m.device_id === f.devices.a && m.target_type === 'uint64');
    expect(counter).toBeDefined();
    const counterId = counter!.tag_id;
    await ui(`circuit-profile-role-${counterId}`).selectOption('counter_raw');
    await ui('circuit-profile-sources-toggle').click();

    await expect(ui('circuit-counter-confirmation')).toBeVisible();
    await expect(ui('circuit-profile-prepare')).toBeDisabled();
    await ui('circuit-confirm-counter-semantics').selectOption('counter');
    await ui('circuit-confirm-counter-unit').fill('kWh');
    await ui('circuit-confirm-counter-submit').click();
    await expect(ui('circuit-counter-confirmation')).toHaveCount(0);

    await ui('circuit-profile-toggle-cumulative').check();
    await expect(ui('circuit-profile-interval-cumulative')).toHaveValue('10');
    await ui('circuit-profile-interval-cumulative').fill('1');

    await ui('circuit-profile-sources-toggle').click();
    for (const mapping of f.mappings.filter((m) => m.device_id === f.devices.a && m.tag_id !== counterId)) {
      const box = ui(`circuit-profile-member-cum-${mapping.tag_id}`);
      await box.evaluate((el) => el.scrollIntoView({ block: 'center', inline: 'nearest' }));
      await box.uncheck();
    }
    await expect(ui(`circuit-profile-member-cum-${counterId}`)).toBeChecked();
    await ui('circuit-profile-sources-toggle').click();

    // 4. Schema Preview, Confirm, Repair & Error Recovery
    await expect(ui('circuit-profile-prepare')).toBeEnabled();
    await ui('circuit-profile-prepare').click();
    await expect(ui('group-schema-preview-result')).toHaveCount(2);

    // In native mode: verify aligned partition scheme and function in advanced SQL preview
    if (storageMode === 'native_monthly') {
      const advancedSql = await ui('circuit-prepared-group-instantaneous').locator('pre').textContent();
      expect(advancedSql).toContain('CREATE PARTITION FUNCTION');
      expect(advancedSql).toContain('CREATE PARTITION SCHEME');
    }

    // Confirm instantaneous with lost-reply retry recovery seam
    const instSchemaPanel = ui('circuit-prepared-group-instantaneous');
    let schemaApplyCalled = false;
    await page.route('**/schema-apply', async (route) => {
      if (!schemaApplyCalled) {
        schemaApplyCalled = true;
        const resp = await route.fetch();
        expect(resp.status()).toBe(200);
        await route.abort('failed'); // lost HTTP reply after commit
        return;
      }
      await route.continue();
    }, { times: 1 });

    await instSchemaPanel.getByTestId('group-schema-apply').click();
    await expect(instSchemaPanel.getByTestId('group-schema-operation-pending')).toBeVisible();
    await instSchemaPanel.getByTestId('group-schema-check').click();

    // Confirm cumulative schema normally
    const cumSchemaPanel = ui('circuit-prepared-group-cumulative');
    const cumApplyResponse = page.waitForResponse((r) => r.url().endsWith('/schema-apply') && r.request().method() === 'POST');
    await cumSchemaPanel.getByTestId('group-schema-apply').click();
    expect((await cumApplyResponse).status()).toBe(200);

    await expect.poll(async () =>
      (await groups(f)).every((g) => Boolean(g.destination.schema_digest)),
    ).toBe(true);

    const schemaOps = schemaApplyOperations(f);
    expect(schemaOps).toHaveLength(2);

    // 5. Positive Group Test Write Lifecycle (verified + cleaned + no residuals)
    const instGroupPanel = ui('circuit-prepared-group-instantaneous');
    const testWritePreviewBtn = instGroupPanel.getByTestId('group-test-write-preview');
    const connConfig = validateConnectionConfigFile(f.connection_config_file, f.database, f.schema);

    // Positive test requires verified capability: preview enabled, full execution
    await expect(testWritePreviewBtn).toBeEnabled();
    const previewResponse = page.waitForResponse((r) => r.request().method() === 'POST' && r.url().endsWith('/test-write-preview'));
    await testWritePreviewBtn.click();
    const testPreview = ((await (await previewResponse).json()) as { data: WriteGroupTestWritePreview }).data;
    expect(testPreview.record_ids?.length).toBeGreaterThan(0);
    for (const id of testPreview.record_ids ?? []) expect(id).toMatch(/^[0-9a-f]{64}$/);
    const previewCard = instGroupPanel.getByTestId('group-test-write-preview-card');
    await expect(previewCard).toBeVisible();
    await expect(instGroupPanel.getByTestId('group-test-write-row-count')).toBeVisible();

    const confirmBtn = instGroupPanel.getByTestId('group-test-write-confirm');
    const confirmationResponse = page.waitForResponse((r) => r.request().method() === 'POST' && r.url().endsWith('/test-write'));
    await confirmBtn.click();
    const testOperation = ((await (await confirmationResponse).json()) as { data: WriteGroupTestWriteOperation }).data;
    expect(testOperation.operation_id).toBe(testPreview.operation_id);
    expect(testOperation.write_outcome).toBe('written_verified');
    expect(testOperation.cleanup_status).toBe('cleaned');

    // Positive test requires verified write outcome and cleaned status
    await expect(instGroupPanel.getByTestId('group-test-write-outcome')).toBeVisible();
    await expect(instGroupPanel.getByTestId('group-test-write-outcome')).toContainText('已寫入且讀回正確');

    await expect(instGroupPanel.getByTestId('group-test-write-cleanup')).toBeVisible();
    await expect(instGroupPanel.getByTestId('group-test-write-cleanup')).toContainText('測試列已移除');

    // Verify test-write production cleanup leaves zero residuals in target table
    const targetTableInstant = sqlserverIdentifier(f.schema, testPreview.target.table);
    const testLedger = readConfigSQL(f.config_db,
      "SELECT detail FROM managed_schema_operations WHERE workspace_id=? AND operation_id=? AND action='test_write'",
      [f.workspace_id, testPreview.operation_id]);
    expect(testLedger).toHaveLength(1);
    const progress = JSON.parse(String(testLedger[0].detail)) as { record_ids: string[]; effect_keys: string[] };
    expect(progress.record_ids).toEqual(testPreview.record_ids);
    expect(progress.effect_keys).toHaveLength(progress.record_ids.length);
    for (const key of progress.effect_keys) expect(key).toMatch(/^[0-9a-f]{64}$/);
    const recordIDsSQL = progress.record_ids.map((id) => `'${id}'`).join(',');
    const effectKeysSQL = progress.effect_keys.map((key) => `'${key}'`).join(',');
    const residualCheck = querySQLServerReadback(
      connConfig,
      `SELECT CONVERT(NVARCHAR(30),COUNT_BIG(*)) AS remaining FROM ${targetTableInstant} WHERE record_id IN (${recordIDsSQL}) FOR JSON PATH,WITHOUT_ARRAY_WRAPPER`,
    );
    expect((JSON.parse(residualCheck.trim()) as { remaining: string }).remaining).toBe('0');
    const receiptResidual = querySQLServerReadback(connConfig,
      `SELECT CONVERT(NVARCHAR(30),COUNT_BIG(*)) AS remaining FROM ${sqlserverIdentifier(f.schema, 'gw_effect_receipts')} WHERE effect_key IN (${effectKeysSQL}) FOR JSON PATH,WITHOUT_ARRAY_WRAPPER`);
    expect((JSON.parse(receiptResidual.trim()) as { remaining: string }).remaining).toBe('0');

    // 6. Recording Start with Lost HTTP Reply Recovery
    let startRequest: RecordingStartRequest | undefined;
    let startCommitted: RecordingStartOperation | undefined;
    await page.route('**/studio-v2/workspace/recording-start', async (route) => {
      startRequest = route.request().postDataJSON() as RecordingStartRequest;
      expect(startRequest.request_id).toBeTruthy();
      expect(startRequest.device_ids).toEqual([f.devices.a]);
      const response = await route.fetch();
      expect(response.status()).toBe(200);
      startCommitted = ((await response.json()) as { data: RecordingStartOperation }).data;
      expect(startCommitted.status).toBe('succeeded');
      await route.abort('failed'); // simulate lost HTTP reply
    }, { times: 1 });

    await ui('circuit-profile-start').click();
    await expect(ui('basic-recording-retry')).toBeVisible();
    await expect(ui('circuit-profile-start')).toBeDisabled();
    expect(startRequest).toBeDefined();

    const resent = page.waitForRequest((r) => r.url().endsWith('/recording-start') && r.method() === 'POST');
    await ui('basic-recording-retry').click();
    expect((await resent).postDataJSON()).toEqual(startRequest);
    await page.waitForURL('**/studio/runtime**');

    // 7. Verify Runtime Versions and Ledger Operations
    const versions = runtimeVersions(f);
    expect(versions).toHaveLength(2);
    const startOps = recordingStartOperations(f);
    expect(startOps).toHaveLength(1);
    expect(startOps[0].operation_id).toBe(startCommitted!.operation_id);

    // 8. Refetch Groups After Start: Applied Revisions Must Be Present
    const appliedGroups = await groups(f);
    expect(appliedGroups).toHaveLength(2);
    for (const g of appliedGroups) {
      expect(g.applied_revision).toBeTruthy();
    }

    // 9. Mandatory Outbox and Target SQL Row + Receipt Readback (Exact Strings, No Number conversion)
    for (const group of appliedGroups) {
      await expect.poll(() => deliveryCommittedOutbox(f, group.id, group.applied_revision).length, { timeout: 120_000 }).toBeGreaterThan(0);
      const outbox = deliveryCommittedOutbox(f, group.id, group.applied_revision);
      expect(outbox[0].table_name).toBe(group.destination.table_name);
      expect(outbox[0].connector_revision).toBe(group.destination.connector_revision);
      expect(outbox[0].state).toBe('sql_committed');

      const recordId = String(outbox[0].record_id);
      const effectKey = String(outbox[0].effect_key);
      const expectedDigest = String(outbox[0].payload_digest);
      expect(createHash('sha256').update(String(outbox[0].payload)).digest('hex')).toBe(expectedDigest);

      // Parse actual frozen outbox wire payload strictly matching payload.go
      const wirePayload = parseWireRowPayload(String(outbox[0].payload));
      expect(wirePayload.record_id).toBe(recordId);
      expect(wirePayload.effect_key).toBe(effectKey);

      // Query ALL frozen column names with typed scalar NVARCHAR encoding via FOR JSON PATH
      const selectSql = buildSQLServerRowSelectSql(
        f.schema,
        group.destination.table_name,
        recordId,
        wirePayload.cells,
      );
      const rowRaw = querySQLServerReadback(connConfig, selectSql);
      const sqlRow = JSON.parse(rowRaw.trim()) as Record<string, string | null>;
      expect(sqlRow.record_id).toBe(recordId);
      verifyWirePayloadAgainstRow(wirePayload, sqlRow);

      // Matching effect receipt readback via FOR JSON PATH
      const receiptTable = sqlserverIdentifier(f.schema, 'gw_effect_receipts');
      const safeEffectKey = effectKey.replace(/'/g, "''");
      const receiptSql = `SELECT TOP 1 effect_key, payload_digest FROM ${receiptTable} WHERE effect_key = '${safeEffectKey}' FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER`;
      const receiptRaw = querySQLServerReadback(connConfig, receiptSql);
      const receiptObj = JSON.parse(receiptRaw.trim()) as { effect_key?: string; payload_digest?: string };
      expect(receiptObj.effect_key).toBe(effectKey);
      expect(receiptObj.payload_digest).toBe(expectedDigest);
    }

    // 10. Final Visual Evidence Capture at Dual Viewports
    await page.goto(`${f.frontend}/studio/v2?step=4`);
    await expect(ui('circuit-profile-basic-panel')).toBeVisible();
    await tile(f.devices.a).click();

    for (const viewport of [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }]) {
      await page.setViewportSize(viewport);
      await noOverflow();
      await screenshot(`sqlserver-${storageMode}-acceptance-${viewport.width}`);
    }

    await attachJSON(`sqlserver-${storageMode}-evidence-final`, {
      storage_mode: storageMode,
      workspace_id: f.workspace_id,
      target_database: f.database,
      target_schema: f.schema,
      groups: appliedGroups.map((g) => ({
        id: g.id,
        layout: g.layout_version,
        applied_revision: g.applied_revision,
        table: g.destination.table_name,
      })),
    });
  });
}
