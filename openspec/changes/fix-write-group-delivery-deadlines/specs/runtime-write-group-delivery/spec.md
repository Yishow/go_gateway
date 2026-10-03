## MODIFIED Requirements

### Requirement: Bounded isolated recovery and capacity
Delivery SHALL bound concurrency, retry/backoff, batch sizes and shutdown time, with durable exclusive claims or lease/fencing. One destination failure MUST NOT block healthy independent destinations or independently valid Share output. Retry exhaustion SHALL retain blocked data. Quota exhaustion MUST stop new ACKs and expose affected scope and loss risk without deleting accepted pending records. Each destination attempt result SHALL receive a separate bounded local-settlement budget beginning after that attempt completes, without inheriting caller cancellation.

#### Scenario: Outage and restart
- **WHEN** database A is offline while database B and Share are healthy and the gateway restarts
- **THEN** A retains and resumes its backlog while healthy outputs continue within configured bounds.

#### Scenario: Disk full
- **WHEN** local durable storage reaches its configured hard limit
- **THEN** new intake is explicitly blocked without a success ACK; accepted journal/outbox evidence is retained.

#### Scenario: Poison row
- **WHEN** one row has a permanent SQL type or constraint failure
- **THEN** it is quarantined with a safe reason, its ordered group/entity successors wait for explicit resolution, and other partitions continue.

#### Scenario: Overlapping workers
- **WHEN** a stale worker and a new worker claim the same partition after restart
- **THEN** transactional ownership/fencing prevents conflicting progress and shutdown respects its deadline.

#### Scenario: Slow successful destination attempt
- **GIVEN** local settlement budget S and destination attempt budget D where S is less than D
- **WHEN** destination commit succeeds after S but before D and local storage is available
- **THEN** a fresh S budget is available to record sql_committed rather than inheriting an already expired settlement context.

#### Scenario: Slow failed destination attempt
- **WHEN** a destination attempt returns a transient or permanent error after the earlier S duration
- **THEN** retrying, blocked or quarantined classification is durably recorded using its fresh local budget, subject to the existing fencing rules.

#### Scenario: Caller cancelled after commit
- **WHEN** the caller is cancelled after a destination commit
- **THEN** local result recording remains bounded but is not cancelled by that caller; a genuine local write failure remains recoverable and is not falsely reported as sql_committed.
