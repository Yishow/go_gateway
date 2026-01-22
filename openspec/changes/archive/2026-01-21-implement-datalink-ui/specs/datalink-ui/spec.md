## MODIFIED Requirements

### Requirement: Guided workflow

The UI SHALL provide a guided workflow from device setup to mapping activation using a 6-step Wizard with the following stages:

1. **Device Selection**: Select existing device or create new
2. **Point Configuration**: Configure address and data type
3. **Tag Selection**: Select target tag or create new
4. **Transform Configuration**: Build transform pipeline
5. **Preview Confirmation**: Verify raw → transformed values via SSE
6. **Activation**: Save and activate mapping

The workflow SHALL support:

- Progress persistence (draft save)
- Step validation before proceeding
- Back navigation to any completed step
- Side panel with live preview

#### Scenario: Complete workflow

- WHEN a user completes all 6 steps
- THEN the mapping is activated and data starts flowing

#### Scenario: Save draft and resume

- WHEN a user saves a draft at any step
- THEN the user can resume from the same step later

#### Scenario: Step validation failure

- WHEN a user attempts to proceed with invalid data
- THEN the UI displays validation errors and blocks progression

---

### Requirement: Drag-drop mapping canvas

The UI SHALL provide a drag-drop canvas to connect source points to target tags and mapping steps.

#### Scenario: Create mapping by drag

- WHEN a user drags a point onto a tag
- THEN a mapping is created and editable

#### Scenario: Edit existing mapping

- WHEN a user clicks on an existing mapping
- THEN the mapping details are shown for editing

---

### Requirement: Transform builder

The UI SHALL provide a transform builder with ordered steps, parameters, and validation.

The builder SHALL support the following transform types:

- `decode`: Decode raw bytes to typed value
- `cast`: Type conversion
- `scale`: Linear scaling with parameters (multiplier, offset)
- `lookup`: Table lookup mapping
- `conditional`: Conditional branching
- `formula`: Expression-based calculation

#### Scenario: Configure scaling step

- WHEN a user adds a scaling step with parameters
- THEN the UI validates and saves the step

#### Scenario: Reorder transform steps

- WHEN a user drags a step to a new position
- THEN the pipeline order is updated

#### Scenario: Pipeline validation

- WHEN a user saves the transform pipeline
- THEN the UI calls `/mappings/validate-pipeline` to verify

---

### Requirement: Live preview

The UI SHALL provide a live preview of raw and transformed values for a selected mapping using Server-Sent Events (SSE).

The preview panel SHALL display:

- Raw value from device
- Intermediate step results
- Final transformed value
- Quality indicator

#### Scenario: View preview

- WHEN a user opens preview
- THEN the UI displays raw, step results, and final value

#### Scenario: Live update via SSE

- WHEN device values change
- THEN the preview panel updates in real-time

#### Scenario: SSE reconnection

- WHEN the SSE connection is lost
- THEN the UI automatically reconnects

---

### Requirement: Write precision settings

The UI SHALL allow operators to configure write timestamp precision (seconds or milliseconds).

#### Scenario: Select millisecond precision

- WHEN a user selects millisecond precision
- THEN the UI saves the setting for subsequent writes

#### Scenario: Configure partition interval

- WHEN a user selects a partition interval (daily/weekly/monthly)
- THEN the UI saves the setting for time-series storage
