import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { ErrorBoundary } from '../../../src/components/ErrorBoundary';

// Component that throws an error
function ThrowError({ shouldThrow }: { shouldThrow: boolean }) {
  if (shouldThrow) {
    throw new Error('Test error message');
  }
  return <div>No error</div>;
}

describe('ErrorBoundary', () => {
  // Suppress console.error for these tests
  const originalError = console.error;
  beforeAll(() => {
    console.error = vi.fn();
  });
  afterAll(() => {
    console.error = originalError;
  });

  it('renders children when no error', () => {
    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ErrorBoundary>,
    );
    expect(screen.getByText('No error')).toBeInTheDocument();
  });

  it('renders default error UI when error occurs', () => {
    render(
      <ErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>,
    );
    expect(screen.getByText('發生錯誤')).toBeInTheDocument();
    expect(screen.getByText('Test error message')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '重新載入' })).toBeInTheDocument();
  });

  it('renders custom fallback when provided', () => {
    const customFallback = (error: Error) => (
      <div>Custom error: {error.message}</div>
    );

    render(
      <ErrorBoundary fallback={customFallback}>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>,
    );
    expect(screen.getByText('Custom error: Test error message')).toBeInTheDocument();
  });

  it('resets error state when reset is called', () => {
    const { rerender } = render(
      <ErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ErrorBoundary>,
    );

    expect(screen.getByText('發生錯誤')).toBeInTheDocument();

    const resetButton = screen.getByRole('button', { name: '重新載入' });
    resetButton.click();

    // After reset, re-render with no error
    rerender(
      <ErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ErrorBoundary>,
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });
});
