import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { savedGroupState } from '../../fixtures/writeGroupState';
import type { DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => ({ data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] }, isError: false, isLoading: false, refetch: vi.fn() }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: vi.fn() }),
  useWriteGroupDeliveryQuery: () => ({ data: undefined }),
}));

function mount(saveState: DbConnector['save_state']) {
  const state = savedGroupState();
  state.db.connector = { ...state.db.connector, status: 'ready', save_state: saveState, name: 'Unsaved long name', save_error: 'password=SECRET /private/backend/path' };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const page = () => <QueryClientProvider client={client}><Step4Database state={state} dispatch={vi.fn()} onCommit={vi.fn()} /></QueryClientProvider>;
  const view = render(page());
  return { ...view, changeSaveState: (next: DbConnector['save_state']) => { state.db.connector = { ...state.db.connector, save_state: next }; view.rerender(page()); } };
}

describe('database connector save presentation', () => {
  it('shows a safe failed-save state, not historical ready, while preserving the editable draft', () => {
    mount('save-error');
    expect(screen.getByRole('alert')).toHaveTextContent('step4.connector_save_failed');
    expect(screen.getByLabelText('step4.field_name')).toHaveValue('Unsaved long name');
    expect(screen.queryByText('ready', { exact: true })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'step4.delivery_expand_btn' }));
    expect(screen.queryByText('step4.delivery_readiness_ready')).not.toBeInTheDocument();
    expect(document.body.textContent).not.toContain('SECRET');
    expect(document.body.textContent).not.toContain('/private/backend/path');
  });

  it('keeps existing recording draft values when a later connector save fails', () => {
    const view = mount('saved');
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '23' } });
    view.changeSaveState('save-error');
    expect(screen.getByTestId('step4-basic-recording')).toBeVisible();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(23);
    expect(screen.getByTestId('basic-recording-interval')).toBeDisabled();
  });

  it.each(['save-error', 'saving', 'draft-invalid'] as const)('keeps saved recording context visible but readonly during %s', (saveState) => {
    mount(saveState);
    expect(screen.getByTestId('step4-basic-recording')).toBeVisible();
    expect(screen.getByTestId('basic-recording-interval')).toBeDisabled();
    expect(screen.getByTestId('basic-recording-completeness')).toBeDisabled();
    expect(document.querySelector('[data-testid="basic-recording-start"]:not([disabled])')).toBeNull();
  });
});
