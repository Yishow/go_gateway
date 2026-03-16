export const WORKBENCH_STEPS = ['device', 'source', 'tag', 'output'] as const;

export type WorkbenchStep = (typeof WORKBENCH_STEPS)[number];

export type WorkbenchStepMeta = {
  labelKey: string;
  placeholderKey: string;
};

export const WORKBENCH_STEP_META: Record<WorkbenchStep, WorkbenchStepMeta> = {
  device: {
    labelKey: 'workbench.steps.device',
    placeholderKey: 'workbench.placeholders.device',
  },
  source: {
    labelKey: 'workbench.steps.source',
    placeholderKey: 'workbench.placeholders.source',
  },
  tag: {
    labelKey: 'workbench.steps.tag',
    placeholderKey: 'workbench.placeholders.tag',
  },
  output: {
    labelKey: 'workbench.steps.output',
    placeholderKey: 'workbench.placeholders.output',
  },
};
