## ADDED Requirements

### Requirement: Compact home source observations preserve canonical evidence

The recording home SHALL use a workspace-scoped read-only compact source observation derived from the existing canonical source authorities. The endpoint MUST NOT create, repair, reconcile or activate configuration. It SHALL compare complete source projections twice and return available only when identity, membership and content checks agree. Its digest MUST NOT serve as a command or readiness token or imply an atomic snapshot.

#### Scenario: Complete unchanged sources

- **WHEN** GET /api/v1/datalink/studio-v2/recording-home/sources receives an existing workspace_id and both complete projections agree
- **THEN** it returns that workspace identity, an opaque source digest and the same device/candidate/blocker projection as the existing home logic
- **AND** it excludes credentials, connection configuration, live point values and raw backend errors

#### Scenario: Inconsistent or unavailable sources

- **WHEN** source content changes between passes, duplicate identities occur, members are missing or foreign, or a source read fails
- **THEN** the endpoint returns unknown or an actionable read error instead of a partial available summary
- **AND** navigation performs no mutation

##### Example: mismatched mapping owner

- **GIVEN** workspace W1 owns device D1 while mapping M1 claims point P1 on foreign device D2
- **WHEN** W1 source observations are requested
- **THEN** M1 does not become an available candidate and the observation is unknown

### Requirement: Return navigation reuses observed presentation without reusing authority

The recording home SHALL preserve its same-workspace filter, page, scroll and focus while revalidating a cached observation on return. Cached evidence SHALL be labeled observed or stale until the current visit validates it. Every visit SHALL retain a separate cancellation and late-response fence. Cache entries SHALL be isolated by QueryClient/backend context and current API origin, retain only the latest observation per workspace, and preserve the existing eight-workspace bound and oldest-scope eviction.

#### Scenario: Same workspace returns

- **WHEN** a user returns to a previously observed workspace
- **THEN** the cached presentation is available with a refreshing indicator and one new source-summary request
- **AND** pending or failed validation cannot label the sources fresh or enable commands requiring fresh proof

#### Scenario: Scope changes and old replies arrive

- **WHEN** navigation follows workspace A, workspace B, then workspace A and an earlier request finishes late
- **THEN** only the latest matching visit updates displayed evidence
- **AND** B evidence and older A evidence cannot overwrite the current A observation

#### Scenario: Refresh fails or sources mutate

- **WHEN** revalidation fails or an affected source mutation invalidates the observation
- **THEN** the last observed view is explicitly stale or failed, with read-only retry
- **AND** fresh locator, CAS, revision and readiness checks remain required by existing action owners

### Requirement: Home data demand stays bounded and complete

One home source revalidation SHALL use one compact summary request rather than the previous ten raw source requests. Existing bounded visible-page detail demand SHALL remain bounded, and the summary MUST NOT infer completeness from a truncated source list.

#### Scenario: Large source fixture is projected

- **WHEN** the same synthetic 20-device and 2000-point fixture is read before and after migration
- **THEN** compact source bytes are less than the original 3.06 MB replay and canonical candidates and excluded are identical
- **AND** unrelated groups do not receive new lifecycle or delivery reads merely because sources were summarized

#### Scenario: Source transport migration

- **WHEN** the new endpoint and projection parity tests pass and the home is migrated
- **THEN** the home does not run both old and new source transports
- **AND** full source APIs remain available for existing editor flows

#### Scenario: Cache context isolation and eviction

- **WHEN** two backend contexts contain the same workspace identity, or a context observes a ninth workspace
- **THEN** neither backend receives the other's cached evidence, and the oldest scope is evicted while each retained workspace has only its latest observation
- **AND** invalidation tests cover successful device, rule, point and mapping mutations without changing their save barriers
