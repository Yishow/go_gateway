
import { render, screen, fireEvent } from '@testing-library/react';
import { SlidePanel } from '../SlidePanel';
import { describe, it, expect, vi } from 'vitest';

describe('SlidePanel', () => {
  const defaultProps = {
    isOpen: true,
    title: 'Test Panel',
    onClose: vi.fn(),
    children: <div>Content</div>,
  };

  it('should render when open', () => {
    render(<SlidePanel {...defaultProps} />);
    expect(screen.getByTestId('slide-panel')).toBeInTheDocument();
    expect(screen.getByText('Test Panel')).toBeInTheDocument(); // Will fail until implementation
  });

  it('should not render when closed (or be hidden)', () => {
    // If we implement mounting/unmounting or CSS hiding
    // For now assuming CSS transform or conditional render
    const { container } = render(<SlidePanel {...defaultProps} isOpen={false} />);
    expect(screen.queryByTestId('slide-panel')).not.toBeInTheDocument();
  });
});
