## ADDED Requirements

### Requirement: Persisted source-rule lifecycle
The system SHALL persist each source rule as a first-class record with device association, planning definition, and enabled/disabled state.

#### Scenario: Save rule definition
- **WHEN** an operator creates or updates a source rule
- **THEN** the system stores the rule in the database with its device association, planning inputs, and current enabled/disabled state
- **AND** the saved rule can be reloaded independently of the current browser session

#### Scenario: Restore rule state after restart
- **WHEN** the service restarts
- **THEN** the system reloads persisted source rules from the database
- **AND** restores each rule's enabled/disabled state before resuming runtime behavior

### Requirement: Rule enablement controls collection without deleting derived relationships
The system SHALL treat rule enablement as collection control only and SHALL preserve Point, Tag, Mapping, and Output relationships when a rule is disabled.

#### Scenario: Disable rule preserves derived records
- **WHEN** an operator disables a rule
- **THEN** the system stops collection for that rule
- **AND** preserves the derived Point, Tag, Mapping, and Output relationships in storage

#### Scenario: Re-enable rule resumes existing definition
- **WHEN** an operator re-enables a previously disabled rule
- **THEN** the system resumes collection using the stored rule definition and derived relationships
- **AND** does not require the operator to recreate the rule or downstream bindings

### Requirement: Rule activation is gated by protocol probe success
The system SHALL block rule activation and collection when the associated device has not passed protocol probe validation.

#### Scenario: Connect succeeds but probe fails
- **WHEN** a device test reports transport connect success and protocol probe failure
- **THEN** the operator MAY save the device
- **AND** the system SHALL reject rule activation and data collection until probe succeeds

### Requirement: Rule-driven live values appear in the source grid
The system SHALL display parsed live values for enabled rules in the Step 2 source grid.

#### Scenario: Enabled rule updates grid
- **WHEN** runtime collection succeeds for an enabled source rule
- **THEN** the Step 2 grid shows the parsed value, timestamp, and error state for the rule-derived span
- **AND** the displayed value follows the rule's configured data-type and merge semantics
