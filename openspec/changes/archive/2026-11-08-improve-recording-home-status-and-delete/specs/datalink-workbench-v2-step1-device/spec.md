## MODIFIED Requirements

### Requirement: Persisted device deletion cascade

Deleting a persisted Step 1 device from the v2 workspace SHALL remove that device and its v2-owned related data from the backend. Before any detach or delete, the backend SHALL check within one transaction whether any non-tombstoned recording has a member or canonical key referencing the device; when one does, it SHALL return 409 `DEVICE_REFERENCED_BY_RECORDINGS` with action `remove_recordings_first`, the total count and up to 20 referencing recording identifiers and names, and SHALL change nothing. Tombstoned recordings that reference only this device SHALL be permanently deleted under the recording deletion contract before the device is deleted; if any of them cannot be deleted the device deletion SHALL fail without change. A tombstoned recording that also references another device SHALL refuse the deletion with `DEVICE_REFERENCED_BY_RECORDINGS`. Reference check, tombstone purge, workspace detach and device and source-rule deletion SHALL commit in one configuration transaction; the Modbus Share membership invalidation SHALL run after that commit through its existing fail-closed path, and when it fails the response SHALL be 500 `DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED` stating that the device was deleted, never a plain success. A failure to close the device connection after commit SHALL return 500 `DEVICE_DELETED_CONNECTION_CLOSE_FAILED`. A successful deletion SHALL record a `device_deleted` audit entry containing only the device identifier, name and protocol.

#### Scenario: Delete persisted device

- **GIVEN** a persisted device belongs to the singleton v2 workspace and no non-tombstoned recording references it
- **WHEN** the operator deletes that device
- **THEN** the backend removes the device from the workspace
- **AND** the device no longer appears after a page reload

#### Scenario: Device referenced by recordings

- **GIVEN** recordings `MC (Instantaneous)` and `100.76.76.1(mc)` are drafts with members on device MC
- **WHEN** a delete request for MC is sent
- **THEN** the response is 409 `DEVICE_REFERENCED_BY_RECORDINGS` listing both recordings with total 2
- **AND** the device, its source rules and both recordings are unchanged

#### Scenario: Device referenced only by a tombstone

- **GIVEN** removed recording `Test` is the only recording referencing device D1 and it has no pending delivery
- **WHEN** a delete request for D1 is sent
- **THEN** `Test` is permanently deleted, then D1 is deleted, and a `device_deleted` audit entry exists

#### Scenario: Tombstone with pending delivery blocks device deletion

- **GIVEN** removed recording `Test` references only device D1 and still has 2 retrying outbox rows
- **WHEN** a delete request for D1 is sent
- **THEN** the response is 409 `WRITE_GROUP_DELETE_BACKLOG` with the categorized counts, and D1, its source rules and `Test` are unchanged

#### Scenario: Share refresh fails after the device is deleted

- **WHEN** the configuration transaction commits but the Modbus Share membership invalidation fails
- **THEN** the response is 500 `DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED`, the device stays deleted and Modbus Share remains fail-closed
