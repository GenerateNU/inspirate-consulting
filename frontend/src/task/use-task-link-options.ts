import { useQueries } from "@tanstack/react-query";

import {
  getListPersonalCollegeApplicationsQueryOptions,
  useGetEssaysFromStudent,
  useGetTodoItems,
  useListGlobalColleges,
  useListMedia,
} from "../api/endpoints";
import {
  GetTodoItemsLinkedTo,
  GetTodoItemsStatus,
  type BulkCreateTodoItemsBodyLinkType,
} from "../api/models";

export type TaskLinkType = BulkCreateTodoItemsBodyLinkType;

export interface LinkOption {
  id: string;
  label: string;
}

// essays shown under their college in the essay picker
export interface EssayGroup {
  collegeId: number | null;
  collegeName: string;
  essays: LinkOption[];
}

export interface TaskLinkOptions {
  colleges: LinkOption[];
  essayGroups: EssayGroup[];
  modules: LinkOption[];
  isLoading: boolean;
}

const MODULE_LIMIT = 100;

// Items a task can link to, following the task assignment flowchart:
// one student sees their own colleges, essays and unfinished modules;
// a group sees the colleges every student shares and all modules.
// Essays and extracurriculars can't be bulk-assigned, so groups get none.
export function useTaskLinkOptions(
  type: TaskLinkType | null,
  studentIds: string[],
): TaskLinkOptions {
  const individualId = studentIds.length === 1 ? studentIds[0] : "";

  const colleges = useListGlobalColleges({
    query: { enabled: type === "college" || type === "essay" },
  });
  const collegeName = new Map(
    (colleges.data ?? []).map((college) => [college.id, college.school_name]),
  );

  // Colleges: one list per student, then keep the ones they all have
  const applications = useQueries({
    queries: studentIds.map((id) =>
      getListPersonalCollegeApplicationsQueryOptions(
        { student_id: id },
        { query: { enabled: type === "college" } },
      ),
    ),
  });
  const sharedCollegeIds = applications.reduce<number[] | null>(
    (shared, result) => {
      const ids = (result.data ?? []).map((app) => app.global_college_id);
      return shared === null ? ids : shared.filter((id) => ids.includes(id));
    },
    null,
  );

  // Essays: only for one student
  const essays = useGetEssaysFromStudent(individualId, undefined, {
    query: { enabled: type === "essay" && individualId !== "" },
  });

  // Modules: everything, minus what an individual student already finished
  const media = useListMedia(
    { limit: MODULE_LIMIT, offset: 0 },
    { query: { enabled: type === "module" } },
  );
  const completedModules = useGetTodoItems(
    {
      student_id: individualId,
      linked_to: GetTodoItemsLinkedTo.module,
      status: GetTodoItemsStatus.completed,
      limit: MODULE_LIMIT,
    },
    { query: { enabled: type === "module" && individualId !== "" } },
  );
  const completedMediaIds = new Set(
    (completedModules.data?.items ?? []).flatMap((task) =>
      task.media_id ? [task.media_id] : [],
    ),
  );

  const essayGroups: EssayGroup[] = [];
  for (const essay of essays.data?.essays ?? []) {
    let group = essayGroups.find((g) => g.collegeId === essay.college_id);
    if (!group) {
      group = {
        collegeId: essay.college_id,
        collegeName:
          essay.college_id === null
            ? "Other essays"
            : (collegeName.get(essay.college_id) ??
              `College #${essay.college_id}`),
        essays: [],
      };
      essayGroups.push(group);
    }
    group.essays.push({ id: essay.id, label: essay.type });
  }

  return {
    colleges: (sharedCollegeIds ?? []).map((id) => ({
      id: String(id),
      label: collegeName.get(id) ?? `College #${id}`,
    })),
    essayGroups,
    modules: (media.data ?? [])
      .filter((m) => !completedMediaIds.has(m.id))
      .map((m) => ({ id: m.id, label: m.title })),
    isLoading:
      colleges.isLoading ||
      applications.some((result) => result.isLoading) ||
      essays.isLoading ||
      media.isLoading ||
      completedModules.isLoading,
  };
}
