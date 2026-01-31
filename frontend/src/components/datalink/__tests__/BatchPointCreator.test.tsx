
import { render, screen, fireEvent } from '@testing-library/react';
import { BatchPointCreator } from '../BatchPointCreator';
import { describe, it, expect, vi } from 'vitest';
import { ToastProvider } from '../../../contexts/ToastContext';

const renderWithProviders = (ui: React.ReactElement) => {
  return render(<ToastProvider>{ui}</ToastProvider>);
};

describe('BatchPointCreator', () => {
  const defaultProps = {
    deviceId: 'dev1',
    protocol: 'modbus_tcp' as const,
    preselectedAddresses: ['40001', '40002', '40003'],
    pollingGroups: [{ id: 'pg1', name: 'Fast Poll' }],
    onCreated: vi.fn(),
    onCancel: vi.fn(),
  };

  it('should render form with preselected addresses', () => {
    renderWithProviders(<BatchPointCreator {...defaultProps} />);
    expect(screen.getByTestId('batch-point-creator')).toBeInTheDocument();
  });

  it('should show preview of generated points', async () => {
    renderWithProviders(<BatchPointCreator {...defaultProps} />);
    
    // We expect input for template and quantity
    // Template: Pump_{index}
    // Result: Pump_0, Pump_1, Pump_2
    
    // This will fail until implementation
    const templateInput = screen.getByLabelText(/命名模板/i);
    fireEvent.change(templateInput, { target: { value: 'Pump_{index}' } });
    
    expect(screen.getByText('Pump_0')).toBeInTheDocument();
    expect(screen.getByText('40001')).toBeInTheDocument();
  });
});
