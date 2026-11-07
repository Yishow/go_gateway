## ADDED Requirements

### Requirement: Home shows live device connection state

The recording home SHALL derive one connection state for every workspace device from a single workspace-wide `GET /api/v1/datalink/runtime/status` read, polled every 5 seconds only while the page is visible. The state SHALL be the first matching rule of: `unavailable` (no successful read within the last 15 seconds, `snapshot_state.unavailable` is true, or no collector exists for the device), `disconnected` (`breaker_state` is not `closed` or `status` is `error`), `idle` (`running` is false, availability is `unavailable`, or `points_total` is 0), `degraded` (stale or error points exist, healthy points differ from total points, `status` is `warning`, `last_error` is present, or `last_read_at` is missing or older than the maximum of 3 × `effective_interval_ms` and 5000 ms), otherwise `collecting`. Each state SHALL be shown with text, an icon and a semantic color, together with healthy/total point counts and the relative last-read time. After a failed read the last successful states MAY remain visible for at most 15 seconds marked as outdated. The home MUST NOT show `collecting` when the status read is unavailable, and MUST NOT display raw backend errors.

#### Scenario: Device states are derived from collectors

- **WHEN** the runtime status read succeeds
- **THEN** each device shows the state selected by the ordered rules and its point and freshness facts

##### Example: state derivation at now = 12:00:10, interval 1000 ms

| Collector facts | Shown state |
| --- | --- |
| running, 48/48 healthy, last read 12:00:09, breaker closed | collecting |
| running, 46/48 healthy with 2 stale, status warning | degraded |
| running, 48/48 healthy, last read 12:00:04 (6 s old, threshold 5 s) | degraded |
| running, 48/48 healthy, last read 12:00:05 (exactly 5 s old) | collecting |
| running, 47/48 healthy with 0 stale and 0 error points | degraded |
| running false, status idle | idle |
| breaker open | disconnected |
| no collector for the device | unavailable |

#### Scenario: Status read fails

- **WHEN** the runtime status read fails after a previous success
- **THEN** the last observed states remain visible for at most 15 seconds marked as outdated, after which every device shows `unavailable`

#### Scenario: Page is hidden

- **WHEN** the browser tab is hidden or the operator leaves the home
- **THEN** status polling stops and resumes with one immediate read when the home becomes visible again

### Requirement: Home header summarizes devices, recordings and runtime writes

The recording home header SHALL show the device total with counts per connection state, the recording total with counts per saved status computed from all loaded recordings independent of filters, and `metrics.write_success_total` and `metrics.write_error_total` labelled as runtime cumulative counts that MUST NOT be presented as SQL-confirmed writes. Below the summary a device status strip SHALL list every device with its name, protocol, connection state and point facts; activating a device item SHALL apply the existing device filter to the recording list and activating it again SHALL clear that filter. The separate bottom device summary SHALL be removed from the recordings view, and the devices view SHALL show the same connection state per device.

#### Scenario: Summary counts match loaded data

- **WHEN** the home has loaded 3 devices and 4 recordings
- **THEN** the summary counts equal the per-state and per-status counts of those loaded items regardless of active filters

##### Example: mixed workspace

- **GIVEN** devices MC (degraded), mc_new (idle), M1 (idle) and recordings with statuses draft, draft, draft, removed
- **WHEN** the home renders
- **THEN** the summary reads devices 3 (degraded 1, idle 2) and recordings 4 (draft 3, removed 1)

#### Scenario: Device strip filters recordings

- **WHEN** the operator activates device MC in the status strip
- **THEN** only recordings with a member on MC remain listed and the device filter control shows MC
- **AND** activating MC again restores the unfiltered list

#### Scenario: Snapshot is stale

- **WHEN** `snapshot_state.stale` is true
- **THEN** device states remain shown and the summary displays an outdated-status notice

### Requirement: Recording rows show destination identity and write stages inline

Each recording row SHALL show its destination `schema.table` directly in a monospace secondary line with a copy control, and SHALL show write progress as four labelled stages — configuration, activation, SQL write, readback — each in exactly one of complete, unconfirmed or attention, derived from the existing home result summary. Historical or previous-revision SQL commits SHALL be labelled historical and treated as unconfirmed for the current revision. Stage details formerly behind disclosure controls (observed time, effective revision, requested boundary, owner inspection, test-write cleanup) SHALL be reachable by hover and keyboard focus on the stage. Recording rows SHALL NOT contain disclosure (`details`) elements for destination identity or result evidence. The connector identifier SHALL NOT be rendered as visible row text.

#### Scenario: Stage mapping

- **WHEN** a row renders with a result summary
- **THEN** each stage shows the state given by the mapping and an accessible text label

##### Example: stage mapping

| Summary fields | configuration | activation | SQL write | readback |
| --- | --- | --- | --- | --- |
| configuration=not_requested, effective=unknown, sql=unknown, readback=unknown | unconfirmed | unconfirmed | unconfirmed | unconfirmed |
| configuration=requested, effective=requested_effective, sql=revision_receipt, readback=test_write_verified | complete | complete | complete | complete |
| configuration=requested, effective=requested_not_effective, sql=historical_commit | complete | attention | unconfirmed (historical) | unconfirmed |
| configuration=requested, effective=requested_effective, current=delivery_attention | complete | complete | attention | unconfirmed |

#### Scenario: Copy destination table

- **WHEN** the operator activates the copy control for `public.gw_long_3fe66f8c5c98acb3`
- **THEN** that exact text is written to the clipboard and a status message confirms it, or when the clipboard is refused the text is selected and a status message states it could not be copied

### Requirement: Home offers recording and device actions with explicit confirmation

Each recording row SHALL keep its open action and add a keyboard-operable actions menu offering stop when the recording is running and delete when server action facts allow it; stale action facts SHALL disable both with a reload hint. Each device in the status strip and devices view SHALL offer delete through the same menu pattern. The recording delete dialog SHALL state what is permanently deleted, what external data is retained, that the action cannot be undone, and SHALL require an explicit acknowledgement before submission. The device delete dialog SHALL list recordings referencing the device; when any exist it SHALL block confirmation and explain they must be deleted first, otherwise it SHALL require typing the exact device name. Server conflicts SHALL keep the dialog open with an actionable explanation per error code and SHALL NOT be shown as success.

#### Scenario: Delete a recording from the home

- **WHEN** the operator opens the row menu, chooses delete, acknowledges and confirms, and the server returns success
- **THEN** the recording disappears from the list and the summary counts update without a page reload

#### Scenario: Delete is refused because of pending delivery

- **WHEN** the server returns `WRITE_GROUP_DELETE_BACKLOG` with 3 unconsumed samples and 2 pending outbox rows
- **THEN** the dialog shows the samples and outbox rows as separately labelled counts, offers stop and refresh actions, and the recording remains listed

#### Scenario: Device still referenced

- **WHEN** the operator chooses delete for a device that 2 non-deleted recordings reference
- **THEN** the dialog names both recordings, the confirm action is disabled, and no delete request is sent

#### Scenario: Device name mismatch

- **WHEN** an unreferenced device named `mc_new` is being deleted and the typed name is `mc-new`
- **THEN** the confirm action remains disabled

### Requirement: Home header has one primary action and visible freshness

The recording home header SHALL present add recording as the only primary action, with settings, logs and refresh as secondary actions. The two observation explanation lines SHALL be replaced by a last-updated time next to refresh, with the full explanation available from an adjacent help control. A resume-draft banner SHALL name the draft when one exists and SHALL be absent otherwise. The unconfirmed or incomplete recordings region SHALL be collapsible with a heading stating the item count, expanded by default when any item is not a tombstone and collapsed when all items are tombstones. At widths of 1440, 1024, 768 and 390 pixels the page SHALL NOT scroll horizontally, and at 390 pixels each recording SHALL render as a card whose secondary actions are in the actions menu.

#### Scenario: Only tombstones are unconfirmed

- **WHEN** the only unconfirmed item is a removed recording
- **THEN** the region is collapsed and its heading reads 1 item needs attention

#### Scenario: Narrow viewport

- **WHEN** the home renders at 390 pixels wide
- **THEN** the document scroll width equals the viewport width and recordings render as cards
