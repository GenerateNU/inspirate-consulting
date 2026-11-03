import { useInfiniteQuery } from "@tanstack/react-query";
import { type FetchStudentsParams } from "../api/models";
import { fetchStudents, getFetchStudentsQueryKey } from "../api/endpoints";

// the endpoint allows at most 10 students per page
const PAGE_SIZE = 10;

type StudentFilters = Omit<FetchStudentsParams, "cursor" | "limit">;

// Pages through fetch-students with its next_cursor, keeping every loaded page.
// Call loadMore() (e.g. when the end of the list scrolls into view) to fetch the next one.
// Named differently from Orval's useFetchStudents (a single page) so the two can't be mixed up.
export function useInfiniteStudents(params: StudentFilters) {
    const { data, error, isLoading, hasNextPage, isFetchingNextPage, fetchNextPage } = useInfiniteQuery({
        // own key so this cache entry ({ pages }) can't collide with the plain useFetchStudents query;
        // it still starts with getFetchStudentsQueryKey(), so invalidating that refreshes it
        queryKey: [...getFetchStudentsQueryKey(params), "infinite"],
        queryFn: ({ pageParam, signal }) =>
            fetchStudents({ ...params, cursor: pageParam, limit: PAGE_SIZE }, { signal }),
        initialPageParam: undefined as string | undefined,
        // no next_cursor means this was the last page
        getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    });

    return {
        students: data?.pages.flatMap((page) => page.students ?? []) ?? [],
        error,
        isLoading,
        hasMore: hasNextPage,
        isLoadingMore: isFetchingNextPage,
        loadMore: () => fetchNextPage(),
    };
}
