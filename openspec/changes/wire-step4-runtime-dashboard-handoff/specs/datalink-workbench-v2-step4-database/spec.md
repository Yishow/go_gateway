## MODIFIED Requirements

### Requirement: Commit completion card

After `state.committed === true`, the right column SHALL render an emerald-tone success card containing: a circular check icon, a heading equivalent to `設定已套用 · 開始收集資料`, a subline equivalent to `Scheduler 已啟動 · 第一筆資料預計在 ~{write_interval_seconds}s 後寫入`, and a secondary button `前往 Runtime Dashboard` that invokes the `onCommit` callback. The CommitSummary and CommitProgress views MUST be hidden when the success card is shown.

When the runtime dashboard handoff callback is wired by the shell, clicking the secondary button SHALL navigate to the post-setup runtime dashboard route. If the shell can resolve a single handoff device from the current workbench state, the destination SHALL include `device_id=<resolved-id>`. If the shell cannot resolve a single handoff device, the destination SHALL fall back to the runtime dashboard route without `device_id`.

#### Scenario: Success card visible

- **GIVEN** `state.committed === true` and `state.commit.status === 'success'`
- **WHEN** the right column renders
- **THEN** the emerald success card is present
- **AND** CommitSummary (the 5-row summary) is NOT in the DOM
- **AND** the secondary button text equals `前往 Runtime Dashboard` (or i18n equivalent)
- **AND** clicking the button invokes `onCommit` exactly once

#### Scenario: Subline reflects write_interval_seconds

- **GIVEN** the success card is visible and `state.db.connector.write_interval_seconds === 10`
- **WHEN** the subline renders
- **THEN** the subline text contains `~10s`

#### Scenario: Handoff includes resolved device id

- **GIVEN** the success card is visible
- **AND** the shell resolves `device_id=d-1` from the current workbench state
- **WHEN** the operator clicks `前往 Runtime Dashboard`
- **THEN** the system navigates to `/studio/runtime?device_id=d-1`

#### Scenario: Handoff falls back when device cannot be resolved

- **GIVEN** the success card is visible
- **AND** the shell cannot resolve a single handoff device from the current workbench state
- **WHEN** the operator clicks `前往 Runtime Dashboard`
- **THEN** the system navigates to `/studio/runtime`
- **AND** the handoff does not degrade to a no-op or console-only side effect
