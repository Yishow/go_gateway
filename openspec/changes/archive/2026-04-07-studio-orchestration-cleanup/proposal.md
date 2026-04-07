## Why

Under umbrella `rule-centric-studio-refactor`, once `rule-driven-foundation`, `tag-review-first-flow`, `database-output-review-apply`, and `local-modbus-review-apply` are all defined, the last remaining problem is product orchestration drift. `/studio` must become the sole primary workflow, legacy `/datalink/workbench` behavior must collapse into redirects, and manual-first tools must stop competing with the rule-centric path as equal product entry points.

## What Changes

- Execute this change as child phase 5 under umbrella `rule-centric-studio-refactor`.
- Treat `rule-driven-foundation`, `tag-review-first-flow`, `database-output-review-apply`, and `local-modbus-review-apply` as prerequisites for this final cleanup phase.
- Make `/studio` the sole primary product route for the datalink workflow.
- Redirect legacy `/datalink/workbench` entry into `/studio` with preserved workflow context where possible.
- Demote drag-drop point/tag/mapping tooling to secondary engineering surfaces instead of the main product path.
- Remove remaining manual-first assumptions from the desktop workbench shell and guided workflow copy.
- **BREAKING**: legacy workbench entry points no longer remain as parallel primary workflow routes.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `datalink-ui`: Make the guided workflow explicitly `/studio`-primary and demote manual-first tooling from the normal path.
- `datalink-workbench-desktop`: Consolidate the shell and output route behavior onto `/studio` and remove legacy primary-route assumptions.

## Impact

- **Frontend**: `/studio` routing, navigation, workbench shell composition, and workflow copy.
- **UX**: primary-vs-secondary workflow positioning for manual engineering tools.
- **Tests**: route redirect coverage, guided workflow assertions, and shell consistency checks.
