## ADDED Requirements

### Requirement: Tag is the primary semantic data object

The workbench SHALL treat `Tag` as the primary operator-facing semantic data object, while `Point` remains an internal acquisition unit.

#### Scenario: Destination planning uses tags instead of points
- **WHEN** an operator enters the destination flow
- **THEN** the primary planning unit is derived from tags and delivery groups
- **AND** the operator is not forced into a point-first flow

### Requirement: Delivery groups drive destination planning

The workbench SHALL use `Delivery Group` as the primary destination planning unit for database and share destinations.

#### Scenario: Delivery groups are suggested from tag groups
- **WHEN** an operator completes tag refinement
- **THEN** the system suggests delivery groups based on tag groups
- **AND** allows the operator to adjust them before applying a destination projection
