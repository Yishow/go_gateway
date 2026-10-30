import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import RecordingStylePreview from '@/dev/RecordingStylePreview';
import { I18nextProvider } from 'react-i18next';
import { previewI18n, previewTranslationsReady } from '@/dev/previewI18n';
await previewTranslationsReady;
import { settingsAPI } from '@/services/datalink';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';

afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });
describe('isolated actual recording style preview', () => {
  it('exercises actual home, inputs, statuses and portal with no network or persisted preference writes', async () => {
    const fetch = vi.fn(); vi.stubGlobal('fetch', fetch);
    const xhr = vi.spyOn(XMLHttpRequest.prototype, 'open');
    const storage = vi.spyOn(Storage.prototype, 'setItem');
    const update = vi.spyOn(settingsAPI, 'updateKey');
    const commands = (['create', 'update', 'apply', 'schemaApply', 'remove'] as const).map(key => vi.spyOn(studioV2WorkspaceWriteGroupsAPI, key));
    render(<I18nextProvider i18n={previewI18n}><RecordingStylePreview /></I18nextProvider>);
    expect(screen.getByTestId('recording-home')).toBeInTheDocument();
    expect(screen.getByText(/SIMULATED · No backend/)).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('Preview theme'), { target: { value: 'dark' } });
    expect(document.documentElement).toHaveClass('dark');
    fireEvent.change(screen.getByLabelText('Preview state'), { target: { value: 'error' } });
    expect(screen.getByTestId('recording-home-guidance')).toHaveAttribute('data-guidance-state', 'error');
    fireEvent.change(screen.getByLabelText('Sample input'), { target: { value: 'sample' } });
    expect(screen.getByRole('button', { name: 'Disabled action' })).toBeDisabled();
    fireEvent.click(screen.getByRole('button', { name: 'Preview confirmation' }));
    expect(await screen.findByRole('alertdialog')).toHaveClass('gw-modal');
    fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/password|credential/i)).not.toBeInTheDocument();
    expect(fetch).not.toHaveBeenCalled(); expect(xhr).not.toHaveBeenCalled();
    expect(update).not.toHaveBeenCalled(); expect(storage).not.toHaveBeenCalled();
    commands.forEach(command => expect(command).not.toHaveBeenCalled());
  });
  it('renders each explicitly simulated state without implying a live workspace', () => {
    render(<I18nextProvider i18n={previewI18n}><RecordingStylePreview /></I18nextProvider>);
    for (const value of ['loading', 'error', 'unknown', 'missing_devices', 'draft', 'source_attention', 'recordings']) {
      fireEvent.change(screen.getByLabelText('Preview state'), { target: { value } });
      expect(screen.getByTestId('recording-home-guidance')).toBeInTheDocument();
      expect(screen.getByText(/SIMULATED · No backend/)).toBeInTheDocument();
    }
    fireEvent.click(screen.getByRole('button', { name: 'Preview confirmation' }));
    const dialog = screen.getByRole('alertdialog');
    expect(dialog.parentElement).toBe(document.body);
    fireEvent.click(screen.getByTestId('group-delete-cancel'));
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  });
  it('restores its own document theme on unmount and supports light portal styling without a theme owner', () => {
    document.documentElement.className = 'light preserved';
    const { unmount } = render(<I18nextProvider i18n={previewI18n}><RecordingStylePreview /></I18nextProvider>);
    fireEvent.change(screen.getByLabelText('Preview theme'), { target: { value: 'dark' } });
    fireEvent.change(screen.getByLabelText('Preview theme'), { target: { value: 'light' } });
    expect(document.documentElement).toHaveClass('light');
    unmount(); expect(document.documentElement.className).toBe('light preserved');
  });
});
