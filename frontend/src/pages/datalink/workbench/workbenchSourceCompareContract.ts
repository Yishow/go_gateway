export const WORKBENCH_SOURCE_COMPARE_SCENARIOS = [
  { id: 'source-create-rule', label: '建立 / 編輯規則' },
  { id: 'source-apply-template', label: '套用 template' },
  { id: 'source-plan-live-link', label: '切換 plan / live / link' },
  { id: 'source-stale-preview-recovery', label: '處理 stale preview / retry / recovery' },
  { id: 'source-handoff-tag', label: 'handoff 到 Tag' },
] as const;

export const WORKBENCH_SOURCE_COMPARE_ACCEPTANCE = [
  'Source phase shared logic / handoff semantics remain unchanged',
  'v1 / v2 / v3 必須在 Source phase 仍然 clearly distinct',
  'action placement、primary work surface、preview framing、visual language 都必須一起改變',
  '只有換 kit、加 scoped CSS、或只調整密度都不算完成',
  '若 review 仍認為三版除了 Device 幾乎一樣，Phase 2 視為不通過',
] as const;

export const WORKBENCH_SOURCE_COMPARE_ARCHETYPES = {
  v1: {
    label: 'Linear Control Room',
    interactionModel: 'control-room editing skeleton with calmer task zoning',
  },
  v2: {
    label: 'Sentry Incident Desk',
    interactionModel: 'incident-desk command surface with blocker-first flow',
  },
  v3: {
    label: 'ClickHouse Data Cockpit',
    interactionModel: 'cockpit banner / KPI / board composition',
  },
} as const;

export const WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK = {
  id: 'source-edit-to-tag-handoff',
  label: '完成 Source rule 編輯並 handoff 到 Tag',
  measures: [
    'time-to-open-rule-worksurface',
    'time-to-confirm-preview-scope',
    'time-to-reach-tag-handoff',
  ],
} as const;
