import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

// pill colors from the hi-fi "Tags" component
export type TagTone = "grey" | "green" | "yellow" | "blue";

const toneClassName: Record<TagTone, string> = {
  grey: "border border-greyish-brown bg-white text-dark-greyish-brown",
  green: "bg-lightest-green text-dark-green",
  yellow: "bg-lightest-yellow text-dark-yellow",
  blue: "bg-lightest-blue text-dark-blue",
};

export const tagClassName =
  "inline-flex shrink-0 items-center justify-center rounded-full px-2 py-1 font-host text-xs leading-[normal] font-medium whitespace-nowrap";

interface TagProps {
  tone: TagTone;
  className?: string;
  title?: string;
  children: ReactNode;
}

// read-only pill, e.g. a card's cycle, interest, category, or essay progress
export default function Tag({ tone, className, title, children }: TagProps) {
  return (
    <span
      title={title}
      className={cn(tagClassName, toneClassName[tone], className)}
    >
      {children}
    </span>
  );
}
