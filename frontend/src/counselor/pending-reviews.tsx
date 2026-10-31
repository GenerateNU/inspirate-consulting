import { errorMessage } from "../api/client";
import {
  useListCounselorPendingReviews,
  useCompleteEssayReview,
} from "../api/endpoints";
import type { PendingReview } from "../api/models";

export default function PendingReviews() {
  const {
    data: reviewsData,
    error,
    isLoading,
    refetch,
  } = useListCounselorPendingReviews();

  const completeReview = useCompleteEssayReview();

  if (isLoading) {
    return (
      <div className="form-container">
        <h2>Pending Reviews</h2>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="form-container">
        <h2>Pending Reviews</h2>
        <p>Error: {errorMessage(error)}</p>
      </div>
    );
  }

  // Returned oldest first, so the longest-waiting review is at the top
  const reviews: PendingReview[] = reviewsData ?? [];

  // A completed review is no longer open, so the list is refetched to drop it
  const handleComplete = (review: PendingReview) => {
    completeReview.reset();

    completeReview.mutate(
      { id: review.transaction_id },
      { onSuccess: () => refetch() },
    );
  };

  return (
    <div className="form-container">
      <h2>Pending Reviews</h2>
      {completeReview.error && (
        <p role="alert">Error: {errorMessage(completeReview.error)}</p>
      )}

      {reviews.length === 0 && <p>No reviews waiting on you.</p>}
      <ul>
        {reviews.map((review) => (
          <li key={review.transaction_id}>
            <strong>{review.student_name}</strong> · {review.essay_type} ·
            requested {new Date(review.requested_at).toLocaleDateString()}{" "}
            <a href={review.essay_link} target="_blank" rel="noreferrer">
              Open essay
            </a>{" "}
            <button
              type="button"
              disabled={completeReview.isPending}
              onClick={() => handleComplete(review)}
            >
              Mark complete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
