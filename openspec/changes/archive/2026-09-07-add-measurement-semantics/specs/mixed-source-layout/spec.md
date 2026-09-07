## ADDED Requirements

### Requirement: Explicit mixed source layout
The system SHALL support homogeneous legacy rules and versioned mixed rules containing independently typed items with explicit spans, word order and effective transforms.

#### Scenario: Eight registers with one double-word item
- **WHEN** the configured MC example uses D0, D1, D2, D3-D4, D5, D6 and D7
- **THEN** seven primary decoded items are produced and D4 is not allocated as another independent point.

#### Scenario: Overlapping items
- **WHEN** a uint32 at D3-D4 overlaps another independent item at D4
- **THEN** saving or activating that layout is rejected with the conflicting items and span.

### Requirement: Safe decoding and single effective transform
The system MUST validate type width, order, signedness and conversion before activation and apply the effective engineering transform exactly once.

#### Scenario: Temperature scaling
- **WHEN** raw int16 -123 uses multiplier 0.1 and offset 0
- **THEN** the engineering value is -12.3 rather than -1.23 and its transform provenance is retained.

#### Scenario: Unsupported format
- **WHEN** a requested string length or word order is not supported by the selected adapter
- **THEN** the layout is blocked before collection rather than partially decoded.

### Requirement: Transport batching preserves item provenance
The system SHALL reuse compatible read optimization without merging item identity, quality or pretending cross-request atomicity.

#### Scenario: One optimized read serves several items
- **WHEN** compatible mixed items share a successful protocol read
- **THEN** the raw response is decoded by each item's declared layout and each output preserves its own identity.

### Requirement: Versioned templates with per-device review
The system SHALL provide three-phase meter and mixed-sensor templates that require per-device confirmation of roles, spans, units and scaling before apply.

#### Scenario: Reuse across different meter models
- **WHEN** an operator applies a meter template to a device with unconfirmed register positions
- **THEN** the system previews suggested bindings as unconfirmed and does not start collection or usage calculation.

### Requirement: Legacy lifecycle compatibility
The system MUST preserve existing rule ownership, probe gates, candidate revisions and homogeneous decoding when the new mode is not selected.

#### Scenario: Re-enable an existing rule
- **WHEN** an existing homogeneous rule is disabled and re-enabled after migration
- **THEN** it retains the same points, mappings, scaling and output bindings without requiring reconfiguration.
