## ADDED Requirements

### Requirement: Confirmed template intent and recoverable failures
The system SHALL apply a real unconfirmed template preview only after an explicit current-device confirmation, invalidate previews on device/template changes, and expose safe actionable load, preview and apply errors.

#### Scenario: Explicit confirmation of a real preview
- **WHEN** a current preview has confirmed=false and the operator confirms it
- **THEN** the client SHALL submit confirmed=true and the backend SHALL continue rejecting an unconfirmed submission.

#### Scenario: Changed scope or failed operation
- **WHEN** device/template changes or load/preview/apply fails through disconnection, conflict or server failure
- **THEN** stale previews and replies SHALL NOT be applied and the operator SHALL see a safe retry action without raw exception details.

#### Scenario: Offline intent and template reload
- **WHEN** an explicit template operation is attempted offline or a failed template load is retried
- **THEN** the operation SHALL fail visibly rather than queue consent for automatic reconnect submission, and loading SHALL preserve the expanded repair section.
