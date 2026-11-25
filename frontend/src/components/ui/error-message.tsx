import type { ReactNode } from "react";
import { CircleAlert } from "lucide-react";
import { cn } from "cn";
import { errorMessage } from "@/api/client";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

type ErrorMessageProps = {
  /** Any thrown/API error; its text comes from errorMessage(). */
  error?: unknown;
  /** Text to show instead of the error's own message (e.g. client-side validation). */
  message?: ReactNode;
  /** "inline" sits under a form or action; "banner" replaces a section that failed to load. */
  variant?: "inline" | "banner";
  /** Inline only: lead-in such as "Couldn't save", rendered as "Couldn't save: <detail>". */
  prefix?: string;
  /** Inline only: text size, to match the surrounding form. */
  size?: "xs" | "sm";
  /** Banner only: heading such as "Couldn't load your tasks". */
  title?: string;
  /** Banner only: shows a "Try again" button. */
  onRetry?: () => unknown;
  /** Banner only: extra action, e.g. a back link. */
  action?: ReactNode;
  className?: string;
};

// Renders nothing when there is no error, so callers can pass query/mutation errors directly.
export function ErrorMessage({
  error,
  message,
  variant = "inline",
  prefix,
  size = "sm",
  title,
  onRetry,
  action,
  className,
}: ErrorMessageProps) {
  if (!message && error == null) return null;
  const detail = message || errorMessage(error);

  if (variant === "inline") {
    return (
      <p
        role="alert"
        className={cn(
          "font-host text-danger",
          size === "xs" ? "text-xs" : "text-sm",
          className,
        )}
      >
        {prefix && `${prefix}: `}
        {detail}
      </p>
    );
  }

  return (
    <Alert variant="danger" className={className}>
      <CircleAlert aria-hidden />
      {title && <AlertTitle>{title}</AlertTitle>}
      <AlertDescription>
        <p>{detail}</p>
        {(onRetry || action) && (
          <div className="mt-1 flex items-center gap-4">
            {onRetry && (
              <Button
                type="button"
                variant="link"
                onClick={() => void onRetry()}
                className="h-auto p-0 text-sm text-danger underline"
              >
                Try again
              </Button>
            )}
            {action}
          </div>
        )}
      </AlertDescription>
    </Alert>
  );
}
