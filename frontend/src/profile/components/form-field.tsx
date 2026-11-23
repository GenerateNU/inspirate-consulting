import { useId, type ReactNode } from "react";

export const fieldClass =
  "w-full rounded-lg border border-greyish-brown bg-white px-4 py-3 font-host text-sm leading-5 text-dark-greyish-brown outline-none transition-[border-color,box-shadow] focus-visible:border-dark-green focus-visible:ring-[3px] focus-visible:ring-dark-green/20 disabled:cursor-not-allowed disabled:opacity-60 aria-invalid:border-danger";

interface FormFieldProps {
  label: string;
  error?: string;
  children: (id: string) => ReactNode;
}

export function FormField({ label, error, children }: FormFieldProps) {
  const id = useId();
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="font-host text-sm leading-5 text-charcoal">
        {label}
      </label>
      {children(id)}
      {error && <p className="font-host text-xs text-danger">{error}</p>}
    </div>
  );
}
