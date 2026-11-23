import type { LucideIcon } from "lucide-react";

import { CardHeader, CardTitle } from "@/components/ui/card";

interface CardIconHeaderProps {
  icon: LucideIcon;
  title: string;
}

// a module page card's header: icon tile and title, e.g. "Video" or "Description"
export default function CardIconHeader({
  icon: Icon,
  title,
}: CardIconHeaderProps) {
  return (
    <CardHeader className="flex items-center gap-4 px-5">
      <div className="flex size-9 items-center justify-center rounded-lg bg-lightest-yellow">
        <Icon className="size-5 text-charcoal" strokeWidth={1.75} />
      </div>
      <CardTitle className="font-jakarta text-lg font-medium text-charcoal">
        {title}
      </CardTitle>
    </CardHeader>
  );
}
