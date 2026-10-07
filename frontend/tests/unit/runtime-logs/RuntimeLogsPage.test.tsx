import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogsPage from '../../../src/features/runtime-logs/RuntimeLogsPage';
import { runtimeLogTransport, RuntimeLogRequestError } from '../../../src/services/runtimeLogs';
import type { RuntimeLogEvent, RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';
import i18n from '../../../src/i18n/config';
import { snapshot as fixture, record } from './fixtures';
vi.mock('../../../src/services/runtimeLogs', async (original) => ({
  ...await original<typeof import('../../../src/services/runtimeLogs')>(),
  runtimeLogTransport: { snapshot: vi.fn(), stream: vi.fn() },
}));
const text = '<img src=x onerror="alert(1)"><script>alert(2)</script>';
const snapshot: RuntimeLogSnapshot = { ...fixture, records: [{ ...record('1', text), truncated: true }], sink_health: 'degraded' };
let emit: (event: RuntimeLogEvent) => void;
beforeEach(async () => {
  await i18n.changeLanguage('en');
  vi.mocked(runtimeLogTransport.snapshot).mockResolvedValue(snapshot);
  vi.mocked(runtimeLogTransport.stream).mockImplementation((_filters, _cursor, _signal, callback) => {
    emit = callback; return new Promise(() => undefined);
  });
});
afterEach(() => vi.clearAllMocks());
describe('runtime log operator view', () => {
  it('renders literal hostile text and accessible controls with independent storage/truncation states', async () => {
    const view = render(<RuntimeLogsPage />); await screen.findByText(text);
    expect(view.container.querySelector('img,script')).toBeNull();
    expect(screen.getByRole('combobox', { name: 'Minimum severity' })).toBeInTheDocument();
    expect(screen.getByLabelText('Search safe text')).toHaveAttribute('maxlength', '256');
    expect(screen.getByText(/Diagnostic storage is degraded/)).toBeInTheDocument();
    expect(screen.getByText(/Some displayed content was truncated/)).toBeInTheDocument();
    act(() => emit({ type: 'handshake', metadata: snapshot, cursor: 'instance-a:1' }));
    expect(screen.getByRole('status')).toHaveTextContent('Connected');
    fireEvent.click(screen.getByRole('button', { name: 'Pause view' }));
    fireEvent.click(screen.getByRole('button', { name: 'Clear view' }));
    expect(screen.getByRole('status')).toHaveTextContent('Paused'); expect(screen.queryByText(text)).toBeNull();
    fireEvent.click(screen.getByRole('button', { name: 'Resume view' }));
    await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalledTimes(2));
    expect(runtimeLogTransport.snapshot).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByLabelText('Follow tail')); expect(screen.getByLabelText('Follow tail')).not.toBeChecked();
  });
  it('distinguishes empty, capture gaps, file-only gaps and reset', async () => {
    vi.mocked(runtimeLogTransport.snapshot).mockResolvedValue({ ...snapshot, records: [] });
    render(<RuntimeLogsPage />); await waitFor(() => expect(runtimeLogTransport.stream).toHaveBeenCalled());
    act(() => emit({ type: 'handshake', metadata: snapshot, cursor: 'instance-a:1' }));
    expect(screen.getByText('No matching retained records')).toBeInTheDocument();
    act(() => emit({ type: 'gap', reason: 'file_loss', cursor: '' }));
    expect(screen.getByText(/Diagnostic files have missing history/)).toBeInTheDocument();
    expect(screen.queryByText(/Some capture history is missing/)).toBeNull();
    act(() => emit({ type: 'gap', reason: 'retention', cursor: '' }));
    expect(screen.getByText(/Some capture history is missing/)).toBeInTheDocument();
    act(() => emit({ type: 'reset', reason: 'reset', cursor: '' }));
    expect(screen.getByText(/The gateway process restarted/)).toBeInTheDocument();
  });
  it('ignores late filter A and cancels filter B when navigating away', async () => {
    let finish!: (value: RuntimeLogSnapshot) => void;
    vi.mocked(runtimeLogTransport.snapshot).mockImplementationOnce(() => new Promise((resolve) => { finish = resolve; }));
    const view = render(<RuntimeLogsPage />);
    fireEvent.change(screen.getByLabelText('Source (exact ID)'), { target: { value: 'http' } });
    fireEvent.change(screen.getByLabelText('Search safe text'), { target: { value: 'literal.*' } });
    fireEvent.click(screen.getByRole('button', { name: 'Apply filters' })); await screen.findByText(text);
    expect(vi.mocked(runtimeLogTransport.snapshot).mock.calls[0][1].aborted).toBe(true);
    expect(runtimeLogTransport.stream).toHaveBeenCalledWith({ level: 'debug', source: 'http', q: 'literal.*' },
      'instance-a:1', expect.any(AbortSignal), expect.any(Function));
    await act(async () => finish({ ...snapshot, records: [record('99', 'late A')] }));
    expect(screen.queryByText('late A')).toBeNull(); expect(runtimeLogTransport.stream).toHaveBeenCalledTimes(1);
    view.unmount(); expect(vi.mocked(runtimeLogTransport.stream).mock.calls[0][2].aborted).toBe(true);
  });
  it('shows safe failed state rather than empty and supports explicit retry', async () => {
    vi.mocked(runtimeLogTransport.snapshot).mockRejectedValueOnce(new RuntimeLogRequestError(400));
    render(<RuntimeLogsPage />); await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('Failed'));
    expect(screen.queryByText('No matching retained records')).toBeNull();
    fireEvent.click(screen.getByRole('button', { name: 'Retry connection' })); await screen.findByText(text);
  });
  it('translates denial and scope controls into Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW');
    vi.mocked(runtimeLogTransport.snapshot).mockRejectedValue(new RuntimeLogRequestError(403));
    render(<RuntimeLogsPage />); expect(await screen.findByText(/僅限本機存取/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '套用篩選' })).toBeInTheDocument();
    expect(screen.getByText(/僅搜尋本次程序保留的有限日誌/)).toBeInTheDocument();
  });
});
