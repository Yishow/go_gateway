## ADDED Requirements

### Requirement: Immutable counter proof for accepted samples

The system SHALL persist complete, versioned counter meaning in the immutable applied runtime snapshot before acknowledging cumulative or counter_raw samples, including unscaled members. Each accepted input SHALL be durably associated with that proof. The proof SHALL bind workspace, group and applied revision, source and member identity, measurement identity, definition and source-binding revisions, series epoch, kind, unit, source and result types, exact scaling and rollover/rate policy. Recovery, reconciliation and replay MUST NOT derive those facts from editable current definitions. This requirement supplements Revision-safe delivery and apply without replacing destination freezing or atomic acceptance.

#### Scenario: Definition changes before replay

- **WHEN** exact readings 100 and 105 were accepted under epoch-1, definition def-1 and unit kWh, the public measurement update changes the definition to epoch-2 and MWh without changing the old applied revision, and the gateway restarts
- **THEN** replay SHALL preserve epoch-1, def-1, kWh and delta 5 for the old accepted set, or quarantine it if its original proof cannot be established; it MUST NOT emit delta 5 MWh or attach epoch-2 to that set.

#### Scenario: Definition is unavailable after durable acceptance

- **WHEN** a confirmed definition is edited, removed or unreadable after samples with an immutable proof were acknowledged
- **THEN** replay SHALL use the original proof and original destination, preserve exact values and counter policy, and require no current-definition read.

##### Example: Frozen definition outage

- **GIVEN** sample c-1 was accepted with proof p-1 for def-1, epoch-1 and kWh
- **WHEN** the current definition repository returns unavailable during restart
- **THEN** c-1 SHALL replay with p-1 without reading a replacement definition.

#### Scenario: Concurrent definition edit and apply

- **WHEN** a definition changes while a new applied revision is freezing its counter proof
- **THEN** apply SHALL commit a complete revision-matched proof or fail without starting that revision; it MUST NOT acknowledge samples against a partial or mixed proof.

##### Example: Revision race

- **GIVEN** review references def-1 and a concurrent edit persists def-2
- **WHEN** apply attempts to save the proof after def-2 is committed
- **THEN** the stale apply SHALL fail without recording a mixed def-1/def-2 proof.

#### Scenario: Incoming claims name another epoch

- **WHEN** a new sample declares a measurement, epoch, definition revision or source-binding revision inconsistent with the active immutable proof
- **THEN** intake SHALL refuse that sample with a scoped actionable reason and SHALL NOT relabel it, advance the old baseline or silently substitute current semantics.

##### Example: Foreign epoch claim

- **GIVEN** the active proof names epoch-1 and the incoming sample names epoch-2
- **WHEN** the sample reaches acceptance
- **THEN** it SHALL be refused and the old checkpoint SHALL remain unchanged.

### Requirement: Evidence-preserving legacy counter replay

The system SHALL version the proof reference and journal identity format without rewriting historical payloads or their identity digests. Legacy unfinalized counter inputs SHALL replay only when immutable evidence proves the entire original binding for the accepted set. Current definitions, matching labels or incomplete checkpoints MUST NOT be treated as proof. Missing, conflicting or unsupported proof SHALL quarantine the affected group revision with safe reason, affected scope, pending count and repair guidance while retaining inputs and checkpoints. This requirement supplements Recovery and retention protect evidence; it grants no permission to delete or reinterpret historical data.

#### Scenario: Proven legacy exact source

- **WHEN** a legacy journal has no new proof-reference field but its immutable applied exact-source descriptor proves all required counter facts
- **THEN** migration SHALL attach a traceable proof association without modifying the original payload or digest, and replay SHALL preserve the original values and semantics.

##### Example: Old journal with complete scaled proof

- **GIVEN** legacy sample c-1 has original digest d-1 and a complete immutable def-1/epoch-1/kWh scale-0.1 descriptor
- **WHEN** the legacy proof association is migrated
- **THEN** c-1 SHALL retain digest d-1 and its original payload while replay uses that descriptor.

#### Scenario: Legacy proof cannot be established

- **WHEN** an unfinalized legacy journal has no complete immutable binding, contradictory evidence, or an unsupported proof version
- **THEN** recovery SHALL expose respectively counter-replay-proof-missing, counter-replay-proof-conflict or counter-replay-proof-unsupported, retain the affected durable data, block its calculation and disclose how to restore trusted revision evidence; it SHALL NOT use the current definition, reset the baseline or silently discard inputs.

##### Example: Absent original proof

- **GIVEN** an unfinalized legacy journal has no complete immutable definition evidence and the current editable definition is def-2
- **WHEN** recovery inspects the journal
- **THEN** it SHALL report counter-replay-proof-missing, retain every pending input and reject use of def-2 as a substitute.

#### Scenario: Quarantine does not rewrite finalized effects

- **WHEN** unfinalized inputs in one revision are quarantined while finalized effects already exist or another group is healthy
- **THEN** finalized effects SHALL retain their original bytes and destination and remain independently deliverable, while the unrelated healthy group continues under its own proof.

##### Example: Mixed backlog

- **GIVEN** group A has a quarantined unfinalized sample and finalized effect e-1 with digest d-1, while group B has valid proof
- **WHEN** delivery and intake continue
- **THEN** e-1 SHALL retain d-1 and its frozen destination, and B SHALL continue independently.

### Requirement: Replay proof preserves epoch and idempotency boundaries

Counter proof persistence, sample acceptance, checkpoint advancement and finalized output SHALL retain the existing durable transaction and idempotency guarantees. An explicitly applied new epoch SHALL start a new baseline; old accepted inputs MUST NOT be joined to that baseline. Existing scaled-source proof and supported exact encodings SHALL remain unchanged in meaning.

#### Scenario: New epoch after old backlog

- **WHEN** old epoch-1 samples remain pending and a confirmed new revision accepts its first epoch-2 reading
- **THEN** the old set SHALL drain using epoch-1 proof and its frozen destination, and the new reading SHALL have NULL delta and baseline status without consuming the old checkpoint.

#### Scenario: Crash and resend across codec versions

- **WHEN** the process fails before or after proof/sample commit or closure commit and an identical sample is resent under its original format
- **THEN** acknowledged inputs SHALL remain recoverable once, uncommitted inputs SHALL not be reported accepted, historical digest verification SHALL remain valid, and finalized payload retries SHALL not recompute or advance consumption.

##### Example: Legacy resend after new reader deployment

- **GIVEN** legacy sample c-1 with digest d-1 was durably acknowledged before deployment
- **WHEN** the new reader restarts and receives the identical c-1 resend
- **THEN** it SHALL recognize the existing input without rewriting d-1 or creating another consumed input.
