## ADDED Requirements

### Requirement: Destination locks are scoped and explained

The editor SHALL decide destination field locks from one rule: a per-point (mixed) recording SHALL keep its schema and table editable until its first save, choosable from the destination catalog or entered manually, while its storage kind SHALL stay managed; saved per-point, profile and long-table recordings SHALL keep storage kind, schema, table and layout fixed. Catalog tables owned by another recording SHALL be unavailable for selection, and the server ownership checks SHALL remain authoritative. Every locked destination field SHALL show a plain-language reason next to it, and a saved locked recording SHALL offer the existing copy-as-new-recording action.

#### Scenario: New per-point recording chooses its table

- **WHEN** an operator creates a per-point recording and has not saved it
- **THEN** the schema and table can be selected from the catalog or typed, and the managed-only storage kind is shown with its reason

#### Scenario: Saved long-table recording explains the lock

- **WHEN** an operator opens a saved long-table recording
- **THEN** storage kind, schema, table and layout are disabled, a sentence states that long-table recordings share one table per schema and are fixed after saving, and copy-as-new is offered

### Requirement: Recording method is visible where it is named

For a per-point recording, the recording method section SHALL list the members grouped by calculation and interval with point counts, and each group SHALL offer an action that scrolls to and focuses the first matching member's calculation control in the source table. For a recording without a recording profile, the section SHALL state "latest (legacy wide table)" with its saved row interval and SHALL offer creating a new recording prefilled with the same points and destination connector; it SHALL NOT convert the recording in place.

#### Scenario: Grouped per-point methods

| Members | Shown groups |
| ------- | ------------ |
| 3 latest every 15 s, 1 mean every 60 s | "latest · 15 s · 3 points", "mean · 60 s · 1 point" |

- **WHEN** the operator activates the "mean · 60 s" group's action
- **THEN** the source table scrolls to that member and its calculation control receives focus

#### Scenario: Legacy recording is summarized, not converted

- **WHEN** a recording without a recording profile and a 60 second interval is opened
- **THEN** the method section shows latest (legacy wide table) every 60 seconds and a shortcut to create a new recording, and the saved recording is unchanged

### Requirement: Instantaneous recordings edit their reduction explicitly

For a recording whose profile kind is instantaneous, the editor SHALL offer the reduction (latest or arithmetic mean) and the recording interval within the existing instantaneous bounds. Choosing mean SHALL require the existing gauge-use confirmation. Saving SHALL go through the existing save with revision checks, and the preparation area SHALL then require a fresh schema check and confirmation before the change takes effect. Cumulative recordings SHALL NOT offer mean.

#### Scenario: Switching to mean

- **GIVEN** a saved instantaneous recording using latest every 15 seconds
- **WHEN** the operator selects mean, confirms gauge use and saves
- **THEN** the update carries `sample_reduction` mean, a new revision is saved, and the preparation area asks for a schema check before recording uses it

### Requirement: Hidden and pending areas state why

For a per-point recording, the advanced row-rule area SHALL explain that intervals are set per point, incomplete rows are fixed to partial results, and managed-table duplicate protection, entity key and provenance columns are system-defined. For an unsaved managed recording, the preparation area SHALL state that the table structure can be checked after saving. When the destination catalog is not being read, the editor SHALL NOT show a loading message; when the read fails, it SHALL state that the destination database could not be reached and offer retry.

#### Scenario: Unsaved managed recording

- **WHEN** a new managed recording has not been saved
- **THEN** the preparation area shows that the structure check becomes available after saving

#### Scenario: Unreachable destination

- **WHEN** the destination catalog request fails
- **THEN** the editor shows that the destination database could not be reached with a retry action instead of a loading message

### Requirement: Readiness issues are shown in plain language

The recording editor and the runtime dashboard SHALL present readiness issues with a localized explanation keyed by issue code, covering at least `schema-confirmation-required`, `write-group-apply-required`, `database-target-missing`, `source-revision-stale` and `schema-check-blocked`; unknown codes SHALL use a generic localized explanation. The raw code and server message SHALL appear only inside diagnostics, together with the localized explanation. The recording member search input SHALL disable browser autofill.

#### Scenario: Schema confirmation required

- **WHEN** a saved managed recording reports `schema-confirmation-required`
- **THEN** the issue reads as a localized instruction to check and confirm the table structure before starting, and the code with its English message appears only in diagnostics
