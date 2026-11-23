import { Link } from "react-router-dom";

import type { CounselorStudent } from "../api/models";
import UserAvatar from "@/components/user-avatar";
import { gradeForGraduationYear } from "@/lib/graduation-year";
import { cn } from "@/lib/utils";
import Tag from "../college/utils/components/tag";

// avatar colors from the hi-fi, picked per student so they stay the same
const AVATAR_COLORS = ["bg-warm-grey-dark", "bg-tan", "bg-blue"];

function avatarColor(studentId: string) {
  let hash = 0;
  for (const char of studentId) hash = (hash * 31 + char.charCodeAt(0)) >>> 0;
  return AVATAR_COLORS[hash % AVATAR_COLORS.length];
}

function gradeLabel(graduationYear: number) {
  const grade = gradeForGraduationYear(graduationYear);
  return grade === null ? `Class of ${graduationYear}` : `${grade}th`;
}

const plural = (count: number, one: string, many: string) =>
  `${count} ${count === 1 ? one : many}`;

// the hi-fi "My Students" card: accent bar, avatar with a count badge, name, grade and pills
export default function StudentCard({
  student,
}: {
  student: CounselorStudent;
}) {
  const reviews = student.pending_review_count;
  const messages = student.unread_message_count;
  const notifications = reviews + messages;

  return (
    <li>
      <Link
        to={`/counselor/students/${student.student_id}`}
        className="relative flex w-full items-center gap-4 overflow-clip rounded-lg border border-cream bg-white py-6 pr-6 pl-8 shadow-card transition-colors hover:bg-off-white"
      >
        <div
          aria-hidden="true"
          className={cn(
            "absolute top-[-1px] bottom-[-1px] left-[-1px] w-1",
            reviews > 0 ? "bg-danger" : "bg-blue",
          )}
        />

        <div className="relative shrink-0">
          <UserAvatar
            name={student.name}
            className="size-12 ring-4 ring-cream"
            fallbackClassName={cn(
              "font-host text-xl font-normal text-white",
              avatarColor(student.student_id),
            )}
          />
          {notifications > 0 && (
            <span
              aria-label={plural(
                notifications,
                "notification",
                "notifications",
              )}
              className={cn(
                "absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full font-host text-[9px] text-white",
                reviews > 0 ? "bg-danger" : "bg-dark-blue",
              )}
            >
              {notifications}
            </span>
          )}
        </div>

        <div className="flex min-w-0 flex-1 flex-col gap-1">
          <div className="flex min-w-0 flex-wrap items-center gap-2">
            <p className="truncate font-host text-lg leading-[normal] text-charcoal">
              {student.name}
            </p>
            {reviews > 0 && (
              <Tag tone="red" className="py-0.5 text-[10px]">
                {plural(reviews, "review requested", "reviews requested")}
              </Tag>
            )}
            {messages > 0 && (
              <Tag tone="blue" className="py-0.5 text-[10px]">
                {plural(messages, "new message", "new messages")}
              </Tag>
            )}
          </div>
          <p className="font-host text-xs text-warm-grey-dark">
            {gradeLabel(student.graduation_year)}
          </p>
          <Tag
            tone="grey"
            className="mt-1 self-start border-none bg-cream/60 py-0.5 text-[10px]"
          >
            {plural(student.open_task_count, "open task", "open tasks")}
          </Tag>
        </div>
      </Link>
    </li>
  );
}
