import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { GroupEditorDestinationSection, type GroupEditorDestinationSectionProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorDestinationSection';
import { GroupEditorMethodSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection';
import { newMixedDraft } from '@/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { createBasicRecordingRequestId } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { DbConnector } from '@/features/datalink/workbench-v2/state/types';
import type { EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';

const destination = { connector_id: 'c1', connector_revision: 'r1', table_schema: 'public', table_name: '' };
const mixedDraft = (): EditorDraft => ({ ...newMixedDraft(destination, createBasicRecordingRequestId()), name: 'Test' });
// Not persisted: the destination catalog stays idle, so no request leaves the test.
const connector = { connector_id: 'c1', identity_revision: 'r1', kind: 'postgres', name: 'Plant DB', schema: 'public', table: '', persisted: false } as unknown as DbConnector;

function renderWithQuery(node: React.ReactNode) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}>{node}</QueryClientProvider>);
}

function destinationProps(overrides: Partial<GroupEditorDestinationSectionProps> = {}): GroupEditorDestinationSectionProps {
  const draft = mixedDraft();
  return { presentation: 'simple', scopedPresentation: false, workspaceId: 'w1', workspaceRevision: 'wr1', connector, group: undefined, saved: false,
    destinationMatches: true, draft, readonly: false, disabled: false, saveUnknown: false, saveReadbackPending: false, metadataStatus: 'exists', onRetryMetadata: vi.fn(),
    configured: { context: { workspaceId: 'w1', connector, draft }, invalid: false, onChoose: vi.fn() }, ownedWideConflict: false,
    canChooseNewDestination: () => true, onChange: vi.fn(), onStorageStrategy: vi.fn(), onReplaceDraft: vi.fn(), ...overrides };
}

beforeEach(async () => { await i18n.changeLanguage('zh-TW'); });
afterEach(cleanup);

describe('recording editor constraints namespace', () => {
  it('is registered for both operator languages', () => {
    expect(i18n.hasResourceBundle('en', 'recording-editor-constraints')).toBe(true);
    expect(i18n.hasResourceBundle('zh-TW', 'recording-editor-constraints')).toBe(true);
  });
});

describe('GroupEditorMethodSection', () => {
  it('shows the per-point method editor without the profile summary for a mixed record', () => {
    renderWithQuery(<GroupEditorMethodSection presentation="simple" draft={mixedDraft()} saved={false} disabled={false} dialect="postgres" onReplaceDraft={vi.fn()} />);
    expect(screen.getByTestId('recording-mixed-intent')).toBeInTheDocument();
    expect(screen.queryByTestId('group-profile-summary')).not.toBeInTheDocument();
  });

  it('adds the profile summary for an instantaneous profile record', () => {
    const draft = { ...mixedDraft(), layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' } } as unknown as EditorDraft;
    renderWithQuery(<GroupEditorMethodSection presentation="simple" draft={draft} saved disabled={false} dialect="postgres" onReplaceDraft={vi.fn()} />);
    expect(screen.getByTestId('group-profile-summary')).toBeInTheDocument();
    expect(screen.getByTestId('group-profile-interval')).toHaveTextContent('15');
  });
});

describe('GroupEditorDestinationSection', () => {
  it('assembles the simple destination fields with storage fixed for a new per-point record', () => {
    renderWithQuery(<GroupEditorDestinationSection {...destinationProps()} />);
    expect(screen.getByTestId('group-storage-managed')).toBeDisabled();
    expect(screen.getByTestId('group-schema')).toBeEnabled();
    expect(screen.getByTestId('recording-mixed-destination')).toBeInTheDocument();
    expect(screen.queryByTestId('recording-destination-metadata-note')).not.toBeInTheDocument();
  });

  it('explains an unreadable destination table and retries the same inspection', () => {
    const onRetryMetadata = vi.fn();
    renderWithQuery(<GroupEditorDestinationSection {...destinationProps({ metadataStatus: 'failed', onRetryMetadata })} />);
    expect(screen.getByTestId('recording-destination-metadata-note')).toHaveTextContent('無法確認目的資料表');
    fireEvent.click(screen.getByRole('button', { name: '重試' }));
    expect(onRetryMetadata).toHaveBeenCalledTimes(1);
  });

  it('offers the current destination when the draft points elsewhere', () => {
    const onChange = vi.fn();
    const draft = { ...mixedDraft(), connector_id: 'other', connector_revision: 'old' };
    renderWithQuery(<GroupEditorDestinationSection {...destinationProps({ draft, saved: true, destinationMatches: false, onChange })} />);
    fireEvent.click(screen.getByTestId('group-use-current-destination'));
    expect(onChange).toHaveBeenCalledWith({ connector_id: 'c1', connector_revision: 'r1', table_schema: 'public', table_name: '' });
  });

  it('reports an owned wide table as a blocking alert', () => {
    renderWithQuery(<GroupEditorDestinationSection {...destinationProps({ ownedWideConflict: true })} />);
    expect(screen.getByRole('alert')).toHaveTextContent(i18n.t('recording-editor:destinations.reason.wideOwner'));
  });
});
