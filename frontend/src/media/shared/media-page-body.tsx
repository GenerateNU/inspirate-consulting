import type { ReactNode } from "react";

interface MediaPageBodyProps {
  children: ReactNode;
}

// everything under a module page's header. 72px = back button + gap, so the
// content lines up under the header's module icon
export default function MediaPageBody({ children }: MediaPageBodyProps) {
  return (
    <div className="flex w-full flex-col gap-4 pl-[72px] max-lg:pl-0">
      {children}
    </div>
  );
}
