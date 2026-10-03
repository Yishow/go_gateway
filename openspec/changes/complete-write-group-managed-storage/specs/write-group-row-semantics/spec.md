## ADDED Requirements

### Requirement: Basic managed SQL preserves record time and origin identity
New basic per-device managed rows SHALL store stable record/group/device identities, UTC bucket time and member-level observed_at, quality and reason from the frozen row. Metadata bindings SHALL persist with the applied layout. Destination defaults, current UI scope and replay time MUST NOT replace acquisition facts. Existing custom-table metadata remains governed by its selected capabilities.

#### Scenario: Delayed delivery retains original time
- **GIVEN** a managed row accepted at its original acquisition and bucket times
- **WHEN** the destination receives it after an outage and gateway restart
- **THEN** SQL stores those original times and source/group identities, not the later flush or INSERT time.

#### Scenario: Missing and bad values are not fabricated
- **WHEN** the existing completeness policy yields an incomplete or explicitly allowed partial result
- **THEN** skipped/no_data remains visible or permitted SQL NULL carries the member reason; no zero, stale value or invented good quality is substituted.

### Requirement: Stable exact managed column generation
Managed columns SHALL be generated once from persisted member identities and verified exact-type capabilities, with collision-safe names and a frozen mapping. Display renames MUST NOT silently rebind data. Managed storage SHALL use the existing receipt-capable delivery contract; this change MUST NOT introduce a new selectable deduplication or sampling mode.

#### Scenario: Same labels and exact uint64 values
- **WHEN** two persisted Tags share a display label or a uint64 exceeds JavaScript's exact integer range
- **THEN** generated columns remain distinct and the SQL representation preserves the complete value without floating-point conversion.

#### Scenario: Replay and display rename
- **WHEN** a saved display name changes or an accepted row is retried
- **THEN** its frozen columns, record identity, payload digest and receipt-scoped effect remain unchanged, and target acknowledgement loss does not create a second effect.
