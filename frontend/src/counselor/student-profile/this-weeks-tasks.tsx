import { useState } from "react";
import { Link } from "react-router-dom";
import { ListTodo, Plus } from "lucide-react";

import type { TodoItem } from "../../api/models";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ErrorMessage } from "@/components/ui/error-message";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import DeleteTaskDialog from "./delete-task-dialog";
import { useThisWeeksTasks } from "./hooks/use-this-weeks-tasks";
import ProfileTaskRow from "./profile-task-row";
import {
  outlineButtonClassName,
  panelClassName,
  primaryButtonClassName,
  sectionTitleClassName,
} from "./profile-styles";

interface ThisWeeksTasksProps {
  studentId: string;
}

// "This Week's Tasks": what's due this week, what's been done, and links to assign more
export default function ThisWeeksTasks({ studentId }: ThisWeeksTasksProps) {
  const {
    toDo,
    done,
    isLoading,
    error,
    refetch,
    toggleCompleted,
    toggleError,
    deleteTask,
  } = useThisWeeksTasks(studentId);
  const [pendingDelete, setPendingDelete] = useState<TodoItem | null>(null);

  const confirmDelete = () => {
    if (!pendingDelete) return;
    deleteTask.mutate(
      { id: pendingDelete.id },
      { onSuccess: () => setPendingDelete(null) },
    );
  };

  const renderRows = (tasks: TodoItem[]) => (
    <ul className="flex flex-col gap-2">
      {tasks.map((task) => (
        <ProfileTaskRow
          key={task.id}
          task={task}
          studentId={studentId}
          onToggleComplete={toggleCompleted}
          onDelete={(task) => {
            deleteTask.reset();
            setPendingDelete(task);
          }}
        />
      ))}
    </ul>
  );

  let body;
  if (isLoading) {
    body = (
      <div className="flex flex-col gap-2">
        {[0, 1, 2].map((i) => (
          <Skeleton key={i} className="h-14 w-full rounded-md" />
        ))}
      </div>
    );
  } else if (error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load this student's tasks"
        error={error}
        onRetry={refetch}
      />
    );
  } else if (toDo.length === 0 && done.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        Nothing is due or done this week.
      </p>
    );
  } else {
    body = (
      <>
        {toDo.length > 0 && (
          <div className="flex flex-col gap-2">
            <h3 className="font-host text-[10px] tracking-wide text-dark-greyish-brown uppercase">
              To do
            </h3>
            {renderRows(toDo)}
          </div>
        )}
        {done.length > 0 && (
          <div className="flex flex-col gap-2">
            <h3 className="font-host text-[10px] tracking-wide text-dark-greyish-brown uppercase">
              Completed
            </h3>
            {renderRows(done)}
          </div>
        )}
      </>
    );
  }

  return (
    <Card className={cn(panelClassName, "gap-4")}>
      <div className="flex flex-wrap items-center gap-3">
        <ListTodo className="size-5 text-charcoal" strokeWidth={1.75} />
        <h2 className={sectionTitleClassName}>This Week&apos;s Tasks</h2>
        {!isLoading && !error && (
          <span className="font-host text-[11px] text-dark-greyish-brown">
            {toDo.length} to do · {done.length} done
          </span>
        )}
        <div className="ml-auto flex gap-2">
          <Button asChild variant="outline" className={outlineButtonClassName}>
            <Link to={`/counselor/students/${studentId}/tasks`}>
              View Full Dashboard
            </Link>
          </Button>
          <Button asChild className={primaryButtonClassName}>
            <Link to={`/assign-task?studentId=${studentId}`}>
              <Plus className="size-3.5" strokeWidth={2} />
              Assign task
            </Link>
          </Button>
        </div>
      </div>

      <ErrorMessage error={toggleError} prefix="Couldn't update the task" />
      {body}

      <DeleteTaskDialog
        open={pendingDelete != null}
        onOpenChange={(open) => !open && setPendingDelete(null)}
        onDelete={confirmDelete}
        deleting={deleteTask.isPending}
        failed={deleteTask.isError}
      />
    </Card>
  );
}
