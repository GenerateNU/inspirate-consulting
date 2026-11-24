import { useMemo } from "react";
import type { TodoItem } from "../api/models";

export type TaskGroup = {
  // local date of the deadline, or "none"
  key: string;
  heading: string;
  tasks: TodoItem[];
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

const dueTime = (task: TodoItem) =>
  task.deadline ? new Date(task.deadline).getTime() : Infinity;

// Tasks grouped by the day they're due, soonest first; tasks with no due date go last
export function useTaskGroups(tasks: TodoItem[]): TaskGroup[] {
  return useMemo(() => {
    const sorted = [...tasks].sort((a, b) => dueTime(a) - dueTime(b));

    const groups: TaskGroup[] = [];
    for (const task of sorted) {
      const due = task.deadline ? new Date(task.deadline) : null;
      const key = due ? due.toDateString() : "none";
      const last = groups[groups.length - 1];
      if (last && last.key === key) {
        last.tasks.push(task);
      } else {
        groups.push({
          key,
          heading: due ? formatDateHeading(due) : "No due date",
          tasks: [task],
        });
      }
    }
    return groups;
  }, [tasks]);
}
