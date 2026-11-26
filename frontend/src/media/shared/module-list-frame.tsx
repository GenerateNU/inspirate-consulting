import type { ReactNode } from "react";

interface ModuleListFrameProps {
  // shown right of the "Modules" heading, e.g. the search box and Create Module
  actions: ReactNode;
  children: ReactNode;
}

// page layout for a Modules list: background, "Modules" heading row, and the content below
export default function ModuleListFrame({
  actions,
  children,
}: ModuleListFrameProps) {
  return (
    <div className="flex min-h-full w-full flex-col gap-6 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <div className="flex w-full flex-wrap items-center gap-4 pl-[59px] max-lg:pl-0">
        <h1 className="min-w-0 flex-1 font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          Modules
        </h1>
        {actions}
      </div>
      <div className="flex w-full flex-col gap-6 pl-[57px] max-lg:pl-0">
        {children}
      </div>
    </div>
  );
}
