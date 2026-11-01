import { RecordingSourceNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSourceNavigation';
import { useState } from 'react';
import { RecordingDraftNavigation, useRecordingNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { RecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointOrder';
import { useRecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
import { pointFormatAPI } from '../../../src/services/pointFormat';
import type { PointFormatSnapshot } from '../../../src/types/pointFormat';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/services/pointFormat', async original => ({ ...await original<typeof import('../../../src/services/pointFormat')>(), pointFormatAPI: { read: vi.fn(), save: vi.fn(), inspect: vi.fn() } }));
const snapshot = (): PointFormatSnapshot => ({ workspace_id: 'W', device_id: 'D', point_id: 'P', rule_id: 'R', workspace_revision: 'source-7', workspace_revision_origin: 'persisted', point_revision: '', rule_revision: 'rule-2', connection_revision: 'conn-3', editable: true,
  resolution: { resolved: true, format_mode: 'legacy', configured_format: '', effective_format: 'ABCD', effective_source: 'connection', effective_policy: 'legacy_mc3e', native_compatibility: false, applicable: true, format_version: 0, workspace_revision: 'source-7', point_revision: '', rule_revision: 'rule-2', connection_revision: 'conn-3' } });
function Page() {
  const port = useRecordingPointOrder('W', true);
  return <RecordingDraftNavigation><RecordingSourceNavigation workspaceId="W" pointOrder={port}><View port={port} /></RecordingSourceNavigation></RecordingDraftNavigation>;
}
function View({ port }: { port: ReturnType<typeof useRecordingPointOrder> }) {
  const [open, setOpen] = useState(true); const navigation = useRecordingNavigation();
  return <><button onClick={() => navigation.request(() => setOpen(false))}>Leave record</button><button onClick={() => setOpen(true)}>Reopen</button><button onClick={() => setOpen(false)}>Hide row</button>
    {open && <RecordingPointOrder port={port} workspaceId="W" deviceId="D" pointId="P" disabled={false} />}</>;
}
beforeEach(() => { vi.clearAllMocks(); vi.mocked(pointFormatAPI.read).mockImplementation(async () => snapshot()); });
it('unsent point-format intent uses the leave guard and discard only clears that unsent intent', async () => {
  render(<Page />); await screen.findByLabelText('editor.order'); fireEvent.click(screen.getByText('editor.title', { exact: false }));
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'CDAB' } });
  fireEvent.click(screen.getByText('Leave record')); expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.savedSources');
  fireEvent.click(screen.getByText('navigation.keep')); expect(screen.getByLabelText('editor.order')).toHaveValue('CDAB');
  fireEvent.click(screen.getByText('Leave record')); fireEvent.click(screen.getByText('navigation.discard')); fireEvent.click(screen.getByText('Reopen'));
  expect(screen.getByLabelText('editor.order')).toHaveValue('legacy'); expect(pointFormatAPI.save).not.toHaveBeenCalled();
});
it('a submitted point format remains recoverable after leaving rather than being discarded or resubmitted', async () => {
  vi.mocked(pointFormatAPI.save).mockRejectedValue(new Error('unknown'));
  render(<Page />); await screen.findByLabelText('editor.order'); fireEvent.click(screen.getByText('editor.title', { exact: false }));
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'CDAB' } }); fireEvent.click(screen.getByText('editor.save'));
  await screen.findByText('editor.inspect'); fireEvent.click(screen.getByText('Leave record'));
  expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
  fireEvent.click(screen.getByText('navigation.discard')); fireEvent.click(screen.getByText('Reopen'));
  expect(await screen.findByText('editor.inspect')).toBeInTheDocument(); expect(pointFormatAPI.save).toHaveBeenCalledOnce();
});

it.each(['unsent', 'unknown'])('keeps %s source intent guarded after its member row is removed from view', async status => {
  vi.mocked(pointFormatAPI.save).mockRejectedValue(new Error('unknown'));
  render(<Page />); await screen.findByLabelText('editor.order'); fireEvent.click(screen.getByText('editor.title', { exact: false }));
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'CDAB' } });
  if (status === 'unknown') { fireEvent.click(screen.getByText('editor.save')); await screen.findByText('editor.inspect'); }
  fireEvent.click(screen.getByText('Hide row')); expect(screen.queryByLabelText('editor.order')).not.toBeInTheDocument();
  fireEvent.click(screen.getByText('Leave record')); expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
  fireEvent.click(screen.getByText('navigation.keep')); fireEvent.click(screen.getByText('Reopen'));
  if (status === 'unknown') expect(await screen.findByText('editor.inspect')).toBeInTheDocument();
  else expect(screen.getByLabelText('editor.order')).toHaveValue('CDAB');
});
it('a confirmed format save with runtime apply pending is saved source state, not an endless in-flight draft', async () => {
  vi.mocked(pointFormatAPI.save).mockResolvedValue({ point: { id: 'P', device_id: 'D', revision_id: 'point-2', format_mode: 'override', format_override: 'CDAB', format_version: 1, data_format: 'CDAB' },
    resolution: { ...snapshot().resolution, format_mode: 'override', configured_format: 'CDAB', effective_format: 'CDAB', point_revision: 'point-2', workspace_revision: 'source-8' },
    workspace_revision: 'source-8', runtime_apply_status: 'pending', invalidation_status: 'invalidated' });
  render(<Page />); await screen.findByLabelText('editor.order'); fireEvent.click(screen.getByText('editor.title', { exact: false }));
  fireEvent.change(screen.getByLabelText('editor.order'), { target: { value: 'CDAB' } }); fireEvent.click(screen.getByText('editor.save'));
  await screen.findByText('editor.recover'); fireEvent.click(screen.getByText('Leave record'));
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  fireEvent.click(screen.getByText('Reopen')); expect(await screen.findByText('editor.recover')).toBeInTheDocument();
  expect(pointFormatAPI.save).toHaveBeenCalledOnce();
});
