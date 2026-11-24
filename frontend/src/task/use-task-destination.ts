import {
  useGetEssaysFromStudent,
  useGetExtracurricular,
  useGetGlobalCollege,
  useGetMedia,
  useListPersonalCollegeApplications,
} from "../api/endpoints";
import type { TodoItem } from "../api/models";
import { useAuthContext } from "@/contexts/use-auth-context";

export interface TaskDestination {
  // the linked item's name, e.g. "Track and Field"; undefined while loading or unlinked
  name: string | undefined;
  // where "Start" goes; null for a task with no link
  path: string | null;
}

// Where a task's "Start" button goes, and what it's linked to.
// A task linked to a type but no specific item goes to that type's page
// (or, for extracurriculars, straight to adding one).
export function useTaskDestination(task: TodoItem | null): TaskDestination {
  const { studentID } = useAuthContext();
  const type = task?.link_type ?? null;
  const collegeId = task?.global_college_id ?? null;
  const extracurricularId = task?.extracurricular_id ?? null;
  const mediaId = task?.media_id ?? null;
  const essayId = task?.essay_id ?? null;

  const college = useGetGlobalCollege(collegeId ?? 0, {
    query: { enabled: type === "college" && collegeId !== null },
  });
  // school pages are keyed by the student's application, not the college
  const applications = useListPersonalCollegeApplications(undefined, {
    query: { enabled: type === "college" && collegeId !== null },
  });
  const extracurricular = useGetExtracurricular(extracurricularId ?? 0, {
    query: {
      enabled: type === "extracurricular" && extracurricularId !== null,
    },
  });
  const media = useGetMedia(mediaId ?? "", {
    query: { enabled: type === "module" && mediaId !== null },
  });
  const essays = useGetEssaysFromStudent(studentID ?? "", undefined, {
    query: { enabled: type === "essay" && essayId !== null && !!studentID },
  });

  switch (type) {
    case "college": {
      const application = (applications.data ?? []).find(
        (app) => app.global_college_id === collegeId,
      );
      return {
        name: college.data?.school_name,
        path: application
          ? `/college-applications/${application.id}`
          : "/college-applications",
      };
    }
    case "extracurricular":
      return extracurricularId === null
        ? { name: undefined, path: "/extracurriculars/new" }
        : {
            name: extracurricular.data?.name,
            path: `/extracurriculars/${extracurricularId}`,
          };
    case "module":
      return { name: media.data?.title, path: "/videos" };
    case "essay":
      return {
        name: essays.data?.essays?.find((e) => e.id === essayId)?.type,
        path: "/essay",
      };
    default:
      return { name: undefined, path: null };
  }
}
