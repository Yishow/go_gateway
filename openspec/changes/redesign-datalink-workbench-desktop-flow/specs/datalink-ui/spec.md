## MODIFIED Requirements

### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps a persistent visual context of the complete datalink flow inside one workbench workspace: **Device -> Source Plan -> Tag Binding -> Output Targets**.

The workflow SHALL ensure:

1. Operators establish device capability context before source planning.
2. Source, tag, and output context remain visible without leaving `/datalink/workbench`.
3. Local Modbus and Database output readiness remain visible in the same workspace.
4. Failures are localized to a specific step, selection, or output target.

#### Scenario: End-to-end guided configuration in one workbench
- **WHEN** an operator selects a device, applies source planning, completes tag binding, and configures an output target
- **THEN** the UI keeps the flow inside `/datalink/workbench`
- **AND** the operator does not need to switch to a separate page to complete the main path

#### Scenario: Step-localized error handling
- **WHEN** planning, binding, or output validation fails
- **THEN** the UI marks the owning step or selected object as failed or blocked
- **AND** provides actionable retry or edit guidance in context

### Requirement: Flow-first workspace visualization

The UI SHALL provide a flow-first workspace with a desktop shell that combines `StepRail`, `ContextBar`, `PrimaryWorkArea`, `InspectorPanel`, and `BottomSummaryBar`.

The workspace SHALL preserve the following flow-reading order:
- Device context
- Source planning and value visualization
- Tag linkage
- Output readiness and target mapping

#### Scenario: Persistent desktop workbench shell
- **WHEN** the operator changes selected devices, source rules, tag bindings, or output targets
- **THEN** the workbench updates the relevant shell regions cohesively
- **AND** preserves a stable reading order across the full flow

#### Scenario: Source-to-output linkage visibility
- **WHEN** a source span becomes linked to a tag and prepared for output
- **THEN** the UI shows that linkage in source, tag, and output contexts
- **AND** displays readiness as `draft`, `ready`, `partial`, `blocked`, or `applied`

### Requirement: Accessible and responsive operator workspace

The UI SHALL remain fully operable by keyboard, preserve i18n compatibility, and avoid horizontal overflow at supported desktop breakpoints.

#### Scenario: Keyboard-only operation in workbench
- **WHEN** an operator uses keyboard-only navigation in `/datalink/workbench`
- **THEN** the operator can switch steps, operate the current toolbar, review selections, and trigger inspector actions
- **AND** all interactive controls expose visible focus state and screen-reader-readable state changes

#### Scenario: Desktop responsive stability
- **WHEN** the workbench is rendered at 1280px or above
- **THEN** the desktop shell shows all core workflow regions without horizontal scrolling
- **AND** the main working surface remains usable at 1920×1080 without overflow traps or collapsed critical controls

### Requirement: Source template library

The UI SHALL allow operators to save, load, update, and delete source planning templates for Step 2 rule groups.

Each template SHALL store at least:
- source rules
- rule ordering and lock state
- preferred Step 2 view mode
- naming defaults when applicable

Template persistence for this workflow SHALL be browser-local in this redesign round.

#### Scenario: Save and reuse source template
- **WHEN** an operator saves the current Step 2 rule set as a template
- **THEN** the template is persisted locally
- **AND** the operator can later load it to prefill planning controls

#### Scenario: Warn on capability mismatch
- **WHEN** an operator applies a saved template to a device with different capability assumptions
- **THEN** the UI warns about mismatches such as `address base` or `word order`
- **AND** requires the operator to confirm before applying the template
