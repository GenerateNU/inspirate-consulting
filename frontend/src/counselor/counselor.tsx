import { useState } from 'react';
import { completeEssayReview, setStudentReviewBalance } from '../api/endpoints';

export default function Counselor() {
    // Complete a review
    const [reviewId, setReviewId] = useState('');
    const [completeResult, setCompleteResult] = useState('');

    // Assign a review balance
    const [studentId, setStudentId] = useState('');
    const [balance, setBalance] = useState('0');
    const [balanceResult, setBalanceResult] = useState('');

    const handleComplete = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        setCompleteResult('');
        try {
            const response = await completeEssayReview(reviewId);

            if (response.status === 200) {
                setCompleteResult(`Completed at ${response.data.completed_at}`);
            } else {
                setCompleteResult(`Error ${response.status}: ${JSON.stringify(response.data)}`);
            }
        } catch (err) {
            console.error(err);
            setCompleteResult('Something went wrong.');
        }
    };

    const handleAssignBalance = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        setBalanceResult('');
        try {
            const response = await setStudentReviewBalance(studentId, {
                review_balance: Number(balance),
            });

            if (response.status === 200) {
                setBalanceResult(`Balance is now ${response.data.review_balance}`);
            } else {
                setBalanceResult(`Error ${response.status}: ${JSON.stringify(response.data)}`);
            }
        } catch (err) {
            console.error(err);
            setBalanceResult('Something went wrong.');
        }
    };

    return (
        <div className="form-container">
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
        </div>
    );
}
