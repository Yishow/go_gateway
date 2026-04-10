export const WORKBENCH_SHELL_COMPARE_SCENARIOS = [
  { id: 'shell-readiness-summary', label: 'show readiness summary and correct next action' },
  { id: 'shell-active-blocker', label: 'surface global blocker and shortest valid return path' },
  {
    id: 'shell-diagnostics-refresh',
    label: 'refresh diagnostics and handle failure / retry / recovery',
  },
  { id: 'shell-return-to-mainline', label: 'take the shortest return action to the blocked step' },
] as const;

export const WORKBENCH_SHELL_COMPARE_ACCEPTANCE = [
  'Shell shared readiness / blocker / diagnostics / return semantics remain unchanged',
  'shell 只可擁有 readiness summary、active blocker summary、diagnostics refresh status、shortest return action',
  'shell 不得接管 step-local edit / mutation / validation ownership',
  'v1 / v2 / v3 必須在 Shell / diagnostics phase 仍然 clearly distinct',
  '若 review 仍認為三版除了 Device 幾乎一樣，Phase 5 視為不通過',
] as const;

export const WORKBENCH_SHELL_COMPARE_ARCHETYPES = {
  v1: {
    label: 'Linear Control Room',
    interactionModel: 'control-room shell summary with calmer blocker framing and lower interference',
  },
  v2: {
    label: 'Sentry Incident Desk',
    interactionModel: 'incident-desk shell diagnostics with global blocker and repair hop',
  },
  v3: {
    label: 'ClickHouse Data Cockpit',
    interactionModel: 'cockpit summary and alert shell without step-local workflow takeover',
  },
} as const;

export const WORKBENCH_SHELL_COMPARE_CRITICAL_TASK = {
  id: 'shell-blocker-to-mainline',
  label: '從 shell 讀取 readiness / blocker、刷新 diagnostics，並回到最短主線修復',
  measures: [
    'time-to-read-shell-status',
    'time-to-refresh-diagnostics',
    'time-to-return-to-blocked-step',
  ],
} as const;
