# tag-dictionary Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Global tag dictionary
The system SHALL maintain a global tag dictionary with unique tag keys and optional namespace segments.

#### Scenario: Create tag
- WHEN an operator creates a tag key
- THEN the tag is stored and enforced as unique

### Requirement: Tag key format
The system SHALL accept tag keys composed of ASCII letters, digits, underscore, dash, dot, and slash, with length 1 to 128 characters.

#### Scenario: Reject invalid characters
- WHEN a tag key contains a space or an unsupported character
- THEN the system rejects the tag with a validation error

### Requirement: Tag key uniqueness
The system SHALL enforce tag key uniqueness in a case-insensitive manner.

#### Scenario: Reject case-variant tag
- WHEN a tag key already exists with different letter casing
- THEN the system rejects the new tag as a duplicate

### Requirement: Tag metadata
The system SHALL store metadata for each tag including display name, unit, data type, and labels.

#### Scenario: Update metadata
- WHEN an operator updates a tag unit and labels
- THEN the metadata is saved and visible in UI and API

### Requirement: Tag lifecycle
The system SHALL allow tags to be set to draft, active, or retired.

#### Scenario: Retire tag
- WHEN a tag is retired
- THEN new mappings to that tag are prevented

### Requirement: Tag-to-source linkage
The system SHALL allow each tag to link to exactly one source point through one active mapping, and each source point SHALL link to at most one active tag through one active mapping.

#### Scenario: Link tag to point
- **WHEN** the system or an operator creates a mapping between an unlinked point and an unlinked tag
- **THEN** the tag SHALL reference that one source point and mapping id
- **AND** the point SHALL be marked as linked to that one active tag

#### Scenario: Reject second point for same tag
- **WHEN** another point attempts to create an active mapping to a tag that already has an active source mapping
- **THEN** the system SHALL reject the new mapping as a cardinality violation
- **AND** SHALL preserve the existing active linkage

### Requirement: In-context global tag editing

The system SHALL allow operators to edit global tag metadata directly from Pipeline Studio linkage context.

Any saved change SHALL update the global tag dictionary, not only local draft state.

#### Scenario: Edit tag in grid linkage context
- **WHEN** an operator edits display name, unit, or labels for a linked tag from Pipeline Studio
- **THEN** the global tag record is updated
- **AND** all references show updated metadata after refresh

#### Scenario: Show global impact before save
- **WHEN** a tag edit affects existing mappings
- **THEN** the UI shows affected mapping count before confirmation
- **AND** requires explicit confirmation for save

### Requirement: Rules can generate tags and mappings automatically
The system SHALL generate rule-derived tag candidates and pending mapping intent from the primary source-rule flow before apply, rather than immediately creating active downstream tag-and-mapping state on rule save.

#### Scenario: Rule save produces reviewable tag candidates
- **WHEN** an operator creates or enables a source rule in the primary flow
- **THEN** the system SHALL generate the corresponding tag candidate and pending mapping intent
- **AND** SHALL defer active downstream persistence until the operator completes tag review/apply

#### Scenario: Apply approved tag candidate
- **WHEN** an operator approves a rule-generated tag candidate for apply
- **THEN** the system creates or updates the active tag and active mapping state from the approved candidate
- **AND** preserves the approved review decision as part of the effective tag state

#### Scenario: Generated tag identity uses canonical rule-owned identity
- **WHEN** the system generates a tag candidate from a rule
- **THEN** the tag candidate SHALL use the rule-owned identity contract for tag candidates
- **AND** SHALL remain unique under the global tag dictionary rules

### Requirement: Rule-generated tags preserve review decisions across compatible revisions
The system SHALL preserve rule-scoped rename, skip, and override decisions for rule-generated tags while the candidate identity remains materially the same.

For this workflow, a rule-generated tag candidate is materially the same only when all of the following remain unchanged:

- `source_rule_id`
- `derived_from_rule_address`
- derived target data type
- candidate kind (`tag`)

#### Scenario: Rename survives a compatible rule revision
- **WHEN** an operator renames a rule-generated tag and a later rule revision keeps the same rule-owned tag identity
- **THEN** the rename remains as the effective tag decision
- **AND** the system recomputes downstream output candidates from that effective tag state

#### Scenario: Identity change invalidates prior review decisions
- **WHEN** a later rule revision changes a rule-generated tag candidate so it is no longer materially the same
- **THEN** the system marks prior review decisions as stale for that candidate
- **AND** requires the operator to review the new candidate before apply

