import type { WorkbenchStep } from './workbenchTypes';

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
    scenarioFocus: ['create-rule', 'apply-template', 'edit-range', 'plan-live-link'],
  },
  {
    id: 'phase3',
    task: '3.1',
    label: 'Tag',
    step: 'tag',
    scenarioFocus: ['review-candidates', 'batch-actions', 'diff-preview', 'apply-feedback'],
  },
  {
    id: 'phase4',
    task: '4.1',
    label: 'Output',
    step: 'output',
    scenarioFocus: ['mapping-visibility', 'dry-run', 'apply-output', 'blocker-diagnosis'],
  },
  {
    id: 'phase5',
    task: '5.1',
    label: 'Shell / Diagnostics',
    scenarioFocus: ['readiness', 'global-blockers', 'repair-hop', 'return-to-flow'],
  },
  {
    id: 'phase6',
    task: '6.1',
    label: 'Final recommendation',
    scenarioFocus: ['end-to-end-run', 'full-compare', 'winner-recommendation'],
  },
] as const;
