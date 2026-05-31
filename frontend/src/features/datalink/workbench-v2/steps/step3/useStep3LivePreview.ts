import * as React from 'react';
import { mappingAPI } from '../../../../../services/datalink';
import type { MappingPreviewResponse } from '../../../../../types/datalink';
import type { Mapping, Point } from '../../state/types';
import { buildPayload } from '../../state/transformPipeline';

export interface Step3LivePreviewState {
  status: 'idle' | 'loading' | 'success' | 'error';
  data: MappingPreviewResponse | null;
  error: string | null;
}

const PREVIEW_DEBOUNCE_MS = 250;

export function useStep3LivePreview(
  point: Point | null,
  mapping: Mapping | null,
  rawValue: unknown | null,
): Step3LivePreviewState {
  const [state, setState] = React.useState<Step3LivePreviewState>({
    status: 'idle',
    data: null,
    error: null,
  });
  const requestIdRef = React.useRef(0);

  React.useEffect(() => {
    if (!point || !mapping || rawValue === null || rawValue === undefined) {
      setState({
        status: 'idle',
        data: null,
        error: null,
      });
      return;
    }

    const currentRequestId = ++requestIdRef.current;
    let disposed = false;

    setState((prev) => ({
      status: 'loading',
      data: prev.data,
      error: null,
    }));

    const timer = window.setTimeout(() => {
      void mappingAPI
        .preview({
          raw_value: rawValue,
          transform_pipeline: (buildPayload(mapping, point) as {
            transform_pipeline: MappingPreviewResponse['pipeline'];
          }).transform_pipeline,
        })
        .then((data) => {
          if (disposed || currentRequestId !== requestIdRef.current) {
            return;
          }

          if (data.error) {
            setState({
              status: 'error',
              data,
              error: data.error,
            });
            return;
          }

          setState({
            status: 'success',
            data,
            error: null,
          });
        })
        .catch((error: unknown) => {
          if (disposed || currentRequestId !== requestIdRef.current) {
            return;
          }

          setState({
            status: 'error',
            data: null,
            error: error instanceof Error ? error.message : 'preview failed',
          });
        });
    }, PREVIEW_DEBOUNCE_MS);

    return () => {
      disposed = true;
      window.clearTimeout(timer);
    };
  }, [mapping, point, rawValue]);

  return state;
}
