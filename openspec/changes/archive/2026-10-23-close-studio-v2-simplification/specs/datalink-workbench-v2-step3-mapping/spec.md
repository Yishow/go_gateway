## MODIFIED Requirements

### Requirement: Inline-editable mapping table

The Step 3 mapping table SHALL render one row per acquired, non-skipped point. Basic presentation SHALL expose source/address, display name, unit, current final value and mapping enabled state. Advanced presentation SHALL retain Tag Key, target type, scale, offset, intermediate values and API payload inspection. Switching presentation SHALL preserve every mapping draft and its type precision without rebuilding or clearing configuration. Editing or focusing a mapping field SHALL identify that row as the explicit bulk source; checking a bulk target SHALL NOT silently change the source. The mapping enabled control SHALL remain distinct from point acquisition and group recording membership.

#### Scenario: Default row rendering
- **GIVEN** eight acquired points with persisted mappings
- **WHEN** Step 3 mounts in Basic
- **THEN** eight rows expose the five basic field groups with persisted values and source-aware defaults, without demo-specific names, units or forced float64 conversion.

#### Scenario: Row selection by clicking blank area
- **WHEN** the operator clicks a row body
- **THEN** that row becomes the visibly identified source for an explicit bulk action.

#### Scenario: Editing identifies the actual source
- **GIVEN** the current source is row 0
- **WHEN** the operator focuses or edits a mapping input in row 3
- **THEN** row 3 becomes the source and its edit updates only that mapping draft.

#### Scenario: Basic and Advanced retain drafts
- **WHEN** the operator edits Tag Key, target type, scale or offset in Advanced and returns to Basic and back
- **THEN** the same values and current preview remain intact, including exact int64 and uint64 representations.

#### Scenario: Toggle enabled does not select row
- **GIVEN** the current source is row 0
- **WHEN** the operator clicks the enabled control or bulk-target checkbox in row 2
- **THEN** the source remains row 0 and only the explicitly requested enabled or target-selection state changes.


### Requirement: Bulk apply transform parameters

The mapping table SHALL provide an explicit source row, selectable enabled target rows and a before/after difference preview for bulk unit, scale, offset and target-type application. All enabled targets SHALL remain selectable through an explicitly labelled all-target action; disabled mappings and the source SHALL be excluded. No bulk mutation SHALL occur before review and explicit apply. Changes to source, targets or their relevant values SHALL invalidate an earlier review. Cancelling SHALL preserve all original values. The reducer SHALL require explicit target point IDs and preserve unrelated fields.

#### Scenario: Bulk apply propagates fields
- **GIVEN** a temperature source with scale 0.5 and offset 2, plus pressure and energy rows
- **WHEN** the operator selects only the pressure row, reviews its old/new unit and transform, and explicitly applies
- **THEN** only that target receives the source unit, scale, offset and target type; energy, disabled rows, Tag Keys, display names and enabled states remain unchanged.

#### Scenario: Bulk apply with no selection
- **WHEN** no source or no target is selected
- **THEN** no mapping is modified and the UI explains the missing selection.

#### Scenario: Explicit all enabled targets
- **WHEN** the operator selects all eligible enabled targets
- **THEN** the review identifies that complete scope and differences before confirmation, without silently adding disabled rows.

#### Scenario: Cancel or stale review
- **WHEN** the operator cancels review or a relevant source/target value changes before apply
- **THEN** no stale bulk mutation occurs and another explicit review is required.


### Requirement: Target type quick actions

Step 3 SHALL retain single-row target type editing through its Advanced target-type control, including bool, int16, float64 and string and all other supported types. Applying only target_type to multiple rows SHALL use the explicit source, selected enabled targets and reviewed confirmation contract of Bulk apply transform parameters; it SHALL NOT immediately update all rows.

#### Scenario: Quick action updates the current row
- **GIVEN** a mapping row is selected
- **WHEN** the operator changes its target type
- **THEN** only that row changes through the existing reducer/autosave path.

#### Scenario: Apply current target type to selected rows
- **GIVEN** the explicit source has target_type int16 and eligible targets are selected
- **WHEN** the operator reviews and confirms the target-type-only action
- **THEN** only those targets receive int16 and all other mapping fields remain unchanged.

#### Scenario: No selected row means no quick target type action
- **WHEN** no source row is selected
- **THEN** no bulk type action is active.


### Requirement: Automatic mapping initialization

The system SHALL automatically initialize a mapping for every enabled point that does not yet have one. Initialization MUST be triggered by Step 3 mounting or by the set of enabled point IDs changing. For each new mapping, defaults SHALL derive Tag Key from device identity, point name and protocol-aware address; display name SHALL use the source name/display and address; unit SHALL use the source unit or remain empty. Neutral mappings SHALL preserve the supported source type, including int64, uint64 and string. A non-neutral numeric scale or offset SHALL default to float64; bool and string SHALL retain their source types. Scale/offset SHALL retain the point's rule values and enabled SHALL default true. Demo semantic arrays and raw seeds SHALL NOT determine a new persisted mapping. Existing mappings MUST NOT be overwritten. Mappings whose point no longer exists in `state.points` MUST be removed.

#### Scenario: Initialize new mappings on first mount

- **GIVEN** `state.mappings === {}` and `state.points` contains 8 enabled points
- **WHEN** Step 3 mounts and the `initMappingsForPoints` action runs
- **THEN** `state.mappings` contains 8 entries keyed by each point's `id`
- **AND** each mapping's defaults follow its source identity/type/unit and rule transform, without cycling demo templates

#### Scenario: Preserve existing mappings

- **GIVEN** the operator has edited `state.mappings[p3].tag_key` to `custom.tag.x`
- **WHEN** Step 3 remounts (or any other trigger of `initMappingsForPoints`)
- **THEN** `state.mappings[p3].tag_key` remains `custom.tag.x`
- **AND** no other mapping fields for `p3` are reset

#### Scenario: Remove orphan mappings

- **GIVEN** `state.mappings` contains entry for `pt-X` but `state.points` no longer has any point with id `pt-X` (e.g. Step 2 deleted the source rule)
- **WHEN** `initMappingsForPoints` runs
- **THEN** `state.mappings.pt-X` is removed

### Requirement: Transform pipeline preview

The Step 3 mapping preview surfaces SHALL request backend preview results for each displayed mapping instead of computing the operator-visible transform steps entirely on the client. The request SHALL use the actual scoped point sample as raw_value with the current draft transform pipeline, preserving exact integer representations. Missing source values SHALL show waiting/unavailable rather than a seeded success; stale/disconnected values SHALL retain explicit degraded state. Each row preview SHALL update reactively when its mapping scale, offset, target type or actual sample changes. The UI SHALL render explicit loading, success and error states. Basic SHALL expose the final value; Advanced SHALL retain available backend transform intermediates and payload inspection.

#### Scenario: Selected row triggers backend preview request

- **GIVEN** `selectedIdx === 0`, the row 0 mapping has `scale=0.1`, `offset=0`, `target_type='float64'`, and an actual source raw_value of `243`
- **WHEN** the preview updates
- **THEN** the system sends `POST /api/v1/datalink/mappings/preview`
- **AND** the request includes that actual raw_value `243`
- **AND** the request includes a transform pipeline derived from the current draft mapping values

#### Scenario: Preview loading state

- **GIVEN** a selected mapping exists
- **WHEN** a new preview request is in flight
- **THEN** the preview panel shows a loading state instead of stale success copy

#### Scenario: Preview renders backend result

- **GIVEN** backend preview returns successful `step_results` and `final_value`
- **WHEN** the response is received
- **THEN** the row preview renders the available backend transform intermediates and final value according to Basic/Advanced presentation
- **AND** the final value shown to the operator comes from backend preview output

#### Scenario: Preview failure remains visible

- **GIVEN** backend preview returns an execution error
- **WHEN** the response is received
- **THEN** the preview panel shows an actionable error state
- **AND** the UI does not synthesize a local success preview

#### Scenario: Older preview response is ignored

- **GIVEN** the operator changes the selected mapping twice in quick succession
- **WHEN** the earlier preview response arrives after the later one
- **THEN** the UI keeps the latest preview result
- **AND** the earlier response does not overwrite the newer preview state
