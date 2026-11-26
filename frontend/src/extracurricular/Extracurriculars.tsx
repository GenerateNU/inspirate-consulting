import { useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Plus } from "lucide-react";
import {
  getListExtracurricularsQueryKey,
  useDeleteExtracurricular,
  useListExtracurriculars,
} from "../api/endpoints";
import type { Extracurricular } from "../api/models";
import { Button } from "@/components/ui/button";
import ExtracurricularCard from "./ExtracurricularCard";
import { ErrorMessage } from "@/components/ui/error-message";

export default function Extracurriculars() {
  const queryClient = useQueryClient();
  const { data, isLoading, error, refetch } = useListExtracurriculars();
  const {
    mutate: deleteItem,
    isPending: isDeleting,
    error: deleteError,
  } = useDeleteExtracurricular({
    mutation: {
      onSuccess: () =>
        queryClient.invalidateQueries({
          queryKey: getListExtracurricularsQueryKey(),
        }),
    },
  });
  const list = data ?? [];

  const handleDelete = (item: Extracurricular) => {
    if (!window.confirm(`Delete ${item.name}?`)) return;
    deleteItem({ id: item.id });
  };

  return (
    <div className="flex min-h-full flex-1 flex-col gap-8 bg-[#fbfaf7] px-6 py-10 font-host-grotesk md:px-[108px] md:py-16">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="font-jakarta text-[30px] font-bold text-charcoal">
          Extracurriculars
        </h1>
        <Button
          asChild
          className="h-auto w-[200px] rounded-lg bg-[#19664c] px-6 py-2 text-sm leading-5 font-normal text-[#fbfaf7] hover:bg-[#19664c]/90 [&_svg:not([class*='size-'])]:size-5"
        >
          <Link to="/extracurriculars/new">
            <Plus aria-hidden="true" />
            Add Extracurricular
          </Link>
        </Button>
      </div>

      {isLoading && <p>Loading...</p>}
      <ErrorMessage
        variant="banner"
        title="Couldn't load your extracurriculars"
        error={error}
        onRetry={refetch}
      />
      <ErrorMessage error={deleteError} prefix="Couldn't delete" />
      {!isLoading && !error && list.length === 0 && (
        <p>
          No extracurriculars yet.{" "}
          <Link to="/extracurriculars/new" className="underline">
            Add your first one
          </Link>
          .
        </p>
      )}

      {list.length > 0 && (
        <ul className="flex flex-col gap-4">
          {list.map((item) => (
            <li key={item.id}>
              <ExtracurricularCard
                item={item}
                to={`/extracurriculars/${item.id}`}
                onDelete={() => handleDelete(item)}
                deleteDisabled={isDeleting}
              />
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
