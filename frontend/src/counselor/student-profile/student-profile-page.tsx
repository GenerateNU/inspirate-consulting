import { Link, useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, MessageCircle } from "lucide-react";

import type { ApiError } from "../../api/client";
import { Button } from "@/components/ui/button";
import { ErrorMessage } from "@/components/ui/error-message";
import { Skeleton } from "@/components/ui/skeleton";
import { useStudentProfile } from "../hooks/use-student-profile";
import CollegeListSection from "./college-list-section";
import ExtracurricularsSection from "./extracurriculars-section";
import { outlineButtonClassName } from "./profile-styles";
import RecentEssaysSection from "./recent-essays-section";
import StudentDetailsCard from "./student-details-card";
import ThisWeeksTasks from "./this-weeks-tasks";

// a counselor's view of one assigned student: details, this week's tasks,
// and previews of their college list, essays, and extracurriculars
export default function CounselorStudentProfile() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { student, name, isLoading, error, refetch } = useStudentProfile(id);

  let content;
  if (isLoading) {
    content = (
      <>
        <Skeleton className="h-[180px] w-full rounded-lg" />
        <Skeleton className="h-[320px] w-full rounded-lg" />
      </>
    );
  } else if (error || !student) {
    content =
      (error as ApiError | null)?.status === 404 || !error ? (
        <p className="font-host text-sm">Student not found.</p>
      ) : (
        <ErrorMessage
          variant="banner"
          title="Couldn't load this student"
          error={error}
          onRetry={refetch}
        />
      );
  } else {
    content = (
      <>
        <StudentDetailsCard student={student} />
        <ThisWeeksTasks studentId={id} />
        <CollegeListSection studentId={id} />
        <RecentEssaysSection studentId={id} />
        <ExtracurricularsSection studentId={id} />
      </>
    );
  }

  return (
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-12 max-lg:px-6">
      <div className="flex w-full items-center gap-6 max-lg:gap-4">
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Back to dashboard"
          onClick={() => navigate("/counselor")}
          className="text-charcoal hover:bg-black/5 [&_svg:not([class*='size-'])]:size-6"
        >
          <ArrowLeft strokeWidth={1.5} />
        </Button>
        <h1 className="min-w-0 flex-1 truncate font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          {isLoading ? "" : name}
        </h1>
        {student && (
          <Button asChild variant="outline" className={outlineButtonClassName}>
            <Link to={`/chats/${student.user_id}`}>
              <MessageCircle className="size-3.5" strokeWidth={1.75} />
              Message student
            </Link>
          </Button>
        )}
      </div>

      {/* 60px = back button + gap, so the cards line up under the name */}
      <div className="flex w-full flex-col gap-5 pl-[60px] max-lg:pl-0">
        {content}
      </div>
    </div>
  );
}
