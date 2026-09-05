import { Component, type ErrorInfo, type ReactNode } from "react";
import { AlertTriangle, RotateCcw } from "lucide-react";

interface ErrorBoundaryProps {
  children: ReactNode;
}

interface ErrorBoundaryState {
  error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    const errStr = error?.message || (typeof error === 'object' ? JSON.stringify(error) : String(error));
    console.error("[ErrorBoundary] Caught render error:", errStr, info.componentStack);
  }

  handleReset = () => {
    this.setState({ error: null });
  };

  render() {
    const { error } = this.state;

    if (error) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-[#f2efe9] p-6">
          <div className="max-w-lg w-full bg-white rounded-2xl border border-slate-300 shadow-sm p-8 text-center space-y-4">
            <div className="w-14 h-14 rounded-full bg-red-100 flex items-center justify-center mx-auto">
              <AlertTriangle size={28} className="text-red-500" />
            </div>
            <div>
              <h1 className="text-lg font-bold text-slate-800">Something went wrong</h1>
              <p className="text-sm text-slate-500 mt-1">
                This page ran into an unexpected error. You can try reloading it below.
              </p>
            </div>
            <pre className="text-left text-xs bg-slate-100 rounded-lg p-3 overflow-auto max-h-40 text-red-600 whitespace-pre-wrap">
              {error?.message || (typeof error === 'object' ? JSON.stringify(error) : String(error))}
            </pre>
            <div className="flex items-center justify-center gap-3 pt-2">
              <button
                onClick={this.handleReset}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-orange text-white text-sm font-medium hover:bg-orange/90 transition-colors"
              >
                <RotateCcw size={14} />
                Try Again
              </button>
              <button
                onClick={() => window.location.assign("/dashboard")}
                className="px-4 py-2 rounded-lg border border-slate-300 text-sm font-medium hover:bg-slate-100 transition-colors"
              >
                Go to Dashboard
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
