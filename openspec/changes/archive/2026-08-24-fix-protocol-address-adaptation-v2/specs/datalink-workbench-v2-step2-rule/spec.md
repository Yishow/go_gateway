# datalink-workbench-v2-step2-rule Specification Delta

## MODIFIED Requirements

### Requirement: Multi-rule tab management

The Step 2 source-rule workspace SHALL render a horizontally scrollable tab rail listing all rules in `state.rules`. The tab rail MUST allow the operator to add a new rule, select an existing rule for editing, rename a rule inline, toggle a rule's `enabled` flag, and remove a non-only rule. Each rule tab MUST display the rule color dot, name, `{start_address} · {data_type}` summary, enabled point count, and (when the workspace has 2+ devices) the owning device's color dot and name. When creating a new rule, the system SHALL determine the default start address based on the target device's communication protocol (`40001` for Modbus TCP/RTU/UDP, and `D0` for FATEK FBs / Mitsubishi MC 3E).

#### Scenario: Default rule on first render
- **WHEN** the operator opens Step 2 with default state
- **THEN** the tab rail shows exactly one rule with start address corresponding to the default device protocol
- **AND** for a Modbus device, start address is `40001` with data type `int16` and enabled point count `8/8`
- **AND** for an MC 3E or FATEK device, start address defaults to `D0`

#### Scenario: Add rule
- **WHEN** the operator clicks `+ 新增規則`
- **THEN** a new rule is appended to `state.rules` with name `規則 N` (N is next ordinal), default count 8, naming_prefix `BLOCK{N}_`, and device_id defaulting to the currently selected rule's device_id or `state.devices[0].id`
- **AND** `start_address` is populated using `getDefaultPlannerStartAddress(device.protocol)` (`40001` for Modbus, `D0` for FATEK/MC 3E)
- **AND** the new rule is selected

#### Scenario: Remove rule
- **WHEN** the operator hovers a rule tab and clicks the close icon while `state.rules.length >= 2`
- **THEN** the rule is removed from `state.rules`
- **AND** if the removed rule was selected, the first remaining rule becomes selected
- **AND** the close icon is hidden when `state.rules.length === 1`

#### Scenario: Multi-device device-row visibility
- **WHEN** `state.devices.length === 1`
- **THEN** the rule tab does NOT render the device color dot or device name row

#### Scenario: Multi-device device-row appears with 2+ devices
- **WHEN** `state.devices.length >= 2`
- **THEN** every rule tab renders a third row showing the owning device's color dot and name

### Requirement: Rule editor with linked reset

The Step 2 rule editor SHALL render input controls for rule properties (name, owning device, start address, point count, data type, naming prefix, scale multiplier, scale offset, data format, skipped addresses, and Modbus Share settings). All placeholder hints and range calculation summaries MUST dynamically adapt to the communication protocol of the owning device.

#### Scenario: Start address placeholder adapts to protocol
- **WHEN** editing a rule owned by an MC 3E or FATEK device
- **THEN** the start address input displays a hint indicating alphanumeric register notation (e.g. `D0`, `M0`, `W0`)
- **AND** the address range preview shows valid alphanumeric bounds (e.g. `D0 ~ D7`)

#### Scenario: Start address placeholder for Modbus
- **WHEN** editing a rule owned by a Modbus device
- **THEN** the start address input displays a hint for 5-digit Modbus registers (e.g. `40001`)
- **AND** the address range preview shows numeric bounds (e.g. `40001 ~ 40008`)

#### Scenario: Reset skipped on count change
- **GIVEN** a rule with `count=8`, `skipped_addresses=['40003','40005']`
- **WHEN** the operator changes count to 6
- **THEN** the rule's `count` becomes 6
- **AND** `skipped_addresses` becomes `[]`
- **AND** the point grid below renders 6 cells, all enabled

#### Scenario: Reset skipped on data_type change
- **GIVEN** a rule with `data_type='int16'`, `skipped_addresses=['40002']`
- **WHEN** the operator changes data_type to `int32`
- **THEN** the rule's `data_type` becomes `int32`
- **AND** `skipped_addresses` becomes `[]`
- **AND** the grid renders cells using stride 2 (int32 width)

#### Scenario: Count clamped to 1-64
- **WHEN** the operator enters count 0 or count 100
- **THEN** the reducer clamps `count` to `max(1, min(64, input))`
- **AND** the grid reflects the clamped value

#### Scenario: Device select disabled with single device
- **WHEN** `state.devices.length === 1`
- **THEN** the owning-device select is disabled
- **AND** the rule is auto-assigned to that device

#### Scenario: Invalid start address blocks continuation
- **WHEN** the operator enters an address that is invalid for the owning device protocol
- **THEN** the editor displays a localized validation error and marks the input invalid
- **AND** no derived points are produced from the invalid address
- **AND** Step 2 does not allow continuation to Step 3

#### Scenario: One invalid enabled rule blocks a multi-rule workspace
- **GIVEN** a workspace has at least one enabled rule with valid derived points
- **AND** another enabled rule has an address that is invalid for its owning device protocol
- **WHEN** the operator attempts to continue from Step 2
- **THEN** Step 2 remains active and identifies the invalid rule
- **AND** valid points from other rules do not mask the invalid rule

## ADDED Requirements

### Requirement: Protocol-aware point address derivation

The system SHALL derive point addresses by taking the rule's `start_address`, `count`, `data_type` (width stride), and owning device's `protocol`. The address calculation MUST preserve alphanumeric register area prefixes (such as `D`, `M`, `W`, `X`, `Y` for Mitsubishi/FATEK) and correctly advance address offsets without stripping characters or falling back to Modbus defaults.

#### Scenario: Derive points for MC 3E device
- **WHEN** a rule for an MC 3E device is defined with start address `D100`, count `4`, and data type `int16` (width `1`)
- **THEN** the system derives 4 points with addresses `D100`, `D101`, `D102`, and `D103`
- **AND** no Modbus `40001` fallback is applied

##### Example: Multi-protocol point derivation
| Protocol | Start Address | Count | Data Type | Expected Addresses |
| --- | --- | --- | --- | --- |
| `modbus_tcp` | `40001` | `3` | `int16` | `40001, 40002, 40003` |
| `mc_3e` | `D0` | `4` | `int16` | `D0, D1, D2, D3` |
| `mc_3e` | `D100` | `2` | `int32` | `D100, D102` |
| `mc_3e` | `X0` | `3` | `bool` | `X0, X1, X2` |
| `fatek_fbs` | `R0` | `3` | `int16` | `R0, R1, R2` |

#### Scenario: Backend address offset matches protocol radix
- **WHEN** the source-rule service offsets MC 3E address `X0` by `16`
- **THEN** it returns `X10`
- **AND** MC 3E areas `X`, `Y`, and `B` use hexadecimal suffixes while supported decimal areas retain decimal suffixes

#### Scenario: Backend rejects invalid protocol address
- **WHEN** source-rule validation or address offset receives an address invalid for the selected protocol
- **THEN** the operation returns an explicit validation error
- **AND** it does not return the original input as a successful derived address
