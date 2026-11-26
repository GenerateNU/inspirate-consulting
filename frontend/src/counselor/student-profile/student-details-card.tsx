import type { Student } from "../../api/models";
import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { panelClassName, sectionTitleClassName } from "./profile-styles";

interface DetailProps {
  label: string;
  value: string | number | null;
}

function Detail({ label, value }: DetailProps) {
  return (
    <div className="flex min-w-0 flex-col gap-1">
      <dt className="font-host text-[10px] text-dark-greyish-brown">{label}</dt>
      <dd
        className={cn(
          "truncate font-host text-xs",
          value == null || value === "" ? "text-warm-grey" : "text-charcoal",
        )}
      >
        {value == null || value === "" ? "Not set" : value}
      </dd>
    </div>
  );
}

interface StudentDetailsCardProps {
  student: Student;
}

// "Student details": academic overview and college preferences
export default function StudentDetailsCard({
  student,
}: StudentDetailsCardProps) {
  return (
    <Card className={cn(panelClassName, "gap-4")}>
      <h2 className={sectionTitleClassName}>Student details</h2>

      <section className="flex flex-col gap-3">
        <h3 className="font-host text-sm text-charcoal">Academic overview</h3>
        <dl className="grid grid-cols-3 gap-4 max-md:grid-cols-1">
          <Detail label="High school" value={student.high_school} />
          <Detail label="Graduation year" value={student.graduation_year} />
          <Detail label="GPA" value={`${student.gpa.toFixed(2)}/4`} />
        </dl>
      </section>

      <hr className="border-cream" />

      {/* the hi-fi's campus setting preference isn't stored yet, so it's left out */}
      <section className="flex flex-col gap-3">
        <h3 className="font-host text-sm text-charcoal">
          College preferences &amp; specializations
        </h3>
        <dl className="grid grid-cols-3 gap-4 max-md:grid-cols-1">
          <Detail label="Intended major" value={student.intended_major} />
          <Detail
            label="Minor / secondary academic interests"
            value={student.intended_minor}
          />
          <Detail
            label="Geographic region preference"
            value={student.geographic_preferences}
          />
        </dl>
      </section>
    </Card>
  );
}
