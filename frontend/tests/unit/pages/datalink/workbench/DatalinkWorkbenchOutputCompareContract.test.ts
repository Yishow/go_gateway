import { describe, expect, it } from 'vitest';
import {
  WORKBENCH_EXPERIMENT_PHASES,
  WORKBENCH_EXPERIMENT_SHARED_RULES,
} from '@/pages/datalink/workbench/workbenchExperimentContract';

describe('workbench Output compare contract', () => {
  it('defines the reopened Phase 4 shared matrix on the main contract surface', async () => {
    const contract = (await import('@/pages/datalink/workbench/workbenchExperimentContract')) as Record<string, unknown>;
    const scenarios = contract.WORKBENCH_OUTPUT_COMPARE_SCENARIOS as
      | Array<{ id: string; label: string }>
      | undefined;
    const acceptance = contract.WORKBENCH_OUTPUT_COMPARE_ACCEPTANCE as
      | readonly string[]
      | undefined;
    const archetypes = contract.WORKBENCH_OUTPUT_COMPARE_ARCHETYPES as
      | Record<string, { label: string; interactionModel: string }>
      | undefined;
    const criticalTask = contract.WORKBENCH_OUTPUT_COMPARE_CRITICAL_TASK as
      | { id: string; label: string; measures: readonly string[] }
      | undefined;

    expect(WORKBENCH_EXPERIMENT_SHARED_RULES).toContain(
      'compare 完成前不得進下一個 phase',
    );
    expect(scenarios).toEqual([
      { id: 'output-readiness', label: 'inspect readiness for Local Modbus + Database' },
      { id: 'output-dry-run', label: 'run dry-run before apply' },
      { id: 'output-apply', label: 'apply output bindings' },
      { id: 'output-blocker-diagnosis', label: 'diagnose blocker / retry / recovery' },
    ]);
    expect(WORKBENCH_EXPERIMENT_PHASES.find((phase) => phase.id === 'phase4')?.scenarioFocus).toEqual(
      scenarios?.map(({ id }) => id.replace(/^output-/, '')),
    );
    expect(acceptance).toContain(
      'Output phase shared readiness / dry-run / apply / blocker semantics remain unchanged',
    );
    expect(acceptance).toContain(
      'Local Modbus 與 Database 兩個 target families 都必須被 shared compare contract 覆蓋',
    );
    expect(acceptance).toContain(
      'v1 / v2 / v3 必須在 Output phase 仍然 clearly distinct',
    );
    expect(archetypes).toEqual({
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
    });
    expect(criticalTask).toEqual({
      id: 'output-readiness-to-apply',
      label: '確認 Local Modbus / Database readiness、執行 dry-run，並完成 Output apply',
      measures: [
        'time-to-read-readiness',
        'time-to-confirm-dry-run',
        'time-to-complete-output-apply',
      ],
    });
  });
});
