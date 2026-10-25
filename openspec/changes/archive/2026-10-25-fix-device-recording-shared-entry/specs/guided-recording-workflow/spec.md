## MODIFIED Requirements

### Requirement: Purpose-led layout with one advanced authority

The system SHALL label the two new managed paths as whole-device snapshot (wide) and selected measurements (long), explain their row granularity, and allow instantaneous and cumulative roles in either layout. Advanced SHALL open the same saved group with the same profile/reduction/interval and explicit saved-versus-applied state. It SHALL not expose a competing interval or start authority. Every database device SHALL enter this shared recording setup by default, including a device with a saved legacy managed wide group. The presence of that group SHALL NOT automatically replace the shared setup with its legacy editor or block configuration of a new recording profile. Legacy wide/custom settings SHALL remain editable under their own label without silent conversion; the legacy managed editor SHALL require an explicit operator action and SHALL retain its saved values and scoped drafts independently of new profile drafts. A legacy applied revision SHALL be described as saved applied settings rather than proof of current collection. New starts SHALL use the existing scoped request and server ownership checks; layout changes and takeover of still-owned sources SHALL use reviewed version/ownership transitions instead of changing IDs, deleting groups or disabling an existing writer automatically. The existing explicit preparation-only context SHALL remain unable to restore or send a recording start.

#### Scenario: Full meter and mixed selected sources
- **WHEN** one device selects full voltage/current/power/counter snapshots and another selects temperature plus water/gas/electric counters
- **THEN** the first can use managed wide and the second managed long, each preserves explicitly selected roles and periods, and neither layout implies that counters are averaged.

#### Scenario: Advanced round trip
- **WHEN** an operator opens a profile from Basic in Advanced and returns without saving
- **THEN** group ID, reduction, interval, members and destination remain identical and no mutation or additional start occurs.

#### Scenario: Legacy and new devices share the default setup
- **WHEN** A has a saved legacy managed wide group and B has none, and the operator switches between them
- **THEN** both devices show the shared purpose and instantaneous/cumulative setup, their new drafts remain distinct, and viewing either device changes no saved group or writer.

#### Scenario: Explicit legacy editing preserves independent drafts
- **WHEN** the operator opens legacy editing on A, changes its 60-second draft to 73 seconds, returns to the new setup with a 47-second draft, or reloads
- **THEN** the corresponding draft is restored only in its own scope and editor, and neither value is automatically saved or copied to the other draft or device.

#### Scenario: Never-applied legacy draft does not join a new start
- **WHEN** A retains a never-applied legacy draft and the operator starts an already prepared new profile
- **THEN** the existing scoped request includes only A and its reviewed profile groups; the legacy draft is not updated, applied, disabled or deleted.

#### Scenario: Applied legacy ownership remains authoritative
- **WHEN** a device has an applied legacy revision and the server refuses a new start that would take over its owned sources
- **THEN** the UI exposes the refusal with the existing settings and a review-cutover entry, preserves the old group, and does not bypass ownership through a second Apply or workspace activation; explicit preparation-only mode still performs no recording start.
