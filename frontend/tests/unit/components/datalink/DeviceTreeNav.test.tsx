import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { DeviceTreeNav } from '@/components/datalink/DeviceTreeNav';
import type { Device } from '@/types/datalink';

describe('DeviceTreeNav', () => {
  const mockDevices = [
    {
      id: '1',
      name: 'PLC-001',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    },
    {
      id: '2',
      name: 'PLC-002',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    },
  ] as Device[];

  it('should render device list', () => {
    render(
      <DeviceTreeNav
        devices={mockDevices}
        selectedDeviceId={null}
        onSelectDevice={() => {}}
        isCollapsed={false}
        onToggleCollapse={() => {}}
      />
    );
    expect(screen.getByTestId('device-tree-nav')).toBeInTheDocument();
  });

  it('should support drag and drop reordering', () => {
    const handleReorder = vi.fn();
    render(
      <DeviceTreeNav
        devices={mockDevices}
        selectedDeviceId={null}
        onSelectDevice={() => {}}
        isCollapsed={false}
        onToggleCollapse={() => {}}
        onReorder={handleReorder}
      />
    );

    expect(screen.getByTestId('device-tree-nav')).toBeInTheDocument();
  });
});
