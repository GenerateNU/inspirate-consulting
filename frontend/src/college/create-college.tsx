import { useState } from 'react';
import { useCreateGlobalCollege } from '../api/endpoints';

export default function CreateCollege() {
    const [schoolName, setSchoolName] = useState('');
    const [schoolLocation, setSchoolLocation] = useState('');
    const [eaDeadline, setEaDeadline] = useState('');
    const [edDeadline, setEdDeadline] = useState('');
    const [rdDeadline, setRdDeadline] = useState('');
    const [status, setStatus] = useState('');

    const { trigger, isMutating, error } = useCreateGlobalCollege();

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const toISODateTime = (dateStr: string) => dateStr ? `${dateStr}T00:00:00Z` : undefined;

            const response = await trigger({
                school_name: schoolName,
                school_location: schoolLocation,
                ea_deadline: toISODateTime(eaDeadline),
                ed_deadline: toISODateTime(edDeadline),
                rd_deadline: toISODateTime(rdDeadline),
            });

            // result is { GlobalCollege (or error), status, headers }
            if (response.status === 200) {
                setStatus(`Created: ${response.data.school_name}`);
            } else {
                setStatus('Something went wrong.');
            }
        } catch (err) {
            console.error(err);
            setStatus('Something went wrong.');
        }
    };

    return (
        <div className="form-container">
            <h2>Create a College</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>School Name:</label>
                    <input
                        type="text"
                        value={schoolName}
                        onChange={(e) => setSchoolName(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label>School Location:</label>
                    <input
                        type="text"
                        value={schoolLocation}
                        onChange={(e) => setSchoolLocation(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label>EA Deadline:</label>
                    <input
                        type="date"
                        value={eaDeadline}
                        onChange={(e) => setEaDeadline(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <label>ED Deadline:</label>
                    <input
                        type="date"
                        value={edDeadline}
                        onChange={(e) => setEdDeadline(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <label>RD Deadline:</label>
                    <input
                        type="date"
                        value={rdDeadline}
                        onChange={(e) => setRdDeadline(e.target.value)}
                    />
                </div>
                <button type="submit" disabled={isMutating}>
                    {isMutating ? 'Creating...' : 'Create College'}
                </button>
            </form>
            {status && <p>{status}</p>}
            {error && <p>Error: {String(error)}</p>}
        </div>
    );
}