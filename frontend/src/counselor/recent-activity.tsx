import { Link } from "react-router-dom";

import { errorMessage } from "../api/client";
import { useGetActivities, useListCounselorStudents } from "../api/endpoints";
import {
  ActivityAction,
  ActivityEssayStatus,
  CurrentUserRole,
  type Activity,
  type CounselorStudent,
} from "../api/models";
import { useAuthContext } from "../contexts/use-auth-context";
import { formatTime } from "../lib/format-relative-time";
import { Skeleton } from "@/components/ui/skeleton";
import UserAvatar from "@/components/user-avatar";

function activityCopy(activity: Activity, name: string): string {
  switch (activity.action) {
    case ActivityAction.added_school:
      return `${name} added a school`;
    case ActivityAction.removed_school:
      return `${name} removed a school`;
    case ActivityAction.added_extracurricular:
      return `${name} added an extracurricular`;
    case ActivityAction.removed_extracurricular:
      return `${name} removed an extracurricular`;
    case ActivityAction.task_past_due:
      return `${name} has a task past due`;
    case ActivityAction.essay_status_changed:
      switch (activity.essay_status) {
        case ActivityEssayStatus.Review:
          return `${name} requested a review`;
        case ActivityEssayStatus.Submitted:
          return `${name} submitted an essay`;
        case ActivityEssayStatus.Draft:
          return `${name} moved an essay to draft`;
        case ActivityEssayStatus.Archived:
          return `${name} archived an essay`;
        default:
          return `${name} updated an essay`;
      }
    default:
      return `${name} had activity`;
  }
}

function activityLink(activity: Activity): string {
  return activity.action === ActivityAction.essay_status_changed
    ? `/counselor/students/${activity.student_id}/essays`
    : `/counselor/students/${activity.student_id}`;
}

export default function RecentActivity() {
  const { role } = useAuthContext();
  const isCounselor = role === CurrentUserRole.counselor;

  const {
    data: activities,
    error: activitiesError,
    isLoading: activitiesLoading,
  } = useGetActivities(undefined, { query: { enabled: isCounselor } });

  const { data: studentsData, error: studentsError } = useListCounselorStudents(
    undefined,
    { query: { enabled: isCounselor } },
  );

  // Only counselors can see this panel - and the backend 403s anyone else,
  // so skip fetching and rendering entirely otherwise.
  if (!isCounselor) return null;

  const nameByStudentId = new Map<string, string>(
    (studentsData ?? []).map((student: CounselorStudent) => [
      student.student_id,
      student.name,
    ]),
  );

  const error = activitiesError ?? studentsError;
  const items = activities ?? [];

  return (
    <section className="flex flex-col gap-4">
      <h2 className="font-jakarta text-xl leading-[normal] font-semibold text-charcoal">
        Recent Activity
      </h2>
      <div className="max-h-[480px] overflow-y-auto rounded-lg border border-cream bg-white px-4 py-2 shadow-card">
        {activitiesLoading &&
          Array.from({ length: 4 }, (_, i) => (
            <div key={i} className="flex items-center gap-3 py-3">
              <Skeleton className="size-9 rounded-full bg-cream/60" />
              <div className="flex flex-1 flex-col gap-2">
                <Skeleton className="h-4 w-3/4 bg-cream/60" />
                <Skeleton className="h-3 w-1/3 bg-cream/60" />
              </div>
            </div>
          ))}
        {error && (
          <p className="py-3 font-host text-sm text-danger">
            {errorMessage(error)}
          </p>
        )}
        {!activitiesLoading && !error && items.length === 0 && (
          <p className="py-3 font-host text-sm text-warm-grey-dark">
            No recent activity yet.
          </p>
        )}
        {items.map((activity) => {
          const name =
            nameByStudentId.get(activity.student_id) ?? "Unknown student";
          return (
            <div
              key={activity.id}
              className="flex items-center gap-3 border-b border-cream py-3 last:border-0"
            >
              <UserAvatar
                name={name}
                className="size-9"
                fallbackClassName="bg-blue font-host text-white"
              />
              <div className="min-w-0 flex-1">
                <p className="truncate font-host text-xs text-charcoal">
                  {activityCopy(activity, name)}
                </p>
                <p className="font-host text-[10px] text-warm-grey-dark">
                  {formatTime(new Date(activity.created_at))} ·{" "}
                  <Link to={activityLink(activity)} className="underline">
                    View
                  </Link>
                </p>
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
