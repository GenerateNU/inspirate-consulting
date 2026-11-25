import { useEffect, useRef, useState } from "react";
import { ChevronDown, ChevronUp } from "lucide-react";

import type { TodoItem } from "../api/models";
import { errorMessage } from "../api/client";
import PageLinks, { type PageLink } from "@/components/page-links";
import { Skeleton } from "@/components/ui/skeleton";
import TaskDrawer from "./task-drawer";
import TaskRow from "./task-row";
import TaskStatCard from "./task-stat-card";
import { useTaskGroups, type TaskGroup } from "./use-task-groups";
import { useStudentTasks, useToggleTaskCompleted } from "./use-student-tasks";

// Temporary links to this page's sub-pages (see the Figma sitemap).
const pageLinks: PageLink[] = [
  { to: "/counselor-students", label: "My Students", role: "Counselor" },
  { to: "/search", label: "Claim Students", role: "Counselor" },
];

// how long a just-completed task stays visible (to finish its animation) before hiding
const COMPLETE_LINGER_MS = 900;

// PLACEHOLDER: the numbers come from the todo stats endpoint (#221); wire them in once it lands
function TaskStatCards() {
  return (
    <div className="flex flex-wrap gap-4">
      <TaskStatCard
        label="This Week's Tasks"
        value="–"
        valueSuffix="/– Completed"
        progress={0}
        caption="–% finished"
      />
      <TaskStatCard label="Applications due this week" value="–" />
    </div>
  );
}

function CompletedToggle({
  shown,
  count,
  onToggle,
}: {
  shown: boolean;
  count: number;
  onToggle: () => void;
}) {
  const Chevron = shown ? ChevronUp : ChevronDown;
  return (
    <div className="flex items-center gap-3 px-2">
      <div className="h-px flex-1 bg-cream" />
      <button
        type="button"
        data-slot="icon-button"
        aria-expanded={shown}
        aria-label={
          shown
            ? "Hide completed tasks"
            : `Show ${count} completed ${count === 1 ? "task" : "tasks"}`
        }
        title={shown ? "Hide completed" : "Show completed"}
        onClick={onToggle}
        className="flex size-6 items-center justify-center rounded-md text-light-black hover:bg-cream"
      >
        <Chevron className="size-4" />
      </button>
      <div className="h-px flex-1 bg-cream" />
    </div>
  );
}

function TaskGroupSection({
  group,
  showCompleted,
  lingering,
  onToggleShown,
  onToggleComplete,
  onOpen,
}: {
  group: TaskGroup;
  showCompleted: boolean;
  lingering: Set<string>;
  onToggleShown: () => void;
  onToggleComplete: (task: TodoItem) => void;
  onOpen: (task: TodoItem) => void;
}) {
  const completedCount = group.tasks.filter(
    (task) => task.completed_at != null,
  ).length;

  return (
    <section className="flex flex-col gap-3">
      <h2 className="font-jakarta text-base leading-[normal] font-semibold text-charcoal">
        {group.heading}
      </h2>
      <ul className="flex flex-col">
        {group.tasks.map((task) => (
          <TaskRow
            key={task.id}
            task={task}
            hidden={
              task.completed_at != null &&
              !showCompleted &&
              !lingering.has(task.id)
            }
            onToggleComplete={onToggleComplete}
            onOpen={onOpen}
          />
        ))}
      </ul>
      {completedCount > 0 && (
        <CompletedToggle
          shown={showCompleted}
          count={completedCount}
          onToggle={onToggleShown}
        />
      )}
    </section>
  );
}

// The student's task dashboard: tasks grouped by due date, with completed ones shown or hidden per day
export default function StudentTaskList() {
  const [selectedTask, setSelectedTask] = useState<TodoItem | null>(null);
  // days whose completed tasks are hidden (all shown by default)
  const [hiddenDays, setHiddenDays] = useState<Set<string>>(new Set());
  // tasks just checked off in a hidden day, kept visible while they animate
  const [lingering, setLingering] = useState<Set<string>>(new Set());
  const timers = useRef<number[]>([]);

  const tasksQuery = useStudentTasks();
  const toggleCompleted = useToggleTaskCompleted();
  const groups = useTaskGroups(tasksQuery.data ?? []);

  useEffect(() => {
    const pending = timers.current;
    return () => pending.forEach((id) => window.clearTimeout(id));
  }, []);

  const toggleDay = (key: string) =>
    setHiddenDays((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });

  const handleToggleComplete = (task: TodoItem) => {
    if (task.completed_at == null) {
      setLingering((prev) => new Set(prev).add(task.id));
      timers.current.push(
        window.setTimeout(
          () =>
            setLingering((prev) => {
              const next = new Set(prev);
              next.delete(task.id);
              return next;
            }),
          COMPLETE_LINGER_MS,
        ),
      );
    }
    toggleCompleted.mutate(task);
  };

  // keep the drawer showing the latest copy of the open task
  const openTask =
    selectedTask &&
    (tasksQuery.data?.find((task) => task.id === selectedTask.id) ??
      selectedTask);

  return (
    <>
      <PageLinks links={pageLinks} />
      <div className="flex min-h-full w-full flex-col gap-8 bg-off-white p-12 max-lg:px-6">
        <h1 className="font-jakarta text-2xl leading-[normal] font-semibold text-charcoal">
          Tasks
        </h1>

        <TaskStatCards />

        {toggleCompleted.error && (
          <p role="alert" className="font-host text-sm text-danger">
            Couldn&apos;t update that task:{" "}
            {errorMessage(toggleCompleted.error)}
          </p>
        )}

        <div className="flex max-w-228 flex-col gap-8">
          {tasksQuery.isLoading &&
            Array.from({ length: 3 }, (_, i) => (
              <div key={i} className="flex flex-col gap-3">
                <Skeleton className="h-5 w-40 bg-cream/60" />
                <Skeleton className="h-20 w-full rounded-md bg-cream/60" />
              </div>
            ))}

          {tasksQuery.error && (
            <p className="font-host text-sm text-danger">
              Couldn&apos;t load your tasks: {errorMessage(tasksQuery.error)}
            </p>
          )}

          {!tasksQuery.isLoading &&
            !tasksQuery.error &&
            groups.length === 0 && (
              <p className="font-host text-sm text-warm-grey-dark">
                No tasks yet. Tasks your counselor assigns will show up here.
              </p>
            )}

          {groups.map((group) => (
            <TaskGroupSection
              key={group.key}
              group={group}
              showCompleted={!hiddenDays.has(group.key)}
              lingering={lingering}
              onToggleShown={() => toggleDay(group.key)}
              onToggleComplete={handleToggleComplete}
              onOpen={setSelectedTask}
            />
          ))}
        </div>

        <TaskDrawer task={openTask} onClose={() => setSelectedTask(null)} />
      </div>
    </>
  );
}
