import { describe, expect, it } from 'vitest';
import {
  WORKBENCH_EXPERIMENT_PHASES,
  WORKBENCH_EXPERIMENT_SHARED_RULES,
} from '@/pages/datalink/workbench/workbenchExperimentContract';

describe('workbench shell compare contract', () => {
  it('defines the reopened Phase 5 shared matrix on the main contract surface', async () => {
    const contract = (await import('@/pages/datalink/workbench/workbenchExperimentContract')) as Record<string, unknown>;
    const scenarios = contract.WORKBENCH_SHELL_COMPARE_SCENARIOS as
      | Array<{ id: string; label: string }>
      | undefined;
    const acceptance = contract.WORKBENCH_SHELL_COMPARE_ACCEPTANCE as
      | readonly string[]
      | undefined;
    const archetypes = contract.WORKBENCH_SHELL_COMPARE_ARCHETYPES as
      | Record<string, { label: string; interactionModel: string }>
      | undefined;
    const criticalTask = contract.WORKBENCH_SHELL_COMPARE_CRITICAL_TASK as
      | { id: string; label: string; measures: readonly string[] }
      | undefined;

    expect(WORKBENCH_EXPERIMENT_SHARED_RULES).toContain(
      'compare 完成前不得進下一個 phase',
    );
    expect(scenarios).toEqual([
      { id: 'shell-readiness-summary', label: 'show readiness summary and correct next action' },
      { id: 'shell-active-blocker', label: 'surface global blocker and shortest valid return path' },
      {
        id: 'shell-diagnostics-refresh',
        label: 'refresh diagnostics and handle failure / retry / recovery',
      },
      { id: 'shell-return-to-mainline', label: 'take the shortest return action to the blocked step' },
    ]);
    expect(WORKBENCH_EXPERIMENT_PHASES.find((phase) => phase.id === 'phase5')?.scenarioFocus).toEqual(
      scenarios?.map(({ id }) => id.replace(/^shell-/, '')),
    );
    expect(acceptance).toContain(
      'Shell shared readiness / blocker / diagnostics / return semantics remain unchanged',
    );
    expect(acceptance).toContain(
      'shell 只可擁有 readiness summary、active blocker summary、diagnostics refresh status、shortest return action',
    );
    expect(acceptance).toContain(
      'shell 不得接管 step-local edit / mutation / validation ownership',
    );
    expect(archetypes).toEqual({
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
    });
    expect(criticalTask).toEqual({
      id: 'shell-blocker-to-mainline',
      label: '從 shell 讀取 readiness / blocker、刷新 diagnostics，並回到最短主線修復',
      measures: [
        'time-to-read-shell-status',
        'time-to-refresh-diagnostics',
        'time-to-return-to-blocked-step',
      ],
    });
  });
});
