import { useState } from "react";
import { Link } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import {
  useDeleteGlobalCollege,
  useListGlobalColleges,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import trashIcon from "@/assets/college/trash.svg";
import RemoveSchoolDialog from "../student/remove-school-dialog";
import SchoolCard from "../utils/components/school-card";
import Tag from "../utils/components/tag";
import { cycleLabel } from "../utils/shared";
import { editorCycleOrder } from "./cycles-editor";
import { primaryButtonClassName } from "./college-details";

// counselor's "Colleges" page: every global college with the cycles it offers
export default function CounselorCollegeList() {
  const collegesQuery = useListGlobalColleges();
  const deleteCollege = useDeleteGlobalCollege();
  const [pendingDelete, setPendingDelete] = useState<GlobalCollege | null>(
    null,
  );

  const colleges = [...(collegesQuery.data ?? [])].sort((a, b) =>
    a.school_name.localeCompare(b.school_name),
  );

  // the backend refuses (400 "still in use") while any student has the school
  // on their list or has essays for it
  const deleteBlocked =
    (deleteCollege.error as ApiError | null)?.status === 400 &&
    errorMessage(deleteCollege.error).includes("still in use");

  const confirmDelete = () => {
    if (!pendingDelete) return;
    deleteCollege.mutate(
      { id: pendingDelete.id },
      {
        onSuccess: async () => {
          await collegesQuery.refetch();
          setPendingDelete(null);
        },
      },
    );
  };

  let body;
  if (collegesQuery.isLoading) {
    body = <p className="font-host text-sm">Loading...</p>;
  } else if (collegesQuery.error) {
    body = (
      <p className="font-host text-sm">
        Error: {errorMessage(collegesQuery.error)}
      </p>
    );
  } else if (colleges.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        No colleges yet. Choose Add College to add one.
      </p>
    );
  } else {
    body = (
      <ul className="flex w-full flex-col gap-5">
        {colleges.map((college) => (
          <SchoolCard
            key={college.id}
            name={college.school_name}
            detailPath={`/counselor-college-list/${college.id}`}
            location={college.school_location}
            tags={
              (college.cycles ?? []).length === 0 ? (
                <span className="font-host text-xs text-warm-grey-dark">
                  No application cycles yet
                </span>
              ) : (
                editorCycleOrder
                  .filter((cycle) =>
                    college.cycles?.some((c) => c.cycle === cycle),
                  )
                  .map((cycle) => (
                    <Tag key={cycle} tone="blue">
                      {cycleLabel(cycle)}
                    </Tag>
                  ))
              )
            }
            action={
              <button
                type="button"
                data-slot="icon-button"
                aria-label={`Delete ${college.school_name}`}
                title="Delete college"
                onClick={() => {
                  deleteCollege.reset();
                  setPendingDelete(college);
                }}
                className="flex size-9 shrink-0 items-center justify-center rounded-lg hover:bg-off-white"
              >
                <img src={trashIcon} alt="" width={18} height={18} />
              </button>
            }
          />
        ))}
      </ul>
    );
  }

  return (
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <div className="flex w-full items-start gap-9 pl-[59px] max-lg:pl-0">
        <h1 className="min-w-0 flex-1 font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          Colleges
        </h1>
        <Link
          to="/counselor-college-list/new"
          data-slot="primary-action"
          className={primaryButtonClassName}
        >
          Add College
        </Link>
      </div>
      <div className="flex w-full flex-col pl-[57px] max-lg:pl-0">{body}</div>

      <RemoveSchoolDialog
        collegeName={pendingDelete?.school_name ?? null}
        description="It will be removed from the college list. A school that's on a student's list or has essays can't be deleted."
        blockedMessage={
          deleteBlocked
            ? `${pendingDelete?.school_name} is on a student's list or has essays, so it can't be deleted. It can be deleted once no student has it on their list.`
            : null
        }
        removing={deleteCollege.isPending}
        error={deleteCollege.error}
        onCancel={() => setPendingDelete(null)}
        onConfirm={confirmDelete}
      />
    </div>
  );
}
