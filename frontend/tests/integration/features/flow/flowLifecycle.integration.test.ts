import { act, renderHook } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { useFlowLifecycle } from '@/features/flow/stateMachine';

describe('useFlowLifecycle integration', () => {
  it('supports draft -> validated -> active flow with diagnostics updates', () => {
    const { result } = renderHook(() => useFlowLifecycle());

    act(() => {
      result.current.setSource('device-1', 'D100', 'point-1');
      result.current.setTag('tag-1');
      result.current.setDiagnostics({
        source: {
          latestValue: '23.1',
          quality: 'good',
          timestamp: '2026-02-13T00:00:00Z',
          error: '',
        },
        grid: {
          latestValue: 'D100',
          quality: 'good',
          timestamp: '2026-02-13T00:00:01Z',
          error: '',
        },
        tag: {
          latestValue: 'line.temp',
          quality: 'good',
          timestamp: '2026-02-13T00:00:02Z',
          error: '',
        },
      });
    });

    expect(result.current.canValidate).toBe(true);
    expect(result.current.state.status).toBe('draft');

    act(() => {
      result.current.markValidated();
    });

    expect(result.current.state.status).toBe('validated');
    expect(result.current.canActivate).toBe(true);

    act(() => {
      result.current.markActive();
    });

    expect(result.current.state.status).toBe('active');
    expect(result.current.state.diagnostics.tag.latestValue).toBe('line.temp');
  });

  it('captures sink failure and supports recovery', () => {
    const { result } = renderHook(() => useFlowLifecycle());

    act(() => {
      result.current.markError('sink', 'write timeout');
    });

    expect(result.current.state.status).toBe('error');
    expect(result.current.state.diagnostics.sink.error).toBe('write timeout');
    expect(result.current.hasError).toBe(true);

    act(() => {
      result.current.resetDraft();
      result.current.setDiagnostics({
        source: {
          latestValue: '-',
          quality: 'unknown',
          timestamp: '-',
          error: '',
        },
        grid: {
          latestValue: '-',
          quality: 'unknown',
          timestamp: '-',
          error: '',
        },
        tag: {
          latestValue: '-',
          quality: 'unknown',
          timestamp: '-',
          error: '',
        },
        sink: {
          latestValue: '-',
          quality: 'unknown',
          timestamp: '-',
          error: '',
        },
      });
    });

    expect(result.current.state.status).toBe('draft');
    expect(result.current.state.diagnostics.sink.error).toBe('');
    expect(result.current.state.diagnostics.sink.quality).toBe('unknown');
  });
});
