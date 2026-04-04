# datalink-workbench-desktop Specification Delta

## ADDED Requirements

### Requirement: Source rule builder exposes target data type selection

The system SHALL provide, within Step 2 source planning, a control to set the **target / intended tag data type** alongside the protocol read data type when creating or editing a source rule. The control SHALL default to matching the protocol read data type and SHALL clearly distinguish protocol read semantics from target output semantics.

#### Scenario: Operator sets distinct target type before applying rule

- **WHEN** an operator configures a new rule and selects a target data type different from the protocol read type
- **THEN** the workbench SHALL persist that intent with the rule according to backend contract
- **AND** SHALL show validation feedback when the pair is unsupported

#### Scenario: Defaults avoid extra friction

- **WHEN** an operator does not change the default target data type
- **THEN** the workbench SHALL behave as today with a single data type selection driving protocol read planning

### Requirement: Device connection form exposes data format when the protocol uses multi-register decoding

The workbench SHALL render **`data_format`** (or the documented equivalent) in the device connection editor for each protocol whose connector uses multi-register / floating-point byte ordering (for example **Modbus TCP, Modbus RTU, Modbus UDP**, and **Mitsubishi MC 3E**). The control SHALL list the supported orderings (for example **ABCD, BADC, CDAB, DCBA**) and SHALL persist the value in `connection_config` for the backend to consume during reads.

#### Scenario: Modbus connection shows byte order

- **WHEN** an operator edits a Modbus-class device connection
- **THEN** the UI SHALL expose `data_format` selection
- **AND** saving the device SHALL persist the chosen value

### Requirement: Source rule builder MAY expose scale or offset for unit conversion

The system SHALL provide optional inputs for **linear scaling** (multiplier and/or offset) on the source rule form when the product enables engineering-unit conversion from rule planning, with clear labels that distinguish scaling from protocol read type and from target tag type.

#### Scenario: Operator sets scale without changing target type

- **WHEN** an operator enters scale parameters only
- **THEN** the workbench SHALL persist scale intent with the rule
- **AND** SHALL not require a target type different from the read type
