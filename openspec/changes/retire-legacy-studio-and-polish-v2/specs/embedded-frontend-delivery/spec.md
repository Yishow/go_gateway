## MODIFIED Requirements

### Requirement: Existing route product semantics remain unchanged during delivery refactoring

The embedded delivery implementation SHALL retain the canonical route contract: the default entry resolves to `/studio/v2`, `/studio/runtime` remains the focused monitor, `/test` remains an independent engineering tool, and `/gateway/*` remains experimental. The `/studio` route SHALL be removed immediately with its dedicated registration, handler/tombstone, and proven legacy-only assets. After deletion, `/studio` SHALL be handled exactly like an arbitrary unknown route under the repository generic unknown-route policy, without loading legacy workspace code. Lazy loading and chunk grouping SHALL NOT alter supported route parameters or non-target route identities.

#### Scenario: Default and deleted routes retain their identities

- **WHEN** an operator opens the root entry, generic datalink landing, `/studio`, and `/studio/v2` after deletion
- **THEN** root and generic landing resolve to the v2 entry contract, `/studio` follows the same generic unknown-route policy as an arbitrary unknown path, and `/studio/v2` renders the guided user-facing entry
- **AND** no dedicated `/studio` route/handler/tombstone, special redirect, legacy mount, or legacy-only asset request occurs

#### Scenario: Monitor, tool, and experimental routes remain distinct

- **WHEN** an operator opens `/studio/runtime` with its existing device context, `/test`, and one `/gateway/*` route
- **THEN** `/studio/runtime` renders the focused runtime monitor contract
- **AND** `/test` remains an independent engineering tool
- **AND** `/gateway/*` remains an experimental surface
- **AND** lazy loading introduces no redirect from these non-target routes to an unrelated product flow

##### Example: Route family identity checks

- **GIVEN** the runtime URL includes `device_id=device-A` and the gateway URL is `/gateway/quick-setup`
- **WHEN** the operator opens `/studio/runtime?device_id=device-A`, `/test`, and `/gateway/quick-setup`
- **THEN** the runtime page remains the focused monitor for device-A
- **AND** `/test` remains the independent engineering tool
- **AND** `/gateway/quick-setup` remains experimental without a redirect to `/studio/v2`

#### Scenario: Deleted route is equivalent to an arbitrary unknown path

- **WHEN** a route test requests `/studio` and `/unknown-route-for-retirement` after deletion
- **THEN** both requests produce the same response/navigation, status/fallback, error boundary, and network asset pattern
- **AND** neither request loads a legacy route chunk or legacy-only asset

### Requirement: Embedded browser smoke proves the full asset graph is loadable

A built embedded binary served with a synchronized static tree SHALL load the canonical and non-target route families without a JavaScript chunk or vendor asset returning HTTP 404. The smoke verification SHALL check network responses and generic unknown-route equivalence for `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and `/gateway/*`.

#### Scenario: Embedded binary loads all required route families

- **GIVEN** a binary built from a fresh checkout after frontend assets are synchronized
- **WHEN** browser smoke opens the five route families and waits for generic unknown-route handling to settle
- **THEN** every required JavaScript, CSS, and route asset request returns a non-404 response
- **AND** each preserved route identity and `/studio` unknown-route equivalence result matches the product contract
- **AND** the smoke check exits successfully

#### Scenario: Stale or missing asset fails the smoke check

- **GIVEN** one referenced v2, runtime, tool, experimental, or shared asset is absent from `cmd/test_ui/static`
- **WHEN** browser smoke opens the route that imports the missing asset
- **THEN** the check reports the failed asset request
- **AND** the check exits nonzero
- **AND** the missing-asset condition is not hidden by a loading fallback

## ADDED Requirements

### Requirement: Final removal preserves non-target embedded delivery

After immediate deletion the embedded asset graph SHALL contain the `/studio/v2`, runtime, test, and gateway route assets required by their contracts and SHALL NOT contain a legacy-only workspace chunk, import, or asset. The `/studio` request SHALL be indistinguishable from an arbitrary unknown route and SHALL NOT load a dedicated route/handler/tombstone, legacy component, or legacy-only asset.

#### Scenario: Legacy asset is absent after removal

- **WHEN** an embedded browser requests `/studio` alongside an arbitrary unknown path after deletion
- **THEN** both requests produce the same response/navigation without a dedicated `/studio` handler or tombstone
- **AND** network evidence shows no legacy-specific redirect logic, legacy workspace component, or legacy-only asset request

#### Scenario: Non-target assets remain available

- **WHEN** an embedded browser opens `/test` and `/gateway/quick-setup` after deletion
- **THEN** their required route chunks and assets load successfully
- **AND** the removal does not redirect or delete those non-target surfaces
