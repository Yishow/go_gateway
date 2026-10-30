## Purpose

Provide an evidence-led recording home inside Studio V2 that helps operators identify the next safe action while retaining access to all existing recording configurations.

## ADDED Requirements

### Requirement: Home guidance projects canonical scoped evidence
The system SHALL derive recording home guidance from the existing workspace, source, group, lifecycle and delivery authorities. It MUST NOT introduce a second persisted configuration authority or mutate configuration through navigation.

#### Scenario: Opening the recording home
- **WHEN** an operator opens /studio/v2?view=recordings
- **THEN** the page reads the current scoped evidence and offers navigation actions without saving, activating, creating a table or restarting a recording.

#### Scenario: A previous scope returns late
- **WHEN** a request for an older workspace or group revision completes after the operator has changed scope
- **THEN** that response does not replace the current home evidence or enable an action in the new scope.

### Requirement: Guidance states coexist with deterministic priority
The system SHALL display loading, failed bootstrap and unprovable scope independently from empty or healthy data. For available data, unresolved operations or actionable faults SHALL precede drafts, missing devices, missing points, missing recordings and healthy summaries as the primary guidance, while every other item remains reachable.

#### Scenario: A fresh workspace has no devices
- **WHEN** a successfully loaded workspace has no devices or recording rules
- **THEN** the primary action invites creating the first recording and begins with device connection.

#### Scenario: A device lacks points
- **WHEN** a workspace has a device but no usable source points and no higher-priority issue
- **THEN** guidance opens that device's point configuration rather than requiring its connection to be recreated.

#### Scenario: Existing points are ready for a new recording
- **WHEN** usable points exist without recordings or higher-priority issues
- **THEN** the primary action selects existing points to create a recording.

#### Scenario: Draft and failed recording coexist
- **WHEN** recording A has an unresolved operation, recording B has historical committed evidence and draft C lacks a destination
- **THEN** the primary action checks A's operation, B remains inspectable, and C offers a named continuation at the missing destination section.

#### Scenario: Bootstrap fails
- **WHEN** workspace evidence cannot be loaded
- **THEN** the page shows an actionable unavailable state and never renders the workspace as empty or healthy.

### Requirement: Healthy summaries distinguish history from present evidence
The system SHALL label saving, start acknowledgement, collection, SQL commit and verified readback distinctly. A historical SQL commit SHALL NOT establish present freshness, and general delivery evidence SHALL NOT be presented as raw row values unless the existing readback contract supplies them.

#### Scenario: A previously successful device disconnects
- **WHEN** a recording has a historical commit at 10:32 but current source evidence is stale or disconnected
- **THEN** the home retains the last successful time and destination while exposing the current fault and a scoped repair action.

#### Scenario: First start has no SQL receipt
- **WHEN** start has been acknowledged but no matching commit evidence exists
- **THEN** the home states that it is awaiting the first write rather than claiming successful recording.

### Requirement: Optional home preserves existing Studio navigation
The system SHALL expose the recording home as an explicit view of /studio/v2 without changing the default entry in this change. Valid entity locators and runtime return intents SHALL retain precedence and scope. The system SHALL validate raw query encoding and uniqueness before removing only a single recognized presentation hint for the existing strict locator parser; invalid locator input SHALL remain invalid and MUST NOT fall back to the general home. Existing runtime, logs, test and advanced paths SHALL remain accessible.

#### Scenario: A deep link includes a recording home hint
- **WHEN** a valid entity locator or runtime repair focus is accompanied by view=recordings
- **THEN** the exact existing locator or repair intent is honored without losing its query, hash or scope.

#### Scenario: An unknown view is supplied
- **WHEN** the URL contains an unsupported view value without any locator intent
- **THEN** the application uses the existing entry behavior without mutating data or restoring the retired /studio surface.

#### Scenario: Presentation hints do not weaken locator validation
- **WHEN** a locator request contains malformed encoding, duplicate view values, an unknown presentation value or an invalid locator key
- **THEN** it remains invalid with no fallback to an unscoped home, while a valid locator with one recognized view retains its original scoped destination.

### Requirement: Theme preference resolves without inventing persisted intent
The system SHALL use one existing ThemeProvider with separate system/light/dark preference and resolved light/dark theme. The existing server theme SHALL remain global rather than introducing per-workspace preferences. Initial resolution SHALL prefer a valid actually stored server theme key, then a versioned local preference, then a valid legacy light/dark value, then system. Mapper-injected defaults MUST NOT count as stored intent. Legacy values SHALL be preserved as compatibility settings without claiming they prove manual choice. The auto setting SHALL map to system.

#### Scenario: No stored preference exists
- **WHEN** raw server settings have no theme key and no valid local preference exists
- **THEN** the provider follows the current OS theme, subscribes to subsequent OS changes, and does not persist the resolved color as an override.

#### Scenario: Explicit override and return to system
- **WHEN** an operator selects light while the OS is dark, then later selects system
- **THEN** OS changes do not override the explicit light choice, and selecting system immediately follows the current OS and future changes.

#### Scenario: Legacy value and mapper default coexist
- **WHEN** legacy local storage contains light, the raw server theme key is absent, and the settings mapper supplies dark
- **THEN** light remains the compatibility choice with an available system option rather than treating the mapper default as a stored server preference.

#### Scenario: Persistence or browser APIs are unavailable
- **WHEN** storage, matchMedia or a theme-save operation fails
- **THEN** the interface remains usable, unavailable OS resolution falls back to light without saving an override, and save failure is not represented as persisted success.

#### Scenario: An old hydration response arrives
- **WHEN** an operator changes preference or backend/app context before an earlier settings request completes
- **THEN** stale hydration does not overwrite the new choice or import another backend context's preference, and no recording state or device operation is changed.

#### Scenario: Theme-only choice does not flush business settings
- **WHEN** an operator changes theme while scheduler or Modbus Share edits remain unsaved
- **THEN** only updateKey for theme may persist, other setting keys remain unwritten, workspace navigation does not reset the global preference, and OS-only resolution changes perform no persistence at all.

### Requirement: Semantic design tokens govern affected surfaces
The new recording surfaces and reused affected controls SHALL implement the root DESIGN.md primitive, semantic light/dark and component-state token contract through one CSS variable authority. Typography, spacing, radii, elevation, focus, contrast and reduced motion SHALL be consistent across pages and portals. Raw color literals MUST NOT be duplicated in new component styling. Unaffected legacy pages MUST NOT require a whole-application redesign.

#### Scenario: Theme changes while a modal is open
- **WHEN** the resolved theme changes with an unsaved configuration modal or pending operation visible
- **THEN** the page, portal, native controls and status messages use the same theme while preserving field values, focus and operation identity without recording, device or database-schema commands. Explicit preference selection may only use the existing theme-key update mutation, not a full settings snapshot; OS resolution changes do not save a preference.

#### Scenario: Information must remain readable
- **WHEN** either theme renders secondary text, placeholders, disabled explanations, status or keyboard focus
- **THEN** informative text meets 4.5:1 contrast and required controls/focus meet 3:1 against their actual adjacent backgrounds, status has a non-color label, and reduced motion removes nonessential animation.

### Requirement: New configuration popups are centered viewport-sized modals
New simple-recording configuration popups SHALL use the centered modal size classifications in DESIGN.md: normal form popups at 75 percent of viewport width and height, approved near-fullscreen mobile/short-viewport presentation, and approved smaller auto-height short confirmations. The two main pages, including the recording editor, SHALL remain independently responsive and MUST NOT be converted into modals merely because they contain configuration.

#### Scenario: Opening a configuration popup
- **WHEN** an operator opens a device, point or supplemental configuration popup from the simple interface
- **THEN** a centered modal uses the appropriate DESIGN.md size classification rather than a side drawer, preserves the existing unsaved-change and focus-return contracts, and does not resize the main page to 75 percent.

### Requirement: Classified modals and contextual return state preserve safety
Normal form popups SHALL use 75vw by 75dvh when viewport width is at least 640px and height at least 480px. Below either boundary, approved near-fullscreen presentation SHALL leave an 8px outer inset and respect safe areas and the visible viewport. Short confirmations without complex forms SHALL use a maximum 420px width, 16px outer inset and auto height with a bounded scrollable body. Modal headers and action footers SHALL remain visible while the body scrolls. Contextual help MUST NOT hide required fields, validation failures, disabled reasons, schema/history impacts, unresolved risks or consent. Return navigation SHALL preserve search, filters and scroll only within the same workspace, entity and view context and SHALL revalidate identity when revisions change.

#### Scenario: A narrow or short viewport opens a form
- **WHEN** a form popup is opened at a width below 640px or height below 480px
- **THEN** the near-fullscreen classification keeps the title, actions, focused field and validation errors reachable without clipping, while preserving the existing dirty guard and operation state.

#### Scenario: A short confirmation is displayed
- **WHEN** stop, remove or discard requires a short confirmation
- **THEN** the smaller centered auto-height dialog identifies the object and consequence, retains explicit consent and does not substitute for a detailed schema-impact review.

#### Scenario: Returning after an item was removed
- **WHEN** an operator returns to the same filtered list after its entity or revision changed
- **THEN** search, filter and position restore only within that context, unavailable selections are cleared with a reason, focus moves to a safe target and stale responses cannot restore old authority or action availability.

#### Scenario: Help is collapsed
- **WHEN** optional contextual help is closed
- **THEN** critical risk, errors, disabled reasons and consent remain visible and the help control remains keyboard accessible.

### Requirement: Development style preview is isolated from operator navigation
A development/test-only style preview SHALL reuse the actual components and semantic tokens to exercise light/dark, inputs, actions, statuses, errors, disabled, focus and modal states. It MUST NOT appear in the operator primary menu or become a third primary workflow. The existing build-mode guard SHALL exclude this development entry from production bundles. The preview MUST NOT write configuration or data, activate resources, collect credentials or fabricate production bootstrap. A mock demo MAY expose a dedicated non-primary link or hash and SHALL label its fixtures as simulated.

#### Scenario: Preview is inspected in development
- **WHEN** a developer opens the dedicated development preview
- **THEN** real shared components and tokens render the selected states without backend mutations, and ordinary operator navigation has no preview menu item.

#### Scenario: Production build is generated
- **WHEN** the production frontend bundle is built
- **THEN** the development preview entry is excluded and cannot become a reachable third operator page, while the actual recording surfaces retain their shared component styling.

### Requirement: Large recording collections remain searchable and bounded
The home SHALL provide compact management lists with keyword search across recording, device and configured table labels, device/status/destination filters, deterministic ordering with canonical identity tie-breaks and 25-or-50-row pagination as design limits. The device summary SHALL show at most four rows with a same-page management view, not a third primary page. Attention evidence SHALL remain available independently of list filters. Client-side projection over loaded canonical data MUST NOT be described as an existing bounded cross-device backend search service.

#### Scenario: Many devices and recordings are loaded
- **WHEN** the available dataset exceeds one page
- **THEN** only the chosen page renders in the complete management list, the device summary remains bounded, and search/filter can locate a named recording without traversing an unbounded card wall.

#### Scenario: A filtered return targets a removed item
- **WHEN** an operator returns from editing and the matching scope revision has removed the selected item or reduced the number of pages
- **THEN** filters and sort restore only in the same context, selection is revalidated, the page clamps to a valid page, and focus/scroll return safely without stale authority.

#### Scenario: Search has no matches
- **WHEN** valid loaded data exists but the current search or filters match no entries
- **THEN** the page shows a no-results state with a clear-filter action rather than claiming the workspace is empty or healthy.
