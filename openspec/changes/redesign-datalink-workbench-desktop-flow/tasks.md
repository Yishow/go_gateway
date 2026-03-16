## 1. Workbench shell and shared state

- [ ] 1.1 Replace the current workbench desktop frame with the approved five-region shell (`StepRail`, `ContextBar`, `PrimaryWorkArea`, `InspectorPanel`, `BottomSummaryBar`).
- [ ] 1.2 Move page-level summary/readiness behavior out of `ActionDock` and into shared inspector/bottom-bar state.
- [ ] 1.3 Extend workbench state/selectors to support shared readiness, active output target, Step 2 rule metadata, and cross-step traceability.

## 2. Device workspace redesign

- [ ] 2.1 Build the new `DeviceWorkspace` browser surface with search, filters, create, refresh, and clone entry points.
- [ ] 2.2 Implement the Step 1 inspector and context-bar capability summary (`address base`, `word order`, `unit id`, protocol traits, recent test history).
- [ ] 2.3 Add the workbench-specific clone flow and cover Step 1 desktop interactions with tests.

## 3. Source workspace redesign

- [ ] 3.1 Introduce the formal Step 2 rule model and `RuleLayerBar` with enable/disable, lock, reorder, and focus behaviors.
- [ ] 3.2 Rebuild `AddressCanvas` to render a continuous 16-bit lattice with multi-rule merged spans and visible gaps.
- [ ] 3.3 Implement `Plan / Live / Link` view modes plus value-format switch, freeze/snapshot compare, jump-to-address, and coverage overview.
- [ ] 3.4 Add local template persistence for rule groups and capability-mismatch warnings when applying templates.
- [ ] 3.5 Add the secondary audit surface and Step 2 inspector flows for rule selection, span inspection, and batch point creation.

## 4. Tag binding workspace redesign

- [ ] 4.1 Rebuild Step 3 as a dense `TagBindingBoard` that exposes naming preview, source span, raw/transformed value, bit width, and binding state in the main surface.
- [ ] 4.2 Implement existing/new tag branching, filter/search, batch diff preview, and batch result summary flows.
- [ ] 4.3 Build the Step 3 inspector for full metadata editing, conflict detail, and single-item overrides.

## 5. Unified output workspace redesign

- [ ] 5.1 Create the shared `OutputCandidateBoard` and output readiness/filter model for Local Modbus and Database targets.
- [ ] 5.2 Build the Local Modbus studio with register map visualization, base/offset context, auto-map strategies, dry-run/test write, health summary, and sync-result feedback.
- [ ] 5.3 Build the Database studio with schema snapshot, type badges, required/missing highlighting, write-row preview, and sync-result feedback.
- [ ] 5.4 Build the Step 4 inspector trace for source -> tag -> output status and target-specific blocked reasons.

## 6. Quality, compatibility, and rollout

- [ ] 6.1 Update i18n, accessibility, and keyboard coverage for the redesigned shell and Step 1~4 workspaces.
- [ ] 6.2 Add or update desktop regression coverage for 1920×1080 layout stability and the new Step 2/3/4 information surfaces.
- [ ] 6.3 Decide and implement the compatibility strategy for `/datalink/local-modbus` and other legacy workbench entry points during rollout.
