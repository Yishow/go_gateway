import { describe, expect, it } from 'vitest';
import {
  WORKBENCH_EXPERIMENT_TOKEN_SOURCES,
  workbenchExperimentTokens,
} from '@/styles/workbench-experiment-tokens';

describe('workbench experiment token contract', () => {
  it('keeps the approved design-source set in one shared contract', () => {
    expect(WORKBENCH_EXPERIMENT_TOKEN_SOURCES).toEqual([
      'linear.app',
      'sentry',
      'clickhouse',
    ]);
    expect(workbenchExperimentTokens.meta.phase1r).toBe('three-archetypes-reboot');
  });

  it('exposes distinct archetype palettes without version-only files', () => {
    expect(workbenchExperimentTokens.archetype.linear.accent).toBe('#7170ff');
    expect(workbenchExperimentTokens.archetype.sentry.highlight).toBe('#c2ef4e');
    expect(workbenchExperimentTokens.archetype.clickhouse.accent).toBe('#faff69');

    expect(
      new Set([
        workbenchExperimentTokens.archetype.linear.canvas,
        workbenchExperimentTokens.archetype.sentry.canvas,
        workbenchExperimentTokens.archetype.clickhouse.canvas,
      ]).size,
    ).toBe(3);
  });

  it('includes shared typography and treatment semantics for the reboot round', () => {
    expect(workbenchExperimentTokens.typography.weight.linearUi).toBe(510);
    expect(workbenchExperimentTokens.typography.weight.heavy).toBe(900);
    expect(workbenchExperimentTokens.treatment.glassPanel.backdropFilter).toContain('blur');
    expect(workbenchExperimentTokens.treatment.insetButton).toContain('inset');
    expect(workbenchExperimentTokens.treatment.neonGlow).toContain('rgba');
  });
});
