## MODIFIED Requirements

### Requirement: Placeholder step and settings surfaces

The v2 shell SHALL render placeholder content for each step or page that has not yet had its dedicated change applied. Once a step-specific or settings-specific change lands, the placeholder for that surface MUST be replaced with the real content delivered by that change. Placeholders MUST NOT make backend API calls, MUST NOT consume `useQuery` or `useMutation` hooks, and MUST NOT block the `onContinue` action for development purposes. After this change applies, Steps 1, 2, 3, and 4 MUST render their dedicated content (capabilities `datalink-workbench-v2-step1-device`, `datalink-workbench-v2-step2-rule`, `datalink-workbench-v2-step3-mapping`, `datalink-workbench-v2-step4-database`); Settings MUST continue to render a placeholder until the `datalink-workbench-v2-settings` change applies.

#### Scenario: Settings placeholder remains during step4-database delivery

- **WHEN** the view state is `settings` after the `datalink-workbench-v2-step4-database` change has applied
- **THEN** the central content displays text identifying it as the Settings page
- **AND** the placeholder text indicates that full Settings functionality is delivered by a subsequent change

#### Scenario: No backend calls from the settings placeholder

- **WHEN** the v2 page mounts and the operator navigates to Settings
- **THEN** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the placeholder

---

## ADDED Requirements

### Requirement: Scheduler indicator reflects committed state

The shell top bar SHALL render a scheduler status indicator whose tone reflects `state.committed`. When `state.committed === false` the indicator MUST render with an amber dot and the text `scheduler idle` (or i18n equivalent). When `state.committed === true` the indicator MUST render with an emerald dot, a `pulse-dot` animation, and the text `scheduler running` (or i18n equivalent). The transition MUST occur on the same render tick as the reducer applies the `completeCommit` action.

#### Scenario: Idle indicator before commit

- **WHEN** the v2 page mounts with `state.committed === false`
- **THEN** the top bar shows `scheduler idle` text
- **AND** the indicator dot uses an amber tone with no `pulse-dot` class

#### Scenario: Running indicator after commit

- **WHEN** the Step 4 commit animation completes and `state.committed` becomes `true`
- **THEN** the top bar shows `scheduler running` text
- **AND** the indicator dot uses an emerald tone with the `pulse-dot` animation class

#### Scenario: Indicator persists across step navigation

- **GIVEN** `state.committed === true`
- **WHEN** the operator navigates from Step 4 to Step 1 via the StepRail
- **THEN** the top bar continues to show `scheduler running`
