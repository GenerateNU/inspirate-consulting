import { cn } from "@/lib/utils";
import { tagClassName } from "./tag";

export interface ChoiceOption<T extends string> {
  value: T;
  label: string;
}

interface ChoiceChipsProps<T extends string> {
  label: string;
  options: readonly ChoiceOption<T>[];
  value: T | null;
  onChange: (value: T | null) => void;
  // clicking the selected chip again clears it (for optional fields)
  allowDeselect?: boolean;
  disabled?: boolean;
  chipClassName?: string;
}

// a row of pills where one can be picked: dark green when selected, outlined otherwise
export default function ChoiceChips<T extends string>({
  label,
  options,
  value,
  onChange,
  allowDeselect = false,
  disabled = false,
  chipClassName,
}: ChoiceChipsProps<T>) {
  return (
    <div
      role="radiogroup"
      aria-label={label}
      className="flex flex-wrap items-center gap-2"
    >
      {options.map((option) => {
        const selected = option.value === value;
        return (
          <button
            key={option.value}
            type="button"
            role="radio"
            aria-checked={selected}
            data-slot="choice-chip"
            disabled={disabled}
            onClick={() =>
              onChange(selected && allowDeselect ? null : option.value)
            }
            className={cn(
              tagClassName,
              "transition-colors disabled:cursor-not-allowed disabled:opacity-50",
              selected
                ? "bg-dark-green text-white"
                : "border border-greyish-brown bg-white text-dark-greyish-brown hover:bg-off-white",
              chipClassName,
            )}
          >
            {option.label}
          </button>
        );
      })}
    </div>
  );
}
