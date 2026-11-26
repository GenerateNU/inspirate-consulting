import type { ReactNode } from "react";
import { cn } from "@/lib/utils";
import { Skeleton } from "@/components/ui/skeleton";

const EMPTY_VALUE = "—";

interface ReadOnlyFieldProps {
  label: string;
  value: ReactNode;
}

const isEmpty = (value: ReactNode) =>
  value === null || value === undefined || value === "";

export function ReadOnlyField({ label, value }: ReadOnlyFieldProps) {
  const empty = isEmpty(value);
  return (
    <div className="flex min-w-0 flex-col gap-1">
      <dt className="font-host text-xs leading-4 text-dark-greyish-brown">
        {label}
      </dt>
      <dd
        className={cn(
          "font-host text-sm leading-5 break-words",
          empty ? "text-warm-grey" : "text-charcoal",
        )}
      >
        {empty ? EMPTY_VALUE : value}
      </dd>
    </div>
  );
}

export function ReadOnlyFieldSkeleton() {
  return (
    <div className="flex flex-col gap-1">
      <Skeleton className="h-4 w-24 bg-cream" />
      <Skeleton className="h-5 w-40 bg-cream" />
    </div>
  );
}
