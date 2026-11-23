import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { SquarePlay } from "lucide-react";

import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { panelClassName } from "./media-styles";

interface ModuleCardProps {
  // where clicking the title goes
  to: string;
  title: string;
  description: string;
  // shown after the title, e.g. the student's status badge and due date
  details?: ReactNode;
  // shown on the right, e.g. the counselor's delete button
  action?: ReactNode;
  // completed modules get the hi-fi's dashed, see-through card
  muted?: boolean;
}

// one row in a Modules list: tan accent bar, module icon, title, and description
export default function ModuleCard({
  to,
  title,
  description,
  details,
  action,
  muted = false,
}: ModuleCardProps) {
  return (
    <Card
      role="listitem"
      className={cn(
        panelClassName,
        "relative min-h-[115px] flex-row items-center gap-4 overflow-clip py-6 pr-6 pl-7",
        muted &&
          "border-dashed border-greyish-brown bg-transparent shadow-none",
      )}
    >
      <div
        aria-hidden="true"
        className="absolute top-[-1px] bottom-[-1px] left-[-1px] w-1 bg-tan"
      />
      <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-lightest-yellow">
        <SquarePlay
          className="size-5 text-dark-greyish-brown"
          strokeWidth={1.75}
        />
      </div>
      {/* inline-size containment: long titles/descriptions no longer
          widen the page; they truncate to fit the card instead */}
      <div className="flex min-w-0 flex-1 flex-col gap-1 [contain:inline-size]">
        <div className="flex min-w-0 items-center gap-3">
          <Link
            to={to}
            className="truncate font-jakarta text-lg leading-[normal] font-medium text-charcoal hover:underline"
          >
            {title}
          </Link>
          {details}
        </div>
        <p className="truncate font-host text-xs text-dark-greyish-brown">
          {description}
        </p>
      </div>
      {action}
    </Card>
  );
}
