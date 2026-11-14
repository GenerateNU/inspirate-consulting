import { useState } from "react";

import {
  useListPersonalCollegeApplications,
  useListGlobalColleges,
  useCreatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplication,
  useDeletePersonalCollegeApplication,
} from "../../api/endpoints";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  UpdatePersonalCollegeApplicationRequestBodyApplicationStatus as ApplicationStatus,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { useAuthContext } from "@/contexts/use-auth-context";
import { cycleCloseDate, findCycle } from "../utils/shared";

export const sortOptions = [
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

// one school picked in the Add School modal
export interface NewSchool {
  globalCollegeId: number;
  applicationType: ApplicationType;
  category: Category;
  interest: Interest | null;
}

// all of My List's data, page state, and changes
// the page only renders what this returns
export function useMyCollegeList() {
  const [sortKey, setSortKey] = useState<SortKey>("deadline");
  // the application waiting on the remove confirmation
  const [pendingRemoval, setPendingRemoval] =
    useState<PersonalCollegeApplication | null>(null);

  // only students have a list; a counselor calling this without student_id gets a 400
  const isStudent = useAuthContext().role === "student";

  const {
    data: applicationsData,
    error,
    isLoading,
    refetch,
  } = useListPersonalCollegeApplications(undefined, {
    query: { enabled: isStudent },
  });

  const { data: collegesData, isLoading: collegesLoading } =
    useListGlobalColleges();

  const createApplication = useCreatePersonalCollegeApplication();
  const updateApplication = useUpdatePersonalCollegeApplication();
  const deleteApplication = useDeletePersonalCollegeApplication();
  const [adding, setAdding] = useState(false);
  const [addError, setAddError] = useState<unknown>(null);

  const applications: PersonalCollegeApplication[] = applicationsData ?? [];
  const colleges: GlobalCollege[] = collegesData ?? [];

  const collegeById = new Map(colleges.map((college) => [college.id, college]));

  const collegeName = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "Unknown college";

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

  // saves every school picked in the Add School modal. Create doesn't take
  // interest, so a picked interest is set with a follow-up update.
  // returns the schools that failed so the modal can keep them selected
  const addSchools = async (schools: NewSchool[]): Promise<NewSchool[]> => {
    setAdding(true);
    setAddError(null);
    const failed: NewSchool[] = [];

    for (const school of schools) {
      try {
        const created = await createApplication.mutateAsync({
          data: {
            global_college_id: school.globalCollegeId,
            application_type: school.applicationType,
            category: school.category,
          },
        });
        if (school.interest) {
          await updateApplication.mutateAsync({
            id: created.id,
            data: {
              global_college_id: school.globalCollegeId,
              application_type: school.applicationType,
              category: school.category,
              interest: school.interest,
              application_status: ApplicationStatus.not_started,
            },
          });
        }
      } catch (err) {
        failed.push(school);
        setAddError(err);
      }
    }

    await refetch();
    setAdding(false);
    return failed;
  };

  const confirmRemoval = () => {
    if (!pendingRemoval) return;
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
    isStudent,
    isLoading: isStudent && (isLoading || collegesLoading),
    error,

    // what's on the list
    applications: sorted,
    colleges,
    collegeById,
    collegeName,
    addedCollegeIds: new Set(applications.map((app) => app.global_college_id)),

    // sorting
    sortKey,
    setSortKey,

    // add school modal
    addSchools,
    adding,
    addError,

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
  };
}
