## ADDED Requirements

### Requirement: Settings page layout

The Settings page SHALL render a header banner card, five SectionCards in vertical order (`Connector Pool`, `Timeseries Strategy` + `Scheduler Defaults` side by side, `Local Modbus Share`, `UI` + `API` + `Diagnostics` side by side), and a sticky bottom save bar. The right summary rail of the shell MUST be hidden while the view is `settings`. The save bar MUST remain visible during page scroll.

#### Scenario: Default render

- **WHEN** the operator switches the shell view to `settings`
- **THEN** the header banner card is rendered first with the slider icon and `系統設定` heading
- **AND** the five SectionCards appear in the specified order
- **AND** the right summary rail is not present in the DOM
- **AND** the save bar at the bottom uses `sticky bottom-4 z-10` positioning

#### Scenario: Save bar visible during scroll

- **WHEN** the operator scrolls the settings page downward
- **THEN** the save bar remains visible at the bottom of the viewport (sticky positioning intact)
- **AND** the save bar contains the `重設為預設` ghost button and `儲存所有設定` success button

---

### Requirement: Connector pool CRUD

The Connector Pool section SHALL list every entry in `state.settings.connectors`, allow the operator to add a new connector via the `+ 新增連接器` aside button, edit per-row fields (name, kind, host, port, database, username, default_write_interval_seconds, enabled), test the connection, and remove a connector. Adding a connector MUST append a new entry with default name `新連線 N`, kind `postgres`, status `unknown`. Editing `kind`, `host`, or `port` MUST set the connector's status back to `unknown`. Removing a connector MUST be immediate (no confirmation in this change).

#### Scenario: Default connector

- **WHEN** the settings page mounts with the application's default settings
- **THEN** `state.settings.connectors` contains exactly one entry with id `conn-prod`, name `TimeSeries Prod`, kind `postgres`, host `tsdb.internal`, status `ready`

#### Scenario: Add connector

- **GIVEN** the pool has 1 connector
- **WHEN** the operator clicks `+ 新增連接器`
- **THEN** `state.settings.connectors.length === 2`
- **AND** the new entry has name `新連線 2`, kind `postgres`, status `unknown`, port `5432`, database `metrics`, username `gw_writer`, default_write_interval_seconds `5`

#### Scenario: Edit kind resets status

- **GIVEN** a connector has status `ready` and last_check_at set
- **WHEN** the operator changes the connector's `kind` from `postgres` to `mysql`
- **THEN** the connector's status becomes `unknown`
- **AND** the StatusChip displays `未測試` (or i18n equivalent)

#### Scenario: Remove connector

- **GIVEN** the pool has 2 connectors
- **WHEN** the operator clicks the remove button on connector 2
- **THEN** `state.settings.connectors.length === 1`
- **AND** no confirmation dialog is shown

---

### Requirement: Connector mock test

The system SHALL provide a mock connector test triggered by the test button on each connector row. When triggered, the connector's status MUST transition to `testing` and the test button MUST become disabled with a spinning refresh icon. After exactly 900ms, the connector's status MUST become either `ready` (85% probability) with `last_check_at` set to the current ISO timestamp, or `unreachable` (15% probability) with `last_check_at` set and `last_check_error` set to `connection refused`. The test MUST NOT issue any real network request.

#### Scenario: Success path

- **GIVEN** `Math.random` is mocked to return `0.5` (i.e. > 0.15)
- **WHEN** the operator clicks the test button on a connector and `vi.useFakeTimers()` advances 900ms
- **THEN** the connector's status becomes `ready`
- **AND** `last_check_at` is set to a valid ISO timestamp
- **AND** the StatusChip displays `已就緒` (or i18n equivalent) in success tone

#### Scenario: Failure path

- **GIVEN** `Math.random` is mocked to return `0.05` (i.e. ≤ 0.15)
- **WHEN** the operator clicks the test button and 900ms elapses
- **THEN** the connector's status becomes `unreachable`
- **AND** `last_check_error` becomes `connection refused`
- **AND** the StatusChip displays `無法連線` in error tone
- **AND** the error text `connection refused` is rendered next to the chip

#### Scenario: No network call

- **WHEN** the operator clicks the test button
- **THEN** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the settings page

---

### Requirement: Timeseries strategy fields

The Timeseries section SHALL provide four fields with the following defaults and constraints:

- `write_precision`: select with options `second` and `millisecond`, default `millisecond`.
- `partition_interval`: select with options `daily`, `weekly`, `monthly`, default `daily`.
- `batch_size`: number, default `500`, hint `一次最多寫入的筆數`.
- `retention_days`: number, default `90`, hint `超過天數的資料將被歸檔`.

Changes MUST write through to `state.settings.timeseries` via `updateSettingsSection`.

#### Scenario: Default values

- **WHEN** the settings page mounts with default state
- **THEN** the timeseries section selects show `millisecond` and `daily`
- **AND** the number inputs show `500` and `90`

#### Scenario: Update write_precision

- **WHEN** the operator selects `second` from the write_precision select
- **THEN** `state.settings.timeseries.write_precision` becomes `second`

---

### Requirement: Scheduler defaults fields

The Scheduler section SHALL provide five fields:

- `default_interval_ms`: number, default `1000`, hint `毫秒 (ms)`.
- `default_retry_count`: number, default `3`.
- `default_retry_delay_ms`: number, default `500`, hint `毫秒 (ms)`.
- `breaker_threshold`: number, default `10`, hint `連續錯誤次數`.
- `auto_start`: toggle, default `true`, label `開機時自動啟動 collector`.

Changes MUST write through to `state.settings.scheduler`.

#### Scenario: Toggle auto_start

- **WHEN** the operator toggles the `auto_start` switch from on to off
- **THEN** `state.settings.scheduler.auto_start` becomes `false`

---

### Requirement: Local Modbus Share configuration

The Local Modbus Share section SHALL provide a master toggle in the aside slot and a 4-column form below: `bind_address` (text mono, default `0.0.0.0`, hint `0.0.0.0 = 全部介面`), `port` (number, default `5020`), `slave_id` (number, default `1`), `base_register` (number, default `40001`). The form MUST be hidden when the master toggle is off and a help text MUST be shown in its place. Toggling the master toggle MUST update `state.settings.modbus_share.enabled` and MUST propagate to Step 2's share layout computation via the existing `useShareLayout` selector.

#### Scenario: Toggle hides the form

- **GIVEN** the master toggle is on and the form is visible
- **WHEN** the operator toggles it off
- **THEN** the 4-column form is removed from the DOM
- **AND** the help text `啟用後可將所有 Tag 再次經由 Modbus TCP 提供給其他系統。` (or i18n equivalent) is shown

#### Scenario: base_register change affects Step 2

- **GIVEN** `state.settings.modbus_share.enabled === true` and `base_register === 40001`
- **WHEN** the operator changes `base_register` to `50001`
- **AND** the operator navigates to Step 2 with rules that have `share_start_register === null`
- **THEN** the Step 2 share layout for those rules starts at `50001` (the new base)

---

### Requirement: UI / API / Diagnostics options

The General settings SHALL be split across three side-by-side SectionCards:

- **UI**: `theme` select (`dark` / `light` / `auto`, default `dark`), `locale` select (`zh-TW` / `en`, default `zh-TW`), `addr_format` select (`modbus` / `hex` / `raw`, default `modbus`).
- **API**: `api_base` text (mono, default `http://localhost:8080`, hint `不含尾斜線`), `api_version` select (`v1` / `v2`, default `v1`), `timeout_seconds` number (default `30`, hint `秒`).
- **Diagnostics**: `log_level` select (`trace`/`debug`/`info`/`warn`/`error`, default `info`), `sse_heartbeat_seconds` number (default `15`, hint `秒`), `enable_debug_panel` toggle (default `false`), `enable_audit_log` toggle (default `true`).

Changes MUST write through to `state.settings.general`. The theme select MUST NOT actually toggle the UI theme during this change (dark-only); writing the value is sufficient.

#### Scenario: Change theme stored but UI stays dark

- **WHEN** the operator selects `light` from the theme select
- **THEN** `state.settings.general.theme` becomes `light`
- **AND** the v2 page subtree continues to use dark tokens (the visual theme does not change)

#### Scenario: Toggle audit log

- **WHEN** the operator toggles `enable_audit_log` from on to off
- **THEN** `state.settings.general.enable_audit_log` becomes `false`

---

### Requirement: Save bar actions

The sticky save bar SHALL contain a left-side info text equivalent to `設定會立即套用，並於下次重啟後生效。`, a `重設為預設` ghost button, and a `儲存所有設定` success button. The `重設為預設` button MUST trigger a native `window.confirm`; on confirmation, `state.settings` MUST be replaced with `DEFAULT_SETTINGS`. The `儲存所有設定` button MUST be a no-op in this change and MUST log a warning (`console.warn('pending backend-wiring')`) when clicked; it MUST NOT issue any network request.

#### Scenario: Reset with confirm

- **GIVEN** the operator has modified several settings fields
- **WHEN** the operator clicks `重設為預設` and confirms the native dialog
- **THEN** `state.settings` equals `DEFAULT_SETTINGS`

#### Scenario: Reset cancelled

- **GIVEN** the operator has modified settings fields
- **WHEN** the operator clicks `重設為預設` and cancels the dialog
- **THEN** `state.settings` remains unchanged

#### Scenario: Save is no-op

- **WHEN** the operator clicks `儲存所有設定`
- **THEN** `console.warn` is called once with a message indicating backend-wiring is pending
- **AND** no `fetch` or other network call is initiated
- **AND** `state.settings` remains unchanged
