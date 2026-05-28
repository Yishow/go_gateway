## MODIFIED Requirements

### Requirement: Placeholder step and settings surfaces

The v2 shell SHALL render placeholder content for each step or page that has not yet had its dedicated change applied. Once a step-specific or settings-specific change lands, the placeholder for that surface MUST be replaced with the real content delivered by that change. Placeholders MUST NOT make backend API calls, MUST NOT consume `useQuery` or `useMutation` hooks, and MUST NOT block the `onContinue` action for development purposes. After this change applies, Step 1 (`新增裝置`) MUST render its dedicated content defined by capability `datalink-workbench-v2-step1-device`, while Step 2, Step 3, Step 4, and Settings MUST continue to render placeholders until their corresponding changes apply.

#### Scenario: Step 2 placeholder remains during step1-device delivery

- **WHEN** the v2 page mounts with current step 2 after the `datalink-workbench-v2-step1-device` change has applied but before the Step 2 change has applied
- **THEN** the central step content displays text identifying it as Step 2 (`接入規則`)
- **AND** the placeholder text indicates that full Step 2 functionality is delivered by a subsequent change

#### Scenario: Settings placeholder remains during step1-device delivery

- **WHEN** the view state is `settings` after the `datalink-workbench-v2-step1-device` change has applied
- **THEN** the central content displays text identifying it as the Settings page
- **AND** the placeholder text indicates that full Settings functionality is delivered by a subsequent change

#### Scenario: No backend calls from any active placeholder

- **WHEN** the v2 page mounts and the operator navigates among Step 2, Step 3, Step 4, and Settings
- **THEN** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by any placeholder surface
