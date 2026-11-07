import { useState } from "react";
import { useSearchParams } from "react-router-dom";

import type { ApiError } from "../../api/client";
import {
  useListPersonalCollegeApplications,
  useListGlobalColleges,
  useCreatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplicationRank,
  useDeletePersonalCollegeApplication,
} from "../../api/endpoints";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { deadlineOptions } from "../utils/shared";

// all of My List's data, page state, and changes
// the page only renders what this returns
export function useMyCollegeList() {
  const [draggedId, setDraggedId] = useState<number | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);

  // set when the backend rejects a deadline the college doesn't offer -> shows the popup
  const [deadlineError, setDeadlineError] = useState<{
    globalCollegeId: number;
    applicationType: string;
  } | null>(null);

  // ?college=<id> preselects a college from the colleges list and opens the form
  const [searchParams] = useSearchParams();
  const [showAddForm, setShowAddForm] = useState(searchParams.has("college"));

  const {
    data: applicationsData,
    error,
    isLoading,
    refetch,
  } = useListPersonalCollegeApplications();

  const { data: collegesData, isLoading: collegesLoading } =
    useListGlobalColleges();

  // one hook per kind of change, the row's id is passed when calling mutate
  const createApplication = useCreatePersonalCollegeApplication();
  const updateApplication = useUpdatePersonalCollegeApplication();
  const updateRank = useUpdatePersonalCollegeApplicationRank();
  const deleteApplication = useDeletePersonalCollegeApplication();

  const busy =
    createApplication.isPending ||
    updateApplication.isPending ||
    updateRank.isPending ||
    deleteApplication.isPending;

  const mutationError =
    createApplication.error ??
    updateApplication.error ??
    updateRank.error ??
    deleteApplication.error;

  const applications: PersonalCollegeApplication[] = applicationsData ?? [];
  const colleges: GlobalCollege[] = collegesData ?? [];

  const collegeById = new Map(colleges.map((college) => [college.id, college]));

  const collegeName = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "Unknown college";

  // ranked applications (non null) are shown in order, unranked (null) come after them
  const ranked = applications
    .filter((application) => application.rank != null)
    .sort((a, b) => (a.rank ?? 0) - (b.rank ?? 0));

  const unranked = applications.filter(
    (application) => application.rank == null,
  );

  // colleges already on the list can't be added again
  const appliedIds = new Set(applications.map((app) => app.global_college_id));
  const availableColleges = colleges.filter(
    (college) => !appliedIds.has(college.id),
  );

  // skip the preselect once that college is already on the list
  const requestedId = Number(searchParams.get("college")) || null;
  const preselectedId = availableColleges.some(
    (college) => college.id === requestedId,
  )
    ? requestedId
    : null;

  // backend answers 400 when the college has no date for the chosen deadline
  const isMissingDeadline = (
    error: ApiError,
    globalCollegeId: number,
    applicationType: string,
  ) => {
    const key = deadlineOptions.find(
      ({ value }) => value === applicationType,
    )?.key;
    const college = collegeById.get(globalCollegeId);
    return error.status === 400 && !!key && !!college && !college[key];
  };

  // clears the last mutation's error so it doesn't linger once a new change starts
  const resetMutations = () => {
    createApplication.reset();
    updateApplication.reset();
    updateRank.reset();
    deleteApplication.reset();
  };

  // refetch full list after rank change, since the backend shifts other applications too
  const setRank = (id: number, rank: number | null) => {
    resetMutations();
    updateRank.mutate({ id, data: { rank } }, { onSuccess: () => refetch() });
  };

  // ranks are 1-based, the position in ranked is 0-based
  const rankIndex = (app: PersonalCollegeApplication) => ranked.indexOf(app);
  const moveUp = (app: PersonalCollegeApplication) =>
    setRank(app.id, rankIndex(app));
  const moveDown = (app: PersonalCollegeApplication) =>
    setRank(app.id, rankIndex(app) + 2);
  const addToRanking = (app: PersonalCollegeApplication) =>
    setRank(app.id, ranked.length + 1);
  const unrank = (app: PersonalCollegeApplication) => setRank(app.id, null);

  // dropping on a ranked card moves the dragged application into that card's
  // position (the backend shifts everything in between).
  // dropping on an unranked card unranks it.
  const handleDrop = (target: PersonalCollegeApplication) => {
    if (draggedId == null || draggedId === target.id) return;

    if (target.rank != null) {
      setRank(draggedId, target.rank);
      return;
    }

    const dragged = applications.find(
      (application) => application.id === draggedId,
    );
    if (dragged?.rank != null) setRank(draggedId, null);
  };

  const startEditing = (app: PersonalCollegeApplication) => {
    resetMutations();
    setEditingId(app.id);
  };

  // save category/type changes for an application, then refetch the full list to show the updated values
  const saveEdit = (
    app: PersonalCollegeApplication,
    applicationType: ApplicationType,
    category: Category,
  ) => {
    resetMutations();

    updateApplication.mutate(
      {
        id: app.id,
        data: {
          global_college_id: app.global_college_id,
          application_type: applicationType,
          category,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setEditingId(null);
        },
        // keep the card in edit mode so the student can pick another deadline
        onError: (error) => {
          if (
            isMissingDeadline(error, app.global_college_id, applicationType)
          ) {
            updateApplication.reset();
            setDeadlineError({
              globalCollegeId: app.global_college_id,
              applicationType,
            });
          }
        },
      },
    );
  };

  // add a college to the list as unranked, then refetch the full list to show it
  const createApplicationFor = (
    globalCollegeId: number,
    applicationType: ApplicationType,
    category: Category,
  ) => {
    resetMutations();

    createApplication.mutate(
      {
        data: {
          global_college_id: globalCollegeId,
          application_type: applicationType,
          category,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setShowAddForm(false);
        },
        // keep the form open so the student can pick another deadline
        onError: (error) => {
          if (isMissingDeadline(error, globalCollegeId, applicationType)) {
            createApplication.reset();
            setDeadlineError({ globalCollegeId, applicationType });
          }
        },
      },
    );
  };

  // remove an application from the list, then refetch the full list to show the updated values
  const removeApplication = (app: PersonalCollegeApplication) => {
    if (!window.confirm(`Remove ${collegeName(app)} from your list?`)) return;

    resetMutations();

    deleteApplication.mutate({ id: app.id }, { onSuccess: () => refetch() });
  };

  return {
    // loading the list
    isLoading: isLoading || collegesLoading,
    error,

    // what's on the list
    applications: [...ranked, ...unranked],
    rankedCount: ranked.length,
    rankIndex,
    collegeById,
    collegeName,

    // add form
    showAddForm,
    setShowAddForm,
    availableColleges,
    preselectedId,
    createApplication: createApplicationFor,
    creating: createApplication.isPending,

    // editing a card
    editingId,
    startEditing,
    cancelEditing: () => setEditingId(null),
    saveEdit,
    saving: updateApplication.isPending,
    removeApplication,

    // ranking
    moveUp,
    moveDown,
    addToRanking,
    unrank,
    draggedId,
    setDraggedId,
    handleDrop,

    // any change in progress / the last change's error
    busy,
    mutationError,

    // deadline popup
    deadlineError,
    closeDeadlineError: () => setDeadlineError(null),
  };
}
