import { Check } from "lucide-react";

import type { TodoItem } from "../api/models";
import { cn } from "@/lib/utils";
import Tag from "../college/utils/components/tag";
import { UNLINKED_ACCENT_CLASS_NAME, taskTypeStyle } from "./task-type-styles";

type TaskRowProps = {
  task: TodoItem;
  // collapsed away (a completed task while completed ones are hidden)
  hidden: boolean;
  onToggleComplete: (task: TodoItem) => void;
  onOpen: (task: TodoItem) => void;
};

// The hi-fi task card: type color on the left, a round checkbox, then the tag, title and description
export default function TaskRow({
  task,
  hidden,
  onToggleComplete,
  onOpen,
}: TaskRowProps) {
  const completed = task.completed_at != null;
  const style = taskTypeStyle(task);
  const Icon = style?.icon;

  return (
    // grid rows 1fr -> 0fr animates the height, so hiding collapses smoothly
    <li
      aria-hidden={hidden}
      inert={hidden}
      className={cn(
        "grid transition-[grid-template-rows,opacity] duration-300 ease-out motion-reduce:transition-none",
        hidden ? "grid-rows-[0fr] opacity-0" : "grid-rows-[1fr] opacity-100",
      )}
    >
      <div className="min-h-0 overflow-hidden">
        <div className="pb-3">
          <div
            className={cn(
              "relative flex items-center gap-4 overflow-clip rounded-md border py-4 pr-6 pl-6 transition-colors duration-300 motion-reduce:transition-none",
              completed
                ? "border-dashed border-blue/60 bg-off-white"
                : "border-cream bg-white shadow-card",
            )}
          >
            <div
              aria-hidden="true"
              className={cn(
                "absolute top-[-1px] bottom-[-1px] left-[-1px] w-1",
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
                "flex size-5 shrink-0 items-center justify-center rounded-full border transition-colors duration-300 motion-reduce:transition-none",
                completed
                  ? "border-light-black bg-light-black text-white"
                  : "border-warm-grey bg-white hover:border-dark-green",
              )}
            >
              {completed && (
                <Check
                  className="size-3 animate-in duration-300 zoom-in-50 motion-reduce:animate-none"
                  strokeWidth={4}
                />
              )}
            </button>

            <button
              type="button"
              data-slot="task-open"
              onClick={() => onOpen(task)}
              className="flex min-w-0 flex-1 items-center gap-3 text-left"
            >
              {Icon && (
                <span
                  className={cn(
                    "flex size-8 shrink-0 items-center justify-center rounded-md transition-opacity duration-300",
                    style.iconClassName,
                    completed && "opacity-50",
                  )}
                >
                  <Icon className="size-4" aria-hidden />
                </span>
              )}
              <span className="flex min-w-0 flex-col gap-1">
                {style && (
                  <Tag
                    tone={style.tagTone}
                    className={cn(
                      "self-start px-2 py-0.5 text-[10px] transition-opacity duration-300",
                      style.tagClassName,
                      completed && "opacity-50",
                    )}
                  >
                    {style.label}
                  </Tag>
                )}
                <span
                  className={cn(
                    "truncate font-host text-base leading-[normal] transition-colors duration-300",
                    completed ? "text-warm-grey" : "text-charcoal",
                  )}
                >
                  {task.title}
                </span>
                {task.todo_description && (
                  <span
                    className={cn(
                      "truncate font-host text-xs leading-[normal]",
                      completed ? "text-greyish-brown" : "text-warm-grey-dark",
                    )}
                  >
                    {task.todo_description}
                  </span>
                )}
              </span>
            </button>
          </div>
        </div>
      </div>
    </li>
  );
}
