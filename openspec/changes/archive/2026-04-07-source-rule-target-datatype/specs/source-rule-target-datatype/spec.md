# source-rule-target-datatype Specification (Change Delta)

## Purpose

定義「來源規則」在規劃階段可宣告 **協議讀取資料型態** 與 **目標 Tag／輸出資料型態** 的產品契約，以及與 Mapping 轉換管線的銜接方式。

## ADDED Requirements

### Requirement: Operator MAY declare a target data type distinct from protocol read type

The system SHALL allow an optional **target data type** on a source rule that MAY differ from the rule’s protocol read (`point`) data type. When omitted or null, the system SHALL treat the target as identical to the protocol read type.

#### Scenario: Default matches protocol type

- **WHEN** an operator creates a rule without setting a separate target data type
- **THEN** the system SHALL persist the rule such that target equals protocol read type
- **AND** downstream Tag and Mapping behavior SHALL match today’s same-type, empty-pipeline behavior

#### Scenario: Explicit target differs from protocol read type

- **WHEN** an operator sets target data type to a value different from the protocol read data type
- **THEN** the system SHALL accept the rule only if the pair is supported by validated cast semantics
- **AND** SHALL persist both values on the rule record

### Requirement: Unsupported type pairs are rejected at rule save time

The system SHALL reject rule persistence when the chosen protocol-read and target type pair cannot be satisfied by a validated `cast` transform according to pipeline validation rules.

#### Scenario: Rejection surfaces actionable error

- **WHEN** an operator submits a rule with an unsupported type pair
- **THEN** the system SHALL reject the save with a validation error
- **AND** the error SHALL identify the incompatible pair without exposing internal secrets

### Requirement: Transformed values use the mapping pipeline at runtime

The system SHALL NOT introduce a parallel transformation path for rule-target typing. Final values written to tags and outputs SHALL continue to flow through the persisted mapping `transform_pipeline` executed by the runtime ingestor.

#### Scenario: Cast appears as mapping pipeline step

- **WHEN** a rule declares a target type that differs from the point read type
- **THEN** the derived mapping SHALL include a `cast` step consistent with the declared target type
- **AND** runtime execution SHALL use the same `ExecutePipeline` behavior as manually authored mappings

### Requirement: Supported example — uint16 read with float64 tag semantics

The system SHALL support the common case where the protocol read type is **`uint16`** (single-register raw integer) and the operator declares a target type of **`float64`** for the Tag, provided the `cast` is valid per pipeline rules.

#### Scenario: Cast from uint16 to float64

- **WHEN** a rule sets protocol read type `uint16` and target type `float64`
- **THEN** the system SHALL persist both intents and derive a mapping pipeline that includes a validated `cast` to `float64`
- **AND** the ingested raw value SHALL be interpreted first as the point read type before casting

### Requirement: Operator MAY declare linear scaling for engineering units

The system SHALL allow optional **scale** parameters on a source rule (for example multiplier and offset aligned with `TransformParamsScale`) when the operator needs **unit conversion** beyond type casting. When set, the derived mapping SHALL include a validated **`scale`** transform step in addition to any `cast` step, in an order defined by the implementation and validated by the pipeline.

#### Scenario: Scale applies after type alignment

- **WHEN** an operator sets non-default scale parameters on a rule that also declares a target type
- **THEN** the persisted mapping pipeline SHALL validate successfully
- **AND** runtime output SHALL reflect both cast and scaling as separate traceable steps
