import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { readRecordingRangeSchedule } from '@/features/datalink/workbench-v2/recording-editor/recordingRangeSchedule';
import { RecordingRangeSchedule } from '@/features/datalink/workbench-v2/recording-editor/RecordingRangeSchedule';
import type { RangeSourceOutcome } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
import { pointAPI, pollingGroupAPI } from '@/services/datalink';
vi.mock('@/services/datalink', () => ({ pointAPI: { get: vi.fn() }, pollingGroupAPI: { list: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, params?: Record<string, unknown>) => params ? `${key} ${JSON.stringify(params)}` : key }) }));
const sources = [{ device_id: 'D', point_id: 'P', tag_id: 'T' }];
const point = { id: 'P', device_id: 'D', polling_group_id: 'G' } as Awaited<ReturnType<typeof pointAPI.get>>;
const group = { id: 'G', name: 'Configured group', interval_ms: 2500, enabled: true, priority: 50, created_at: '', updated_at: '' };
const outcome = { request: { workspaceId: 'W', deviceId: 'D', requestId: 'R', sessionId: 'S', generation: 1 }, status: 'saved-inactive', savedPoints: sources } as RangeSourceOutcome;
beforeEach(() => { vi.clearAllMocks(); vi.mocked(pointAPI.get).mockResolvedValue(point); vi.mocked(pollingGroupAPI.list).mockResolvedValue([group]); });
it('uses exact persisted point and polling-group IDs, never a guessed default interval', async () => {
  expect(await readRecordingRangeSchedule('D', sources)).toEqual([{ id: 'G', name: 'Configured group', intervalMs: 2500, enabled: true, points: 1 }]);
  expect(pointAPI.get).toHaveBeenCalledWith('P');
  for (const changed of [{ ...point, device_id: 'foreign' }, { ...point, id: 'foreign' }, { ...point, polling_group_id: '' }]) {
    vi.mocked(pointAPI.get).mockResolvedValue(changed);
    expect(await readRecordingRangeSchedule('D', sources)).toBeUndefined();
  }
  vi.mocked(pointAPI.get).mockResolvedValue(point);
  vi.mocked(pollingGroupAPI.list).mockResolvedValue([{ ...group, id: 'other' }]);
  expect(await readRecordingRangeSchedule('D', sources)).toBeUndefined();
});
it('does not fetch before persisted identities or for a different workspace, and never polls', async () => {
  const view = render(<RecordingRangeSchedule workspaceId="W" deviceId="D" />);
  expect(screen.getByText('range.scheduleUnassigned')).toBeInTheDocument();
  expect(pointAPI.get).not.toHaveBeenCalled();
  view.rerender(<RecordingRangeSchedule workspaceId="OTHER" deviceId="D" outcome={outcome} />);
  expect(pointAPI.get).not.toHaveBeenCalled();
  view.rerender(<RecordingRangeSchedule workspaceId="W" deviceId="D" outcome={outcome} />);
  await screen.findByText(/Configured group/);
  expect(pointAPI.get).toHaveBeenCalledTimes(1);
  view.rerender(<RecordingRangeSchedule workspaceId="W" deviceId="D" outcome={outcome} />);
  expect(pointAPI.get).toHaveBeenCalledTimes(1);
});
it('keeps failed reads honest and allows a display-only retry', async () => {
  vi.mocked(pointAPI.get).mockRejectedValueOnce(new Error('private transport detail'));
  render(<RecordingRangeSchedule workspaceId="W" deviceId="D" outcome={outcome} />);
  await screen.findByText('range.scheduleUnavailable');
  expect(screen.queryByText(/private transport/)).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'range.scheduleRefresh' }));
  await screen.findByText(/Configured group/);
  expect(pointAPI.get).toHaveBeenCalledTimes(2);
});
it('ignores a late read when the displayed source scope changes', async () => {
  let resolve!: (value: typeof point) => void;
  vi.mocked(pointAPI.get).mockImplementationOnce(() => new Promise(yes => { resolve = yes; }));
  const view = render(<RecordingRangeSchedule workspaceId="W" deviceId="D" outcome={outcome} />);
  await waitFor(() => expect(pointAPI.get).toHaveBeenCalledTimes(1));
  view.rerender(<RecordingRangeSchedule workspaceId="OTHER" deviceId="D" />);
  await act(async () => { resolve(point); });
  expect(screen.getByText('range.scheduleUnassigned')).toBeInTheDocument();
  expect(screen.queryByText(/Configured group/)).not.toBeInTheDocument();
});
