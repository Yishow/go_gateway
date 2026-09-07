## MODIFIED Requirements

### Requirement: Continue gate
The bottom action bar SHALL show total devices, the explicitly selected setup scope and successful current tests. At least one selected device and successful tests for every selected device MUST be required to continue. Unselected devices SHALL remain visible drafts, not be silently enabled or deleted. A missing selection SHALL preserve the legacy all-device scope. Actual activation SHALL still revalidate current device readiness and the existing workspace barrier.

#### Scenario: Disabled when not all tested
- **WHEN** two devices are selected and only the first has a successful current test
- **THEN** the continue button is disabled and the untested selected device and repair action are shown.

#### Scenario: Enabled when all tested
- **WHEN** every device in the selected scope has a successful current test
- **THEN** clicking continue invokes onContinue, completes Step 1 and advances to Step 2 for that scope.

#### Scenario: Offline device explicitly excluded
- **WHEN** a tested device remains selected and the user explicitly excludes an offline device
- **THEN** setup can continue for the tested device and the excluded device remains a clearly labelled draft.
