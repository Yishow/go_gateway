# datalink-workbench-v2-step4-database Specification Delta

## ADDED Requirements

### Requirement: Cross-protocol Modbus Share and database target binding

The system SHALL allow points derived from any supported protocol (Modbus, FATEK, MC 3E) to be mapped to database columns and optional Modbus Share slave registers. When Modbus Share is enabled for a non-Modbus source rule, the system MUST allocate continuous 40001+ Modbus holding registers for the output while maintaining source point linkage.

#### Scenario: Non-Modbus source with Modbus Share output
- **WHEN** an MC 3E rule with 4 points (`D0 ~ D3`) has Modbus Share enabled starting at `40001`
- **THEN** the Step 4 summary displays Modbus Slave mapping as `40001 ~ 40004`
- **AND** database column mapping operates independently without address collision

#### Scenario: Rule-level Modbus Share survives hydration
- **WHEN** a persisted source rule has Modbus Share enabled with a manual start register and stride
- **THEN** the rule create/update/list/get contract preserves `share_enabled`, `share_start_register`, and `share_stride`
- **AND** reloading the Studio V2 workspace hydrates the same Share settings instead of resetting them to disabled

#### Scenario: Activation projects Share mappings fail-closed
- **WHEN** the user activates a workspace containing persisted tags for enabled Modbus Share rules
- **THEN** the system synchronizes those tags to runtime Modbus Share mappings before workspace activation
- **AND** human holding register `40001` is sent to the runtime mapping API as zero-based register `0`
- **AND** stale mappings are deleted only when durable source-rule/tag relationships prove that their tag belongs to the current workspace
- **AND** a missing persisted tag or any synchronization failure prevents workspace activation

#### Scenario: Global Share disable overrides rule-level enablement
- **GIVEN** one or more source rules have `share_enabled=true`
- **WHEN** `settings.modbus_share.enabled` is false
- **THEN** Step 4 does not display active Share output mappings
- **AND** workspace activation does not create, update, or delete runtime Share mappings for those rules

#### Scenario: Unproven runtime mapping ownership is preserved
- **GIVEN** the runtime contains a mapping whose `tag_id` is not present in the current workspace's durable persisted-tag ownership set
- **WHEN** workspace Share mappings are synchronized
- **THEN** the system does not delete that runtime mapping
- **AND** a non-empty frontend mapping `tag_id` alone is not treated as ownership proof

#### Scenario: Automatic target column matching for protocol tags
- **WHEN** running autoAssignTargets on points with protocol tag keys (e.g. `line01.sensor.d0`)
- **THEN** the algorithm matches columns matching the tag suffix (e.g. `sensor_d0` or `d0`)
