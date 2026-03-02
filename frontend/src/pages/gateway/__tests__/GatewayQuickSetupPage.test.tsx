import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import GatewayQuickSetupPage from '../GatewayQuickSetupPage';

describe('GatewayQuickSetupPage', () => {
  const renderComponent = () => {
    render(
      <MemoryRouter>
        <GatewayQuickSetupPage />
      </MemoryRouter>
    );
  };

  it('renders the core components', () => {
    renderComponent();

    // Check IntentForm fields
    expect(screen.getByText('1. 設備連線 (Intent Form)')).toBeInTheDocument();
    expect(screen.getByTestId('protocol-select')).toBeInTheDocument();
    expect(screen.getByTestId('host-input')).toBeInTheDocument();
    expect(screen.getByTestId('port-input')).toBeInTheDocument();
    expect(screen.getByTestId('unit-id-input')).toBeInTheDocument(); // default is modbus-tcp

    // Check SimpleRouteBuilder fields
    expect(screen.getByText('2. 路由設定 (Route Builder)')).toBeInTheDocument();
    expect(screen.getByTestId('route-input')).toBeInTheDocument();

    // Check BasicAuthToggle
    expect(screen.getByText('3. 安全設定 (Auth)')).toBeInTheDocument();
    expect(screen.getByTestId('auth-checkbox')).toBeInTheDocument();

    // Check Payload Preview
    expect(screen.getByText('Payload 預覽 (供驗證)')).toBeInTheDocument();
    expect(screen.getByTestId('payload-preview')).toBeInTheDocument();
  });

  it('updates payload preview when changing fields', () => {
    renderComponent();

    const hostInput = screen.getByTestId('host-input');
    fireEvent.change(hostInput, { target: { value: '10.0.0.1' } });

    const portInput = screen.getByTestId('port-input');
    fireEvent.change(portInput, { target: { value: '8080' } });

    const routeInput = screen.getByTestId('route-input');
    fireEvent.change(routeInput, { target: { value: '/custom/api' } });

    const authCheckbox = screen.getByTestId('auth-checkbox');
    fireEvent.click(authCheckbox);

    const payloadPreview = screen.getByTestId('payload-preview');
    const payloadContent = JSON.parse(payloadPreview.textContent || '{}');

    expect(payloadContent).toMatchObject({
      protocol: 'modbus-tcp',
      config: {
        host: '10.0.0.1',
        port: 8080,
        unitID: 1,
        route: '/custom/api',
        auth: true,
      },
    });
  });

  it('conditionally renders unitID or station based on protocol', () => {
    renderComponent();

    const protocolSelect = screen.getByTestId('protocol-select');
    
    // Default Modbus: shows Unit ID
    expect(screen.getByTestId('unit-id-input')).toBeInTheDocument();
    expect(screen.queryByTestId('station-input')).not.toBeInTheDocument();

    // Change to Fatek
    fireEvent.change(protocolSelect, { target: { value: 'fatek-tcp' } });
    expect(screen.queryByTestId('unit-id-input')).not.toBeInTheDocument();
    expect(screen.getByTestId('station-input')).toBeInTheDocument();

    // Change to MC Protocol (neither)
    fireEvent.change(protocolSelect, { target: { value: 'mc-tcp' } });
    expect(screen.queryByTestId('unit-id-input')).not.toBeInTheDocument();
    expect(screen.queryByTestId('station-input')).not.toBeInTheDocument();
  });
});