# source-rule-target-datatype Specification

## Purpose
Define the canonical contract for declaring a source rule’s protocol read type, optional target Tag or output data type, optional linear scaling, and the validation/runtime semantics that connect that intent to the mapping pipeline.

## Requirements

### Requirement: Operator MAY declare a target data type distinct from protocol read type
The system SHALL allow an optional **target data type** on a source rule that MAY differ from the rule’s protocol read (`point`) data type. When omitted or null, the system SHALL treat the target as identical to the protocol read type.

#### Scenario: Default matches protocol type
- **WHEN** an operator creates a rule without setting a separate target data type
- **THEN** the system SHALL persist the rule such that target equals protocol read type
- **AND** downstream Tag and Mapping behavior SHALL match today’s same-type, empty-pipeline behavior

#### Scenario: Explicit target differs from protocol read type
- **WHEN** an operator sets target data type to a value different from the protocol read data type
- **THEN** the system SHALL accept the rule only if the pair is supported by validated cast semantics
- **AND** SHALL persist both values on the rule record

---
### Requirement: Unsupported type pairs are rejected at rule save time
The system SHALL reject rule persistence when the chosen protocol-read and target type pair cannot be satisfied by a validated `cast` transform according to pipeline validation rules.

#### Scenario: Rejection surfaces actionable error
- **WHEN** an operator submits a rule with an unsupported type pair
- **THEN** the system SHALL reject the save with a validation error
- **AND** the error SHALL identify the incompatible pair without exposing internal secrets

---
### Requirement: Transformed values use the mapping pipeline at runtime
The system SHALL NOT introduce a parallel transformation path for rule-target typing. Final values written to tags and outputs SHALL continue to flow through the persisted mapping `transform_pipeline` executed by the runtime ingestor.

#### Scenario: Cast appears as mapping pipeline step
- **WHEN** a rule declares a target type that differs from the point read type
- **THEN** the derived mapping SHALL include a `cast` step consistent with the declared target type
- **AND** runtime execution SHALL use the same `ExecutePipeline` behavior as manually authored mappings

---
### Requirement: Supported example — uint16 read with float64 tag semantics
The system SHALL support the common case where the protocol read type is **`uint16`** (single-register raw integer) and the operator declares a target type of **`float64`** for the Tag, provided the `cast` is valid per pipeline rules.

#### Scenario: Cast from uint16 to float64
- **WHEN** a rule sets protocol read type `uint16` and target type `float64`
- **THEN** the system SHALL persist both intents and derive a mapping pipeline that includes a validated `cast` to `float64`
- **AND** the ingested raw value SHALL be interpreted first as the point read type before casting

---
### Requirement: Operator MAY declare linear scaling for engineering units
The system SHALL allow optional **scale** parameters on a source rule (for example multiplier and offset aligned with `TransformParamsScale`) when the operator needs **unit conversion** beyond type casting. When set, the derived mapping SHALL include a validated **`scale`** transform step in addition to any `cast` step, in an order defined by the implementation and validated by the pipeline.

#### Scenario: Scale applies after type alignment
- **WHEN** an operator sets non-default scale parameters on a rule that also declares a target type
- **THEN** the persisted mapping pipeline SHALL validate successfully
- **AND** runtime output SHALL reflect both cast and scaling as separate traceable steps

---
### Requirement: Database context MAY influence Source-rule defaults without overriding operator intent
The system SHALL allow the active `database` output target to influence Source-step defaults for target data type, scale, and naming guidance when the operator is planning a source rule.

These defaults SHALL remain advisory. The system SHALL NOT silently replace an explicit operator choice only because the output target is `database`.

#### Scenario: Database target adjusts planning defaults
- **WHEN** an operator opens or creates a source rule while the active output target is `database`
- **THEN** the Source-step form MAY prefill or highlight target data type, scale, or naming defaults that better support downstream database review
- **AND** the operator can still change those values before saving the rule

#### Scenario: Explicit operator value wins over database default
- **WHEN** an operator manually changes target data type, scale, or naming-related rule input after a database-aware default is suggested
- **THEN** the system preserves the operator’s explicit value on save
- **AND** the UI does not silently revert that choice to a database-oriented default

<!-- @trace
source: source-step-rule-first-database-preview
updated: 2026-05-09
code:
  - .github/prompts/spectra-ask.prompt.md
  - .github/prompts/spectra-archive.prompt.md
  - tests/shell/start-port-management.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/skills/spectra-ask/SKILL.md
  - .github/skills/spectra-debug/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - .github/skills/spectra-audit/SKILL.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - frontend/src/i18n/locales/zh-TW/common.json
  - tests/shell/start-backend-before-frontend.sh
  - .github/prompts/spectra-apply.prompt.md
  - CLAUDE.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-discuss.prompt.md
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/skills/spectra-discuss/SKILL.md
  - .github/prompts/opsx-propose.prompt.md
  - start.sh
  - .github/skills/spectra-ingest/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - .spectra.yaml
  - .github/skills/spectra-commit/SKILL.md
  - frontend/src/main.tsx
  - .github/prompts/opsx-apply.prompt.md
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - frontend/package.json
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/prompts/opsx-archive.prompt.md
  - frontend/src/App.tsx
  - .github/skills/spectra-propose/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - .github/prompts/spectra-commit.prompt.md
  - start.ps1
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - frontend/src/components/DevAgentation.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
  - tests/shell/start-backend-logfile.sh
  - .github/skills/spectra-archive/SKILL.md
  - .github/prompts/spectra-ingest.prompt.md
tests:
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
-->