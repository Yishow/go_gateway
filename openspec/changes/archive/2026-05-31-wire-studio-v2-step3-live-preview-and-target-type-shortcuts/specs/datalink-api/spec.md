## ADDED Requirements

### Requirement: Mapping preview request API

The system SHALL expose `POST /api/v1/datalink/mappings/preview` for draft transform preview requests used by `/studio/v2` Step 3 and similar clients.

#### Scenario: Preview a draft transform pipeline

- **WHEN** a client posts `raw_value` and `transform_pipeline` to `POST /api/v1/datalink/mappings/preview`
- **THEN** the API returns preview data containing `raw_value`, `final_value`, and ordered `step_results`
- **AND** the response can be rendered directly as transform preview output

#### Scenario: Preview failure returns actionable error

- **WHEN** preview execution fails for the submitted draft pipeline
- **THEN** the API response includes an actionable error message for the client
- **AND** the client can surface that failure without synthesizing a local success result
