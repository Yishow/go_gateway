## MODIFIED Requirements

### Requirement: Tag step is review-first for rule-derived mappings
The workbench SHALL treat Step 3 as a review and exception-handling surface for rule-derived tag candidates and pending mapping intent, not as an immediate active tag-and-mapping persistence step.

#### Scenario: Rule creation pre-populates review surface
- **WHEN** an operator creates or restores a rule in the primary flow
- **THEN** Step 3 SHALL load the generated tag candidates and pending mapping intent without requiring a manual first-pass bind
- **AND** SHALL show the generated status for operator review

#### Scenario: Review step supports exception handling
- **WHEN** automatic tag generation needs correction or approval
- **THEN** Step 3 SHALL surface rename, skip, override, and apply actions inline
- **AND** SHALL allow corrective review actions without reverting the entire workflow to manual binding

## ADDED Requirements

### Requirement: Tag review consumes source-rule revision candidates
The workbench SHALL bind Step 3 to tag candidates derived from the active source-rule revision.

#### Scenario: Tag review uses the active revision
- **WHEN** the operator opens Step 3 for a source rule
- **THEN** the workbench loads the tag candidate set for the current active revision
- **AND** does not mix candidates from older revisions into the same review session

### Requirement: Tag review detects stale candidate views
The workbench SHALL detect when Step 3 is showing candidates from an older source-rule revision.

#### Scenario: Stale tag review blocks apply
- **WHEN** a newer source-rule revision exists for the rule currently open in Step 3
- **THEN** the workbench marks the review state as stale
- **AND** blocks tag apply until the operator refreshes to the latest candidate revision
