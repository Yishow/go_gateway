## ADDED Requirements

### Requirement: Aggregated device connection validation

Device connection validation SHALL report every missing required field and every null or empty-string required field in one error, preserving schema order within each category. When only one category exists, the existing single-category message format SHALL remain intact. Missing fields SHALL precede empty fields in a mixed message. Existing FATEK serial/network conditional requirements and numeric range/enum checks SHALL remain effective; range/enum errors SHALL be reported after required problems are resolved. Validation SHALL NOT mutate settings or include missing/empty field values in the error. The existing HTTP error envelope SHALL remain unchanged.

#### Scenario: All missing required fields
- **WHEN** Modbus TCP receives an empty or nil config
- **THEN** one error identifies host, port and slave_id in schema order.

#### Scenario: Mixed missing and empty fields
- **WHEN** Modbus TCP receives host="", slave_id=1 and no port
- **THEN** the same error identifies missing port and empty host without requiring another request.

#### Scenario: Conditional protocol fields
- **WHEN** a FATEK config has an empty mode and no station_no
- **THEN** one error identifies both problems, while complete serial/network configs retain their existing conditional requirement behavior.

#### Scenario: Valid and invalid existing values
- **WHEN** required fields are complete
- **THEN** valid configurations remain accepted and existing range/enum errors remain rejected without changing their prior message or protocol behavior.
