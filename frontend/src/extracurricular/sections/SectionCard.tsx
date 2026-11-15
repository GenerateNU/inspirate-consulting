import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { cn } from "cn";
import { edit } from "../ExtracurricularShared";

export function SectionCard({
  icon: Icon,
  title,
  className,
  children,
}: {
  icon: LucideIcon;
  title: string;
  className?: string;
  children: ReactNode;
}) {
  return (
    <section className={cn(edit.card, className)}>
      <div className="flex items-center gap-5">
        <CategoryIcon icon={Icon} />
        <h2 className={edit.cardTitle}>{title}</h2>
      </div>
      {children}
    </section>
  );
}

export function CategoryIcon({ icon: Icon }: { icon: LucideIcon }) {
  return (
    <span
      aria-hidden="true"
      className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-[#f8efe4] text-charcoal"
    >
      <Icon className="size-6" strokeWidth={1.75} />
    </span>
  );
}
