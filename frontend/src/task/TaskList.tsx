import { useEffect, useState } from "react";

import { useUpdateTodoItemCompleted } from "../api/endpoints";
import {
  GetTodoItemsSortBy,
  GetTodoItemsSortOrder,
  GetTodoItemsStatus,
  type TodoItem,
} from "../api/models";
import {
  LINK_TYPES,
  LINK_TYPE_LABELS,
  errorMessage,
  formatDate,
  formatDeadline,
  useTodoLinkOptions,
  type LinkType,
} from "./todoShared";
import { useTaskFilterParams, useTodoPages } from "./useTodoList";

const SORT_LABELS: Record<GetTodoItemsSortBy, string> = {
  created_at: "Date created",
  deadline: "Due date",
  completed_at: "Date completed",
  updated_at: "Last updated",
};

function TaskItem({
  task,
  linkLabel,
  onToggled,
}: {
  task: TodoItem;
  linkLabel: string | null;
  onToggled: () => void;
}) {
  const { trigger } = useUpdateTodoItemCompleted(task.id);

  const handleToggle = async () => {
    try {
      const res = await trigger({ completed: task.completed_at == null });
      if (res.status !== 200) console.error(errorMessage(res.data));
      onToggled();
    } catch (err) {
      console.error(err);
    }
  };

  const overdue =
    task.completed_at == null &&
    task.deadline != null &&
    new Date(task.deadline) < new Date();

  return (
    <li>
      <input
        type="checkbox"
        checked={task.completed_at !== null}
        onChange={handleToggle}
      />
      {task.todo_description}
      {task.deadline && (
        <span>
          {" "}
          — due {formatDeadline(task.deadline)}
          {overdue && <strong> (overdue)</strong>}
        </span>
      )}
      {task.completed_at && (
        <span> — completed {formatDate(task.completed_at)}</span>
      )}
      {linkLabel && <span> [{linkLabel}]</span>}
    </li>
  );
}

export default function TaskList() {
  const { filters, q, linkType, linkId, update, setLink } =
    useTaskFilterParams();
  const { options, labelFor } = useTodoLinkOptions();
  const { items, error, isLoading, hasMore, isLoadingMore, loadMore, refresh } =
    useTodoPages(filters, q);
  const [search, setSearch] = useState(q);

  // Debounce typing before it becomes a request
  useEffect(() => {
    const timer = setTimeout(() => {
      if (search !== q) update({ q: search });
    }, 300);
    return () => clearTimeout(timer);
  }, [search, q, update]);

  return (
    <div>
      <h2>To-Do List</h2>
      <div>
        <label>
          Search:{" "}
          <input
            type="search"
            value={search}
            placeholder="Search descriptions"
            maxLength={200}
            onChange={(e) => setSearch(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Status:{" "}
          <select
            value={filters.status}
            onChange={(e) => update({ status: e.target.value })}
          >
            <option value={GetTodoItemsStatus.all}>All</option>
            <option value={GetTodoItemsStatus.incomplete}>Incomplete</option>
            <option value={GetTodoItemsStatus.completed}>Completed</option>
          </select>
        </label>{" "}
        <label>
          Linked to:{" "}
          <select
            value={linkType}
            onChange={(e) => setLink(e.target.value as LinkType | "", "")}
          >
            <option value="">Anything</option>
            {LINK_TYPES.map((type) => (
              <option key={type} value={type}>
                {LINK_TYPE_LABELS[type]}
              </option>
            ))}
          </select>
        </label>{" "}
        {linkType && (
          <select
            value={linkId}
            onChange={(e) => setLink(linkType, e.target.value)}
          >
            <option value="">
              Any {LINK_TYPE_LABELS[linkType].toLowerCase()}
            </option>
            {options[linkType].map((o) => (
              <option key={o.id} value={o.id}>
                {o.label}
              </option>
            ))}
          </select>
        )}
      </div>
      <div>
        <label>
          Sort by:{" "}
          <select
            value={filters.sort_by}
            onChange={(e) => update({ sort_by: e.target.value })}
          >
            {(Object.keys(SORT_LABELS) as GetTodoItemsSortBy[]).map((key) => (
              <option key={key} value={key}>
                {SORT_LABELS[key]}
              </option>
            ))}
          </select>
        </label>{" "}
        <select
          value={filters.sort_order}
          onChange={(e) => update({ sort_order: e.target.value })}
        >
          <option value={GetTodoItemsSortOrder.asc}>Ascending</option>
          <option value={GetTodoItemsSortOrder.desc}>Descending</option>
        </select>
      </div>

      {isLoading && <p>Loading...</p>}
      {error && <p>Error: {error.message}</p>}
      {!isLoading && !error && items.length === 0 && (
        <p>{q.trim() ? "No tasks match your search." : "No tasks yet."}</p>
      )}
      <ul>
        {items.map((task) => (
          <TaskItem
            key={task.id}
            task={task}
            linkLabel={labelFor(task)}
            onToggled={refresh}
          />
        ))}
      </ul>
      {hasMore && (
        <button onClick={loadMore} disabled={isLoadingMore}>
          {isLoadingMore ? "Loading..." : "Load more"}
        </button>
      )}
    </div>
  );
}
