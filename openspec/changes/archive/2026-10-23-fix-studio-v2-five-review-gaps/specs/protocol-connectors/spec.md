## ADDED Requirements

### Requirement: Scalar point reads preserve protocol count units
The system SHALL interpret single-point requests as one typed value for MC 3E, while preserving register-count semantics for Modbus and FATEK and the existing typed-value count of direct MC batch requests.

#### Scenario: Multiword MC point across entrypoints
- **WHEN** periodic, manual or read-now acquisition reads one uint32 point at D3 with words 34464 and 1 in CDAB order
- **THEN** it SHALL request two words and return scalar 100000 rather than an array; a single 64-bit value SHALL occupy four words and retain exact encoding.

#### Scenario: Direct batch and other protocols remain unchanged
- **WHEN** a direct MC batch requests two uint32 values or a Modbus/FATEK point is polled
- **THEN** MC SHALL retain two-value batch decoding and other protocols SHALL retain their existing word counts.
