## MODIFIED Requirements

### Requirement: Safe basic group lifecycle
The system SHALL support group rename, member and policy edits, disable and logical deletion under CAS. Rename MUST preserve the stable group and existing record/receipt identities without restarting the applied writer. Member/column/policy/destination changes SHALL save a draft and take effect only on explicit valid apply at the next bucket boundary. Disable SHALL stop new intake while retaining accepted backlog. Delete MUST tombstone the group, stop new intake and preserve immutable payloads, revisions, ownership, receipts and operation/backlog lookup until every accepted record is safely resolved; physical purge is outside this change. A saved draft MUST NOT retire an otherwise eligible applied version, including after restart.

#### Scenario: Rename and member edits
- **WHEN** a group is renamed or members are added or removed while it is running
- **THEN** rename keeps identity stable, semantic edits remain draft until apply, the old bucket closes under its original snapshot, and the new membership starts at the next bucket boundary.

#### Scenario: Delete with backlog
- **WHEN** the operator deletes a group with accepted pending records
- **THEN** new intake stops and the tombstone retains queryable ownership and immutable delivery evidence; if the legacy format cannot preserve that ownership, deletion is blocked with a disable alternative. Completing accepted records remains the responsibility of the production durable delivery path.

#### Scenario: Saved draft leaves the running revision active
- **GIVEN** a running group and a newly saved name or member draft
- **WHEN** two record intervals pass without Apply, including a gateway restart
- **THEN** the old applied version continues accepting and recording under its original membership; no draft value is used and no retirement is caused by draft status.

#### Scenario: Concurrent cutoff and acquisition
- **WHEN** acquisition races with a validated version cutoff or Disable
- **THEN** each sample is accepted or refused consistently against one cutoff, without a data race or acceptance by competing writer owners.

#### Scenario: Re-enable after a completed retirement
- **WHEN** an operator explicitly re-enables a disabled group with unchanged row semantics
- **THEN** eligible new intake resumes without reviving closed buckets, retaining an expired in-memory cutoff or duplicating accepted effects.
