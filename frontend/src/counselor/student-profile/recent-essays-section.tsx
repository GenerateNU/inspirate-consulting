import { FileText } from "lucide-react";

import { useGetEssaysFromStudent } from "../../api/endpoints";
import type { EssayReviewStatus, Essays } from "../../api/models";
import { ErrorMessage } from "@/components/ui/error-message";
import { Skeleton } from "@/components/ui/skeleton";
import Tag, { type TagTone } from "../../college/utils/components/tag";
import {
  isReviewOpen,
  useCollegeNames,
  useReviewStatuses,
} from "../../essays/essay-sections";
import { TASK_TYPE_STYLES } from "../../task/task-type-styles";
import ProfileListRow from "./profile-list-row";
import ProfileSection from "./profile-section";

// how many essays the section previews; "View all essays" shows the rest
const PREVIEW_COUNT = 4;

// an open review request wins over the essay's own Draft/Submitted status
function essayStatus(
  essay: Essays,
  review: EssayReviewStatus | undefined,
): { label: string; tone: TagTone } {
  if (isReviewOpen(review)) return { label: "Needs review", tone: "yellow" };
  if (essay.status === "Submitted")
    return { label: "Completed", tone: "green" };
  return { label: "Drafting", tone: "blue" };
}

interface RecentEssaysSectionProps {
  studentId: string;
}

// "Recent essays": a few of the student's essays with where each one stands
export default function RecentEssaysSection({
  studentId,
}: RecentEssaysSectionProps) {
  const essaysQuery = useGetEssaysFromStudent(studentId, undefined, {
    query: { enabled: studentId !== "" },
  });
  const essays = essaysQuery.data?.essays ?? [];
  const shown = essays.slice(0, PREVIEW_COUNT);

  const { reviews } = useReviewStatuses(shown);
  const collegeNames = useCollegeNames(
    shown.flatMap((essay) => (essay.college_id ? [essay.college_id] : [])),
  );
  const style = TASK_TYPE_STYLES.essay;

  let body;
  if (essaysQuery.isLoading) {
    body = <Skeleton className="h-28 w-full rounded-md" />;
  } else if (essaysQuery.error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's essays"
        error={essaysQuery.error}
        onRetry={essaysQuery.refetch}
      />
    );
  } else if (shown.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        No essays yet.
      </p>
    );
  } else {
    body = (
      <ul className="flex flex-col gap-2">
        {shown.map((essay) => {
          const status = essayStatus(essay, reviews[essay.id]);
          const title = essay.college_id
            ? `${collegeNames[essay.college_id]}: ${essay.type}`
            : essay.type;
          return (
            <ProfileListRow
              key={essay.id}
              icon={FileText}
              accentClassName={style.accentClassName}
              iconClassName={style.iconClassName}
              title={title}
              muted={status.label === "Completed"}
              titleTag={
                <Tag tone={status.tone} className="px-2 py-0.5 text-[10px]">
                  {status.label}
                </Tag>
              }
            />
          );
        })}
      </ul>
    );
  }

  return (
    <ProfileSection
      icon={FileText}
      title="Recent essays"
      count={
        essaysQuery.data ? `${shown.length} of ${essays.length}` : undefined
      }
      viewAllLabel="View all essays"
      viewAllTo={`/counselor/students/${studentId}/essays`}
    >
      {body}
    </ProfileSection>
  );
}
