import { useEffect, useRef, useState } from "react";
import type { StudentListItem } from "../api/models";
import { useAssignCounselor, getFetchStudentsQueryKey } from "../api/endpoints";
import { useQueryClient } from "@tanstack/react-query";
import { Check, Plus } from "lucide-react";

// how each student year is shown in the Grade Cohort column
const GRADE_LABELS: Record<string, string> = {
  freshman: "9th",
  sophomore: "10th",
  junior: "11th",
  senior: "12th",
};

export function ClaimStudentTable({
  input,
  students,
  counselorID,
  error,
  hasMore,
  isLoadingMore,
  onLoadMore,
}: {
  input: string;
  students: StudentListItem[];
  counselorID: string;
  error: unknown;
  hasMore: boolean;
  isLoadingMore: boolean;
  onLoadMore: () => void;
}) {
  const { mutate: assignCounselor } = useAssignCounselor();
  const queryClient = useQueryClient();
  // students claimed on this page, so the row updates right away instead of after a refetch
  const [claimedIDs, setClaimedIDs] = useState<Set<string>>(new Set());

  const isClaimed = (student: StudentListItem) =>
    student.counselor_id === counselorID || claimedIDs.has(student.id);

  const claimStudent = (student: StudentListItem) => {
    setClaimedIDs((prev) => new Set(prev).add(student.id));
    assignCounselor(
      { data: { counselor_id: counselorID, user_id: student.user_id } },
      {
        onSuccess: () =>
          queryClient.invalidateQueries({
            queryKey: getFetchStudentsQueryKey(),
          }),
        onError: () =>
          setClaimedIDs((prev) => {
            const next = new Set(prev);
            next.delete(student.id);
            return next;
          }),
      },
    );
  };

  const studentRows = students.map((student) => {
    const claimed = isClaimed(student);
    return (
      <div
        key={student.id}
        className={`grid grid-cols-12 items-center border-b border-[#c6c6cd] px-4 pt-4 pb-[17px] ${claimed ? "bg-[rgba(240,244,250,0.7)]" : ""}`}
      >
        {/* select square: filled once the student is claimed */}
        <div className="col-span-1 flex items-center">
          <span
            aria-hidden="true"
            className={`size-4 rounded-[2.5px] ${claimed ? "bg-black" : "border border-[#767676] bg-white"}`}
          />
        </div>

        <p className="col-span-4 font-inter text-base leading-6 font-semibold text-[#171c20]">
          {student.name}
        </p>

        <p className="col-span-2 font-inter text-sm leading-5 text-[#171c20]">
          {GRADE_LABELS[student.year] ?? student.year}
        </p>

        {/* target pathway column: not set up yet */}
        <div className="col-span-3" />

        <div className="col-span-2 flex justify-end">
          {claimed ? (
            <span className="inline-flex h-9 items-center gap-1.5 border border-black bg-black px-[13px] font-inter text-[11px] leading-[14px] font-semibold tracking-[0.44px] text-white">
              Student Claimed
              <Check size={12} aria-hidden="true" />
            </span>
          ) : (
            <button
              type="button"
              onClick={() => claimStudent(student)}
              className="inline-flex h-9 items-center gap-1 rounded-none border border-[#76777d] bg-white px-[9px] font-inter text-sm leading-5 font-medium tracking-[0.14px] text-[#171c20]"
            >
              <Plus size={10} aria-hidden="true" />
              Claim Student
            </button>
          )}
        </div>
      </div>
    );
  });

  const hasSearch = input.trim() !== "";

  // infinite scroll: load the next page when the end of the list scrolls into view
  const bottomRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const el = bottomRef.current;
    if (!el || !hasSearch) return;
    const observer = new IntersectionObserver(([entry]) => {
      if (entry.isIntersecting && hasMore && !isLoadingMore) onLoadMore();
    });
    observer.observe(el);
    return () => observer.disconnect();
  }, [hasSearch, hasMore, isLoadingMore, onLoadMore]);

  return (
    <div className="min-h-[326px] border border-[#c6c6cd] bg-white p-px">
      <div className="grid grid-cols-12 items-center border-b border-[#c6c6cd] bg-[#eaeef4] px-4 pt-2 pb-[9px] font-inter text-[11px] leading-[14px] font-semibold tracking-[0.55px] text-[#171c20] uppercase">
        <span className="col-span-1 tracking-[0.44px]">Select</span>
        <span className="col-span-4">Student Identifier / ID</span>
        <span className="col-span-2">Grade Cohort</span>
        {/* target pathway column: not set up yet */}
        <span className="col-span-3" />
        <span className="col-span-2 text-right">Assignment Status</span>
      </div>
      {hasSearch && studentRows}
      {hasSearch && error ? (
        <p className="p-4 font-inter text-sm text-[#45464c]">
          Could not load students.
        </p>
      ) : null}
      {hasSearch && isLoadingMore ? (
        <p className="p-4 font-inter text-sm text-[#45464c]">
          Loading more students…
        </p>
      ) : null}
      <div ref={bottomRef} />
    </div>
  );
}
