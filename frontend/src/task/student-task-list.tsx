import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { TodoItem } from "../api/models";
import {
  useGetTodoItems,
  useUpdateTodoItemCompleted,
  getGetTodoItemsQueryKey,
} from "../api/endpoints";
import TaskStatCard from "./task-stat-card";
import TaskRow from "./task-row";
import TaskDrawer from "./task-drawer";

const DAY_MS = 24 * 60 * 60 * 1000;

// "September 23rd"
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

export default function StudentTaskList() {
  const queryClient = useQueryClient();
  const [selectedTask, setSelectedTask] = useState<TodoItem | null>(null);
  const [sortAscending, setSortAscending] = useState(true);

  const todosQuery = useGetTodoItems();
  const updateCompleted = useUpdateTodoItemCompleted();

  if (todosQuery.isLoading) return <p>Loading...</p>;
  if (todosQuery.error) return <p>Something went wrong.</p>;

  const todos: TodoItem[] = Array.isArray(todosQuery.data)
    ? todosQuery.data
    : [];

  // ---- Stats ----
  const now = Date.now();
  const completedCount = todos.filter((t) => t.completed_at != null).length;
  const percentFinished = todos.length
    ? Math.round((completedCount / todos.length) * 100)
    : 0;

  const incomplete = todos.filter((t) => t.completed_at == null);
  const dueThisWeek = incomplete.filter((t) => {
    if (!t.deadline) return false;
    const due = new Date(t.deadline).getTime();
    return due >= now && due <= now + 7 * DAY_MS;
  }).length;

  // Days left until the end of the week (Sunday)
  const daysRemaining = 7 - new Date().getDay();

  // Placeholder until there's a real priority field: incomplete tasks due within 24 hours
  const highPriority = incomplete.filter((t) => {
    if (!t.deadline) return false;
    const due = new Date(t.deadline).getTime();
    return due >= now && due <= now + DAY_MS;
  }).length;

  // ---- Sort, then group by due date ----
  const sorted = [...todos].sort((a, b) => {
    const aTime = a.deadline ? new Date(a.deadline).getTime() : Infinity;
    const bTime = b.deadline ? new Date(b.deadline).getTime() : Infinity;
    return sortAscending ? aTime - bTime : bTime - aTime;
  });

  const groups: { heading: string; tasks: TodoItem[] }[] = [];
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

  // ---- Toggle completion ----
  const handleToggleComplete = (task: TodoItem) => {
    updateCompleted.mutate(
      { id: task.id, data: { completed: task.completed_at == null } },
      {
        onSuccess: () => 
          queryClient.invalidateQueries({
            queryKey: getGetTodoItemsQueryKey(),
        }),
      },
    );
  };

  return (
    <div className="flex flex-col gap-6 p-10">
      {/* Header + legend */}
      <div>
        <h1 className="text-3xl font-semibold">Tasks</h1>
        <div className="mt-2 flex gap-16 text-xs">
          <span>Extracurriculars</span>
          <span>Essays</span>
          <span>Modules</span>
        </div>
      </div>

      {/* Stat cards */}
      <div className="flex gap-4 pl-10">
        <TaskStatCard
          label="Tasks completed"
          value={`${completedCount}/${todos.length}`}
          caption={`${percentFinished}% finished`}
        />
        <TaskStatCard
          label="Tasks due this week"
          value={String(dueThisWeek)}
          caption={`${daysRemaining} days remaining`}
        />
        <TaskStatCard
          label="High priority"
          value={String(highPriority)}
          caption="due within 24 hours"
        />
      </div>

      {/* Date groups */}
      <div className="flex flex-col gap-6">
        {groups.length === 0 && <p>No tasks yet.</p>}

        {groups.map((group, index) => (
          <section key={group.heading} className="flex flex-col gap-3">
            <div className="flex items-center gap-4">
              <h2 className="shrink-0 text-xl font-semibold">
                {group.heading}
              </h2>
              <div className="h-px flex-1 bg-gray-400" />
              {/* Sort toggle sits on the first heading's line, like the Figma */}
              {index === 0 && (
                <button
                  type="button"
                  onClick={() => setSortAscending((prev) => !prev)}
                  aria-label={
                    sortAscending ? "Sort latest first" : "Sort earliest first"
                  }
                  className="rounded border px-2 py-1"
                >
                  ↓↑
                </button>
              )}
            </div>

            <div className="flex flex-col gap-3 pl-10">
              {group.tasks.map((task) => (
                <TaskRow
                  key={task.id}
                  task={task}
                  title={task.todo_description}
                  onToggleComplete={handleToggleComplete}
                  onOpen={setSelectedTask}
                />
              ))}
            </div>
          </section>
        ))}
      </div>

      <TaskDrawer task={selectedTask} onClose={() => setSelectedTask(null)} />
    </div>
  );
}
