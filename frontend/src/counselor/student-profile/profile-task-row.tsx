import { Link } from "react-router-dom";
import { Check, Pencil, Trash2 } from "lucide-react";

import type { TodoItem } from "../../api/models";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import Tag from "../../college/utils/components/tag";
import {
  UNLINKED_ACCENT_CLASS_NAME,
  taskTypeStyle,
  type TaskType,
} from "../../task/task-type-styles";
import { formatShortDate } from "./utils/format-short-date";

// the empty checkbox's ring takes the task type's color, per the hi-fi
const checkboxBorderClassName: Record<TaskType, string> = {
  college: "border-mint",
  essay: "border-blue",
  extracurricular: "border-tan",
  module: "border-greyish-brown",
};

interface ProfileTaskRowProps {
  task: TodoItem;
  studentId: string;
  onToggleComplete: (task: TodoItem) => void;
  onDelete: (task: TodoItem) => void;
}

// one task on the counselor's view of a student: type edge, checkbox, icon,
// title, type pill and date, then edit and delete
export default function ProfileTaskRow({
  task,
  studentId,
  onToggleComplete,
  onDelete,
}: ProfileTaskRowProps) {
  const completed = task.completed_at != null;
  const style = taskTypeStyle(task);
  const Icon = style?.icon;

  let date = null;
  if (completed && task.completed_at) {
    date = `Completed ${formatShortDate(task.completed_at)}`;
  } else if (task.deadline) {
    date = `Due ${formatShortDate(task.deadline)}`;
  }

  return (
    <li
      className={cn(
        "relative flex items-center gap-3 overflow-clip rounded-md border py-3 pr-4 pl-5",
        completed
          ? "border-dashed border-greyish-brown bg-off-white"
          : "border-cream bg-white",
      )}
    >
      <div
        aria-hidden="true"
        className={cn(
          "absolute -top-px -bottom-px -left-px w-1",
          style?.accentClassName ?? UNLINKED_ACCENT_CLASS_NAME,
        )}
      />

      <button
        type="button"
        role="checkbox"
        data-slot="task-checkbox"
        aria-checked={completed}
        aria-label={`Mark "${task.title}" as ${completed ? "not done" : "done"}`}
        onClick={() => onToggleComplete(task)}
        className={cn(
          "flex size-[18px] shrink-0 items-center justify-center rounded-full border-[1.5px]",
          completed
            ? "border-light-black bg-light-black text-white"
            : cn(
                "bg-white",
                task.link_type
                  ? checkboxBorderClassName[task.link_type]
                  : "border-warm-grey",
              ),
        )}
      >
        {completed && <Check className="size-3" strokeWidth={4} />}
      </button>

      {Icon && (
        <span
          aria-hidden="true"
          className={cn(
            "flex size-7 shrink-0 items-center justify-center rounded-md",
            style.iconClassName,
            completed && "opacity-50",
          )}
        >
          <Icon className="size-4" />
        </span>
      )}

      <div className="flex min-w-0 flex-1 flex-col gap-1 [contain:inline-size]">
        <p
          className={cn(
            "truncate font-host text-sm leading-[normal]",
            completed ? "text-dark-greyish-brown" : "text-charcoal",
          )}
        >
          {task.title}
        </p>
        <div className="flex items-center gap-2">
          {style && (
            <Tag
              tone={style.tagTone}
              className={cn(
                "px-2 py-0.5 text-[10px]",
                style.tagClassName,
                completed && "opacity-70",
              )}
            >
              {style.label}
            </Tag>
          )}
          {date && (
            <span className="font-host text-[11px] text-dark-greyish-brown">
              {date}
            </span>
          )}
        </div>
      </div>

      <div className={cn("flex shrink-0 gap-1", completed && "opacity-50")}>
        <Button
          asChild
          variant="ghost"
          size="icon"
          className="size-7 text-dark-greyish-brown hover:bg-off-white"
        >
          <Link
            to={`/counselor/students/${studentId}/tasks/${task.id}/edit`}
            aria-label={`Edit "${task.title}"`}
          >
            <Pencil className="size-3.5" strokeWidth={1.75} />
          </Link>
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label={`Delete "${task.title}"`}
          onClick={() => onDelete(task)}
          className="size-7 text-dark-greyish-brown hover:bg-off-white"
        >
          <Trash2 className="size-3.5" strokeWidth={1.75} />
        </Button>
      </div>
    </li>
  );
}
