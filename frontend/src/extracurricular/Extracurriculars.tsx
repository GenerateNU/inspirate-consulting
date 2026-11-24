import type { ReactNode } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Palette, Plus, Trash2 } from "lucide-react";
import { errorMessage } from "../api/client";
import {
  getListExtracurricularsQueryKey,
  useDeleteExtracurricular,
  useListExtracurriculars,
} from "../api/endpoints";
import type { Extracurricular } from "../api/models";
import { Button } from "@/components/ui/button";
import { categoryLabel, dateRangeLabel } from "./ExtracurricularShared";

export default function Extracurriculars() {
  const queryClient = useQueryClient();
  const { data, isLoading, error } = useListExtracurriculars();
  const { mutate: deleteItem, isPending: isDeleting } =
    useDeleteExtracurricular({
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
      {error && <p>Error: {errorMessage(error)}</p>}
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

function ExtracurricularCard({
  item,
  onDelete,
  deleteDisabled,
}: {
  item: Extracurricular;
  onDelete: () => void;
  deleteDisabled: boolean;
}) {
  const dates = dateRangeLabel(item);

  return (
    <div className="relative flex min-h-[136px] items-center gap-4 overflow-hidden rounded-lg border border-cream bg-white py-6 pr-6 pl-8 transition-shadow focus-within:ring-2 focus-within:ring-ring/50 hover:shadow-sm">
      <span
        aria-hidden="true"
        className="absolute inset-y-0 left-0 w-1.5 bg-tan"
      />
      <span
        aria-hidden="true"
        className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-[#f8efe4] text-charcoal"
      >
        <Palette className="size-6" strokeWidth={1.75} />
      </span>

      <div className="flex min-w-0 flex-1 flex-col gap-2">
        {/* The link's ::after covers the card so the whole card is clickable. */}
        <Link
          to={`/extracurriculars/${item.id}`}
          className="truncate text-xl font-medium text-charcoal outline-none after:absolute after:inset-0"
        >
          {item.name}
        </Link>
        <div className="flex flex-wrap items-center gap-2">
          {item.leadership_role && (
            <Pill className="bg-[#e7edef] text-[#415763]">
              {item.leadership_role}
            </Pill>
          )}
          {dates && (
            <Pill className="bg-[#e6f5ef] text-[#19664c]">{dates}</Pill>
          )}
          <Pill className="bg-[#f8efe4] text-[#786246]">
            {categoryLabel(item.type)}
          </Pill>
        </div>
      </div>

      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="relative z-10 rounded-lg text-charcoal [&_svg:not([class*='size-'])]:size-[18px]"
        onClick={onDelete}
        disabled={deleteDisabled}
        aria-label={`Delete ${item.name}`}
      >
        <Trash2 aria-hidden="true" />
      </Button>
    </div>
  );
}

function Pill({
  className,
  children,
}: {
  className: string;
  children: ReactNode;
}) {
  return (
    <span
      className={`rounded-full px-2 py-1 text-xs leading-none whitespace-nowrap ${className}`}
    >
      {children}
    </span>
  );
}
