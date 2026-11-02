## Purpose

Provide a guided Step 4 presentation within Studio V2 that makes scope, rules, prerequisites, and the next safe action clear. Preserve complete advanced details and the existing canonical configuration and operation authorities.

## ADDED Requirements

### Requirement: One selected scope and next safe primary action

Step 4 SHALL guide source selection, recording mode, rule review, and preparation or start within one selected workspace/device/circuit scope. It SHALL derive the current primary action from revision-bound server evidence rather than persist an independent readiness model. Unresolved operations SHALL take priority over starting a new operation. Every disabled dependent action SHALL expose an actionable reason, while offline draft navigation and unrelated read-only views SHALL remain available.

#### Scenario: Ordinary typed recording without advanced semantics

- **WHEN** an operator configures valid non-counter instantaneous Tags
- **THEN** the guided path reaches schema review and scoped start without requiring reports, retention policies, physical measurement definitions, or manual advanced group editing.

##### Example: Ordinary typed recording without advanced semantics

- **GIVEN** device A has two saved non-counter gauge Tags
- **WHEN** the operator follows Step 4
- **THEN** the path reaches schema review and scoped start without opening Advanced or creating measurement definitions

#### Scenario: An operation is unresolved

- **WHEN** the current scope has an unresolved preparation, start, or test-write operation
- **THEN** the primary action leads to querying that operation instead of creating another one and the reason remains visible.

##### Example: An operation is unresolved

- **GIVEN** start operation O1 is unknown for device A
- **WHEN** Step 4 renders the next action
- **THEN** query O1 is primary and a new start is not offered as the primary action


### Requirement: Rule summaries distinguish saved and effective configuration

The guided summary SHALL expose selected sources, layout, profile kind, reduction, interval and unit, incomplete and lateness rules, destination table, and safe ownership context. Saved draft and applied rules SHALL be separately labeled when different. Detailed views SHALL retain source bindings, types, columns, counter policies and epochs, revisions, schema evidence, operations, and delivery stages. Saved, applied, accepted, SQL committed, readback verified, and cleanup completed MUST NOT be represented as interchangeable success states.

#### Scenario: Saved period differs from the running period

- **WHEN** a group has a new saved interval but an older applied revision
- **THEN** the summary labels both periods and does not claim that the new interval is currently executing.

##### Example: Saved period differs from the running period

- **GIVEN** saved draft R2 is 30 seconds and applied R1 is 15 seconds
- **WHEN** the rule summary renders
- **THEN** 30 seconds is labeled saved and 15 seconds is labeled effective

#### Scenario: Inspect complete technical details

- **WHEN** an operator opens the summary's detailed rules or delivery evidence
- **THEN** the details retain the same selected scope and expose the existing policy and evidence without requiring a second configuration model.

##### Example: Inspect complete technical details

- **GIVEN** G1 in workspace W1 has saved source binding B1 and applied revision R1
- **WHEN** its detail drawer is opened
- **THEN** B1 and R1 remain accessible under G1/W1 and closing restores focus to the opener


### Requirement: Recording management is discoverable without a second activation path

Saved recording cards SHALL expose server-authorized edit, stop-intake, remove-configuration, transition, and recreate actions with their blockers. Advanced content SHALL be labeled as advanced settings and workspace status, not workspace activation. Managed profiles SHALL retain one reviewed scoped-start path. The UI SHALL NOT add a workspace-wide start action or convert legacy/custom groups by navigation.

#### Scenario: Stop or remove a saved recording

- **WHEN** an operator opens a saved wide or long recording card
- **THEN** the appropriate lifecycle actions and their history or pending-delivery impact are visible without requiring the operator to infer what a write group is.

##### Example: Stop or remove a saved recording

- **GIVEN** saved wide recording G1 has queued row Q1
- **WHEN** its recording card is opened
- **THEN** stop/remove actions disclose that Q1 keeps its original destination and SQL history is retained

#### Scenario: Open advanced workspace status

- **WHEN** an operator expands advanced settings
- **THEN** the workspace overview remains a status view, no workspace-wide start occurs, and an existing managed profile's effective action returns to its own reviewed scope.

##### Example: Open advanced workspace status

- **GIVEN** device A is selected and device B also exists
- **WHEN** Advanced is expanded
- **THEN** no activation request is sent and the managed-profile action remains scoped to A


### Requirement: Navigation preserves drafts and safety boundaries

Changing Basic/Advanced presentation, expanding details, changing viewed devices, using Back or Cancel, or reloading SHALL NOT implicitly save, apply, start, convert, or discard unresolved operations. The selected scope SHALL NOT silently switch to another device because its source list is empty. Keyboard interaction, focusable repair links, accessible status messages, and usable source-list scrolling SHALL be preserved across narrow and desktop layouts.

#### Scenario: Dismiss a review and switch devices

- **WHEN** an operator cancels an unsubmitted review and selects another device
- **THEN** no mutation is sent, the original draft remains associated with its scope, and returning restores its pending evidence and selection.

##### Example: Dismiss a review and switch devices

- **GIVEN** A has an unsaved 30-second draft and an unsubmitted review
- **WHEN** Cancel is selected and the view changes to B then A
- **THEN** A retains its draft and no save/apply/start request was sent

#### Scenario: Partial start success

- **WHEN** only a subset of the reviewed device scope is confirmed activated
- **THEN** the summary identifies the confirmed subset and the runtime handoff contains only those activated device IDs rather than all workspace devices.

##### Example: Partial start success

- **GIVEN** reviewed devices are A and B but only A is confirmed activated
- **WHEN** the operator continues to runtime
- **THEN** the handoff contains A only and B is visibly incomplete


### Requirement: Canonical setup selection revalidates ownership

The setup selection receiver SHALL accept a typed `StudioSelectionRequest` containing workspace and step plus applicable device, rule, mapping, and group identifiers. An optional group revision hint SHALL be treated as historical context and SHALL NOT authorize mutations. The receiver SHALL revalidate current workspace ownership and backend revisions before selecting an editable target. Selection SHALL NOT trigger saves, applies, or starts. URL parsing and diagnostic resolver links SHALL remain the responsibility of the subsequent runtime-rules capability and SHALL invoke this URL-independent interface.

#### Scenario: Selection opens an owned saved group

- **WHEN** a selection request names an accessible group and its workspace/device scope
- **THEN** Studio selects that group's current saved state, distinguishes an older revision hint if present, and waits for explicit operator action.

##### Example: Selection opens an owned saved group

- **GIVEN** StudioSelectionRequest names W1/A/G1 with historical revision hint R1 while G1 is now R2
- **WHEN** the receiver validates the request
- **THEN** G1/R2 is selected, R1 is only a historical hint, and mutation waits for an explicit action

#### Scenario: Selection is stale or foreign

- **WHEN** a selection request contains a deleted, missing, conflicting, or foreign identity
- **THEN** Studio displays a safe unavailable or historical-only state without selecting an unrelated group, exposing foreign data, or performing a mutation.

##### Example: Selection is stale or foreign

- **GIVEN** the current workspace is W1 and a selection request names W2/G2
- **WHEN** the receiver validates ownership
- **THEN** a safe unavailable state appears without foreign data or mutation

### Requirement: Browser verification exceptions remain explicit

Under the user-authorized exception for the existing cloud environment, delivery SHALL retain formal browser test cases, test discovery, E2E-specific typechecking, lint, and scenario coverage. Only actual browser execution and screenshots SHALL be eligible for a conditional delivery-gate exception when the verified environment blockers remain unchanged and execution remains unavailable. The exception SHALL NOT weaken product behavior, dependency order, ownership, backend database or frozen-proof checks, other tests or reviews, or platform and field verification boundaries. It SHALL NOT authorize a new environment or a security workaround.

#### Scenario: Blocked browser execution is reported without a false pass

- **WHEN** unchanged verified cloud blockers prevent actual browser execution or screenshot capture during this change's acceptance
- **THEN** the delivery record SHALL identify the affected cases and screenshots as NOT RUN, SKIP, or unverified, retain the blocker evidence and rerun commands, and report preparation checks separately from browser execution
- **AND** the delivery SHALL NOT report browser or screenshot PASS from discovery, typechecking, lint, simulated DOM results, or synthetic images

##### Example: Discovered case without a browser result

- **GIVEN** browser case B1 is discoverable and its E2E-specific typecheck and lint pass, but the verified launch blocker remains
- **WHEN** acceptance evidence is recorded under this authorized exception
- **THEN** B1 and its screenshot remain NOT RUN with the blocker reason; only the executed preparation checks are PASS
