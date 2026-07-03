## ADDED Requirements

### Requirement: Post-setup runtime route shows truthful empty and degraded states

The post-setup runtime route SHALL show truthful empty and degraded states when runtime data is missing or unavailable.

#### Scenario: Route opens without real runtime data

- **WHEN** an operator opens the runtime route for a selected device that has no current runtime data yet
- **THEN** the page shows a truthful waiting, empty, or degraded state
- **AND** the page SHALL NOT present a synthetic ready dashboard just because the route resolved successfully

##### Example: direct link resolves but selected device is still cold

- **GIVEN** the operator opens /studio/runtime for dev-A immediately after activation but runtime has not produced real data yet
- **WHEN** the route resolves successfully
- **THEN** the page shows waiting or empty state for dev-A instead of a synthetic ready dashboard
