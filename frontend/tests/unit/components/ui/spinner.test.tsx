import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Spinner } from '../../../../src/components/ui/spinner';

describe('Spinner', () => {
  it('renders with default size', () => {
    render(<Spinner />);
    const spinner = screen.getByRole('status');
    expect(spinner).toBeInTheDocument();
    expect(spinner).toHaveAttribute('aria-label', '載入中');
  });

  it('renders with small size', () => {
    render(<Spinner size="sm" />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('h-4', 'w-4');
  });

  it('renders with medium size', () => {
    render(<Spinner size="md" />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('h-6', 'w-6');
  });

  it('renders with large size', () => {
    render(<Spinner size="lg" />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('h-8', 'w-8');
  });

  it('applies custom className', () => {
    render(<Spinner className="custom-class" />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('custom-class');
  });

  it('has spinning animation', () => {
    render(<Spinner />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('animate-spin');
  });
});
