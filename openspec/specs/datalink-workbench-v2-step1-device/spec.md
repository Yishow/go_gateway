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

The bottom action bar SHALL render `{N} 個設備，{M} 已通過測試` plus a `全部建立並繼續` button. The button MUST be disabled when any device has not yet completed a successful test. When enabled and clicked, the button MUST invoke the shell `onContinue` callback for Step 1.

#### Scenario: Disabled when not all tested

- **WHEN** there are two devices and only the first has `test.status === 'success'`
- **THEN** the continue button is disabled
- **AND** a warning chip `尚有 1 個設備未通過測試` is rendered next to the count text

#### Scenario: Enabled when all tested

- **WHEN** every device in `state.devices` has `test.status === 'success'`
- **THEN** the continue button is enabled
- **AND** clicking it invokes the shell `onContinue` callback, which causes the shell to add step 1 to `completed` and advance `current` to 2


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