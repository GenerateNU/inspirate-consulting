import type { ReactNode } from "react";

interface PageShellProps {
  children: ReactNode;
}

export function PageShell({ children }: PageShellProps) {
  return (
    <div className="flex min-h-full w-full flex-col bg-off-white py-12 pr-20 pl-12 max-lg:px-6">
      <div className="flex w-full max-w-[1090px] flex-col gap-6">
        {children}
      </div>
    </div>
  );
}
