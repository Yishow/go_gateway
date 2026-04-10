export const WORKBENCH_OUTPUT_COMPARE_SCENARIOS = [
  { id: 'output-readiness', label: 'inspect readiness for Local Modbus + Database' },
  { id: 'output-dry-run', label: 'run dry-run before apply' },
  { id: 'output-apply', label: 'apply output bindings' },
  { id: 'output-blocker-diagnosis', label: 'diagnose blocker / retry / recovery' },
] as const;

export const WORKBENCH_OUTPUT_COMPARE_ACCEPTANCE = [
  'Output phase shared readiness / dry-run / apply / blocker semantics remain unchanged',
  'Local Modbus 與 Database 兩個 target families 都必須被 shared compare contract 覆蓋',
  'v1 / v2 / v3 必須在 Output phase 仍然 clearly distinct',
  'readiness、dry-run、apply、blocker diagnosis surface 都必須一起改變',
  '只有換 kit、加 scoped CSS、或只調整單一 target family 都不算完成',
  '若 review 仍認為三版除了 Device 幾乎一樣，Phase 4 視為不通過',
] as const;

export const WORKBENCH_OUTPUT_COMPARE_ARCHETYPES = {
  v1: {
    label: 'Linear Control Room',
    interactionModel: 'control-room output console with calmer blocker tracing across targets',
  },
  v2: {
    label: 'Sentry Incident Desk',
    interactionModel: 'incident-desk output command surface with target diagnosis and apply flow',
  },
  v3: {
    label: 'ClickHouse Data Cockpit',
    interactionModel: 'cockpit mapping/state surface spanning Local Modbus and Database',
  },
} as const;

export const WORKBENCH_OUTPUT_COMPARE_CRITICAL_TASK = {
  id: 'output-readiness-to-apply',
  label: '確認 Local Modbus / Database readiness、執行 dry-run，並完成 Output apply',
  measures: [
    'time-to-read-readiness',
    'time-to-confirm-dry-run',
    'time-to-complete-output-apply',
  ],
} as const;
