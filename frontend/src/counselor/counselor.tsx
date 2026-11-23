import { useState } from "react";
import { Link } from "react-router-dom";
import { Plus } from "lucide-react";

import { errorMessage } from "../api/client";
import {
  useListCounselorStudents,
  useSetStudentReviewBalance,
} from "../api/endpoints";
import {
  ListCounselorStudentsSortBy as SortBy,
  ListCounselorStudentsSortOrder as SortOrder,
} from "../api/models";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import RecentActivity from "./recent-activity";
import StudentCard from "./student-card";

const sectionHeadingClassName =
  "font-jakarta text-xl leading-[normal] font-semibold text-charcoal";

function EssaysToReview({ count }: { count: number }) {
  return (
    <div className="flex flex-col gap-3 rounded-lg border border-cream bg-white p-6 shadow-card">
      <div className="flex items-center justify-between gap-2">
        <p className="font-host text-[10px] text-light-black">
          Essays to review
        </p>
        <Button
          asChild
          variant="outline"
          size="xs"
          className="rounded-md border-greyish-brown font-host text-[10px] font-normal text-light-black shadow-none hover:bg-off-white"
        >
          <Link to="/pending-reviews">View essays</Link>
        </Button>
      </div>
      <p className="font-jakarta text-2xl font-semibold text-charcoal">
        {count}
      </p>
    </div>
  );
}

// Set how many essay reviews a student can request
function ReviewBalanceForm() {
  const [studentId, setStudentId] = useState("");
  const [balance, setBalance] = useState("0");
  const [result, setResult] = useState("");
  const assignBalance = useSetStudentReviewBalance();

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setResult("");
    assignBalance.mutate(
      { id: studentId, data: { review_balance: Number(balance) } },
      {
        onSuccess: (student) =>
          setResult(`Balance is now ${student.review_balance}`),
      },
    );
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex flex-col gap-3 rounded-lg border border-cream bg-white p-6 shadow-card"
    >
      <p className="font-host text-sm text-charcoal">Assign review balance</p>
      <Input
        aria-label="Student ID"
        placeholder="Student ID"
        value={studentId}
        onChange={(e) => setStudentId(e.target.value)}
        required
        className="h-8 border-greyish-brown font-host text-xs md:text-xs"
      />
      <Input
        aria-label="Amount"
        type="number"
        min="0"
        value={balance}
        onChange={(e) => setBalance(e.target.value)}
        required
        className="h-8 border-greyish-brown font-host text-xs md:text-xs"
      />
      <Button
        type="submit"
        size="sm"
        disabled={assignBalance.isPending}
        className="self-start bg-dark-green font-host text-xs font-normal text-white hover:bg-forest"
      >
        Assign balance
      </Button>
      {result && <p className="font-host text-xs text-dark-green">{result}</p>}
      {assignBalance.error && (
        <p className="font-host text-xs text-danger">
          {errorMessage(assignBalance.error)}
        </p>
      )}
    </form>
  );
}

// The counselor dashboard: their students, essays waiting on review, and recent activity
export default function Counselor() {
  // students with requested reviews first, then by name
  const {
    data: studentsData,
    error,
    isLoading,
  } = useListCounselorStudents({
    sort_by: SortBy.name,
    sort_order: SortOrder.asc,
  });
  const students = [...(studentsData ?? [])].sort(
    (a, b) => b.pending_review_count - a.pending_review_count,
  );
  const essaysToReview = students.reduce(
    (total, student) => total + student.pending_review_count,
    0,
  );

  return (
    <div className="flex min-h-full w-full flex-col gap-8 bg-off-white p-12 max-lg:px-6">
      <h1 className="font-jakarta text-[28px] leading-[normal] font-bold text-charcoal">
        Dashboard
      </h1>

      <div className="flex gap-12 max-lg:flex-col">
        <section className="flex min-w-0 flex-1 flex-col gap-5">
          <div className="flex items-center justify-between gap-4">
            <h2 className={sectionHeadingClassName}>
              My Students{" "}
              <span className="ml-2 font-host text-base font-normal text-warm-grey-dark">
                {isLoading ? "" : students.length}
              </span>
            </h2>
            <Button
              asChild
              className="h-9 rounded-md bg-dark-green px-5 font-host text-xs font-normal text-white hover:bg-forest"
            >
              <Link to="/assign-task">
                <Plus aria-hidden />
                Assign task
              </Link>
            </Button>
          </div>

          {error && (
            <p className="font-host text-sm text-danger">
              {errorMessage(error)}
            </p>
          )}

          {!isLoading && !error && students.length === 0 && (
            <p className="font-host text-sm text-warm-grey-dark">
              You don&apos;t have any students yet.{" "}
              <Link to="/search" className="underline">
                Claim students
              </Link>
            </p>
          )}

          <ul className="flex flex-col gap-4">
            {isLoading &&
              Array.from({ length: 4 }, (_, i) => (
                <li key={i}>
                  <Skeleton className="h-[98px] w-full rounded-lg bg-cream/60" />
                </li>
              ))}
            {students.map((student) => (
              <StudentCard key={student.student_id} student={student} />
            ))}
          </ul>
        </section>

        <aside className="flex w-[285px] shrink-0 flex-col gap-8 max-lg:w-full">
          <EssaysToReview count={essaysToReview} />
          <RecentActivity />
          <ReviewBalanceForm />
        </aside>
      </div>
    </div>
  );
}
