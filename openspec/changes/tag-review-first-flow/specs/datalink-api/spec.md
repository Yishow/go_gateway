## ADDED Requirements

### Requirement: Rule-scoped tag apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved tag candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/tags/apply` - Apply approved tag candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the tag candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply approved tag candidates
- **WHEN** a client applies tag candidates for a source rule
- **THEN** the API applies only the approved tag candidates for that revision
- **AND** leaves database and Local Modbus outputs as downstream review targets instead of silently applying them

#### Scenario: Partial tag apply returns per-item results
- **WHEN** a client applies multiple tag candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result
