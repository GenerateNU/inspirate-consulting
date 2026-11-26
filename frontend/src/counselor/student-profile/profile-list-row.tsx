import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import type { LucideIcon } from "lucide-react";

import { cn } from "@/lib/utils";

interface ProfileListRowProps {
  icon: LucideIcon;
  // the colored edge, e.g. "bg-mint"
  accentClassName: string;
  // the square behind the icon, e.g. "bg-lightest-green text-dark-green"
  iconClassName: string;
  title: string;
  // where the title links; plain text when left out
  to?: string;
  // pill right after the title, e.g. a status
  titleTag?: ReactNode;
  // small line under the title, e.g. a location
  subtitle?: ReactNode;
  // pills on the right
  tags?: ReactNode;
  // done items get the hi-fi's dashed card
  muted?: boolean;
}

// one compact row in a student profile section (a college, essay, or extracurricular)
export default function ProfileListRow({
  icon: Icon,
  accentClassName,
  iconClassName,
  title,
  to,
  titleTag,
  subtitle,
  tags,
  muted = false,
}: ProfileListRowProps) {
  const titleClassName = cn(
    "truncate font-host text-sm leading-[normal]",
    muted ? "text-dark-greyish-brown" : "text-charcoal",
  );

  return (
    <li
      className={cn(
        "relative flex items-center gap-3 overflow-clip rounded-md border py-3 pr-4 pl-5",
        muted ? "border-dashed border-blue bg-off-white" : "border-cream bg-white",
      )}
    >
      <div
        aria-hidden="true"
        className={cn("absolute -top-px -bottom-px -left-px w-1", accentClassName)}
      />
      <span
        aria-hidden="true"
        className={cn(
          "flex size-7 shrink-0 items-center justify-center rounded-md",
          iconClassName,
        )}
      >
        <Icon className="size-4" strokeWidth={1.75} />
      </span>

      <div className="flex min-w-0 flex-1 flex-col gap-1 [contain:inline-size]">
        <div className="flex min-w-0 items-center gap-2">
          {to ? (
            <Link to={to} className={cn(titleClassName, "hover:underline")}>
              {title}
            </Link>
          ) : (
            <p className={titleClassName}>{title}</p>
          )}
          {titleTag}
        </div>
        {subtitle && (
          <div className="flex items-center gap-1 font-host text-[11px] text-warm-grey">
            {subtitle}
          </div>
        )}
      </div>

      {tags && <div className="flex shrink-0 items-center gap-2">{tags}</div>}
    </li>
  );
}
