## MODIFIED Requirements

### Requirement: Coexisting v2 workbench route

The system SHALL expose the guided workbench route at `/studio/v2` as the canonical product setup shell. The dedicated legacy `/studio` route SHALL be removed immediately with its proven legacy-only imports/chunks/assets, obsolete tests, and obsolete docs references. After deletion, `/studio` SHALL be indistinguishable from an arbitrary unknown route under the repository's generic unknown-route policy and SHALL NOT have a dedicated handler, tombstone, special redirect, or legacy workspace mount. The two product surfaces SHALL NOT be treated as an indefinite parallel compatibility contract.

#### Scenario: Operator opens `/studio/v2`

- **WHEN** the operator navigates to `/studio/v2`
- **THEN** the system loads the Workbench v2 shell with top bar, collapsible step rail, central step content, summary rail, and tweaks panel
- **AND** the route remains the canonical guided setup surface

#### Scenario: Deleted legacy route follows unknown-route policy

- **WHEN** the operator navigates to `/studio` and an arbitrary unknown path after deletion
- **THEN** both paths are handled identically by the repository generic unknown-route policy
- **AND** neither request loads legacy workspace code, a legacy-only chunk, or a `/studio`-specific redirect/handler

#### Scenario: Direct deep link to v2

- **WHEN** the operator opens a fresh browser tab with URL ending in `/studio/v2`
- **THEN** the v2 shell mounts with the default state: Step 1 highlighted, view mode flow, no completed steps, and summary rail visible on viewports at least 1280px

### Requirement: Default v2 entry route

The system SHALL treat `/studio/v2` as the default guided entry for the datalink setup flow. Removing `/studio` SHALL NOT change the existing generic unknown-route policy for unrelated paths.

#### Scenario: Root path enters v2

- **WHEN** an operator opens `/`
- **THEN** the system redirects to `/studio/v2`
- **AND** the first rendered setup surface is the v2 shell rather than a legacy page

#### Scenario: Unknown path policy is preserved

- **WHEN** an operator opens `/unknown-route-for-retirement`
- **THEN** the system applies the existing generic unknown-route policy
- **AND** the behavior is the same when the operator opens `/studio` after deletion

### Requirement: Legacy landing paths converge on /studio/v2

The system SHALL route generic legacy datalink landing paths to `/studio/v2` using their existing product semantics. Generic landing routes SHALL NOT navigate to the removed legacy workbench or introduce a new `/studio` compatibility route.

#### Scenario: Generic legacy datalink root lands on v2

- **WHEN** an operator opens `/datalink`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`

#### Scenario: Generic legacy workbench path lands on v2

- **WHEN** an operator opens `/datalink/workbench`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`

#### Scenario: Generic landing does not recreate legacy workspace

- **WHEN** a generic landing route is opened after `/studio` deletion
- **THEN** the route resolves through its existing v2 entry behavior
- **AND** no legacy workspace component or legacy-only chunk is mounted

### Requirement: i18n namespace for v2

The v2 shell SHALL load all user-facing strings from the dedicated i18n namespace `workbench-v2`, available in zh-TW and en. All button labels, breadcrumb fragments, step titles, step subtitles, placeholder copy, summary labels, typed error messages, retry actions, accessibility names, Share states, and CommitProgress status text MUST resolve through i18n and MUST NOT be hardcoded in the component.

#### Scenario: zh-TW locale shows complete operator copy

- **WHEN** the locale is zh-TW and the v2 page mounts
- **THEN** step titles, settings, error actions, Share labels, and progress statuses render in Traditional Chinese
- **AND** no new user-facing English or raw backend string is shown

#### Scenario: en locale shows complete operator copy

- **WHEN** the locale is en and the v2 page mounts
- **THEN** step titles, settings, error actions, Share labels, and progress statuses render in English
- **AND** no raw backend diagnostic string is shown

## ADDED Requirements

### Requirement: V2 error presentation separates operator copy from diagnostics

The v2 shell SHALL map typed backend/runtime error codes to localized actionable copy and SHALL keep raw exception details in a separate diagnostics channel. An unknown code SHALL render a safe generic message with an opaque request identifier.

#### Scenario: Typed error has localized retry action

- **WHEN** a v2 API response has code `runtime_stream_unavailable` and `retryable` true
- **THEN** the shell renders the localized stream-unavailable message and retry action
- **AND** the shell does not render the raw response message or stack trace

#### Scenario: Unknown error code remains safe

- **WHEN** a v2 API response has an unknown code
- **THEN** the shell renders a localized generic error and request identifier
- **AND** raw diagnostic text is absent from the DOM

### Requirement: Share toggle and CommitProgress use truthful accessible states

The v2 Share/global toggle SHALL expose a localized accessible name, `aria-pressed`, visible focus, and keyboard activation through Enter and Space. The control SHALL be disabled during its pending UI mutation and SHALL issue at most one request per activation gesture. CommitProgress SHALL render pending, success, failed, and skipped states from actual results; failed entries SHALL show safe error copy and SHALL NOT show a fixed 200 success value. This requirement SHALL NOT add backend CAS or Share concurrency semantics.

#### Scenario: Share toggle is keyboard accessible and single-flight

- **WHEN** an operator activates the Share/global toggle by keyboard or pointer
- **THEN** the control exposes its localized name and `aria-pressed` state, becomes disabled while pending, and sends one request
- **AND** another activation during pending sends no second request

#### Scenario: Failed progress remains failed

- **WHEN** a CommitProgress item returns a failed result
- **THEN** its row shows failed state, safe error copy, and an actionable retry
- **AND** the row does not show a fixed 200 or success presentation
