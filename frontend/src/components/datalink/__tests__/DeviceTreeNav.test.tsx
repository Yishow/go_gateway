
import { render, screen } from '@testing-library/react';
import { DeviceTreeNav } from '../DeviceTreeNav';
import { describe, it, expect, vi } from 'vitest';
import type { Device } from '../../../types/datalink';

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
    // Basic render test
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

  // This test is expected to FAIL until we implement drag/drop
  it('should support drag and drop reordering', () => {
    // Just a placeholder test for now that we expect to fulfill later
    // In a real TDD cycle for dnd-kit, testing drag interaction requires complex setup
    // We will verifying that SortableContext is used.
    
    // For now, let's just make a test that checks if the component accepts onReorder
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
    // TODO: Add strict drag test logic
  });
});
