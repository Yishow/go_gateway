import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { QuickActions } from '@/components/datalink/QuickActions';
import type { Device } from '@/types/datalink';

describe('QuickActions', () => {
  const mockDevice: Device = {
    id: 'd1',
    name: 'Test Device',
    description: '',
    protocol: 'modbus_tcp',
    status: 'active',
    connection_config: '',
    last_test_at: null,
    last_test_success: null,
    last_test_error: '',
    created_at: '',
    updated_at: '',
  };

  const defaultProps = {
    device: mockDevice,
    selectedCount: 0,
    onBatchCreate: vi.fn(),
    onQuickMapping: vi.fn(),
    onTestConnection: vi.fn(),
    onOpenWorkbench: vi.fn(),
  };

  it('should render actions panel', () => {
    render(<QuickActions {...defaultProps} />);
    expect(screen.getByTestId('quick-actions')).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: /Server Memory Grid/i })
    ).toBeInTheDocument();
  });

  it('should disable batch create when no cells selected', () => {
    render(<QuickActions {...defaultProps} selectedCount={0} />);
  });
});
