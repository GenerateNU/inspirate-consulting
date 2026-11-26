import type { Student } from "@/api/models";
import { Panel } from "@/profile/components/panel";
import { ReadOnlyField } from "@/profile/components/read-only-field";

interface StudentSummaryProps {
  student: Student;
}

export function StudentSummary({ student }: StudentSummaryProps) {
  return (
    <>
      <Panel title="Academic Overview" readOnly>
        <ReadOnlyField label="High School" value={student.high_school} />
        <ReadOnlyField
          label="Graduation Year"
          value={student.graduation_year}
        />
        <ReadOnlyField label="GPA (Unweighted)" value={student.gpa} />
      </Panel>

      <Panel title="College Preferences & Specializations" readOnly>
        <ReadOnlyField label="Intended Major" value={student.intended_major} />
        <ReadOnlyField
          label="Minor / Secondary Academic Interests"
          value={student.intended_minor}
        />
        <ReadOnlyField
          label="Geographic Region Preference"
          value={student.geographic_preferences}
        />
      </Panel>
    </>
  );
}
