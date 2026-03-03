import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import {
  getGatewayDraftState,
  resetGatewayDraftStore,
  updateGatewayExpertDraft,
} from '../../../features/gateway/gatewayDraftStore';
import GatewayEntryPage from '../GatewayEntryPage';
import GatewayQuickSetupPage from '../GatewayQuickSetupPage';
import GatewayExpertWorkbenchPage from '../GatewayExpertWorkbenchPage';

describe('GatewayExpertWorkbenchPage', () => {
  beforeEach(() => {
    resetGatewayDraftStore();
    vi.restoreAllMocks();
  });

  function renderComponent() {
    render(
      <MemoryRouter>
        <GatewayExpertWorkbenchPage />
      </MemoryRouter>,
    );
  }

  it('renders core expert workbench sections', () => {
    renderComponent();

    expect(screen.getByText('Expert Workbench')).toBeInTheDocument();
    expect(screen.getByTestId('expert-route-table')).toBeInTheDocument();
    expect(screen.getByTestId('expert-plugin-chain')).toBeInTheDocument();
    expect(screen.getByTestId('expert-raw-manifest')).toBeInTheDocument();
    expect(screen.getByTestId('expert-payload-preview')).toBeInTheDocument();
  });

  it('adds a route row when clicking add route', () => {
    renderComponent();

    const before = screen.getAllByText('刪除').length;
    fireEvent.click(screen.getByText('+ 新增路由'));
    const after = screen.getAllByText('刪除').length;

    expect(after).toBe(before + 1);
  });

  it('shows manifest parse error when raw JSON is invalid', () => {
    renderComponent();

    const textarea = screen.getByTestId('expert-raw-manifest');
    fireEvent.change(textarea, { target: { value: '{ invalid-json }' } });

    expect(screen.getByTestId('expert-manifest-error')).toBeInTheDocument();
    expect(screen.getByTestId('expert-submit-btn')).toBeDisabled();
  });

  it('shows downgrade warning and blocks switching to quick when not confirmed', () => {
    updateGatewayExpertDraft({
      protocol: 'modbus-tcp',
      configJson: '{"host":"10.10.10.10","port":502,"routes":[{"path":"/x"}]}',
    });

    const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false);
    render(
      <MemoryRouter initialEntries={['/gateway/expert-workbench']}>
        <Routes>
          <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchPage />} />
          <Route path="/gateway/quick-setup" element={<GatewayQuickSetupPage />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByTestId('expert-downgrade-warning')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('expert-switch-quick-link'));

    expect(confirmSpy).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('expert-route-table')).toBeInTheDocument();
    expect(screen.queryByTestId('host-input')).not.toBeInTheDocument();
  });

  it('downgrades and switches to quick after confirm', () => {
    updateGatewayExpertDraft({
      protocol: 'modbus-tcp',
      configJson: '{"host":"10.0.9.9","port":1502,"routes":[{"path":"/x"}]}',
    });

    vi.spyOn(window, 'confirm').mockReturnValue(true);
    render(
      <MemoryRouter initialEntries={['/gateway/expert-workbench']}>
        <Routes>
          <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchPage />} />
          <Route path="/gateway/quick-setup" element={<GatewayQuickSetupPage />} />
        </Routes>
      </MemoryRouter>,
    );

    fireEvent.click(screen.getByTestId('expert-switch-quick-link'));

    expect(screen.getByTestId('host-input')).toHaveValue('10.0.9.9');
    expect(getGatewayDraftState().expertDraft.configJson).not.toContain('routes');
    expect(getGatewayDraftState().expertQuickCompatibility).toBe('compatible');
  });

  it('keeps expert builder draft when switching expert -> entry -> expert', () => {
    render(
      <MemoryRouter initialEntries={['/gateway/expert-workbench']}>
        <Routes>
          <Route path="/gateway/entry" element={<GatewayEntryPage />} />
          <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchPage />} />
        </Routes>
      </MemoryRouter>,
    );

    fireEvent.click(screen.getByText('+ 新增路由'));
    expect(screen.getAllByText('刪除').length).toBe(3);

    fireEvent.click(screen.getByTestId('expert-back-entry-link'));
    fireEvent.click(screen.getByTestId('gateway-entry-expert-link'));

    expect(screen.getAllByText('刪除').length).toBe(3);
  });
});
