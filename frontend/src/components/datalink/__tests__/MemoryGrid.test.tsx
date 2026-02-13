
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryGrid } from '../MemoryGrid';
import { describe, it, expect, vi } from 'vitest';
import type { Point } from '../../../types/datalink';

describe('MemoryGrid', () => {
  const createPoint = (address: string, name = 'P1'): Point => ({
    id: `point-${address}`,
    device_id: 'd1',
    name,
    description: '',
    data_type: 'int16',
    address,
    enabled: true,
    polling_group_id: 'pg-1',
    last_value: null,
    last_read_at: '',
    last_error: '',
    error_count: 0,
    created_at: '',
    updated_at: '',
  });

  const defaultProps = {
    deviceId: 'dev1',
    protocol: 'modbus_tcp' as const,
    centerAddress: '40001',
    existingPoints: [],
    selectedAddresses: [],
    onSelect: vi.fn(),
    onCellClick: vi.fn(),
  };

  it('should render grid container', () => {
    render(<MemoryGrid {...defaultProps} />);
    expect(screen.getByTestId('memory-grid')).toBeInTheDocument();
  });

  it('should handle single selection on click', () => {
    const onSelect = vi.fn();
    render(<MemoryGrid {...defaultProps} range={10} onSelect={onSelect} />);
    
    // Click header
    const cells = screen.getAllByTestId('grid-cell');
    fireEvent.click(cells[0]);
    
    // Should select just that one
    expect(onSelect).toHaveBeenCalledWith(expect.arrayContaining(['40001']));
    expect(onSelect.mock.calls[0][0]).toHaveLength(1);
  });

  it('should handle multi-selection with Ctrl click', () => {
    const onSelect = vi.fn();
    // Start with one selected
    render(<MemoryGrid 
      {...defaultProps} 
      range={10} 
      onSelect={onSelect} 
      selectedAddresses={['40001']} 
    />);
    
    const cells = screen.getAllByTestId('grid-cell');
    
    // Ctrl+Click second cell
    fireEvent.click(cells[1], { ctrlKey: true });
    
    // Should keep 40001 and add 40002
    const lastCall = onSelect.mock.calls[onSelect.mock.calls.length - 1];
    expect(lastCall[0]).toContain('40001');
    expect(lastCall[0]).toContain('40002');
    expect(lastCall[0]).toHaveLength(2);
  });

  it('should handle range selection with Shift click', () => {
    const onSelect = vi.fn();
    // Start with 40001 selected
    render(<MemoryGrid 
      {...defaultProps} 
      range={10} 
      onSelect={onSelect} 
      selectedAddresses={['40001']} 
    />);
    
    const cells = screen.getAllByTestId('grid-cell');
    
    // Shift+Click 40004 (index 3)
    // Should select 40001, 40002, 40003, 40004
    fireEvent.click(cells[3], { shiftKey: true });
    
    const lastCall = onSelect.mock.calls[onSelect.mock.calls.length - 1];
    expect(lastCall[0]).toHaveLength(4);
    expect(lastCall[0]).toContain('40001');
    expect(lastCall[0]).toContain('40004');
  });

  it('should call onCellClick with point info', () => {
    const onCellClick = vi.fn();
    const point = createPoint('40001');
    
    render(<MemoryGrid 
      {...defaultProps} 
      range={10} 
      existingPoints={[point]}
      onCellClick={onCellClick} 
    />);
    
    const cells = screen.getAllByTestId('grid-cell');
    fireEvent.click(cells[0]);
    
    expect(onCellClick).toHaveBeenCalledWith('40001', point);
  });

  it('should match int16 count to exact planned cells', () => {
    const plannedAllocations = Array.from({ length: 5 }).map((_, index) => ({
      id: `int-${index}`,
      dataType: 'int16' as const,
      addresses: [`${40001 + index}`],
      label: `INT${index + 1}`,
    }));

    render(
      <MemoryGrid
        {...defaultProps}
        range={30}
        plannedAllocations={plannedAllocations}
      />
    );

    const intLabels = screen.getAllByText(/^INT\d+$/);
    expect(intLabels).toHaveLength(5);
  });

  it('should render float32 planned count as paired cells', () => {
    const plannedAllocations = Array.from({ length: 10 }).map((_, index) => {
      const base = 40001 + index * 2;
      return {
        id: `float-${index}`,
        dataType: 'float32' as const,
        addresses: [`${base}`, `${base + 1}`],
        label: `F${index + 1}`,
      };
    });

    render(
      <MemoryGrid
        {...defaultProps}
        range={50}
        plannedAllocations={plannedAllocations}
      />
    );

    const pairedMarkers = screen.getAllByText(/\/2$/);
    expect(pairedMarkers).toHaveLength(20);
  });

  it('should render linked occupancy state', () => {
    render(
      <MemoryGrid
        {...defaultProps}
        range={5}
        linkedAddresses={['40001']}
      />
    );

    const address = screen.getByText('40001');
    const cell = address.closest('[data-testid="grid-cell"]');
    expect(cell).toHaveAttribute('data-status', 'linked');
  });

  it('should classify conflict severity as hard when planned overlaps used point', () => {
    const point = createPoint('40001');
    const plannedAllocations = [
      {
        id: 'plan-1',
        dataType: 'int16' as const,
        addresses: ['40001'],
        label: 'S1',
      },
    ];

    render(
      <MemoryGrid
        {...defaultProps}
        range={5}
        existingPoints={[point]}
        plannedAllocations={plannedAllocations}
      />
    );

    const address = screen.getByText('40001');
    const cell = address.closest('[data-testid="grid-cell"]');
    expect(cell).toHaveAttribute('data-status', 'conflict');
    expect(cell).toHaveAttribute('data-conflict-severity', 'hard');
    expect(screen.getByText('硬衝突')).toBeInTheDocument();
  });

  it('should classify conflict severity as soft when planned overlaps linked address', () => {
    const plannedAllocations = [
      {
        id: 'plan-1',
        dataType: 'int16' as const,
        addresses: ['40001'],
        label: 'S1',
      },
    ];

    render(
      <MemoryGrid
        {...defaultProps}
        range={5}
        linkedAddresses={['40001']}
        plannedAllocations={plannedAllocations}
      />
    );

    const address = screen.getByText('40001');
    const cell = address.closest('[data-testid="grid-cell"]');
    expect(cell).toHaveAttribute('data-status', 'conflict');
    expect(cell).toHaveAttribute('data-conflict-severity', 'soft');
    expect(screen.getByText('軟衝突')).toBeInTheDocument();
  });

  it('should handle edge collision when planned cells exceed visible range', () => {
    const plannedAllocations = [
      {
        id: 'float-edge',
        dataType: 'float32' as const,
        addresses: ['40003', '40004'],
        label: 'EDGE',
      },
    ];

    render(
      <MemoryGrid
        {...defaultProps}
        range={3}
        plannedAllocations={plannedAllocations}
      />
    );

    expect(screen.getByText('40003')).toBeInTheDocument();
    expect(screen.queryByText('40004')).not.toBeInTheDocument();
    expect(screen.getByText('EDGE 1/2')).toBeInTheDocument();
  });

  it('should keep neighbor context when conflicts-only filter is enabled', () => {
    const point = createPoint('40002');
    const plannedAllocations = [
      {
        id: 'plan-1',
        dataType: 'int16' as const,
        addresses: ['40002'],
        label: 'S1',
      },
    ];

    render(
      <MemoryGrid
        {...defaultProps}
        range={5}
        existingPoints={[point]}
        plannedAllocations={plannedAllocations}
        showConflictsOnly
      />
    );

    expect(screen.getByText('40001')).toBeInTheDocument();
    expect(screen.getByText('40002')).toBeInTheDocument();
    expect(screen.getByText('40003')).toBeInTheDocument();
    expect(screen.queryByText('40004')).not.toBeInTheDocument();
  });

  it('should render cells as keyboard-focusable buttons', () => {
    render(<MemoryGrid {...defaultProps} range={2} />);
    const cells = screen.getAllByTestId('grid-cell');
    expect(cells[0].tagName).toBe('BUTTON');
    expect(cells[0]).toHaveAttribute('aria-label', expect.stringContaining('Address 40001'));
  });
});
