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
The system SHALL allow tags to link to one or more source points through mapping rules.

#### Scenario: Link tag to point
- WHEN a mapping is created
- THEN the tag references the source point and mapping id

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

