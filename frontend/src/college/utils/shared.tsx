import type { ReactNode } from "react";
import { Link } from "react-router-dom";

import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  type GlobalCollege,
} from "../../api/models";

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

// the backend wants full datetimes, date inputs give YYYY-MM-DD
export const toISODateTime = (dateStr: string): string | undefined =>
  dateStr ? `${dateStr}T00:00:00Z` : undefined;

// date inputs need YYYY-MM-DD
export const toDateInput = (dateTime?: string | null): string =>
  dateTime ? dateTime.slice(0, 10) : "";

export const inputClassName =
  "min-w-0 flex-1 rounded-md border border-gray-400 px-3 py-2 text-sm disabled:cursor-not-allowed disabled:bg-gray-100";

export const primaryButtonClassName =
  "rounded-md bg-black px-5 py-2 text-sm text-white hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-40";

// gray text button, e.g. Edit, Unrank, and the rank arrows on My List
export const smallButtonClassName =
  "text-xs text-gray-500 hover:text-black hover:underline disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:no-underline";

const circleClassName =
  "flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-black";

// round icon button on the right of a card (e.g. Plus to add or X to delete)
interface CircleButtonProps {
  label: string;
  title?: string;
  disabled?: boolean;
  onClick: () => void;
  children: ReactNode;
}

export function CircleButton({
  label,
  title,
  disabled,
  onClick,
  children,
}: CircleButtonProps) {
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

// same circle as CircleButton but not clickable (e.g. Check for already added)
interface CircleBadgeProps {
  label: string;
  title?: string;
  children: ReactNode;
}

export function CircleBadge({ label, title, children }: CircleBadgeProps) {
  return (
    <span
      role="img"
      aria-label={label}
      title={title}
      className={circleClassName}
    >
      {children}
    </span>
  );
}

// one labeled row on a college's page
interface FieldProps {
  label: string;
  children: ReactNode;
}

export function Field({ label, children }: FieldProps) {
  return (
    <div className="flex flex-col gap-1 border-b border-gray-200 py-3 sm:flex-row sm:items-center sm:gap-4">
      <dt className="w-40 shrink-0 text-sm text-gray-500">{label}</dt>
      <dd className="flex flex-1 text-base">{children}</dd>
    </div>
  );
}

// one search result: name links to the college's page,
// children go under the summary line, action on the right
interface CollegeCardProps {
  college: GlobalCollege;
  detailPath: string;
  action?: ReactNode;
  children?: ReactNode;
}

export function CollegeCard({
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
