import { useMemo } from "react";
import type { TodoItem } from "../api/models";

const DAY_MS = 24 * 60 * 60 * 1000;

export type TaskGroup = { heading: string; tasks: TodoItem[] };

export type TaskSummary = {
  completedCount: number;
  totalCount: number;
  percentFinished: number;
  dueThisWeek: number;
  daysRemaining: number;
  highPriority: number;
  groups: TaskGroup[];
};

// Formats date in month first then day using ordinal suffix e.g. "September 23rd"
function formatDateHeading(date: Date) {
  const day = date.getDate();
  const suffix =
    day % 10 === 1 && day !== 11
      ? "st"
      : day % 10 === 2 && day !== 12
        ? "nd"
        : day % 10 === 3 && day !== 13
          ? "rd"
          : "th";
  return `${date.toLocaleDateString("en-US", { month: "long" })} ${day}${suffix}`;
}

function isDueWithin(task: TodoItem, now: number, windowMs: number) {
  if (!task.deadline) return false;
  const due = new Date(task.deadline).getTime();
  return due >= now && due <= now + windowMs;
}

export function useTaskSummary(
  todos: TodoItem[],
  sortAscending: boolean,
): TaskSummary {
  return useMemo(() => {
    const now = Date.now();

    // Stats
    const completedCount = todos.filter((t) => t.completed_at != null).length;
    const percentFinished = todos.length
      ? Math.round((completedCount / todos.length) * 100)
      : 0;

    const incomplete = todos.filter((t) => t.completed_at == null);
    const dueThisWeek = incomplete.filter((t) =>
      isDueWithin(t, now, 7 * DAY_MS),
    ).length;

    // Days left until the end of the week (Sunday)
    const daysRemaining = 7 - new Date().getDay();

    // Placeholder until there's a real priority field: incomplete tasks due within 24 hours
    const highPriority = incomplete.filter((t) =>
      isDueWithin(t, now, DAY_MS),
    ).length;

    // Sort
    const sorted = [...todos].sort((a, b) => {
      const aTime = a.deadline ? new Date(a.deadline).getTime() : Infinity;
      const bTime = b.deadline ? new Date(b.deadline).getTime() : Infinity;
      return sortAscending ? aTime - bTime : bTime - aTime;
    });

    const groups: TaskGroup[] = [];
    for (const task of sorted) {
      const heading = task.deadline
        ? formatDateHeading(new Date(task.deadline))
        : "No due date";
      const last = groups[groups.length - 1];
      if (last && last.heading === heading) {
        last.tasks.push(task);
      } else {
        groups.push({ heading, tasks: [task] });
      }
    }

    return {
      completedCount,
      totalCount: todos.length,
      percentFinished,
      dueThisWeek,
      daysRemaining,
      highPriority,
      groups,
    };
  }, [todos, sortAscending]);
}
