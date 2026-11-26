import type { ReactNode } from "react";
import { ArrowLeft, SquarePlay } from "lucide-react";

import { Button } from "@/components/ui/button";

interface MediaHeaderProps {
  title: string;
  onBack: () => void;
  // shown right after the title, e.g. the student's status badge
  badge?: ReactNode;
  // the page's button on the right, e.g. Publish, Save, or Mark as Complete
  action?: ReactNode;
}

// header row: back arrow, module icon, title, and the page's button on the right
export default function MediaHeader({
  title,
  onBack,
  badge,
  action,
}: MediaHeaderProps) {
  return (
    <div className="flex w-full items-center gap-9 max-lg:gap-4">
      <Button
        type="button"
        variant="ghost"
        size="icon"
        aria-label="Back to modules"
        onClick={onBack}
        className="text-charcoal hover:bg-black/5 [&_svg:not([class*='size-'])]:size-8"
      >
        <ArrowLeft strokeWidth={1.5} />
      </Button>
      <div className="flex min-w-0 flex-1 items-center gap-5">
        <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-lightest-yellow">
          <SquarePlay className="size-6 text-charcoal" strokeWidth={1.75} />
        </div>
        <h1 className="min-w-0 truncate font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          {title}
        </h1>
        {badge}
      </div>
      {action}
    </div>
  );
}
