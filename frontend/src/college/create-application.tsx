import { useState } from 'react';
import {
    useListGlobalColleges,
    useCreateGlobalCollege,
    useCreatePersonalCollegeApplication,
} from '../api/endpoints';

type CollegeMode = 'select' | 'create';

const toISODateTime = (dateStr: string): string | undefined =>
    dateStr ? `${dateStr}T00:00:00Z` : undefined;

export default function CreateCollegeApplication() {
    const [mode, setMode] = useState<CollegeMode>('select');

    // If selecting an existing college (select from list of global)
    const { data: collegesResponse } = useListGlobalColleges();
    const colleges = collegesResponse?.status === 200 ? collegesResponse.data ?? [] : [];
    const [selectedCollegeId, setSelectedCollegeId] = useState<number | null>(null);

    // If creating a new college
    const [newSchoolName, setNewSchoolName] = useState('');
    const [newSchoolLocation, setNewSchoolLocation] = useState('');
    const [newEaDeadline, setNewEaDeadline] = useState('');
    const [newEdDeadline, setNewEdDeadline] = useState('');
    const [newRdDeadline, setNewRdDeadline] = useState('');

    // Application-specific fields
    const [applicationType, setApplicationType] = useState<'EA' | 'ED' | 'RD' | ''>('');
    const [category, setCategory] = useState<'safety' | 'target' | 'reach' | ''>('');

    const [status, setStatus] = useState('');

    const { trigger: createCollege, isMutating: isCreatingCollege } = useCreateGlobalCollege();
    const { trigger: createApplication, isMutating: isCreatingApplication } = useCreatePersonalCollegeApplication();

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (!applicationType || !category) {
            setStatus('Please choose a deadline type and a category.');
            return;
        }

        try {
            let globalCollegeId = selectedCollegeId;

            if (mode === 'create') {
                const collegeResult = await createCollege({
                    school_name: newSchoolName,
                    school_location: newSchoolLocation,
                    ea_deadline: toISODateTime(newEaDeadline),
                    ed_deadline: toISODateTime(newEdDeadline),
                    rd_deadline: toISODateTime(newRdDeadline),
                });

                if (collegeResult.status !== 200) {
                    setStatus('Could not create the college.');
                    return;
                }
                globalCollegeId = collegeResult.data.id;
            }

            if (globalCollegeId === null) {
                setStatus('Please select or create a college first.');
                return;
            }

            const applicationResult = await createApplication({
                global_college_id: globalCollegeId,
                application_type: applicationType,
                category,
            });

            if (applicationResult.status === 200) {
                const submittedCollege = colleges.find((college) => college.id === globalCollegeId);
                const collegeName =
                    mode === 'create'
                        ? newSchoolName
                        : submittedCollege?.school_name ?? `College #${globalCollegeId}`;

                setStatus(`Application submitted for ${collegeName}.`);
            } else {
                setStatus('Something went wrong submitting the application.');
            }
        } catch (err) {
            console.error(err);
            setStatus('Something went wrong.');
        }
    };

    const isSubmitting = isCreatingCollege || isCreatingApplication;

    return (
        <div className="form-container">
            <h2>Apply to a College</h2>

            <div className="form-group">
                <label>
                    <input
                        type="radio"
                        checked={mode === 'select'}
                        onChange={() => setMode('select')}
                    />
                    Select an existing college
                </label>
                <label>
                    <input
                        type="radio"
                        checked={mode === 'create'}
                        onChange={() => setMode('create')}
                    />
                    Create a new college
                </label>
            </div>

            <form onSubmit={handleSubmit}>
                {mode === 'select' ? (
                    <div className="form-group">
                        <label>College:</label>
                        <select
                            value={selectedCollegeId ?? ''}
                            onChange={(e) => setSelectedCollegeId(Number(e.target.value) || null)}
                            required
                        >
                            <option value="">-- Select a college --</option>
                            {colleges.map((college) => (
                                <option key={college.id} value={college.id}>
                                    {college.school_name} ({college.school_location})
                                </option>
                            ))}
                        </select>
                        {colleges.length === 0 && <p>No colleges available yet.</p>}
                    </div>
                ) : (
                    <>
                        <div className="form-group">
                            <label>School Name:</label>
                            <input
                                type="text"
                                value={newSchoolName}
                                onChange={(e) => setNewSchoolName(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>School Location:</label>
                            <input
                                type="text"
                                value={newSchoolLocation}
                                onChange={(e) => setNewSchoolLocation(e.target.value)}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>EA Deadline:</label>
                            <input
                                type="date"
                                value={newEaDeadline}
                                onChange={(e) => setNewEaDeadline(e.target.value)}
                            />
                        </div>
                        <div className="form-group">
                            <label>ED Deadline:</label>
                            <input
                                type="date"
                                value={newEdDeadline}
                                onChange={(e) => setNewEdDeadline(e.target.value)}
                            />
                        </div>
                        <div className="form-group">
                            <label>RD Deadline:</label>
                            <input
                                type="date"
                                value={newRdDeadline}
                                onChange={(e) => setNewRdDeadline(e.target.value)}
                            />
                        </div>
                    </>
                )}

                <div className="form-group">
                    <label>Applying by:</label>
                    <select
                        value={applicationType}
                        onChange={(e) => setApplicationType(e.target.value as typeof applicationType)}
                        required
                    >
                        <option value="">-- Select a deadline type --</option>
                        <option value="EA">Early Action (EA)</option>
                        <option value="ED">Early Decision (ED)</option>
                        <option value="RD">Regular Decision (RD)</option>
                    </select>
                </div>

                <div className="form-group">
                    <label>Category:</label>
                    <select
                        value={category}
                        onChange={(e) => setCategory(e.target.value as typeof category)}
                        required
                    >
                        <option value="">-- Select a category --</option>
                        <option value="safety">Safety</option>
                        <option value="target">Target</option>
                        <option value="reach">Reach</option>
                    </select>
                </div>

                <button type="submit" disabled={isSubmitting}>
                    {isSubmitting ? 'Submitting...' : 'Submit Application'}
                </button>
            </form>
            {status && <p>{status}</p>}
        </div>
    );
}