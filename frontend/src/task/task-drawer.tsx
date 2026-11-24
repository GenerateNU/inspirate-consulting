import { useNavigate } from "react-router-dom";

import type { TodoItem } from "../api/models";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { useTaskDestination } from "./use-task-destination";

type TaskDrawerProps = {
  task: TodoItem | null;
  onClose: () => void;
};

// e.g. "Due 10/13 at 11:59 PM EST"
function formatDueDateTime(deadline: string | null) {
  if (!deadline) return "No due date";
  const due = new Date(deadline);
  const date = due.toLocaleDateString("en-US", {
    month: "numeric",
    day: "numeric",
  });
  const time = due.toLocaleTimeString("en-US", {
    hour: "numeric",
    minute: "2-digit",
    timeZoneName: "short",
  });
  return `Due ${date} at ${time}`;
}

// The hi-fi description drawer: the task's details, and Start to go to what it's linked to
export default function TaskDrawer({ task, onClose }: TaskDrawerProps) {
  const navigate = useNavigate();
  const { name, path } = useTaskDestination(task);

  return (
    <Sheet open={task !== null} onOpenChange={(open) => !open && onClose()}>
      <SheetContent className="w-full gap-0 border-l-cream bg-white sm:max-w-80.5">
        {task && (
          <>
            <SheetHeader className="gap-1 px-4 pt-10 pb-0">
              <SheetTitle className="font-host text-lg leading-[normal] font-normal text-charcoal">
                {task.title}
              </SheetTitle>
              {name && (
                <p className="font-host text-sm leading-[normal] text-charcoal">
                  {name}
                </p>
              )}
              <p className="mt-2 font-host text-xs text-warm-grey">
                {formatDueDateTime(task.deadline)}
              </p>
            </SheetHeader>

            <SheetDescription className="px-4 pt-3 font-host text-xs leading-normal whitespace-pre-line text-warm-grey-dark">
              {task.todo_description || "No description."}
            </SheetDescription>

            {path && (
              <SheetFooter className="items-center px-4 pb-8">
                <Button
                  type="button"
                  onClick={() => void navigate(path)}
                  className="h-8 w-19.5 rounded-md bg-dark-green font-host text-xs font-normal text-white hover:bg-forest"
                >
                  Start
                </Button>
              </SheetFooter>
            )}
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
