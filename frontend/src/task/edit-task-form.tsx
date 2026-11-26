import { Link, useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  getGetTodoItemsQueryKey,
  getListCounselorStudentsQueryKey,
  getTodoItems,
  updateTodoItemDeadline,
  updateTodoItemDescription,
  updateTodoItemLink,
  updateTodoItemTitle,
} from "../api/endpoints";
import type { TodoItem, UpdateTodoItemLinkBody } from "../api/models";
import { dueDateToDeadline } from "./assign-task-form";
import TaskForm, { type TaskFormValues } from "./task-form";
import { ErrorMessage } from "@/components/ui/error-message";

const PAGE_SIZE = 100;

// there's no single-task endpoint, so page through the student's tasks to find it
function useStudentTask(studentId: string, taskId: string) {
  return useQuery({
    queryKey: [...getGetTodoItemsQueryKey({ student_id: studentId }), taskId],
    queryFn: async ({ signal }) => {
      let cursor: string | undefined;
      do {
        const page = await getTodoItems(
          { student_id: studentId, cursor, limit: PAGE_SIZE },
          { signal },
        );
        const task = page.items?.find((item) => item.id === taskId);
        if (task) return task;
        cursor = page.next_cursor ?? undefined;
      } while (cursor);
      return null;
    },
  });
}

// the date input wants YYYY-MM-DD in local time
const deadlineToDueDate = (deadline: string | null) => {
  if (!deadline) return "";
  const date = new Date(deadline);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
};

const linkedItemId = (task: TodoItem) =>
  task.essay_id ??
  task.media_id ??
  (task.global_college_id != null ? String(task.global_college_id) : "");

function toFormValues(task: TodoItem): TaskFormValues {
  return {
    title: task.title,
    description: task.todo_description,
    dueDate: deadlineToDueDate(task.deadline),
    assignTo: {
      all: false,
      graduationYears: [],
      studentIds: [task.student_id],
    },
    connectTo: { type: task.link_type, itemId: linkedItemId(task) },
  };
}

// Edit one student's task. Group tasks can't be edited from here.
export default function EditTaskForm() {
  const { studentId = "", taskId = "" } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const taskQuery = useStudentTask(studentId, taskId);
  const task = taskQuery.data;

  const save = useMutation({
    mutationFn: async ({
      before,
      values,
    }: {
      before: TaskFormValues;
      values: TaskFormValues;
    }) => {
      // each field has its own endpoint, so only send what changed
      const title = values.title.trim();
      if (title !== before.title) await updateTodoItemTitle(taskId, { title });

      const description = values.description.trim();
      if (description !== before.description)
        await updateTodoItemDescription(taskId, {
          new_description: description,
        });

      if (values.dueDate !== before.dueDate)
        await updateTodoItemDeadline(taskId, {
          deadline: dueDateToDeadline(values.dueDate) ?? null,
        });

      const { type, itemId } = values.connectTo;
      if (
        type !== before.connectTo.type ||
        itemId !== before.connectTo.itemId
      ) {
        const body: UpdateTodoItemLinkBody = { link_type: type };
        if (itemId && type === "essay") body.essay_id = itemId;
        if (itemId && type === "college")
          body.global_college_id = Number(itemId);
        if (itemId && type === "module") body.media_id = itemId;
        await updateTodoItemLink(taskId, body);
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: getGetTodoItemsQueryKey(),
      });
      void queryClient.invalidateQueries({
        queryKey: getListCounselorStudentsQueryKey(),
      });
      void navigate(`/counselor/students/${studentId}`);
    },
  });

  if (taskQuery.isLoading)
    return <p className="p-12 font-host text-sm">Loading task...</p>;
  if (taskQuery.error || !task) {
    return (
      <div className="p-12">
        <ErrorMessage
          variant="banner"
          title={taskQuery.error ? "Couldn't load this task" : undefined}
          error={taskQuery.error}
          message={taskQuery.error ? undefined : "This task doesn't exist."}
          onRetry={taskQuery.error ? taskQuery.refetch : undefined}
          action={
            <Link to={`/counselor/students/${studentId}`} className="underline">
              Back to student
            </Link>
          }
        />
      </div>
    );
  }

  const initialValues = toFormValues(task);

  return (
    <TaskForm
      // reset the form if a refetch brings a different version of the task
      key={task.updated_at}
      heading="Edit Task"
      initialValues={initialValues}
      submitLabel="Save task"
      pendingLabel="Saving..."
      isPending={save.isPending}
      error={save.error}
      assigneeReadOnly
      // the backend won't clear a description, but one that was empty can stay empty
      requireDescription={task.todo_description.trim() !== ""}
      onSubmit={(values) => save.mutate({ before: initialValues, values })}
    />
  );
}
