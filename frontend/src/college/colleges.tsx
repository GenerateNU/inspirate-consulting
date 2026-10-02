import { useListPersonalCollegeApplications, useListGlobalColleges } from '../api/endpoints';

export default function CollegeApplicationList() {
    const { data: applicationsResponse, error: applicationsError, isLoading: applicationsLoading } = useListPersonalCollegeApplications();
    const { data: collegesResponse } = useListGlobalColleges();

    if (applicationsLoading) {
        return (
            <div className="form-container">
                <h2>My College Applications</h2>
                <p>Loading...</p>
            </div>
        );
    }

    if (applicationsError || !applicationsResponse || applicationsResponse.status !== 200) {
        return (
            <div className="form-container">
                <h2>My College Applications</h2>
                <p>Something went wrong.</p>
            </div>
        );
    }

    const applications = applicationsResponse.data ?? [];
    const colleges = collegesResponse?.status === 200 ? collegesResponse.data ?? [] : [];

    // PersonalCollegeApplication only stores global_college_id, so have to pull full list of colleges to get the school_name and school_location
    const collegeById = new Map(colleges.map((college) => [college.id, college]));

    return (
        <div className="form-container">
            <h2>My College Applications</h2>
            {applications.length === 0 && <p>No applications yet.</p>}
            <ul>
                {applications.map((application) => {
                    const college = collegeById.get(application.global_college_id);
                    return (
                        <li key={application.id}>
                            <strong>
                                {college ? college.school_name : `College #${application.global_college_id}`}
                            </strong>
                            {college && ` — ${college.school_location}`}
                            <div>
                                Applying: {application.application_type} | Category: {application.category}
                            </div>
                        </li>
                    );
                })}
            </ul>
        </div>
    );
}