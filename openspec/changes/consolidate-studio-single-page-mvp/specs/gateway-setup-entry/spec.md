## ADDED Requirements

### Requirement: Setup entry page presents additive route choices
The system SHALL provide a setup entry page at `/gateway/entry` that helps operators choose between the existing setup workflows without removing or replacing them.

The entry page SHALL expose at least these route choices:
1. A primary `quick complete setup` action that launches the existing quick setup workflow.
2. A direct entry to the existing `/studio` workbench.
3. A direct entry to the existing expert workbench.

#### Scenario: Entry page shows the available setup paths
- **WHEN** an operator opens `/gateway/entry`
- **THEN** the system shows a setup entry page
- **AND** the page exposes route choices for quick complete setup, direct workbench access, and expert workbench access

#### Scenario: Quick complete setup uses the existing quick flow
- **WHEN** an operator chooses the primary quick complete setup action
- **THEN** the system routes the operator to `/gateway/quick-setup`
- **AND** the operator continues with the existing quick setup workflow instead of a duplicate form

#### Scenario: Existing advanced routes remain visible from the entry page
- **WHEN** an operator wants the existing full or advanced flow
- **THEN** the entry page provides direct links to `/studio` and `/gateway/expert-workbench`
- **AND** those links do not warn that the old routes were removed or deprecated

### Requirement: Entry page rollout does not gate existing workflows
The system SHALL treat the setup entry page as an additive front door.

Any rollout flag or route guard for `/gateway/entry` SHALL NOT disable direct access to `/studio`, `/gateway/quick-setup`, or `/gateway/expert-workbench`.

#### Scenario: Entry page flag is disabled
- **WHEN** the rollout flag for `/gateway/entry` is disabled
- **THEN** `/gateway/entry` falls back to an existing supported flow
- **AND** `/studio`, `/gateway/quick-setup`, and `/gateway/expert-workbench` remain directly accessible

#### Scenario: Operator bypasses the entry page
- **WHEN** an operator opens `/studio`, `/gateway/quick-setup`, or `/gateway/expert-workbench` directly
- **THEN** the selected page loads normally
- **AND** the workflow does not depend on visiting `/gateway/entry` first
