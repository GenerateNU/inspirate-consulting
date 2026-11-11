import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Check, Plus } from "lucide-react";

import { useListPersonalCollegeApplications } from "../../api/endpoints";
import CollegeSearch from "../utils/college-search";
import CircleBadge from "../utils/components/circle-badge";
import CircleButton from "../utils/components/circle-button";
import CollegeCard from "../utils/components/college-card";
import { inputClassName, primaryButtonClassName } from "../utils/shared";
import { useCollegeSearch } from "../utils/use-college-search";
import PageLinks, { type PageLink } from "@/components/page-links";

// Temporary links to this page's sub-pages (see the Figma sitemap).
const pageLinks: PageLink[] = [
  { to: "/college-applications", label: "My College List", role: "Student" },
  {
    to: "/counselor-college-list",
    label: "Counselor College List",
    role: "Counselor",
  },
];

// list of colleges
// sends the student to the applications page to apply
export default function StudentCollegeList() {
  // request form fields, not submitted anywhere until the backend supports requests
  const [requestName, setRequestName] = useState("");
  const [requestLocation, setRequestLocation] = useState("");

  const search = useCollegeSearch();
  const navigate = useNavigate();
  const { data: applications } = useListPersonalCollegeApplications();

  // colleges the student already applied to show as added
  const appliedIds = new Set(
    (applications ?? []).map((app) => app.global_college_id),
  );

  return (
    <>
      <PageLinks links={pageLinks} />
      <CollegeSearch
        title="Student Colleges"
        search={search}
        renderCollege={(college) => (
          <CollegeCard
            key={college.id}
            college={college}
            detailPath={`/student-college-list/${college.id}`}
            action={
              appliedIds.has(college.id) ? (
                <CircleBadge
                  label={`${college.school_name} added`}
                  title="Already on your list"
                >
                  <Check aria-hidden="true" className="h-4 w-4" />
                </CircleBadge>
              ) : (
                <CircleButton
                  label={`Apply to ${college.school_name}`}
                  title="Create an application"
                  onClick={() =>
                    navigate(`/college-applications?college=${college.id}`)
                  }
                >
                  <Plus aria-hidden="true" className="h-4 w-4" />
                </CircleButton>
              )
            }
          />
        )}
        footer={
          <>
            <p className="text-xs">
              Can&apos;t find your school? Requesting new colleges is coming
              soon.
            </p>
            {/* TODO: submit to a college request endpoint once the backend has one */}
            <form
              onSubmit={(e) => e.preventDefault()}
              className="flex w-full max-w-xl flex-wrap items-center justify-center gap-2"
            >
              <input
                type="text"
                value={requestName}
                onChange={(e) => setRequestName(e.target.value)}
                placeholder="School name"
                aria-label="School name"
                disabled
                className={inputClassName}
              />
              <input
                type="text"
                value={requestLocation}
                onChange={(e) => setRequestLocation(e.target.value)}
                placeholder="School location"
                aria-label="School location"
                disabled
                className={inputClassName}
              />
              <button type="submit" disabled className={primaryButtonClassName}>
                Request College
              </button>
            </form>
          </>
        }
      />
    </>
  );
}
