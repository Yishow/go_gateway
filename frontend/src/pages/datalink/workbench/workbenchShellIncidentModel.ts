import type { StepReadinessState, WorkbenchReadiness, WorkbenchStep } from './workbenchTypes';

type WorkbenchShellIncidentTone = 'ok' | 'warning' | 'critical';

export type WorkbenchShellIncidentModel = {
  blockerKey: string;
  returnLabelKey: string;
  returnStep: WorkbenchStep;
  tone: WorkbenchShellIncidentTone;
};

export type WorkbenchShellReadinessStats = {
  readyCount: number;
  attentionCount: number;
  stepCount: number;
};

const READY_STATUSES: ReadonlySet<WorkbenchReadiness> = new Set(['ready', 'applied']);
const ATTENTION_STATUSES: ReadonlySet<WorkbenchReadiness> = new Set(['partial', 'blocked']);

export function getWorkbenchReadinessExplanationKey(
  state: StepReadinessState,
): string | null {
  switch (state.reason) {
    case 'no-device-selected':
      return 'workbench.shell.blockers.noDevice';
    case 'no-points':
    case 'no-source-points':
      return 'workbench.shell.blockers.sourceRequired';
    case 'no-tags-linked':
      return 'workbench.shell.blockers.tagRequired';
    case 'some-tags-blocked':
      return 'workbench.shell.blockers.tagAttention';
    case 'output-rule-required':
      return 'workbench.shell.blockers.outputRuleRequired';
    case 'output-not-applied':
      return 'workbench.shell.blockers.outputPending';
    default:
      return null;
  }
}

export function getWorkbenchShellReadinessStats(
  readinessStates: ReadonlyArray<StepReadinessState>,
): WorkbenchShellReadinessStats {
  return {
    readyCount: readinessStates.filter(({ status }) => READY_STATUSES.has(status)).length,
    attentionCount: readinessStates.filter(({ status }) => ATTENTION_STATUSES.has(status)).length,
    stepCount: readinessStates.length,
  };
}

export function getWorkbenchShellIncidentModel(input: {
  hasSelectedDevice: boolean;
  activeStep: WorkbenchStep;
  sourceReadiness: StepReadinessState;
  tagReadiness: StepReadinessState;
  outputReadiness: StepReadinessState;
}): WorkbenchShellIncidentModel {
  if (!input.hasSelectedDevice) {
    return {
      blockerKey: 'workbench.shell.blockers.noDevice',
      returnLabelKey: 'workbench.shell.actions.returnDevice',
      returnStep: 'device',
      tone: 'warning',
    };
  }

  if (
    input.sourceReadiness.reason === 'no-points' ||
    input.tagReadiness.reason === 'no-source-points' ||
    input.outputReadiness.reason === 'no-source-points'
  ) {
    return {
      blockerKey: 'workbench.shell.blockers.sourceRequired',
      returnLabelKey: 'workbench.shell.actions.returnSource',
      returnStep: 'source',
      tone: 'critical',
    };
  }

  if (
    input.tagReadiness.reason === 'no-tags-linked' ||
    input.outputReadiness.reason === 'no-tags-linked'
  ) {
    return {
      blockerKey: 'workbench.shell.blockers.tagRequired',
      returnLabelKey: 'workbench.shell.actions.returnTag',
      returnStep: 'tag',
      tone: 'warning',
    };
  }

  if (input.tagReadiness.reason === 'some-tags-blocked') {
    return {
      blockerKey: 'workbench.shell.blockers.tagAttention',
      returnLabelKey: 'workbench.shell.actions.returnTag',
      returnStep: 'tag',
      tone: 'warning',
    };
  }

  if (input.outputReadiness.reason === 'output-rule-required') {
    return {
      blockerKey: 'workbench.shell.blockers.outputRuleRequired',
      returnLabelKey: 'workbench.shell.actions.returnTag',
      returnStep: 'tag',
      tone: 'warning',
    };
  }

  if (input.outputReadiness.reason === 'output-not-applied') {
    return {
      blockerKey: 'workbench.shell.blockers.outputPending',
      returnLabelKey: 'workbench.shell.actions.returnOutput',
      returnStep: 'output',
      tone: 'warning',
    };
  }

  return {
    blockerKey: 'workbench.shell.blockers.stable',
    returnLabelKey: 'workbench.shell.actions.focusCurrent',
    returnStep: input.activeStep,
    tone: 'ok',
  };
}
