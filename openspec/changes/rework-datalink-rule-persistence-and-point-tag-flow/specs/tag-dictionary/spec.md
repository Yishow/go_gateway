## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Rules can generate tags and mappings automatically
The system SHALL support automatic tag creation and mapping generation from the primary source-rule flow.

#### Scenario: Auto-create tag and mapping from rule
- **WHEN** an operator creates or enables a source rule in the primary flow
- **THEN** the system SHALL create or synchronize the corresponding tag
- **AND** SHALL create the active mapping to the derived point automatically

#### Scenario: Generated tag identity combines rule name and address
- **WHEN** the system generates a tag identity from a rule
- **THEN** the generated tag name or key SHALL combine the rule name and source address context
- **AND** SHALL remain unique under the global tag dictionary rules
