import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

import { errorMessage } from "../../api/client";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
} from "../../api/models";
import chatCircleIcon from "@/assets/college/chat-circle.svg";
import closeIcon from "@/assets/college/close-x.svg";
import plusIcon from "@/assets/college/plus.svg";
import searchIcon from "@/assets/college/search.svg";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog";
import ChoiceChips from "../utils/components/choice-chips";
import Tag from "../utils/components/tag";
import { cycleLabel } from "../utils/shared";
import { useCollegeSearch } from "../utils/use-college-search";
import { categoryOptions } from "./application-labels";

// what the student has picked on one search result
interface Selection {
  cycle: ApplicationType | null;
  category: Category | null;
}

const emptySelection: Selection = { cycle: null, category: null };

interface AddSchoolModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  addedCollegeIds: Set<number>;
  adding: boolean;
  error: unknown;
  onAdd: (
    globalCollegeId: number,
    applicationType: ApplicationType,
    category: Category,
    onSuccess: () => void,
  ) => void;
}

// search for a school, pick the cycle (then classification) and add it to My List
// stays open so several schools can be added in a row
export default function AddSchoolModal({
  open,
  onOpenChange,
  addedCollegeIds,
  adding,
  error,
  onAdd,
}: AddSchoolModalProps) {
  const search = useCollegeSearch();
  const navigate = useNavigate();
  const [selections, setSelections] = useState<Record<number, Selection>>({});
  const [addingId, setAddingId] = useState<number | null>(null);
  const searchRef = useRef<HTMLInputElement>(null);

  const selectionFor = (id: number) => selections[id] ?? emptySelection;
  const select = (id: number, change: Partial<Selection>) =>
    setSelections((prev) => ({
      ...prev,
      [id]: { ...selectionFor(id), ...change },
    }));

  const add = (college: GlobalCollege) => {
    const { cycle, category } = selectionFor(college.id);
    if (!cycle || !category) return;
    setAddingId(college.id);
    onAdd(college.id, cycle, category, () =>
      setSelections((prev) => {
        const next = { ...prev };
        delete next[college.id];
        return next;
      }),
    );
  };

  const renderResult = (college: GlobalCollege) => {
    const { cycle, category } = selectionFor(college.id);
    const offered = college.cycles ?? [];
    const added = addedCollegeIds.has(college.id);
    const ready = cycle != null && category != null;

    return (
      <li key={college.id} className="flex flex-col gap-2.5">
        <div className="flex items-center justify-between gap-4 px-5 py-2.5">
          <div className="flex min-w-0 flex-1 flex-col gap-3">
            <div className="flex flex-col gap-1">
              <p className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
                {college.school_name}
              </p>
              <p className="font-host text-xs leading-[normal] text-charcoal/50">
                {college.school_location}
              </p>
            </div>

            {added ? null : offered.length === 0 ? (
              <p className="font-host text-xs text-warm-grey-dark">
                No application cycles yet. Ask your counselor to add them.
              </p>
            ) : (
              <>
                <ChoiceChips
                  label={`Application cycle for ${college.school_name}`}
                  options={offered.map(({ cycle: value }) => ({
                    value: value as ApplicationType,
                    label: cycleLabel(value),
                  }))}
                  value={cycle}
                  onChange={(value) => select(college.id, { cycle: value })}
                />
                {cycle && (
                  <ChoiceChips
                    label={`Classification for ${college.school_name}`}
                    options={categoryOptions}
                    value={category}
                    onChange={(value) =>
                      select(college.id, { category: value })
                    }
                    chipClassName="min-w-[60px]"
                  />
                )}
              </>
            )}
          </div>

          {added ? (
            <Tag tone="green">On your list</Tag>
          ) : (
            <button
              type="button"
              data-slot="icon-button"
              aria-label={`Add ${college.school_name} to your list`}
              title={
                ready ? "Add to your list" : "Pick a cycle and classification"
              }
              disabled={!ready || adding}
              onClick={() => add(college)}
              className="flex size-8 shrink-0 items-center justify-center rounded-full border border-greyish-brown bg-off-white hover:bg-lightest-green disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-off-white"
            >
              {adding && addingId === college.id ? (
                <span className="font-host text-[10px]">...</span>
              ) : (
                <img src={plusIcon} alt="" width={16} height={16} />
              )}
            </button>
          )}
        </div>
        <div aria-hidden="true" className="h-[2px] w-full bg-off-white" />
      </li>
    );
  };

  let results;
  if (search.isLoading || (search.searching && search.searchLoading)) {
    results = <p className="py-6 text-center font-host text-sm">Loading...</p>;
  } else if (search.error || (search.searching && search.searchError)) {
    results = (
      <p className="py-6 text-center font-host text-sm">
        Error: {errorMessage(search.error ?? search.searchError)}
      </p>
    );
  } else if (search.colleges.length === 0) {
    results = (
      <p className="py-6 text-center font-host text-sm text-warm-grey-dark">
        {search.searching ? "No matching schools." : "No schools yet."}
      </p>
    );
  } else {
    results = (
      <ul className="flex flex-col gap-2.5">
        {search.colleges.map(renderResult)}
      </ul>
    );
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        showCloseButton={false}
        // start in the search box instead of on the close button
        onOpenAutoFocus={(e) => {
          e.preventDefault();
          searchRef.current?.focus();
        }}
        className="flex w-[712px] flex-col items-center gap-2 rounded-lg border-0 bg-off-white px-6 pt-5 pb-[25px] drop-shadow-[0px_1px_1.5px_rgba(22,25,29,0.1)] sm:max-w-[712px]"
      >
        <div className="flex w-full justify-end">
          <DialogClose asChild>
            <button
              type="button"
              data-slot="icon-button"
              aria-label="Close"
              className="rounded-sm hover:opacity-70"
            >
              <img src={closeIcon} alt="" width={20} height={20} />
            </button>
          </DialogClose>
        </div>
        <DialogTitle className="font-jakarta text-xl leading-[normal] font-medium text-black">
          Add School
        </DialogTitle>
        <DialogDescription className="sr-only">
          Search for a school, pick how you&apos;re applying, and add it to your
          list.
        </DialogDescription>

        <div className="flex h-[500px] max-h-[60vh] w-full flex-col gap-2.5 px-[30px] py-2.5">
          <label className="flex h-10 w-full shrink-0 items-center gap-2 rounded-full border border-greyish-brown bg-white px-3">
            <span className="flex size-4 shrink-0 items-center justify-center">
              <img src={searchIcon} alt="" width={12.6667} height={12.6667} />
            </span>
            <input
              type="text"
              ref={searchRef}
              data-slot="search-input"
              value={search.query}
              onChange={(e) => search.setQuery(e.target.value)}
              placeholder="Search by name or location"
              aria-label="Search schools"
              className="min-w-0 flex-1 bg-transparent font-host text-sm text-charcoal outline-none placeholder:text-warm-grey"
            />
          </label>

          {error != null && (
            <p role="alert" className="font-host text-sm text-danger">
              Error: {errorMessage(error)}
            </p>
          )}

          <div className="min-h-0 flex-1 overflow-y-auto rounded-lg bg-white py-5 shadow-card">
            {results}
          </div>
        </div>

        <div className="flex flex-col items-center gap-[3px]">
          <p className="font-host text-xs leading-[normal] text-warm-grey-dark">
            Can&apos;t find your school?
          </p>
          {/* TODO: open the chat with this student's counselor directly once
              /user/me returns the counselor's user id */}
          <button
            type="button"
            data-slot="chat-counselor"
            onClick={() => navigate("/chats")}
            className="flex items-center justify-center gap-[3px] rounded-[10px] border border-white bg-black px-4 py-2 font-host text-[13px] leading-[normal] text-white hover:bg-black/85"
          >
            <img src={chatCircleIcon} alt="" width={15} height={15} />
            Chat with your counselor
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
