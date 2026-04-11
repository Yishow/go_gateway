import type {
  SourceRuleCandidateSnapshotView,
  SourceRuleDatabaseOutputCandidateView,
  SourceRuleLocalModbusOutputCandidateView,
} from '../../../types/sourceRuleCandidates';
import type { StepReadinessState } from './workbenchTypes';

export type WorkbenchOutputMainlineState = {
  outputReady: boolean;
  outputReadiness: StepReadinessState;
  databaseConfigured: boolean;
  localModbusConfigured: boolean;
};

export function hasAppliedDatabaseOutput(
  candidates: ReadonlyArray<SourceRuleDatabaseOutputCandidateView>,
): boolean {
  return candidates.some((candidate) => Boolean(candidate.mapping_id));
}

export function hasAppliedLocalModbusOutput(
  candidates: ReadonlyArray<SourceRuleLocalModbusOutputCandidateView>,
): boolean {
  return candidates.some((candidate) => typeof candidate.register === 'number');
}

export function buildWorkbenchOutputMainlineState(input: {
  hasSelectedDevice: boolean;
  sourceReady: boolean;
  tagReady: boolean;
  activeRuleId: string | null;
  snapshot?: SourceRuleCandidateSnapshotView | null;
}): WorkbenchOutputMainlineState {
  if (!input.hasSelectedDevice) {
    return {
      outputReady: false,
      outputReadiness: { status: 'draft', reason: 'no-device-selected' },
      databaseConfigured: false,
      localModbusConfigured: false,
    };
  }

  if (!input.sourceReady) {
    return {
      outputReady: false,
      outputReadiness: { status: 'blocked', reason: 'no-source-points' },
      databaseConfigured: false,
      localModbusConfigured: false,
    };
  }

  if (!input.tagReady) {
    return {
      outputReady: false,
      outputReadiness: { status: 'blocked', reason: 'no-tags-linked' },
      databaseConfigured: false,
      localModbusConfigured: false,
    };
  }

  if (!input.activeRuleId) {
    return {
      outputReady: false,
      outputReadiness: { status: 'partial', reason: 'output-rule-required' },
      databaseConfigured: false,
      localModbusConfigured: false,
    };
  }

  const databaseConfigured = hasAppliedDatabaseOutput(
    input.snapshot?.database_outputs.candidates ?? [],
  );
  const localModbusConfigured = hasAppliedLocalModbusOutput(
    input.snapshot?.local_modbus_outputs.candidates ?? [],
  );
  const outputReady = databaseConfigured || localModbusConfigured;

  return {
    outputReady,
    outputReadiness: outputReady
      ? { status: 'ready' }
      : { status: 'partial', reason: 'output-not-applied' },
    databaseConfigured,
    localModbusConfigured,
  };
}
