import { Panel } from "@/profile/components/panel";
import { TextField } from "@/profile/components/text-field";
import {
  tomorrowDateInput,
  type ProfileErrors,
  type ProfileValues,
  type SetProfileValue,
} from "@/profile/profile-form-values";

interface CounselorFieldsProps {
  values: ProfileValues;
  errors: ProfileErrors;
  setValue: SetProfileValue;
}

export function CounselorFields({
  values,
  errors,
  setValue,
}: CounselorFieldsProps) {
  return (
    <Panel title="Counselor Details">
      <TextField
        label="Custom Link (optional)"
        type="url"
        placeholder="https://"
        value={values.customLink}
        onChange={setValue("customLink")}
        error={errors.customLink}
        autoComplete="url"
      />
      <TextField
        label="Account Expiration (blank for none)"
        type="date"
        min={tomorrowDateInput()}
        value={values.expiresAt}
        onChange={setValue("expiresAt")}
        error={errors.expiresAt}
      />
    </Panel>
  );
}
