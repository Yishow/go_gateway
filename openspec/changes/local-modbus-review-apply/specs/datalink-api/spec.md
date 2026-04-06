## ADDED Requirements

### Requirement: Rule-scoped Local Modbus apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved Local Modbus candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/local-modbus/apply` - Apply approved Local Modbus candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the Local Modbus candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply Local Modbus outputs without affecting Database
- **WHEN** a client applies Local Modbus candidates for a source rule
- **THEN** the API applies only the Local Modbus target changes
- **AND** preserves database candidate and apply state for the same rule

#### Scenario: Partial Local Modbus apply returns per-item results
- **WHEN** a client applies multiple Local Modbus candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result
