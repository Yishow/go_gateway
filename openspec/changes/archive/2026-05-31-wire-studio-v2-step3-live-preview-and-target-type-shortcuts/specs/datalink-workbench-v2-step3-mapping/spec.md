## MODIFIED Requirements

### Requirement: Transform pipeline preview

The right-side preview panel SHALL request backend preview results for the currently selected mapping instead of computing the operator-visible transform steps entirely on the client. The request SHALL use the deterministic raw seed array `[243, 251, 1024, 985, 67, 542, 18, 1450]` indexed by selected row position modulo 8 as `raw_value`, together with the current draft transform pipeline. The panel MUST update reactively whenever the selected row's mapping `scale`, `offset`, or `target_type` changes. The panel MUST render explicit loading, success, and error states.

#### Scenario: Selected row triggers backend preview request

- **GIVEN** `selectedIdx === 0`, the row 0 mapping has `scale=0.1`, `offset=0`, `target_type='float64'`, and `tag_key='line01.temp.inlet'`
- **WHEN** the preview updates
- **THEN** the system sends `POST /api/v1/datalink/mappings/preview`
- **AND** the request includes raw seed `243`
- **AND** the request includes a transform pipeline derived from the current draft mapping values

#### Scenario: Preview loading state

- **GIVEN** a selected mapping exists
- **WHEN** a new preview request is in flight
- **THEN** the preview panel shows a loading state instead of stale success copy

#### Scenario: Preview renders backend result

- **GIVEN** backend preview returns successful `step_results` and `final_value`
- **WHEN** the response is received
- **THEN** the preview panel renders ordered decode / scale / cast / final steps from that backend result
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

### Requirement: API payload preview

The preview panel SHALL render a `POST /mappings` JSON payload preview below the transform preview. The payload MUST include `point_id`, `tag_id` (`tag.${tag_key}`), `transform_pipeline` (three steps: decode, scale with params `{scale, offset}`, cast with param `{to: target_type}`), and `enabled`. The payload preview MUST update on every change to the selected mapping. The payload preview MUST remain display-only and MUST NOT initiate its own network request.

#### Scenario: Payload preview remains display-only

- **WHEN** any value in the selected mapping changes
- **THEN** the payload preview updates to reflect the latest draft values
- **AND** no additional `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the payload preview itself

## ADDED Requirements

### Requirement: Target type quick actions

The Step 3 preview panel SHALL provide quick target type actions for the currently selected mapping and a one-click apply-all action for `target_type` only. The quick actions MUST include at least `bool`, `int16`, `float64`, and `string`.

#### Scenario: Quick action updates the current row

- **GIVEN** a mapping row is selected
- **WHEN** the operator clicks one of the target type quick actions
- **THEN** the selected row's `target_type` is updated to that value
- **AND** the change flows through the same Step 3 reducer/autosave path as the inline select

#### Scenario: Apply current target type to all rows

- **GIVEN** `selectedIdx === 0` and the selected row has `target_type='int16'`
- **WHEN** the operator clicks `套用到全部列`
- **THEN** every mapping in `state.mappings` has `target_type === 'int16'`
- **AND** other mapping fields remain unchanged

#### Scenario: No selected row means no quick target type action

- **GIVEN** `selectedIdx === null`
- **WHEN** the operator views the preview area
- **THEN** no target type quick action is offered as an active control
