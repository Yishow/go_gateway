## REMOVED Requirements

### Requirement: Placeholder step and settings surfaces

**Reason**: After the `datalink-workbench-v2-settings` change applies, every step and the settings page have their dedicated content delivered by `datalink-workbench-v2-step1-device`, `datalink-workbench-v2-step2-rule`, `datalink-workbench-v2-step3-mapping`, `datalink-workbench-v2-step4-database`, and `datalink-workbench-v2-settings` respectively. No placeholder surface remains in the v2 shell after this change, so the placeholder requirement no longer applies.

**Migration**: Any future surface that should temporarily render placeholder content during a multi-change rollout SHALL define its own scoped requirement under the change that introduces it, rather than relying on the generic shell-wide placeholder. The placeholder TSX files (`Step1DevicePlaceholder.tsx`, `Step2RulePlaceholder.tsx`, `Step3MappingPlaceholder.tsx`, `Step4DatabasePlaceholder.tsx`, `SettingsPlaceholder.tsx`) remain in the repository as rollback-only references and SHALL NOT be imported by the live shell.
