import * as React from 'react';
import type { CommitLog } from '../../state/types';
import type {
  StudioV2ActivationRecovery,
  StudioV2ActivationResponse,
} from '../../../../../types/studioV2Activation';
import { StudioV2ActivationBarrierError } from '../../../../../services/studioV2WorkspaceActivation';
import { normalizeTypedEnvelope, parseStudioV2ActivationResponse } from '../../../../../utils/safeJson';

export type Step4ActivationPhase = 'idle' | 'activating' | 'done';

type ActivationOutcome = 'failed' | 'unconfirmed';

interface Step4ActivationState {
  phase: Step4ActivationPhase;
  response: StudioV2ActivationResponse | null;
  recovery: StudioV2ActivationRecovery | null;
  recoveryUnavailable: boolean;
}

function activationOutcomeOf(error: unknown): ActivationOutcome {
  if (typeof error === 'object' && error !== null && !Array.isArray(error)) {
    const outcome = (error as { outcome?: unknown }).outcome;
    if (outcome === 'failed' || outcome === 'unconfirmed') {
      return outcome;
    }
  }

  return normalizeTypedEnvelope(error).code ? 'failed' : 'unconfirmed';
}

export function useStep4Activation(
  activateWorkspace: (() => Promise<StudioV2ActivationResponse>) | undefined,
  recoverActivationStatus?: () => Promise<StudioV2ActivationRecovery>,
  workspaceId?: string,
) {
  const [state, setState] = React.useState<Step4ActivationState>({
    phase: 'idle',
    response: null,
    recovery: null,
    recoveryUnavailable: false,
  });
  const mountedRef = React.useRef(true);
  const generationRef = React.useRef(0);
  const inFlightRef = React.useRef(false);
  const workspaceIdRef = React.useRef(workspaceId);

  React.useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
      generationRef.current += 1;
    };
  }, []);

  React.useLayoutEffect(() => {
    if (workspaceIdRef.current === workspaceId) {
      return;
    }
    workspaceIdRef.current = workspaceId;
    generationRef.current += 1;
    setState((previous) => ({
      ...previous,
      phase: 'idle',
      response: null,
      recovery: null,
      recoveryUnavailable: false,
    }));
  }, [workspaceId]);

  React.useEffect(() => {
    if (!recoverActivationStatus) {
      return;
    }

    let active = true;
    const recoveryGeneration = generationRef.current;
    void recoverActivationStatus()
      .then((value) => {
        if (!active || !mountedRef.current || generationRef.current !== recoveryGeneration) {
          return;
        }
        const recovery = workspaceId && value.workspace_id !== workspaceId ? null : value;
        setState((previous) => ({
          ...previous,
          recovery,
          recoveryUnavailable: recovery === null,
        }));
      })
      .catch(() => {
        if (!active || !mountedRef.current || generationRef.current !== recoveryGeneration) {
          return;
        }
        setState((previous) => ({
          ...previous,
          recovery: null,
          recoveryUnavailable: true,
        }));
      });

    return () => {
      active = false;
    };
  }, [recoverActivationStatus, workspaceId]);

  const start = React.useCallback(async () => {
    if (!activateWorkspace || inFlightRef.current) {
      return;
    }

    inFlightRef.current = true;
    const attemptWorkspaceId = workspaceId;
    const attemptGeneration = ++generationRef.current;
    setState((previous) => ({ ...previous, phase: 'activating', response: null }));
    try {
      const rawResponse = await activateWorkspace();
      const response = parseStudioV2ActivationResponse(rawResponse);
      if (!response) {
        const envelope = normalizeTypedEnvelope(rawResponse);
        throw new StudioV2ActivationBarrierError(
          'activation_failed',
          false,
          envelope.action,
          envelope.requestId,
          'unconfirmed',
          envelope.operationId,
        );
      }
      if (workspaceId && response.workspace_id !== workspaceId) {
        throw new StudioV2ActivationBarrierError(
          'activation_failed',
          false,
          response.action,
          response.request_id,
          'unconfirmed',
          response.operation_id,
        );
      }
      if (!mountedRef.current || workspaceIdRef.current !== attemptWorkspaceId || generationRef.current !== attemptGeneration) {
        return;
      }
      setState((previous) => ({ ...previous, phase: 'done', response }));
    } catch (error) {
      const typed = normalizeTypedEnvelope(error);
      const outcome = activationOutcomeOf(error);
      const operationId = typed.operationId;
      if (!mountedRef.current || workspaceIdRef.current !== attemptWorkspaceId || generationRef.current !== attemptGeneration) {
        return;
      }
      setState((previous) => ({
        ...previous,
        phase: 'done',
        response: {
          workspace_id: '',
          results: [],
          code: typed.code ?? 'activation_failed',
          action: typed.action,
          request_id: typed.requestId,
          retryable: outcome === 'unconfirmed' ? false : typed.retryable ?? true,
          outcome,
          ...(operationId ? { operation_id: operationId } : {}),
        },
      }));
    } finally {
      inFlightRef.current = false;
    }
  }, [activateWorkspace, workspaceId]);

  const reset = React.useCallback(() => {
    generationRef.current += 1;
    setState((previous) => ({ ...previous, phase: 'idle', response: null }));
  }, []);

  const logs = React.useMemo<CommitLog[]>(() => {
    const response = state.response;
    if (!response) {
      return [];
    }
    if (response.outcome === 'unconfirmed' && response.results.length === 0) {
      return [{
        label: 'POST /studio-v2/workspace/activate',
        detail: '',
        status: 'pending' as const,
        code: response.code,
        retryable: false,
        action: response.action,
        request_id: response.request_id,
      }];
    }
    const resultLogs = response.results.map((result) => ({
      label: `POST /studio-v2/workspace/activate → ${result.device_id}`,
      detail: '',
      status: result.status,
      code: result.code,
      retryable: result.retryable,
      action: result.action,
      request_id: result.request_id,
    }));
    if (resultLogs.length === 0 && response.code) {
      return [{
        label: 'POST /studio-v2/workspace/activate',
        detail: '',
        status: 'failed' as const,
        code: response.code,
        retryable: response.retryable,
        action: response.action,
        request_id: response.request_id,
      }];
    }
    return resultLogs;
  }, [state.response]);

  const canContinue = Boolean(
    state.response?.outcome !== 'unconfirmed' &&
    state.response?.results.some((result) => result.status === 'success'),
  );
  return {
    ...state,
    logs,
    canContinue,
    start,
    reset,
  };
}
