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
import { useTaskSummary } from "./useTaskSummary";

export default function StudentTaskList() {
  const queryClient = useQueryClient();
  const [selectedTask, setSelectedTask] = useState<TodoItem | null>(null);
  const [sortAscending, setSortAscending] = useState(true);

  const todosQuery = useGetTodoItems();
  const updateCompleted = useUpdateTodoItemCompleted();

  // Hooks must run before any early return
  const todos: TodoItem[] = Array.isArray(todosQuery.data)
    ? todosQuery.data
    : [];
  const summary = useTaskSummary(todos, sortAscending);

  if (todosQuery.isLoading) return <p>Loading...</p>;
  if (todosQuery.error) return <p>Something went wrong.</p>;

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
          value={`${summary.completedCount}/${summary.totalCount}`}
          caption={`${summary.percentFinished}% finished`}
        />
        <TaskStatCard
          label="Tasks due this week"
          value={String(summary.dueThisWeek)}
          caption={`${summary.daysRemaining} days remaining`}
        />
        <TaskStatCard
          label="High priority"
          value={String(summary.highPriority)}
          caption="due within 24 hours"
        />
      </div>

      {/* Date groups */}
      <div className="flex flex-col gap-6">
        {summary.groups.length === 0 && <p>No tasks yet.</p>}

        {summary.groups.map((group, index) => (
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
