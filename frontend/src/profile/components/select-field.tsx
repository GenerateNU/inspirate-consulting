import { ChevronDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { FormField, fieldClass } from "@/profile/components/form-field";

export interface SelectOption {
  value: string;
  label: string;
}

interface SelectFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[];
  placeholder: string;
  disabled?: boolean;
}

export function SelectField({ label, ...selectProps }: SelectFieldProps) {
  return (
    <FormField label={label}>
      {(id) => <Select id={id} {...selectProps} />}
    </FormField>
  );
}

interface SelectProps extends Omit<SelectFieldProps, "label"> {
  id: string;
}

const withValue = (options: SelectOption[], value: string) =>
  value && !options.some((o) => o.value === value)
    ? [{ value, label: value }, ...options]
    : options;

function Select({
  id,
  value,
  onChange,
  options,
  placeholder,
  disabled,
}: SelectProps) {
  return (
    <div className="relative">
      <select
        id={id}
        data-slot="profile-select"
        value={value}
        disabled={disabled}
        onChange={(e) => onChange(e.target.value)}
        className={cn(fieldClass, "appearance-none pr-12")}
      >
        <option value="" disabled>
          {placeholder}
        </option>
        {withValue(options, value).map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      <ChevronDown
        aria-hidden="true"
        className="pointer-events-none absolute top-1/2 right-4 size-6 -translate-y-1/2 text-charcoal"
        strokeWidth={1.5}
      />
    </div>
  );
}
