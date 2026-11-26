import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  getGetTodoItemsQueryKey,
  getTodoItems,
  updateTodoItemCompleted,
} from "../api/endpoints";
import {
  GetTodoItemsSortBy,
  GetTodoItemsSortOrder,
  type GetTodoItemsParams,
  type TodoItem,
} from "../api/models";

const PAGE_SIZE = 100;

const params: GetTodoItemsParams = {
  sort_by: GetTodoItemsSortBy.deadline,
  sort_order: GetTodoItemsSortOrder.asc,
};

// starts with getGetTodoItemsQueryKey(), so invalidating that refreshes it
const queryKey = [...getGetTodoItemsQueryKey(params), "all"];

// Every one of the signed-in student's tasks, soonest due first.
// The dashboard groups them by day, so it needs all pages, not just the first.
export function useStudentTasks() {
  return useQuery({
    queryKey,
    queryFn: async ({ signal }) => {
      const tasks: TodoItem[] = [];
      let cursor: string | undefined;
      do {
        const page = await getTodoItems(
          { ...params, cursor, limit: PAGE_SIZE },
          { signal },
        );
        tasks.push(...(page.items ?? []));
        cursor = page.next_cursor ?? undefined;
      } while (cursor);
      return tasks;
    },
  });
}

// Checks a task off (or back on) right away, then saves; undoes the check if saving fails
export function useToggleTaskCompleted() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (task: TodoItem) =>
      updateTodoItemCompleted(task.id, {
        completed: task.completed_at == null,
      }),
    onMutate: async (task) => {
      await queryClient.cancelQueries({ queryKey });
      const previous = queryClient.getQueryData<TodoItem[]>(queryKey);
      queryClient.setQueryData<TodoItem[]>(queryKey, (tasks) =>
        tasks?.map((t) =>
          t.id === task.id
            ? {
                ...t,
                completed_at:
                  t.completed_at == null ? new Date().toISOString() : null,
              }
            : t,
        ),
      );
      return { previous };
    },
    onError: (_error, _task, context) => {
      if (context?.previous)
        queryClient.setQueryData(queryKey, context.previous);
    },
    onSettled: () =>
      queryClient.invalidateQueries({ queryKey: getGetTodoItemsQueryKey() }),
  });
}
