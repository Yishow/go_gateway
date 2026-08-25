import { modbusShareAPI } from '../../../../services/datalink';
import type {
  ModbusShareDesiredMapping,
  ModbusShareReconcileOutcome,
  ModbusShareCanonicalPlan,
} from '../../../../types/datalink';
import type { Rule, WorkbenchV2State } from './types';
import { normalizeTypedEnvelope } from '../../../../utils/safeJson';

export interface StudioV2ShareActivationContext {
  workspace_id: string;
  workspace_revision: string;
  settings_revision: string;
  readiness_token: string;
  configured_enabled?: boolean;
  /** Backend-produced canonical plan; the browser must not derive cursor geometry. */
  canonical_plan?: ModbusShareCanonicalPlan;
}

export interface ModbusShareProjectionAPI {
  reconcile: (request: {
    workspace_id: string;
    expected_workspace_revision: string;
    expected_settings_revision: string;
    readiness_token: string;
    canonical_plan_signature?: string;
    desired_mappings: ModbusShareDesiredMapping[];
  }) => Promise<ModbusShareReconcileOutcome>;
}

type SharePlanningState = Pick<WorkbenchV2State, 'rules' | 'settings'>;

export function getPersistedShareBaseRegister(rules: Rule[]): number | null {
  const registers = rules
    .filter((rule) => rule.persisted === true && rule.share_enabled && rule.share_start_register !== null)
    .map((rule) => rule.share_start_register as number)
    .filter((register) => Number.isInteger(register));
  return registers.length > 0 ? Math.min(...registers) : null;
}

export class StudioV2ShareProjectionError extends Error {
  readonly code: string;
  readonly retryable: boolean;
  readonly action?: string;
  readonly request_id?: string;
  readonly outcome?: ModbusShareReconcileOutcome['outcome'];
  readonly invalidated_count?: number;
  readonly invalidated_spans?: ModbusShareReconcileOutcome['invalidated_spans'];
  readonly diagnostics?: ModbusShareReconcileOutcome['diagnostics'];

  constructor(
    code: string,
    retryable: boolean,
    action?: string,
    requestId?: string,
    details?: Pick<ModbusShareReconcileOutcome, 'outcome' | 'invalidated_count' | 'invalidated_spans' | 'diagnostics'>,
  ) {
    super(code);
    this.name = 'StudioV2ShareProjectionError';
    this.code = code;
    this.retryable = retryable;
    this.action = action;
    this.request_id = requestId;
    this.outcome = details?.outcome;
    this.invalidated_count = details?.invalidated_count;
    this.invalidated_spans = details?.invalidated_spans;
    this.diagnostics = details?.diagnostics;
  }
}

function projectionErrorFromUnknown(error: unknown): StudioV2ShareProjectionError {
  if (error instanceof StudioV2ShareProjectionError) {
    return error;
  }
  if (typeof error === 'object' && error !== null && 'response' in error) {
    const response = error.response;
    if (typeof response === 'object' && response !== null && 'data' in response) {
      const data = response.data;
      if (typeof data === 'object' && data !== null && 'error' in data) {
        const payload = data.error;
        if (typeof payload === 'object' && payload !== null) {
          const envelope = normalizeTypedEnvelope(payload);
          return new StudioV2ShareProjectionError(envelope.code ?? 'modbus_share_reconcile_failed', envelope.retryable ?? true, envelope.action, envelope.requestId);
        }
      }
    }
  }
  return new StudioV2ShareProjectionError('modbus_share_reconcile_failed', true);
}

/** Reconciles the complete desired Share set through the backend projection seam. */
export async function syncStudioV2ShareMappings(
  state: SharePlanningState,
  context?: StudioV2ShareActivationContext,
  projectionAPI: ModbusShareProjectionAPI = modbusShareAPI,
): Promise<ModbusShareReconcileOutcome | undefined> {
  if (!state.settings.modbus_share.enabled || context?.configured_enabled === false) {
    return undefined;
  }
  if (!context?.canonical_plan) {
    throw new StudioV2ShareProjectionError('modbus_share_projection_required', true);
  }
  const plan = context.canonical_plan;
  if (
    plan.workspace_id !== context.workspace_id ||
    plan.workspace_revision !== context.workspace_revision ||
    plan.settings_revision !== context.settings_revision ||
    !plan.signature
  ) {
    throw new StudioV2ShareProjectionError('modbus_share_projection_required', true);
  }
  if (!context || !context.workspace_id || !context.workspace_revision || !context.settings_revision || !context.readiness_token) {
    throw new StudioV2ShareProjectionError('modbus_share_hydration_required', true);
  }

  let outcome: ModbusShareReconcileOutcome;
  try {
    outcome = await projectionAPI.reconcile({
      workspace_id: context.workspace_id,
      expected_workspace_revision: context.workspace_revision,
      expected_settings_revision: context.settings_revision,
      readiness_token: context.readiness_token,
      canonical_plan_signature: plan.signature,
      desired_mappings: plan.desired_mappings,
    });
  } catch (error) {
    throw projectionErrorFromUnknown(error);
  }
  if (outcome.outcome === 'failed' || outcome.outcome === 'dirty_unknown' || outcome.outcome === 'invalidated_unknown') {
    const diagnostic = outcome.diagnostics?.find((item) => item.severity === 'error') ?? outcome.diagnostics?.[0];
    throw new StudioV2ShareProjectionError(
      diagnostic?.code ?? 'modbus_share_reconcile_failed',
      diagnostic?.retryable ?? true,
      diagnostic?.action,
      diagnostic?.request_id,
      {
        outcome: outcome.outcome,
        invalidated_count: outcome.invalidated_count,
        invalidated_spans: outcome.invalidated_spans,
        diagnostics: outcome.diagnostics,
      },
    );
  }
  return outcome;
}

/** Projects Share mappings before activation so synchronization failures stop the workspace transition. */
export async function activateStudioV2WorkspaceWithShare<T>(
  state: SharePlanningState,
  activateWorkspace: (projection?: ModbusShareReconcileOutcome) => Promise<T>,
  context?: StudioV2ShareActivationContext,
  projectionAPI: ModbusShareProjectionAPI = modbusShareAPI,
): Promise<T> {
  const projection = await syncStudioV2ShareMappings(state, context, projectionAPI);
  return activateWorkspace(projection);
}
