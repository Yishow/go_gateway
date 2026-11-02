## MODIFIED Requirements

### Requirement: Persisted reduction and independent value roles

For existing instantaneous and cumulative profiles, the system SHALL persist recording_profile.sample_reduction as latest or mean and each instantaneous member's value_role as gauge, counter_raw or state. An omitted reduction SHALL decode as latest for existing groups and snapshots. New instantaneous recording SHALL default to mean with a 15-second interval and retain the existing integer 1–300-second bounds. Cumulative recording SHALL retain the existing integer 1–600-minute bounds and 10-minute default and SHALL reject mean. Gauge mean SHALL require a numeric source and explicit operator selection without requiring a fabricated measurement identity. Counter roles SHALL require the existing confirmed counter binding. Bool/text states SHALL retain latest eligible values and SHALL NOT enter an average. Names and units MUST NOT infer roles.

#### Scenario: Legacy reload and explicit edit
- **WHEN** an existing group without sample_reduction is loaded, or its interval is outside new instantaneous bounds
- **THEN** its saved latest semantics and interval remain unchanged; only an explicit compatible new revision enters mean recording.

#### Scenario: Invalid combinations
- **WHEN** mean is requested for a cumulative profile, bool/text gauge, unknown role or unconfirmed counter role
- **THEN** validation rejects the incompatible intent with an actionable typed error before Apply and no writer starts.

#### Scenario: One source in both profiles
- **WHEN** a confirmed energy counter is selected as counter_raw in an instantaneous snapshot and as cumulative in its separate profile
- **THEN** the instantaneous record stores its latest eligible raw value, the cumulative record uses its own existing interval and baseline, and neither averages the counter.

#### Scenario: An explicit mixed recording is requested
- **WHEN** a new versioned mixed profile assigns mean, mean and delta to individual members
- **THEN** its per-member calculation contract applies rather than the legacy group sample_reduction restriction, while unsupported combinations fail closed and all existing profiles retain their original semantics.
