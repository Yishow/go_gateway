## ADDED Requirements
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
