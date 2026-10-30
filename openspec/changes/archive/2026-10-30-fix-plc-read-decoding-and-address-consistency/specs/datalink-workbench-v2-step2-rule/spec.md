## MODIFIED Requirements

### Requirement: Protocol-aware point address derivation

The system SHALL derive planned point addresses by taking the rule's `start_address`, `count`, `data_type` (width stride), and owning device's `protocol`. The address calculation MUST preserve alphanumeric register area prefixes and advance offsets without stripping characters or falling back to Modbus defaults. Persisted MC point identities and addresses SHALL remain authoritative for an existing rule; derivation SHALL NOT silently replace them.

#### Scenario: Derive points for MC 3E device

- **WHEN** a new rule for an MC 3E device is defined with start address `D100`, count `4`, and data type `int16` (width `1`)
- **THEN** the system SHALL derive 4 points with addresses `D100`, `D101`, `D102`, and `D103`
- **AND** no Modbus `40001` fallback SHALL be applied.

##### Example: Multi-protocol point derivation

| Protocol | Start Address | Count | Data Type | Expected Addresses |
| --- | --- | --- | --- | --- |
| `modbus_tcp` | `40001` | `3` | `int16` | `40001, 40002, 40003` |
| `modbus_tcp` | `40001` | `3` | `float32` | `40001, 40003, 40005` |
| `mc_3e` | `D0` | `4` | `int16` | `D0, D1, D2, D3` |
| `mc_3e` | `D100` | `2` | `int32` | `D100, D102` |
| `mc_3e` | `D0` | `3` | `float32` | `D0, D2, D4` |
| `mc_3e` | `X0` | `3` | `bool` | `X0, X1, X2` |
| `fatek_fbs` | `R0` | `3` | `int16` | `R0, R1, R2` |

#### Scenario: Backend address offset matches protocol radix

- **WHEN** the source-rule service offsets MC 3E address `X0` by `16`
- **THEN** it SHALL return `X20`
- **AND** MC 3E X/Y SHALL use octal suffixes, W/B SHALL use hexadecimal suffixes, and D/M SHALL retain decimal suffixes.

#### Scenario: Backend rejects invalid protocol address

- **WHEN** source-rule validation or address offset receives an address invalid for the selected protocol
- **THEN** the operation SHALL return an explicit validation error
- **AND** it SHALL NOT return the original input as a successful derived address.

### Requirement: Protocol-radix source address validation and expansion

The Step 2 source-rule planner SHALL validate and expand an address using the owning device protocol at both frontend and backend boundaries. For new MC 3E plans, X/Y SHALL use strictly octal suffixes, W/B SHALL use hexadecimal suffixes, and D/M SHALL use decimal suffixes. The supported planner area set SHALL remain X, Y, W, B, D and M. Invalid input SHALL produce an explicit rule-scoped error, SHALL NOT produce derived points and SHALL NOT return the original raw address as a successful result. Existing direct MC reads SHALL retain the runtime parser's octal-first X/Y parsing with decimal fallback; this fallback SHALL NOT make X8, X9 or hexadecimal X/Y input valid for a new plan.

#### Scenario: MC 3E hexadecimal offset crosses a digit boundary

- **GIVEN** an enabled new MC 3E rule starts at `W0`
- **WHEN** the planner or source-rule service offsets that address by `10` and by `16`
- **THEN** the resulting addresses SHALL be `WA` and `W10` respectively
- **AND** the same hexadecimal formatting SHALL apply to B.

##### Example: Protocol radix boundary

| Protocol | Start | Offset | Expected |
| --- | --- | ---: | --- |
| `mc_3e` | `X0` | 16 | `X20` |
| `mc_3e` | `Y7` | 1 | `Y10` |
| `mc_3e` | `W0` | 10 | `WA` |
| `mc_3e` | `W0` | 16 | `W10` |
| `mc_3e` | `WA` | 1 | `WB` |
| `mc_3e` | `BFF` | 1 | `B100` |
| `mc_3e` | `D0` | 16 | `D16` |
| `mc_3e` | `M0` | 16 | `M16` |
| `fatek_fbs` | `R0` | 16 | `R16` |
| `modbus_tcp` | `40001` | 2 | `40003` |

#### Scenario: Invalid protocol address fails closed

- **GIVEN** a new MC rule contains X8, X9, YA, WG or Z0, or another rule contains input invalid for its owning device protocol
- **WHEN** frontend validation, point expansion or backend source-rule validation runs
- **THEN** the operation SHALL return a validation error that identifies the rule and address field
- **AND** the point collection SHALL be empty or the request SHALL be rejected
- **AND** the invalid raw address SHALL NOT be returned as a successful derived point.

#### Scenario: Hexadecimal suffixes survive point creation

- **WHEN** a new W0 int16 rule with count 17 is created, or WA or BFF is used as a new start address
- **THEN** backend planning and point storage SHALL accept the supported hexadecimal suffixes including all-letter suffixes
- **AND** the W0 rule SHALL persist W0 through W9, WA through WF and W10
- **AND** the runtime parser SHALL interpret those positions as consecutive physical offsets 0 through 16
- **AND** the change SHALL NOT add another MC area.

## ADDED Requirements

### Requirement: MC persisted source identities survive radix correction

Startup, reload and ordinary source-rule updates SHALL NOT silently migrate existing MC addresses to the corrected planner grammar. Before any mutation in an MC rule update, the service SHALL compare the address set produced from the stored rule's range and skipped addresses under the corrected grammar with its persisted links. A grammar error or set mismatch SHALL fail with a rule-scoped validation conflict identifiable as `mc_address_plan_conflict`. The update SHALL perform no point, link, tag, mapping, rule or revision mutation, including orphan cleanup. This change SHALL provide no implicit migration or bypass.

#### Scenario: Metadata update cannot recreate legacy points

- **GIVEN** an existing W0 count-17 rule has persisted addresses W0 through W16 with point, tag and mapping identities
- **WHEN** an update changes only the naming prefix or scale
- **THEN** it SHALL fail with the address-plan conflict before any write
- **AND** all rule/link/point/tag/mapping IDs, addresses, relationships and revision values SHALL equal the pre-request snapshot.

#### Scenario: A range patch does not authorize legacy migration

- **GIVEN** the stored MC rule's persisted set conflicts with its corrected stored-range plan, or its old start address is invalid under the corrected grammar
- **WHEN** an update also supplies a new start address, count or data type
- **THEN** the service SHALL still reject the request before mutation
- **AND** SHALL NOT infer migration authorization from that patch.

#### Scenario: Startup uses persisted points

- **WHEN** a valid existing MC rule and its complete persisted links are synchronized on startup or reload
- **THEN** synchronization SHALL retain point addresses and point/tag/mapping/link identities
- **AND** SHALL NOT regenerate point addresses from the corrected planner
- **AND** existing enabled-state and readiness synchronization SHALL remain effective.

#### Scenario: Compatible rules retain normal updates

- **GIVEN** a stored MC rule's corrected stored-range plan matches its persisted links
- **WHEN** a valid ordinary update is submitted
- **THEN** it SHALL retain the existing normal update behavior without applying the legacy-conflict guard to unrelated protocols.

### Requirement: Persisted MC points remain authoritative in Step 2

For persisted MC rules, Step 2 hydration, point selectors, mapping alignment and summaries SHALL use persisted point identities and addresses obtained through the existing scoped candidates GET. Each read SHALL carry workspace ID, expected workspace revision and rule revision, and its cache identity SHALL distinguish those scopes. Source metadata SHALL use the rule's source data type rather than the candidate target type. New drafts SHALL continue using the corrected planner. Checking, missing, stale, invalid or mismatched persisted evidence SHALL block automatic rule saves and mapping creation, rebinding or deletion for the affected rule. During hydration, the frontend SHALL NOT call POST candidates/recompute or any candidate apply action. The existing GET endpoint's internal database-output snapshot refresh SHALL remain unchanged and SHALL NOT rebuild the persisted MC tag-candidate addresses or identities.

#### Scenario: Reload does not reinterpret legacy addresses

- **GIVEN** a persisted W0 count-17 rule has scoped snapshot addresses W0 through W16 with existing point, tag and mapping IDs
- **WHEN** Step 2 reloads after the radix correction
- **THEN** its point list SHALL retain those persisted addresses and identities
- **AND** SHALL report the corrected-plan mismatch without presenting WA through WF as already stored
- **AND** loading or editing a non-address field SHALL NOT issue an automatic rule or mapping mutation for that conflicting rule.

#### Scenario: Unavailable or stale evidence blocks synthetic success

- **WHEN** a persisted MC rule's candidates read is pending, absent, incomplete, rejected with 403/409/503, or bound to a different workspace or rule revision
- **THEN** the rule SHALL be shown as checking or unavailable and SHALL remain blocked for dependent settings mutations
- **AND** SHALL NOT fall back to synthetic points, saved state or ready state
- **AND** unaffected devices and local drafts SHALL remain intact.

#### Scenario: Summaries do not assert a fabricated continuous range

- **WHEN** persisted MC points are noncontiguous or an old address cannot be parsed under the new-plan grammar
- **THEN** the point grid and summary SHALL retain the persisted address list or show an explicit blocked/unknown state
- **AND** occupied-range and conflict results SHALL NOT claim a newly computed continuous range or a clean result without sufficient evidence.

#### Scenario: Scope changes revoke writable evidence

- **WHEN** the workspace revision or rule revision changes after a compatible snapshot was loaded
- **THEN** Step 2 SHALL revoke that snapshot's mutation eligibility before reading the new scope
- **AND** a late response for the old scope SHALL NOT re-enable rule autosave or mapping reconciliation.
