import type { DeviceStatus, TagStatus } from '../../../types/datalink';

// ---------------------------------------------------------------------------
// Step definitions
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Readiness model  (spec §8.2)
// ---------------------------------------------------------------------------

export type WorkbenchReadiness = 'draft' | 'ready' | 'partial' | 'blocked' | 'applied';

export type StepReadinessState = {
  status: WorkbenchReadiness;
  reason?: string;
};

/**
 * Map domain DeviceStatus → WorkbenchReadiness as defined in spec §8.2.
 */
export function deviceStatusToReadiness(status: DeviceStatus): WorkbenchReadiness {
  switch (status) {
    case 'active':
      return 'ready';
    case 'draft':
      return 'draft';
    case 'disabled':
      return 'blocked';
  }
}

/**
 * Map domain TagStatus + output binding state → WorkbenchReadiness (spec §8.2).
 */
export function tagStatusToReadiness(
  status: TagStatus,
  hasOutputMapping: boolean,
): WorkbenchReadiness {
  switch (status) {
    case 'active':
      return hasOutputMapping ? 'ready' : 'partial';
    case 'draft':
      return 'draft';
    case 'retired':
      return 'blocked';
  }
}

// ---------------------------------------------------------------------------
// Output target
// ---------------------------------------------------------------------------

export type OutputTarget = 'modbus' | 'database';

// ---------------------------------------------------------------------------
// Inspector selection context  (spec §8.1)
// ---------------------------------------------------------------------------

export type InspectorSelection =
  | { kind: 'none' }
  | { kind: 'device'; deviceId: string }
  | { kind: 'rule'; ruleId: string }
  | { kind: 'span'; spanAddress: string; ruleId?: string }
  | { kind: 'tag'; tagId: string; pointId?: string }
  | { kind: 'outputCandidate'; tagId: string; target?: OutputTarget };

export const INSPECTOR_SELECTION_NONE: InspectorSelection = { kind: 'none' } as const;

// ---------------------------------------------------------------------------
// Cross-step traceability context  (spec §8.1)
// ---------------------------------------------------------------------------
//
// Unlike `inspectorSelection` (which clears on every step navigation),
// `WorkbenchCrossStepContext` deliberately survives step transitions so that
// later workspaces can surface the originating source context.
//
// Example flows:
//   Step 2 → Step 3 : focusedRuleId lets the tag-binding board know which
//                     source rule the operator was working on, enabling rule-
//                     scoped filtering and pre-selection in the binding board.
//   Step 3 → Step 4 : focusedTagIds lets the output workspace pre-select the
//                     same tags that were active in the tag-binding step.

export type WorkbenchCrossStepContext = {
  /** The rule the operator was focused on in Step 2 when transitioning to Step 3. */
  focusedRuleId: string | null;
  /**
   * Tag IDs that were selected or most recently acted on in Step 3 and should
   * be highlighted in Step 4 (output workspace) for continuity.
   */
  focusedTagIds: ReadonlyArray<string>;
};

export const WORKBENCH_CROSS_STEP_CONTEXT_INITIAL: WorkbenchCrossStepContext = {
  focusedRuleId: null,
  focusedTagIds: [],
} as const;

export type DevicePanelState =
  | { mode: 'create' }
  | { mode: 'edit'; deviceId: string }
  | { mode: 'clone'; sourceDeviceId: string }
  | null;

export type DeviceTestHistoryEntry = {
  id: string;
  testedAt: string;
  success: boolean;
  message: string;
  latencyMs: number | null;
};
