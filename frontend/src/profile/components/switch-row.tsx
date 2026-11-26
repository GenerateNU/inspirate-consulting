import { useId } from "react";
import { cn } from "@/lib/utils";

interface SwitchRowProps {
  label: string;
  description: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

export function SwitchRow({
  label,
  description,
  checked,
  onChange,
  disabled,
}: SwitchRowProps) {
  const id = useId();
  const descriptionId = `${id}-desc`;
  return (
    <div
      className={cn(
        "flex items-center justify-between gap-4 sm:col-span-2",
        disabled && "opacity-60",
      )}
    >
      <div className="flex min-w-0 flex-col gap-0.5">
        <label
          htmlFor={id}
          className="font-host text-sm leading-5 text-charcoal"
        >
          {label}
        </label>
        <p
          id={descriptionId}
          className="font-host text-xs text-dark-greyish-brown"
        >
          {description}
        </p>
      </div>
      <Switch
        id={id}
        describedBy={descriptionId}
        checked={checked}
        onChange={onChange}
        disabled={disabled}
      />
    </div>
  );
}

interface SwitchProps {
  id: string;
  describedBy: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}

function Switch({ id, describedBy, checked, onChange, disabled }: SwitchProps) {
  return (
    <button
      id={id}
      type="button"
      role="switch"
      data-slot="profile-switch"
      aria-checked={checked}
      aria-describedby={describedBy}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={cn(
        "relative inline-flex h-6 w-11 shrink-0 items-center rounded-full outline-none transition-colors focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed",
        checked ? "bg-dark-green" : "bg-greyish-brown",
      )}
    >
      <span
        aria-hidden="true"
        className={cn(
          "size-5 rounded-full bg-white shadow-sm transition-transform motion-reduce:transition-none",
          checked ? "translate-x-[22px]" : "translate-x-0.5",
        )}
      />
    </button>
  );
}
