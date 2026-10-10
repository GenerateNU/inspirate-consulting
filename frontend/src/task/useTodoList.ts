import { useSearchParams } from "react-router-dom";
import useSWRInfinite from "swr/infinite";

import { getTodoItems, searchTodoItems } from "../api/endpoints";
import {
  GetTodoItemsSortBy,
  GetTodoItemsSortOrder,
  GetTodoItemsStatus,
  type GetTodoItemsParams,
  type TodoItemPage,
} from "../api/models";
import { LINK_TYPES, errorMessage, type LinkType } from "./todoShared";

const PAGE_SIZE = 5;

const LINK_ID_PARAMS = {
  essay: "essay_id",
  media: "media_id",
  college: "global_college_id",
} as const;

type TodoFilters = Omit<GetTodoItemsParams, "cursor" | "limit">;

const oneOf = <T extends string>(
  values: Record<string, T>,
  value: string | null,
) => Object.values(values).find((v) => v === value);

export function useTaskFilterParams() {
  const [params, setParams] = useSearchParams();

  const linkedType = LINK_TYPES.find((t) => params.get(LINK_ID_PARAMS[t]));
  const linkType: LinkType | "" =
    linkedType ?? LINK_TYPES.find((t) => t === params.get("linked_to")) ?? "";
  const linkId = linkedType
    ? (params.get(LINK_ID_PARAMS[linkedType]) ?? "")
    : "";

  const filters: TodoFilters = {
    status:
      oneOf(GetTodoItemsStatus, params.get("status")) ?? GetTodoItemsStatus.all,
    sort_by:
      oneOf(GetTodoItemsSortBy, params.get("sort_by")) ??
      GetTodoItemsSortBy.created_at,
    sort_order:
      oneOf(GetTodoItemsSortOrder, params.get("sort_order")) ??
      GetTodoItemsSortOrder.desc,
  };
  if (linkType === "essay" && linkId) filters.essay_id = linkId;
  else if (linkType === "media" && linkId) filters.media_id = linkId;
  else if (linkType === "college" && linkId)
    filters.global_college_id = Number(linkId);
  else if (linkType) filters.linked_to = linkType;

  const update = (changes: Record<string, string>) =>
    setParams(
      (prev) => {
        const next = new URLSearchParams(prev);
        for (const [key, value] of Object.entries(changes)) {
          if (value) next.set(key, value);
          else next.delete(key);
        }
        return next;
      },
      { replace: true },
    );

  const setLink = (type: LinkType | "", id: string) =>
    update({
      linked_to: type && !id ? type : "",
      essay_id: type === "essay" ? id : "",
      media_id: type === "media" ? id : "",
      global_college_id: type === "college" ? id : "",
    });

  return {
    filters,
    q: params.get("q") ?? "",
    linkType,
    linkId,
    update,
    setLink,
  };
}

export function useTodoPages(filters: TodoFilters, q: string) {
  const search = q.trim();

  const { data, error, isLoading, isValidating, size, setSize, mutate } =
    useSWRInfinite<TodoItemPage, Error>(
      (pageIndex, previous: TodoItemPage | null) => {
        if (previous && !previous.next_cursor) return null;
        const cursor =
          pageIndex === 0 ? undefined : (previous?.next_cursor ?? undefined);
        return ["todo-items", search, filters, cursor] as const;
      },
      async ([, search, filters, cursor]: readonly [
        string,
        string,
        TodoFilters,
        string | undefined,
      ]) => {
        const res = search
          ? await searchTodoItems({
              ...filters,
              q: search,
              cursor,
              limit: PAGE_SIZE,
            })
          : await getTodoItems({ ...filters, cursor, limit: PAGE_SIZE });
        if (res.status !== 200) throw new Error(errorMessage(res.data));
        return res.data;
      },
    );

  return {
    items: data?.flatMap((page) => page.items ?? []) ?? [],
    error,
    isLoading,
    hasMore: !!data?.[data.length - 1]?.next_cursor,
    isLoadingMore: isValidating && !!data && size > data.length,
    loadMore: () => setSize(size + 1),
    refresh: () => mutate(),
  };
}
