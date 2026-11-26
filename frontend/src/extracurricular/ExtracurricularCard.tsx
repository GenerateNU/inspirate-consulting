import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { Palette, Trash2 } from "lucide-react";
import type { Extracurricular } from "../api/models";
import { Button } from "@/components/ui/button";
import { categoryLabel, dateRangeLabel } from "./ExtracurricularShared";

interface ExtracurricularCardProps {
  item: Extracurricular;
  // where clicking the card goes
  to: string;
  // shows the trash button; left out on the counselor's read-only list
  onDelete?: () => void;
  deleteDisabled?: boolean;
}

// one activity on an extracurriculars list: tan edge, icon, name, and pills
export default function ExtracurricularCard({
  item,
  to,
  onDelete,
  deleteDisabled = false,
}: ExtracurricularCardProps) {
  const dates = dateRangeLabel(item);

  return (
    <div className="relative flex min-h-[136px] items-center gap-4 overflow-hidden rounded-lg border border-cream bg-white py-6 pr-6 pl-8 transition-shadow focus-within:ring-2 focus-within:ring-ring/50 hover:shadow-sm">
      <span
        aria-hidden="true"
        className="absolute inset-y-0 left-0 w-1.5 bg-tan"
      />
      <span
        aria-hidden="true"
        className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-[#f8efe4] text-charcoal"
      >
        <Palette className="size-6" strokeWidth={1.75} />
      </span>

      <div className="flex min-w-0 flex-1 flex-col gap-2">
        {/* The link's ::after covers the card so the whole card is clickable. */}
        <Link
          to={to}
          className="truncate text-xl font-medium text-charcoal outline-none after:absolute after:inset-0"
        >
          {item.name}
        </Link>
        <div className="flex flex-wrap items-center gap-2">
          {item.leadership_role && (
            <Pill className="bg-[#e7edef] text-[#415763]">
              {item.leadership_role}
            </Pill>
          )}
          {dates && (
            <Pill className="bg-[#e6f5ef] text-[#19664c]">{dates}</Pill>
          )}
          <Pill className="bg-[#f8efe4] text-[#786246]">
            {categoryLabel(item.type)}
          </Pill>
        </div>
      </div>

      {onDelete && (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="relative z-10 rounded-lg text-charcoal [&_svg:not([class*='size-'])]:size-[18px]"
          onClick={onDelete}
          disabled={deleteDisabled}
          aria-label={`Delete ${item.name}`}
        >
          <Trash2 aria-hidden="true" />
        </Button>
      )}
    </div>
  );
}

interface PillProps {
  className: string;
  children: ReactNode;
}

function Pill({ className, children }: PillProps) {
  return (
    <span
      className={`rounded-full px-2 py-1 text-xs leading-none whitespace-nowrap ${className}`}
    >
      {children}
    </span>
  );
}
