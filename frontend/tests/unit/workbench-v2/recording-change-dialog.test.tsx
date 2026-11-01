import * as React from 'react';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { RecordingChangeDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog';
import type { RecordingDestinationChoice } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeTypes';
import { fixture } from '../../fixtures/recordingChange';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
function setup(destinations: RecordingDestinationChoice[] = [{ id: 'C', revision: 'C1', name: 'Safe destination' }]) {
  const f = fixture(); const controller = new RecordingChangeController(f.workspace, f.deps);
  function Host() {
    const state = React.useSyncExternalStore(controller.subscribe, controller.snapshot);
    return <><button onClick={() => controller.requestAction(f.intent)}>Open change</button><RecordingChangeDialog controller={{ ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: controller.unresolved }} destinations={destinations} /></>;
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

it('clears only the derived long table when selecting a new connector before original review', async () => {
  const { f } = setup([{ id: 'C2', revision: 'C2R', name: 'Fresh destination' }]);
  Object.assign(f.seed, { layout_version: 'managed_long_v1', storage_mode: 'plain', record_id: '11111111-1111-4111-8111-111111111111',
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
    row_policy: { allowed_lateness_seconds: 2, incomplete_policy: 'partial', record_key_column: 'record_id', provenance_column: 'provenance', bucket_start_column: 'bucket_start' },
  });
  f.seed.destination.table_name = 'gw_long_original_scope';
  Object.assign(f.seed.members[0], { calculation: 'delta', cadence: { interval: 10, unit: 'second', alignment: 'utc_epoch' },
    source_value_type: 'uint64', result_value_type: 'decimal', value_role: 'counter_raw', measurement_id: 'MEASUREMENT',
    exact_recording_source: { measurement_id: 'MEASUREMENT', definition_revision: 'DEF', source_binding_revision: 'BINDING', series_epoch: 'EPOCH', original_data_type: 'uint64', scale: '0.1', unit: 'kWh', semantic_kind: 'counter' },
  });
  const before = structuredClone(f.seed);
  fireEvent.click(screen.getByText('Open change'));
  fireEvent.change(await screen.findByLabelText('dialog.destination'), { target: { value: 'C2' } });
  fireEvent.click(screen.getByTestId('recording-change-review'));await waitFor(() => expect(f.api.review).toHaveBeenCalledOnce());
  const sent = f.api.review.mock.calls[0][1].group;
  expect(sent.destination).toEqual({ ...before.destination, connector_id: 'C2', connector_revision: 'C2R', database: '', table_name: '' });
  expect(sent.record_id).toBe(before.record_id);expect(sent.layout_version).toBe(before.layout_version);expect(sent.storage_mode).toBe('plain');
  expect(sent.members).toEqual(before.members);expect(sent.recording_profile).toEqual(before.recording_profile);expect(sent.row_policy).toEqual(before.row_policy);
  expect(f.seed).toEqual(before);expect(f.api.reserve).not.toHaveBeenCalled();expect(f.api.schemaConfirm).not.toHaveBeenCalled();
});
it('does not clear an explicit wide table when selecting a new connector', async () => {
  const { f } = setup([{ id: 'C2', revision: 'C2R', name: 'Fresh destination' }]);
  f.seed.layout_version = 'managed_wide_profile_v1';f.seed.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second' };f.seed.destination.table_name = 'explicit_wide_table';
  fireEvent.click(screen.getByText('Open change'));fireEvent.change(await screen.findByLabelText('dialog.destination'), { target: { value: 'C2' } });
  fireEvent.click(screen.getByTestId('recording-change-review'));await waitFor(() => expect(f.api.review).toHaveBeenCalledOnce());
  expect(f.api.review.mock.calls[0][1].group.destination).toMatchObject({ connector_id: 'C2', connector_revision: 'C2R', database: '', table_name: 'explicit_wide_table' });
  expect(f.api.reserve).not.toHaveBeenCalled();expect(f.api.schemaConfirm).not.toHaveBeenCalled();
});

it('collapsed optional prior-preview help never hides required risk, schema impact or the current unknown result', async () => {
  const { controller, unmount } = setup(); fireEvent.click(screen.getByText('Open change'));
  fireEvent.click(await screen.findByTestId('recording-change-review')); fireEvent.click(await screen.findByTestId('recording-change-reserve_preview'));
  await screen.findByText('dialog.reserved');
  const state = { ...controller.snapshot(), recovery: { ...controller.snapshot().recovery!, priorPreviewIds: ['older-preview'] }, unknown: true, error: 'unknown' as const };
  unmount(); render(<RecordingChangeDialog controller={{ ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: true }} />);
  const dialog = screen.getByRole('dialog'); const help = within(dialog).getByText('dialog.previous_previews').closest('details')!;
  expect(help).not.toHaveAttribute('open');
  for (const key of ['dialog.history_kept', 'dialog.schema_confirmation_required', 'dialog.unknown_no_replacement', 'dialog.errors.unknown']) {
    const message = within(dialog).getByText(key); expect(message).toBeVisible(); expect(help).not.toContainElement(message);
  }
});

describe('recording change impact summary', () => {
  it('distinguishes the same table name in different schemas without displaying a database path', async () => {
    const { f, controller } = setup([{ id: 'C', revision: 'C1', name: 'Original destination' }, { id: 'C2', revision: 'C2R', name: 'New destination' }]);
    Object.assign(f.seed, { layout_version: 'managed_long_v1', storage_mode: 'plain', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } });
    Object.assign(f.seed.destination, { table_schema: 'old_schema', table_name: 'same_table' });
    f.api.review.mockImplementation(async (_id, request) => ({ ...f.review, envelope: { ...f.envelope, request } }));
    await act(async () => {
      await controller.requestAction(f.intent);
      controller.setTarget({ ...f.seed, destination: { ...f.seed.destination, connector_id: 'C2', connector_revision: 'C2R', table_schema: 'new_schema' } });
      await controller.review();
    });
    expect(screen.getByTestId('recording-change-original-table')).toHaveTextContent('old_schema.same_table');
    expect(screen.getByTestId('recording-change-new-table')).toHaveTextContent('new_schema.same_table');
    expect(screen.queryByText('private.db')).toBeNull(); expect(f.api.reserve).not.toHaveBeenCalled(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  });
  it('compares reviewed layouts and safe destinations without borrowing the edited target as the original table', async () => {
    const { f, controller } = setup([{ id: 'C', revision: 'C1', name: 'Original destination' }, { id: 'C2', revision: 'C2R', name: 'New destination' }]);
    Object.assign(f.seed, { layout_version: 'managed_long_v1', storage_mode: 'plain', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' } });
    f.seed.destination.table_name = 'original_long_table';
    f.api.review.mockImplementation(async (_id, request) => ({ ...f.review, envelope: { ...f.envelope, request } }));
    fireEvent.click(screen.getByText('Open change'));
    fireEvent.change(await screen.findByLabelText('dialog.layout'), { target: { value: 'managed_wide_profile_v1' } });
    fireEvent.change(screen.getByLabelText('dialog.destination'), { target: { value: 'C2' } });
    fireEvent.click(screen.getByTestId('recording-change-review')); await screen.findByTestId('recording-change-reserve_preview');
    expect(screen.getByTestId('recording-change-original-layout')).toHaveTextContent('dialog.layout_long');
    expect(screen.getByTestId('recording-change-new-layout')).toHaveTextContent('dialog.layout_wide');
    expect(screen.getByTestId('recording-change-original-table')).toHaveTextContent('original_long_table');
    expect(screen.getByTestId('recording-change-new-table')).toHaveTextContent('dialog.new_table_pending');
    expect(screen.getByTestId('recording-change-original-destination')).toHaveTextContent('Original destination');
    expect(screen.getByTestId('recording-change-new-destination')).toHaveTextContent('New destination');
    expect(controller.snapshot().target?.destination.table_name).toBe('');
    expect(screen.getByText('dialog.history_kept')).toBeVisible(); expect(f.api.reserve).not.toHaveBeenCalled(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  });
  it('states missing table names honestly until a matching preparation returns the new table', async () => {
    const { f } = setup();
    fireEvent.click(screen.getByText('Open change')); fireEvent.click(await screen.findByTestId('recording-change-review'));
    await screen.findByTestId('recording-change-reserve_preview');
    expect(screen.getByTestId('recording-change-original-table')).toHaveTextContent('dialog.original_table_unavailable');
    expect(screen.getByTestId('recording-change-new-table')).toHaveTextContent('dialog.new_table_pending');
    expect(screen.queryByText('recording_G1')).toBeNull(); expect(f.api.reserve).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('recording-change-reserve_preview')); await screen.findByTestId('recording-change-confirm_schema');
    expect(screen.getByTestId('recording-change-new-table')).toHaveTextContent('recording_G2');
    expect(screen.getByTestId('recording-change-original-table')).toHaveTextContent('dialog.original_table_unavailable');
    expect(screen.getByText('dialog.history_kept')).toBeVisible(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  });
  it.each(['source', 'preparation'] as const)('does not present a foreign %s table as matching impact evidence', async foreign => {
    const { f, controller, unmount } = setup();
    await act(async () => { await controller.requestAction(f.intent); await controller.review(); await controller.reserveAndPreview(); });
    const state = structuredClone(controller.snapshot());
    if (foreign === 'source' && state.lifecycle?.successor_seed) {
      state.lifecycle.group_id = 'OTHER'; state.lifecycle.successor_seed.group.destination.table_name = 'foreign_table';
    }
    if (foreign === 'preparation' && state.recovery?.preparation) {
      state.recovery.preparation.operation_id = 'OTHER'; state.recovery.preparation.group.destination.table_name = 'foreign_table';
      state.review!.envelope.request.group.destination.table_name = ''; state.preview = undefined;
    }
    unmount(); render(<RecordingChangeDialog controller={{ ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: true }} />);
    expect(screen.getByTestId(`recording-change-${foreign === 'source' ? 'original' : 'new'}-table`)).toHaveTextContent(foreign === 'source' ? 'dialog.original_table_unavailable' : 'dialog.new_table_pending');
    expect(screen.queryByText('foreign_table')).toBeNull(); expect(f.api.reserve).toHaveBeenCalledOnce(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  });
});
