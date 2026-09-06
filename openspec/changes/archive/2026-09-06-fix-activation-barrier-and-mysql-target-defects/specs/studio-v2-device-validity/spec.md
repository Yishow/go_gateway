## ADDED Requirements

### Requirement: Displayed connection defaults are backed by device configuration

The device connection form SHALL render required connection fields from the device configuration only. It SHALL NOT display a fallback value that has not been written into the device configuration. Protocol default values for required connection fields SHALL be written into the device configuration during device load normalization, using the existing protocol default table as the single source of those values.

#### Scenario: A device loaded without a station number is normalized

- **GIVEN** a device using the MC protocol was created through the API without a station number
- **WHEN** the workspace loads that device
- **THEN** the device configuration carries the protocol default station number
- **AND** the connection form displays that value from the configuration

#### Scenario: Form display matches validity evaluation

- **GIVEN** a device loaded without a station number
- **WHEN** the workspace evaluates device validity after load normalization
- **THEN** the device is not marked draft-invalid on account of the station number
- **AND** the value shown in the station field is the value validity was evaluated against

#### Scenario: A missing value is shown as empty when no default applies

- **GIVEN** a device configuration field that has no protocol default
- **WHEN** the connection form renders that field
- **THEN** the field is displayed empty rather than pre-filled with an unwritten value

##### Example: MC protocol device loaded without a station number

- **GIVEN** a device with protocol mc_3e and configuration containing host and port but no station number
- **WHEN** the workspace loads and normalizes the device
- **THEN** the configuration station number is 0
- **AND** device validity for that device is valid
- **AND** the station field displays 0

#### Scenario: Normalization does not mark unchanged devices as pending save

- **GIVEN** a set of devices loaded into the workspace
- **WHEN** load normalization writes protocol default values into device configurations
- **THEN** no autosave request is issued for a device the operator has not edited
