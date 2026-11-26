import { useNavigate, useParams } from "react-router-dom";

import {
  useListCounselorStudents,
  useListGlobalColleges,
  useListPersonalCollegeApplications,
} from "../../api/endpoints";
import arrowLeftIcon from "@/assets/college/arrow-left.svg";
import ApplicationCard from "../student/application-card";
import { sortApplications } from "../utils/sort-applications";
import { ErrorMessage } from "@/components/ui/error-message";

// "Anna's College List": a counselor's read-only view of one student's list
export default function StudentCollegeList() {
  const { studentId = "" } = useParams();
  const navigate = useNavigate();

  const applicationsQuery = useListPersonalCollegeApplications(
    { student_id: studentId },
    { query: { enabled: !!studentId } },
  );
  const collegesQuery = useListGlobalColleges();
  const studentsQuery = useListCounselorStudents();

  const student = studentsQuery.data?.find((s) => s.student_id === studentId);
  const collegeById = new Map(
    (collegesQuery.data ?? []).map((college) => [college.id, college]),
  );
  const applications = sortApplications(
    applicationsQuery.data ?? [],
    "deadline",
    collegeById,
  );

  let body;
  if (applicationsQuery.isLoading || collegesQuery.isLoading) {
    body = <p className="font-host text-sm">Loading...</p>;
  } else if (applicationsQuery.error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's schools"
        error={applicationsQuery.error}
        onRetry={applicationsQuery.refetch}
      />
    );
  } else if (applications.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        {student?.name ?? "This student"} hasn&apos;t added any schools yet.
      </p>
    );
  } else {
    body = (
      <ol className="flex w-full flex-col gap-5">
        {applications.map((app) => {
          const college = collegeById.get(app.global_college_id);
          return (
            <ApplicationCard
              key={app.id}
              app={app}
              college={college}
              collegeName={college?.school_name ?? "Unknown college"}
              detailPath={`/counselor-college-list/${app.global_college_id}`}
            />
          );
        })}
      </ol>
    );
  }

  return (
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <div className="flex w-full items-start gap-9">
        <button
          type="button"
          data-slot="icon-button"
          aria-label="Back"
          onClick={() => navigate(-1)}
          className="shrink-0 rounded-md hover:bg-black/5"
        >
          <img src={arrowLeftIcon} alt="" width={32} height={32} />
        </button>
        <h1 className="min-w-0 flex-1 truncate font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          {student ? `${student.name}’s College List` : "College List"}
        </h1>
      </div>
      <div className="flex w-full flex-col pl-[57px] max-lg:pl-0">{body}</div>
    </div>
  );
}
