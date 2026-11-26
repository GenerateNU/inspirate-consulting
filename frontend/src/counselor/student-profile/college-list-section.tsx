import { MapPin, School } from "lucide-react";

import {
  useListGlobalColleges,
  useListPersonalCollegeApplications,
} from "../../api/endpoints";
import { PersonalCollegeApplicationApplicationStatus as ApplicationStatus } from "../../api/models";
import { ErrorMessage } from "@/components/ui/error-message";
import { Skeleton } from "@/components/ui/skeleton";
import { categoryLabel } from "../../college/student/application-labels";
import Tag, { type TagTone } from "../../college/utils/components/tag";
import { cycleLabel } from "../../college/utils/shared";
import { TASK_TYPE_STYLES } from "../../task/task-type-styles";
import ProfileListRow from "./profile-list-row";
import ProfileSection from "./profile-section";

// how many schools the section previews; "View full list" shows the rest
const PREVIEW_COUNT = 4;

const statusTag: Record<string, { label: string; tone: TagTone }> = {
  [ApplicationStatus.not_started]: { label: "Not started", tone: "grey" },
  [ApplicationStatus.in_progress]: { label: "In progress", tone: "blue" },
  [ApplicationStatus.submitted]: { label: "Submitted", tone: "green" },
};

interface CollegeListSectionProps {
  studentId: string;
}

// "College list": the student's top-ranked schools with status, cycle, and category
export default function CollegeListSection({
  studentId,
}: CollegeListSectionProps) {
  const applicationsQuery = useListPersonalCollegeApplications(
    { student_id: studentId },
    { query: { enabled: studentId !== "" } },
  );
  const collegesQuery = useListGlobalColleges();

  const collegeById = new Map(
    (collegesQuery.data ?? []).map((college) => [college.id, college]),
  );
  // the backend returns ranked schools first
  const applications = applicationsQuery.data ?? [];
  const shown = applications.slice(0, PREVIEW_COUNT);
  const style = TASK_TYPE_STYLES.college;

  let body;
  if (applicationsQuery.isLoading || collegesQuery.isLoading) {
    body = <Skeleton className="h-28 w-full rounded-md" />;
  } else if (applicationsQuery.error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's schools"
        error={applicationsQuery.error}
        onRetry={applicationsQuery.refetch}
      />
    );
  } else if (shown.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        No schools added yet.
      </p>
    );
  } else {
    body = (
      <ul className="flex flex-col gap-2">
        {shown.map((app) => {
          const college = collegeById.get(app.global_college_id);
          const status = statusTag[app.application_status];
          return (
            <ProfileListRow
              key={app.id}
              icon={School}
              accentClassName={style.accentClassName}
              iconClassName={style.iconClassName}
              title={college?.school_name ?? "Unknown college"}
              to={`/counselor-college-list/${app.global_college_id}`}
              titleTag={
                status && (
                  <Tag tone={status.tone} className="px-2 py-0.5 text-[10px]">
                    {status.label}
                  </Tag>
                )
              }
              subtitle={
                college?.school_location && (
                  <>
                    <MapPin className="size-3" aria-hidden="true" />
                    {college.school_location}
                  </>
                )
              }
              tags={
                <>
                  <Tag tone="blue" className="px-2 py-0.5 text-[10px]">
                    {cycleLabel(app.application_type)}
                  </Tag>
                  <Tag tone="grey" className="px-2 py-0.5 text-[10px]">
                    {categoryLabel(app.category)}
                  </Tag>
                </>
              }
            />
          );
        })}
      </ul>
    );
  }

  return (
    <ProfileSection
      icon={School}
      title="College list"
      count={
        applicationsQuery.data
          ? `${shown.length} of ${applications.length}`
          : undefined
      }
      viewAllLabel="View full list"
      viewAllTo={`/counselor/students/${studentId}/colleges`}
    >
      {body}
    </ProfileSection>
  );
}
