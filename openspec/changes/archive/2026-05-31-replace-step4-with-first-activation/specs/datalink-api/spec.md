## ADDED Requirements

### Requirement: Studio v2 workspace first activation API

The system SHALL expose `POST /api/v1/datalink/studio-v2/workspace/activate` as the first activation endpoint for Step 4.

#### Scenario: Activate eligible workspace devices

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/activate`
- **THEN** the API activates all valid, available, not-yet-running devices in the singleton workspace
- **AND** the response returns one result entry per attempted device

#### Scenario: Partial failure does not roll back successes

- **WHEN** one device activation succeeds and another fails in the same request
- **THEN** the response reports both outcomes separately
- **AND** the already successful activation remains active

##### Example: one success and one timeout

- **GIVEN** `device-A` is eligible and starts successfully while `device-B` is eligible but times out during activation
- **WHEN** the client calls `POST /api/v1/datalink/studio-v2/workspace/activate`
- **THEN** the response includes one `success` result for `device-A`
- **AND** one `failed` result for `device-B`
- **AND** `device-A` remains running after the request returns
