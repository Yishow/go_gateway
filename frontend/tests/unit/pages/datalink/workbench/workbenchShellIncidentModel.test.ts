import { describe, expect, it } from 'vitest';
import { getWorkbenchShellIncidentModel } from '@/pages/datalink/workbench/workbenchShellIncidentModel';

describe('getWorkbenchShellIncidentModel', () => {
  it('returns Tag review as the recovery surface when Output lost the active rule handoff', () => {
    expect(
      getWorkbenchShellIncidentModel({
        hasSelectedDevice: true,
        activeStep: 'output',
        sourceReadiness: { status: 'ready' },
        tagReadiness: { status: 'ready' },
        outputReadiness: { status: 'partial', reason: 'output-rule-required' },
      }),
    ).toEqual({
      blockerKey: 'workbench.shell.blockers.outputRuleRequired',
      returnLabelKey: 'workbench.shell.actions.returnTag',
      returnStep: 'tag',
      tone: 'warning',
    });
  });

  it('keeps Output as the recovery surface when a rule exists but no target has been applied yet', () => {
    expect(
      getWorkbenchShellIncidentModel({
        hasSelectedDevice: true,
        activeStep: 'output',
        sourceReadiness: { status: 'ready' },
        tagReadiness: { status: 'ready' },
        outputReadiness: { status: 'partial', reason: 'output-not-applied' },
      }),
    ).toEqual({
      blockerKey: 'workbench.shell.blockers.outputPending',
      returnLabelKey: 'workbench.shell.actions.returnOutput',
      returnStep: 'output',
      tone: 'warning',
    });
  });
});
