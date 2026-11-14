import { useState } from "react";

import {
  useListPersonalCollegeApplications,
  useListGlobalColleges,
  useCreatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplicationRank,
  useDeletePersonalCollegeApplication,
} from "../../api/endpoints";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { cycleCloseDate, findCycle } from "../utils/shared";

export const sortOptions = [
  { value: "rank", label: "My ranking" },
  { value: "deadline", label: "Deadline" },
  { value: "name", label: "Name" },
  { value: "category", label: "Classification" },
  { value: "interest", label: "Interest" },
] as const;

export type SortKey = (typeof sortOptions)[number]["value"];

const categoryOrder: string[] = [
  Category.reach,
  Category.target,
  Category.safety,
];
const interestOrder: string[] = [
  Interest.top_choice,
  Interest.strong_contender,
  Interest.considering,
];

// position in a fixed order; values not in it (e.g. no interest) go last
const orderOf = (order: string[], value: string | null) => {
  const index = value == null ? -1 : order.indexOf(value);
  return index === -1 ? order.length : index;
};

// all of My List's data, page state, and changes
// the page only renders what this returns
export function useMyCollegeList() {
  const [sortKey, setSortKey] = useState<SortKey>("rank");
  const [draggedId, setDraggedId] = useState<number | null>(null);
  // the application waiting on the remove confirmation
  const [pendingRemoval, setPendingRemoval] =
    useState<PersonalCollegeApplication | null>(null);

  const {
    data: applicationsData,
    error,
    isLoading,
    refetch,
  } = useListPersonalCollegeApplications();

  const { data: collegesData, isLoading: collegesLoading } =
    useListGlobalColleges();

  const createApplication = useCreatePersonalCollegeApplication();
  const updateRank = useUpdatePersonalCollegeApplicationRank();
  const deleteApplication = useDeletePersonalCollegeApplication();

  const busy =
    createApplication.isPending ||
    updateRank.isPending ||
    deleteApplication.isPending;

  const mutationError =
    createApplication.error ?? updateRank.error ?? deleteApplication.error;

  const applications: PersonalCollegeApplication[] = applicationsData ?? [];
  const colleges: GlobalCollege[] = collegesData ?? [];

  const collegeById = new Map(colleges.map((college) => [college.id, college]));

  const collegeName = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "Unknown college";

  // ranked applications in rank order, then unranked ones
  const ranked = applications
    .filter((application) => application.rank != null)
    .sort((a, b) => (a.rank ?? 0) - (b.rank ?? 0));
  const unranked = applications.filter(
    (application) => application.rank == null,
  );
  const byRank = [...ranked, ...unranked];

  const byName = (
    a: PersonalCollegeApplication,
    b: PersonalCollegeApplication,
  ) => collegeName(a).localeCompare(collegeName(b));

  // closing date of the cycle the student picked; schools without one go last
  const closeDate = (app: PersonalCollegeApplication) =>
    cycleCloseDate(
      findCycle(collegeById.get(app.global_college_id), app.application_type),
    ) ?? Number.POSITIVE_INFINITY;

  const sorted = (() => {
    switch (sortKey) {
      case "rank":
        return byRank;
      case "name":
        return [...applications].sort(byName);
      case "deadline":
        return [...applications].sort(
          (a, b) => closeDate(a) - closeDate(b) || byName(a, b),
        );
      case "category":
        return [...applications].sort(
          (a, b) =>
            orderOf(categoryOrder, a.category) -
              orderOf(categoryOrder, b.category) || byName(a, b),
        );
      case "interest":
        return [...applications].sort(
          (a, b) =>
            orderOf(interestOrder, a.interest) -
              orderOf(interestOrder, b.interest) || byName(a, b),
        );
    }
  })();

  // clears the last mutation's error so it doesn't linger once a new change starts
  const resetMutations = () => {
    createApplication.reset();
    updateRank.reset();
    deleteApplication.reset();
  };

  // refetch the full list after a rank change, since the backend shifts other applications too
  const setRank = (id: number, rank: number | null) => {
    resetMutations();
    updateRank.mutate({ id, data: { rank } }, { onSuccess: () => refetch() });
  };

  // only "My ranking" order can be dragged.
  // dropping on a ranked card moves the dragged school into that position
  // (the backend shifts everything in between); dropping on an unranked card
  // puts it at the end of the ranking
  const canReorder = sortKey === "rank" && !busy;
  const handleDrop = (target: PersonalCollegeApplication) => {
    if (draggedId == null || draggedId === target.id) return;
    const dragged = applications.find((app) => app.id === draggedId);
    if (!dragged) return;

    if (target.rank != null) {
      setRank(draggedId, target.rank);
    } else if (dragged.rank == null) {
      setRank(draggedId, ranked.length + 1);
    } else {
      setRank(draggedId, ranked.length);
    }
  };

  // add a school, ranked at the bottom of the list
  const addSchool = (
    globalCollegeId: number,
    applicationType: ApplicationType,
    category: Category,
    onSuccess?: () => void,
  ) => {
    resetMutations();
    const nextRank = ranked.length + 1;

    createApplication.mutate(
      {
        data: {
          global_college_id: globalCollegeId,
          application_type: applicationType,
          category,
        },
      },
      {
        onSuccess: (created) => {
          onSuccess?.();
          updateRank.mutate(
            { id: created.id, data: { rank: nextRank } },
            { onSettled: () => refetch() },
          );
        },
      },
    );
  };

  const confirmRemoval = () => {
    if (!pendingRemoval) return;
    resetMutations();
    deleteApplication.mutate(
      { id: pendingRemoval.id },
      {
        onSuccess: async () => {
          await refetch();
          setPendingRemoval(null);
        },
      },
    );
  };

  return {
    // loading the list
    isLoading: isLoading || collegesLoading,
    error,

    // what's on the list
    applications: sorted,
    colleges,
    collegeById,
    collegeName,
    addedCollegeIds: new Set(applications.map((app) => app.global_college_id)),

    // sorting and drag to reorder
    sortKey,
    setSortKey,
    canReorder,
    draggedId,
    setDraggedId,
    handleDrop,

    // add school modal
    addSchool,
    adding: createApplication.isPending,
    addError: createApplication.error,

    // remove confirmation
    pendingRemoval,
    requestRemoval: (app: PersonalCollegeApplication) => {
      deleteApplication.reset();
      setPendingRemoval(app);
    },
    cancelRemoval: () => setPendingRemoval(null),
    confirmRemoval,
    removing: deleteApplication.isPending,
    removeError: deleteApplication.error,

    // any change in progress / the last change's error
    busy,
    mutationError,
  };
}
