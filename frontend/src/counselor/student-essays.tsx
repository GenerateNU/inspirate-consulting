import { useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, CheckCircle, ExternalLink, IdCard } from "lucide-react";
import { Link, useParams } from "react-router-dom";

import { errorMessage } from "../api/client";
import {
  getGetEssayReviewStatusQueryKey,
  getListCounselorPendingReviewsQueryKey,
  useCompleteEssayReview,
  useGetEssayGroupsFromStudent,
  useGetEssaysFromStudent,
  useListCounselorStudents,
} from "../api/endpoints";
import type { Essays } from "../api/models";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  EssaySections,
  RequestedAt,
  StatGrid,
  isReviewCompleted,
  isReviewOpen,
  reviewStats,
  useReviewStatuses,
} from "../essays/essay-sections";

const initials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .map((part) => part[0])
    .slice(0, 2)
    .join("")
    .toUpperCase();

export default function StudentEssays() {
  const { studentId = "" } = useParams();
  const queryClient = useQueryClient();

  const essaysQuery = useGetEssaysFromStudent(studentId, {
    query: { enabled: !!studentId },
  });
  const groupsQuery = useGetEssayGroupsFromStudent(studentId, {
    query: { enabled: !!studentId },
  });
  const studentsQuery = useListCounselorStudents();
  const completeReview = useCompleteEssayReview();

  const essays = essaysQuery.data?.essays ?? [];
  const groups = groupsQuery.data?.essay_groups ?? [];
  const { reviews, error: reviewsError } = useReviewStatuses(essays);
  const student = studentsQuery.data?.find((s) => s.student_id === studentId);

  const error =
    essaysQuery.error ||
    groupsQuery.error ||
    reviewsError ||
    completeReview.error;

  const stats = reviewStats(essays, reviews);

  const handleComplete = (essay: Essays) => {
    const review = reviews[essay.id];
    if (!review) return;
    completeReview.mutate(
      { id: review.transaction_id },
      {
        onSuccess: () => {
          void queryClient.invalidateQueries({
            queryKey: getGetEssayReviewStatusQueryKey(essay.id),
          });
          void queryClient.invalidateQueries({
            queryKey: getListCounselorPendingReviewsQueryKey(),
          });
        },
      },
    );
  };

  return (
    <div className="flex-1 bg-muted/40 px-4 py-8">
      <div className="mx-auto flex max-w-5xl flex-col gap-6">
        <Button variant="outline" className="w-fit" asChild>
          <Link to="/counselor-students">
            <ArrowLeft className="size-4" aria-hidden />
            Back to all students
          </Link>
        </Button>

        <Card className="flex-row flex-wrap items-center justify-between gap-4 p-6 sm:p-8">
          <div className="flex flex-wrap items-center gap-4">
            <div
              aria-hidden
              className="flex size-16 items-center justify-center rounded-md bg-muted text-2xl font-semibold"
            >
              {student ? initials(student.name) : "?"}
            </div>
            <div>
              <h1 className="text-3xl font-bold tracking-tight">
                {student?.name ?? "Student"}
              </h1>
              {student && (
                <p className="text-sm text-muted-foreground capitalize">
                  {student.year}
                </p>
              )}
            </div>
          </div>
          <Button variant="secondary" asChild>
            <Link to={`/students/${studentId}`}>
              <IdCard className="size-4" aria-hidden />
              View student profile
            </Link>
          </Button>
        </Card>

        <StatGrid
          stats={[
            { label: "Total essays", value: stats.total },
            { label: "Awaiting review", value: stats.open, highlight: true },
            { label: "Reviewed", value: stats.completed },
            { label: "Not reviewed", value: stats.notReviewed },
          ]}
        />

        {error && (
          <p
            role="alert"
            className="rounded-md bg-destructive/10 p-4 text-destructive"
          >
            {errorMessage(error)}
          </p>
        )}

        {essaysQuery.isLoading || groupsQuery.isLoading ? (
          <p>Loading...</p>
        ) : essays.length === 0 ? (
          <Card className="p-8 text-center text-muted-foreground">
            This student has no essays yet.
          </Card>
        ) : (
          <EssaySections
            essays={essays}
            groups={groups}
            renderBadges={(essay) => {
              const review = reviews[essay.id];
              return (
                <>
                  {isReviewOpen(review) && (
                    <Badge variant="destructive">Review requested</Badge>
                  )}
                  {isReviewCompleted(review) && (
                    <Badge variant="outline">Review returned</Badge>
                  )}
                  {!isReviewOpen(review) && !isReviewCompleted(review) && (
                    <Badge variant="secondary">Not reviewed</Badge>
                  )}
                  {isReviewOpen(review) && <RequestedAt review={review!} />}
                </>
              );
            }}
            renderActions={(essay) => (
              <>
                {isReviewOpen(reviews[essay.id]) && (
                  <Button
                    variant="secondary"
                    disabled={completeReview.isPending}
                    onClick={() => handleComplete(essay)}
                  >
                    <CheckCircle className="size-4" aria-hidden />
                    Mark review complete
                  </Button>
                )}
                <Button asChild>
                  <a
                    href={essay.link_to_content}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {isReviewOpen(reviews[essay.id])
                      ? "Review & leave feedback"
                      : "Open essay"}
                    <ExternalLink className="size-4" aria-hidden />
                  </a>
                </Button>
              </>
            )}
          />
        )}
      </div>
    </div>
  );
}
