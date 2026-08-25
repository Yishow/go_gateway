import { act, render, renderHook, screen, waitFor } from '@testing-library/react';
import { beforeEach, afterEach, describe, expect, it, vi } from 'vitest';
import { LivePreviewPanel } from '../../../src/components/datalink/wizard/LivePreviewPanel';
import { parsePreviewEvent, usePreviewStream } from '../../../src/hooks/usePreviewStream';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { defaultValue?: string }) => {
      if (key === 'errors.preview_invalid_request') return 'Invalid preview request.';
      if (key === 'errors.preview_unavailable') return 'Live preview is unavailable.';
      if (key === 'errors.request_id') return 'Request ID';
      return options?.defaultValue ?? key;
    },
  }),
}));

class MockEventSource {
  static instances: MockEventSource[] = [];
  static CLOSED = 2;
  static OPEN = 1;
  readonly url: string;
  onopen: (() => void) | null = null;
  onmessage: ((event: MessageEvent<string>) => void) | null = null;
  onerror: (() => void) | null = null;
  closed = false;

  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }

  close() {
    this.closed = true;
  }

  emitOpen() {
    this.onopen?.();
  }

  emitMessage(payload: unknown) {
    this.onmessage?.({ data: JSON.stringify(payload) } as MessageEvent<string>);
  }

  emitError() {
    this.onerror?.();
  }
}

describe('usePreviewStream', () => {
  beforeEach(() => {
    MockEventSource.instances = [];
    vi.stubGlobal('EventSource', MockEventSource);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('parses top-level and transitional nested typed errors without retaining raw details', () => {
    expect(parsePreviewEvent({
      type: 'error',
      code: 'preview_stream_closed',
      action: 'retry the preview stream',
      request_id: 'req-top-level-1',
      retryable: true,
      error: 'dial tcp 10.0.0.9:502 with password secret',
      timestamp: '2026-08-25T00:00:00Z',
    })).toMatchObject({
      type: 'error',
      code: 'preview_stream_closed',
      action: 'retry the preview stream',
      request_id: 'req-top-level-1',
      retryable: true,
    });
    expect(parsePreviewEvent({
      type: 'error',
      error: {
        code: 'preview_unavailable',
        action: 'retry preview',
        request_id: 'req-nested-1',
        retryable: true,
      },
      timestamp: '2026-08-25T00:00:00Z',
    })).toMatchObject({
      code: 'preview_unavailable',
      action: 'retry preview',
      request_id: 'req-nested-1',
      retryable: true,
      error: { code: 'preview_unavailable', request_id: 'req-nested-1' },
    });

    const parsed = parsePreviewEvent({
      type: 'error',
      error: 'raw backend detail must not become UI state',
      timestamp: '2026-08-25T00:00:00Z',
    });
    expect(parsed).not.toHaveProperty('error');
    expect(parsed).not.toHaveProperty('code');
    expect(parsePreviewEvent({
      type: 'preview',
      timestamp: '2026-08-25T00:00:00Z',
      steps: [{ step_index: 0, step_type: 'scale', input: 1, output: 2, error: 'secret transform details' }],
    })).toMatchObject({ steps: [{ step_index: 0, step_type: 'scale', input: 1, output: 2 }] });
    expect(parsePreviewEvent({
      type: 'preview',
      timestamp: '2026-08-25T00:00:00Z',
      steps: [{ step_index: 0, step_type: 'scale', input: 1, output: 2, error: 'secret transform details' }],
    })?.steps?.[0]).not.toHaveProperty('error');
  });

  it('parses typed preview close and error metadata while rejecting untrusted shapes', () => {
    expect(parsePreviewEvent({
      type: 'close',
      code: 'preview_unavailable',
      request_id: 'req-close-1',
      timestamp: '2026-08-25T00:00:00Z',
    })).toMatchObject({
      type: 'close',
      code: 'preview_unavailable',
      request_id: 'req-close-1',
    });

    const malformedError = parsePreviewEvent({
      type: 'error',
      code: { secret: 'must not become UI text' },
      request_id: 42,
      timestamp: '2026-08-25T00:00:00Z',
    });
    expect(malformedError?.type).toBe('error');
    expect(malformedError).not.toHaveProperty('code');
    expect(malformedError).not.toHaveProperty('request_id');

    expect(parsePreviewEvent({ type: 'unknown', timestamp: '2026-08-25T00:00:00Z' })).toBeNull();
    expect(parsePreviewEvent('raw backend detail')).toBeNull();
  });

  it('stores typed SSE error metadata and closes on a typed close event', async () => {
    const { result } = renderHook(() => usePreviewStream({
      mappingId: 'mapping-1',
      workspaceId: 'workspace/1',
      autoConnect: true,
      baseUrl: '/api/v1',
      maxReconnectAttempts: 0,
    }));
    const source = MockEventSource.instances[0];
    act(() => {
      source.emitOpen();
      source.emitMessage({
        type: 'error',
        code: 'preview_invalid_request',
        request_id: 'req-error-1',
        error: 'raw backend detail must not be rendered',
        timestamp: '2026-08-25T00:00:00Z',
      });
    });

    await waitFor(() => expect(result.current.error).toBe('preview_invalid_request'));
    expect(result.current.requestId).toBe('req-error-1');

    act(() => {
      source.emitMessage({
        type: 'close',
        code: 'preview_unavailable',
        request_id: 'req-close-2',
        timestamp: '2026-08-25T00:00:01Z',
      });
    });

    await waitFor(() => expect(result.current.connectionState).toBe('disconnected'));
    expect(result.current.error).toBe('preview_unavailable');
    expect(result.current.requestId).toBe('req-close-2');
    expect(source.closed).toBe(true);
  });

  it('does not reconnect or mutate state after an explicit disconnect', async () => {
    const { result } = renderHook(() => usePreviewStream({
      mappingId: 'mapping-1',
      workspaceId: 'workspace-1',
      autoConnect: true,
      reconnectDelay: 1,
      maxReconnectAttempts: 1,
    }));
    const source = MockEventSource.instances[0];
    act(() => {
      result.current.disconnect();
      source.emitError();
      source.emitMessage({
        type: 'error',
        code: 'preview_unavailable',
        timestamp: '2026-08-25T00:00:00Z',
      });
    });

    await new Promise((resolve) => setTimeout(resolve, 5));
    expect(MockEventSource.instances).toHaveLength(1);
    expect(result.current.connectionState).toBe('disconnected');
    expect(result.current.error).toBeNull();
  });

  it('renders localized safe preview errors without raw SSE details', async () => {
    render(<LivePreviewPanel mappingId="mapping-1" workspaceId="workspace-1" autoStart />);
    const source = MockEventSource.instances[0];
    act(() => {
      source.emitMessage({
        type: 'error',
        code: 'preview_invalid_request',
        request_id: 'req-ui-1',
        error: 'dial tcp 10.0.0.9:502 with password secret',
        timestamp: '2026-08-25T00:00:00Z',
      });
    });

    await waitFor(() => expect(screen.getByText(/Invalid preview request\./)).toBeInTheDocument());
    expect(screen.getByText(/請求識別碼|Request ID/)).toHaveTextContent('req-ui-1');
    expect(screen.queryByText(/dial tcp|password secret/)).not.toBeInTheDocument();
  });

  it('does not open a preview stream without a workspace scope and encodes both scopes', () => {
    renderHook(() => usePreviewStream({
      mappingId: 'mapping/1',
      workspaceId: '',
      autoConnect: true,
    }));
    expect(MockEventSource.instances).toHaveLength(0);

    renderHook(() => usePreviewStream({
      mappingId: 'mapping/1',
      workspaceId: 'workspace/1',
      autoConnect: true,
    }));
    expect(MockEventSource.instances[0].url).toContain('workspace_id=workspace%2F1');
    expect(MockEventSource.instances[0].url).toContain('mapping_id=mapping%2F1');
  });
});
