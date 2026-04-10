export const WORKBENCH_TAG_COMPARE_SCENARIOS = [
  { id: 'tag-review-candidates', label: 'review candidates queue' },
  { id: 'tag-diff-preview', label: 'inspect diff preview' },
  { id: 'tag-apply-decision', label: 'choose action and apply' },
  { id: 'tag-failure-retry-recovery', label: 'handle failure / retry / recovery' },
  { id: 'tag-handoff-output', label: 'handoff 到 Output' },
] as const;

export const WORKBENCH_TAG_COMPARE_ACCEPTANCE = [
  'Tag phase shared review/apply logic and Output handoff semantics remain unchanged',
  'v1 / v2 / v3 必須在 Tag phase 仍然 clearly distinct',
  'review queue、diff preview、batch decision、apply / retry / recovery surface 都必須一起改變',
  '只有換 kit、加 scoped CSS、或只調整密度都不算完成',
  '若 review 仍認為三版除了 Device 幾乎一樣，Phase 3 視為不通過',
] as const;

export const WORKBENCH_TAG_COMPARE_ARCHETYPES = {
  v1: {
    label: 'Linear Control Room',
    interactionModel: 'control-room review skeleton optimized for long batch sessions',
  },
  v2: {
    label: 'Sentry Incident Desk',
    interactionModel: 'incident-desk review command deck with blocker-first apply flow',
  },
  v3: {
    label: 'ClickHouse Data Cockpit',
    interactionModel: 'cockpit board / summary density with additive apply surfaces',
  },
} as const;

export const WORKBENCH_TAG_COMPARE_CRITICAL_TASK = {
  id: 'tag-review-to-output-handoff',
  label: '完成 Tag review/apply 並 handoff 到 Output',
  measures: [
    'time-to-open-review-queue',
    'time-to-confirm-diff-preview',
    'time-to-reach-output-handoff',
  ],
} as const;
