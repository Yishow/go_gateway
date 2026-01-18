## ADDED Requirements
### Requirement: Mapping definition
The system SHALL define mappings from source points to target tags with an ordered transform pipeline.

#### Scenario: Create mapping
- WHEN an operator maps a point to a tag
- THEN the mapping is stored with its transform steps

### Requirement: Transform steps
The system SHALL support transform steps including decode, type cast, scaling, lookup table, conditional rules, and formula expressions.

#### Scenario: Scale and offset
- WHEN a raw value is processed with scale and offset
- THEN the final value is computed and stored

### Requirement: Formula function set
The system SHALL support formula expressions with operators (+, -, *, /, %), comparisons (>, <, >=, <=, ==, !=), boolean operators (and, or, not), conditional function (if), math functions (abs, min, max, clamp, round, floor, ceil, sqrt, pow), bitwise functions (bitand, bitor, bitxor, shiftl, shiftr), and conversion functions (int, float, bool, string, coalesce).

#### Scenario: Clamp and round
- WHEN a formula applies clamp and round to a value
- THEN the system returns the expected numeric result

### Requirement: Formula expression syntax
The system SHALL accept formula expressions using infix operators, snake_case function names, and boolean operators and/or/not.

#### Scenario: Evaluate infix expression
- WHEN a formula uses infix operators with snake_case functions
- THEN the system evaluates the expression without syntax errors

### Requirement: Validation of transforms
The system SHALL validate mapping steps for type compatibility and required parameters.

#### Scenario: Invalid cast
- WHEN a mapping attempts to cast text to a numeric type without a parser
- THEN the system rejects the mapping with a validation error

### Requirement: Raw and final value capture
The system SHALL preserve raw values alongside final transformed values.

#### Scenario: Preserve raw
- WHEN a value is transformed
- THEN the raw and final values are available for preview and storage

### Requirement: Mapping preview
The system SHALL provide a preview result for a mapping using live or sample values.

#### Scenario: Preview with sample
- WHEN an operator runs preview
- THEN the system shows raw, step-by-step, and final values
