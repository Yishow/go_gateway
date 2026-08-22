import { describe, expect, it } from 'vitest';
import {
  WORKBENCH_EXPERIMENT_PHASES,
  WORKBENCH_EXPERIMENT_SHARED_RULES,
} from '@/pages/datalink/workbench/workbenchExperimentContract';

describe('workbench Tag compare contract', () => {
  it('defines the reopened Phase 3 shared matrix on the main contract surface', async () => {
    const contract = (await import('@/pages/datalink/workbench/workbenchExperimentContract')) as Record<string, unknown>;
    const scenarios = contract.WORKBENCH_TAG_COMPARE_SCENARIOS as
      | Array<{ id: string; label: string }>
      | undefined;
    const acceptance = contract.WORKBENCH_TAG_COMPARE_ACCEPTANCE as
      | readonly string[]
      | undefined;
    const archetypes = contract.WORKBENCH_TAG_COMPARE_ARCHETYPES as
      | Record<string, { label: string; interactionModel: string }>
      | undefined;
    const criticalTask = contract.WORKBENCH_TAG_COMPARE_CRITICAL_TASK as
      | { id: string; label: string; measures: readonly string[] }
      | undefined;

    expect(WORKBENCH_EXPERIMENT_SHARED_RULES).toContain(
      'compare 完成前不得進下一個 phase',
    );
    expect(scenarios).toEqual([
      { id: 'tag-review-candidates', label: 'review candidates queue' },
      { id: 'tag-diff-preview', label: 'inspect diff preview' },
      { id: 'tag-apply-decision', label: 'choose action and apply' },
      { id: 'tag-failure-retry-recovery', label: 'handle failure / retry / recovery' },
      { id: 'tag-handoff-output', label: 'handoff 到 Output' },
    ]);
    expect(WORKBENCH_EXPERIMENT_PHASES.find((phase) => phase.id === 'phase3')?.scenarioFocus).toEqual(
      scenarios?.map(({ id }) => id.replace(/^tag-/, '')),
    );
    expect(acceptance).toContain(
      'Tag phase shared review/apply logic and Output handoff semantics remain unchanged',
    );
    expect(acceptance).toContain(
      'v1 / v2 / v3 必須在 Tag phase 仍然 clearly distinct',
    );
    expect(archetypes).toEqual({
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
    });
    expect(criticalTask).toEqual({
      id: 'tag-review-to-output-handoff',
      label: '完成 Tag review/apply 並 handoff 到 Output',
      measures: [
        'time-to-open-review-queue',
        'time-to-confirm-diff-preview',
        'time-to-reach-output-handoff',
      ],
    });
  });
});
