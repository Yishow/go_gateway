import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Skeleton, SkeletonText, SkeletonCard } from '../../../../src/components/ui/skeleton';

describe('Skeleton', () => {
  it('renders with default variant', () => {
    const { container } = render(<Skeleton />);
    const skeleton = container.firstChild as HTMLElement;
    expect(skeleton).toHaveClass('rounded-lg');
    expect(skeleton).toHaveAttribute('aria-hidden', 'true');
  });

  it('renders with text variant', () => {
    const { container } = render(<Skeleton variant="text" />);
    const skeleton = container.firstChild as HTMLElement;
    expect(skeleton).toHaveClass('h-4', 'w-full');
  });

  it('renders with circular variant', () => {
    const { container } = render(<Skeleton variant="circular" />);
    const skeleton = container.firstChild as HTMLElement;
    expect(skeleton).toHaveClass('rounded-full');
  });

  it('applies custom className', () => {
    const { container } = render(<Skeleton className="custom-class" />);
    const skeleton = container.firstChild as HTMLElement;
    expect(skeleton).toHaveClass('custom-class');
  });

  it('has pulse animation', () => {
    const { container } = render(<Skeleton />);
    const skeleton = container.firstChild as HTMLElement;
    expect(skeleton).toHaveClass('animate-pulse');
  });
});

describe('SkeletonText', () => {
  it('renders default 3 lines', () => {
    const { container } = render(<SkeletonText />);
    // Count only the Skeleton components (not the wrapper div)
    const lines = container.querySelectorAll('.animate-pulse');
    expect(lines).toHaveLength(3);
  });

  it('renders custom number of lines', () => {
    const { container } = render(<SkeletonText lines={5} />);
    const lines = container.querySelectorAll('.animate-pulse');
    expect(lines).toHaveLength(5);
  });

  it('last line is shorter', () => {
    const { container } = render(<SkeletonText lines={3} />);
    const lines = container.querySelectorAll('.animate-pulse');
    const lastLine = lines[lines.length - 1];
    expect(lastLine).toHaveClass('w-3/4');
  });
});

describe('SkeletonCard', () => {
  it('renders card structure', () => {
    const { container } = render(<SkeletonCard />);
    const card = container.firstChild as HTMLElement;
    expect(card).toHaveClass('rounded-2xl', 'border', 'border-slate-800');
    expect(card).toHaveAttribute('aria-hidden', 'true');
  });

  it('contains circular and text skeletons', () => {
    const { container } = render(<SkeletonCard />);
    const skeletons = container.querySelectorAll('[aria-hidden="true"]');
    expect(skeletons.length).toBeGreaterThan(1);
  });

  it('applies custom className', () => {
    const { container } = render(<SkeletonCard className="custom-class" />);
    const card = container.firstChild as HTMLElement;
    expect(card).toHaveClass('custom-class');
  });
});
