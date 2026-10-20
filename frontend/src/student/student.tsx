import { useState } from "react";
import { errorMessage } from "../api/client";
import {
  useRequestEssayReview,
  useGetEssayReviewStatus,
  useRefundEssayReview,
} from "../api/endpoints";

export default function Student() {
  // Request a review
  const [essayId, setEssayId] = useState("");
  const [amount, setAmount] = useState("1");
  const [requestResult, setRequestResult] = useState("");

  // Look up the review on an essay
  const [lookupEssayId, setLookupEssayId] = useState("");
  const [submittedEssayId, setSubmittedEssayId] = useState("");
  const [lookupResult, setLookupResult] = useState("");

  const requestReview = useRequestEssayReview();
  const lookup = useGetEssayReviewStatus(submittedEssayId, {
    query: { enabled: !!submittedEssayId },
  });
  const refund = useRefundEssayReview();
  const review = submittedEssayId ? lookup.data : undefined;

  const handleRequest = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setRequestResult("");
    requestReview.mutate(
      {
        data: {
          essay_id: essayId,
          amount: Number(amount),
        },
      },
      {
        onSuccess: (transaction) =>
          setRequestResult(`Requested. Review id: ${transaction.id}`),
      },
    );
  };

  const handleLookup = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLookupResult("");
    refund.reset();
    if (lookupEssayId === submittedEssayId) void lookup.refetch();
    else setSubmittedEssayId(lookupEssayId);
  };

  const handleRefund = () => {
    if (!review) return;
    setLookupResult("");
    refund.mutate(
      { id: review.transaction_id },
      {
        onSuccess: () => {
          setLookupResult(`Refunded. Balance credited back.`);
          setSubmittedEssayId("");
        },
      },
    );
  };

  const lookupError = submittedEssayId ? lookup.error : null;

  return (
    <div className="form-container">
      <h2>Request a Review</h2>
      <form onSubmit={handleRequest}>
        <div className="form-group">
          <label>Essay ID:</label>
          <input
            type="text"
            value={essayId}
            onChange={(e) => setEssayId(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Amount:</label>
          <input
            type="number"
            min="1"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            required
          />
        </div>
        <button type="submit">Request Review</button>
      </form>
      {requestResult && <p>{requestResult}</p>}
      {requestReview.error && <p>Error: {errorMessage(requestReview.error)}</p>}

      <h2>Review Status for an Essay</h2>
      <form onSubmit={handleLookup}>
        <div className="form-group">
          <label>Essay ID:</label>
          <input
            type="text"
            value={lookupEssayId}
            onChange={(e) => setLookupEssayId(e.target.value)}
            required
          />
        </div>
        <button type="submit">Get Status</button>
      </form>
      {review && (
        <div>
          <p>Status: {review.status}</p>
          <p>Review ID: {review.transaction_id}</p>
          <p>Requested: {review.requested_at}</p>
          {review.completed_at && <p>Completed: {review.completed_at}</p>}
          <button type="button" onClick={handleRefund}>
            Request Refund
          </button>
        </div>
      )}
      {lookupResult && <p>{lookupResult}</p>}
      {lookupError && <p>Error: {errorMessage(lookupError)}</p>}
      {refund.error && <p>Error: {errorMessage(refund.error)}</p>}
    </div>
  );
}
