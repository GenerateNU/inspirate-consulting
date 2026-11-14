import { useState } from "react";

import type { ApiError } from "../../api/client";
import { useUpdateGlobalCollege } from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import {
  cyclesDraftError,
  cyclesFromDraft,
  draftFromCycles,
  type CyclesDraft,
} from "./cycles-editor";

// fields and save for the counselor's School Details edit mode
export function useCollegeEditForm(
  college: GlobalCollege,
  // refetches the college after a save
  onSaved: () => Promise<unknown>,
) {
  const [name, setName] = useState(college.school_name);
  const [location, setLocation] = useState(college.school_location);
  const [cycles, setCycles] = useState<CyclesDraft>(() =>
    draftFromCycles(college.cycles),
  );
  const [validationError, setValidationError] = useState<string | null>(null);

  const update = useUpdateGlobalCollege();

  // 409 means the edit is a near-duplicate of another college
  const isDuplicateWarning = (update.error as ApiError | null)?.status === 409;

  // back to the saved values
  const discard = () => {
    setName(college.school_name);
    setLocation(college.school_location);
    setCycles(draftFromCycles(college.cycles));
    setValidationError(null);
    update.reset();
  };

  // saves, refetches, then calls onDone (e.g. to leave edit mode)
  const save = (onDone: () => void, confirmDuplicate = false) => {
    update.reset();
    const problem =
      !name.trim() || !location.trim()
        ? "School name and location are required."
        : cyclesDraftError(cycles);
    setValidationError(problem);
    if (problem) return;

    update.mutate(
      {
        id: college.id,
        data: {
          school_name: name.trim(),
          school_location: location.trim(),
          // the full set: cycles left out are removed
          cycles: cyclesFromDraft(cycles),
          confirm_duplicate: confirmDuplicate,
        },
      },
      {
        onSuccess: async () => {
          await onSaved();
          onDone();
        },
      },
    );
  };

  // any edit clears the last "fix this first" message
  const edited =
    <T>(set: (value: T) => void) =>
    (value: T) => {
      setValidationError(null);
      set(value);
    };

  return {
    name,
    setName: edited(setName),
    location,
    setLocation: edited(setLocation),
    cycles,
    setCycles: edited(setCycles),
    discard,
    save,
    saving: update.isPending,
    validationError,
    isDuplicateWarning,
    error: isDuplicateWarning ? null : update.error,
  };
}
