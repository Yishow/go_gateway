# datalink-workbench-v2-step2-rule Specification Delta

## ADDED Requirements

### Requirement: Protocol-radix source address validation and expansion

The Step 2 source-rule planner SHALL validate and expand an address using the owning device protocol at both frontend and backend boundaries. MC 3E `X`, `Y`, and `B` areas SHALL use hexadecimal suffixes, MC 3E `D`, `W`, and `M` areas SHALL use decimal suffixes, and invalid input SHALL produce an explicit rule-scoped error. Invalid input SHALL NOT produce derived points and SHALL NOT return the original raw address as a successful result.

#### Scenario: MC 3E hexadecimal offset crosses a digit boundary

- **GIVEN** an enabled MC 3E rule starts at `X0`
- **WHEN** the planner or source-rule service offsets that address by `16`
- **THEN** the resulting address is `X10`
- **AND** the same rule treats `Y0` and `B0` with hexadecimal suffixes

##### Example: Protocol radix boundary

| Protocol | Start | Offset | Expected |
| --- | --- | ---: | --- |
| `mc_3e` | `X0` | 16 | `X10` |
| `mc_3e` | `D0` | 16 | `D16` |
| `fatek_fbs` | `R0` | 16 | `R16` |
| `modbus_tcp` | `40001` | 2 | `40003` |

#### Scenario: Invalid protocol address fails closed

- **GIVEN** an enabled rule contains an address that does not match its owning device protocol
- **WHEN** frontend validation, point expansion, or backend source-rule validation runs
- **THEN** the operation returns a validation error that names the rule and address field
- **AND** the point collection is empty or the request is rejected
- **AND** the invalid raw address is not returned as a successful derived point

### Requirement: Step 2 continuation evaluates every enabled rule

The Step 2 continuation gate SHALL evaluate every enabled rule in the workspace, resolve each rule to its owning device protocol, and require a valid address and derived point set for each rule. One invalid enabled rule SHALL block continuation even when another enabled rule has valid points. The gate SHALL expose the failing rule identifier and a concrete address or ownership reason in the readiness result.

#### Scenario: A valid sibling does not mask an invalid rule

- **GIVEN** a workspace contains one enabled rule with valid MC 3E `D0` points
- **AND** a second enabled rule contains an invalid address for its owning protocol
- **WHEN** the operator activates the Step 2 continue action
- **THEN** Step 2 remains active
- **AND** the readiness result identifies the second rule by stable rule identifier
- **AND** the valid sibling's points do not permit navigation to Step 3

#### Scenario: Disabled invalid rule does not create an enabled-rule blocker

- **GIVEN** a workspace contains a disabled rule with an invalid address and an enabled rule with a valid address
- **WHEN** the operator activates the Step 2 continue action
- **THEN** the disabled rule is excluded from the enabled-rule continuation check
- **AND** the enabled valid rule is evaluated normally
- **AND** the disabled rule remains visible with its own validation state for later editing
