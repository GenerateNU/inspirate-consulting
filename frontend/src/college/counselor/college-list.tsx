import { useState } from "react";
import { X } from "lucide-react";

import { errorMessage } from "../../api/client";
import { useDeleteGlobalCollege } from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import CollegeSearch from "../utils/college-search";
import {
  CircleButton,
  CollegeCard,
  primaryButtonClassName,
} from "../utils/shared";
import { useCollegeSearch } from "../utils/use-college-search";
import AddCollegeModal from "./add-college-modal";

export default function CounselorCollegeList() {
  const search = useCollegeSearch();

  const deleteCollege = useDeleteGlobalCollege();

  const [showAddModal, setShowAddModal] = useState(false);

  // delete a college, then refetch to show the updated list
  const handleDelete = (college: GlobalCollege) => {
    if (!window.confirm(`Delete ${college.school_name}?`)) return;

    deleteCollege.reset();

    deleteCollege.mutate(
      { id: college.id },
      { onSuccess: () => search.refetch() },
    );
  };

  return (
    <>
      <CollegeSearch
        title="Global Colleges"
        search={search}
        renderCollege={(college) => (
          <CollegeCard
            key={college.id}
            college={college}
            detailPath={`/counselor-college-list/${college.id}`}
            action={
              <CircleButton
                label={`Delete ${college.school_name}`}
                title="Delete college"
                disabled={deleteCollege.isPending}
                onClick={() => handleDelete(college)}
              >
                <X aria-hidden="true" className="h-4 w-4" />
              </CircleButton>
            }
          />
        )}
        notice={
          deleteCollege.error && (
            <p role="alert" className="mt-3 text-sm text-red-600">
              Error: {errorMessage(deleteCollege.error)}
            </p>
          )
        }
        footer={
          <>
            <p className="text-xs">Missing a school?</p>
            <button
              type="button"
              onClick={() => setShowAddModal(true)}
              className={primaryButtonClassName}
            >
              Add a college
            </button>
          </>
        }
      />

      {showAddModal && (
        <AddCollegeModal
          onClose={() => setShowAddModal(false)}
          onCreated={search.refetch}
        />
      )}
    </>
  );
}
