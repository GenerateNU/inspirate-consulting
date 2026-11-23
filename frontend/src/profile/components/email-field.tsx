import { useEffect, useState, type ReactNode } from "react";
import { Link } from "react-router-dom";
import { Check, Copy, MessageCircle } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const iconButtonClass =
  "flex size-7 shrink-0 items-center justify-center rounded-md text-dark-greyish-brown outline-none transition-colors hover:bg-dark-green/5 hover:text-dark-green focus-visible:ring-[3px] focus-visible:ring-dark-green/30";

function IconTooltip({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>{children}</TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  );
}

function CopyEmailButton({ email }: { email: string }) {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!copied) return;
    const timer = setTimeout(() => setCopied(false), 1500);
    return () => clearTimeout(timer);
  }, [copied]);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(email);
      setCopied(true);
    } catch {}
  };

  const label = copied ? "Copied!" : "Copy email";
  return (
    <IconTooltip label={label}>
      <button
        type="button"
        data-slot="copy-email"
        aria-label={label}
        onClick={copy}
        className={iconButtonClass}
      >
        {copied ? (
          <Check className="size-4 text-dark-green" strokeWidth={1.75} />
        ) : (
          <Copy className="size-4" strokeWidth={1.75} />
        )}
      </button>
    </IconTooltip>
  );
}

export function EmailValue({
  email,
  userId,
}: {
  email: string;
  userId: string;
}) {
  return (
    <span className="flex min-w-0 items-center gap-1">
      <span className="min-w-0 break-words">{email}</span>
      <CopyEmailButton email={email} />
      <IconTooltip label="Send message">
        <Link
          to={`/chats/${userId}`}
          data-slot="message-user"
          aria-label="Send message"
          className={iconButtonClass}
        >
          <MessageCircle className="size-4" strokeWidth={1.75} />
        </Link>
      </IconTooltip>
    </span>
  );
}
