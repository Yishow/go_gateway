## ADDED Requirements

### Requirement: Live backend diagnostics for Step 1

The Step 1 device test panel SHALL use real backend diagnostics instead of frontend mock animation.

#### Scenario: Real diagnostics request on execute test

- **GIVEN** the operator is editing a device in Step 1
- **WHEN** the operator clicks `執行測試`
- **THEN** the system sends a real backend diagnostics request for the current draft or persisted device state
- **AND** the UI does not synthesize success from local timers or random values

#### Scenario: Diagnostics success marks the device as tested

- **GIVEN** backend diagnostics returns successful connect and probe outcomes
- **WHEN** the response is received
- **THEN** the device enters `tested` state
- **AND** the continue gate may treat that device as passed

#### Scenario: Diagnostics failure remains actionable

- **GIVEN** backend diagnostics reports connect failure or probe failure
- **WHEN** the response is received
- **THEN** the UI shows the failing stage and message
- **AND** the device does not enter `tested` state
- **AND** the continue gate remains blocked for that device
