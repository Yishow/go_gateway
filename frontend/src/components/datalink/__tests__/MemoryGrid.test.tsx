
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

  // Tests expected to fail until implemented
  it('should render grid cells', () => {
    // We expect it to render cells based on range (default 100 potentially)
    render(<MemoryGrid {...defaultProps} range={10} />);
    // This assumes implementation will use data-testid="grid-cell"
    const cells = screen.queryAllByTestId('grid-cell');
    expect(cells.length).toBeGreaterThan(0);
  });
  
  it('should handle cell click', () => {
    const onCellClick = vi.fn();
    render(<MemoryGrid {...defaultProps} range={10} onCellClick={onCellClick} />);
    const cell = screen.queryAllByTestId('grid-cell')[0];
    if (cell) {
        fireEvent.click(cell);
        expect(onCellClick).toHaveBeenCalled();
    } else {
        // Fail if no cells
        expect(true).toBe(false); 
    }
  });
});
