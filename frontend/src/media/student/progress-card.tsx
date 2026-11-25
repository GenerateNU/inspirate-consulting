import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { panelClassName } from "../shared/media-styles";

interface ProgressCardProps {
  completed: number;
  total: number;
}

// "Your Progress": how many assigned modules the student has completed
export default function ProgressCard({ completed, total }: ProgressCardProps) {
  const percent = total === 0 ? 0 : Math.round((completed / total) * 100);

  return (
    <Card className={cn(panelClassName, "gap-2 px-3 py-3")}>
      <h2 className="font-jakarta text-base leading-[normal] text-charcoal">
        Your Progress
      </h2>
      <p className="font-host text-xs text-dark-greyish-brown">
        {completed} of {total} modules completed
      </p>
      <div
        role="progressbar"
        aria-label="Modules completed"
        aria-valuemin={0}
        aria-valuemax={total}
        aria-valuenow={completed}
        className="h-1 w-full overflow-hidden rounded-full bg-cream"
      >
        <div
          className="h-full rounded-full bg-mint"
          style={{ width: `${percent}%` }}
        />
      </div>
    </Card>
  );
}
