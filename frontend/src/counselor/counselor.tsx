import { useState } from "react";
import { errorMessage } from "../api/client";
import {
  useCompleteEssayReview,
  useSetStudentReviewBalance,
} from "../api/endpoints";
import { Link } from "react-router-dom";

export default function Counselor() {
  // Complete a review
  const [reviewId, setReviewId] = useState("");
  const [completeResult, setCompleteResult] = useState("");

  // Assign a review balance
  const [studentId, setStudentId] = useState("");
  const [balance, setBalance] = useState("0");
  const [balanceResult, setBalanceResult] = useState("");

  const completeReview = useCompleteEssayReview();
  const assignBalance = useSetStudentReviewBalance();

  const handleComplete = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setCompleteResult("");
    completeReview.mutate(
      { id: reviewId },
      {
        onSuccess: (review) =>
          setCompleteResult(`Completed at ${review.completed_at}`),
      },
    );
  };

  const handleAssignBalance = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setBalanceResult("");
    assignBalance.mutate(
      { id: studentId, data: { review_balance: Number(balance) } },
      {
        onSuccess: (student) =>
          setBalanceResult(`Balance is now ${student.review_balance}`),
      },
    );
  };

  return (
    <div className="form-container">
      <Link to="/assign-task" className="rounded bg-black px-6 py-2 text-white">
        Assign Task →
      </Link>
      <h2>Complete a Review</h2>
      <form onSubmit={handleComplete}>
        <div className="form-group">
          <label>Review ID:</label>
          <input
            type="text"
            value={reviewId}
            onChange={(e) => setReviewId(e.target.value)}
            required
          />
        </div>
        <button type="submit">Complete Review</button>
      </form>
      {completeResult && <p>{completeResult}</p>}
      {completeReview.error && (
        <p>Error: {errorMessage(completeReview.error)}</p>
      )}

      <h2>Assign Review Balance</h2>
      <form onSubmit={handleAssignBalance}>
        <div className="form-group">
          <label>Student ID:</label>
          <input
            type="text"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Amount:</label>
          <input
            type="number"
            min="0"
            value={balance}
            onChange={(e) => setBalance(e.target.value)}
            required
          />
        </div>
        <button type="submit">Assign Balance</button>
      </form>
      {balanceResult && <p>{balanceResult}</p>}
      {assignBalance.error && <p>Error: {errorMessage(assignBalance.error)}</p>}
    </div>
  );
}
