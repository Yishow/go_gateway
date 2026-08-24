# studio-v2-device-validity Specification Delta

## ADDED Requirements

### Requirement: Source-rule device ownership is fail closed

A Studio V2 source rule SHALL be valid only when its `device_id` resolves to a current device owned by the active workspace. Source-rule ownership is single-workspace and indirect: the rule belongs to exactly one device through `device_id`, and the device is the ownership boundary; the persisted `source_rules` schema intentionally has no `workspace_id` or tombstone. A successful normal deletion of a workspace device SHALL remove its device-owned source rules through the existing foreign-key cascade, making them unavailable rather than leaving persisted orphan records for repair. Cleanup or repair of pre-existing orphan rows is explicitly outside this change. An unknown, deleted, empty, or otherwise unavailable device reference that is exposed by a handler/API response or retained in stale workspace state SHALL produce an ownership issue and SHALL NOT be rebound to the first workspace device, the first Modbus device, or any other fallback device. A newly constructed default rule SHALL select the current default device during creation, while hydration, editing, derivation, and continuation SHALL preserve and validate the persisted ownership identity. Workspace listing SHALL remain scoped to its ordered device IDs and SHALL NOT search globally for cascaded-away rules.

#### Scenario: Unknown device blocks the rule

- **GIVEN** a workspace-scoped handler/API response or stale in-memory workspace state contains a source rule with `device_id=deleted-device` and the current workspace device set does not contain that ID
- **WHEN** Studio V2 hydrates that returned state or evaluates Step 2 readiness
- **THEN** the rule is marked invalid with an ownership reason containing the rule ID and device ID
- **AND** no point is derived against another device
- **AND** the continue action remains blocked while the rule is enabled
- **AND** this defensive scenario does not require the List API to bypass `OrderedDeviceIDs`/`ListByDeviceIDs` or enumerate globally stored rules

#### Scenario: Normal device deletion removes its owned rules without reassignment

- **GIVEN** a source rule is owned by workspace device `device-A`
- **WHEN** the normal deletion of `device-A` succeeds
- **THEN** the existing foreign-key cascade removes the device-owned source rule and it is unavailable in subsequent workspace projections
- **AND** no `source_rules.workspace_id` or tombstone is created, and no global orphan lookup is performed
- **AND** if a stale/inconsistent reference is nevertheless returned by an API or retained in workspace state, the UI marks it deleted and fails closed rather than assigning it to `devices[0]` or the first Modbus device

<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/db.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/migrations/005_source_rules_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
tests:
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/datalink/db_sqlite_fk_test.go
  - frontend/tests/unit/workbench-v2/studioV2RuleAutosave.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule-readiness.test.ts
-->

#### Scenario: A valid owned device remains the only ownership source

- **GIVEN** a workspace contains `device-A` and `device-B` with different protocols
- **AND** a rule has `device_id=device-B`
- **WHEN** the rule is expanded and validated
- **THEN** the system uses `device-B`'s protocol and identity
- **AND** it does not use `device-A` merely because `device-A` is first in the collection
