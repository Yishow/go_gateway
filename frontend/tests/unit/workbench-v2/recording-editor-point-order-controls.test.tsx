import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { RecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointOrder';
import { useRecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
import { pointFormatAPI } from '../../../src/services/pointFormat';
import type { PointFormatSnapshot } from '../../../src/types/pointFormat';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/services/pointFormat', async original => ({ ...await original<typeof import('../../../src/services/pointFormat')>(), pointFormatAPI: { read: vi.fn(), save: vi.fn(), inspect: vi.fn() } }));
const snapshot = (): PointFormatSnapshot => ({ workspace_id: 'W', device_id: 'D', point_id: 'P', rule_id: 'R', workspace_revision: 'source-7', workspace_revision_origin: 'persisted', point_revision: '', rule_revision: 'rule-2', connection_revision: 'conn-3', editable: true,
  resolution: { resolved: true, format_mode: 'legacy', configured_format: '', effective_format: 'ABCD', effective_source: 'connection', effective_policy: 'legacy_mc3e', native_compatibility: false, applicable: true, format_version: 0, workspace_revision: 'source-7', point_revision: '', rule_revision: 'rule-2', connection_revision: 'conn-3' } });
function Page({ open = true, disabled = false, enabled = true }: { open?: boolean; disabled?: boolean; enabled?: boolean }) {
  const port = useRecordingPointOrder('W', enabled);
  return open ? <RecordingPointOrder port={port} workspaceId="W" deviceId="D" pointId="P" disabled={disabled} /> : null;
}
beforeEach(() => { vi.clearAllMocks(); vi.mocked(pointFormatAPI.read).mockImplementation(async () => snapshot()); });
async function ready() { await screen.findByLabelText('editor.order'); fireEvent.click(screen.getByText('editor.title', { exact: false })); }
it('shows legacy as preserved-only and saves explicit override equal to default without changing it to inherit', async () => {
  render(<Page />); await ready();
  expect(screen.getByRole('option', { name: 'mode.legacy' })).toBeDisabled();
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'ABCD' } });
  expect(pointFormatAPI.save).not.toHaveBeenCalled();
  fireEvent.click(screen.getByText('editor.save'));
  await waitFor(() => expect(pointFormatAPI.save).toHaveBeenCalledWith(expect.objectContaining({ format_mode: 'override', configured_format: 'ABCD', expected_workspace_revision: 'source-7', expected_point_revision: '' })));
});
it('reset sends explicit inheritance and native empty effective order is not presented as ABCD', async () => {
  const native = snapshot(); Object.assign(native.resolution, { applicable: false, effective_format: '', effective_source: 'native', effective_policy: 'legacy_fatek_components', native_compatibility: true });
  vi.mocked(pointFormatAPI.read).mockResolvedValue(native);
  render(<Page />); await ready();
  expect(screen.getByTestId('point-order-effective')).toHaveTextContent('editor.native');
  expect(screen.getByRole('option', { name: 'ABCD · mode.override' })).toBeDisabled();
  fireEvent.click(screen.getByText('editor.reset')); fireEvent.click(screen.getByText('editor.save'));
  await waitFor(() => expect(pointFormatAPI.save).toHaveBeenCalledWith(expect.objectContaining({ format_mode: 'inherit', configured_format: '' })));
});
it('noneditable and read-only source details never submit an override', async () => {
  const locked = snapshot(); locked.editable = false; vi.mocked(pointFormatAPI.read).mockResolvedValue(locked);
  render(<Page disabled />); await ready();
  expect(screen.getByLabelText('editor.order')).toBeDisabled(); expect(screen.getByText('editor.reset')).toBeDisabled(); expect(screen.getByText('editor.save')).toBeDisabled();
  fireEvent.click(screen.getByText('editor.save')); expect(pointFormatAPI.save).not.toHaveBeenCalled();
});
it('closing and reopening the real row preserves its original pending operation and offers inspect, never another Save', async () => {
  let reject!: (error: Error) => void;
  vi.mocked(pointFormatAPI.save).mockImplementation(() => new Promise((_resolve, fail) => { reject = fail; }));
  const view = render(<Page />); await ready();
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'CDAB' } }); fireEvent.click(screen.getByText('editor.save'));
  await waitFor(() => expect(pointFormatAPI.save).toHaveBeenCalledOnce());
  const original = vi.mocked(pointFormatAPI.save).mock.calls[0][0];
  view.rerender(<Page open={false} />); reject(new Error('timeout'));
  view.rerender(<Page />); await ready();
  fireEvent.click(await screen.findByText('editor.inspect'));
  await waitFor(() => expect(pointFormatAPI.inspect).toHaveBeenCalledWith(original));
  expect(pointFormatAPI.save).toHaveBeenCalledOnce(); expect(pointFormatAPI.read).toHaveBeenCalledOnce();
});

it('loads once when page hydration becomes available after the row already mounted', async () => {
  const view = render(<Page enabled={false} />);
  expect(pointFormatAPI.read).not.toHaveBeenCalled();
  view.rerender(<Page enabled />);
  await screen.findByLabelText('editor.order');
  expect(pointFormatAPI.read).toHaveBeenCalledOnce();
});
