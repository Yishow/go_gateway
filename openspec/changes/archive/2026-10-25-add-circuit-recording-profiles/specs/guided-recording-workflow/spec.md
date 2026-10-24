## MODIFIED Requirements

### Requirement: Four-step intent-led setup

Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps labelled 連接設備 → 選要讀的資料 → 命名與換算 → 儲存並開始記錄. Basic/Advanced presentation changes SHALL preserve configuration values and protocol/type capabilities. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route. The basic path SHALL derive and persist routine mappings. New basic recording SHALL offer circuit-oriented instantaneous and cumulative profiles backed by managed long groups, while existing per-device wide and advanced configurations SHALL retain their identity and behavior until explicit reviewed cutover. Raw unassigned Tags SHALL remain selectable without inventing circuit or measurement identity.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator maps temperature, pressure, flow and counter source values as basic typed Tags
- **THEN** the operator can select those Tags for managed instantaneous recording without hand-written SQL or mandatory derived-value policies; choosing cumulative usage then requests the required confirmed roles.

#### Scenario: Reload and edited connection
- **WHEN** an operator reloads setup or changes a saved connection identity
- **THEN** the same persisted validity and revision checks determine progress, and the old probe does not keep the changed device ready.

#### Scenario: Unavailable protocol
- **WHEN** MQTT or another protocol lacks the complete V2 setup/parser/collector capability
- **THEN** the UI clearly marks it unavailable with a reason rather than letting the user enter an unsupported range flow.

#### Scenario: Same address on different devices
- **WHEN** two devices both expose address 40001
- **THEN** the address conflict checker keeps them separate while still detecting true overlapping ranges within the same device/address space.

#### Scenario: Missing advanced measurement
- **WHEN** a saved Point-to-Tag mapping has no persisted MeasurementDefinition
- **THEN** basic instantaneous setup proceeds without fake measurement IDs; only explicitly selected cumulative semantics remain blocked.

#### Scenario: Offline draft navigation
- **WHEN** a device is offline or its saved configuration has no successful probe
- **THEN** the operator can save a draft and use Next/Back to configure later steps, including after reload, while the device remains visibly unverified and its dependent activation is blocked with a repair action.

#### Scenario: Routine mapping does not require duplicate entry
- **WHEN** the operator confirms known points and supported types in the basic path
- **THEN** the existing mapping services persist real scoped Tag identities once, display actual reads and allow correction without a second manual mapping exercise.

#### Scenario: Existing advanced configuration is reopened
- **WHEN** an existing custom-table, multi-entity or cross-device group is loaded
- **THEN** its persisted identity and layout remain intact with accessible advanced controls; the basic default neither splits nor overwrites it.

#### Scenario: New basic setup is resumed
- **WHEN** a new circuit-profile managed draft is reloaded or its save response is retried
- **THEN** the persisted workspace, circuit, profile-kind and canonical-role key reuses the same group and mappings without duplicate writers. A changed destination requires an explicit CAS-checked draft and new preparation/apply evidence; accepted historical data keeps its original destination. Existing per-device wide keys remain valid and advanced groups remain outside automatic conversion.

### Requirement: Complete basic preparation with one effective configuration

Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, storage and profile interval. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics. New long profiles SHALL provide plain or supported Timescale storage selection, source/profile review and schema confirmation in Basic. Existing wide groups SHALL retain their preparation path.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides source/profile selection, schema preview and explicit confirmation and proceeds to timed typed long records per selected source without Advanced.

#### Scenario: Inspect effective settings
- **WHEN** an operator prepares a canonical group
- **THEN** no competing connector table, INSERT/UPSERT or recording interval control is presented as effective, and technical identifiers remain inspectable through diagnostics.

#### Scenario: Destination before recording choices
- **WHEN** no database destination is saved and Share is not enabled
- **THEN** Basic first presents destination setup before recording members, frequency and start; a saved destination is shown as a summary with an edit entry preserving advanced connection and password semantics.

#### Scenario: Prepare action and focused repair
- **WHEN** Basic preparation is requested or a prepare_schema error needs repair
- **THEN** one preparation action loads the read-only preview in the Basic schema card, with explicit DDL confirmation still required and no forced switch into all Advanced settings.

### Requirement: Basic missing-value policy reveals required gating

Existing wide-group Basic recording SHALL state that required missing or invalid members still block a partial row, show the effective required state, and provide a narrow entry to existing per-member required controls without silently changing any member. New long-profile Basic recording SHALL explain per-source NULL and quality records and MUST NOT present wide-row skip/partial controls as effective long-profile settings.

#### Scenario: Required gap remains blocked
- **WHEN** an existing wide basic group selects partial and a required member is missing or invalid
- **THEN** the UI SHALL explain that the row is skipped and SHALL expose the existing member settings without changing required automatically.

#### Scenario: Optional gap and recovery
- **WHEN** an explicitly optional wide-group member is missing in a partial-capable destination and subsequently recovers
- **THEN** the system SHALL preserve NULL and its reason for the partial row and write its actual value after recovery; skip_row and required protection SHALL remain unchanged.

#### Scenario: Edit existing exact managed members
- **WHEN** an operator opens required settings for saved SQLite managed wide uint16 or uint32 members in canonical generated TEXT columns
- **THEN** the editor SHALL retain their exact storage compatibility and require fresh schema confirmation after a group revision change without weakening custom-table checks.

#### Scenario: Confirm unchanged schema after member policy edit
- **WHEN** a saved managed group changes required or completeness policy without changing its physical columns
- **THEN** a compatible preview SHALL return an empty statements array, retain its revision-bound token and operation identity, and require explicit confirmation before applying the new group revision.

#### Scenario: Long profile has a missing source
- **WHEN** one source in a long profile has no usable sample while another source is valid
- **THEN** Basic SHALL show the missing source's NULL/reason and the other source's independent recording state without offering whole-wide-row completeness controls.

## ADDED Requirements

### Requirement: Simple circuit choices and scoped batch edits

Basic SHALL provide independently labelled instantaneous and cumulative controls, their interval units and the 15-second and 10-minute defaults. Operators SHALL be able to select several circuits and apply a shared profile or interval while retaining unselected circuits. Existing individual overrides SHALL be listed before the operator explicitly confirms their replacement. Cumulative selection SHALL explain that both raw readings and differences are stored, with baseline and uncertain differences displayed as unknown rather than zero. Controls and errors SHALL be keyboard accessible and localized in en and zh-TW.

#### Scenario: Batch edit with override
- **WHEN** circuits A and B are selected for a 30-second instantaneous interval, B already has an individual override, and C is unselected
- **THEN** Basic previews B's override replacement, saves A and B only after confirmation, preserves C, and retains these choices after reload.

#### Scenario: Cumulative first record
- **WHEN** the first cumulative raw reading has reached SQL but no comparable prior baseline exists
- **THEN** the interface distinguishes successful raw recording from an unknown first interval difference and does not display zero consumption or a write failure.

### Requirement: Explain advanced groups and workspace responsibilities

Basic SHALL identify Advanced as the editor for canonical source mapping, custom layouts and confirmed signed/rollover policies, and Workspace as the scope for readiness and activation. It MUST NOT describe them as alternative data-table formats or competing recording settings. A blocked recording SHALL identify the affected circuit, source or destination and offer its repair action. Identical blocker causes SHALL not be repeated as independent failures, and collecting or local acceptance SHALL not be shown as SQL commit evidence.

#### Scenario: Advanced and workspace review
- **WHEN** an operator opens Advanced or workspace confirmation after selecting circuit profiles
- **THEN** Advanced shows the same saved profile settings, workspace shows scoped readiness and activation, and neither silently changes intervals or storage.

#### Scenario: Database preparation failed
- **WHEN** a selected destination lacks schema confirmation or Timescale capability
- **THEN** the operator sees one actionable explanation for that cause with a destination repair entry rather than an unexplained disabled start button or raw driver error.
