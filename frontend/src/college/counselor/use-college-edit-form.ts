import { useState } from "react";

import type { ApiError } from "../../api/client";
import {
  useCreateGlobalCollege,
  useUpdateGlobalCollege,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import {
  cyclesDraftError,
  cyclesFromDraft,
  draftFromCycles,
  type CyclesDraft,
} from "./cycles-editor";

// fields and save for the counselor's add/edit school form.
// college is null when adding a new one
export function useCollegeEditForm(
  college: GlobalCollege | null,
  // runs after a successful save with the saved college (e.g. refetch or navigate)
  onSaved: (saved: GlobalCollege) => Promise<unknown> | void,
) {
  const [name, setName] = useState(college?.school_name ?? "");
  const [location, setLocation] = useState(college?.school_location ?? "");
  const [cycles, setCycles] = useState<CyclesDraft>(() =>
    draftFromCycles(college?.cycles),
  );
  const [validationError, setValidationError] = useState<string | null>(null);

  const create = useCreateGlobalCollege();
  const update = useUpdateGlobalCollege();
  const mutation = college ? update : create;

  // 409 means the school is a near-duplicate of another college
  const isDuplicateWarning =
    (mutation.error as ApiError | null)?.status === 409;

  // saves, then calls onSaved with the result
  const save = (confirmDuplicate = false) => {
    mutation.reset();
    const problem =
      !name.trim() || !location.trim()
        ? "School name and location are required."
        : cyclesDraftError(cycles);
    setValidationError(problem);
    if (problem) return;

    const data = {
      school_name: name.trim(),
      school_location: location.trim(),
      // the full set: on edit, cycles left out are removed
      cycles: cyclesFromDraft(cycles),
      confirm_duplicate: confirmDuplicate,
    };
    const handlers = { onSuccess: (saved: GlobalCollege) => onSaved(saved) };

    if (college) {
      update.mutate({ id: college.id, data }, handlers);
    } else {
      create.mutate({ data }, handlers);
    }
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
    save,
    saving: mutation.isPending,
    validationError,
    isDuplicateWarning,
    error: isDuplicateWarning ? null : mutation.error,
  };
}
