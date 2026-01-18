# datalink-ui Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Guided workflow
The UI SHALL provide a guided workflow from device setup to mapping activation.

#### Scenario: Complete workflow
- WHEN a user completes the steps
- THEN the mapping is activated and data starts flowing

### Requirement: Drag-drop mapping canvas
The UI SHALL provide a drag-drop canvas to connect source points to target tags and mapping steps.

#### Scenario: Create mapping by drag
- WHEN a user drags a point onto a tag
- THEN a mapping is created and editable

### Requirement: Transform builder
The UI SHALL provide a transform builder with ordered steps, parameters, and validation.

#### Scenario: Configure scaling step
- WHEN a user adds a scaling step with parameters
- THEN the UI validates and saves the step

### Requirement: Live preview
The UI SHALL provide a live preview of raw and transformed values for a selected mapping.

#### Scenario: View preview
- WHEN a user opens preview
- THEN the UI displays raw, step results, and final value

### Requirement: Write precision settings
The UI SHALL allow operators to configure write timestamp precision (seconds or milliseconds).

#### Scenario: Select millisecond precision
- WHEN a user selects millisecond precision
- THEN the UI saves the setting for subsequent writes

