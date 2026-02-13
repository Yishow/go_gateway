import { useCallback, useMemo, useReducer } from 'react';

export type FlowStatus = 'draft' | 'validated' | 'active' | 'error';
export type FlowSegment = 'source' | 'grid' | 'tag' | 'sink';

export interface FlowDiagnostic {
  latestValue: string;
  quality: 'good' | 'warning' | 'bad' | 'unknown';
  timestamp: string;
  error: string;
}

export interface FlowState {
  status: FlowStatus;
  sourceDeviceId: string;
  sourceAddress: string;
  pointId: string;
  tagId: string;
  sinkTarget: string;
  diagnostics: Record<FlowSegment, FlowDiagnostic>;
}

type FlowAction =
  | { type: 'set_source'; payload: { sourceDeviceId: string; sourceAddress: string; pointId: string } }
  | { type: 'set_tag'; payload: { tagId: string } }
  | { type: 'set_sink'; payload: { sinkTarget: string } }
  | { type: 'set_diagnostics'; payload: Partial<Record<FlowSegment, Partial<FlowDiagnostic>>> }
  | { type: 'mark_validated' }
  | { type: 'mark_active' }
  | { type: 'mark_error'; payload: { segment: FlowSegment; message: string } }
  | { type: 'reset_draft' };

const EMPTY_DIAGNOSTIC: FlowDiagnostic = {
  latestValue: '-',
  quality: 'unknown',
  timestamp: '-',
  error: '',
};

export const INITIAL_FLOW_STATE: FlowState = {
  status: 'draft',
  sourceDeviceId: '',
  sourceAddress: '',
  pointId: '',
  tagId: '',
  sinkTarget: 'timeseries',
  diagnostics: {
    source: { ...EMPTY_DIAGNOSTIC },
    grid: { ...EMPTY_DIAGNOSTIC },
    tag: { ...EMPTY_DIAGNOSTIC },
    sink: { ...EMPTY_DIAGNOSTIC },
  },
};

export function flowReducer(state: FlowState, action: FlowAction): FlowState {
  switch (action.type) {
    case 'set_source':
      return {
        ...state,
        status: 'draft',
        sourceDeviceId: action.payload.sourceDeviceId,
        sourceAddress: action.payload.sourceAddress,
        pointId: action.payload.pointId,
      };
    case 'set_tag':
      return {
        ...state,
        status: 'draft',
        tagId: action.payload.tagId,
      };
    case 'set_sink':
      return {
        ...state,
        sinkTarget: action.payload.sinkTarget,
      };
    case 'set_diagnostics': {
      const nextDiagnostics = { ...state.diagnostics };
      (Object.keys(action.payload) as FlowSegment[]).forEach((segment) => {
        const partial = action.payload[segment];
        if (!partial) return;
        nextDiagnostics[segment] = { ...nextDiagnostics[segment], ...partial };
      });
      return { ...state, diagnostics: nextDiagnostics };
    }
    case 'mark_validated':
      if (!state.sourceDeviceId || !state.sourceAddress || !state.tagId) return state;
      return { ...state, status: 'validated' };
    case 'mark_active':
      if (state.status !== 'validated') return state;
      return { ...state, status: 'active' };
    case 'mark_error':
      return {
        ...state,
        status: 'error',
        diagnostics: {
          ...state.diagnostics,
          [action.payload.segment]: {
            ...state.diagnostics[action.payload.segment],
            quality: 'bad',
            error: action.payload.message,
          },
        },
      };
    case 'reset_draft':
      return {
        ...state,
        status: 'draft',
      };
    default:
      return state;
  }
}

export function useFlowLifecycle(initialState: FlowState = INITIAL_FLOW_STATE) {
  const [state, dispatch] = useReducer(flowReducer, initialState);

  const canValidate = useMemo(
    () => Boolean(state.sourceDeviceId && state.sourceAddress && state.tagId),
    [state.sourceAddress, state.sourceDeviceId, state.tagId]
  );

  const canActivate = useMemo(() => state.status === 'validated', [state.status]);
  const hasError = useMemo(() => state.status === 'error', [state.status]);

  const setSource = useCallback(
    (sourceDeviceId: string, sourceAddress: string, pointId: string) =>
      dispatch({ type: 'set_source', payload: { sourceDeviceId, sourceAddress, pointId } }),
    []
  );
  const setTag = useCallback((tagId: string) => dispatch({ type: 'set_tag', payload: { tagId } }), []);
  const setSink = useCallback((sinkTarget: string) => dispatch({ type: 'set_sink', payload: { sinkTarget } }), []);
  const setDiagnostics = useCallback(
    (payload: Partial<Record<FlowSegment, Partial<FlowDiagnostic>>>) =>
      dispatch({ type: 'set_diagnostics', payload }),
    []
  );
  const markValidated = useCallback(() => dispatch({ type: 'mark_validated' }), []);
  const markActive = useCallback(() => dispatch({ type: 'mark_active' }), []);
  const markError = useCallback(
    (segment: FlowSegment, message: string) =>
      dispatch({ type: 'mark_error', payload: { segment, message } }),
    []
  );
  const resetDraft = useCallback(() => dispatch({ type: 'reset_draft' }), []);

  return useMemo(
    () => ({
      state,
      canValidate,
      canActivate,
      hasError,
      setSource,
      setTag,
      setSink,
      setDiagnostics,
      markValidated,
      markActive,
      markError,
      resetDraft,
    }),
    [
      canActivate,
      canValidate,
      hasError,
      markActive,
      markError,
      markValidated,
      resetDraft,
      setDiagnostics,
      setSink,
      setSource,
      setTag,
      state,
    ]
  );
}
