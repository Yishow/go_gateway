import * as React from 'react';
import { mappingAPI } from '../../../../../services/datalink';
import type { MappingPreviewResponse } from '../../../../../types/datalink';
import type { RuntimeStreamConnectionState } from '../../../../../types/datalink';
import type { Mapping, Point } from '../../state/types';
import { normalizeTypedEnvelope, parseMappingPreviewResponse } from '../../../../../utils/safeJson';

export interface Step3LivePreviewState {
  status: 'idle' | 'connecting' | 'live' | 'reconnecting' | 'degraded' | 'error' | 'stale';
  data: MappingPreviewResponse | null;
  error: string | null;
  errorCode?: string;
  requestId?: string;
  retryable: boolean;
  retry: () => void;
  lastSuccessAt: string | null;
}

const PREVIEW_DEBOUNCE_MS = 250;
export const PREVIEW_RETRY_DELAYS_MS = [1000, 2000, 4000, 8000, 16000, 30000] as const;

export function useStep3LivePreview(
  point: Point | null,
  mapping: Mapping | null,
  rawValue: unknown | null,
  connectionState: RuntimeStreamConnectionState = 'connected',
  workspaceId?: string,
): Step3LivePreviewState {
  const [state, setState] = React.useState<Step3LivePreviewState>({
    status: 'idle',
    data: null,
    error: null,
    retryable: false,
    retry: () => undefined,
    lastSuccessAt: null,
  });
  const [retryNonce, setRetryNonce] = React.useState(0);
  const retry = React.useCallback(() => {
    setRetryNonce((current) => current + 1);
  }, []);
  const pointDataType = point?.data_type;
  const pointId = point?.id;
  const mappingPointId = mapping?.point_id;
  const mappingScale = mapping?.scale;
  const mappingOffset = mapping?.offset;
  const mappingTargetType = mapping?.target_type;
  const previewWorkspaceId = workspaceId ?? mapping?.workspace_id;
  const previewPayload = React.useMemo(
    () => pointDataType !== undefined && mappingPointId !== undefined && mappingScale !== undefined &&
      mappingOffset !== undefined && mappingTargetType !== undefined
      ? {
        transform_pipeline: [
          { type: 'decode', order: 1, params: { data_type: pointDataType } },
          { type: 'scale', order: 2, params: { scale: mappingScale, offset: mappingOffset } },
          { type: 'cast', order: 3, params: { target_type: mappingTargetType } },
        ],
      } as { transform_pipeline: NonNullable<MappingPreviewResponse['pipeline']> }
      : null,
    [mappingOffset, mappingPointId, mappingScale, mappingTargetType, pointDataType],
  );
  const previewKey = React.useMemo(
    () => pointId !== undefined && mappingPointId !== undefined && rawValue !== null && rawValue !== undefined && previewPayload
      ? JSON.stringify({
        point_id: pointId,
        raw_value: rawValue,
        transform_pipeline: previewPayload.transform_pipeline,
      })
      : null,
    [mappingPointId, pointId, previewPayload, rawValue],
  );
  const requestIdRef = React.useRef(0);
  const lastSuccessKeyRef = React.useRef<string | null>(null);
  const lastSuccessRef = React.useRef<Pick<Step3LivePreviewState, 'data' | 'lastSuccessAt'>>({
    data: null,
    lastSuccessAt: null,
  });

  React.useEffect(() => {
    if (!previewPayload || rawValue === null || rawValue === undefined) {
      lastSuccessKeyRef.current = null;
      lastSuccessRef.current = { data: null, lastSuccessAt: null };
      setState({
        status: 'idle',
        data: null,
        error: null,
        retryable: false,
        retry,
        lastSuccessAt: null,
      });
      return;
    }

    if (lastSuccessKeyRef.current !== previewKey) {
      lastSuccessKeyRef.current = previewKey;
      lastSuccessRef.current = { data: null, lastSuccessAt: null };
    }

    const payload = previewPayload;
    const currentRequestId = ++requestIdRef.current;
    let disposed = false;
    let retryIndex = 0;
    let retryTimer: number | null = null;
    const lastSuccess = lastSuccessRef.current;

    const requestPreview = () => {
      if (disposed || currentRequestId !== requestIdRef.current) {
        return;
      }

      setState((previous) => ({
        ...previous,
        status: lastSuccess.data ? 'reconnecting' : 'connecting',
        data: lastSuccess.data,
        error: null,
      }));

      void mappingAPI
        .preview({
          ...(previewWorkspaceId ? { workspace_id: previewWorkspaceId } : {}),
          raw_value: rawValue,
          transform_pipeline: payload!.transform_pipeline,
        })
        .then((data) => {
          if (disposed || currentRequestId !== requestIdRef.current) {
            return;
          }

          const safeData = parseMappingPreviewResponse(data);
          if (!safeData || safeData.error) {
            handleFailure();
            return;
          }

          const success = {
            data: safeData,
            lastSuccessAt: new Date().toISOString(),
          };
          lastSuccessRef.current = success;
          setState({
            status: 'live',
            data: safeData,
            error: null,
            retryable: false,
            retry,
            lastSuccessAt: success.lastSuccessAt,
          });
        })
        .catch((error: unknown) => {
          if (disposed || currentRequestId !== requestIdRef.current) {
            return;
          }

          handleFailure(extractErrorMetadata(error));
        });
    };

    const handleFailure = (metadata: ErrorMetadata = {}) => {
      if (disposed || currentRequestId !== requestIdRef.current) {
        return;
      }

      const hasLastSuccess = lastSuccess.data !== null && lastSuccess.lastSuccessAt !== null;
      const exhausted = retryIndex >= PREVIEW_RETRY_DELAYS_MS.length;
      setState({
        status: hasLastSuccess ? (exhausted ? 'degraded' : 'stale') : 'error',
        data: lastSuccess.data,
        error: 'preview_unavailable',
        errorCode: metadata.code ?? 'preview_unavailable',
        requestId: metadata.requestId,
        retryable: metadata.retryable ?? true,
        retry,
        lastSuccessAt: lastSuccess.lastSuccessAt,
      });

      if (!exhausted) {
        const delay = PREVIEW_RETRY_DELAYS_MS[retryIndex];
        retryIndex += 1;
        retryTimer = window.setTimeout(requestPreview, delay);
      }
    };

    setState(() => ({
      status: lastSuccess.data ? 'reconnecting' : 'connecting',
      data: lastSuccess.data,
      error: null,
      errorCode: undefined,
      requestId: undefined,
      retryable: false,
      retry,
      lastSuccessAt: lastSuccess.lastSuccessAt,
    }));

    const timer = window.setTimeout(() => {
      requestPreview();
    }, PREVIEW_DEBOUNCE_MS);

    return () => {
      disposed = true;
      window.clearTimeout(timer);
      if (retryTimer !== null) {
        window.clearTimeout(retryTimer);
      }
    };
  }, [previewKey, previewPayload, previewWorkspaceId, rawValue, retry, retryNonce]);

  const streamStatus = connectionState === 'error' || connectionState === 'disconnected'
    ? (state.data ? 'degraded' : 'error')
    : connectionState === 'stale' || connectionState === 'reconnecting'
      ? connectionState
      : connectionState === 'degraded'
        ? 'degraded'
        : null;

  return streamStatus ? { ...state, status: streamStatus, retry } : { ...state, retry };
}

interface ErrorMetadata {
  code?: string;
  requestId?: string;
  retryable?: boolean;
}

function extractErrorMetadata(error: unknown): ErrorMetadata {
  if (typeof error !== 'object' || error === null) {
    return {};
  }
  const record = error as Record<string, unknown>;
  const response = record.response as Record<string, unknown> | undefined;
  const data = response?.data as Record<string, unknown> | undefined;
  const payload = data?.error;
  const envelope = normalizeTypedEnvelope(payload);
  return {
    code: envelope.code,
    requestId: envelope.requestId,
    retryable: envelope.retryable,
  };
}
