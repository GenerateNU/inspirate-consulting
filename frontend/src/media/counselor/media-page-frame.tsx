import type { ReactNode } from "react";

interface MediaPageFrameProps {
  children: ReactNode;
}

// page background and padding shared by the create and edit pages
export default function MediaPageFrame({ children }: MediaPageFrameProps) {
  return (
    <div className="flex min-h-full w-full flex-col gap-4 bg-off-white py-12 pr-20 pl-12 max-lg:px-6">
      {children}
    </div>
  );
}
