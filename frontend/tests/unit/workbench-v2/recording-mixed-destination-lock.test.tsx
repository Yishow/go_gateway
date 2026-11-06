import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { useState } from 'react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { destinationLock } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorCapabilities';
import { GroupIdentityFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields';
import { newMixedDraft } from '@/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { createBasicRecordingRequestId } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';

const destination = { connector_id: 'c1', connector_revision: 'r1', table_schema: 'public', table_name: '' };
const mixed = (layout: 'managed_wide_profile_v1' | 'managed_long_v1') => ({ ...newMixedDraft(destination, createBasicRecordingRequestId(), layout), name: 'Test' });
const catalog = { status: 'ready' as const, defaultSchema: 'public', retry: vi.fn(),
  schemas: [{ name: 'public', tables: ['existing_a', 'owned_b'] }, { name: 'plant', tables: ['line1'] }], blockedTables: [{ schema: 'public', name: 'owned_b' }] };

beforeEach(async () => { await i18n.changeLanguage('zh-TW'); });
afterEach(cleanup);

describe('per-point recording destination before the first save', () => {
  it('keeps schema and table editable until saved, with storage fixed to managed', () => {
    for (const layout of ['managed_wide_profile_v1', 'managed_long_v1'] as const) {
      expect(destinationLock(mixed(layout), false)).toEqual({ locked: false, storageLocked: true });
      expect(destinationLock(mixed(layout), true)).toEqual({ locked: true, storageLocked: true });
    }
  });

  it('keeps saved profile and long-table records locked as before', () => {
    const profile = { ...mixed('managed_long_v1'), recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' } } as unknown as EditorDraft;
    expect(destinationLock(profile, false)).toEqual({ locked: true, storageLocked: true });
    const legacy = { ...mixed('managed_long_v1'), layout_version: undefined, recording_profile: undefined } as unknown as EditorDraft;
    expect(destinationLock(legacy, false)).toEqual({ locked: false, storageLocked: false });
  });

  it('lets a new per-point record pick or type schema and table, but not an owned table', () => {
    const onChange = vi.fn();
    function DestinationHarness() {
      const [draft, setDraft] = useState(() => mixed('managed_wide_profile_v1'));
      return <GroupIdentityFields section="destination" scopedPresentation draft={draft} readonly={false} {...destinationLock(draft, false)} catalog={catalog} onChange={patch => {
        onChange(patch); setDraft(current => ({ ...current, ...patch }));
      }} onStorageStrategy={vi.fn()} />;
    }
    render(<DestinationHarness />);
    expect(screen.getByTestId('group-storage-managed')).toBeDisabled();
    expect(screen.getByTestId('group-storage-mixed-note')).toHaveTextContent('受管資料表');
    const schema = screen.getByTestId('group-schema'), table = screen.getByTestId('group-table');
    expect(schema).toBeEnabled(); expect(table).toBeEnabled();
    fireEvent.change(schema, { target: { value: 'listed:plant' } });
    expect(onChange).toHaveBeenLastCalledWith({ table_schema: 'plant' });
    expect(schema).toHaveValue('listed:plant');
    expect(screen.getByRole('option', { name: 'line1' })).toBeEnabled();
    fireEvent.change(schema, { target: { value: 'manual' } });
    const manualSchema = screen.getByTestId('group-schema-manual');
    expect(manualSchema).toBeEnabled();
    fireEvent.change(manualSchema, { target: { value: 'my_new_schema' } });
    expect(onChange).toHaveBeenLastCalledWith({ table_schema: 'my_new_schema' });
    expect(manualSchema).toHaveValue('my_new_schema');
    expect(screen.queryByRole('option', { name: 'existing_a' })).not.toBeInTheDocument();
    fireEvent.change(schema, { target: { value: 'listed:public' } });
    expect(schema).toHaveValue('listed:public');
    fireEvent.change(table, { target: { value: 'listed:existing_a' } });
    expect(onChange).toHaveBeenLastCalledWith({ table_name: 'existing_a' });
    expect(screen.getByRole('option', { name: /owned_b/ })).toBeDisabled();
    fireEvent.change(table, { target: { value: 'manual' } });
    fireEvent.change(screen.getByTestId('group-table-manual'), { target: { value: 'my_new_table' } });
    expect(onChange).toHaveBeenLastCalledWith({ table_name: 'my_new_table' });
    expect(screen.getByTestId('group-table-manual')).toHaveValue('my_new_table');
  });
});
