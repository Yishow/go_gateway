import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { resetGatewayDraftStore } from '../../../features/gateway/gatewayDraftStore';
import GatewayEntryPage from '../GatewayEntryPage';
import GatewayExpertWorkbenchPage from '../GatewayExpertWorkbenchPage';
import GatewayQuickSetupPage from '../GatewayQuickSetupPage';

const connectMock = vi.fn();

vi.mock('../../../services/api', () => ({
  useTestAPI: () => ({
    connect: connectMock,
  }),
}));

describe('GatewayQuickSetupPage', () => {
  const renderComponent = () => {
    render(
      <MemoryRouter>
        <GatewayQuickSetupPage />
      </MemoryRouter>
    );
  };

  beforeEach(() => {
    resetGatewayDraftStore();
    connectMock.mockReset();
    connectMock.mockResolvedValue({
      connection_id: 'conn-001',
      status: 'connected',
    });
  });

  it('supports step navigation with next/back controls', () => {
    renderComponent();

    expect(screen.getByRole('heading', { name: 'Step1 設備連線' })).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('step-next-btn'));
    expect(screen.getByRole('heading', { name: 'Step2 路由/來源規劃' })).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('step-back-btn'));
    expect(screen.getByRole('heading', { name: 'Step1 設備連線' })).toBeInTheDocument();
    expect(screen.getByTestId('host-input')).toBeInTheDocument();
    expect(screen.getByTestId('payload-preview')).toBeInTheDocument();
  });

  it('blocks moving to next step when required fields are invalid', () => {
    renderComponent();

    fireEvent.click(screen.getByTestId('step-next-btn'));
    expect(screen.getByRole('heading', { name: 'Step2 路由/來源規劃' })).toBeInTheDocument();

    fireEvent.change(screen.getByTestId('route-input'), { target: { value: '' } });
    fireEvent.change(screen.getByTestId('source-input'), { target: { value: '' } });
    fireEvent.click(screen.getByTestId('step-next-btn'));

    expect(screen.getByTestId('step-error-message')).toHaveTextContent('請輸入 API 路徑');
    expect(screen.getByRole('heading', { name: 'Step2 路由/來源規劃' })).toBeInTheDocument();
  });

  it('retains entered state when switching between steps', () => {
    renderComponent();

    fireEvent.change(screen.getByTestId('host-input'), { target: { value: '10.0.0.1' } });
    fireEvent.click(screen.getByTestId('step-next-btn'));
    fireEvent.click(screen.getByTestId('step-back-btn'));

    expect(screen.getByTestId('host-input')).toHaveValue('10.0.0.1');
    const payloadContent = JSON.parse(screen.getByTestId('payload-preview').textContent || '{}');
    expect(payloadContent.config.host).toBe('10.0.0.1');
  });

  it('enforces submit conditions before calling connect', async () => {
    renderComponent();

    fireEvent.click(screen.getByTestId('step-next-btn'));
    fireEvent.click(screen.getByTestId('step-next-btn'));

    fireEvent.change(screen.getByTestId('tag-input'), { target: { value: 'temperature' } });
    fireEvent.click(screen.getByTestId('security-review-checkbox'));
    fireEvent.click(screen.getByTestId('step-next-btn'));

    fireEvent.click(screen.getByTestId('validation-confirm-checkbox'));
    fireEvent.click(screen.getByTestId('step-next-btn'));

    const saveButton = screen.getByTestId('save-btn');
    expect(saveButton).toBeDisabled();
    fireEvent.click(saveButton);
    expect(connectMock).not.toHaveBeenCalled();

    fireEvent.click(screen.getByTestId('submit-confirm-checkbox'));
    await waitFor(() => expect(screen.getByTestId('save-btn')).toBeEnabled());
    fireEvent.click(screen.getByTestId('save-btn'));

    await waitFor(() => expect(connectMock).toHaveBeenCalledTimes(1));
  });

  it('shares quick draft state with expert payload preview', () => {
    const { unmount } = render(
      <MemoryRouter>
        <GatewayQuickSetupPage />
      </MemoryRouter>,
    );

    fireEvent.change(screen.getByTestId('host-input'), { target: { value: '10.20.30.40' } });
    unmount();

    render(
      <MemoryRouter>
        <GatewayExpertWorkbenchPage />
      </MemoryRouter>,
    );

    expect(screen.getByTestId('expert-payload-preview')).toHaveTextContent('10.20.30.40');
  });

  it('keeps draft data when switching quick -> entry -> expert -> entry -> quick', () => {
    render(
      <MemoryRouter initialEntries={['/gateway/quick-setup']}>
        <Routes>
          <Route path="/gateway/entry" element={<GatewayEntryPage />} />
          <Route path="/gateway/quick-setup" element={<GatewayQuickSetupPage />} />
          <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchPage />} />
        </Routes>
      </MemoryRouter>,
    );

    fireEvent.change(screen.getByTestId('host-input'), { target: { value: '172.16.1.20' } });
    fireEvent.click(screen.getByText('返回模式選擇'));
    fireEvent.click(screen.getByTestId('gateway-entry-expert-link'));
    fireEvent.click(screen.getByTestId('expert-back-entry-link'));
    fireEvent.click(screen.getByTestId('gateway-entry-quick-link'));

    expect(screen.getByTestId('host-input')).toHaveValue('172.16.1.20');
  });
});
