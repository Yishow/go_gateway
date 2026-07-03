## ADDED Requirements

### Requirement: Runtime dashboard does not invent unsupported runtime values

The runtime dashboard SHALL NOT invent runtime values, statuses, or summaries that are not supplied by the backend contract.

#### Scenario: Unsupported field stays absent instead of guessed

- **WHEN** a dashboard panel does not receive a required backend-supported field
- **THEN** that panel renders an empty or degraded state for that field
- **AND** it SHALL NOT guess a replacement from unrelated client state

##### Example: missing write error metric stays absent

- **GIVEN** the backend response for dev-A omits a write-error detail field that one panel expects
- **WHEN** the panel renders
- **THEN** it shows that field as unavailable instead of deriving a substitute from another counter
