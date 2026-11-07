import type { ReactNode } from "react";

// one labeled row on a college's page
interface FieldProps {
  label: string;
  children: ReactNode;
}

export default function Field({ label, children }: FieldProps) {
  return (
    <div className="flex flex-col gap-1 border-b border-gray-200 py-3 sm:flex-row sm:items-center sm:gap-4">
      <dt className="w-40 shrink-0 text-sm text-gray-500">{label}</dt>
      <dd className="flex flex-1 text-base">{children}</dd>
    </div>
  );
}
