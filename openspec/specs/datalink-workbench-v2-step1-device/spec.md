# datalink-workbench-v2-step1-device Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-step1-device'. Update Purpose after archive.

## Requirements

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


<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

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


<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

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


<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

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


<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

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


<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

---
### Requirement: Continue gate
The bottom action bar SHALL show total devices, the explicitly selected setup scope and successful current tests. At least one selected device and successful tests for every selected device MUST be required to continue. Unselected devices SHALL remain visible drafts, not be silently enabled or deleted. A missing selection SHALL preserve the legacy all-device scope. Actual activation SHALL still revalidate current device readiness and the existing workspace barrier.

#### Scenario: Disabled when not all tested
- **WHEN** two devices are selected and only the first has a successful current test
- **THEN** the continue button is disabled and the untested selected device and repair action are shown.

#### Scenario: Enabled when all tested
- **WHEN** every device in the selected scope has a successful current test
- **THEN** clicking continue invokes onContinue, completes Step 1 and advances to Step 2 for that scope.

#### Scenario: Offline device explicitly excluded
- **WHEN** a tested device remains selected and the user explicitly excludes an offline device
- **THEN** setup can continue for the tested device and the excluded device remains a clearly labelled draft.

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

<!-- @trace
source: datalink-workbench-v2-step1-device
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/i18n/config.ts
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
-->

---
### Requirement: Valid-only device autosave

The Step 1 device workspace SHALL autosave only valid device content to the backend workspace.

#### Scenario: Valid device change saves immediately

- **GIVEN** a device row already belongs to the v2 workspace
- **WHEN** the operator edits the device into a valid form state
- **THEN** the system sends an immediate save request for that device
- **AND** the saved backend version reflects the latest valid values

#### Scenario: Invalid device change stays local

- **GIVEN** a device row already has a last successful backend version
- **WHEN** the operator edits the device into an invalid form state
- **THEN** the system does not overwrite the backend version
- **AND** the UI keeps the invalid local values visible with an unsaved marker


<!-- @trace
source: wire-studio-v2-step1-device-autosave
updated: 2026-05-31
code:
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_status.go
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/sourcerule/service.go
  - task_plan.md
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/migrator.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/router.go
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/device/service_point_read.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/services/studioV2Mappings.ts
  - .air.toml
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - go.sum
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - scripts/start-process-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/errors.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/index.ts
  - internal/datalink/workspace/service_activation.go
  - progress.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service_devices.go
  - scripts/start-log-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - internal/datalink/workspace/memory_repository.go
tests:
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_device_order_test.go
-->

---
### Requirement: Per-device save isolation

The Step 1 device workspace SHALL isolate save success and failure per device tab.

#### Scenario: One invalid device does not block another valid device

- **GIVEN** the workspace shows multiple devices
- **WHEN** one device remains invalid while another device is edited into a valid state
- **THEN** the valid device still saves successfully
- **AND** the invalid device remains local with its own error state


<!-- @trace
source: wire-studio-v2-step1-device-autosave
updated: 2026-05-31
code:
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_status.go
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/sourcerule/service.go
  - task_plan.md
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/migrator.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/router.go
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/device/service_point_read.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/services/studioV2Mappings.ts
  - .air.toml
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - go.sum
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - scripts/start-process-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/errors.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/index.ts
  - internal/datalink/workspace/service_activation.go
  - progress.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service_devices.go
  - scripts/start-log-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - internal/datalink/workspace/memory_repository.go
tests:
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_device_order_test.go
-->

---
### Requirement: Persisted device deletion cascade

Deleting a persisted Step 1 device from the v2 workspace SHALL remove that device and its v2-owned related data from the backend.

#### Scenario: Delete persisted device

- **GIVEN** a persisted device belongs to the singleton v2 workspace
- **WHEN** the operator deletes that device from Step 1
- **THEN** the backend removes the device from the workspace
- **AND** the device no longer appears after a page reload

<!-- @trace
source: wire-studio-v2-step1-device-autosave
updated: 2026-05-31
code:
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_status.go
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/sourcerule/service.go
  - task_plan.md
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/migrator.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/router.go
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/device/service_point_read.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/services/studioV2Mappings.ts
  - .air.toml
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - go.sum
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - scripts/start-process-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/errors.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/index.ts
  - internal/datalink/workspace/service_activation.go
  - progress.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service_devices.go
  - scripts/start-log-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - internal/datalink/workspace/memory_repository.go
tests:
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_device_order_test.go
-->

---
### Requirement: Live backend diagnostics for Step 1

The Step 1 device test panel SHALL use real backend diagnostics instead of frontend mock animation.

#### Scenario: Real diagnostics request on execute test

- **GIVEN** the operator is editing a device in Step 1
- **WHEN** the operator clicks `執行測試`
- **THEN** the system sends a real backend diagnostics request for the current draft or persisted device state
- **AND** the UI does not synthesize success from local timers or random values

#### Scenario: Diagnostics success marks the device as tested

- **GIVEN** backend diagnostics returns successful connect and probe outcomes
- **WHEN** the response is received
- **THEN** the device enters `tested` state
- **AND** the continue gate may treat that device as passed

#### Scenario: Diagnostics failure remains actionable

- **GIVEN** backend diagnostics reports connect failure or probe failure
- **WHEN** the response is received
- **THEN** the UI shows the failing stage and message
- **AND** the device does not enter `tested` state
- **AND** the continue gate remains blocked for that device

<!-- @trace
source: replace-step1-mock-test-with-live-diagnostics
updated: 2026-05-31
code:
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/migrator.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/workspace/service_devices.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - findings.md
  - internal/api/handlers/point_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/api/handlers/point_direct_reader.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - docs/goal.md
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/datalink/mapping/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/service.go
  - internal/api/router.go
  - scripts/start-process-utils.sh
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/services/studioV2Workspace.ts
  - start.ps1
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/runtime_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2Mappings.ts
  - go.mod
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/datalink/dbtarget/types.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/sourcerule/service_links.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/db.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - scripts/start-log-utils.sh
  - go.sum
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - start.sh
  - .air.toml
  - gateway.db
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/workspace/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
tests:
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/availability_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
-->

---
### Requirement: Step 1 reload restores persisted device definitions and save truth

The Step 1 device workspace SHALL restore persisted device definitions and their persisted save truth after reload.

#### Scenario: Reload shows persisted device state instead of defaults

- **WHEN** the operator reloads Studio V2 after one or more devices have already been persisted
- **THEN** Step 1 renders those persisted device definitions, protocol settings, and persisted save state indicators
- **AND** Step 1 SHALL NOT show the default single-device form for devices that already exist in persisted storage

<!-- @trace
source: preserve-studio-v2-setup-state
updated: 2026-06-09
code:
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/dbtarget/live_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/dbtarget/writer.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/studio-v2-real-check.cjs
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - Makefile
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/api/router.go
  - frontend/src/types/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/runtime/service_device_sync.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/audit/types.go
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/workspace/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/dbtarget/tooling_service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - tests/shell/start-frontend-install-failure.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/api/handlers/studio_v2_workspace_audit.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/audit/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/types/databaseDelivery.ts
  - internal/datalink/runtime/service_workspace_projection.go
  - internal/datalink/runtime/status.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/schema/schema_dbtarget_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - tests/shell/start-frontend-readiness.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/audit/sql_repository.go
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/App.tsx
  - internal/datalink/dbtarget/delivery_outcome.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/datalink/dbtarget/sql_repository.go
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - tests/shell/start-frontend-deps.sh
tests:
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - internal/datalink/runtime/service_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/datalink/audit/service_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/api/router_runtime_database_delivery_test.go
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/runtime/ingestor_test.go
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
-->

---
### Requirement: Protocol switch auto-populates dedicated communication presets

When adding or editing a device in Step 1, selecting any supported communication protocol SHALL auto-populate dedicated default communication parameters for that protocol:
- `mc_3e` (MC Protocol): `port` SHALL default to `6000`, `station` SHALL default to `0`, `network_no` SHALL default to `0`, `pc_no` SHALL default to `255`, `module_io_no` SHALL default to `1023`, `module_station_no` SHALL default to `0`.
- `fatek_fbs`: `port` SHALL default to `500`, `station` SHALL default to `1`.
- `modbus_tcp`: `port` SHALL default to `502`, `slave_id` SHALL default to `1`, `timeout` SHALL default to `5`.
- `modbus_udp`: `port` SHALL default to `502`, `slave_id` SHALL default to `1`, `timeout` SHALL default to `5`.
- `modbus_rtu`: `port` SHALL default to `/dev/ttyUSB0`, `baud` SHALL default to `9600`, `parity` SHALL default to `N`, `slave_id` SHALL default to `1`.
- `mqtt`: `broker` SHALL default to `mqtt://127.0.0.1:1883`, `client_id` SHALL default to `gw-01`.

#### Scenario: Selecting MC Protocol in Step 1
- **WHEN** the operator selects protocol `mc_3e` in the Device creation modal or editor
- **THEN** the port field SHALL default to `6000`
- **AND** the station field SHALL default to `0`

#### Scenario: Selecting Fatek Protocol in Step 1
- **WHEN** the operator selects protocol `fatek_fbs`
- **THEN** the port field SHALL default to `500`
- **AND** the station field SHALL default to `1`

#### Scenario: Selecting Modbus RTU Protocol in Step 1
- **WHEN** the operator selects protocol `modbus_rtu`
- **THEN** the serial port field SHALL default to `/dev/ttyUSB0`
- **AND** the baud rate SHALL default to `9600`
- **AND** parity SHALL default to `N`


<!-- @trace
source: fix-studio-v2-setup-and-db-flow
updated: 2026-09-06
code:
  - scripts/lib/b10_focused_suite.py
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/migrator_database_target_mysql_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/dbtarget/tooling_service.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/modbusshare/settings_lifecycle.go
  - internal/datalink/runtime/status.go
  - docs/swagger/swagger.yaml
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - internal/datalink/sourcerule/candidate_api.go
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/api/router.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/settings/errors.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - internal/datalink/modbusshare/service_projection.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - docs/swagger/docs.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/api/handlers/device_health_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/sourcerule/service_state.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/sourcerule/share_desired_mappings.go
  - internal/virtual/server/modbus/server.go
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/geometry.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/api/handlers/mapping_handler.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/datalink/modbusshare/types.go
  - internal/datalink/workspace/service_readiness_device.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/datalink/settings/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/device/repository_memory.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/typed_errors.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/collector/scheduler_dispatch.go
  - docs/swagger/swagger.json
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/hooks/previewStreamEvents.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - frontend/src/services/datalink.ts
  - internal/api/handlers/response_keys.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/datalink/modbusshare/errors.go
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/modbusshare/canonical_plan.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/template.go
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/connection.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/types/modbusShare.test-d.ts
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/types/sourceRuleCandidates.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/utils/typedErrors.ts
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/sourcerule/interfaces.go
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/modbusshare/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/vite.config.ts
  - scripts/lib/b10_acceptance_helpers.py
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/datalink/modbusshare/service_persistence.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/types/datalink.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/services/sourceRuleCandidates.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/dashboard_handler.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/tag/service_crud.go
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - internal/api/handlers/test.go
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - internal/datalink/modbusshare/sql_revision_store.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/api/handlers/settings_handler.go
  - internal/api/handlers/source_rule_handler_output_apply.go
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/datalink/device/service_readiness_safe.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/virtual/server/modbus/server_connection.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - internal/datalink/sourcerule/validation.go
  - internal/api/handlers/transport_wrapper.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - internal/api/router_modbus_share.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/sourcerule/tag_apply_service.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - internal/api/handlers/modbus_share_handler_status.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/services/datalinkClient.ts
  - internal/datalink/collector/scheduler.go
  - scripts/lib/b10_projection_assertions.py
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/virtual/server/modbus/server_register_handlers.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - scripts/check_file_lines.sh
  - internal/api/handlers/datalink_sse_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/types/modbusShare.ts
  - internal/datalink/dbtarget/service_validate.go
  - internal/api/handlers/modbus_share_handler.go
  - scripts/lib/b10_negative_steps.py
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/service.go
  - frontend/src/services/modbusShare.ts
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - go.mod
tests:
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/modbusshare/service_test.go
  - internal/datalink/modbusshare/status_contract_test.go
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - internal/datalink/device/service_readiness_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - internal/api/handlers/polling_group_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - scripts/lib/test_check_file_lines.py
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
  - internal/datalink/collector/scheduler_refactor_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/api/handlers/test_client_factory.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - internal/datalink/sourcerule/candidate_scope_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - cmd/test_ui/harness_config_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/datalink_sse_handler_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - cmd/test_ui/harness_config.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/modbusshare/reconciler_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - cmd/test_ui/main.go
  - internal/virtual/server/modbus/server_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/utils/safeJson.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/device/service_refactor_test.go
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/api/handlers/test_connection_handler.go
  - cmd/test_ui/share_startup.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - internal/api/handlers/test_monitor_handler.go
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - cmd/test_ui/server_runtime.go
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/api/device_handler_test.go
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - internal/api/modbus_share_swagger_contract_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - internal/api/handlers/test_request_helpers_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/migrator_test.go
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/sourcerule/share_gates_test.go
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
-->

---
### Requirement: Protocol form frontend-backend field alignment

The Device Connection form in Step 1 SHALL accurately render protocol-specific fields and labels, and SHALL serialize parameters according to each backend protocol driver specification without missing or mismatched field names.

#### Scenario: Submitting MC Protocol device
- **GIVEN** an MC Protocol device with port `6000` and station `0`
- **WHEN** the operator saves the device
- **THEN** the payload sent to backend SHALL contain `port: 6000` and station configuration matching the MC Protocol backend driver specification

<!-- @trace
source: fix-studio-v2-setup-and-db-flow
updated: 2026-09-06
code:
  - scripts/lib/b10_focused_suite.py
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/migrator_database_target_mysql_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/dbtarget/tooling_service.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/modbusshare/settings_lifecycle.go
  - internal/datalink/runtime/status.go
  - docs/swagger/swagger.yaml
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - internal/datalink/sourcerule/candidate_api.go
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/api/router.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/settings/errors.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - internal/datalink/modbusshare/service_projection.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - docs/swagger/docs.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/api/handlers/device_health_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/sourcerule/service_state.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/sourcerule/share_desired_mappings.go
  - internal/virtual/server/modbus/server.go
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/geometry.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/api/handlers/mapping_handler.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/datalink/modbusshare/types.go
  - internal/datalink/workspace/service_readiness_device.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/datalink/settings/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/device/repository_memory.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/typed_errors.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/collector/scheduler_dispatch.go
  - docs/swagger/swagger.json
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/hooks/previewStreamEvents.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - frontend/src/services/datalink.ts
  - internal/api/handlers/response_keys.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/datalink/modbusshare/errors.go
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/modbusshare/canonical_plan.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/template.go
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/connection.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/types/modbusShare.test-d.ts
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/types/sourceRuleCandidates.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/utils/typedErrors.ts
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/sourcerule/interfaces.go
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/modbusshare/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/vite.config.ts
  - scripts/lib/b10_acceptance_helpers.py
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/datalink/modbusshare/service_persistence.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/types/datalink.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/services/sourceRuleCandidates.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/dashboard_handler.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/tag/service_crud.go
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - internal/api/handlers/test.go
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - internal/datalink/modbusshare/sql_revision_store.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/api/handlers/settings_handler.go
  - internal/api/handlers/source_rule_handler_output_apply.go
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/datalink/device/service_readiness_safe.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/virtual/server/modbus/server_connection.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - internal/datalink/sourcerule/validation.go
  - internal/api/handlers/transport_wrapper.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - internal/api/router_modbus_share.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/sourcerule/tag_apply_service.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - internal/api/handlers/modbus_share_handler_status.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/services/datalinkClient.ts
  - internal/datalink/collector/scheduler.go
  - scripts/lib/b10_projection_assertions.py
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/virtual/server/modbus/server_register_handlers.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - scripts/check_file_lines.sh
  - internal/api/handlers/datalink_sse_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/types/modbusShare.ts
  - internal/datalink/dbtarget/service_validate.go
  - internal/api/handlers/modbus_share_handler.go
  - scripts/lib/b10_negative_steps.py
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/service.go
  - frontend/src/services/modbusShare.ts
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - go.mod
tests:
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/modbusshare/service_test.go
  - internal/datalink/modbusshare/status_contract_test.go
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - internal/datalink/device/service_readiness_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - internal/api/handlers/polling_group_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - scripts/lib/test_check_file_lines.py
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
  - internal/datalink/collector/scheduler_refactor_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/api/handlers/test_client_factory.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - internal/datalink/sourcerule/candidate_scope_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - cmd/test_ui/harness_config_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/datalink_sse_handler_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - cmd/test_ui/harness_config.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/modbusshare/reconciler_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - cmd/test_ui/main.go
  - internal/virtual/server/modbus/server_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/utils/safeJson.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/device/service_refactor_test.go
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/api/handlers/test_connection_handler.go
  - cmd/test_ui/share_startup.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - internal/api/handlers/test_monitor_handler.go
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - cmd/test_ui/server_runtime.go
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/api/device_handler_test.go
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - internal/api/modbus_share_swagger_contract_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - internal/api/handlers/test_request_helpers_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/migrator_test.go
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/sourcerule/share_gates_test.go
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
-->
