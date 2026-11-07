import type { ReactNode } from "react";
import { Link } from "react-router-dom";

import type { GlobalCollege } from "../../../api/models";
import { deadlineOptions, formatDate } from "../shared";

// one search result: name links to the college's page,
// children go under the summary line, action on the right
interface CollegeCardProps {
  college: GlobalCollege;
  detailPath: string;
  action?: ReactNode;
  children?: ReactNode;
}

export default function CollegeCard({
  college,
  detailPath,
  action,
  children,
}: CollegeCardProps) {
  const deadlines = deadlineOptions.filter(({ key }) => college[key]);

  return (
    <li className="flex items-center gap-4 border border-black px-4 py-6">
      <div className="min-w-0 flex-1">
        <h3 className="text-2xl font-bold">
          <Link to={detailPath} className="hover:underline">
            {college.school_name}
          </Link>
        </h3>
        <p className="mt-1 text-sm">
          {college.school_location}
          {deadlines.map(({ value, key }) => (
            <span key={value}>
              {" · "}
              {value}: {formatDate(college[key]!)}
            </span>
          ))}
        </p>
        {children}
      </div>
      {action}
    </li>
  );
}
