## Purpose

Provide a single progressive recording editor that reuses existing Studio device, point, mapping, preparation and scoped-start contracts without removing advanced capabilities.

## ADDED Requirements

### Requirement: Progressive editing preserves source capabilities
The system SHALL guide source selection, recording intent and destination on one task screen. Existing devices and points SHALL be reusable. Inline device and point fields SHALL remain in a local pending buffer until an explicit save dispatches the existing autosave action and queue. Cancellation SHALL discard only unsent fields. Inline configuration SHALL reuse canonical service identities and preserve supported address spaces, data types, word order, scale and units. Reading frequency SHALL remain distinct from recording frequency.

#### Scenario: A recording needs additional points
- **WHEN** an operator opens point settings while a recording name, selected sources and destination are already entered
- **THEN** closing or saving the settings returns to the same recording draft and focus without losing those values, and newly selected points use confirmed canonical identities.

#### Scenario: Item count differs from register span
- **WHEN** a supported word-addressed range contains 10 float32 items
- **THEN** the editor distinguishes 10 decoded items from the 20 occupied 16-bit registers and preserves the existing protocol-specific address validation.

#### Scenario: Duplicate names exist across devices
- **WHEN** points on different devices share a name or displayed address
- **THEN** selection preserves distinct scoped identities and shows device context rather than merging by name or address.

### Requirement: Basic and advanced recording semantics remain faithful
The system SHALL allow basic typed raw recording without mandatory measurement semantics. Explicit derived operations SHALL require their existing semantics. Layout SHALL remain independent from measurement meaning, and configurations not representable by the simple editor SHALL retain an accurate summary and scoped advanced access without automatic conversion.

#### Scenario: Mixed raw sensors are selected
- **WHEN** temperature, pressure and a raw counter are selected for a supported basic latest-value recording
- **THEN** the editor does not require usage-calculation policies merely because one value is cumulative.

#### Scenario: Unsupported simple-mode combination is opened
- **WHEN** an existing cross-device, custom-table or derived configuration cannot be represented by the simple editor
- **THEN** its identities and stored configuration remain unchanged and the operator can open the matching advanced editor with an explanation.

### Requirement: Draft navigation preserves or explicitly discards unsaved intent
The system SHALL distinguish persisted draft state from unsaved page state and effective configuration. Close, Escape, Back, Forward, scope changes and return-home actions SHALL use consistent retention or explicit discard behavior. Only server-recoverable content SHALL be described as surviving reload.

#### Scenario: A modal is cancelled
- **WHEN** an operator cancels device or point editing inside a recording draft
- **THEN** the original recording draft remains intact and focus returns to the invoking control.

#### Scenario: An unsaved edit is dismissed
- **WHEN** an operator leaves an edited recording through Escape, browser navigation or return home
- **THEN** the same dirty-state policy offers retention or explicit discard, and no unconfirmed edit is represented as the effective recording.

#### Scenario: Autosave fails and the page reloads
- **WHEN** a draft write fails before a reload
- **THEN** only the last confirmed persisted state is restored and the application does not claim the failed fields were saved.

### Requirement: One next action retains every safety boundary
The system SHALL reuse existing saving, validation, schema review and scoped-start contracts. A change to relevant inputs SHALL invalidate prior readiness or confirmation. Displaying a preview or saving configuration SHALL NOT start recording, and an unresolved command SHALL be recovered by status lookup rather than an automatic replacement command.

#### Scenario: Input changes after schema confirmation
- **WHEN** the selected source, layout or destination changes after confirmation
- **THEN** the old confirmation cannot authorize the new configuration and the operator must review the new impact before starting.

#### Scenario: Starting is acknowledged without a receipt
- **WHEN** the existing scoped-start operation is acknowledged
- **THEN** the editor returns to a waiting-first-write state until matching SQL evidence is received, without requiring workspace-wide activation or claiming SQL success.

#### Scenario: Repeated confirmation or late response
- **WHEN** the operator clicks the same action repeatedly or receives a stale response after changing scope
- **THEN** the existing in-flight, operation-identity and revision guards prevent duplicate commands and cross-scope updates.

### Requirement: Closing an editor does not roll back submitted source changes
The system SHALL distinguish unsent modal fields, submitted source changes and confirmed saved source changes. A successful source save SHALL remain saved when the recording form is discarded, and a pending or unknown save SHALL use existing settlement or recovery rather than being labelled cancelled.

#### Scenario: Source save succeeds before recording cancellation
- **WHEN** a new point range has been explicitly saved and the operator later discards the unsaved recording form
- **THEN** the interface states that the saved range remains, discards only unsent recording intent and does not issue an implicit compensating delete.

#### Scenario: A source save outcome is unknown
- **WHEN** the operator closes a modal after submitting a save whose result is unconfirmed
- **THEN** the unresolved save is retained for settlement or status recovery and is not described as rolled back or cancelled.

### Requirement: Protocol and parameter parity is explicit before implementation
The simple editor SHALL preserve the existing usable Modbus TCP, Modbus UDP, Modbus RTU, FATEK FBs and MC 3E source paths, while MQTT remains explicitly unsupported by the current V2 source planner. Each existing required or advanced parameter SHALL be directly editable, reachable through a scoped advanced path with lossless preservation, or explicitly unsupported. Existing backend validators, address grammar and identity guards SHALL remain authoritative; no driver or scheduler capability SHALL be added to fill a demo gap.

#### Scenario: Serial configuration is reloaded and changed
- **WHEN** an existing RTU or FATEK serial configuration is loaded and one visible field is edited
- **THEN** data bits, stop bits, parity, timeout, data format and other retained parameters round-trip without loss, FATEK does not silently become TCP, and server baud validation is enforced.

#### Scenario: MC and Modbus addresses are edited
- **WHEN** an operator edits MC D/M, X/Y or W/B addresses or supported Modbus coil/input/register areas
- **THEN** existing decimal/octal/hex and protocol validation applies rather than demo-only D or 40001 restrictions, aliases normalize through the existing adapter, and persisted MC compatibility remains fail-closed.

#### Scenario: Point details and sampling cadence are inspected
- **WHEN** an operator opens point details
- **THEN** the existing datatype, count, word order, scale, skipped addresses, naming, enabled and Share controls remain available with their guards, mapping unit/scale/target types stay in their proper layer, and sampling cadence is shown from existing settings or a scoped advanced handoff rather than a fabricated per-rule poll field.

#### Scenario: An unavailable protocol is selected in a demo
- **WHEN** MQTT or an unimplemented demo parameter is displayed
- **THEN** unsupported or simulation-only status is explicit, no real credential is requested, and no successful real recording capability is implied.

### Requirement: Recent connection defaults are session-scoped and explicitly clearable
The new-device flow SHALL offer the most recent successfully saved non-sensitive connection parameters for the same workspace, protocol and applicable subtype during the currently mounted page session. It SHALL visibly identify prefilled values and remind the operator to verify host or serial port before an explicit save. Cache updates SHALL derive from the successfully saved request snapshot or canonical server response, not the current editable form or a saved-status observer, and SHALL match the scope/session epoch, protocol/subtype and request generation. Clearing SHALL invalidate earlier request generations. Only a protocol-schema whitelist SHALL be cached; names, identities, credentials, unknown server keys, points, mappings, activation state and readiness, receipt, operation or ownership evidence MUST NOT be copied. MQTT SHALL remain unsupported. The cache MUST NOT use browser persistent storage or introduce a backend preference model or an invented fixed-duration TTL.

#### Scenario: A second device uses the same protocol
- **WHEN** a new-device save succeeds in the current workspace and another new-device form is opened for the same protocol and subtype
- **THEN** eligible connection parameters prefill with a visible reuse notice and endpoint reminder, while the device name and identity remain new and no probe or activation occurs automatically.

#### Scenario: A protocol or FATEK mode changes
- **WHEN** an operator selects another protocol or switches FATEK between TCP and serial
- **THEN** only the matching protocol/subtype cache may prefill, incompatible fields are cleared, and no secret or unknown field is carried across.

#### Scenario: Saving does not establish success
- **WHEN** a new-device form is cancelled, its save fails, reports runtime_apply_status=apply_failed or remains pending/unknown, or an older scope returns late
- **THEN** those values do not replace the last confirmed successful cache; loading or editing an existing device also does not populate this new-device cache.

#### Scenario: The operator clears reused parameters
- **WHEN** the operator confirms clearing the current new form's reused connection values
- **THEN** all cached subtypes for that protocol are removed, connection fields return to the selected protocol/mode defaults or empty values where no default exists, the reuse notice disappears, and neither the new name nor any existing device is changed or saved.

#### Scenario: Clearing would discard edited connection fields
- **WHEN** clearing would replace connection values edited after prefill
- **THEN** the existing unsaved-change guard asks before resetting, and cancelling leaves both the form and cache unchanged.

#### Scenario: The page session ends
- **WHEN** the page reloads or unmounts, or the workspace changes
- **THEN** all recent-connection cache entries are cleared; closing only the new-device modal retains prior successful entries but never caches the cancelled input.

#### Scenario: An earlier save resolves after editing or clearing
- **WHEN** a successful save response arrives after additional unsaved edits or after the operator cleared that protocol cache
- **THEN** unsaved current form values never enter the cache, and an invalidated generation cannot repopulate cleared entries.

#### Scenario: Continue adding after a successful save
- **WHEN** an operator chooses to add another device in the same page session without a preexisting explicit protocol or typed fields
- **THEN** the last confirmed successful protocol/subtype may be preselected with its eligible cached parameters, while protocol switching remains available and no save or connection is automatic.

#### Scenario: A late save arrives while another new form is open
- **WHEN** the operator has already chosen a protocol or typed values in a new form before another save response arrives
- **THEN** that response does not change the open form's protocol, subtype or entered values even if it is eligible to update the session cache for later forms.

### Requirement: Parameter help is accessible and source-grounded
Each displayed configuration parameter control SHALL have an adjacent circular information button with a field-specific accessible name. The button SHALL open lightweight help by click, touch or keyboard, not hover alone, and MUST NOT submit the form. Help SHALL provide purpose, units or allowed range, examples and when to change the value where relevant, grounded in existing schema, validators and protocol contracts rather than invented limits or safety assurances. Required markers, errors, disabled reasons, critical risks and consent SHALL remain inline. Help popovers SHALL be distinguished from form modals and MUST NOT inherit the 75-percent modal sizing rule.

#### Scenario: An operator opens and closes help
- **WHEN** an operator activates the information button with pointer, touch, Enter or Space
- **THEN** readable theme-consistent help opens within the visible viewport without changing fields or issuing save, probe or activation commands; Escape or outside activation closes it, keyboard dismissal returns focus to its trigger, and outside dismissal does not steal focus from the new target.

#### Scenario: A parameter has protocol-specific limits
- **WHEN** help describes a protocol parameter's unit, range, example or default
- **THEN** the description matches the existing canonical schema or validator and does not imply that an example is valid for every device or weaken validation.

#### Scenario: Coverage differs between demo and product
- **WHEN** parameter-help coverage is reviewed
- **THEN** each displayed parameter maps to its canonical field, help source and visible demo location, while formal fields absent from the demo are explicitly listed in parity or advanced handoff rather than represented as implemented; critical risks remain visible with help closed.

### Requirement: Configured destinations offer safe table reuse
The editor SHALL offer previously configured destinations from the existing workspace write-group data, identified by connector, database, schema and table rather than display name alone. This chooser MUST NOT imply discovery of all physical database tables. Compatible canonical shared long-table destinations MAY be selected according to the existing scope/layout contract. A wide table owned by another group MUST NOT be silently reused by a new group; the interface SHALL explain the ownership restriction and offer scoped editing of the existing recording. Selection SHALL populate intent only and MUST NOT copy another group's identity, schema proof, readiness or ownership evidence.

#### Scenario: A compatible configured long table is selected
- **WHEN** an operator selects a canonical long-table destination compatible with the current scope and layout
- **THEN** its destination identity is prefilled but fresh connector, schema, revision and ownership checks remain necessary before any explicit start or write.

#### Scenario: The same table label appears in different databases
- **WHEN** configured destinations share a table label but have different connector or database identities
- **THEN** they remain distinguishable and cannot be merged by label, while loading failures or stale revisions disable unproven selections with a reason.

#### Scenario: A new recording selects another group's wide table
- **WHEN** the existing ownership contract does not permit reuse
- **THEN** selection cannot bypass ownership, the reason and original recording's edit path are available, and no table or recording is changed by opening the chooser.

### Requirement: Compact forms and short confirmations preserve useful density
Common device connection parameters SHALL remain visible in a compact responsive form rather than being hidden as advanced details. Point configuration SHALL use compact comparable rows for address, datatype, span and effective decoding order, while uncommon details MAY be disclosed progressively. Short discard, unsaved-change, stop and remove confirmations SHALL use semantic short-confirmation sizing with auto or fit-content height and bounded body scrolling, not the regular 75-percent form-modal height.

#### Scenario: Unsaved edits need confirmation
- **WHEN** a dirty editor requests discard confirmation
- **THEN** a compact centered short confirmation identifies unsaved changes without stretching to the regular modal height, while retaining explicit choice, focus return and dirty-state safety.

#### Scenario: An operator configures common parameters
- **WHEN** device or point configuration is opened
- **THEN** common connection fields and compact point rows are immediately visible with their labels, help and validation rather than hidden solely to lower apparent information density.

### Requirement: Point order controls consume the canonical override contract
Per-point order controls SHALL depend on the completed preserve-per-point-word-order contract rather than treating existing point fields as sufficient runtime support. Each point SHALL show its explicit override, inherited or legacy-preserved intent, effective decoding order and source. Supported choices SHALL come from the protocol/datatype capability contract. Reset-to-inherit SHALL be explicit; changing a range default MUST NOT overwrite individual overrides. Legacy intent MUST NOT be guessed or silently relabeled as inherited. The canonical Point SHALL remain the sole intent authority, with rule links used only for ownership. The UI SHALL project API mode values legacy/inherit/override and configured/effective/source/applicable/revision readback; legacy-preserved is only a display label. New points SHALL inherit, while an explicit override equal to the current default SHALL remain explicit. Saving SHALL retain canonical scope and revision guards.

#### Scenario: Points on one PLC need different orders
- **WHEN** supported multi-register points on the same PLC are configured with distinct explicit ABCD and CDAB orders
- **THEN** each row shows its own effective order and saves through the canonical override API rather than applying one shared range value to all points.

#### Scenario: Range defaults or supported types change
- **WHEN** a range default changes or a point's datatype/protocol no longer supports an order
- **THEN** explicit overrides are preserved according to the prerequisite contract, unsupported choices are blocked with a reason, and resetting to inherited order requires an explicit operator choice.

#### Scenario: The prerequisite capability is incomplete
- **WHEN** the canonical override or decoder capability is not available or the point revision is stale
- **THEN** the UI does not claim the order was saved or applied; a demo labels the interaction as simulated and product mutation waits for the prerequisite contract.

#### Scenario: A single-word point retains its existing byte policy
- **WHEN** the capability marks a single-word point inapplicable for a new multi-register override
- **THEN** the UI prevents that new override without forcing ABCD, clearing existing metadata or discarding the driver's legacy byte-swap behavior; legacy MC3E and FATEK effective behavior remains resolved by the prerequisite contract rather than inferred from stored data_format.

### Requirement: Recording choices explain outcomes without inventing restrictions
The editor SHALL describe recording choices by operator intent, concise plain-language labels and illustrative values, rather than the phrase single calculation per rule. Latest SHALL mean the eligible latest value in the recording window; mean SHALL describe an equal-weight arithmetic mean of eligible observations for the applicable gauge member, not a cross-point or time-weighted mean. Counter usage SHALL describe change against the confirmed reliable baseline under the existing counter policy. Per-member effective roles and results SHALL remain faithful: a mean profile can average gauges while retaining latest counter_raw/state values. The UI MUST NOT claim a universal same-algorithm restriction, invent unsupported algorithms beyond the separately specified mixed member contract, silently change saved profiles or present samples as real SQL receipts.

#### Scenario: An operator compares recording choices
- **WHEN** the editor explains eligible readings 20, 22 and 24
- **THEN** it illustrates latest as 24 and arithmetic mean as 22 with understandable labels and notes that examples are simulated, while each selected member's actual role remains visible.

#### Scenario: A cumulative reading has no proven change
- **WHEN** a counter has no baseline or continuity cannot be established
- **THEN** the explanation states that usage is unavailable with its actual reason rather than zero; confirmed signed-counter or rollover policies are not incorrectly described as universally invalid negative changes.

### Requirement: Successful read-range saving closes only its matching modal
An explicit read-range save SHALL keep the modal pending until matching scope, request, modal session, submitted buffer generation, persistence and required identities are confirmed. A success for an older buffer generation SHALL settle only its submitted changes and MUST NOT clear or close newer dirty input. Matching success SHALL close that modal, clear its submitted-buffer dirty state and return to the original editor or selection context with updated sources and restored focus. Validation failure, failed persistence or unknown outcome SHALL retain the modal and entered values with actionable feedback. A stale response MUST NOT close another modal or scope, and saving MUST NOT imply activation or SQL success.

#### Scenario: Range persistence succeeds
- **WHEN** the current read-range save returns matching confirmed success
- **THEN** its modal closes without a false unsaved-change prompt, the originating recording draft is preserved and the updated source is available without automatically starting collection.

#### Scenario: Range persistence fails or returns late
- **WHEN** saving fails or remains unknown, or an earlier request returns after a new modal or scope has opened
- **THEN** current input is retained with the appropriate recovery state and the late response cannot close or overwrite the newer form.

#### Scenario: The range form changes while its save is pending
- **WHEN** the operator edits the range form after submitting a save and that older request succeeds
- **THEN** only the submitted generation is settled, the newer input remains dirty and visible in the open modal, and the interface identifies the remaining unsaved edits instead of silently discarding them.

### Requirement: Source editing begins with selected items and bounded discovery
Opening a recording for editing SHALL show only its draft members with per-member method and interval controls and one Select Points entry. Discovery SHALL use one centered modal with one search/device/point-type toolbar controlling the same result table and matching counts. Query or filter changes SHALL reset the result page, while canonical pending selection persists across queries, filters and pages in the same scope. Large results SHALL use deterministic canonical identity ordering and bounded 25-or-50-row pages; short results SHALL NOT require pagination. Point type SHALL distinguish gauge, counter and state from decoding datatype. Loading, failure, empty data and no matches SHALL remain distinct, and loaded zero-point devices SHALL remain filterable without claiming a new backend search index.

The picker SHALL stage additive pending membership separately from recording draft members and source-form buffers. Existing members SHALL be marked already added and SHALL NOT be removable or duplicated in the picker. A persistent footer SHALL show the pending-new count and Cancel/Add N actions derived from the same canonical set, including how many pending points are outside the current result. Add SHALL be disabled at zero and SHALL merge only validated pending identities into the draft, not save or start the recording. Cancel, Close and Escape SHALL discard only pending additions, preserve draft values and invalid input text, and return focus without prompting to discard the entire recording. Backdrop activation SHALL NOT silently dismiss the picker. Source-form dirty guards and unresolved-operation recovery SHALL remain intact, with only one active modal or confirmation view.

Removal SHALL occur on the main draft with undo and SHALL NOT delete shared sources. Undo or re-adding the same canonical identity within the draft session SHALL preserve prior member method, interval and edited input intent; only genuinely new members SHALL receive defaults. Scope changes or invalid source evidence SHALL NOT revive old member authority.

#### Scenario: Filters and the action count describe the same selection
- **WHEN** five draft members already exist and the operator stages three new points across two devices, then changes the device filter so only one pending point is visible
- **THEN** the result table and matching count follow that filter, the footer reports three pending additions with two outside the result, the five existing members remain locked as already added, and Add merges exactly those three validated points without overwriting existing member configuration.

#### Scenario: Pending additions are cancelled
- **WHEN** an operator stages new points and cancels, closes or presses Escape
- **THEN** only those pending additions are discarded, original members and invalid interval text remain unchanged, no recording save or activation occurs, and focus returns to Select Points.

#### Scenario: A member is removed and restored
- **WHEN** the operator removes a member with an edited mean method and interval, then undoes or explicitly re-adds the same canonical identity in the same draft session
- **THEN** its previous method and interval intent are restored without deleting the source, substituting a same-named identity or applying new-member defaults.

#### Scenario: Source creation has one explicit target
- **WHEN** an operator enters the secondary source-creation flow from all-device discovery or a specific-device filter
- **THEN** the same modal replaces the picker with setup, all-device context requires an explicit existing-or-new device choice with no first-device default, and specific-device context visibly fixes that target without a competing selector; no filter silently changes an active setup target.

#### Scenario: A recording is opened in a large workspace
- **WHEN** the workspace contains 100 devices and 300 rules as a test fixture
- **THEN** the initial editor shows selected sources only, and explicit discovery renders bounded pages rather than all devices or points; selected counts remain available while searching.

#### Scenario: A selected point leaves the current search result
- **WHEN** a query or page no longer displays an already selected point
- **THEN** that point stays selected unless explicitly removed, no entire device is selected implicitly, and returning to the matching query shows the same canonical selection.

#### Scenario: A new recording starts without selected points
- **WHEN** an operator opens source discovery for a new recording with no selected points
- **THEN** the single secondary source flow exposes new-device creation and an existing loaded device with no points remains filterable with a contextual add-read-range action; filtered no-match states offer clearing conditions before creating a duplicate source, without requiring a previously selected point or conflating empty data with loading failure.

#### Scenario: A newly created range returns to discoverable points
- **WHEN** matching confirmed persistence creates a new range from the recording chooser, including on an empty existing device, a device with paginated points, or a newly created device
- **THEN** the setup step closes into bounded discovery focused on that device and range, conflicting text or point-type filters are cleared with a visible explanation and the page containing the new canonical points is shown, so the operator can select them without reopening the editor; existing selections and per-member calculation and interval settings remain unchanged, and no whole device or range is automatically selected.

#### Scenario: Source setup is cancelled from the chooser
- **WHEN** an operator enters device or range setup from discovery and cancels its unsent edits through the dirty guard
- **THEN** a single modal view stack returns to the prior chooser query, filters, page, scroll and focus with the recording name, destination, other unsaved fields, draft membership, pending additions, per-member configuration and invalid input text intact, rather than stacking dialogs or resetting the recording draft.

#### Scenario: A saved device still needs its first range
- **WHEN** a new device is confirmed saved but no range has been saved for it
- **THEN** the interface explicitly states that the device is saved but points are not yet created, and its empty-source state provides an immediate add-read-range action and a return path to discovery without inventing points; cancelling range setup preserves the saved device, and later cancelling the picker or recording does not roll back an explicitly saved device or range; device and range persistence remain separate explicit operations rather than a claimed atomic transaction.

#### Scenario: A source is actually removed or the workspace changes
- **WHEN** canonical source evidence becomes invalid or the operator changes workspace
- **THEN** invalid selection is identified and blocks unsafe saving without substituting a same-named source, cross-workspace state is isolated and stale responses cannot restore old selection authority.

### Requirement: Per-record member controls consume mixed calculation intent
After support-per-member-recording-calculation passes its production gates, the editor SHALL expose latest, mean or delta for each selected recording member with its own interval and UTC-epoch aligned mixed-profile window. This intent SHALL belong to the recording mapping, not the global Point, and SHALL remain separate from acquisition word order. Applicable choices, result types, revisions and completeness behavior SHALL come from the canonical mixed contract. Group defaults MAY initialize new members only. Legacy profiles MUST NOT silently become mixed or reinterpret raw columns as delta results.

#### Scenario: Two means and one delta are configured
- **WHEN** an operator assigns mean, mean and delta to three eligible members in one managed-wide recording
- **THEN** the editor represents independently aligned member windows and a union-boundary wide row with independent results, illustrates 22,12,8, and does not split the request into separate records or modify those sources' use in another recording.

#### Scenario: Calculation or result availability changes
- **WHEN** a member calculation is changed or its delta has no trustworthy current-window result
- **THEN** stale readiness and consent are invalidated, the UI explains required wide-row skip or long-member-row skip or confirmed optional NULL with a reason, and no first-baseline zero, raw substitute or cross-window span is presented as current-window usage.

#### Scenario: Members have different due boundaries
- **WHEN** a record combines 15-second means and a 60-second delta
- **THEN** the UI and acceptance evidence show wide NULL with not_due at :15/:30/:45, both due calculations with their own [start,end) at :00, and long rows only for due members, without carrying delta or inserting zero; required-not-due does not skip a row and due failures retain distinct reasons and coverage.

### Requirement: Acquisition request preview consumes the canonical planner
The consumer SHALL depend on integrate-protocol-aware-acquisition-batching for canonical scoped revision-bound acquisition-plan previews and actual observations, not derive its own batching plan. It SHALL label planned protocol read requests per scan excluding fault retries as estimates, distinguish acquisition Point count and unique wire units, and retain capability, coverage and fallback reasons. Actual attempts, successes, failures, retries, fallback and cancellations SHALL be separately labelled with matching scope/revision; absent observations MUST NOT become zero or measured speedup. Preview SHALL perform no PLC I/O, reject stale revisions, and MUST NOT be called TCP packet counts or SQL writes. This dependency SHALL NOT redefine member recording windows or make recording selection the acquisition set.

#### Scenario: Estimate exists before actual acquisition evidence
- **WHEN** a canonical preview reports three planned protocol requests but no matching actual scan observations exist
- **THEN** the UI labels three as an estimate, reports actual evidence unavailable rather than zero, and a later stale-revision response invalidates that preview instead of retaining a false current count.
