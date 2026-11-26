import { useQuery, useQueryClient } from "@tanstack/react-query";

import {
  getGetTodoItemsQueryKey,
  getTodoItems,
  useDeleteTodoItem,
  useUpdateTodoItemCompleted,
} from "../../../api/endpoints";
import {
  GetTodoItemsSortBy,
  GetTodoItemsSortOrder,
  GetTodoItemsStatus,
  type TodoItem,
} from "../../../api/models";
import { isDueByEndOfWeek, isThisWeek } from "../utils/this-week";

// 100 is the backend's max page size
const PAGE_SIZE = 100;

// One student's tasks for "This Week's Tasks": to do = not done and due by the
// end of this week (overdue included); done = checked off this week
export function useThisWeeksTasks(studentId: string) {
  const queryClient = useQueryClient();
  const params = {
    student_id: studentId,
    status: GetTodoItemsStatus.all,
    sort_by: GetTodoItemsSortBy.deadline,
    sort_order: GetTodoItemsSortOrder.asc,
  };

  // every page, since tasks due this week can be anywhere in the list.
  // Starts with getGetTodoItemsQueryKey(), so invalidating that refreshes it
  const tasksQuery = useQuery({
    queryKey: [...getGetTodoItemsQueryKey(params), "all"],
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
    enabled: studentId !== "",
  });

  const refreshTasks = () =>
    queryClient.invalidateQueries({ queryKey: getGetTodoItemsQueryKey() });
  const toggleCompleted = useUpdateTodoItemCompleted({
    mutation: { onSuccess: refreshTasks },
  });
  const deleteTask = useDeleteTodoItem({
    mutation: { onSuccess: refreshTasks },
  });

  const tasks = tasksQuery.data ?? [];
  const toDo = tasks.filter(
    (task) =>
      task.completed_at == null &&
      task.deadline != null &&
      isDueByEndOfWeek(task.deadline),
  );
  const done = tasks.filter(
    (task) => task.completed_at != null && isThisWeek(task.completed_at),
  );

  return {
    toDo,
    done,
    isLoading: tasksQuery.isLoading,
    error: tasksQuery.error,
    refetch: tasksQuery.refetch,
    toggleCompleted: (task: TodoItem) =>
      toggleCompleted.mutate({
        id: task.id,
        data: { completed: task.completed_at == null },
      }),
    toggleError: toggleCompleted.error,
    deleteTask,
  };
}
