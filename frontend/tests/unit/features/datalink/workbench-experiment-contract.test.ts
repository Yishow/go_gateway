import { describe, expect, it } from 'vitest';
import * as workbenchExperimentContract from '@/pages/datalink/workbench/workbenchExperimentContract';

describe('workbench source compare shared contract', () => {
  it('defines the reopened source scenario matrix', () => {
    expect(workbenchExperimentContract).toHaveProperty('WORKBENCH_SOURCE_COMPARE_SCENARIOS');

    const scenarios = Reflect.get(
      workbenchExperimentContract,
      'WORKBENCH_SOURCE_COMPARE_SCENARIOS',
    ) as ReadonlyArray<{ id: string }>;

    expect(scenarios.map((scenario) => scenario.id)).toEqual([
      'source-create-rule',
      'source-apply-template',
      'source-plan-live-link',
      'source-stale-preview-recovery',
      'source-handoff-tag',
    ]);

    const phase2 = workbenchExperimentContract.WORKBENCH_EXPERIMENT_PHASES.find(
      (phase) => phase.id === 'phase2',
    );

    expect(phase2?.scenarioFocus).toEqual([
      'create-rule',
      'apply-template',
      'plan-live-link',
      'stale-preview-recovery',
      'handoff-tag',
    ]);
  });

  it('keeps shared logic invariant while forcing archetype-level differences', () => {
    expect(workbenchExperimentContract).toHaveProperty('WORKBENCH_SOURCE_COMPARE_ACCEPTANCE');

    const acceptance = Reflect.get(
      workbenchExperimentContract,
      'WORKBENCH_SOURCE_COMPARE_ACCEPTANCE',
    ) as readonly string[];

    expect(acceptance).toContain('Source phase shared logic / handoff semantics remain unchanged');
    expect(acceptance).toContain('v1 / v2 / v3 必須在 Source phase 仍然 clearly distinct');
    expect(acceptance).toContain('若 review 仍認為三版除了 Device 幾乎一樣，Phase 2 視為不通過');
  });

  it('captures the three source archetypes explicitly', () => {
    expect(workbenchExperimentContract).toHaveProperty('WORKBENCH_SOURCE_COMPARE_ARCHETYPES');

    const archetypes = Reflect.get(
      workbenchExperimentContract,
      'WORKBENCH_SOURCE_COMPARE_ARCHETYPES',
    ) as Record<string, { label: string; interactionModel: string }>;

    expect(archetypes.v1.label).toBe('Linear Control Room');
    expect(archetypes.v1.interactionModel).toContain('control-room');
    expect(archetypes.v2.label).toBe('Sentry Incident Desk');
    expect(archetypes.v2.interactionModel).toContain('incident-desk');
    expect(archetypes.v3.label).toBe('ClickHouse Data Cockpit');
    expect(archetypes.v3.interactionModel).toContain('cockpit');
  });

  it('records the critical source handoff task for compare timing', () => {
    expect(workbenchExperimentContract).toHaveProperty('WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK');

    const criticalTask = Reflect.get(
      workbenchExperimentContract,
      'WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK',
    ) as { id: string; measures: readonly string[] };

    expect(criticalTask.id).toBe('source-edit-to-tag-handoff');
    expect(criticalTask.measures).toEqual([
      'time-to-open-rule-worksurface',
      'time-to-confirm-preview-scope',
      'time-to-reach-tag-handoff',
    ]);
  });
});
