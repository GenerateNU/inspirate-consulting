import type { ReactNode } from "react";

interface PanelProps {
  title: string;
  aside?: ReactNode;
  readOnly?: boolean;
  children: ReactNode;
}

export function Panel({ title, aside, readOnly, children }: PanelProps) {
  const Grid = readOnly ? "dl" : "div";
  return (
    <section className="flex w-full flex-col gap-3 rounded-md border border-cream bg-white p-6 shadow-card">
      <h2 className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
        {title}
      </h2>
      <div className="flex flex-col gap-6 pt-2 sm:flex-row sm:items-center">
        {aside}
        <Grid className="grid min-w-0 flex-1 grid-cols-1 gap-4 sm:grid-cols-2">
          {children}
        </Grid>
      </div>
    </section>
  );
}
