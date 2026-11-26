import { useParams } from "react-router-dom";
import { useListExtracurriculars } from "../api/endpoints";
import { useStudentProfile } from "../counselor/hooks/use-student-profile";
import StudentPageHeader from "../counselor/student-page-header";
import ExtracurricularCard from "./ExtracurricularCard";
import { ErrorMessage } from "@/components/ui/error-message";

// a counselor's read-only list of one assigned student's extracurriculars
export default function CounselorExtracurriculars() {
  const { studentId = "" } = useParams();
  const { name } = useStudentProfile(studentId);
  const { data, isLoading, error, refetch } = useListExtracurriculars(
    { student_id: studentId },
    { query: { enabled: studentId !== "" } },
  );
  const list = data ?? [];

  let body;
  if (isLoading) {
    body = <p className="font-host text-sm">Loading...</p>;
  } else if (error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's extracurriculars"
        error={error}
        onRetry={refetch}
      />
    );
  } else if (list.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        {name || "This student"} hasn&apos;t added any extracurriculars yet.
      </p>
    );
  } else {
    body = (
      <ul className="flex flex-col gap-4">
        {list.map((item) => (
          <li key={item.id}>
            <ExtracurricularCard
              item={item}
              to={`/counselor/students/${studentId}/extracurriculars/${item.id}`}
            />
          </li>
        ))}
      </ul>
    );
  }

  // same layout as the counselor's view of a student's college list
  return (
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <StudentPageHeader
        title={name ? `${name}’s Extracurriculars` : "Extracurriculars"}
        backTo={`/counselor/students/${studentId}`}
      />
      <div className="flex w-full flex-col pl-[57px] max-lg:pl-0">{body}</div>
    </div>
  );
}
