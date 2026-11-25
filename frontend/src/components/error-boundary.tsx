import { Component, type ErrorInfo, type ReactNode } from "react";
import { ErrorMessage } from "@/components/ui/error-message";

type Props = { children: ReactNode; className?: string };
type State = { error: Error | null };

// Catches render crashes so a bug shows a message instead of a blank screen.
export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(error, info.componentStack);
  }

  render() {
    if (!this.state.error) return this.props.children;
    return (
      <div className={this.props.className ?? "p-12 max-lg:px-6"}>
        <ErrorMessage
          variant="banner"
          title="Something went wrong on this page"
          message="An unexpected error occurred. Reloading usually fixes it."
          onRetry={() => window.location.reload()}
        />
      </div>
    );
  }
}
