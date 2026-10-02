import { useState } from 'react';
import { requestEssayReview, getEssayReviewStatus, refundEssayReview } from '../api/endpoints';
import type { EssayReviewStatus } from '../api/models';

export default function Student() {
    // Request a review
    const [studentId, setStudentId] = useState('');
    const [essayId, setEssayId] = useState('');
    const [amount, setAmount] = useState('1');
    const [requestResult, setRequestResult] = useState('');

    // Look up the review on an essay
    const [lookupEssayId, setLookupEssayId] = useState('');
    const [review, setReview] = useState<EssayReviewStatus | null>(null);
    const [lookupResult, setLookupResult] = useState('');

    const handleRequest = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        setRequestResult('');
        try {
            const response = await requestEssayReview({
                student_id: studentId,
                essay_id: essayId,
                amount: Number(amount),
            });

            if (response.status === 200) {
                setRequestResult(`Requested. Review id: ${response.data.id}`);
            } else {
                setRequestResult(`Error ${response.status}: ${JSON.stringify(response.data)}`);
            }
        } catch (err) {
            console.error(err);
            setRequestResult('Something went wrong.');
        }
    };

    const handleLookup = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLookupResult('');
        setReview(null);
        try {
            const response = await getEssayReviewStatus(lookupEssayId);

            if (response.status === 200) {
                setReview(response.data);
            } else {
                setLookupResult(`Error ${response.status}: ${JSON.stringify(response.data)}`);
            }
        } catch (err) {
            console.error(err);
            setLookupResult('Something went wrong.');
        }
    };

    const handleRefund = async () => {
        if (!review) return;
        setLookupResult('');
        try {
            const response = await refundEssayReview(review.transaction_id);

            if (response.status === 200) {
                setLookupResult(`Refunded. Balance credited back.`);
                setReview(null);
            } else {
                setLookupResult(`Error ${response.status}: ${JSON.stringify(response.data)}`);
            }
        } catch (err) {
            console.error(err);
            setLookupResult('Something went wrong.');
        }
    };

    return (
        <div className="form-container">
            <h2>Request a Review</h2>
            <form onSubmit={handleRequest}>
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
                    <button type="button" onClick={handleRefund}>Request Refund</button>
                </div>
            )}
            {lookupResult && <p>{lookupResult}</p>}
        </div>
    );
}
