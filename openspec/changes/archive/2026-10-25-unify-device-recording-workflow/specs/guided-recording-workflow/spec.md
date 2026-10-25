## MODIFIED Requirements

### Requirement: Four-step intent-led setup
Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps labelled 連接設備 → 選要讀的資料 → 命名與換算 → 儲存並開始記錄. Basic/Advanced presentation changes SHALL preserve configuration values and protocol/type capabilities. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route. The basic path SHALL derive and persist routine mappings, show only necessary operator choices, and offer whole-device managed wide snapshots or selected-source managed long recording in the selected device scope without converting existing configurations. Table layout SHALL be independent of instantaneous or cumulative semantics. Existing canonical same-kind ownership SHALL remain unique and existing latest snapshots SHALL not become means through navigation.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator maps temperature, pressure, flow and counter source values as basic typed Tags
- **THEN** the operator can select those Tags for a managed write group without hand-written SQL or mandatory derived-value policies; choosing derived usage then requests the required confirmed roles.

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
- **THEN** basic raw group setup proceeds without fake measurement IDs; only explicitly selected advanced semantics remain blocked.

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
- **WHEN** a new per-device managed draft is reloaded or its save response is retried
- **THEN** the persisted (workspace_id, device_id, canonical managed role) key reuses the same group and mappings without duplicate writers or selecting another group. A changed destination requires an explicit CAS-checked draft and new preparation/apply evidence; accepted historical data keeps its original destination. Existing advanced groups remain outside automatic conversion.

### Requirement: Complete basic preparation with one effective configuration
Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, layout, storage, interval and reduction. New numeric gauge recording SHALL default to the explicitly labelled arithmetic mean; existing latest and counter roles SHALL retain their persisted interpretation. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides schema preview and explicit confirmation and proceeds through the selected wide or long recording flow without Advanced and without an additional workspace activation action.

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
Wide and custom recording SHALL state that required missing or invalid members still block a partial row, show the effective required state, and provide a narrow entry to existing per-member required controls without silently changing any existing member. New mean wide groups SHALL default to optional partial recording only with verified nullable/provenance support. Long recording SHALL show independent per-source missing/coverage results and SHALL NOT offer a whole-row skip policy. Counter baseline or unknown delta SHALL not be described as a zero value or a generic save failure.

#### Scenario: Required gap remains blocked
- **WHEN** a basic group selects partial and a required member is missing or invalid
- **THEN** the UI SHALL explain that the row is skipped and SHALL expose the existing member settings without changing required automatically.

#### Scenario: Optional gap and recovery
- **WHEN** an explicitly optional member is missing in a partial-capable destination and subsequently recovers
- **THEN** the system SHALL preserve NULL and its reason for the partial row and write its actual value after recovery; skip_row and required protection SHALL remain unchanged.

#### Scenario: Edit existing exact managed members
- **WHEN** an operator opens required settings for saved SQLite managed uint16 or uint32 members in canonical generated TEXT columns
- **THEN** the editor SHALL retain their exact storage compatibility and require fresh schema confirmation after a group revision change without weakening custom-table checks.

#### Scenario: Confirm unchanged schema after member policy edit
- **WHEN** a saved managed group changes required or completeness policy without changing its physical columns
- **THEN** a compatible preview SHALL return an empty statements array, retain its revision-bound token and operation identity, and require explicit confirmation before applying the new group revision.

## ADDED Requirements

### Requirement: Device tiles change only the viewed scope

Step 4 SHALL provide named device tiles with one visibly selected persisted device ID, keyboard focus and native activation. It SHALL NOT use a device dropdown. Tile selection SHALL filter the current device's recording configuration, sources, preparation and status without creating, applying, disabling, activating or deleting resources. Drafts SHALL remain scoped by workspace/device/purpose/circuit; late replies and unmounted requests SHALL not update another scope. A legacy wide group on device A SHALL not force device B onto A's legacy workflow. Cross-device custom groups SHALL retain their existing advanced identity and SHALL not be automatically split.

#### Scenario: Same address and late reply
- **WHEN** A and B share the same address/display name, the operator edits A, switches to B and A's delayed save/read reply arrives
- **THEN** B retains its own draft and status, A retains its pending work under its own IDs, and switching back restores A without creating a writer or changing acquisition.

#### Scenario: Keyboard and reload
- **WHEN** the operator selects a device tile using keyboard and reloads or navigates Back/Next
- **THEN** the selected valid device and its scoped drafts are restored without a dropdown, false saved state or automatic Apply.

### Requirement: Purpose-led layout with one advanced authority

The system SHALL label the two new managed paths as whole-device snapshot (wide) and selected measurements (long), explain their row granularity, and allow instantaneous and cumulative roles in either layout. Advanced SHALL open the same saved group with the same profile/reduction/interval and explicit saved-versus-applied state. It SHALL not expose a competing interval or start authority. Legacy wide/custom settings SHALL remain editable under their own label without silent conversion; layout changes SHALL use reviewed version/ownership transitions instead of changing IDs or deleting groups.

#### Scenario: Full meter and mixed selected sources
- **WHEN** one device selects full voltage/current/power/counter snapshots and another selects temperature plus water/gas/electric counters
- **THEN** the first can use managed wide and the second managed long, each preserves explicitly selected roles and periods, and neither layout implies that counters are averaged.

#### Scenario: Advanced round trip
- **WHEN** an operator opens a profile from Basic in Advanced and returns without saving
- **THEN** group ID, reduction, interval, members and destination remain identical and no mutation or additional start occurs.

### Requirement: Compact source review preserves actionable exceptions

The source review SHALL default to a collapsed item count with an explicit expand action and an internally scrollable expanded list. Names SHALL use saved display names then source addresses, with technical keys restricted to diagnostics. Missing/foreign-source counts and repair actions SHALL remain visible while collapsed. Acquisition, recording intervals and later report windows SHALL be explained separately without adding a report engine or new acquisition cadence controls.

#### Scenario: Large source list and an unavailable member
- **WHEN** a saved group contains 32 members including one unavailable source
- **THEN** the default view shows the count and actionable unavailable-source summary without 32 permanently expanded rows, and expansion reveals the scoped source details without changing configuration.

#### Scenario: Recording versus report interval
- **WHEN** the selected instantaneous interval is 15 seconds and the cumulative interval is 10 minutes
- **THEN** the summary explains that accepted acquisition feeds those independent recording windows and does not claim to restore 15-second history from only 10-minute records or create reports.
