## ADDED Requirements

### Requirement: Multi-device tab management

The Step 1 device workspace SHALL render a horizontally scrollable tab rail listing all devices in `state.devices`. The tab rail MUST allow the operator to add a new device, select an existing device for editing, rename a device inline, and remove a non-only device. The rail MUST display each device's name, protocol, host, and connection test status.

#### Scenario: Default device on first render

- **WHEN** the operator opens `/studio/v2` with default state
- **THEN** the tab rail shows exactly one device with name `PLC-生產線-01`, protocol `modbus_tcp`, host `192.168.1.100`, and test status badge `未測試`
- **AND** the device is selected and its editor is rendered below the tab rail

#### Scenario: Add device

- **WHEN** the operator clicks the `+ 新增設備` button at the end of the tab rail
- **THEN** a new device is appended to `state.devices` with default name `設備 N` where N is the next ordinal, protocol `modbus_tcp`, and host `192.168.1.10X` where X is incremented per existing device
- **AND** the new device is selected
- **AND** the editor below the tab rail shows the new device's fields

#### Scenario: Inline rename

- **WHEN** the operator types into the name input inside a device tab
- **THEN** the corresponding `state.devices[i].name` updates on every change
- **AND** the tab text reflects the new name
- **AND** other tabs are not selected as a side effect

#### Scenario: Remove device with confirmation

- **WHEN** the operator hovers a device tab with `state.devices.length >= 2`, the close icon is revealed, and the operator clicks it and confirms the dialog
- **THEN** the device is removed from `state.devices`
- **AND** the device cannot be deleted if it is the only device in the list (close icon is not shown when `state.devices.length === 1`)

#### Scenario: Cascade clean-up on device removal

- **WHEN** the operator removes a device that has associated source rules in `state.rules`
- **THEN** every rule with `device_id` matching the removed device is also removed from `state.rules`
- **AND** every point in `state.points` derived from a removed rule is removed
- **AND** every mapping in `state.mappings` keyed by a removed point is removed
- **AND** every database target in `state.db.targets` keyed by a removed point is removed

##### Example: cascade chain

- **GIVEN** state has device `dev-A` with two rules `rule-A1`, `rule-A2`, each rule has 4 derived points, all points have mappings, half the points have db targets
- **WHEN** the operator removes `dev-A`
- **THEN** `state.devices`, `state.rules`, `state.points`, `state.mappings`, `state.db.targets` no longer contain any entry referencing `dev-A`, `rule-A1`, `rule-A2`, or any of their derived point IDs

---

### Requirement: Protocol switching with reset

The device editor SHALL provide a protocol selector with six options: `modbus_tcp`, `modbus_rtu`, `modbus_udp`, `fatek_fbs`, `mc_3e`, `mqtt`. Switching the protocol MUST replace the device's `config` with the default configuration for the new protocol, clear `device.test` to `null`, and reset `device.status` to `'draft'`.

#### Scenario: Switch from modbus_tcp to mqtt

- **WHEN** the operator selects `mqtt` from the protocol selector while the device has protocol `modbus_tcp` and a successful `test`
- **THEN** `device.protocol` becomes `mqtt`
- **AND** `device.config` is replaced with mqtt defaults (broker `mqtts://broker.local:8883`, empty username, client_id `gw-01`)
- **AND** `device.test` becomes `null`
- **AND** `device.status` becomes `'draft'`
- **AND** the connection config form swaps to the mqtt layout (Broker / Username / Client ID inputs)

#### Scenario: Per-protocol default configuration

- **WHEN** the operator switches the protocol to any of the six options
- **THEN** the device's config matches the default configuration table below

##### Example: defaults table

| Protocol     | Default config (JSON)                                                                  |
| ------------ | -------------------------------------------------------------------------------------- |
| modbus_tcp   | `{ host: "192.168.1.100", port: 502, slave_id: 1, timeout: 5 }`                       |
| modbus_udp   | `{ host: "192.168.1.100", port: 502, slave_id: 1, timeout: 5 }`                       |
| fatek_fbs    | `{ host: "192.168.1.100", port: 500, slave_id: 1, timeout: 5 }`                       |
| mc_3e        | `{ host: "192.168.1.100", port: 5007, slave_id: 1, timeout: 5 }`                      |
| modbus_rtu   | `{ port: "/dev/ttyUSB0", baud: 9600, parity: "N", slave_id: 1 }`                      |
| mqtt         | `{ broker: "mqtts://broker.local:8883", username: "", client_id: "gw-01" }`           |

---

### Requirement: Readiness check with explicit connect / probe separation

The connection test panel SHALL display a readiness check that groups stages into two semantic groups: `connect` (network or channel establishment) and `probe` (protocol-level handshake or subscription). Each stage MUST be tagged with its group and rendered under a labeled section separator. The operator MUST be able to distinguish whether a failure occurred at the channel layer or the protocol layer by visual grouping.

#### Scenario: TCP-family stages

- **WHEN** the device protocol is `modbus_tcp`, `modbus_udp`, `fatek_fbs`, or `mc_3e`
- **THEN** the readiness check shows three stages in order: `resolve` (group connect), `connect` (group connect), `probe` (group probe)
- **AND** the UI renders a section separator labeled "連線通道" (or i18n equivalent) above the connect group and "協議握手" above the probe group

#### Scenario: RTU stages

- **WHEN** the device protocol is `modbus_rtu`
- **THEN** the readiness check shows three stages in order: `open_port` (group connect), `handshake` (group connect), `probe` (group probe)

#### Scenario: MQTT stages

- **WHEN** the device protocol is `mqtt`
- **THEN** the readiness check shows three stages in order: `resolve` (group connect), `connect` (group connect), `subscribe` (group probe)

#### Scenario: Stage progression during test

- **WHEN** the operator clicks `執行測試`
- **THEN** the `執行測試` button becomes disabled
- **AND** every 380ms the next stage transitions from `pending` to `success` with a per-stage latency between 10ms and 28ms
- **AND** once all stages complete, `device.test.status` becomes `success`, `device.test.latency_ms` is the sum of stage latencies plus jitter (42-60ms range), `device.status` becomes `tested`
- **AND** the success readiness panel renders below the stage list listing `Connection Configuration ✓`, `Protocol Probe ✓ ({latency_ms}ms)`, `可進行下一步：設定接入規則`

##### Example: mqtt success sequence

- **GIVEN** device with protocol `mqtt` and no prior test
- **WHEN** the operator clicks `執行測試` and 220ms + 3 × 380ms elapse
- **THEN** stages progress: resolve → success → connect → success → subscribe → success
- **AND** `device.test.status === 'success'` with `device.status === 'tested'`

---

### Requirement: Connection config form per protocol

The device editor SHALL render protocol-specific connection configuration fields. The form layout MUST adapt to the active protocol: TCP-family protocols (`modbus_tcp`, `modbus_udp`, `fatek_fbs`, `mc_3e`) render Host / Port / Slave ID / Timeout(s) in a 4-column grid; `modbus_rtu` renders Port / Baud Rate / Parity / Slave ID; `mqtt` renders Broker / Username / Client ID.

#### Scenario: TCP layout

- **WHEN** the device protocol is `modbus_tcp`
- **THEN** the connection config form renders inputs labeled `Host`, `Port`, `Slave ID`, `Timeout (s)`
- **AND** all numeric inputs use mono font and accept positive numbers

#### Scenario: RTU layout

- **WHEN** the device protocol is `modbus_rtu`
- **THEN** the connection config form renders inputs labeled `Port`, `Baud Rate`, `Parity`, `Slave ID`
- **AND** `Parity` is a select with options `N`, `E`, `O`

#### Scenario: MQTT layout

- **WHEN** the device protocol is `mqtt`
- **THEN** the connection config form renders inputs labeled `Broker`, `Username`, `Client ID`

---

### Requirement: Device color cycle

The system SHALL expose a deterministic 6-color cycle (`blue`, `emerald`, `amber`, `fuchsia`, `cyan`, `rose`) keyed by the device's position in `state.devices`. A helper hook `useDeviceColor(deviceId)` MUST return the color name plus Tailwind class fragments (bg, border, text, solid). This color cycle MUST be the single source of truth used by Step 2 (rule tabs), the merged point table, and the right SummaryRail when those step changes land.

#### Scenario: First six devices map to fixed colors

- **WHEN** the operator creates devices at positions 0 through 5 (zero-indexed)
- **THEN** `useDeviceColor(devices[0].id).name === 'blue'`
- **AND** `useDeviceColor(devices[1].id).name === 'emerald'`
- **AND** `useDeviceColor(devices[2].id).name === 'amber'`
- **AND** `useDeviceColor(devices[3].id).name === 'fuchsia'`
- **AND** `useDeviceColor(devices[4].id).name === 'cyan'`
- **AND** `useDeviceColor(devices[5].id).name === 'rose'`

#### Scenario: Wrap-around on 7th device

- **WHEN** the operator creates a 7th device (position 6)
- **THEN** `useDeviceColor(devices[6].id).name === 'blue'`

---

### Requirement: Continue gate

The bottom action bar SHALL render `{N} 個設備，{M} 已通過測試` plus a `全部建立並繼續` button. The button MUST be disabled when any device has not yet completed a successful test. When enabled and clicked, the button MUST invoke the shell `onContinue` callback for Step 1.

#### Scenario: Disabled when not all tested

- **WHEN** there are two devices and only the first has `test.status === 'success'`
- **THEN** the continue button is disabled
- **AND** a warning chip `尚有 1 個設備未通過測試` is rendered next to the count text

#### Scenario: Enabled when all tested

- **WHEN** every device in `state.devices` has `test.status === 'success'`
- **THEN** the continue button is enabled
- **AND** clicking it invokes the shell `onContinue` callback, which causes the shell to add step 1 to `completed` and advance `current` to 2

---

### Requirement: Send payload preview block

The connection test panel SHALL render a payload preview block showing the JSON body that would be sent to `POST /api/v1/datalink/devices/test`. The preview MUST update reactively as the operator edits the protocol or connection config fields. The preview MUST NOT cause any actual network request to be issued.

#### Scenario: Payload reflects current form

- **WHEN** the operator changes the host to `10.0.0.5`
- **THEN** the payload preview's `connection_config.host` field shows `"10.0.0.5"`
- **AND** no fetch or XHR request to `/api/v1/datalink/devices/test` is initiated

#### Scenario: Payload mock notice

- **WHEN** the payload preview is rendered
- **THEN** the preview block clearly indicates the payload is a mock preview (for example via section label `送出 payload`)
- **AND** running the `執行測試` button triggers the mock animation instead of a real request
