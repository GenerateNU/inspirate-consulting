import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { ModuleStatus } from "./utils/module-status";

const statusLabel: Record<ModuleStatus, string> = {
  notStarted: "Not Started",
  inProgress: "In Progress",
  completed: "Completed",
};

// hi-fi pill colors
const statusClassName: Record<ModuleStatus, string> = {
  notStarted: "bg-danger/10 text-danger",
  inProgress: "bg-cream text-dark-greyish-brown",
  completed: "bg-lightest-green text-dark-green",
};

interface StatusBadgeProps {
  status: ModuleStatus;
}

// Not Started / In Progress / Completed pill next to a module's title
export default function StatusBadge({ status }: StatusBadgeProps) {
  return (
    <Badge
      className={cn(
        "rounded-full px-2 py-0.5 font-host text-[10px] font-normal",
        statusClassName[status],
      )}
    >
      {statusLabel[status]}
    </Badge>
  );
}
