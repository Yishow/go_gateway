## MODIFIED Requirements

### Requirement: Rules can generate tags and mappings automatically
The system SHALL generate rule-derived tag candidates and pending mapping intent from the primary source-rule flow before apply, rather than immediately creating active downstream tag-and-mapping state on rule save.

#### Scenario: Rule save produces reviewable tag candidates
- **WHEN** an operator creates or enables a source rule in the primary flow
- **THEN** the system SHALL generate the corresponding tag candidate and pending mapping intent
- **AND** SHALL defer active downstream persistence until the operator completes review/apply

#### Scenario: Generated tag identity uses canonical rule-owned identity
- **WHEN** the system generates a tag candidate from a rule
- **THEN** the tag candidate SHALL use the rule-owned identity contract for tag candidates
- **AND** SHALL remain unique under the global tag dictionary rules

## ADDED Requirements

### Requirement: Rule-generated tags preserve review decisions across compatible revisions
The system SHALL preserve rule-scoped rename, skip, and override decisions for rule-generated tags while the candidate identity remains materially the same.

For this workflow, a rule-generated tag candidate is materially the same only when all of the following remain unchanged:

- `source_rule_id`
- `derived_from_rule_address`
- derived target data type
- candidate kind (`tag`)

#### Scenario: Rename survives a compatible rule revision
- **WHEN** an operator renames a rule-generated tag and a later rule revision keeps the same rule-owned tag identity
- **THEN** the rename remains as the effective tag decision
- **AND** the system recomputes downstream output candidates from that effective tag state

#### Scenario: Identity change invalidates prior review decisions
- **WHEN** a later rule revision changes a rule-generated tag candidate so it is no longer materially the same
- **THEN** the system marks prior review decisions as stale for that candidate
- **AND** requires the operator to review the new candidate before apply
