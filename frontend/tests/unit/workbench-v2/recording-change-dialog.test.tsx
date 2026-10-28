import * as React from 'react';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { RecordingChangeDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog';
import { fixture } from '../../fixtures/recordingChange';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
function setup() {
  const f = fixture(); const controller = new RecordingChangeController(f.workspace, f.deps);
  function Host() {
    const state = React.useSyncExternalStore(controller.subscribe, controller.snapshot);
    return <><button onClick={() => controller.requestAction(f.intent)}>Open change</button><RecordingChangeDialog controller={{ ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: controller.unresolved }} destinations={[{ id: 'C', revision: 'C1', name: 'Safe destination' }]} /></>;
  }
  const host = render(<Host />); return { ...host, f, controller };
}
describe('Radix saved recording change dialog', () => {
  it('has a name, three steps, keyboard focus trap, Escape cancellation and focus restoration', async () => {
    const { f } = setup(); const trigger = screen.getByRole('button', { name: 'Open change' }); trigger.focus(); fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog', { name: 'dialog.title' });
    expect(within(dialog).getAllByRole('listitem')).toHaveLength(3);
    const first = within(dialog).getByRole('textbox'); first.focus(); fireEvent.keyDown(first, { key: 'Tab', shiftKey: true }); expect(within(dialog).getByTestId('recording-change-review')).toHaveFocus();
    fireEvent.keyDown(document.activeElement ?? document, { key: 'Escape' }); await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull()); expect(trigger).toHaveFocus();
    expect(f.api.review).not.toHaveBeenCalled(); expect(f.api.reserve).not.toHaveBeenCalled(); expect(f.api.schemaPreview).not.toHaveBeenCalled();
  });
  it('shows one primary operation per stage and never claims reserved settings are running', async () => {
    const { f } = setup(); fireEvent.click(screen.getByText('Open change')); fireEvent.click(await screen.findByTestId('recording-change-review'));
    const reserve = await screen.findByTestId('recording-change-reserve_preview'); fireEvent.click(reserve);
    expect(await screen.findByText('dialog.reserved')).toBeVisible(); expect(await screen.findByTestId('recording-change-confirm_schema')).toBeEnabled();
    expect(f.api.schemaConfirm).not.toHaveBeenCalled(); expect(screen.queryByText('PRIVATE SQL')).not.toBeInTheDocument(); expect(screen.queryByText('private.db')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('recording-change-confirm_schema')); expect(await screen.findByText('dialog.schema_confirmed')).toBeVisible(); expect(screen.getByTestId('recording-change-review_activation')).toBeDisabled();
  });
  it('after reservation Close only dismisses and readonly still permits original status recovery', async () => {
    const { f, controller } = setup(); fireEvent.click(screen.getByText('Open change')); fireEvent.click(await screen.findByTestId('recording-change-review')); fireEvent.click(await screen.findByTestId('recording-change-reserve_preview'));
    await screen.findByText('dialog.reserved'); fireEvent.click(screen.getByText('dialog.close')); expect(screen.getByTestId('recording-change-recovery')).toBeVisible();
    controller.readonly = true; fireEvent.click(screen.getByText('dialog.resume')); fireEvent.click(await screen.findByTestId('recording-change-check_status'));
    await waitFor(() => expect(f.api.operationStatus).toHaveBeenCalledTimes(1)); expect(f.api.reserve).toHaveBeenCalledTimes(1);
  });
  it('exposes separate pending units in review rather than one total', async () => {
    const { f, controller } = setup(); f.review.history.generations[0].revisions.push({ group_revision: 'A1', effective_at: '2026-01-01T00:00:00Z', sources: [], destination: { connector_id: 'C', connector_revision: 'C1', destination_scope: 'safe' }, pending: { samples: 3, outbox: { outstanding: 7, total: 7, blocked: 2, unknown: 1, queued: 4, retrying: 0, quarantined: 0, sql_committed: 0, skipped: 0 }, handoffs: [] } });
    await act(async () => { await controller.requestAction(f.intent); await controller.review(); });
    expect(screen.getByTestId('lifecycle-pending-samples')).toHaveTextContent('3'); expect(screen.getByTestId('lifecycle-pending-outbox')).toHaveTextContent('7'); expect(screen.getByTestId('lifecycle-pending-unknown')).toHaveTextContent('1');
  });
});
