# datalink-workbench-desktop Specification

## Purpose
TBD - created by archiving change redesign-datalink-workbench-desktop-flow. Update Purpose after archive.
## Requirements
### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/studio` as the primary product route with persistent workbench regions for step navigation, context, main review surfaces, inspection, and summary feedback.

Any legacy `/datalink/workbench` entry SHALL redirect to `/studio` instead of remaining a separate primary workflow.

#### Scenario: Primary workbench opens through `/studio`
- **WHEN** the operator opens the primary datalink workbench
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not require a separate primary `/datalink/workbench` route for the same workflow

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
The system SHALL provide one Output workspace inside `/studio` that keeps Local Modbus and Database outputs in the same workbench step.

#### Scenario: Operator changes output target without leaving `/studio`
- **WHEN** the operator switches between `Local Modbus` and `Database` in the Output step
- **THEN** the output review surface remains in the same `/studio` route
- **AND** the target-specific studio changes without requiring navigation to a separate primary page

#### Scenario: Output target switching preserves one primary workflow
- **WHEN** legacy navigation or bookmarks would previously open a separate workbench route for output work
- **THEN** the system keeps the operator inside `/studio`
- **AND** treats any retained legacy entry as a compatibility redirect rather than a second primary workflow

### Requirement: Device testing distinguishes transport and protocol phases
The workbench SHALL show and judge transport connectivity and protocol probe results as separate outcomes in Step 1.

#### Scenario: Connect succeeds but probe fails
- **WHEN** an operator tests a device and transport connect succeeds while protocol probe fails
- **THEN** Step 1 SHALL show connect as successful and probe as failed
- **AND** SHALL allow the device to be saved
- **AND** SHALL mark rule activation and data collection as blocked until probe succeeds

#### Scenario: Route failure is classified as connect-stage failure
- **WHEN** a device test fails with a network error such as `no route to host`
- **THEN** Step 1 SHALL classify the failure as a connect-stage failure
- **AND** SHALL present diagnostics without mislabeling the protocol probe as the failing stage

### Requirement: Source planning reflects device capability and persisted rule state
The workbench SHALL align Step 2 planning behavior with the selected device capability context and persisted rule state.

#### Scenario: Device capability updates planner behavior
- **WHEN** an operator selects a device with different address model or protocol traits
- **THEN** Step 2 SHALL update available rule options, address semantics, merge behavior, and warning/block behavior to match that device

#### Scenario: Existing unmanaged points are explained in the grid
- **WHEN** Step 2 loads points that exist without an active persisted rule association
- **THEN** the grid SHALL distinguish them from rule-planned spans
- **AND** SHALL explain whether each address is unmanaged used state, conflict state, or rule-derived state

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

### Requirement: Source rule builder exposes target data type selection
The system SHALL provide, within Step 2 source planning, a control to set the **target / intended tag data type** alongside the protocol read data type when creating or editing a source rule. The control SHALL default to matching the protocol read data type and SHALL clearly distinguish protocol read semantics from target output semantics.

#### Scenario: Operator sets distinct target type before applying rule
- **WHEN** an operator configures a new rule and selects a target data type different from the protocol read type
- **THEN** the workbench SHALL persist that intent with the rule according to backend contract
- **AND** SHALL show validation feedback when the pair is unsupported

#### Scenario: Defaults avoid extra friction
- **WHEN** an operator does not change the default target data type
- **THEN** the workbench SHALL behave as today with a single data type selection driving protocol read planning

### Requirement: Device connection form exposes data format when the protocol uses multi-register decoding
The workbench SHALL render **`data_format`** (or the documented equivalent) in the device connection editor for each protocol whose connector uses multi-register / floating-point byte ordering (for example **Modbus TCP, Modbus RTU, Modbus UDP**, and **Mitsubishi MC 3E**). The control SHALL list the supported orderings (for example **ABCD, BADC, CDAB, DCBA**) and SHALL persist the value in `connection_config` for the backend to consume during reads.

#### Scenario: Modbus connection shows byte order
- **WHEN** an operator edits a Modbus-class device connection
- **THEN** the UI SHALL expose `data_format` selection
- **AND** saving the device SHALL persist the chosen value

### Requirement: Source rule builder MAY expose scale or offset for unit conversion
The system SHALL provide optional inputs for **linear scaling** (multiplier and/or offset) on the source rule form when the product enables engineering-unit conversion from rule planning, with clear labels that distinguish scaling from protocol read type and from target tag type.

#### Scenario: Operator sets scale without changing target type
- **WHEN** an operator enters scale parameters only
- **THEN** the workbench SHALL persist scale intent with the rule
- **AND** SHALL not require a target type different from the read type
