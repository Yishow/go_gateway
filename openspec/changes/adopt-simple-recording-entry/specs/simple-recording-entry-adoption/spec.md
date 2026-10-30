## Purpose

Adopt the simple recording home as the default Studio entry only after its complete agreed journeys are verified and reviewed, while preserving existing explicit routes and safe fallback behavior.

## ADDED Requirements

### Requirement: Default adoption follows explicit acceptance gates
Product implementation of these dependent changes SHALL begin only after the owner has confirmed the completed demo checkpoint for this batch; demo confirmation SHALL NOT count as formal browser or device verification. The system SHALL retain the opt-in entry until the home, editor and lifecycle changes meet their agreed acceptance conditions and the owner has reviewed the demo checkpoint. Required verification gaps SHALL block default adoption unless their limitations are explicitly accepted for this batch. Silence or a historical waiver SHALL NOT establish acceptance.

#### Scenario: A required browser journey remains unverified
- **WHEN** current browser verification cannot run and no current acceptance of that limitation exists
- **THEN** the default remains the original entry, the new view remains opt-in, and the limitation is reported as NOT RUN rather than passed.

#### Scenario: Complete entry is accepted
- **WHEN** the dependent changes satisfy the acceptance matrix and the owner checkpoint is recorded
- **THEN** an unqualified /studio/v2 entry opens the recording home while view=advanced retains the original guided presentation.

### Requirement: Explicit existing intents survive default adoption
The system SHALL preserve valid entity locators, runtime repair focus and supported legacy step, target, section and modal intents before applying the default view. Requests with entry or canonical locator identity fields SHALL undergo strict locator validation and SHALL NOT fall back to legacy routing on error. Only requests without locator identity fields SHALL use bounded existing legacy intent allowlists after raw query encoding and uniqueness checks. Root and unknown routes SHALL retain their generic Studio V2 redirection and search/hash preservation. The retired /studio surface SHALL NOT be restored.

#### Scenario: A legacy Modbus setup link is opened
- **WHEN** an existing supported step=output and target=modbus intent reaches Studio V2
- **THEN** the matching existing setup is opened instead of the general recording home, without changing workspace settings.

#### Scenario: Returning from a runtime repair
- **WHEN** a valid scoped repair intent is opened after default adoption
- **THEN** the correct configuration context and return focus remain available without redirecting the operator through a fresh setup flow.

#### Scenario: A malformed locator is mixed with a legacy intent
- **WHEN** a request combines invalid locator identity or entry fields with step=output and target=modbus
- **THEN** it remains an invalid locator rather than being accepted as a legacy Modbus request.

### Requirement: Acceptance records distinguish evidence layers
The system's release evidence SHALL identify each agreed journey, actual verification command or manual assertion, final source version and PASS, FAIL or NOT RUN status. Demo, mocked UI, DOM tests, browser tests, SQL tests and field-device acceptance SHALL be identified separately. Results from older versions SHALL NOT count as current verification.

#### Scenario: Mocked UI journeys pass
- **WHEN** UI fixtures and DOM tests pass without real SQL or device execution
- **THEN** the release report states the UI coverage and does not claim real database writes or PLC acceptance.

#### Scenario: Final source changes after a full run
- **WHEN** a subsequent change affects a previously checked route, build or behavior
- **THEN** the affected checks are rerun against the final batch and the report does not present the older result as final evidence.

### Requirement: Entry rollback preserves configuration and operations
The system SHALL allow reverting the entry presentation to the original flow without undoing stored recording settings, deleting history or resubmitting operations. Advanced configuration and runtime, logs and test routes SHALL remain reachable throughout adoption.

#### Scenario: Entry adoption is reverted
- **WHEN** a regression requires the original entry presentation
- **THEN** only the UI entry selection changes and existing recording identities, pending operations and historical database data remain intact.

### Requirement: Revised design and protocol parity precede acceptance
Acceptance SHALL cover each existing supported protocol's required and advanced parameter paths, lossless reload/edit/save, both themes, system live changes, override persistence, portals, typography, contrast, reduced motion and responsive modal behavior. The owner checkpoint SHALL use the revised fixed-version demo; prior viewing SHALL NOT imply approval. Unselected optional features or modal size exceptions MUST NOT enter the implementation checklist.

#### Scenario: The revised point-picker journey is reviewed
- **WHEN** the revised fixed-version demo and eventual product journey are evaluated
- **THEN** evidence separately covers one additive picker, synchronized filters and counts, cancellation preserving member configuration and invalid input, main-page removal and restoration, explicit source targets, same-modal setup, distinct device/range persistence and saved-source survival after cancellation; mock success alone does not authorize product apply or establish browser verification.

#### Scenario: A design tool or browser gate cannot run
- **WHEN** Taste scripts, screenshots or required browser journeys have not run
- **THEN** the record states NOT RUN and distinguishes adopted text guidance, static checks and owner feedback from actual visual or hardware acceptance, without silently enabling the new default.

### Requirement: Mixed member acceptance preserves one-row and legacy contracts
Adoption SHALL depend on the completed per-member calculation prerequisite and editor consumer, with evidence for one managed-wide record using mean, mean and delta at their own aligned intervals with actual sparse-wide SQL rows and due-only long rows. First-baseline, required skip, optional NULL, gap spans, method proof invalidation, replay and unchanged legacy output SHALL be covered. Unsupported shapes MUST NOT be claimed as supported or replaced by silently split records.

#### Scenario: Only the demo shows mixed results
- **WHEN** a demo displays 22,12,8 but canonical runtime or SQL readback evidence is absent
- **THEN** the feature remains simulated and adoption does not claim production completion or bypass the revised owner checkpoint.

#### Scenario: Members have different due boundaries
- **WHEN** a record combines 15-second means and a 60-second delta
- **THEN** the UI and acceptance evidence show wide NULL with not_due at :15/:30/:45, both due calculations with their own [start,end) at :00, and long rows only for due members, without carrying delta or inserting zero; required-not-due does not skip a row and due failures retain distinct reasons and coverage.

### Requirement: Adoption distinguishes planned and actual acquisition counts
The consumer SHALL depend on integrate-protocol-aware-acquisition-batching for canonical scoped revision-bound acquisition-plan previews and actual observations, not derive its own batching plan. It SHALL label planned protocol read requests per scan excluding fault retries as estimates, distinguish acquisition Point count and unique wire units, and retain capability, coverage and fallback reasons. Actual attempts, successes, failures, retries, fallback and cancellations SHALL be separately labelled with matching scope/revision; absent observations MUST NOT become zero or measured speedup. Preview SHALL perform no PLC I/O, reject stale revisions, and MUST NOT be called TCP packet counts or SQL writes. This dependency SHALL NOT redefine member recording windows or make recording selection the acquisition set.

#### Scenario: Estimate exists before actual acquisition evidence
- **WHEN** a canonical preview reports three planned protocol requests but no matching actual scan observations exist
- **THEN** the UI labels three as an estimate, reports actual evidence unavailable rather than zero, and a later stale-revision response invalidates that preview instead of retaining a false current count.
