# simple-recording-editor Specification

## Purpose

Provide a single progressive recording editor that reuses existing Studio device, point, mapping, preparation and scoped-start contracts without removing advanced capabilities.

## Requirements

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




---
### Requirement: Basic and advanced recording semantics remain faithful
The system SHALL allow basic typed raw recording without mandatory measurement semantics. Explicit derived operations SHALL require their existing semantics. Layout SHALL remain independent from measurement meaning, and configurations not representable by the simple editor SHALL retain an accurate summary and scoped advanced access without automatic conversion.

#### Scenario: Mixed raw sensors are selected
- **WHEN** temperature, pressure and a raw counter are selected for a supported basic latest-value recording
- **THEN** the editor does not require usage-calculation policies merely because one value is cumulative.

#### Scenario: Unsupported simple-mode combination is opened
- **WHEN** an existing cross-device, custom-table or derived configuration cannot be represented by the simple editor
- **THEN** its identities and stored configuration remain unchanged and the operator can open the matching advanced editor with an explanation.




---
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




---
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

#### Scenario: A saved mixed recording has sources on multiple devices
- **WHEN** an operator explicitly starts one saved mixed recording whose members require two selected devices
- **THEN** the existing scoped-start owner includes exactly those required devices, leaves an unrelated third workspace device untouched, and preserves the recording identity without creating a Basic record or converting legacy configurations.
- **AND** legacy Basic intents without the explicit saved-record discriminator retain their original single-device validation.

#### Scenario: A start reply is lost before its operation ID is received
- **WHEN** the original scoped start request is retained but its operation ID is unknown
- **THEN** the operator SHALL be able to query the existing ledger using that original request through a read-only lookup, with the same workspace, request-ID and normalized-payload-digest checks.
- **AND** lookup SHALL NOT claim, resume, apply or activate any work, and a found operation SHALL return to the original operation-ID recovery owner.

#### Scenario: Original-request lookup cannot establish an operation
- **WHEN** lookup finds no operation, fails, or receives a different workspace or normalized payload for the same request ID
- **THEN** it SHALL remain unknown or explicitly reject the mismatch, without treating absence as cancellation or permission to send another start request.




---
### Requirement: Closing an editor does not roll back submitted source changes
The system SHALL distinguish unsent modal fields, submitted source changes and confirmed saved source changes. A successful source save SHALL remain saved when the recording form is discarded, and a pending or unknown save SHALL use existing settlement or recovery rather than being labelled cancelled.

#### Scenario: Source save succeeds before recording cancellation
- **WHEN** a new point range has been explicitly saved and the operator later discards the unsaved recording form
- **THEN** the interface states that the saved range remains, discards only unsent recording intent and does not issue an implicit compensating delete.

#### Scenario: A source save outcome is unknown
- **WHEN** the operator closes a modal after submitting a save whose result is unconfirmed
- **THEN** the unresolved save is retained for settlement or status recovery and is not described as rolled back or cancelled.




---
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




---
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




---
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




---
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




---
### Requirement: Compact forms and short confirmations preserve useful density
Common device connection parameters SHALL remain visible in a compact responsive form rather than being hidden as advanced details. Point configuration SHALL use compact comparable rows for address, datatype, span and effective decoding order, while uncommon details MAY be disclosed progressively. Short discard, unsaved-change, stop and remove confirmations SHALL use semantic short-confirmation sizing with auto or fit-content height and bounded body scrolling, not the regular 75-percent form-modal height.

#### Scenario: Unsaved edits need confirmation
- **WHEN** a dirty editor requests discard confirmation
- **THEN** a compact centered short confirmation identifies unsaved changes without stretching to the regular modal height, while retaining explicit choice, focus return and dirty-state safety.

#### Scenario: An operator configures common parameters
- **WHEN** device or point configuration is opened
- **THEN** common connection fields and compact point rows are immediately visible with their labels, help and validation rather than hidden solely to lower apparent information density.




---
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




---
### Requirement: Recording choices explain outcomes without inventing restrictions
The editor SHALL describe recording choices by operator intent, concise plain-language labels and illustrative values, rather than the phrase single calculation per rule. Latest SHALL mean the eligible latest value in the recording window; mean SHALL describe an equal-weight arithmetic mean of eligible observations for the applicable gauge member, not a cross-point or time-weighted mean. Counter usage SHALL describe change against the confirmed reliable baseline under the existing counter policy. Per-member effective roles and results SHALL remain faithful: a mean profile can average gauges while retaining latest counter_raw/state values. The UI MUST NOT claim a universal same-algorithm restriction, invent unsupported algorithms beyond the separately specified mixed member contract, silently change saved profiles or present samples as real SQL receipts.

#### Scenario: An operator compares recording choices
- **WHEN** the editor explains eligible readings 20, 22 and 24
- **THEN** it illustrates latest as 24 and arithmetic mean as 22 with understandable labels and notes that examples are simulated, while each selected member's actual role remains visible.

#### Scenario: A cumulative reading has no proven change
- **WHEN** a counter has no baseline or continuity cannot be established
- **THEN** the explanation states that usage is unavailable with its actual reason rather than zero; confirmed signed-counter or rollover policies are not incorrectly described as universally invalid negative changes.




---
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




---
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




---
### Requirement: Per-record member controls consume mixed calculation intent
After support-per-member-recording-calculation passes its production gates, the editor SHALL expose latest, mean or delta for each selected recording member with its own interval and UTC-epoch aligned mixed-profile window. This intent SHALL belong to the recording mapping, not the global Point, and SHALL remain separate from acquisition word order. Applicable choices, result types, revisions and completeness behavior SHALL come from the canonical mixed contract. Group defaults MAY initialize new members only. Legacy profiles MUST NOT silently become mixed or reinterpret raw columns as delta results.

The draft capability SHALL authorize only canonical draft operations, not substitute for a false editor-support capability. The integrated consumer SHALL be validated before the server advertises editor support for release; the capability and its contract tests SHALL be updated together. Editor support SHALL remain distinct from each record's fresh readiness, schema consent, source validation, CAS and ownership. No global capability SHALL authorize activation by itself.

#### Scenario: Draft support precedes an integrated editor
- **WHEN** canonical mixed drafts are supported but the editor consumer has not passed its integration gates
- **THEN** draft support does not claim editor support; after the consumer passes, the server advertises the integrated editor while records still require their own fresh activation proofs.

Consumer authoring acceptance SHALL distinguish Save, schema consent, record Apply and source activation, and SHALL preserve canonical skipped/not_due/no-source readback. End-to-end successful value/window evidence SHALL be completed through the existing owned scoped-start flow before overall editor acceptance. The mixed recording SHALL NOT require an auxiliary Basic recording solely to start its sources, and record Apply alone MUST NOT be represented as source activation.

#### Scenario: Two means and one delta are configured
- **WHEN** an operator assigns mean, mean and delta to three eligible members in one managed-wide recording
- **THEN** the editor represents independently aligned member windows and a union-boundary wide row with independent results, illustrates 22,12,8, and does not split the request into separate records or modify those sources' use in another recording.

#### Scenario: Calculation or result availability changes
- **WHEN** a member calculation is changed or its delta has no trustworthy current-window result
- **THEN** stale readiness and consent are invalidated, the UI explains required wide-row skip or long-member-row skip or confirmed optional NULL with a reason, and no first-baseline zero, raw substitute or cross-window span is presented as current-window usage.

#### Scenario: Members have different due boundaries
- **WHEN** a record combines 15-second means and a 60-second delta
- **THEN** the UI and acceptance evidence show wide NULL with not_due at :15/:30/:45, both due calculations with their own [start,end) at :00, and long rows only for due members, without carrying delta or inserting zero; required-not-due does not skip a row and due failures retain distinct reasons and coverage.




---
### Requirement: Acquisition request preview consumes the canonical planner
The consumer SHALL depend on integrate-protocol-aware-acquisition-batching for canonical scoped revision-bound acquisition-plan previews and actual observations, not derive its own batching plan. It SHALL label planned protocol read requests per scan excluding fault retries as estimates, distinguish acquisition Point count and unique wire units, and retain capability, coverage and fallback reasons. Actual attempts, successes, failures, retries, fallback and cancellations SHALL be separately labelled with matching scope/revision; absent observations MUST NOT become zero or measured speedup. Preview SHALL perform no PLC I/O, reject stale revisions, and MUST NOT be called TCP packet counts or SQL writes. This dependency SHALL NOT redefine member recording windows or make recording selection the acquisition set.

The consumer SHALL display per-group estimates and unavailable states without deriving plans from selected recording members. It SHALL distinguish admitted wire attempts and overlapping retry/fallback subsets from logical request final outcomes, and shall not assert equality between those units or label them TCP packets, SQL writes, or measured speedup. Uninstalled or stale previews SHALL not remain usable; an installed estimate with unavailable actual evidence SHALL remain explicitly an estimate.

#### Scenario: Estimate exists before actual acquisition evidence
- **WHEN** a canonical preview reports three planned protocol requests but no matching actual scan observations exist
- **THEN** the UI labels three as an estimate, reports actual evidence unavailable rather than zero, and a later stale-revision response invalidates that preview instead of retaining a false current count.

---
### Requirement: Destination locks are scoped and explained

The editor SHALL decide destination field locks from one rule: a per-point (mixed) recording SHALL keep its schema and table editable until its first save, choosable from the destination catalog or entered manually, while its storage kind SHALL stay managed; saved per-point, profile and long-table recordings SHALL keep storage kind, schema, table and layout fixed. Catalog tables owned by another recording SHALL be unavailable for selection, and the server ownership checks SHALL remain authoritative. Every locked destination field SHALL show a plain-language reason next to it, and a saved locked recording SHALL offer the existing copy-as-new-recording action.

#### Scenario: New per-point recording chooses its table

- **WHEN** an operator creates a per-point recording and has not saved it
- **THEN** the schema and table can be selected from the catalog or typed, and the managed-only storage kind is shown with its reason

#### Scenario: Saved long-table recording explains the lock

- **WHEN** an operator opens a saved long-table recording
- **THEN** storage kind, schema, table and layout are disabled, a sentence states that long-table recordings share one table per schema and are fixed after saving, and copy-as-new is offered


<!-- @trace
source: clarify-recording-editor-constraints
updated: 2026-11-06
code:
  - internal/datalink/collector/scheduler_acquisition_preview.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection.tsx
  - internal/datalink/dbtarget/writer_idle.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields.tsx
  - internal/datalink/collector/scheduler_acquisition_dispatch.go
  - frontend/src/types/acquisitionAnomaly.ts
  - internal/datalink/collector/scheduler_auto_interval.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/anomaly/types.go
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice.tsx
  - internal/api/handlers/device_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/protocol/mcprotocol/random_read.go
  - internal/protocol/mcprotocol/client_planned_exact.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/point/service_device_polling_group.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/types/acquisitionSettings.ts
  - internal/datalink/device/acquisition_validation.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingAcquisitionSettings.css
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary.tsx
  - internal/datalink/dbtarget/writer_config_busy.go
  - frontend/src/i18n/locales/en/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/scatteredRange.ts
  - internal/datalink/sourcerule/point_runtime_compensation.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - internal/protocol/mcprotocol/const.go
  - internal/datalink/optimizer/acquisition_contiguous.go
  - internal/datalink/dbtarget/writer_live_check.go
  - internal/datalink/dbtarget/writer_retry_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel.tsx
  - frontend/src/i18n/locales/en/recording-home.json
  - internal/datalink/dbtarget/writer_workers.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/hooks/datalink/useAcquisitionAnomalies.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/i18n/locales/zh-TW/parameter-help-advanced.json
  - internal/datalink/sourcerule/sql_repo.go
  - internal/datalink/runtime/anomaly_monitor.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/runtime_rule_status_projection.go
  - internal/datalink/sourcerule/service_requests.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection.tsx
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups.tsx
  - internal/datalink/measurement/source_rule_format.go
  - internal/datalink/sourcerule/format_preflight.go
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/optimizer/acquisition_coverage_validation.go
  - internal/datalink/device/service.go
  - frontend/src/services/studioV2LocatorEntry.ts
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupHelpers.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/instantaneousReduction.ts
  - internal/api/router_acquisition_anomaly.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/pollinggroup/acquisition_storage.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingEditorLayout.css
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAcknowledgement.ts
  - frontend/src/pages/datalink/workbench-v2/deviceSourceSubmission.ts
  - internal/datalink/collector/scheduler_acquisition_publication.go
  - internal/datalink/collector/scheduler_acquisition_execute.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation.tsx
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/anomaly/detector.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/services/acquisitionAnomalies.ts
  - internal/datalink/connector/adapters/mc3e_planned_read.go
  - internal/datalink/device/repository_memory.go
  - internal/datalink/tag/workspace_save_rollback.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/service_read_profile.go
  - internal/datalink/optimizer/acquisition_mc_profile.go
  - internal/datalink/runtime/service_lifecycle.go
  - frontend/src/i18n/locales/en/parameter-help-advanced.json
  - internal/datalink/sourcerule/mixed_identity_retirement.go
  - internal/datalink/dbtarget/writer_refresh_tags.go
  - internal/datalink/point/memory_repo_materialized_batch.go
  - internal/datalink/mapping/source_rule_retirement.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/sourcerule/mutation_rollback_helpers.go
  - internal/datalink/optimizer/acquisition_candidates.go
  - frontend/src/i18n/locales/en/recording-editor-constraints.json
  - internal/datalink/connector/read_profile_identity.go
  - internal/protocol/mcprotocol/client.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingInstantaneousMethod.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker.tsx
  - internal/datalink/optimizer/acquisition_builder.go
  - docs/swagger/docs.go
  - internal/datalink/pollinggroup/sql_repo.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/acquisition_anomaly_types.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer.ts
  - frontend/src/types/measurement.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
  - internal/datalink/dbtarget/tooling_types.go
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/tag_mapping_sync_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/sql_repo_acquisition.go
  - internal/datalink/connector/acquisition_profile_decode.go
  - internal/datalink/device/errors.go
  - internal/datalink/schema/schema_point_models.go
  - internal/datalink/sourcerule/mixed_identity_mapping.go
  - internal/datalink/sourcerule/repository_rule_delete_revision.go
  - internal/datalink/sourcerule/mc_profile_address_plan.go
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/sourcerule/point_update_plan.go
  - internal/datalink/tag/source_rule_retirement.go
  - internal/datalink/sourcerule/candidate_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer.ts
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/dbtarget/write_history.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_group_snapshot.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyDialog.tsx
  - frontend/src/pages/datalink/workbench-v2/useDeviceSourceSubmission.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/migrator_acquisition_anomalies.go
  - frontend/src/services/runtimeStatusProjection.ts
  - internal/datalink/point/sql_repo_device_acquisition.go
  - internal/datalink/anomaly/detector_recovered.go
  - internal/datalink/sourcerule/mixed_identity_links.go
  - internal/datalink/collector/scheduler_dispatch.go
  - internal/datalink/sourcerule/acquisition_profile_source_identity.go
  - frontend/src/i18n/locales/zh-TW/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts
  - internal/datalink/runtime/mixed_source_layout.go
  - internal/datalink/dbtarget/writer_destination.go
  - internal/datalink/connector/acquisition_profile.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/router.go
  - internal/datalink/anomaly/detector_scan.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingAcquisitionSettings.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/workspace/write_group_source_format_identity.go
  - internal/datalink/point/repository_memory.go
  - internal/datalink/point/format_mutation_scope.go
  - internal/datalink/sourcerule/service_update.go
  - internal/datalink/schema/migrations/053_device_acquisition_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/pages/datalink/workbench-v2/rangeSourceReadback.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScatteredRangeSection.tsx
  - internal/protocol/mcprotocol/client_word_read.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/anomaly/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/recording-editor-constraints.json
  - docs/swagger/swagger.json
  - internal/api/handlers/acquisition_anomaly_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/recording-home/acquisitionAnomalies.css
  - internal/datalink/optimizer/acquisition_random.go
  - internal/datalink/point/service_materialized_batch.go
  - internal/datalink/point/format_recovery_snapshot.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/components/parameterHelpContent.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton.tsx
  - internal/datalink/collector/scheduler_acquisition_snapshot.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/schema/migrations/054_acquisition_anomalies_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - internal/datalink/optimizer/acquisition_plan.go
  - internal/datalink/connector/acquisition_profile_schema.go
  - internal/datalink/point/source_rule_retirement.go
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/sourcerule/mixed_review_decision_cas.go
  - internal/datalink/sourcerule/share_recovery_sql_reads.go
  - internal/datalink/dbtarget/writer_skip.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent.tsx
  - internal/datalink/point/service_device_acquisition_memory.go
  - internal/datalink/sourcerule/mixed_rule.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingCompactMemberTable.tsx
  - internal/datalink/point/sql_repo_materialized_batch.go
  - internal/datalink/connector/read_profile.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/sourcerule/tag_review_decision_stale.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/point/service_polling_group.go
  - internal/datalink/sourcerule/acquisition_source_exclusions.go
tests:
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - cmd/test_ui/target_writer_async_test.go
  - cmd/test_ui/acquisition_scheduled_safe_mixed_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-sections.test.tsx
  - internal/datalink/anomaly/detector_test.go
  - internal/api/handlers/mixed_address_profile_regression_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - internal/datalink/sourcerule/service_device_polling_group_test.go
  - frontend/tests/unit/runtime-dashboard/diagnostics-no-write-claim.test.tsx
  - internal/datalink/connector/adapters/mc3e_planned_peer_test.go
  - internal/api/handlers/acquisition_anomaly_retention_regression_test.go
  - internal/datalink/sourcerule/mixed_identity_retirement_runtime_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.tsx
  - frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx
  - internal/datalink/sourcerule/acquisition_qiqr_geometry_regression_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/datalink/migrator_device_acquisition_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.tsx
  - internal/datalink/connector/acquisition_profile_radix_review_test.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/workbench-v2/recording-instantaneous-method.test.tsx
  - internal/api/handlers/runtime_rule_status_safety_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/connector/adapters/mc3e_planned_fixture_test.go
  - internal/datalink/dbtarget/writer_ownership_async_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - internal/datalink/connector/acquisition_profile_random_limit_test.go
  - frontend/tests/unit/workbench-v2/recording-hidden-and-pending-areas.test.tsx
  - internal/datalink/pollinggroup/sql_repo_test.go
  - internal/datalink/dbtarget/writer_review_busy_test.go
  - frontend/tests/unit/workbench-v2/recording-confirmation-readable.test.tsx
  - internal/datalink/sourcerule/mixed_rule_test.go
  - internal/api/router_acquisition_anomaly_test.go
  - internal/datalink/optimizer/acquisition_random_long_timer_regression_test.go
  - internal/datalink/point/format_mutation_sql_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/sourcerule/mixed_review_decision_memory_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/api/handlers/acquisition_anomaly_handler_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/writer_review_refresh_test.go
  - internal/datalink/pollinggroup/sql_repo_owner_test.go
  - cmd/test_ui/acquisition_scheduled_mixed_test.go
  - internal/datalink/collector/scheduler_auto_interval_test.go
  - frontend/tests/unit/workbench-v2/scattered-metadata-roundtrip.test.ts
  - internal/datalink/optimizer/acquisition_random_contract_test.go
  - frontend/tests/unit/services/runtimeStatusSafety.test.ts
  - internal/datalink/runtime/status_auto_interval_test.go
  - frontend/tests/unit/scatteredRangeForm.test.tsx
  - internal/datalink/runtime/mixed_rule_format_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-unknown-codes.test.ts
  - internal/api/handlers/mixed_item_replacement_concurrency_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-destination-lock-notice.test.tsx
  - internal/datalink/anomaly/sql_repo_test.go
  - internal/api/handlers/mixed_item_replacement_regression_test.go
  - internal/protocol/mcprotocol/mcprotocol_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/device/service_acquisition_test.go
  - internal/datalink/optimizer/acquisition_review_regression_test.go
  - internal/datalink/point/service_polling_group_owner_test.go
  - internal/datalink/dbtarget/writer_backoff_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-editor.test.tsx
  - internal/datalink/sourcerule/acquisition_source_mixed_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-copy.test.tsx
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - frontend/tests/unit/services/studio-v2-locator-acquisition.test.ts
  - frontend/tests/unit/workbench-v2/recording-method-groups.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-plain-language.test.tsx
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/anomaly/detector_restart_regression_test.go
  - internal/protocol/mcprotocol/client_planned_exact_test.go
  - internal/datalink/pollinggroup/acquisition_mutation_test.go
  - frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts
  - internal/datalink/migrator_acquisition_anomalies_test.go
  - internal/datalink/dbtarget/writer_review_timeout_test.go
  - internal/datalink/workspace/mixed_rule_format_regression_test.go
  - internal/datalink/anomaly/testhelpers_test.go
  - internal/datalink/collector/scheduler_qiqr_shape_regression_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - internal/datalink/runtime/service_device_sync_test.go
  - internal/api/handlers/mixed_item_review_decision_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx
  - frontend/tests/unit/workbench-v2/recording-schema-instruction.test.tsx
  - internal/datalink/point/service_device_acquisition_concurrency_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - cmd/test_ui/target_writer.go
  - cmd/test_ui/bug_a_activation_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.html
  - internal/datalink/dbtarget/writer_recovery_isolation_test.go
  - internal/datalink/point/mixed_rule_format_regression_test.go
  - internal/datalink/connector/acquisition_profile_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - frontend/tests/unit/workbench-v2/recording-instant-role-label.test.tsx
  - internal/datalink/dbtarget/writer_starvation_test.go
  - internal/datalink/point/service_device_acquisition_atomic_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-source-refresh.test.tsx
  - internal/protocol/mcprotocol/client_validation_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.html
  - frontend/tests/unit/hooks/useAcquisitionAnomalies.test.tsx
  - frontend/tests/unit/scatteredRangeReviewFixes.test.tsx
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - frontend/tests/unit/services/acquisitionAnomalies.test.ts
  - internal/datalink/runtime/ingestor_database_decoupling_test.go
  - internal/datalink/device/sql_repo_acquisition_test.go
  - internal/api/handlers/dbtarget_handler_tooling_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/datalink/runtime/ingestor_async_acceptance_test.go
  - internal/protocol/mcprotocol/client_word_read_test.go
  - internal/datalink/anomaly/detector_observation_regression_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-method.test.tsx
  - internal/api/handlers/studio_v2_workspace_rule_switch_test.go
  - internal/datalink/runtime/anomaly_monitor_test.go
  - internal/datalink/workspace/runtime_acquisition_mixed_test.go
-->

---
### Requirement: Recording method is visible where it is named

For a per-point recording, the recording method section SHALL list the members grouped by calculation and interval with point counts, and each group SHALL offer an action that scrolls to and focuses the first matching member's calculation control in the source table. For a recording without a recording profile, the section SHALL state "latest (legacy wide table)" with its saved row interval and SHALL offer creating a new recording prefilled with the same points and destination connector; it SHALL NOT convert the recording in place.

#### Scenario: Grouped per-point methods

| Members | Shown groups |
| ------- | ------------ |
| 3 latest every 15 s, 1 mean every 60 s | "latest · 15 s · 3 points", "mean · 60 s · 1 point" |

- **WHEN** the operator activates the "mean · 60 s" group's action
- **THEN** the source table scrolls to that member and its calculation control receives focus

#### Scenario: Legacy recording is summarized, not converted

- **WHEN** a recording without a recording profile and a 60 second interval is opened
- **THEN** the method section shows latest (legacy wide table) every 60 seconds and a shortcut to create a new recording, and the saved recording is unchanged


<!-- @trace
source: clarify-recording-editor-constraints
updated: 2026-11-06
code:
  - internal/datalink/collector/scheduler_acquisition_preview.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection.tsx
  - internal/datalink/dbtarget/writer_idle.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields.tsx
  - internal/datalink/collector/scheduler_acquisition_dispatch.go
  - frontend/src/types/acquisitionAnomaly.ts
  - internal/datalink/collector/scheduler_auto_interval.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/anomaly/types.go
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice.tsx
  - internal/api/handlers/device_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/protocol/mcprotocol/random_read.go
  - internal/protocol/mcprotocol/client_planned_exact.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/point/service_device_polling_group.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/types/acquisitionSettings.ts
  - internal/datalink/device/acquisition_validation.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingAcquisitionSettings.css
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary.tsx
  - internal/datalink/dbtarget/writer_config_busy.go
  - frontend/src/i18n/locales/en/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/scatteredRange.ts
  - internal/datalink/sourcerule/point_runtime_compensation.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - internal/protocol/mcprotocol/const.go
  - internal/datalink/optimizer/acquisition_contiguous.go
  - internal/datalink/dbtarget/writer_live_check.go
  - internal/datalink/dbtarget/writer_retry_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel.tsx
  - frontend/src/i18n/locales/en/recording-home.json
  - internal/datalink/dbtarget/writer_workers.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/hooks/datalink/useAcquisitionAnomalies.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/i18n/locales/zh-TW/parameter-help-advanced.json
  - internal/datalink/sourcerule/sql_repo.go
  - internal/datalink/runtime/anomaly_monitor.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/runtime_rule_status_projection.go
  - internal/datalink/sourcerule/service_requests.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection.tsx
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups.tsx
  - internal/datalink/measurement/source_rule_format.go
  - internal/datalink/sourcerule/format_preflight.go
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/optimizer/acquisition_coverage_validation.go
  - internal/datalink/device/service.go
  - frontend/src/services/studioV2LocatorEntry.ts
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupHelpers.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/instantaneousReduction.ts
  - internal/api/router_acquisition_anomaly.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/pollinggroup/acquisition_storage.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingEditorLayout.css
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAcknowledgement.ts
  - frontend/src/pages/datalink/workbench-v2/deviceSourceSubmission.ts
  - internal/datalink/collector/scheduler_acquisition_publication.go
  - internal/datalink/collector/scheduler_acquisition_execute.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation.tsx
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/anomaly/detector.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/services/acquisitionAnomalies.ts
  - internal/datalink/connector/adapters/mc3e_planned_read.go
  - internal/datalink/device/repository_memory.go
  - internal/datalink/tag/workspace_save_rollback.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/service_read_profile.go
  - internal/datalink/optimizer/acquisition_mc_profile.go
  - internal/datalink/runtime/service_lifecycle.go
  - frontend/src/i18n/locales/en/parameter-help-advanced.json
  - internal/datalink/sourcerule/mixed_identity_retirement.go
  - internal/datalink/dbtarget/writer_refresh_tags.go
  - internal/datalink/point/memory_repo_materialized_batch.go
  - internal/datalink/mapping/source_rule_retirement.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/sourcerule/mutation_rollback_helpers.go
  - internal/datalink/optimizer/acquisition_candidates.go
  - frontend/src/i18n/locales/en/recording-editor-constraints.json
  - internal/datalink/connector/read_profile_identity.go
  - internal/protocol/mcprotocol/client.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingInstantaneousMethod.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker.tsx
  - internal/datalink/optimizer/acquisition_builder.go
  - docs/swagger/docs.go
  - internal/datalink/pollinggroup/sql_repo.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/acquisition_anomaly_types.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer.ts
  - frontend/src/types/measurement.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
  - internal/datalink/dbtarget/tooling_types.go
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/tag_mapping_sync_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/sql_repo_acquisition.go
  - internal/datalink/connector/acquisition_profile_decode.go
  - internal/datalink/device/errors.go
  - internal/datalink/schema/schema_point_models.go
  - internal/datalink/sourcerule/mixed_identity_mapping.go
  - internal/datalink/sourcerule/repository_rule_delete_revision.go
  - internal/datalink/sourcerule/mc_profile_address_plan.go
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/sourcerule/point_update_plan.go
  - internal/datalink/tag/source_rule_retirement.go
  - internal/datalink/sourcerule/candidate_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer.ts
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/dbtarget/write_history.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_group_snapshot.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyDialog.tsx
  - frontend/src/pages/datalink/workbench-v2/useDeviceSourceSubmission.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/migrator_acquisition_anomalies.go
  - frontend/src/services/runtimeStatusProjection.ts
  - internal/datalink/point/sql_repo_device_acquisition.go
  - internal/datalink/anomaly/detector_recovered.go
  - internal/datalink/sourcerule/mixed_identity_links.go
  - internal/datalink/collector/scheduler_dispatch.go
  - internal/datalink/sourcerule/acquisition_profile_source_identity.go
  - frontend/src/i18n/locales/zh-TW/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts
  - internal/datalink/runtime/mixed_source_layout.go
  - internal/datalink/dbtarget/writer_destination.go
  - internal/datalink/connector/acquisition_profile.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/router.go
  - internal/datalink/anomaly/detector_scan.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingAcquisitionSettings.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/workspace/write_group_source_format_identity.go
  - internal/datalink/point/repository_memory.go
  - internal/datalink/point/format_mutation_scope.go
  - internal/datalink/sourcerule/service_update.go
  - internal/datalink/schema/migrations/053_device_acquisition_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/pages/datalink/workbench-v2/rangeSourceReadback.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScatteredRangeSection.tsx
  - internal/protocol/mcprotocol/client_word_read.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/anomaly/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/recording-editor-constraints.json
  - docs/swagger/swagger.json
  - internal/api/handlers/acquisition_anomaly_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/recording-home/acquisitionAnomalies.css
  - internal/datalink/optimizer/acquisition_random.go
  - internal/datalink/point/service_materialized_batch.go
  - internal/datalink/point/format_recovery_snapshot.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/components/parameterHelpContent.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton.tsx
  - internal/datalink/collector/scheduler_acquisition_snapshot.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/schema/migrations/054_acquisition_anomalies_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - internal/datalink/optimizer/acquisition_plan.go
  - internal/datalink/connector/acquisition_profile_schema.go
  - internal/datalink/point/source_rule_retirement.go
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/sourcerule/mixed_review_decision_cas.go
  - internal/datalink/sourcerule/share_recovery_sql_reads.go
  - internal/datalink/dbtarget/writer_skip.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent.tsx
  - internal/datalink/point/service_device_acquisition_memory.go
  - internal/datalink/sourcerule/mixed_rule.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingCompactMemberTable.tsx
  - internal/datalink/point/sql_repo_materialized_batch.go
  - internal/datalink/connector/read_profile.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/sourcerule/tag_review_decision_stale.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/point/service_polling_group.go
  - internal/datalink/sourcerule/acquisition_source_exclusions.go
tests:
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - cmd/test_ui/target_writer_async_test.go
  - cmd/test_ui/acquisition_scheduled_safe_mixed_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-sections.test.tsx
  - internal/datalink/anomaly/detector_test.go
  - internal/api/handlers/mixed_address_profile_regression_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - internal/datalink/sourcerule/service_device_polling_group_test.go
  - frontend/tests/unit/runtime-dashboard/diagnostics-no-write-claim.test.tsx
  - internal/datalink/connector/adapters/mc3e_planned_peer_test.go
  - internal/api/handlers/acquisition_anomaly_retention_regression_test.go
  - internal/datalink/sourcerule/mixed_identity_retirement_runtime_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.tsx
  - frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx
  - internal/datalink/sourcerule/acquisition_qiqr_geometry_regression_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/datalink/migrator_device_acquisition_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.tsx
  - internal/datalink/connector/acquisition_profile_radix_review_test.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/workbench-v2/recording-instantaneous-method.test.tsx
  - internal/api/handlers/runtime_rule_status_safety_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/connector/adapters/mc3e_planned_fixture_test.go
  - internal/datalink/dbtarget/writer_ownership_async_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - internal/datalink/connector/acquisition_profile_random_limit_test.go
  - frontend/tests/unit/workbench-v2/recording-hidden-and-pending-areas.test.tsx
  - internal/datalink/pollinggroup/sql_repo_test.go
  - internal/datalink/dbtarget/writer_review_busy_test.go
  - frontend/tests/unit/workbench-v2/recording-confirmation-readable.test.tsx
  - internal/datalink/sourcerule/mixed_rule_test.go
  - internal/api/router_acquisition_anomaly_test.go
  - internal/datalink/optimizer/acquisition_random_long_timer_regression_test.go
  - internal/datalink/point/format_mutation_sql_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/sourcerule/mixed_review_decision_memory_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/api/handlers/acquisition_anomaly_handler_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/writer_review_refresh_test.go
  - internal/datalink/pollinggroup/sql_repo_owner_test.go
  - cmd/test_ui/acquisition_scheduled_mixed_test.go
  - internal/datalink/collector/scheduler_auto_interval_test.go
  - frontend/tests/unit/workbench-v2/scattered-metadata-roundtrip.test.ts
  - internal/datalink/optimizer/acquisition_random_contract_test.go
  - frontend/tests/unit/services/runtimeStatusSafety.test.ts
  - internal/datalink/runtime/status_auto_interval_test.go
  - frontend/tests/unit/scatteredRangeForm.test.tsx
  - internal/datalink/runtime/mixed_rule_format_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-unknown-codes.test.ts
  - internal/api/handlers/mixed_item_replacement_concurrency_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-destination-lock-notice.test.tsx
  - internal/datalink/anomaly/sql_repo_test.go
  - internal/api/handlers/mixed_item_replacement_regression_test.go
  - internal/protocol/mcprotocol/mcprotocol_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/device/service_acquisition_test.go
  - internal/datalink/optimizer/acquisition_review_regression_test.go
  - internal/datalink/point/service_polling_group_owner_test.go
  - internal/datalink/dbtarget/writer_backoff_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-editor.test.tsx
  - internal/datalink/sourcerule/acquisition_source_mixed_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-copy.test.tsx
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - frontend/tests/unit/services/studio-v2-locator-acquisition.test.ts
  - frontend/tests/unit/workbench-v2/recording-method-groups.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-plain-language.test.tsx
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/anomaly/detector_restart_regression_test.go
  - internal/protocol/mcprotocol/client_planned_exact_test.go
  - internal/datalink/pollinggroup/acquisition_mutation_test.go
  - frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts
  - internal/datalink/migrator_acquisition_anomalies_test.go
  - internal/datalink/dbtarget/writer_review_timeout_test.go
  - internal/datalink/workspace/mixed_rule_format_regression_test.go
  - internal/datalink/anomaly/testhelpers_test.go
  - internal/datalink/collector/scheduler_qiqr_shape_regression_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - internal/datalink/runtime/service_device_sync_test.go
  - internal/api/handlers/mixed_item_review_decision_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx
  - frontend/tests/unit/workbench-v2/recording-schema-instruction.test.tsx
  - internal/datalink/point/service_device_acquisition_concurrency_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - cmd/test_ui/target_writer.go
  - cmd/test_ui/bug_a_activation_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.html
  - internal/datalink/dbtarget/writer_recovery_isolation_test.go
  - internal/datalink/point/mixed_rule_format_regression_test.go
  - internal/datalink/connector/acquisition_profile_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - frontend/tests/unit/workbench-v2/recording-instant-role-label.test.tsx
  - internal/datalink/dbtarget/writer_starvation_test.go
  - internal/datalink/point/service_device_acquisition_atomic_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-source-refresh.test.tsx
  - internal/protocol/mcprotocol/client_validation_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.html
  - frontend/tests/unit/hooks/useAcquisitionAnomalies.test.tsx
  - frontend/tests/unit/scatteredRangeReviewFixes.test.tsx
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - frontend/tests/unit/services/acquisitionAnomalies.test.ts
  - internal/datalink/runtime/ingestor_database_decoupling_test.go
  - internal/datalink/device/sql_repo_acquisition_test.go
  - internal/api/handlers/dbtarget_handler_tooling_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/datalink/runtime/ingestor_async_acceptance_test.go
  - internal/protocol/mcprotocol/client_word_read_test.go
  - internal/datalink/anomaly/detector_observation_regression_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-method.test.tsx
  - internal/api/handlers/studio_v2_workspace_rule_switch_test.go
  - internal/datalink/runtime/anomaly_monitor_test.go
  - internal/datalink/workspace/runtime_acquisition_mixed_test.go
-->

---
### Requirement: Instantaneous recordings edit their reduction explicitly

For a recording whose profile kind is instantaneous, the editor SHALL offer the reduction (latest or arithmetic mean) and the recording interval within the existing instantaneous bounds. Choosing mean SHALL require the existing gauge-use confirmation. Saving SHALL go through the existing save with revision checks, and the preparation area SHALL then require a fresh schema check and confirmation before the change takes effect. Cumulative recordings SHALL NOT offer mean.

#### Scenario: Switching to mean

- **GIVEN** a saved instantaneous recording using latest every 15 seconds
- **WHEN** the operator selects mean, confirms gauge use and saves
- **THEN** the update carries `sample_reduction` mean, a new revision is saved, and the preparation area asks for a schema check before recording uses it


<!-- @trace
source: clarify-recording-editor-constraints
updated: 2026-11-06
code:
  - internal/datalink/collector/scheduler_acquisition_preview.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection.tsx
  - internal/datalink/dbtarget/writer_idle.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields.tsx
  - internal/datalink/collector/scheduler_acquisition_dispatch.go
  - frontend/src/types/acquisitionAnomaly.ts
  - internal/datalink/collector/scheduler_auto_interval.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/anomaly/types.go
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice.tsx
  - internal/api/handlers/device_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/protocol/mcprotocol/random_read.go
  - internal/protocol/mcprotocol/client_planned_exact.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/point/service_device_polling_group.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/types/acquisitionSettings.ts
  - internal/datalink/device/acquisition_validation.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingAcquisitionSettings.css
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary.tsx
  - internal/datalink/dbtarget/writer_config_busy.go
  - frontend/src/i18n/locales/en/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/scatteredRange.ts
  - internal/datalink/sourcerule/point_runtime_compensation.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - internal/protocol/mcprotocol/const.go
  - internal/datalink/optimizer/acquisition_contiguous.go
  - internal/datalink/dbtarget/writer_live_check.go
  - internal/datalink/dbtarget/writer_retry_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel.tsx
  - frontend/src/i18n/locales/en/recording-home.json
  - internal/datalink/dbtarget/writer_workers.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/hooks/datalink/useAcquisitionAnomalies.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/i18n/locales/zh-TW/parameter-help-advanced.json
  - internal/datalink/sourcerule/sql_repo.go
  - internal/datalink/runtime/anomaly_monitor.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/runtime_rule_status_projection.go
  - internal/datalink/sourcerule/service_requests.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection.tsx
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups.tsx
  - internal/datalink/measurement/source_rule_format.go
  - internal/datalink/sourcerule/format_preflight.go
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/optimizer/acquisition_coverage_validation.go
  - internal/datalink/device/service.go
  - frontend/src/services/studioV2LocatorEntry.ts
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupHelpers.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/instantaneousReduction.ts
  - internal/api/router_acquisition_anomaly.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/pollinggroup/acquisition_storage.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingEditorLayout.css
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAcknowledgement.ts
  - frontend/src/pages/datalink/workbench-v2/deviceSourceSubmission.ts
  - internal/datalink/collector/scheduler_acquisition_publication.go
  - internal/datalink/collector/scheduler_acquisition_execute.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation.tsx
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/anomaly/detector.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/services/acquisitionAnomalies.ts
  - internal/datalink/connector/adapters/mc3e_planned_read.go
  - internal/datalink/device/repository_memory.go
  - internal/datalink/tag/workspace_save_rollback.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/service_read_profile.go
  - internal/datalink/optimizer/acquisition_mc_profile.go
  - internal/datalink/runtime/service_lifecycle.go
  - frontend/src/i18n/locales/en/parameter-help-advanced.json
  - internal/datalink/sourcerule/mixed_identity_retirement.go
  - internal/datalink/dbtarget/writer_refresh_tags.go
  - internal/datalink/point/memory_repo_materialized_batch.go
  - internal/datalink/mapping/source_rule_retirement.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/sourcerule/mutation_rollback_helpers.go
  - internal/datalink/optimizer/acquisition_candidates.go
  - frontend/src/i18n/locales/en/recording-editor-constraints.json
  - internal/datalink/connector/read_profile_identity.go
  - internal/protocol/mcprotocol/client.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingInstantaneousMethod.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker.tsx
  - internal/datalink/optimizer/acquisition_builder.go
  - docs/swagger/docs.go
  - internal/datalink/pollinggroup/sql_repo.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/acquisition_anomaly_types.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer.ts
  - frontend/src/types/measurement.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
  - internal/datalink/dbtarget/tooling_types.go
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/tag_mapping_sync_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/sql_repo_acquisition.go
  - internal/datalink/connector/acquisition_profile_decode.go
  - internal/datalink/device/errors.go
  - internal/datalink/schema/schema_point_models.go
  - internal/datalink/sourcerule/mixed_identity_mapping.go
  - internal/datalink/sourcerule/repository_rule_delete_revision.go
  - internal/datalink/sourcerule/mc_profile_address_plan.go
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/sourcerule/point_update_plan.go
  - internal/datalink/tag/source_rule_retirement.go
  - internal/datalink/sourcerule/candidate_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer.ts
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/dbtarget/write_history.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_group_snapshot.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyDialog.tsx
  - frontend/src/pages/datalink/workbench-v2/useDeviceSourceSubmission.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/migrator_acquisition_anomalies.go
  - frontend/src/services/runtimeStatusProjection.ts
  - internal/datalink/point/sql_repo_device_acquisition.go
  - internal/datalink/anomaly/detector_recovered.go
  - internal/datalink/sourcerule/mixed_identity_links.go
  - internal/datalink/collector/scheduler_dispatch.go
  - internal/datalink/sourcerule/acquisition_profile_source_identity.go
  - frontend/src/i18n/locales/zh-TW/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts
  - internal/datalink/runtime/mixed_source_layout.go
  - internal/datalink/dbtarget/writer_destination.go
  - internal/datalink/connector/acquisition_profile.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/router.go
  - internal/datalink/anomaly/detector_scan.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingAcquisitionSettings.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/workspace/write_group_source_format_identity.go
  - internal/datalink/point/repository_memory.go
  - internal/datalink/point/format_mutation_scope.go
  - internal/datalink/sourcerule/service_update.go
  - internal/datalink/schema/migrations/053_device_acquisition_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/pages/datalink/workbench-v2/rangeSourceReadback.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScatteredRangeSection.tsx
  - internal/protocol/mcprotocol/client_word_read.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/anomaly/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/recording-editor-constraints.json
  - docs/swagger/swagger.json
  - internal/api/handlers/acquisition_anomaly_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/recording-home/acquisitionAnomalies.css
  - internal/datalink/optimizer/acquisition_random.go
  - internal/datalink/point/service_materialized_batch.go
  - internal/datalink/point/format_recovery_snapshot.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/components/parameterHelpContent.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton.tsx
  - internal/datalink/collector/scheduler_acquisition_snapshot.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/schema/migrations/054_acquisition_anomalies_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - internal/datalink/optimizer/acquisition_plan.go
  - internal/datalink/connector/acquisition_profile_schema.go
  - internal/datalink/point/source_rule_retirement.go
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/sourcerule/mixed_review_decision_cas.go
  - internal/datalink/sourcerule/share_recovery_sql_reads.go
  - internal/datalink/dbtarget/writer_skip.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent.tsx
  - internal/datalink/point/service_device_acquisition_memory.go
  - internal/datalink/sourcerule/mixed_rule.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingCompactMemberTable.tsx
  - internal/datalink/point/sql_repo_materialized_batch.go
  - internal/datalink/connector/read_profile.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/sourcerule/tag_review_decision_stale.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/point/service_polling_group.go
  - internal/datalink/sourcerule/acquisition_source_exclusions.go
tests:
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - cmd/test_ui/target_writer_async_test.go
  - cmd/test_ui/acquisition_scheduled_safe_mixed_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-sections.test.tsx
  - internal/datalink/anomaly/detector_test.go
  - internal/api/handlers/mixed_address_profile_regression_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - internal/datalink/sourcerule/service_device_polling_group_test.go
  - frontend/tests/unit/runtime-dashboard/diagnostics-no-write-claim.test.tsx
  - internal/datalink/connector/adapters/mc3e_planned_peer_test.go
  - internal/api/handlers/acquisition_anomaly_retention_regression_test.go
  - internal/datalink/sourcerule/mixed_identity_retirement_runtime_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.tsx
  - frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx
  - internal/datalink/sourcerule/acquisition_qiqr_geometry_regression_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/datalink/migrator_device_acquisition_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.tsx
  - internal/datalink/connector/acquisition_profile_radix_review_test.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/workbench-v2/recording-instantaneous-method.test.tsx
  - internal/api/handlers/runtime_rule_status_safety_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/connector/adapters/mc3e_planned_fixture_test.go
  - internal/datalink/dbtarget/writer_ownership_async_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - internal/datalink/connector/acquisition_profile_random_limit_test.go
  - frontend/tests/unit/workbench-v2/recording-hidden-and-pending-areas.test.tsx
  - internal/datalink/pollinggroup/sql_repo_test.go
  - internal/datalink/dbtarget/writer_review_busy_test.go
  - frontend/tests/unit/workbench-v2/recording-confirmation-readable.test.tsx
  - internal/datalink/sourcerule/mixed_rule_test.go
  - internal/api/router_acquisition_anomaly_test.go
  - internal/datalink/optimizer/acquisition_random_long_timer_regression_test.go
  - internal/datalink/point/format_mutation_sql_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/sourcerule/mixed_review_decision_memory_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/api/handlers/acquisition_anomaly_handler_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/writer_review_refresh_test.go
  - internal/datalink/pollinggroup/sql_repo_owner_test.go
  - cmd/test_ui/acquisition_scheduled_mixed_test.go
  - internal/datalink/collector/scheduler_auto_interval_test.go
  - frontend/tests/unit/workbench-v2/scattered-metadata-roundtrip.test.ts
  - internal/datalink/optimizer/acquisition_random_contract_test.go
  - frontend/tests/unit/services/runtimeStatusSafety.test.ts
  - internal/datalink/runtime/status_auto_interval_test.go
  - frontend/tests/unit/scatteredRangeForm.test.tsx
  - internal/datalink/runtime/mixed_rule_format_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-unknown-codes.test.ts
  - internal/api/handlers/mixed_item_replacement_concurrency_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-destination-lock-notice.test.tsx
  - internal/datalink/anomaly/sql_repo_test.go
  - internal/api/handlers/mixed_item_replacement_regression_test.go
  - internal/protocol/mcprotocol/mcprotocol_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/device/service_acquisition_test.go
  - internal/datalink/optimizer/acquisition_review_regression_test.go
  - internal/datalink/point/service_polling_group_owner_test.go
  - internal/datalink/dbtarget/writer_backoff_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-editor.test.tsx
  - internal/datalink/sourcerule/acquisition_source_mixed_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-copy.test.tsx
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - frontend/tests/unit/services/studio-v2-locator-acquisition.test.ts
  - frontend/tests/unit/workbench-v2/recording-method-groups.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-plain-language.test.tsx
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/anomaly/detector_restart_regression_test.go
  - internal/protocol/mcprotocol/client_planned_exact_test.go
  - internal/datalink/pollinggroup/acquisition_mutation_test.go
  - frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts
  - internal/datalink/migrator_acquisition_anomalies_test.go
  - internal/datalink/dbtarget/writer_review_timeout_test.go
  - internal/datalink/workspace/mixed_rule_format_regression_test.go
  - internal/datalink/anomaly/testhelpers_test.go
  - internal/datalink/collector/scheduler_qiqr_shape_regression_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - internal/datalink/runtime/service_device_sync_test.go
  - internal/api/handlers/mixed_item_review_decision_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx
  - frontend/tests/unit/workbench-v2/recording-schema-instruction.test.tsx
  - internal/datalink/point/service_device_acquisition_concurrency_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - cmd/test_ui/target_writer.go
  - cmd/test_ui/bug_a_activation_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.html
  - internal/datalink/dbtarget/writer_recovery_isolation_test.go
  - internal/datalink/point/mixed_rule_format_regression_test.go
  - internal/datalink/connector/acquisition_profile_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - frontend/tests/unit/workbench-v2/recording-instant-role-label.test.tsx
  - internal/datalink/dbtarget/writer_starvation_test.go
  - internal/datalink/point/service_device_acquisition_atomic_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-source-refresh.test.tsx
  - internal/protocol/mcprotocol/client_validation_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.html
  - frontend/tests/unit/hooks/useAcquisitionAnomalies.test.tsx
  - frontend/tests/unit/scatteredRangeReviewFixes.test.tsx
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - frontend/tests/unit/services/acquisitionAnomalies.test.ts
  - internal/datalink/runtime/ingestor_database_decoupling_test.go
  - internal/datalink/device/sql_repo_acquisition_test.go
  - internal/api/handlers/dbtarget_handler_tooling_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/datalink/runtime/ingestor_async_acceptance_test.go
  - internal/protocol/mcprotocol/client_word_read_test.go
  - internal/datalink/anomaly/detector_observation_regression_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-method.test.tsx
  - internal/api/handlers/studio_v2_workspace_rule_switch_test.go
  - internal/datalink/runtime/anomaly_monitor_test.go
  - internal/datalink/workspace/runtime_acquisition_mixed_test.go
-->

---
### Requirement: Hidden and pending areas state why

For a per-point recording, the advanced row-rule area SHALL explain that intervals are set per point, incomplete rows are fixed to partial results, and managed-table duplicate protection, entity key and provenance columns are system-defined. For an unsaved managed recording, the preparation area SHALL state that the table structure can be checked after saving. When the destination catalog is not being read, the editor SHALL NOT show a loading message; when the read fails, it SHALL state that the destination database could not be reached and offer retry.

#### Scenario: Unsaved managed recording

- **WHEN** a new managed recording has not been saved
- **THEN** the preparation area shows that the structure check becomes available after saving

#### Scenario: Unreachable destination

- **WHEN** the destination catalog request fails
- **THEN** the editor shows that the destination database could not be reached with a retry action instead of a loading message


<!-- @trace
source: clarify-recording-editor-constraints
updated: 2026-11-06
code:
  - internal/datalink/collector/scheduler_acquisition_preview.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection.tsx
  - internal/datalink/dbtarget/writer_idle.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields.tsx
  - internal/datalink/collector/scheduler_acquisition_dispatch.go
  - frontend/src/types/acquisitionAnomaly.ts
  - internal/datalink/collector/scheduler_auto_interval.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/anomaly/types.go
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice.tsx
  - internal/api/handlers/device_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/protocol/mcprotocol/random_read.go
  - internal/protocol/mcprotocol/client_planned_exact.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/point/service_device_polling_group.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/types/acquisitionSettings.ts
  - internal/datalink/device/acquisition_validation.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingAcquisitionSettings.css
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary.tsx
  - internal/datalink/dbtarget/writer_config_busy.go
  - frontend/src/i18n/locales/en/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/scatteredRange.ts
  - internal/datalink/sourcerule/point_runtime_compensation.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - internal/protocol/mcprotocol/const.go
  - internal/datalink/optimizer/acquisition_contiguous.go
  - internal/datalink/dbtarget/writer_live_check.go
  - internal/datalink/dbtarget/writer_retry_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel.tsx
  - frontend/src/i18n/locales/en/recording-home.json
  - internal/datalink/dbtarget/writer_workers.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/hooks/datalink/useAcquisitionAnomalies.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/i18n/locales/zh-TW/parameter-help-advanced.json
  - internal/datalink/sourcerule/sql_repo.go
  - internal/datalink/runtime/anomaly_monitor.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/runtime_rule_status_projection.go
  - internal/datalink/sourcerule/service_requests.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection.tsx
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups.tsx
  - internal/datalink/measurement/source_rule_format.go
  - internal/datalink/sourcerule/format_preflight.go
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/optimizer/acquisition_coverage_validation.go
  - internal/datalink/device/service.go
  - frontend/src/services/studioV2LocatorEntry.ts
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupHelpers.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/instantaneousReduction.ts
  - internal/api/router_acquisition_anomaly.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/pollinggroup/acquisition_storage.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingEditorLayout.css
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAcknowledgement.ts
  - frontend/src/pages/datalink/workbench-v2/deviceSourceSubmission.ts
  - internal/datalink/collector/scheduler_acquisition_publication.go
  - internal/datalink/collector/scheduler_acquisition_execute.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation.tsx
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/anomaly/detector.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/services/acquisitionAnomalies.ts
  - internal/datalink/connector/adapters/mc3e_planned_read.go
  - internal/datalink/device/repository_memory.go
  - internal/datalink/tag/workspace_save_rollback.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/service_read_profile.go
  - internal/datalink/optimizer/acquisition_mc_profile.go
  - internal/datalink/runtime/service_lifecycle.go
  - frontend/src/i18n/locales/en/parameter-help-advanced.json
  - internal/datalink/sourcerule/mixed_identity_retirement.go
  - internal/datalink/dbtarget/writer_refresh_tags.go
  - internal/datalink/point/memory_repo_materialized_batch.go
  - internal/datalink/mapping/source_rule_retirement.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/sourcerule/mutation_rollback_helpers.go
  - internal/datalink/optimizer/acquisition_candidates.go
  - frontend/src/i18n/locales/en/recording-editor-constraints.json
  - internal/datalink/connector/read_profile_identity.go
  - internal/protocol/mcprotocol/client.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingInstantaneousMethod.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker.tsx
  - internal/datalink/optimizer/acquisition_builder.go
  - docs/swagger/docs.go
  - internal/datalink/pollinggroup/sql_repo.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/acquisition_anomaly_types.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer.ts
  - frontend/src/types/measurement.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
  - internal/datalink/dbtarget/tooling_types.go
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/tag_mapping_sync_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/sql_repo_acquisition.go
  - internal/datalink/connector/acquisition_profile_decode.go
  - internal/datalink/device/errors.go
  - internal/datalink/schema/schema_point_models.go
  - internal/datalink/sourcerule/mixed_identity_mapping.go
  - internal/datalink/sourcerule/repository_rule_delete_revision.go
  - internal/datalink/sourcerule/mc_profile_address_plan.go
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/sourcerule/point_update_plan.go
  - internal/datalink/tag/source_rule_retirement.go
  - internal/datalink/sourcerule/candidate_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer.ts
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/dbtarget/write_history.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_group_snapshot.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyDialog.tsx
  - frontend/src/pages/datalink/workbench-v2/useDeviceSourceSubmission.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/migrator_acquisition_anomalies.go
  - frontend/src/services/runtimeStatusProjection.ts
  - internal/datalink/point/sql_repo_device_acquisition.go
  - internal/datalink/anomaly/detector_recovered.go
  - internal/datalink/sourcerule/mixed_identity_links.go
  - internal/datalink/collector/scheduler_dispatch.go
  - internal/datalink/sourcerule/acquisition_profile_source_identity.go
  - frontend/src/i18n/locales/zh-TW/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts
  - internal/datalink/runtime/mixed_source_layout.go
  - internal/datalink/dbtarget/writer_destination.go
  - internal/datalink/connector/acquisition_profile.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/router.go
  - internal/datalink/anomaly/detector_scan.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingAcquisitionSettings.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/workspace/write_group_source_format_identity.go
  - internal/datalink/point/repository_memory.go
  - internal/datalink/point/format_mutation_scope.go
  - internal/datalink/sourcerule/service_update.go
  - internal/datalink/schema/migrations/053_device_acquisition_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/pages/datalink/workbench-v2/rangeSourceReadback.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScatteredRangeSection.tsx
  - internal/protocol/mcprotocol/client_word_read.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/anomaly/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/recording-editor-constraints.json
  - docs/swagger/swagger.json
  - internal/api/handlers/acquisition_anomaly_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/recording-home/acquisitionAnomalies.css
  - internal/datalink/optimizer/acquisition_random.go
  - internal/datalink/point/service_materialized_batch.go
  - internal/datalink/point/format_recovery_snapshot.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/components/parameterHelpContent.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton.tsx
  - internal/datalink/collector/scheduler_acquisition_snapshot.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/schema/migrations/054_acquisition_anomalies_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - internal/datalink/optimizer/acquisition_plan.go
  - internal/datalink/connector/acquisition_profile_schema.go
  - internal/datalink/point/source_rule_retirement.go
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/sourcerule/mixed_review_decision_cas.go
  - internal/datalink/sourcerule/share_recovery_sql_reads.go
  - internal/datalink/dbtarget/writer_skip.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent.tsx
  - internal/datalink/point/service_device_acquisition_memory.go
  - internal/datalink/sourcerule/mixed_rule.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingCompactMemberTable.tsx
  - internal/datalink/point/sql_repo_materialized_batch.go
  - internal/datalink/connector/read_profile.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/sourcerule/tag_review_decision_stale.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/point/service_polling_group.go
  - internal/datalink/sourcerule/acquisition_source_exclusions.go
tests:
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - cmd/test_ui/target_writer_async_test.go
  - cmd/test_ui/acquisition_scheduled_safe_mixed_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-sections.test.tsx
  - internal/datalink/anomaly/detector_test.go
  - internal/api/handlers/mixed_address_profile_regression_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - internal/datalink/sourcerule/service_device_polling_group_test.go
  - frontend/tests/unit/runtime-dashboard/diagnostics-no-write-claim.test.tsx
  - internal/datalink/connector/adapters/mc3e_planned_peer_test.go
  - internal/api/handlers/acquisition_anomaly_retention_regression_test.go
  - internal/datalink/sourcerule/mixed_identity_retirement_runtime_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.tsx
  - frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx
  - internal/datalink/sourcerule/acquisition_qiqr_geometry_regression_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/datalink/migrator_device_acquisition_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.tsx
  - internal/datalink/connector/acquisition_profile_radix_review_test.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/workbench-v2/recording-instantaneous-method.test.tsx
  - internal/api/handlers/runtime_rule_status_safety_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/connector/adapters/mc3e_planned_fixture_test.go
  - internal/datalink/dbtarget/writer_ownership_async_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - internal/datalink/connector/acquisition_profile_random_limit_test.go
  - frontend/tests/unit/workbench-v2/recording-hidden-and-pending-areas.test.tsx
  - internal/datalink/pollinggroup/sql_repo_test.go
  - internal/datalink/dbtarget/writer_review_busy_test.go
  - frontend/tests/unit/workbench-v2/recording-confirmation-readable.test.tsx
  - internal/datalink/sourcerule/mixed_rule_test.go
  - internal/api/router_acquisition_anomaly_test.go
  - internal/datalink/optimizer/acquisition_random_long_timer_regression_test.go
  - internal/datalink/point/format_mutation_sql_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/sourcerule/mixed_review_decision_memory_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/api/handlers/acquisition_anomaly_handler_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/writer_review_refresh_test.go
  - internal/datalink/pollinggroup/sql_repo_owner_test.go
  - cmd/test_ui/acquisition_scheduled_mixed_test.go
  - internal/datalink/collector/scheduler_auto_interval_test.go
  - frontend/tests/unit/workbench-v2/scattered-metadata-roundtrip.test.ts
  - internal/datalink/optimizer/acquisition_random_contract_test.go
  - frontend/tests/unit/services/runtimeStatusSafety.test.ts
  - internal/datalink/runtime/status_auto_interval_test.go
  - frontend/tests/unit/scatteredRangeForm.test.tsx
  - internal/datalink/runtime/mixed_rule_format_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-unknown-codes.test.ts
  - internal/api/handlers/mixed_item_replacement_concurrency_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-destination-lock-notice.test.tsx
  - internal/datalink/anomaly/sql_repo_test.go
  - internal/api/handlers/mixed_item_replacement_regression_test.go
  - internal/protocol/mcprotocol/mcprotocol_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/device/service_acquisition_test.go
  - internal/datalink/optimizer/acquisition_review_regression_test.go
  - internal/datalink/point/service_polling_group_owner_test.go
  - internal/datalink/dbtarget/writer_backoff_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-editor.test.tsx
  - internal/datalink/sourcerule/acquisition_source_mixed_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-copy.test.tsx
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - frontend/tests/unit/services/studio-v2-locator-acquisition.test.ts
  - frontend/tests/unit/workbench-v2/recording-method-groups.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-plain-language.test.tsx
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/anomaly/detector_restart_regression_test.go
  - internal/protocol/mcprotocol/client_planned_exact_test.go
  - internal/datalink/pollinggroup/acquisition_mutation_test.go
  - frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts
  - internal/datalink/migrator_acquisition_anomalies_test.go
  - internal/datalink/dbtarget/writer_review_timeout_test.go
  - internal/datalink/workspace/mixed_rule_format_regression_test.go
  - internal/datalink/anomaly/testhelpers_test.go
  - internal/datalink/collector/scheduler_qiqr_shape_regression_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - internal/datalink/runtime/service_device_sync_test.go
  - internal/api/handlers/mixed_item_review_decision_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx
  - frontend/tests/unit/workbench-v2/recording-schema-instruction.test.tsx
  - internal/datalink/point/service_device_acquisition_concurrency_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - cmd/test_ui/target_writer.go
  - cmd/test_ui/bug_a_activation_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.html
  - internal/datalink/dbtarget/writer_recovery_isolation_test.go
  - internal/datalink/point/mixed_rule_format_regression_test.go
  - internal/datalink/connector/acquisition_profile_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - frontend/tests/unit/workbench-v2/recording-instant-role-label.test.tsx
  - internal/datalink/dbtarget/writer_starvation_test.go
  - internal/datalink/point/service_device_acquisition_atomic_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-source-refresh.test.tsx
  - internal/protocol/mcprotocol/client_validation_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.html
  - frontend/tests/unit/hooks/useAcquisitionAnomalies.test.tsx
  - frontend/tests/unit/scatteredRangeReviewFixes.test.tsx
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - frontend/tests/unit/services/acquisitionAnomalies.test.ts
  - internal/datalink/runtime/ingestor_database_decoupling_test.go
  - internal/datalink/device/sql_repo_acquisition_test.go
  - internal/api/handlers/dbtarget_handler_tooling_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/datalink/runtime/ingestor_async_acceptance_test.go
  - internal/protocol/mcprotocol/client_word_read_test.go
  - internal/datalink/anomaly/detector_observation_regression_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-method.test.tsx
  - internal/api/handlers/studio_v2_workspace_rule_switch_test.go
  - internal/datalink/runtime/anomaly_monitor_test.go
  - internal/datalink/workspace/runtime_acquisition_mixed_test.go
-->

---
### Requirement: Readiness issues are shown in plain language

The recording editor and the runtime dashboard SHALL present readiness issues with a localized explanation keyed by issue code, covering at least `schema-confirmation-required`, `write-group-apply-required`, `database-target-missing`, `source-revision-stale` and `schema-check-blocked`; unknown codes SHALL use a generic localized explanation. The raw code and server message SHALL appear only inside diagnostics, together with the localized explanation. The recording member search input SHALL disable browser autofill.

#### Scenario: Schema confirmation required

- **WHEN** a saved managed recording reports `schema-confirmation-required`
- **THEN** the issue reads as a localized instruction to check and confirm the table structure before starting, and the code with its English message appears only in diagnostics

<!-- @trace
source: clarify-recording-editor-constraints
updated: 2026-11-06
code:
  - internal/datalink/collector/scheduler_acquisition_preview.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection.tsx
  - internal/datalink/dbtarget/writer_idle.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields.tsx
  - internal/datalink/collector/scheduler_acquisition_dispatch.go
  - frontend/src/types/acquisitionAnomaly.ts
  - internal/datalink/collector/scheduler_auto_interval.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/anomaly/types.go
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice.tsx
  - internal/api/handlers/device_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/protocol/mcprotocol/random_read.go
  - internal/protocol/mcprotocol/client_planned_exact.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/point/service_device_polling_group.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/types/acquisitionSettings.ts
  - internal/datalink/device/acquisition_validation.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingAcquisitionSettings.css
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary.tsx
  - internal/datalink/dbtarget/writer_config_busy.go
  - frontend/src/i18n/locales/en/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/scatteredRange.ts
  - internal/datalink/sourcerule/point_runtime_compensation.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - internal/protocol/mcprotocol/const.go
  - internal/datalink/optimizer/acquisition_contiguous.go
  - internal/datalink/dbtarget/writer_live_check.go
  - internal/datalink/dbtarget/writer_retry_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel.tsx
  - frontend/src/i18n/locales/en/recording-home.json
  - internal/datalink/dbtarget/writer_workers.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/hooks/datalink/useAcquisitionAnomalies.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/i18n/locales/zh-TW/parameter-help-advanced.json
  - internal/datalink/sourcerule/sql_repo.go
  - internal/datalink/runtime/anomaly_monitor.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/runtime_rule_status_projection.go
  - internal/datalink/sourcerule/service_requests.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection.tsx
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups.tsx
  - internal/datalink/measurement/source_rule_format.go
  - internal/datalink/sourcerule/format_preflight.go
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/optimizer/acquisition_coverage_validation.go
  - internal/datalink/device/service.go
  - frontend/src/services/studioV2LocatorEntry.ts
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupHelpers.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSourceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/instantaneousReduction.ts
  - internal/api/router_acquisition_anomaly.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/pollinggroup/acquisition_storage.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingEditorLayout.css
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAcknowledgement.ts
  - frontend/src/pages/datalink/workbench-v2/deviceSourceSubmission.ts
  - internal/datalink/collector/scheduler_acquisition_publication.go
  - internal/datalink/collector/scheduler_acquisition_execute.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation.tsx
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/anomaly/detector.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/services/acquisitionAnomalies.ts
  - internal/datalink/connector/adapters/mc3e_planned_read.go
  - internal/datalink/device/repository_memory.go
  - internal/datalink/tag/workspace_save_rollback.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/service_read_profile.go
  - internal/datalink/optimizer/acquisition_mc_profile.go
  - internal/datalink/runtime/service_lifecycle.go
  - frontend/src/i18n/locales/en/parameter-help-advanced.json
  - internal/datalink/sourcerule/mixed_identity_retirement.go
  - internal/datalink/dbtarget/writer_refresh_tags.go
  - internal/datalink/point/memory_repo_materialized_batch.go
  - internal/datalink/mapping/source_rule_retirement.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/sourcerule/mutation_rollback_helpers.go
  - internal/datalink/optimizer/acquisition_candidates.go
  - frontend/src/i18n/locales/en/recording-editor-constraints.json
  - internal/datalink/connector/read_profile_identity.go
  - internal/protocol/mcprotocol/client.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingInstantaneousMethod.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingPointPicker.tsx
  - internal/datalink/optimizer/acquisition_builder.go
  - docs/swagger/docs.go
  - internal/datalink/pollinggroup/sql_repo.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/acquisition_anomaly_types.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer.ts
  - frontend/src/types/measurement.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
  - internal/datalink/dbtarget/tooling_types.go
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/tag_mapping_sync_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/sql_repo_acquisition.go
  - internal/datalink/connector/acquisition_profile_decode.go
  - internal/datalink/device/errors.go
  - internal/datalink/schema/schema_point_models.go
  - internal/datalink/sourcerule/mixed_identity_mapping.go
  - internal/datalink/sourcerule/repository_rule_delete_revision.go
  - internal/datalink/sourcerule/mc_profile_address_plan.go
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/sourcerule/point_update_plan.go
  - internal/datalink/tag/source_rule_retirement.go
  - internal/datalink/sourcerule/candidate_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/recordingRangeBuffer.ts
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/dbtarget/write_history.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_group_snapshot.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyDialog.tsx
  - frontend/src/pages/datalink/workbench-v2/useDeviceSourceSubmission.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/migrator_acquisition_anomalies.go
  - frontend/src/services/runtimeStatusProjection.ts
  - internal/datalink/point/sql_repo_device_acquisition.go
  - internal/datalink/anomaly/detector_recovered.go
  - internal/datalink/sourcerule/mixed_identity_links.go
  - internal/datalink/collector/scheduler_dispatch.go
  - internal/datalink/sourcerule/acquisition_profile_source_identity.go
  - frontend/src/i18n/locales/zh-TW/acquisition-settings.json
  - frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts
  - internal/datalink/runtime/mixed_source_layout.go
  - internal/datalink/dbtarget/writer_destination.go
  - internal/datalink/connector/acquisition_profile.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/router.go
  - internal/datalink/anomaly/detector_scan.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingAcquisitionSettings.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/workspace/write_group_source_format_identity.go
  - internal/datalink/point/repository_memory.go
  - internal/datalink/point/format_mutation_scope.go
  - internal/datalink/sourcerule/service_update.go
  - internal/datalink/schema/migrations/053_device_acquisition_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/pages/datalink/workbench-v2/rangeSourceReadback.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScatteredRangeSection.tsx
  - internal/protocol/mcprotocol/client_word_read.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/anomaly/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/recording-editor-constraints.json
  - docs/swagger/swagger.json
  - internal/api/handlers/acquisition_anomaly_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/recording-home/acquisitionAnomalies.css
  - internal/datalink/optimizer/acquisition_random.go
  - internal/datalink/point/service_materialized_batch.go
  - internal/datalink/point/format_recovery_snapshot.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent.tsx
  - frontend/src/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/components/parameterHelpContent.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/AcquisitionAnomalyButton.tsx
  - internal/datalink/collector/scheduler_acquisition_snapshot.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/schema/migrations/054_acquisition_anomalies_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - internal/datalink/optimizer/acquisition_plan.go
  - internal/datalink/connector/acquisition_profile_schema.go
  - internal/datalink/point/source_rule_retirement.go
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/sourcerule/mixed_review_decision_cas.go
  - internal/datalink/sourcerule/share_recovery_sql_reads.go
  - internal/datalink/dbtarget/writer_skip.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingRangeContent.tsx
  - internal/datalink/point/service_device_acquisition_memory.go
  - internal/datalink/sourcerule/mixed_rule.go
  - frontend/src/features/datalink/workbench-v2/recording-editor/RecordingCompactMemberTable.tsx
  - internal/datalink/point/sql_repo_materialized_batch.go
  - internal/datalink/connector/read_profile.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/sourcerule/tag_review_decision_stale.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/point/service_polling_group.go
  - internal/datalink/sourcerule/acquisition_source_exclusions.go
tests:
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - cmd/test_ui/target_writer_async_test.go
  - cmd/test_ui/acquisition_scheduled_safe_mixed_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-sections.test.tsx
  - internal/datalink/anomaly/detector_test.go
  - internal/api/handlers/mixed_address_profile_regression_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - internal/datalink/sourcerule/service_device_polling_group_test.go
  - frontend/tests/unit/runtime-dashboard/diagnostics-no-write-claim.test.tsx
  - internal/datalink/connector/adapters/mc3e_planned_peer_test.go
  - internal/api/handlers/acquisition_anomaly_retention_regression_test.go
  - internal/datalink/sourcerule/mixed_identity_retirement_runtime_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.tsx
  - frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx
  - internal/datalink/sourcerule/acquisition_qiqr_geometry_regression_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/datalink/migrator_device_acquisition_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.tsx
  - internal/datalink/connector/acquisition_profile_radix_review_test.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/workbench-v2/recording-instantaneous-method.test.tsx
  - internal/api/handlers/runtime_rule_status_safety_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/connector/adapters/mc3e_planned_fixture_test.go
  - internal/datalink/dbtarget/writer_ownership_async_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - internal/datalink/connector/acquisition_profile_random_limit_test.go
  - frontend/tests/unit/workbench-v2/recording-hidden-and-pending-areas.test.tsx
  - internal/datalink/pollinggroup/sql_repo_test.go
  - internal/datalink/dbtarget/writer_review_busy_test.go
  - frontend/tests/unit/workbench-v2/recording-confirmation-readable.test.tsx
  - internal/datalink/sourcerule/mixed_rule_test.go
  - internal/api/router_acquisition_anomaly_test.go
  - internal/datalink/optimizer/acquisition_random_long_timer_regression_test.go
  - internal/datalink/point/format_mutation_sql_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/sourcerule/mixed_review_decision_memory_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/api/handlers/acquisition_anomaly_handler_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/writer_review_refresh_test.go
  - internal/datalink/pollinggroup/sql_repo_owner_test.go
  - cmd/test_ui/acquisition_scheduled_mixed_test.go
  - internal/datalink/collector/scheduler_auto_interval_test.go
  - frontend/tests/unit/workbench-v2/scattered-metadata-roundtrip.test.ts
  - internal/datalink/optimizer/acquisition_random_contract_test.go
  - frontend/tests/unit/services/runtimeStatusSafety.test.ts
  - internal/datalink/runtime/status_auto_interval_test.go
  - frontend/tests/unit/scatteredRangeForm.test.tsx
  - internal/datalink/runtime/mixed_rule_format_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-unknown-codes.test.ts
  - internal/api/handlers/mixed_item_replacement_concurrency_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-destination-lock-notice.test.tsx
  - internal/datalink/anomaly/sql_repo_test.go
  - internal/api/handlers/mixed_item_replacement_regression_test.go
  - internal/protocol/mcprotocol/mcprotocol_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/device/service_acquisition_test.go
  - internal/datalink/optimizer/acquisition_review_regression_test.go
  - internal/datalink/point/service_polling_group_owner_test.go
  - internal/datalink/dbtarget/writer_backoff_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-editor.test.tsx
  - internal/datalink/sourcerule/acquisition_source_mixed_test.go
  - frontend/tests/unit/workbench-v2/recording-simple-copy.test.tsx
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - frontend/tests/unit/services/studio-v2-locator-acquisition.test.ts
  - frontend/tests/unit/workbench-v2/recording-method-groups.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/recording-readiness-plain-language.test.tsx
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/anomaly/detector_restart_regression_test.go
  - internal/protocol/mcprotocol/client_planned_exact_test.go
  - internal/datalink/pollinggroup/acquisition_mutation_test.go
  - frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts
  - internal/datalink/migrator_acquisition_anomalies_test.go
  - internal/datalink/dbtarget/writer_review_timeout_test.go
  - internal/datalink/workspace/mixed_rule_format_regression_test.go
  - internal/datalink/anomaly/testhelpers_test.go
  - internal/datalink/collector/scheduler_qiqr_shape_regression_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - internal/datalink/runtime/service_device_sync_test.go
  - internal/api/handlers/mixed_item_review_decision_regression_test.go
  - frontend/tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx
  - frontend/tests/unit/workbench-v2/recording-schema-instruction.test.tsx
  - internal/datalink/point/service_device_acquisition_concurrency_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - cmd/test_ui/target_writer.go
  - cmd/test_ui/bug_a_activation_test.go
  - frontend/tests/e2e/fixtures/recording-anomalies-host.html
  - internal/datalink/dbtarget/writer_recovery_isolation_test.go
  - internal/datalink/point/mixed_rule_format_regression_test.go
  - internal/datalink/connector/acquisition_profile_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - frontend/tests/unit/workbench-v2/recording-instant-role-label.test.tsx
  - internal/datalink/dbtarget/writer_starvation_test.go
  - internal/datalink/point/service_device_acquisition_atomic_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-source-refresh.test.tsx
  - internal/protocol/mcprotocol/client_validation_test.go
  - frontend/tests/e2e/fixtures/recording-constraints-host.html
  - frontend/tests/unit/hooks/useAcquisitionAnomalies.test.tsx
  - frontend/tests/unit/scatteredRangeReviewFixes.test.tsx
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - frontend/tests/unit/services/acquisitionAnomalies.test.ts
  - internal/datalink/runtime/ingestor_database_decoupling_test.go
  - internal/datalink/device/sql_repo_acquisition_test.go
  - internal/api/handlers/dbtarget_handler_tooling_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/datalink/runtime/ingestor_async_acceptance_test.go
  - internal/protocol/mcprotocol/client_word_read_test.go
  - internal/datalink/anomaly/detector_observation_regression_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - frontend/tests/unit/workbench-v2/recording-legacy-method.test.tsx
  - internal/api/handlers/studio_v2_workspace_rule_switch_test.go
  - internal/datalink/runtime/anomaly_monitor_test.go
  - internal/datalink/workspace/runtime_acquisition_mixed_test.go
-->