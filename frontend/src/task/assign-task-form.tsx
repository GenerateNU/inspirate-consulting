import { useNavigate, useSearchParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";

import {
  getGetTodoItemsQueryKey,
  getListCounselorStudentsQueryKey,
  useBulkCreateTodoItems,
  useCreateTodoItem,
} from "../api/endpoints";
import type {
  BulkCreateTodoItemsBody,
  CreateTodoItemRequestBody,
} from "../api/models";
import { isIndividual } from "./assign-to-section";
import TaskForm, { type TaskFormValues } from "./task-form";
import { DEFAULT_DUE_TIME } from "./todoShared";

export const dueDateToDeadline = (dueDate: string) =>
  dueDate
    ? new Date(`${dueDate}T${DEFAULT_DUE_TIME}`).toISOString()
    : undefined;

// Create a task for one student, or a bulk task for classes and/or several students
export default function AssignTaskForm() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  // If we came from a student's page, preselect that student
  const [searchParams] = useSearchParams();
  const preselectedStudentId = searchParams.get("studentId");

  const createOne = useCreateTodoItem();
  const createBulk = useBulkCreateTodoItems();

  const initialValues: TaskFormValues = {
    title: "",
    description: "",
    dueDate: "",
    assignTo: {
      all: false,
      graduationYears: [],
      studentIds: preselectedStudentId ? [preselectedStudentId] : [],
    },
    connectTo: { type: null, itemId: "" },
  };

  const onSuccess = () => {
    void queryClient.invalidateQueries({ queryKey: getGetTodoItemsQueryKey() });
    void queryClient.invalidateQueries({
      queryKey: getListCounselorStudentsQueryKey(),
    });
    void navigate(
      preselectedStudentId
        ? `/counselor/students/${preselectedStudentId}`
        : "/counselor",
    );
  };

  const handleSubmit = (values: TaskFormValues) => {
    const { title, description, dueDate, assignTo, connectTo } = values;
    const { type, itemId } = connectTo;
    const deadline = dueDateToDeadline(dueDate);

    if (isIndividual(assignTo)) {
      const body: CreateTodoItemRequestBody = {
        student_id: assignTo.studentIds[0],
        title: title.trim(),
        todo_description: description.trim(),
        deadline,
      };
      if (type) body.link_type = type;
      if (itemId && type === "essay") body.essay_id = itemId;
      if (itemId && type === "college") body.global_college_id = Number(itemId);
      if (itemId && type === "module") body.media_id = itemId;
      createOne.mutate({ data: body }, { onSuccess });
      return;
    }

    const body: BulkCreateTodoItemsBody = {
      title: title.trim(),
      description: description.trim(),
      deadline,
    };
    if (assignTo.all) body.all = true;
    else {
      if (assignTo.graduationYears.length > 0)
        body.graduation_years = assignTo.graduationYears;
      if (assignTo.studentIds.length > 0)
        body.student_ids = assignTo.studentIds;
    }
    // groups can only link to a college or a module
    if (type === "college" || type === "module") body.link_type = type;
    if (itemId && type === "college") body.college_id = Number(itemId);
    if (itemId && type === "module") body.media_id = itemId;
    createBulk.mutate({ data: body }, { onSuccess });
  };

  return (
    <TaskForm
      heading="Create Task"
      initialValues={initialValues}
      submitLabel="Assign task"
      pendingLabel="Assigning..."
      isPending={createOne.isPending || createBulk.isPending}
      error={createOne.error ?? createBulk.error}
      onSubmit={handleSubmit}
    />
  );
}
