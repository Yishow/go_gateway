import * as React from 'react';
import type { CommitLog } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import type { StudioV2ActivationResponse } from '../../../../../types/studioV2Activation';
import { normalizeTypedEnvelope, parseStudioV2ActivationResponse } from '../../../../../utils/safeJson';

export type Step4ActivationPhase = 'idle' | 'activating' | 'done';

export function useStep4Activation(
  activateWorkspace: (() => Promise<StudioV2ActivationResponse>) | undefined,
  dispatch: React.Dispatch<WorkbenchV2Action>,
) {
  const [state, setState] = React.useState<{
    phase: Step4ActivationPhase;
    response: StudioV2ActivationResponse | null;
  }>({ phase: 'idle', response: null });

  const start = React.useCallback(async () => {
    if (!activateWorkspace) {
      return;
    }

    setState({ phase: 'activating', response: null });
    try {
      const response = parseStudioV2ActivationResponse(await activateWorkspace());
      if (!response) throw new Error('activation response invalid');
      response.results.forEach((result) => {
        dispatch({
          type: 'updateDevice',
          deviceId: result.device_id,
          patch: result.status === 'success'
            ? { status: 'active', running: true }
            : { running: false },
        });
      });
      setState({ phase: 'done', response });
    } catch (error) {
      const typed = normalizeTypedEnvelope(error);
      setState({
        phase: 'done',
        response: {
          workspace_id: '',
          results: [],
          code: typed.code ?? 'activation_failed',
          action: typed.action,
          request_id: typed.requestId,
          retryable: typed.retryable ?? true,
        },
      });
    }
  }, [activateWorkspace, dispatch]);

  const reset = React.useCallback(() => {
    setState({ phase: 'idle', response: null });
  }, []);

  const logs = React.useMemo<CommitLog[]>(() => {
    const response = state.response;
    if (!response) {
      return [];
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

  const canContinue = Boolean(state.response?.results.some((result) => result.status === 'success'));
  return {
    ...state,
    logs,
    canContinue,
    start,
    reset,
  };
}
