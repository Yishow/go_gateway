## ADDED Requirements

### Requirement: Activation autosave settlement is time-bounded

The frontend activation flow SHALL wait for in-flight autosaves for a bounded period. When the bound elapses while saves are still reported as in flight, the wait SHALL resolve with the current autosave barrier snapshot instead of remaining pending, and activation SHALL fail with the existing retryable save-incomplete error. The activation UI SHALL NOT remain in an activating state with no outcome.

#### Scenario: Settled saves resolve the wait immediately

- **WHEN** activation starts and no autosave reports an in-flight state
- **THEN** the wait resolves without delay with a barrier snapshot taken from the current state
- **AND** activation proceeds to the remaining activation gates

#### Scenario: A stuck in-flight save resolves as save-incomplete after the bound

- **GIVEN** an autosave entry reports an in-flight save state that never settles
- **WHEN** the operator starts activation and the bounded wait elapses
- **THEN** the wait resolves with a barrier snapshot reporting at least one pending save
- **AND** activation fails with modbus_share_save_incomplete marked retryable
- **AND** the activation UI leaves the activating state and offers a retry

##### Example: bounded wait outcomes

| In-flight saves at start | Saves settle before bound | Wait outcome | Activation result |
| ------------------------ | ------------------------- | ------------ | ----------------- |
| none | not applicable | resolves immediately | proceeds to remaining gates |
| one | yes | resolves when the last save settles | proceeds to remaining gates |
| one | no | resolves at the bound with pending_saves >= 1 | modbus_share_save_incomplete, retryable |

#### Scenario: Resolved waiters are not resolved twice

- **GIVEN** a bounded wait has already resolved by timeout
- **WHEN** the autosave state later settles
- **THEN** the previously resolved wait SHALL NOT be resolved again
- **AND** no pending timer for that wait remains active

### Requirement: Disabled Share preserves workspace-level activation gates

The activation barrier validator SHALL separate Share-specific gates from workspace-level gates. When the global Share setting is disabled and Share hydration is ready, the validator SHALL still enforce the readiness token and the workspace revision, and SHALL skip only the Share settings-revision gate. When the global Share setting is disabled and Share hydration is not ready, the validator SHALL allow activation so that an unbootstrapped Share does not block workspace activation.

#### Scenario: Disabled Share with ready hydration still rejects a stale readiness token

- **GIVEN** the global Share setting is disabled and Share hydration reports ready with a readiness token
- **WHEN** an activation request carries a missing or non-matching readiness token
- **THEN** the request is rejected with modbus_share_save_incomplete marked retryable
- **AND** no workspace activation is attempted

#### Scenario: Disabled Share with ready hydration still rejects a stale workspace revision

- **GIVEN** the global Share setting is disabled and Share hydration reports ready
- **WHEN** an activation request carries a workspace revision that does not match the hydrated workspace revision
- **THEN** the request is rejected with modbus_share_revision_conflict marked retryable

#### Scenario: Disabled Share with unready hydration activates the workspace

- **GIVEN** the global Share setting is disabled and Share hydration reports a failed or not-ready state
- **WHEN** an activation request arrives
- **THEN** the workspace activation proceeds
- **AND** no Share restore is invoked

#### Scenario: Enabled Share keeps the full gate sequence

- **GIVEN** the global Share setting is enabled
- **WHEN** an activation request arrives
- **THEN** hydration readiness, readiness token, settings revision, and workspace revision are all enforced as before

##### Example: gate coverage by Share state

| Share enabled | Hydration ready | Hydration gate | Readiness token gate | Settings revision gate | Workspace revision gate |
| ------------- | --------------- | -------------- | -------------------- | ---------------------- | ----------------------- |
| true | any | enforced | enforced | enforced | enforced |
| false | true | skipped | enforced | skipped | enforced |
| false | false | skipped | skipped | skipped | skipped |

### Requirement: Share settings read failure is not reported as disabled

When the activation barrier validator fails to read the global Share settings, it SHALL report a retryable internal precondition failure and SHALL NOT report the Share-disabled error code. The validator SHALL NOT read revision fields from an unpopulated settings value when the read failed.

#### Scenario: Settings read failure surfaces as a retryable precondition failure

- **GIVEN** the Share settings read returns an error
- **WHEN** an activation request arrives
- **THEN** the response error code is not the Share-disabled code
- **AND** the error is marked retryable
- **AND** no settings revision value from an unpopulated settings value is included in the response
