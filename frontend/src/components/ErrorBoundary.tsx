import { Component, type ReactNode } from 'react';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: (error: Error, reset: () => void) => ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
  }

  reset = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError && this.state.error) {
      if (this.props.fallback) {
        return this.props.fallback(this.state.error, this.reset);
      }

      return (
        <div className="flex min-h-screen items-center justify-center bg-slate-950 p-6">
          <div className="max-w-md rounded-2xl border border-red-500/30 bg-red-500/10 p-6 text-center">
            <div className="mb-4 text-4xl">⚠️</div>
            <h1 className="mb-2 text-xl font-semibold text-red-200">發生錯誤</h1>
            <p className="mb-4 text-sm text-red-300">
              {this.state.error.message || '應用程式遇到未預期的錯誤'}
            </p>
            <button
              onClick={this.reset}
              className="rounded-lg bg-red-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-red-600"
            >
              重新載入
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
