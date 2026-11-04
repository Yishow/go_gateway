import { useState } from 'react';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingDraftNavigation, useRecordingDraftGuard } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

it('discards only the picker name and selection, keeping a separately guarded recording draft', async () => {
  const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'p', label: 'Point', device_name: 'PLC', address: 'D0', target_type: 'int16' };
  const session = { id: 'S', identity: candidate, displayName: 'Point' };
  const name = { prepare: vi.fn(async () => ({ session, status: 'ready' as const })), submit: vi.fn(), inspect: vi.fn(), release: vi.fn() };
  function Editor() {
    const [draft, setDraft] = useState('Keep the recording');
    useRecordingDraftGuard({ version: draft, dirty: true, discard: () => setDraft('Discarded recording') });
    return <><input aria-label="Recording draft" value={draft} onChange={event => setDraft(event.target.value)} />
      <RecordingPointSelection scope="W/G" workspaceId="W" candidates={[candidate]} members={[]} onChange={vi.fn()} disabled={false}
        devices={[{ id: 'D', name: 'PLC' }]} deviceSource={{ name, submit: vi.fn(), inspect: vi.fn(), observations: () => [] }} /></>;
  }
  render(<RecordingDraftNavigation><Editor /></RecordingDraftNavigation>);
  fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point' }));
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus());
  await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Unsent point name' } });
  fireEvent.click(within(screen.getByRole('contentinfo')).getByRole('button', { name: 'picker.cancel' }));
  fireEvent.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: 'name.cancelDiscard' }));
  expect(screen.getByLabelText('Recording draft')).toHaveValue('Keep the recording');
  expect(name.submit).not.toHaveBeenCalled();
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
});
