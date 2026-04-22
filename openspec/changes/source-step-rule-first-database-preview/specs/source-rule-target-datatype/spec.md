## ADDED Requirements

### Requirement: Database context MAY influence Source-rule defaults without overriding operator intent
The system SHALL allow the active `database` output target to influence Source-step defaults for target data type, scale, and naming guidance when the operator is planning a source rule.

These defaults SHALL remain advisory. The system SHALL NOT silently replace an explicit operator choice only because the output target is `database`.

#### Scenario: Database target adjusts planning defaults
- **WHEN** an operator opens or creates a source rule while the active output target is `database`
- **THEN** the Source-step form MAY prefill or highlight target data type, scale, or naming defaults that better support downstream database review
- **AND** the operator can still change those values before saving the rule

#### Scenario: Explicit operator value wins over database default
- **WHEN** an operator manually changes target data type, scale, or naming-related rule input after a database-aware default is suggested
- **THEN** the system preserves the operator’s explicit value on save
- **AND** the UI does not silently revert that choice to a database-oriented default
