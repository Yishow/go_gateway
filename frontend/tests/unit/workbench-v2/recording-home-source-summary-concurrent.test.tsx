import { Suspense, startTransition, useState, type ReactNode } from 'react';
import { act, cleanup, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import axios from 'axios';
import { afterEach, expect, it, vi } from 'vitest';
import { useRecordingHomeSources } from '@/hooks/datalink/useRecordingHomeSources';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
const value = (workspace_id = 'W', name = 'observed') => ({ workspace_id, state: 'available' as const, source_digest: 'a'.repeat(64),
  devices: [{ id: 'D', name, protocol: 'modbus-tcp', availability: 'available' }], candidates: [], excluded: [] });
afterEach(() => { cleanup(); vi.restoreAllMocks(); });
it('an abandoned newer visit cannot poison the still-committed visit refresh', async () => {
  const client = new QueryClient(), backend = axios.create({ baseURL: '/api' });
  const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(value());
  const never = new Promise<void>(() => {}), attempted = vi.fn();
  function Visit({ epoch }: { epoch: number }) {
    const sources = useRecordingHomeSources({ workspaceId: 'W', workspaceRevision: '', epoch }, { backend });
    if (epoch === 2) { attempted(); throw never; }
    return <><p>{sources.status}</p><p>{sources.observed?.devices[0].name}</p><button onClick={() => { void sources.refetch(); }}>Retry committed visit</button></>;
  }
  function App() {
    const [epoch, setEpoch] = useState(1);
    return <><button onClick={() => startTransition(() => setEpoch(2))}>Speculate</button><Suspense fallback={<p>Suspended</p>}><Visit epoch={epoch} /></Suspense></>;
  }
  const view = render(<QueryClientProvider client={client}><App /></QueryClientProvider>);
  await screen.findByText('fresh'); fireEvent.click(screen.getByText('Speculate'));
  await waitFor(() => expect(attempted).toHaveBeenCalled());
  expect(read).toHaveBeenCalledTimes(1); expect(screen.queryByText('Suspended')).not.toBeInTheDocument();
  read.mockResolvedValue(value('W', 'updated'));
  fireEvent.click(screen.getByText('Retry committed visit'));
  await screen.findByText('updated'); await screen.findByText('fresh');
  view.unmount(); client.clear();
});
it('an abandoned ninth workspace render cannot evict a previously observed scope', async () => {
  const client = new QueryClient(), backend = axios.create({ baseURL: '/api' });
  const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockImplementation(async id => value(id));
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  for (let i = 0; i < 8; i++) {
    const view = renderHook(() => useRecordingHomeSources({ workspaceId: `W${i}`, workspaceRevision: '', epoch: i + 1 }, { backend }), { wrapper });
    await waitFor(() => expect(view.result.current.fresh).toBe(true)); view.unmount();
  }
  const never = new Promise<void>(() => {}), attempted = vi.fn();
  function Abandoned(): ReactNode {
    useRecordingHomeSources({ workspaceId: 'W8', workspaceRevision: '', epoch: 9 }, { backend });
    attempted(); throw never;
  }
  const suspended = render(<Suspense fallback={<p>Waiting</p>}><Abandoned /></Suspense>, { wrapper });
  await waitFor(() => expect(attempted).toHaveBeenCalled()); suspended.unmount();
  expect(read).toHaveBeenCalledTimes(8);
  read.mockReturnValue(new Promise(() => {}));
  const first = renderHook(() => useRecordingHomeSources({ workspaceId: 'W0', workspaceRevision: '', epoch: 10 }, { backend }), { wrapper });
  expect(first.result.current.observed?.workspace_id).toBe('W0');
  await act(async () => { first.unmount(); }); client.clear();
});
