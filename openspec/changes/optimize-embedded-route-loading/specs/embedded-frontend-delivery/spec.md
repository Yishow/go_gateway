## MODIFIED Requirements

### Requirement: Frontend build emits a complete embedded asset graph

The frontend build SHALL emit frontend/dist/index.html and every JavaScript, CSS, route chunk, vendor chunk, and static asset required by its direct and transitive imports. References SHALL resolve within the embedded static root. Vendor grouping SHALL preserve stable logical React/runtime, charting, router/i18n/query, residual dependency groups, and Material UI/Emotion groups only when those dependencies are present. The retired /studio route SHALL use the generic unknown-route policy and MUST NOT regain a dedicated route or asset.

#### Scenario: Build emits entry and lazy route assets

- **WHEN** the repository frontend build runs with its declared dependencies
- **THEN** frontend/dist/index.html exists and the asset graph resolves for /studio/v2, /studio/runtime, /studio/logs, /test, and /gateway/*
- **AND** every discovered asset reference resolves to a file under frontend/dist

##### Example: vendor group names

- **GIVEN** a production build includes React, charting, router, i18n and query dependencies
- **WHEN** Vite writes the production output
- **THEN** matching dependencies have stable logical react-vendor, charts-vendor, app-vendor or misc-vendor groups, and mui-vendor exists only when its matching dependencies exist
- **AND** application route modules remain separately loadable without changing route semantics

## ADDED Requirements

### Requirement: Initial recording load excludes unused editors

The Studio frontend SHALL render its stable reading frame before deferred advanced editor code settles. The initial recording-home dependency graph MUST NOT eagerly require advanced-only editor modules. Loading and failure SHALL remain distinguishable without losing draft navigation guards.

#### Scenario: Home opens before advanced code

- **WHEN** the user opens /studio/v2 in recording-home mode while advanced editor imports are delayed
- **THEN** the reading frame is visible and the advanced editor is not mounted
- **AND** the home does not show an empty successful workspace while evidence is pending

#### Scenario: Deep link or failed import

- **WHEN** the user opens an existing advanced deep link or a required chunk rejects
- **THEN** query/hash intent reaches the same destination after loading, or the existing actionable error boundary is shown
- **AND** Back/Forward and dirty-draft cancellation do not clear the draft or create duplicate navigation

### Requirement: Static representations preserve HTTP and version semantics

The embedded static server SHALL negotiate gzip or identity for hashed JS/CSS without compressing API or SSE responses. It SHALL preserve correct MIME, representation-specific validators, Vary, GET/HEAD behavior and identity Range responses only when identity is acceptable. HTML SHALL revalidate and only content-hashed assets SHALL use immutable caching.

#### Scenario: Encoding negotiation and range

- **WHEN** a client requests a hashed JS asset with gzip accepted, gzip forbidden, or a Range header
- **THEN** a supported compressed representation is selected only when accepted, gzip;q=0 receives identity when accepted, and Range receives identity only when accepted, otherwise 406
- **AND** HEAD sends the same selected representation headers as GET with no response body

##### Example: no acceptable encoding

- **GIVEN** an existing hashed CSS asset and Accept-Encoding: identity;q=0, gzip;q=0
- **WHEN** the asset is requested
- **THEN** the response is 406 with no substituted representation

#### Scenario: Upgrade and missing asset

- **WHEN** a new build replaces the embedded asset graph
- **THEN** HTML is revalidated, unchanged hashed assets remain cacheable, and validators cannot match different encoded bytes
- **AND** a missing asset returns 404 rather than SPA HTML or a successful route

### Requirement: Performance evidence separates transfer and interaction

The change SHALL record build-graph bytes, encoded transfer bytes, request counts, and browser interaction timings separately against the same fixture and toolchain. It MUST NOT present local replay or DOM tests as real-browser acceptance.

#### Scenario: Before and after comparison

- **WHEN** the implementation is assessed on the same clean build and cold/warm fixture
- **THEN** home eager JavaScript bytes and eligible encoded transfer bytes are smaller than their recorded baselines, with the full route graph still loadable
- **AND** actual browser runs cover cold/warm loading, deep links, chunk failure and version changes, or remain explicitly blocked without a passed claim
