import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { SlidePanel } from '@/components/datalink/SlidePanel';

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
    expect(screen.getByText('Test Panel')).toBeInTheDocument();
  });

  it('should not render when closed (or be hidden)', () => {
    render(<SlidePanel {...defaultProps} isOpen={false} />);
    expect(screen.queryByTestId('slide-panel')).not.toBeInTheDocument();
  });
});
