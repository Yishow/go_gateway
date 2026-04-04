# mapping-pipeline Specification Delta

## ADDED Requirements

### Requirement: Source rule MAY seed a default transform pipeline

The system SHALL allow a mapping’s initial `transform_pipeline` to be populated automatically when created from a source rule that declares a target data type different from the point read type. The automatic content SHALL be limited to validated steps required to satisfy that intent (typically a single `cast` step) unless a broader policy is explicitly defined elsewhere.

#### Scenario: Default cast is validatable

- **WHEN** a mapping is created or updated from source-rule synchronization with differing point and target types
- **THEN** the automatically inserted steps SHALL pass the same transform validation as manually authored pipelines
- **AND** execution SHALL produce values compatible with the tag’s declared data type

### Requirement: Automatic pipeline updates SHALL NOT silently destroy unrelated manual steps

The system SHALL define and enforce a deterministic policy when a mapping already contains non-empty transform steps that were not derived from the rule (for example manual scale or formula). The policy MUST either preserve operator intent with explicit conflict resolution or refuse automatic updates until resolved.

#### Scenario: Conflict is visible to the operator

- **WHEN** automatic rule synchronization would replace or prepend transforms on a mapping that has non-default manual steps
- **THEN** the system SHALL NOT silently drop manual steps
- **AND** SHALL surface a conflict or required confirmation according to the chosen policy

### Requirement: Source rule MAY seed a default scale step

The system SHALL allow a mapping’s `transform_pipeline` to include an automatically generated **`scale`** step when the source rule declares non-default scale parameters. The step SHALL use the same `scale` transform type and validation rules as manually authored pipelines.

#### Scenario: Scale-only rule without cast

- **WHEN** a rule keeps protocol read type equal to target type but sets scale parameters
- **THEN** the derived mapping MAY contain only a `scale` step (or `scale` combined with other allowed steps)
- **AND** validation SHALL succeed
