## ADDED Requirements

### Requirement: Basic recording controls hide incidental complexity
The basic database path SHALL expose destination, selected data, recording interval, completeness behavior and one start action. Revisions and advanced bindings SHALL remain internally enforced and available for diagnosis, not required operator inputs. Managed preparation/readiness MUST NOT depend on an unrelated legacy target or row-group collection.

#### Scenario: Ready managed group without legacy targets
- **WHEN** the canonical group and confirmed managed schema are ready while legacy targets are empty
- **THEN** the basic start flow is available and does not require duplicate manual column configuration.

#### Scenario: First bucket is not yet due
- **WHEN** recording has started but the first complete UTC bucket has not closed
- **THEN** the UI explains the interval and waiting state rather than fabricating a row, shortening the first bucket or declaring a write failure.

### Requirement: Recoverable scoped recording start
One explicit basic start intent SHALL coordinate existing save, readiness, group Apply and device activation services under the owning scope and revisions. It MUST require already confirmed schema preparation, retain progress in the existing operation ledger and reuse operation identity on retries. It MUST NOT silently execute DDL, bypass Share gates or claim atomic rollback across resources.

#### Scenario: Missing schema confirmation
- **WHEN** start finds a required structure without a valid confirmed preparation result
- **THEN** it directs the operator to preparation without creating the structure or reporting recording success.

#### Scenario: Duplicate start or response loss
- **WHEN** the same scoped start request is repeated, including after reload or process restart
- **THEN** the same operation is returned or resumed from recorded progress without duplicate groups, repeated effects or a second writer.

#### Scenario: Partial start and changed intent
- **WHEN** some selected resources have applied or started before a later failure, or the operator changes the intent
- **THEN** the product exposes actual per-scope progress, preserves healthy unrelated resources and revalidates changed intent before further actions. Recovery SHALL use the originally recorded device/group set and recorded revisions rather than the current whole workspace.

#### Scenario: Selected device scope is independent
- **GIVEN** selected device A is ready and unrelated B is incomplete or already healthy
- **WHEN** the operator starts A
- **THEN** A is checked and activated in its recorded scope; B is neither newly activated nor stopped or reconfigured, and incomplete B does not block A. Existing shared Share barriers remain enforced.

#### Scenario: Share-only setup
- **WHEN** only a valid Local Modbus output is selected
- **THEN** database preparation and database writers remain unused while all existing Share readiness and activation barriers remain enforced.

### Requirement: First recording result is evidence-backed
Completion SHALL distinguish live acquisition, local durable acceptance and confirmed SQL commit for the selected group and revision. A committed indication SHALL identify the corresponding recorded effect; verified SHALL remain reserved for actual readback. Running devices, successful test writes and buffered rows MUST NOT substitute for production SQL evidence.

#### Scenario: Destination offline after acquisition
- **WHEN** samples are durably accepted while destination writes remain queued
- **THEN** the UI shows acquisition and local acceptance with pending delivery, not SQL success.

#### Scenario: Production first row committed
- **WHEN** the production delivery path confirms its first SQL effect
- **THEN** the selected group view shows the actual committed stage and associated identity without claiming unperformed readback or cleanup.
