## MODIFIED Requirements
### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps a persistent visual context of the complete data flow in one workspace: **Source Plan -> Memory Grid Allocation -> Tag Mapping -> Database Commit**.

The workflow SHALL ensure:

1. Operators can define source quantity and source data type before allocation.
2. Mapping and tag operations remain visible without page switches.
3. Validation and commit states are visible in the same workspace.
4. Failures are localized to a specific flow segment.

#### Scenario: End-to-end guided configuration in one workspace
- **WHEN** an operator defines source type/count and completes allocation, tag, and commit actions
- **THEN** the UI keeps source, grid, tag, and commit context visible together
- **AND** the operator does not need to switch pages to complete DB commit

#### Scenario: Segment-localized error handling
- **WHEN** validation or commit fails
- **THEN** the UI marks the failed segment
- **AND** provides actionable retry or edit guidance for that segment

### Requirement: Sidebar navigation improvements

The system SHALL provide sidebar navigation aligned to the one-screen operator workflow.

#### Scenario: Pipeline Studio as primary entry
- **WHEN** a user views the sidebar
- **THEN** the primary operation entry is the Pipeline Studio dashboard (`/datalink`)
- **AND** mapping operations are completed inside this entry

#### Scenario: Legacy mapping pages decommissioned
- **WHEN** a user accesses `/datalink/points`, `/datalink/mappings`, or `/datalink/wizard`
- **THEN** the system redirects to `/datalink`
- **AND** shows a migration notice that the workflow has moved to Pipeline Studio

## ADDED Requirements
### Requirement: Source template library

The UI SHALL allow operators to save, load, update, and delete source planning templates.

Each template SHALL store at least protocol context, source data type, source count, and naming defaults.

#### Scenario: Save and reuse source template
- **WHEN** an operator saves a source plan as a template
- **THEN** the template is persisted
- **AND** the operator can later load it to prefill planning controls

#### Scenario: Edit existing source template
- **WHEN** an operator updates a saved template
- **THEN** the updated defaults are used in subsequent planning sessions

### Requirement: Motion-guided operator flow

The UI SHALL use motion cues to guide stage transitions between source planning, grid allocation, tag linkage, and DB commit.

Animations MUST use short transitions (150-300ms) and MUST support reduced-motion preference.

#### Scenario: Guided transition after source planning
- **WHEN** an operator confirms source plan
- **THEN** the grid allocation region receives a transition cue indicating next action

#### Scenario: Reduced-motion mode
- **WHEN** user preference is `prefers-reduced-motion`
- **THEN** motion cues are replaced with static visual emphasis without animation

### Requirement: Planning intelligence and safety checks

The UI SHALL provide planning assistance and safety checks for one-screen operator execution.

#### Scenario: Naming preview and duplicate detection
- **WHEN** an operator defines batch naming rules
- **THEN** the UI previews generated names
- **AND** flags duplicates or naming conflicts before commit

#### Scenario: Two-stage validation execution
- **WHEN** an operator validates pending changes
- **THEN** the UI runs structural validation first
- **AND** runs executable validation only after structural validation succeeds

## REMOVED Requirements
### Requirement: Points management page
**Reason**: Points and mapping operations move to a single-screen Pipeline Studio workflow.
**Migration**: Redirect `/datalink/points` to `/datalink` with a one-time migration notice.
