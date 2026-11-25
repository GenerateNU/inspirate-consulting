import type { ComponentProps } from "react";
import { FormField, fieldClass } from "@/profile/components/form-field";

interface TextFieldProps extends Omit<
  ComponentProps<"input">,
  "id" | "value" | "onChange"
> {
  label: string;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function TextField({
  label,
  value,
  onChange,
  error,
  ...inputProps
}: TextFieldProps) {
  return (
    <FormField label={label} error={error}>
      {(id) => (
        <input
          {...inputProps}
          id={id}
          data-slot="profile-input"
          className={fieldClass}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          aria-invalid={!!error}
        />
      )}
    </FormField>
  );
}
