## ADDED Requirements

### Requirement: Unimplemented recording operations fail closed
The recording setup API SHALL reject unimplemented schema-apply and test-write operations with HTTP 501 and distinct stable error codes. It MUST NOT report a successful operation, fabricate delivery evidence, or mutate the target database. Existing request validation SHALL remain effective.

#### Scenario: Nonempty schema token reaches an unimplemented operation
- **WHEN** a syntactically valid schema-apply request reaches an unimplemented service
- **THEN** the response is 501 with code RECORDING_SCHEMA_NOT_IMPLEMENTED and success false
- **AND** it contains neither applied true nor evidence of executed statements.

#### Scenario: Well-formed test request reaches an unimplemented operation
- **WHEN** a syntactically valid test-write request reaches an unimplemented service
- **THEN** the response is 501 with code RECORDING_TEST_WRITE_NOT_IMPLEMENTED
- **AND** no record identifier is fabricated and no target write occurs.

#### Scenario: Invalid input remains invalid
- **WHEN** a request is missing a required body or required field
- **THEN** existing request-validation rejection occurs before operation dispatch
- **AND** neither a success response nor a target mutation is produced.

### Requirement: Capability flags describe connected implementations
The API SHALL advertise each operation as available only when both its adapter and its connected implementation support it. The UI MUST explain unavailable actions, and direct API calls MUST enforce the same limitation independently of the UI.

#### Scenario: Connection succeeds while test writes are unimplemented
- **WHEN** a valid connector supports connections but its recording test-write operation is not connected
- **THEN** connection success remains visible while supports_test_writes is false
- **AND** the UI does not enable a successful-test claim.

#### Scenario: Existing custom-table operation remains implemented
- **WHEN** the recording-plan path is disabled but a separate custom-table operation is implemented and authorized
- **THEN** the custom-table operation retains its applicable validation and availability
- **AND** the recording-plan limitation does not disable all database output.

### Requirement: Safe errors and uncertain outcomes remain distinguishable
Operator errors SHALL contain a stable code, safe message, retryable flag, repair action and opaque request identifier without secrets or raw database diagnostics. An unimplemented operation SHALL be non-retryable. Transport uncertainty MUST NOT be interpreted as either a confirmed success or proof of no side effect.

#### Scenario: Operation is not implemented
- **WHEN** the client receives a recording-operation 501 error
- **THEN** it displays an unavailable-action explanation with no automatic retries
- **AND** it does not change saved, applied or verified status to successful.

#### Scenario: Connection is lost after dispatch
- **WHEN** the client loses the response to a potentially mutating operation
- **THEN** it displays an unconfirmed outcome and retains the operation identity when supplied
- **AND** it does not silently submit another mutation under a new identity.

#### Scenario: Malformed or unfamiliar result
- **WHEN** a nominally successful response lacks required evidence or contains an unknown status
- **THEN** the UI displays an unconfirmed or invalid-response state rather than a green success state.
