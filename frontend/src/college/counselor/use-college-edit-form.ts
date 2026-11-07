import { useState } from "react";
import { useNavigate } from "react-router-dom";

import type { ApiError } from "../../api/client";
import {
  useDeleteGlobalCollege,
  useUpdateGlobalCollege,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import { deadlineOptions, toDateInput, toISODateTime } from "../utils/shared";

// "ea_deadline" | "ed_deadline" | "rd_deadline"
type DeadlineKey = (typeof deadlineOptions)[number]["key"];

// props for CollegeEditForm, which passes them straight to this hook
export interface CollegeEditFormProps {
  // the loaded college, fields start from its saved values
  college: GlobalCollege;
  // refetches the college after a save
  onSaved: () => Promise<unknown>;
}

// fields, save, and delete for CollegeEditForm
export function useCollegeEditForm({ college, onSaved }: CollegeEditFormProps) {
  const [name, setName] = useState(college.school_name);
  const [location, setLocation] = useState(college.school_location);
  const [deadlines, setDeadlines] = useState({
    ea_deadline: toDateInput(college.ea_deadline),
    ed_deadline: toDateInput(college.ed_deadline),
    rd_deadline: toDateInput(college.rd_deadline),
  });
  const [status, setStatus] = useState("");

  const navigate = useNavigate();

  const updateMutation = useUpdateGlobalCollege();
  const deleteMutation = useDeleteGlobalCollege();

  const busy = updateMutation.isPending || deleteMutation.isPending;
  const mutationError = updateMutation.error ?? deleteMutation.error;

  // 409 means the edit is a near-duplicate of another college
  const isDuplicateWarning =
    (updateMutation.error as ApiError | null)?.status === 409;

  const setDeadline = (key: DeadlineKey, value: string) =>
    setDeadlines({ ...deadlines, [key]: value });

  // clears the last mutation's error
  const resetMutations = () => {
    setStatus("");
    updateMutation.reset();
    deleteMutation.reset();
  };

  // save changes, then refetch
  const save = (confirmDuplicate = false) => {
    resetMutations();

    updateMutation.mutate(
      {
        id: college.id,
        data: {
          school_name: name,
          school_location: location,
          ea_deadline: toISODateTime(deadlines.ea_deadline),
          ed_deadline: toISODateTime(deadlines.ed_deadline),
          rd_deadline: toISODateTime(deadlines.rd_deadline),
          confirm_duplicate: confirmDuplicate,
        },
      },
      {
        onSuccess: async () => {
          await onSaved();
          setStatus("Saved.");
        },
      },
    );
  };

  // delete the college, then go back a page since this one is gone
  const deleteCollege = () => {
    if (!window.confirm(`Delete ${college.school_name}?`)) return;

    resetMutations();

    deleteMutation.mutate(
      { id: college.id },
      { onSuccess: () => navigate(-1) },
    );
  };

  return {
    // fields
    name,
    setName,
    location,
    setLocation,
    deadlines,
    setDeadline,

    // actions
    save,
    deleteCollege,

    // status
    busy,
    saving: updateMutation.isPending,
    deleting: deleteMutation.isPending,
    status,
    isDuplicateWarning,
    mutationError,
  };
}
