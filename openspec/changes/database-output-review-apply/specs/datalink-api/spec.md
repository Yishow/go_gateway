## ADDED Requirements

### Requirement: Rule-scoped database output apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved database output candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/database-outputs/apply` - Apply approved database output candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the database candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply database outputs without affecting Local Modbus
- **WHEN** a client applies database output candidates for a source rule
- **THEN** the API applies only the database target changes
- **AND** preserves Local Modbus candidate and apply state for the same rule

#### Scenario: Partial database apply returns per-item results
- **WHEN** a client applies multiple database candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result
