## MODIFIED Requirements

### Requirement: New configuration popups are centered viewport-sized modals
New simple-recording configuration popups SHALL use the modal size classifications in DESIGN.md: centered desktop form popups, approved full-viewport mobile and tablet presentation, and smaller auto-height desktop confirmations. The same responsive modal contract SHALL apply to existing modal primitives; non-modal informational popovers MUST remain excluded. The two main pages, including the recording editor, SHALL remain independently responsive and MUST NOT be converted into modals merely because they contain configuration.

#### Scenario: Opening a configuration popup
- **WHEN** an operator opens a device, point or supplemental configuration popup from the simple interface
- **THEN** the modal uses the appropriate DESIGN.md size classification, preserves the existing unsaved-change and focus-return contracts, and does not resize the main page to 75 percent.

### Requirement: Classified modals and contextual return state preserve safety
Normal desktop form popups SHALL retain their existing 75vw by 75dvh classification, with approved source and memory exceptions in DESIGN.md. At viewport widths up to 1024 CSS pixels, or with a coarse pointer up to 1366 CSS pixels, both form and confirmation modals SHALL use the full visible viewport with safe areas inside the modal. Low-height viewport handling SHALL keep controls reachable. Desktop short confirmations without complex forms SHALL retain the smaller auto-height classification. Modal headers and action footers SHALL remain visible while the body scrolls. Contextual help MUST NOT hide required fields, validation failures, disabled reasons, schema/history impacts, unresolved risks or consent. Return navigation SHALL preserve search, filters and scroll only within the same workspace, entity and view context and SHALL revalidate identity when revisions change.

#### Scenario: A mobile or tablet viewport opens a modal
- **WHEN** any classified form, confirmation or existing modal primitive opens on a mobile or tablet viewport
- **THEN** it uses the shared full-viewport presentation while keeping the title, actions, focused field and validation errors reachable
- **AND** orientation changes preserve the existing dirty guard, operation state and focus context.

#### Scenario: A short confirmation is displayed
- **WHEN** stop, remove or discard requires a short confirmation
- **THEN** the dialog identifies the object and consequence, retains explicit consent and does not substitute for a detailed schema-impact review
- **AND** desktop retains the smaller centered classification while mobile and tablet use the shared full-viewport presentation.

#### Scenario: Returning after an item was removed
- **WHEN** an operator returns to the same filtered list after its entity or revision changed
- **THEN** search, filter and position restore only within that context, unavailable selections are cleared with a reason, focus moves to a safe target and stale responses cannot restore old authority or action availability.

#### Scenario: Help is collapsed
- **WHEN** optional contextual help is closed
- **THEN** critical risk, errors, disabled reasons and consent remain visible and the help control remains keyboard accessible.
