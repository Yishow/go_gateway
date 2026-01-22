
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryGrid } from '../MemoryGrid';
import { describe, it, expect, vi } from 'vitest';

describe('MemoryGrid', () => {
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
    const point = { id: '1', name: 'P1', address: '40001', device_id: 'd1' } as any;
    
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
});
