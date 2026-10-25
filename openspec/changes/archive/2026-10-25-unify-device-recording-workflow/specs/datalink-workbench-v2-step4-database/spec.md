## MODIFIED Requirements

### Requirement: First activation targets all eligible devices

The primary Step 4 database action SHALL use the explicitly reviewed selected device/group scope and SHALL NOT infer the whole singleton workspace from a device filter or status overview. It SHALL coordinate confirmed preparation, group Apply and scoped activation using the existing recording-start ledger. The workspace overview SHALL NOT present a second activation action. The existing global activation API SHALL remain compatible for explicit legacy callers and SHALL continue to activate all eligible not-yet-running devices only when that global operation is explicitly invoked.

#### Scenario: Activate all eligible devices

- **GIVEN** the singleton workspace contains multiple devices
- **AND** some devices are valid, available, and not yet running
- **WHEN** an explicit legacy caller invokes the existing global activation API
- **THEN** the system activates every eligible device in that set
- **AND** already running devices are not activated a second time

#### Scenario: Selected device is the only primary action scope
- **WHEN** A is selected and ready while B is incomplete or already collecting
- **THEN** the primary start validates/applies/activates only its frozen A/group scope, B is not started or stopped, unrelated B does not become a database readiness blocker, and existing shared Share barriers remain enforced.

### Requirement: Basic recording controls hide incidental complexity
The basic database path SHALL expose destination, device tiles, purpose/layout, selected data, native recording interval/reduction and one scoped start action. Whole-row completeness controls SHALL appear only for applicable wide/custom layouts; long recording SHALL show per-source missing evidence. The full member list SHALL be collapsed by default with visible exception summaries. Revisions and advanced bindings SHALL remain internally enforced and available for diagnosis, not required operator inputs. Managed preparation/readiness MUST NOT depend on an unrelated legacy target or row-group collection.

#### Scenario: Ready managed group without legacy targets
- **WHEN** the canonical group and confirmed managed schema are ready while legacy targets are empty
- **THEN** the basic start flow is available and does not require duplicate manual column configuration.

#### Scenario: First bucket is not yet due
- **WHEN** recording has started but the first complete UTC bucket has not closed
- **THEN** the UI explains the interval and waiting state rather than fabricating a row, shortening the first bucket or declaring a write failure.

### Requirement: Recoverable scoped recording start
One explicit primary start intent shared by Basic and Advanced SHALL coordinate existing save, readiness, group Apply and device activation services under the owning scope and revisions for either managed layout. Advanced SHALL save the same draft and route supported managed recording effectiveness through this shared scoped intent rather than provide an independent workspace activation. Existing cross-device custom Apply/lifecycle SHALL remain an engineering path outside automatic single-device orchestration and SHALL NOT be silently split or included. Existing lifecycle APIs, explicit stop/delete and legacy global activation callers SHALL remain compatible. It MUST require already confirmed schema preparation, retain progress in the existing operation ledger and reuse operation identity on retries. It MUST NOT silently execute DDL, bypass Share gates or claim atomic rollback across resources.

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

#### Scenario: Uncertain start while switching device
- **WHEN** A starts with a frozen request and its response is lost before the operator views B
- **THEN** A retains its original recoverable intent and progress, B does not inherit A's success/error or retry controls, and checking or retrying A reuses its original scope and request identity.
