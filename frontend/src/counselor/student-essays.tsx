import { useQueryClient } from "@tanstack/react-query";
import { CheckCircle, ExternalLink } from "lucide-react";
import { useParams } from "react-router-dom";

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
import { ErrorMessage } from "@/components/ui/error-message";
import StudentPageHeader from "./student-page-header";

export default function StudentEssays() {
  const { studentId = "" } = useParams();
  const queryClient = useQueryClient();

  const essaysQuery = useGetEssaysFromStudent(studentId, undefined, {
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
    // same layout as the counselor's view of a student's college list
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <StudentPageHeader
        title={student ? `${student.name}’s Essays` : "Essays"}
        backTo={`/counselor/students/${studentId}`}
      />
      <div className="flex w-full flex-col gap-6 pl-[57px] max-lg:pl-0">

        <StatGrid
          stats={[
            { label: "Total essays", value: stats.total },
            { label: "Awaiting review", value: stats.open, highlight: true },
            { label: "Reviewed", value: stats.completed },
            { label: "Not reviewed", value: stats.notReviewed },
          ]}
        />

        <ErrorMessage variant="banner" error={error} />

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
