import { Palette } from "lucide-react";

import { useListExtracurriculars } from "../../api/endpoints";
import { ErrorMessage } from "@/components/ui/error-message";
import { Skeleton } from "@/components/ui/skeleton";
import Tag from "../../college/utils/components/tag";
import {
  categoryLabel,
  dateRangeLabel,
} from "../../extracurricular/ExtracurricularShared";
import { TASK_TYPE_STYLES } from "../../task/task-type-styles";
import ProfileListRow from "./profile-list-row";
import ProfileSection from "./profile-section";

// how many activities the section previews; "View all extracurriculars" shows the rest
const PREVIEW_COUNT = 3;

interface ExtracurricularsSectionProps {
  studentId: string;
}

// "Extracurriculars": a few of the student's activities with role, years, and category
export default function ExtracurricularsSection({
  studentId,
}: ExtracurricularsSectionProps) {
  const extracurricularsQuery = useListExtracurriculars(
    { student_id: studentId },
    { query: { enabled: studentId !== "" } },
  );
  const activities = extracurricularsQuery.data ?? [];
  const shown = activities.slice(0, PREVIEW_COUNT);
  const style = TASK_TYPE_STYLES.extracurricular;

  let body;
  if (extracurricularsQuery.isLoading) {
    body = <Skeleton className="h-28 w-full rounded-md" />;
  } else if (extracurricularsQuery.error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's extracurriculars"
        error={extracurricularsQuery.error}
        onRetry={extracurricularsQuery.refetch}
      />
    );
  } else if (shown.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        No extracurriculars yet.
      </p>
    );
  } else {
    body = (
      <ul className="flex flex-col gap-2">
        {shown.map((item) => {
          const dates = dateRangeLabel(item);
          return (
            <ProfileListRow
              key={item.id}
              icon={Palette}
              accentClassName={style.accentClassName}
              iconClassName={style.iconClassName}
              title={item.name}
              to={`/counselor/students/${studentId}/extracurriculars/${item.id}`}
              tags={
                <>
                  {item.leadership_role && (
                    <Tag tone="blue" className="px-2 py-0.5 text-[10px]">
                      {item.leadership_role}
                    </Tag>
                  )}
                  {dates && (
                    <Tag tone="green" className="px-2 py-0.5 text-[10px]">
                      {dates}
                    </Tag>
                  )}
                  <Tag tone="yellow" className="px-2 py-0.5 text-[10px]">
                    {categoryLabel(item.type)}
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
      icon={Palette}
      title="Extracurriculars"
      count={
        extracurricularsQuery.data
          ? `${shown.length} of ${activities.length}`
          : undefined
      }
      viewAllLabel="View all extracurriculars"
      viewAllTo={`/counselor/students/${studentId}/extracurriculars`}
    >
      {body}
    </ProfileSection>
  );
}
