import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ArrowLeft } from "lucide-react";

import { errorMessage, type ApiError } from "../../api/client";
import { useGetGlobalCollege } from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";

// loads the college whose id is in the page's URL (/...-college-list/:id)
export function useCollegeFromRoute() {
  const id = Number(useParams().id);
  return useGetGlobalCollege(id);
}

interface CollegeDetailsLayoutProps {
  college: GlobalCollege | undefined;
  isLoading: boolean;
  error: unknown;
  children: ReactNode;
}

// frame for a college's page: back button, then the college's name and children,
// or a loading/error message until the college has loaded
export default function CollegeDetailsLayout({
  college,
  isLoading,
  error,
  children,
}: CollegeDetailsLayoutProps) {
  const navigate = useNavigate();

  let body;
  if (isLoading) {
    body = <p className="text-sm">Loading...</p>;
  } else if (error || !college) {
    body = (
      <p className="text-sm">
        {(error as ApiError | null)?.status === 404
          ? "College not found."
          : `Error: ${errorMessage(error)}`}
      </p>
    );
  } else {
    body = (
      <>
        <h3 className="mb-4 text-3xl font-bold">{college.school_name}</h3>
        {children}
      </>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-10">
      <button
        type="button"
        onClick={() => navigate(-1)}
        className="mb-4 inline-flex items-center gap-1 border-0 bg-transparent text-sm hover:underline"
      >
        <ArrowLeft aria-hidden="true" className="h-4 w-4" />
        Back
      </button>
      <div className="border border-black p-6">{body}</div>
    </div>
  );
}
