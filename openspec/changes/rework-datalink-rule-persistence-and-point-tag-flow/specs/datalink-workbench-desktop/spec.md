## ADDED Requirements

### Requirement: Device testing distinguishes transport and protocol phases
The workbench SHALL show and judge transport connectivity and protocol probe results as separate outcomes in Step 1.

#### Scenario: Connect succeeds but probe fails
- **WHEN** an operator tests a device and transport connect succeeds while protocol probe fails
- **THEN** Step 1 SHALL show connect as successful and probe as failed
- **AND** SHALL allow the device to be saved
- **AND** SHALL mark rule activation and data collection as blocked until probe succeeds

#### Scenario: Route failure is classified as connect-stage failure
- **WHEN** a device test fails with a network error such as `no route to host`
- **THEN** Step 1 SHALL classify the failure as a connect-stage failure
- **AND** SHALL present diagnostics without mislabeling the protocol probe as the failing stage

### Requirement: Source planning reflects device capability and persisted rule state
The workbench SHALL align Step 2 planning behavior with the selected device capability context and persisted rule state.

#### Scenario: Device capability updates planner behavior
- **WHEN** an operator selects a device with different address model or protocol traits
- **THEN** Step 2 SHALL update available rule options, address semantics, merge behavior, and warning/block behavior to match that device

#### Scenario: Existing unmanaged points are explained in the grid
- **WHEN** Step 2 loads points that exist without an active persisted rule association
- **THEN** the grid SHALL distinguish them from rule-planned spans
- **AND** SHALL explain whether each address is unmanaged used state, conflict state, or rule-derived state

### Requirement: Tag step is review-first for rule-derived mappings
The workbench SHALL treat Step 3 as a review and exception-handling surface for automatically created rule-derived tags and mappings.

#### Scenario: Rule creation pre-populates review surface
- **WHEN** an operator creates or restores a rule in the primary flow
- **THEN** Step 3 SHALL load the generated tag identity and mapping state without requiring a manual first-pass bind
- **AND** SHALL show the generated status for operator review

#### Scenario: Review step supports exception handling
- **WHEN** automatic tag or mapping generation fails or needs correction
- **THEN** Step 3 SHALL surface the mismatch or failure inline
- **AND** SHALL allow corrective review actions without reverting the entire workflow to manual binding
