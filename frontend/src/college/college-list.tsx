import { useListGlobalColleges } from '../api/endpoints';

export default function CollegeList() {
    const { data: response, error, isLoading } = useListGlobalColleges();

    if (isLoading) {
        return (
            <div className="form-container">
                <h2>Global Colleges</h2>
                <p>Loading...</p>
            </div>
        );
    }

    if (error || !response || response.status !== 200) {
        return (
            <div className="form-container">
                <h2>Global Colleges</h2>
                <p>Something went wrong.</p>
            </div>
        );
    }

    // response is { data: GlobalCollege[] | null, status, headers }
    const colleges = response.data ?? [];

    return (
        <div className="form-container">
            <h2>Global Colleges</h2>
            {colleges.length === 0 && <p>No colleges yet.</p>}
            <ul>
                {colleges.map((college) => (
                    <li key={college.id}>
                        <strong>{college.school_name}</strong> — {college.school_location}
                        <div>
                            EA: {college.ea_deadline ?? 'N/A'} | ED: {college.ed_deadline ?? 'N/A'} | RD: {college.rd_deadline ?? 'N/A'}
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
}