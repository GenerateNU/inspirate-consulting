import { CircleAlert } from "lucide-react";
import { errorMessage } from "@/api/client";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

interface ErrorAlertProps {
  error: unknown;
  attempt: number;
  label: string;
}

export function ErrorAlert({ error, attempt, label }: ErrorAlertProps) {
  const failed = error !== null && error !== undefined;
  return (
    <>
      {failed && (
        <ErrorBubble
          message={errorMessage(error)}
          attempt={attempt}
          label={label}
        />
      )}
      <p role="alert" className="sr-only">
        {failed && `${label}: ${errorMessage(error)}`}
      </p>
    </>
  );
}

interface ErrorBubbleProps {
  message: string;
  attempt: number;
  label: string;
}

function ErrorBubble({ message, attempt, label }: ErrorBubbleProps) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <button
          type="button"
          data-slot="save-error"
          aria-label={label}
          className="flex size-8 cursor-help items-center justify-center rounded-full bg-danger/10 text-danger outline-none hover:bg-danger/15 focus-visible:ring-[3px] focus-visible:ring-danger/30"
        >
          {/* Remounting the icon restarts its animation */}
          <CircleAlert
            key={attempt}
            className="save-error-shake size-5"
            strokeWidth={1.75}
          />
        </button>
      </TooltipTrigger>
      <TooltipContent side="bottom" sideOffset={6} className="max-w-xs">
        {message}
      </TooltipContent>
    </Tooltip>
  );
}
