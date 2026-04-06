## MODIFIED Requirements

### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/studio` as the primary product route with persistent workbench regions for step navigation, context, main review surfaces, inspection, and summary feedback.

Any legacy `/datalink/workbench` entry SHALL redirect to `/studio` instead of remaining a separate primary workflow.

#### Scenario: Primary workbench opens through `/studio`
- **WHEN** the operator opens the primary datalink workbench
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not require a separate primary `/datalink/workbench` route for the same workflow

### Requirement: Tag step is review-first for rule-derived mappings
The workbench SHALL treat Step 3 as a review and exception-handling surface for rule-derived tag candidates and pending mapping intent, not as an immediate active tag-and-mapping persistence step.

#### Scenario: Rule save populates Step 3 with candidates
- **WHEN** an operator creates or restores a rule in the primary flow
- **THEN** Step 3 SHALL load rule-derived tag candidates and pending mapping intent
- **AND** SHALL NOT require the system to create the active downstream tag/mapping state before operator review

## ADDED Requirements

### Requirement: Workbench review steps consume rule revision candidates
The workbench SHALL bind Step 3 and Step 4 to candidate snapshots derived from the active source-rule revision.

#### Scenario: Rule revision refreshes downstream review steps
- **WHEN** an operator saves a source rule revision
- **THEN** the workbench refreshes tag and output review surfaces from the new candidate snapshot
- **AND** any previously applied downstream item with the same rule-owned identity but a different proposed signature is shown as `out_of_sync` instead of being silently overwritten

### Requirement: Workbench detects stale candidate views
The workbench SHALL detect when an open Step 3 or Step 4 review surface is showing candidates from an older source-rule revision.

#### Scenario: Stale review state blocks apply
- **WHEN** a newer source-rule revision exists for the rule currently open in Step 3 or Step 4
- **THEN** the workbench marks the open review state as stale
- **AND** blocks apply until the operator refreshes to the latest candidate revision
