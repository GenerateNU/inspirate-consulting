import { useEffect } from "react";

import { queryClient } from "../../../api/client";
import {
  getGetMediaAccessQueryKey,
  useGetMediaAccess,
  useRecordMediaProgress,
} from "../../../api/endpoints";
import { RecordMediaProgressInputBodyStatus } from "../../../api/models";

// one assigned module plus the student's progress on it.
// Opening the page marks it started; markComplete marks it done
export function useStudentModule(id: string) {
  // the list endpoint is the one that carries progress and due dates, and the
  // Modules page has usually cached it already. 100 is the backend's max page size
  const modulesQuery = useGetMediaAccess({ limit: 100 });
  const module = modulesQuery.data?.find((item) => item.id === id) ?? null;

  const progress = useRecordMediaProgress({
    mutation: {
      // every cached list (searched or not) now has stale progress
      onSuccess: () =>
        queryClient.invalidateQueries({
          queryKey: getGetMediaAccessQueryKey(),
        }),
    },
  });
  const { mutate: recordProgress } = progress;

  const notStartedId = module && !module.started_at ? module.id : null;
  useEffect(() => {
    // the backend keeps the first timestamp, so a repeat call is harmless
    if (notStartedId) {
      recordProgress({
        id: notStartedId,
        data: { status: RecordMediaProgressInputBodyStatus.started },
      });
    }
  }, [notStartedId, recordProgress]);

  const markComplete = () => {
    if (!module) return;
    progress.mutate({
      id: module.id,
      data: { status: RecordMediaProgressInputBodyStatus.completed },
    });
  };

  return {
    module,
    isLoading: modulesQuery.isLoading,
    error: modulesQuery.error,
    markComplete,
    marking: progress.isPending,
    markError: progress.error,
  };
}
