## ADDED Requirements

### Requirement: Cutover ownership is enforced at durable acceptance and closure

For every currently supported forward cutover and rollback, the system SHALL use persisted cutover intent and revision ownership as the authority shared by cutover commit, sample acceptance and window closure. A stale reconciliation cache MUST NOT authorize an incompatible durable write. The old owner SHALL own observations strictly before the immutable boundary; the new owner SHALL own observations at or after it within the reviewed source, destination-ownership and profile scope. A successful cutover SHALL establish this fence before success is reported. This requirement supplements Explicit forward cutover preserves wide history and does not expand its supported conversion matrix.

#### Scenario: Successful cutover before delayed reconciliation

- **WHEN** a reviewed cutover commits before its scheduled effective boundary at relative time 22.5 seconds, a value 6 has observed_at 22.499 seconds and a value 7 has observed_at 23 seconds, and runtime reconciliation is delayed
- **THEN** the old owner's truncated window SHALL contain 6 rather than 7, and 7 SHALL be durably assigned only to the new owner or retained in its durable pending handoff; absence of a ready worker SHALL NOT make the old owner eligible.

#### Scenario: Observation exactly at the boundary

- **WHEN** observations arrive before, exactly at and after the cutover boundary in any arrival order
- **THEN** only the observation strictly before the boundary SHALL belong to the old owner, and boundary-or-later observations SHALL belong only to the new owner within each affected ownership scope.

##### Example: Half-open ownership fixture

- **GIVEN** the scheduled boundary is 22.5 seconds and observations are at 22.499, 22.5 and 23 seconds
- **WHEN** they arrive in reverse order before their closure deadlines
- **THEN** 22.499 SHALL belong to the old owner and both 22.5 and 23 SHALL belong only to the new owner.

#### Scenario: Cutover races with sample or closure commit

- **WHEN** cutover commit overlaps acceptance or closure using a stale owner-fence revision
- **THEN** transactions SHALL serialize on the persisted ownership authority or reject and retry the stale work without advancing consumed input or checkpoints; a successful stale-cache write SHALL NOT be possible.

##### Example: Stale fence fixture

- **GIVEN** an acceptance plan holds fence f-1 and cutover commits fence f-2 before that plan commits
- **WHEN** the stale acceptance transaction attempts its write
- **THEN** it SHALL retry or refuse under f-2 without committing an f-1-owned cross-boundary sample.

#### Scenario: Candidate boundary overlaps finalized history

- **WHEN** already finalized effects conflict with the proposed ownership boundary
- **THEN** cutover SHALL fail with recording-cutover-finalized-overlap and a fresh-review action, preserve the previous owner and existing effects, and SHALL NOT rewrite finalized rows to make the cutover succeed.

##### Example: Finalized ownership conflict fixture

- **GIVEN** effect e-1 already assigns an observation at second 23 to the old owner
- **WHEN** a reviewed boundary at second 22.5 is proposed
- **THEN** the candidate SHALL be rejected with recording-cutover-finalized-overlap and e-1 SHALL remain unchanged.

### Requirement: Buffered cutover inputs are transferred without loss or reinterpretation

The system SHALL re-evaluate every affected unfinalized accepted observation against actual half-open owner bounds, including observations accepted before the cutover was committed. Moving an observation SHALL atomically persist its target acceptance or durable pending handoff, source transfer disposition and stable transfer identity while retaining original sample identity, payload, digest and measurement proof. The old owner MUST NOT consume a transferred observation, and the target MUST NOT consume it twice. Missing ownership or incompatible semantic proof SHALL block the affected handoff with safe scope and repair guidance, without deleting or relabelling accepted data.

#### Scenario: Future observation was buffered before cutover

- **WHEN** the old long-profile buffer contains value 100 before second 12.5 and value 200 observed at second 13 before learning that its end is second 12.5
- **THEN** its finalized row SHALL retain 100 with an observation time inside its actual window, and 200 SHALL remain durably owned by the successor or pending handoff rather than disappearing or appearing in both owners.

#### Scenario: Crash during durable transfer

- **WHEN** the process terminates before or after the transfer transaction and recovery retries the same accepted sample
- **THEN** the source disposition and target acceptance SHALL recover consistently under the same transfer identity, with exactly one consuming owner and no change to sample values, proof, original digest or destination ownership history.

##### Example: Transfer retry fixture

- **GIVEN** transfer t-1 moves sample s-1 from owner A to owner B
- **WHEN** the process exits after the transfer commit but before reporting success and t-1 is retried
- **THEN** the transfer SHALL still identify one s-1 consumer B with no second target acceptance.

#### Scenario: Handoff cannot prove compatible semantics

- **WHEN** an accepted sample's immutable counter or source proof conflicts with the proposed successor proof, or ownership evidence cannot be read
- **THEN** the affected handoff SHALL remain durable and blocked with recording-handoff-proof-conflict or recording-ownership-unavailable; the system SHALL NOT use current definitions, silently drop the sample or acknowledge an unpersisted transfer.

##### Example: Epoch conflict fixture

- **GIVEN** sample s-1 carries epoch-1/kWh proof while the successor only proves epoch-2/MWh
- **WHEN** handoff attempts to assign s-1
- **THEN** it SHALL remain blocked with recording-handoff-proof-conflict and SHALL NOT become an epoch-2/MWh sample.

### Requirement: Actual cutover bounds govern all supported reductions and recovery

Within the existing supported lifecycle matrix, actual owner bounds SHALL govern legacy wide/custom selection, managed-wide-profile and managed-long latest, mean, state/counter_raw and cumulative input sets. Rebuilding a truncated window SHALL retain earlier eligible candidates, recompute mean sum/count from its eligible set and exclude successor observations from old counter baselines. Closure SHALL wait for durably accepted pending handoffs affecting that window, preserve original acceptance-time eligibility and validate the current ownership fence. Boundary coverage SHALL remain partial, a new cumulative owner SHALL start with NULL delta and baseline status, and finalized outputs SHALL drain unchanged to frozen destinations.

#### Scenario: Latest mean and cumulative use the same owner partition

- **WHEN** a supported cutover truncates a window containing eligible pre-boundary and post-boundary samples
- **THEN** latest SHALL select only a pre-boundary candidate for the old owner, mean SHALL use only its unique pre-boundary samples in sum and count, and cumulative SHALL not compare an old checkpoint to a successor-owned observation; the successor's first cumulative result SHALL be baseline with NULL delta.

##### Example: Reduction partition fixture

- **GIVEN** a second-10 boundary divides values 10 at second 8, 20 at second 9 and 100 at second 11
- **WHEN** a supported latest or mean profile closes its old window
- **THEN** old latest SHALL equal 20 and old mean SHALL equal 15 with count 2; value 100 SHALL belong only to the successor.

#### Scenario: Restart before delayed handoff completion

- **WHEN** a sample was durably accepted before its closure deadline but handoff processing or reconciliation is delayed past that deadline and the process restarts
- **THEN** recovery SHALL preserve the accepted sample for the correct owner, prevent premature closure while its durable handoff is unresolved, and produce the same eligible result as uninterrupted execution without treating that accepted sample as newly late.

##### Example: Accepted before deadline fixture

- **GIVEN** sample s-1 belongs to successor window [10,15), was durably accepted by that successor at second 14 and has closure deadline second 16
- **WHEN** ordinary recovery and reconciliation run after a physical restart at second 20
- **THEN** s-1 SHALL remain eligible and recovery SHALL produce its original value without treating it as newly late.

##### Example: Delayed handoff preserves original acceptance time

- **GIVEN** source owner A durably accepts sample s-1 with value 200 and observed_at second 11 at second 9, and a cutover committed at second 9.5 establishes future boundary second 10 and successor window [10,15) with closure deadline second 16
- **WHEN** handoff remains unresolved at second 16 and recovery completes the handoff after a physical restart at second 20
- **THEN** closure of the affected window SHALL remain blocked while the handoff is unresolved; recovery SHALL preserve the original acceptance time second 9, keep s-1 eligible and produce value 200 for the successor without treating the delayed handoff as a newly late input.

#### Scenario: Truly late input and finalized retry

- **WHEN** an input is first received after the existing closure deadline, or an already finalized old effect is retried after cutover
- **THEN** the late input SHALL retain the existing late-after-close diagnosis without reopening a finalized window, and the finalized effect SHALL retain its original payload, effect identity and frozen destination without recomputation.

##### Example: First arrival after deadline fixture

- **GIVEN** window [10,15) finalized with deadline second 16 and effect e-1
- **WHEN** a new sample for that window first arrives at second 17 and e-1 is retried
- **THEN** the new sample SHALL be diagnosed late and e-1 SHALL retain its payload and effect identity.

#### Scenario: Rollback and independent profiles

- **WHEN** an existing supported rollback establishes a new boundary while instantaneous and cumulative profiles coexist
- **THEN** each affected profile scope SHALL have one owner at each observation instant, valid independent profile outputs SHALL remain distinct, and rollback SHALL preserve already accepted history with the same fence and transfer guarantees.

##### Example: Independent profile fixture

- **GIVEN** source s-1 has independently enabled instantaneous profile I and cumulative profile C
- **WHEN** a supported rollback changes ownership at second 40
- **THEN** each I/C scope SHALL have one owner at second 40, while I and C SHALL retain distinct output identities.

##### Example: Compatible same-origin fanout converges on rollback

- **GIVEN** the existing supported public setup has managed-wide instantaneous `counter_raw` owner A/R1 and independent I/C target profiles on distinct source/target canonical circuits and real SQLite connectors, all proving the same immutable measurement, epoch and unit; A/R1 accepts sample s-25 at second 8 with observed_at second 25, and a forward cutover committed at second 9 takes effect at second 10, transferring s-25 to both I and C with the same original A/R1 identity
- **WHEN** a supported rollback committed at second 19 takes effect at second 20 and both independently verified branches transfer that original observation into restored owner A/R3
- **THEN** before ordinary reclamation there SHALL be four journal rows (A/R1, I, C and A/R3) and four transfer receipts, with both rollback edges referencing one canonical A/R3 journal sequence and exactly one eventual A/R3 consumer; original payload bytes, digest, measurement proof and acceptance time SHALL remain unchanged, while independent profile outputs outside this converged ownership scope SHALL remain distinct
- **AND** all converging peers SHALL agree on full origin scope/sequence, original sample identity, bytes, target bucket, immutable proof and runtime digests; their source scope, fence, source bucket and target sets SHALL remain independently verified branch facts, and equal SampleIDs from different origins SHALL NOT merge
- **AND** the live target row's immutable handoff identity SHALL select the canonical audit, not SQL result order; alias edges SHALL retain an immutable canonical-transfer reference, legacy receipts SHALL remain compatible, and recovery or recognition SHALL reject a missing or mismatched canonical audit or corrupt peer, including after ordinary reclamation removes shared raw bytes and only compact meaning and terminal-edge evidence remain
- **AND** failure while accepting the second edge SHALL preserve accepted data and required pins, block affected closure and permit a corrected retry to complete without duplicate consumption or effects; physical reopen, replay and normal reclamation SHALL NOT change that result or fabricate consumed/finalized state to bypass validation.
