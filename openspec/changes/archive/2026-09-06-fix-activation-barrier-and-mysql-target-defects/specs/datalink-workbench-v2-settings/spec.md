## ADDED Requirements

### Requirement: Settings save state reflects only real save requests

The settings page SHALL set the Modbus Share save state to the in-flight value only when a save request is actually dispatched. Editing a settings field and resetting settings to defaults SHALL NOT set the in-flight save state. Unsaved-change indication SHALL be carried by the save bar, not by the save state field consumed by the activation barrier.

#### Scenario: Editing a Modbus Share field does not fake an in-flight save

- **WHEN** the operator edits any Modbus Share field
- **THEN** the Modbus Share save state is not set to the in-flight value
- **AND** any previously recorded save error for that section is cleared

#### Scenario: Resetting settings does not fake an in-flight save

- **WHEN** the operator resets settings to defaults
- **THEN** the Modbus Share save state is not set to the in-flight value

#### Scenario: Saving sets and clears the in-flight save state

- **WHEN** the operator triggers a settings save
- **THEN** the Modbus Share save state becomes the in-flight value while the request is dispatched
- **AND** the save state becomes the saved value once every setting entry is persisted

##### Example: save state transitions

| Action | Save state after action |
| ------ | ----------------------- |
| edit a Modbus Share field | unchanged by the edit |
| reset settings to defaults | unchanged by the reset |
| trigger save, request pending | in-flight |
| trigger save, request succeeds | saved |
| trigger save, request fails | save-error |

### Requirement: Settings save always converges to a terminal state

A dispatched settings save SHALL drive the Modbus Share save state to a terminal value — saved on success, save-error on failure — even when the operator edits other settings fields while the request is in flight. On failure the operator SHALL see the error and SHALL be offered a retry. The settings revision returned by the server is a server-assigned concurrency token, not operator-editable data: it SHALL be adopted unconditionally so that later saves cannot be permanently rejected as revision conflicts.

#### Scenario: Editing during an in-flight save does not swallow the failure

- **GIVEN** a settings save request is in flight
- **WHEN** the operator edits another settings field and the in-flight save then fails
- **THEN** the Modbus Share save state becomes save-error
- **AND** the operator-visible operation error is populated
- **AND** a retry for the save operation is offered

#### Scenario: Editing during an in-flight save does not swallow the success

- **GIVEN** a settings save request is in flight
- **WHEN** the operator edits another settings field and the in-flight save then succeeds
- **THEN** the Modbus Share save state becomes the saved value

#### Scenario: The server-assigned revision is always adopted

- **GIVEN** a settings save request is in flight and the operator edits another settings field
- **WHEN** the response carrying a new settings revision returns
- **THEN** the settings revision and expected settings revision are updated from that response
- **AND** the next save request carries the adopted revision rather than a stale one
- **AND** the save state still converges to a terminal value

### Requirement: Deleted connectors do not reappear after later additions

The connector pool SHALL treat a deleted connector as permanently removed from the optimistic add bookkeeping. Adding a connector after an earlier connector was deleted SHALL NOT reintroduce the deleted connector into the displayed list.

#### Scenario: Add, delete, then add again

- **GIVEN** the operator adds connector A and then deletes connector A
- **WHEN** the operator adds connector B
- **THEN** the connector list contains only connector B
- **AND** connector A is not present in the list

##### Example: connector list after each step

| Step | Backend connectors | Displayed connectors |
| ---- | ------------------ | -------------------- |
| add A | A | A |
| delete A | none | none |
| add B | B | B |
