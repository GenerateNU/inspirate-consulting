import type { ReactNode } from "react";
import { Link } from "react-router-dom";

import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  type GlobalCollege,
} from "../api/models";

// single list of deadline types: the create and update application enums share these values
export const deadlineOptions = [
  { value: ApplicationType.EA, key: "ea_deadline", label: "Early Action" },
  { value: ApplicationType.ED, key: "ed_deadline", label: "Early Decision" },
  { value: ApplicationType.RD, key: "rd_deadline", label: "Regular Decision" },
] as const;

// only need date of full datetimes
export const formatDate = (dateTime: string) =>
  new Date(dateTime).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  });

export const inputClassName =
  "min-w-0 flex-1 rounded-md border border-gray-400 px-3 py-2 text-sm disabled:cursor-not-allowed disabled:bg-gray-100";

export const primaryButtonClassName =
  "rounded-md bg-black px-5 py-2 text-sm text-white hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-40";

const circleClassName =
  "flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-black text-lg leading-none";

// round icon button on the right of a card (e.g. + for add or × for delete)
export function CircleButton({
  label,
  title,
  disabled,
  onClick,
  children,
}: {
  label: string;
  title?: string;
  disabled?: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      aria-label={label}
      title={title}
      disabled={disabled}
      onClick={onClick}
      className={`${circleClassName} hover:bg-gray-100 disabled:cursor-not-allowed disabled:hover:bg-transparent`}
    >
      {children}
    </button>
  );
}

// same circle as CircleButton but not clickable (e.g. ✓ for already added)
export function CircleBadge({
  label,
  title,
  children,
}: {
  label: string;
  title?: string;
  children: ReactNode;
}) {
  return (
    <span role="img" aria-label={label} title={title} className={circleClassName}>
      {children}
    </span>
  );
}

// one labeled row on a college's page
export function Field({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <div className="flex flex-col gap-1 border-b border-gray-200 py-3 sm:flex-row sm:items-center sm:gap-4">
      <dt className="w-40 shrink-0 text-sm text-gray-500">{label}</dt>
      <dd className="flex flex-1 text-base">{children}</dd>
    </div>
  );
}

// one search result: name links to the college's page, 
// children go under the summary line, action sits on the right
export function CollegeCard({
  college,
  detailPath,
  action,
  children,
}: {
  college: GlobalCollege;
  detailPath: string;
  action?: ReactNode;
  children?: ReactNode;
}) {
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

// small inline icons

export function SearchIcon() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      className="h-5 w-5 shrink-0"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
    >
      <circle cx="10.5" cy="10.5" r="6.5" />
      <path d="M20 20l-4.5-4.5" />
    </svg>
  );
}

export function PinIcon() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      className="h-3.5 w-3.5"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path d="M12 21s-7-6.2-7-12a7 7 0 0 1 14 0c0 5.8-7 12-7 12z" />
      <circle cx="12" cy="9" r="2.5" />
    </svg>
  );
}

export function PencilIcon() {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      className="h-3 w-3"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinejoin="round"
    >
      <path d="M4 20l1-4L16 5l3 3L8 19l-4 1z" />
    </svg>
  );
}
