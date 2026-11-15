import { cn } from "cn";

export function Toggle({
  checked,
  onChange,
  labelledBy,
}: {
  checked: boolean;
  // Leave out to render a read-only switch.
  onChange?: (checked: boolean) => void;
  labelledBy: string;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-labelledby={labelledBy}
      aria-readonly={onChange ? undefined : true}
      onClick={onChange ? () => onChange(!checked) : undefined}
      className={cn(
        "relative h-8 w-[52px] shrink-0 rounded-full transition-colors outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50",
        checked ? "bg-[#19664c]" : "bg-charcoal",
        !onChange && "cursor-default",
      )}
    >
      <span
        aria-hidden="true"
        className={cn(
          "absolute top-1 left-1 size-6 rounded-full bg-white transition-transform",
          checked && "translate-x-5",
        )}
      />
    </button>
  );
}
