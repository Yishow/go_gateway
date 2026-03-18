# datalink-workbench-desktop Specification

## Purpose
TBD - created by archiving change redesign-datalink-workbench-desktop-flow. Update Purpose after archive.
## Requirements
### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/datalink/workbench` with five persistent regions: `StepRail`, `ContextBar`, `PrimaryWorkArea`, `InspectorPanel`, and `BottomSummaryBar`.

#### Scenario: 1920 desktop shell stays readable
- **WHEN** the operator opens `/datalink/workbench` on a 1920×1080 desktop viewport
- **THEN** the workbench shows the five-region shell without horizontal overflow
- **AND** the main working surface is not constrained by a centered `max-width` content box

#### Scenario: Shared inspector replaces page-level summary dock
- **WHEN** the operator switches between Device, Source, Tag, and Output steps
- **THEN** the same desktop `InspectorPanel` hosts the active detail surface for the current selection
- **AND** the page does not rely on a separate persistent right-side action dock as a shell region

### Requirement: Device workspace establishes capability context
The system SHALL provide a `DeviceWorkspace` that lets the operator select or create a source device and understand the capability context required for source planning.

#### Scenario: Selected device exposes planning-critical capabilities
- **WHEN** the operator selects a device in Step 1
- **THEN** the UI shows `address base`, `word order`, `unit id`, and protocol traits in the workbench context
- **AND** those capability hints remain visible in the device inspector and context bar

#### Scenario: Clone flow starts from selected device defaults
- **WHEN** the operator chooses to clone the selected device
- **THEN** the UI opens a workbench-specific create flow prefilled with reusable connection defaults
- **AND** still requires a distinct device identity before save

### Requirement: Address canvas supports rule-driven continuous planning
The system SHALL render source planning as a continuous 16-bit lattice that can expand multiple rules into merged spans on one shared canvas.

#### Scenario: Typed rules expand into merged spans
- **WHEN** the operator applies a rule with `32-bit`, `float32`, or `64-bit` data width
- **THEN** the canvas merges the correct number of adjacent 16-bit cells per item
- **AND** the operator can visually distinguish merged spans from empty cells

#### Scenario: Multiple rules share one continuous address range
- **WHEN** the operator defines multiple source rules with different start addresses
- **THEN** the canvas spans from the minimum rule start address to the maximum rule end address
- **AND** gaps between rule-covered ranges remain visible as unplanned address space

### Requirement: Address canvas uses stable multi-mode visualization
The system SHALL let the operator switch between `Plan`, `Live`, and `Link` modes without replacing the underlying address geometry.

#### Scenario: View-mode switching preserves geometry
- **WHEN** the operator switches between `Plan`, `Live`, and `Link`
- **THEN** the address positions, merged spans, and current selection remain geometrically stable
- **AND** only the information layer changes between planning, value validation, and downstream linkage states

#### Scenario: Toolbar exposes live validation controls
- **WHEN** the operator uses Step 2 toolbar controls
- **THEN** the toolbar provides value-format switching, live freeze, snapshot compare, jump-to-address, and template load/save actions
- **AND** those controls apply to the current canvas state without leaving the workbench route

### Requirement: Tag binding board preserves source-to-tag review context
The system SHALL provide a `TagBindingBoard` that exposes source context, proposed tag identity, and binding status densely enough for batch review.

#### Scenario: Source-to-tag information is visible before inspector open
- **WHEN** the operator reviews items in Step 3
- **THEN** each row or card shows tag key preview, source address/span, raw vs transformed value, data width context, and binding status
- **AND** the operator can compare multiple items without opening them one by one

#### Scenario: Batch review includes diff and result feedback
- **WHEN** the operator performs a batch tag action
- **THEN** the UI shows a pre-apply diff preview before confirmation
- **AND** shows a post-apply result summary with created, linked, skipped, and failed items

### Requirement: Unified output workspace supports both targets
The system SHALL provide one `OutputWorkspace` that keeps Local Modbus and Database outputs inside the same workbench step.

#### Scenario: Operator changes output target without leaving workbench
- **WHEN** the operator switches between `Local Modbus` and `Database` in Step 4
- **THEN** the output candidate board remains in the same workbench route
- **AND** the target-specific studio changes without requiring navigation to another page

#### Scenario: Candidate readiness tracks both targets independently
- **WHEN** an output candidate is ready for one target but not the other
- **THEN** the UI displays a distinct readiness state for each target and an overall readiness summary
- **AND** the inspector explains why a candidate is `partial` or `blocked`

