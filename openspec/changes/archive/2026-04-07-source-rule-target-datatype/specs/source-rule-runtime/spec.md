# source-rule-runtime Specification Delta

## ADDED Requirements

### Requirement: Persisted rule MAY store intended target data type

The system SHALL persist an optional **target data type** field on each source rule record. The field SHALL be reloadable after service restart and SHALL participate in rule enable/disable lifecycle without being silently discarded.

#### Scenario: Reload after restart preserves target intent

- **WHEN** the service restarts after a rule with a non-null target data type was saved
- **THEN** the system SHALL reload the target data type with the rest of the rule definition
- **AND** SHALL apply the same synchronization behavior as before restart

### Requirement: Rule synchronization derives tag type and default mapping pipeline from target intent

When synchronizing a rule to Points, Tags, and Mappings, the system SHALL set **Tag.data_type** to the rule’s target data type (or protocol read type when target is unset). When target differs from the derived Point read type, the system SHALL ensure the mapping includes a validated default **`cast`** transform pipeline bridging point semantics to the tag type, without requiring manual pipeline entry for the cast alone.

#### Scenario: New tag uses target type when provided

- **WHEN** a rule declares a target data type and synchronization creates a new tag for a derived point
- **THEN** the created tag SHALL use the target data type
- **AND** the point SHALL retain the protocol read data type

#### Scenario: Mapping receives default cast when types differ

- **WHEN** target data type differs from the point read data type for a rule-derived mapping
- **THEN** the mapping SHALL include a `cast` step appropriate to the declared target type
- **AND** the pipeline SHALL pass mapping validation

### Requirement: Persisted rule MAY store optional scale parameters

The system SHALL persist optional **linear scale** fields on a source rule (multiplier and/or offset, or an equivalent single structure) when the operator configures engineering-unit conversion. When unset, the system SHALL treat scaling as identity and SHALL NOT insert a redundant `scale` step.

#### Scenario: Scale fields persist across restarts

- **WHEN** a rule saved with non-default scale parameters is reloaded after restart
- **THEN** the system SHALL restore those parameters
- **AND** SHALL regenerate or validate the derived mapping `scale` step consistently with the saved values
