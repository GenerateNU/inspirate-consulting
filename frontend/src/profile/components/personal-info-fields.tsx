import type { CurrentUser } from "@/api/models";
import { graduationYearOptions } from "@/lib/graduation-year";
import { Panel } from "@/profile/components/panel";
import { ProfilePictureField } from "@/profile/components/pfp-dialog";
import { SelectField } from "@/profile/components/select-field";
import { TextField } from "@/profile/components/text-field";
import type {
  ProfileErrors,
  ProfileValues,
  SetProfileValue,
} from "@/profile/profile-form-values";

interface PersonalInfoFieldsProps {
  user: CurrentUser;
  isStudent: boolean;
  values: ProfileValues;
  errors: ProfileErrors;
  setValue: SetProfileValue;
}

export function PersonalInfoFields({
  user,
  isStudent,
  values,
  errors,
  setValue,
}: PersonalInfoFieldsProps) {
  return (
    <Panel
      title="Personal Information"
      aside={<ProfilePictureField user={user} />}
    >
      <TextField
        label="Name"
        value={values.name}
        onChange={setValue("name")}
        error={errors.name}
        autoComplete="name"
      />
      {isStudent && (
        <AcademicFields values={values} errors={errors} setValue={setValue} />
      )}
    </Panel>
  );
}

type AcademicFieldsProps = Omit<PersonalInfoFieldsProps, "user" | "isStudent">;

const yearOptions = () =>
  graduationYearOptions().map(({ value, label }) => ({
    value: String(value),
    label,
  }));

function AcademicFields({ values, errors, setValue }: AcademicFieldsProps) {
  return (
    <>
      <TextField
        label="High School"
        value={values.highSchool}
        onChange={setValue("highSchool")}
      />
      <SelectField
        label="Graduation Year"
        value={values.graduationYear}
        onChange={setValue("graduationYear")}
        placeholder="Select year"
        options={yearOptions()}
      />
      <TextField
        label="GPA (Unweighted)"
        value={values.gpa}
        onChange={setValue("gpa")}
        error={errors.gpa}
        inputMode="decimal"
      />
    </>
  );
}
