## ADDED Requirements

### Requirement: Workspace runtime view stays aligned to the selected device truth

The workspace-scoped runtime view SHALL stay aligned to the selected device truth from runtime backend data.

#### Scenario: Workspace runtime view refuses cross-device synthetic summary

- **WHEN** the selected device lacks runtime data while another device has runtime data
- **THEN** the workspace runtime view reports the selected device as missing or degraded
- **AND** it SHALL NOT display another device's runtime summary as a substitute

##### Example: runtime view keeps selected scope truthful

- **GIVEN** selected device dev-A has no runtime summary and dev-B has one
- **WHEN** the workspace runtime view loads for dev-A
- **THEN** the view keeps dev-A in missing or degraded state instead of borrowing dev-B data
