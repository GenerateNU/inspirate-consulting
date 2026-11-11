import { Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";

export type PageLink = {
  to: string;
  label: string;
  role: "Student" | "Counselor";
};

/** Temporary buttons to a page's sub-pages until the real UI links to them. */
export default function PageLinks({ links }: { links: PageLink[] }) {
  return (
    <div className="flex flex-wrap items-center gap-x-6 gap-y-3 border-b bg-background px-4 py-3 md:px-6">
      {(["Student", "Counselor"] as const).map((role) => {
        const group = links.filter((link) => link.role === role);
        if (group.length === 0) return null;
        return (
          <div key={role} className="flex flex-wrap items-center gap-2">
            <span className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
              {role}
            </span>
            {group.map((link) => (
              <Link
                key={link.to}
                to={link.to}
                className="inline-flex h-8 items-center gap-1 rounded-full border bg-card px-3 text-sm whitespace-nowrap outline-none transition-colors hover:border-primary/40 hover:bg-accent focus-visible:ring-[3px] focus-visible:ring-ring/50"
              >
                {link.label}
                <ChevronRight className="size-3.5 text-muted-foreground" />
              </Link>
            ))}
          </div>
        );
      })}
    </div>
  );
}
