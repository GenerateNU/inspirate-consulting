import {
  useGetEssaysFromStudent,
  useGetMediaAccess,
  useListGlobalColleges,
} from "../api/endpoints";
import type { TodoItem } from "../api/models";

const PLACEHOLDER_STUDENT_ID = "00000000-0000-0000-0000-000000000002";

export const DEFAULT_DUE_TIME = "23:59:59";

export type LinkType = "essay" | "media" | "college";

export const LINK_TYPE_LABELS: Record<LinkType, string> = {
  essay: "Essay",
  media: "Video",
  college: "College",
};
export const LINK_TYPES = Object.keys(LINK_TYPE_LABELS) as LinkType[];

export const errorMessage = (body: unknown) => {
  const { message, detail } = (body ?? {}) as {
    message?: unknown;
    detail?: unknown;
  };
  if (typeof message === "string") return message;
  if (typeof detail === "string") return detail;
  return "Something went wrong.";
};

export const formatDate = (value: string) =>
  new Date(value).toLocaleDateString();

export const formatDeadline = (value: string) => {
  const date = new Date(value);
  if (date.getHours() === 23 && date.getMinutes() === 59)
    return date.toLocaleDateString();
  return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" })}`;
};

export function useTodoLinkOptions() {
  const essays = useGetEssaysFromStudent(PLACEHOLDER_STUDENT_ID);
  const media = useGetMediaAccess({ limit: 100 });
  const colleges = useListGlobalColleges();

  const options: Record<LinkType, { id: string; label: string }[]> = {
    essay:
      essays.data?.status === 200
        ? (essays.data.data.essays ?? []).map((e) => ({
            id: e.id,
            label: e.type,
          }))
        : [],
    media:
      media.data?.status === 200
        ? (media.data.data ?? []).map((m) => ({ id: m.id, label: m.title }))
        : [],
    college:
      colleges.data?.status === 200
        ? (colleges.data.data ?? []).map((c) => ({
            id: String(c.id),
            label: c.school_name,
          }))
        : [],
  };

  const labelFor = (task: TodoItem) => {
    const link: [LinkType, string] | null = task.essay_id
      ? ["essay", task.essay_id]
      : task.media_id
        ? ["media", task.media_id]
        : task.global_college_id != null
          ? ["college", String(task.global_college_id)]
          : null;
    if (!link) return null;
    const [type, id] = link;
    const name = options[type].find((o) => o.id === id)?.label;
    return name ? `${LINK_TYPE_LABELS[type]}: ${name}` : LINK_TYPE_LABELS[type];
  };

  return { options, labelFor };
}
