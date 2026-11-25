import { Panel } from "@/profile/components/panel";
import { SelectField } from "@/profile/components/select-field";
import { TextField } from "@/profile/components/text-field";
import type {
  ProfileValues,
  SetProfileValue,
} from "@/profile/profile-form-values";

const REGION_OPTIONS = [
  "Northeast",
  "Southeast",
  "Midwest",
  "Southwest",
  "West",
  "No Preference",
].map((region) => ({ value: region, label: region }));

interface CollegePreferenceFieldsProps {
  values: ProfileValues;
  setValue: SetProfileValue;
}

export function CollegePreferenceFields({
  values,
  setValue,
}: CollegePreferenceFieldsProps) {
  return (
    <Panel title="College Preferences & Specializations">
      <TextField
        label="Intended Major"
        value={values.intendedMajor}
        onChange={setValue("intendedMajor")}
      />
      <TextField
        label="Minor / Secondary Academic Interests"
        value={values.intendedMinor}
        onChange={setValue("intendedMinor")}
      />
      <SelectField
        label="Geographic Region Preference"
        value={values.region}
        onChange={setValue("region")}
        placeholder="Select region"
        options={REGION_OPTIONS}
      />
    </Panel>
  );
}
