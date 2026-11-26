import { useState, type ReactNode } from "react";
import { Link } from "react-router-dom";
import { ChevronDown, type LucideIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";
import {
  outlineButtonClassName,
  panelClassName,
  sectionTitleClassName,
} from "./profile-styles";

interface ProfileSectionProps {
  icon: LucideIcon;
  title: string;
  // e.g. "4 of 8"; left out while loading
  count?: string;
  viewAllLabel: string;
  viewAllTo: string;
  children: ReactNode;
}

// a collapsible card on the student profile: icon, title, count, "View all" link, and chevron
export default function ProfileSection({
  icon: Icon,
  title,
  count,
  viewAllLabel,
  viewAllTo,
  children,
}: ProfileSectionProps) {
  const [open, setOpen] = useState(true);

  return (
    <Collapsible open={open} onOpenChange={setOpen} asChild>
      <Card className={panelClassName}>
        <div className="flex items-center gap-3">
          <Icon className="size-5 text-charcoal" strokeWidth={1.75} />
          <h2 className={sectionTitleClassName}>{title}</h2>
          {count && (
            <span className="font-host text-[11px] text-dark-greyish-brown">
              {count}
            </span>
          )}
          <div className="ml-auto flex items-center gap-2">
            <Button
              asChild
              variant="outline"
              className={outlineButtonClassName}
            >
              <Link to={viewAllTo}>{viewAllLabel}</Link>
            </Button>
            <CollapsibleTrigger asChild>
              <Button
                type="button"
                variant="outline"
                size="icon"
                aria-label={open ? `Collapse ${title}` : `Expand ${title}`}
                className="size-7 rounded-full border-greyish-brown bg-white text-charcoal shadow-none hover:bg-off-white"
              >
                <ChevronDown
                  className={cn(
                    "size-4 transition-transform duration-200",
                    open && "rotate-180",
                  )}
                  strokeWidth={1.75}
                />
              </Button>
            </CollapsibleTrigger>
          </div>
        </div>
        <CollapsibleContent className="pt-4">{children}</CollapsibleContent>
      </Card>
    </Collapsible>
  );
}
