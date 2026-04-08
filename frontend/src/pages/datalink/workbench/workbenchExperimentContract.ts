import type { WorkbenchStep } from './workbenchTypes';
export {
  WORKBENCH_SOURCE_COMPARE_ACCEPTANCE,
  WORKBENCH_SOURCE_COMPARE_ARCHETYPES,
  WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK,
  WORKBENCH_SOURCE_COMPARE_SCENARIOS,
} from './workbenchSourceCompareContract';
export {
  WORKBENCH_OUTPUT_COMPARE_ACCEPTANCE,
  WORKBENCH_OUTPUT_COMPARE_ARCHETYPES,
  WORKBENCH_OUTPUT_COMPARE_CRITICAL_TASK,
  WORKBENCH_OUTPUT_COMPARE_SCENARIOS,
} from './workbenchOutputCompareContract';
export {
  WORKBENCH_SHELL_COMPARE_ACCEPTANCE,
  WORKBENCH_SHELL_COMPARE_ARCHETYPES,
  WORKBENCH_SHELL_COMPARE_CRITICAL_TASK,
  WORKBENCH_SHELL_COMPARE_SCENARIOS,
} from './workbenchShellCompareContract';
export {
  WORKBENCH_TAG_COMPARE_ACCEPTANCE,
  WORKBENCH_TAG_COMPARE_ARCHETYPES,
  WORKBENCH_TAG_COMPARE_CRITICAL_TASK,
  WORKBENCH_TAG_COMPARE_SCENARIOS,
} from './workbenchTagCompareContract';

export const WORKBENCH_EXPERIMENT_ROUTE = '/studio' as const;

export const WORKBENCH_EXPERIMENT_SHARED_RULES = [
  'baseline 與三個新版本共用同一個 /studio route',
  'baseline 與三個新版本共用同一組真實 API',
  '三個新版本共用同一組 shared semantic tokens',
  'compare 完成前不得進下一個 phase',
  '禁止 mock、禁止版本專屬後端契約、禁止平行產品 route',
] as const;

export const WORKBENCH_EXPERIMENT_VARIANTS = {
  baseline: {
    label: 'Current Studio',
    role: 'baseline',
    kit: 'current',
    branch: 'woe-base-current-ui',
    worktree: '.worktrees/woe-base-current-ui',
    port: 4173,
  },
  v1: {
    label: 'shadcn/Radix',
    role: 'variant',
    kit: 'shadcn-radix',
    branch: 'woe-v1-radix',
    worktree: '.worktrees/woe-v1-radix',
    port: 4174,
  },
  v2: {
    label: 'MUI',
    role: 'variant',
    kit: 'mui',
    branch: 'woe-v2-mui',
    worktree: '.worktrees/woe-v2-mui',
    port: 4175,
  },
  v3: {
    label: 'Ant Design',
    role: 'variant',
    kit: 'antd',
    branch: 'woe-v3-antd',
    worktree: '.worktrees/woe-v3-antd',
    port: 4176,
  },
} as const;

export const WORKBENCH_EXPERIMENT_COMPARE_METRICS = [
  {
    id: 'operator-flow',
    label: '操作順暢度',
    description: '點擊、輸入、確認與切換是否精簡且不中斷。',
  },
  {
    id: 'logic-clarity',
    label: '邏輯清晰度',
    description: '畫面結構、狀態關係與下一步是否直觀。',
  },
  {
    id: 'system-completeness',
    label: '對系統的完整性',
    description: '是否完整承接 Device / Source / Tag / Output 主線。',
  },
  {
    id: 'first-screen-density',
    label: '首屏資訊密度',
    description: '首屏是否能看見必要狀態、上下文與下一步。',
  },
  {
    id: 'critical-time',
    label: '關鍵操作時間',
    description: '完成指定任務的耗時與往返次數。',
  },
  {
    id: 'implementation-risk',
    label: '實作 / 維護風險',
    description: '元件複雜度、檔案膨脹風險與後續維護成本。',
  },
] as const;

type WorkbenchExperimentPhase = {
  id: string;
  task: string;
  label: string;
  step?: WorkbenchStep;
  scenarioFocus: readonly string[];
};

export const WORKBENCH_EXPERIMENT_PHASES: readonly WorkbenchExperimentPhase[] = [
  {
    id: 'phase0',
    task: '0.x',
    label: 'baseline、tokens 與工作樹隔離',
    scenarioFocus: ['shared-route', 'shared-api', 'shared-tokens', 'compare-gate'],
  },
  {
    id: 'phase1',
    task: '1.1',
    label: 'Device',
    step: 'device',
    scenarioFocus: ['create', 'edit', 'clone', 'connect', 'probe', 'diagnostics'],
  },
  {
    id: 'phase2',
    task: '2.1',
    label: 'Source',
    step: 'source',
    scenarioFocus: [
      'create-rule',
      'apply-template',
      'plan-live-link',
      'stale-preview-recovery',
      'handoff-tag',
    ],
  },
  {
    id: 'phase3',
    task: '3.1',
    label: 'Tag',
    step: 'tag',
    scenarioFocus: [
      'review-candidates',
      'diff-preview',
      'apply-decision',
      'failure-retry-recovery',
      'handoff-output',
    ],
  },
  {
    id: 'phase4',
    task: '4.1',
    label: 'Output',
    step: 'output',
    scenarioFocus: ['readiness', 'dry-run', 'apply', 'blocker-diagnosis'],
  },
  {
    id: 'phase5',
    task: '5.1',
    label: 'Shell / Diagnostics',
    scenarioFocus: [
      'readiness-summary',
      'active-blocker',
      'diagnostics-refresh',
      'return-to-mainline',
    ],
  },
  {
    id: 'phase6',
    task: '6.1',
    label: 'Final recommendation',
    scenarioFocus: ['end-to-end-run', 'full-compare', 'winner-recommendation'],
  },
] as const;

export const WORKBENCH_DEVICE_COMPARE_SCENARIOS = [
  {
    id: 'device-create',
    label: '建立設備',
    goal: '從 Device step 建立新的 PLC device，保留在同一個 /studio 流程內。',
    actions: ['open-create', 'fill-identity', 'fill-connection', 'save-device'],
    successSignals: ['device-created', 'device-listed', 'editor-state-preserved'],
  },
  {
    id: 'device-edit',
    label: '編輯設備',
    goal: '修改既有設備的連線欄位，確認 edit path 不會遺失現有設定。',
    actions: ['open-edit', 'inspect-current-config', 'change-connection', 'save-device'],
    successSignals: ['existing-config-loaded', 'mutation-submitted', 'updated-state-visible'],
  },
  {
    id: 'device-clone',
    label: '複製設備',
    goal: '從既有設備 clone 出新設備，並確認名稱衝突規則與初始值承接。',
    actions: ['open-clone', 'adjust-name', 'verify-carried-config', 'save-device'],
    successSignals: ['clone-mode-visible', 'distinct-name-enforced', 'cloned-config-persisted'],
  },
  {
    id: 'device-connect',
    label: '連線檢查',
    goal: '觸發 connect 段落診斷，確認結果區與可行動訊息清楚。',
    actions: ['run-connect', 'observe-connect-stage', 'inspect-result-panel'],
    successSignals: ['connect-stage-visible', 'latency-or-status-visible', 'actionable-feedback-visible'],
  },
  {
    id: 'device-probe',
    label: '探測檢查',
    goal: '觸發 probe 段落診斷，確認與 connect 結果明確分離。',
    actions: ['run-probe', 'observe-probe-stage', 'inspect-result-panel'],
    successSignals: ['probe-stage-visible', 'connect-and-probe-separated', 'capability-feedback-visible'],
  },
  {
    id: 'device-diagnostics',
    label: '錯誤診斷',
    goal: '在失敗條件下看見可行動的 blocker 與下一步。',
    actions: ['trigger-failure', 'read-diagnostic-copy', 'find-recovery-action'],
    successSignals: ['blocking-reason-visible', 'next-step-clear', 'no-silent-failure'],
  },
] as const;

export const WORKBENCH_DEVICE_PHASE1R_REBOOT = {
  task: '1R.1',
  label: 'Device 三 archetypes reboot',
  reason: '上一輪 Device compare 被直接 review 判定為差異不足且過度接近 baseline。',
  domainLayers: ['selection', 'task', 'diagnostics', 'system'],
} as const;

export const WORKBENCH_DEVICE_PHASE1R_ARCHETYPES = {
  v1: {
    label: 'Linear Control Room',
    sourceMapping: ['linear.app', 'sentry-badges', 'clickhouse-stats'],
    interactionModel: 'device rail / task canvas / live inspector',
    visualGoal: 'near-black low-noise surface with indigo emphasis',
  },
  v2: {
    label: 'Sentry Incident Desk',
    sourceMapping: ['sentry', 'clickhouse-diagnostics', 'linear-structure'],
    interactionModel: 'incident feed / diagnostics desk / system intel',
    visualGoal: 'warm purple-black command center with uppercase labels and lime/coral emphasis',
  },
  v3: {
    label: 'ClickHouse Data Cockpit',
    sourceMapping: ['clickhouse', 'linear-structure', 'sentry-alert-semantics'],
    interactionModel: 'KPI strip / fleet grid / inline cockpit expansion',
    visualGoal: 'pure-black high-density console with neon throughput signals',
  },
} as const;

export const WORKBENCH_DEVICE_COMPARE_ACCEPTANCE = [
  'create / edit / clone / connect / probe / diagnostics 六種動作都必須留在同一個 /studio workflow 中',
  'Selection / Task / Diagnostics / System 四層必須在 Device surface 可辨識，不可再次被折疊成單一 editor 心智模型',
  'v1 / v2 / v3 必須形成不同 archetype，而不是同一個畫面只換 kit 或小幅重排',
  '整個 Device phase 首屏必須一眼就能辨識成不同產品語言，不得只是沿用 baseline 骨架再換色、換字重或換少量元件',
  'Device phase 的工具列、主工作區、檢查區與摘要節奏都必須跟 archetype 一起重寫；若主要 chrome 仍與目前 Studio 八分像，視為不通過',
  'connect 與 probe 的結果必須可分辨，不得把兩者折疊成單一模糊成功/失敗訊號',
  'connect / probe 必須是第一公民，不可再次埋進單一 diagnostics blob 或只有 editor 內結果可見',
  'Device form 必須保留既有 protocol-specific connection fields 與 draft preservation',
  '畫面必須能讓操作員快速看出目前是在 create、edit 還是 clone 模式',
  '失敗訊息必須帶可行動診斷，不可只有抽象錯誤字串或 silent failure',
  'compare 時四版都必須使用同一組真實 API 與同一批 device data',
  '其餘 workbench pages 即使暫不改邏輯，也必須以同 archetype 的設計語言重新表達，避免 Device 很像另一個產品、其餘頁面卻仍是 baseline 皮膚',
] as const;

export const WORKBENCH_DEVICE_COMPARE_CRITICAL_TASK = {
  id: 'device-connect-probe-pass',
  label: '建立或選取設備後完成 connect / probe 診斷判讀',
  measures: ['time-to-open-editor', 'time-to-submit-check', 'time-to-understand-result'],
} as const;
